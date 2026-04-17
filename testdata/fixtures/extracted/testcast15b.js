for (var i = 0; i < 10; i++) {
var x = /** @type {{foo:number}}|number} */ ({foo: 3});
/** @param {number} x */ function f(x) {}
f(x.foo);
f([].foo);
}
