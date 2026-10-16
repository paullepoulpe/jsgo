a = function() { this.logger_ = function() { return 1; }; }
a.prototype.b = function() {
  var one = this.logger_.foo.bar(); if (one) foo() }
