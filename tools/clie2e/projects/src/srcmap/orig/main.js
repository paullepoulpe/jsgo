/**
 * @fileoverview clie2e projects suite: the second original source of
 * prebuilt/bundle.js. It reads an undeclared global, so that a compile of
 * the bundle reports a warning whose position maps back here.
 */

var shapes = [new Box(2, 3), new Box(4, 5)];

function report() {
  var area = totalArea(shapes);
  if (area > LIMIT) {
    console.log('too large: ' + area);
  }
  return area;
}

window['report'] = report;
