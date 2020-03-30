import axios from 'axios'

export interface Record {
  id: number
  recipeID: number
  qty: number
}

type AddOrderResult = [number, string]

export default class OrderClient {
  private endpoint: string
  private orders: Record[]

  constructor(url: string = '', initial: Record[] = []) {
    this.endpoint = `${url}/order`
    this.orders = initial
  }

  private async apiFetchOrders(): Promise<any> {
    try {
      let data = await axios.get(this.endpoint)
      return data.data.data
    } catch (error) {
      return Promise.reject(error.response.data.trim())
    }
  }

  async fetchOrders(): Promise<any> {
    this.orders = await this.apiFetchOrders()
  }

  private async apiAddOrder(recipeID: number, qty: number): Promise<AddOrderResult> {
    try {
      let res = await axios.post(this.endpoint, { id: Number(recipeID), qty: Number(qty) })
      return [res.data.data.id, '']
    } catch (error) {
      return [0, error.response.data.trim()]
    }
  }

  async addOrder(recipeID: number, qty: number): Promise<string> {
    if (!qty || qty === 0) {
      return 'quantity mandatory'
    }

    let result = await this.apiAddOrder(recipeID, qty)

    let err = result[1]

    switch (err) {
      case 'not enough stock':
        return 'not enough stock'
      case '':
        break
      default:
        throw `unknown error: ${err}`
    }

    this.orders.push({
      id: result[0],
      recipeID: recipeID,
      qty: qty
    })

    return ''
  }

  getOrders(): Record[] {
    return [...this.orders]
  }
}
