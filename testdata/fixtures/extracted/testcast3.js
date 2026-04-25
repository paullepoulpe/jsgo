/** @constructor */function base() {}
/** @constructor @extends {base} */function derived() {}
/** @type {!derived} */ var baz = new base();
