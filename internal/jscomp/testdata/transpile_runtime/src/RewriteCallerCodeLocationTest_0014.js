function signal(val, here = goog.callerLocation()) {}
const obj = {prop: signal(0)};
