// Port of com/google/javascript/jscomp/Es6ConvertSuperTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/Es6ConvertSuperTest.java)

/*
 * Copyright 2014 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type Es6ConvertSuperTest struct {
	CompilerTestCaseBase
}

// Java: Es6ConvertSuperTest()
func NewEs6ConvertSuperTest() *Es6ConvertSuperTest {
	c := &Es6ConvertSuperTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, CompilerTestCase_MINIMAL_EXTERNS)
	return c
}

func TestEs6ConvertSuperTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewEs6ConvertSuperTest() })
}

// Java: Es6ConvertSuperTest.setUp()
func (c *Es6ConvertSuperTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableNormalize()
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2016)
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.EnableTypeCheck()
	c.EnableTypeInfoValidation()
	c.EnableScriptFeatureValidation()
	c.ReplaceTypesWithColors()
	c.EnableMultistageCompilation()
}

// Java: Es6ConvertSuperTest.getProcessor(Compiler)
func (c *Es6ConvertSuperTest) GetProcessor(compiler *Compiler) CompilerPass {
	optimizer := NewPhaseOptimizer(compiler, nil)
	optimizer.AddOneTimePass(CompilerTestCase_MakePassFactory("es6NormalizeClasses", func(c AbstractCompiler) CompilerPass { return NewEs6NormalizeClasses(c) }))
	optimizer.AddOneTimePass(CompilerTestCase_MakePassFactory("es6ConvertSuper", func(c AbstractCompiler) CompilerPass { return NewEs6ConvertSuper(c) }))
	return optimizer
}

// Instance `super` resolution
//
// Java: Es6ConvertSuperTest.testCallingSuperInstanceProperty()
func (c *Es6ConvertSuperTest) TestCallingSuperInstanceProperty() {
	c.Test(
		c.Externs(
			`class A {
  constructor() { }

  /** @param {number} x @return {string} */
  g(x) { }
}
`),
		c.Srcs(
			`class B extends A {
  constructor() { super(); }

  f() { super.g(3); }
}
`),
		c.Expected(
			`class B extends A {
  constructor() { super(); }

  f() { A.prototype.g.call(this, 3); }
}
`))

	// get types we need to check
	classAPrototypeType := colors.Color_CreateUnion(
		CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "A").GetRootNode().GetColor().GetPrototypes())

	aDotGMethodType := CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "A").
		FindMethodDefinition("g").
		GetRootNode().
		GetFirstChild().
		GetColor()

	classBInstanceType := colors.Color_CreateUnion(
		CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "B").
			GetRootNode().
			GetColor().
			GetInstanceColors())

	// A.prototype.g.call(this, 3)
	callNode := CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "B").
		FindMethodDefinition("f").
		FindMatchingQNameReferences("A.prototype.g.call")[0].
		// GETPROP node for A.prototype.g.call
		GetParent() // CALL node within the statement
	AssertNode(callNode).
		HasToken(rhino.Token_CALL).
		HasLineno(4). // position and length of `super.g(3)`
		HasCharno(8)
	AssertNode(callNode).HasColorThat().IsEqualTo(colors.StandardColors_STRING)

	// A.prototype.g.call
	callee := callNode.GetFirstChild()
	AssertNode(callee).
		MatchesQualifiedName("A.prototype.g.call").
		HasLineno(4). // position and length of `super.g`
		HasCharno(14)
	AssertNode(callee).HasColorThat().IsEqualTo(colors.StandardColors_UNKNOWN)

	// A.prototype.g
	superDotGReplacement := callee.GetFirstChild()
	AssertNode(superDotGReplacement).
		MatchesQualifiedName("A.prototype.g").
		HasLineno(4). // position and length of `super.g`
		HasCharno(14)
	AssertNode(superDotGReplacement).HasColorThat().IsEqualTo(aDotGMethodType)

	// A.prototype
	superReplacement := superDotGReplacement.GetFirstChild()
	AssertNode(superReplacement).
		MatchesQualifiedName("A.prototype").
		HasLineno(4). // position and length of `super`
		HasCharno(8).
		HasOriginalName("super")
	AssertNode(superReplacement).HasColorThat().IsEqualTo(classAPrototypeType)

	// `this` node from `A.prototype.g.call(this, 3)`
	thisNode := callee.GetNext()
	AssertNode(thisNode).
		HasToken(rhino.Token_THIS).
		HasLineno(4). // position and length of `super.g`
		HasCharno(14).
		IsIndexable(false) // there's no direct correlation with text in the original source
	AssertNode(thisNode).HasColorThat().IsEqualTo(classBInstanceType)
}

// Java: Es6ConvertSuperTest.testCallingSuperInstanceElement()
func (c *Es6ConvertSuperTest) TestCallingSuperInstanceElement() {
	c.Test(
		c.Externs(
			`/** @dict */
class A {
  constructor() { }

