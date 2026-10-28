// Port of com/google/javascript/jscomp/GlobalNamespaceTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/GlobalNamespaceTest.java)

/*
 * Copyright 2008 The Closure Compiler Authors.
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
	"slices"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

// GlobalNamespaceTest: tests for GlobalNamespace.
type GlobalNamespaceTest struct {
	// Go-only: test-helper — the Go test running the current test method.
	t *testing.T

	lastCompiler                     *Compiler // nil is null
	assumeStaticInheritanceIsNotUsed bool
}

// Java: GlobalNamespaceTest()
func NewGlobalNamespaceTest() *GlobalNamespaceTest {
	return &GlobalNamespaceTest{lastCompiler: nil, assumeStaticInheritanceIsNotUsed: true}
}

func TestGlobalNamespaceTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewGlobalNamespaceTest() })
}

// JUnitSetT gives the instance the Go test running its test method.
func (c *GlobalNamespaceTest) JUnitSetT(t *testing.T) { c.t = t }

// JUnitBefore: the class has no @Before method.
func (c *GlobalNamespaceTest) JUnitBefore() {}

// JUnitAfter: the class has no @After method.
func (c *GlobalNamespaceTest) JUnitAfter() {}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.detectsPropertySetsInAssignmentOperators()
func (c *GlobalNamespaceTest) TestDetectsPropertySetsInAssignmentOperators() {
	namespace := c.Parse("const a = {b: 0}; a.b += 1; a.b = 2;")
	AssertThat(globalNamespaceTest_name(namespace.GetSlot("a.b")).GetGlobalSets()).IsEqualTo(3)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.detectsPropertySetsInDestructuring()
func (c *GlobalNamespaceTest) TestDetectsPropertySetsInDestructuring() {
	namespace := c.Parse("const a = {b: 0}; [a.b] = [1]; ({b: a.b} = {b: 2});")
	// TODO(b/120303257): this should be 3
	AssertThat(globalNamespaceTest_name(namespace.GetSlot("a.b")).GetGlobalSets()).IsEqualTo(1)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.detectsPropertySetsInIncDecOperators()
func (c *GlobalNamespaceTest) TestDetectsPropertySetsInIncDecOperators() {
	namespace := c.Parse("const a = {b: 0}; a.b++; a.b--;")
	AssertThat(globalNamespaceTest_name(namespace.GetSlot("a.b")).GetGlobalSets()).IsEqualTo(3)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.firstGlobalAssignmentIsConsideredDeclaration()
func (c *GlobalNamespaceTest) TestFirstGlobalAssignmentIsConsideredDeclaration() {
	namespace := c.Parse("")
	n := namespace.CreateNameForTesting("a")
	set1 := n.AddSingleRefForTesting(rhino.IR_Name("set1"), GlobalNamespace_Ref_Type_SET_FROM_GLOBAL)
	set2 := n.AddSingleRefForTesting(rhino.IR_Name("set2"), GlobalNamespace_Ref_Type_SET_FROM_GLOBAL)
	AssertThat(n.GetRefs()).ContainsExactly(set1, set2).InOrder()
	AssertThat(n.getDeclarationRef()).IsEqualTo(set1)
	AssertThat(n.GetGlobalSets()).IsEqualTo(2)
	n.RemoveRef(set1)
	// declaration moves to next global assignment when first is removed
	AssertThat(n.getDeclarationRef()).IsEqualTo(set2)
	AssertThat(n.GetGlobalSets()).IsEqualTo(1)
	AssertThat(n.GetRefs()).ContainsExactly(set2)
}

// Java: GlobalNamespaceTest.testReferencesToUndefinedRootName()
func (c *GlobalNamespaceTest) TestReferencesToUndefinedRootName() {
	namespace := c.Parse("a; a.b = 0; a.b; a?.b")
	AssertThat(globalNamespaceTest_name(namespace.GetSlot("a"))).IsNull()
	AssertThat(globalNamespaceTest_name(namespace.GetSlot("a.b"))).IsNull()
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.nullishCoalesce()
func (c *GlobalNamespaceTest) TestNullishCoalesce() {
	namespace := c.Parse("var a = a ?? {};")
	a := globalNamespaceTest_name(namespace.GetSlot("a"))
	AssertThat(a).IsNotNull()
	AssertThat(a.GetRefs()).HasSize(2)
	AssertThat(a.GetLocalSets()).IsEqualTo(0)
	AssertThat(a.GetGlobalSets()).IsEqualTo(1)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.logicalAssignment1()
func (c *GlobalNamespaceTest) TestLogicalAssignment1() {
	namespace := c.Parse("var a = a ||= {};")
	a := globalNamespaceTest_name(namespace.GetSlot("a"))
	AssertThat(a).IsNotNull()
	AssertThat(a.GetRefs()).HasSize(2)
	AssertThat(a.GetLocalSets()).IsEqualTo(0)
	AssertThat(a.GetGlobalSets()).IsEqualTo(2)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.logicalAssignment2()
func (c *GlobalNamespaceTest) TestLogicalAssignment2() {
	namespace := c.Parse("var a = a || (a = {});")
	a := globalNamespaceTest_name(namespace.GetSlot("a"))
	AssertThat(a).IsNotNull()
	AssertThat(a.GetRefs()).HasSize(3)
	AssertThat(a.GetLocalSets()).IsEqualTo(0)
	AssertThat(a.GetGlobalSets()).IsEqualTo(2)
}

// Java: GlobalNamespaceTest.testlogicalAssignmentGets1()
func (c *GlobalNamespaceTest) TestlogicalAssignmentGets1() {
	namespace := c.Parse(
		`const ns = {};
ns.n1 = 1;
ns.n2 = 2;
ns.n1 ??= ns.n2;
ns.n1 &&= ns.n2;
`)
	bar := globalNamespaceTest_name(namespace.GetSlot("ns"))
	AssertThat(bar.GetGlobalSets()).IsEqualTo(1)
	AssertThat(bar.GetAliasingGets()).IsEqualTo(0)
	AssertThat(bar.GetTotalGets()).IsEqualTo(0)
	n1 := globalNamespaceTest_name(namespace.GetSlot("ns.n1"))
	AssertThat(n1.GetGlobalSets()).IsEqualTo(3)
	AssertThat(n1.GetAliasingGets()).IsEqualTo(0)
	AssertThat(n1.GetTotalGets()).IsEqualTo(0)
	n2 := globalNamespaceTest_name(namespace.GetSlot("ns.n2"))
	AssertThat(n2.GetGlobalSets()).IsEqualTo(1)
	AssertThat(n2.GetAliasingGets()).IsEqualTo(2)
	AssertThat(n2.GetTotalGets()).IsEqualTo(2)
}

// Java: GlobalNamespaceTest.testlogicalAssignmentGets2()
func (c *GlobalNamespaceTest) TestlogicalAssignmentGets2() {
	namespace := c.Parse(
		`const ns = {};
ns.n1 = 1;
ns.n2 = 2;
ns.n1 ?? (ns.n1 = ns.n2);
ns.n1 && (ns.n1 = ns.n2);
`)
	bar := globalNamespaceTest_name(namespace.GetSlot("ns"))
	AssertThat(bar.GetGlobalSets()).IsEqualTo(1)
	AssertThat(bar.GetAliasingGets()).IsEqualTo(0)
	AssertThat(bar.GetTotalGets()).IsEqualTo(0)
	n1 := globalNamespaceTest_name(namespace.GetSlot("ns.n1"))
	AssertThat(n1.GetGlobalSets()).IsEqualTo(3)
	AssertThat(n1.GetAliasingGets()).IsEqualTo(2)
	AssertThat(n1.GetTotalGets()).IsEqualTo(4)
	n2 := globalNamespaceTest_name(namespace.GetSlot("ns.n2"))
	AssertThat(n2.GetGlobalSets()).IsEqualTo(1)
	AssertThat(n2.GetAliasingGets()).IsEqualTo(2)
	AssertThat(n2.GetTotalGets()).IsEqualTo(2)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.detectsPropertySetsInLogicalAssignmentOperators1()
func (c *GlobalNamespaceTest) TestDetectsPropertySetsInLogicalAssignmentOperators1() {
	namespace := c.Parse("const a = {b: 0}; a.b ||= 1; a.b = 2;")
	AssertThat(globalNamespaceTest_name(namespace.GetSlot("a.b")).GetGlobalSets()).IsEqualTo(3)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.detectsPropertySetsInLogicalAssignmentOperators2()
func (c *GlobalNamespaceTest) TestDetectsPropertySetsInLogicalAssignmentOperators2() {
	namespace := c.Parse("const a = {b: 0}; a.b || (a.b = 1); a.b = 2;")
	AssertThat(globalNamespaceTest_name(namespace.GetSlot("a.b")).GetGlobalSets()).IsEqualTo(3)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.hook()
func (c *GlobalNamespaceTest) TestHook() {
	namespace := c.Parse("var a = a ? a : {}")
	a := globalNamespaceTest_name(namespace.GetSlot("a"))
	AssertThat(a).IsNotNull()
	AssertThat(a.GetRefs()).HasSize(3)
	AssertThat(a.GetLocalSets()).IsEqualTo(0)
	AssertThat(a.GetGlobalSets()).IsEqualTo(1)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.localAssignmentWillNotBeConsideredADeclaration()
func (c *GlobalNamespaceTest) TestLocalAssignmentWillNotBeConsideredADeclaration() {
	namespace := c.Parse("")
	n := namespace.CreateNameForTesting("a")
	set1 := n.AddSingleRefForTesting(rhino.IR_Name("set1"), GlobalNamespace_Ref_Type_SET_FROM_GLOBAL)
	localSet := n.AddSingleRefForTesting(rhino.IR_Name("localSet"), GlobalNamespace_Ref_Type_SET_FROM_LOCAL)
	AssertThat(n.GetRefs()).ContainsExactly(set1, localSet).InOrder()
	AssertThat(n.getDeclarationRef()).IsEqualTo(set1)
	AssertThat(n.GetGlobalSets()).IsEqualTo(1)
	AssertThat(n.GetLocalSets()).IsEqualTo(1)
	n.RemoveRef(set1)
	// local set will not be used as the declaration
	AssertThat(n.getDeclarationRef()).IsNull()
	AssertThat(n.GetGlobalSets()).IsEqualTo(0)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.firstDeclarationJSDocAlwaysWins()
func (c *GlobalNamespaceTest) TestFirstDeclarationJSDocAlwaysWins() {
	namespace := c.Parse(
		`const X = {};
/** @type {symbol} */ // later assignment should win
X.number;
/** @type {number} */ // this is the JSDoc we should use
X.number = 3;
/** @type {string} */
X.number = 'hi';
/** @type {Object} */
X.number;
`)
	nameX := globalNamespaceTest_name(namespace.GetOwnSlot("X.number"))
	declarationRef := nameX.getDeclarationRef()
	AssertThat(declarationRef).IsNotNull()
	// make sure first assignment is considered to be the declaration
	declarationNode := declarationRef.GetNode()
	AssertNode(declarationNode).MatchesQualifiedName("X.number")
	assignNode := declarationNode.GetParent()
	AssertNode(assignNode).IsAssign()
	valueNode := declarationNode.GetNext()
	AssertNode(valueNode).IsNumber(3)
	// Make sure JSDoc on the first assignment is the JSDoc for the name
	jsDocInfo := nameX.GetJSDocInfo()
	AssertThat(jsDocInfo).IsNotNull()
	jsTypeExpression := jsDocInfo.GetType()
	AssertThat(jsTypeExpression).IsNotNull()
	jsType := jsTypeExpression.Evaluate( /* scope= */ nil, c.lastCompiler.GetTypeRegistry())
	rhinotesting.NewTruth(c.t).AssertType(jsType).IsNumber()
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.withoutAssignmentFirstQnameDeclarationStatementJSDocWins()
func (c *GlobalNamespaceTest) TestWithoutAssignmentFirstQnameDeclarationStatementJSDocWins() {
	namespace := c.Parse(
		`const X = {};
