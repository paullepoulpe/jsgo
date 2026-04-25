/** @constructor */ function Bindow() {}
/** @param {string} msg */
Bindow.prototype.alert = function(msg) {};
this.alert = 3;
(new Bindow()).alert(this.alert)