  /** @param {number} x */
  ['g'](x) { };
}
`),
		c.Srcs(
			`/** @dict */
class B extends A {
  constructor() { super(); }

  ['f']() { super['g'](4); }
}
`),
		c.Expected(
			`class B extends A {
  constructor() {
    super();
  }

  ['f']() { A.prototype['g'].call(this, 4); }
}
`))

	// get types we need to check
	classAPrototypeType := colors.Color_CreateUnion(
		CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "A").GetRootNode().GetColor().GetPrototypes())

	classBInstanceType := colors.Color_CreateUnion(
		CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "B").
			GetRootNode().
			GetColor().
			GetInstanceColors())

	// A.prototype['g'].call(this, 4)
	callNode := c.GetLastCompiler().
		GetJsRoot().      // root
		GetFirstChild().  // script
		GetFirstChild().  // class B
		GetLastChild().   // class B's body
		GetLastChild().   // ['f']() { ... }
		GetSecondChild(). // () { ... }
		GetLastChild().   // { ... }
		GetFirstChild().  // statement `A.prototype['g'].call(this, 4);`
		GetOnlyChild()    // CALL node within the statement
	AssertNode(callNode).
		HasToken(rhino.Token_CALL).
		HasLineno(5). // position and length of `super['g'](4)`
		HasCharno(12)
	// computed property prevents doing any better than StandardColors.UNKNOWN
	AssertNode(callNode).HasColorThat().IsEqualTo(colors.StandardColors_UNKNOWN)

	// A.prototype['g'].call
	callee := callNode.GetFirstChild()
	AssertNode(callee).
		HasToken(rhino.Token_GETPROP).
		HasLineno(5). // position and length of `super['g']`
		HasCharno(12)
	AssertNode(callee).HasColorThat().IsEqualTo(colors.StandardColors_UNKNOWN)

	// A.prototype['g']
	superGetelemReplacement := callee.GetFirstChild()
	AssertNode(superGetelemReplacement).
		HasToken(rhino.Token_GETELEM).
		HasLineno(5). // position and length of `super['g']`
		HasCharno(12)
	// computed property prevents doing any better than StandardColors.UNKNOWN
	AssertNode(superGetelemReplacement).HasColorThat().IsEqualTo(colors.StandardColors_UNKNOWN)

	// A.prototype
	superReplacement := superGetelemReplacement.GetFirstChild()
	AssertNode(superReplacement).
		MatchesQualifiedName("A.prototype").
		HasLineno(5). // position and length of `super`
		HasCharno(12).
		HasOriginalName("super")
	AssertNode(superReplacement).HasColorThat().IsEqualTo(classAPrototypeType)

	// `this` node from `A.prototype['g'].call(this, 3)`
	thisNode := callee.GetNext()
	AssertNode(thisNode).
		HasToken(rhino.Token_THIS).
		HasLineno(5). // position and length of `super['g']`
		HasCharno(12).
		IsIndexable(false) // there's no direct correlation with text in the original source
	AssertNode(thisNode).HasColorThat().IsEqualTo(classBInstanceType)
}

// Java: Es6ConvertSuperTest.testAccessingSuperInstanceProperty()
func (c *Es6ConvertSuperTest) TestAccessingSuperInstanceProperty() {
	c.Test(
		c.Externs(
			`class A {
  constructor() { }

  /** @param {number} x */
  g(x) { }
}
`),
		c.Srcs(
			`class B extends A {
  constructor() { super(); }

  f() { var t = super.g; }
}
`),
		c.Expected(
			`class B extends A {
  constructor() { super(); }

  f() { var t = A.prototype.g; }
}
`))

	// get types we need to check
	classAPrototypeType := colors.Color_CreateUnion(
		CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "A").GetRootNode().GetColor().GetPrototypes())
	aDotGMethodType := CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "A").
		FindMethodDefinition("g").
		GetRootNode().
		GetFirstChild().
		GetColor()

	// A.prototype.g
	superDotGReplacement := c.GetLastCompiler().
		GetJsRoot().     // root
		GetFirstChild(). // script
		GetFirstChild(). // class B
		GetLastChild().  // class B's body
		GetLastChild().  // f
		GetOnlyChild().  // function that implements f() {}
		GetLastChild().  // method body
		GetFirstChild(). // `var t = A.prototype.g;`
		GetFirstChild(). // `t = A.prototype.g`
		GetOnlyChild()   // A.prototype.g
	AssertNode(superDotGReplacement).
		MatchesQualifiedName("A.prototype.g").
		HasLineno(4). // position and length of `super.g`
		HasCharno(22)
	AssertNode(superDotGReplacement).HasColorThat().IsEqualTo(aDotGMethodType)

	// A.prototype
	superReplacement := superDotGReplacement.GetFirstChild()
	AssertNode(superReplacement).
		MatchesQualifiedName("A.prototype").
		HasLineno(4). // position and length of `super`
		HasCharno(16).
		HasOriginalName("super")
	AssertNode(superReplacement).HasColorThat().IsEqualTo(classAPrototypeType)
}

// Java: Es6ConvertSuperTest.testAccessingSuperInstanceElement()
func (c *Es6ConvertSuperTest) TestAccessingSuperInstanceElement() {
	c.Test(
		c.Externs(
			`/** @dict */
