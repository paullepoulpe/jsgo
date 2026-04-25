/** @return {Object} */
function f() { return {}; }
function g() { var x = f(); if (x.p) x.a = 'a'; else x.a = 'b'; }
function h() { var x = f(); x.a = false; }
