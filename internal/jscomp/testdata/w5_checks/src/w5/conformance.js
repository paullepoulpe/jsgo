//// file: w5/conformance/a.js
goog.dom = {};
/** @param {string} id @return {?Element} */
goog.dom.getElement = function(id) { return null; };
/** @param {string} tag @param {?Object=} attrs @return {!Element} */
goog.dom.createDom = function(tag, attrs) { return document.createElement(tag); };
goog.string = {};
/** @constructor */ goog.string.Const = function() {};
/** @param {string} s @return {!goog.string.Const} */
goog.string.Const.from = function(s) { return new goog.string.Const(); };

eval('1');
goog.dom.getElement('x');
goog.global;
goog.isDef;
setTimeout(function() {}, 1);
alert('x');
function usesCallee() { return arguments.callee; }
window.location.href;
var href = window.location.href;
var el = document.createElement('div');
el.innerHTML = 'x';
el.innerHTML = 'x' + href;
el.outerHTML = 'x';
el.outerHTML = href;
document.write('x');
({}).hasOwnProperty('x');
parseInt('1');
parseInt('1', 10);
parseInt(1, 10);
[].slice(1);
[1].forEach(function(v) {});
document.createElement('script');
document.createElement('iframe');
document.createElement('div');
goog.dom.createDom('iframe', {'src': 'x'});
goog.dom.createDom('div', {'class': 'c'});
goog.dom.createDom('script', {'type': 't'});
el.setAttribute('src', 'x');
el.setAttribute('HREF', 'y');
el.setAttribute('title', 'z');
el.attr = 'q';
document.execCommand('insertHTML', false, 'x');
document.execCommand('bold');
var s = 'javascript:void(0)';
var d = '2024-01-01';
var t = `javascript:${s}`;
for (var k of []) {}
/** @type {?Object} */ var maybe = null;
maybe.x;
function unknownThis() { return this.x; }
/** @constructor */ function C() {}
C.prototype.m = function() { return this.unknownProp; };
/** @type {!Unresolved} */ var u;
/** @type {?Unresolved.Inner} */ var u2 = null;
var g1 = 1;
class S { static m() { return this.x; } }
function unknownFn() {}
/** @const */ var inferred = unknownFn();
goog.string.Const.from(s);
goog.string.Const.from('literal');
var str = 'abc';
str.blink;
typeof u == 'undefined';
function thrower() { throw 'str'; }
function thrower2() { throw new Error('e'); }
document.domain = 'x';
Object.prototype.hasOwnProperty.call({}, 'x');
function looseWrites(t, h) {
  t.innerHTML = h.__html;
  t.innerHTML = '';
  t.outerHTML = h.x;
  t.attr = h.y;
  t.setAttribute('src', h.z);
  t.write('x');
  return t.href + t.stack;
}
/** @return {{name: string}} */ function getRec() { return {name: 'src'}; }
el.setAttribute(getRec().name, 'x');
/** @const {string} */ var ATTR = 'src';
el.setAttribute(ATTR, 'x');
/** @const */ var ns = {};
/** @const {string} */ ns.ATTR = 'href';
el.setAttribute(ns.ATTR, 'x');

//// file: w5/conformance/banned/b.js
var bannedDep = 1;

//// file: w5/conformance/fail/c.js
eval('3');

//// file: w5/conformance/x.min.js
eval('2');

//// file: w5/conformance/legacy/d.js
eval('4');
goog.isDef;

//// file: w5/conformance/enhance.js
/**
 * @fileoverview
 * @enhance {goog.foo}
 * @mods {w5.foo_bar}
 */
var enhanced = 1;

//// file: w5/conformance/enhance2.js
/**
 * @fileoverview
 * @enhance {goog.bar}
 * @mods {w5.base}
 */
var enhanced2 = 1;

//// file: w5/conformance/strict.js
'use strict';
var strictVar = 1;
