class Student {
  fullName: string;
  constructor(
    public firstName: string,
    public middleInitial: string,
    public lastName: string
  ) {
    this.fullName = firstName + " " + middleInitial + " " + lastName;
  }

  doGreet() {
    return "Hello, " + this.fullName;
  }
}

document.addEventListener("DOMContentLoaded", function() {
  let user = new Student("Jane", "M.", "User");
  console.log(user.doGreet());
  let textInputElement = document.querySelector("#message");
  textInputElement.addEventListener("keyup", function() {
    console.log("Got keyup event.");
  });
});
