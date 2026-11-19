/*
 * Copyright 2012 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/AngularPass.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.AngularPass
type AngularPass struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler    AbstractCompiler           `gson:"compiler"`
	injectables []*AngularPass_NodeContext `gson:"injectables"`
}

func (this_ *AngularPass) BindSelf_(s any) {
	_ = s
}
func (this_ *AngularPass) IsJavaNil_() bool { return this_ == nil }
func (*AngularPass) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AngularPass", "AngularPass"
}
func (this_ *AngularPass) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*AngularPass) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AngularPass", "com.google.javascript.jscomp.AngularPass", func(v_ any) bool {
		return j2gort.InstanceOf[*AngularPass](v_)
	})
}
func (*AngularPass) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*AngularPass) Implements_jscomp_CompilerPass()                            {}
func (*AngularPass) Implements_jscomp_NodeTraversal_Callback()                  {}

func StaticGet_AngularPass_INJECT_PROPERTY_NAME() string { return "$inject" }

var static_AngularPass_INJECT_IN_NON_GLOBAL_OR_BLOCK_ERROR *DiagnosticType

func StaticGet_AngularPass_INJECT_IN_NON_GLOBAL_OR_BLOCK_ERROR() *DiagnosticType {
	clinit_AngularPass()
	return static_AngularPass_INJECT_IN_NON_GLOBAL_OR_BLOCK_ERROR
}
func StaticSet_AngularPass_INJECT_IN_NON_GLOBAL_OR_BLOCK_ERROR(v *DiagnosticType) {
	clinit_AngularPass()
	static_AngularPass_INJECT_IN_NON_GLOBAL_OR_BLOCK_ERROR = v
}

var static_AngularPass_INJECT_NON_FUNCTION_ERROR *DiagnosticType

func StaticGet_AngularPass_INJECT_NON_FUNCTION_ERROR() *DiagnosticType {
	clinit_AngularPass()
	return static_AngularPass_INJECT_NON_FUNCTION_ERROR
}
func StaticSet_AngularPass_INJECT_NON_FUNCTION_ERROR(v *DiagnosticType) {
	clinit_AngularPass()
	static_AngularPass_INJECT_NON_FUNCTION_ERROR = v
}

var static_AngularPass_INJECTED_FUNCTION_HAS_DESTRUCTURED_PARAM *DiagnosticType

func StaticGet_AngularPass_INJECTED_FUNCTION_HAS_DESTRUCTURED_PARAM() *DiagnosticType {
	clinit_AngularPass()
	return static_AngularPass_INJECTED_FUNCTION_HAS_DESTRUCTURED_PARAM
}
func StaticSet_AngularPass_INJECTED_FUNCTION_HAS_DESTRUCTURED_PARAM(v *DiagnosticType) {
	clinit_AngularPass()
	static_AngularPass_INJECTED_FUNCTION_HAS_DESTRUCTURED_PARAM = v
}

var static_AngularPass_INJECTED_FUNCTION_HAS_DEFAULT_VALUE *DiagnosticType

func StaticGet_AngularPass_INJECTED_FUNCTION_HAS_DEFAULT_VALUE() *DiagnosticType {
	clinit_AngularPass()
	return static_AngularPass_INJECTED_FUNCTION_HAS_DEFAULT_VALUE
}
func StaticSet_AngularPass_INJECTED_FUNCTION_HAS_DEFAULT_VALUE(v *DiagnosticType) {
	clinit_AngularPass()
	static_AngularPass_INJECTED_FUNCTION_HAS_DEFAULT_VALUE = v
}

var static_AngularPass_INJECTED_FUNCTION_ON_NON_QNAME *DiagnosticType

func StaticGet_AngularPass_INJECTED_FUNCTION_ON_NON_QNAME() *DiagnosticType {
	clinit_AngularPass()
	return static_AngularPass_INJECTED_FUNCTION_ON_NON_QNAME
}
func StaticSet_AngularPass_INJECTED_FUNCTION_ON_NON_QNAME(v *DiagnosticType) {
	clinit_AngularPass()
	static_AngularPass_INJECTED_FUNCTION_ON_NON_QNAME = v
}

var clinitDone_AngularPass uint32
var clinitState_AngularPass j2gort.ClinitState

func clinit_AngularPass() {
	if atomic.LoadUint32(&clinitDone_AngularPass) != 0 {
		return
	}
	clinitSlow_AngularPass()
}

func clinitSlow_AngularPass() {
	if !j2gort.EnterClinit(&clinitState_AngularPass, &clinitDone_AngularPass) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_AngularPass, &clinitDone_AngularPass)
	static_AngularPass_INJECT_IN_NON_GLOBAL_OR_BLOCK_ERROR = DiagnosticType_Error("JSC_INJECT_IN_NON_GLOBAL_OR_BLOCK_ERROR", "@ngInject only applies to functions defined in blocks or global scope.")
	static_AngularPass_INJECT_NON_FUNCTION_ERROR = DiagnosticType_Error("JSC_INJECT_NON_FUNCTION_ERROR", "@ngInject can only be used when defining a function or assigning a function expression.")
	static_AngularPass_INJECTED_FUNCTION_HAS_DESTRUCTURED_PARAM = DiagnosticType_Error("JSC_INJECTED_FUNCTION_HAS_DESTRUCTURED_PARAM", "@ngInject cannot be used on functions containing destructured parameter.")
	static_AngularPass_INJECTED_FUNCTION_HAS_DEFAULT_VALUE = DiagnosticType_Error("JSC_INJECTED_FUNCTION_HAS_DEFAULT_VALUE", "@ngInject cannot be used on functions containing default value.")
	static_AngularPass_INJECTED_FUNCTION_ON_NON_QNAME = DiagnosticType_Error("JSC_INJECTED_FUNCTION_ON_NON_QNAME", "@ngInject can only be used on qualified names.")
}

func (this_ *AngularPass) initFields_AngularPass() {
	_ = this_
	this_.injectables = []*AngularPass_NodeContext{}
}

// Java constructor: AngularPass(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewAngularPass(compiler AbstractCompiler) *AngularPass {
	clinit_AngularPass()
	this_ := &AngularPass{}
	InitAngularPass(this_, compiler)
	return this_
}

func InitAngularPass(this_ *AngularPass, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_AngularPass()
	this_.compiler = compiler
}

// Java method: AngularPass.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *AngularPass) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
	for _, entry := range this_.injectables {
		_ = entry
		var name string = entry.GetName()
		_ = name
		var fn pkg_rhino.Node = entry.GetFunctionNode()
		_ = fn
		var dependencies j2gort.Ext_java_util_List = this_.priv_AngularPass_createDependenciesList_593(fn)
		_ = dependencies
		if dependencies.IsEmpty() {
			continue
		}
		var dependenciesArray pkg_rhino.Node = pkg_rhino.IR_Arraylit_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](dependencies))
		_ = dependenciesArray
		var statement pkg_rhino.Node = pkg_rhino.IR_ExprResult(pkg_rhino.IR_Assign(pkg_rhino.IR_Getelem(NodeUtil_NewQName_2_AbstractCompiler_string(this_.compiler, name), pkg_rhino.IR_String("$inject")), dependenciesArray))
		_ = statement
		statement.SrcrefTree(entry.GetNode())
		statement.SetOriginalName(name)
		var newPropertyDoc *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_S_Builder()
		_ = newPropertyDoc
		newPropertyDoc.RecordVisibility(pkg_rhino.StaticGet_JSDocInfo_Visibility_PUBLIC())
		statement.GetFirstChild().SetJSDocInfo(newPropertyDoc.Build_0())
		var insertionPoint pkg_rhino.Node = entry.GetTarget()
		_ = insertionPoint
		var next pkg_rhino.Node = insertionPoint.GetNext()
		_ = next
		for ((next != 0) && NodeUtil_IsExprCall(next)) && (this_.compiler.GetCodingConvention().GetClassesDefinedByCall(next.GetFirstChild()) != nil) {
			insertionPoint = next
			next = insertionPoint.GetNext()
		}
		statement.InsertAfter(insertionPoint)
		this_.compiler.ReportChangeToEnclosingScope(statement)
	}

}

// Java method: AngularPass.createDependenciesList(com.google.javascript.rhino.Node n)
func (this_ *AngularPass) priv_AngularPass_createDependenciesList_593(n pkg_rhino.Node) j2gort.Ext_java_util_List {
	_ = n
	if !(n.IsFunction()) {
		j2gort.CheckArgument(false)
	}
	var params pkg_rhino.Node = NodeUtil_GetFunctionParameters(n)
	_ = params
	if params != 0 {
		return this_.priv_AngularPass_createStringsFromParamList_594(params)
	}
	return j2gort.NewExt_java_util_ArrayList()
}

// Java method: AngularPass.createStringsFromParamList(com.google.javascript.rhino.Node params)
func (this_ *AngularPass) priv_AngularPass_createStringsFromParamList_594(params pkg_rhino.Node) j2gort.Ext_java_util_List {
	_ = params
	var param pkg_rhino.Node = params.GetFirstChild()
	_ = param
	var names []pkg_rhino.Node = []pkg_rhino.Node{}
	_ = names
	for param != 0 {
		if param.IsName() {
			names = append(names, pkg_rhino.IR_String(param.GetString()).Srcref(param))
		} else if param.IsDestructuringPattern() {
			this_.compiler.Report_1(JSError_Make_3(param, StaticGet_AngularPass_INJECTED_FUNCTION_HAS_DESTRUCTURED_PARAM(), []string{}))
			return j2gort.NewExt_java_util_ArrayList()
		} else if param.IsDefaultValue() {
			this_.compiler.Report_1(JSError_Make_3(param, StaticGet_AngularPass_INJECTED_FUNCTION_HAS_DEFAULT_VALUE(), []string{}))
			return j2gort.NewExt_java_util_ArrayList()
		}
		param = param.GetNext()
	}
	return j2gort.WrapTypedArrayList(names)
}

// Java method: AngularPass.visit(com.google.javascript.jscomp.NodeTraversal unused, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *AngularPass) Visit_3(unused *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = unused
	_ = n
	_ = parent
	var docInfo *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = docInfo
	if (docInfo != nil) && docInfo.IsNgInject() {
		this_.priv_AngularPass_addNode_595(n)
	}
}

// Java method: AngularPass.addNode(com.google.javascript.rhino.Node n)
func (this_ *AngularPass) priv_AngularPass_addNode_595(n pkg_rhino.Node) {
	_ = n
	var injectAfter pkg_rhino.Node = pkg_rhino.Node(0)
	_ = injectAfter
	var fn pkg_rhino.Node = pkg_rhino.Node(0)
	_ = fn
	var name string = j2gort.NullString
	_ = name
	switch n.TokenSwitchOrd_() {
	case int32(51) /* Token.ASSIGN */ :
		if !(n.GetFirstChild().IsQualifiedName()) {
			this_.compiler.Report_1(JSError_Make_3(n, StaticGet_AngularPass_INJECTED_FUNCTION_ON_NON_QNAME(), []string{}))
			return
		}
		name = n.GetFirstChild().GetQualifiedName()
		fn = n
		for fn.IsAssign() {
			fn = fn.GetLastChild()
		}
		injectAfter = n.GetParent()
	case int32(73) /* Token.FUNCTION */ :
		name = NodeUtil_GetName(n)
		fn = n
		injectAfter = n
		if n.GetParent().IsAssign() && n.GetParent().GetJSDocInfo().IsNgInject() {
			return
		}
	case int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */ :
		name = n.GetFirstChild().GetString()
		fn = AngularPass_GetDeclarationRValue(n)
		injectAfter = n
	case int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
		var parent pkg_rhino.Node = n.GetParent()
		_ = parent
		if parent.IsClassMembers() {
			var classNode pkg_rhino.Node = parent.GetParent()
			_ = classNode
			var midPart string
			_ = midPart
			if n.IsStaticMember() {
				midPart = "."
			} else {
				midPart = ".prototype."
			}
			name = (j2gort.StrValString(NodeUtil_GetName(classNode)) + j2gort.StrValString(midPart) + j2gort.StrValString(n.GetString()))
			if NodeUtil_IsEs6ConstructorMemberFunctionDef(n) {
				name = NodeUtil_GetName(classNode)
			}
			fn = n.GetFirstChild()
			injectAfter = NodeUtil_GetEnclosingStatement(classNode)
		}
	default:
	}
	if (fn == 0) || !(fn.IsFunction()) {
		this_.compiler.Report_1(JSError_Make_3(n, StaticGet_AngularPass_INJECT_NON_FUNCTION_ERROR(), []string{}))
		return
	}
	if injectAfter.GetParent().IsExport() {
		injectAfter = injectAfter.GetParent()
	}
	if !(NodeUtil_IsStatementBlock(injectAfter.GetParent())) {
		this_.compiler.Report_1(JSError_Make_3(n, StaticGet_AngularPass_INJECT_IN_NON_GLOBAL_OR_BLOCK_ERROR(), []string{}))
		return
	}
	j2gort.CheckNotNullString(name)
	this_.injectables = append(this_.injectables, NewAngularPass_NodeContext(name, n, fn, injectAfter))
}

