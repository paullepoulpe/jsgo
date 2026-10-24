// Port of com/google/javascript/jscomp/disambiguate/ColorGraphBuilderTest.java @ closure-compiler v20261015

/*
 * Copyright 2019 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// The Java @Test methods are not named test*; their Go names get a "Test"
// prefix so that RunJUnit4 runs them.
//
// Java: public final class ColorGraphBuilderTest extends CompilerTestCase
type ColorGraphBuilderTest struct {
	CompilerTestCaseBase

	compiler *Compiler

	// This registry is only used to lookup box colors so it doesn't have to be the same one used by
	// the graph builder.
	graphNodeFactory *ColorGraphNodeFactory

	processor CompilerPass // @Nullable
	result    DiGraph[*ColorGraphNode, any]
	labelToId *jutil.LinkedHashMap[string, colors.ColorId]
}

// Java: ColorGraphBuilderTest()
func NewColorGraphBuilderTest() *ColorGraphBuilderTest {
	c := &ColorGraphBuilderTest{}
	c.InitCompilerTestCaseBase(c)
	c.compiler = NewCompiler()
	c.graphNodeFactory =
		ColorGraphNodeFactory_CreateFactory(
			colors.ColorRegistry_Builder_().SetDefaultNativeColorsForTesting().Build())
	return c
}

func TestColorGraphBuilderTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewColorGraphBuilderTest() })
}

// Java: ColorGraphBuilderTest.createCompiler()
func (c *ColorGraphBuilderTest) CreateCompiler() *Compiler {
	return c.compiler
}

// Java: ColorGraphBuilderTest.getProcessor(Compiler)
func (c *ColorGraphBuilderTest) GetProcessor(compiler *Compiler) CompilerPass {
	AssertThat(compiler).IsSameInstanceAs(c.compiler)
	return c.processor
}

// Java: ColorGraphBuilderTest.setUp()
func (c *ColorGraphBuilderTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()

	c.EnableTypeCheck()
	c.ReplaceTypesWithColors()
	c.DisableCompareJsDoc() // type -> color conversion also erases JSDoc
}

// TearDown runs the class's @After methods, then CompilerTestCase's
// (JUnit runs a subclass's @After methods before its superclass's).
func (c *ColorGraphBuilderTest) TearDown() {
	c.VerifyResult_topNodeIsOnlyRoot()
	c.VerifyResult_hasNoSelfEdges()
	c.VerifyResult_hasNoParallelEdges()
	c.RenderResultGraph()
	c.CompilerTestCaseBase.TearDown()
}

// Java: ColorGraphBuilderTest.createBuilder(StubLcaFinder, ColorRegistry)
func (c *ColorGraphBuilderTest) createBuilder(
	optLcaFinder *colorGraphBuilderTest_StubLcaFinder, registry *colors.ColorRegistry) *ColorGraphBuilder {
	lcaFinder := optLcaFinder
	if lcaFinder == nil {
		lcaFinder = newColorGraphBuilderTest_StubLcaFinder()
	}
	return NewColorGraphBuilder(
		c.graphNodeFactory,
		func(graph DiGraph[*ColorGraphNode, any]) LowestCommonAncestorFinder[*ColorGraphNode, any] {
			return lcaFinder.setGraph(graph)
		},
		registry)
}

// colorGraphBuilderTest_collector is the anonymous AbstractPostOrderCallback
// of createBuilderIncludingCode.
type colorGraphBuilderTest_collector struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	graphNodeFactory *ColorGraphNodeFactory
	testTypes        *jutil.LinkedHashMap[string, *ColorGraphNode]
	labelToId        *jutil.LinkedHashMap[string, colors.ColorId]
}

// Java: ColorGraphBuilderTest.createBuilderIncludingCode(...).new AbstractPostOrderCallback() {...}.visit(NodeTraversal, Node, Node)
func (cb *colorGraphBuilderTest_collector) Visit(t *NodeTraversal, n rhino.Node, unused rhino.Node) {
	if n.IsName() && jutil.StringStartsWith(n.GetString(), "test") {
		cb.testTypes.Put(n.GetString(), cb.graphNodeFactory.CreateNode(n.GetColor()))
	}

	if n.IsLabel() {
		labelName := n.GetFirstChild().GetString()
		var labeledExpr rhino.Node
		if n.GetSecondChild().IsExprResult() {
			labeledExpr = n.GetSecondChild().GetOnlyChild()
		} else {
			labeledExpr = n.GetSecondChild()
		}
		cb.labelToId.Put(labelName, labeledExpr.GetColor().GetId())
	}
}

// Java: ColorGraphBuilderTest.createBuilderIncludingCode(Supplier<StubLcaFinder>, String)
func (c *ColorGraphBuilderTest) createBuilderIncludingCode(
	// takes a Supplier<StubLcaFinder> so that callers may lazily create this StubLcaFinder after
	// the labelToId map has been populated.
	optLcaFinder func() *colorGraphBuilderTest_StubLcaFinder, src string) *ColorGraphBuilder {
	graphNodeFactory := c.graphNodeFactory
	testTypes := jutil.NewLinkedHashMap[string, *ColorGraphNode]()
	c.labelToId = jutil.NewLinkedHashMap[string, colors.ColorId]()

	/* Flatten and collect the types of all NAMEs that start with "test". */
	c.processor = CompilerPassFunc(func(externs rhino.Node, main rhino.Node) {
		NodeTraversal_Traverse(
			c.compiler,
			main,
			&colorGraphBuilderTest_collector{
				graphNodeFactory: graphNodeFactory, testTypes: testTypes, labelToId: c.labelToId})
	})
	c.TestSame(c.Srcs(src))

	c.processor = nil

	var lcaFinder *colorGraphBuilderTest_StubLcaFinder
	if optLcaFinder != nil {
		lcaFinder = optLcaFinder()
	}
	graphBuilder := c.createBuilder(lcaFinder, c.compiler.GetColorRegistry())
	graphBuilder.AddAll(testTypes.Values())
	return graphBuilder
}

