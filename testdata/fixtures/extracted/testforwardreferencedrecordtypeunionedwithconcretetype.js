class Concrete {}

/**
 * @param {!ForwardReferenced|!Concrete} arg
 */
var bar = (arg) => {};

/**
 * @typedef {{baz}}
 */
var ForwardReferenced;

/**
 * @param {!Concrete} arg
 */
var foo = (arg) => bar(arg);
