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

import axios from 'axios'

$(async () => {
  let apiURL = process.env.API_URL

  if (!apiURL) {
    console.error('missing api URL')
    return
  }

  let healthcheckURL = `${apiURL}/health`
  let apiStatus = await apiIsHealthy(healthcheckURL)

  if (!apiStatus) {
    return
  }

  $('#alert').hide()

  let httpClient = axios.create({
    baseURL: apiURL
  })

  let inventory = new InventoryClient(httpClient)
  let recipe = new RecipeClient(httpClient)
  let order = new OrderClient(httpClient)
  let stock = new StockClient(httpClient)

  initializeInventory(inventory, stock)

  initializeRecipe(inventory, recipe)

  initializeStock(inventory, stock)

  initializeOrder(stock, recipe, order)
})
