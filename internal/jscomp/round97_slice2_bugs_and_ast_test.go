/*
 * Copyright 2026 The Closure Compiler Authors.
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
	"strings"
	"testing"
)

func round97CompileEveryPass(
	t *testing.T,
	code string,
	langOut CompilerOptions_LanguageMode,
	configure func(opts *CompilerOptions),
) (*Compiler, *Result) {
	t.Helper()
	compiler := NewCompiler()
	opts := NewCompilerOptions()
	opts.SetDevMode(CompilerOptions_DevMode_EVERY_PASS)
	opts.SetLanguageIn(CompilerOptions_LanguageMode_UNSUPPORTED)
	opts.SetLanguageOut(langOut)
	if configure != nil {
		configure(opts)
	}

	externs := []*SourceFile{
		SourceFile_FromCode(
			"externs.js",
			NewTestExternsBuilder().
				AddObject().
				AddArray().
				AddReflect().
				AddPromise().
				AddAsyncIterable().
				AddJSCompLibraries().
				AddExtra(`
/** @constructor @param {*=} msg */
function Error(msg) {}
/** @type {string} */ Error.prototype.message;
/** @type {string} */ Error.prototype.name;
/** @type {string} */ Error.prototype.stack;
function use(x) {}
function getObj() {}
`).
				Build(),
		),
	}
	inputs := []*SourceFile{SourceFile_FromCode("input.js", code)}
	res := compiler.Compile(externs, inputs, opts)
	return compiler, res
}

func TestRound97PrivateMemberReceiverTempConciseArrowHoistingAndScopeChange(t *testing.T) {
	t.Run("ConciseArrowBlockifiesInsteadOfLeakingTempToOuterScope", func(t *testing.T) {
		src := `
class Box {
  #val = 10;
  #inc(x) { return this.#val + x; }
  static makeHelpers() {
    const callMethod = () => getObj().#inc(1);
    const optField = () => getObj()?.#val;
    const optMethod = () => getObj()?.#inc(2);
    return [callMethod, optField, optMethod];
  }
}
use(Box.makeHelpers());
`
		for _, langOut := range []CompilerOptions_LanguageMode{
			CompilerOptions_LanguageMode_ECMASCRIPT_2020,
			CompilerOptions_LanguageMode_ECMASCRIPT5,
		} {
			compiler, res := round97CompileEveryPass(t, src, langOut, nil)
			if len(res.Errors) > 0 {
				t.Fatalf("langOut=%v unexpected errors: %v", langOut, formatJSErrors(res.Errors))
			}
			out := compiler.ToSource()
			if langOut == CompilerOptions_LanguageMode_ECMASCRIPT_2020 {
				// In ES2020 output, arrows remain arrows and must have block bodies
				// containing their own local receiver temporaries rather than hoisting
				// `let $jscomp$tmp` into `makeHelpers`!
				idxHelpers := strings.Index(out, "makeHelpers")
				idxFirstArrow := strings.Index(out, "()=>")
				if idxFirstArrow == -1 {
					idxFirstArrow = strings.Index(out, "() =>")
				}
				if idxHelpers >= 0 && idxFirstArrow > idxHelpers {
					prefix := out[idxHelpers:idxFirstArrow]
					if strings.Contains(prefix, "$jscomp$tmp") {
						t.Fatalf("receiver temp leaked outside concise arrow into outer function:\n%s", out)
					}
				}
			}
		}
	})

	t.Run("InsertedAfterHoistedInnerFunctionDeclarations", func(t *testing.T) {
		src := `
