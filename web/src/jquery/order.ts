import $ = require('jquery')

import OrderClient, { Record } from '../app/order/client'
import Order, { Page, tableRowDTO } from '../app/order/order'
import RecipeClient, { Recipe } from '../app/recipe/client'
import StockClient from '../app/stock/client'

export function initializeOrder(stock: StockClient, recipe: RecipeClient, order: OrderClient) {
  let recipeInput: JQuery<HTMLElement> = $('#orderRecipe'),
    qtyInput: JQuery<HTMLElement> = $('#orderQty'),
    placeOrderBtn: JQuery<HTMLElement> = $('#placeOrder')

  let page: Page = {
    toggleSubmitButtonState: (v: boolean): void => toggleSubmitButtonState(v),
    getRecipeID: (): number => getNumberValue(recipeInput),
    getQty: (): number => getNumberValue(qtyInput),
    resetQty: (): void => resetInput(qtyInput),
    toggleQtyError: (v: boolean) => toggleQtyErr(v),
    toggleNotEnoughStockError: (v: boolean) => toggleNotEnoughStockErr(v),
    populateDropdown: (rows: Recipe[]): void => populateDropdown(rows, recipeInput),
    populateTable: (rows: tableRowDTO[]): void => populateTable(rows)
  }

  let app = new Order(stock, order, recipe, page)

  recipeInput.on('change', () => {
    app.recipeInputChanged()
  })

  qtyInput.on('change', () => {
    app.qtyInputChanged()
  })

  placeOrderBtn.on('click', () => {
    app.placeOrder()
  })

  $('#order-tab').on('click', () => {
    app.show()
  })

  app.init()
}

function toggleSubmitButtonState(v: boolean): void {
  let btn = $('button#placeOrder.btn.btn-secondary')
  btn.prop('disabled', !v)
}

function toggleNotEnoughStockErr(v: boolean): void {
  if (v) {
    showNotEnoughStockErr()
    return
  }
  resetNotEnoughStockErr()
}

function toggleQtyErr(v: boolean): void {
  if (v) {
    showQtyError()
    return
  }
  resetQtyError()
}

function getNumberValue(input: JQuery<HTMLElement>): number {
  return Number(input.val())
}

function populateDropdown(recipes: Recipe[], input: JQuery<HTMLElement>): void {
  input.empty()
  recipes.map(item => {
    input.append(new Option(item.name, String(item.id)))
  })
}

const byRecipeID = (i1: tableRowDTO, i2: tableRowDTO) => Number(i1.id) - Number(i2.id)

function populateTable(data: tableRowDTO[]): void {
  let rows = data.sort(byRecipeID)
  let table = <HTMLTableElement>$('#orderTable tbody')[0]
  $('#orderTable tbody tr').remove()
  rows.forEach((element: tableRowDTO) => {
    let row = <HTMLTableRowElement>table.insertRow(0)
    row.insertCell(0).innerHTML = element.id
    row.insertCell(1).innerHTML = element.recipe
    row.insertCell(2).innerHTML = element.qty
  })
}

function placeOrder(app: Order, order: OrderClient): void {
  /*
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
    */
}

function resetInput(input: JQuery<HTMLElement>): void {
  input.val('')
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
