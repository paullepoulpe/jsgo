function foo(val, here = goog.callerLocation()) {}
function bar(val, here = goog.callerLocation()) {
  foo(val, goog.callerLocationIdInternalDoNotCallOrElse('testcode:3:2'));
}
bar(0, goog.callerLocationIdInternalDoNotCallOrElse('testcode:7:0'));
