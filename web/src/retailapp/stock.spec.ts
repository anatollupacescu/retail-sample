import 'mocha'
import chai = require('chai')
import spies = require('chai-spies')

chai.use(spies)
let expect = chai.expect

describe('provision stock', () => {
  it('should succeed', function() {
    //here we should spy on api methods and all the rest...
    expect(true).to.be.true
  })
})
