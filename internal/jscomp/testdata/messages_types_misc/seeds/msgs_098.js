/** @desc something */
const MSG_X =
    goog.getMsg(
        'Hello, {$name}',
        { 'name': getName() },
        {
          original_code: {
            'unknownPlaceholder': 'something' // not a valid placeholder name
          }
        });
