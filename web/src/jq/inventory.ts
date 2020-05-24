import $ = require('jquery')

import Client, { inventoryItem } from '../app/inventory/client'
import App, { Page } from '../app/inventory/app'
import StockClient from '../app/stock/client'

export function initializeInventory(client: Client, stock: StockClient) {
  let nameInput = $('#name')
  let form = $('#mainForm')

  let page: Page = {
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

  let openModalBtn = $('#openModalBtn')

  openModalBtn.on('click', () => {
    app.showModal()
  })

  let closeModalClass = $('.close-modal')

  closeModalClass.on('click', () => {
    app.closeModal()
  })

  let tableRow = $('#inventoryTable tbody')

  tableRow.on('click', 'tr', function(){
    //TODO use marker css classes to highligh specific row in the Page component

    let dark = 'list-group-item-dark'
    
    $('#inventoryTable tbody tr').removeClass(dark)
    
    let el = $(this)
    let cells = el[0].cells
    let id = cells[0].innerHTML
    
    app.onRowClick(id)
    
    el.toggleClass(dark)
  })

  app.init()
}

function populateModal(i: inventoryItem): void {
  $('#modalID').html(String(i.id))
  $('#modalName').html(i.name)
  $('#modalStatus').html(String(false))
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
  })
}