/** @type {string} */
X.number;
/** @type {Object} */
X.number;
`)
	nameX := globalNamespaceTest_name(namespace.GetOwnSlot("X.number"))
	declarationRef := nameX.getDeclarationRef()
	AssertThat(declarationRef).IsNull()
	// Make sure JSDoc on the first assignment is the JSDoc for the name
	jsDocInfo := nameX.GetJSDocInfo()
	AssertThat(jsDocInfo).IsNotNull()
	jsTypeExpression := jsDocInfo.GetType()
	AssertThat(jsTypeExpression).IsNotNull()
	jsType := jsTypeExpression.Evaluate( /* scope= */ nil, c.lastCompiler.GetTypeRegistry())
	rhinotesting.NewTruth(c.t).AssertType(jsType).IsString()
}

// Java: GlobalNamespaceTest.testSimpleSubclassingRefCollection()
func (c *GlobalNamespaceTest) TestSimpleSubclassingRefCollection() {
	namespace := c.Parse(
		`class Superclass {}
class Subclass extends Superclass {}
`)
	superclass := globalNamespaceTest_name(namespace.GetOwnSlot("Superclass"))
	AssertThat(superclass.GetRefs()).HasSize(2)
	AssertThat(superclass.GetSubclassingGets()).IsEqualTo(1)
}

// Java: GlobalNamespaceTest.testStaticInheritedReferencesDontReferToSuperclass()
func (c *GlobalNamespaceTest) TestStaticInheritedReferencesDontReferToSuperclass() {
	namespace := c.Parse(
		`class Superclass {
  static staticMethod() {}
}
class Subclass extends Superclass {}
Subclass.staticMethod();
Subclass.staticMethod?.();
Subclass?.staticMethod();
`)
	superclass := globalNamespaceTest_name(namespace.GetOwnSlot("Superclass"))
	AssertThat(superclass.GetSubclassingGets()).IsEqualTo(1)
	superclassStaticMethod := globalNamespaceTest_name(namespace.GetOwnSlot("Superclass.staticMethod"))
	AssertThat(superclassStaticMethod.GetRefs()).HasSize(1)
	AssertThat(superclassStaticMethod.getDeclarationRef()).IsNotNull()
	subclassStaticMethod := globalNamespaceTest_name(namespace.GetOwnSlot("Subclass.staticMethod"))
	// 2 references:
	// `Subclass.staticMethod()`
	// `Subclass.staticMethod?.()`
	// `SubClass?.staticMethod()` is a reference to `SubClass`, but not
	// to `SubClass.staticmethod`.
	AssertThat(subclassStaticMethod.GetRefs()).HasSize(2)
	AssertThat(subclassStaticMethod.getDeclarationRef()).IsNull()
	AssertThat(subclassStaticMethod.GetCallGets()).IsEqualTo(2)
	subclass := globalNamespaceTest_name(namespace.GetOwnSlot("Subclass"))
	AssertThat(subclass.GetRefs()).HasSize(2)
	// `class Subclass` is the declaration reference
	AssertThat(subclass.getDeclarationRef()).IsNotNull()
	// `SubClass?.staticMethod` is an aliasing get on `SubClass`
	AssertThat(subclass.GetAliasingGets()).IsEqualTo(1)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.updateRefNodeRejectsRedundantUpdate()
func (c *GlobalNamespaceTest) TestUpdateRefNodeRejectsRedundantUpdate() {
	namespace := c.Parse("const A = 3;")
	nameA := globalNamespaceTest_name(namespace.GetOwnSlot("A"))
	refA := nameA.GetFirstRef()
	AssertThrows[jutil.IIllegalArgumentException](func() { nameA.UpdateRefNode(refA, refA.GetNode()) })
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.updateRefNodeMovesRefFromOldNodeToNewNode()
func (c *GlobalNamespaceTest) TestUpdateRefNodeMovesRefFromOldNodeToNewNode() {
	namespace := c.Parse("const A = 3;")
	nameA := globalNamespaceTest_name(namespace.GetOwnSlot("A"))
	refA := nameA.GetFirstRef()
	oldNode := refA.GetNode()
	newNode := rhino.IR_Name("A")
	AssertThat(nameA.GetRefForNode(oldNode)).IsEqualTo(refA)
	nameA.UpdateRefNode(refA, newNode)
	AssertThat(refA.GetNode()).IsEqualTo(newNode)
	AssertThat(nameA.GetRefForNode(oldNode)).IsNull()
	AssertThat(nameA.GetRefForNode(newNode)).IsEqualTo(refA)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.updateRefNodeCanSetNodeToNullButPreventsFurtherUpdates()
func (c *GlobalNamespaceTest) TestUpdateRefNodeCanSetNodeToNullButPreventsFurtherUpdates() {
	namespace := c.Parse("const A = 3;")
	nameA := globalNamespaceTest_name(namespace.GetOwnSlot("A"))
	refA := nameA.GetFirstRef()
	oldNode := refA.GetNode()
	AssertThat(nameA.GetRefForNode(oldNode)).IsEqualTo(refA)
	nameA.UpdateRefNode(refA, rhino.NilNode)
	AssertThat(refA.GetNode()).IsNull()
	AssertThat(nameA.GetRefForNode(oldNode)).IsNull()
	// cannot get refs for null
	AssertThrows[jutil.INullPointerException](func() { nameA.GetRefForNode(rhino.NilNode) })
	// cannot update the node again once it's been set to null
	AssertThrows[jutil.IIllegalArgumentException](func() { nameA.UpdateRefNode(refA, rhino.NilNode) })
	AssertThrows[jutil.IIllegalStateException](func() { nameA.UpdateRefNode(refA, oldNode) })
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.updateRefNodeRejectsNodeWithExistingRefs()
func (c *GlobalNamespaceTest) TestUpdateRefNodeRejectsNodeWithExistingRefs() {
	namespace := c.Parse(
		`const A = 3; // declaration ref
A;
`)
	// use ref
	nameA := globalNamespaceTest_name(namespace.GetOwnSlot("A"))
	declarationRef := nameA.getDeclarationRef()
	useRef := nameA.GetRefs()[1]
	// use ref is 2nd
	useNode := useRef.GetNode()
	AssertThrows[jutil.IIllegalArgumentException](func() { nameA.UpdateRefNode(declarationRef, useNode) })
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.confirmTwinsAreCreated()
func (c *GlobalNamespaceTest) TestConfirmTwinsAreCreated() {
	namespace := c.Parse(
		`let A;
const B = A = 3;
`)
	// A will have twin refs here
	nameA := globalNamespaceTest_name(namespace.GetOwnSlot("A"))
	// first ref is declaration of A
	twinRef := nameA.GetRefs()[1]
	// second is the GET_AND_SET twin
	// confirm that they start as twins
	AssertThat(twinRef.IsTwin()).IsTrue()
	AssertThat(twinRef.IsSetFromGlobal()).IsTrue()
	AssertThat(twinRef.IsAliasingGet()).IsTrue()
	oldNode := twinRef.GetNode()
	// confirm that it is associated with oldNode
	AssertThat(nameA.GetRefForNode(oldNode)).IsEqualTo(twinRef)
	// confirm that nameA correctly tracks its aliasingGets and globalSets
	AssertThat(nameA.GetGlobalSets()).IsEqualTo(2)
	AssertThat(nameA.GetAliasingGets()).IsEqualTo(1)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.updateRefNodeCanRemoveTwinRefs()
func (c *GlobalNamespaceTest) TestUpdateRefNodeCanRemoveTwinRefs() {
	namespace := c.Parse(
		`let A;
const B = A = 3;
`)
	// A will have twin refs here
	nameA := globalNamespaceTest_name(namespace.GetOwnSlot("A"))
	// first ref is declaration of A
	twinRef := nameA.GetRefs()[1]
	// second is the GET_AND_SET twin
	oldNode := twinRef.GetNode()
	// move the getTwinRef
	newNode := rhino.IR_Name("A")
	nameA.UpdateRefNode(twinRef, newNode)
	// see confirmTwinsAreCreated() for verification of the original twin relationship
	// confirm that getTwinRef has been updated
	AssertThat(twinRef.GetNode()).IsEqualTo(newNode)
	AssertThat(nameA.GetRefForNode(newNode)).IsEqualTo(twinRef)
	AssertThat(twinRef.IsTwin()).IsTrue()
	// confirm that all references to oldNode are removed.
	AssertThat(nameA.GetRefForNode(oldNode)).IsNull()
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.removeTwinRef()
func (c *GlobalNamespaceTest) TestRemoveTwinRef() {
	namespace := c.Parse(
		`let A;
