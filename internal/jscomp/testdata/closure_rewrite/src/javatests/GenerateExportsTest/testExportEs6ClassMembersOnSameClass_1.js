class G {
  /** @export */ methodA() {}
  /** @export */ methodB() {}
}
goog.exportProperty(G.prototype, 'methodB', G.prototype.methodB);
goog.exportProperty(G.prototype, 'methodA', G.prototype.methodA);
