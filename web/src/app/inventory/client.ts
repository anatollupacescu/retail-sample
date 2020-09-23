import { AxiosInstance } from 'axios'

export interface inventoryItem {
  id: number
  name: string
  enabled: boolean
}

export default class Client {
  private state: inventoryItem[]
  private httpClient: AxiosInstance

  constructor(httpClient: AxiosInstance, initial: inventoryItem[] = []) {
    this.httpClient = httpClient
    this.state = [...initial]
  }

  private async apiFetchState(): Promise<any> {
    try {
      let response = await this.httpClient.get("/inventory")
      return response.data.data
    } catch (error) {
      throw error.response.data.trim()
    }
  }

  async fetchState(): Promise<inventoryItem[]> {
    let apiData = await this.apiFetchState()
    this.state = [...apiData]
    return apiData
  }

  private async apiToggleItemStatus(id: string, enabled: boolean): Promise<inventoryItem> {
    try {
      let payload = {
        enabled: enabled
      }
      let response = await this.httpClient.patch(`/inventory/${id}`, payload)
      return response.data.data
    } catch (error) {
      throw error.response.data.trim()
    }
  }

  async toggleItemStatus(id: string, enabled: boolean): Promise<inventoryItem> {
    try {
      let item = await this.apiToggleItemStatus(id, enabled)
      this.saveToState(item)
      return item
    } catch (error) {
      throw `got error changing item status: ${error}`
    }
  }

  private saveToState(item: inventoryItem) {
    this.state = this.state.filter(i => i.id !== item.id)
    this.state.push(item)
  }

  private async apiAddItem(name: string): Promise<any> {
    try {
      let payload = {
        name: name
      }
      let res = await this.httpClient.post("/inventory", payload)
      return res.data.data
    } catch (error) {
      throw error.response.data.trim()
    }
  }

  async addItem(itemName: string): Promise<inventoryItem> {
    const errMsgEmptyName = 'name empty',
      errMsgNamePresent = 'name present'

    if (!itemName || itemName.length === 0) {
      throw errMsgEmptyName
    }

    if (!this.isUnique(itemName)) {
      throw errMsgNamePresent
    }

    try {
      let newItem = await this.apiAddItem(itemName)
      this.state.push(newItem)
      return newItem
    } catch (error) {
      switch (error) {
        case 'name not provided':
          throw errMsgEmptyName
        case 'item type already present':
          throw errMsgNamePresent
        default:
          throw 'unexpected response from the server'
      }
    }
  }

  private isUnique(name: string) {
    let found = this.state.find(item => item.name === name)
    return found === undefined
  }

  getState(): inventoryItem[] {
    return [...this.state]
  }

  getEnabledItems(): inventoryItem[] {
    return this.getState().filter(i => i.enabled)
  }

  findByID(id: string): inventoryItem {
    let toMatch = Number.parseInt(id)
    let found = this.state.find(i => i.id === toMatch)
    if (!found) {
      throw 'not found'
    }
    return found
  }
}
