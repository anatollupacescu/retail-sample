import $ = require('jquery')

import InventoryClient from '../client/inventory'
import RetailApp from '../retailapp/app'

let nameInput: JQuery<HTMLElement>

export function initializeInventory(app: RetailApp, inv: InventoryClient) {
  nameInput = $('#name')

  onTableRowClick_highlight_row()
  onChangeNameInput_resetErrorMessage()
  onSaveNewItem_submit(app, inv)

  inv.fetchState().then(() => {
    populateTable(inv)
  })
}

function onSaveNewItem_submit(app: RetailApp, inv: InventoryClient): void {
  $('#mainForm').on('submit', async (e: Event) => {
    e.preventDefault()
    let itemName = <string>nameInput.val()
    doSave(itemName, app, inv)
  })
}

async function doSave(itemName: string, app: RetailApp, inv: InventoryClient): Promise<void> {
  let rsp = await app.addInventoryItem(itemName)
  switch (rsp) {
    case 'name empty':
      $('#nonempty.invalid-feedback').addClass('d-block')
      return
    case 'name present':
      $('#unique.invalid-feedback').addClass('d-block')
      return
    default:
      populateTable(inv)
      nameInput.val('')
  }
}

function onChangeNameInput_resetErrorMessage(): void {
  nameInput.keyup(function() {
    $('#nonempty.invalid-feedback').removeClass('d-block')
    $('#unique.invalid-feedback').removeClass('d-block')
  })
}
function onTableRowClick_highlight_row(): void {
  $('#inventoryTable tbody').on('click', 'tr', function() {
    $(this).toggleClass('list-group-item-dark')
  })
}

interface inventoryItem {
  id: number
  name: string
}

function populateTable(inv: InventoryClient) {
  $('#inventoryTable tbody tr').remove()
  let table = <HTMLTableElement>$('#inventoryTable tbody')[0]

  let byID = (i1: inventoryItem, i2: inventoryItem) => i1.id - i2.id

  let rows = inv.getInventory().sort(byID)

  rows.forEach((element: inventoryItem) => {
    let row = <HTMLTableRowElement>table.insertRow(0)
    row.insertCell(0).innerHTML = `${element.id}`
    row.insertCell(1).innerHTML = element.name
  })
}
