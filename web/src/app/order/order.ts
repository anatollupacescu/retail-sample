import OrderClient from './client'
import StockClient from '../stock/client'
import RecipeClient from '../recipe/client'

export default class Order {
  private stock: StockClient
  private order: OrderClient
  private recipe: RecipeClient

  constructor(stockClient: StockClient, orderClient: OrderClient, recipeClient: RecipeClient) {
    this.stock = stockClient
    this.order = orderClient
    this.recipe = recipeClient
  }

  async placeOrder(recipeID: number, qty: number): Promise<any> {
    await this.order.addOrder(recipeID, qty)
    this.updateStock(recipeID, qty)
  }

  private updateStock(recipeID: number, orderSize: number): void {
    let recipeItems = this.recipe.getByID(recipeID).items

    recipeItems.forEach(item => {
      this.stock.substractFromPosition(item.id, item.qty * orderSize)
    })
  }
}
