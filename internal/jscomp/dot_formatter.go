// Port of com/google/javascript/jscomp/DotFormatter.java @ closure-compiler v20261021

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
	"slices"
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/graph"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const dotFormatter_INDENT = "  "
const dotFormatter_ARROW = " -> "
const dotFormatter_LINE = " -- "
const dotFormatter_MAX_LABEL_NAME_LENGTH = 10

// DotFormatter prints out a dot file of the Abstract Syntax Tree. For a
// detailed description of the dot format and visualization tool refer to
// Graphviz (http://www.graphviz.org).
//
// Typical usage of this class
// `logfile.write(new DotFormatter().toDot(node));`
//
// This class is not thread safe and should not be used without proper
// external synchronization.
//
// Go-only: language-adapter — the IOException that Java's Appendable may
// throw is a panic (see jutil.Appendable), so the `throws IOException`
// clauses have no Go counterpart.
type DotFormatter struct {
	// stores the current assignment of node to keys
	assignments map[rhino.Node]int

	// key count in order to assign a unique key to each node
	keyCount int

	// the builder used to generate the dot diagram
	builder jutil.Appendable

	cfg *ControlFlowGraph[rhino.Node] // nil is null

	printAnnotations bool
}

// NewDotFormatter is the constructor for testing only.
//
// Java: DotFormatter()
func NewDotFormatter() *DotFormatter {
	return &DotFormatter{
		assignments: make(map[rhino.Node]int),
		builder:     jutil.NewStringBuilder(),
	}
}

// Java: DotFormatter(Node, ControlFlowGraph<Node>, Appendable, boolean)
func NewDotFormatterWithGraph(
	n rhino.Node, cfg *ControlFlowGraph[rhino.Node], builder jutil.Appendable, printAnnotations bool) *DotFormatter {
	f := &DotFormatter{
		assignments:      make(map[rhino.Node]int),
		cfg:              cfg,
		builder:          builder,
		printAnnotations: printAnnotations,
	}

	f.FormatPreamble()
	f.TraverseNodes(n)
	f.FormatConclusion()
	return f
}

// DotFormatter_ToDot converts an AST to dot representation.
//
// n is the root of the AST described in the dot formatted string; returns
// the dot representation of the AST.
//
// Java: DotFormatter.toDot(Node)
func DotFormatter_ToDot(n rhino.Node) string {
	return DotFormatter_ToDotWithCfg(n, nil)
}

// DotFormatter_ToDotWithCfg converts an AST to dot representation.
//
// n is the root of the AST described in the dot formatted string; inCFG is
// the Control Flow Graph. Returns the dot representation of the AST.
//
// Java: DotFormatter.toDot(Node, ControlFlowGraph<Node>)
func DotFormatter_ToDotWithCfg(n rhino.Node, inCFG *ControlFlowGraph[rhino.Node]) string {
	builder := jutil.NewStringBuilder()
	_ = NewDotFormatterWithGraph(n, inCFG, builder, false)
	return builder.String()
}

// DotFormatter_ToDotOfGraphvizGraph outputs a string in DOT format that
// presents the graph.
//
// graph is the input graph. Returns a string in Dot format that presents
// the graph.
//
// Java: DotFormatter.toDot(GraphvizGraph)
func DotFormatter_ToDotOfGraphvizGraph(gvGraph graph.GraphvizGraph) string {
	builder := jutil.NewStringBuilder()
	if gvGraph.IsDirected() {
		builder.Append("digraph")
	} else {
		builder.Append("graph")
	}
	builder.Append(dotFormatter_INDENT)
	builder.Append(gvGraph.GetName())
	builder.Append(" {\n")
	builder.Append(dotFormatter_INDENT)
	builder.Append("node [color=lightblue2, style=filled];\n")

	edgeSymbol := dotFormatter_LINE
	if gvGraph.IsDirected() {
		edgeSymbol = dotFormatter_ARROW
	}

	nodes := gvGraph.GetGraphvizNodes()

	nodeNames := make([]string, len(nodes))

	for i := 0; i < len(nodeNames); i++ {
		gNode := nodes[i]
		nodeNames[i] = gNode.GetId() +
			" [label=\"" +
			gNode.GetLabel() +
			"\" color=\"" +
			gNode.GetColor() +
			"\"]"
	}

	// We sort the nodes so we get a deterministic output every time
	// regardless of the implementation of the graph data structure.
	slices.SortStableFunc(nodeNames, jutil.StringCompareTo)

	for _, nodeName := range nodeNames {
		builder.Append(dotFormatter_INDENT)
		builder.Append(nodeName)
		builder.Append(";\n")
	}

	// Sort the list of edges so the output is deterministic
	edges := gvGraph.GetGraphvizEdges()
	slices.SortStableFunc(edges, func(a, b graph.GraphvizGraph_GraphvizEdge) int {
		if c := jutil.StringCompareTo(a.GetNode1Id(), b.GetNode1Id()); c != 0 {
			return c
		}
		return jutil.StringCompareTo(a.GetNode2Id(), b.GetNode2Id())
	})

	for _, edge := range edges {
		builder.Append(dotFormatter_INDENT)
		builder.Append(edge.GetNode1Id())
		builder.Append(edgeSymbol)
		builder.Append(edge.GetNode2Id())
		builder.
			Append(" [label=\"").
			Append(edge.GetLabel()).
			Append("\" color=\"").
			Append(edge.GetColor()).
			Append("\"]")
		builder.Append(";\n")
	}

	builder.Append("}\n")
	return builder.String()
}

