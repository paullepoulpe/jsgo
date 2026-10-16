goog.provide('a.b.c');
goog.provide('a.b.d');
goog.provide('a.x');
goog.provide('CONST_NAME');
goog.provide('simple');
goog.require('a.b.c');
goog.requireType('a.x');
goog.forwardDeclare('fwd.Decl');
a.b.c = function() {};
/** @typedef {string} */ a.b.d;
simple = 5;
/** @const */ var CONST_NAME = {};
/** @provideAlreadyProvided */ goog.provide('a.already');
function a() {}
class simple2 {}
