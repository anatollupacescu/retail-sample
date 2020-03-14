import InventoryClient from '../retailapp/inventory'
import RecipeClient from '../retailapp/recipe'

import $ = require('jquery')

export function initializeRecipe(inv: InventoryClient, recipe: RecipeClient) {
  let itemNameDropdown = $('#recipeItemName')

  $('#recipe-tab').on('click', () => {
    populateDropdown(itemNameDropdown, inv.getInventory())
  })

  recipe.fetchRecipes().then(() => {
    populateTable(recipe.getRecipes())
  })

  $('#addRecipeItem').on('click', () => {
    let id = <number>itemNameDropdown.val()
    let qty = <number>$('#recipeItemQty').val()
    let res = recipe.addIngredient(Number(id), Number(qty))
    if (res) {
      showAddIngredientError(res)
    }
  })

  //link buttons etc...
}

function showAddIngredientError(_s: string): void {}

function populateDropdown(component: any, items: any[]) {
  component.empty()
  items.map(item => {
    component.append(new Option(item.name, item.id))
  })
}

function populateTable(_items: any[]) {}
