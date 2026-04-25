/** @const */ var ns1 = {};
/** @constructor */ ns1.Foo = function() {};
const ns2 = ns1;
// Verify this works although we have never seen any assignments to `ns2.Foo`
/** @const {typeof ns2.Foo} */ var Foo2 = /** @type {?} */ (x);
/** @type {null} */ var x = new Foo2();
