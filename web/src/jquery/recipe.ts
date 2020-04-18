import Client from '../app/inventory/client'
import Client from '../app/recipe/client'
import App, { optionDTO, recipeDTO, ingredientDTO } from '../app/recipe/app'

import $ = require('jquery')
import { Page } from '../app/recipe/app'

export function initializeRecipe(inv: Client, recipe: Client): void {
  let itemNameDropdown = $('#recipeItemName'),
    itemQtyPicker = $('#recipeItemQty'),
    addItemBtn = $('#addRecipeItem'),
    saveRecipeBtn = $('#saveRecipe'),
    recipeNameInput = $('#recipeName')

  let page: Page = {
    ingredientID: (): number => getNumberValue(itemNameDropdown),
    removeIngredientFromDropdown: (s: string): void => removeIngredientNameFromTheList(s),
    ingredientQty: (): number => getNumberValue(itemQtyPicker),

    resetQty: () => resetValue(itemQtyPicker),
    toggleQtyError: (v: boolean) => toggleQtyError(v),

    toggleAddToListBtnDisabledState: (v: boolean): void => toggleDisabledState(v, addItemBtn),

    recipeName: (): string => getStringValue(recipeNameInput),
    resetRecipeName: () => resetValue(recipeNameInput),
    toggleRecipeNameError: (v: boolean): void => toggleRecipeNameError(v),
    toggleAddRecipeButtonState: (v: boolean) => toggleDisabledState(v, saveRecipeBtn),
    toggleNoIngredientsError: (v: boolean): void => toggleNoIngredientsError(v),
    toggleNoUniqueNameError: (v: boolean): void => toggleNoUniqueNameErr(v),
    populateIngredientsDropdown: (dtos: optionDTO[]): void => populateDropdown(itemNameDropdown, dtos),
    populateIngredientsTable: (dtos: ingredientDTO[]): void => populateIngredientsTable(dtos),
    populateTable: (rows: recipeDTO[]): void => populateRecipeTable(rows)
  }

  let app = new App(inv, recipe, page)

  $('#recipe-tab').on('click', () => {
    app.show()
  })

  itemQtyPicker.on('change', () => {
    app.onIngredientQtyChange()
  })

  addItemBtn.on('click', () => {
    app.onAddIngredient()
  })

  recipeNameInput.on('keyup', () => {
    app.onRecipeNameChange()
  })

  saveRecipeBtn.on('click', () => {
    app.onSaveRecipe()
  })

  app.init()
}

function toggleDisabledState(v: boolean, input: JQuery): void {
  input.prop('disabled', v)
}

function resetValue(input: JQuery): void {
  input.val('')
}

function getNumberValue(input: JQuery): number {
  return Number(input.val())
}

function getStringValue(input: JQuery): string {
  return String(input.val())
}

function toggleQtyError(v: boolean): void {
  if (v) {
    $('#recipeItemQtyErr.invalid-feedback').addClass('d-block')
    return
  }
  $('#recipeItemQtyErr.invalid-feedback').removeClass('d-block')
}

function toggleNoIngredientsError(v: boolean): void {
  if (v) {
    $('#noRowsErr.invalid-feedback').addClass('d-block')
    return
  }
  $('#noRowsErr.invalid-feedback').removeClass('d-block')
}

function toggleNoUniqueNameErr(v: boolean): void {
  if (v) {
    $('#noUniqueNameErr.invalid-feedback').addClass('d-block')
    return
  }
  $('#noUniqueNameErr.invalid-feedback').removeClass('d-block')
}

function toggleRecipeNameError(v: boolean): void {
  if (v) {
    $('#noNameErr.invalid-feedback').addClass('d-block')
    return
  }
  $('#noNameErr.invalid-feedback').removeClass('d-block')
}

function removeIngredientNameFromTheList(op: string): void {
  $(`#recipeItemName option[value='${op}']`).remove()
}

const byID = (i1: { id: number }, i2: { id: number }) => i1.id - i2.id

function populateRecipeTable(recipes: recipeDTO[]): void {
  let rows = recipes.sort(byID)
  let table = <HTMLTableElement>$('#recipes tbody')[0]
  $('#recipes tbody tr').remove()
  rows.forEach((element: recipeDTO) => {
    let row = <HTMLTableRowElement>table.insertRow(0)
    row.insertCell(0).innerHTML = String(element.id)
    row.insertCell(1).innerHTML = element.name
  })
}

const byName = (i1: { name: string }, i2: { name: string }) => i1.name.localeCompare(i2.name)

function populateIngredientsTable(dtos: ingredientDTO[]): void {
  let rows = dtos.sort(byName)
  let recipeItemsTable = <HTMLTableElement>$('#recipeItems tbody')[0]
  $('#recipeItems tbody tr').remove()
  rows.forEach((item: ingredientDTO) => {
    let row = <HTMLTableRowElement>recipeItemsTable.insertRow(0)
    row.insertCell(0).innerHTML = item.name
    row.insertCell(1).innerHTML = String(item.qty)
  })
}

function populateDropdown(input: JQuery, options: optionDTO[]) {
  input.empty()
  options.forEach((item: optionDTO) => {
    input.append(new Option(item.name, String(item.id)))
  })
}
