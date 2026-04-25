/**
 * @template T
 * @param {function(new: T)} clazz
 */
function f(clazz) {
  /**
   * @constructor
   * @extends {clazz}
   */
  function Foo() {}
}
