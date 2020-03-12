import RetailUI from './main'
import { expect } from 'chai'
import 'mocha'

describe('First test', () => {
  let counter = 0
  let adapters = {
    fetchInventoryState: () => {
      counter++
    }
  }

  let app = new RetailUI('', adapters)
  app.fetchInventoryState()

  it('should return true', () => {
    expect(counter).to.equal(1)
  })
})
