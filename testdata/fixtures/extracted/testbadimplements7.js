/** @interface */
MyInterface = function() {};
MyInterface.prototype = {
  /** @return {number} */
  method: function() {}
}
/** @implements {MyInterface}
 * @constructor */
MyClass = function() {};
MyClass.prototype = {
  /** @return {string}
   * @override */
  method: function() { return ''; }
}