// Java: ColorGraphBuilderTest.top_isAboveInterface()
func (c *ColorGraphBuilderTest) TestTop_isAboveInterface() {
	// Given
	builder :=
		c.createBuilderIncludingCode(
			nil,
			`/** @interface */
class IFoo { }

let /** !IFoo */ test;

IFOO_PROTOTYPE: IFoo.prototype
IFOO: test;
`)

	// When
	c.result = builder.Build()

	// Then
	c.assertThatResultAsTable().
		ContainsCell(colors.StandardColors_TOP_OBJECT.GetId(), c.id("IFOO_PROTOTYPE"), ColorGraphBuilder_EdgeReason_CAN_HOLD)
	c.assertThatResultAsTable().ContainsCell(c.id("IFOO_PROTOTYPE"), c.id("IFOO"), ColorGraphBuilder_EdgeReason_CAN_HOLD)
}

// Java: ColorGraphBuilderTest.prototypeChain_isInserted()
func (c *ColorGraphBuilderTest) TestPrototypeChain_isInserted() {
	// Given
	builder :=
		c.createBuilderIncludingCode(
			nil,
			`class Foo { }

const test = new Foo();

FOO_PROTOTYPE: Foo.prototype;
FOO: new Foo();
`)

	// When
	c.result = builder.Build()

	// Then
	c.assertThatResultAsTable().
		ContainsCell(colors.StandardColors_TOP_OBJECT.GetId(), c.id("FOO_PROTOTYPE"), ColorGraphBuilder_EdgeReason_CAN_HOLD)
	c.assertThatResultAsTable().ContainsCell(c.id("FOO_PROTOTYPE"), c.id("FOO"), ColorGraphBuilder_EdgeReason_CAN_HOLD)
}

