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
      qty: 1
    }
  ]

  let app = new StockClient('', data)

  let apiProvision = chai.spy.on(app, 'apiProvision', () =>
    Promise.resolve({
      1: 2
    })
  )

  it('calls the api and stores the state locally', async () => {
    await app.provision('1', 2)
    expect(apiProvision).to.have.been.called.once
    expect(app.getData()).to.have.length(1)
    expect(app.getData()[0].qty).to.equal(2)
    expect(app.getData()[0].id).to.equal(1)
  })
})

describe('fetching state', () => {
  let app = new StockClient()
  let mockApi = chai.spy.on(app, 'apiFetchState', () => [
    {
      id: 1,
      qty: 2
    }
  ])

  it('calls the api and stores the state locally', async () => {
    await app.fetchState()
    expect(mockApi).to.have.been.called.once
    expect(app.getData()).to.have.length(1)
    expect(app.getData()[0].qty).to.equal(2)
    expect(app.getData()[0].id).to.equal(1)
  })
})
