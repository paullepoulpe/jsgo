/** @constructor */ function Bindow() {}
/** @param {string} msg */
Bindow.prototype.alert = function(msg) {};
/** @return {number} */ this.alert = function() { return 3; };
(new Bindow()).alert(this.alert())
