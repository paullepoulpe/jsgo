function signal() {}
function outter() {
  function signal(here = goog.callerLocation()) {}
  signal(); // this signal() call should be rewritten
}
signal();
