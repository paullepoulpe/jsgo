// Window is not marked as a constructor, so the
// inheritance doesn't happen.
function Window() {}
Window.prototype.alert = function() {};
this.alert();
