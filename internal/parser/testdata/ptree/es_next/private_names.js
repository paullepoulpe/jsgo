class P {
  #a;
  #b = 1;
  static #c = 2;
  get #d() { return 1; }
  set #d(v) {}
  #e() {}
  async #f() {}
  *#g() {}
  #\u0068() {}
  m() { return this.#a + P.#c + #a in this; }
}
var x = this.#nope;
var y = #;