// Java: ColorGraphBuilderTest.prototypeChain_isInserted_throughExtends()
func (c *ColorGraphBuilderTest) TestPrototypeChain_isInserted_throughExtends() {
	// Given
	builder :=
		c.createBuilderIncludingCode(
			nil,
			`class Foo { }
class Bar extends Foo { }

const test = new Bar();

FOO_PROTOTYPE: Foo.prototype;
FOO: new Foo();
BAR_PROTOTYPE: Bar.prototype;
FOO_PROTOTYPE: Foo.prototype;
`)

	// When
	c.result = builder.Build()

	// Then
	c.assertThatResultAsTable().
		ContainsCell(colors.StandardColors_TOP_OBJECT.GetId(), c.id("FOO_PROTOTYPE"), ColorGraphBuilder_EdgeReason_CAN_HOLD)
	c.assertThatResultAsTable().ContainsCell(c.id("FOO_PROTOTYPE"), c.id("FOO"), ColorGraphBuilder_EdgeReason_CAN_HOLD)
	c.assertThatResultAsTable().ContainsCell(c.id("FOO"), c.id("BAR_PROTOTYPE"), ColorGraphBuilder_EdgeReason_CAN_HOLD)
}

// Java: ColorGraphBuilderTest.prototypeChain_canBranch()
func (c *ColorGraphBuilderTest) TestPrototypeChain_canBranch() {
	// Given
	builder :=
		c.createBuilderIncludingCode(
			nil,
			`class Foo { }
class Bar extends Foo { }
class Qux extends Foo { }

const testBar = new Bar();
const testQux = new Qux();

FOO_PROTOTYPE: Foo.prototype;
FOO: new Foo();
BAR_PROTOTYPE: Bar.prototype;
BAR: new Bar();
QUX_PROTOTYPE: Qux.prototype;
QUX: new Qux();
`)

	// When
	c.result = builder.Build()

	// Then
	c.assertThatResultAsTable().ContainsCell(c.id("FOO_PROTOTYPE"), c.id("FOO"), ColorGraphBuilder_EdgeReason_CAN_HOLD)
	c.assertThatResultAsTable().ContainsCell(c.id("FOO"), c.id("BAR_PROTOTYPE"), ColorGraphBuilder_EdgeReason_CAN_HOLD)
	c.assertThatResultAsTable().ContainsCell(c.id("BAR_PROTOTYPE"), c.id("BAR"), ColorGraphBuilder_EdgeReason_CAN_HOLD)
	c.assertThatResultAsTable().ContainsCell(c.id("FOO"), c.id("QUX_PROTOTYPE"), ColorGraphBuilder_EdgeReason_CAN_HOLD)
	c.assertThatResultAsTable().ContainsCell(c.id("QUX_PROTOTYPE"), c.id("QUX"), ColorGraphBuilder_EdgeReason_CAN_HOLD)
}

// Java: ColorGraphBuilderTest.constructorDef_includesPrototypeAndInstanceType_evenIfUnused()
func (c *ColorGraphBuilderTest) TestConstructorDef_includesPrototypeAndInstanceType_evenIfUnused() {
	// Given
	builder :=
		c.createBuilderIncludingCode(
			nil,
			`class Foo { }

const test = Foo;

FOO_PROTOTYPE: Foo.prototype;
FOO: new Foo();
BAR_PROTOTYPE: Bar.prototype;
`)

	// When
	c.result = builder.Build()

	// Then
	c.assertThatResultAsTable().
		ContainsCell(colors.StandardColors_TOP_OBJECT.GetId(), c.id("FOO_PROTOTYPE"), ColorGraphBuilder_EdgeReason_CAN_HOLD)
	c.assertThatResultAsTable().ContainsCell(c.id("FOO_PROTOTYPE"), c.id("FOO"), ColorGraphBuilder_EdgeReason_CAN_HOLD)
}

// Java: ColorGraphBuilderTest.prototypeChain_connectsInterfaces()
func (c *ColorGraphBuilderTest) TestPrototypeChain_connectsInterfaces() {
	// Given
	builder :=
		c.createBuilderIncludingCode(
			nil,
			`/** @interface */
class IFoo { }

let /** !IFoo */ test;

IFOO: test;
IFOO_PROTOTYPE: IFoo.prototype;
`)

	// When
	c.result = builder.Build()

	// Then
	c.assertThatResultAsTable().ContainsCell(c.id("IFOO_PROTOTYPE"), c.id("IFOO"), ColorGraphBuilder_EdgeReason_CAN_HOLD)
}

