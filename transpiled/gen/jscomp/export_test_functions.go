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
// Port of com/google/javascript/jscomp/ExportTestFunctions.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.ExportTestFunctions
type ExportTestFunctions struct {
	compiler               AbstractCompiler `gson:"compiler"`
	exportSymbolFunction   string           `gson:"exportSymbolFunction"`
	exportPropertyFunction string           `gson:"exportPropertyFunction"`
	idHash_                atomic.Int32
}

func (this_ *ExportTestFunctions) BindSelf_(s any) {
	_ = s
}
func (this_ *ExportTestFunctions) IsJavaNil_() bool { return this_ == nil }
func (*ExportTestFunctions) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ExportTestFunctions", "ExportTestFunctions"
}
func (this_ *ExportTestFunctions) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ExportTestFunctions) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ExportTestFunctions", "com.google.javascript.jscomp.ExportTestFunctions", func(v_ any) bool {
		return j2gort.InstanceOf[*ExportTestFunctions](v_)
	})
}
func (*ExportTestFunctions) Implements_jscomp_CompilerPass() {}

var static_ExportTestFunctions_TEST_FUNCTIONS_NAME_PATTERN *j2gort.Ext_java_util_regex_Pattern

func StaticGet_ExportTestFunctions_TEST_FUNCTIONS_NAME_PATTERN() *j2gort.Ext_java_util_regex_Pattern {
	clinit_ExportTestFunctions()
	return static_ExportTestFunctions_TEST_FUNCTIONS_NAME_PATTERN
}
func StaticSet_ExportTestFunctions_TEST_FUNCTIONS_NAME_PATTERN(v *j2gort.Ext_java_util_regex_Pattern) {
	clinit_ExportTestFunctions()
	static_ExportTestFunctions_TEST_FUNCTIONS_NAME_PATTERN = v
}
func StaticGet_ExportTestFunctions_GOOG_TESTING_TEST_SUITE() string { return "goog.testing.testSuite" }

var clinitDone_ExportTestFunctions uint32
var clinitState_ExportTestFunctions j2gort.ClinitState

func clinit_ExportTestFunctions() {
	if atomic.LoadUint32(&clinitDone_ExportTestFunctions) != 0 {
		return
	}
	clinitSlow_ExportTestFunctions()
}

func clinitSlow_ExportTestFunctions() {
	if !j2gort.EnterClinit(&clinitState_ExportTestFunctions, &clinitDone_ExportTestFunctions) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ExportTestFunctions, &clinitDone_ExportTestFunctions)
	static_ExportTestFunctions_TEST_FUNCTIONS_NAME_PATTERN = j2gort.Ext_java_util_regex_Pattern_Compile_1("^(?:((\\w+\\.)+prototype\\.||window\\.)*(setUpPage|setUp|shouldRunTests|tearDown|tearDownPage|test[\\w\\$]+))$")
}

func (this_ *ExportTestFunctions) initFields_ExportTestFunctions() {
	_ = this_
	this_.exportSymbolFunction = j2gort.NullString
	this_.exportPropertyFunction = j2gort.NullString
}

// Java constructor: ExportTestFunctions(com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.String exportSymbolFunction, java.lang.String exportPropertyFunction)
func NewExportTestFunctions(compiler AbstractCompiler, exportSymbolFunction string, exportPropertyFunction string) *ExportTestFunctions {
	clinit_ExportTestFunctions()
	this_ := &ExportTestFunctions{}
	InitExportTestFunctions(this_, compiler, exportSymbolFunction, exportPropertyFunction)
	return this_
}

func InitExportTestFunctions(this_ *ExportTestFunctions, compiler AbstractCompiler, exportSymbolFunction string, exportPropertyFunction string) {
	_ = compiler
	_ = exportSymbolFunction
	_ = exportPropertyFunction
	this_.initFields_ExportTestFunctions()
	j2gort.CheckNotNullIface(compiler)
	this_.compiler = compiler
	this_.exportSymbolFunction = exportSymbolFunction
	this_.exportPropertyFunction = exportPropertyFunction
}

// Java method: ExportTestFunctions.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ExportTestFunctions) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, NewExportTestFunctions_ExportTestFunctionsNodes(this_))
}

