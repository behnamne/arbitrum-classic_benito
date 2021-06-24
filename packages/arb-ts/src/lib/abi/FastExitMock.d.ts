/* Autogenerated file. Do not edit manually. */
/* tslint:disable */
/* eslint-disable */

import {
  ethers,
  EventFilter,
  Signer,
  BigNumber,
  BigNumberish,
  PopulatedTransaction,
} from 'ethers'
import {
  Contract,
  ContractTransaction,
  Overrides,
  CallOverrides,
} from '@ethersproject/contracts'
import { BytesLike } from '@ethersproject/bytes'
import { Listener, Provider } from '@ethersproject/providers'
import { FunctionFragment, EventFragment, Result } from '@ethersproject/abi'

interface FastExitMockInterface extends ethers.utils.Interface {
  functions: {
    'onTokenTransfer(address,uint256,bytes)': FunctionFragment
    'requestLiquidity(address,address,uint256,uint256,bytes)': FunctionFragment
    'setFee(uint256)': FunctionFragment
  }

  encodeFunctionData(
    functionFragment: 'onTokenTransfer',
    values: [string, BigNumberish, BytesLike]
  ): string
  encodeFunctionData(
    functionFragment: 'requestLiquidity',
    values: [string, string, BigNumberish, BigNumberish, BytesLike]
  ): string
  encodeFunctionData(functionFragment: 'setFee', values: [BigNumberish]): string

  decodeFunctionResult(
    functionFragment: 'onTokenTransfer',
    data: BytesLike
  ): Result
  decodeFunctionResult(
    functionFragment: 'requestLiquidity',
    data: BytesLike
  ): Result
  decodeFunctionResult(functionFragment: 'setFee', data: BytesLike): Result

  events: {
    'Triggered()': EventFragment
  }

  getEvent(nameOrSignatureOrTopic: 'Triggered'): EventFragment
}

export class FastExitMock extends Contract {
  connect(signerOrProvider: Signer | Provider | string): this
  attach(addressOrName: string): this
  deployed(): Promise<this>

  on(event: EventFilter | string, listener: Listener): this
  once(event: EventFilter | string, listener: Listener): this
  addListener(eventName: EventFilter | string, listener: Listener): this
  removeAllListeners(eventName: EventFilter | string): this
  removeListener(eventName: any, listener: Listener): this

  interface: FastExitMockInterface

  functions: {
    onTokenTransfer(
      _sender: string,
      _value: BigNumberish,
      _data: BytesLike,
      overrides?: Overrides
    ): Promise<ContractTransaction>

    'onTokenTransfer(address,uint256,bytes)'(
      _sender: string,
      _value: BigNumberish,
      _data: BytesLike,
      overrides?: Overrides
    ): Promise<ContractTransaction>

    requestLiquidity(
      dest: string,
      erc20: string,
      amount: BigNumberish,
      exitNum: BigNumberish,
      liquidityProof: BytesLike,
      overrides?: Overrides
    ): Promise<ContractTransaction>

    'requestLiquidity(address,address,uint256,uint256,bytes)'(
      dest: string,
      erc20: string,
      amount: BigNumberish,
      exitNum: BigNumberish,
      liquidityProof: BytesLike,
      overrides?: Overrides
    ): Promise<ContractTransaction>

    setFee(
      _fee: BigNumberish,
      overrides?: Overrides
    ): Promise<ContractTransaction>

    'setFee(uint256)'(
      _fee: BigNumberish,
      overrides?: Overrides
    ): Promise<ContractTransaction>
  }

  onTokenTransfer(
    _sender: string,
    _value: BigNumberish,
    _data: BytesLike,
    overrides?: Overrides
  ): Promise<ContractTransaction>

  'onTokenTransfer(address,uint256,bytes)'(
    _sender: string,
    _value: BigNumberish,
    _data: BytesLike,
    overrides?: Overrides
  ): Promise<ContractTransaction>

  requestLiquidity(
    dest: string,
    erc20: string,
    amount: BigNumberish,
    exitNum: BigNumberish,
    liquidityProof: BytesLike,
    overrides?: Overrides
  ): Promise<ContractTransaction>

