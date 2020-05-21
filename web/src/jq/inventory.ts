import $ = require('jquery')

import Client, { inventoryItem } from '../app/inventory/client'
import App, { Page } from '../app/inventory/app'
import StockClient from '../app/stock/client'

export function initializeInventory(client: Client, stock: StockClient) {
  let nameInput = $('#name')
  let form = $('#mainForm')

  let page: Page = {
    name: () => getStringValue(nameInput),
    resetName: () => setNameEmpty(nameInput),
    toggleNameError: (v: boolean) => toggleEmptyNameError(v),
    toggleUniqueError: (v: boolean) => toggleUniqueNameError(v),
    addBtnEnabled: (v: boolean) => setAddBtnState(v),
    renderTable: (data: inventoryItem[]) => populateTable(data)
  }

  let app = new App(client, stock, page)

  form.on('submit', (e: Event) => {
    e.preventDefault()
    app.onSubmit()
  })

  nameInput.on('keyup', () => {
    app.onNameChange()
  })

  app.init()
}

function showModal() {
  $('#showModal').on('shown.bs.modal', function () {
    $('#myInput').trigger('focus')
  })
}

function getStringValue(nameInput: JQuery): string {
  return String(nameInput.val())
}

function setNameEmpty(nameInput: JQuery): void {
  nameInput.val('')
}

function setAddBtnState(v: boolean): void {
  let btn = $('button#inventoryPageSubmitButton.btn.btn-secondary')
  btn.prop('disabled', !v)
}

function toggleUniqueNameError(v: boolean): void {
  if (v) {
    $('#unique.invalid-feedback').addClass('d-block')
    return
  }
  $('#unique.invalid-feedback').removeClass('d-block')
}

function toggleEmptyNameError(v: boolean): void {
  if (v) {
    $('#nonempty.invalid-feedback').addClass('d-block')
    return
  }
  $('#nonempty.invalid-feedback').removeClass('d-block')
}

/* for later
function onTableRowClick_highlight_row(): void {
  $('#inventoryTable tbody').on('click', 'tr', function() {
    $(this).toggleClass('list-group-item-dark')
  })
}
*/

const byID = (i1: inventoryItem, i2: inventoryItem) => i1.id - i2.id

function populateTable(data: inventoryItem[]) {
  let table = <HTMLTableElement>$('#inventoryTable tbody')[0]

  $('#inventoryTable tbody tr').remove()
  data.sort(byID).forEach((element: inventoryItem) => {
    let row = <HTMLTableRowElement>table.insertRow(0)
    row.insertCell(0).innerHTML = `${element.id}`
    row.insertCell(1).innerHTML = element.name
  })
}