// Java method: ExportTestFunctions.exportTestFunctionAsSymbol(java.lang.String testFunctionName, com.google.javascript.rhino.Node node)
func (this_ *ExportTestFunctions) priv_ExportTestFunctions_exportTestFunctionAsSymbol_2786(testFunctionName string, node pkg_rhino.Node) {
	_ = testFunctionName
	_ = node
	var exportCallTarget pkg_rhino.Node = NodeUtil_NewQName_4_AbstractCompiler_string_Node_string(this_.compiler, this_.exportSymbolFunction, node, testFunctionName)
	_ = exportCallTarget
	var call pkg_rhino.Node = pkg_rhino.IR_Call(exportCallTarget, []pkg_rhino.Node{})
	_ = call
	if exportCallTarget.IsName() {
		call.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
	}
	call.AddChildToBack(pkg_rhino.IR_String(testFunctionName))
	call.AddChildToBack(NodeUtil_NewQName_4_AbstractCompiler_string_Node_string(this_.compiler, testFunctionName, node, testFunctionName))
	var expression pkg_rhino.Node = pkg_rhino.IR_ExprResult(call).SrcrefTreeIfMissing(node)
	_ = expression
	expression.InsertAfter(node)
	this_.compiler.ReportChangeToEnclosingScope(expression)
}

// Java method: ExportTestFunctions.exportTestFunctionAsProperty(com.google.javascript.rhino.Node fullyQualifiedFunctionName, com.google.javascript.rhino.Node node)
func (this_ *ExportTestFunctions) priv_ExportTestFunctions_exportTestFunctionAsProperty_2787(fullyQualifiedFunctionName pkg_rhino.Node, node pkg_rhino.Node) {
	_ = fullyQualifiedFunctionName
	_ = node
	if !(fullyQualifiedFunctionName.IsGetProp()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(fullyQualifiedFunctionName))
	}
	var testFunctionName string = NodeUtil_GetPrototypePropertyName(node.GetFirstChild())
	_ = testFunctionName
	if javalang.StringStartsWith(j2gort.CheckNotNullString(node.GetFirstChild().GetQualifiedName()), "window.") {
		testFunctionName = j2gort.JStrSubFrom(node.GetFirstChild().GetQualifiedName(), j2gort.JStrLen("window."))
	}
	var exportCall pkg_rhino.Node = pkg_rhino.IR_Call(NodeUtil_NewQName_2_AbstractCompiler_string(this_.compiler, this_.exportPropertyFunction), []pkg_rhino.Node{fullyQualifiedFunctionName.GetOnlyChild().CloneTree_0(), pkg_rhino.IR_String(testFunctionName), fullyQualifiedFunctionName.CloneTree_0()})
	_ = exportCall
	exportCall.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), exportCall.GetFirstChild().IsName())
	var export pkg_rhino.Node = pkg_rhino.IR_ExprResult(exportCall).SrcrefTree(node)
	_ = export
	export.InsertAfter(node.GetParent())
	this_.compiler.ReportChangeToEnclosingScope(export)
}

// Java static method: ExportTestFunctions.isTestSuiteArgument(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.NodeTraversal t)
func ExportTestFunctions_IsTestSuiteArgument(n pkg_rhino.Node, t *NodeTraversal) bool {
	clinit_ExportTestFunctions()
	_ = n
	_ = t
	return (((n.IsObjectLit() && n.GetParent().IsCall()) && n.IsSecondChildOf(n.GetParent())) && ExportTestFunctions_IsGoogTestingTestSuite(t, n.GetPrevious()))
}

// Java static method: ExportTestFunctions.isGoogTestingTestSuite(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node qname)
func ExportTestFunctions_IsGoogTestingTestSuite(t *NodeTraversal, qname pkg_rhino.Node) bool {
	clinit_ExportTestFunctions()
	_ = t
	_ = qname
	qname = NodeUtil_GetCallTargetResolvingIndirectCalls(qname.GetParent())
	if !(qname.IsQualifiedName()) {
		return false
	}
	var root pkg_rhino.Node = NodeUtil_GetRootOfQualifiedName_1_Node(qname)
	_ = root
	var rootName string = root.GetString()
	_ = rootName
	var rootVar *Var = j2gort.CastPtr[Var](t.GetScope().GetSlot(rootName))
	_ = rootVar
	if (rootVar == nil) || rootVar.IsGlobal() {
		return qname.MatchesQualifiedName_1_string("goog.testing.testSuite")
	} else if j2gort.CastPtr[Scope](rootVar.GetScope()).IsModuleScope() {
		var originalValue pkg_rhino.Node = rootVar.GetInitialValue()
		_ = originalValue
		return ((((originalValue != 0) && NodeUtil_IsGoogRequireCall(originalValue)) && originalValue.HasTwoChildren()) && func() bool {
			sr_1 := originalValue.GetSecondChild().GetString()
			se_2 := "goog.testing.testSuite"
			return j2gort.CheckNotNullString(sr_1) == se_2
		}())
	}
	return false
}

