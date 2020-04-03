import Client, { inventoryItem } from './client'
import StockClient from '../stock/client'

export interface Page {
  toggleNameError(v: boolean): void
  name(): string
  resetName(): void
  toggleUniqueError(v: boolean): void
  addBtnEnabled(v: boolean): void
  renderTable(data: inventoryItem[]): void
}

export default class App {
  private client: Client
  private stock: StockClient
  private page: Page

  constructor(inv: Client, stock: StockClient, page: Page) {
    this.client = inv
    this.stock = stock
    this.page = page
  }

  init(): Promise<any> {
    return this.client.fetchState()
  }

  show() {
    let data = this.client.getState()
    this.page.renderTable(data)
  }

  onNameChange() {
    let n = this.page.name()

    if (n && n.trim().length > 0) {
      this.page.addBtnEnabled(true)
    } else {
      this.page.addBtnEnabled(false)
    }

    this.page.toggleNameError(false)
    this.page.toggleUniqueError(false)
  }

  onSubmit() {
    let name: string = this.page.name()

    this.client.addItem(name).then(res => {
      let addedItem = res[0]
      let error = res[1]

      switch (error) {
        case 'name empty':
          this.page.toggleNameError(true)
          this.page.addBtnEnabled(false)
          return

        case 'name present':
          this.page.toggleUniqueError(true)
          this.page.addBtnEnabled(false)
          return
      }

      this.stock.addPosition(addedItem.id)

      let data = this.client.getState()
      this.page.renderTable(data)
      this.page.resetName()
      this.page.addBtnEnabled(false)
    })
  }
}