class A {
  constructor() { }

  /** @param {number} x */
  ['g'](x) { };
}
`),
		c.Srcs(
			`/** @dict */
class B extends A {
  constructor() { super(); }

  ['f']() { var t = super['g']; }
}
`),
		c.Expected(
			`class B extends A {
  constructor() {
    super();
  }

  ['f']() { var t = A.prototype['g']; }
}
`))

	// get types we need to check
	classAPrototypeType := colors.Color_CreateUnion(
		CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "A").GetRootNode().GetColor().GetPrototypes())

	// A.prototype['g']
	superGetelemReplacement := c.GetLastCompiler().
		GetJsRoot().      // root
		GetFirstChild().  // script
		GetFirstChild().  // class B
		GetLastChild().   // class B's body
		GetLastChild().   // ['f']
		GetSecondChild(). // () { ... }
		GetLastChild().   // { ... }
		GetFirstChild().  // `var t = A.prototype['g'];`
		GetFirstChild().  // `t = A.prototype['g']`
		GetOnlyChild()    // A.prototype['g']
	AssertNode(superGetelemReplacement).
		HasToken(rhino.Token_GETELEM).
		HasLineno(5). // position and length of `super['g']`
		HasCharno(20)
	// getelem prevents us doing any better than unknown type
	AssertNode(superGetelemReplacement).HasColorThat().IsEqualTo(colors.StandardColors_UNKNOWN)

	// A.prototype
	superReplacement := superGetelemReplacement.GetFirstChild()
	AssertNode(superReplacement).
		MatchesQualifiedName("A.prototype").
		HasLineno(5). // position and length of `super`
		HasCharno(20).
		HasOriginalName("super")
	AssertNode(superReplacement).HasColorThat().IsEqualTo(classAPrototypeType)
}

// Java: Es6ConvertSuperTest.testCannotAssignToSuperInstanceProperty()
func (c *Es6ConvertSuperTest) TestCannotAssignToSuperInstanceProperty() {
	c.TestError(
		`class A {
  constructor() { }

  /** @param {number} x */
  g(x) { }
}

class B extends A {
  constructor() { super(); }

  f() { super.g = 5; }
}
`,
		TranspilationUtil_CANNOT_CONVERT_YET)
}

// Java: Es6ConvertSuperTest.testCannotAssignToSuperInstanceElement()
func (c *Es6ConvertSuperTest) TestCannotAssignToSuperInstanceElement() {
	c.TestError(
		`/** @dict */
class A {
  constructor() { }

  /** @param {number} x */
  ['g'](x) { }
}

class B extends A {
  constructor() { super(); }

  ['f']() { super['g'] = 5; }
}
`,
		TranspilationUtil_CANNOT_CONVERT_YET)
}

// Static `super` resolution
//
// Java: Es6ConvertSuperTest.testCallingSuperStaticProperty()
func (c *Es6ConvertSuperTest) TestCallingSuperStaticProperty() {
	c.Test(
		c.Externs(
			`class A {
  constructor() { }

  /** @param {number} x @return {string} */
  static g(x) { }
}
`),
		c.Srcs(
			`class B extends A {
  constructor() { super(); }

  static f() { super.g(3); }
}
`),
		c.Expected(
			`class B extends A {
  constructor() {
    super();
  }
  /** @nocollapse */
  static f() {
    A.g.call(this, 3);
  }
}
`))

	// get types we need to check
	classAType := CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "A").GetRootNode().GetColor()
	aDotGMethodType := CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "A").
		FindMethodDefinition("g").
		GetRootNode().
		GetFirstChild().
		GetColor()
	aDotGDotCallType := colors.StandardColors_UNKNOWN // colors do not track ".call" type

	classBType := CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "B").GetRootNode().GetColor()

	// A.g.call(this, 3)
	callNode := c.GetLastCompiler().
		GetJsRoot().     // root
		GetFirstChild(). // script
		GetFirstChild(). // class B
		GetLastChild().  // class B's body
		GetLastChild().  // static f
		GetOnlyChild().  // function that implements static f() {}
		GetLastChild().  // method body
		GetFirstChild(). // `A.g.call(this, 3);`
		GetOnlyChild()   // CALL node within the statement
	AssertNode(callNode).
		HasToken(rhino.Token_CALL).
		HasLineno(4). // position and length of `super.g(3)`
		HasCharno(15)
	AssertNode(callNode).HasColorThat().IsEqualTo(colors.StandardColors_STRING)

	// A.g.call
	callee := callNode.GetFirstChild()
	AssertNode(callee).
		MatchesQualifiedName("A.g.call").
		HasLineno(4). // position and length of `super.g`
		HasCharno(21)
	AssertNode(callee).HasColorThat().IsEqualTo(aDotGDotCallType)

	// A.g
	superDotGReplacement := callee.GetFirstChild()
	AssertNode(superDotGReplacement).
		MatchesQualifiedName("A.g").
		HasLineno(4). // position and length of `super.g`
		HasCharno(21)
	AssertNode(superDotGReplacement).HasColorThat().IsEqualTo(aDotGMethodType)

	// A
	superReplacement := superDotGReplacement.GetFirstChild()
	AssertNode(superReplacement).
		MatchesQualifiedName("A").
		HasLineno(4). // position and length of `super`
		HasCharno(15).
		HasOriginalName("super")
	AssertNode(superReplacement).HasColorThat().IsEqualTo(classAType)

	thisNode := callee.GetNext()
	AssertNode(thisNode).HasType(rhino.Token_THIS)
	AssertNode(thisNode).HasColorThat().IsEqualTo(classBType)
}

// Java: Es6ConvertSuperTest.testCallingSuperStaticElement()
func (c *Es6ConvertSuperTest) TestCallingSuperStaticElement() {
	c.Test(
		c.Externs(
			`/** @dict */
