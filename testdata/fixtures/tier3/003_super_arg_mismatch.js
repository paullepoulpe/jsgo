class Base {
  /** @param {number} x */
  constructor(x) {}
}
class Derived extends Base {
  constructor() { super('hi'); }
}
