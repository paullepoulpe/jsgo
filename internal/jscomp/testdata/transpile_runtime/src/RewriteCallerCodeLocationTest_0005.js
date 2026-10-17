function localSignal(val, here = goog.callerLocation()) {
  return val;
}
exports.signal = localSignal;
signal(0);
