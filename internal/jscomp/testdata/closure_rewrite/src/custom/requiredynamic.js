goog.module('rd.client');
async function load() {
  const {Foo} = await goog.requireDynamic('rd.lib');
  const lib = await goog.requireDynamic('rd.lib');
  goog.requireDynamic('rd.lib').then(({Foo}) => { new Foo(); });
  goog.requireDynamic('rd.lib').then((lib) => lib.Foo);
  return [Foo, lib];
}
exports.load = load;
