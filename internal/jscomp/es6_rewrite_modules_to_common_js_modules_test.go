// Port of com/google/javascript/jscomp/Es6RewriteModulesToCommonJsModulesTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/Es6RewriteModulesToCommonJsModulesTest.java)

/*
 * Copyright 2018 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"testing"
)

type Es6RewriteModulesToCommonJsModulesTest struct {
	CompilerTestCaseBase

	moduleRoots        []string
	resolutionMode     ModuleLoader_ResolutionMode
	prefixReplacements *jutil.LinkedHashMap[string, string]
	pathEscaper        ModuleLoader_PathEscaper
}

func NewEs6RewriteModulesToCommonJsModulesTest() *Es6RewriteModulesToCommonJsModulesTest {
	c := &Es6RewriteModulesToCommonJsModulesTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestEs6RewriteModulesToCommonJsModulesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewEs6RewriteModulesToCommonJsModulesTest() })
}

// Java: Es6RewriteModulesToCommonJsModulesTest.setUp()
func (c *Es6RewriteModulesToCommonJsModulesTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	// ECMASCRIPT5 to trigger module processing after parsing.
	c.SetLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2015, CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.moduleRoots = []string{}
	c.resolutionMode = ModuleLoader_ResolutionMode_BROWSER
	c.prefixReplacements = esModulesTest_mapOf()
	c.pathEscaper = ModuleLoader_PathEscaper_ESCAPE
}

// Java: Es6RewriteModulesToCommonJsModulesTest.getOptions()
func (c *Es6RewriteModulesToCommonJsModulesTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	// ECMASCRIPT5 to Trigger module processing after parsing.
	options.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_ERROR)
	options.SetModuleRoots(c.moduleRoots)
	options.SetModuleResolutionMode(c.resolutionMode)
	options.SetBrowserResolverPrefixReplacements(c.prefixReplacements)
	options.SetPathEscaper(c.pathEscaper)
	return options
}

// Java: Es6RewriteModulesToCommonJsModulesTest.getProcessor(Compiler)
func (c *Es6RewriteModulesToCommonJsModulesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewEs6RewriteModulesToCommonJsModules(compiler)
}

// Java: Es6RewriteModulesToCommonJsModulesTest.testExports()
func (c *Es6RewriteModulesToCommonJsModulesTest) TestExports() {
	c.Test(
		"export var x;",
		`$jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {
  Object.defineProperties($$exports, {
    x: {
      enumerable: true,
      get: function() {
        return x;
      },
    },
  });
  var x;
}, 'testcode', []);
`)

	c.Test(
		"var x;\nexport {x}",
		`$jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {
  Object.defineProperties($$exports, {
    x: {
      enumerable: true,
      get: function() {
        return x;
      },
    },
  });
  var x;
}, 'testcode', []);
`)

	c.Test(
		"var x;\nexport {x as y}",
		`$jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {
  Object.defineProperties($$exports, {
    y: {
      enumerable: true,
      get: function() {
       return x;
      },
    },
  });
  var x;
}, 'testcode', []);
`)

	c.Test(
		"export function f() {}",
		`$jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {
  Object.defineProperties($$exports, {
    f: {
      enumerable: true,
      get: function() {
        return f;
      },
    },
  });
  function f() {}
}, 'testcode', []);
`)

	c.Test(
		"export class c {}",
		`$jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {
  Object.defineProperties($$exports, {
    c: {
      enumerable: true,
      get: function() {
        return c;
      },
    },
  });
  class c {}
}, 'testcode', []);
`)

	c.Test(
		"export default 123;",
		`$jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {
  Object.defineProperties($$exports, {
    default: {
      enumerable: true,
      get: function() {
        return $$default;
      },
    },
  });
  const $$default = 123;
}, 'testcode', []);
`)

	c.Test(
		`const x = 0;
export default x;
x++;
`,
		`$jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {
  Object.defineProperties($$exports, {
    default: {
      enumerable: true,
      get: function() {
        return $$default;
      },
    },
  });
  const x = 0;
  const $$default = x;
  x++
}, 'testcode', []);
`)

	c.Test(
		`export default function f() { return 5; }
f = () => 0;
`,
		`$jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {
  Object.defineProperties($$exports, {
    default: {
      enumerable: true,
      get: function() {
        return f;
      },
    },
  });
  function f() { return 5; }
  f = () => 0;
}, 'testcode', []);
`)

	c.Test(
		"export var x, y, z;",
		`$jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {
  Object.defineProperties($$exports, {
    x: {
      enumerable: true,
      get: function() {
        return x;
      },
    },
    y: {
      enumerable: true,
      get: function() {
        return y;
      },
    },
    z: {
      enumerable: true,
      get: function() {
        return z;
      },
    },
  });
  var x, y, z;
}, 'testcode', []);
`)

	// Exports need to be ordered with the natural string ordering.
	c.Test(
		"export var z, y, x;",
		`$jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {
  Object.defineProperties($$exports, {
    x: {
      enumerable: true,
      get: function() {
        return x;
      },
    },
    y: {
      enumerable: true,
      get: function() {
        return y;
      },
    },
    z: {
      enumerable: true,
      get: function() {
        return z;
      },
    },
  });
  var z, y, x;
}, 'testcode', []);
`)
}

// Java: Es6RewriteModulesToCommonJsModulesTest.testExportDestructureDeclaration()
func (c *Es6RewriteModulesToCommonJsModulesTest) TestExportDestructureDeclaration() {
	c.Test(
		"export let {a, c:b} = obj;",
		`$jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {
  Object.defineProperties($$exports, {
    a: {
      enumerable: true,
      get: function() {
        return a;
      },
    },
    b: {
      enumerable: true,
      get: function() {
        return b;
      },
    },
  });
  let {a, c:b} = obj;
}, 'testcode', []);
`)

	c.Test(
		"export let [a, b] = obj;",
		`$jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {
  Object.defineProperties($$exports, {
    a: {
      enumerable: true,
      get: function() {
        return a;
      },
    },
    b: {
      enumerable: true,
      get: function() {
        return b;
      },
    },
  });
  let [a, b] = obj;
}, 'testcode', []);
`)

	c.Test(
		"export let {a, b:[c,d]}  = obj;",
		`$jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {
  Object.defineProperties($$exports, {
    a: {
      enumerable: true,
      get: function() {
        return a;
      },
    },
    c: {
      enumerable: true,
      get: function() {
        return c;
      },
    },
    d: {
      enumerable: true,
      get: function() {
        return d;
      },
    },
  });
  let {a, b:[c,d]}  = obj;
}, 'testcode', []);
`)
}

// Java: Es6RewriteModulesToCommonJsModulesTest.testImport()
func (c *Es6RewriteModulesToCommonJsModulesTest) TestImport() {
	c.Test(
		"import * as x from 'other.js'; use(x, x.y);",
		`$jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {
  var x = $$require('other.js');
  use(x, x.y);
}, 'testcode', ['other.js']);
`)

	c.Test(
		"import Default, {x, y as z} from 'bogus.js'; use(x, z, Default);",
		`$jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {
  var module$bogus = $$require('bogus.js');
  use(module$bogus.x, module$bogus.y, module$bogus.default);
}, 'testcode', ['bogus.js']);
`)

	c.Test(
		"import Default, * as x from 'other.js'; use(x, x.y, Default);",
		`$jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {
  var x = $$require('other.js');
  use(x, x.y, x.default);
}, 'testcode', ['other.js']);
`)

	c.Test(
		"import First from 'other.js'; import {Second} from 'other.js'; use(First, Second);",
		`$jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {
  var module$other = $$require('other.js');
  use(module$other.default, module$other.Second);
}, 'testcode', ['other.js']);
`)

	c.Test(
		"import First from 'first.js'; import {Second} from 'second.js'; use(First, Second);",
		`$jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {
  var module$first = $$require('first.js');
  var module$second = $$require('second.js');
  use(module$first.default, module$second.Second);
}, 'testcode', ['first.js', 'second.js']);
`)
}

// Java: Es6RewriteModulesToCommonJsModulesTest.testImportAndExport()
func (c *Es6RewriteModulesToCommonJsModulesTest) TestImportAndExport() {
	c.Test(
		"export var x; import {y} from 'other.js';",
		`$jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {
  Object.defineProperties($$exports, {
    x: {
      enumerable: true,
      get: function() {
        return x;
      },
    },
  });
  var module$other = $$require('other.js');
  var x;
}, 'testcode', ['other.js']);
`)

	c.Test(
		"import {y} from 'other.js'; export var x;",
		`$jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {
  Object.defineProperties($$exports, {
    x: {
      enumerable: true,
      get: function() {
        return x;
      },
    },
  });
  var module$other = $$require('other.js');
  var x;
}, 'testcode', ['other.js']);
`)

	c.Test(
		"import {y as Y} from 'other.js'; export {Y as X};",
		`$jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {
  Object.defineProperties($$exports, {
    X: {
      enumerable: true,
      get: function() {
        return module$other.y;
      },
    },
  });
  var module$other = $$require('other.js');
}, 'testcode', ['other.js']);
`)
}

// Java: Es6RewriteModulesToCommonJsModulesTest.testExportFrom()
func (c *Es6RewriteModulesToCommonJsModulesTest) TestExportFrom() {
	c.Test(
		"export {x, y as z} from 'other.js';",
		`$jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {
  Object.defineProperties($$exports, {
    x: {
      enumerable: true,
      get: function() {
        return module$other.x;
      },
    },
    z: {
      enumerable: true,
      get: function() {
        return module$other.y;
      },
    },
  });
  var module$other = $$require('other.js');
}, 'testcode', ['other.js']);
`)
}

// Java: Es6RewriteModulesToCommonJsModulesTest.testExportWithArguments()
func (c *Es6RewriteModulesToCommonJsModulesTest) TestExportWithArguments() {
	c.Test(
		"export default function f() { return arguments[1]; }",
		`$jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {
  Object.defineProperties($$exports, {
    default: {
      enumerable: true,
      get: function() {
        return f;
      },
    },
  });
function f() { return arguments[1]; }
}, 'testcode', []);
`)
}

// Java: Es6RewriteModulesToCommonJsModulesTest.testProtocolAndDomainAreRemovedInRegisteredPathWhenNotEscaping()
func (c *Es6RewriteModulesToCommonJsModulesTest) TestProtocolAndDomainAreRemovedInRegisteredPathWhenNotEscaping() {
	c.pathEscaper = ModuleLoader_PathEscaper_CANONICALIZE_ONLY
	c.Test(
		c.Srcs(SourceFile_FromCode("https://example.domain.google.com/test.js", "export var x;")),
		c.Expected(
			SourceFile_FromCode(
				"https://example.domain.google.com/test.js",
				`$jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {
  Object.defineProperties($$exports, {
    x: {
      enumerable: true,
      get: function() {
        return x;
      },
    },
  });
  var x;
}, 'test.js', []);
`)))
}

// Java: Es6RewriteModulesToCommonJsModulesTest.testProtocolInImportPathIsError()
func (c *Es6RewriteModulesToCommonJsModulesTest) TestProtocolInImportPathIsError() {
	c.TestError("import * as foo from 'file://imported.js';", TranspilationUtil_CANNOT_CONVERT)
}

// Java: Es6RewriteModulesToCommonJsModulesTest.testRegisteredPathDoesNotIncludeModuleRoot()
func (c *Es6RewriteModulesToCommonJsModulesTest) TestRegisteredPathDoesNotIncludeModuleRoot() {
	c.moduleRoots = []string{"module/root/"}

	c.Test(
		c.Srcs(SourceFile_FromCode("module/root/test.js", "export {};")),
		c.Expected(
			SourceFile_FromCode(
				"module/root/test.js",
				`$jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {
}, 'test.js', []);
`)))
}

// Java: Es6RewriteModulesToCommonJsModulesTest.testImportPathDoesNotIncludeModuleRoot()
func (c *Es6RewriteModulesToCommonJsModulesTest) TestImportPathDoesNotIncludeModuleRoot() {
	c.moduleRoots = []string{"module/root/"}

	c.Test(
		c.Srcs(SourceFile_FromCode("not/root/test.js", "import * as foo from 'module/root/foo.js';")),
		c.Expected(
			SourceFile_FromCode(
				"not/root/test.js",
				`$jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {
  var foo = $$require('foo.js');
}, 'not/root/test.js', ['foo.js']);
`)))
}

// Java: Es6RewriteModulesToCommonJsModulesTest.testImportPathWithBrowserPrefixReplacementResolution()
func (c *Es6RewriteModulesToCommonJsModulesTest) TestImportPathWithBrowserPrefixReplacementResolution() {
	c.resolutionMode = ModuleLoader_ResolutionMode_BROWSER_WITH_TRANSFORMED_PREFIXES
	c.prefixReplacements = esModulesTest_mapOf("@root/", "")

	c.Test(
		c.Srcs(SourceFile_FromCode("not/root/test.js", "import * as foo from '@root/foo.js';")),
		c.Expected(
			SourceFile_FromCode(
				"not/root/test.js",
				`$jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {
  var foo = $$require('foo.js');
}, 'not/root/test.js', ['foo.js']);
`)))
}

// Java: Es6RewriteModulesToCommonJsModulesTest.testExportStarFrom()
func (c *Es6RewriteModulesToCommonJsModulesTest) TestExportStarFrom() {
	c.Test(
		"export * from './other.js';",
		`$jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {
  var module$other = $$require('other.js');
  $$module.exportAllFrom(module$other);
}, 'testcode', ['other.js']);
`)
}
