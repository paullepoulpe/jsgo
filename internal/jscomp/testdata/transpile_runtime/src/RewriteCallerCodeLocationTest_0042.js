function signal(strings, here = goog.callerLocation()) {}
const tpl = signal`not passing a here argument explicitly`;
const tpl = signal`passing a here argument ${'here!'}`;