const B = A = 3;
`)
	// A will have twin refs here
	nameA := globalNamespaceTest_name(namespace.GetOwnSlot("A"))
	// first ref is declaration of A
	twinRef := nameA.GetRefs()[1]
	// second is the GET_AND_SET twin
	// see confirmTwinsAreCreated() for verification of the original twin relationship
	oldNode := twinRef.GetNode()
	// confirm that they are both associated with oldNode
	AssertThat(nameA.GetRefForNode(oldNode)).IsEqualTo(twinRef)
	nameA.RemoveRef(twinRef)
	AssertThat(nameA.GetRefs()).DoesNotContain(twinRef)
	AssertThat(nameA.GetRefForNode(oldNode)).IsNull()
	AssertThat(nameA.GetAliasingGets()).IsEqualTo(0)
	AssertThat(nameA.GetGlobalSets()).IsEqualTo(1)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.rescanningExistingNodesDoesNotCreateDuplicateRefs()
func (c *GlobalNamespaceTest) TestRescanningExistingNodesDoesNotCreateDuplicateRefs() {
	namespace := c.Parse("class Foo {} const Bar = Foo; const Baz = Bar;")
	foo := globalNamespaceTest_name(namespace.GetOwnSlot("Foo"))
	bar := globalNamespaceTest_name(namespace.GetOwnSlot("Bar"))
	baz := globalNamespaceTest_name(namespace.GetOwnSlot("Baz"))
	originalFooRefs := slices.Clone(foo.GetRefs())
	originalBarRefs := slices.Clone(bar.GetRefs())
	originalBazRefs := slices.Clone(baz.GetRefs())
	// Rescan all of the nodes for which we got refs as if they were newly added
	root := c.lastCompiler.GetJsRoot()
	var astChanges astChangeSet
	for _, name := range []*GlobalNamespace_Name{foo, bar, baz} {
		for _, ref := range name.GetRefs() {
			astChanges.add(c.CreateGlobalAstChangeForNode(root, ref.GetNode()))
		}
	}
	namespace.ScanNewNodes(&astChanges)
	// We should get the same Name objects
	AssertThat(globalNamespaceTest_name(namespace.GetOwnSlot("Foo"))).IsEqualTo(foo)
	AssertThat(globalNamespaceTest_name(namespace.GetOwnSlot("Bar"))).IsEqualTo(bar)
	AssertThat(globalNamespaceTest_name(namespace.GetOwnSlot("Baz"))).IsEqualTo(baz)
	// ...and they should contain the same refs with no duplicates added
	AssertThat(foo.GetRefs()).ContainsExactlyElementsIn(originalFooRefs).InOrder()
	AssertThat(bar.GetRefs()).ContainsExactlyElementsIn(originalBarRefs).InOrder()
	AssertThat(baz.GetRefs()).ContainsExactlyElementsIn(originalBazRefs).InOrder()
}

// Java: GlobalNamespaceTest.testScanFromNodeDoesntDuplicateVarDeclarationSets()
func (c *GlobalNamespaceTest) TestScanFromNodeDoesntDuplicateVarDeclarationSets() {
	namespace := c.Parse("class Foo {} const Bar = Foo; const Baz = Bar;")
	foo := globalNamespaceTest_name(namespace.GetOwnSlot("Foo"))
	AssertThat(foo.GetAliasingGets()).IsEqualTo(1)
	baz := globalNamespaceTest_name(namespace.GetOwnSlot("Baz"))
	AssertThat(baz.GetGlobalSets()).IsEqualTo(1)
	// Replace "const Baz = Bar" with "const Baz = Foo"
	root := c.lastCompiler.GetJsRoot()
	barRef := root.GetFirstChild().GetLastChild().GetFirstFirstChild()
	jutil.CheckStateMsg1(barRef.GetString() == "Bar", "%s", barRef)
	fooName := rhino.IR_Name("Foo")
	barRef.ReplaceWith(fooName)
	// Rescan the new nodes
	namespace.ScanNewNodes(globalNamespaceTest_setOf(c.CreateGlobalAstChangeForNode(root, fooName)))
	AssertThat(foo.GetAliasingGets()).IsEqualTo(2)
	// A bug in scanFromNode used to make this `2`
	AssertThat(baz.GetGlobalSets()).IsEqualTo(1)
}

// Java: GlobalNamespaceTest.testScanFromNodeAddsReferenceToParentGetprop()
func (c *GlobalNamespaceTest) TestScanFromNodeAddsReferenceToParentGetprop() {
	namespace := c.Parse("const x = {bar: 0}; const y = x; const baz = y.bar;")
	xbar := globalNamespaceTest_name(namespace.GetOwnSlot("x.bar"))
	AssertThat(xbar.GetAliasingGets()).IsEqualTo(0)
	baz := globalNamespaceTest_name(namespace.GetOwnSlot("baz"))
	AssertThat(baz.GetGlobalSets()).IsEqualTo(1)
	// Replace "const baz = y.bar" with "const baz = x.bar"
	root := c.lastCompiler.GetJsRoot()
	yRef := root.GetFirstChild().GetLastChild().GetFirstFirstChild().GetFirstChild()
	jutil.CheckStateMsg1(yRef.GetString() == "y", "%s", yRef)
	xName := rhino.IR_Name("x")
	yRef.ReplaceWith(xName)
	// Rescan the new nodes
	namespace.ScanNewNodes(globalNamespaceTest_setOf(c.CreateGlobalAstChangeForNode(root, xName)))
	AssertThat(xbar.GetAliasingGets()).IsEqualTo(1)
	AssertThat(baz.GetGlobalSets()).IsEqualTo(1)
	var xBarGet *GlobalNamespace_Ref
	for _, ref := range xbar.GetRefs() {
		if ref.IsAliasingGet() {
			xBarGet = ref
			break
		}
	}
	jutil.CheckNotNullMsg(xBarGet, "No value present") // Optional.get()
	AssertThat(xBarGet.GetNode()).IsEqualTo(xName.GetParent())
	AssertThat(xBarGet.IsAliasingGet()).IsTrue()
	AssertThat(xBarGet.GetChunk()).IsEqualTo(xbar.getDeclarationRef().GetChunk())
}

// Java: GlobalNamespaceTest.testScanFromNodeNoticesHasOwnProperty()
func (c *GlobalNamespaceTest) TestScanFromNodeNoticesHasOwnProperty() {
	namespace := c.Parse("const x = {bar: 0}; const y = x; y.hasOwnProperty('bar');")
	xName := globalNamespaceTest_name(namespace.GetOwnSlot("x"))
	yName := globalNamespaceTest_name(namespace.GetOwnSlot("y"))
	AssertThat(xName.UsesHasOwnProperty()).IsFalse()
	AssertThat(yName.UsesHasOwnProperty()).IsTrue()
	// Replace "const baz = y.bar" with "const baz = x.bar"
	root := c.lastCompiler.GetJsRoot()
	yDotHasOwnProperty := root.GetFirstChild(). // SCRIPT
							GetLastChild(). // EXPR_RESULT `y.hasOwnProperty('bar');`
							GetFirstFirstChild()
	// `y.hasOwnProperty`
	yNode := yDotHasOwnProperty.GetFirstChild()
	// `y`
	AssertNode(yDotHasOwnProperty).MatchesQualifiedName("y.hasOwnProperty")
	xNode := rhino.IR_Name("x")
	yNode.ReplaceWith(xNode)
	// Rescan the new nodes
	// In this case the new node is `x.hasOwnProperty`, since that's the full, new qualified name.
	namespace.ScanNewNodes(globalNamespaceTest_setOf(c.CreateGlobalAstChangeForNode(root, yDotHasOwnProperty)))
	AssertThat(xName.UsesHasOwnProperty()).IsTrue()
}

// Java: GlobalNamespaceTest.createGlobalAstChangeForNode(Node, Node)
func (c *GlobalNamespaceTest) CreateGlobalAstChangeForNode(jsRoot rhino.Node, n rhino.Node) GlobalNamespace_AstChange {
	// This only creates a global scope, so don't use this with local nodes
	globalScope := NewSyntacticScopeCreator(c.lastCompiler).CreateScopeOfScope(jsRoot, nil)
	var firstChunk *JSChunk
	if chunks := c.lastCompiler.GetChunks(); len(chunks) > 0 {
		firstChunk = chunks[0]
	}
	return NewGlobalNamespace_SimpleAstChange(n, firstChunk, globalScope)
}

// Java: GlobalNamespaceTest.testCollapsing_forEscapedConstructor_ignoringStaticInheritance()
func (c *GlobalNamespaceTest) TestCollapsing_forEscapedConstructor_ignoringStaticInheritance() {
	namespace := c.Parse(
		`/** @constructor */
function Bar() {}
use(Bar);
`)
	bar := globalNamespaceTest_name(namespace.GetSlot("Bar"))
	AssertThat(bar.CanCollapse()).IsTrue()
	// trivially true, already collapsed
	// we collapse properties of Bar even though it's escaped, intentionally unsafe.
	// this is mostly to support minification for goog.provide namespaces containing @constructors
	AssertThat(bar.CanCollapseUnannotatedChildNames()).IsTrue()
}

// Java: GlobalNamespaceTest.testCollapsing_forEscapedConstructor_consideringStaticInheritance()
func (c *GlobalNamespaceTest) TestCollapsing_forEscapedConstructor_consideringStaticInheritance() {
	c.assumeStaticInheritanceIsNotUsed = false
	namespace := c.Parse(
		`/** @constructor */
function Bar() {}
use(Bar);
`)
	bar := globalNamespaceTest_name(namespace.GetSlot("Bar"))
	AssertThat(bar.CanCollapse()).IsTrue()
	// trivially true, already collapsed
	AssertThat(bar.CanCollapseUnannotatedChildNames()).IsFalse()
}

// Java: GlobalNamespaceTest.testInlinability_forAliasingPropertyOnEscapedConstructor_ignoringStaticInheritance()
func (c *GlobalNamespaceTest) TestInlinability_forAliasingPropertyOnEscapedConstructor_ignoringStaticInheritance() {
	namespace := c.Parse(
		`var prop = 1;
