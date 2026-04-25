/**\s
 * @constructor\s
 * @struct\s
 */\s
function Square(side) {\s
  this.side = side;\s
}\s
Square.prototype = /** @struct */ {
  area: function() { return this.side * this.side; }
};
Square.prototype.id = function(x) { return x; };
