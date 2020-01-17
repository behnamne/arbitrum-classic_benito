/*
 * Copyright 2020, Offchain Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ethbridge

import (
	"bytes"
	"context"
	"math/big"
	"strings"

	errors2 "github.com/pkg/errors"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/offchainlabs/arbitrum/packages/arb-util/common"
	"github.com/offchainlabs/arbitrum/packages/arb-util/hashing"
	"github.com/offchainlabs/arbitrum/packages/arb-util/protocol"
	"github.com/offchainlabs/arbitrum/packages/arb-util/value"
	"github.com/offchainlabs/arbitrum/packages/arb-validator/arbbridge"
	"github.com/offchainlabs/arbitrum/packages/arb-validator/ethbridge/rollup"
	"github.com/offchainlabs/arbitrum/packages/arb-validator/structures"
	"github.com/offchainlabs/arbitrum/packages/arb-validator/valprotocol"
)

var rollupStakeCreatedID ethcommon.Hash
var rollupChallengeStartedID ethcommon.Hash
var rollupChallengeCompletedID ethcommon.Hash
var rollupRefundedID ethcommon.Hash
var rollupPrunedID ethcommon.Hash
var rollupStakeMovedID ethcommon.Hash
var rollupAssertedID ethcommon.Hash
var rollupConfirmedID ethcommon.Hash
var confirmedAssertionID ethcommon.Hash
var messageDeliveredID ethcommon.Hash

func init() {
	parsedRollup, err := abi.JSON(strings.NewReader(rollup.ArbRollupABI))
	if err != nil {
		panic(err)
	}
	parsedInbox, err := abi.JSON(strings.NewReader(rollup.IGlobalPendingInboxABI))
	if err != nil {
		panic(err)
	}
	rollupStakeCreatedID = parsedRollup.Events["RollupStakeCreated"].ID()
	rollupChallengeStartedID = parsedRollup.Events["RollupChallengeStarted"].ID()
	rollupChallengeCompletedID = parsedRollup.Events["RollupChallengeCompleted"].ID()
	rollupRefundedID = parsedRollup.Events["RollupStakeRefunded"].ID()
	rollupPrunedID = parsedRollup.Events["RollupPruned"].ID()
	rollupStakeMovedID = parsedRollup.Events["RollupStakeMoved"].ID()
	rollupAssertedID = parsedRollup.Events["RollupAsserted"].ID()
	rollupConfirmedID = parsedRollup.Events["RollupConfirmed"].ID()
	confirmedAssertionID = parsedRollup.Events["ConfirmedAssertion"].ID()

	messageDeliveredID = parsedInbox.Events["MessageDelivered"].ID()
}

type ethRollupWatcher struct {
	ArbRollup          *rollup.ArbRollup
	GlobalPendingInbox *rollup.IGlobalPendingInbox

	rollupAddress ethcommon.Address
	inboxAddress  ethcommon.Address
	client        *ethclient.Client
}

func newRollupWatcher(rollupAddress ethcommon.Address, client *ethclient.Client) (*ethRollupWatcher, error) {
	arbitrumRollupContract, err := rollup.NewArbRollup(rollupAddress, client)
	if err != nil {
		return nil, errors2.Wrap(err, "Failed to connect to arbRollup")
	}

	globalPendingInboxAddress, err := arbitrumRollupContract.GlobalInbox(&bind.CallOpts{
		Pending: false,
		Context: context.Background(),
	})
	if err != nil {
		return nil, errors2.Wrap(err, "Failed to get GlobalPendingInbox address")
	}
	globalPendingContract, err := rollup.NewIGlobalPendingInbox(globalPendingInboxAddress, client)
	if err != nil {
		return nil, errors2.Wrap(err, "Failed to connect to GlobalPendingInbox")
	}

	return &ethRollupWatcher{
		ArbRollup:          arbitrumRollupContract,
		GlobalPendingInbox: globalPendingContract,
		rollupAddress:      rollupAddress,
		inboxAddress:       globalPendingInboxAddress,
		client:             client,
	}, nil
}

func (vm *ethRollupWatcher) GetEvents(ctx context.Context, blockId *structures.BlockId) ([]arbbridge.Event, error) {
	bh := blockId.HeaderHash.ToEthHash()
	addressIndex := ethcommon.Hash{}
	copy(addressIndex[:], ethcommon.LeftPadBytes(vm.rollupAddress.Bytes(), 32))
	inboxLogs, err := vm.client.FilterLogs(ctx, ethereum.FilterQuery{
		BlockHash: &bh,
		Addresses: []ethcommon.Address{vm.inboxAddress},
		Topics: [][]ethcommon.Hash{
			{messageDeliveredID},
			{addressIndex},
		},
	})
	if err != nil {
		return nil, err
	}
	rollupLogs, err := vm.client.FilterLogs(ctx, ethereum.FilterQuery{
		BlockHash: &bh,
		Addresses: []ethcommon.Address{vm.rollupAddress},
		Topics: [][]ethcommon.Hash{
			{
				rollupStakeCreatedID,
				rollupChallengeStartedID,
				rollupChallengeCompletedID,
				rollupRefundedID,
				rollupPrunedID,
				rollupStakeMovedID,
				rollupAssertedID,
				rollupConfirmedID,
				confirmedAssertionID,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	events := make([]arbbridge.Event, 0, len(inboxLogs)+len(inboxLogs))
	for _, evmLog := range inboxLogs {
		event, err := vm.processEvents(getChainInfo2(evmLog, blockId), evmLog)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	for _, evmLog := range rollupLogs {
		event, err := vm.processEvents(getChainInfo2(evmLog, blockId), evmLog)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (vm *ethRollupWatcher) processEvents(chainInfo arbbridge.ChainInfo, ethLog types.Log) (arbbridge.Event, error) {
	if ethLog.Topics[0] == rollupStakeCreatedID {
		eventVal, err := vm.ArbRollup.ParseRollupStakeCreated(ethLog)
		if err != nil {
			return nil, err
		}
		return arbbridge.StakeCreatedEvent{
			ChainInfo: chainInfo,
			Staker:    common.NewAddressFromEth(eventVal.Staker),
			NodeHash:  eventVal.NodeHash,
		}, nil
	} else if ethLog.Topics[0] == rollupChallengeStartedID {
		eventVal, err := vm.ArbRollup.ParseRollupChallengeStarted(ethLog)
		if err != nil {
			return nil, err
		}
		return arbbridge.ChallengeStartedEvent{
			ChainInfo:         chainInfo,
			Asserter:          common.NewAddressFromEth(eventVal.Asserter),
			Challenger:        common.NewAddressFromEth(eventVal.Challenger),
			ChallengeType:     structures.ChildType(eventVal.ChallengeType.Uint64()),
			ChallengeContract: common.NewAddressFromEth(eventVal.ChallengeContract),
		}, nil
	} else if ethLog.Topics[0] == rollupChallengeCompletedID {
		eventVal, err := vm.ArbRollup.ParseRollupChallengeCompleted(ethLog)
		if err != nil {
			return nil, err
		}
		return arbbridge.ChallengeCompletedEvent{
			ChainInfo:         chainInfo,
			Winner:            common.NewAddressFromEth(eventVal.Winner),
			Loser:             common.NewAddressFromEth(eventVal.Loser),
			ChallengeContract: common.NewAddressFromEth(eventVal.ChallengeContract),
		}, nil
	} else if ethLog.Topics[0] == rollupRefundedID {
		eventVal, err := vm.ArbRollup.ParseRollupStakeRefunded(ethLog)
		if err != nil {
			return nil, err
		}
		return arbbridge.StakeRefundedEvent{
			ChainInfo: chainInfo,
			Staker:    common.NewAddressFromEth(eventVal.Staker),
		}, nil
	} else if ethLog.Topics[0] == rollupPrunedID {
		eventVal, err := vm.ArbRollup.ParseRollupPruned(ethLog)
		if err != nil {
			return nil, err
		}
		return arbbridge.PrunedEvent{
			ChainInfo: chainInfo,
			Leaf:      eventVal.Leaf,
		}, nil
	} else if ethLog.Topics[0] == rollupStakeMovedID {
		eventVal, err := vm.ArbRollup.ParseRollupStakeMoved(ethLog)
		if err != nil {
			return nil, err
		}
		return arbbridge.StakeMovedEvent{
			ChainInfo: chainInfo,
			Staker:    common.NewAddressFromEth(eventVal.Staker),
			Location:  eventVal.ToNodeHash,
		}, nil
	} else if ethLog.Topics[0] == rollupAssertedID {
		eventVal, err := vm.ArbRollup.ParseRollupAsserted(ethLog)
		if err != nil {
			return nil, err
		}
		return arbbridge.AssertedEvent{
			ChainInfo:    chainInfo,
			PrevLeafHash: eventVal.Fields[0],
			Params: &structures.AssertionParams{
				NumSteps: eventVal.NumSteps,
				TimeBounds: &protocol.TimeBoundsBlocks{
					common.NewTimeBlocks(eventVal.TimeBoundsBlocks[0]),
					common.NewTimeBlocks(eventVal.TimeBoundsBlocks[1]),
				},
				ImportedMessageCount: eventVal.ImportedMessageCount,
			},
			Claim: &structures.AssertionClaim{
				AfterPendingTop:       eventVal.Fields[2],
				ImportedMessagesSlice: eventVal.Fields[3],
				AssertionStub: &valprotocol.ExecutionAssertionStub{
					AfterHash:        eventVal.Fields[4],
					DidInboxInsn:     eventVal.DidInboxInsn,
					NumGas:           eventVal.NumArbGas,
					FirstMessageHash: [32]byte{},
					LastMessageHash:  eventVal.Fields[5],
					FirstLogHash:     [32]byte{},
					LastLogHash:      eventVal.Fields[6],
				},
			},
			MaxPendingTop:   eventVal.Fields[1],
			MaxPendingCount: eventVal.PendingCount,
		}, nil
	} else if ethLog.Topics[0] == rollupConfirmedID {
		eventVal, err := vm.ArbRollup.ParseRollupConfirmed(ethLog)
		if err != nil {
			return nil, err
		}
		return arbbridge.ConfirmedEvent{
			ChainInfo: chainInfo,
			NodeHash:  eventVal.NodeHash,
		}, nil
	} else if ethLog.Topics[0] == confirmedAssertionID {
		eventVal, err := vm.ArbRollup.ParseConfirmedAssertion(ethLog)
		if err != nil {
			return nil, err
		}
		return arbbridge.ConfirmedAssertionEvent{
			ChainInfo:   chainInfo,
			LogsAccHash: eventVal.LogsAccHash,
		}, nil
	} else if ethLog.Topics[0] == messageDeliveredID {
		val, err := vm.GlobalPendingInbox.ParseMessageDelivered(ethLog)
		if err != nil {
			return nil, err
		}

		rd := bytes.NewReader(val.Data)
		msgData, err := value.UnmarshalValue(rd)
		if err != nil {
			return nil, err
		}

		messageHash := hashing.SoliditySHA3(
			hashing.Address(common.NewAddressFromEth(val.VmId)),
			hashing.Bytes32(msgData.Hash()),
			hashing.Uint256(val.Value),
			val.TokenType[:],
		)
		msgHashInt := new(big.Int).SetBytes(messageHash.Bytes())

		msgVal, _ := value.NewTupleFromSlice([]value.Value{
			msgData,
			value.NewIntValue(new(big.Int).SetUint64(ethLog.BlockNumber)),
			value.NewIntValue(msgHashInt),
		})

		msg := valprotocol.NewSimpleMessage(msgVal, val.TokenType, val.Value, common.NewAddressFromEth(val.Sender))
		return arbbridge.MessageDeliveredEvent{
			ChainInfo: chainInfo,
			Msg:       msg,
		}, nil
	}
	return nil, errors2.New("unknown arbitrum event type")
}

func (vm *ethRollupWatcher) GetParams(ctx context.Context) (structures.ChainParams, error) {
	rawParams, err := vm.ArbRollup.VmParams(nil)
	if err != nil {
		return structures.ChainParams{}, err
	}
	stakeRequired, err := vm.ArbRollup.GetStakeRequired(nil)
	if err != nil {
		return structures.ChainParams{}, err
	}
	return structures.ChainParams{
		StakeRequirement:        stakeRequired,
		GracePeriod:             common.TimeTicks{rawParams.GracePeriodTicks},
		MaxExecutionSteps:       rawParams.MaxExecutionSteps,
		ArbGasSpeedLimitPerTick: rawParams.ArbGasSpeedLimitPerTick.Uint64(),
	}, nil
}

func (vm *ethRollupWatcher) InboxAddress(ctx context.Context) (common.Address, error) {
	addr, err := vm.ArbRollup.GlobalInbox(nil)
	return common.NewAddressFromEth(addr), err
}
