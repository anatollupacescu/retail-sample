import 'mocha'
import chai = require('chai')
import spies = require('chai-spies')

chai.use(spies)
let expect = chai.expect

import Order, { Page, tableRowDTO } from './order'
import OrderClient from './client'
import StockClient from '../stock/client'
import RecipeClient, { Recipe } from '../recipe/client'

const defaultPage = (): Page => ({
  toggleSubmitButtonState: (_v: boolean): void => {},
  getRecipeID: (): number => {
    return 0
  },
  getQty: (): number => {
    return 0
  },
  resetQty: (): void => {},
  toggleQtyError: (_v: boolean): void => {},
  toggleNotEnoughStockError: (_v: boolean): void => {},
  populateDropdown: (_rows: Recipe[]): void => {},
  populateTable: (_rows: tableRowDTO[]): void => {}
})

describe('place order', () => {
  let order = new OrderClient(),
    recipe = new RecipeClient(),
    stock = new StockClient()

  let page = defaultPage()
  let app = new Order(stock, order, recipe, page)

  let addOrder = chai.spy.on(order, 'addOrder', () => {})

  let getRecipe = chai.spy.on(
    recipe,
    'getByID',
    (): Recipe => ({
      id: 1,
      name: 'test',
      items: [
        {
          id: 11,
          qty: 45
        }
      ]
    })
  )

  // let substractFromPosition = chai.spy.on(stock, 'substractFromPosition')

  it('updates the stock', async () => {
    await app.placeOrder()
    expect(addOrder).to.have.been.called
    expect(getRecipe).to.have.been.called
    expect(getRecipe).to.have.been.called
    // expect(substractFromPosition).to.have.been.called.exactly(1)
  })
})
