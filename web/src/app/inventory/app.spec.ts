import 'mocha'
import chai = require('chai')
import spies = require('chai-spies')
import { Page } from './app'
import { inventoryItem } from './client'

chai.use(spies)
let expect = chai.expect

describe('on name change', () => {
  describe('when new value is empty', () => {
    it('should disable the ADD button', () => {})
  })
})

describe('on submit', () => {
  it('test', () => {
    noOpPage()
    expect(true).to.be.true
  })
})

function noOpPage(): Page {
  return {
    clearRow: (): void => {},
    highlightRow: (_: string): void => {},
    populateModal: (_: inventoryItem): void => {},
    toggleModal: (_: boolean): void => {},
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
