function Box(a, b) {
  this.w = a;
  this.h = b;
}
Box.prototype.area = function() {
  return this.w * this.h;
};
function totalArea(a) {
  for (var b = 0, c = 0; c < a.length; c++) {
    b += a[c].area();
  }
  return b;
}
;var shapes = [new Box(2, 3), new Box(4, 5)];
function report() {
  var a = totalArea(shapes);
  a > LIMIT && console.log("too large: " + a);
  return a;
}
window.report = report;

//# sourceMappingURL=bundle.js.map
