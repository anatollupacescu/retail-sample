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

  constructor(url: string = '', initial: Recipe[] = []) {
    this.endpoint = `${url}/recipe`
    this.recipes = [...initial]
  }

  async apiSaveRecipe(name: string, ingredients: RecipeItem[]): Promise<any> {
    let items: any = {}
    ingredients.forEach(i => {
      items[i.id] = i.qty
    })
    let payload = {
      name: name,
      items: items
    }
    try {
      let res = await axios.post(this.endpoint, payload)
      return res.data.data
    } catch (error) {
      throw error.response.data.trim()
    }
  }

  async saveRecipe(name: string, ingredients: RecipeItem[]): Promise<string> {
    let found = this.recipes.find(r => r.name === name)

    if (found) {
      throw 'name present'
    }

    let data = await this.apiSaveRecipe(name, ingredients)

    Object.keys(data).forEach((name: any) => {
      let id = data[name]
      this.recipes.push({
        id: Number(id),
        name: String(name),
        items: ingredients
      })
    })

    return ''
  }

  async apiFetchRecipes(): Promise<any> {
    try {
      let res = await axios.get(this.endpoint)
      return res.data.data
    } catch (error) {
      throw error.response.data.trim()
    }
  }

  async fetchRecipes(): Promise<any> {
    let data = await this.apiFetchRecipes()
    this.recipes = data
    return this.recipes
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
