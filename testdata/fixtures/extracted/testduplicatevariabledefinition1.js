/** @record */
function A() {}
/** @type {number} */
A.prototype.prop;
/** @record */
function B() {}
/** @type {number} */
B.prototype.prop;
/** @constructor */
function C() {}
/** @type {number} */
C.prototype.prop;
/** @return {(A|B|C)} */
function fun () {}
/** @return {(B|A|C)} */
function fun () {}
