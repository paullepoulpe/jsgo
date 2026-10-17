/**
 * @desc A message with lots of stuff.
 */
var MSG_HUGE = goog.getMsg(
    '{$startLink_1}Google{$endLink}' +
    '{$startLink_2}blah{$endLink}{$boo}{$foo_001}{$boo}' +
    '{$foo_002}{$xxx_001}{$image}{$image_001}{$xxx_002}',
    {'startLink_1': '<a href=http://www.google.com/>',
     'endLink': '</a>',
     'startLink_2': '<a href="' + opt_data.url + '">',
     'boo': opt_data.boo,
     'foo_001': opt_data.foo,
     'foo_002': opt_data.boo.foo,
     'xxx_001': opt_data.boo + opt_data.foo,
     'image': htmlTag7,
     'image_001': opt_data.image,
     'xxx_002': foo.callWithOnlyTopLevelKeys(
         bogusFn, opt_data, null, 'bogusKey1',
         opt_data.moo, 'bogusKey2', param10)});
