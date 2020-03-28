import 'mocha'
import chai = require('chai')
import spies = require('chai-spies')
import InventoryClient, { inventoryItem } from './client'
import Inventory, { InventoryPage } from './inventory'
import StockClient from '../stock/client'

chai.use(spies)
let expect = chai.expect

describe('init', () => {
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

    let page: InventoryPage = noOpPage()

    let getNameValue = chai.spy.on(page, 'getNameValue', () => '')
    let nameError = chai.spy.on(page, 'nameError')
    let addBtnEnabled = chai.spy.on(page, 'addBtnEnabled')

    let app = new Inventory(client, stock, page)

    it('should show error ', () => {
      app.onSubmit()
      expect(getNameValue).to.have.been.called
      expect(nameError).to.have.been.called.with(true)
    })

    it('should disable ADD button', () => {
      expect(addBtnEnabled).to.have.been.called.with(false)
    })
  })

  xdescribe('when ', () => {})
})

function noOpPage(): InventoryPage {
  return {
    nameError: (_v: boolean | undefined): void => {},
    getNameValue: (): string => {
      return ''
    },
    setNameEmpty: (): void => {},
    uniqueError: (_v: boolean): void => {},
    addBtnEnabled: (_v: boolean): void => {},
    isAddBtnEnabled: (): boolean => {
      return false
    },
    renderTable: (_data: inventoryItem[]): void => {}
  }
}
