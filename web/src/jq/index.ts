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

import createAuth0Client from '@auth0/auth0-spa-js'
import axios from 'axios'

$(async () => {
  let token = null
  /* 
  if (process.env.AUTH == "1") {
    let auth0 = await createAuth0Client({
      domain: process.env.AUTH_DOMAIN || "",
      client_id: process.env.AUTH_CLIENT_ID || ""
    });

    let status = await auth0.isAuthenticated();

    if (!status) {
      let query = window.location.search;
      if (query.includes("code=") && query.includes("state=")) {
        await auth0.handleRedirectCallback();
        window.history.replaceState({}, document.title, "/");
        token = await auth0.getTokenSilently()
      } else {
        await auth0.loginWithRedirect({
          redirect_uri: "http://localhost:8080"
        });
      }
    }

    let user = await auth0.getUser()

    $("#userInfo").text(`You are logged in as ${user.nickname}, `)
    $("#userInfo").append("<a href='#' id='logout'>Logout</a>")
    $("#userInfo").show()

    $("#logout").on("click", ()=>{
      auth0.logout();
    })
  }
 */
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

  let httpClient = axios.create({
    baseURL: `${apiUrl}:${apiPort}`,
    headers: token ? { Authorization: `Bearer ${token}` } : {}
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
