function f(arr) {
  const fns = [];
  for (let i = 0; i < arr.length; i++) {
    const v = arr[i];
    fns.push(function() { return v + i; });
  }
  for (const k in arr) {
    fns.push(() => k);
  }
  let x, y = 2;
  while (y--) {
    let z;
    fns.push({get p() { return z; }});
  }
  label: for (let j = 0; j < 2; j++) { setTimeout(function() { return j; }); continue label; }
  {
    function inner() { return x; }
    inner();
  }
  if (x) { function g() {} g(); }
  return fns;
}
const c = 1;
let d;
{ let e = c; d = e; }
