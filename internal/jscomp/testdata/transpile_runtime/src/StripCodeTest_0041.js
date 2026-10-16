let a$b$c = function() {}
a$b$c.prototype = {
  get f() {return this.x;},
  set f(a) {this.x = null}
}
