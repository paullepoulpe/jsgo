/** @interface */ function T() {};
/** @type {number} */T.prototype.x;
/** @constructor\s
 * @implements {T}\s
 */
function C() {}
/** @override */
C.prototype.x = 'foo';
