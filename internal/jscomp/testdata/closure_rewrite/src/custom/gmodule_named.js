goog.module('app.util');

/** @return {number} */
function add(a, b) { return a + b; }
const PI = 3.14;
/** @enum {string} */
const Color = {RED: 'red', BLUE: 'blue'};
/** @typedef {string|number} */
let Key;
class Box { constructor() { /** @type {?Key} */ this.k = null; } }

exports = {add, PI, Color, Key, Box};
