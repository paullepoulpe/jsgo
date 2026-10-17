/**
 * @fileoverview clie2e projects suite: an ES module that declares a Closure
 * module id, imports Closure namespaces with goog: paths, and is required
 * by a goog.module.
 */

goog.declareModuleId('clie2e.esm');

import googArray from 'goog:goog.array';
import Point from 'goog:clie2e.legacy.Point';

/** A polygon over legacy points. */
export class Polygon {
  /** @param {!Array<!Point>} points */
  constructor(points) {
    /** @const {!Array<!Point>} */
    this.points = points;
  }

  /**
   * @return {number}
   * @export
   */
  perimeter() {
    const closed = googArray.concat(this.points, [this.points[0]]);
    let total = 0;
    for (let i = 1; i < closed.length; i++) {
      total += closed[i - 1].distanceTo(closed[i]);
    }
    return total;
  }
}

/**
 * @param {number} sides
 * @return {!Polygon}
 */
export function regular(sides) {
  const points = [];
  for (let i = 0; i < sides; i++) {
    points.push(Point.onUnitCircle(360 * i / sides));
  }
  return new Polygon(points);
}
