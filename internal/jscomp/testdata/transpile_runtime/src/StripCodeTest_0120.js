const x = undefined, y = goog.debug.Logger.getLogger('a.b.c');
const logger = x;
const curlevel = logger.level_ ? logger.getLevel().name : 3;
