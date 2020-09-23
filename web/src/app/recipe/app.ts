import InventoryClient, { inventoryItem } from '../inventory/client'
import Client, { Recipe, RecipeItem } from './client'

export interface optionDTO {
  id: number
  name: string
}

export interface recipeRecordDTO {
  id: number
  name: string
  enabled: boolean
}

export interface Page {
  clearRow(): void
  highlightRow(s: string): void

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
  populateTable(rows: recipeRecordDTO[]): void
  ingredientQty(): number
  resetQty(): void
  toggleQtyError(v: boolean): void
}

interface ingredient {
  id: number
  qty: number
}

export interface Modal {
  toggle(v: boolean): void
  populate(i: recipeDTO): void
}

export interface ingredientDTO {
  name: string
  qty: string
}

export interface recipeDTO {
  id: number
  name: string
  items: ingredientDTO[]
  enabled: boolean
}

export default class App {
  private inventory: InventoryClient
  private client: Client
  private page: Page
  private modal: Modal

  // belongs here because does not need persistence
  private ingredients: ingredient[] = []
  private selectedID: string = ''

  constructor(inv: InventoryClient, recipe: Client, page: Page, modal: Modal) {
    this.inventory = inv
    this.client = recipe
    this.page = page
    this.modal = modal
  }

  show() {
    this.renderIngredientsDropdown()
  }

  init() {
    this.client.fetchRecipes().then(() => {
      let recipes: Recipe[] = this.client.getState()
      let rows: recipeRecordDTO[] = this.toRows(recipes)
      this.page.populateTable(rows)
    })
  }

  private toRows(recipes: Recipe[]): recipeRecordDTO[] {
    return recipes.map((r) => ({
      id: r.id,
      name: r.name,
      enabled: r.enabled
    }))
  }

  private renderIngredientsDropdown() {
    let allInventory = this.inventory.getEnabledItems()
    let filteredOptions = this.removeExisting(allInventory)
    let dtos = this.toOptionDTO(filteredOptions)
    this.page.populateIngredientsDropdown(dtos)
  }

  private toOptionDTO(filteredOptions: inventoryItem[]): optionDTO[] {
    return filteredOptions.map((i) => ({
      id: i.id,
      name: i.name
    }))
  }

  private removeExisting(items: inventoryItem[]): inventoryItem[] {
    return items.filter((item) => {
      let found = this.ingredients.find((i) => i.id === item.id)
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

  showModal() {
    if (!this.selectedID) {
      throw 'no row selected'
    }

    let id = Number.parseInt(this.selectedID)
    let recipe = this.client.getByID(id)
    let recipeDTO = this.toRecipeDTO(recipe)

    this.modal.populate(recipeDTO)
    this.openModal()
  }

  toRecipeDTO(r: Recipe): recipeDTO {
    let dtos = []

    for (let i in r.items) {
      let item = r.items[i]
      let name = this.inventory.findByID(String(item.id)).name
      dtos.push({
        name: name,
        qty: String(item.qty)
      })
    }

    return {
      id: r.id,
      name: r.name,
      items: dtos,
      enabled: r.enabled
    }
  }

  closeModal() {
    this.modal.toggle(false)
  }

  openModal() {
    this.modal.toggle(true)
  }

  onRowClick(id: string) {
    this.page.clearRow()

    if (id === this.selectedID) {
      this.selectedID = ''
      return
    }

    this.selectedID = id
    this.page.highlightRow(id)
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

    this.client.saveRecipe(name, recipeItems).then((msg) => {
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
    let dtos = this.ingredients.map((i) => {
      let id = i.id.toString()
      let item = this.inventory.findByID(id)
      return {
        qty: String(i.qty),
        name: item.name
      }
    })

    this.page.populateIngredientsTable(dtos)
  }

  private toRecipeItems(ingredients: ingredient[]): RecipeItem[] {
    return ingredients.map((i) => ({
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

  toggleRecipeStatus(b: boolean): void {
    let id = this.selectedID

    if (!id) {
      throw 'no row selected'
    }

    let item = this.client.getByID(Number(id))

    if (item.enabled === b) {
      throw 'already in the expected state'
    }

    this.client
      .toggleRecipeStatus(id, b)
      .then((r) => this.toRecipeDTO(r))
      .then(this.modal.populate)
      .then(() => {
        let data = this.client.getState()
        this.page.populateTable(data)
      })
      .then(() => this.page.highlightRow(id))
  }
}