// Java: ColorGraphBuilderTest.prototypeChain_connectsSubclasses_viaClassSideInheritance()
func (c *ColorGraphBuilderTest) TestPrototypeChain_connectsSubclasses_viaClassSideInheritance() {
	// Given
	builder :=
		c.createBuilderIncludingCode(
			nil,
			`/** @constructor */ function Foo0 () { }
class Foo1 extends Foo0 { }
class Foo2 extends Foo1 { }

let /** !(typeof Foo2) */ test;
FOO0_CTOR: Foo0;
FOO1_CTOR: Foo1;
FOO2_CTOR: Foo2;
`)

	// When
	c.result = builder.Build()

	// Then
	c.assertThatResultAsTable().
		ContainsCell(colors.StandardColors_TOP_OBJECT.GetId(), c.id("FOO0_CTOR"), ColorGraphBuilder_EdgeReason_CAN_HOLD)
	c.assertThatResultAsTable().ContainsCell(c.id("FOO0_CTOR"), c.id("FOO1_CTOR"), ColorGraphBuilder_EdgeReason_CAN_HOLD)
	c.assertThatResultAsTable().ContainsCell(c.id("FOO1_CTOR"), c.id("FOO2_CTOR"), ColorGraphBuilder_EdgeReason_CAN_HOLD)
}

// Java: ColorGraphBuilderTest.disambiguationSupertypes_createConnection()
func (c *ColorGraphBuilderTest) TestDisambiguationSupertypes_createConnection() {
	// Given
	parent := colorGraphBuilderTest_colorWithId(colors.ColorId_FromUnsigned(100)).Build()
	child := colorGraphBuilderTest_colorWithId(colors.ColorId_FromUnsigned(101)).Build()

	builder :=
		c.createBuilder(
			nil,
			colors.ColorRegistry_Builder_().
				SetDefaultNativeColorsForTesting().
				AddDisambiguationEdge(child, parent).
				Build())

	builder.Add(c.graphNodeFactory.CreateNode(child))

	// When
	c.result = builder.Build()

	// Then
	c.assertThatResultAsTable().ContainsCell(parent.GetId(), child.GetId(), ColorGraphBuilder_EdgeReason_CAN_HOLD)
}

// Java: ColorGraphBuilderTest.unions_connectedAboveMembers()
func (c *ColorGraphBuilderTest) TestUnions_connectedAboveMembers() {
	// Given
	stubFinder :=
		func() *colorGraphBuilderTest_StubLcaFinder {
			return newColorGraphBuilderTest_StubLcaFinder().
				addStub(
					jutil.NewLinkedHashSet(c.id("FOO"), c.id("BAR"), c.id("QUX")),
					jutil.NewLinkedHashSet(c.ids("FOO", "BAR", "QUX"), colors.StandardColors_UNKNOWN.GetId()))
		}
	builder :=
		c.createBuilderIncludingCode(
			stubFinder,
			`class Foo { }
class Bar { }
class Qux { }

let /** (!Foo|!Bar|!Qux) */ test;

FOO: new Foo();
BAR: new Bar();
QUX: new Qux();
`)

	// When
	c.result = builder.Build()

	// Then
	c.assertThatResultAsTable().ContainsCell(c.ids("FOO", "BAR", "QUX"), c.id("BAR"), ColorGraphBuilder_EdgeReason_ALGEBRAIC)
	c.assertThatResultAsTable().ContainsCell(c.ids("FOO", "BAR", "QUX"), c.id("FOO"), ColorGraphBuilder_EdgeReason_ALGEBRAIC)
	c.assertThatResultAsTable().ContainsCell(c.ids("FOO", "BAR", "QUX"), c.id("QUX"), ColorGraphBuilder_EdgeReason_ALGEBRAIC)
}

