/** @return {void} */
function foo(){
  /** @type {undefined|number} */
  var a;
  if (a == foo())
    return;
}
