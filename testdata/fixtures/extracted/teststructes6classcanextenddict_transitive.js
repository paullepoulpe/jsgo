/** @dict */
class DictParent {}
/** @unrestricted */
class UnannotatedMiddleClass extends DictParent {}
/** @struct */
class StructChild extends UnannotatedMiddleClass {
  ['prop']() {}
}
