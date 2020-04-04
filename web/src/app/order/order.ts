import OrderClient, { Record } from './client'
import StockClient from '../stock/client'
import RecipeClient, { Recipe } from '../recipe/client'

export interface tableRowDTO {
  id: string
  recipe: string
  qty: string
}

export interface Page {
  toggleSubmitButtonState(v: boolean): void
  getRecipeID(): number
  getQty(): number
  resetQty(): void
  toggleQtyError(v: boolean): void
  toggleNotEnoughStockError(v: boolean): void
  populateDropdown(rows: Recipe[]): void
  populateTable(rows: tableRowDTO[]): void
}

export default class Order {
  private page: Page

  private stock: StockClient
  private recipe: RecipeClient
  private client: OrderClient

  constructor(stockClient: StockClient, orderClient: OrderClient, recipeClient: RecipeClient, page: Page) {
    this.stock = stockClient
    this.client = orderClient
    this.recipe = recipeClient
    this.page = page
  }

  qtyInputChanged() {
    let qty = this.page.getQty()

    if (this.badQuantity(qty)) {
      this.page.toggleQtyError(true)
      this.page.toggleSubmitButtonState(false)
      return
    }

    this.page.toggleQtyError(false)
    this.page.toggleNotEnoughStockError(false)

    this.page.toggleSubmitButtonState(true)
  }

  recipeInputChanged() {
    this.page.toggleNotEnoughStockError(false)
  }

  show() {
    let items = this.recipe.getRecipes()
    this.page.populateDropdown(items)
  }

  init() {
    this.client
      .fetchOrders()
      .then(() => {
        let data: Record[] = this.client.getOrders()
        let rows = this.toRows(data)
        this.page.populateTable(rows)
      })
      .catch(err => console.error(err))
  }

  badQuantity(qty: any): boolean {
    if (!qty || Number(qty) <= 0) {
      return true
    }
    return false
  }

  async placeOrder(): Promise<any> {
    let recipeID: number = this.page.getRecipeID(),
      qty: number = this.page.getQty()

    if (this.badQuantity(qty)) {
      this.page.toggleQtyError(true)
      this.page.toggleSubmitButtonState(false)
      return
    }

    try {
      let result = await this.client.addOrder(recipeID, qty)
      console.log('got new order with id', result)
    } catch (error) {
      switch (error) {
        case 'not enough stock':
          this.page.toggleNotEnoughStockError(true)
          this.page.toggleSubmitButtonState(false)
          return
      }
    }

    let data: Record[] = this.client.getOrders()

    let rows = this.toRows(data)
    this.page.populateTable(rows)
    this.page.resetQty()
    this.page.toggleSubmitButtonState(false)
    this.page.toggleNotEnoughStockError(false)

    this.updateStock(recipeID, qty)
  }

  toRows(data: Record[]): tableRowDTO[] {
    return data.map(record => ({
      id: String(record.id),
      recipe: this.recipe.getByID(record.recipeID).name,
      qty: String(record.qty)
    }))
  }

  private updateStock(recipeID: number, orderSize: number): void {
    let recipeItems = this.recipe.getByID(recipeID).items

    recipeItems.forEach(item => {
      this.stock.substractFromPosition(item.id, item.qty * orderSize)
    })
  }
}
