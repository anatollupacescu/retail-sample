var Student = /** @class */ (function () {
    function Student(firstName, middleInitial, lastName) {
        this.firstName = firstName;
        this.middleInitial = middleInitial;
        this.lastName = lastName;
        this.fullName = firstName + " " + middleInitial + " " + lastName;
    }
    Student.prototype.doGreet = function () {
        return "Hello, " + this.fullName;
    };
    return Student;
}());
document.addEventListener("DOMContentLoaded", function () {
    var user = new Student("Jane", "M.", "User");
    console.log(user.doGreet());
    var textInputElement = document.querySelector("#message");
    textInputElement.addEventListener("keyup", function () {
        console.log("Got keyup event.");
    });
});
