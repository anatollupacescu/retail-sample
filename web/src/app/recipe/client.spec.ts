import 'mocha'
import chai = require('chai')
import spies = require('chai-spies')

import chaiAsPromised = require('chai-as-promised')
chai.use(chaiAsPromised)

import axios from "axios"
import RecipeClient from './client'

chai.use(spies)
let expect = chai.expect

describe('saving a recipe', () => {
  describe('when recipe is missing ingredients', () => {
    let app = new RecipeClient(axios, [])

    let mockApi = chai.spy.on(app, 'apiSaveRecipe')

    it('should err', () => {
      expect(app.saveRecipe('test', [])).to.be.rejectedWith('no ingredients')
      expect(mockApi).to.have.been.called.exactly(0)
      expect(app.getState()).to.have.length(0)
    })
  })

  describe('when recipe name is already taken', () => {
    let state = [
      {
        id: 1,
        name: 'test',
        items: [],
        enabled: true
      }
    ]

    let app = new RecipeClient(axios, state)

    let mockApi = chai.spy.on(app, 'apiSaveRecipe')

    it('should err', () => {
      expect(app.saveRecipe('test', [{ id: 1, qty: 2 }])).to.be.rejectedWith('name present')
      expect(mockApi).to.have.been.called.exactly(0)
      expect(app.getState()).to.have.length(1)
    })
  })

  describe('when recipe is correct', () => {
    let app = new RecipeClient(axios)

    let mockApi = chai.spy.on(app, 'apiSaveRecipe', () => ({
      test: 1
    }))

    it('makes the api call', async () => {
      await app.saveRecipe('name', [{ id: 1, qty: 2 }])
      expect(mockApi).to.have.been.called.once
      expect(app.getState()).to.have.length(1)
    })
  })
})

describe('fetching recipes', () => {
  let app = new RecipeClient(axios)

  var mockApi = chai.spy.on(app, 'apiFetchRecipes', () => [
    {
      name: 'item1',
      id: 2
    }
  ])

  it('should make the api call', async () => {
    await app.fetchRecipes()
    expect(mockApi).to.have.been.called.exactly(1)
    expect(app.getState()).to.have.length(1)
  })
})
