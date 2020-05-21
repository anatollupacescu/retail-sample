import $ = require('jquery')
import { apiIsHealthy } from '../health'

import InventoryClient from '../app/inventory/client'
import { initializeInventory } from './inventory'

import RecipeClient from '../app/recipe/client'
import { initializeRecipe } from './recipe'

import { initializeStock } from './stock'
import StockClient from '../app/stock/client'

import { initializeOrder } from './order'
import OrderClient from '../app/order/client'

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

  $('#alert').hide()

  let apiEndpoint = `${apiUrl}:${apiPort}`
  let inventory = new InventoryClient(`${apiEndpoint}/inventory`)
  let recipe = new RecipeClient(`${apiEndpoint}/recipe`)
  let order = new OrderClient(`${apiEndpoint}/order`)
  let stock = new StockClient(`${apiEndpoint}/stock`)

  initializeInventory(inventory, stock)

  initializeRecipe(inventory, recipe)

  initializeStock(inventory, stock)

  initializeOrder(stock, recipe, order)
})
