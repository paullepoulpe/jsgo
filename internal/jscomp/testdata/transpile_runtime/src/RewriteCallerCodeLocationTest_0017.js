function signal(val, here = goog.callerLocation()) {}
const maybeSignal = Math.random() ? signal(0, goog.callerLocationIdInternalDoNotCallOrElse('testcode:2:36')) : null;
