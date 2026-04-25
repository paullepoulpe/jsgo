class A {
  /**
   * @param {string} a
   */
  constructor(a) {
    this.a = a;
  }
}
class B extends A {
  constructor() {
    super('b');
  }
}
