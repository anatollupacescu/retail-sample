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
    populateTable(inv.getInventory())
  })
}

function onSaveNewItem_submit(app: RetailApp, inv: InventoryClient): void {
  $('#mainForm').on('submit', function(e) {
    e.preventDefault()

    var data = <string>nameInput.val()

    app
      .addInventoryItem(data)
      .then(rsp => {
        switch (rsp) {
          case 'ERR_EMPTY':
          case 'name empty':
            $('#nonempty.invalid-feedback').addClass('d-block')
            return
          case 'ERR_UNIQUE':
          case 'name present':
            $('#unique.invalid-feedback').addClass('d-block')
            return
          default:
            populateTable(inv.getInventory())
            nameInput.val('')
        }
      })
      .catch(err => {
        console.error(err)
      })
  })
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

function populateTable(items: inventoryItem[]) {
  $('#inventoryTable tbody tr').remove()
  let table = <HTMLTableElement>$('#inventoryTable tbody')[0]
  let rows = items.sort((i1: inventoryItem, i2: inventoryItem) => {
    return i1.id - i2.id
  })
  rows.forEach((element: inventoryItem) => {
    let row = <HTMLTableRowElement>table.insertRow(0)
    row.insertCell(0).innerHTML = `${element.id}`
    row.insertCell(1).innerHTML = element.name
  })
}
