var a = {};
a.b = function() {};
a.b.prototype = {
  logger(a) {this.x = goog.debug.Logger.getLogger('a.b.c')}
}
