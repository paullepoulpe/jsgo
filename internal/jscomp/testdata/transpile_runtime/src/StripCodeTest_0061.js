let a$b$c = {
  x: null,
  get logger() {return this.x},
  set logger(a) {this.x  = goog$debug$Logger.getLogger(a)}
};