// Java static method: ExportTestFunctions.isTestFunction(java.lang.String functionName)
func ExportTestFunctions_IsTestFunction(functionName string) bool {
	clinit_ExportTestFunctions()
	_ = functionName
	return ((functionName != j2gort.NullString) && StaticGet_ExportTestFunctions_TEST_FUNCTIONS_NAME_PATTERN().Matcher(j2gort.BoxString(functionName)).Matches())
}

// Java class: com.google.javascript.jscomp.ExportTestFunctions.ExportTestFunctionsNodes
type ExportTestFunctions_ExportTestFunctionsNodes struct {
	NodeTraversal_AbstractShallowCallbackBase
	outer_ *ExportTestFunctions
}

func (this_ *ExportTestFunctions_ExportTestFunctionsNodes) BindSelf_(s any) {
	_ = s
}
func (this_ *ExportTestFunctions_ExportTestFunctionsNodes) IsJavaNil_() bool { return this_ == nil }
func (*ExportTestFunctions_ExportTestFunctionsNodes) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ExportTestFunctions$ExportTestFunctionsNodes", "ExportTestFunctionsNodes"
}
func (this_ *ExportTestFunctions_ExportTestFunctionsNodes) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractShallowCallbackBase.IdentityHash_()
}
func (*ExportTestFunctions_ExportTestFunctionsNodes) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ExportTestFunctions.ExportTestFunctionsNodes", "com.google.javascript.jscomp.ExportTestFunctions$ExportTestFunctionsNodes", func(v_ any) bool {
		return j2gort.InstanceOf[*ExportTestFunctions_ExportTestFunctionsNodes](v_)
	})
}
func (this_ *ExportTestFunctions_ExportTestFunctionsNodes) SetOuter_(o *ExportTestFunctions) {
	this_.outer_ = o
}
func (this_ *ExportTestFunctions_ExportTestFunctionsNodes) GetOuter_() *ExportTestFunctions {
	return this_.outer_
}
func (*ExportTestFunctions_ExportTestFunctionsNodes) Implements_jscomp_NodeTraversal_AbstractShallowCallback() {
}
func (*ExportTestFunctions_ExportTestFunctionsNodes) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *ExportTestFunctions_ExportTestFunctionsNodes) initFields_ExportTestFunctions_ExportTestFunctionsNodes() {
	_ = this_
}

// Java constructor: ExportTestFunctionsNodes()
func NewExportTestFunctions_ExportTestFunctionsNodes(outer_ *ExportTestFunctions) *ExportTestFunctions_ExportTestFunctionsNodes {
	this_ := &ExportTestFunctions_ExportTestFunctionsNodes{}
	this_.outer_ = outer_
	InitExportTestFunctions_ExportTestFunctionsNodes(this_)
	return this_
}

func InitExportTestFunctions_ExportTestFunctionsNodes(this_ *ExportTestFunctions_ExportTestFunctionsNodes) {
	InitNodeTraversal_AbstractShallowCallback(&this_.NodeTraversal_AbstractShallowCallbackBase)
	this_.initFields_ExportTestFunctions_ExportTestFunctionsNodes()
}

