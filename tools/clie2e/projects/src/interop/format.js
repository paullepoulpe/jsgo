/**
 * @fileoverview clie2e projects suite: a plain ES module.
 */
import padStart from 'lodash/padStart';

export const LABEL = 'count';

/**
 * @param {string} label
 * @param {number} n
 * @return {string}
 */
export function formatCount(label, n) {
  return label + ': ' + padStart(String(n), 6, '0');
}