/** @constructor */
var Foo = function() {}

Foo.prop = prop;

/** @constructor */
function Bar() {}
Bar.aliasOfFoo = Foo; // alias Foo
use(Bar); // uninlinable alias of Bar
const BarAlias = Bar; // inlinable alias of Bar
alert(Bar.aliasOfFoo.prop);
alert(BarAlias.aliasOfFoo.prop);
`)
	barAliasOfFoo := globalNamespaceTest_name(namespace.GetSlot("Bar.aliasOfFoo"))
	barAliasInlinability := barAliasOfFoo.CalculateInlinability()
	// We should convert references to `Bar.aliasOfFoo.prop` to become `Foo.prop`
	// because...
	AssertThat(barAliasInlinability.ShouldInlineUsages()).IsTrue()
	// However, we should not remove the assignment (`Bar.aliasOfFoo = Foo`) that creates the alias,
	// because "BarAlias" still needs to be inlined to "Bar", which will create another usage of
	// "Bar.aliasOfFoo" in the last line, We will locate the value to inline Bar.aliasOfFoo
	// again from `Bar.aliasOfFoo = Foo`.
	AssertThat(barAliasInlinability.ShouldRemoveDeclaration()).IsFalse()
}

// Java: GlobalNamespaceTest.testInlinability_forAliasingPropertyOnEscapedConstructor_consideringStaticInheritance()
func (c *GlobalNamespaceTest) TestInlinability_forAliasingPropertyOnEscapedConstructor_consideringStaticInheritance() {
	c.assumeStaticInheritanceIsNotUsed = false
	namespace := c.Parse(
		`var prop = 1;
/** @constructor */
var Foo = function() {}

Foo.prop = prop;

/** @constructor */
function Bar() {}
Bar.aliasOfFoo = Foo; // alias Foo
use(Bar); // uninlinable alias of Bar
const BarAlias = Bar; // inlinable alias of Bar
alert(Bar.aliasOfFoo.prop);
alert(BarAlias.aliasOfFoo.prop);
`)
	barAliasOfFoo := globalNamespaceTest_name(namespace.GetSlot("Bar.aliasOfFoo"))
	barAliasInlinability := barAliasOfFoo.CalculateInlinability()
	AssertThat(barAliasInlinability.ShouldInlineUsages()).IsFalse()
	AssertThat(barAliasInlinability.ShouldRemoveDeclaration()).IsFalse()
}

// Java: GlobalNamespaceTest.testClassPrototypeProp()
func (c *GlobalNamespaceTest) TestClassPrototypeProp() {
	ns := c.Parse("class C { x() {} }")
	AssertThat(globalNamespaceTest_name(ns.GetSlot("C.x"))).IsNull()
}

// Java: GlobalNamespaceTest.testClassStaticField_withInitializer()
func (c *GlobalNamespaceTest) TestClassStaticField_withInitializer() {
	ns := c.Parse(
		`class C {
  static x = 1;
}
`)
	c_ := globalNamespaceTest_name(ns.GetSlot("C"))
	cDotX := globalNamespaceTest_name(ns.GetSlot("C.x"))
	AssertThat(c_.GetGlobalSets()).IsEqualTo(1)
	AssertThat(c_.PropsSlice()).ContainsExactly(cDotX)
	AssertThat(cDotX.GetGlobalSets()).IsEqualTo(1)
	AssertThat(cDotX.GetParent()).IsEqualTo(c_)
	AssertThat(cDotX.CanCollapse()).IsTrue()
}

// Java: GlobalNamespaceTest.testClassStaticField_withoutInitializer()
func (c *GlobalNamespaceTest) TestClassStaticField_withoutInitializer() {
	ns := c.Parse(
		`class C {
  /** @type {number} */
  static x;
}
`)
	c_ := globalNamespaceTest_name(ns.GetSlot("C"))
	cDotX := globalNamespaceTest_name(ns.GetSlot("C.x"))
	AssertThat(c_.GetGlobalSets()).IsEqualTo(1)
	AssertThat(c_.PropsSlice()).ContainsExactly(cDotX)
	AssertThat(cDotX.GetGlobalSets()).IsEqualTo(1)
	AssertThat(cDotX.GetParent()).IsEqualTo(c_)
	AssertThat(cDotX.CanCollapse()).IsTrue()
	jsDocInfo := cDotX.GetJSDocInfo()
	AssertThat(jsDocInfo).IsNotNull()
	jsTypeExpression := jsDocInfo.GetType()
	AssertThat(jsTypeExpression).IsNotNull()
	jsType := jsTypeExpression.Evaluate( /* scope= */ nil, c.lastCompiler.GetTypeRegistry())
	rhinotesting.NewTruth(c.t).AssertType(jsType).IsNumber()
}

// Java: GlobalNamespaceTest.testClassStaticField_withSuper()
func (c *GlobalNamespaceTest) TestClassStaticField_withSuper() {
	c.assumeStaticInheritanceIsNotUsed = false
	ns := c.Parse(
		`class A {
  static y = 1;
}
class C extends A {
  static x = super.y;
}
`)
	cDotX := globalNamespaceTest_name(ns.GetSlot("C.x"))
	AssertThat(cDotX.CanCollapse()).IsFalse()
}

// Java: GlobalNamespaceTest.testClassStaticAndPrototypePropWithSameName()
func (c *GlobalNamespaceTest) TestClassStaticAndPrototypePropWithSameName() {
	ns := c.Parse("class C { x() {} static x() {} }")
	c_ := globalNamespaceTest_name(ns.GetSlot("C"))
	cDotX := globalNamespaceTest_name(ns.GetSlot("C.x"))
	AssertThat(c_.GetGlobalSets()).IsEqualTo(1)
	AssertThat(c_.PropsSlice()).ContainsExactly(cDotX)
	AssertThat(cDotX.GetGlobalSets()).IsEqualTo(1)
	AssertThat(cDotX.GetParent()).IsEqualTo(c_)
}

// Java: GlobalNamespaceTest.testLocalVarsDefinedinStaticBlocks()
func (c *GlobalNamespaceTest) TestLocalVarsDefinedinStaticBlocks() {
	namespace := c.Parse("class C{ static{ var x; }}")
	AssertThat(globalNamespaceTest_name(namespace.GetSlot("x"))).IsNull()
}

// Java: GlobalNamespaceTest.testAddPropertytoGlobalObjectinClassStaticBlock()
func (c *GlobalNamespaceTest) TestAddPropertytoGlobalObjectinClassStaticBlock() {
	namespace := c.Parse("const a = {}; class C{ static { a.b = 1;}}")
	a := globalNamespaceTest_name(namespace.GetSlot("a"))
	AssertThat(a.GetGlobalSets()).IsEqualTo(1)
	AssertThat(a.GetTotalSets()).IsEqualTo(1)
	ab := globalNamespaceTest_name(namespace.GetSlot("a.b"))
	AssertThat(ab.GetParent()).IsEqualTo(a)
	AssertThat(ab.GetGlobalSets()).IsEqualTo(0)
	AssertThat(ab.GetLocalSets()).IsEqualTo(1)
}

// Java: GlobalNamespaceTest.testDirectGets()
func (c *GlobalNamespaceTest) TestDirectGets() {
	// None of the symbol uses here should be considered aliasing gets.
	namespace := c.Parse(
		`const ns = {};
ns.n1 = 1;
ns.n2 = 2;
ns.n1 === ns.n2;
ns.n1 == ns.n2;
ns.n1 !== ns.n2;
ns.n1 != ns.n2;
ns.n1 <  ns.n2;
ns.n1 <= ns.n2;
ns.n1 >  ns.n2;
ns.n1 >= ns.n2;
ns.n1 + ns.n2;
ns.n1 - ns.n2;
ns.n1 * ns.n2;
ns.n1 / ns.n2;
ns.n1 % ns.n2;
ns.n1 ** ns.n2;
ns.n1 & ns.n2;
ns.n1 | ns.n2;
ns.n1 ^ ns.n2;
ns.n1 << ns.n2;
ns.n1 >> ns.n2;
ns.n1 >>> ns.n2;
ns.n1 && ns.n2;
ns.n1 || ns.n2;
`)
	bar := globalNamespaceTest_name(namespace.GetSlot("ns"))
	AssertThat(bar.GetGlobalSets()).IsEqualTo(1)
	AssertThat(bar.GetAliasingGets()).IsEqualTo(0)
	// getting `ns.n1` doesn't count as a get on `ns`
	AssertThat(bar.GetTotalGets()).IsEqualTo(0)
	n1 := globalNamespaceTest_name(namespace.GetSlot("ns.n1"))
	AssertThat(n1.GetGlobalSets()).IsEqualTo(1)
	AssertThat(n1.GetAliasingGets()).IsEqualTo(0)
	AssertThat(n1.GetTotalGets()).IsEqualTo(22)
	n2 := globalNamespaceTest_name(namespace.GetSlot("ns.n2"))
	AssertThat(n2.GetGlobalSets()).IsEqualTo(1)
	AssertThat(n2.GetAliasingGets()).IsEqualTo(0)
	AssertThat(n2.GetTotalGets()).IsEqualTo(22)
}

// Java: GlobalNamespaceTest.testObjectPatternAliasInDeclaration()
func (c *GlobalNamespaceTest) TestObjectPatternAliasInDeclaration() {
	namespace := c.Parse("const ns = {a: 3}; const {a: b} = ns;")
	bar := globalNamespaceTest_name(namespace.GetSlot("ns"))
	AssertThat(bar.GetGlobalSets()).IsEqualTo(1)
	AssertThat(bar.GetAliasingGets()).IsEqualTo(1)
	AssertThat(bar.GetTotalGets()).IsEqualTo(1)
	nsA := globalNamespaceTest_name(namespace.GetSlot("ns.a"))
	AssertThat(nsA.GetGlobalSets()).IsEqualTo(1)
	AssertThat(nsA.GetAliasingGets()).IsEqualTo(1)
	b := globalNamespaceTest_name(namespace.GetSlot("b"))
	AssertThat(b.GetGlobalSets()).IsEqualTo(1)
	AssertThat(b.GetTotalGets()).IsEqualTo(0)
}

// Java: GlobalNamespaceTest.testConditionalDestructuringDoesNotHideAliasingGet()
func (c *GlobalNamespaceTest) TestConditionalDestructuringDoesNotHideAliasingGet() {
	namespace := c.Parse(
		`const ns1 = {a: 3};
