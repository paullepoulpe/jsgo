goog.loadModule(function(exports) {
  goog.module('rd.lib');
  exports.Foo = class {};
  return exports;
});
goog.loadModule(function(exports) {
  goog.module('rd.client2');
  async function load() {
    const {Foo} = await goog.requireDynamic('rd.lib');
    const lib = await goog.requireDynamic('rd.lib');
    goog.requireDynamic('rd.lib').then(({Foo}) => { new Foo(); });
    goog.requireDynamic('rd.lib').then((lib) => lib.Foo);
    let bad; bad = await goog.requireDynamic('rd.lib');
    return [Foo, lib, bad];
  }
  exports.load = load;
  return exports;
});
