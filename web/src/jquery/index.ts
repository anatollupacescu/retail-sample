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

  //register page listeners and load initial data
  initializeInventory(inventory)

  let recipe = new RecipeClient(apiEndpoint)
  initializeRecipe(inventory, recipe)

  let stock = new StockClient(apiEndpoint)
  initializeStock(inventory, stock)

  let order = new OrderClient(apiEndpoint)

  let app = new RetailApp(stock, order, recipe)
  initializeOrder(app, recipe, order)
})
