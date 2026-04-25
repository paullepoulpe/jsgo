/** @const */ var ns1 = {};
/** @type {string} */ ns1.foo;
/** @const */ var ns2 = {};
/** @type {number} */ ns2.bar;

/** @const {typeof ns2} */ var ns = /** @type {?} */ (ns1);
/** @type {null} */ var x = ns.foo;
