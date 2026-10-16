function signal(val, here = goog.callerLocation()) {}
const maybeSignal = Math.random() ? signal(0) : null;
