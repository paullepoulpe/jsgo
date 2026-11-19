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
// Port of com/google/javascript/jscomp/DotFormatter.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_graph "github.com/paullepoulpe/jsgo/transpiled/gen/graph"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_dot_formatter_1 func(any) string
var lambdaSingleton_dot_formatter_2 func(any) string

func init() {
	lambdaSingleton_dot_formatter_1 = func(mrp0_ any) string {
		return j2gort.AsIface[pkg_graph.GraphvizGraph_GraphvizEdge](mrp0_).GetNode1Id()
	}
	lambdaSingleton_dot_formatter_2 = func(mrp0_ any) string {
		return j2gort.AsIface[pkg_graph.GraphvizGraph_GraphvizEdge](mrp0_).GetNode2Id()
	}
}

// Java class: com.google.javascript.jscomp.DotFormatter
type DotFormatter struct {
	assignments      j2gort.OrderedMap[pkg_rhino.Node, *j2gort.BoxInt] `gson:"assignments"`
	builder          j2gort.Ext_java_lang_Appendable                   `gson:"builder"`
	cfg              ControlFlowGraph[pkg_rhino.Node]                  `gson:"cfg"`
	keyCount         int32                                             `gson:"keyCount"`
	idHash_          atomic.Int32
	printAnnotations bool `gson:"printAnnotations"`
}

func (this_ *DotFormatter) BindSelf_(s any) {
	_ = s
}
func (this_ *DotFormatter) IsJavaNil_() bool { return this_ == nil }
func (*DotFormatter) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DotFormatter", "DotFormatter"
}
func (this_ *DotFormatter) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DotFormatter) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DotFormatter", "com.google.javascript.jscomp.DotFormatter", func(v_ any) bool {
		return j2gort.InstanceOf[*DotFormatter](v_)
	})
}

func StaticGet_DotFormatter_INDENT() string               { return "  " }
func StaticGet_DotFormatter_ARROW() string                { return " -> " }
func StaticGet_DotFormatter_LINE() string                 { return " -- " }
func StaticGet_DotFormatter_MAX_LABEL_NAME_LENGTH() int32 { return int32(10) }
func (this_ *DotFormatter) initFields_DotFormatter() {
	_ = this_
	this_.assignments = j2gort.NewOrderedMap[pkg_rhino.Node, *j2gort.BoxInt]()

}

// Java constructor: DotFormatter()
func NewDotFormatter() *DotFormatter {
	this_ := &DotFormatter{}
	InitDotFormatter(this_)
	return this_
}

func InitDotFormatter(this_ *DotFormatter) {
	this_.initFields_DotFormatter()
	this_.builder = j2gort.NewExt_java_lang_StringBuilder()
	this_.cfg = nil
	this_.printAnnotations = false
}

// Java constructor: DotFormatter(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.ControlFlowGraph<com.google.javascript.rhino.Node> cfg, java.lang.Appendable builder, boolean printAnnotations)
func NewDotFormatter_4(n pkg_rhino.Node, cfg ControlFlowGraph[pkg_rhino.Node], builder j2gort.Ext_java_lang_Appendable, printAnnotations bool) *DotFormatter {
	this_ := &DotFormatter{}
	InitDotFormatter_4(this_, n, cfg, builder, printAnnotations)
	return this_
}

func InitDotFormatter_4(this_ *DotFormatter, n pkg_rhino.Node, cfg ControlFlowGraph[pkg_rhino.Node], builder j2gort.Ext_java_lang_Appendable, printAnnotations bool) {
	_ = n
	_ = cfg
	_ = builder
	_ = printAnnotations
	this_.initFields_DotFormatter()
	this_.cfg = cfg
	this_.builder = builder
	this_.printAnnotations = printAnnotations
	this_.priv_DotFormatter_formatPreamble_2453()
	this_.priv_DotFormatter_traverseNodes_2451(n)
	this_.priv_DotFormatter_formatConclusion_2454()
}

// Java static method: DotFormatter.toDot(com.google.javascript.rhino.Node n)
func DotFormatter_ToDot_1_Node(n pkg_rhino.Node) string {
	_ = n
	return DotFormatter_ToDot_2(n, nil)
}

// Java static method: DotFormatter.toDot(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable ControlFlowGraph<com.google.javascript.rhino.Node> inCFG)
func DotFormatter_ToDot_2(n pkg_rhino.Node, inCFG ControlFlowGraph[pkg_rhino.Node]) string {
	_ = n
	_ = inCFG
	var builder *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = builder
	var val_1 DotFormatter
	InitDotFormatter_4(&val_1, n, inCFG, builder, false)
	var unused *DotFormatter = &val_1
	_ = unused
	return builder.String()
}

