import InventoryClient from '../app/inventory/client'
import Client from '../app/recipe/client'
import App, { optionDTO, recipeRecordDTO, ingredientDTO, Modal } from '../app/recipe/app'

import $ = require('jquery')
import { Page, recipeDTO } from '../app/recipe/app'

export function initializeRecipe(inv: InventoryClient, recipe: Client): void {
  let itemNameDropdown = $('#recipeItemName'),
    itemQtyPicker = $('#recipeItemQty'),
    addItemBtn = $('#addRecipeItem'),
    saveRecipeBtn = $('#saveRecipe'),
    recipeNameInput = $('#recipeName')

  let page: Page = {
    clearRow: (): void => clearRowSelection(),
    highlightRow: (s: string): void => selectTableRow(s),

    ingredientID: (): number => getNumberValue(itemNameDropdown),
    removeIngredientFromDropdown: (s: string): void => removeIngredientNameFromTheList(s),
    ingredientQty: (): number => getNumberValue(itemQtyPicker),

    resetQty: () => resetValue(itemQtyPicker),
    toggleQtyError: (v: boolean) => toggleQtyError(v),

    toggleAddToListBtnState: (v: boolean): void => toggleDisabledState(v, addItemBtn),

    recipeName: (): string => getStringValue(recipeNameInput),
    resetRecipeName: () => resetValue(recipeNameInput),
    toggleRecipeNameError: (v: boolean): void => toggleRecipeNameError(v),
    toggleAddRecipeButtonState: (v: boolean) => toggleDisabledState(v, saveRecipeBtn),
    toggleNoIngredientsError: (v: boolean): void => toggleNoIngredientsError(v),
    toggleNoUniqueNameError: (v: boolean): void => toggleNoUniqueNameErr(v),
    populateIngredientsDropdown: (dtos: optionDTO[]): void => populateDropdown(itemNameDropdown, dtos),
    populateIngredientsTable: (dtos: ingredientDTO[]): void => populateIngredientsTable(dtos),
    populateTable: (rows: recipeRecordDTO[]): void => populateRecipeTable(rows)
  }

  let modal: Modal = {
    toggle: (v: boolean): void => toggleModal(v),
    populate: (i: recipeDTO): void => populateModal(i)
  }

  let app = new App(inv, recipe, page, modal)

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

  let openModalBtn = $('#openRecipeModalBtn')

  openModalBtn.on('click', () => {
    app.showModal()
  })

  let closeModalClass = $('.close-modal')

  closeModalClass.on('click', () => {
    app.closeModal()
  })

  let table = $('#recipes tbody')

  table.on('click', 'tr', function () {
    let id = $(this).find('td:eq(0)').text()

    app.onRowClick(id)
  })

  let enableRecipe = $('#recipeModalEnable')

  enableRecipe.on('click', () => {
    app.toggleRecipeStatus(true)
  })

  let disableRecipe = $('#recipeModalDisable')

  disableRecipe.on('click', () => {
    app.toggleRecipeStatus(false)
  })

  app.init()
}

let dark = 'list-group-item-dark'

function clearRowSelection(): void {
  $('#recipes tbody tr').removeClass(dark)
}

function selectTableRow(id: string): void {
  $('#recipes tbody tr').each(function () {
    let currentRow = $(this)
    let currentID = currentRow.find('td:eq(0)').text()

    if (id === currentID) {
      currentRow.addClass(dark)
    }
  })
}

function toggleModal(v: boolean): void {
  let el = $('#recipeModal')
  if (v) {
    el.addClass('show')
    el.addClass('d-block')
    return
  }
  el.removeClass('show')
  el.removeClass('d-block')
}

function populateModal(i: recipeDTO): void {
  $('#recipeModalID').html(String(i.id))
  $('#recipeModalName').html(i.name)
  let table = <HTMLTableElement>$('#modalRecipeItems tbody')[0]
  $('#modalRecipeItems tbody tr').remove()
  i.items.forEach((element: ingredientDTO) => {
    let row = <HTMLTableRowElement>table.insertRow(0)
    row.insertCell(0).innerHTML = element.name
    row.insertCell(1).innerHTML = String(element.qty)
  })

  let enbl = $('#recipeModalEnable'),
    dsbl = $('#recipeModalDisable')

  enbl.removeClass('btn-primary')
  dsbl.removeClass('btn-primary')

  if (i.enabled) {
    enbl.addClass('btn-primary')
    return
  }

  dsbl.addClass('btn-primary')
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

function populateRecipeTable(recipes: recipeRecordDTO[]): void {
  let rows = recipes.sort(byID)
  let table = <HTMLTableElement>$('#recipes tbody')[0]
  $('#recipes tbody tr').remove()
  rows.forEach((element: recipeRecordDTO) => {
    let row = <HTMLTableRowElement>table.insertRow(0)
    row.insertCell(0).innerHTML = String(element.id)
    row.insertCell(1).innerHTML = element.name
    if (!element.enabled) {
      row.classList.add('disabled')
    }
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
