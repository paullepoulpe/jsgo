/** @desc Hello */
var MSG_HELLO =
    goog.getMsg(
        'Hello, {$name}',
        {
          'name': getName(),
        },
        {
          html: false,
          unescapeHtmlEntities: true,
          example: {
            'name': 'George',
          },
          original_code: {
            'name': 'getName()',
          },
        })
