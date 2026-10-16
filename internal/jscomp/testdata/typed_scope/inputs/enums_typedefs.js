/** @enum {string} */
var Color = {RED: 'red', GREEN: 'green'};
/** @enum */
const Num = {ONE: 1, TWO: 2};
/** @const */ var ColorAlias = Color;
/** @enum {string} */ var ColorAlias2 = Color;
/** @typedef {{a: number, b: string}} */
var Rec;
var ns = {};
/** @typedef {Array<number>} */
ns.NumArr;
/** @const */ ns.RecAlias = Rec;
/** @type {ns.NumArr} */ var arr = [1];
/** @enum {number} */ ns.E = {X: 1, [Symbol.iterator]: 2};
/** @enum {number} */ ns.E2;
/** @typedef */ var Malformed;
/** @param {Color} c @param {Rec} r */
function useIt(c, r) {}
/** @typedef {number|string} */ let LetDef;
