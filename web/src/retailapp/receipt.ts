interface ReceiptItem {
  id: number
  qty: number
}

interface ReceiptCandidate {
  name: string
  items: ReceiptItem[]
}

interface Receipt {
  id: number
  name: string
}

export default class ReceiptClient {
  private pendingReceipt: ReceiptCandidate = this.emptyReceipt()
  private endpoint: string
  private receipts: Receipt[]

  constructor(url: string, receipts: Receipt[] = []) {
    this.endpoint = `${url}/receipt`
    this.receipts = receipts
  }

  private emptyReceipt(): ReceiptCandidate {
    return {
      name: '',
      items: []
    }
  }

  setName(newReceiptName: string) {
    this.pendingReceipt.name = newReceiptName
  }

  addIngredient(id: number, qty: number): string {
    if (qty === 0) {
      return 'zero quantity'
    }

    let item: ReceiptItem = {
      id: id,
      qty: qty
    }

    this.pendingReceipt.items.push(item)
    return ''
  }

  apiSaveReceipt(): Promise<any> {
    console.log(this.endpoint)
    return Promise.resolve()
  }

  async saveReceipt(): Promise<any> {
    let rName = this.pendingReceipt.name
    if (!rName || rName.length === 0) {
      return Promise.resolve('name empty')
    }

    let found = this.receipts.find(r => r.name === rName)

    if (found) {
      return Promise.resolve('name present')
    }

    if (this.pendingReceipt.items.length === 0) {
      return Promise.resolve('no ingredients')
    }

    let data = await this.apiSaveReceipt()
    this.receipts.push(data.data.data)

    this.pendingReceipt = this.emptyReceipt()
  }

  apiFetchReceipts(): Promise<any> {
    return Promise.resolve([])
  }

  async fetchReceipts(): Promise<any> {
    let data = await this.apiFetchReceipts()
    this.receipts = data.data.data
  }

  getReceipts(): Receipt[] {
    return this.receipts
  }

  listItems(): ReceiptItem[] {
    return this.pendingReceipt.items
  }
}
