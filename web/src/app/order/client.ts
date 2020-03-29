import axios from 'axios'

export interface OrderDTO {
  date: string
  recipeID: number
  qty: number
}

export default class OrderClient {
  private endpoint: string
  private orders: OrderDTO[]

  constructor(url: string = '', initial: OrderDTO[] = []) {
    this.endpoint = `${url}/order`
    this.orders = initial
  }

  async fetchOrders(): Promise<any> {
    let data = await this.apiFetchOrders()
    this.orders = data.data.data
  }

  private apiFetchOrders(): Promise<any> {
    return axios.get(this.endpoint)
  }

  private async apiAddOrder(recipeID: number, qty: number): Promise<string> {
    try {
      await axios.post(this.endpoint, { id: Number(recipeID), qty: Number(qty) })
    } catch (error) {
      return error.response.data.trim()
    }

    return ''
  }

  async addOrder(recipeID: number, qty: number): Promise<string> {
    if (!qty || qty === 0) {
      return 'quantity mandatory'
    }

    let msg = await this.apiAddOrder(recipeID, qty)

    switch (msg) {
      case 'not enough stock':
        return 'not enough stock'
      case '':
        this.orders.push({
          date: '',
          recipeID: recipeID,
          qty: qty
        })
        break
      default:
        throw `unknown error: ${msg}`
    }

    return ''
  }

  getOrders(): OrderDTO[] {
    return this.orders
  }
}
