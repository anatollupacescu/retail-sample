import 'mocha'
import chai = require('chai')
import spies = require('chai-spies')
import InventoryClient from '../client/inventory'
import Inventory from './inventory'
import StockClient from '../client/stock'

chai.use(spies)
let expect = chai.expect

describe('init', () => {
  // let client = new InventoryClient()
  // let stock = new StockClient()

  // let page = new Inventory(client, stock)

  it('should fetch state and render table', () => {})
})

describe('on name change', () => {
  describe('when new value is empty', () => {
    it('should disable the ADD button', () => {})
  })
})

describe('on submit', () => {
  describe('when new name is empty', () => {
    let client = new InventoryClient()
    let stock = new StockClient()
    let page = new Inventory(client, stock)
    //  let addItem = chai.spy.on(client, 'addItem', () => [{ id: 1, name: 'test' }])
    let getNameValue = chai.spy.on(page, 'getNameValue', () => '')
    let nameError = chai.spy.on(page, 'nameError')
    let addBtnEnabled = chai.spy.on(page, 'addBtnEnabled')

    // let addStockPos = chai.spy.on(stock, "addPosition")

    it('should show error ', () => {
      page.onSubmit()
      expect(getNameValue).to.have.been.called
      expect(nameError).to.have.been.called.with(true)
      expect(addBtnEnabled).to.have.been.called.with(false)
      // expect(addStockPos).to.have.been.called
    })

    it('should disable ADD button', () => {})
  })

  xdescribe('when ', () => {})
})