// Java static method: DotFormatter.toDot(com.google.javascript.jscomp.graph.GraphvizGraph graph)
func DotFormatter_ToDot_1_GraphvizGraph(graph pkg_graph.GraphvizGraph) string {
	_ = graph
	var builder j2gort.LocalStringBuilder

	builder.AppendStr(func() string {
		if graph.IsDirected() {
			return "digraph"
		}
		return "graph"
	}())
	builder.AppendStr("  ")
	builder.AppendStr(graph.GetName())
	builder.AppendStr(" {\n")
	builder.AppendStr("  ")
	builder.AppendStr("node [color=lightblue2, style=filled];\n")
	var edgeSymbol string
	_ = edgeSymbol
	if graph.IsDirected() {
		edgeSymbol = " -> "
	} else {
		edgeSymbol = " -- "
	}
	var nodes j2gort.Ext_java_util_List = graph.GetGraphvizNodes()
	_ = nodes
	var nodeNames []string = j2gort.NewStringArray(nodes.Size())
	_ = nodeNames
	{
		var i int32 = int32(0)
		_ = i
		for ; i < int32(len(nodeNames)); i = i + 1 {
			var gNode pkg_graph.GraphvizGraph_GraphvizNode = j2gort.AsIface[pkg_graph.GraphvizGraph_GraphvizNode](nodes.Get(i))
			_ = gNode
			{
				a_2 := nodeNames
				i_3 := i
				a_2[i_3] = (j2gort.StrValString(gNode.GetId()) + " [label=\"" + j2gort.StrValString(gNode.GetLabel()) + "\" color=\"" + j2gort.StrValString(gNode.GetColor()) + "\"]")
			}
		}
	}
	j2gort.ArraysSortStrings(nodeNames)
	for _, el_4 := range j2gort.CheckNotNullSlice(nodeNames) {
		nodeName := el_4
		_ = nodeName
		builder.AppendStr("  ")
		builder.AppendStr(nodeName)
		builder.AppendStr(";\n")
	}
	var edges j2gort.Ext_java_util_List = graph.GetGraphvizEdges()
	_ = edges
	edges.Sort(j2gort.ComparatorThenComparingStr(j2gort.ComparatorComparingStr(lambdaSingleton_dot_formatter_1), lambdaSingleton_dot_formatter_2))
	{
		cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(edges))
		for cur_5.HasNext() {
			rawElem_6 := cur_5.Next()
			edge := j2gort.AsIface[pkg_graph.GraphvizGraph_GraphvizEdge](rawElem_6)
			_ = edge
			builder.AppendStr("  ")
			builder.AppendStr(edge.GetNode1Id())
			builder.AppendStr(edgeSymbol)
			builder.AppendStr(edge.GetNode2Id())
			builder.AppendStr(" [label=\"").AppendStr(edge.GetLabel()).AppendStr("\" color=\"").AppendStr(edge.GetColor()).AppendStr("\"]")
			builder.AppendStr(";\n")
		}
	}
	builder.AppendStr("}\n")
	return builder.String()
}

// Java static method: DotFormatter.appendDot(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.ControlFlowGraph<com.google.javascript.rhino.Node> inCFG, java.lang.Appendable builder)
func DotFormatter_AppendDot(n pkg_rhino.Node, inCFG ControlFlowGraph[pkg_rhino.Node], builder j2gort.Ext_java_lang_Appendable) {
	_ = n
	_ = inCFG
	_ = builder
	var val_7 DotFormatter
	InitDotFormatter_4(&val_7, n, inCFG, builder, false)
	var unused *DotFormatter = &val_7
	_ = unused
}

// Java static method: DotFormatter.newInstanceForTesting()
func DotFormatter_NewInstanceForTesting() *DotFormatter {
	return NewDotFormatter()
}

