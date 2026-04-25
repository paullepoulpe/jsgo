/** @const */ var ns = {};
/**
 * @constructor
 * @template T1
 */
ns.Foo = function() {};
/**
 * @interface
 * @template T2
 */
ns.High = function() {};
/** @type {!ns.Foo<T2>} */
ns.High.prototype.myprop;
/**
 * @interface
 * @template T3
 * @extends {ns.High<T3>}
 */
ns.Med1 = function() {};
/**
 * @interface
 * @template T4
 * @extends {ns.High<T4>}
 */
ns.Med2 = function() {};
/**
 * @interface
 * @template T5
 * @extends {ns.Med1<T5>}
 * @extends {ns.Med2<T5>}
 */
ns.Low = function() {};
