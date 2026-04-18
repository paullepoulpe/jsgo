/** @param {string?} a
@constructor */ var T = function(a) {this.a = a};
/** @return {string} */ T.prototype.f = function() {
return this.a ? this.a : 'a'; }
