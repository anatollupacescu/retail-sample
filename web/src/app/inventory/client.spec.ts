import 'mocha'
import chai = require('chai')
import spies = require('chai-spies')

import Client from './client'

chai.use(spies)
let expect = chai.expect

describe('saving a new item', () => {
  describe('when item name is empty', () => {
    let app = new Client()
    let mockApi = chai.spy.on(app, 'apiAddItem')

    it('should err', async () => {
      // no done
      let result = await app.addItem('')
      expect(mockApi).to.not.have.been.called
      expect(result[1]).to.equal('name empty')
      expect(app.getState()).to.have.length(0)
    })
  })

  describe('when server says item name is empty', () => {
    let app = new Client()
    let mockApi = chai.spy.on(app, 'apiAddItem')

    it('should err', async () => {
      // no done
      let result = await app.addItem('')
      expect(mockApi).to.not.have.been.called
      expect(result[1]).to.equal('name empty')
      expect(app.getState()).to.have.length(0)
    })
  })

  describe('when item name is already present', () => {
    let initialData = [
      {
        id: 1,
        name: 'test'
      }
    ]

    let app = new Client('', initialData)

    let mockApi = chai.spy.on(app, 'apiAddItem')

    it('should err without calling api', async () => {
      let result = await app.addItem('test')
      expect(mockApi).to.have.not.been.called
      expect(result[1]).to.equal('name present')
      expect(app.getState()).to.have.length(1)
    })
  })

  describe('when server says item name is already present', () => {
    let app = new Client()

    let mockApi = chai.spy.on(app, 'apiAddItem', () => [null, 'ERR_UNIQUE'])

    it('should err', async () => {
      let result = await app.addItem('test')
      expect(mockApi).to.have.been.called.exactly(1)
      expect(result[1]).to.equal('name present')
      expect(app.getState()).to.be.empty
    })
  })

  describe('when item name is unique', () => {
    let initialData = [
      {
        id: 1,
        name: 'test1'
      }
    ]
    let app = new Client('', initialData)

    let apiResponse = {
      id: 2,
      name: 'test2'
    }
    let mockApi = chai.spy.on(app, 'apiAddItem', () => [apiResponse, ''])

    it('should make the api call', async () => {
      let result = await app.addItem('test2')
      expect(mockApi).to.have.been.called
      expect(result[1]).to.be.empty
      expect(app.getState()).to.have.length(2)
      expect(app.getState()).to.have.members([...initialData, apiResponse])
    })
  })
})

describe('fetching inventory state', () => {
  let app = new Client()

  let apiItem = {
    id: 1,
    name: 'item1'
  }

  var mockApi = chai.spy.on(app, 'apiFetchState', () => ({
    data: {
      data: [apiItem]
    }
  }))

  it('should make the api call', async () => {
    await app.fetchState()
    expect(mockApi).to.have.been.called.exactly(1)
    expect(app.getState()).to.have.length(1)
    expect(app.getState()[0]).to.equal(apiItem)
  })
})
