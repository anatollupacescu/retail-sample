import InventoryClient from '../retailapp/inventory'
import RecipeClient from '../retailapp/recipe'

import $ = require('jquery')

export function initializeRecipe(inv: InventoryClient, recipe: RecipeClient) {
  $('#recipe-tab').on('click', () => {
    populateDropdown(inv.getInventory())
  })

  recipe.fetchRecipes().then(() => {
    populateTable(recipe.getRecipes())
  })

  //link buttons etc...
}

function populateDropdown(items: any[]) {
  $('#itemType').empty()
  items.map(item => {
    $('#itemType').append(new Option(item.name))
  })
}

function populateTable(_items: any[]) {}
