/** @desc Hello */ var MSG_A = goog.getMsg('hello');
/** @desc Hello */ var MSG_B = goog.getMsg('hello!');
var x = goog.getMsgWithFallback(messages.MSG_A, MSG_B);
