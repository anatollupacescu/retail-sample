import $ = require('jquery')

import InventoryClient from '../app/inventory/client'
import Inventory from '../app/inventory/inventory'
import StockClient from '../app/stock/client'

interface inventoryItem {
  id: number
  name: string
}

export function initializeInventory(client: InventoryClient, stock: StockClient) {
  let nameInput: JQuery<HTMLElement> = $('#name')
  let form: JQuery<HTMLElement> = $('#mainForm')

  let page = {
    getNameValue: () => getNameValue(nameInput),
    setNameEmpty: () => setNameEmpty(nameInput),
    nameError: (v: boolean) => toggleEmptyNameError(v),
    uniqueError: (v: boolean) => toggleUniqueNameError(v),
    addBtnEnabled: (v: boolean) => setAddBtnState(v),
    isAddBtnEnabled: (): boolean => getAddBtnState(),
    renderTable: (data: inventoryItem[]) => populateTable(data)
  }

  let app = new Inventory(client, stock, page)

  form.on('submit', (e: Event) => {
    e.preventDefault()
    app.onSubmit()
  })

  nameInput.on('keyup', () => {
    app.onNameChange()
  })

  app.init()
}

function getNameValue(nameInput: JQuery<HTMLElement>): string {
  return <string>nameInput.val()
}

function setNameEmpty(nameInput: JQuery<HTMLElement>): void {
  nameInput.val('')
}

function getAddBtnState(): boolean {
  let btn = $('button#inventoryPageSubmitButton.btn.btn-secondary')
  return !btn.is(':disabled')
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
