import InventoryClient, { inventoryItem } from '../inventory/client'
import RecipeClient, { Recipe } from './client'

export interface ingredientDTO {
  id: number
  name: string
  qty: number
}

export interface optionDTO {
  id: number
  name: string
}

export interface recipeDTO {
  id: number
  name: string
}

export interface Page {
  ingredientID(): number
  removeIngredientFromDropdown(s: string): void
  toggleIngredientNameError(v: boolean): void
  toggleAddIngredientButtonState(v: boolean): void
  toggleNoIngredientsError(v: boolean): void
  recipeID(): number
  resetRecipeName(): void
  toggleAddRecipeButtonState(v: boolean): void
  toggleRecipeNameError(v: boolean): void
  populateIngredientsDropdown(dtos: optionDTO[]): void
  populateIngredientsTable(dtos: ingredientDTO[]): void
  populateTable(rows: recipeDTO[]): void
  ingredientQty(): number
  resetQty(): void
  toggleQtyError(v: boolean): void
}

export default class App {
  private inventory: InventoryClient
  private client: RecipeClient
  private page: Page

  constructor(inv: InventoryClient, recipe: RecipeClient, page: Page) {
    this.inventory = inv
    this.client = recipe
    this.page = page
  }

  init() {
    this.client.fetchRecipes().then(() => {
      let recipes: Recipe[] = this.client.getRecipes()
      let rows: recipeDTO[] = this.toRows(recipes)
      this.page.populateTable(rows)
    })
  }

  toRows(recipes: Recipe[]): recipeDTO[] {
    return []
  }

  show() {
    let dropdownOptions = this.inventory.getInventory()
    let filteredOptions = this.removeExisting(dropdownOptions)
    let dtos = this.toOptionDTO(filteredOptions)
    this.page.populateIngredientsDropdown(dtos)
  }

  toOptionDTO(filteredOptions: inventoryItem[]): optionDTO[] {
    return []
  }

  removeExisting(dropdownOptions: inventoryItem[]): inventoryItem[] {
    return []
  }

  onRecipeNameChange() {
    throw new Error('Method not implemented.')
  }

  onSaveRecipe() {
    let recipeID = this.page.recipeID()
    this.client.setName('WIP')
    this.client.saveRecipe().then(msg => {
      switch (msg) {
        case 'ERR_EMPTY':
        case 'name empty': {
          this.page.toggleRecipeNameError(true)
          return
        }
        case 'ERR_PRESENT':
        case 'name present': {
          alert('duplicate name')
          return
        }
        case 'no ingredients': {
          this.page.toggleNoIngredientsError(true)
          return
        }
        case '':
          break
        default: {
          throw 'unknown error'
        }
      }

      //  success
    })
  }

  onAddIngredient() {
    throw new Error('Method not implemented.')
  }

  onIngredientQtyChange() {
    throw new Error('Method not implemented.')
  }
}
