function signal(val, here = goog.callerLocation()) {}
const mySignal = signal(
  () => foo() % 2 === 0,
  xid('path/to/file.ts:25'),
);