const ns2 = {b: 3};
// Creates an aliasing get for both ns1 and ns2
const {a, b} = Math.random() ? ns1 : ns2;
`)
	ns1 := globalNamespaceTest_name(namespace.GetSlot("ns1"))
	AssertThat(ns1.GetAliasingGets()).IsEqualTo(1)
	AssertThat(ns1.GetTotalGets()).IsEqualTo(1)
	ns2 := globalNamespaceTest_name(namespace.GetSlot("ns2"))
	AssertThat(ns2.GetAliasingGets()).IsEqualTo(1)
	AssertThat(ns2.GetTotalGets()).IsEqualTo(1)
}

// Java: GlobalNamespaceTest.testNestedObjectPatternAliasInDeclaration()
func (c *GlobalNamespaceTest) TestNestedObjectPatternAliasInDeclaration() {
	namespace := c.Parse("const ns = {a: {b: 3}}; const {a: {b}} = ns;")
	bar := globalNamespaceTest_name(namespace.GetSlot("ns"))
	AssertThat(bar.GetGlobalSets()).IsEqualTo(1)
	AssertThat(bar.GetAliasingGets()).IsEqualTo(1)
	// we treat ns.a as having an 'aliasing' get since we don't traverse into the nested pattern
	nsA := globalNamespaceTest_name(namespace.GetSlot("ns.a"))
	AssertThat(nsA.GetGlobalSets()).IsEqualTo(1)
	AssertThat(nsA.GetAliasingGets()).IsEqualTo(1)
	nsAB := globalNamespaceTest_name(namespace.GetSlot("ns.a.b"))
	AssertThat(nsAB.GetGlobalSets()).IsEqualTo(1)
	// we don't consider this an 'aliasing get' because it's in a nested pattern
	AssertThat(nsAB.GetAliasingGets()).IsEqualTo(0)
	b := globalNamespaceTest_name(namespace.GetSlot("b"))
	AssertThat(b.GetGlobalSets()).IsEqualTo(1)
	AssertThat(b.GetTotalGets()).IsEqualTo(0)
}

// Java: GlobalNamespaceTest.testObjectPatternAliasInAssign()
func (c *GlobalNamespaceTest) TestObjectPatternAliasInAssign() {
	namespace := c.Parse("const ns = {a: 3}; const x = {}; ({a: x.y} = ns);")
	bar := globalNamespaceTest_name(namespace.GetSlot("ns"))
	AssertThat(bar.GetGlobalSets()).IsEqualTo(1)
	AssertThat(bar.GetAliasingGets()).IsEqualTo(1)
	nsA := globalNamespaceTest_name(namespace.GetSlot("ns.a"))
	AssertThat(nsA.GetGlobalSets()).IsEqualTo(1)
	AssertThat(nsA.GetAliasingGets()).IsEqualTo(1)
	xY := globalNamespaceTest_name(namespace.GetSlot("x.y"))
	// TODO(b/117673791): this should be 1
	AssertThat(xY.GetGlobalSets()).IsEqualTo(0)
}

// Java: GlobalNamespaceTest.testObjectPatternRestInDeclaration()
func (c *GlobalNamespaceTest) TestObjectPatternRestInDeclaration() {
	namespace := c.Parse("const ns = {a: 3}; const {a, ...b} = ns;")
	ns := globalNamespaceTest_name(namespace.GetSlot("ns"))
	AssertThat(ns.GetGlobalSets()).IsEqualTo(1)
	AssertThat(ns.GetTotalGets()).IsEqualTo(1)
	AssertThat(ns.GetAliasingGets()).IsEqualTo(1)
	nsA := globalNamespaceTest_name(namespace.GetSlot("ns.a"))
	AssertThat(nsA.GetGlobalSets()).IsEqualTo(1)
	AssertThat(nsA.GetTotalGets()).IsEqualTo(1)
	AssertThat(nsA.GetAliasingGets()).IsEqualTo(1)
	b := globalNamespaceTest_name(namespace.GetSlot("b"))
	AssertThat(b.GetGlobalSets()).IsEqualTo(1)
	AssertThat(b.GetTotalGets()).IsEqualTo(0)
}

// Java: GlobalNamespaceTest.testObjectPatternRestNestedInDeclaration()
func (c *GlobalNamespaceTest) TestObjectPatternRestNestedInDeclaration() {
	namespace := c.Parse("const ns = {a: 3, b: {}}; const {a, b: {...c}} = ns;")
	ns := globalNamespaceTest_name(namespace.GetSlot("ns"))
	AssertThat(ns.GetGlobalSets()).IsEqualTo(1)
	AssertThat(ns.GetTotalGets()).IsEqualTo(1)
	AssertThat(ns.GetAliasingGets()).IsEqualTo(1)
	nsB := globalNamespaceTest_name(namespace.GetSlot("ns.b"))
	AssertThat(nsB.GetGlobalSets()).IsEqualTo(1)
	AssertThat(nsB.GetTotalGets()).IsEqualTo(1)
	AssertThat(nsB.GetAliasingGets()).IsEqualTo(1)
}

// Java: GlobalNamespaceTest.testObjectPatternRestAliasInAssign()
func (c *GlobalNamespaceTest) TestObjectPatternRestAliasInAssign() {
	namespace := c.Parse("const ns = {a: 3}; const x = {}; ({a, ...x.y} = ns);")
	ns := globalNamespaceTest_name(namespace.GetSlot("ns"))
	AssertThat(ns.GetGlobalSets()).IsEqualTo(1)
	AssertThat(ns.GetAliasingGets()).IsEqualTo(1)
	nsA := globalNamespaceTest_name(namespace.GetSlot("ns.a"))
	AssertThat(nsA.GetGlobalSets()).IsEqualTo(1)
	AssertThat(nsA.GetAliasingGets()).IsEqualTo(1)
	xY := globalNamespaceTest_name(namespace.GetSlot("x.y"))
	// TODO(b/117673791): this should be 1
	AssertThat(xY.GetGlobalSets()).IsEqualTo(0)
}

// Java: GlobalNamespaceTest.testObjectPatternAliasInForOf()
func (c *GlobalNamespaceTest) TestObjectPatternAliasInForOf() {
	namespace := c.Parse("const ns = {a: 3}; for (const {a: b} of [ns]) {}")
	bar := globalNamespaceTest_name(namespace.GetSlot("ns"))
	AssertThat(bar.GetGlobalSets()).IsEqualTo(1)
	AssertThat(bar.GetAliasingGets()).IsEqualTo(1)
	// GlobalNamespace ignores for-of and array literals, not realizing that `b` reads `ns.a`
	nsA := globalNamespaceTest_name(namespace.GetSlot("ns.a"))
	AssertThat(nsA.GetGlobalSets()).IsEqualTo(1)
	AssertThat(nsA.GetAliasingGets()).IsEqualTo(0)
}

// Java: GlobalNamespaceTest.testObjectLitSpreadAliasInDeclaration()
func (c *GlobalNamespaceTest) TestObjectLitSpreadAliasInDeclaration() {
	namespace := c.Parse("const ns = {a: 3}; const {a} = {...ns};")
	ns := globalNamespaceTest_name(namespace.GetSlot("ns"))
	AssertThat(ns.GetGlobalSets()).IsEqualTo(1)
	AssertThat(ns.GetAliasingGets()).IsEqualTo(1)
	nsA := globalNamespaceTest_name(namespace.GetSlot("ns.a"))
	AssertThat(nsA.GetGlobalSets()).IsEqualTo(1)
	AssertThat(nsA.GetAliasingGets()).IsEqualTo(0)
	a := globalNamespaceTest_name(namespace.GetSlot("a"))
	AssertThat(a.GetGlobalSets()).IsEqualTo(1)
}

// Java: GlobalNamespaceTest.testObjectAssignOntoAGetProp()
func (c *GlobalNamespaceTest) TestObjectAssignOntoAGetProp() {
	namespace := c.Parse("const obj = {a:3}; const ns = {}; ns.a = Object.assign({}, obj); ")
	obj := globalNamespaceTest_name(namespace.GetSlot("obj"))
	AssertThat(obj.GetGlobalSets()).IsEqualTo(1)
	AssertThat(obj.GetAliasingGets()).IsEqualTo(1)
	objA := globalNamespaceTest_name(namespace.GetSlot("obj.a"))
	AssertThat(objA.GetGlobalSets()).IsEqualTo(1)
	AssertThat(objA.GetAliasingGets()).IsEqualTo(0)
	ns := globalNamespaceTest_name(namespace.GetSlot("ns"))
	AssertThat(ns.GetGlobalSets()).IsEqualTo(1)
	AssertThat(ns.GetAliasingGets()).IsEqualTo(0)
	nsA := globalNamespaceTest_name(namespace.GetSlot("ns.a"))
	AssertThat(nsA.GetGlobalSets()).IsEqualTo(1)
	AssertThat(nsA.GetAliasingGets()).IsEqualTo(0)
	AssertThat(nsA.IsObjectLiteral()).IsFalse()
	// `ns.a` is considered an "OTHER" type
}

// Java: GlobalNamespaceTest.testObjectLitSpreadOntoAGetProp()
func (c *GlobalNamespaceTest) TestObjectLitSpreadOntoAGetProp() {
	namespace := c.Parse("const obj = {a:3}; const ns = {}; ns.a = {...obj}")
	obj := globalNamespaceTest_name(namespace.GetSlot("obj"))
	AssertThat(obj.GetGlobalSets()).IsEqualTo(1)
	AssertThat(obj.GetAliasingGets()).IsEqualTo(1)
	objA := globalNamespaceTest_name(namespace.GetSlot("obj.a"))
	AssertThat(objA.GetGlobalSets()).IsEqualTo(1)
	AssertThat(objA.GetAliasingGets()).IsEqualTo(0)
	ns := globalNamespaceTest_name(namespace.GetSlot("ns"))
	AssertThat(ns.GetGlobalSets()).IsEqualTo(1)
	AssertThat(ns.GetAliasingGets()).IsEqualTo(0)
	nsA := globalNamespaceTest_name(namespace.GetSlot("ns.a"))
	AssertThat(nsA.GetGlobalSets()).IsEqualTo(1)
	AssertThat(nsA.GetAliasingGets()).IsEqualTo(0)
	AssertThat(nsA.IsObjectLiteral()).
		IsTrue()
	// `ns.a` is an "OBJECTLIT" type despite containing spread.
}

// Java: GlobalNamespaceTest.testObjectLitSpreadAliasInAssign()
func (c *GlobalNamespaceTest) TestObjectLitSpreadAliasInAssign() {
	namespace := c.Parse("const ns = {a: 3}; const x = {}; ({a: x.y} = {...ns});")
	ns := globalNamespaceTest_name(namespace.GetSlot("ns"))
	AssertThat(ns.GetGlobalSets()).IsEqualTo(1)
	AssertThat(ns.GetAliasingGets()).IsEqualTo(1)
	nsA := globalNamespaceTest_name(namespace.GetSlot("ns.a"))
	AssertThat(nsA.GetGlobalSets()).IsEqualTo(1)
	AssertThat(nsA.GetAliasingGets()).IsEqualTo(0)
	xY := globalNamespaceTest_name(namespace.GetSlot("x.y"))
	// TODO(b/117673791): this should be 1
	AssertThat(xY.GetGlobalSets()).IsEqualTo(0)
}

// Java: GlobalNamespaceTest.testLhsCastInAssignment()
func (c *GlobalNamespaceTest) TestLhsCastInAssignment() {
	// The type of the cast doesn't matter.
	// Casting is only legal JS syntax in simple assignments, not with destructuring or declaration.
	namespace := c.Parse("const ns = {}; const b = 5; /** @type {*} */ (ns.a) = b;")
	ns := globalNamespaceTest_name(namespace.GetSlot("ns"))
	AssertThat(ns.GetGlobalSets()).IsEqualTo(1)
	AssertThat(ns.GetTotalGets()).IsEqualTo(0)
	nsA := globalNamespaceTest_name(namespace.GetSlot("ns.a"))
	AssertThat(nsA.GetGlobalSets()).IsEqualTo(0)
	// TODO(b/127505242): Should be 1.
	AssertThat(nsA.GetTotalGets()).IsEqualTo(1)
	// TODO(b/127505242): Should be 0.
	b := globalNamespaceTest_name(namespace.GetSlot("b"))
	AssertThat(b.GetGlobalSets()).IsEqualTo(1)
	AssertThat(b.GetAliasingGets()).IsEqualTo(1)
}

// Java: GlobalNamespaceTest.testDoubleLhsCastInAssignment_doesNotCrash()
func (c *GlobalNamespaceTest) TestDoubleLhsCastInAssignment_doesNotCrash() {
	// The type of the cast doesn't matter.
	// Casting is only legal JS syntax in simple assignments, not with destructuring or declaration.
	namespace := c.Parse(
		`const ns = {};
 const b = 5;
 /** @type {*} */ (/** @type {*} */ (ns.a)) = b;
