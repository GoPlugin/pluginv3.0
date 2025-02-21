import { ethers } from 'hardhat'
import { assert } from 'chai'
import { Contract, ContractFactory } from 'ethers'
import { getUsers, Roles } from '../test-helpers/setup'
import {
  convertFufillParams,
  decodeCCRequest,
  decodeRunRequest,
  RunRequest,
} from '../test-helpers/oracle'
import { decodeDietCBOR } from '../test-helpers/helpers'
import { evmRevert } from '../test-helpers/matchers'

let concretePluginClientFactory: ContractFactory
let emptyOracleFactory: ContractFactory
let getterSetterFactory: ContractFactory
let operatorFactory: ContractFactory
let pliTokenFactory: ContractFactory

let roles: Roles

before(async () => {
  roles = (await getUsers()).roles

  concretePluginClientFactory = await ethers.getContractFactory(
    'src/v0.8/tests/PluginClientTestHelper.sol:PluginClientTestHelper',
    roles.defaultAccount,
  )
  emptyOracleFactory = await ethers.getContractFactory(
    'src/v0.8/operatorforwarder/test/testhelpers/EmptyOracle.sol:EmptyOracle',
    roles.defaultAccount,
  )
  getterSetterFactory = await ethers.getContractFactory(
    'src/v0.8/operatorforwarder/test/testhelpers/GetterSetter.sol:GetterSetter',
    roles.defaultAccount,
  )
  operatorFactory = await ethers.getContractFactory(
    'src/v0.8/operatorforwarder/Operator.sol:Operator',
    roles.defaultAccount,
  )
  pliTokenFactory = await ethers.getContractFactory(
    'src/v0.8/shared/test/helpers/PliTokenTestHelper.sol:PliTokenTestHelper',
    roles.defaultAccount,
  )
})

