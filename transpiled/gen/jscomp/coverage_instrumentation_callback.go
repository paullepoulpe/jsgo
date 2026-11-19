/*
 * Copyright 2009 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/instrumentation/CoverageInstrumentationCallback.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.instrumentation.CoverageInstrumentationCallback
type CoverageInstrumentationCallback struct {
	instrumentationData *j2gort.OrderedMap[string, *FileInstrumentationData] `gson:"instrumentationData"`
	reach               CoverageInstrumentationPass_CoverageReach            `gson:"reach"`
	idHash_             atomic.Int32
}

func (this_ *CoverageInstrumentationCallback) BindSelf_(s any) {
	_ = s
}
func (this_ *CoverageInstrumentationCallback) IsJavaNil_() bool { return this_ == nil }
func (*CoverageInstrumentationCallback) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.instrumentation.CoverageInstrumentationCallback", "CoverageInstrumentationCallback"
}
func (this_ *CoverageInstrumentationCallback) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CoverageInstrumentationCallback) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.instrumentation.CoverageInstrumentationCallback", "com.google.javascript.jscomp.instrumentation.CoverageInstrumentationCallback", func(v_ any) bool {
		return j2gort.InstanceOf[*CoverageInstrumentationCallback](v_)
	})
}
func (*CoverageInstrumentationCallback) Implements_jscomp_NodeTraversal_Callback() {}

func StaticGet_CoverageInstrumentationCallback_ARRAY_NAME_PREFIX() string {
	return "JSCompiler_lcov_data_"
}
func (this_ *CoverageInstrumentationCallback) initFields_CoverageInstrumentationCallback() {
	_ = this_
}

// Java constructor: CoverageInstrumentationCallback(java.util.Map<java.lang.String,com.google.javascript.jscomp.instrumentation.FileInstrumentationData> instrumentationData, com.google.javascript.jscomp.instrumentation.CoverageInstrumentationPass.CoverageReach reach)
func NewCoverageInstrumentationCallback(instrumentationData *j2gort.OrderedMap[string, *FileInstrumentationData], reach CoverageInstrumentationPass_CoverageReach) *CoverageInstrumentationCallback {
	this_ := &CoverageInstrumentationCallback{}
	InitCoverageInstrumentationCallback(this_, instrumentationData, reach)
	return this_
}

func InitCoverageInstrumentationCallback(this_ *CoverageInstrumentationCallback, instrumentationData *j2gort.OrderedMap[string, *FileInstrumentationData], reach CoverageInstrumentationPass_CoverageReach) {
	_ = instrumentationData
	_ = reach
	this_.initFields_CoverageInstrumentationCallback()
	this_.instrumentationData = instrumentationData
	this_.reach = reach
}

// Java static method: CoverageInstrumentationCallback.getFileName(com.google.javascript.jscomp.NodeTraversal traversal)
func CoverageInstrumentationCallback_GetFileName(traversal *NodeTraversal) string {
	_ = traversal
	return traversal.GetSourceName()
}

// Java method: CoverageInstrumentationCallback.createArrayName(com.google.javascript.jscomp.NodeTraversal traversal)
func (this_ *CoverageInstrumentationCallback) priv_CoverageInstrumentationCallback_createArrayName_6467(traversal *NodeTraversal) string {
	_ = traversal
	return ("JSCompiler_lcov_data_" + j2gort.StrValString(CoverageUtil_CreateIdentifierFromText(CoverageInstrumentationCallback_GetFileName(traversal))))
}

// Java method: CoverageInstrumentationCallback.newInstrumentationNode(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node node)
func (this_ *CoverageInstrumentationCallback) priv_CoverageInstrumentationCallback_newInstrumentationNode_6468(traversal *NodeTraversal, node pkg_rhino.Node) pkg_rhino.Node {
	_ = traversal
	_ = node
	var lineNumber int32 = node.GetLineno()
	_ = lineNumber
	var arrayName string = this_.priv_CoverageInstrumentationCallback_createArrayName_6467(traversal)
	_ = arrayName
	var getElemNode pkg_rhino.Node = pkg_rhino.IR_Getelem(pkg_rhino.IR_Name(arrayName), pkg_rhino.IR_Number(float64((lineNumber - int32(1)))))
	_ = getElemNode
	var exprNode pkg_rhino.Node = pkg_rhino.IR_ExprResult(pkg_rhino.IR_Assign(getElemNode, pkg_rhino.IR_TrueNode()))
	_ = exprNode
	var fileName string = CoverageInstrumentationCallback_GetFileName(traversal)
	_ = fileName
	j2gort.OrderedMapComputeIfAbsent(this_.instrumentationData, fileName, nil, func(lp_5 string) *FileInstrumentationData {
		k := lp_5
		_ = k
		return NewFileInstrumentationData(k, arrayName)
	})

	this_.instrumentationData.Get(fileName).SetLineAsInstrumented(lineNumber)
	return exprNode.SrcrefTreeIfMissing(node)
}

// Java method: CoverageInstrumentationCallback.newArrayDeclarationNode(com.google.javascript.jscomp.NodeTraversal traversal)
func (this_ *CoverageInstrumentationCallback) priv_CoverageInstrumentationCallback_newArrayDeclarationNode_6469(traversal *NodeTraversal) pkg_rhino.Node {
	_ = traversal
	return pkg_rhino.IR_Var_2(pkg_rhino.IR_Name(this_.priv_CoverageInstrumentationCallback_createArrayName_6467(traversal)), pkg_rhino.IR_Arraylit_1_Node([]pkg_rhino.Node{}))
}

// Java method: CoverageInstrumentationCallback.newHeaderNode(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node srcref)
func (this_ *CoverageInstrumentationCallback) priv_CoverageInstrumentationCallback_newHeaderNode_6470(traversal *NodeTraversal, srcref pkg_rhino.Node) pkg_rhino.Node {
	_ = traversal
	_ = srcref
	var fileName string = CoverageInstrumentationCallback_GetFileName(traversal)
	_ = fileName
	var arrayName string = this_.priv_CoverageInstrumentationCallback_createArrayName_6467(traversal)
	_ = arrayName
	var data *FileInstrumentationData = this_.instrumentationData.Get(fileName)
	_ = data
	j2gort.CheckNotNullPtr(data)
	var objName string = "__jscov"
	_ = objName
	return pkg_rhino.IR_Block_1_Node_2([]pkg_rhino.Node{this_.priv_CoverageInstrumentationCallback_newArrayDeclarationNode_6469(traversal), pkg_rhino.IR_ExprResult(pkg_rhino.IR_Call(pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Getelem(pkg_rhino.IR_Name(objName), pkg_rhino.IR_String("executedLines")), "push"), []pkg_rhino.Node{pkg_rhino.IR_Name(arrayName)})), pkg_rhino.IR_ExprResult(pkg_rhino.IR_Call(pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Getelem(pkg_rhino.IR_Name(objName), pkg_rhino.IR_String("instrumentedLines")), "push"), []pkg_rhino.Node{pkg_rhino.IR_String(data.GetInstrumentedLinesAsHexString())})), pkg_rhino.IR_ExprResult(pkg_rhino.IR_Call(pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Getelem(pkg_rhino.IR_Name(objName), pkg_rhino.IR_String("fileNames")), "push"), []pkg_rhino.Node{pkg_rhino.IR_String(fileName)}))}).SrcrefTreeIfMissing(srcref)
}

// Java method: CoverageInstrumentationCallback.shouldTraverse(com.google.javascript.jscomp.NodeTraversal nodeTraversal, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CoverageInstrumentationCallback) ShouldTraverse_3(nodeTraversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = nodeTraversal
	_ = n
	_ = parent
	return ((!(n.IsScript()) || (n.GetJSDocInfo() == nil)) || !(n.GetJSDocInfo().IsNoCoverage()))
}

// Java method: CoverageInstrumentationCallback.visit(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *CoverageInstrumentationCallback) Visit_3(traversal *NodeTraversal, node pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = traversal
	_ = node
	_ = parent
	if node.IsScript() {
		var fileName string = CoverageInstrumentationCallback_GetFileName(traversal)
		_ = fileName
		if this_.instrumentationData.Get(fileName) != nil {
			if node.HasChildren() && node.GetFirstChild().IsModuleBody() {
				node = node.GetFirstChild()
			}
			node.AddChildrenToFront(this_.priv_CoverageInstrumentationCallback_newHeaderNode_6470(traversal, node).RemoveChildren())
		}
		traversal.ReportCodeChange_0()
		return
	}
	if ((this_.reach == StaticGet_CoverageInstrumentationPass_CoverageReach_CONDITIONAL()) && (parent != 0)) && parent.IsScript() {
		return
	}
	if node.IsFunction() && !(NodeUtil_GetFunctionBody(node).IsBlock()) {
		var returnValue pkg_rhino.Node = NodeUtil_GetFunctionBody(node)
		_ = returnValue
		var body pkg_rhino.Node = pkg_rhino.IR_Block_1_Node(pkg_rhino.IR_ReturnNode_1(returnValue.Detach()))
		_ = body
		body.SrcrefTreeIfMissing(returnValue)
		node.AddChildToBack(body)
	}
	if (((node.IsFunction() || node.IsWith()) || node.IsCase()) || node.IsDefaultCase()) || node.IsCatch() {
		var codeBlock pkg_rhino.Node = node.GetLastChild()
		_ = codeBlock
		codeBlock.AddChildToFront(this_.priv_CoverageInstrumentationCallback_newInstrumentationNode_6468(traversal, node))
		traversal.ReportCodeChange_0()
		return
	}
	if node.IsTry() {
		var firstChild pkg_rhino.Node = node.GetFirstChild()
		_ = firstChild
		firstChild.AddChildToFront(this_.priv_CoverageInstrumentationCallback_newInstrumentationNode_6468(traversal, node))
		traversal.ReportCodeChange_0()
		return
	}
	if ((parent != 0) && NodeUtil_IsStatementBlock(parent)) && !(node.IsModuleBody()) {
		this_.priv_CoverageInstrumentationCallback_newInstrumentationNode_6468(traversal, node).InsertBefore(node)
		traversal.ReportCodeChange_0()
		return
	}
}
func (*CoverageInstrumentationCallback) CanReleaseSyntacticScopes_() {}
