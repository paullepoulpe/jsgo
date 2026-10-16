/** @constructor */ function Widget() {}
var extend = function(a, b) {};
extend(Widget.prototype, /** @lends {Widget.prototype} */ ({
  /** @return {number} */ size: function() { return 1; },
  name: 'w'
}));
extend(Widget, /** @lends {Widget} */ ({ create: function() { return new Widget(); } }));
var x = /** @lends {Unknown} */ ({a: 1});
/** @type {number} */ var num = 3;
var y = /** @lends {num} */ ({b: 2});
