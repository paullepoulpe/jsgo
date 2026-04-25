/**
 * @constructor
 * @struct
 */
function Foo() {}
(true ? new Foo() : {}).x = 123;
