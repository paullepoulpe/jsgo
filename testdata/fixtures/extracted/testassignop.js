function fn(someUnknown) {
  var x = someUnknown;
  x *= 2; // infer 'x' to now be 'number'
  var /** null */ y = x;
}
