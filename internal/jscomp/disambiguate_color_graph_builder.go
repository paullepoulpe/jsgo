// Port of com/google/javascript/jscomp/disambiguate/ColorGraphBuilder.java @ closure-compiler v20260415

/*
 * Copyright 2020 The Closure Compiler Authors.
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
	"iter"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// ColorGraphBuilder builds a graph of the Colors on the AST from a
// specified set of seed colors.
type ColorGraphBuilder struct {
	registry    *colors.ColorRegistry
	nodeFactory *ColorGraphNodeFactory
	lcaFinder   LowestCommonAncestorFinder[*ColorGraphNode, any]
	topNode     *LinkedDirectedGraph_LinkedDiGraphNode[*ColorGraphNode, any]

	// The graph of colors as defined by `holdsInstanceOf`.
	//
	// We use `holdsInstanceOf` rather than `isSupertypeOf` because we only
	// actually care about edges that were "used" in the program. If
	// instances never flow over an edge at runtime, then properties also
	// don't need to be tracked across that edge either. Of course, since we
	// don't track the types of all assignments in a program, many of the
	// edges are `isSupertypeOf` edges that we include to be conservative.
	//
	// This graph, when fully constructed, is still only an approximation.
	// This is the due to both memory and time constraints. The following
	// quirks are expected:
	//
	//   - Some colors, such as primitives, will not have nodes.
	//   - Transitive edges (shortcut edges for which there exist alternate
	//     paths) are kept minimal.
	//
	// In practice, this could be declared as taking an EdgeReason instead of
	// an Object, but Object is used to indicate that EdgeReasons are only
	// meant for debugging and not any actual logic in (dis)ambiguation.
	colorHoldsInstanceGraph *LinkedDirectedGraph[*ColorGraphNode, any]
}

// ColorGraphBuilder_EdgeReason is the relationship that caused an edge to
// be created.
//
// This information is only retained for diagnostics, not correctness.
type ColorGraphBuilder_EdgeReason int32

const (
	ColorGraphBuilder_EdgeReason_ALGEBRAIC ColorGraphBuilder_EdgeReason = iota
	ColorGraphBuilder_EdgeReason_CAN_HOLD
)

var colorGraphBuilder_EdgeReason_names = [...]string{"ALGEBRAIC", "CAN_HOLD"}

// Java: ColorGraphBuilder.EdgeReason.values()
func ColorGraphBuilder_EdgeReason_values() []ColorGraphBuilder_EdgeReason {
	return []ColorGraphBuilder_EdgeReason{
		ColorGraphBuilder_EdgeReason_ALGEBRAIC,
		ColorGraphBuilder_EdgeReason_CAN_HOLD,
	}
}

// Java: ColorGraphBuilder.EdgeReason.name()
func (r ColorGraphBuilder_EdgeReason) Name() string { return colorGraphBuilder_EdgeReason_names[r] }

// Java: ColorGraphBuilder.EdgeReason.toString()
func (r ColorGraphBuilder_EdgeReason) String() string { return r.Name() }

// Java: ColorGraphBuilder.EdgeReason.ordinal()
func (r ColorGraphBuilder_EdgeReason) Ordinal() int { return int(r) }

// Java: ColorGraphBuilder(ColorGraphNodeFactory, LowestCommonAncestorFinder.Factory<ColorGraphNode, Object>, ColorRegistry)
func NewColorGraphBuilder(
	nodeFactory *ColorGraphNodeFactory,
	lcaFinderFactory LowestCommonAncestorFinder_Factory[*ColorGraphNode, any],
	registry *colors.ColorRegistry) *ColorGraphBuilder {
	b := &ColorGraphBuilder{
		colorHoldsInstanceGraph: LinkedDirectedGraph_CreateWithoutAnnotations[*ColorGraphNode, any](),
	}
	b.registry = registry
	b.nodeFactory = nodeFactory
	b.lcaFinder = lcaFinderFactory(b.colorHoldsInstanceGraph)

	b.topNode =
		b.colorHoldsInstanceGraph.CreateLinkedDiGraphNode(
			b.nodeFactory.CreateNode(colors.StandardColors_UNKNOWN))
	return b
}

// Java: ColorGraphBuilder.add(ColorGraphNode)
func (b *ColorGraphBuilder) Add(flat *ColorGraphNode) {
	b.addInternal(flat)
}

// Java: ColorGraphBuilder.addAll(Collection<ColorGraphNode>)
func (b *ColorGraphBuilder) AddAll(flats iter.Seq[*ColorGraphNode]) {
	for flat := range flats {
		b.Add(flat)
	}
}

// Java: ColorGraphBuilder.build()
func (b *ColorGraphBuilder) Build() *LinkedDirectedGraph[*ColorGraphNode, any] {
	// connectUnionWithAncestors adds edges only, never nodes, so iterating a
	// snapshot of the nodes matches Java's forEach over the live view.
	for _, node := range b.colorHoldsInstanceGraph.GetLinkedDiGraphNodes() {
		b.connectUnionWithAncestors(node)
	}

	temp := b.colorHoldsInstanceGraph
	b.colorHoldsInstanceGraph = nil
	return temp
}

// connectUnionWithAncestors: during initial lattice construction unions
// were only given outbound edges. Here we add any necessary inbound ones.
//
// We defer this operation because adding union-to-union and
// common-supertype-to-union edges, is a hard problem. Solving it after all
// other colors are in place makes it easier.
//
// Java: ColorGraphBuilder.connectUnionWithAncestors(LinkedDiGraphNode<ColorGraphNode, Object>)
func (b *ColorGraphBuilder) connectUnionWithAncestors(unionNode *LinkedDirectedGraph_LinkedDiGraphNode[*ColorGraphNode, any]) {
	flatUnion := unionNode.GetValue()
	if !flatUnion.GetColor().IsUnion() {
		return
	}

	// Connect the LCAs to the union.
	//
	// The union itself will be found in most cases, but since we don't add
	// self-edges, that won't matter.
	//
	// Some of these edges may pollute the "lattice-ness" of the graph, but
	// all the invariants we actually care about will be maintained. If
	// disambiguation is too slow and stricter invariants would help, we
	// could be more careful.
	jutil.CheckState(len(unionNode.GetOutEdges()) != 0)
	graphNodes := jutil.NewLinkedHashSet[*ColorGraphNode]() // toImmutableSet()
	for alt := range flatUnion.GetColor().GetUnionElements().All() {
		graphNodes.Add(b.nodeFactory.CreateNode(alt))
	}
	for lca := range b.lcaFinder.FindAll(graphNodes).All() {
		b.connectSourceToDest(
			jutil.CheckNotNull(b.colorHoldsInstanceGraph.GetLinkedDiGraphNode(lca)), ColorGraphBuilder_EdgeReason_ALGEBRAIC, unionNode)
	}
}

// addInternalColor inserts color and all necessary related colors into the
// datastructures of this pass.
//
// Java: ColorGraphBuilder.addInternal(Color)
func (b *ColorGraphBuilder) addInternalColor(color *colors.Color) *LinkedDirectedGraph_LinkedDiGraphNode[*ColorGraphNode, any] {
	return b.addInternal(b.nodeFactory.CreateNode(color))
}

// addInternal inserts node and all necessary related colors into the
// datastructures of this pass.
//
// Java: ColorGraphBuilder.addInternal(ColorGraphNode)
func (b *ColorGraphBuilder) addInternal(node *ColorGraphNode) *LinkedDirectedGraph_LinkedDiGraphNode[*ColorGraphNode, any] {
	flatNode := b.colorHoldsInstanceGraph.GetLinkedDiGraphNode(node)
	if flatNode != nil {
		return flatNode
	}
	flatNode = b.colorHoldsInstanceGraph.CreateLinkedDiGraphNode(node)

	if node.GetColor().IsUnion() {
		for alt := range node.GetColor().GetUnionElements().All() {
			b.connectSourceToDest(flatNode, ColorGraphBuilder_EdgeReason_ALGEBRAIC, b.addInternalColor(alt))
		}
		return flatNode
	}

	color := node.GetColor()
	supertypes := b.registry.GetDisambiguationSupertypes(color)
	if supertypes.IsEmpty() {
		b.connectSourceToDest(b.topNode, ColorGraphBuilder_EdgeReason_ALGEBRAIC, flatNode)
	} else {
		for supertype := range supertypes.All() {
			b.connectSourceToDest(b.addInternalColor(supertype), ColorGraphBuilder_EdgeReason_CAN_HOLD, flatNode)
		}
	}

	// Add all instance and prototype colors when visiting a constructor. We
	// won't necessarily see all possible instance colors that exist at
	// runtime during an AST traversal.
	//
	// For example, a subclass constructor may never be explicitly
	// initialized but instead passed to some function expecting
	// `function(new:Parent)`. See
	// AmbiguatePropertiesTest#testImplementsAndExtends_respectsUndeclaredProperties().
	for prototype := range color.GetPrototypes().All() {
		b.addInternalColor(prototype)
	}
	for instanceColor := range color.GetInstanceColors().All() {
		b.addInternalColor(instanceColor)
	}
	return flatNode
}

// Java: ColorGraphBuilder.connectSourceToDest(LinkedDiGraphNode<ColorGraphNode, Object>, EdgeReason, LinkedDiGraphNode<ColorGraphNode, Object>)
func (b *ColorGraphBuilder) connectSourceToDest(
	source *LinkedDirectedGraph_LinkedDiGraphNode[*ColorGraphNode, any],
	reason ColorGraphBuilder_EdgeReason,
	dest *LinkedDirectedGraph_LinkedDiGraphNode[*ColorGraphNode, any]) {
	if source == dest ||
		b.colorHoldsInstanceGraph.IsConnectedInDirectionOfNodes(source, func(any) bool { return true }, dest) {
		return
	}

	b.colorHoldsInstanceGraph.ConnectNodes(source, reason, dest)
}
