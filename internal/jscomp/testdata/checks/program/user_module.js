goog.module('user.mod');
const legacy = goog.require('a.legacy');
const {X} = goog.require('a.mod');
const {Thing} = goog.requireType('a.legacy.ns');
const ab = goog.require('a.b');
new a.b.C();
new ab.Sub();
/** @type {!a.legacy} */ let q;
/** @type {!legacy.Inner} */ let r;
/** @param {T} t @template T */ function f(t) {}
new legacy.Inner();
use(X, Thing);
