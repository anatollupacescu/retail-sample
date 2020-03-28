import $ = require('jquery')

import OrderClient, { OrderDTO } from '../app/order/client'
import Order from '../app/order/order'
import RecipeClient from '../app/recipe/client'
import StockClient from '../app/stock/client'

let recipeInput: JQuery<HTMLElement>, qtyInput: JQuery<HTMLElement>, placeOrderBtn: JQuery<HTMLElement>

export function initializeOrder(stock: StockClient, recipe: RecipeClient, order: OrderClient) {
  recipeInput = $('#orderRecipe')
  qtyInput = $('#orderQty')
  placeOrderBtn = $('#placeOrder')

  let app = new Order(stock, order, recipe)

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
  let rows = order.getOrders().sort((i1: OrderDTO, i2: OrderDTO) => i1.recipeID - i2.recipeID)
  rows.forEach((element: OrderDTO) => {
    let row = <HTMLTableRowElement>table.insertRow(0)
    row.insertCell(0).innerHTML = String(element.recipeID)
    row.insertCell(1).innerHTML = String(element.qty)
  })
}

function onPlaceBtnClick_placeOrder(app: Order, order: OrderClient): void {
  placeOrderBtn.on('click', () => {
    placeOrder(app, order)
  })
}

function placeOrder(app: Order, order: OrderClient): void {
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