// Java: ColorGraphBuilderTest.unions_connectBelowLca()
func (c *ColorGraphBuilderTest) TestUnions_connectBelowLca() {
	// Given
	kifId := colors.ColorId_FromUnsigned(100)
	stubFinder :=
		func() *colorGraphBuilderTest_StubLcaFinder {
			return newColorGraphBuilderTest_StubLcaFinder().
				addStub(
					jutil.NewLinkedHashSet(c.id("FOO"), c.id("BAR"), c.id("QUX")),
					jutil.NewLinkedHashSet(c.ids("FOO", "BAR", "QUX"), kifId))
		}
	builder :=
		c.createBuilderIncludingCode(
			stubFinder,
			`class Foo { }
class Bar { }
class Qux { }

let /** (!Foo|!Bar|!Qux) */ test;

FOO: new Foo();
BAR: new Bar();
QUX: new Qux();
`)
	builder.Add(c.graphNodeFactory.CreateNode(colorGraphBuilderTest_colorWithId(kifId).Build()))

	// When
	c.result = builder.Build()

	// Then
	c.assertThatResultAsTable().ContainsCell(kifId, c.ids("FOO", "BAR", "QUX"), ColorGraphBuilder_EdgeReason_ALGEBRAIC)
}

// Java: ColorGraphBuilderTest.unions_connectBelowLca_withMultipleLcas()
func (c *ColorGraphBuilderTest) TestUnions_connectBelowLca_withMultipleLcas() {
	// Given
	kixId := colors.ColorId_FromUnsigned(100)
	lopId := colors.ColorId_FromUnsigned(101)
	stubFinder :=
		func() *colorGraphBuilderTest_StubLcaFinder {
			return newColorGraphBuilderTest_StubLcaFinder().
				addStub(
					jutil.NewLinkedHashSet(c.id("FOO"), c.id("BAR"), c.id("QUX")),
					jutil.NewLinkedHashSet(c.ids("FOO", "BAR", "QUX"), kixId, lopId))
		}
	builder :=
		c.createBuilderIncludingCode(
			stubFinder,
			`class Foo { }
class Bar { }
class Qux { }

let /** (!Foo|!Bar|!Qux) */ test;

FOO: new Foo();
BAR: new Bar();
QUX: new Qux();
`)
	builder.Add(c.graphNodeFactory.CreateNode(colorGraphBuilderTest_colorWithId(kixId).Build()))
	builder.Add(c.graphNodeFactory.CreateNode(colorGraphBuilderTest_colorWithId(lopId).Build()))

	// When
	c.result = builder.Build()

	// Then
	c.assertThatResultAsTable().ContainsCell(kixId, c.ids("BAR", "FOO", "QUX"), ColorGraphBuilder_EdgeReason_ALGEBRAIC)
	c.assertThatResultAsTable().ContainsCell(lopId, c.ids("BAR", "FOO", "QUX"), ColorGraphBuilder_EdgeReason_ALGEBRAIC)
}

// Java: ColorGraphBuilderTest.unions_connectBelowLca_whichIsAlsoUnion()
func (c *ColorGraphBuilderTest) TestUnions_connectBelowLca_whichIsAlsoUnion() {
	// Given
	kifId := colors.ColorId_FromUnsigned(100)
	stubFinder :=
		func() *colorGraphBuilderTest_StubLcaFinder {
			return newColorGraphBuilderTest_StubLcaFinder().
				addStub(
					jutil.NewLinkedHashSet(c.id("FOO"), c.id("BAR"), c.id("QUX")),
					jutil.NewLinkedHashSet(c.ids("FOO", "BAR", "QUX"), kifId)).
				addStub(
					jutil.NewLinkedHashSet(c.id("FOO"), c.id("BAR")),
					jutil.NewLinkedHashSet(c.ids("FOO", "BAR"), c.ids("FOO", "BAR", "QUX")))
		}
	builder :=
		c.createBuilderIncludingCode(
			stubFinder,
			`class Foo { }
class Bar { }
class Qux { }

let /** (!Foo|!Bar|!Qux) */ testA;
let /** (!Foo|!Bar) */ testB;

FOO: new Foo();
BAR: new Bar();
QUX: new Qux();
`)
	builder.Add(c.graphNodeFactory.CreateNode(colorGraphBuilderTest_colorWithId(kifId).Build()))

	// When
	c.result = builder.Build()

	// Then
	c.assertThatResultAsTable().ContainsCell(kifId, c.ids("FOO", "BAR", "QUX"), ColorGraphBuilder_EdgeReason_ALGEBRAIC)
	c.assertThatResultAsTable().
		ContainsCell(c.ids("FOO", "BAR", "QUX"), c.ids("FOO", "BAR"), ColorGraphBuilder_EdgeReason_ALGEBRAIC)
}

