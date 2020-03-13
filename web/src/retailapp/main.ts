interface inventoryItem {
  name: string
  id: number
}

export default class RetailUI {
  private inventory: inventoryItem[] = []

  getInventory (): inventoryItem[] {
    return this.inventory
  }

  apiFetchInventoryState () {
    return []
  }

  apiAddInventoryItem (name: string): inventoryItem {
    console.log(name)
    return {
      name: '',
      id: 1
    }
  }

  fetchInventoryState () {
    this.inventory = this.apiFetchInventoryState()
  }

  addInventoryItem (itemName: string) {
    if (!itemName || itemName.length === 0) {
      return 'name empty'
    }
    if (!this.nameIsUnique(itemName)) {
      return 'name present'
    }
    let addedItem = this.apiAddInventoryItem(itemName)

    this.inventory.push(addedItem)
    return ''
  }

  private nameIsUnique (name: string) {
    let found = this.inventory.find(item => item.name === name)
    return found === undefined
  }
}
