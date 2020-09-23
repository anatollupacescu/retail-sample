import { apiIsHealthy } from '../health'

class Checker {
  private endpoint: string

  constructor(public url: string) {
    this.endpoint = `${url}/health`
  }

  check() {
    return apiIsHealthy(this.endpoint)
  }
}

document.addEventListener('DOMContentLoaded', function () {
  document.querySelector('#message')?.addEventListener('keyup', function () {
    // let msg = document.querySelector('#message') as HTMLInputElement
    let healthChecker = new Checker('http://localhost:8080')
    console.log(healthChecker.check())
  })
})
