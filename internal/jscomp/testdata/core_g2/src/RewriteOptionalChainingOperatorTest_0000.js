class TestObject {
  constructor() {
    /** @const {!TestObject} */
    this.obj = this;
    /** @const {!Array<!TestObject>} */
    this.ary = [this];
    /** @const {function(number): !TestObject} */
    this.fun = (num) => this.ary[num];
    /** @const */
    this.num = 0;
  }
  /** @return {!TestObject} */
  getObj() { return this; }
  /** @return {!Array<!TestObject>} */
  getArr() { return this.ary; }
  /** @return {function(number): !TestObject} */
  getFun() { return this.fun; }
  /** @return {number} */
  getNum() { return 0; }
}

const obj = new TestObject();
const ary = obj.ary;
const fun = obj.fun;

/** @return {!TestObject} */
function getObj() {
  return obj;
}

/** @return {!Array<!TestObject>} */
function getAry() {
  return ary;
}

/** @return {function(number): !TestObject} */
function getFun() {
  return fun;
}
