/** @record */
class R {
  constructor() {
    /** @type {string} */
    this.str;
  }
}
/** @record */
class S {
  constructor() {
    /** @type {string} */
    this.str;
    /** @type {number} */
    this.num;
  }
}
 /**
 * @param {!R} r
 * @param {!S} s
 */
function f(r, s) {
  switch (r) {
    case s:
      return true;
  }
}
