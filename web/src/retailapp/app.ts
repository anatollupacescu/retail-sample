import OrderClient from '../client/order'
import StockClient from '../client/stock'
import RecipeClient from '../client/recipe'
import InventoryClient, { inventoryItem } from '../client/inventory'

export default class RetailApp {
  private stock: StockClient
  private order: OrderClient
  private recipe: RecipeClient
  private inventory: InventoryClient

  constructor(stockClient: StockClient, orderClient: OrderClient, recipeClient: RecipeClient, inv: InventoryClient) {
    this.stock = stockClient
    this.order = orderClient
    this.recipe = recipeClient
    this.inventory = inv
  }

  async addInventoryItem(name: string): Promise<any> {
    let newItems = await this.inventory.addItem(name)

    newItems.forEach((i: inventoryItem) => {
      this.stock.addPosition(i.id)
    })
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
