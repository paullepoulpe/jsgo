/** @const */ var ns = {};
/** @param {number} x */ ns.foo = function(x) {};
(function() {
    ns.foo = function(x) {};
    ns.foo(true);
})();