// Java method: DotFormatter.traverseNodes(com.google.javascript.rhino.Node parent)
func (this_ *DotFormatter) priv_DotFormatter_traverseNodes_2451(parent pkg_rhino.Node) {
	_ = parent
	var keyParent int32 = this_.Key(parent)
	_ = keyParent
	{
		var child pkg_rhino.Node = parent.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			var keyChild int32 = this_.Key(child)
			_ = keyChild
			this_.builder.Append_1_any("  ")
			this_.builder.Append_1_any(j2gort.BoxString(DotFormatter_FormatNodeName(j2gort.ValueOfInt(keyParent))))
			this_.builder.Append_1_any(" -> ")
			this_.builder.Append_1_any(j2gort.BoxString(DotFormatter_FormatNodeName(j2gort.ValueOfInt(keyChild))))
			this_.builder.Append_1_any(" [weight=1];\n")
			this_.priv_DotFormatter_traverseNodes_2451(child)
		}
	}
	if (this_.cfg != nil && !this_.cfg.IsJavaNil_()) && this_.cfg.HasNode(parent) {
		var outEdges j2gort.Ext_java_util_List = this_.cfg.GetOutEdges(parent)
		_ = outEdges
		var edgeList []string = j2gort.NewStringArray(outEdges.Size())
		_ = edgeList
		{
			var i int32 = int32(0)
			_ = i
			for ; i < int32(len(edgeList)); i = i + 1 {
				var edge pkg_graph.DiGraph_DiGraphEdge[pkg_rhino.Node, ControlFlowGraph_Branch] = j2gort.AsIface[pkg_graph.DiGraph_DiGraphEdge[pkg_rhino.Node, ControlFlowGraph_Branch]](outEdges.Get(i))
				_ = edge
				var succ pkg_graph.DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch] = edge.GetDestination()
				_ = succ
				var toNode string = j2gort.NullString
				_ = toNode
				if j2gort.IfaceEquals[pkg_graph.DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch]](succ, this_.cfg.GetImplicitReturn()) {
					toNode = "RETURN"
				} else {
					var keySucc int32 = this_.Key(succ.GetValue())
					_ = keySucc
					toNode = DotFormatter_FormatNodeName(j2gort.ValueOfInt(keySucc))
				}
				{
					a_8 := edgeList
					i_9 := i
					a_8[i_9] = (j2gort.StrValString(DotFormatter_FormatNodeName(j2gort.ValueOfInt(keyParent))) + " -> " + j2gort.StrValString(toNode) + " [label=\"" + j2gort.StrVal(edge.GetValue()) + "\", fontcolor=\"red\", weight=0.01, color=\"red\"];\n")
				}
			}
		}
		j2gort.ArraysSortStrings(edgeList)
		for _, el_10 := range j2gort.CheckNotNullSlice(edgeList) {
			element := el_10
			_ = element
			this_.builder.Append_1_any("  ")
			this_.builder.Append_1_any(j2gort.BoxString(element))
		}
	}
}

// Java method: DotFormatter.key(com.google.javascript.rhino.Node n)
func (this_ *DotFormatter) Key(n pkg_rhino.Node) int32 {
	_ = n
	var key *j2gort.BoxInt = this_.assignments.Get(n)
	_ = key
	if key == nil {
		var iifeArg_1 int32
		{

			old_11 := this_.keyCount
			nw_12 := old_11 + 1
			this_.keyCount = nw_12
			iifeArg_1 = old_11
		}
		key = j2gort.ValueOfInt(iifeArg_1)
		this_.assignments.Put(n, key)
		this_.builder.Append_1_any("  ")
		this_.builder.Append_1_any(j2gort.BoxString(DotFormatter_FormatNodeName(key)))
		this_.builder.Append_1_any(" [label=\"")
		this_.builder.Append_1_any(j2gort.BoxString(j2gort.CheckNotNullEnum(n.GetToken()).String()))
		if (((n.IsName() || n.IsStringLit()) || n.IsGetProp()) || n.IsImportStar()) || n.IsStringKey() {
			this_.builder.Append_1_any(j2gort.BoxString(DotFormatter_GetNodeLabel(n)))
		}
		var type_ pkg_rhino.JSType = n.GetJSType()
		_ = type_
		if type_ != nil && !type_.IsJavaNil_() {
			this_.builder.Append_1_any(" : ")
			this_.builder.Append_1_any(j2gort.BoxString(type_.GetJSTypeBase().String()))
		}
		if (this_.printAnnotations && (this_.cfg != nil && !this_.cfg.IsJavaNil_())) && this_.cfg.HasNode(n) {
			var annotation any = j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch]](this_.cfg.GetNode(n)).GetAnnotation()
			_ = annotation
			if annotation != nil {
				this_.builder.Append_1_any("\\n")
				this_.builder.Append_1_any(javalang.ToJavaString(j2gort.CheckNotNull(annotation)))
			}
		}
		this_.builder.Append_1_any("\"")
		if n.GetJSDocInfo() != nil {
			this_.builder.Append_1_any(" color=\"green\"")
		}
		this_.builder.Append_1_any("];\n")
	}
	return j2gort.UnboxInt(key)
}

// Java static method: DotFormatter.getNodeLabel(com.google.javascript.rhino.Node n)
func DotFormatter_GetNodeLabel(n pkg_rhino.Node) string {
	_ = n
	var content string = n.GetString()
	_ = content
	if j2gort.JStrLen(content) > int32(10) {
		content = j2gort.JStrSub(content, int32(0), int32(10))
	}
	if !(j2gort.CheckNotNullString(content) == "") {
		content = ("(" + j2gort.StrValString(content) + ")")
	}
	return content
}

// Java static method: DotFormatter.formatNodeName(java.lang.Integer key)
func DotFormatter_FormatNodeName(key *j2gort.BoxInt) string {
	_ = key
	return ("node" + j2gort.StrVal(key))
}

// Java method: DotFormatter.formatPreamble()
func (this_ *DotFormatter) priv_DotFormatter_formatPreamble_2453() {
	this_.builder.Append_1_any("digraph AST {\n")
	this_.builder.Append_1_any("  ")
	this_.builder.Append_1_any("node [color=lightblue2, style=filled];\n")
}

// Java method: DotFormatter.formatConclusion()
func (this_ *DotFormatter) priv_DotFormatter_formatConclusion_2454() {
	this_.builder.Append_1_any("}\n")
}
