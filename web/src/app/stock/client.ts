import { AxiosResponse, AxiosInstance } from 'axios'

export interface Position {
  id: number
  qty: number
}

export type StockDict = Record<number, number>

export default class Client {
  private httpClient: AxiosInstance
  private state: StockDict

  constructor(httpClient: AxiosInstance, initial: Position[] = []) {
    this.httpClient = httpClient
    this.state = {}
    initial.forEach((v) => {
      this.state[v.id] = v.qty
    })
  }

  addPosition(id: number): void {
    this.state[id] = 0
  }

  private async apiProvision(id: string, qty: number): Promise<any> {
    try {
      let data = {
        qty: qty
      }
      let res: AxiosResponse = await this.httpClient.post(`/stock/${id}`, data)
      return res.data.data
    } catch (error) {
      throw error.response.data.trim()
    }
  }

  async provision(id: string, qty: number): Promise<void> {
    let stock = await this.apiProvision(id, qty)
    this.state[stock.id] = stock.qty
  }

  substractFromPosition(ingredientID: number, toSubstract: number): void {
    this.state[ingredientID] = this.state[ingredientID] - toSubstract
  }

  private async apiFetchState(): Promise<Position[]> {
    try {
      let res = await this.httpClient.get('/stock')
      return res.data.data
    } catch (error) {
      throw error.response.data.trim()
    }
  }

  async fetchState(): Promise<void> {
    let statePositions = await this.apiFetchState()
    this.state = {}
    statePositions.forEach((position) => {
      this.state[position.id] = position.qty
    })
  }

  getState(): StockDict {
    return this.state
  }
}