`)
	ns := globalNamespaceTest_name(namespace.GetSlot("ns"))
	AssertThat(ns.GetGlobalSets()).IsEqualTo(1)
	AssertThat(ns.GetTotalGets()).IsEqualTo(0)
	nsA := globalNamespaceTest_name(namespace.GetSlot("ns.a"))
	AssertThat(nsA.GetGlobalSets()).IsEqualTo(0)
	// TODO(b/127505242): Should be 1.
	AssertThat(nsA.GetTotalGets()).IsEqualTo(1)
	// TODO(b/127505242): Should be 0.
	b := globalNamespaceTest_name(namespace.GetSlot("b"))
	AssertThat(b.GetGlobalSets()).IsEqualTo(1)
	AssertThat(b.GetAliasingGets()).IsEqualTo(1)
}

// Java: GlobalNamespaceTest.testCannotCollapseAliasedObjectLitProperty()
func (c *GlobalNamespaceTest) TestCannotCollapseAliasedObjectLitProperty() {
	namespace := c.Parse("var foo = {prop: 0}; use(foo);")
	fooProp := globalNamespaceTest_name(namespace.GetSlot("foo.prop"))
	// We should not convert foo.prop -> foo$prop because use(foo) might read foo.prop
	AssertThat(fooProp.CanCollapse()).IsFalse()
}

// Java: GlobalNamespaceTest.testCannotCollapseConditionalObjectLitProperty()
func (c *GlobalNamespaceTest) TestCannotCollapseConditionalObjectLitProperty() {
	namespace := c.Parse("var foo = x || {prop: 0}; use(foo.prop);")
	fooProp := globalNamespaceTest_name(namespace.GetSlot("foo.prop"))
	// We should not convert foo.prop -> foo$prop because use(foo) might read foo.prop
	AssertThat(fooProp.CanCollapse()).IsFalse()
}

// Java: GlobalNamespaceTest.testCannotCollapseConditionalObjectLitNestedProperty()
func (c *GlobalNamespaceTest) TestCannotCollapseConditionalObjectLitNestedProperty() {
	namespace := c.Parse("var foo = x || {prop: {nested: 0}}; use(foo.prop.nested);")
	fooProp := globalNamespaceTest_name(namespace.GetSlot("foo.prop.nested"))
	// We should not convert foo.prop -> foo$prop because use(foo) might read foo.prop
	AssertThat(fooProp.CanCollapse()).IsFalse()
}

// Java: GlobalNamespaceTest.testCannotCollapseConditionalAndObjectLitProperty()
func (c *GlobalNamespaceTest) TestCannotCollapseConditionalAndObjectLitProperty() {
	namespace := c.Parse("var foo = x && {prop: 0}; use(foo.prop);")

	fooProp := globalNamespaceTest_name(namespace.GetSlot("foo.prop"))

	// We should not convert foo.prop -> foo$prop because use(foo) might read foo.prop
	AssertThat(fooProp.CanCollapse()).IsFalse()
}

// Java: GlobalNamespaceTest.testCannotCollapseConditionalAndObjectLitNestedProperty()
func (c *GlobalNamespaceTest) TestCannotCollapseConditionalAndObjectLitNestedProperty() {
	namespace := c.Parse("var foo = x && {prop: {nested: 0}}; use(foo.prop.nested);")

	fooProp := globalNamespaceTest_name(namespace.GetSlot("foo.prop.nested"))

	// We should not convert foo.prop -> foo$prop because use(foo) might read foo.prop
	AssertThat(fooProp.CanCollapse()).IsFalse()
}

// Java: GlobalNamespaceTest.testGitHubIssue3733()
func (c *GlobalNamespaceTest) TestGitHubIssue3733() {
	namespace := c.Parse(
		`const X = {Y: 1};

function fn(a) {
  if (a) {
// Before issue #3733 was fixed GlobalNamespace failed to see this reference
// as creating an alias for X due to a switch statement that failed to check
// for the RETURN node type, so X.Y was incorrectly collapsed.
    return a ? X : {};
  }
}

console.log(fn(true).Y);
`)
	nameX := globalNamespaceTest_name(namespace.GetSlot("X"))
	AssertThat(nameX.CanCollapseUnannotatedChildNames()).IsFalse()
	propY := globalNamespaceTest_name(namespace.GetSlot("X.Y"))
	AssertThat(propY.CanCollapse()).IsFalse()
}

// Java: GlobalNamespaceTest.testThrowPreventsCollapsingChildNames()
func (c *GlobalNamespaceTest) TestThrowPreventsCollapsingChildNames() {
	namespace := c.Parse(
		`const X = {Y: 1};

function fn(a) {
// This is specifically testing a bugfix closely related to GitHub issue
// #3733. A quirk of the implementation hides the bug when the throw isn't
// inside an if statement or the thrown value isn't a conditional expression.
  if (a) {
    throw a ? X : {};
  }
}

console.log(fn(true).Y);
`)
	nameX := globalNamespaceTest_name(namespace.GetSlot("X"))
	AssertThat(nameX.CanCollapseUnannotatedChildNames()).IsFalse()
	propY := globalNamespaceTest_name(namespace.GetSlot("X.Y"))
	AssertThat(propY.CanCollapse()).IsFalse()
}

// Java: GlobalNamespaceTest.testCannotCollapseObjectLitPropertyEscapedWithOptChainCall()
func (c *GlobalNamespaceTest) TestCannotCollapseObjectLitPropertyEscapedWithOptChainCall() {
	namespace := c.Parse("var foo = {prop: 0}; use?.(foo);")
	fooProp := globalNamespaceTest_name(namespace.GetSlot("foo.prop"))
	// We should not convert foo.prop -> foo$prop because use(foo) might read foo.prop
	AssertThat(fooProp.CanCollapse()).IsFalse()
}

// Java: GlobalNamespaceTest.testCanCollapseAliasedConstructorProperty_ignoringStaticInheritance()
func (c *GlobalNamespaceTest) TestCanCollapseAliasedConstructorProperty_ignoringStaticInheritance() {
	namespace := c.Parse(
		`/** @constructor */
var Foo = function() {}

Foo.prop = prop;
use(Foo);
`)
	fooProp := globalNamespaceTest_name(namespace.GetSlot("Foo.prop"))
	// We should still convert Foo.prop -> Foo$prop, even though use(Foo) might read Foo.prop,
	// because Foo is a constructor
	AssertThat(fooProp.CanCollapse()).IsTrue()
}

// Java: GlobalNamespaceTest.testCannotCollapseAliasedConstructorProperty_consideringStaticInheritance()
func (c *GlobalNamespaceTest) TestCannotCollapseAliasedConstructorProperty_consideringStaticInheritance() {
	c.assumeStaticInheritanceIsNotUsed = false
	namespace := c.Parse(
		`/** @constructor */
var Foo = function() {}

