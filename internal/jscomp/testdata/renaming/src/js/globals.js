var counter = 0;
let limit = 10;
const PREFIX = 'id-';
function nextId() {
  counter++;
  return PREFIX + counter;
}
class Registry {
  constructor() {
    this.items = [];
  }
  add(item) {
    this.items.push(item);
    return nextId();
  }
}
var registry = new Registry();
function fill() {
  for (var i = 0; i < limit; i++) {
    registry.add({index: i});
  }
  helper = function() {
    return counter;
  };
}
var helper;
fill();
window.exported = helper();
if (counter > 5) {
  var late = counter * 2;
}
var a = 1, b = a + 1;
(function() {
  var local = b;
  counter += local;
})();
