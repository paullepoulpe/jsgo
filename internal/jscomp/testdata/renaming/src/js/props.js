/** @constructor */
function Point(x, y) {
  /** @type {number} */
  this.x = x;
  /** @type {number} */
  this.y = y;
}
Point.prototype.length = function() {
  return Math.sqrt(this.x * this.x + this.y * this.y);
};
Point.prototype.scale = function(factor) {
  return new Point(this.x * factor, this.y * factor);
};
Point.origin = new Point(0, 0);

class Shape {
  constructor(name) {
    this.name = name;
    this.sides = 0;
  }
  describe() {
    return this.name + ' with ' + this.sides + ' sides';
  }
  get label() {
    return this.name.toUpperCase();
  }
  static create(name) {
    return new Shape(name);
  }
}

class Square extends Shape {
  constructor(size) {
    super('square');
    this.sides = 4;
    this.size = size;
  }
  area() {
    return this.size * this.size;
  }
}

const config = {
  width: 10,
  height: 20,
  'quotedKey': 30,
  nested: {depth: 1, label: 'n'},
};
config.extra = config.width + config.height + config['quotedKey'];

function use(p, s) {
  const total = p.length() + s.area() + config.nested.depth;
  document.title = s.describe() + s.label + total;
  return {result: total, origin: Point.origin.x};
}
use(new Point(3, 4), new Square(2));
Shape.create('triangle').sides = 3;
