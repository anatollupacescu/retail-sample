import 'mocha'
import chai = require('chai')
import spies = require('chai-spies')

import RecipeClient from './client'

chai.use(spies)
let expect = chai.expect

describe('add ingredients', () => {
  describe('when qty is 0', () => {
    let app = new RecipeClient('')
    let result = app.addIngredient(1, 0)

    it('should err', async () => {
      expect(result).to.equal('zero quantity')
      expect(app.listItems()).to.have.length(0)
    })
  })

  describe('when ingredient is duplicate', () => {
    let app = new RecipeClient('')
    app.addIngredient(1, 1)
    let result = app.addIngredient(1, 2)

    it('should err', async () => {
      expect(result).to.equal('duplicate id')
      expect(app.listItems()).to.have.length(1)
    })
  })

  describe('when ok', () => {
    //wip
    let app = new RecipeClient('')
    app.addIngredient(1, 1)

    it('should save the ingredient', () => {
      expect(app.listItems()).to.have.length(1)
    })
  })
})

describe('saving a recipe', () => {
  describe('when name is empty', () => {
    let app = new RecipeClient('')

    it('Should err', async () => {
      let result = await app.saveRecipe()
      expect(result).to.equal('name empty')
      expect(app.getRecipes()).to.have.length(0)
    })
  })

  describe('when recipe name is already taken', () => {
    let state = [
      {
        id: 1,
        name: 'test',
        items: []
      }
    ]

    let app = new RecipeClient('', state)
    app.setName('test')

    let mockApi = chai.spy.on(app, 'apiSaveRecipe')

    it('should err', async () => {
      let result = await app.saveRecipe()
      expect(mockApi).to.not.have.been.called
      expect(result).to.equal('name present')
      expect(app.getRecipes()).to.have.length(1)
    })
  })

  describe('when recipe does not have ingredients', () => {
    let app = new RecipeClient('')
    app.setName('test')

    let mockApi = chai.spy.on(app, 'apiSaveRecipe')

    it('should err', async () => {
      let result = await app.saveRecipe()
      expect(mockApi).to.not.have.been.called
      expect(result).to.equal('no ingredients')
      expect(app.getRecipes()).to.have.length(0)
    })
  })

  describe('when recipe is correct', () => {
    let app = new RecipeClient('')
    app.setName('test')
    app.addIngredient(1, 1)

    let mockApi = chai.spy.on(app, 'apiSaveRecipe', () => ({
      data: {
        data: {
          test: 1
        }
      }
    }))

    it('makes the api call', async () => {
      await app.saveRecipe()
      expect(mockApi).to.have.been.called
      expect(app.getRecipes()).to.have.length(1)
    })
  })
})

describe('fetching recipes', () => {
  let app = new RecipeClient('')

  var mockApi = chai.spy.on(app, 'apiFetchRecipes', () => ({
    data: {
      data: [
        {
          name: 'item1',
          id: 2
        }
      ]
    }
  }))

  it('should make the api call', async () => {
    await app.fetchRecipes()
    expect(mockApi).to.have.been.called.exactly(1)
    expect(app.getRecipes()).to.have.length(1)
  })
})
