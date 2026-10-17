/** @desc something */
const MSG_X =
    goog.getMsg(
        'Hello, {$name}',
        {
          [name]: getName() // a computed key is not allowed
        });
