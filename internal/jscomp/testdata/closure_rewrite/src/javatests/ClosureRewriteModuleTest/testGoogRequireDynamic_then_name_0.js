async function test() {
  goog.requireDynamic('a.b.c').then((foo) => {console.log(foo.Foo);});
}
