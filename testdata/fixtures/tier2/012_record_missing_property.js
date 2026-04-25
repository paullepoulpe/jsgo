/** @record */
function Shape() {}
/** @type {number} */
Shape.prototype.area;

/** @param {Shape} s */
function take(s) {
  return s.perimeter;
}
