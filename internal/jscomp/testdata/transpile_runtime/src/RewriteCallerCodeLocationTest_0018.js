function signal(val, here = goog.callerLocation()) {}
const intermediateFunction = (() => signal(0))();
