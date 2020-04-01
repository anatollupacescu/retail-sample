import 'mocha'
import chai = require('chai')
import spies = require('chai-spies')

import RecipeClient from './client'

chai.use(spies)
let expect = chai.expect

describe('saving a recipe', () => {
  describe('when recipe name is already taken', () => {
    let state = [
      {
        id: 1,
        name: 'test',
        items: []
      }
    ]

    let app = new RecipeClient('', state)

    let mockApi = chai.spy.on(app, 'apiSaveRecipe')

    it('should err', async () => {
      let result = await app.saveRecipe('test', [])
      expect(mockApi).to.not.have.been.called
      expect(result).to.equal('name present')
      expect(app.getRecipes()).to.have.length(1)
    })
  })

  describe('when recipe is correct', () => {
    let app = new RecipeClient('')

    let mockApi = chai.spy.on(app, 'apiSaveRecipe', () => ({
      data: {
        data: {
          test: 1
        }
      }
    }))

    it('makes the api call', async () => {
      await app.saveRecipe('name', [{ id: 1, qty: 2 }])
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
