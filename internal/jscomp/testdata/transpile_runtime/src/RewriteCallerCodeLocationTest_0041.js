function foo(val1 = 0, val2 = 5, here = goog.callerLocation(), val3 = 10, val4 = 15) {}
foo()
foo(1)
foo(1, 2)
foo(1, 2, goog.callerLocationIdInternalDoNotCallOrElse('customString'))
foo(1, 2, goog.callerLocationIdInternalDoNotCallOrElse('customString'), 3)
foo(1, 2, goog.callerLocationIdInternalDoNotCallOrElse('customString'), 3, 4)
