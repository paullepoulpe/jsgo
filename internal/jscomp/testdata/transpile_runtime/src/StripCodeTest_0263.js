let a$b$C = function() {
  this.nonStripped_ = a$b$C$nonStripped_;
}
a$b$C.prototype.method = function() {
  goog$log(this.nonStripped_, 'some message');
}
var a$b$C$nonStripped_ = goog$debug$Logger$getLogger();
