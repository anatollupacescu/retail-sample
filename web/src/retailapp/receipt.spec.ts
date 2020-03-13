import 'mocha'
import chai = require('chai')
import spies = require('chai-spies')

import ReceiptClient from './receipt'

chai.use(spies)
let expect = chai.expect

describe('add ingredients', () => {
  describe('when qty is 0', () => {
    let app = new ReceiptClient('')
    let result = app.addIngredient(1, 0)

    it('should err', async () => {
      expect(result).to.equal('zero quantity')
      expect(app.listItems()).to.have.length(0)
    })
  })

  describe('when ok', () => {
    //wip
    let app = new ReceiptClient('')
    app.addIngredient(1, 1)

    it('should save the ingredient', () => {
      expect(app.listItems()).to.have.length(1)
    })
  })
})

describe('saving a receipt', () => {
  describe('when name is empty', () => {
    let app = new ReceiptClient('')

    it('Should err', async () => {
      let result = await app.saveReceipt()
      expect(result).to.equal('name empty')
      expect(app.getReceipts()).to.have.length(0)
    })
  })

  describe('when receipt name is already taken', () => {
    let state = [
      {
        id: 1,
        name: 'test'
      }
    ]

    let app = new ReceiptClient('', state)
    app.setName('test')

    let mockApi = chai.spy.on(app, 'apiSaveReceipt')

    it('Should err', async () => {
      let result = await app.saveReceipt()
      expect(mockApi).to.not.have.been.called
      expect(result).to.equal('name present')
      expect(app.getReceipts()).to.have.length(1)
    })
  })

  describe('when receipt does not have ingredients', () => {
    let app = new ReceiptClient('')
    app.setName('test')

    let mockApi = chai.spy.on(app, 'apiSaveReceipt')

    it('should err', async () => {
      let result = await app.saveReceipt()
      expect(mockApi).to.not.have.been.called
      expect(result).to.equal('no ingredients')
      expect(app.getReceipts()).to.have.length(0)
    })
  })

  describe('when receipt is correct', () => {
    let app = new ReceiptClient('')
    app.setName('test')
    app.addIngredient(1, 1)

    let mockApi = chai.spy.on(app, 'apiSaveReceipt', () => ({
      data: {
        data: {
          test: 1
        }
      }
    }))

    it('makes the api call', async () => {
      await app.saveReceipt()
      expect(mockApi).to.have.been.called
      expect(app.getReceipts()).to.have.length(1)
    })
  })
})

describe('fetching receipts', () => {
  let app = new ReceiptClient('')

  var mockApi = chai.spy.on(app, 'apiFetchReceipts', () => ({
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
    await app.fetchReceipts()
    expect(mockApi).to.have.been.called.exactly(1)
    expect(app.getReceipts()).to.have.length(1)
  })
})
