let a$b$c = function(level) {
  if (!a$b$c$logger.isLoggable(level)) {
    a$b$c$logger.setLevel(level);
  }
  a$b$c$logger.log(level, 'hi');
};
let a$b$c$logger = goog$debug$Logger.getLogger('a.b.c');
