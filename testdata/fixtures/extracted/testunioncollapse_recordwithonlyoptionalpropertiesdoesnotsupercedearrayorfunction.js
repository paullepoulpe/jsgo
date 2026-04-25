/** @typedef {{x: (string|undefined)}} */
var Options;
/** @constructor */
function Foo() {}
/** @typedef {function(new: Foo)|!Array<function(new: Foo)>} */
var FooCtor;
/** @type {FooCtor|Options} */
var x;
/** @type {null} */
var y = x;
