var a = {};
a.b = {
  x: 0,
  ['logger']: goog.debug.Logger.getLogger('a.b.c')
};
a.b['logger']();
