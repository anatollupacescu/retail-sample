import axios from 'axios'

export interface inventoryItem {
  id: number
  name: string
}

export default class Client {
  private endpoint: string
  private inventory: inventoryItem[]

  constructor(url: string = '', initial: inventoryItem[] = []) {
    this.endpoint = `${url}/inventory`
    this.inventory = [...initial]
  }

  private async apiFetchState(): Promise<any> {
    try {
      let response = await axios.get(this.endpoint)
      return response.data.data
    } catch (error) {
      throw error.response.data
    }
  }

  async fetchState(): Promise<inventoryItem[]> {
    const data = await this.apiFetchState()
    this.inventory = data.data.data
    return this.inventory
  }

  async apiAddItem(name: string): Promise<any> {
    try {
      let payload = {
        name: name
      }
      let res = await axios.post(this.endpoint, payload)
      return res.data.data
    } catch (error) {
      throw error.response.data
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
      this.inventory.push(newItem)
      return newItem
    } catch (error) {
      switch (error) {
        case 'ERR_EMPTY':
          throw errMsgEmptyName
        case 'ERR_UNIQUE':
          throw errMsgNamePresent
        default:
          throw new Error('unexpected response from the server')
      }
    }
  }

  private isUnique(name: string) {
    let found = this.inventory.find(item => item.name === name)
    return found === undefined
  }

  getState(): inventoryItem[] {
    return [...this.inventory]
  }

  getName(id: number): string {
    let item = this.inventory.find(i => i.id === id)
    if (item) {
      return item.name
    }
    throw `inventory item with ${id} not found`
  }
}
