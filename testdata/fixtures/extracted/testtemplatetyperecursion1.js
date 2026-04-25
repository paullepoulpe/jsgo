/** @typedef {{a: D2}} */
var D1;

/** @typedef {{b: D1}} */
var D2;

fn(x);


/**
 * @param {!D1} s
 * @template T
 */
var fn = function(s) {};
