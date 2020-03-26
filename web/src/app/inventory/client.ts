import axios from 'axios'

export interface inventoryItem {
  id: number
  name: string
}

type AddItemResult = [inventoryItem, string]

const zeroValueItem: inventoryItem = {
  id: 0,
  name: ''
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
    return this.inventory
  }

  async apiAddItem(name: string): Promise<AddItemResult> {
    try {
      let res = await axios.post(this.endpoint, {
        name: name
      })
      let createdItem: inventoryItem = res.data.data
      return [createdItem, '']
    } catch (error) {
      return [zeroValueItem, error.response.data]
    }
  }

  async addItem(itemName: string): Promise<AddItemResult> {
    const errMsgEmptyName = 'name empty',
      errMsgNamePresent = 'name present'

    if (!itemName || itemName.length === 0) {
      return [zeroValueItem, errMsgEmptyName]
    }

    if (!this.isUnique(itemName)) {
      return [zeroValueItem, errMsgNamePresent]
    }

    let apiResponse = await this.apiAddItem(itemName)

    switch (apiResponse[1]) {
      case 'ERR_EMPTY':
        return [zeroValueItem, errMsgEmptyName]
      case 'ERR_UNIQUE':
        return [zeroValueItem, errMsgNamePresent]
      case '':
        let newItem = apiResponse[0]
        this.inventory.push(newItem)
        return [zeroValueItem, '']
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
