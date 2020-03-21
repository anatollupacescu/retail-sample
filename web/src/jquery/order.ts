import $ = require('jquery')

import OrderClient, { Order } from '../client/order'
import RecipeClient from '../client/recipe'
import RetailApp from '../retailapp/app'

let recipeInput: JQuery<HTMLElement>, qtyInput: JQuery<HTMLElement>, placeOrderBtn: JQuery<HTMLElement>

export function initializeOrder(app: RetailApp, recipe: RecipeClient, order: OrderClient) {
  recipeInput = $('#orderRecipe')
  qtyInput = $('#orderQty')
  placeOrderBtn = $('#placeOrder')

  onQtyChange_resetErr()
  onRecipeInputChange_resetErr()
  onPlaceBtnClick_placeOrder(app, order)
  onTabClick_populateRecipeList(recipe)

  order.fetchOrders().then(() => {
    populateTable(order)
  })
}

function onRecipeInputChange_resetErr(): void {
  recipeInput.on('change', () => {
    resetNotEnoughStockErr()
  })
}

function onQtyChange_resetErr(): void {
  qtyInput.on('change', () => {
    resetQtyError()
    resetNotEnoughStockErr()
  })
}

function onTabClick_populateRecipeList(recipe: RecipeClient) {
  $('#order-tab').on('click', () => {
    populateDropdown(recipe)
  })
}

function populateDropdown(recipe: RecipeClient): void {
  recipeInput.empty()
  recipe.getRecipes().map(item => {
    recipeInput.append(new Option(item.name, String(item.id)))
  })
}

function populateTable(order: OrderClient): void {
  $('#orderTable tbody tr').remove()
  let table = <HTMLTableElement>$('#orderTable tbody')[0]
  let rows = order.getOrders().sort((i1: Order, i2: Order) => i1.recipeID - i2.recipeID)
  rows.forEach((element: Order) => {
    let row = <HTMLTableRowElement>table.insertRow(0)
    row.insertCell(0).innerHTML = String(element.recipeID)
    row.insertCell(1).innerHTML = String(element.qty)
  })
}

function onPlaceBtnClick_placeOrder(app: RetailApp, order: OrderClient): void {
  placeOrderBtn.on('click', () => {
    placeOrder(app, order)
  })
}

function placeOrder(app: RetailApp, order: OrderClient): void {
  let qty = qtyInput.val()
  if (!qty || Number(qty) === 0) {
    showQtyError()
    return
  }

  resetQtyError()

  let recipeID = recipeInput.val()

  app
    .placeOrder(Number(recipeID), Number(qty))
    .then(() => {
      populateTable(order)
      resetNotEnoughStockErr()
      resetQty()
    })
    .catch(res => {
      let errMsg = res.response.data
      if (errMsg && errMsg.startsWith('not enough stock')) {
        showNotEnoughStockErr()
      }
    })
}

function resetQty(): void {
  qtyInput.val('')
}

function showNotEnoughStockErr(): void {
  $('#errStock.invalid-feedback').addClass('d-block')
}

function resetNotEnoughStockErr(): void {
  $('#errStock.invalid-feedback').removeClass('d-block')
}

function showQtyError(): void {
  $('#errQty.invalid-feedback').addClass('d-block')
}

function resetQtyError(): void {
  $('#errQty.invalid-feedback').removeClass('d-block')
}
