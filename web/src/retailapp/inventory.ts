import InventoryClient, { inventoryItem } from '../client/inventory'
import StockClient from '../client/stock'

export default class Inventory {
  private client: InventoryClient
  private stock: StockClient

  constructor(inv: InventoryClient, stock: StockClient) {
    this.client = inv
    this.stock = stock
  }

  nameError(_v: boolean | undefined) {}
  getNameValue(): string {
    return ''
  }
  setNameEmpty(): void {}
  uniqueError(_v: boolean) {}
  addBtnEnabled(_v: boolean) {}
  isAddBtnEnabled(): boolean {
    return false
  }
  renderTable(_data: inventoryItem[]) {}

  async init() {
    let data = await this.client.fetchState()
    this.renderTable(data)
  }

  onNameChange() {
    let n: string = this.getNameValue()

    if (n && n.trim().length > 0) {
      this.addBtnEnabled(true)
    } else {
      this.addBtnEnabled(false)
    }

    this.nameError(false)
    this.uniqueError(false)
  }

  onSubmit() {
    let n: string = this.getNameValue()

    if (!n || n.trim().length === 0) {
      this.nameError(true)
      this.addBtnEnabled(false)
      return
    }

    this.client.addItem(n).then(res => {
      let addedItem = res[0]
      let error = res[1]

      switch (error) {
        case 'name empty':
          this.nameError(true)
          this.addBtnEnabled(false)
          return

        case 'name present':
          this.uniqueError(true)
          this.addBtnEnabled(false)
          return
      }

      this.stock.addPosition(addedItem.id)

      let data = this.client.getInventory()
      this.renderTable(data)
      this.setNameEmpty()
      this.addBtnEnabled(false)
    })
  }
}
