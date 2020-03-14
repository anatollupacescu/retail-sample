import InventoryClient from '../retailapp/inventory'
import RecipeClient from '../retailapp/recipe'

import $ = require('jquery')

export function initializeRecipe(inv: InventoryClient, recipe: RecipeClient) {
  let itemNameComponent = $('#recipeItemName')

  $('#recipe-tab').on('click', () => {
    populateDropdown(itemNameComponent, inv.getInventory())
  })

  recipe.fetchRecipes().then(() => {
    populateTable(recipe.getRecipes())
  })

  $('#addRecipeItem').on('click', () => {
    let id = <number>itemNameComponent.val()
    let qty = <number>$('#recipeItemQty').val()
    let res = recipe.addIngredient(id, qty)
    console.log(res)
  })
  //link buttons etc...
}

function populateDropdown(component: any, items: any[]) {
  component.empty()
  items.map(item => {
    component.append(new Option(item.name, item.id))
  })
}

function populateTable(_items: any[]) {}
