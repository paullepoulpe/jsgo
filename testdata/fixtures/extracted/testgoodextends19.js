/** @constructor */
var MyType = function() {};
MyType.prototype = {
  /** @return {number} */
  method: function() {}
}
/** @constructor
 * @extends {MyType}
 */
var MyOtherType = function() {};
MyOtherType.prototype = {
  /** @return {number}
   * @override */
  method: function() {}
}