Foo.prop = prop;
use(Foo);
`)
	fooProp := globalNamespaceTest_name(namespace.GetSlot("Foo.prop"))
	AssertThat(fooProp.CanCollapse()).IsFalse()
}

// Java: GlobalNamespaceTest.testCanCollapseAliasedInterfaceProperty_ignoringStaticInheritance()
func (c *GlobalNamespaceTest) TestCanCollapseAliasedInterfaceProperty_ignoringStaticInheritance() {
	namespace := c.Parse(
		`/** @interface */
var Foo = function() {}

Foo.prop = prop;
use(Foo);
`)
	fooProp := globalNamespaceTest_name(namespace.GetSlot("Foo.prop"))
	// We should still convert Foo.prop -> Foo$prop, even though use(Foo) might read Foo.prop,
	// because Foo is a constructor
	AssertThat(fooProp.CanCollapse()).IsTrue()
}

// Java: GlobalNamespaceTest.testCannotCollapseAliasedInterfaceProperty_consideringStaticInheritance()
func (c *GlobalNamespaceTest) TestCannotCollapseAliasedInterfaceProperty_consideringStaticInheritance() {
	c.assumeStaticInheritanceIsNotUsed = false
	namespace := c.Parse(
		`/** @interface */
var Foo = function() {}

Foo.prop = prop;
use(Foo);
`)
	fooProp := globalNamespaceTest_name(namespace.GetSlot("Foo.prop"))
	AssertThat(fooProp.CanCollapse()).IsFalse()
}

// Java: GlobalNamespaceTest.testCanCollapseAliasedClassProperty_ignoringStaticInheritance()
func (c *GlobalNamespaceTest) TestCanCollapseAliasedClassProperty_ignoringStaticInheritance() {
	namespace := c.Parse("class Foo {} Foo.prop = prop; use(Foo);")
	fooProp := globalNamespaceTest_name(namespace.GetSlot("Foo.prop"))
	// We should still convert Foo.prop -> Foo$prop, even though use(Foo) might read Foo.prop,
	// because Foo is a constructor
	AssertThat(fooProp.CanCollapse()).IsTrue()
}

// Java: GlobalNamespaceTest.testCanCollapseAliasedClassProperty_consideringStaticInheritance()
func (c *GlobalNamespaceTest) TestCanCollapseAliasedClassProperty_consideringStaticInheritance() {
	c.assumeStaticInheritanceIsNotUsed = false
	namespace := c.Parse("class Foo {} Foo.prop = prop; use(Foo);")
	fooProp := globalNamespaceTest_name(namespace.GetSlot("Foo.prop"))
	AssertThat(fooProp.CanCollapse()).IsFalse()
}

// Java: GlobalNamespaceTest.testCannotCollapseOrInlineDeletedProperty()
func (c *GlobalNamespaceTest) TestCannotCollapseOrInlineDeletedProperty() {
	namespace := c.Parse(
		`const global = window;
delete global.HTMLElement;
global.HTMLElement = (class {});
`)
	deletedProp := globalNamespaceTest_name(namespace.GetSlot("global.HTMLElement"))
	AssertThat(deletedProp.CanCollapseOrInline()).IsEqualTo(GlobalNamespace_Inlinability_DO_NOT_INLINE)
}

// Java: GlobalNamespaceTest.testCanCollapse_objectLitProperty_declaredBeforeASpread()
func (c *GlobalNamespaceTest) TestCanCollapse_objectLitProperty_declaredBeforeASpread() {
	namespace := c.Parse("var foo = {prop: 0, ...bar}; use(foo.prop);")
	fooProp := globalNamespaceTest_name(namespace.GetSlot("foo.prop"))
	AssertThat(fooProp.CanCollapse()).IsFalse()
}

// Java: GlobalNamespaceTest.testGoogProvideName()
func (c *GlobalNamespaceTest) TestGoogProvideName() {
	namespace := c.Parse("goog.provide('a'); var a = {};")
	a := globalNamespaceTest_name(namespace.GetSlot("a"))
	AssertThat(a).IsNotNull()
	AssertThat(a.GetGlobalSets()).IsEqualTo(1)
	// The VAR, not the goog.provide, is considered the 'declaration' of `a`.
	AssertNode(a.getDeclarationRef().GetNode().GetParent()).HasToken(rhino.Token_VAR)
}

// Java: GlobalNamespaceTest.testGoogProvideNamespace_noExplicitAssignment()
func (c *GlobalNamespaceTest) TestGoogProvideNamespace_noExplicitAssignment() {
	namespace := c.Parse("goog.provide('a.b');")
	a := globalNamespaceTest_name(namespace.GetSlot("a"))
	AssertThat(a).IsNotNull()
	AssertThat(a.GetGlobalSets()).IsEqualTo(0)
	ab := globalNamespaceTest_name(namespace.GetSlot("a.b"))
	AssertThat(ab).IsNotNull()
	AssertThat(ab.GetGlobalSets()).IsEqualTo(0)
	AssertThat(a.getDeclarationRef()).IsNull()
	AssertThat(ab.getDeclarationRef()).IsNull()
	AssertThat(ab.GetParent()).IsEqualTo(a)
}

// Java: GlobalNamespaceTest.testGoogProvideLongNamespace()
func (c *GlobalNamespaceTest) TestGoogProvideLongNamespace() {
	namespace := c.Parse("goog.provide('a.b.c.d');")
	AssertThat(globalNamespaceTest_name(namespace.GetSlot("a.b.c.d"))).IsNotNull()
}

// Java: GlobalNamespaceTest.testGoogProvideNamespace_explicitAssignment()
func (c *GlobalNamespaceTest) TestGoogProvideNamespace_explicitAssignment() {
	namespace := c.Parse("goog.provide('a.b'); /** @const */ a.b = {};")
	a := globalNamespaceTest_name(namespace.GetSlot("a"))
	AssertThat(a).IsNotNull()
	AssertThat(a.GetGlobalSets()).IsEqualTo(0)
	ab := globalNamespaceTest_name(namespace.GetSlot("a.b"))
	AssertThat(ab).IsNotNull()
	AssertThat(ab.GetGlobalSets()).IsEqualTo(1)
}

// Java: GlobalNamespaceTest.testGoogProvideNamespace_assignmentToProperty()
func (c *GlobalNamespaceTest) TestGoogProvideNamespace_assignmentToProperty() {
	namespace := c.Parse("goog.provide('a.b'); a.b.Class = class {};")
	abClass := globalNamespaceTest_name(namespace.GetSlot("a.b.Class"))
	AssertThat(abClass).IsNotNull()
	AssertThat(abClass.GetGlobalSets()).IsEqualTo(1)
	AssertThat(abClass.GetParent()).IsEqualTo(globalNamespaceTest_name(namespace.GetSlot("a.b")))
}

// Java: GlobalNamespaceTest.testGoogProvideName_multipleProvidesForName()
func (c *GlobalNamespaceTest) TestGoogProvideName_multipleProvidesForName() {
	namespace := c.Parse("goog.provide('a.b'); goog.provide('a.c');")
	a := globalNamespaceTest_name(namespace.GetSlot("a"))
	AssertThat(a).IsNotNull()
	AssertThat(a.GetGlobalSets()).IsEqualTo(0)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.googModuleLevelNamesAreCaptured()
func (c *GlobalNamespaceTest) TestGoogModuleLevelNamesAreCaptured() {
	namespace := c.ParseAndGatherModuleData("goog.module('m'); const x = 0;")
	metadata := c.lastCompiler.GetModuleMetadataMap().GetModulesByGoogNamespace().Get("m")
	x := namespace.GetNameFromModule(metadata, "x")
	AssertThat(x).IsNotNull()
	AssertThat(x.getDeclarationRef()).IsNotNull()
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.googModuleLevelQualifiedNamesAreCaptured()
func (c *GlobalNamespaceTest) TestGoogModuleLevelQualifiedNamesAreCaptured() {
	namespace := c.ParseAndGatherModuleData("goog.module('m'); class Foo {} Foo.Bar = 0;")
	metadata := c.lastCompiler.GetModuleMetadataMap().GetModulesByGoogNamespace().Get("m")
	x := namespace.GetNameFromModule(metadata, "Foo.Bar")
	AssertThat(x).IsNotNull()
	AssertThat(x.getDeclarationRef()).IsNotNull()
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.googModule_containsExports()
func (c *GlobalNamespaceTest) TestGoogModule_containsExports() {
	namespace := c.ParseAndGatherModuleData("goog.module('m'); const x = 0;")
	metadata := c.lastCompiler.GetModuleMetadataMap().GetModulesByGoogNamespace().Get("m")
	exports := namespace.GetNameFromModule(metadata, "exports")
	AssertThat(exports.GetGlobalSets()).IsEqualTo(0)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.googLoadModule_containsExports()
func (c *GlobalNamespaceTest) TestGoogLoadModule_containsExports() {
	namespace := c.ParseAndGatherModuleData(
		`goog.loadModule(function(exports) {
  goog.module('m');
  const x = 0;
  return exports;
});
`)
	metadata := c.lastCompiler.GetModuleMetadataMap().GetModulesByGoogNamespace().Get("m")
	exports := namespace.GetNameFromModule(metadata, "exports")
	AssertThat(exports.GetGlobalSets()).IsEqualTo(0)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.googLoadModule_capturesQualifiedNames()
func (c *GlobalNamespaceTest) TestGoogLoadModule_capturesQualifiedNames() {
	namespace := c.ParseAndGatherModuleData(
		`goog.loadModule(function(exports) {
  goog.module('m');
  class Foo {}
  Foo.Bar = class {};
  return exports;
});
`)
	metadata := c.lastCompiler.GetModuleMetadataMap().GetModulesByGoogNamespace().Get("m")
	foo := namespace.GetNameFromModule(metadata, "Foo")
	fooBar := namespace.GetNameFromModule(metadata, "Foo.Bar")
	AssertThat(fooBar.GetParent()).IsEqualTo(foo)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.googLoadModule_containsExportsPropertyAssignments()
func (c *GlobalNamespaceTest) TestGoogLoadModule_containsExportsPropertyAssignments() {
	namespace := c.ParseAndGatherModuleData(
		`goog.loadModule(function(exports) {
  goog.module('m');
  exports.Foo = class {};
  return exports;
});
`)
	metadata := c.lastCompiler.GetModuleMetadataMap().GetModulesByGoogNamespace().Get("m")
	exportsFoo := namespace.GetNameFromModule(metadata, "exports.Foo")
	AssertThat(exportsFoo.GetGlobalSets()).IsEqualTo(1)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.googModule_containsExports_explicitAssign()
func (c *GlobalNamespaceTest) TestGoogModule_containsExports_explicitAssign() {
	namespace := c.ParseAndGatherModuleData("goog.module('m'); const x = 0; exports = {x};")
	metadata := c.lastCompiler.GetModuleMetadataMap().GetModulesByGoogNamespace().Get("m")
	exports := namespace.GetNameFromModule(metadata, "exports")
	AssertThat(exports.GetGlobalSets()).IsEqualTo(1)
	AssertThat(namespace.GetNameFromModule(metadata, "x").GetGlobalSets()).IsEqualTo(1)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.assignToGlobalNameInLoadModule_doesNotCreateModuleName()
func (c *GlobalNamespaceTest) TestAssignToGlobalNameInLoadModule_doesNotCreateModuleName() {
	namespace := c.ParseAndGatherModuleData(
		`class Foo {}
