// Leading line comment.
/* Leading block comment. */
var a = 1; // trailing line comment
var b = 2; /* trailing block comment */
/* before */ var c = 3;
function f(x /* inline param */, y) {
  // inside function
  return x + y; // after return
}
/**
 * JSDoc and a comment.
 */
// line comment after jsdoc
var d = 4;
if (a) {
  // in then
} else {
  // in else
}
var arr = [
  1, // one
  2, // two
];
f(/* arg */ 1);
/* block */ /* two blocks */ var e;
var obj = {
  // comment in object
  a: 1,
  /* before b */ b: 2,
};
switch (a) {
  // before case
  case 1: // after case
    break;
}
// final comment
