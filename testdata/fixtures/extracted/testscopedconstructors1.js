function foo1() {
  /** @constructor */ function Bar() {
    /** @type {number} */ this.x = 3;
  }
}
function foo2() {
  /** @constructor */ function Bar() {
    /** @type {string} */ this.x = 'y';
  }
  /** * @param {Bar} b
   * @return {number}
   */
  function baz(b) { return b.x; }
}
