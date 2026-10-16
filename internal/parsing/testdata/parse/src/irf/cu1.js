/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('a.b');
/** @closureUnaware */
(function() {
  /** @license inner license */
  /** @type {number} */ const x = 5;
  // a line comment
  /* block */ x;
  /*! important */
}).call(globalThis);
let string1 = "ClosureAwareString";
/** @closureUnaware */
(function() {
  const y = 5;
}).apply(globalThis, []);
/** @closureUnaware */
(function() { return 1; })();
