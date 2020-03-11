import $ = require('jquery')
import { initInventory } from './inventory'
import { apiIsHealthy } from './health'

$(document).ready(function () {
  // $('#myTab a').on('click', function (e) {
  //   e.preventDefault()
  //   console.log('click')
  //   $(this).tab('show')
  // })

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
  initInventory(apiEndpoint)
})
