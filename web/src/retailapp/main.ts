interface inventoryItem {
  name: string
  id: number
}

export class apiAdapter {
  fetchInventoryState!: Function
  addInventoryItem!: Function
}

export default class RetailUI {
  private inventory: inventoryItem[]
  apiAdapter: apiAdapter

  constructor (apiAdapter: apiAdapter) {
    this.inventory = []
    this.apiAdapter = apiAdapter
  }

  fetchInventoryState () {
    return this.apiAdapter.fetchInventoryState()
  }

  addInventoryItem (itemName: string) {
    if (!itemName || itemName.length === 0) {
      return 'name empty'
    }
    if (!this.nameIsUnique(itemName)) {
      return 'name present'
    }
    let addedItem = this.apiAdapter.addInventoryItem(itemName)

    this.inventory.push(addedItem)
    return ''
  }

  private nameIsUnique (name: string) {
    let found = this.inventory.find(item => item.name === name)
    return found === undefined
  }
}
