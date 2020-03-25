import $ = require('jquery')

import InventoryClient from '../client/inventory'
import Inventory from '../retailapp/inventory'
import StockClient from '../client/stock'

interface inventoryItem {
  id: number
  name: string
}

export function initializeInventory(client: InventoryClient, stock: StockClient) {
  let nameInput: JQuery<HTMLElement> = $('#name')
  let form: JQuery<HTMLElement> = $('#mainForm')

  let page = new Inventory(client, stock)

  page.name = (newVal: string) => manageName(nameInput, newVal)
  page.nameError = (v: boolean) => showEmptyNameError(v)
  page.uniqueError = (v: boolean) => showUniqueNameError(v)
  page.addBtnEnabled = (v: boolean) => setAddBtnState(v)
  page.renderTable = (data: inventoryItem[]) => populateTable(data)

  form.on('submit', (e: Event) => {
    e.preventDefault()
    page.onSubmit()
  })

  nameInput.on('keyup', () => {
    page.onNameChange()
  })

  page.init()
}

function manageName(nameInput: JQuery<HTMLElement>, newVal: string): string {
  if (newVal !== undefined) {
    nameInput.val(newVal)
  }

  return <string>nameInput.val()
}

function setAddBtnState(v: boolean): boolean {
  let btn = $('button#inventoryPageSubmitButton.btn.btn-secondary')
  if (v === false) {
    btn.prop('disabled', true)
    return v
  }
  if (v === true) {
    btn.prop('disabled', false)
    return v
  }

  return !btn.is(':disabled')
}

function showUniqueNameError(v: boolean): void {
  if (v) {
    $('#unique.invalid-feedback').addClass('d-block')
    return
  }
  $('#unique.invalid-feedback').removeClass('d-block')
}

function showEmptyNameError(v: boolean): void {
  if (v) {
    $('#nonempty.invalid-feedback').addClass('d-block')
    return
  }
  $('#nonempty.invalid-feedback').removeClass('d-block')
}

function onTableRowClick_highlight_row(): void {
  $('#inventoryTable tbody').on('click', 'tr', function() {
    $(this).toggleClass('list-group-item-dark')
  })
}

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
