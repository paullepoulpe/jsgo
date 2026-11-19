/*
 * Copyright 2016 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/instrumentation/BranchCoverageInstrumentationCallback.java (generated mechanically by j2go).

package jscomp

import (
	pkg_graph "github.com/paullepoulpe/jsgo/transpiled/gen/graph"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.instrumentation.BranchCoverageInstrumentationCallback
type BranchCoverageInstrumentationCallback struct {
	NodeTraversal_AbstractCfgCallbackBase
	compiler            AbstractCompiler                                     `gson:"compiler"`
	instrumentationData *j2gort.OrderedMap[string, *FileInstrumentationData] `gson:"instrumentationData"`
}

func (this_ *BranchCoverageInstrumentationCallback) BindSelf_(s any) {
	_ = s
	this_.NodeTraversal_AbstractCfgCallbackBase.BindSelf_(s)
}
func (this_ *BranchCoverageInstrumentationCallback) IsJavaNil_() bool { return this_ == nil }
func (*BranchCoverageInstrumentationCallback) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.instrumentation.BranchCoverageInstrumentationCallback", "BranchCoverageInstrumentationCallback"
}
func (this_ *BranchCoverageInstrumentationCallback) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractCfgCallbackBase.IdentityHash_()
}
func (*BranchCoverageInstrumentationCallback) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.instrumentation.BranchCoverageInstrumentationCallback", "com.google.javascript.jscomp.instrumentation.BranchCoverageInstrumentationCallback", func(v_ any) bool {
		return j2gort.InstanceOf[*BranchCoverageInstrumentationCallback](v_)
	})
}
func (*BranchCoverageInstrumentationCallback) Implements_jscomp_NodeTraversal_AbstractCfgCallback() {}
func (*BranchCoverageInstrumentationCallback) Implements_jscomp_NodeTraversal_ScopedCallback()      {}
func (*BranchCoverageInstrumentationCallback) Implements_jscomp_NodeTraversal_Callback()            {}

func StaticGet_BranchCoverageInstrumentationCallback_BRANCH_ARRAY_NAME_PREFIX() string {
	return "JSCompiler_lcov_branch_data_"
}
func (this_ *BranchCoverageInstrumentationCallback) initFields_BranchCoverageInstrumentationCallback() {
	_ = this_
}

// Java constructor: BranchCoverageInstrumentationCallback(com.google.javascript.jscomp.AbstractCompiler compiler, java.util.Map<java.lang.String,com.google.javascript.jscomp.instrumentation.FileInstrumentationData> instrumentationData)
func NewBranchCoverageInstrumentationCallback(compiler AbstractCompiler, instrumentationData *j2gort.OrderedMap[string, *FileInstrumentationData]) *BranchCoverageInstrumentationCallback {
	this_ := &BranchCoverageInstrumentationCallback{}
	this_.BindSelf_(this_)
	InitBranchCoverageInstrumentationCallback(this_, compiler, instrumentationData)
	return this_
}

func InitBranchCoverageInstrumentationCallback(this_ *BranchCoverageInstrumentationCallback, compiler AbstractCompiler, instrumentationData *j2gort.OrderedMap[string, *FileInstrumentationData]) {
	_ = compiler
	_ = instrumentationData
	InitNodeTraversal_AbstractCfgCallback(&this_.NodeTraversal_AbstractCfgCallbackBase)
	this_.initFields_BranchCoverageInstrumentationCallback()
	this_.compiler = compiler
	this_.instrumentationData = instrumentationData
}

// Java static method: BranchCoverageInstrumentationCallback.createArrayName(com.google.javascript.jscomp.NodeTraversal traversal)
func BranchCoverageInstrumentationCallback_CreateArrayName(traversal *NodeTraversal) string {
	_ = traversal
	return ("JSCompiler_lcov_branch_data_" + j2gort.StrValString(CoverageUtil_CreateIdentifierFromText(traversal.GetSourceName())))
}

// Java method: BranchCoverageInstrumentationCallback.visit(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *BranchCoverageInstrumentationCallback) Visit_3(traversal *NodeTraversal, node pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = traversal
	_ = node
	_ = parent
	var fileName string = traversal.GetSourceName()
	_ = fileName
	if !(fileName == node.GetSourceFileName()) {
		return
	}
	if node.IsScript() {
		if this_.instrumentationData.Get(fileName) != nil {
			var toAddTo pkg_rhino.Node
			_ = toAddTo
			if node.HasChildren() && node.GetFirstChild().IsModuleBody() {
				toAddTo = node.GetFirstChild()
			} else {
				toAddTo = node
			}
			toAddTo.AddChildrenToFront(this_.priv_BranchCoverageInstrumentationCallback_newHeaderNode_6427(traversal, toAddTo).RemoveChildren())
			this_.compiler.ReportChangeToEnclosingScope(node)
			this_.priv_BranchCoverageInstrumentationCallback_instrumentBranchCoverage_6423(traversal, this_.instrumentationData.Get(fileName))
		}
	}
	if node.IsIf() {
		var cfg ControlFlowGraph[pkg_rhino.Node] = this_.GetControlFlowGraph(this_.compiler)
		_ = cfg
		var hasDefaultBlock bool = false
		_ = hasDefaultBlock
		{
			cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(cfg.GetOutEdges(node)))
			for cur_1.HasNext() {
				rawElem_2 := cur_1.Next()
				outEdge := j2gort.AsIface[pkg_graph.DiGraph_DiGraphEdge[pkg_rhino.Node, ControlFlowGraph_Branch]](rawElem_2)
				_ = outEdge
				if outEdge.GetValue() == StaticGet_ControlFlowGraph_Branch_ON_FALSE() {
					var destination pkg_rhino.Node = outEdge.GetDestination().GetValue()
					_ = destination
					if (((destination != 0) && destination.IsBlock()) && (destination.GetParent() != 0)) && destination.GetParent().IsIf() {
						hasDefaultBlock = true
					}
					break
				}
			}
		}
		if !(hasDefaultBlock) {
			this_.priv_BranchCoverageInstrumentationCallback_addDefaultBlock_6426(node)
		}
		j2gort.OrderedMapComputeIfAbsent(this_.instrumentationData, fileName, nil, func(lp_7 string) *FileInstrumentationData {
			k := lp_7
			_ = k
			return NewFileInstrumentationData(k, BranchCoverageInstrumentationCallback_CreateArrayName(traversal))
		})

		this_.priv_BranchCoverageInstrumentationCallback_processBranchInfo_6425(node, this_.instrumentationData.Get(fileName), this_.priv_BranchCoverageInstrumentationCallback_getChildrenBlocks_6422(node))
	} else if NodeUtil_IsLoopStructure(node) {
		var blocks j2gort.Ext_java_util_List = this_.priv_BranchCoverageInstrumentationCallback_getChildrenBlocks_6422(node)
		_ = blocks
		var cfg_2 ControlFlowGraph[pkg_rhino.Node] = this_.GetControlFlowGraph(this_.compiler)
		_ = cfg_2
		{
			cur_8 := j2gort.NewIterCursor(j2gort.CheckNotNull(cfg_2.GetOutEdges(node)))
			for cur_8.HasNext() {
				rawElem_9 := cur_8.Next()
				outEdge_2 := j2gort.AsIface[pkg_graph.DiGraph_DiGraphEdge[pkg_rhino.Node, ControlFlowGraph_Branch]](rawElem_9)
				_ = outEdge_2
				if outEdge_2.GetValue() == StaticGet_ControlFlowGraph_Branch_ON_FALSE() {
					var destination_2 pkg_rhino.Node = outEdge_2.GetDestination().GetValue()
					_ = destination_2
					if (destination_2 != 0) && destination_2.IsBlock() {
						blocks.Add_1(j2gort.BoxNode(destination_2))
					} else {
						var exitBlock pkg_rhino.Node = pkg_rhino.IR_Block_0()
						_ = exitBlock
						if (destination_2 != 0) && destination_2.GetParent().IsBlock() {
							exitBlock.InsertBefore(destination_2)
						} else {
							exitBlock.InsertAfter(outEdge_2.GetSource().GetValue())
						}
						blocks.Add_1(j2gort.BoxNode(exitBlock))
					}
				}
			}
		}
		j2gort.OrderedMapComputeIfAbsent(this_.instrumentationData, fileName, nil, func(lp_14 string) *FileInstrumentationData {
			k := lp_14
			_ = k
			return NewFileInstrumentationData(k, BranchCoverageInstrumentationCallback_CreateArrayName(traversal))
		})

		this_.priv_BranchCoverageInstrumentationCallback_processBranchInfo_6425(node, this_.instrumentationData.Get(fileName), blocks)
	}
}

// Java method: BranchCoverageInstrumentationCallback.getChildrenBlocks(com.google.javascript.rhino.Node node)
func (this_ *BranchCoverageInstrumentationCallback) priv_BranchCoverageInstrumentationCallback_getChildrenBlocks_6422(node pkg_rhino.Node) j2gort.Ext_java_util_List {
	_ = node
	var blocks []pkg_rhino.Node = []pkg_rhino.Node{}
	_ = blocks
	{
		var child pkg_rhino.Node = node.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			if child.IsBlock() {
				blocks = append(blocks, child)
			}
		}
	}
	return j2gort.WrapTypedArrayList(blocks)
}

// Java method: BranchCoverageInstrumentationCallback.instrumentBranchCoverage(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.jscomp.instrumentation.FileInstrumentationData data)
func (this_ *BranchCoverageInstrumentationCallback) priv_BranchCoverageInstrumentationCallback_instrumentBranchCoverage_6423(traversal *NodeTraversal, data *FileInstrumentationData) {
	_ = traversal
	_ = data
	var maxLine int32 = data.MaxBranchPresentLine()
	_ = maxLine
	var branchCoverageOffset int32 = int32(0)
	_ = branchCoverageOffset
	{
		var lineIdx int32 = int32(1)
		_ = lineIdx
		for ; lineIdx <= maxLine; lineIdx = lineIdx + 1 {
			var numBranches *j2gort.BoxInt = j2gort.ValueOfInt(data.GetNumBranches(lineIdx))
			_ = numBranches
			if numBranches != nil {
				{
					var branchIdx int32 = int32(1)
					_ = branchIdx
					for ; branchIdx <= j2gort.UnboxInt(numBranches); branchIdx = branchIdx + 1 {
						var block pkg_rhino.Node = data.GetBranchNode(lineIdx, branchIdx)
						_ = block
						block.AddChildToFront(this_.priv_BranchCoverageInstrumentationCallback_newBranchInstrumentationNode_6424(traversal, block, ((branchCoverageOffset + branchIdx) - int32(1))))
						this_.compiler.ReportChangeToEnclosingScope(block)
					}
				}
				branchCoverageOffset = (branchCoverageOffset + j2gort.UnboxInt(numBranches))
			}
		}
	}
}

// Java method: BranchCoverageInstrumentationCallback.newBranchInstrumentationNode(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node node, int idx)
func (this_ *BranchCoverageInstrumentationCallback) priv_BranchCoverageInstrumentationCallback_newBranchInstrumentationNode_6424(traversal *NodeTraversal, node pkg_rhino.Node, idx int32) pkg_rhino.Node {
	_ = traversal
	_ = node
	_ = idx
	var arrayName string = BranchCoverageInstrumentationCallback_CreateArrayName(traversal)
	_ = arrayName
	var getElemNode pkg_rhino.Node = pkg_rhino.IR_Getelem(pkg_rhino.IR_Name(arrayName), pkg_rhino.IR_Number(float64(idx)))
	_ = getElemNode
	var exprNode pkg_rhino.Node = pkg_rhino.IR_ExprResult(pkg_rhino.IR_Assign(getElemNode, pkg_rhino.IR_TrueNode()))
	_ = exprNode
	var fileName string = traversal.GetSourceName()
	_ = fileName
	j2gort.OrderedMapComputeIfAbsent(this_.instrumentationData, fileName, nil, func(lp_19 string) *FileInstrumentationData {
		k := lp_19
		_ = k
		return NewFileInstrumentationData(k, arrayName)
	})

	return exprNode.SrcrefTreeIfMissing(node)
}

// Java method: BranchCoverageInstrumentationCallback.processBranchInfo(com.google.javascript.rhino.Node branchNode, com.google.javascript.jscomp.instrumentation.FileInstrumentationData data, java.util.List<com.google.javascript.rhino.Node> blocks)
func (this_ *BranchCoverageInstrumentationCallback) priv_BranchCoverageInstrumentationCallback_processBranchInfo_6425(branchNode pkg_rhino.Node, data *FileInstrumentationData, blocks j2gort.Ext_java_util_List) {
	_ = branchNode
	_ = data
	_ = blocks
	var lineNumber int32 = branchNode.GetLineno()
	_ = lineNumber
	data.SetBranchPresent(lineNumber)
	var numBranches int32 = int32(0)
	_ = numBranches
	{
		cur_20 := j2gort.NewIterCursor(j2gort.CheckNotNull(blocks))
		for cur_20.HasNext() {
			child := pkg_rhino.Node(cur_20.NextNode())
			_ = child
			data.PutBranchNode(lineNumber, (numBranches + int32(1)), child)
			numBranches = numBranches + 1
		}
	}
	data.AddBranches(lineNumber, numBranches)
}

// Java method: BranchCoverageInstrumentationCallback.addDefaultBlock(com.google.javascript.rhino.Node node)
func (this_ *BranchCoverageInstrumentationCallback) priv_BranchCoverageInstrumentationCallback_addDefaultBlock_6426(node pkg_rhino.Node) pkg_rhino.Node {
	_ = node
	var defaultBlock pkg_rhino.Node = pkg_rhino.IR_Block_0()
	_ = defaultBlock
	node.AddChildToBack(defaultBlock)
	return defaultBlock.SrcrefTreeIfMissing(node)
}

// Java method: BranchCoverageInstrumentationCallback.newHeaderNode(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node srcref)
func (this_ *BranchCoverageInstrumentationCallback) priv_BranchCoverageInstrumentationCallback_newHeaderNode_6427(traversal *NodeTraversal, srcref pkg_rhino.Node) pkg_rhino.Node {
	_ = traversal
	_ = srcref
	var fileName string = traversal.GetSourceName()
	_ = fileName
	var data *FileInstrumentationData = this_.instrumentationData.Get(fileName)
	_ = data
	j2gort.CheckNotNullPtr(data)
	var objName string = "__jscov"
	_ = objName
	var nodes j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = nodes
	nodes.Add_1(j2gort.BoxNode(this_.priv_BranchCoverageInstrumentationCallback_newArrayDeclarationNode_6428(traversal)))
	nodes.Add_1(j2gort.BoxNode(pkg_rhino.IR_ExprResult(pkg_rhino.IR_Call(pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Getelem(pkg_rhino.IR_Name(objName), pkg_rhino.IR_String("branchesTaken")), "push"), []pkg_rhino.Node{pkg_rhino.IR_Name(BranchCoverageInstrumentationCallback_CreateArrayName(traversal))}))))
	nodes.Add_1(j2gort.BoxNode(pkg_rhino.IR_ExprResult(pkg_rhino.IR_Call(pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Getelem(pkg_rhino.IR_Name(objName), pkg_rhino.IR_String("branchPresent")), "push"), []pkg_rhino.Node{pkg_rhino.IR_String(data.GetBranchPresentAsHexString())}))))
	nodes.Add_1(j2gort.BoxNode(this_.priv_BranchCoverageInstrumentationCallback_newBranchesInLineNode_6429("JSCompiler_lcov_branchesInLine", data)))
	nodes.Add_1(j2gort.BoxNode(pkg_rhino.IR_ExprResult(pkg_rhino.IR_Call(pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Getelem(pkg_rhino.IR_Name(objName), pkg_rhino.IR_String("branchesInLine")), "push"), []pkg_rhino.Node{pkg_rhino.IR_Name("JSCompiler_lcov_branchesInLine")}))))
	nodes.Add_1(j2gort.BoxNode(pkg_rhino.IR_ExprResult(pkg_rhino.IR_Call(pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Getelem(pkg_rhino.IR_Name(objName), pkg_rhino.IR_String("fileNames")), "push"), []pkg_rhino.Node{pkg_rhino.IR_String(fileName)}))))
	return pkg_rhino.IR_Block_1_Ext_java_util_List(nodes).SrcrefTreeIfMissing(srcref)
}

// Java method: BranchCoverageInstrumentationCallback.newArrayDeclarationNode(com.google.javascript.jscomp.NodeTraversal traversal)
func (this_ *BranchCoverageInstrumentationCallback) priv_BranchCoverageInstrumentationCallback_newArrayDeclarationNode_6428(traversal *NodeTraversal) pkg_rhino.Node {
	_ = traversal
	return pkg_rhino.IR_Var_2(pkg_rhino.IR_Name(BranchCoverageInstrumentationCallback_CreateArrayName(traversal)), pkg_rhino.IR_Arraylit_1_Node([]pkg_rhino.Node{}))
}

// Java method: BranchCoverageInstrumentationCallback.newBranchesInLineNode(java.lang.String name, com.google.javascript.jscomp.instrumentation.FileInstrumentationData data)
func (this_ *BranchCoverageInstrumentationCallback) priv_BranchCoverageInstrumentationCallback_newBranchesInLineNode_6429(name string, data *FileInstrumentationData) pkg_rhino.Node {
	_ = name
	_ = data
	var assignments j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = assignments
	assignments.Add_1(j2gort.BoxNode(pkg_rhino.IR_Var_2(pkg_rhino.IR_Name(name), pkg_rhino.IR_Arraylit_1_Node([]pkg_rhino.Node{}))))
	var lineWithBranch int32 = int32(0)
	_ = lineWithBranch
	{
		var lineIdx int32 = int32(1)
		_ = lineIdx
		for ; lineIdx <= data.MaxBranchPresentLine(); lineIdx = lineIdx + 1 {
			var numBranches *j2gort.BoxInt = j2gort.ValueOfInt(data.GetNumBranches(lineIdx))
			_ = numBranches
			if (numBranches != nil) && (j2gort.UnboxInt(numBranches) > int32(0)) {
				var assignment pkg_rhino.Node = pkg_rhino.IR_ExprResult(pkg_rhino.IR_Assign(pkg_rhino.IR_Getelem(pkg_rhino.IR_Name(name), pkg_rhino.IR_Number(float64(func() int32 { old_21 := lineWithBranch; nw_22 := old_21 + 1; lineWithBranch = nw_22; return old_21 }()))), pkg_rhino.IR_Number(float64(j2gort.UnboxInt(numBranches)))))
				_ = assignment
				assignments.Add_1(j2gort.BoxNode(assignment.SrcrefTreeIfMissing(assignment)))
			}
		}
	}
	return pkg_rhino.IR_Block_1_Ext_java_util_List(assignments)
}
func (*BranchCoverageInstrumentationCallback) AlwaysShouldTraverse_() {}
