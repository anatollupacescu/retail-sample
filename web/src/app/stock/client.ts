import axios, { AxiosResponse } from 'axios'

export interface Position {
  id: number
  qty: number
}

export default class Client {
  private endpoint: string
  private state: Position[]

  constructor(url: string = '', initial: Position[] = []) {
    this.endpoint = url
    this.state = [...initial]
  }

  addPosition(id: number): void {
    this.state.push({
      id: id,
      qty: 0
    })
  }

  private async apiProvision(id: string, qty: number): Promise<any> {
    try {
      let data = {
        qty: qty
      }
      let res: AxiosResponse = await axios.post(`${this.endpoint}/${id}`, data)
      return res.data.data
    } catch (error) {
      throw error.response.data.trim()
    }
  }

  async provision(id: string, qty: number): Promise<any> {
    let stock = await this.apiProvision(id, qty)
    for (let id in stock) {
      this.updatePosition(Number(id), Number(stock[id]))
    }
  }

  updatePosition(id: number, newValue: number): void {
    for (let p in this.state) {
      if (id === this.state[p].id) {
        this.state[p].qty = newValue
        return
      }
    }
  }

  substractFromPosition(ingredientID: number, toSubstract: number): void {
    let ingredient = this.state.filter(p => p.id === ingredientID)
    if (ingredient && ingredient.length > 0) {
      ingredient[0].qty = ingredient[0].qty - toSubstract
    }
  }

  private async apiFetchState(): Promise<Position[]> {
    try {
      let res = await axios.get(this.endpoint)
      return res.data.data
    } catch (error) {
      throw error.response.data.trim()
    }
  }

  async fetchState(): Promise<void> {
    this.state = await this.apiFetchState()
  }

  getState(): Position[] {
    return [...this.state]
  }
}
