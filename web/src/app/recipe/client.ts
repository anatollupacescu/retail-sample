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

export default class Client {
  private endpoint: string
  private state: Recipe[]

  constructor(url: string = '', initial: Recipe[] = []) {
    this.endpoint = `${url}/recipe`
    this.state = [...initial]
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
    let found = this.state.find(r => r.name === name)

    if (found) {
      throw 'name present'
    }

    let data = await this.apiSaveRecipe(name, ingredients)

    Object.keys(data).forEach((name: any) => {
      let id = data[name]
      this.state.push({
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
    this.state = [...data]
    return data
  }

  getState(): Recipe[] {
    return [...this.state]
  }

  getByID(id: number): Recipe {
    let r = this.state.filter(r => r.id === Number(id))
    if (!r || r.length === 0) {
      throw `recipe with id ${id} not found`
    }
    return r[0]
  }
}
