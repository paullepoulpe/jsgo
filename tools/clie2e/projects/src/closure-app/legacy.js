/**
 * @fileoverview clie2e projects suite: a goog.provide namespace with @export
 * annotations on a constructor, its prototype, its local properties and a
 * static function, plus a @define.
 */

goog.provide('clie2e.legacy');
goog.provide('clie2e.legacy.Point');

goog.require('goog.math');
goog.require('goog.string');

/** @define {boolean} Whether to round coordinates. */
clie2e.legacy.ROUND = goog.define('clie2e.legacy.ROUND', true);

/**
 * A point whose public API survives ADVANCED renaming.
 * @param {number} x
 * @param {number} y
 * @constructor
 * @export
 */
clie2e.legacy.Point = function(x, y) {
  /** @export {number} */
  this.x = clie2e.legacy.ROUND ? Math.round(x) : x;
  /** @export {number} */
  this.y = clie2e.legacy.ROUND ? Math.round(y) : y;
  /** @private {string} */
  this.label_ = goog.string.padNumber(this.x, 3) + ',' + goog.string.padNumber(this.y, 3);
};

/**
 * @param {!clie2e.legacy.Point} other
 * @return {number}
 * @export
 */
clie2e.legacy.Point.prototype.distanceTo = function(other) {
  const dx = this.x - other.x;
  const dy = this.y - other.y;
  return Math.sqrt(dx * dx + dy * dy);
};

/** @return {string} */
clie2e.legacy.Point.prototype.toString = function() {
  return '(' + this.label_ + ')';
};

/**
 * @param {number} angle In degrees.
 * @return {!clie2e.legacy.Point}
 * @export
 */
clie2e.legacy.Point.onUnitCircle = function(angle) {
  return new clie2e.legacy.Point(
      goog.math.angleDx(angle, 100), goog.math.angleDy(angle, 100));
};

/**
 * @param {!Array<!clie2e.legacy.Point>} points
 * @return {number} The length of the path through the points.
 */
clie2e.legacy.pathLength = function(points) {
  let total = 0;
  for (let i = 1; i < points.length; i++) {
    total += points[i - 1].distanceTo(points[i]);
  }
  return total;
};
