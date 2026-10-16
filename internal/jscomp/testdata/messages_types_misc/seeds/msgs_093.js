/** @desc something */
const MSG_X =
    goog.getMsg(
        'Hello, {$name}',
        { 'name': getName() },
        {
          html: 'true' // boolean option value must be a boolean literal
        });
