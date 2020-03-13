import RetailUI from './main'
import 'mocha'

let chai = require('chai')
let spies = require('chai-spies')

chai.use(spies)
let expect = chai.expect

describe('saving a new item', () => {
  describe('when item name is empty', () => {
    let app = new RetailUI()

    let result = app.addInventoryItem('')

    it('should return error', () => {
      expect(result).to.equal('name empty')
      expect(app.getInventory()).to.have.length(0)
    })
  })

  describe('when item name is already present', () => {
    let app = new RetailUI()

    let mockApi = chai.spy.on(app, 'apiAddInventoryItem', () => ({
      name: 'test',
      id: 1
    }))

    app.addInventoryItem('test')
    let result = app.addInventoryItem('test')

    it('should return error', () => {
      expect(mockApi).to.have.been.called.once
      expect(result).to.equal('name present')
      expect(app.getInventory()).to.have.length(1)
    })
  })

  describe('when item name is unique', () => {
    let app = new RetailUI()

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

    let mockApi = chai.spy.on(app, 'apiAddInventoryItem', () =>
      mockApiResponses.pop()
    )

    app.addInventoryItem('test')
    let result = app.addInventoryItem('test1')

    it('should make the api call', () => {
      expect(mockApi).to.have.been.called.exactly(2)
      expect(result).to.equal('')
      expect(app.getInventory()).to.have.length(2)
    })
  })
})

describe('fetching inventory state', () => {
  let app = new RetailUI()

  var mockApi = chai.spy.on(app, 'apiFetchInventoryState', () => [
    {
      name: 'item1',
      id: 2
    }
  ])

  app.fetchInventoryState()

  it('should make the api call', () => {
    expect(mockApi).to.have.been.called.exactly(1)
    expect(app.getInventory()).to.have.length(1)
  })
})
