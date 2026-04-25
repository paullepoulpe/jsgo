/** @const */ var ns1 = {};
/** @type {string} */ ns1.foo;

/** @const {typeof ns1} */ var ns = /** @type {?} */ (x);

/** @type {string} */ ns1.bar;

/** @type {null} */ var x = ns.bar;
