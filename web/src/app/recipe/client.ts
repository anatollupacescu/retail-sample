import axios from 'axios'

export interface RecipeItem {
  id: number
  qty: number
}

export interface Recipe {
  id: number
  name: string
  items: RecipeItem[]
}

export default class RecipeClient {
  private endpoint: string
  private recipes: Recipe[]

  constructor(url: string = '', recipe: Recipe[] = []) {
    this.endpoint = `${url}/recipe`
    this.recipes = recipe
  }

  apiSaveRecipe(name: string, ingredients: RecipeItem[]): Promise<any> {
    let items: any = {}
    ingredients.forEach(i => {
      items[i.id] = i.qty
    })
    let payload = {
      name: name,
      items: items
    }
    return axios.post(this.endpoint, payload)
  }

  async saveRecipe(name: string, ingredients: RecipeItem[]): Promise<string> {
    let found = this.recipes.find(r => r.name === name)

    if (found) {
      return 'name present'
    }

    let data = await this.apiSaveRecipe(name, ingredients)

    Object.keys(data.data.data).forEach((i: any) => {
      this.recipes.push({
        id: Number(data.data.data[i]),
        name: String(i),
        items: ingredients
      })
    })

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
    return [...this.recipes]
  }

  getByID(id: number): Recipe {
    let r = this.recipes.filter(r => r.id === Number(id))
    if (!r || r.length === 0) {
      throw `recipe with id ${id} not found`
    }
    return r[0]
  }
}