describe('PluginClientTestHelper', () => {
  const specId =
    '0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000000'
  let cc: Contract
  let gs: Contract
  let oc: Contract
  let newoc: Contract
  let pli: Contract

  beforeEach(async () => {
    pli = await pliTokenFactory.connect(roles.defaultAccount).deploy()
    oc = await operatorFactory
      .connect(roles.defaultAccount)
      .deploy(pli.address, await roles.defaultAccount.getAddress())
    newoc = await operatorFactory
      .connect(roles.defaultAccount)
      .deploy(pli.address, await roles.defaultAccount.getAddress())
    gs = await getterSetterFactory.connect(roles.defaultAccount).deploy()
    cc = await concretePluginClientFactory
      .connect(roles.defaultAccount)
      .deploy(pli.address, oc.address)
  })

  describe('#newRequest', () => {
    it('forwards the information to the oracle contract through the pli token', async () => {
      const tx = await cc.publicNewRequest(
        specId,
        gs.address,
        ethers.utils.toUtf8Bytes('requestedBytes32(bytes32,bytes32)'),
      )
      const receipt = await tx.wait()

      assert.equal(1, receipt.logs?.length)
      const [jId, cbAddr, cbFId, cborData] = receipt.logs
        ? decodeCCRequest(receipt.logs[0])
        : []
      const params = decodeDietCBOR(cborData ?? '')

      assert.equal(specId, jId)
      assert.equal(gs.address, cbAddr)
      assert.equal('0xed53e511', cbFId)
      assert.deepEqual({}, params)
    })
  })

  describe('#pluginRequest(Request)', () => {
    it('emits an event from the contract showing the run ID', async () => {
      const tx = await cc.publicRequest(
        specId,
        cc.address,
        ethers.utils.toUtf8Bytes('fulfillRequest(bytes32,bytes32)'),
        0,
      )

      const { events, logs } = await tx.wait()

      assert.equal(4, events?.length)

      assert.equal(logs?.[0].address, cc.address)
      assert.equal(events?.[0].event, 'PluginRequested')
    })
  })

  describe('#pluginRequestTo(Request)', () => {
    it('emits an event from the contract showing the run ID', async () => {
      const tx = await cc.publicRequestRunTo(
        newoc.address,
        specId,
        cc.address,
        ethers.utils.toUtf8Bytes('fulfillRequest(bytes32,bytes32)'),
        0,
      )
      const { events } = await tx.wait()

      assert.equal(4, events?.length)
      assert.equal(events?.[0].event, 'PluginRequested')
    })

    it('emits an event on the target oracle contract', async () => {
      const tx = await cc.publicRequestRunTo(
        newoc.address,
        specId,
        cc.address,
        ethers.utils.toUtf8Bytes('fulfillRequest(bytes32,bytes32)'),
        0,
      )
      const { logs } = await tx.wait()
      const event = logs && newoc.interface.parseLog(logs[3])

      assert.equal(4, logs?.length)
      assert.equal(event?.name, 'OracleRequest')
    })

    it('does not modify the stored oracle address', async () => {
      await cc.publicRequestRunTo(
        newoc.address,
        specId,
        cc.address,
        ethers.utils.toUtf8Bytes('fulfillRequest(bytes32,bytes32)'),
        0,
      )

      const actualOracleAddress = await cc.publicOracleAddress()
      assert.equal(oc.address, actualOracleAddress)
    })
  })

  describe('#requestOracleData', () => {
    it('emits an event from the contract showing the run ID', async () => {
      const tx = await cc.publicRequestOracleData(
        specId,
        ethers.utils.toUtf8Bytes('fulfillRequest(bytes32,bytes32)'),
        0,
      )

      const { events, logs } = await tx.wait()

      assert.equal(4, events?.length)

      assert.equal(logs?.[0].address, cc.address)
      assert.equal(events?.[0].event, 'PluginRequested')
    })
  })

  describe('#requestOracleDataFrom', () => {
    it('emits an event from the contract showing the run ID', async () => {
      const tx = await cc.publicRequestOracleDataFrom(
        newoc.address,
        specId,
        ethers.utils.toUtf8Bytes('fulfillRequest(bytes32,bytes32)'),
        0,
      )
      const { events } = await tx.wait()

      assert.equal(4, events?.length)
      assert.equal(events?.[0].event, 'PluginRequested')
    })

    it('emits an event on the target oracle contract', async () => {
      const tx = await cc.publicRequestOracleDataFrom(
        newoc.address,
        specId,
        ethers.utils.toUtf8Bytes('fulfillRequest(bytes32,bytes32)'),
        0,
      )
      const { logs } = await tx.wait()
      const event = logs && newoc.interface.parseLog(logs[3])

      assert.equal(4, logs?.length)
      assert.equal(event?.name, 'OracleRequest')
    })

    it('does not modify the stored oracle address', async () => {
      await cc.publicRequestOracleDataFrom(
        newoc.address,
        specId,
        ethers.utils.toUtf8Bytes('fulfillRequest(bytes32,bytes32)'),
        0,
      )

      const actualOracleAddress = await cc.publicOracleAddress()
      assert.equal(oc.address, actualOracleAddress)
    })
  })

  describe('#cancelPluginRequest', () => {
    let requestId: string
    // a concrete plugin attached to an empty oracle
    let ecc: Contract

    beforeEach(async () => {
      const emptyOracle = await emptyOracleFactory
        .connect(roles.defaultAccount)
        .deploy()
      ecc = await concretePluginClientFactory
        .connect(roles.defaultAccount)
        .deploy(pli.address, emptyOracle.address)

      const tx = await ecc.publicRequest(
        specId,
        ecc.address,
        ethers.utils.toUtf8Bytes('fulfillRequest(bytes32,bytes32)'),
        0,
      )
      const { events } = await tx.wait()
      requestId = (events?.[0]?.args as any).id
    })

    it('emits an event from the contract showing the run was cancelled', async () => {
      const tx = await ecc.publicCancelRequest(
        requestId,
        0,
        ethers.utils.hexZeroPad('0x', 4),
        0,
      )
      const { events } = await tx.wait()

      assert.equal(1, events?.length)
      assert.equal(events?.[0].event, 'PluginCancelled')
      assert.equal(requestId, (events?.[0].args as any).id)
    })

    it('throws if given a bogus event ID', async () => {
      await evmRevert(
        ecc.publicCancelRequest(
          ethers.utils.formatBytes32String('bogusId'),
          0,
          ethers.utils.hexZeroPad('0x', 4),
          0,
        ),
      )
    })
  })

  describe('#recordPluginFulfillment(modifier)', () => {
    let request: RunRequest

    beforeEach(async () => {
      await oc.setAuthorizedSenders([await roles.defaultAccount.getAddress()])
      const tx = await cc.publicRequest(
        specId,
        cc.address,
        ethers.utils.toUtf8Bytes('fulfillRequest(bytes32,bytes32)'),
        0,
      )
      const { logs } = await tx.wait()

      request = decodeRunRequest(logs?.[3])
    })

    it('emits an event marking the request fulfilled', async () => {
      const tx = await oc
        .connect(roles.defaultAccount)
        .fulfillOracleRequest(
          ...convertFufillParams(
            request,
            ethers.utils.formatBytes32String('hi mom!'),
          ),
        )
      const { logs } = await tx.wait()

      const event = logs && cc.interface.parseLog(logs[1])

      assert.equal(2, logs?.length)
      assert.equal(event?.name, 'PluginFulfilled')
      assert.equal(request.requestId, event?.args.id)
    })

    it('should only allow one fulfillment per id', async () => {
      await oc
        .connect(roles.defaultAccount)
        .fulfillOracleRequest(
          ...convertFufillParams(
            request,
            ethers.utils.formatBytes32String('hi mom!'),
          ),
        )

      await evmRevert(
        oc
          .connect(roles.defaultAccount)
          .fulfillOracleRequest(
            ...convertFufillParams(
              request,
              ethers.utils.formatBytes32String('hi mom!'),
            ),
          ),
        'Must have a valid requestId',
      )
    })

    it('should only allow the oracle to fulfill the request', async () => {
      await evmRevert(
        oc
          .connect(roles.stranger)
          .fulfillOracleRequest(
            ...convertFufillParams(
              request,
              ethers.utils.formatBytes32String('hi mom!'),
            ),
          ),
        'Not authorized sender',
      )
    })
  })

  describe('#fulfillPluginRequest(function)', () => {
    let request: RunRequest

    beforeEach(async () => {
      await oc.setAuthorizedSenders([await roles.defaultAccount.getAddress()])
      const tx = await cc.publicRequest(
        specId,
        cc.address,
        ethers.utils.toUtf8Bytes(
          'publicFulfillPluginRequest(bytes32,bytes32)',
        ),
        0,
      )
      const { logs } = await tx.wait()

      request = decodeRunRequest(logs?.[3])
    })

    it('emits an event marking the request fulfilled', async () => {
      await oc.setAuthorizedSenders([await roles.defaultAccount.getAddress()])
      const tx = await oc
        .connect(roles.defaultAccount)
        .fulfillOracleRequest(
          ...convertFufillParams(
            request,
            ethers.utils.formatBytes32String('hi mom!'),
          ),
        )

      const { logs } = await tx.wait()
      const event = logs && cc.interface.parseLog(logs[1])

      assert.equal(2, logs?.length)
      assert.equal(event?.name, 'PluginFulfilled')
      assert.equal(request.requestId, event?.args?.id)
    })

    it('should only allow one fulfillment per id', async () => {
      await oc
        .connect(roles.defaultAccount)
        .fulfillOracleRequest(
          ...convertFufillParams(
            request,
            ethers.utils.formatBytes32String('hi mom!'),
          ),
        )

      await evmRevert(
        oc
          .connect(roles.defaultAccount)
          .fulfillOracleRequest(
            ...convertFufillParams(
              request,
              ethers.utils.formatBytes32String('hi mom!'),
            ),
          ),
        'Must have a valid requestId',
      )
    })

    it('should only allow the oracle to fulfill the request', async () => {
      await evmRevert(
        oc
          .connect(roles.stranger)
          .fulfillOracleRequest(
            ...convertFufillParams(
              request,
              ethers.utils.formatBytes32String('hi mom!'),
            ),
          ),
        'Not authorized sender',
      )
    })
  })

  describe('#pluginToken', () => {
    it('returns the Pli Token address', async () => {
      const addr = await cc.publicPluginToken()
      assert.equal(addr, pli.address)
    })
  })

  describe('#addExternalRequest', () => {
    let mock: Contract
    let request: RunRequest

    beforeEach(async () => {
      mock = await concretePluginClientFactory
        .connect(roles.defaultAccount)
        .deploy(pli.address, oc.address)

      const tx = await cc.publicRequest(
        specId,
        mock.address,
        ethers.utils.toUtf8Bytes('fulfillRequest(bytes32,bytes32)'),
        0,
      )
      const receipt = await tx.wait()

      request = decodeRunRequest(receipt.logs?.[3])
      await mock.publicAddExternalRequest(oc.address, request.requestId)
    })

    it('allows the external request to be fulfilled', async () => {
      await oc.setAuthorizedSenders([await roles.defaultAccount.getAddress()])
      await oc.fulfillOracleRequest(
        ...convertFufillParams(
          request,
          ethers.utils.formatBytes32String('hi mom!'),
        ),
      )
    })

    it('does not allow the same requestId to be used', async () => {
      await evmRevert(
        cc.publicAddExternalRequest(newoc.address, request.requestId),
      )
    })
  })
})
