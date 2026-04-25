/** @enum {symbol} */
var ES = {A: Symbol('a'), B: Symbol('b')};

/** @const */
var o = {};

o[ES.A] = 1;
