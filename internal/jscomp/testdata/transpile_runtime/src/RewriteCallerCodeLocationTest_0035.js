function foo(val, here = goog.callerLocation()) {}
function bar(val) {
  foo(val, goog.callerLocationIdInternalDoNotCallOrElse('testcode:3:2'));
}
bar(0)
