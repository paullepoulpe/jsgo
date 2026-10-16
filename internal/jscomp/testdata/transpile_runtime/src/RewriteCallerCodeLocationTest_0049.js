function signal() {}
function outter() {
  function signal(here = goog.callerLocation()) {}
  signal(goog.callerLocationIdInternalDoNotCallOrElse('testcode:4:2')); // rewritten
}
signal();
