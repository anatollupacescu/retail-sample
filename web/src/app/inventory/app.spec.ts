import 'mocha'
import chai = require('chai')
import spies = require('chai-spies')
import Client, { inventoryItem } from './client'
import Inventory, { Page } from './app'
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
    let client = new Client()
    let stock = new StockClient()

    let page: Page = noOpPage()

    let getNameValue = chai.spy.on(page, 'name', () => '')
    let nameError = chai.spy.on(page, 'toggleNameError')
    let addBtnEnabled = chai.spy.on(page, 'addBtnEnabled')

    let app = new Inventory(client, stock, page)
    app.onSubmit()

    it('should show error', () => {
      expect(getNameValue).to.have.been.called.once
      expect(nameError).to.have.been.called.with(true)
    })

    it('should disable ADD button', () => {
      expect(addBtnEnabled).to.have.been.called.with(false)
    })
  })

  xdescribe('when ', () => {})
})

function noOpPage(): Page {
  return {
    toggleNameError: (_v: boolean): void => {},
    name: (): string => {
      return ''
    },
    resetName: (): void => {},
    toggleUniqueError: (_v: boolean): void => {},
    addBtnEnabled: (_v: boolean): void => {},
    renderTable: (_data: inventoryItem[]): void => {}
  }
}
