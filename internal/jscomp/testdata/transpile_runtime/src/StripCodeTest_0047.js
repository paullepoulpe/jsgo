var a = {};
a.b = function() {};
a.b.prototype = {
  ['logger']: null
}
