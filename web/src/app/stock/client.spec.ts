import 'mocha'
import chai = require('chai')
import spies = require('chai-spies')

import StockClient from './client'

chai.use(spies)
let expect = chai.expect

describe('provision stock', () => {
  let data = [
    {
      id: 1,
      qty: 1,
      items: []
    }
  ]

  let app = new StockClient('', data)
  let mockApi = chai.spy.on(app, 'apiProvision', () => ({
    status: 201,
    data: {
      data: {
        '1': 2
      }
    }
  }))

  it('calls the api and stores the state locally', async () => {
    await app.provision('1', 2)
    expect(mockApi).to.have.been.called
    expect(app.getData()).to.have.length(1)
  })
})

describe('fetching state', () => {
  let app = new StockClient('')
  let mockApi = chai.spy.on(app, 'apiFetchState', () => ({
    data: {
      data: [
        {
          id: 1,
          qty: 2
        }
      ]
    }
  }))

  it('calls the api and stores the state locally', async () => {
    await app.fetchState()
    expect(mockApi).to.have.been.called
    expect(app.getData()).to.have.length(1)
  })
})
