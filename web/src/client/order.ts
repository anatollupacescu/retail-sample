import axios from 'axios'

export interface Order {
  date: string
  recipeID: number
  qty: number
}

export default class OrderClient {
  private endpoint: string
  private orders: Order[]

  constructor(url: string = '', initial: Order[] = []) {
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

  private apiAddOrder(recipeID: number, qty: number): Promise<any> {
    return axios.post(this.endpoint, { id: recipeID, qty: qty })
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

  getOrders(): Order[] {
    return this.orders
  }
}
