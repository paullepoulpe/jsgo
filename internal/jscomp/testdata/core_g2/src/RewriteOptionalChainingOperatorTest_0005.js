let tmp0;
(tmp0 = eval) == null
    ? void 0
// The spec says that `eval?.()` must behave like an indirect
// eval, so it is important that `eval?.()` not be transpiled to
// anything that ends up containing `eval()`.
// We must be sure to call it using the temporary variable.
    : tmp0('foo()');