class A {
  constructor() { }

  /** @param {number} x */
  static ['g'](x) { };
}
`),
		c.Srcs(
			`/** @dict */
class B extends A {
  constructor() { super(); }

  static ['f']() { super['g'](4); }
}
`),
		c.Expected(
			`class B extends A {
  constructor() {
    super();
  }
  /** @nocollapse */
  static ["f"]() {
    A["g"].call(this, 4);
  }
}
`))

	// get types we need to check
	classAType := CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "A").GetRootNode().GetColor()
	classBType := CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "B").GetRootNode().GetColor()

	// A['g'].call(this, 4)
	callNode := c.GetLastCompiler().
		GetJsRoot().      // root
		GetFirstChild().  // script
		GetFirstChild().  // class B
		GetLastChild().   // class B's body
		GetLastChild().   // static ['f']() { ... }
		GetSecondChild(). // () { ... }
		GetLastChild().   // { ... }
		GetFirstChild().  // statement `A['g'].call(this, 4);`
		GetOnlyChild()    // CALL node within the statement
	AssertNode(callNode).
		HasToken(rhino.Token_CALL).
		HasLineno(5). // position and length of `super['g'](4)`
		HasCharno(19)
	// computed property prevents doing any better than StandardColors.UNKNOWN
	AssertNode(callNode).HasColorThat().IsEqualTo(colors.StandardColors_UNKNOWN)

	// A['g'].call
	callee := callNode.GetFirstChild()
	AssertNode(callee).
		HasToken(rhino.Token_GETPROP).
		HasLineno(5). // position and length of `super['g']`
		HasCharno(19)
	AssertNode(callee).HasColorThat().IsEqualTo(colors.StandardColors_UNKNOWN)

	// A['g']
	superGetelemReplacement := callee.GetFirstChild()
	AssertNode(superGetelemReplacement).
		HasToken(rhino.Token_GETELEM).
		HasLineno(5). // position and length of `super['g']`
		HasCharno(19)
	// computed property prevents doing any better than StandardColors.UNKNOWN
	AssertNode(superGetelemReplacement).HasColorThat().IsEqualTo(colors.StandardColors_UNKNOWN)

	// A
	superReplacement := superGetelemReplacement.GetFirstChild()
	AssertNode(superReplacement).
		MatchesName("A").
		HasLineno(5). // position and length of `super`
		HasCharno(19).
		HasOriginalName("super")
	AssertNode(superReplacement).HasColorThat().IsEqualTo(classAType)

	// `this` node from `A['g'].call(this, 3)`
	thisNode := callee.GetNext()
	AssertNode(thisNode).
		HasToken(rhino.Token_THIS).
		HasLineno(5). // position and length of `super['g']`
		HasCharno(19).
		IsIndexable(false) // there's no direct correlation with text in the original source
	AssertNode(thisNode).HasColorThat().IsEqualTo(classBType)
}

// Java: Es6ConvertSuperTest.testAccessingSuperStaticProperty()
func (c *Es6ConvertSuperTest) TestAccessingSuperStaticProperty() {
	c.Test(
		c.Externs(
			`class A {
  constructor() { }

  /** @param {number} x */
  static g(x) { }
}
`),
		c.Srcs(
			`class B extends A {
  constructor() { super(); }

  static f() { var t = super.g; }
}
`),
		c.Expected(
			`class B extends A {
  constructor() { super(); }

  static f() { var t = A.g; }
}
`))

	// get types we need to check
	classAType := CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "A").GetRootNode().GetColor()
	aDotGMethodType := CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "A").
		FindMethodDefinition("g").
		GetRootNode().
		GetFirstChild().
		GetColor()

	// A.prototype.g
	superDotGReplacement := c.GetLastCompiler().
		GetJsRoot().     // root
		GetFirstChild(). // script
		GetFirstChild(). // class B
		GetLastChild().  // class B's body
		GetLastChild().  // f
		GetOnlyChild().  // function that implements f() {}
		GetLastChild().  // method body
		GetFirstChild(). // `var t = A.g;`
		GetFirstChild(). // `t = A.g`
		GetOnlyChild()   // A.g
	AssertNode(superDotGReplacement).
		MatchesQualifiedName("A.g").
		HasLineno(4). // position and length of `super.g`
		HasCharno(29)
	AssertNode(superDotGReplacement).HasColorThat().IsEqualTo(aDotGMethodType)

	// A.prototype
	superReplacement := superDotGReplacement.GetFirstChild()
	AssertNode(superReplacement).
		MatchesQualifiedName("A").
		HasLineno(4). // position and length of `super`
		HasCharno(23).
		HasOriginalName("super")
	AssertNode(superReplacement).HasColorThat().IsEqualTo(classAType)
}

// Java: Es6ConvertSuperTest.testAccessingSuperStaticElement()
func (c *Es6ConvertSuperTest) TestAccessingSuperStaticElement() {
	c.Test(
		c.Externs(
			`/** @dict */
