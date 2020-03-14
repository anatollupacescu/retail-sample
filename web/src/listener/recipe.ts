import InventoryClient from '../retailapp/inventory'
import RecipeClient, { RecipeItem } from '../retailapp/recipe'

import $ = require('jquery')

export function initializeRecipe(inv: InventoryClient, recipe: RecipeClient) {
  let itemNameDropdown = $('#recipeItemName'),
    itemQtyPicker = $('#recipeItemQty'),
    addItemBtn = $('#addRecipeItem'),
    saveRecipeBtn = $('#saveRecipe'),
    recipeNameInput = $('#recipeName')

  onClickRecipeTab_updateIngredientsNameDropdown(itemNameDropdown, inv)
  onChangeQty_resetErrorMessage(itemQtyPicker)
  onClickAddItemBtn_addIngredientToPendingRecipe(addItemBtn, itemNameDropdown, itemQtyPicker, recipe, inv)
  onClickSaveRecipe_saveAndResetStage(saveRecipeBtn, recipeNameInput, recipe)
  onChangeRecipeName_resetNoNameError(recipeNameInput)

  //fetch main table data
  recipe.fetchRecipes().then(() => {
    populateTable(recipe.getRecipes())
  })
}

function onChangeRecipeName_resetNoNameError(recipeNameInput: any): void {
  recipeNameInput.keyup(() => {
    $('#noNameErr.invalid-feedback').removeClass('d-block')
  })
}

function onClickSaveRecipe_saveAndResetStage(btn: any, recipeNameInput: any, recipe: RecipeClient): void {
  btn.on('click', () => {
    let recipeName = <string>recipeNameInput.val()
    recipe.setName(recipeName)
    recipe
      .saveRecipe()
      .then(msg => {
        switch (msg) {
          case 'name empty': {
            showNoNameError()
            return
          }
          case 'name present': {
            alert('duplicate name')
            return
          }
          case 'no ingredients': {
            showNoIngredientsError()
            return
          }
          default: {
            alert('hooray')
          }
        }
      })
      .catch(() => {
        alert('got an error')
      })
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

function onClickAddItemBtn_addIngredientToPendingRecipe(
  btn: any,
  itemNameDropdown: any,
  itemQtyPicker: any,
  recipe: RecipeClient,
  inv: InventoryClient
): void {
  btn.on('click', () => {
    let id = <number>itemNameDropdown.val()
    let qty = <number>itemQtyPicker.val()
    let err = recipe.addIngredient(Number(id), Number(qty))
    if (err) {
      showAddIngredientError(err)
      return
    }
    populateIngredientsTable(recipe, inv)
    resetNoIngredientsError()
  })
}

function populateIngredientsTable(recipe: RecipeClient, inv: InventoryClient): void {
  $('#recipeItems tbody tr').remove()

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

function onChangeQty_resetErrorMessage(component: any): void {
  component.on('change', () => {
    $('#recipeItemQtyErr.invalid-feedback').removeClass('d-block')
  })
}

function onClickRecipeTab_updateIngredientsNameDropdown(component: any, inv: InventoryClient): void {
  $('#recipe-tab').on('click', () => {
    populateDropdown(component, inv.getInventory())
  })
}

function populateDropdown(component: any, items: any[]) {
  component.empty()
  items.map(item => {
    component.append(new Option(item.name, item.id))
  })
}

function populateTable(_items: any[]) {}
