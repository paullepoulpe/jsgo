/** @constructor */
function X() {}
/** @constructor @extends {X} */
function X2() {}
/** @enum {!X} */
var XE = {A:new X};
/** @type {Object<(!XE|!X2), string>} */
var Y = {};
