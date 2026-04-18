/** @template T */
class Base {
 /**
  * @param {(function(function(new: T, ...?)): T)=} instantiate
  * @return {!Array<T>}
  */
 delegates(instantiate = undefined) {
   return [instantiate];
 };
}
class Bar {}
class Foo {}

/** @type {Base<Foo|Bar>} */
const c = new Base();

c.delegates(ctor => new ctor(42));
