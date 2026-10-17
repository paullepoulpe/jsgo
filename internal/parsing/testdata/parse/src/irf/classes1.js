class A {
  #x; #y = 1; static #z; #m() {} get #g() { return 1; } set #g(v) {} static get #sg() {} static set #sg(v) {}
  constructor() {} 
  m() { return this.#x + this.#y + this?.#x + (#x in this) + this.#undefinedProp; }
  n() { delete this.#x; delete this?.#y; this.#nm(); }
  get constructor() {}
  *constructor() {}
  static constructor() {}
  ['computed']() {} static ['cs']() {} get ['cg']() {} set ['cs'](v) {} ['cf'] = 1; static ['scf'];
  'quoted'() {} 1() {} static async *ag() {}
  static { var x = 1; }
}
class B { constructor() {} constructor() {} #dup; #dup; get #gs() {} get #gs() {} }
class Outer { #o; m() { class Inner { #i; f(x) { return x.#o + x.#i + x.#missing; } } } }
function outside(o) { return #x in o; }
o.#notInClass;
