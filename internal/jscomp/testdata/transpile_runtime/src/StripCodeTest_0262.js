a.b.C = function() {
  this.nonStripped_ = a.b.C.nonStripped_;
}
a.b.C.prototype.method = function() {
}
a.b.C.nonStripped_ = null;
