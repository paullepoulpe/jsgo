a.b.c = function(level) {
  if (!this.logger.isLoggable(level)) {
    this.logger.setLevel(level);
  }
  this.logger.log(level, 'hi');
};
a.b.c.prototype.logger = goog.debug.Logger.getLogger('a.b.c');
a.b.c.prototype.go = function() { this.logger.finer('x'); };
