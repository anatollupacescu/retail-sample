import Client, { inventoryItem } from './client'
import StockClient from '../stock/client'

export interface Page {
  clearRow(): void
  highlightRow(s: string): void
  toggleModal(v: boolean): void
  populateModal(i: inventoryItem): void
  toggleNameError(v: boolean): void
  name(): string
  resetName(): void
  toggleUniqueError(v: boolean): void
  addBtnEnabled(v: boolean): void
  renderTable(data: inventoryItem[]): void
}

export default class App {
  private client: Client
  private stock: StockClient
  private page: Page

  private selectedID: string = ''

  constructor(inv: Client, stock: StockClient, page: Page) {
    this.client = inv
    this.stock = stock
    this.page = page
  }

  init() {
    this.client.fetchState().then(data => {
      this.page.renderTable(data)
    })
  }

  toggleItemStatus(s: boolean) {
    let id = this.selectedID

    if (!id) {
      return
    }

    let item = this.client.findByID(id)

    if (item.enabled === s) {
      throw 'already in the expected state'
    }

    this.client
      .toggleItemStatus(id, s)
      .then(this.page.populateModal)
      .then(() => {
        let data = this.client.getState()
        this.page.renderTable(data)
      })
      .then(() => this.page.highlightRow(id))
  }

  showModal() {
    if (!this.selectedID) {
      throw 'no row selected'
    }

    let modalItem = this.client.findByID(this.selectedID)
    this.page.populateModal(modalItem)
    this.openModal()
  }

  closeModal() {
    this.page.toggleModal(false)
  }

  openModal() {
    this.page.toggleModal(true)
  }

  onRowClick(id: string) {
    if (id === this.selectedID) {
      this.selectedID = ''
      return this.page.clearRow()
    }
    this.selectedID = id
    this.page.highlightRow(id)
  }

  onNameChange() {
    let n = this.page.name()

    if (n && n.trim().length > 0) {
      this.page.addBtnEnabled(true)
    } else {
      this.page.addBtnEnabled(false)
    }

    this.page.toggleNameError(false)
    this.page.toggleUniqueError(false)
  }

  onSubmit() {
    let name: string = this.page.name()

    this.client
      .addItem(name)
      .then((addedItem: inventoryItem) => {
        this.stock.addPosition(addedItem.id)

        let data = this.client.getState()
        this.page.renderTable(data)
        this.page.resetName()
        this.page.addBtnEnabled(false)
      })
      .catch(error => {
        switch (error) {
          case 'name empty':
            this.page.toggleNameError(true)
            this.page.addBtnEnabled(false)
            return

          case 'name present':
            this.page.toggleUniqueError(true)
            this.page.addBtnEnabled(false)
            return
        }
      })
  }
}
