function signal(val, here = goog.callerLocation()) {}
const intermediateFunction = (() => signal(0, goog.callerLocationIdInternalDoNotCallOrElse('testcode:2:36')))();
