interface RecipeItem {
  id: number
  qty: number
}

interface RecipeCandidate {
  name: string
  items: RecipeItem[]
}

interface Recipe {
  id: number
  name: string
}

export default class RecipeClient {
  private pendingRecipe: RecipeCandidate = this.emptyRecipe()
  private endpoint: string
  private recipes: Recipe[]

  constructor(url: string, receipts: Recipe[] = []) {
    this.endpoint = `${url}/receipt`
    this.recipes = receipts
  }

  private emptyRecipe(): RecipeCandidate {
    return {
      name: '',
      items: []
    }
  }

  setName(newReceiptName: string) {
    this.pendingRecipe.name = newReceiptName
  }

  addIngredient(id: number, qty: number): string {
    if (qty === 0) {
      return 'zero quantity'
    }

    let item: RecipeItem = {
      id: id,
      qty: qty
    }

    this.pendingRecipe.items.push(item)
    return ''
  }

  apiSaveRecipe(): Promise<any> {
    console.log(this.endpoint)
    return Promise.resolve()
  }

  async saveRecipe(): Promise<any> {
    let rName = this.pendingRecipe.name
    if (!rName || rName.length === 0) {
      return Promise.resolve('name empty')
    }

    let found = this.recipes.find(r => r.name === rName)

    if (found) {
      return Promise.resolve('name present')
    }

    if (this.pendingRecipe.items.length === 0) {
      return Promise.resolve('no ingredients')
    }

    let data = await this.apiSaveRecipe()
    this.recipes.push(data.data.data)

    this.pendingRecipe = this.emptyRecipe()
  }

  apiFetchReceipts(): Promise<any> {
    return Promise.resolve([])
  }

  async fetchRecipes(): Promise<any> {
    await this.apiFetchReceipts()
    this.recipes = []
  }

  getRecipes(): Recipe[] {
    return this.recipes
  }

  listItems(): RecipeItem[] {
    return this.pendingRecipe.items
  }
}
