import axios from 'axios'

export interface RecipeItem {
  id: number
  qty: number
}

interface RecipeCandidate {
  name: string
  items: RecipeItem[]
}

export interface Recipe {
  id: number
  name: string
  items: RecipeItem[]
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
    let items: any = {}
    this.pendingRecipe.items.forEach(i => {
      items[i.id] = i.qty
    })
    let payload = {
      name: this.pendingRecipe.name,
      items: items
    }
    return axios.post(this.endpoint, payload)
  }

  async saveRecipe(): Promise<string> {
    let rName = this.pendingRecipe.name
    if (!rName || rName.length === 0) {
      return 'name empty'
    }

    let found = this.recipes.find(r => r.name === rName)

    if (found) {
      return 'name present'
    }

    if (this.pendingRecipe.items.length === 0) {
      return 'no ingredients'
    }

    let data = await this.apiSaveRecipe()

    Object.keys(data.data.data).forEach((i: any) => {
      this.recipes.push({
        id: Number(data.data.data[i]),
        name: String(i),
        items: this.pendingRecipe.items
      })
    })

    this.pendingRecipe = this.emptyRecipe()

    return ''
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

  getByID(id: number): Recipe {
    let r = this.recipes.filter(r => r.id === id)
    if (!r || r.length === 0) {
      throw `recipe with id ${id} not found`
    }
    return r[0]
  }

  listItems(): RecipeItem[] {
    return this.pendingRecipe.items
  }
}
