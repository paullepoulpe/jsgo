var CLOSURE_TOGGLE_ORDINALS;
var CLOSURE_TOGGLE_ORDINALS = {'foo': 1, 'bar': 30, 'baz': 59};
var CLOSURE_TOGGLE_ORDINALS;
const foo = !!(goog.TOGGLES_[0] & 2);
const bar = !!(goog.TOGGLES_[1] & 1);
const baz = !!(goog.TOGGLES_[1] >> 29 & 1);
