async function test() {
  goog.requireDynamic('a.b.c').then(({Foo}) => {console.log(Foo);});
}
