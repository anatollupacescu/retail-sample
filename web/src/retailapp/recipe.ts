import axios from 'axios'

export interface RecipeItem {
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
  private pendingRecipe: RecipeCandidate
  private endpoint: string
  private recipes: Recipe[]

  constructor(url: string, recipe: Recipe[] = []) {
    this.endpoint = `${url}/recipe`
    this.recipes = recipe
    this.pendingRecipe = this.emptyRecipe()
  }

  private emptyRecipe(): RecipeCandidate {
    return {
      name: '',
      items: []
    }
  }

  setName(newRecipeName: string) {
    this.pendingRecipe.name = newRecipeName
  }

  addIngredient(id: number, qty: number): string {
    if (!qty || qty === 0) {
      return 'zero quantity'
    }

    let found = this.pendingRecipe.items.find(i => i.id === id)

    if (found) {
      return 'duplicate id'
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

  apiFetchRecipes(): Promise<any> {
    return axios.get(this.endpoint)
  }

  async fetchRecipes(): Promise<any> {
    let data = await this.apiFetchRecipes()
    this.recipes = data.data.data
  }

  getRecipes(): Recipe[] {
    return this.recipes
  }

  listItems(): RecipeItem[] {
    return this.pendingRecipe.items
  }
}
