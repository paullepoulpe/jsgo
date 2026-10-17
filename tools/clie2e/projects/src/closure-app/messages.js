/**
 * @fileoverview clie2e projects suite: goog.getMsg messages, translated by
 * the closure-cli project's XTB file.
 */
goog.module('clie2e.messages');

const {htmlEscape} = goog.require('goog.string');

/**
 * @param {string} name
 * @return {string}
 */
function greeting(name) {
  /** @desc Greeting shown when a user opens the shapes page. */
  const MSG_CLIE2E_GREETING = goog.getMsg('Hello, {$name}!', {name: htmlEscape(name)});
  return MSG_CLIE2E_GREETING;
}

/**
 * @param {number} count
 * @return {string}
 */
function shapeCount(count) {
  /**
   * @desc Number of shapes in the registry.
   * @meaning shape count
   */
  const MSG_CLIE2E_SHAPES = goog.getMsg('{$count} shapes registered', {count: String(count)});
  return MSG_CLIE2E_SHAPES;
}

/** @desc Label of the button that clears every shape. */
const MSG_CLIE2E_CLEAR = goog.getMsg('Clear all');

/** @desc Message without a translation in the XTB file. */
const MSG_CLIE2E_UNTRANSLATED = goog.getMsg('Nothing to see here');

exports = {greeting, shapeCount, MSG_CLIE2E_CLEAR, MSG_CLIE2E_UNTRANSLATED};