// Java static method: AngularPass.getDeclarationRValue(com.google.javascript.rhino.Node n)
func AngularPass_GetDeclarationRValue(n pkg_rhino.Node) pkg_rhino.Node {
	clinit_AngularPass()
	_ = n
	j2gort.CheckNotNullNode(n)
	if !(NodeUtil_IsNameDeclaration(n)) {
		j2gort.CheckArgument(false)
	}
	n = n.GetFirstFirstChild()
	if n == 0 {
		return pkg_rhino.Node(0)
	}
	for n.IsAssign() {
		n = n.GetLastChild()
	}
	return n
}

// Java class: com.google.javascript.jscomp.AngularPass.NodeContext
type AngularPass_NodeContext struct {
	name         string         `gson:"name"`
	node         pkg_rhino.Node `gson:"node"`
	functionNode pkg_rhino.Node `gson:"functionNode"`
	target       pkg_rhino.Node `gson:"target"`
	idHash_      atomic.Int32
}

func (this_ *AngularPass_NodeContext) BindSelf_(s any) {
	_ = s
}
func (this_ *AngularPass_NodeContext) IsJavaNil_() bool { return this_ == nil }
func (*AngularPass_NodeContext) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AngularPass$NodeContext", "NodeContext"
}
func (this_ *AngularPass_NodeContext) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AngularPass_NodeContext) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AngularPass.NodeContext", "com.google.javascript.jscomp.AngularPass$NodeContext", func(v_ any) bool {
		return j2gort.InstanceOf[*AngularPass_NodeContext](v_)
	})
}

