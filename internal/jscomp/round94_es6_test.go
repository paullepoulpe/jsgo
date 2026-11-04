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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type round94ClassNormalizationCase struct {
	CompilerTestCaseBase
	languageOut                      CompilerOptions_LanguageMode
	enableNormalize                  bool
	skipNonTranspilation             bool
	assumeStaticInheritanceIsNotUsed bool
}

func newRound94ClassNormalizationCase() *round94ClassNormalizationCase {
	c := &round94ClassNormalizationCase{
		languageOut:     CompilerOptions_LanguageMode_ECMASCRIPT_2020,
		enableNormalize: true,
	}
	c.InitCompilerTestCaseBaseWithExterns(
		c,
		NewTestExternsBuilder().
			AddObject().
			AddReflect().
			AddJSCompLibraries().
			AddExtra("function use(x) {}\n").
			Build(),
	)
	return c
}

func (c *round94ClassNormalizationCase) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	if c.enableNormalize {
		c.EnableNormalize()
	} else {
		c.DisableNormalize()
	}
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_UNSUPPORTED)
	c.SetLanguageOut(c.languageOut)
	c.SetGenericNameReplacements(Es6NormalizeClasses_GENERIC_NAME_REPLACEMENTS)
}

func (c *round94ClassNormalizationCase) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetSkipNonTranspilationPasses(c.skipNonTranspilation)
	options.SetAssumeStaticInheritanceIsNotUsed(c.assumeStaticInheritanceIsNotUsed)
	return options
}

func (c *round94ClassNormalizationCase) GetProcessor(compiler *Compiler) CompilerPass {
	return NewEs6NormalizeClasses(compiler)
}

func TestRound94_Es6NormalizeClasses(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return newRound94ClassNormalizationCase() })
}

// Verifies that `super(this)` in a derived class with instance fields only
// records the `super(...)` call once rather than twice (once for `super` and
// once for the `this` argument).
func (c *round94ClassNormalizationCase) TestSuperCallWithThisArgumentDoesNotDuplicateSuperCalls() {
	c.skipNonTranspilation = true
	c.Test(
		`class Base {
  constructor(arg) {}
}
class Sub extends Base {
  x = 1;
  constructor() {
    super(this);
  }
}
`,
		`class Base {
  constructor(arg) {}
}
class Sub extends Base {
  constructor() {
    super(this);
    this.x = 1;
  }
}
`,
	)
}

// Verifies that an inner class whose `extends` clause references `this` or
// `super` of an enclosing class's static block, static field, or static method
// resolves against the enclosing class record instead of panicking with
// "wrong class node".
func (c *round94ClassNormalizationCase) TestInnerClassExtendsThisOrSuperInOuterStaticContext() {
	c.Test(
		`class Base {}
class Outer {
  static Base = Base;
  static {
    class Inner extends this.Base {}
    use(Inner);
  }
}
`,
		`class Base {}
class Outer {}
Outer.Base = Base;
{
  class Inner extends Outer.Base {}
  use(Inner);
}
`,
	)

	c.Test(
		`class Base {}
class Outer {
  static Base = Base;
  static Sub = class extends this.Base {};
}
`,
		`class Base {}
class Outer {}
Outer.Base = Base;
Outer.Sub = (() => {
  const CLASS_DECL$0 = class extends Outer.Base {};
  return CLASS_DECL$0;
})();
`,
	)

	c.Test(
		`class Outer {
  static f() { return 42; }
  static {
    class Inner extends this {
      static #p() {
        return super.f();
      }
      static m() { return Inner.#p(); }
    }
    use(Inner.m());
  }
}
`,
		`class Outer {
  static f() {
    return 42;
  }
}
{
  const STATIC_PRIVATE_MAP$0 = new $jscomp.PrivateMap();
  class Inner extends Outer {
    static m() {
      return STATIC_PRIVATE_MAP$0.get(Inner).p.call(Inner);
    }
    static STATIC_INIT$2() {
      const PRIVATE$1 = Object.create(null);
      STATIC_PRIVATE_MAP$0.set(Inner, PRIVATE$1);
      PRIVATE$1.p = function() {
        return Outer.f.call(this);
      };
    }
  }
  Inner.STATIC_INIT$2();
  use(Inner.m());
}
`,
	)

	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT_2021
	c.assumeStaticInheritanceIsNotUsed = true
	c.SetLanguageOut(c.languageOut)
	c.Test(
		`class Base {}
class Sup {
  static Base = Base;
}
class Outer extends Sup {
  static m() {
    return class extends super.Base {};
  }
}
`,
		`class Base {}
class Sup {}
Sup.Base = Base;
class Outer extends Sup {
  static m() {
    const CLASS_DECL$0 = class extends Sup.Base {};
    return CLASS_DECL$0;
  }
}
`,
	)
}

