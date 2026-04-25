/** @enum {symbol} */
var ES = {A: Symbol('a'), B: Symbol('b')};

/** @type {!Object<number, number>} */
var o = {};

o[ES.A] = 1;
