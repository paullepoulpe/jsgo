a = function() { this.logger_ = function() { return 1; }; }
a.prototype.b = function() {
  var one = this.logger_(); if (one) foo() }
