const obj = {
  prop: class G {
    /** @export */ method() {}
  }
};
goog.exportProperty(obj.prop.prototype, 'method', obj.prop.prototype.method);
