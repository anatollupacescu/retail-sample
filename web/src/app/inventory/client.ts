import axios from 'axios'

export interface inventoryItem {
  id: number
  name: string
}

export default class Client {
  private endpoint: string
  private state: inventoryItem[]

  constructor(url: string = '', initial: inventoryItem[] = []) {
    this.endpoint = url
    this.state = [...initial]
  }

  private async apiFetchState(): Promise<any> {
    try {
      let response = await axios.get(this.endpoint)
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

  async apiAddItem(name: string): Promise<any> {
    try {
      let payload = {
        name: name
      }
      let res = await axios.post(this.endpoint, payload)
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
        case "name not provided":
          throw errMsgEmptyName
        case "item type already present":
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

  getName(id: number): string {
    let item = this.state.find(i => i.id === id)
    if (item) {
      return item.name
    }
    throw `inventory item with ${id} not found`
  }
}
