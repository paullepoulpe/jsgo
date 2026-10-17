/** @interface */
class Named {
  /** @return {string} */
  getName() {}
}

/** @implements {Named} */
class Person {
  constructor(name) {
    /** @const */
    this.name = name;
    this.age = 0;
  }
  /** @override */
  getName() {
    return this.name;
  }
}

/** @implements {Named} */
class Pet {
  constructor(kind) {
    this.kind = kind;
    this.age = 1;
  }
  /** @override */
  getName() {
    return 'a ' + this.kind;
  }
}

class Engine {
  constructor() {
    this.age = 5;
    this.power = 100;
  }
  start() {
    return this.power;
  }
}

class Car {
  constructor() {
    this.engine = new Engine();
    this.name = 'car';
  }
  start() {
    return this.engine.start();
  }
}

/** @enum {string} */
const Color = {RED: 'red', GREEN: 'green'};

/** @typedef {{red: number, blue: number}} */
let Rgb;

/** @param {!Rgb} c @return {number} */
function brightness(c) {
  return c.red + c.blue;
}

/** @param {!Named} n */
function greet(n) {
  return 'hello ' + n.getName();
}

/** @param {*} anything */
function touch(anything) {
  return anything.age;
}

/** @param {!Person|!Pet} who */
function ageOf(who) {
  return who.age;
}

const people = [new Person('ann'), new Person('bob')];
let sum = 0;
for (const p of people) {
  sum += p.age + ageOf(p);
}
greet(new Pet('dog'));
touch(new Engine());
new Car().start();
brightness({red: 1, blue: 2});
const colorName = Color.RED + Color.GREEN;
/** @struct */
class Box {
  constructor() {
    this.value = 1;
  }
}
/** @dict */
class Bag {
  constructor() {
    this['value'] = 2;
  }
}
const both = new Box().value + new Bag()['value'] + sum + colorName.length;
