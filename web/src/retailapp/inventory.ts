import axios from 'axios'

interface inventoryItem {
  name: string
  id: number
}

export default class InventoryClient {
  private endpoint: string
  private inventory: inventoryItem[] = []

  constructor(url: string) {
    this.endpoint = `${url}/inventory`
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

  apiFetchInventoryState(): Promise<any> {
    return axios.get(this.endpoint)
  }

  apiAddInventoryItem(name: string): Promise<any> {
    return axios.post(this.endpoint, [name])
  }

  async fetchInventoryState() {
    const data = await this.apiFetchInventoryState()
    this.inventory = data.data.data
  }

  async addInventoryItem(itemName: string): Promise<any> {
    if (!itemName || itemName.length === 0) {
      return Promise.resolve('name empty')
    }
    if (!this.nameIsUnique(itemName)) {
      return Promise.resolve('name present')
    }
    const data = await this.apiAddInventoryItem(itemName)

    Object.keys(data.data.data).forEach((key: string) => {
      this.inventory.push({
        name: key,
        id: data.data.data[key]
      })
    })
  }

  private nameIsUnique(name: string) {
    let found = this.inventory.find(item => item.name === name)
    return found === undefined
  }
}
