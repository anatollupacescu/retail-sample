import { apiIsHealthy } from '../health'

class Student {
  fullName: string

  constructor(public firstName: string, public middleInitial: string, public lastName: string) {
    this.fullName = firstName + ' ' + middleInitial + ' ' + lastName
  }

  doGreet() {
    return this.fullName + ' says ' + apiIsHealthy('localhost:8080/health')
  }
}

document.addEventListener('DOMContentLoaded', function() {
  document.querySelector('#message')?.addEventListener('keyup', function() {
    let msg = document.querySelector('#message') as HTMLInputElement
    let user = new Student('Jane', msg.value, 'Shapokleak')
    console.log(user.doGreet())
  })
})
