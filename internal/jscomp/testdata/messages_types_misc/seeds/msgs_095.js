/** @desc something */
const MSG_X =
    goog.getMsg(
        'Hello, {$name}',
        { 'name': getName() },
        {
          example: 'name: something' // not an object literal
        });
