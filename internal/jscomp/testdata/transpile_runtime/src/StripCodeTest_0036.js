a.b.c = function() {}
a.b.c.prototype = {
  set logger(a) {this.x = goog.debug.Logger.getLogger('a.b.c')}
}