// Verifies that temporary receiver declarations extracted for private member
// accesses inside labeled loops or labeled statements climb past enclosing
// LABEL nodes rather than corrupting LABEL with a 3rd child.
func (c *round94ClassNormalizationCase) TestPrivateAccessInLabeledLoopAndStatement() {
	c.Test(
		`class C {
  #m() { return false; }
  test(getObj) {
    outer: inner: while (getObj()?.#m()) {
      break outer;
    }
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$2 = Object.create(null, {
  m: {
    value: function() {
      return false;
    }
  }
});
class C {
  constructor() {
    const PRIVATE$3 = Object.create(PRIVATE_PROTO$2);
    PRIVATE_MAP$0.set(this, PRIVATE$3);
  }
  test(getObj) {
    let TMP$1;
    outer: inner: for (; TMP$1 = getObj(), TMP$1 == null ? void 0 : PRIVATE_MAP$0.get(TMP$1).m.call(TMP$1); ) {
      break outer;
    }
  }
}
`,
	)

	c.enableNormalize = false
	c.DisableNormalize()
	c.Test(
		`class C {
  #m() {}
  test(getObj) {
    lbl: getObj()?.#m();
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$2 = Object.create(null, {
  m: {
    value: function() {}
  }
});
class C {
  constructor() {
    const PRIVATE$3 = Object.create(PRIVATE_PROTO$2);
    PRIVATE_MAP$0.set(this, PRIVATE$3);
  }
  test(getObj) {
    let TMP$1;
    lbl: TMP$1 = getObj(), TMP$1 == null ? void 0 : PRIVATE_MAP$0.get(TMP$1).m.call(TMP$1);
  }
}
`,
	)
}

func TestRound94_Es6NormalizeClasses_NilNodeParentSafety(t *testing.T) {
	compiler := NewCompiler()
	compiler.Init([]*SourceFile{}, []*SourceFile{}, NewCompilerOptions())
	pass := NewEs6NormalizeClasses(compiler)

	// Traversing a detached block containing `this` and an object literal with a
	// computed property must not panic on NilNode parents.
	script := compiler.ParseTestCode("this; ({['k']: 1});")
	detached := rhino.IR_Block()
	for script.HasChildren() {
		detached.AddChildToBack(script.RemoveFirstChild())
	}
	NodeTraversal_Traverse(compiler, detached, pass)
	ast.Check(detached, func(n rhino.Node, problem string) {
		t.Fatalf("ast.Check failed on %v: %s", n, problem)
	})
}

type round94DestructuringCase struct {
	CompilerTestCaseBase
}

func newRound94DestructuringCase() *round94DestructuringCase {
	c := &round94DestructuringCase{}
	c.InitCompilerTestCaseBaseWithExterns(c, Es6RewriteDestructuringTest_EXTERNS_BASE)
	return c
}

func (c *round94DestructuringCase) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.DisableNormalize()
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
}

func (c *round94DestructuringCase) GetProcessor(compiler *Compiler) CompilerPass {
	optimizer := NewPhaseOptimizer(compiler, nil)
	optimizer.AddOneTimePass(
		CompilerTestCase_MakePassFactory(
			"injectTranspilationRuntimeLibraries",
			func(abstractCompiler AbstractCompiler) CompilerPass {
				return NewInjectTranspilationRuntimeLibraries(abstractCompiler)
			},
		),
	)
	optimizer.AddOneTimePass(
		CompilerTestCase_MakePassFactory(
			"es6RewriteDestructuring",
			func(abstractCompiler AbstractCompiler) CompilerPass {
				return NewEs6RewriteDestructuring(
					abstractCompiler,
					Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_ALL_OBJECT_PATTERNS,
				)
			},
		),
	)
	return optimizer
}

func TestRound94_Es6RewriteDestructuring_LabeledStatements(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return newRound94DestructuringCase() })
}

func (c *round94DestructuringCase) TestUnnormalizedLabeledDestructuringWrapsInBlock() {
	c.Test(
		`var a, b; lbl: [a, b] = [1, 2];`,
		`var a, b;
lbl: {
  var $jscomp$destructuring$var0 = (0, $jscomp.makeIterator)([1, 2]);
  a = $jscomp$destructuring$var0.next().value;
  b = $jscomp$destructuring$var0.next().value;
}
`,
	)

	c.Test(
		`var a, b; lbl: ({a: a, b: b} = {a: 1, b: 2});`,
		`var a, b;
lbl: {
  var $jscomp$destructuring$var0 = {a: 1, b: 2};
  a = $jscomp$destructuring$var0.a;
  b = $jscomp$destructuring$var0.b;
}
`,
	)

	c.Test(
		`lbl: var [a, b] = [1, 2];`,
		`lbl: {
  var $jscomp$destructuring$var0 = (0, $jscomp.makeIterator)([1, 2]);
  var a = $jscomp$destructuring$var0.next().value;
  var b = $jscomp$destructuring$var0.next().value;
}
`,
	)
}

