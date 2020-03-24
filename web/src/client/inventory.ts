import axios from 'axios'

export interface inventoryItem {
  id: number
  name: string
}

export default class InventoryClient {
  private endpoint: string
  private inventory: inventoryItem[]

  constructor(url: string = '', initial: inventoryItem[] = []) {
    this.endpoint = `${url}/inventory`
    this.inventory = initial.slice(0)
  }

  apiFetchState(): Promise<any> {
    return axios.get(this.endpoint)
  }

  async fetchState(): Promise<inventoryItem[]> {
    const data = await this.apiFetchState()
    this.inventory = data.data.data
    return data
  }

  async apiAddItem(name: string): Promise<any> {
    try {
      let res = await axios.post(this.endpoint, {
        name: name
      })
      return [res.data.data, '']
    } catch (error) {
      return [[], error.response.data]
    }
  }

  async addItem(itemName: string /*TODO shoud accept an array*/): Promise<[inventoryItem[], string]> {
    const errMsgEmptyName = 'name empty',
      errMsgNamePresent = 'name present'

    if (!itemName || itemName.length === 0) {
      return [[], errMsgEmptyName]
    }

    if (!this.isUnique(itemName)) {
      return [[], errMsgNamePresent]
    }

    let apiResponse: [inventoryItem[], string] = await this.apiAddItem(itemName)

    let newItems = apiResponse[0]

    this.inventory.push(...newItems)

    switch (apiResponse[1]) {
      case 'ERR_EMPTY':
        return [newItems, errMsgEmptyName]
      case 'ERR_UNIQUE':
        return [newItems, errMsgNamePresent]
      case '':
        return apiResponse
      default:
        throw new Error('unexpected response from the server')
    }
  }

  private isUnique(name: string) {
    let found = this.inventory.find(item => item.name === name)
    return found === undefined
  }

  getInventory(): inventoryItem[] {
    return this.inventory
  }

  getName(id: number): string {
    let item = this.inventory.find(i => i.id === id)
    if (item) {
      return item.name
    }
    return ''
  }
}