class Secret {
  #f = 41;
  #m() { return this.#f + 1; }
  static run() {
    use(getObj().#m());
    use(getObj()?.#f);
    function hoistedInner() {
      return 1;
    }
    {
      use(getObj().#m());
      function blockHoisted() {
        return 2;
      }
      use(blockHoisted());
    }
    return hoistedInner();
  }
}
use(Secret.run());
`
		for _, langOut := range []CompilerOptions_LanguageMode{
			CompilerOptions_LanguageMode_ECMASCRIPT_2020,
			CompilerOptions_LanguageMode_ECMASCRIPT5,
		} {
			_, res := round97CompileEveryPass(t, src, langOut, nil)
			if len(res.Errors) > 0 {
				t.Fatalf("langOut=%v ValidityCheck/ChangeVerifier failed: %v", langOut, formatJSErrors(res.Errors))
			}
		}
	})

	t.Run("DefaultParameterScopeChangeReported", func(t *testing.T) {
		src := `
class Cfg {
  #x = 7;
  #get() { return this.#x; }
  static withFnDefault(a = getObj().#get(), b = getObj()?.#x) {
    return a + b;
  }
  static withArrowDefault = (a = getObj().#get(), b = getObj()?.#x) => a + b;
}
use(Cfg.withFnDefault());
use(Cfg.withArrowDefault());
`
		for _, langOut := range []CompilerOptions_LanguageMode{
			CompilerOptions_LanguageMode_ECMASCRIPT_2020,
			CompilerOptions_LanguageMode_ECMASCRIPT5,
		} {
			_, res := round97CompileEveryPass(t, src, langOut, nil)
			if len(res.Errors) > 0 {
				t.Fatalf("langOut=%v ChangeVerifier/ValidityCheck failed: %v", langOut, formatJSErrors(res.Errors))
			}
		}
	})
}

func TestRound97Es6ConvertSuperConstructorCallsBugs(t *testing.T) {
	t.Run("SuperBracketAccessWithoutExtendsReportsCannotConvertYet", func(t *testing.T) {
		src := `
class C {
  m() {
    return super["x"];
  }
}
use(new C().m());
`
		// Verify Es6RewriteClass / Es6ConvertSuperConstructorCalls reports CANNOT_CONVERT_YET
		// instead of crashing with IllegalStateException: GETELEM.
		tc := NewEs6RewriteClassTest()
		tc.JUnitSetT(t)
		tc.JUnitBefore()
		tc.TestError(src, TranspilationUtil_CANNOT_CONVERT_YET)

		// Also verify full Compiler.Compile reports CANNOT_CONVERT_YET without panicking.
		_, res := round97CompileEveryPass(t, src, CompilerOptions_LanguageMode_ECMASCRIPT5, func(opts *CompilerOptions) {
			opts.SetDevMode(CompilerOptions_DevMode_OFF)
		})
		if len(res.Errors) == 0 {
			t.Fatalf("expected CANNOT_CONVERT_YET error for super[\"x\"] without extends, got none")
		}
		found := false
		for _, err := range res.Errors {
			if err.GetType() == TranspilationUtil_CANNOT_CONVERT_YET {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected TranspilationUtil_CANNOT_CONVERT_YET, got: %v", formatJSErrors(res.Errors))
		}
	})

	t.Run("LabeledNativeErrorSuperCallClimbsPastLabels", func(t *testing.T) {
		src := `
class MyErr extends Error {
  constructor(m) {
    outer: inner: while (super(m)) {
      break outer;
    }
  }
}
use(new MyErr("boom"));
`
		compiler, res := round97CompileEveryPass(t, src, CompilerOptions_LanguageMode_ECMASCRIPT5, nil)
		if len(res.Errors) > 0 {
			t.Fatalf("unexpected errors: %v", formatJSErrors(res.Errors))
		}
		out := compiler.ToSource()
		if !strings.Contains(out, "Error.call") && !strings.Contains(out, "Error.apply") {
			t.Fatalf("expected Error super call rewrite in output:\n%s", out)
		}
	})
}

func TestRound97RewriteJsonToModulePackageJsonEdgeCases(t *testing.T) {
	t.Run("ObjectValuedMainFieldIgnoredWithoutPanic", func(t *testing.T) {
		tc := NewRewriteJsonToModuleTest()
		tc.JUnitSetT(t)
		tc.JUnitBefore()

		tc.Test(
			tc.Srcs(SourceFile_FromCode("/package.json", `{"main": {"nested": "./index.js"}}`)),
			tc.Expected(`/** @fileoverview @suppress {undefinedVars} */
goog.provide('module$package_json')
var module$package_json = {"main": {"nested": "./index.js"}};
`))
		AssertThat(tc.GetLastCompiler().GetModuleLoader().GetPackageJsonMainEntries()).IsEmpty()
	})

	t.Run("ObjectValuedCustomEntrySkippedFallsBackToStringMain", func(t *testing.T) {
		compiler := NewCompiler()
		opts := NewCompilerOptions()
		opts.SetDevMode(CompilerOptions_DevMode_EVERY_PASS)
		opts.SetProcessCommonJSModules(true)
		opts.SetModuleResolutionMode(ModuleLoader_ResolutionMode_NODE)
		opts.SetPackageJsonEntryNames([]string{"exports", "browser", "main"})

		externs := []*SourceFile{SourceFile_FromCode("externs.js", "")}
		inputs := []*SourceFile{
			SourceFile_FromCode(
				"/package.json",
				`{
  "exports": {".": {"import": "./esm.js", "default": "./cjs.js"}},
  "browser": {
    "./enabled.js": true,
    "./cond.js": {"browser": "./b.js"},
    "./skip.js": false,
    "./sub.js": "./replacement.js"
  },
  "main": "./cjs.js"
}`),
		}
		res := compiler.Compile(externs, inputs, opts)
		if len(res.Errors) > 0 {
			t.Fatalf("unexpected errors: %v", formatJSErrors(res.Errors))
		}
		entries := compiler.GetModuleLoader().GetPackageJsonMainEntries()
		AssertThat(entries).ContainsExactlyEntriesIn(rewriteJsonToModuleTest_mapOf(
			"/package.json", "/./cjs.js",
			"/skip.js", ModuleLoader_JSC_BROWSER_SKIPLISTED_MARKER,
			"/sub.js", "/./replacement.js",
		))
	})
}

func TestRound97Slice2EASTViewMigrationsUnderEveryPass(t *testing.T) {
	t.Run("SuperSpreadsObjectAndArrayConstructors", func(t *testing.T) {
		src := `
class ObjSub extends Object {
  constructor(...args) {
    super(...args);
  }
}
class ArrSpread extends Array {
  constructor() {
    super(...arguments);
  }
}
class ArrMulti extends Array {
  constructor(a, b) {
    super(a, b);
  }
}
use([new ObjSub(1), new ArrSpread(1, 2), new ArrMulti(3, 4)]);
`
		_, res := round97CompileEveryPass(t, src, CompilerOptions_LanguageMode_ECMASCRIPT5, nil)
		if len(res.Errors) > 0 {
			t.Fatalf("unexpected errors: %v", formatJSErrors(res.Errors))
		}
	})

	t.Run("AsyncGeneratorsAndForAwaitOf", func(t *testing.T) {
		src := `
async function* gen(src) {
  for await (const {x} of src) {
    yield await x;
    yield;
  }
  yield* src;
  return 42;
}
use(gen);
`
		_, res := round97CompileEveryPass(t, src, CompilerOptions_LanguageMode_ECMASCRIPT_2017, nil)
		if len(res.Errors) > 0 {
			t.Fatalf("unexpected errors: %v", formatJSErrors(res.Errors))
		}
	})

	t.Run("InstrumentAsyncContextWithComplexParamsAndGenerators", func(t *testing.T) {
		src := `
class Base {
  step(v) { return v; }
}
class Runner extends Base {
  *genMethod(a = 1, {b} = {b: 2}, ...rest) {
    const x = yield super.step(a + b + rest.length);
    return x;
  }
}
async function runAsync(first = 10, ...more) {
  return (await first) + arguments.length + more.length;
}
use([new Runner(), runAsync]);
`
		_, res := round97CompileEveryPass(
			t,
			src,
			CompilerOptions_LanguageMode_ECMASCRIPT_2020,
			func(opts *CompilerOptions) {
				opts.SetInstrumentAsyncContext(true)
			},
		)
		if len(res.Errors) > 0 {
			t.Fatalf("unexpected errors: %v", formatJSErrors(res.Errors))
		}
	})

	t.Run("TemplateAndTaggedTemplateLiterals", func(t *testing.T) {
		src := "function tag(strings, ...vals) { return strings.raw[0] + vals[0]; }\n" +
			"const a = `${getObj()} suffix`;\n" +
			"const b = `prefix ${getObj()}`;\n" +
			"const c = tag`cooked\\n${a}`;\n" +
			"const d = tag`same${b}`;\n" +
			"use([a, b, c, d]);\n"
		_, res := round97CompileEveryPass(t, src, CompilerOptions_LanguageMode_ECMASCRIPT5, nil)
		if len(res.Errors) > 0 {
			t.Fatalf("unexpected errors: %v", formatJSErrors(res.Errors))
		}
	})
}
