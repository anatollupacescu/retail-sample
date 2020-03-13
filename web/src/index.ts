import $ = require('jquery')
import { initializeInventory } from './listener/inventory'
import { apiIsHealthy } from './health'
import RetailInventory from './retailapp/inventory'

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
  let retailInventory = new RetailInventory(apiEndpoint)

  //register page listeners and load initial data
  initializeInventory(retailInventory)
})