// Java method: ExportTestFunctionsNodes.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ExportTestFunctions_ExportTestFunctionsNodes) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if parent == 0 {
		return
	}
	if parent.IsScript() || parent.IsModuleBody() {
		if NodeUtil_IsFunctionDeclaration(n) {
			var functionName string = NodeUtil_GetName(n)
			_ = functionName
			if ExportTestFunctions_IsTestFunction(functionName) {
				this_.outer_.priv_ExportTestFunctions_exportTestFunctionAsSymbol_2786(functionName, n)
			}
		} else if this_.priv_ExportTestFunctionsNodes_isNameDeclaredFunction_2792(n) {
			var functionNode pkg_rhino.Node = n.GetFirstFirstChild()
			_ = functionNode
			var functionName_2 string = NodeUtil_GetName(functionNode)
			_ = functionName_2
			if ExportTestFunctions_IsTestFunction(functionName_2) {
				this_.outer_.priv_ExportTestFunctions_exportTestFunctionAsSymbol_2786(functionName_2, n)
			}
		} else if this_.priv_ExportTestFunctionsNodes_isNameDeclaredClass_2793(n) {
			var classNode pkg_rhino.Node = n.GetFirstFirstChild()
			_ = classNode
			var className string = NodeUtil_GetName(classNode)
			_ = className
			this_.priv_ExportTestFunctionsNodes_exportClass_2789(classNode, className, n)
		} else if n.IsClass() {
			this_.priv_ExportTestFunctionsNodes_exportClass_2788(n)
		}
	} else if NodeUtil_IsExprAssign(parent) {
		var grandparent pkg_rhino.Node = parent.GetParent()
		_ = grandparent
		if (grandparent != 0) && (grandparent.IsScript() || grandparent.IsModuleBody()) {
			var firstChild pkg_rhino.Node = n.GetFirstChild()
			_ = firstChild
			var lastChild pkg_rhino.Node = n.GetLastChild()
			_ = lastChild
			var nodeName string = firstChild.GetQualifiedName()
			_ = nodeName
			if lastChild.IsFunction() {
				if ExportTestFunctions_IsTestFunction(nodeName) {
					if n.GetFirstChild().IsName() {
						this_.outer_.priv_ExportTestFunctions_exportTestFunctionAsSymbol_2786(nodeName, parent)
					} else {
						this_.outer_.priv_ExportTestFunctions_exportTestFunctionAsProperty_2787(firstChild, n)
					}
				}
			} else if lastChild.IsClass() {
				this_.priv_ExportTestFunctionsNodes_exportClass_2789(lastChild, nodeName, parent)
			}
		}
	} else if ExportTestFunctions_IsTestSuiteArgument(n, t) {
		{
			var c pkg_rhino.Node = n.GetFirstChild()
			_ = c
			for c != 0 {
				var next pkg_rhino.Node = c.GetNext()
				_ = next
				if c.IsStringKey() && !(c.IsQuotedStringKey()) {
					c.SetQuotedStringKey()
					this_.outer_.compiler.ReportChangeToEnclosingScope(c)
				} else if c.IsMemberFunctionDef() {
					this_.priv_ExportTestFunctionsNodes_rewriteMemberDefInObjLit_2790(c, n)
				}
				c = next
			}
		}
	}
}

// Java method: ExportTestFunctionsNodes.exportClass(com.google.javascript.rhino.Node classNode)
func (this_ *ExportTestFunctions_ExportTestFunctionsNodes) priv_ExportTestFunctionsNodes_exportClass_2788(classNode pkg_rhino.Node) {
	_ = classNode
	var className string = NodeUtil_GetName(classNode)
	_ = className
	this_.priv_ExportTestFunctionsNodes_exportClass_2789(classNode, className, classNode)
}

