import 'mocha'
import chai = require('chai')
import spies = require('chai-spies')

import OrderClient from './client'

chai.use(spies)
let expect = chai.expect

describe('provision stock', () => {
  describe('when quantity not provided', () => {
    let app = new OrderClient('')

    it('should fail', async () => {
      let res = await app.addOrder(1, 0)
      expect(res).to.equal('quantity mandatory')
    })
  })

  describe('when all good', () => {
    let app = new OrderClient('')
    let mockApi = chai.spy.on(app, 'apiAddOrder', () => ({
      status: 201,
      data: {
        data: {
          '1': 2
        }
      }
    }))

    it('calls the api and stores the state locally', async () => {
      await app.addOrder(1, 2)
      expect(mockApi).to.have.been.called
      expect(app.getOrders()).to.have.length(1)
    })
  })
})
