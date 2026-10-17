goog.module('err.mod');
goog.provide('err.provided');
const {NotExported} = goog.require('err.mod');
const x = goog.require(123);
goog.module.get(foo);
goog.forwardDeclare();
exports = {[computed]: 1};
exports = 2;
function inner() {
  const y = goog.require('err.inner');
  var exports = 3;
  return exports;
}
let Fwd = goog.forwardDeclare('err.fwd');
Fwd = goog.module.get('err.fwd');
Other = goog.module.get('err.other');
