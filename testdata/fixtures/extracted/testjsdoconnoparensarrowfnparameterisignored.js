/** @param {function(string)} callback */
function foo(callback) {}
// The `/** number */` JSDoc is attached to the entire arrow  function, not the
// parameter `x`
foo(/** number */ x =>  { var /** number */ y = x; });
