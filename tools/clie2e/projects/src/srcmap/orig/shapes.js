/**
 * @fileoverview clie2e projects suite: the original source of
 * prebuilt/bundle.js (see README in the project file).
 */

/**
 * @param {number} w
 * @param {number} h
 * @constructor
 */
function Box(w, h) {
  this.w = w;
  this.h = h;
}

/** @return {number} */
Box.prototype.area = function() {
  return this.w * this.h;
};

/**
 * @param {!Array<!Box>} rects
 * @return {number}
 */
function totalArea(rects) {
  var sum = 0;
  for (var i = 0; i < rects.length; i++) {
    sum += rects[i].area();
  }
  return sum;
}
