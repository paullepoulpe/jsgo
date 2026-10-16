goog.provide('sc.Main');
goog.require('goog.dom');
goog.require('goog.events.EventType');

goog.scope(function() {
  var dom = goog.dom;
  var EventType = goog.events.EventType;
  const Main = sc.Main;
  let ns = goog.module.get('sc.lazy');
  var missing = goog.module.get('sc.missing');

  /** @param {!EventType} t @constructor */
  sc.Main = function(t) { this.el = dom.getElement('x'); this.t = t; };

  /** @return {!Main} */
  Main.create = function() { return new Main(EventType.CLICK); };

  function localHelper() { return dom; }
  class LocalClass {}
  sc.Main.prototype.h = localHelper;
  sc.Main.L = LocalClass;
  var bleed = function named() { return named; };
  sc.Main.b = bleed;
  /** @type {missing} */
  sc.Main.m = missing;
  {
    function blockFn() {}
  }
  for (var i = 0; i < 2; i++) {}
});
