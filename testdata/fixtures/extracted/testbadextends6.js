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
  /** @return {string}
   * @override */
  method: function() { return ''; }
}
