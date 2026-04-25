class Foo {};
/** @enum {number} */ Foo.E = {A: 1};
const F = Foo;
const E = F.E;
/** @type {E} */ let e = undefined;
