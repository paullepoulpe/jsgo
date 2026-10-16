// Port of com/google/javascript/jscomp/disambiguate/ColorGraphBuilderTest.java @ closure-compiler v20260415

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
)

// The tests built from JS source (createBuilderIncludingCode) need the
// Compiler and CompilerTestCase to produce colors and are skipped; the
// registry-only test and the @After checks are ported.

// colorGraphBuilderTest holds the per-test fields of ColorGraphBuilderTest.
type colorGraphBuilderTest struct {
	t *testing.T

	// This registry is only used to lookup box colors so it doesn't have to
	// be the same one used by the graph builder.
	graphNodeFactory *ColorGraphNodeFactory

	result DiGraph[*ColorGraphNode, any]
}

func newColorGraphBuilderTest(t *testing.T) *colorGraphBuilderTest {
	return &colorGraphBuilderTest{
		t: t,
		graphNodeFactory: ColorGraphNodeFactory_CreateFactory(
			colors.ColorRegistry_Builder_().SetDefaultNativeColorsForTesting().Build()),
	}
}

// after runs the @After methods.
func (c *colorGraphBuilderTest) after() {
	c.t.Helper()
	c.verifyResult_topNodeIsOnlyRoot()
	c.verifyResult_hasNoSelfEdges()
	c.verifyResult_hasNoParallelEdges()
}

func (c *colorGraphBuilderTest) createBuilder(optLcaFinder *colorGraphBuilderTest_StubLcaFinder, registry *colors.ColorRegistry) *ColorGraphBuilder {
	lcaFinder := optLcaFinder
	if lcaFinder == nil {
		lcaFinder = newColorGraphBuilderTest_StubLcaFinder()
	}
	return NewColorGraphBuilder(c.graphNodeFactory,
		func(graph DiGraph[*ColorGraphNode, any]) LowestCommonAncestorFinder[*ColorGraphNode, any] {
			return lcaFinder.setGraph(graph)
		},
		registry)
}

func TestColorGraphBuilder_Top_isAboveInterface(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase (type-checked colors from source)")
}

func TestColorGraphBuilder_PrototypeChain_isInserted(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase (type-checked colors from source)")
}

func TestColorGraphBuilder_PrototypeChain_isInserted_throughExtends(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase (type-checked colors from source)")
}

func TestColorGraphBuilder_PrototypeChain_canBranch(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase (type-checked colors from source)")
}

func TestColorGraphBuilder_ConstructorDef_includesPrototypeAndInstanceType_evenIfUnused(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase (type-checked colors from source)")
}

func TestColorGraphBuilder_PrototypeChain_connectsInterfaces(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase (type-checked colors from source)")
}

func TestColorGraphBuilder_PrototypeChain_connectsSubclasses_viaClassSideInheritance(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase (type-checked colors from source)")
}

func TestColorGraphBuilder_Unions_connectedAboveMembers(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase (type-checked colors from source)")
}

func TestColorGraphBuilder_Unions_connectBelowLca(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase (type-checked colors from source)")
}

func TestColorGraphBuilder_Unions_connectBelowLca_withMultipleLcas(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase (type-checked colors from source)")
}

func TestColorGraphBuilder_Unions_connectBelowLca_whichIsAlsoUnion(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase (type-checked colors from source)")
}

func TestColorGraphBuilder_Unions_connectBelowLac_whichHasSameDescendantCount(t *testing.T) {
	t.Skip("needs Compiler, CompilerTestCase (type-checked colors from source)")
}

func TestColorGraphBuilder_DisambiguationSupertypes_createConnection(t *testing.T) {
	c := newColorGraphBuilderTest(t)
	defer c.after()

	// Given
	parent := colorGraphBuilderTest_colorWithId(colors.ColorId_FromUnsigned(100)).Build()
	child := colorGraphBuilderTest_colorWithId(colors.ColorId_FromUnsigned(101)).Build()

	builder := c.createBuilder(
		nil,
		colors.ColorRegistry_Builder_().
			SetDefaultNativeColorsForTesting().
			AddDisambiguationEdge(child, parent).
			Build())

	builder.Add(c.graphNodeFactory.CreateNode(child))

	// When
	c.result = builder.Build()

	// Then
	c.assertResultTableContainsCell(parent.GetId(), child.GetId(), ColorGraphBuilder_EdgeReason_CAN_HOLD)
}

