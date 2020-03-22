import axios from 'axios'

export interface inventoryItem {
  name: string
  id: number
}

export default class InventoryClient {
  private endpoint: string
  private inventory: inventoryItem[] = []

  constructor(url: string) {
    this.endpoint = `${url}/inventory`
  }

  apiFetchState(): Promise<any> {
    return axios.get(this.endpoint)
  }

  apiAddItem(name: string): Promise<any> {
    return axios.post(this.endpoint, [name])
  }

  async fetchState(): Promise<any> {
    const data = await this.apiFetchState()
    this.inventory = data.data.data
    return data
  }

  async addItem(itemName: string /*TODO shoud accept an array*/): Promise<any> {
    if (!itemName || itemName.length === 0) {
      return Promise.resolve('name empty')
    }
    if (!this.nameIsUnique(itemName)) {
      return Promise.resolve('name present')
    }
    const data = await this.apiAddItem(itemName)

    let newItems: inventoryItem[] = []

    Object.keys(data.data.data).forEach((key: string) => {
      newItems.push({
        name: key,
        id: data.data.data[key]
      })
    })

    this.inventory.push(...newItems)

    return newItems
  }

  private nameIsUnique(name: string) {
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
