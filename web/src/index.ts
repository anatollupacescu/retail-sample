import $ = require('jquery')
import { registerInventoryListeners } from './listener/inventory'
import { apiIsHealthy } from './health'
import RetailUI from './retailapp/main'

$(document).ready(function () {
  $('a[data-toggle="tab"]').on('click', function (e) {
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
  if (!apiIsHealthy(diagEndpoint)) {
    console.error('diagnostic check failed', diagEndpoint)
    return
  }

  let apiEndpoint = `${apiUrl}:${apiPort}`
  let retailApp = new RetailUI(apiEndpoint)

  //register page listeners...
  registerInventoryListeners(retailApp)

  //fetch initial state
  retailApp.fetchInventoryState()
})
