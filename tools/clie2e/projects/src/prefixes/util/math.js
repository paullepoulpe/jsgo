/** @fileoverview clie2e projects suite: a leaf ES module. */

/**
 * @param {number} x
 * @return {number}
 */
export function square(x) {
  return x * x;
}

/**
 * @param {number} x
 * @return {number}
 */
export function cube(x) {
  return x * square(x);
}