class A {
  constructor() { }

  /** @param {number} x */
  static ['g'](x) { };
}
`),
		c.Srcs(
			`/** @dict */
class B extends A {
  constructor() { super(); }

  static ['f']() { var t = super['g']; }
}
`),
		c.Expected(
			`class B extends A {
  constructor() {
    super();
  }

  static ['f']() { var t = A['g']; }
}
`))

	// get types we need to check
	classAType := CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "A").GetRootNode().GetColor()

	// A.prototype.g
	superGetElemReplacement := c.GetLastCompiler().
		GetJsRoot().      // root
		GetFirstChild().  // script
		GetFirstChild().  // class B
		GetLastChild().   // class B's body
		GetLastChild().   // ['f']() { ... }
		GetSecondChild(). // () { ... }
		GetLastChild().   // { ... }
		GetFirstChild().  // `var t = A['g'];`
		GetFirstChild().  // `t = A['g']`
		GetOnlyChild()    // `A['g']`
	AssertNode(superGetElemReplacement).
		HasToken(rhino.Token_GETELEM).
		HasLineno(5). // position and length of `super['g']`
		HasCharno(27)
	// computed property prevents doing any better than StandardColors.UNKNOWN
	AssertNode(superGetElemReplacement).HasColorThat().IsEqualTo(colors.StandardColors_UNKNOWN)

	// A.prototype
	superReplacement := superGetElemReplacement.GetFirstChild()
	AssertNode(superReplacement).
		MatchesQualifiedName("A").
		HasLineno(5). // position and length of `super`
		HasCharno(27).
		HasOriginalName("super")
	AssertNode(superReplacement).HasColorThat().IsEqualTo(classAType)
}

// Getters and setters
//
// Java: Es6ConvertSuperTest.testResolvingSuperInGetter()
func (c *Es6ConvertSuperTest) TestResolvingSuperInGetter() {
	c.Test(
		c.Externs(
			`class A {
  constructor() { }

  /** @param {number} x @return {number} */
  g(x) { }
}
`),
		c.Srcs(
			`class B extends A {
  constructor() { super(); }

  get f() { super.g(3); }
}
`),
		c.Expected(
			`class B extends A {
  constructor() { super(); }

  get f() { A.prototype.g.call(this, 3); }
}
`))

	// get types we need to check
	classAPrototypeType := colors.Color_CreateUnion(
		CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "A").GetRootNode().GetColor().GetPrototypes())
	aDotGMethodType := CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "A").
		FindMethodDefinition("g").
		GetRootNode().
		GetFirstChild().
		GetColor()
	aDotGDotCallType := colors.StandardColors_UNKNOWN // colors do not track ".call" type

	classBInstanceType := colors.Color_CreateUnion(
		CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "B").
			GetRootNode().
			GetColor().
			GetInstanceColors())

	// A.prototype.g.call(this, 3)
	callNode := c.GetLastCompiler().
		GetJsRoot().     // root
		GetFirstChild(). // script
		GetFirstChild(). // class B
		GetLastChild().  // class B's body
		GetLastChild().  // get f() { ... }
		GetOnlyChild().  // () { ... }
		GetLastChild().  // { ... }
		GetFirstChild(). // statement `A.prototype.g.call(this, 3);`
		GetOnlyChild()   // CALL node within the statement
	AssertNode(callNode).
		HasToken(rhino.Token_CALL).
		HasLineno(4). // position and length of `super.g(3)`
		HasCharno(12)
	AssertNode(callNode).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)

	// A.prototype.g.call
	callee := callNode.GetFirstChild()
	AssertNode(callee).
		MatchesQualifiedName("A.prototype.g.call").
		HasLineno(4). // position and length of `super.g`
		HasCharno(18)
	AssertNode(callee).HasColorThat().IsEqualTo(aDotGDotCallType)

	// A.prototype.g
	superDotGReplacement := callee.GetFirstChild()
	AssertNode(superDotGReplacement).
		MatchesQualifiedName("A.prototype.g").
		HasLineno(4). // position and length of `super.g`
		HasCharno(18)
	AssertNode(superDotGReplacement).HasColorThat().IsEqualTo(aDotGMethodType)

	// A.prototype
	superReplacement := superDotGReplacement.GetFirstChild()
	AssertNode(superReplacement).
		MatchesQualifiedName("A.prototype").
		HasLineno(4). // position and length of `super`
		HasCharno(12).
		HasOriginalName("super")
	AssertNode(superReplacement).HasColorThat().IsEqualTo(classAPrototypeType)

	// `this` node from `A.prototype.g.call(this, 3)`
	thisNode := callee.GetNext()
	AssertNode(thisNode).
		HasToken(rhino.Token_THIS).
		HasLineno(4). // position and length of `super.g`
		HasCharno(18).
		IsIndexable(false) // there's no direct correlation with text in the original source
	AssertNode(thisNode).HasColorThat().IsEqualTo(classBInstanceType)
}

// Java: Es6ConvertSuperTest.testResolvingSuperInSetter()
func (c *Es6ConvertSuperTest) TestResolvingSuperInSetter() {
	c.Test(
		c.Externs(
			`class A {
  constructor() { }

  /** @param {number} x @return {string} */
  g(x) { }
}
`),
		c.Srcs(
			`class B extends A {
  constructor() { super(); }

  /** @param {number} y */
  set f(y) { super.g(y); }
}
`),
		c.Expected(
			`class B extends A {
  constructor() { super(); }

  set f(y) { A.prototype.g.call(this, y); }
}
`))

	// get types we need to check
	classAPrototypeType := colors.Color_CreateUnion(
		CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "A").GetRootNode().GetColor().GetPrototypes())
	aDotGMethodType := CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "A").
		FindMethodDefinition("g").
		GetRootNode().
		GetFirstChild().
		GetColor()
	aDotGDotCallType := colors.StandardColors_UNKNOWN // colors do not track ".call" type

	classBInstanceType := colors.Color_CreateUnion(
		CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "B").
			GetRootNode().
			GetColor().
			GetInstanceColors())

	// A.prototype.g.call(this, x)
	callNode := c.GetLastCompiler().
		GetJsRoot().     // root
		GetFirstChild(). // script
		GetFirstChild(). // class B
		GetLastChild().  // class B's body
		GetLastChild().  // set f(x) { ... }
		GetOnlyChild().  // (x) { ... }
		GetLastChild().  // { ... }
		GetFirstChild(). // statement `A.prototype.g.call(this, x);`
		GetOnlyChild()   // CALL node within the statement
	AssertNode(callNode).
		HasToken(rhino.Token_CALL).
		HasLineno(5). // position and length of `super.g(3)`
		HasCharno(13)
	AssertNode(callNode).HasColorThat().IsEqualTo(colors.StandardColors_STRING)

	// A.prototype.g.call
	callee := callNode.GetFirstChild()
	AssertNode(callee).
		MatchesQualifiedName("A.prototype.g.call").
		HasLineno(5). // position and length of `super.g`
		HasCharno(19)
	AssertNode(callee).HasColorThat().IsEqualTo(aDotGDotCallType)

	// A.prototype.g
	superDotGReplacement := callee.GetFirstChild()
	AssertNode(superDotGReplacement).
		MatchesQualifiedName("A.prototype.g").
		HasLineno(5). // position and length of `super.g`
		HasCharno(19)
	AssertNode(superDotGReplacement).HasColorThat().IsEqualTo(aDotGMethodType)

	// A.prototype
	superReplacement := superDotGReplacement.GetFirstChild()
	AssertNode(superReplacement).
		MatchesQualifiedName("A.prototype").
		HasLineno(5). // position and length of `super`
		HasCharno(13).
		HasOriginalName("super")
	AssertNode(superReplacement).HasColorThat().IsEqualTo(classAPrototypeType)

	// `this` node from `A.prototype.g.call(this, 3)`
	thisNode := callee.GetNext()
	AssertNode(thisNode).
		HasToken(rhino.Token_THIS).
		HasLineno(5). // position and length of `super.g`
		HasCharno(19).
		IsIndexable(false) // there's no direct correlation with text in the original source
	AssertNode(thisNode).HasColorThat().IsEqualTo(classBInstanceType)
}

// Java: Es6ConvertSuperTest.testAccessingStaticGetterViaSuper()
func (c *Es6ConvertSuperTest) TestAccessingStaticGetterViaSuper() {
	c.Test(
		c.Externs(
			`class Parent {
  constructor() {}
  /** @return {string} */
  static getName() {}
  /** @return {string} */
  static get greeting() {}
}
`),
		c.Srcs(
			`class Child extends Parent {
  constructor() {
    super();
  }
  static getName() {
    return 'Child';
  }
  static getGreeting() {
    return super.greeting;
  }
}
`),
		c.Expected(
			`class Child extends Parent {
  constructor() {
    super();
  }
  static getName() {
    return 'Child';
  }
  /** @nocollapse */
  static getGreeting() {
    return Reflect.get(
        Parent, JSCompiler_renameProperty('greeting', Parent), this);
  }
}
`))

	// get types we need to check
	parentType := CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "Parent").GetRootNode().GetColor()
	childType := CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "Child").GetRootNode().GetColor()

	// return Reflect.get(Parent, JSCompiler_renameProperty("greeting", Parent), this);
	returnStmt := CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "Child").
		FindMethodDefinition("getGreeting").
		GetRootNode().  // MEMBER_FUNCTION_DEF
		GetOnlyChild(). // FUNCTION
		GetLastChild(). // BLOCK
		GetOnlyChild()  // RETURN
	AssertNode(returnStmt).HasToken(rhino.Token_RETURN)

	// Reflect.get(Parent, JSCompiler_renameProperty("greeting", Parent), this)
	reflectGetCall := returnStmt.GetOnlyChild()
	AssertNode(reflectGetCall).HasToken(rhino.Token_CALL).HasColorThat().IsEqualTo(colors.StandardColors_STRING)

	// Reflect.get
	callee := reflectGetCall.GetFirstChild()
	AssertNode(callee).MatchesQualifiedName("Reflect.get")

	// Parent
	parentName := callee.GetNext()
	AssertNode(parentName).
		MatchesQualifiedName("Parent").
		HasOriginalName("super").
		HasColorThat().
		IsEqualTo(parentType)

	// JSCompiler_renameProperty("greeting", Parent)
	renamePropertyCall := parentName.GetNext()
	AssertNode(renamePropertyCall).
		HasToken(rhino.Token_CALL).
		HasColorThat().
		IsEqualTo(colors.StandardColors_STRING)

	{
		// "greeting"
		greetingString := renamePropertyCall.GetSecondChild()
		AssertNode(greetingString).
			IsString("greeting").
			HasColorThat().
			IsEqualTo(colors.StandardColors_STRING)

		// Parent
		parentNameArg := greetingString.GetNext()
		AssertNode(parentNameArg).MatchesQualifiedName("Parent").HasColorThat().IsEqualTo(parentType)
	}

	// this
	thisNode := renamePropertyCall.GetNext()
	AssertNode(thisNode).HasToken(rhino.Token_THIS).HasColorThat().IsEqualTo(childType)
}

// Constructor synthesis
//
// Java: Es6ConvertSuperTest.testSynthesizingConstructorOfBaseClassInSource()
func (c *Es6ConvertSuperTest) TestSynthesizingConstructorOfBaseClassInSource() {
	c.Test(
		c.Externs(""),
		c.Srcs(
			`class A { } // Force wrapping.

