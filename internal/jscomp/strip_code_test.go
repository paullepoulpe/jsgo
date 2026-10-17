// Port of com/google/javascript/jscomp/StripCodeTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/StripCodeTest.java)

/*
 * Copyright 2007 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// StripCodeTest holds tests for StripCode.
type StripCodeTest struct {
	CompilerTestCaseBase
}

// Java: StripCodeTest.EXTERNS
const StripCodeTest_EXTERNS = ""

// Java: StripCodeTest()
func NewStripCodeTest() *StripCodeTest {
	c := &StripCodeTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, StripCodeTest_EXTERNS)
	return c
}

func TestStripCodeTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewStripCodeTest() })
}

// StripCodeTest_GetStripCodeInstance creates an instance for removing
// logging code.
//
// compiler is the Compiler. Returns a new StripCode instance.
//
// Java: StripCodeTest.getStripCodeInstance(Compiler)
func StripCodeTest_GetStripCodeInstance(compiler *Compiler) *StripCode {
	stripTypes := jutil.NewLinkedHashSet(
		"goog.debug.DebugWindow",
		"goog.debug.FancyWindow",
		"goog.debug.Formatter",
		"goog.debug.HtmlFormatter",
		"goog.debug.TextFormatter",
		"goog.debug.Logger",
		"goog.debug.LogManager",
		"goog.debug.LogRecord",
		"goog.net.BrowserChannel.LogSaver",
		"goog.log",
		"GA_GoogleDebugger")

	stripNameSuffixes := jutil.NewLinkedHashSet(
		"logger", "logger_", "debugWindow", "debugWindow_", "logFormatter_", "logBuffer_")

	stripNamePrefixes := jutil.NewLinkedHashSet("trace")

	return NewStripCode(compiler, stripTypes, stripNameSuffixes, stripNamePrefixes, false)
}

// Java: StripCodeTest.getProcessor(Compiler)
func (c *StripCodeTest) GetProcessor(compiler *Compiler) CompilerPass {
	return StripCodeTest_GetStripCodeInstance(compiler)
}

// Java: StripCodeTest.testLowerCaseSuffixFunctionCallMatch()
func (c *StripCodeTest) TestLowerCaseSuffixFunctionCallMatch() {
	// should lower-case match with "logger"
	c.Test("errorLogger();", "")
}

// Java: StripCodeTest.testChildNameOfStrippedTypeIsStripped()
func (c *StripCodeTest) TestChildNameOfStrippedTypeIsStripped() {
	c.Test(
		`goog.debug.FancyWindow.getStoredKeys_ = function() {};
goog.debug.FancyWindow.getStoredKeys_();
`,
		"")
	// It's safer to replace a variable's initial value with null instead of
	// removing the declaration.
	c.Test(
		`var goog$debug$FancyWindow$getStoredKeys_ = function() {};
goog$debug$FancyWindow$getStoredKeys_();
`,
		"var goog$debug$FancyWindow$getStoredKeys_ = null;")
}

// Java: StripCodeTest.testPartiallyCollapsedTypeIsStripped()
func (c *StripCodeTest) TestPartiallyCollapsedTypeIsStripped() {
	c.Test(
		`goog$debug.FancyWindow.getStoredKeys_ = function() {};
goog$debug.FancyWindow.getStoredKeys_();
`,
		"")
}

// Java: StripCodeTest.testLoggerDefinedInConstructor()
func (c *StripCodeTest) TestLoggerDefinedInConstructor() {
	c.Test(
		`a.b.c = function() {
  this.logger = goog.debug.Logger.getLogger('a.b.c');
};
`,
		"a.b.c=function(){}")
	c.Test(
		`let a$b$c = function() {
  this.logger = goog$debug$Logger.getLogger('a.b.c');
};
`,
		"let a$b$c=function(){}")
}

// Java: StripCodeTest.testLoggerDefinedInConstructorEs6()
func (c *StripCodeTest) TestLoggerDefinedInConstructorEs6() {
	c.Test(
		`class A {
  constructor() {
    this.logger = goog.debug.Logger.getLogger('A');
    this.otherProperty = 3;
  }
}
let a = new A;
a.logger.warning('foobar');
`,
		`class A {
  constructor() {
    this.otherProperty = 3;
  }
}
let a = new A;
`)
	c.Test(
		`class A {
  constructor() {
    this.logger = goog$debug$Logger.getLogger('A');
    this.otherProperty = 3;
  }
}
let a = new A;
a.logger.warning('foobar');
`,
		`class A {
  constructor() {
    this.otherProperty = 3;
  }
}
let a = new A;
`)
}

// Java: StripCodeTest.setUp()
func (c *StripCodeTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableNormalize()
}

// Java: StripCodeTest.testLoggerDefinedInPrototype1()
func (c *StripCodeTest) TestLoggerDefinedInPrototype1() {
	c.Test(
		`a.b.c = function() {}
a.b.c.prototype.logger = goog.debug.Logger.getLogger('a.b.c');
`,
		"a.b.c=function(){}")
	c.Test(
		`a$b$c = function() {};
a$b$c.prototype.logger = goog$debug$Logger.getLogger('a.b.c');
`,
		"a$b$c = function(){}")
}

// Java: StripCodeTest.testLoggerDefinedInPrototype2()
func (c *StripCodeTest) TestLoggerDefinedInPrototype2() {
	c.Test(
		`a.b.c = function() {};
a.b.c.prototype = {logger: goog.debug.Logger.getLogger('a.b.c')}
`,
		`a.b.c = function() {}
a.b.c.prototype = {}
`)
	c.Test(
		`let a$b$c = function() {};
a$b$c.prototype = {logger: goog$debug$Logger.getLogger('a.b.c')}
`,
		`let a$b$c = function() {};
a$b$c.prototype = {}
`)
}

// Java: StripCodeTest.testLoggerDefinedInPrototype3()
func (c *StripCodeTest) TestLoggerDefinedInPrototype3() {
	c.Test(
		`a.b.c = function() {}
a.b.c.prototype = {
  get logger() {return goog.debug.Logger.getLogger('a.b.c')}
}
`,
		`a.b.c = function() {}
a.b.c.prototype = {}
`)
	c.Test(
		`let a$b$c = function() {};
a$b$c.prototype = {
  get logger() {return goog$debug$Logger.getLogger('a.b.c')}
}
`,
		`let a$b$c = function() {};
a$b$c.prototype = {}
`)
}

// Java: StripCodeTest.testLoggerDefinedInPrototype4()
func (c *StripCodeTest) TestLoggerDefinedInPrototype4() {
	c.Test(
		`a.b.c = function() {}
a.b.c.prototype = {
  set logger(a) {this.x = goog.debug.Logger.getLogger('a.b.c')}
}
`,
		`a.b.c = function() {}
a.b.c.prototype = {}
`)
	c.Test(
		`let a$b$c = function() {};
a$b$c.prototype = {
  set logger(a) {this.x = goog$debug$Logger.getLogger('a.b.c')}
}
`,
		`let a$b$c = function() {};
a$b$c.prototype = {}
`)
}

// Java: StripCodeTest.testLoggerDefinedInPrototype5()
func (c *StripCodeTest) TestLoggerDefinedInPrototype5() {
	c.Test(
		`a.b.c = function() {}
a.b.c.prototype = {
  get f() {return this.x;},
  set f(a) {this.x = goog.debug.Logger.getLogger('a.b.c')}
}
`,
		`a.b.c = function() {}
a.b.c.prototype = {
  get f() {return this.x;},
  set f(a) {this.x = null}
}
`)
	c.Test(
		`let a$b$c = function() {}
a$b$c.prototype = {
  get f() {return this.x;},
  set f(a) {this.x = goog$debug$Logger.getLogger('a.b.c')}
}
`,
		`let a$b$c = function() {}
a$b$c.prototype = {
  get f() {return this.x;},
  set f(a) {this.x = null}
}
`)
}

// Java: StripCodeTest.testLoggerDefinedInPrototype6()
func (c *StripCodeTest) TestLoggerDefinedInPrototype6() {
	c.Test(
		`var a = {};
a.b = function() {};
a.b.prototype = {
  logger(a) {this.x = goog.debug.Logger.getLogger('a.b.c')}
}
`,
		"var a = {}; a.b = function() {}; a.b.prototype = {}")
	c.Test(
		`var a = {};
var a$b = function() {};
a$b.prototype = {
  logger(a) {this.x = goog$debug$Logger.getLogger('a.b.c')}
}
`,
		"var a = {}; var a$b = function() {}; a$b.prototype = {}")
}

// Java: StripCodeTest.testLoggerDefinedInPrototype7()
func (c *StripCodeTest) TestLoggerDefinedInPrototype7() {
	c.Test(
		`var a = {};
a.b = function() {};
a.b.prototype = {
  ['logger']: goog.debug.Logger.getLogger('a.b.c')
}
`,
		`var a = {};
a.b = function() {};
a.b.prototype = {
  ['logger']: null
}
`)
	c.Test(
		`var a = {};
var a$b = function() {};
a$b.prototype = {
  ['logger']: goog$debug$Logger.getLogger('a.b.c')
}
`,
		`var a = {};
var a$b = function() {};
a$b.prototype = {
  ['logger']: null
}
`)
}

// Java: StripCodeTest.testLoggerDefinedStatically()
func (c *StripCodeTest) TestLoggerDefinedStatically() {
	c.Test(
		`a.b.c = function() {}
a.b.c.logger = goog.debug.Logger.getLogger('a.b.c');
`,
		"a.b.c=function(){}")
	c.Test(
		`let a$b$c = function() {}
a$b$c.logger = goog$debug$Logger.getLogger('a.b.c');
`,
		"let a$b$c=function(){}")
}

// Java: StripCodeTest.testDeletedScopesAreReported()
func (c *StripCodeTest) TestDeletedScopesAreReported() {
	c.Test("var nodeLogger = function () {};", "")
}

// Java: StripCodeTest.testLoggerDefinedInObjectLiteral1()
func (c *StripCodeTest) TestLoggerDefinedInObjectLiteral1() {
	c.Test(
		`a.b.c = {
  x: 0,
  logger: goog.debug.Logger.getLogger('a.b.c')
};
`,
		"a.b.c={x:0}")
	c.Test(
		`let a$b$c = {
  x: 0,
  logger: goog$debug$Logger.getLogger('a.b.c')
};
`,
		"let a$b$c={x:0}")
}

// Java: StripCodeTest.testLoggerDefinedInObjectLiteral2()
func (c *StripCodeTest) TestLoggerDefinedInObjectLiteral2() {
	c.Test(
		`a.b.c = {
  x: 0,
  get logger() {return goog.debug.Logger.getLogger('a.b.c')}
};
`,
		"a.b.c={x:0}")
	c.Test(
		`let a$b$c = {
  x: 0,
  get logger() {return goog$debug$Logger.getLogger('a.b.c')}
};
`,
		"let a$b$c={x:0}")
}

// Java: StripCodeTest.testLoggerDefinedInObjectLiteral3()
func (c *StripCodeTest) TestLoggerDefinedInObjectLiteral3() {
	c.Test(
		`a.b.c = {
  x: null,
  get logger() {return this.x},
  set logger(a) {this.x  = goog.debug.Logger.getLogger(a)}
};
`,
		"a.b.c={x:null}")
	c.Test(
		`let a$b$c = {
  x: null,
  get logger() {return this.x},
  set logger(a) {this.x  = goog$debug$Logger.getLogger(a)}
};
`,
		"let a$b$c={x:null}")
}

// Java: StripCodeTest.testLoggerDefinedInObjectLiteral4()
func (c *StripCodeTest) TestLoggerDefinedInObjectLiteral4() {
	c.Test(
		`a.b.c = {
  x: null,
  get y() {return this.x},
  set y(a) {this.x  = goog.debug.Logger.getLogger(a)}
};
`,
		`a.b.c = {
  x: null,
  get y() {return this.x},
  set y(a) {this.x  = null}
};
`)
	c.Test(
		`let a$b$c = {
  x: null,
  get y() {return this.x},
  set y(a) {this.x  = goog$debug$Logger.getLogger(a)}
};
`,
		`let a$b$c = {
  x: null,
  get y() {return this.x},
  set y(a) {this.x  = null}
};
`)
}

// Java: StripCodeTest.testLoggerDefinedInObjectLiteral5()
func (c *StripCodeTest) TestLoggerDefinedInObjectLiteral5() {
	c.Test(
		`var a = {};
a.b = {
  x: 0,
  logger() { return goog.debug.Logger.getLogger('a.b.c'); }
};
`,
		"var a = {}; a.b={x:0}")
	c.Test(
		`var a = {};
var a$b = {
  x: 0,
  logger() { return goog$debug$Logger.getLogger('a.b.c'); }
};
`,
		"var a = {}; var a$b={x:0}")
}

// Java: StripCodeTest.testLoggerDefinedInObjectLiteral6()
func (c *StripCodeTest) TestLoggerDefinedInObjectLiteral6() {
	c.Test(
		`var a = {};
a.b = {
  x: 0,
  ['logger']: goog.debug.Logger.getLogger('a.b.c')
};
a.b['logger']();
`,
		`var a = {};
a.b = {
  x: 0,
  ['logger']: null // Don't strip computed properties
};
a.b['logger']()
`)
	c.Test(
		`var a = {};
var a$b = {
  x: 0,
  ['logger']: goog$debug$Logger.getLogger('a.b.c')
};
a$b['logger']();
`,
		`var a = {};
var a$b = {
  x: 0,
  ['logger']: null // Don't strip computed properties
};
a$b['logger']()
`)
}

// Java: StripCodeTest.testLoggerDefinedInPrototypeAndUsedInConstructor()
func (c *StripCodeTest) TestLoggerDefinedInPrototypeAndUsedInConstructor() {
	c.Test(
		`a.b.c = function(level) {
  if (!this.logger.isLoggable(level)) {
    this.logger.setLevel(level);
  }
  this.logger.log(level, 'hi');
};
a.b.c.prototype.logger = goog.debug.Logger.getLogger('a.b.c');
a.b.c.prototype.go = function() { this.logger.finer('x'); };
`,
		`a.b.c=function(level){if(!null);}
a.b.c.prototype.go=function(){}
`)
	c.Test(
		`let a$b$c = function(level) {
  if (!this.logger.isLoggable(level)) {
    this.logger.setLevel(level);
  }
  this.logger.log(level, 'hi');
};
a$b$c.prototype.logger = goog$debug$Logger.getLogger('a.b.c');
a$b$c.prototype.go = function() { this.logger.finer('x'); };
`,
		`let a$b$c=function(level){if(!null);}
a$b$c.prototype.go=function(){}
`)
}

// Java: StripCodeTest.testLoggerDefinedStaticallyAndUsedInConstructor()
func (c *StripCodeTest) TestLoggerDefinedStaticallyAndUsedInConstructor() {
	c.Test(
		`a.b.c = function(level) {
  if (!a.b.c.logger.isLoggable(level)) {
    a.b.c.logger.setLevel(level);
  }
  a.b.c.logger.log(level, 'hi');
};
a.b.c.logger = goog.debug.Logger.getLogger('a.b.c');
`,
		"a.b.c=function(level){if(!null);}")
	c.Test(
		`let a$b$c = function(level) {
  if (!a$b$c$logger.isLoggable(level)) {
    a$b$c$logger.setLevel(level);
  }
  a$b$c$logger.log(level, 'hi');
};
let a$b$c$logger = goog$debug$Logger.getLogger('a.b.c');
`,
		`let a$b$c=function(level){if(!null);}
let a$b$c$logger = null
`)
}

// Java: StripCodeTest.testLoggerVarDeclaration()
func (c *StripCodeTest) TestLoggerVarDeclaration() {
	c.Test("var logger = opt_logger || goog.debug.LogManager.getRoot();", "")
	c.Test("var logger = opt_logger || goog$debug$LogManager.getRoot();", "")
}

// Java: StripCodeTest.testLoggerLetDeclaration()
func (c *StripCodeTest) TestLoggerLetDeclaration() {
	c.Test("let logger = opt_logger || goog.debug.LogManager.getRoot();", "")
	c.Test("let logger = opt_logger || goog$debug$LogManager.getRoot();", "")
}

// Java: StripCodeTest.testLoggerConstDeclaration()
func (c *StripCodeTest) TestLoggerConstDeclaration() {
	c.Test("const logger = opt_logger || goog.debug.LogManager.getRoot();", "")
	c.Test("const logger = opt_logger || goog$debug$LogManager.getRoot();", "")
}

// Java: StripCodeTest.testLoggerDestructuringDeclaration()
func (c *StripCodeTest) TestLoggerDestructuringDeclaration() {
	// NOTE: StripCode currently does not optimize code in destructuring patterns, even if
	// it contains "strip names" or "strip types", since it's unclear what the correct optimization
	// would be or how much it would help users.

	c.TestSame("const {Logger} = goog.debug; Logger;")
	c.TestSame("const {Logger} = goog$debug; Logger;")
	c.TestSame("const {Logger: logger} = goog.debug; logger;")
	c.TestSame("const {Logger: logger} = goog$debug; logger;")
	c.TestSame("const [logger] = [1]; logger;")

	c.Test(
		"const {getLogger} = goog.debug.Logger; const logger = opt_logger || getLogger('A');",
		"const {getLogger} = goog.debug.Logger;")
	c.Test(
		"const {getLogger} = goog$debug$Logger; const logger = opt_logger || getLogger('A');",
		"const {getLogger} = goog$debug$Logger;")
}

// Java: StripCodeTest.testLoggerMethodCallByVariableType_var()
func (c *StripCodeTest) TestLoggerMethodCallByVariableType_var() {
	c.Test("var x = goog.debug.Logger.getLogger('a.b.c'); y.info(a); x.info(a);", "y.info(a)")
	c.Test("var x = goog$debug$Logger.getLogger('a.b.c'); y.info(a); x.info(a);", "y.info(a)")
}

// Java: StripCodeTest.testLoggerMethodCallByVariableType_let()
func (c *StripCodeTest) TestLoggerMethodCallByVariableType_let() {
	c.Test("let x = goog.debug.Logger.getLogger('a.b.c'); y.info(a); x.info(a);", "y.info(a)")
	c.Test("let x = goog$debug$Logger.getLogger('a.b.c'); y.info(a); x.info(a);", "y.info(a)")
}

// Java: StripCodeTest.testLoggerMethodCallByVariableType_const()
func (c *StripCodeTest) TestLoggerMethodCallByVariableType_const() {
	c.Test("const x = goog.debug.Logger.getLogger('a.b.c'); y.info(a); x.info(a);", "y.info(a)")
	c.Test("const x = goog$debug$Logger.getLogger('a.b.c'); y.info(a); x.info(a);", "y.info(a)")
}

// Java: StripCodeTest.testLoggerMethodCallByIndirection()
func (c *StripCodeTest) TestLoggerMethodCallByIndirection() {
	c.Test(
		`const x = goog.debug.Logger.getLogger('a.b.c');
(0, goog.log.error)(x, 'oh noes');
`,
		"")
}

// Java: StripCodeTest.testLoggerMethodCallByIndirection_inVarAssignment()
func (c *StripCodeTest) TestLoggerMethodCallByIndirection_inVarAssignment() {
	c.Test("const x = (0, goog.log.getLogger)('a.b.c');", "")
}

// Java: StripCodeTest.testLoggerMethodCallByIndirection_functionRemovedBeforeCalled()
func (c *StripCodeTest) TestLoggerMethodCallByIndirection_functionRemovedBeforeCalled() {
	c.Test(
		c.Srcs(
			`var goog$log$getLogger = function() {};
var module$x = (0, goog$log$getLogger)('a.b.c');
`),
		c.Expected(
			`var goog$log$getLogger = null;
var module$x = null;
`))
}

// Java: StripCodeTest.testLoggerMethodCallByIndirection_removedFunctionIsNotLastChild()
func (c *StripCodeTest) TestLoggerMethodCallByIndirection_removedFunctionIsNotLastChild() {
	c.Test(
		c.Srcs(
			`var goog$log$getLogger = function() {};
var module$x = (goog$log$getLogger, function() {})('a.b.c');
`),
		c.Expected(
			`var goog$log$getLogger = null;
var module$x = (null, function() {})('a.b.c');
`))
}

// Java: StripCodeTest.testLoggerMethodCallByIndirection_removedFunctionNotInCall()
func (c *StripCodeTest) TestLoggerMethodCallByIndirection_removedFunctionNotInCall() {
	c.Test(
		c.Srcs(
			`var goog$log$getLogger = function() {};
var module$x = (0, goog$log$getLogger);
`),
		c.Expected(
			`var goog$log$getLogger = null;
var module$x = (0, null);
`))
}

// Java: StripCodeTest.testLoggerMethodCallByIndirection_otherCommaExpressionsStayAsTheyAre()
func (c *StripCodeTest) TestLoggerMethodCallByIndirection_otherCommaExpressionsStayAsTheyAre() {
	c.Test(
		c.Srcs(
			`var anotherFunc = function() {};
var x = (0, anotherFunc);
`),
		c.Expected(
			`var anotherFunc = function() {};
var x = (0, anotherFunc);
`))
}

// Java: StripCodeTest.testSubPropertyAccessByVariableName_var()
func (c *StripCodeTest) TestSubPropertyAccessByVariableName_var() {
	c.Test(
		`var x, y = goog.debug.Logger.getLogger('a.b.c')
var logger = x;
var curlevel = logger.level_ ? logger.getLevel().name : 3;
`,
		"var x;var curlevel=null?null:3")
	c.Test(
		`var x, y = goog$debug$Logger.getLogger('a.b.c')
var logger = x;
var curlevel = logger.level_ ? logger.getLevel().name : 3;
`,
		"var x;var curlevel=null?null:3")
}

// Java: StripCodeTest.testSubPropertyAccessByVariableName_let()
func (c *StripCodeTest) TestSubPropertyAccessByVariableName_let() {
	c.Test(
		`let x, y = goog.debug.Logger.getLogger('a.b.c');
let logger = x;
let curlevel = logger.level_ ? logger.getLevel().name : 3;
`,
		"let x; let curlevel=null?null:3")
	c.Test(
		`let x, y = goog$debug$Logger.getLogger('a.b.c');
let logger = x;
let curlevel = logger.level_ ? logger.getLevel().name : 3;
`,
		"let x; let curlevel=null?null:3")
}

// Java: StripCodeTest.testSubPropertyAccessByVariableName_const()
func (c *StripCodeTest) TestSubPropertyAccessByVariableName_const() {
	c.Test(
		`const x = undefined, y = goog.debug.Logger.getLogger('a.b.c');
const logger = x;
const curlevel = logger.level_ ? logger.getLevel().name : 3;
`,
		"const x = undefined; const curlevel=null?null:3")
	c.Test(
		`const x = undefined, y = goog$debug$Logger.getLogger('a.b.c');
const logger = x;
const curlevel = logger.level_ ? logger.getLevel().name : 3;
`,
		"const x = undefined; const curlevel=null?null:3")
}

// Java: StripCodeTest.testPrefixedVariableName()
func (c *StripCodeTest) TestPrefixedVariableName() {
	c.Test(
		`this.blcLogger_ = goog.debug.Logger.getLogger('a.b.c')
this.blcLogger_.fine('Raised dirty states.');
`,
		"")
	c.Test(
		`this.blcLogger_ = goog$debug$Logger.getLogger('a.b.c')
this.blcLogger_.fine('Raised dirty states.');
`,
		"")
}

// Java: StripCodeTest.testPrefixedPropertyName()
func (c *StripCodeTest) TestPrefixedPropertyName() {
	c.Test(
		`a.b.c.staticLogger_ = goog.debug.Logger.getLogger('a.b.c')
a.b.c.staticLogger_.fine('-' + a.b.c.d_())
`,
		"")
	c.Test(
		`a$b$c$staticLogger_ = goog$debug$Logger.getLogger('a.b.c')
a$b$c$staticLogger_.fine('-' + a$b$c$d_())
`,
		"")
}

// Java: StripCodeTest.testPrefixedClassName()
func (c *StripCodeTest) TestPrefixedClassName() {
	c.Test(
		`a.b.MyLogger = function(logger) {
  this.logger_ = logger;
};
a.b.MyLogger.prototype.shout = function(msg, opt_x) {
  this.logger_.log(goog.debug.Logger.Level.SHOUT, msg, opt_x);
};
`,
		`a.b.MyLogger=function(logger){}
a.b.MyLogger.prototype.shout=function(msg,opt_x){}
`)
	c.Test(
		`a$b$MyLogger = function(logger) {
  this.logger_ = logger;
};
a$b$MyLogger.prototype.shout = function(msg, opt_x) {
  this.logger_.log(goog$debug$Logger.Level.SHOUT, msg, opt_x);
};
`,
		`a$b$MyLogger=function(logger){}
a$b$MyLogger.prototype.shout=function(msg,opt_x){}
`)
}

// Java: StripCodeTest.testLoggerClassDefinition()
func (c *StripCodeTest) TestLoggerClassDefinition() {
	c.Test("goog.debug.Logger=function(name){this.name_=name}", "")
	c.Test("goog$debug$Logger=function(name){this.name_=name}", "")
}

// Java: StripCodeTest.testStaticLoggerPropertyDefinition()
func (c *StripCodeTest) TestStaticLoggerPropertyDefinition() {
	c.Test(
		`goog.debug.Logger.Level.SHOUT=
new goog.debug.Logger.Level(x,1200)
`,
		"")
	c.Test(
		`goog$debug$Logger$Level$SHOUT=
new goog$debug$Logger$Level(x,1200)
`,
		"")
}

// Java: StripCodeTest.testStaticLoggerMethodDefinition()
func (c *StripCodeTest) TestStaticLoggerMethodDefinition() {
	c.Test(
		`goog.debug.Logger.getLogger=function(name){
return goog.debug.LogManager.getLogger(name)
};
`,
		"")
	c.Test(
		`goog$debug$Logger$getLogger=function(name){
return goog$debug$LogManager$getLogger(name)
};
`,
		"")
}

// Java: StripCodeTest.testPrototypeFieldDefinition()
func (c *StripCodeTest) TestPrototypeFieldDefinition() {
	c.Test("goog.debug.Logger.prototype.level_=null;", "")
	c.Test("goog$debug$Logger.prototype.level_=null;", "")
}

// Java: StripCodeTest.testPrototypeFieldDefinitionWithoutAssignment()
func (c *StripCodeTest) TestPrototypeFieldDefinitionWithoutAssignment() {
	c.Test("goog.debug.Logger.prototype.level_;", "")
	c.Test("goog$debug$Logger.prototype.level_;", "")
}

// Java: StripCodeTest.testPrototypeMethodDefinition()
func (c *StripCodeTest) TestPrototypeMethodDefinition() {
	c.Test(
		`goog.debug.Logger.prototype.addHandler=
function(handler){this.handlers_.push(handler)};
`,
		"")
	c.Test(
		`goog$debug$Logger.prototype.addHandler=
function(handler){this.handlers_.push(handler)};
`,
		"")
}

// Java: StripCodeTest.testPublicPropertyAssignment()
func (c *StripCodeTest) TestPublicPropertyAssignment() {
	// Eliminate property assignments on vars/properties that we
	// remove as otherwise we create invalid code.
	c.Test("goog.debug.Logger = 1; goog.debug.Logger.prop=2; ", "")
	c.Test("goog$debug$Logger = 1; goog$debug$Logger.prop=2; ", "")
	c.Test("this.blcLogger_.level=x", "")
	c.Test("goog.ui.Component.logger.prop=y", "")
	c.Test("goog$ui$Component$logger$prop=y", "")
	c.Test("goog.ui.Component.logger.prop.foo.bar=baz", "")
	c.Test("goog$ui$Component$logger$prop$foo$bar=baz", "")
}

// Java: StripCodeTest.testGlobalCallWithStrippedType()
func (c *StripCodeTest) TestGlobalCallWithStrippedType() {
	c.TestSame("window.alert(goog.debug.Logger)")
	c.TestSame("window.alert(goog$debug$Logger)")
}

// Java: StripCodeTest.testClassDefiningCallWithStripType1()
func (c *StripCodeTest) TestClassDefiningCallWithStripType1() {
	c.Test("goog.debug.Logger.inherits(Object)", "")
	c.Test("goog$debug$Logger.inherits(Object)", "")
}

// Java: StripCodeTest.testClassDefiningCallWithStripType2()
func (c *StripCodeTest) TestClassDefiningCallWithStripType2() {
	c.Test(
		`goog.formatter=function(){};
goog.inherits(goog.debug.Formatter,goog.formatter)
`,
		"goog.formatter=function(){}")
	c.Test(
		`let goog$formatter=function(){};
goog$inherits(goog$debug$Formatter,goog$formatter)
`,
		"let goog$formatter=function(){}")
}

// Java: StripCodeTest.testClassDefiningCallWithStripType3()
func (c *StripCodeTest) TestClassDefiningCallWithStripType3() {
	c.TestError(
		`goog.formatter=function(){}
goog.inherits(goog.formatter,goog.debug.Formatter)
`,
		StripCode_STRIP_TYPE_INHERIT_ERROR)
	c.TestError(
		`let goog$formatter=function(){}
goog$inherits(goog$formatter,goog$debug$Formatter)
`,
		StripCode_STRIP_TYPE_INHERIT_ERROR)
}

// Java: StripCodeTest.testClassDefiningCallWithStripType4()
func (c *StripCodeTest) TestClassDefiningCallWithStripType4() {
	c.TestSame("goog.formatter=function(){}; goog.formatter.inherits(goog.debug.FormatterFoo)")
	c.TestSame("let goog$formatter=function(){}; goog$formatter$inherits(goog$debug$FormatterFoo)")
}

// Java: StripCodeTest.testClassDefiningCallWithStripType5()
func (c *StripCodeTest) TestClassDefiningCallWithStripType5() {
	c.Test("goog.inherits(goog.debug.TextFormatter, goog.debug.Formatter)", "")
	c.Test("goog$inherits(goog$debug$TextFormatter, goog$debug$Formatter)", "")
}

// Java: StripCodeTest.testClassDefiningCallWithStripType6()
func (c *StripCodeTest) TestClassDefiningCallWithStripType6() {
	// listed types should be removed.
	c.Test("goog.debug.DebugWindow = function(){}", "")
	c.Test("goog$debug$DebugWindow = function(){}", "")
	c.Test("goog.inherits(goog.debug.DebugWindow,Base)", "")
	c.Test("goog$inherits(goog$debug$DebugWindow,Base)", "")
	c.Test("goog.debug.DebugWindow = class {}", "")
	c.Test("goog$debug$DebugWindow = class {}", "")
	c.Test("class GA_GoogleDebugger {}", "")
	c.Test("if (class GA_GoogleDebugger {}) {}", "if (null) {}")

	// types that happen to have strip types as prefix should not be
	// stripped.
	c.TestSame("goog.debug.DebugWindowFoo=function(){}")
	c.TestSame("goog$debug$DebugWindowFoo=function(){}")
	c.TestSame("goog.inherits(goog.debug.DebugWindowFoo,Base)")
	c.TestSame("goog$inherits(goog$debug$DebugWindowFoo,Base)")
	c.TestSame("goog.debug.DebugWindowFoo")
	c.TestSame("goog$debug$DebugWindowFoo")
	c.TestSame("goog.debug.DebugWindowFoo=1")
	c.TestSame("goog$debug$DebugWindowFoo=1")
	c.TestSame("goog.debug.DebugWindowFoo = class {}")
	c.TestSame("goog$debug$DebugWindowFoo = class {}")

	// qualified subtypes should be removed.
	c.Test("goog.debug.DebugWindow.Foo=function(){}", "")
	c.Test("goog$debug$DebugWindow$Foo=function(){}", "")
	c.Test("goog.inherits(goog.debug.DebugWindow.Foo,Base)", "")
	c.Test("goog$inherits(goog$debug$DebugWindow$Foo,Base)", "")
	c.Test("goog.debug.DebugWindow.Foo", "")
	c.Test("goog$debug$DebugWindow$Foo", "")
	c.Test("goog.debug.DebugWindow.Foo=1", "")
	c.Test("goog$debug$DebugWindow$Foo=1", "")
	c.Test("goog.debug.DebugWindow.Foo = class {}", "")
	c.Test("goog$debug$DebugWindow$Foo = class {}", "")
}

// Java: StripCodeTest.testStrippedES5ChildClass()
func (c *StripCodeTest) TestStrippedES5ChildClass() {
	c.Test(
		`goog.debug.Logger.Child = function() {};
goog.inherits(goog.debug.Logger.Child, goog.debug.Logger);
`,
		"")
	c.Test(
		`var goog$debug$Logger$Child = function() {};
goog$inherits(goog$debug$Logger$Child, goog$debug$Logger);
`,
		"var goog$debug$Logger$Child = null;")
}

// Java: StripCodeTest.testClassInheritanceFromStripType1()
func (c *StripCodeTest) TestClassInheritanceFromStripType1() {
	// Formatter is not a strip name or type, so cannot extend a strip type.
	c.TestError(
		"class Formatter extends goog.debug.Formatter {}", StripCode_STRIP_TYPE_INHERIT_ERROR)
	c.TestError(
		"class Formatter extends goog$debug$Formatter {}", StripCode_STRIP_TYPE_INHERIT_ERROR)
}

// Java: StripCodeTest.testClassInheritanceFromStripType2()
func (c *StripCodeTest) TestClassInheritanceFromStripType2() {
	c.TestError(
		"let Formatter = class extends goog.debug.Formatter {}",
		StripCode_STRIP_TYPE_INHERIT_ERROR)
	c.TestError(
		"let Formatter = class extends goog$debug$Formatter {}",
		StripCode_STRIP_TYPE_INHERIT_ERROR)
}

// Java: StripCodeTest.testClassInheritanceFromStripType3()
func (c *StripCodeTest) TestClassInheritanceFromStripType3() {
	// Both subclass and superclass are strip types, so this is okay.
	c.Test("goog.debug.HtmlFormatter = class extends goog.debug.Formatter {}", "")
	c.Test(
		"let goog$debug$HtmlFormatter = class extends goog$debug$Formatter {}",
		"let goog$debug$HtmlFormatter = null;")
}

// Java: StripCodeTest.testPropertyWithEmptyStringKey()
func (c *StripCodeTest) TestPropertyWithEmptyStringKey() {
	c.Test(
		"goog.format.NUMERIC_SCALES_BINARY_ = {'': 1};",
		"goog.format.NUMERIC_SCALES_BINARY_={\"\":1}")
	c.Test(
		"let goog$format$NUMERIC_SCALES_BINARY_ = {'': 1};",
		"let goog$format$NUMERIC_SCALES_BINARY_={\"\":1}")
}

// Java: StripCodeTest.testVarinIf()
func (c *StripCodeTest) TestVarinIf() {
	c.Test("if(x)var logger=null;else foo()", "if(x);else foo()")
}

// Java: StripCodeTest.testGetElemInIf()
func (c *StripCodeTest) TestGetElemInIf() {
	c.Test("var logger=null;if(x)logger[f];else foo()", "if(x);else foo()")
}

// Java: StripCodeTest.testAssignInIf()
func (c *StripCodeTest) TestAssignInIf() {
	c.Test("var logger=null;if(x)logger=1;else foo()", "if(x);else foo()")
}

// Java: StripCodeTest.testNamePrefix()
func (c *StripCodeTest) TestNamePrefix() {
	c.Test(
		`a = function(traceZZZ) {}; a.prototype.traceXXX = {x: 1}
a.prototype.z = function() { this.traceXXX.f(); };
var traceYYY = 0;
`,
		"a=function(traceZZZ){};a.prototype.z=function(){}")
}

// Java: StripCodeTest.testTypePrefix()
func (c *StripCodeTest) TestTypePrefix() {
	c.TestSame(
		`e.f.TraceXXX = function() {};
e.f.TraceXXX.prototype.yyy = 2;
`)
	c.TestSame(
		`let e$f$TraceXXX = function() {};
e$f$TraceXXX.prototype.yyy = 2;
`)
}

// Java: StripCodeTest.testStripCallsToStrippedNames1()
func (c *StripCodeTest) TestStripCallsToStrippedNames1() {
	c.Test(
		`a = function() { this.logger_ = function(msg){}; }
a.prototype.b = function() { this.logger_('hi'); }
`,
		"a=function(){};a.prototype.b=function(){}")
	c.Test(
		`a = function() {}
a.prototype.logger_ = function(msg) {};
a.prototype.b = function() { this.logger_('hi'); }
`,
		"a=function(){};a.prototype.b=function(){}")
}

// Java: StripCodeTest.testStripCallsToStrippedNames2()
func (c *StripCodeTest) TestStripCallsToStrippedNames2() {
	c.Test(
		`a = function() {}
a.prototype.logger_ = function(msg) {};
a.prototype.b = function() { this.logger_('hi'); }
`,
		"a=function(){};a.prototype.b=function(){}")
}

// Java: StripCodeTest.testStripCallsToStrippedNames3()
func (c *StripCodeTest) TestStripCallsToStrippedNames3() {
	c.Test(
		`a = function() { this.logger_ = function(msg){}; }
a.prototype.b = function() { this.logger_('hi').foo = 2; }
`,
		"a=function(){};a.prototype.b=function(){2;}")
}

// Java: StripCodeTest.testStripCallsToStrippedNames4()
func (c *StripCodeTest) TestStripCallsToStrippedNames4() {
	c.Test("a = this.logger_().foo;", "a = null;")
}

// Java: StripCodeTest.testStripVarsInitializedFromStrippedNames1()
func (c *StripCodeTest) TestStripVarsInitializedFromStrippedNames1() {
	c.Test(
		`a = function() { this.logger_ = function() { return 1; }; }
a.prototype.b = function() {
  var one = this.logger_(); if (one) foo() }
`,
		"a=function(){};a.prototype.b=function(){if(null)foo()}")
}

// Java: StripCodeTest.testStripVarsInitializedFromStrippedNames2()
func (c *StripCodeTest) TestStripVarsInitializedFromStrippedNames2() {
	c.Test(
		`a = function() { this.logger_ = function() { return 1; }; }
a.prototype.b = function() {
  var one = this.logger_.foo.bar(); if (one) foo() }
`,
		"a=function(){};a.prototype.b=function(){if(null)foo()}")
}

// Java: StripCodeTest.testReportErrorOnStripInNestedAssignment()
func (c *StripCodeTest) TestReportErrorOnStripInNestedAssignment() {
	// Strip name
	c.TestError("(foo.logger_ = 7) + 8", StripCode_STRIP_ASSIGNMENT_ERROR)

	// Strip namespaced type
	c.TestError("(goog.debug.Logger.foo = 7) + 8", StripCode_STRIP_ASSIGNMENT_ERROR)
	c.TestError("(goog$debug$Logger$foo = 7) + 8", StripCode_STRIP_ASSIGNMENT_ERROR)

	// Strip non-namespaced type
	c.TestError("(GA_GoogleDebugger.foo = 7) + 8", StripCode_STRIP_ASSIGNMENT_ERROR)
	c.TestError("(GA_GoogleDebugger$foo = 7) + 8", StripCode_STRIP_ASSIGNMENT_ERROR)
}

// Java: StripCodeTest.testNewOperator1()
func (c *StripCodeTest) TestNewOperator1() {
	c.Test(
		"function foo() {} foo.bar = new goog.debug.Logger();",
		"function foo() {} foo.bar = null;")
	c.Test(
		"function foo() {} foo.bar = new goog$debug$Logger();",
		"function foo() {} foo.bar = null;")
}

// Java: StripCodeTest.testNewOperator2()
func (c *StripCodeTest) TestNewOperator2() {
	c.Test(
		"function foo() {} foo.bar = (new goog.debug.Logger()).foo();",
		"function foo() {} foo.bar = null;")
	c.Test(
		"function foo() {} let foo$bar = (new goog$debug$Logger()).foo();",
		"function foo() {} let foo$bar = null;")
}

// Java: StripCodeTest.testNewOperator3()
func (c *StripCodeTest) TestNewOperator3() {
	c.Test("(new goog.debug.Logger()).foo().bar = 2;", "2;")
	c.Test("(new goog$debug$Logger()).foo().bar = 2;", "2;")
}

// Java: StripCodeTest.testCrazyNesting1()
func (c *StripCodeTest) TestCrazyNesting1() {
	c.Test("var x = {}; x[new goog.debug.Logger()] = 3;", "var x = {}; x[null] = 3;")
	c.Test("var x = {}; x[new goog$debug$Logger()] = 3;", "var x = {}; x[null] = 3;")
}

// Java: StripCodeTest.testCrazyNesting2()
func (c *StripCodeTest) TestCrazyNesting2() {
	c.Test("var x = {}; x[goog.debug.Logger.getLogger()] = 3;", "var x = {}; x[null] = 3;")
	c.Test("var x = {}; x[goog$debug$Logger.getLogger()] = 3;", "var x = {}; x[null] = 3;")
}

// Java: StripCodeTest.testCrazyNesting3()
func (c *StripCodeTest) TestCrazyNesting3() {
	c.Test("var x = function() {}; x(new goog.debug.Logger());", "var x = function() {}; x(null);")
	c.Test("var x = function() {}; x(new goog$debug$Logger());", "var x = function() {}; x(null);")
}

// Java: StripCodeTest.testCrazyNesting4()
func (c *StripCodeTest) TestCrazyNesting4() {
	c.Test(
		"var x = function() {}; x(goog.debug.Logger.getLogger());",
		"var x = function() {}; x(null);")
	c.Test(
		"var x = function() {}; x(goog$debug$Logger.getLogger());",
		"var x = function() {}; x(null);")
}

// Java: StripCodeTest.testCrazyNesting5()
func (c *StripCodeTest) TestCrazyNesting5() {
	c.Test(
		`var x = function() {}; var y = {};
var z = goog.debug.Logger.getLogger(); x(y[z['foo']]);
`,
		"var x = function() {}; var y = {}; x(y[null]);")
	c.Test(
		`var x = function() {}; var y = {};
var z = goog$debug$Logger.getLogger(); x(y[z['foo']]);
`,
		"var x = function() {}; var y = {}; x(y[null]);")
}

// Java: StripCodeTest.testNamespace1()
func (c *StripCodeTest) TestNamespace1() {
	c.Test("var x = {};x.traceutil = {};x.traceutil.FOO = 1;", "var x = {};")
}

// Java: StripCodeTest.testMethodCallTriggersRemoval()
func (c *StripCodeTest) TestMethodCallTriggersRemoval() {
	c.Test("this.logger_.foo.bar();", "")
}

// Java: StripCodeTest.testRemoveExpressionByName()
func (c *StripCodeTest) TestRemoveExpressionByName() {
	c.Test("this.logger_.foo.bar;", "")
}

// Java: StripCodeTest.testAliasOfRemovedVar()
func (c *StripCodeTest) TestAliasOfRemovedVar() {
	c.Test(
		"var logger_ = goog.debug.Logger.getLogger(); var alias; alias = logger_;",
		"                                             var alias; alias = null;")
	c.Test(
		"var logger_ = goog$debug$Logger.getLogger(); var alias; alias = logger_;",
		"                                             var alias; alias = null;")
}

// Java: StripCodeTest.testComplexExpression()
func (c *StripCodeTest) TestComplexExpression() {
	c.Test(
		"var logger_ = goog.debug.Logger.getLogger(); var alias; alias = (logger_ = 3) + 4;",
		"                                             var alias; alias = 3 + 4;")
	c.Test(
		"var logger_ = goog$debug$Logger.getLogger(); var alias; alias = (logger_ = 3) + 4;",
		"                                             var alias; alias = 3 + 4;")
}

// Java: StripCodeTest.testStrippedFunctionCallAssignedToNonStrippedStaticField()
func (c *StripCodeTest) TestStrippedFunctionCallAssignedToNonStrippedStaticField() {
	c.Test(
		`a.b.C = function() {
  this.nonStripped_ = a.b.C.nonStripped_;
}
a.b.C.prototype.method = function() {
  goog.log(this.nonStripped_, 'some message');
}
a.b.C.nonStripped_ = goog.debug.Logger.getLogger();
`,
		`a.b.C = function() {
  this.nonStripped_ = a.b.C.nonStripped_;
}
a.b.C.prototype.method = function() {
}
a.b.C.nonStripped_ = null;
`)
	c.Test(
		`let a$b$C = function() {
  this.nonStripped_ = a$b$C$nonStripped_;
}
a$b$C.prototype.method = function() {
  goog$log(this.nonStripped_, 'some message');
}
var a$b$C$nonStripped_ = goog$debug$Logger$getLogger();
`,
		`let a$b$C = function() {
  this.nonStripped_ = a$b$C$nonStripped_;
}
a$b$C.prototype.method = function() {
}
var a$b$C$nonStripped_ = null;
`)
}

// Java: StripCodeTest.testVarsIntroducedByCollapsePropertiesInitializedToNull()
func (c *StripCodeTest) TestVarsIntroducedByCollapsePropertiesInitializedToNull() {
	// There's known code where `a.b.strippableName` appears in the AST earlier
	// than where it is assigned a value. CollapseProperties will turn this into
	// a variable before StripCode sees it, so the reference will appear in the AST
	// earlier than the declaration.
	//
	// For such cases, we need to make sure we initialize the variable
	// to `null` instead of removing the declaration. Otherwise, we'll get an
	// error later due to the reference to a variable without a declaration.
	//
	// It would be nicer if we replaced all references to the variable with `null`
	// and removed its declaration, but doing that makes some other things we
	// want to remove harder to recognize when we visit them. (e.g. `goog$inherits(null, null)`
	// where we need to see `goog$inherits(strippable$Child, strippable$Parent)`).
	c.Test(
		`let a$b$C = function() {
  this.nonStripped_ = a$b$C$logger_;
}
a$b$C.prototype.method = function() {
  goog$log(this.nonStripped_, 'some message');
  this.obj[a$b$C$logger_];
}
var a$b$C$logger_ = {};
`,
		`let a$b$C = function() {
  this.nonStripped_ = a$b$C$logger_;
}
a$b$C.prototype.method = function() {
  this.obj[a$b$C$logger_];
}
var a$b$C$logger_ = null;
`)
	// Also handle the case where there's no initial value
	// assigned to the variable.
	c.Test(
		`let a$b$C = function() {
  this.nonStripped_ = a$b$C$logger_;
}
a$b$C.prototype.method = function() {
  goog$log(this.nonStripped_, 'some message');
  this.obj[a$b$C$logger_];
}
var a$b$C$logger_;
`,
		`let a$b$C = function() {
  this.nonStripped_ = a$b$C$logger_;
}
a$b$C.prototype.method = function() {
  this.obj[a$b$C$logger_];
}
var a$b$C$logger_ = null;
`)
}
