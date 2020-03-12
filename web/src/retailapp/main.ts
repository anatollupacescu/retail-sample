interface inventoryItem {
  name: string
  id: number
}

interface adapters {
  fetchInventoryState: Function
}

export default class RetailUI {
  inventory: inventoryItem[]
  apiURL: string
  adapters: adapters

  constructor (apiURL: string, adapters: adapters) {
    this.inventory = []
    this.apiURL = apiURL
    this.adapters = adapters
  }

  fetchInventoryState () {
    this.adapters.fetchInventoryState()
  }

  addInventoryItem (itemName: string) {
    console.log(itemName)
  }
}
