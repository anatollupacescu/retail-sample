import InventoryClient, { inventoryItem } from '../inventory/client'
import Client, { Recipe, RecipeItem } from './client'

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
  toggleAddToListBtnState(v: boolean): void
  toggleNoIngredientsError(v: boolean): void

  recipeName(): string
  resetRecipeName(): void
  toggleAddRecipeButtonState(v: boolean): void
  toggleRecipeNameError(v: boolean): void
  toggleNoUniqueNameError(v: boolean): void

  populateIngredientsDropdown(dtos: optionDTO[]): void
  populateIngredientsTable(dtos: ingredientDTO[]): void
  populateTable(rows: recipeDTO[]): void
  ingredientQty(): number
  resetQty(): void
  toggleQtyError(v: boolean): void
}

interface ingredient {
  id: number
  qty: number
}

export interface ingredientDTO {
  name: string
  qty: string
}

export default class App {
  private inventory: InventoryClient
  private client: Client
  private page: Page

  // belongs here because does not need persistence
  private ingredients: ingredient[] = []

  constructor(inv: InventoryClient, recipe: Client, page: Page) {
    this.inventory = inv
    this.client = recipe
    this.page = page
  }

  show() {
    this.renderIngredientsDropdown()
  }

  init() {
    this.client.fetchRecipes().then(() => {
      let recipes: Recipe[] = this.client.getState()
      let rows: recipeDTO[] = this.toRows(recipes)
      this.page.populateTable(rows)
    })
  }

  private toRows(recipes: Recipe[]): recipeDTO[] {
    return recipes.map(r => ({
      id: r.id,
      name: r.name
    }))
  }

  private renderIngredientsDropdown() {
    let dropdownOptions = this.inventory.getState()
    let filteredOptions = this.removeExisting(dropdownOptions)
    let dtos = this.toOptionDTO(filteredOptions)
    this.page.populateIngredientsDropdown(dtos)
  }

  private toOptionDTO(filteredOptions: inventoryItem[]): optionDTO[] {
    return filteredOptions.map(i => ({
      id: i.id,
      name: i.name
    }))
  }

  private removeExisting(dropdownOptions: inventoryItem[]): inventoryItem[] {
    return dropdownOptions.filter(dop => {
      let found = this.ingredients.find(i => i.id === dop.id)
      return !found
    })
  }

  private nameIsValid(name: string): boolean {
    if (!name || name.trim().length === 0) {
      return false
    }
    return true
  }

  private refreshRecipeNameRelevantUI(name: string): void {
    if (!this.nameIsValid(name)) {
      this.page.toggleRecipeNameError(true)
      this.page.toggleAddRecipeButtonState(true)
      return
    }
    this.page.toggleRecipeNameError(false)
    this.page.toggleAddRecipeButtonState(false)
  }

  onRecipeNameChange() {
    let name = this.page.recipeName()
    this.refreshRecipeNameRelevantUI(name)
  }

  onSaveRecipe() {
    let name = this.page.recipeName()
    this.refreshRecipeNameRelevantUI(name)

    if (this.ingredients.length === 0) {
      this.page.toggleNoIngredientsError(true)
      return
    }

    this.page.toggleNoIngredientsError(false)

    let recipeItems = this.toRecipeItems(this.ingredients)

    this.client.saveRecipe(name, recipeItems).then(msg => {
      switch (msg) {
        case 'name empty': {
          this.page.toggleRecipeNameError(true)
          break
        }
        case 'name present': {
          this.page.toggleNoUniqueNameError(true)
          break
        }
        case 'no ingredients': {
          this.page.toggleNoIngredientsError(true)
          break
        }
        case '': //success
          this.ingredients = []
          this.renderIngredientsDropdown()
          this.populateIngredientsTable()
          this.page.toggleNoUniqueNameError(false)
          this.page.resetRecipeName()
          this.page.toggleAddRecipeButtonState(true)
          let recipes = this.client.getState()
          let rows = this.toRows(recipes)
          this.page.populateTable(rows)
          break
        default: {
          throw 'unknown error'
        }
      }
    })
  }

  private populateIngredientsTable() {
    let dtos = this.ingredients.map(i => ({
      qty: String(i.qty),
      name: this.inventory.getName(i.id)
    }))
    this.page.populateIngredientsTable(dtos)
  }

  private toRecipeItems(ingredients: ingredient[]): RecipeItem[] {
    return ingredients.map(i => ({
      id: i.id,
      qty: i.qty
    }))
  }

  onAddIngredient() {
    let qty = this.page.ingredientQty()

    if (this.badQuantity(qty)) {
      this.page.toggleQtyError(true)
      return
    }

    let id = this.page.ingredientID()

    this.ingredients.push({
      id: id,
      qty: qty
    })

    this.populateIngredientsTable()
    this.page.toggleNoIngredientsError(false)

    this.renderIngredientsDropdown()

    this.page.resetQty()
    this.page.toggleAddToListBtnState(true)
  }

  private badQuantity(qty: number): boolean {
    if (!qty || Number(qty) <= 0) {
      return true
    }
    return false
  }

  onIngredientQtyChange() {
    let qty = this.page.ingredientQty()

    if (this.badQuantity(qty)) {
      this.page.toggleQtyError(true)
      this.page.toggleAddToListBtnState(true)
      return
    }

    let id = this.page.ingredientID()

    if (!id || id === 0) {
      this.page.toggleAddToListBtnState(true)
      return
    }

    this.page.toggleQtyError(false)
    this.page.toggleAddToListBtnState(false)
  }
}
