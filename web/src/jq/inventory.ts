import $ = require('jquery')

import Client, { inventoryItem } from '../app/inventory/client'
import App, { Page } from '../app/inventory/app'
import StockClient from '../app/stock/client'

export function initializeInventory(client: Client, stock: StockClient) {
  let nameInput = $('#name')
  let form = $('#mainForm')

  let page: Page = {
    highlightRow: (id: string) => selectTableRow(id),
    clearRow: () => clearRowSelection(),
    populateModal: (i: inventoryItem) => populateModal(i),
    toggleModal: (v: boolean) => toggleModal(v),
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

  let openModalBtn = $('#openInventoryModalBtn')

  openModalBtn.on('click', () => {
    app.showModal()
  })

  let closeModalClass = $('.close-modal')

  closeModalClass.on('click', () => {
    app.closeModal()
  })

  let table = $('#inventoryTable tbody')

  table.on('click', 'tr', function () {
    let id = $(this).find('td:eq(0)').text()
    app.onRowClick(id)
  })

  let enableItem = $('#enabledItem')

  enableItem.on('click', () => {
    app.toggleItemStatus(true)
  })

  let disableItem = $('#disabledItem')

  disableItem.on('click', () => {
    app.toggleItemStatus(false)
  })

  app.init()
}

let dark = 'list-group-item-dark'

function clearRowSelection(): void {
  $('#inventoryTable tbody tr').removeClass(dark)
}

function selectTableRow(id: string): void {
  $('#inventoryTable tbody tr').each(function () {
    let currentRow = $(this)
    let currentID = currentRow.find('td:eq(0)').text()

    if (id === currentID) {
      currentRow.addClass(dark)
    }
  })
}

function populateModal(i: inventoryItem): void {
  $('#inventoryModalID').html(String(i.id))
  $('#inventoryModalName').html(i.name)

  let enbl = $('#enabledItem'),
    dsbl = $('#disabledItem')

  enbl.removeClass('btn-primary')
  dsbl.removeClass('btn-primary')

  if (i.enabled) {
    enbl.addClass('btn-primary')
    return
  }

  dsbl.addClass('btn-primary')
}

function toggleModal(v: boolean): void {
  let el = $('#inventoryModal')
  if (v) {
    el.addClass('show')
    el.addClass('d-block')
    return
  }
  el.removeClass('show')
  el.removeClass('d-block')
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

const byID = (i1: inventoryItem, i2: inventoryItem) => i1.id - i2.id

function populateTable(data: inventoryItem[]) {
  let table = <HTMLTableElement>$('#inventoryTable tbody')[0]

  $('#inventoryTable tbody tr').remove()
  data.sort(byID).forEach((element: inventoryItem) => {
    let row = <HTMLTableRowElement>table.insertRow(0)
    row.insertCell(0).innerHTML = `${element.id}`
    row.insertCell(1).innerHTML = element.name
    if (!element.enabled) {
      row.classList.add('disabled')
    }
  })
}
