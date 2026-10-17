var f1 = () => 1;
var f2 = (a, b) => { return a + b; };
function outer() {
  var g = () => this.x + arguments.length;
  return [1, 2].map((x) => x * this.factor);
}
class A extends B {
  constructor() {
    super();
    this.h = () => this;
  }
}
