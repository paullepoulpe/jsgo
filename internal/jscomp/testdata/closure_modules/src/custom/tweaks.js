goog.tweak.registerBoolean('TweakA', 'desc');
goog.tweak.registerNumber('TweakB', 'desc', 5);
goog.tweak.registerString('TweakC', 'desc', 'x');
goog.tweak.registerString('Tweak D!', 'desc');
goog.tweak.registerNumber('TweakA', 'again');
goog.tweak.registerNumber('TweakE', 'desc', 'notnum');
alert(goog.tweak.getBoolean('TweakA'));
alert(goog.tweak.getNumber('TweakB'));
alert(goog.tweak.getString('TweakC'));
alert(goog.tweak.getString('TweakA'));
alert(goog.tweak.getNumber('Unregistered'));
alert(goog.tweak.getBoolean(nonLiteral));
function f() { goog.tweak.registerBoolean('Local', 'desc'); }
