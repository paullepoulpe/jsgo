/** @desc something */
const MSG_X =
    goog.getMsg(
        'Hello, {$name}',
        { 'name': getName() },
        {
          unknownOption: true // not a valid option name
        });
