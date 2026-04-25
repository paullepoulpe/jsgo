var ns = {};
(function() {
    /** @param {number} x */ ns.foo = function(x) {}; })();
(function() { ns.foo(true); })();
