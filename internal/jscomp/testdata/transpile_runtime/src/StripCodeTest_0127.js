a.b.MyLogger = function(logger) {
  this.logger_ = logger;
};
a.b.MyLogger.prototype.shout = function(msg, opt_x) {
  this.logger_.log(goog.debug.Logger.Level.SHOUT, msg, opt_x);
};
