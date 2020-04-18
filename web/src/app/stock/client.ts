import axios, { AxiosResponse } from 'axios'

export interface Position {
  id: number
  qty: number
}

export default class Client {
  private endpoint: string
  private state: Position[]

  constructor(url: string = '', initial: Position[] = []) {
    this.endpoint = `${url}/stock`
    this.state = [...initial]
  }

  addPosition(id: number): void {
    this.state.push({
      id: id,
      qty: 0
    })
  }

  private async apiProvision(data: any): Promise<any> {
    try {
      let res: AxiosResponse = await axios.post(this.endpoint, data)
      return res.data.data
    } catch (error) {
      throw error.response.data.trim()
    }
  }

  async provision(id: string, qty: number): Promise<any> {
    let data = { [id]: Number(qty) }
    let stock = await this.apiProvision(data)
    for (let id in stock) {
      this.updatePosition(Number(id), Number(stock[id]))
    }
  }

  updatePosition(id: number, newValue: number): void {
    let updated = false
    this.state.map(p => {
      if (updated) return
      if (p.id === id) {
        p.qty = newValue
        updated = true
      }
    })
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

  getData(): Position[] {
    return [...this.state]
  }
}
