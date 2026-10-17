/**
 * @fileoverview clie2e projects suite: a goog.module entry point with an
 * exported class, exported functions and a legacy namespace, built on a
 * goog.provide file and an ES module.
 */
goog.module('clie2e.exported');
goog.module.declareLegacyNamespace();

const esm = goog.require('clie2e.esm');
const legacy = goog.require('clie2e.legacy');
const {assert} = goog.require('goog.asserts');
const {format} = goog.require('goog.string');
const Point = goog.requireType('clie2e.legacy.Point');

/** Shapes keyed by name. */
class Registry {
  constructor() {
    /** @private @const {!Map<string, !esm.Polygon>} */
    this.shapes_ = new Map();
    /** @export {number} */
    this.version = 1;
  }

  /**
   * @param {string} name
   * @param {number} sides
   * @return {!esm.Polygon}
   * @export
   */
  add(name, sides) {
    assert(sides >= 3, 'a polygon has at least three sides');
    const shape = esm.regular(sides);
    this.shapes_.set(name, shape);
    this.version++;
    return shape;
  }

  /**
   * @param {string} name
   * @return {string}
   * @export
   */
  describe(name) {
    const shape = this.shapes_.get(name);
    if (!shape) {
      return format('%s: unknown', name);
    }
    return format('%s: %s sides, perimeter %s', name, shape.points.length,
                  shape.perimeter().toFixed(2));
  }
}

/**
 * @param {!Array<!Point>} points
 * @return {number}
 */
function pathLength(points) {
  return legacy.pathLength(points);
}
goog.exportSymbol('clie2e.pathLength', pathLength);

/** @const {!Registry} */
const registry = new Registry();
registry.add('triangle', 3);
registry.add('square', 4);

exports = {Registry, pathLength, registry};
