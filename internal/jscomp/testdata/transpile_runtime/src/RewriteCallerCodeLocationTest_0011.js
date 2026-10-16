function signal(val, here = goog.callerLocation()) {}
const [foo, setFoo] = signal(0, goog.callerLocationIdInternalDoNotCallOrElse('testcode:2:22'));
