import 'mocha'
import chai = require('chai')
import spies = require('chai-spies')
import Client from './client'
import axios from 'axios'

import chaiAsPromised = require('chai-as-promised')
chai.use(chaiAsPromised)

chai.use(spies)
let expect = chai.expect
let sandbox = chai.spy.sandbox()

describe('saving a new item', () => {
  describe('when item name is empty', () => {
    let app = new Client(axios)
    let mockApi = sandbox.on(app, 'apiAddItem')

    it('errors simple', () => {
      expect(app.addItem('')).to.be.rejectedWith('name empty')
      expect(mockApi).to.have.been.called.exactly(0)
      expect(app.getState()).to.have.length(0)
    })
  })

  describe('when server says item name is empty', () => {
    let app = new Client(axios)
    let mockApi = chai.spy.on(app, 'apiAddItem', () => {
      throw 'name not provided'
    })

    it('errors', () => {
      expect(app.addItem('')).to.be.rejectedWith('name empty')
      expect(mockApi).to.have.been.called.exactly(0)
      expect(app.getState()).to.have.length(0)
    })
  })

  describe('when item name is already present', () => {
    let initialData = [
      {
        id: 1,
        name: 'test',
        enabled: true
      }
    ]

    let app = new Client(axios, initialData)

    let mockApi = chai.spy.on(app, 'apiAddItem')

    it('errors with the correct message', () => {
      expect(app.addItem('test')).to.be.rejectedWith('name present')
      expect(mockApi).to.have.been.called.exactly(0)
      expect(app.getState()).to.have.length(1)
    })
  })

  describe('when server says item name is already present', () => {
    let app = new Client(axios)

    let mockApi = chai.spy.on(app, 'apiAddItem', () => {
      throw 'item type already present'
    })

    it('errors with the correct message', () => {
      expect(app.addItem('present')).to.be.rejectedWith('name present')
      expect(mockApi).to.have.been.called.exactly(1)
      expect(app.getState()).to.have.length(0)
    })
  })

  describe('when item name is unique', () => {
    let initialData = [
      {
        id: 1,
        name: 'test1',
        enabled: true
      }
    ]
    let app = new Client(axios, initialData)

    let apiResponse = {
      id: 2,
      name: 'test2'
    }

    let mockApi = chai.spy.on(app, 'apiAddItem', () => apiResponse)

    it('should make the api call', async () => {
      let result = await app.addItem('test2')
      expect(result).to.equal(apiResponse)
      expect(mockApi).to.have.been.called.once
      expect(app.getState()).to.have.length(2)
      expect(app.getState()).to.have.members([...initialData, apiResponse])
    })
  })
})

describe('fetching inventory state', () => {
  let app = new Client(axios)

  let apiItem = {
    id: 1,
    name: 'item1'
  }

  let apiItems = [apiItem]
  var mockApi = chai.spy.on(app, 'apiFetchState', () => apiItems)

  it('should make the api call', async () => {
    let result = await app.fetchState()
    expect(result).to.equal(apiItems)
    expect(mockApi).to.have.been.called.once
    expect(app.getState()).to.have.length(1)
    expect(app.getState()[0]).to.equal(apiItem)
  })
})
