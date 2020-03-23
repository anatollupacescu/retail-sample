import axios from 'axios'

export interface Order {
  date: string
  recipeID: number
  qty: number
}

export default class OrderClient {
  private endpoint: string
  private orders: Order[] = []

  constructor(url: string) {
    this.endpoint = `${url}/order`
  }

  async fetchOrders(): Promise<any> {
    let data = await this.apiFetchOrders()
    this.orders = data.data.data
  }

  private apiFetchOrders(): Promise<any> {
    return axios.get(this.endpoint)
  }

  async addOrder(recipeID: number, qty: number): Promise<string> {
    if (!qty || qty === 0) {
      return 'quantity mandatory'
    }
    let result = await this.apiAddOrder(recipeID, qty)
    if (result.status !== 201) {
      throw new Error('bad status')
    }
    this.orders.push({
      date: '',
      recipeID: recipeID,
      qty: qty
    })
    return ''
  }

  private apiAddOrder(recipeID: number, qty: number): Promise<any> {
    return axios.post(this.endpoint, { id: recipeID, qty: qty })
  }

  getOrders(): Order[] {
    return this.orders
  }
}
