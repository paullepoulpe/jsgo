var a = {};
a.b = function() {};
a.b.prototype = {
  ['logger']: goog.debug.Logger.getLogger('a.b.c')
}
