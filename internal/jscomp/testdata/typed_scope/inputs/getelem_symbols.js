/** @constructor */ function It() {}
/** @return {number} */
It.prototype[Symbol.iterator] = function() { return 1; };
It[Symbol.iterator];
var o = {[Symbol.iterator]: 1, ['comp']: 2};
var s1 = Symbol.iterator;
/** @const {symbol} */ var mySym = Symbol('m');
