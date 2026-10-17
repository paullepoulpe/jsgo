function foo(val, here = goog.callerLocation()) {}
function bar(val, here = goog.callerLocation()) {
  foo(val);
}
// `bar(0)` will get rewritten because `bar` has a goog.callerLocation() default
// parameter
bar(0)
