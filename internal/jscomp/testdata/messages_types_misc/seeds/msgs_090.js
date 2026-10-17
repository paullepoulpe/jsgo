/** @desc something */
const MSG_X =
    goog.getMsg(
        'Hello, {$name}',
        { 'name': getName() },
        {
          [options]: true // option names cannot be computed keys
        });
