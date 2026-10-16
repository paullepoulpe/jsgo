class Fields {
  x = 1;
  static y = 2;
  #priv = 3;
  static { Fields.z = 4; }
  m() { return this.#priv; }
}
function NT() { return new.target; }
try { f(); } catch { g(); }
var o = {a, b};
var big = 1_000_000;
