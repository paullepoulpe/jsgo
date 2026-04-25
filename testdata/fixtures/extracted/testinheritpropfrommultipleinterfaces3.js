/**
 * @interface
 * @template T1
 */
function MyCollection() {}
/**
 * @interface
 * @template T2
 * @extends {MyCollection<T2>}
 */
function MySet() {}
/**
 * @interface
 * @template T3,T4
 */
function MyMapEntry() {}
/**
 * @interface
 * @template T5,T6
 */
function MyMultimap() {}
/** @return {MyCollection<MyMapEntry<T5, T6>>} */
MyMultimap.prototype.entries = function() {};
/**
 * @interface
 * @template T7,T8
 * @extends {MyMultimap<T7, T8>}
 */
function MySetMultimap() {}
/** @return {MySet<MyMapEntry<T7, T8>>} */
MySetMultimap.prototype.entries = function() {};
/**
 * @interface
 * @template T9,T10
 * @extends {MyMultimap<T9, T10>}
 */
function MyFilteredMultimap() {}
/**
 * @interface
 * @template T11,T12
 * @extends {MyFilteredMultimap<T11, T12>}
 * @extends {MySetMultimap<T11, T12>}
 */
function MyFilteredSetMultimap() {}
