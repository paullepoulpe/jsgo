goog.provide('ns.provided');
goog.provide('ns.provided.Sub');
goog.require('ns.other');
goog.module('ns.bad');
/** @constructor */
ns.provided.Sub = function() {};
ns.provided.value = goog.module.get('ns.provided');
/** @type {ns.provided.Sub} */
var x = new ns.provided.Sub();
