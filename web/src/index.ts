import $ = require('jquery')
import { apiIsHealthy } from './health'

import InventoryClient from './retailapp/inventory'
import { initializeInventory } from './listener/inventory'

import RecipeClient from './retailapp/recipe'
import { initializeRecipe } from './listener/recipe'

$(document).ready(async () => {
  $('a[data-toggle="tab"]').on('click', function(e) {
    console.log('current tab', e.target.id) // newly activated tab
  })
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
})
