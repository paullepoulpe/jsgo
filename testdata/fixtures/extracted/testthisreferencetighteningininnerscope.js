class C {
  constructor() {
    /** @type {?number} */
    this.size = 0;
// Verify that the inferred type of `this.size`, which is `number`, propagates into
// block scopes.
    { const /** number */ n = this.size; }
  }
}
