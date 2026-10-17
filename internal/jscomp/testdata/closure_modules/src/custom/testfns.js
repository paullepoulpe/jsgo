function testFoo() {}
function setUp() {}
function tearDown() {}
function helper() {}
var testBar = function() {};
let testBaz = () => {};
const testClass = class { testA() {} helper() {} };
var MyTest = class { testB() {} };
window.testWin = function() {};
Foo.prototype.testProto = function() {};
a.b.prototype.setUpPage = function() {};
goog.testing.testSuite({ testX() {}, 'testY': function() {}, testZ: function() {} });
