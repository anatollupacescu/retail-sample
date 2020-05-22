import axios from 'axios'

export interface Record {
  id: number
  recipeID: number
  qty: number
}

export default class Client {
  private endpoint: string
  private state: Record[]

  constructor(url: string = '', initial: Record[] = []) {
    this.endpoint = url
    this.state = initial
  }

  private async apiFetchOrders(): Promise<any> {
    try {
      let data = await axios.get(this.endpoint)
      return data.data.data
    } catch (error) {
      throw error.response.data.trim()
    }
  }

  async fetchState(): Promise<any> {
    this.state = await this.apiFetchOrders()
  }

  private async apiAddOrder(recipeID: number, qty: number): Promise<any> {
    try {
      let res = await axios.post(this.endpoint, { id: Number(recipeID), qty: Number(qty) })
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
      switch (error) {
        case 'not enough stock':
          throw 'not enough stock'
        default:
          throw `unknown error: ${error}`
      }
    }
  }

  getState(): Record[] {
    return [...this.state]
  }
}
