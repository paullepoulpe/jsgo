var Popup = {};
var OtherPopup = {};
goog.scope(function() {
  var Popup = OtherPopup;
  Popup.newMethod = function() { return new Popup(); };
});
