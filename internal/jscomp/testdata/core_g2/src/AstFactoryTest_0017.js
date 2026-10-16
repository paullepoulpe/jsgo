class C {
  /** @return {number} */
  foo() {
// TODO(b/118435472): compiler should be able to infer the return type
    /**
     * @return {number}
     */
    const orig = () => 1; // new arrow function exactly like this one
  }
}
