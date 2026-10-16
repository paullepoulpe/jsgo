function bar() {
  let a; let a$jscomp$1; a + a$jscomp$1;
}
function foo() {
// tests when a$jscomp$2 declared first in the same scope
  var a$jscomp$2; a$jscomp$2; var a; a;
}
