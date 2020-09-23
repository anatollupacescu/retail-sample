import 'mocha'
import chai = require('chai')
import spies = require('chai-spies')

import App, { Page, recipeDTO, ingredientDTO, optionDTO, Modal } from './app'
import InventoryClient from '../inventory/client'
import Client from './client'
import axios from 'axios'

chai.use(spies)
let expect = chai.expect

describe('add ingredient', () => {
  let inv = new InventoryClient(axios),
    client = new Client(axios)

  describe('when given invalid quantity', () => {
    let page = noOpPage()

    let getName = chai.spy.on(page, 'ingredientQty', () => '')
    let qtyErr = chai.spy.on(page, 'toggleQtyError')

    let modal = noOpModal()

    let app = new App(inv, client, page, modal)

    it('errors', () => {
      app.onAddIngredient()
      expect(getName).to.have.been.called.once
      expect(qtyErr).to.have.been.called.once
    })
  })

  describe('when given negative quantity', () => {
    let page = noOpPage()

    let getName = chai.spy.on(page, 'ingredientQty', () => '-1')
    let qtyErr = chai.spy.on(page, 'toggleQtyError')

    let modal = noOpModal()

    let app = new App(inv, client, page, modal)

    it('errors', async () => {
      app.onAddIngredient()
      expect(getName).to.have.been.called
      expect(qtyErr).to.have.been.called
    })
  })

  describe('when quantity is ok', () => {
    let page = noOpPage(),
      inv = new InventoryClient(axios),
      client = new Client(axios)

    let modal = noOpModal()

    let app = new App(inv, client, page, modal)

    let getName = chai.spy.on(page, 'ingredientQty', () => '1'),
      getID = chai.spy.on(page, 'ingredientID', () => '1'),
      qtyErr = chai.spy.on(page, 'toggleQtyError'),
      invLookup = chai.spy.on(inv, 'findByID', (_v) => ({ name: 'test' }))

    it('adds ingredient to list', () => {
      app.onAddIngredient()
      expect(qtyErr).to.not.have.been.called
      expect(getID).to.have.been.called.once
      expect(invLookup).to.have.been.called.once
      expect(getName).to.have.been.called.once
    })
  })
})

describe('add recipe', () => {
  describe('when name is invalid', () => {
    it('errors', () => {})
  })
})

function noOpPage(): Page {
  return {
    clearRow: (): void => {},
    highlightRow: (_: string): void => {},

    ingredientID: (): number => {
      return 0
    },
    removeIngredientFromDropdown: (_s: string): void => {},
    toggleAddToListBtnState: (_v: boolean): void => {},
    toggleNoIngredientsError: (_v: boolean): void => {},
    toggleNoUniqueNameError: (_v: boolean): void => {},
    recipeName: (): string => {
      return ''
    },
    resetRecipeName: (): void => {},
    toggleAddRecipeButtonState: (_v: boolean): void => {},
    toggleRecipeNameError: (_v: boolean): void => {},
    populateIngredientsDropdown: (_dtos: optionDTO[]): void => {},
    populateIngredientsTable: (_dtos: ingredientDTO[]): void => {},
    populateTable: (_rows: recipeDTO[]): void => {},
    ingredientQty: (): number => {
      return 0
    },
    resetQty: (): void => {},
    toggleQtyError: (_v: boolean): void => {}
  }
}

function noOpModal(): Modal {
  return {
    toggle: (_: boolean): void => {},
    populate: (_: recipeDTO): void => {}
  }
}