// Java: ColorGraphBuilderTest.unions_connectBelowLac_whichHasSameDescendantCount()
func (c *ColorGraphBuilderTest) TestUnions_connectBelowLac_whichHasSameDescendantCount() {
	// Given
	stubFinder :=
		func() *colorGraphBuilderTest_StubLcaFinder {
			return newColorGraphBuilderTest_StubLcaFinder().
				addStub(
					jutil.NewLinkedHashSet(c.id("FOO_PROTOTYPE"), c.id("BAR_PROTOTYPE")),
					jutil.NewLinkedHashSet(c.ids("FOO_PROTOTYPE", "BAR_PROTOTYPE"), c.id("KIF")))
		}
	builder :=
		c.createBuilderIncludingCode(
			stubFinder,
			`class Kif { }
class Foo extends Kif { }
class Bar extends Kif { }

let /** ((typeof Foo.prototype)|(typeof Bar.prototype)) */ test;

KIF: new Kif();
FOO_PROTOTYPE: Foo.prototype;
BAR_PROTOTYPE: Bar.prototype;
`)

	// When
	c.result = builder.Build()

	// Then
	c.assertThatResultAsTable().
		ContainsCell(c.id("KIF"), c.ids("FOO_PROTOTYPE", "BAR_PROTOTYPE"), ColorGraphBuilder_EdgeReason_ALGEBRAIC)

	// Also check post-conditions.
}

// Java: ColorGraphBuilderTest.verifyResult_topNodeIsOnlyRoot()
func (c *ColorGraphBuilderTest) VerifyResult_topNodeIsOnlyRoot() {
	roots := jutil.NewObjectLinkedHashSet[DiGraph_DiGraphNode[*ColorGraphNode, any]]()
	for _, n := range c.result.GetDiGraphNodes() {
		if len(n.GetInEdges()) == 0 {
			roots.Add(n)
		}
	}
	AssertThat(roots).
		ComparingElementsUsing(ColorGraphBuilderTest_NODE_HAS_ID).
		ContainsExactly(colors.StandardColors_UNKNOWN.GetId())
}

// Java: ColorGraphBuilderTest.verifyResult_hasNoSelfEdges()
func (c *ColorGraphBuilderTest) VerifyResult_hasNoSelfEdges() {
	selfEdges := jutil.NewObjectLinkedHashSet[DiGraph_DiGraphEdge[*ColorGraphNode, any]]()
	for _, e := range c.result.GetDiGraphEdges() {
		if jutil.ObjectsEquals(e.GetSource(), e.GetDestination()) {
			selfEdges.Add(e)
		}
	}
	AssertThat(selfEdges).IsEmpty()
}

// Java: ColorGraphBuilderTest.verifyResult_hasNoParallelEdges()
func (c *ColorGraphBuilderTest) VerifyResult_hasNoParallelEdges() {
	AssertThat(c.result.GetDiGraphEdges()).
		ComparingElementsUsing(
			Correspondence_From(
				func(a, e DiGraph_DiGraphEdge[*ColorGraphNode, any]) bool {
					return jutil.ObjectsEquals(a.GetSource(), e.GetSource()) &&
						jutil.ObjectsEquals(a.GetDestination(), e.GetDestination()) &&
						!jutil.ObjectsEquals(a, e)
				},
				"is parallel to")).
		ContainsNoneIn(c.result.GetDiGraphEdges())
}

