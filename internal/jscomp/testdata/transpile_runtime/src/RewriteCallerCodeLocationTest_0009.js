function signal(val, here = goog.callerLocation()) {}
const mySignal = signal(0, goog.callerLocationIdInternalDoNotCallOrElse('testcode:2:17'));
const myComputed = computed(() => mySignal[0]() % 2 === 0);
