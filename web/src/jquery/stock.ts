import $ = require('jquery')
import Stock, { Page, stockTableRowDTO, inventoryItemDTO } from '../app/stock/stock'
import Client from '../app/inventory/client'
import Client from '../app/stock/client'

export function initializeStock(inventoryClient: Client, stockClient: Client) {
  let qtyInput: JQuery = $('#provisionQty'),
    nameInput: JQuery = $('#inventoryItemType'),
    addBtn: JQuery = $('#provisionStock'),
    tab: JQuery = $('#stock-tab')

  let page: Page = {
    id: (): string => String(nameInput.val()),
    qty: (): number => Number(qtyInput.val()),
    resetQty: (): void => resetValue(qtyInput),
    toggleError: (v: boolean): void => toggleError(v),
    toggleAddBtnState: (v: boolean): void => toggleDisabledState(v, addBtn),
    populateDropdown: (data: inventoryItemDTO[]): void => populateDropdown(data, nameInput),
    populateTable: (data: stockTableRowDTO[]): void => populateTable(data)
  }

  let app = new Stock(inventoryClient, stockClient, page)

  tab.on('click', () => {
    app.show()
  })

  addBtn.on('click', () => {
    app.onProvision()
  })

  qtyInput.on('change', () => {
    app.onQtyChange()
  })

  app.init()
}

function toggleDisabledState(enabled: boolean, input: JQuery): void {
  input.prop('disabled', !enabled)
}

function toggleError(v: boolean): void {
  if (v) {
    setErrorMessage()
    return
  }
  resetErrorMessage()
}

function setErrorMessage(): void {
  $('#provisionQtyErr.invalid-feedback').addClass('d-block')
}

function resetErrorMessage(): void {
  $('#provisionQtyErr.invalid-feedback').removeClass('d-block')
}

function populateDropdown(items: inventoryItemDTO[], nameInput: JQuery): void {
  nameInput.empty()
  items.forEach(item => {
    nameInput.append(new Option(item.name, item.id))
  })
}

function resetValue(qtyInput: JQuery): void {
  qtyInput.val('')
}

function populateTable(data: stockTableRowDTO[]): void {
  let table = <HTMLTableElement>$('#stock tbody')[0]
  let rows = data.sort((i1: stockTableRowDTO, i2: stockTableRowDTO) => i1.id - i2.id)

  $('#stock tbody tr').remove()

  rows.forEach((element: stockTableRowDTO) => {
    let row = <HTMLTableRowElement>table.insertRow(0)

    row.insertCell(0).innerHTML = String(element.id)
    row.insertCell(1).innerHTML = element.name
    row.insertCell(2).innerHTML = element.qty
  })
}