func TestRound94_Es6RewriteDestructuring_EmptyNestingStackSafety(t *testing.T) {
	compiler := NewCompiler()
	opts := NewCompilerOptions()
	opts.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	compiler.Init(
		[]*SourceFile{SourceFile_FromCode("externs.js", Es6RewriteDestructuringTest_EXTERNS_BASE)},
		[]*SourceFile{SourceFile_FromCode("input.js", "const {a: a} = obj; const {b: b, ...rest} = obj;")},
		opts,
	)
	compiler.ParseInputs()
	NewInjectTranspilationRuntimeLibraries(compiler).Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	pass := NewEs6RewriteDestructuring(
		compiler,
		Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_OBJECT_REST,
	)

	script := ast.MustRoot(compiler.GetJsRoot()).MustFirstItem().MustScript()
	firstDecl := script.MustFirstItem().MustNameDeclaration()
	secondDecl := ast.MustNameDeclaration(firstDecl.Node().GetNext())
	noRestPat := firstDecl.MustOnlyDeclarator().MustDestructuringLhs().Pattern().MustObjectPattern()
	withRestPat := secondDecl.MustOnlyDeclarator().MustDestructuringLhs().Pattern().MustObjectPattern()

	nt := acquireNodeTraversalDirect(compiler, pass, nil, false)
	nt.currentChangeScope = script.Node()
	defer releaseNodeTraversal(nt)
	// Direct Visit without prior ShouldTraverse push must not panic on empty stack
	// in either the hasRest=false or hasRest=true branches.
	pass.Visit(nt, noRestPat.Node(), noRestPat.Node().GetParent())
	pass.Visit(nt, withRestPat.Node(), withRestPat.Node().GetParent())
	ast.Check(script.Node(), func(n rhino.Node, problem string) {
		t.Fatalf("ast.Check failed on %v: %s", n, problem)
	})
}

func TestRound94_Es6RewriteGenerators_CatchWithoutBinding(t *testing.T) {
	// 1. Standalone Es6RewriteGenerators on an un-lowered `catch {}` must report
	// a clean IllegalStateException rather than a Go nil-pointer dereference on
	// NilNode.Detach().
	compiler := NewCompiler()
	options := NewCompilerOptions()
	options.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	compiler.Init(
		[]*SourceFile{
			SourceFile_FromCode(
				"externs.js",
				NewTestExternsBuilder().AddJSCompLibraries().Build(),
			),
		},
		[]*SourceFile{
			SourceFile_FromCode(
				"input.js",
				"function* g() { try { yield 1; } catch { yield 2; } }",
			),
		},
		options,
	)
	compiler.ParseInputs()
	NewInjectTranspilationRuntimeLibraries(compiler).Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)

	ex := AssertThrows[jutil.IRuntimeException](func() {
		NewEs6RewriteGenerators(compiler).Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	})
	AssertThat(ex.GetMessage()).Contains("Expected catch parameter to be normalized to a NAME")

	// 2. Full transpilation pipeline (which runs RewriteCatchWithNoBinding before
	// Es6RewriteGenerators) compiles ES2019 optional catch binding inside
	// generators cleanly.
	fullCompiler := NewCompiler()
	fullOpts := NewCompilerOptions()
	fullOpts.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	fullOpts.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	res := fullCompiler.Compile(
		[]*SourceFile{SourceFile_FromCode("externs.js", "")},
		[]*SourceFile{SourceFile_FromCode("input.js", "function* g() { try { yield 1; } catch { yield 2; } }")},
		fullOpts,
	)
	if len(res.Errors) != 0 {
		t.Fatalf("expected 0 errors in full generator optional-catch compile, got: %v", res.Errors)
	}
}

func TestRound94_RuntimeJsLibManager_NilNodeSafety(t *testing.T) {
	changeTracker := NewChangeTracker()

	// NilNode from resourceProvider in INJECT mode throws NullPointerException cleanly.
	mgrNilResource := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_INJECT,
		func(resourceName, path string) rhino.Node { return rhino.NilNode },
		changeTracker,
		func() rhino.Node { return rhino.IR_Script() },
	)
	_ = AssertThrows[jutil.INullPointerException](func() {
		mgrNilResource.EnsureLibraryInjected("base", false)
	})

	// NilNode from nodeForCodeInsertion in INJECT mode throws NullPointerException cleanly.
	mgrNilTarget := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_INJECT,
		func(resourceName, path string) rhino.Node {
			s := rhino.IR_Script()
			s.AddChildToBack(rhino.IR_Empty())
			return s
		},
		changeTracker,
		func() rhino.Node { return rhino.NilNode },
	)
	_ = AssertThrows[jutil.INullPointerException](func() {
		mgrNilTarget.EnsureLibraryInjected("base", false)
	})

	// NilNode from nodeForCodeInsertion in EXTERN_FIELD_NAMES mode throws NullPointerException cleanly.
	mgrNilExtern := RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_EXTERN_FIELD_NAMES,
		func(resourceName, path string) rhino.Node { return rhino.NilNode },
		changeTracker,
		func() rhino.Node { return rhino.NilNode },
	)
	_ = AssertThrows[jutil.INullPointerException](func() {
		mgrNilExtern.InjectLibForField("$jscomp.makeIterator")
	})
}
