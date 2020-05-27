import InventoryClient, { inventoryItem } from '../inventory/client'
import Client from './client'

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
  id(): string
  qty(): number
  resetQty(): void
  toggleError(v: boolean): void
  toggleAddBtnState(v: boolean): void
  populateTable(data: stockTableRowDTO[]): void
  populateDropdown(data: inventoryItemDTO[]): void
}

type StockDict = Record<string, string>

export default class App {
  private inventory: InventoryClient
  private client: Client
  private page: Page

  constructor(inv: InventoryClient, client: Client, page: Page) {
    this.inventory = inv
    this.client = client
    this.page = page
  }

  init() {
    this.client.fetchState()
  }

  show() {
    this.populateDropdown()
    let data = this.computeTableRows()
    this.page.populateTable(data)
  }

  private populateDropdown() {
    let options = this.inventory.getEnabledItems()
    let dtos: inventoryItemDTO[] = options.map(o => ({
      id: String(o.id),
      name: o.name
    }))
    this.page.populateDropdown(dtos)
  }

  private computeTableRows(): stockTableRowDTO[] {
    let positions: Position[] = this.client.getState()
    let dict = this.toDict(positions)

    let toDTO = (i: inventoryItem) => ({
      id: i.id,
      name: i.name,
      qty: String(dict[i.id] === undefined ? 0 : dict[i.id])
    })

    return this.inventory.getState().map(toDTO)
  }

  private toDict(i: Position[]): StockDict {
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
    let qty = this.page.qty()

    if (this.badQuantity(qty)) {
      this.page.toggleError(true)
      return
    }

    this.page.toggleAddBtnState(true)
    this.page.toggleError(false)
  }

  onProvision() {
    let qty = this.page.qty()

    if (this.badQuantity(qty)) {
      this.page.toggleError(true)
      return
    }

    this.page.toggleError(false)

    let id = this.page.id()
    this.client.provision(id, qty).then(() => {
      this.page.resetQty()
      this.page.toggleAddBtnState(false)

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
