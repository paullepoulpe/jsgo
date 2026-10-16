/** @record */
function D1() { /** @type {function(D2)} */ this.a; }

/** @record */
function D2() { /** @type {D1} */ this.b; }
