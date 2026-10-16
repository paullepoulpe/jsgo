let a$b$C = function() {
  this.nonStripped_ = a$b$C$logger_;
}
a$b$C.prototype.method = function() {
  goog$log(this.nonStripped_, 'some message');
  this.obj[a$b$C$logger_];
}
var a$b$C$logger_;
