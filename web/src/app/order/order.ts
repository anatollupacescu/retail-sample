import OrderClient, { OrderDTO } from './client'
import StockClient from '../stock/client'
import RecipeClient, { Recipe } from '../recipe/client'

export interface OrderPage {
  toggleSubmitButtonState(v: boolean): void
  getRecipeID(): number
  getQty(): number
  toggleQtyError(v: boolean): void
  toggleNotEnoughStockError(v: boolean): void
  populateDropdown(rows: Recipe[]): void
  populateTable(rows: OrderDTO[]): void
}

export default class Order {
  private page: OrderPage

  private stock: StockClient
  private recipe: RecipeClient
  private client: OrderClient

  constructor(stockClient: StockClient, orderClient: OrderClient, recipeClient: RecipeClient, page: OrderPage) {
    this.stock = stockClient
    this.client = orderClient
    this.recipe = recipeClient
    this.page = page
  }

  qtyInputChanged() {
    let qty = this.page.getQty()

    if (!qty || Number(qty) <= 0) {
      this.page.toggleQtyError(true)
      this.page.toggleSubmitButtonState(false)
      return
    }

    this.page.toggleQtyError(false)
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
    console.log('will call api to fetch the state')
  }

  async placeOrder(): Promise<any> {
    let recipeID: number = this.page.getRecipeID(),
      qty: number = this.page.getQty()

    await this.client.addOrder(recipeID, qty)
    this.updateStock(recipeID, qty)
  }

  private updateStock(recipeID: number, orderSize: number): void {
    let recipeItems = this.recipe.getByID(recipeID).items

    recipeItems.forEach(item => {
      this.stock.substractFromPosition(item.id, item.qty * orderSize)
    })
  }
}
