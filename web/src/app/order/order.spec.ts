import 'mocha'
import chai = require('chai')
import spies = require('chai-spies')

chai.use(spies)
let expect = chai.expect

import Order from './order'
import OrderClient from './client'
import StockClient from '../stock/client'
import RecipeClient, { Recipe } from '../recipe/client'

describe('place order', () => {
  let order = new OrderClient(),
    recipe = new RecipeClient(),
    stock = new StockClient()

  let app = new Order(stock, order, recipe)

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