// Java: ColorGraphBuilderTest.renderResultGraph()
func (c *ColorGraphBuilderTest) RenderResultGraph() {
}

// Java: ColorGraphBuilderTest.assertThatResultAsTable()
func (c *ColorGraphBuilderTest) assertThatResultAsTable() *colorGraphBuilderTest_TableSubject {
	table := jutil.NewLinkedHashMap[colorGraphBuilderTest_cell, ColorGraphBuilder_EdgeReason]()
	for _, edge := range c.result.GetDiGraphEdges() {
		cell := colorGraphBuilderTest_cell{
			colorGraphBuilderTest_nameOfNode(edge.GetSource()), colorGraphBuilderTest_nameOfNode(edge.GetDestination())}
		// ImmutableTable.Builder.buildOrThrow() rejects duplicate cells.
		if table.ContainsKey(cell) {
			panic(jutil.NewIllegalArgumentExceptionMsg(
				"Duplicate key: (row=" + cell.row.String() + ", column=" + cell.column.String() + ")"))
		}
		table.Put(cell, edge.GetValue().(ColorGraphBuilder_EdgeReason))
	}
	return &colorGraphBuilderTest_TableSubject{table: table}
}

// colorGraphBuilderTest_cell is an ImmutableTable cell's row and column.
type colorGraphBuilderTest_cell struct {
	row, column colors.ColorId
}

// colorGraphBuilderTest_TableSubject is Truth's TableSubject over an
// ImmutableTable<ColorId, ColorId, EdgeReason>.
//
// Go-only: test-helper — the table is a map from (row, column) cells.
type colorGraphBuilderTest_TableSubject struct {
	table *jutil.LinkedHashMap[colorGraphBuilderTest_cell, ColorGraphBuilder_EdgeReason]
}

// ContainsCell is TableSubject.containsCell(rowKey, colKey, value).
func (s *colorGraphBuilderTest_TableSubject) ContainsCell(
	rowKey colors.ColorId, colKey colors.ColorId, value ColorGraphBuilder_EdgeReason) {
	got, ok := s.table.Lookup(colorGraphBuilderTest_cell{rowKey, colKey})
	if !ok || got != value {
		cells := []string{}
		for k, v := range s.table.All() {
			cells = append(cells, "("+k.row.String()+","+k.column.String()+")="+v.String())
		}
		Fail("value of: table.cellSet()\nexpected to contain: (" + rowKey.String() + "," + colKey.String() +
			")=" + value.String() + "\nbut was: " + jutil.ToJavaString(cells))
	}
}

// A fake implementation of a finder.
//
// Instances allow setting stub responses for findAll calls. Inputs an
// outputs are specfied using sets of type names.
//
// Java: ColorGraphBuilderTest.StubLcaFinder
type colorGraphBuilderTest_StubLcaFinder struct {
	LowestCommonAncestorFinderBase[*ColorGraphNode, any]

	graph DiGraph[*ColorGraphNode, any]
	// Go-only: language-adapter — a LinkedHashMap keyed by ImmutableSet
	// equality, as an insertion-ordered list of entries.
	stubs         []colorGraphBuilderTest_stub
	preconditions []*colorGraphBuilderTest_precondition
}

type colorGraphBuilderTest_stub struct {
	from, to *jutil.LinkedHashSet[colors.ColorId]
}

type colorGraphBuilderTest_precondition struct {
	run   func()
	count int
}

// Java: ColorGraphBuilderTest.StubLcaFinder()
func newColorGraphBuilderTest_StubLcaFinder() *colorGraphBuilderTest_StubLcaFinder {
	f := &colorGraphBuilderTest_StubLcaFinder{}
	f.InitLowestCommonAncestorFinderBase(nil) // super(null)
	return f
}

// Java: ColorGraphBuilderTest.StubLcaFinder.setGraph(DiGraph<ColorGraphNode, Object>)
func (f *colorGraphBuilderTest_StubLcaFinder) setGraph(graph DiGraph[*ColorGraphNode, any]) *colorGraphBuilderTest_StubLcaFinder {
	f.graph = graph
	return f
}

