import RetailUI from './main'
import 'mocha'

let chai = require('chai')
let spies = require('chai-spies')

chai.use(spies)
let expect = chai.expect

const noOpAdapter = () => ({
  fetchInventoryState: () => {},
  addInventoryItem: () => {}
})

describe('saving a new item', () => {
  describe('when item name is empty', () => {
    let mockApiAdapter = noOpAdapter()
    let app = new RetailUI(mockApiAdapter)

    let result = app.addInventoryItem('')

    it('should return error', () => {
      expect(result).to.equal('name empty')
    })
  })

  describe('when item name is already present', () => {
    let mockApiAdapter = noOpAdapter()
    let app = new RetailUI(mockApiAdapter)

    let counter = 0
    mockApiAdapter.addInventoryItem = () => {
      counter++
      return {
        name: 'test'
      }
    }

    app.addInventoryItem('test')
    let result = app.addInventoryItem('test')

    it('should return error', () => {
      expect(counter).to.equal(1)
      expect(result).to.equal('name present')
    })
  })

  describe('when item name is new', () => {
    let mockApiAdapter = noOpAdapter()
    let app = new RetailUI(mockApiAdapter)

    let counter = 0
    let apiNames = ['test', 'test1']
    mockApiAdapter.addInventoryItem = () => {
      let res = {
        name: apiNames[counter]
      }
      counter++
      return res
    }

    app.addInventoryItem('test')
    let result = app.addInventoryItem('test1')

    it('should make the api call', () => {
      expect(counter).to.equal(2)
      expect(result).to.equal('')
    })
  })
})

describe('fetching inventory state', () => {
  let mockApiAdapter = noOpAdapter()

  var mockApi = chai.spy(mockApiAdapter.fetchInventoryState)
  mockApiAdapter.fetchInventoryState = mockApi

  let app = new RetailUI(mockApiAdapter)
  app.fetchInventoryState()

  it('should make the api call', () => {
    expect(mockApi).to.have.been.called.exactly(1)
  })
})
