var a;
class Foo {
  constructor(a) {
    this.a = a;
  }
  f() {
    var x = 1;
    return a + x;
  }
}