// Java: @After verifyResult_topNodeIsOnlyRoot()
func (c *colorGraphBuilderTest) verifyResult_topNodeIsOnlyRoot() {
	c.t.Helper()
	roots := []colors.ColorId{}
	for _, n := range c.result.GetDiGraphNodes() {
		if len(n.GetInEdges()) == 0 {
			roots = append(roots, colorGraphBuilderTest_nameOfNode(n))
		}
	}
	assertSliceEqual(c.t, "root ids", roots, []colors.ColorId{colors.StandardColors_UNKNOWN.GetId()})
}

// Java: @After verifyResult_hasNoSelfEdges()
func (c *colorGraphBuilderTest) verifyResult_hasNoSelfEdges() {
	c.t.Helper()
	for _, e := range c.result.GetDiGraphEdges() {
		if e.GetSource() == e.GetDestination() {
			c.t.Errorf("self edge %v", e)
		}
	}
}

// Java: @After verifyResult_hasNoParallelEdges()
func (c *colorGraphBuilderTest) verifyResult_hasNoParallelEdges() {
	c.t.Helper()
	edges := c.result.GetDiGraphEdges()
	for _, a := range edges {
		for _, e := range edges {
			if a.GetSource() == e.GetSource() && a.GetDestination() == e.GetDestination() && a != e {
				c.t.Errorf("edge %v is parallel to %v", a, e)
			}
		}
	}
}

// assertResultTableContainsCell is
// assertThatResultAsTable().containsCell(row, column, value).
func (c *colorGraphBuilderTest) assertResultTableContainsCell(row, column colors.ColorId, value ColorGraphBuilder_EdgeReason) {
	c.t.Helper()
	// ImmutableTable.Builder.buildOrThrow() rejects duplicate cells.
	type cell struct{ row, col colors.ColorId }
	table := map[cell]ColorGraphBuilder_EdgeReason{}
	for _, edge := range c.result.GetDiGraphEdges() {
		k := cell{colorGraphBuilderTest_nameOfNode(edge.GetSource()), colorGraphBuilderTest_nameOfNode(edge.GetDestination())}
		if _, dup := table[k]; dup {
			c.t.Fatalf("duplicate table cell (%v, %v)", k.row, k.col)
		}
		table[k] = edge.GetValue().(ColorGraphBuilder_EdgeReason)
	}
	if got, ok := table[cell{row, column}]; !ok || got != value {
		c.t.Errorf("table cell (%v, %v): got %v (present=%v), want %v", row, column, got, ok, value)
	}
}

// colorGraphBuilderTest_StubLcaFinder is a fake implementation of a finder.
//
// Instances allow setting stub responses for findAll calls. Inputs an
// outputs are specfied using sets of type names.
type colorGraphBuilderTest_StubLcaFinder struct {
	LowestCommonAncestorFinderBase[*ColorGraphNode, any]

	graph         DiGraph[*ColorGraphNode, any]
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

func newColorGraphBuilderTest_StubLcaFinder() *colorGraphBuilderTest_StubLcaFinder {
	f := &colorGraphBuilderTest_StubLcaFinder{}
	f.InitLowestCommonAncestorFinderBase(nil) // super(null)
	return f
}

func (f *colorGraphBuilderTest_StubLcaFinder) setGraph(graph DiGraph[*ColorGraphNode, any]) *colorGraphBuilderTest_StubLcaFinder {
	f.graph = graph
	return f
}

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
	for _, s := range f.stubs {
		if s.from.Equals(rootIds) {
			resultNames = s.to
		}
	}
	if resultNames == nil {
		panic("stubs does not contain key " + rootIds.String())
	}

	results := jutil.NewLinkedHashSet[*ColorGraphNode]()
	for _, n := range f.graph.GetDiGraphNodes() {
		if resultNames.Contains(colorGraphBuilderTest_nameOf(n.GetValue())) {
			results.Add(n.GetValue())
		}
	}
	if results.Size() != resultNames.Size() {
		panic("unexpected stub result size")
	}

	return results
}

var _ LowestCommonAncestorFinder[*ColorGraphNode, any] = (*colorGraphBuilderTest_StubLcaFinder)(nil)

func colorGraphBuilderTest_nameOfNode(node DiGraph_DiGraphNode[*ColorGraphNode, any]) colors.ColorId {
	return colorGraphBuilderTest_nameOf(node.GetValue())
}

func colorGraphBuilderTest_nameOf(flat *ColorGraphNode) colors.ColorId {
	return flat.GetColor().GetId()
}

func colorGraphBuilderTest_colorWithId(id colors.ColorId) *colors.Color_Builder {
	return colors.Color_SingleBuilder().SetId(id)
}