// Java: ColorGraphBuilderTest.StubLcaFinder.addStub(ImmutableSet<ColorId>, ImmutableSet<ColorId>)
func (f *colorGraphBuilderTest_StubLcaFinder) addStub(from, to *jutil.LinkedHashSet[colors.ColorId]) *colorGraphBuilderTest_StubLcaFinder {
	for i := range f.stubs {
		if f.stubs[i].from.Equals(from) {
			f.stubs[i].to = to // LinkedHashMap.put on an equal key
			return f
		}
	}
	f.stubs = append(f.stubs, colorGraphBuilderTest_stub{from, to})
	return f
}

// Java: ColorGraphBuilderTest.StubLcaFinder.findAll(Set<ColorGraphNode>)
func (f *colorGraphBuilderTest_StubLcaFinder) FindAll(roots jutil.Set[*ColorGraphNode]) *jutil.LinkedHashSet[*ColorGraphNode] {
	for _, entry := range f.preconditions {
		entry.run()
		entry.count--
	}

	rootIds := jutil.NewLinkedHashSet[colors.ColorId]()
	for r := range roots.All() {
		rootIds.Add(colorGraphBuilderTest_nameOf(r))
	}
	var resultNames *jutil.LinkedHashSet[colors.ColorId]
	stubKeys := []*jutil.LinkedHashSet[colors.ColorId]{}
	for _, s := range f.stubs {
		stubKeys = append(stubKeys, s.from)
		if s.from.Equals(rootIds) {
			resultNames = s.to
		}
	}
	if resultNames == nil {
		// assertThat(this.stubs).containsKey(rootIds)
		Fail("value of: map.keySet()\nexpected to contain: " + rootIds.String() +
			"\nbut was: " + jutil.ToJavaString(stubKeys))
	}

	results := jutil.NewLinkedHashSet[*ColorGraphNode]()
	for _, n := range f.graph.GetDiGraphNodes() {
		if resultNames.Contains(colorGraphBuilderTest_nameOf(n.GetValue())) {
			results.Add(n.GetValue())
		}
	}
	AssertThat(results).HasSize(resultNames.Size())

	return results
}

var _ LowestCommonAncestorFinder[*ColorGraphNode, any] = (*colorGraphBuilderTest_StubLcaFinder)(nil)

// Java: ColorGraphBuilderTest.nameOf(DiGraphNode<ColorGraphNode, Object>)
func colorGraphBuilderTest_nameOfNode(node DiGraph_DiGraphNode[*ColorGraphNode, any]) colors.ColorId {
	return colorGraphBuilderTest_nameOf(node.GetValue())
}

// Java: ColorGraphBuilderTest.nameOf(ColorGraphNode)
func colorGraphBuilderTest_nameOf(flat *ColorGraphNode) colors.ColorId {
	return flat.GetColor().GetId()
}

// Java: ColorGraphBuilderTest.colorWithId(ColorId)
func colorGraphBuilderTest_colorWithId(id colors.ColorId) *colors.Color_Builder {
	return colors.Color_SingleBuilder().SetId(id)
}

// Java: ColorGraphBuilderTest.NODE_HAS_ID
var ColorGraphBuilderTest_NODE_HAS_ID = Correspondence_Transforming(colorGraphBuilderTest_nameOfNode, "in a node with type")

// Java: ColorGraphBuilderTest.id(String)
func (c *ColorGraphBuilderTest) id(labelName string) colors.ColorId {
	id, ok := c.labelToId.Lookup(labelName)
	if !ok {
		panic(jutil.NewNullPointerExceptionMsg("Could not find label " + labelName))
	}
	return id
}

// Java: ColorGraphBuilderTest.id(String...)
func (c *ColorGraphBuilderTest) ids(idSources ...string) colors.ColorId {
	ids := jutil.NewLinkedHashSet[colors.ColorId]()
	for _, idSource := range idSources {
		ids.Add(c.id(idSource))
	}
	return colors.ColorId_Union(ids)
}
