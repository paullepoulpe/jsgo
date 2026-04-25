/** @constructor */function base() {
/** @type {!Number}*/
this.baseMember = new Number(4);
}
/** @constructor
  * @extends {base} */function derived() {}
/** @param {!String} x*/
function foo(x){ }
/** @type {!derived}*/var y;
foo(y.baseMember);