class B extends A {
  constructor() { super(); }
}
`),
		c.Expected(
			`class A {
  constructor() { }
}

class B extends A {
  constructor() { super(); }
}
`))

	// class A { ... }
	classANode := CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "A").GetRootNode()
	classAConstructorType := classANode.GetColor()

	// constructor() { }
	constructorMemberFunctionDefForA := classANode.
		GetLastChild(). // { ... }
		GetFirstChild()
	AssertNode(constructorMemberFunctionDefForA).
		IsMemberFunctionDef("constructor").
		HasLineno(1). // synthetic constructor gets position and length of original class definition
		HasCharno(0).
		IsIndexable(false)
	AssertNode(constructorMemberFunctionDefForA).HasColorThat().IsEqualTo(classAConstructorType)

	constructorFunctionForA := constructorMemberFunctionDefForA.GetOnlyChild()
	AssertNode(constructorFunctionForA).
		HasToken(rhino.Token_FUNCTION).
		HasLineno(1). // synthetic constructor gets position and length of original class definition
		HasCharno(0).
		IsIndexable(false)
	AssertNode(constructorFunctionForA).HasColorThat().IsEqualTo(classAConstructorType)
}

// Java: Es6ConvertSuperTest.testSynthesizingConstructorOfDerivedClassInSource()
func (c *Es6ConvertSuperTest) TestSynthesizingConstructorOfDerivedClassInSource() {
	c.Test(
		c.Externs(NewTestExternsBuilder().AddArguments().Build()),
		c.Srcs(
			`class A { // Force wrapping.
  constructor() { }
}

