use(hoisted);
function hoisted() { return later + hoisted; }
var later = 1;
later = 2;
later += 3;
later++;
--later;
var once = {};
var never;
never;
let destr;
({destr} = {destr: 1});
[destr] = [2];
var {p1, p2: {p3}, ...p4} = {};
var [q1, , q2 = q1, ...q5] = [];
for (var fi in {}) { fi; }
for (const fo of []) { fo; }
for (let li = 0; li < 3; li++) { li; }
function params(a, b = a, {c} = {}, [d] = [], ...e) { a = 1; b++; c; d; e; }
var fexpr = function named() { return named; };
var arrow = (x) => x + arguments;
var cond = later ? once : never;
var and = (later && once || never) ?? fexpr;
var opt = once?.p?.[later]?.(never);
if (later) { var inIf = 1; } else { inIf = 2; }
while (later) { later--; }
do { later++; } while (later < 5);
switch (later) { case once: later; case never: break; }
try { later; } catch (err) { err; } finally { later; }
class K extends once { [later] = never; m() { return later; } }
label: for (;;) { break label; }
function nested() {
  function inner() { return innerVar; }
  var innerVar = inner();
  return function() { innerVar = 1; return innerVar; };
}
var weak = goog.weakUsage(once);
delete later.x;
typeof undeclared;
undeclared2 = 1;
