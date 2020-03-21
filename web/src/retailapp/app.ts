import OrderClient from '../client/order'
import StockClient from '../client/stock'
import RecipeClient from '../client/recipe'

export default class RetailApp {
  private stockClient: StockClient
  private orderClient: OrderClient
  private recipeClient: RecipeClient

  constructor(stockClient: StockClient, orderClient: OrderClient, recipeClient: RecipeClient) {
    this.stockClient = stockClient
    this.orderClient = orderClient
    this.recipeClient = recipeClient
  }

  async placeOrder(recipeID: number, qty: number): Promise<any> {
    await this.orderClient.addOrder(recipeID, qty)
    this.updateStock(recipeID, qty)
  }

  private updateStock(recipeID: number, orderSize: number): void {
    let recipeItems = this.recipeClient.getByID(recipeID).items

    recipeItems.forEach(item => {
      this.stockClient.substractFromPosition(item.id, item.qty * orderSize)
    })
  }
}
