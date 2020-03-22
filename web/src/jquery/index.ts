import $ = require('jquery')
import { apiIsHealthy } from '../health'

import InventoryClient from '../client/inventory'
import { initializeInventory } from './inventory'

import RecipeClient from '../client/recipe'
import { initializeRecipe } from './recipe'

import { initializeStock } from './stock'
import StockClient from '../client/stock'

import { initializeOrder } from './order'
import OrderClient from '../client/order'

import RetailApp from '../retailapp/app'

$(document).ready(async () => {
  let apiUrl = process.env.API_URL
  let apiPort = process.env.API_PORT
  let diagPort = process.env.DIAG_PORT

  if (!apiUrl || !apiPort || !diagPort) {
    console.error('missing configuration')
    return
  }

  let diagEndpoint = `${apiUrl}:${diagPort}`
  let apiStatus = await apiIsHealthy(diagEndpoint)

  if (!apiStatus) {
    console.error('diagnostic check failed', diagEndpoint)
    return
  }

  let apiEndpoint = `${apiUrl}:${apiPort}`
  let inventory = new InventoryClient(apiEndpoint)
  let recipe = new RecipeClient(apiEndpoint)

  let order = new OrderClient(apiEndpoint)
  let stock = new StockClient(apiEndpoint)
  let app = new RetailApp(stock, order, recipe, inventory)

  //register page listeners and load initial data
  initializeInventory(app, inventory)

  initializeRecipe(inventory, recipe)

  initializeStock(inventory, stock)

  initializeOrder(app, recipe, order)
})