goog.loadModule(function(exports) {
  goog.module('m');
  Foo.Bar = 0
  return exports;
});
`)
	metadata := c.lastCompiler.GetModuleMetadataMap().GetModulesByGoogNamespace().Get("m")
	AssertThat(namespace.GetNameFromModule(metadata, "Foo.Bar")).IsNull()
	AssertThat(globalNamespaceTest_name(namespace.GetSlot("Foo.Bar"))).IsNotNull()
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.moduleLevelNamesAreCaptured_esExportDecl()
func (c *GlobalNamespaceTest) TestModuleLevelNamesAreCaptured_esExportDecl() {
	namespace := c.ParseAndGatherModuleData("export const x = 0;")
	metadata := c.lastCompiler.GetModuleMetadataMap().GetModulesByPath().Get("test.js")
	x := namespace.GetNameFromModule(metadata, "x")
	AssertThat(x).IsNotNull()
	AssertNode(x.getDeclarationRef().GetNode().GetParent()).HasToken(rhino.Token_CONST)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.moduleLevelNamesAreCaptured_esExportClassDecl()
func (c *GlobalNamespaceTest) TestModuleLevelNamesAreCaptured_esExportClassDecl() {
	namespace := c.ParseAndGatherModuleData("export class Foo {}")
	metadata := c.lastCompiler.GetModuleMetadataMap().GetModulesByPath().Get("test.js")
	x := namespace.GetNameFromModule(metadata, "Foo")
	AssertThat(x).IsNotNull()
	AssertNode(x.getDeclarationRef().GetNode().GetParent()).HasToken(rhino.Token_CLASS)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.moduleLevelNamesAreCaptured_esExportFunctionDecl()
func (c *GlobalNamespaceTest) TestModuleLevelNamesAreCaptured_esExportFunctionDecl() {
	namespace := c.ParseAndGatherModuleData("export function fn() {}")
	metadata := c.lastCompiler.GetModuleMetadataMap().GetModulesByPath().Get("test.js")
	x := namespace.GetNameFromModule(metadata, "fn")
	AssertThat(x).IsNotNull()
	AssertNode(x.getDeclarationRef().GetNode().GetParent()).HasToken(rhino.Token_FUNCTION)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.moduleLevelNamesAreCaptured_esExportDefaultFunctionDecl()
func (c *GlobalNamespaceTest) TestModuleLevelNamesAreCaptured_esExportDefaultFunctionDecl() {
	namespace := c.ParseAndGatherModuleData("export default function fn() {}")
	metadata := c.lastCompiler.GetModuleMetadataMap().GetModulesByPath().Get("test.js")
	x := namespace.GetNameFromModule(metadata, "fn")
	AssertThat(x).IsNotNull()
	AssertNode(x.getDeclarationRef().GetNode().GetParent()).HasToken(rhino.Token_FUNCTION)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.esModuleLevelNamesAreCaptured()
func (c *GlobalNamespaceTest) TestEsModuleLevelNamesAreCaptured() {
	namespace := c.ParseAndGatherModuleData("class Foo {} Foo.Bar = 0; export {Foo};")
	metadata := c.lastCompiler.GetModuleMetadataMap().GetModulesByPath().Get("test.js")
	x := namespace.GetNameFromModule(metadata, "Foo.Bar")
	AssertThat(x).IsNotNull()
	AssertThat(x.getDeclarationRef()).IsNotNull()
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.getCommonAncestorChunk_returnsDeclarationChunk_whenDeclarationLoadedFirst()
func (c *GlobalNamespaceTest) TestGetCommonAncestorChunk_returnsDeclarationChunk_whenDeclarationLoadedFirst() {
	chunks :=
		JSChunkGraphBuilder_ForChain().
			AddChunk("console.log('base');").
			AddChunk("const parent = {};").
			AddChunk("parent.child = {};").
			Build()
	namespace := c.ParseChunks(chunks)
	parentName := globalNamespaceTest_name(namespace.GetSlot("parent"))
	childName := globalNamespaceTest_name(namespace.GetSlot("parent.child"))
	chunkGraph := c.lastCompiler.GetChunkGraph()
	AssertThat(parentName.GetDeepestCommonAncestorChunk(chunkGraph)).IsEqualTo(chunks[1])
	AssertThat(childName.GetDeepestCommonAncestorChunk(chunkGraph)).IsEqualTo(chunks[2])
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: GlobalNamespaceTest.getCommonAncestorChunk_findsCommonAncestorOfSiblingChunks()
func (c *GlobalNamespaceTest) TestGetCommonAncestorChunk_findsCommonAncestorOfSiblingChunks() {
	chunks :=
		JSChunkGraphBuilder_ForBush().
			AddChunk("console.log('base');").
			AddChunk("const parent = {};").      // parent depends on base
			AddChunk("parent.crossChunk = {};"). // depends on parent
			AddChunk(
				"if (parent.crossChunk) { console.log(parent.crossChunk ); }"). // depends on parent
			Build()
	namespace := c.ParseChunks(chunks)
	crossChunkName := globalNamespaceTest_name(namespace.GetSlot("parent.crossChunk"))
	chunkGraph := c.lastCompiler.GetChunkGraph()
	AssertThat(crossChunkName.GetDeepestCommonAncestorChunk(chunkGraph)).IsEqualTo(chunks[1])
	AssertThat(crossChunkName.getDeclarationRef().GetChunk()).IsEqualTo(chunks[2])
}

// ParseAndGatherModuleData exists for testing module metadata lookups.
//
// Java: GlobalNamespaceTest.parseAndGatherModuleData(String)
func (c *GlobalNamespaceTest) ParseAndGatherModuleData(js string) *GlobalNamespace {
	options := c.GetDefaultOptions()
	compiler := c.Compile(js, options)
	// Disabling transpilation also disables these passes that we need to have run when
	// testing behavior related to module metadata.
	NewGatherModuleMetadata(
		compiler, options.GetProcessCommonJSModules(), options.GetModuleResolutionMode()).
		Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	NewModuleMapCreator(compiler, compiler.GetModuleMetadataMap()).
		Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	AssertThat(compiler.GetErrors()).IsEmpty()
	c.lastCompiler = compiler
	return NewGlobalNamespace(compiler, compiler.GetRoot())
}

// Java: GlobalNamespaceTest.parse(String)
func (c *GlobalNamespaceTest) Parse(js string) *GlobalNamespace {
	options := c.GetDefaultOptions()
	c.Compile(js, options)
	return NewGlobalNamespace(c.lastCompiler, c.lastCompiler.GetRoot())
}

// Java: GlobalNamespaceTest.parse(JSChunk[])
func (c *GlobalNamespaceTest) ParseChunks(chunks []*JSChunk) *GlobalNamespace {
	options := c.GetDefaultOptions()
	compiler := NewCompiler()
	result := compiler.CompileChunks([]*SourceFile{}, slices.Clone(chunks), options)
	AssertThat(compiler.GetErrors()).IsEmpty()
	AssertThat(result.Success).IsTrue()
	c.lastCompiler = compiler
	return NewGlobalNamespace(c.lastCompiler, c.lastCompiler.GetRoot())
}

// Java: GlobalNamespaceTest.getDefaultOptions()
func (c *GlobalNamespaceTest) GetDefaultOptions() *CompilerOptions {
	options := NewCompilerOptions()
	options.SetLanguage(CompilerOptions_LanguageMode_UNSUPPORTED)
	// Don't optimize, because we want to know how GlobalNamespace responds to the original code
	// in `js`.
	options.SetSkipNonTranspilationPasses(true)
	options.SetWrapGoogModulesForWhitespaceOnly(false)
	// Test the latest features supported for input and don't transpile, because we want to test how
	// GlobalNamespace deals with the language features actually present in `js`.
	options.SetAssumeStaticInheritanceIsNotUsed(c.assumeStaticInheritanceIsNotUsed)
	return options
}

// Java: GlobalNamespaceTest.compile(String, CompilerOptions)
func (c *GlobalNamespaceTest) Compile(js string, options *CompilerOptions) *Compiler {
	compiler := NewCompiler()
	compiler.CompileOfSourceFiles(SourceFile_FromCode("ex.js", ""), SourceFile_FromCode("test.js", js), options)
	AssertThat(compiler.GetErrors()).IsEmpty()
	c.lastCompiler = compiler
	return compiler
}

// globalNamespaceTest_name is the Name a GlobalNamespace.getSlot /
// getOwnSlot returns (Java narrows StaticScope's StaticSlot to Name).
//
// Go-only: language-adapter — GlobalNamespace.GetSlot keeps
// rhino.StaticScope's signature.
func globalNamespaceTest_name(slot rhino.StaticSlot) *GlobalNamespace_Name {
	name, _ := slot.(*GlobalNamespace_Name)
	return name
}

// globalNamespaceTest_setOf is ImmutableSet.of(change).
func globalNamespaceTest_setOf(change GlobalNamespace_AstChange) *astChangeSet {
	var set astChangeSet
	set.add(change)
	return &set
}
