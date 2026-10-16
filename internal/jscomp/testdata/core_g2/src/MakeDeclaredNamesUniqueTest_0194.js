function bar() {
  let a; let a$jscomp$1;
  a + a$jscomp$1;
}
function foo() {
// tests when a$jscomp$2 declared later in the same scope
  var a; a; var a$jscomp$2; a$jscomp$2;
}
