/** @fileoverview @closureUnaware */
/** @closureUnaware */ (function() {
  /** @closureUnaware */ (function() { var a = `x${1}`; class C { #p = 1; m() { return this.#p; } } })();
  async function g() { await 1; for await (const z of []) {} }
  function* h() { yield 1; }
  label: for (;;) { break label; }
}).call(this);
/** @license outer */
var after = 1;
