/** @record */
function Shape() {}
/** @type {number} */
Shape.prototype.area;

/** @constructor */
function Square() {}
/** @type {number} */
Square.prototype.area = 1;

/** @param {Shape} s */
function take(s) {}
take(new Square());
