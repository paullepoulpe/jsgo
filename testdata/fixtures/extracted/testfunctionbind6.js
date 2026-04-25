/** @constructor */
function MyType() {
  /** @type {number} */
  this.x = 0;
  var f = function() {
    this.x = 'str';
  }.bind(this);
}
