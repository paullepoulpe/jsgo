/**
 * @template T
 * @interface
 */
class IFoo { }

/** @implements {IFoo<number>} */
class FooA { }

/** @implements {IFoo<string>} */
class FooB extends FooA { }

const /** !IFoo<number> */ x = new FooB();
const /** !IFoo<string> */ y = new FooB();
