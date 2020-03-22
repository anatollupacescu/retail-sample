import 'mocha'
import chai = require('chai')
import spies = require('chai-spies')

chai.use(spies)
let expect = chai.expect

import RetailApp from './app'
import OrderClient from '../client/order'
import StockClient from '../client/stock'
import RecipeClient, { Recipe } from '../client/recipe'
import InventoryClient from '../client/inventory'

describe('add inventory item', () => {
  let inv = new InventoryClient(''),
    order = new OrderClient(''),
    recipe = new RecipeClient(''),
    stock = new StockClient('')

  let app = new RetailApp(stock, order, recipe, inv)

  let addItem = chai.spy.on(inv, 'addItem', () => [{ id: 1 }, { id: 2 }])
  let addPosition = chai.spy.on(stock, 'addPosition', () => {})

  it('sets stock position to zero', async () => {
    await app.addInventoryItem('test')
    expect(addItem).to.have.been.called
    expect(addPosition).to.have.been.called.exactly(2)
  })
})

describe('place order', () => {
  let order = new OrderClient(''),
    recipe = new RecipeClient(''),
    stock = new StockClient(''),
    inv = new InventoryClient('')

  let app = new RetailApp(stock, order, recipe, inv)

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

  let substractFromPosition = chai.spy.on(stock, 'substractFromPosition')

  it('updates the stock', async () => {
    await app.placeOrder(1, 1)
    expect(addOrder).to.have.been.called
    expect(getRecipe).to.have.been.called
    expect(getRecipe).to.have.been.called
    expect(substractFromPosition).to.have.been.called.exactly(1)
  })
})
