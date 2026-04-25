/** @enum */ var Foo = {A: 1}
/** @return {*} */ function getFoo() { return Foo; }
var x = /** @type {typeof Foo} */ (getFoo());
