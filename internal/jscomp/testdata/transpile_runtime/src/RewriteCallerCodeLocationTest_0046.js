function signal(here = goog.callerLocation()) {}
class Button {
  method(signal) {
    signal(); // this signal() call shouldn't be rewritten
  }
}
signal();
