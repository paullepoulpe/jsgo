class C {
  constructor() {
    /** @type {?number} */
    this.size = null;
  }
  method() {
// TODO(b/132283774): Redeclaring 'this.size' should cause an error.
    /** @type {number} */
    this.size = unknown;
    const /** number */ num = this.size;
  }
}
