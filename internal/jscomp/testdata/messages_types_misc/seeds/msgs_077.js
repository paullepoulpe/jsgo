function f() {
/** @desc Hello */ var MSG_A = goog.getMsg('hello');
/** @desc Hello */ var MSG_B = goog.getMsg('hello!');
var x = goog.getMsgWithFallback(MSG_A, MSG_B);
}