class B extends A { }
`),
		c.Expected(
			`class A {
  constructor() { }
}

class B extends A {
  constructor() { super(...arguments); }
}
`))

	// class A { ... }
	classANode := CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "A").GetRootNode()
	classAConstructorType := classANode.GetColor()

	// class B extends A { ... }
	classBNode := CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "B").GetRootNode()
	classBConstructorType := classBNode.GetColor()
	classBInstanceType := colors.Color_CreateUnion(classBConstructorType.GetInstanceColors())

	// constructor() { }
	constructorMemberFunctionDefForB := classBNode.
		GetLastChild(). // { ... }
		GetFirstChild()
	AssertNode(constructorMemberFunctionDefForB).
		IsMemberFunctionDef("constructor").
		HasLineno(5). // synthetic constructor gets position and length of original class definition
		HasCharno(0).
		IsIndexable(false)
	AssertNode(constructorMemberFunctionDefForB).HasColorThat().IsEqualTo(classBConstructorType)

	constructorFunctionForB := constructorMemberFunctionDefForB.GetOnlyChild()
	AssertNode(constructorFunctionForB).
		HasToken(rhino.Token_FUNCTION).
		HasLineno(5). // synthetic constructor gets position and length of original class definition
		HasCharno(0).
		IsIndexable(false)
	AssertNode(constructorFunctionForB).HasColorThat().IsEqualTo(classBConstructorType)

	// super(...arguments)
	superConstructorCall := constructorFunctionForB.
		GetLastChild().  // constructor body
		GetFirstChild(). // expr_result statement
		GetOnlyChild()
	AssertNode(superConstructorCall).
		HasToken(rhino.Token_CALL).
		HasLineno(5). // synthetic constructor gets position and length of original class definition
		HasCharno(0).
		IsIndexable(false)
	AssertNode(superConstructorCall).HasColorThat().IsEqualTo(classBInstanceType)

	superNode := superConstructorCall.GetFirstChild()
	AssertNode(superNode).
		HasToken(rhino.Token_SUPER).
		HasLineno(5). // synthetic constructor gets position and length of original class definition
		HasCharno(0).
		IsIndexable(false)
	AssertNode(superNode).HasColorThat().IsEqualTo(classAConstructorType)

	argumentsNode := superNode.
		GetNext(). // ...arguments
		GetOnlyChild()
	AssertNode(argumentsNode).
		IsName("arguments").
		HasLineno(5). // synthetic constructor gets position and length of original class definition
		HasCharno(0).
		IsIndexable(false)
}

// Java: Es6ConvertSuperTest.testSynthesizingConstructorOfBaseInterface()
func (c *Es6ConvertSuperTest) TestSynthesizingConstructorOfBaseInterface() {
	c.Test(
		c.Externs(""),
		c.Srcs("/** @interface */ class A { }"),
		c.Expected("/** @interface */ class A { constructor() { } }"))

	// class A { ... }
	classANode := CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "A").GetRootNode()
	classAConstructorType := classANode.GetColor()

	// constructor() { }
	constructorMemberFunctionDefForA := classANode.
		GetLastChild(). // { ... }
		GetFirstChild()
	AssertNode(constructorMemberFunctionDefForA).
		IsMemberFunctionDef("constructor").
		HasLineno(1). // synthetic constructor gets position and length of original class definition
		HasCharno(18).
		IsIndexable(false)
	AssertNode(constructorMemberFunctionDefForA).HasColorThat().IsEqualTo(classAConstructorType)

	constructorFunctionForA := constructorMemberFunctionDefForA.GetOnlyChild()
	AssertNode(constructorFunctionForA).
		HasToken(rhino.Token_FUNCTION).
		HasLineno(1). // synthetic constructor gets position and length of original class definition
		HasCharno(18).
		IsIndexable(false)
	AssertNode(constructorFunctionForA).HasColorThat().IsEqualTo(classAConstructorType)
}

// Java: Es6ConvertSuperTest.testSynthesizingConstructorOfDerivedInterface()
func (c *Es6ConvertSuperTest) TestSynthesizingConstructorOfDerivedInterface() {
	c.Test(
		c.Externs(
			`/** @interface */ // Force wrapping.
class A {
  constructor() { }
}
`),
		c.Srcs("/** @interface */ class B extends A { }"),
		c.Expected(
			`/** @interface */
class B extends A {
  constructor() { }
}
`))
}

// Java: Es6ConvertSuperTest.testStrippingSuperCallFromConstructorOfDerivedInterface()
func (c *Es6ConvertSuperTest) TestStrippingSuperCallFromConstructorOfDerivedInterface() {
	c.Test(
		c.Externs(
			`const namespace = {};

/** @interface */
namespace.A = class {
  constructor() { }
}
`),
		c.Srcs(
			`/** @interface */
class B extends namespace.A {
  constructor() { super(); }
}
`),
		c.Expected(
			`/** @interface */
class B extends namespace.A {
  constructor() { }
}
`))
}
