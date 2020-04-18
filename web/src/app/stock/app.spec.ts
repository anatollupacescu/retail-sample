import 'mocha'
import chai = require('chai')
import spies = require('chai-spies')

chai.use(spies)
let expect = chai.expect

describe('on provision', () => {
  it('should work', () => {
    expect(true).to.be.true
  })
})
