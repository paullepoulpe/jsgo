/**
 * @template T, U
 * @param {T} t
 * @param {function(T): U} fn
 * @return {U}
 */
function map(t, fn) { return fn(t); }
/**
 * @template R := 'number' =:
 * @return {R}
 */
function ttl() { return 1; }
/** @template K := 'string' =: */
class TtlClass {}
/** @template T */
class Box { /** @param {T} v */ constructor(v) { this.v = v; } /** @template S @param {S} s */ m(s) { var local = s; } }
/** @interface @template T */ function Coll() {}
/** @return {T} */ Coll.prototype.first = function() {};
/** @constructor @implements {Coll<string>} */ function StrColl() {}
StrColl.prototype.first = function() { return 'a'; };
