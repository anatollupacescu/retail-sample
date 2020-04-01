import 'mocha'
import chai = require('chai')
import spies = require('chai-spies')

import App, { Page, recipeDTO, ingredientDTO, optionDTO } from './app'
import InventoryClient from '../inventory/client'
import RecipeClient from './client'

chai.use(spies)
let expect = chai.expect

describe('add ingredient', () => {
  let inv = new InventoryClient(),
    client = new RecipeClient()

  describe('when given invalid quantity', () => {
    let page = noOpPage()

    let getName = chai.spy.on(page, 'ingredientQty', () => '')
    let qtyErr = chai.spy.on(page, 'toggleQtyError')

    let app = new App(inv, client, page)
    app.onAddIngredient()

    it('errors', async () => {
      expect(getName).to.have.been.called.once
      expect(qtyErr).to.have.been.called.once
    })
  })

  describe('when given negative quantity', () => {
    let page = noOpPage()

    let getName = chai.spy.on(page, 'ingredientQty', () => '-1')
    let qtyErr = chai.spy.on(page, 'toggleQtyError')

    let app = new App(inv, client, page)
    app.onAddIngredient()

    it('errors', async () => {
      expect(getName).to.have.been.called
      expect(qtyErr).to.have.been.called
    })
  })

  describe('when quantity is ok', () => {
    let page = noOpPage()

    let getName = chai.spy.on(page, 'ingredientQty', () => '1'),
      qtyErr = chai.spy.on(page, 'toggleQtyError'),
      popTable = chai.spy.on(page, 'populateIngredientsTable')

    let app = new App(inv, client, page)
    app.onAddIngredient()

    it('does not error', () => {
      expect(qtyErr).to.not.have.been.called
    })

    it('adds ingredient to list', () => {
      expect(getName).to.have.been.called.once
      expect(popTable).to.have.been.called.once
    })
  })
})

describe('add recipe', () => {

  describe('when name is invalid', () => {

    it('errors', () => {

    })
  })
})

function noOpPage(): Page {
  return {
    toggleAddToListBtnDisabledState: (_v: boolean): void => { },
    ingredientID: (): number => { return 0 },
    removeIngredientFromDropdown: (_s: string): void => { },
    toggleIngredientNameError: (_v: boolean): void => { },
    toggleAddIngredientButtonState: (_v: boolean): void => { },
    toggleNoIngredientsError: (_v: boolean): void => { },
    recipeName: (): string => { return '' },
    resetRecipeName: (): void => { },
    toggleAddRecipeButtonState: (_v: boolean): void => { },
    toggleRecipeNameError: (_v: boolean): void => { },
    populateIngredientsDropdown: (_dtos: optionDTO[]): void => { },
    populateIngredientsTable: (_dtos: ingredientDTO[]): void => { },
    populateTable: (_rows: recipeDTO[]): void => { },
    ingredientQty: (): number => { return 0 },
    resetQty: (): void => { },
    toggleQtyError: (_v: boolean): void => { }
  }
}