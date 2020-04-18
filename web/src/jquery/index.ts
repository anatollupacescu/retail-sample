import $ = require('jquery')
import { apiIsHealthy } from '../health'

import Client from '../app/inventory/client'
import { initializeInventory } from './inventory'

import Client from '../app/recipe/client'
import { initializeRecipe } from './recipe'

import { initializeStock } from './stock'
import Client from '../app/stock/client'

import { initializeOrder } from './order'
import Client from '../app/order/client'

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
  let inventory = new Client(apiEndpoint)
  let recipe = new Client(apiEndpoint)
  let order = new Client(apiEndpoint)
  let stock = new Client(apiEndpoint)

  initializeInventory(inventory, stock)

  initializeRecipe(inventory, recipe)

  initializeStock(inventory, stock)

  initializeOrder(stock, recipe, order)
})
