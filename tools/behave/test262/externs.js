/**
 * @fileoverview The test262 host functions that tools/behave/runner.js
 * defines, so that ADVANCED neither reports them as undeclared nor renames
 * them. Passed with --externs to compat and fast builds alike.
 * @externs
 */

/**
 * @param {...*} var_args
 * @return {void}
 */
function print(var_args) {}

/** @const */
var $262 = {};

/** @type {!Object} */
$262.global;

/**
 * @param {string} src
 * @return {*}
 */
$262.evalScript = function(src) {};

/**
 * @param {!ArrayBuffer} buffer
 * @return {null}
 */
$262.detachArrayBuffer = function(buffer) {};

/** @type {*} */
$262.gc;

/** @type {*} */
$262.agent;

/** @type {*} */
$262.createRealm;

/** @type {*} */
$262.IsHTMLDDA;
