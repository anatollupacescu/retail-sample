import 'mocha'
import chai = require('chai')
import spies = require('chai-spies')

import InventoryClient from './inventory'

chai.use(spies)
let expect = chai.expect

describe('saving a new item', () => {
  describe('when item name is empty', () => {
    let app = new InventoryClient('')

    it('Should err', async () => {
      // no done
      let result = await app.addInventoryItem('')
      expect(result).to.equal('name empty')
      expect(app.getInventory()).to.have.length(0)
    })
  })

  describe('when item name is already present', () => {
    let app = new InventoryClient('')

    let mockApi = chai.spy.on(app, 'apiAddInventoryItem', () => ({
      data: {
        data: {
          test: 1
        }
      }
    }))

    it('Should err', async () => {
      await app.addInventoryItem('test')
      let result = await app.addInventoryItem('test')
      expect(mockApi).to.have.been.called.exactly(1)
      expect(result).to.equal('name present')
      expect(app.getInventory()).to.have.length(1)
    })
  })

  describe('when item name is unique', () => {
    let app = new InventoryClient('')

    let mockApiResponses = [
      {
        name: 'test1',
        id: 2
      },
      {
        name: 'test',
        id: 1
      }
    ]

    let mockApi = chai.spy.on(app, 'apiAddInventoryItem', () => ({
      data: {
        data: {
          data: mockApiResponses.pop()
        }
      }
    }))

    app.addInventoryItem('test')
    app.addInventoryItem('test1')

    it('should make the api call', () => {
      expect(mockApi).to.have.been.called.exactly(2)
      expect(app.getInventory()).to.have.length(2)
    })
  })
})

describe('fetching inventory state', () => {
  let app = new InventoryClient('')

  var mockApi = chai.spy.on(app, 'apiFetchInventoryState', () => ({
    data: {
      data: [
        {
          name: 'item1',
          id: 2
        }
      ]
    }
  }))

  app.fetchInventoryState()

  it('should make the api call', () => {
    expect(mockApi).to.have.been.called.exactly(1)
    expect(app.getInventory()).to.have.length(1)
  })
})
