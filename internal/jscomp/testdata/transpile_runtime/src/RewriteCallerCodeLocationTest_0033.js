class Foo {
  bar(here = goog.callerLocation()) {}
}
const baz = new Foo();
baz.bar()
