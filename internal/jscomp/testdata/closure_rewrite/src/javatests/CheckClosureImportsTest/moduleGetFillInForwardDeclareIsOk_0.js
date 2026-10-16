goog.module('test');
let symbol = goog.forwardDeclare('symbol');
() => symbol = goog.module.get('symbol');
