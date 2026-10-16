let tmp0;
let tmp1;
((tmp0 = getObj()) == null
    ? void 0
// Ideally we wouldn't generate a temporary to hold
// a temporary we already generated, but our logic
// is simpler if we don't worry about it.
// We will rely on optimizations to clean this up
// after transpilation.
    : (tmp1 = tmp0).getObj).call(tmp1)
