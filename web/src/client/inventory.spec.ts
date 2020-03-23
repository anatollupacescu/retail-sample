import 'mocha'
import chai = require('chai')
import spies = require('chai-spies')

import InventoryClient from './inventory'

chai.use(spies)
let expect = chai.expect

describe('saving a new item', () => {
  describe('when item name is empty', () => {
    let app = new InventoryClient('')
    let mockApi = chai.spy.on(app, 'apiAddItem')

    it('should err', async () => {
      // no done
      let result = await app.addItem('')
      expect(mockApi).to.not.have.been.called
      expect(result[1]).to.equal('name empty')
      expect(app.getInventory()).to.have.length(0)
    })
  })

  describe('when server says item name is empty', () => {
    let app = new InventoryClient('')
    let mockApi = chai.spy.on(app, 'apiAddItem')

    it('should err', async () => {
      // no done
      let result = await app.addItem('')
      expect(mockApi).to.not.have.been.called
      expect(result[1]).to.equal('name empty')
      expect(app.getInventory()).to.have.length(0)
    })
  })

  describe('when item name is already present', () => {
    let presentItem = {
      id: 1,
      name: 'test'
    }

    let app = new InventoryClient('', [presentItem])

    let mockApi = chai.spy.on(app, 'apiAddItem')

    it('should err without calling api', async () => {
      let result = await app.addItem('test')
      expect(mockApi).to.have.not.been.called
      expect(result[1]).to.equal('name present')
      expect(app.getInventory()).to.have.length(1)
    })
  })

  describe('when server says item name is already present', () => {
    let app = new InventoryClient('')

    let mockApi = chai.spy.on(app, 'apiAddItem', () => [[], 'ERR_UNIQUE'])

    it('Should err', async () => {
      let result = await app.addItem('test')
      expect(mockApi).to.have.been.called.exactly(1)
      expect(result[1]).to.equal('name present')
      expect(app.getInventory()).to.be.empty
    })
  })

  describe('when item name is unique', () => {
    let initialData = [
      {
        id: 1,
        name: 'test1'
      }
    ]
    let app = new InventoryClient('', initialData)

    let serverResponse = {
      id: 2,
      name: 'test2'
    }
    let mockApi = chai.spy.on(app, 'apiAddItem', () => [[serverResponse], ''])

    it('should make the api call', async () => {
      let result = await app.addItem('test2')
      expect(mockApi).to.have.been.called
      expect(app.getInventory()).to.have.length(2)
    })
  })
})

describe('fetching inventory state', () => {
  let app = new InventoryClient('')

  var mockApi = chai.spy.on(app, 'apiFetchState', () => ({
    data: {
      data: [
        {
          name: 'item1',
          id: 2
        }
      ]
    }
  }))

  it('should make the api call', async () => {
    await app.fetchState()
    expect(mockApi).to.have.been.called.exactly(1)
    expect(app.getInventory()).to.have.length(1)
  })
})
