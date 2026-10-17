function foo(val, here = goog.callerLocation()) {}
function bar(val) {
  foo(val);
}
// `bar(0)` will not get rewritten because `bar` does not have a goog.callerLocation()
// default parameter
bar(0)
