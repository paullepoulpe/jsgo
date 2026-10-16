/** @desc something */
const MSG_X =
    goog.getMsg(
        'Hello, {$name}',
        { 'name': getName() },
        {
          original_code: {
            'name': getName() // value is not a string
          }
        });
