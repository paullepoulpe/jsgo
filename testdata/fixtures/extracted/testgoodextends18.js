/** @interface */
var MyInterface = function() {};
MyInterface.prototype = {
  /** @return {number} */
  method: function() {}
}
/** @extends {MyInterface}
 * @interface */
var MyOtherInterface = function() {};
MyOtherInterface.prototype = {
  /** @return {number}
   * @override */
  method: function() {}
}
