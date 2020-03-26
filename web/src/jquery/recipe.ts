import InventoryClient, { inventoryItem } from '../app/inventory/client'
import RecipeClient, { RecipeItem, Recipe } from '../app/recipe/client'

import $ = require('jquery')

let itemNameDropdown: JQuery<HTMLElement>,
  itemQtyPicker: JQuery<HTMLElement>,
  addItemBtn: JQuery<HTMLElement>,
  saveRecipeBtn: JQuery<HTMLElement>,
  recipeNameInput: JQuery<HTMLElement>

export function initializeRecipe(inv: InventoryClient, recipe: RecipeClient): void {
  itemNameDropdown = $('#recipeItemName')
  itemQtyPicker = $('#recipeItemQty')
  addItemBtn = $('#addRecipeItem')
  saveRecipeBtn = $('#saveRecipe')
  recipeNameInput = $('#recipeName')

  onClickRecipeTab_updateIngredientsNameDropdown(recipe, inv)
  onChangeQty_resetErrorMessage()
  onClickAddRecipeItem_addIngredientToPendingRecipe(recipe, inv)
  onClickSaveRecipe_saveAndResetStage(recipe, inv)
  onChangeRecipeName_resetNoNameError()

  //fetch main table data
  recipe.fetchRecipes().then(() => {
    populateRecipeTable(recipe.getRecipes())
  })
}

function onChangeRecipeName_resetNoNameError(): void {
  recipeNameInput.keyup(() => {
    $('#noNameErr.invalid-feedback').removeClass('d-block')
  })
}

function onClickSaveRecipe_saveAndResetStage(recipe: RecipeClient, inv: InventoryClient): void {
  saveRecipeBtn.on('click', () => {
    let recipeName = <string>recipeNameInput.val()
    recipe.setName(recipeName)
    recipe
      .saveRecipe()
      .then(msg => {
        switch (msg) {
          case 'ERR_EMPTY':
          case 'name empty': {
            showNoNameError()
            return
          }
          case 'ERR_PRESENT':
          case 'name present': {
            alert('duplicate name')
            return
          }
          case 'no ingredients': {
            showNoIngredientsError()
            return
          }
          default: {
            populateRecipeTable(recipe.getRecipes())
            clearIngredientsTable()
            clearRecipeName()
            populateDropdown(recipe, inv)
          }
        }
      })
      .catch(res => {
        alert('got an error: ' + res)
      })
  })
}

function clearRecipeName(): void {
  recipeNameInput.val('')
}

function clearIngredientsTable(): void {
  $('#recipeItems tbody tr').remove()
}

function resetItemCount(): void {
  itemQtyPicker.val(0)
}

function removeIngredientNameFromTheList(op: string): void {
  $(`#recipeItemName option[value='${op}']`).remove()
}

function populateRecipeTable(recipes: Recipe[]): void {
  $('#recipes tbody tr').remove()
  let table = <HTMLTableElement>$('#recipes tbody')[0]
  let rows = recipes.sort((i1: Recipe, i2: Recipe) => i1.id - i2.id)
  rows.forEach((element: Recipe) => {
    let row = <HTMLTableRowElement>table.insertRow(0)
    row.insertCell(0).innerHTML = String(element.id)
    row.insertCell(1).innerHTML = element.name
  })
}

function showNoIngredientsError(): void {
  $('#noRowsErr.invalid-feedback').addClass('d-block')
}

function resetNoIngredientsError(): void {
  $('#noRowsErr.invalid-feedback').removeClass('d-block')
}

function showNoNameError(): void {
  $('#noNameErr.invalid-feedback').addClass('d-block')
}

function onClickAddRecipeItem_addIngredientToPendingRecipe(recipe: RecipeClient, inv: InventoryClient): void {
  addItemBtn.on('click', () => {
    let id = itemNameDropdown.val()
    let qty = <number>itemQtyPicker.val()
    let err = recipe.addIngredient(Number(id), Number(qty))
    if (err) {
      showAddIngredientError(err)
      return
    }
    populateIngredientsTable(recipe, inv)
    resetNoIngredientsError()
    removeIngredientNameFromTheList(String(id))
    resetItemCount()
  })
}

function populateIngredientsTable(recipe: RecipeClient, inv: InventoryClient): void {
  clearIngredientsTable()

  let recipeItemsTable = <HTMLTableElement>$('#recipeItems tbody')[0]

  let rows = recipe.listItems().sort((i1: RecipeItem, i2: RecipeItem) => i1.id - i2.id)

  rows.forEach((item: RecipeItem) => {
    let row = <HTMLTableRowElement>recipeItemsTable.insertRow(0)
    row.insertCell(0).innerHTML = inv.getName(item.id)
    row.insertCell(1).innerHTML = String(item.qty)
  })
}

function showAddIngredientError(s: string): void {
  if (s === 'zero quantity') {
    $('#recipeItemQtyErr.invalid-feedback').addClass('d-block')
    return
  }
  console.log('got error', s)
}

function onChangeQty_resetErrorMessage(): void {
  itemQtyPicker.on('change', () => {
    $('#recipeItemQtyErr.invalid-feedback').removeClass('d-block')
  })
}

function onClickRecipeTab_updateIngredientsNameDropdown(recipe: RecipeClient, inv: InventoryClient): void {
  $('#recipe-tab').on('click', () => {
    populateDropdown(recipe, inv)
  })
}

function populateDropdown(recipe: RecipeClient, inv: InventoryClient) {
  itemNameDropdown.empty()
  inv.getInventory().map(item => {
    if (!isInRecipe(recipe, item)) {
      itemNameDropdown.append(new Option(item.name, String(item.id)))
    }
  })
}

function isInRecipe(recipe: RecipeClient, item: inventoryItem): boolean {
  let found = recipe.listItems().find(i => i.id === item.id)
  return found !== undefined
}
