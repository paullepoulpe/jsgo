function signal(here = goog.callerLocation()) {}
class Button {
  method(signal) {
    signal(); // not rewritten
  }
}
signal(goog.callerLocationIdInternalDoNotCallOrElse('testcode:7:0'));
