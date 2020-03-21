import $ = require('jquery')

import InventoryClient, { inventoryItem } from '../client/inventory'
import StockClient, { Position } from '../client/stock'

let qtyInput: JQuery<HTMLElement>, addBtn: JQuery<HTMLElement>, nameInput: JQuery<HTMLElement>

export function initializeStock(inv: InventoryClient, stock: StockClient) {
  qtyInput = $('#provisionQty')
  nameInput = $('#inventoryItemType')
  addBtn = $('#provisionStock')

  onTabClick_populateDropdown_populateTable(inv, stock)
  onProvisionBtnClick_doProvision(inv, stock)

  $('#provisionStock').on('click', () => {
    resetCount()
  })

  stock.fetchState().then(() => {
    populateTable(inv, stock)
  })
}

function onProvisionBtnClick_doProvision(inv: InventoryClient, stock: StockClient): void {
  addBtn.on('click', () => {
    let qty = <number>qtyInput.val()
    if (!qty || Number(qty) <= 0) {
      setErrorMessage()
      return
    }
    let id = <string>nameInput.val()
    stock
      .provision(id, qty)
      .then(res => {
        if (res.status !== 201) {
          console.error('got error provisioning')
          return
        }
        resetErrorMessage()
        populateTable(inv, stock)
      })
      .catch(err => {
        console.error(err)
      })
  })
}

function setErrorMessage(): void {
  $('#provisionQtyErr.invalid-feedback').addClass('d-block')
}

function resetErrorMessage(): void {
  $('#provisionQtyErr.invalid-feedback').removeClass('d-block')
}

function onTabClick_populateDropdown_populateTable(inv: InventoryClient, stock: StockClient): void {
  $('#stock-tab').on('click', () => {
    populateDropdown(inv)
    populateTable(inv, stock)
  })
}

function populateDropdown(inv: InventoryClient): void {
  nameInput.empty()
  inv.getInventory().map(item => {
    nameInput.append(new Option(item.name, String(item.id)))
  })
}

function resetCount(): void {
  qtyInput.val('')
}

function populateTable(inventory: InventoryClient, stock: StockClient): void {
  $('#stock tbody tr').remove()
  let table = <HTMLTableElement>$('#stock tbody')[0]

  let items: Position[] = stock.getData()
  let stockDict: StockDict = toDict(items)

  let inventoryItems = inventory.getInventory()
  let rows = inventoryItems.sort((i1: inventoryItem, i2: inventoryItem) => i1.id - i2.id)

  rows.forEach((element: inventoryItem) => {
    let row = <HTMLTableRowElement>table.insertRow(0)

    row.insertCell(0).innerHTML = String(element.id)
    row.insertCell(1).innerHTML = String(inventory.getName(element.id))

    let value: string = stockDict[element.id]

    if (!value) {
      value = '0'
    }

    row.insertCell(2).innerHTML = value
  })
}

type StockDict = Record<string, string>

function toDict(i: Position[]): StockDict {
  let r: StockDict = {}
  i.forEach(e => {
    r = {
      [e.id]: String(e.qty),
      ...r
    }
  })
  return r
}
