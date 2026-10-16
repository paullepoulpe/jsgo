let a$b$c = function() {};
a$b$c.prototype = {
  get logger() {return goog$debug$Logger.getLogger('a.b.c')}
}
