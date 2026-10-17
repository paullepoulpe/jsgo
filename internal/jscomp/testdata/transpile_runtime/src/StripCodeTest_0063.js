a.b.c = {
  x: null,
  get y() {return this.x},
  set y(a) {this.x  = goog.debug.Logger.getLogger(a)}
};
