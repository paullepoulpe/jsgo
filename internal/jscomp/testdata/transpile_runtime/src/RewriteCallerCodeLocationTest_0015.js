function signal(val, here = goog.callerLocation()) {}
const obj = {prop: signal(0, goog.callerLocationIdInternalDoNotCallOrElse('testcode:2:19'))};
