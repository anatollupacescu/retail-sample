import axios from 'axios'

interface item {
  id: number
  qty: number
}

export interface Position {
  id: number
  qty: number
  items: item[]
}

export default class StockClient {
  private endpoint: string
  private data: Position[]

  constructor(url: string, initialData: Position[] = []) {
    this.endpoint = `${url}/stock`
    this.data = initialData
  }

  async provision(id: string, qty: number): Promise<any> {
    let data = { [id]: Number(qty) }
    let resp = await this.apiProvision(data)
    Object.keys(resp.data.data).forEach(k => {
      let updated = this.updatePosition(Number(k), Number(resp.data.data[k]))
      if (!updated) {
        throw 'Not updated'
      }
    })
    return resp
  }

  substractFromPosition(ingredientID: number, toSubstract: number): void {
    let ingredient = this.data.filter(p => p.id === ingredientID)
    if (ingredient && ingredient.length > 0) {
      ingredient[0].qty = ingredient[0].qty - toSubstract
    }
  }

  updatePosition(k: number, v: number): boolean {
    let updated = false
    this.data.map(p => {
      if (updated) return
      if (p.id === k) {
        p.qty = v
        updated = true
      }
    })
    return updated
  }

  private apiProvision(data: any): Promise<any> {
    return axios.post(this.endpoint, data)
  }

  async fetchState(): Promise<void> {
    let data = await this.apiFetchState()
    this.data = data.data.data
  }

  private async apiFetchState(): Promise<any> {
    return axios.get(this.endpoint)
  }

  getData(): Position[] {
    return this.data
  }
}
