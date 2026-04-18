/** @param {string|Array} a
@constructor */ var T = function(a) {this.a = a};
/** @return {string} */ T.prototype.f = function() {
return typeof this.a == 'string' ? this.a : 'a'; }
