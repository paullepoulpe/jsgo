goog.module('user.mod');
const myMod = goog.require('my.mod');
const {fn, Klass} = goog.require('other.dep');
/** @type {myMod.Impl} */ let impl;
class Sub extends goog.module.get('other.dep').Klass {}
/** @param {!Klass} k */ exports.useIt = function(k) {};
