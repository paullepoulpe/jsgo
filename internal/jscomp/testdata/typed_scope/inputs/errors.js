/** @constructor */ var NoInit;
/** @interface */ var NoInitIface;
/** @enum {string} */ var BadEnum = 3;
/** @type {number} */ var m1 = 1, m2 = 2;
/** @extends {Object} */ function NotCtor() {}
/** @struct */ function NotCtor2() {}
/** @dict */ function NotCtor3() {}
/** @implements {Object} */ function NotCtor4() {}
/** @constructor @extends {number} */ function ExtendsNum() {}
/** @constructor @implements {number} */ function ImplNum() {}
/** @constructor @extends {Undeclared} */ function ExtendsUndeclared() {}
/** @interface @implements {Object} */ function IfaceImpl() {}
/** @constructor */ function RealCtor() {}
/** @interface */ var AliasIface = RealCtor;
/** @record */ var AliasRec = RealCtor;
/**
 * @constructor
 * @return {number}
 */
function String(opt_str) {}
/** @constructor @template T @param {T} t */ function Dup(t) {}
