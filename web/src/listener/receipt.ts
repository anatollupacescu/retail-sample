import InventoryClient from '../retailapp/inventory'
import ReceiptClient from '../retailapp/receipt'

export function initializeReceipt(inv: InventoryClient, receipt: ReceiptClient) {
  inv.fetchInventoryState().then(() => {
    populateDropdown(inv.getInventory())
  })

  receipt.fetchReceipts().then(() => {
    populateTable(receipt.listReceipts())
  })

  //link buttons etc...
}

function populateDropdown(items: any[]) {
  console.log(items)
}

function populateTable(items: any[]) {
  console.log(items)
}
