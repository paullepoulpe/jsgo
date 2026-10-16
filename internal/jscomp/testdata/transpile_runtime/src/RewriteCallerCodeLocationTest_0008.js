function signal(val, here = goog.callerLocation()) {}
const mySignal = signal(0); // we will rewrite this `signal(0)` call
// mySignal will not be rewritten because we will not rewrite aliased functions
const myComputed = computed(() => mySignal[0]() % 2 === 0);