func (this_ *AngularPass_NodeContext) initFields_AngularPass_NodeContext() {
	_ = this_
	this_.name = j2gort.NullString
}

// Java constructor: NodeContext(java.lang.String name, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node functionNode, com.google.javascript.rhino.Node target)
func NewAngularPass_NodeContext(name string, node pkg_rhino.Node, functionNode pkg_rhino.Node, target pkg_rhino.Node) *AngularPass_NodeContext {
	this_ := &AngularPass_NodeContext{}
	InitAngularPass_NodeContext(this_, name, node, functionNode, target)
	return this_
}

func InitAngularPass_NodeContext(this_ *AngularPass_NodeContext, name string, node pkg_rhino.Node, functionNode pkg_rhino.Node, target pkg_rhino.Node) {
	_ = name
	_ = node
	_ = functionNode
	_ = target
	this_.initFields_AngularPass_NodeContext()
	this_.name = name
	this_.node = node
	this_.functionNode = functionNode
	this_.target = target
}

// Java method: NodeContext.getName()
func (this_ *AngularPass_NodeContext) GetName() string {
	return this_.name
}

// Java method: NodeContext.getNode()
func (this_ *AngularPass_NodeContext) GetNode() pkg_rhino.Node {
	return this_.node
}

// Java method: NodeContext.getFunctionNode()
func (this_ *AngularPass_NodeContext) GetFunctionNode() pkg_rhino.Node {
	return this_.functionNode
}

// Java method: NodeContext.getTarget()
func (this_ *AngularPass_NodeContext) GetTarget() pkg_rhino.Node {
	return this_.target
}
func (*AngularPass) AlwaysShouldTraverse_()      {}
func (*AngularPass) CanReleaseSyntacticScopes_() {}
