import axios, { AxiosResponse } from 'axios'

export interface Position {
  id: number
  qty: number
}

type ApiResponse = Record<number, number>

export default class StockClient {
  private endpoint: string
  private data: Position[]

  constructor(url: string = '', initialData: Position[] = []) {
    this.endpoint = `${url}/stock`
    this.data = [...initialData]
  }

  addPosition(id: number): void {
    this.data.push({
      id: id,
      qty: 0
    })
  }

  private async apiProvision(data: any): Promise<ApiResponse> {
    try {
      let res: AxiosResponse = await axios.post(this.endpoint, data)
      return res.data.data
    } catch (error) {
      return Promise.reject(`got error provisioning: ${error}`)
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
    this.data.map(p => {
      if (updated) return
      if (p.id === id) {
        p.qty = newValue
        updated = true
      }
    })
  }

  substractFromPosition(ingredientID: number, toSubstract: number): void {
    let ingredient = this.data.filter(p => p.id === ingredientID)
    if (ingredient && ingredient.length > 0) {
      ingredient[0].qty = ingredient[0].qty - toSubstract
    }
  }

  private async apiFetchState(): Promise<Position[]> {
    try {
      let res = await axios.get(this.endpoint)
      return res.data.data
    } catch (error) {
      return Promise.reject(error)
    }
  }

  async fetchState(): Promise<void> {
    this.data = await this.apiFetchState()
  }

  getData(): Position[] {
    return [...this.data]
  }
}
