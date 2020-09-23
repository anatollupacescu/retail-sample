import { AxiosInstance } from 'axios'

export interface RecipeItem {
  id: number
  qty: number
}

export interface Recipe {
  id: number
  name: string
  items: RecipeItem[]
  enabled: boolean
}

export default class Client {
  private httpClient: AxiosInstance
  private state: Recipe[]

  constructor(httpClient: AxiosInstance, initial: Recipe[] = []) {
    this.httpClient = httpClient
    this.state = [...initial]
  }

  private async apiSaveRecipe(name: string, ingredients: RecipeItem[]): Promise<any> {
    let items: any = {}
    ingredients.forEach((i) => {
      items[i.id] = i.qty
    })
    let payload = {
      name: name,
      items: items
    }
    try {
      let res = await this.httpClient.post('/recipe', payload)
      return res.data.data
    } catch (error) {
      throw error.response.data.trim()
    }
  }

  async saveRecipe(name: string, ingredients: RecipeItem[]): Promise<string> {
    var errNameEmpty = 'name empty',
      errNamePresent = 'name present',
      errNoIngredients = 'no ingredients'

    if (name.trim() === '') {
      throw errNameEmpty
    }

    if (ingredients.length === 0) {
      throw errNoIngredients
    }

    let found = this.state.find((r) => r.name === name)

    if (found) {
      throw errNamePresent
    }

    try {
      let data = await this.apiSaveRecipe(name, ingredients)

      Object.keys(data).forEach((name: any) => {
        let id = data[name]
        this.state.push({
          id: Number(id),
          name: String(name),
          items: ingredients,
          enabled: true
        })
      })
    } catch (error) {
      switch (error) {
        case 'empty name':
          throw errNameEmpty
        case 'item type already present':
          throw errNamePresent
        case 'no ingredients provided':
          throw errNoIngredients
        default:
          break
      }
    }

    return ''
  }

  private async apiFetchRecipes(): Promise<any> {
    try {
      let res = await this.httpClient.get('/recipe')
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

  private async apiToggleRecipeStatus(id: string, enabled: boolean): Promise<Recipe> {
    try {
      let payload = {
        enabled: enabled
      }
      let response = await this.httpClient.patch(`/recipe/${id}`, payload)
      return response.data.data
    } catch (error) {
      throw error.response.data.trim()
    }
  }

  async toggleRecipeStatus(id: string, enabled: boolean): Promise<Recipe> {
    try {
      let item = await this.apiToggleRecipeStatus(id, enabled)
      this.saveToState(item)
      return item
    } catch (error) {
      throw `error changing recipe status: ${error}`
    }
  }

  private saveToState(recipe: Recipe) {
    this.state = this.state.filter((i) => i.id !== recipe.id)
    this.state.push(recipe)
  }

  getState(): Recipe[] {
    return [...this.state]
  }

  getEnabledRecipes(): Recipe[] {
    return this.getState().filter((r) => r.enabled)
  }

  getByID(id: number): Recipe {
    let r = this.state.find((r) => r.id === id)

    if (!r) {
      throw `recipe with id ${id} not found`
    }

    return r
  }
}
