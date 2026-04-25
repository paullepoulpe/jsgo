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
  /** @return {number} */
  method: function() { return 0; }
}
