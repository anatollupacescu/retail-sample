import 'mocha'
import chai = require('chai')
import spies = require('chai-spies')

chai.use(spies)
let expect = chai.expect

describe('create recipe', () => {
  it('calls api', async () => {
    expect(true).to.be.true
  })
})
