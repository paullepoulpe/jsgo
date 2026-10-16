function f() {
/** @desc Hello */ var MSG_UNNAMED_1 = goog.getMsg('hello');
/** @desc Hello */ var MSG_UNNAMED_2 = goog.getMsg('hello');
var x = goog.getMsgWithFallback(
    MSG_UNNAMED_1, MSG_UNNAMED_2);
}