// DotFormatter_AppendDot converts an AST to dot representation and appends
// it to the given buffer.
//
// n is the root of the AST described in the dot formatted string; inCFG is
// the Control Flow Graph; builder is a place to dump the graph.
//
// Java: DotFormatter.appendDot(Node, ControlFlowGraph<Node>, Appendable)
func DotFormatter_AppendDot(n rhino.Node, inCFG *ControlFlowGraph[rhino.Node], builder jutil.Appendable) {
	_ = NewDotFormatterWithGraph(n, inCFG, builder, false)
}

// DotFormatter_NewInstanceForTesting creates a DotFormatter purely for
// testing DotFormatter's internal methods.
//
// Java: DotFormatter.newInstanceForTesting()
func DotFormatter_NewInstanceForTesting() *DotFormatter {
	return NewDotFormatter()
}

// Java: DotFormatter.traverseNodes(Node)
func (f *DotFormatter) TraverseNodes(parent rhino.Node) {
	// key
	keyParent := f.Key(parent)

	// edges
	for child := parent.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		keyChild := f.Key(child)
		f.builder.AppendCharSequence(dotFormatter_INDENT)
		f.builder.AppendCharSequence(DotFormatter_FormatNodeName(keyParent))
		f.builder.AppendCharSequence(dotFormatter_ARROW)
		f.builder.AppendCharSequence(DotFormatter_FormatNodeName(keyChild))
		f.builder.AppendCharSequence(" [weight=1];\n")

		f.TraverseNodes(child)
	}

	// Flow Edges
	if f.cfg != nil && f.cfg.HasNode(parent) {
		outEdges := f.cfg.GetOutEdges(parent)
		edgeList := make([]string, len(outEdges))
		for i := 0; i < len(edgeList); i++ {
			edge := outEdges[i]
			succ := edge.GetDestination()

			var toNode string
			if succ == f.cfg.GetImplicitReturn() {
				toNode = "RETURN"
			} else {
				keySucc := f.Key(succ.GetValue())
				toNode = DotFormatter_FormatNodeName(keySucc)
			}

			edgeList[i] = DotFormatter_FormatNodeName(keyParent) +
				dotFormatter_ARROW +
				toNode +
				" [label=\"" +
				edge.GetValue().String() +
				"\", " +
				"fontcolor=\"red\", " +
				"weight=0.01, color=\"red\"];\n"
		}

		slices.SortStableFunc(edgeList, jutil.StringCompareTo)

		for _, element := range edgeList {
			f.builder.AppendCharSequence(dotFormatter_INDENT)
			f.builder.AppendCharSequence(element)
		}
	}
}

// Java: DotFormatter.key(Node)
func (f *DotFormatter) Key(n rhino.Node) int {
	key, ok := f.assignments[n]
	if !ok {
		key = f.keyCount
		f.keyCount++
		f.assignments[n] = key
		f.builder.AppendCharSequence(dotFormatter_INDENT)
		f.builder.AppendCharSequence(DotFormatter_FormatNodeName(key))
		f.builder.AppendCharSequence(" [label=\"")
		f.builder.AppendCharSequence(n.GetToken().String())
		if n.IsName() || n.IsStringLit() || n.IsGetProp() || n.IsImportStar() || n.IsStringKey() {
			f.builder.AppendCharSequence(DotFormatter_GetNodeLabel(n))
		}
		type_ := n.GetJSType()
		if type_ != nil {
			f.builder.AppendCharSequence(" : ")
			f.builder.AppendCharSequence(type_.String())
		}
		if f.printAnnotations && f.cfg != nil && f.cfg.HasNode(n) {
			// Go-only: language-adapter — Java's cfg.getNode(n) is
			// LinkedDirectedGraph's covariant getNode; the annotation is
			// the same on every view of the node.
			annotation := f.cfg.GetDiGraphNode(n).GetAnnotation()
			if annotation != nil {
				f.builder.AppendCharSequence("\\n")
				f.builder.AppendCharSequence(jutil.ToJavaString(annotation))
			}
		}
		f.builder.AppendCharSequence("\"")
		if n.GetJSDocInfo() != nil {
			f.builder.AppendCharSequence(" color=\"green\"")
		}
		f.builder.AppendCharSequence("];\n")
	}
	return key
}

// Java: DotFormatter.getNodeLabel(Node)
func DotFormatter_GetNodeLabel(n rhino.Node) string {
	content := n.GetString()
	if jutil.StringLength(content) > dotFormatter_MAX_LABEL_NAME_LENGTH {
		content = jutil.StringSubstring(content, 0, dotFormatter_MAX_LABEL_NAME_LENGTH)
	}
	if content != "" {
		content = "(" + content + ")"
	}
	return content
}

// Java: DotFormatter.formatNodeName(Integer)
func DotFormatter_FormatNodeName(key int) string {
	return "node" + strconv.Itoa(key)
}

// Java: DotFormatter.formatPreamble()
func (f *DotFormatter) FormatPreamble() {
	f.builder.AppendCharSequence("digraph AST {\n")
	f.builder.AppendCharSequence(dotFormatter_INDENT)
	f.builder.AppendCharSequence("node [color=lightblue2, style=filled];\n")
}

// Java: DotFormatter.formatConclusion()
func (f *DotFormatter) FormatConclusion() {
	f.builder.AppendCharSequence("}\n")
}
