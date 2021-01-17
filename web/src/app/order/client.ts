import { AxiosInstance } from 'axios'

export interface Record {
  id: number
  recipeID: number
  qty: number
}

export default class Client {
  private httpClient: AxiosInstance
  private state: Record[]

  constructor(httpClient: AxiosInstance, initial: Record[] = []) {
    this.httpClient = httpClient
    this.state = initial
  }

  private async apiFetchOrders(): Promise<any> {
    try {
      let data = await this.httpClient.get('/order')
      return data.data.data
    } catch (error) {
      throw error.response.data.trim()
    }
  }

  async fetchState(): Promise<any> {
    this.state = await this.apiFetchOrders()
  }

  private async apiAddOrder(recipeID: number, qty: number): Promise<any> {
    let payload = { id: Number(recipeID), qty: Number(qty) }
    try {
      let res = await this.httpClient.post('/order', payload)
      return res.data.data.id
    } catch (error) {
      throw error.response.data.trim()
    }
  }

  async addOrder(recipeID: number, qty: number): Promise<string> {
    if (!qty || qty === 0) {
      throw 'quantity mandatory'
    }

    try {
      let newID = await this.apiAddOrder(recipeID, qty)

      this.state.push({
        id: newID,
        recipeID: recipeID,
        qty: qty
      })

      return newID
    } catch (error) {
      if (error.endsWith('not enough stock: bad request')) {
        throw 'not enough stock'
      }
      throw `unknown error: ${error}`
    }
  }

  getState(): Record[] {
    return [...this.state]
  }
}
