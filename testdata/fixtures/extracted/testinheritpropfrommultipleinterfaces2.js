/** @interface */
function High1() {}
/** @type {number} */
High1.prototype.prop;
/** @interface */
function High2() {}
/** @type {number|string} */
High2.prototype.prop;
/**
 * @interface
 * @extends {High1}
 * @extends {High2}
 */
function Low() {}
function f(/** !Low */ x) { var /** null */ n = x.prop; }
