goog.module('modB');

var {Foo:importedFoo} = goog.require('modA');

/** @type {importedFoo} */
var f = new importedFoo;
