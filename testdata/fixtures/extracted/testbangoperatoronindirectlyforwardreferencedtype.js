/** @return {!b.Foo} */
function f() {}
/** @const {!b.Foo} */
var x = f();

const a = {};
/** @typedef {?number} */
a.Foo;
const b = a;
