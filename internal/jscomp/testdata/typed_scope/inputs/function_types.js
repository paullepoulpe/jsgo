/** @type {function(number): string} */
var typed = function(n) { return ''; };
/** @type {function(new:Object)} */
var structCtor = function() {};
/** @param {number} a @return {number} */
var stub;
/** @return {string} */
var fnLit = function(/** number */ x) { return ''; };
var local = function() {};
function decl(/** string */ s) /** number */ { return 1; }
var obj = {
  /** @param {string} s */ m(s) {},
  n: function() {},
  /** @type {number} */ p: 1,
  get g() { return 1; },
  set g(v) {}
};
/** @this {Array} */ function withThis() {}
/** @abstract */ function abs() {}
/** @closurePrimitive {asserts.truthy} */ function assertIt(x) {}
