import Client, { inventoryItem } from '../inventory/client'
import StockClient from './client'

import { Position } from '../stock/client'

export interface inventoryItemDTO {
  id: string
  name: string
}

export interface stockTableRowDTO {
  id: number
  name: string
  qty: string
}

export interface Page {
  getID(): string
  getQty(): number
  resetQty(): void
  toggleError(v: boolean): void
  populateTable(data: stockTableRowDTO[]): void
  populateDropdown(data: inventoryItemDTO[]): void
}

type StockDict = Record<string, string>

export default class Stock {
  inventory: Client
  client: StockClient
  page: Page

  constructor(inv: Client, client: StockClient, page: Page) {
    this.inventory = inv
    this.client = client
    this.page = page
  }

  show() {
    let options = this.inventory.getState()
    let dtos: inventoryItemDTO[] = options.map(o => ({
      id: String(o.id),
      name: o.name
    }))
    this.page.populateDropdown(dtos)

    let data = this.computeTableRows()
    this.page.populateTable(data)
  }

  init() {
    this.client.fetchState().then(() => {
      let data = this.computeTableRows()
      this.page.populateTable(data)
    })
  }

  computeTableRows(): stockTableRowDTO[] {
    let positions: Position[] = this.client.getData()
    let dict = this.toDict(positions)

    let toDTO = (i: inventoryItem) => ({
      id: i.id,
      name: i.name,
      qty: String(dict[i.id] === undefined ? 0 : dict[i.id])
    })

    return this.inventory.getState().map(toDTO)
  }

  toDict(i: Position[]): StockDict {
    let r: StockDict = {}
    i.forEach(e => {
      r = {
        [e.id]: String(e.qty),
        ...r
      }
    })
    return r
  }

  onQtyChange() {
    let qty = this.page.getQty()

    if (this.badQuantity(qty)) {
      this.page.toggleError(true)
      return
    }

    this.page.toggleError(false)
  }

  onProvision() {
    let qty = this.page.getQty()

    if (this.badQuantity(qty)) {
      this.page.toggleError(true)
      return
    }

    this.page.toggleError(false)

    let id = this.page.getID()
    this.client.provision(id, qty).then(() => {
      this.page.resetQty()
      let data = this.computeTableRows()
      this.page.populateTable(data)
    })
  }

  private badQuantity(qty: any): boolean {
    if (!qty || Number(qty) <= 0) {
      return true
    }
    return false
  }
}
