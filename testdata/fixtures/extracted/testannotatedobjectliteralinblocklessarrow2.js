function f(/** {g: function(): string} */ x) {}
() => f({/** @constructor */ g: function() {}});
