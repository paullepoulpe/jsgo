/**
  * @constructor
  * @template CLASS
  */
var Class = function() {};

/**
  * @param {function(CLASS):CLASS} a
  * @template T
  */
Class.prototype.foo = function(a) {
  return 'string';
};

/** @param {number} a
  * @return {string} */
var a = function(a) { return '' };

new Class().foo(a);
