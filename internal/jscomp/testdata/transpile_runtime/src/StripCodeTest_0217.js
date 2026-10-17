a = function() {}
a.prototype.logger_ = function(msg) {};
a.prototype.b = function() { this.logger_('hi'); }
