import InventoryClient, { inventoryItem } from './client'
import StockClient from '../stock/client'

export interface InventoryPage {
  nameError(_v: boolean | undefined): void
  getNameValue(): string
  setNameEmpty(): void
  uniqueError(_v: boolean): void
  addBtnEnabled(_v: boolean): void
  isAddBtnEnabled(): boolean
  renderTable(_data: inventoryItem[]): void
}

export default class Inventory {
  private client: InventoryClient
  private stock: StockClient
  private page: InventoryPage

  constructor(inv: InventoryClient, stock: StockClient, page: InventoryPage) {
    this.client = inv
    this.stock = stock
    this.page = page
  }

  async init() {
    let data = await this.client.fetchState()
    this.page.renderTable(data)
  }

  onNameChange() {
    let n: string = this.page.getNameValue()

    if (n && n.trim().length > 0) {
      this.page.addBtnEnabled(true)
    } else {
      this.page.addBtnEnabled(false)
    }

    this.page.nameError(false)
    this.page.uniqueError(false)
  }

  onSubmit() {
    let n: string = this.page.getNameValue()

    if (!n || n.trim().length === 0) {
      this.page.nameError(true)
      this.page.addBtnEnabled(false)
      return
    }

    this.client.addItem(n).then(res => {
      let addedItem = res[0]
      let error = res[1]

      switch (error) {
        case 'name empty':
          this.page.nameError(true)
          this.page.addBtnEnabled(false)
          return

        case 'name present':
          this.page.uniqueError(true)
          this.page.addBtnEnabled(false)
          return
      }

      this.stock.addPosition(addedItem.id)

      let data = this.client.getInventory()
      this.page.renderTable(data)
      this.page.setNameEmpty()
      this.page.addBtnEnabled(false)
    })
  }
}
