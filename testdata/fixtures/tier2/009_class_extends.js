class Base {
  constructor() {}
}
class Derived extends Base {
  constructor() { super(); }
}
/** @param {Base} b */
function take(b) {}
take(new Derived());