  'requestLiquidity(address,address,uint256,uint256,bytes)'(
    dest: string,
    erc20: string,
    amount: BigNumberish,
    exitNum: BigNumberish,
    liquidityProof: BytesLike,
    overrides?: Overrides
  ): Promise<ContractTransaction>

  setFee(
    _fee: BigNumberish,
    overrides?: Overrides
  ): Promise<ContractTransaction>

  'setFee(uint256)'(
    _fee: BigNumberish,
    overrides?: Overrides
  ): Promise<ContractTransaction>

  callStatic: {
    onTokenTransfer(
      _sender: string,
      _value: BigNumberish,
      _data: BytesLike,
      overrides?: CallOverrides
    ): Promise<void>

    'onTokenTransfer(address,uint256,bytes)'(
      _sender: string,
      _value: BigNumberish,
      _data: BytesLike,
      overrides?: CallOverrides
    ): Promise<void>

    requestLiquidity(
      dest: string,
      erc20: string,
      amount: BigNumberish,
      exitNum: BigNumberish,
      liquidityProof: BytesLike,
      overrides?: CallOverrides
    ): Promise<string>

    'requestLiquidity(address,address,uint256,uint256,bytes)'(
      dest: string,
      erc20: string,
      amount: BigNumberish,
      exitNum: BigNumberish,
      liquidityProof: BytesLike,
      overrides?: CallOverrides
    ): Promise<string>

    setFee(_fee: BigNumberish, overrides?: CallOverrides): Promise<void>

    'setFee(uint256)'(
      _fee: BigNumberish,
      overrides?: CallOverrides
    ): Promise<void>
  }

  filters: {
    Triggered(): EventFilter
  }

  estimateGas: {
    onTokenTransfer(
      _sender: string,
      _value: BigNumberish,
      _data: BytesLike,
      overrides?: Overrides
    ): Promise<BigNumber>

    'onTokenTransfer(address,uint256,bytes)'(
      _sender: string,
      _value: BigNumberish,
      _data: BytesLike,
      overrides?: Overrides
    ): Promise<BigNumber>

    requestLiquidity(
      dest: string,
      erc20: string,
      amount: BigNumberish,
      exitNum: BigNumberish,
      liquidityProof: BytesLike,
      overrides?: Overrides
    ): Promise<BigNumber>

    'requestLiquidity(address,address,uint256,uint256,bytes)'(
      dest: string,
      erc20: string,
      amount: BigNumberish,
      exitNum: BigNumberish,
      liquidityProof: BytesLike,
      overrides?: Overrides
    ): Promise<BigNumber>

    setFee(_fee: BigNumberish, overrides?: Overrides): Promise<BigNumber>

    'setFee(uint256)'(
      _fee: BigNumberish,
      overrides?: Overrides
    ): Promise<BigNumber>
  }

  populateTransaction: {
    onTokenTransfer(
      _sender: string,
      _value: BigNumberish,
      _data: BytesLike,
      overrides?: Overrides
    ): Promise<PopulatedTransaction>

    'onTokenTransfer(address,uint256,bytes)'(
      _sender: string,
      _value: BigNumberish,
      _data: BytesLike,
      overrides?: Overrides
    ): Promise<PopulatedTransaction>

    requestLiquidity(
      dest: string,
      erc20: string,
      amount: BigNumberish,
      exitNum: BigNumberish,
      liquidityProof: BytesLike,
      overrides?: Overrides
    ): Promise<PopulatedTransaction>

    'requestLiquidity(address,address,uint256,uint256,bytes)'(
      dest: string,
      erc20: string,
      amount: BigNumberish,
      exitNum: BigNumberish,
      liquidityProof: BytesLike,
      overrides?: Overrides
    ): Promise<PopulatedTransaction>

    setFee(
      _fee: BigNumberish,
      overrides?: Overrides
    ): Promise<PopulatedTransaction>

    'setFee(uint256)'(
      _fee: BigNumberish,
      overrides?: Overrides
    ): Promise<PopulatedTransaction>
  }
}
