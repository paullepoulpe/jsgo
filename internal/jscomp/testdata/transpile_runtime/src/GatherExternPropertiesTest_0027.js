class Foo {
  bar() {
    return this;
  }
}
var baz = new Foo();
var bar = baz.bar;
