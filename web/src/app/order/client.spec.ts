import 'mocha'
import chai = require('chai')
import spies = require('chai-spies')

import chaiAsPromised = require('chai-as-promised')
chai.use(chaiAsPromised)

import axios from 'axios'
import Client from './client'

chai.use(spies)
let expect = chai.expect

describe('provision stock', () => {
  describe('when quantity not provided', () => {
    let app = new Client(axios)

    it('should fail', () => {
      expect(app.addOrder(1, 0)).to.be.rejectedWith('quantity mandatory')
    })
  })

  describe('when all good', () => {
    let app = new Client(axios)
    let mockApi = chai.spy.on(app, 'apiAddOrder', () => 1)

    it('calls the api and stores the state locally', async () => {
      let res = await app.addOrder(1, 2)
      expect(res).to.equal(1)
      expect(mockApi).to.have.been.called.exactly(1)
      expect(app.getState()).to.have.length(1)
    })
  })
})
