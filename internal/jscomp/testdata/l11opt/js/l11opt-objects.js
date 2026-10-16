function useObjects() {
  var o = {x: 1, y: 'two', z: function() { return 3; }};
  var p = {};
  p.a = 1;
  use(o.x + o.y);
  var q = {'@': 5, foo: 6};
  use(q.foo);
  let r = {m: 1};
  use(r.m, r.n);
  const s = {k: 1};
  s.k = 2;
  use(s.k);
  var t = {u: 1};
  t = {u: 2, v: 3};
  use(t.u, t.v);
}
function propsDead(obj) {
  obj.a = 1;
  obj.a = 2;
  obj.b.c = 3;
  obj.b = {};
  obj.b.c = 4;
  obj.d += 1;
  obj.d = 5;
  if (obj.e) { obj.f = 1; }
  obj.f = 2;
  return obj;
}
var s1 = 'repeated string value', s2 = 'repeated string value', s3 = 'repeated string value';
var long1 = 'a very long string that is definitely more than one hundred characters long so that the LARGE mode of AliasStrings sees it';
var long2 = 'a very long string that is definitely more than one hundred characters long so that the LARGE mode of AliasStrings sees it';
var undef = 'undefined', undef2 = 'undefined';
var unicode = 'ééé 😀', unicode2 = 'ééé 😀';