// Java method: ExportTestFunctionsNodes.exportClass(com.google.javascript.rhino.Node classNode, java.lang.String className, com.google.javascript.rhino.Node addAfter)
func (this_ *ExportTestFunctions_ExportTestFunctionsNodes) priv_ExportTestFunctionsNodes_exportClass_2789(classNode pkg_rhino.Node, className string, addAfter pkg_rhino.Node) {
	_ = classNode
	_ = className
	_ = addAfter
	var classMembers pkg_rhino.Node = classNode.GetLastChild()
	_ = classMembers
	{
		var maybeMemberFunctionDef pkg_rhino.Node = classMembers.GetFirstChild()
		_ = maybeMemberFunctionDef
		for ; maybeMemberFunctionDef != 0; maybeMemberFunctionDef = maybeMemberFunctionDef.GetNext() {
			if maybeMemberFunctionDef.IsMemberFunctionDef() {
				var methodName string = maybeMemberFunctionDef.GetString()
				_ = methodName
				if ExportTestFunctions_IsTestFunction(methodName) {
					var functionRef string = (j2gort.StrValString(className) + ".prototype." + j2gort.StrValString(methodName))
					_ = functionRef
					var classRef string = (j2gort.StrValString(className) + ".prototype")
					_ = classRef
					var exportCallTarget pkg_rhino.Node = NodeUtil_NewQName_4_AbstractCompiler_string_Node_string(this_.outer_.compiler, this_.outer_.exportPropertyFunction, maybeMemberFunctionDef, methodName)
					_ = exportCallTarget
					var call pkg_rhino.Node = pkg_rhino.IR_Call(exportCallTarget, []pkg_rhino.Node{})
					_ = call
					if exportCallTarget.IsName() {
						call.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
					}
					call.AddChildToBack(NodeUtil_NewQName_4_AbstractCompiler_string_Node_string(this_.outer_.compiler, classRef, maybeMemberFunctionDef, classRef))
					call.AddChildToBack(pkg_rhino.IR_String(methodName))
					call.AddChildToBack(NodeUtil_NewQName_4_AbstractCompiler_string_Node_string(this_.outer_.compiler, functionRef, maybeMemberFunctionDef, functionRef))
					var expression pkg_rhino.Node = pkg_rhino.IR_ExprResult(call).SrcrefTreeIfMissing(maybeMemberFunctionDef)
					_ = expression
					expression.InsertAfter(addAfter)
					this_.outer_.compiler.ReportChangeToEnclosingScope(expression)
					addAfter = expression
				}
			}
		}
	}
}

// Java method: ExportTestFunctionsNodes.rewriteMemberDefInObjLit(com.google.javascript.rhino.Node memberDef, com.google.javascript.rhino.Node objLit)
func (this_ *ExportTestFunctions_ExportTestFunctionsNodes) priv_ExportTestFunctionsNodes_rewriteMemberDefInObjLit_2790(memberDef pkg_rhino.Node, objLit pkg_rhino.Node) {
	_ = memberDef
	_ = objLit
	var name string = memberDef.GetString()
	_ = name
	var stringKey pkg_rhino.Node = pkg_rhino.IR_StringKey_2(name, memberDef.RemoveFirstChild())
	_ = stringKey
	memberDef.ReplaceWith(stringKey)
	stringKey.SetQuotedStringKey()
	stringKey.SetJSDocInfo(memberDef.GetJSDocInfo())
	this_.outer_.compiler.ReportChangeToEnclosingScope(objLit)
}

// Java method: ExportTestFunctionsNodes.getNameDeclaredGrandchild(com.google.javascript.rhino.Node node)
func (this_ *ExportTestFunctions_ExportTestFunctionsNodes) priv_ExportTestFunctionsNodes_getNameDeclaredGrandchild_2791(node pkg_rhino.Node) pkg_rhino.Node {
	_ = node
	if !(NodeUtil_IsNameDeclaration(node)) {
		return pkg_rhino.Node(0)
	}
	return node.GetFirstFirstChild()
}

// Java method: ExportTestFunctionsNodes.isNameDeclaredFunction(com.google.javascript.rhino.Node node)
func (this_ *ExportTestFunctions_ExportTestFunctionsNodes) priv_ExportTestFunctionsNodes_isNameDeclaredFunction_2792(node pkg_rhino.Node) bool {
	_ = node
	var grandchild pkg_rhino.Node = this_.priv_ExportTestFunctionsNodes_getNameDeclaredGrandchild_2791(node)
	_ = grandchild
	return ((grandchild != 0) && grandchild.IsFunction())
}

// Java method: ExportTestFunctionsNodes.isNameDeclaredClass(com.google.javascript.rhino.Node node)
func (this_ *ExportTestFunctions_ExportTestFunctionsNodes) priv_ExportTestFunctionsNodes_isNameDeclaredClass_2793(node pkg_rhino.Node) bool {
	_ = node
	var grandchild pkg_rhino.Node = this_.priv_ExportTestFunctionsNodes_getNameDeclaredGrandchild_2791(node)
	_ = grandchild
	return ((grandchild != 0) && grandchild.IsClass())
}
func (*ExportTestFunctions_ExportTestFunctionsNodes) CanReleaseSyntacticScopes_() {}
