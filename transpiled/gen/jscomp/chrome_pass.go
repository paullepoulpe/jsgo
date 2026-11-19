/*
 * Copyright 2014 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ChromePass.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.ChromePass
type ChromePass struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler       AbstractCompiler          `gson:"compiler"`
	createdObjects j2gort.OrderedSet[string] `gson:"createdObjects"`
}

func (this_ *ChromePass) BindSelf_(s any) {
	_ = s
}
func (this_ *ChromePass) IsJavaNil_() bool { return this_ == nil }
func (*ChromePass) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ChromePass", "ChromePass"
}
func (this_ *ChromePass) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*ChromePass) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ChromePass", "com.google.javascript.jscomp.ChromePass", func(v_ any) bool {
		return j2gort.InstanceOf[*ChromePass](v_)
	})
}
func (*ChromePass) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*ChromePass) Implements_jscomp_CompilerPass()                            {}
func (*ChromePass) Implements_jscomp_NodeTraversal_Callback()                  {}

func StaticGet_ChromePass_CR_DEFINE() string              { return "cr.define" }
func StaticGet_ChromePass_OBJECT_DEFINE_PROPERTY() string { return "Object.defineProperty" }
func StaticGet_ChromePass_CR_DEFINE_PROPERTY() string     { return "cr.defineProperty" }
func StaticGet_ChromePass_VIRTUAL_FILE() string           { return "<ChromePass.java>" }

var static_ChromePass_VIRTUAL_NODE pkg_rhino.Node

func StaticGet_ChromePass_VIRTUAL_NODE() pkg_rhino.Node {
	clinit_ChromePass()
	return static_ChromePass_VIRTUAL_NODE
}
func StaticSet_ChromePass_VIRTUAL_NODE(v pkg_rhino.Node) {
	clinit_ChromePass()
	static_ChromePass_VIRTUAL_NODE = v
}

var static_ChromePass_CR_PROPERTYKIND_JS pkg_rhino.QualifiedName

func StaticGet_ChromePass_CR_PROPERTYKIND_JS() pkg_rhino.QualifiedName {
	clinit_ChromePass()
	return static_ChromePass_CR_PROPERTYKIND_JS
}
func StaticSet_ChromePass_CR_PROPERTYKIND_JS(v pkg_rhino.QualifiedName) {
	clinit_ChromePass()
	static_ChromePass_CR_PROPERTYKIND_JS = v
}

var static_ChromePass_CR_PROPERTYKIND_ATTR pkg_rhino.QualifiedName

func StaticGet_ChromePass_CR_PROPERTYKIND_ATTR() pkg_rhino.QualifiedName {
	clinit_ChromePass()
	return static_ChromePass_CR_PROPERTYKIND_ATTR
}
func StaticSet_ChromePass_CR_PROPERTYKIND_ATTR(v pkg_rhino.QualifiedName) {
	clinit_ChromePass()
	static_ChromePass_CR_PROPERTYKIND_ATTR = v
}

var static_ChromePass_CR_PROPERTYKIND_BOOL_ATTR pkg_rhino.QualifiedName

func StaticGet_ChromePass_CR_PROPERTYKIND_BOOL_ATTR() pkg_rhino.QualifiedName {
	clinit_ChromePass()
	return static_ChromePass_CR_PROPERTYKIND_BOOL_ATTR
}
func StaticSet_ChromePass_CR_PROPERTYKIND_BOOL_ATTR(v pkg_rhino.QualifiedName) {
	clinit_ChromePass()
	static_ChromePass_CR_PROPERTYKIND_BOOL_ATTR = v
}
func StaticGet_ChromePass_CR_DEFINE_COMMON_EXPLANATION() string {
	return "It should be called like this: cr.define('name.space', function() '{ ... return {Export: Internal}; }');"
}

var static_ChromePass_CR_DEFINE_WRONG_NUMBER_OF_ARGUMENTS *DiagnosticType

func StaticGet_ChromePass_CR_DEFINE_WRONG_NUMBER_OF_ARGUMENTS() *DiagnosticType {
	clinit_ChromePass()
	return static_ChromePass_CR_DEFINE_WRONG_NUMBER_OF_ARGUMENTS
}
func StaticSet_ChromePass_CR_DEFINE_WRONG_NUMBER_OF_ARGUMENTS(v *DiagnosticType) {
	clinit_ChromePass()
	static_ChromePass_CR_DEFINE_WRONG_NUMBER_OF_ARGUMENTS = v
}

var static_ChromePass_CR_DEFINE_INVALID_FIRST_ARGUMENT *DiagnosticType

func StaticGet_ChromePass_CR_DEFINE_INVALID_FIRST_ARGUMENT() *DiagnosticType {
	clinit_ChromePass()
	return static_ChromePass_CR_DEFINE_INVALID_FIRST_ARGUMENT
}
func StaticSet_ChromePass_CR_DEFINE_INVALID_FIRST_ARGUMENT(v *DiagnosticType) {
	clinit_ChromePass()
	static_ChromePass_CR_DEFINE_INVALID_FIRST_ARGUMENT = v
}

var static_ChromePass_CR_DEFINE_INVALID_SECOND_ARGUMENT *DiagnosticType

func StaticGet_ChromePass_CR_DEFINE_INVALID_SECOND_ARGUMENT() *DiagnosticType {
	clinit_ChromePass()
	return static_ChromePass_CR_DEFINE_INVALID_SECOND_ARGUMENT
}
func StaticSet_ChromePass_CR_DEFINE_INVALID_SECOND_ARGUMENT(v *DiagnosticType) {
	clinit_ChromePass()
	static_ChromePass_CR_DEFINE_INVALID_SECOND_ARGUMENT = v
}

var static_ChromePass_CR_DEFINE_INVALID_RETURN_IN_FUNCTION *DiagnosticType

func StaticGet_ChromePass_CR_DEFINE_INVALID_RETURN_IN_FUNCTION() *DiagnosticType {
	clinit_ChromePass()
	return static_ChromePass_CR_DEFINE_INVALID_RETURN_IN_FUNCTION
}
func StaticSet_ChromePass_CR_DEFINE_INVALID_RETURN_IN_FUNCTION(v *DiagnosticType) {
	clinit_ChromePass()
	static_ChromePass_CR_DEFINE_INVALID_RETURN_IN_FUNCTION = v
}

var static_ChromePass_CR_DEFINE_PROPERTY_TOO_FEW_ARGUMENTS *DiagnosticType

func StaticGet_ChromePass_CR_DEFINE_PROPERTY_TOO_FEW_ARGUMENTS() *DiagnosticType {
	clinit_ChromePass()
	return static_ChromePass_CR_DEFINE_PROPERTY_TOO_FEW_ARGUMENTS
}
func StaticSet_ChromePass_CR_DEFINE_PROPERTY_TOO_FEW_ARGUMENTS(v *DiagnosticType) {
	clinit_ChromePass()
	static_ChromePass_CR_DEFINE_PROPERTY_TOO_FEW_ARGUMENTS = v
}

var static_ChromePass_CR_DEFINE_PROPERTY_INVALID_PROPERTY_KIND *DiagnosticType

func StaticGet_ChromePass_CR_DEFINE_PROPERTY_INVALID_PROPERTY_KIND() *DiagnosticType {
	clinit_ChromePass()
	return static_ChromePass_CR_DEFINE_PROPERTY_INVALID_PROPERTY_KIND
}
func StaticSet_ChromePass_CR_DEFINE_PROPERTY_INVALID_PROPERTY_KIND(v *DiagnosticType) {
	clinit_ChromePass()
	static_ChromePass_CR_DEFINE_PROPERTY_INVALID_PROPERTY_KIND = v
}

var clinitDone_ChromePass uint32
var clinitState_ChromePass j2gort.ClinitState

func clinit_ChromePass() {
	if atomic.LoadUint32(&clinitDone_ChromePass) != 0 {
		return
	}
	clinitSlow_ChromePass()
}

func clinitSlow_ChromePass() {
	if !j2gort.EnterClinit(&clinitState_ChromePass, &clinitDone_ChromePass) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ChromePass, &clinitDone_ChromePass)
	static_ChromePass_VIRTUAL_NODE = pkg_rhino.IR_Empty().SetStaticSourceFile(pkg_rhino.ToIface_StaticSourceFile(SourceFile_FromCode_3("<ChromePass.java>", "", pkg_rhino.StaticGet_StaticSourceFile_SourceKind_EXTERN())))
	static_ChromePass_CR_PROPERTYKIND_JS = pkg_rhino.QualifiedName_Of("cr.PropertyKind.JS")
	static_ChromePass_CR_PROPERTYKIND_ATTR = pkg_rhino.QualifiedName_Of("cr.PropertyKind.ATTR")
	static_ChromePass_CR_PROPERTYKIND_BOOL_ATTR = pkg_rhino.QualifiedName_Of("cr.PropertyKind.BOOL_ATTR")
	static_ChromePass_CR_DEFINE_WRONG_NUMBER_OF_ARGUMENTS = DiagnosticType_Error("JSC_CR_DEFINE_WRONG_NUMBER_OF_ARGUMENTS", ("cr.define() should have exactly 2 arguments. " + "It should be called like this: cr.define('name.space', function() '{ ... return {Export: Internal}; }');"))
	static_ChromePass_CR_DEFINE_INVALID_FIRST_ARGUMENT = DiagnosticType_Error("JSC_CR_DEFINE_INVALID_FIRST_ARGUMENT", ("Invalid first argument for cr.define(). " + "It should be called like this: cr.define('name.space', function() '{ ... return {Export: Internal}; }');"))
	static_ChromePass_CR_DEFINE_INVALID_SECOND_ARGUMENT = DiagnosticType_Error("JSC_CR_DEFINE_INVALID_SECOND_ARGUMENT", ("Invalid second argument for cr.define(). " + "It should be called like this: cr.define('name.space', function() '{ ... return {Export: Internal}; }');"))
	static_ChromePass_CR_DEFINE_INVALID_RETURN_IN_FUNCTION = DiagnosticType_Error("JSC_CR_DEFINE_INVALID_RETURN_IN_SECOND_ARGUMENT", ("Function passed as second argument of cr.define() should return the dictionary in its last statement. " + "It should be called like this: cr.define('name.space', function() '{ ... return {Export: Internal}; }');"))
	static_ChromePass_CR_DEFINE_PROPERTY_TOO_FEW_ARGUMENTS = DiagnosticType_Error("JSC_CR_DEFINE_PROPERTY_TOO_FEW_ARGUMENTS", "cr.defineProperty() requires at least 2 arguments.")
	static_ChromePass_CR_DEFINE_PROPERTY_INVALID_PROPERTY_KIND = DiagnosticType_Error("JSC_CR_DEFINE_PROPERTY_INVALID_PROPERTY_KIND", "Invalid cr.PropertyKind passed to cr.defineProperty(): expected ATTR, BOOL_ATTR or JS, found \"{0}\".")
}

func (this_ *ChromePass) initFields_ChromePass() {
	_ = this_
}

// Java constructor: ChromePass(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewChromePass(compiler AbstractCompiler) *ChromePass {
	clinit_ChromePass()
	this_ := &ChromePass{}
	InitChromePass(this_, compiler)
	return this_
}

func InitChromePass(this_ *ChromePass, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_ChromePass()
	this_.compiler = compiler
	this_.createdObjects = j2gort.OrderedSetFromStrCol(j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_java_util_Arrays_AsList([]any{"cr"})))
}

// Java method: ChromePass.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ChromePass) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: ChromePass.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *ChromePass) Visit_3(t *NodeTraversal, node pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = node
	_ = parent
	if node.IsCall() {
		var callee pkg_rhino.Node = node.GetFirstChild()
		_ = callee
		if callee.MatchesQualifiedName_1_string("cr.define") {
			this_.priv_ChromePass_visitNamespaceDefinition_1135(node, parent)
		} else if callee.MatchesQualifiedName_1_string("Object.defineProperty") || callee.MatchesQualifiedName_1_string("cr.defineProperty") {
			this_.priv_ChromePass_visitPropertyDefinition_1132(node, parent)
		}
	}
}

// Java method: ChromePass.visitPropertyDefinition(com.google.javascript.rhino.Node call, com.google.javascript.rhino.Node parent)
func (this_ *ChromePass) priv_ChromePass_visitPropertyDefinition_1132(call pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = call
	_ = parent
	var callee pkg_rhino.Node = call.GetFirstChild()
	_ = callee
	var isCrDefinePropertyCall bool = callee.MatchesQualifiedName_1_string("cr.defineProperty")
	_ = isCrDefinePropertyCall
	if isCrDefinePropertyCall {
		if call.GetChildCount() < int32(3) {
			this_.compiler.Report_1(JSError_Make_3(call, StaticGet_ChromePass_CR_DEFINE_PROPERTY_TOO_FEW_ARGUMENTS(), []string{}))
			return
		}
	} else if call.GetChildCount() < int32(4) {
		return
	}
	var target string = call.GetSecondChild().GetQualifiedName()
	_ = target
	if isCrDefinePropertyCall && !(javalang.StringEndsWith(j2gort.CheckNotNullString(target), ".prototype")) {
		target = (j2gort.StrValString(target) + ".prototype")
	}
	var property pkg_rhino.Node = call.GetChildAtIndex(int32(2))
	_ = property
	var getPropNode pkg_rhino.Node = NodeUtil_NewQName_2_AbstractCompiler_string(this_.compiler, (j2gort.StrValString(target) + "." + j2gort.StrValString(property.GetString()))).SrcrefTree(call)
	_ = getPropNode
	if isCrDefinePropertyCall {
		var propertyType pkg_rhino.Node = this_.priv_ChromePass_getTypeByCrPropertyKind_1133(call.GetChildAtIndex(int32(3)))
		_ = propertyType
		if propertyType != 0 {
			ChromePass_SetJsDocWithType(getPropNode, propertyType)
		} else {
			var sourceJsDocInfo *pkg_rhino.JSDocInfo = call.GetJSDocInfo()
			_ = sourceJsDocInfo
			if (sourceJsDocInfo != nil) && sourceJsDocInfo.HasType() {
				getPropNode.SetJSDocInfo(sourceJsDocInfo)
			} else {
				ChromePass_SetJsDocWithType(getPropNode, pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_QMARK()))
			}
		}
		call.SetJSDocInfo(nil)
	} else {
		ChromePass_SetJsDocWithType(getPropNode, pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_QMARK()))
	}
	var definitionNode pkg_rhino.Node = pkg_rhino.IR_ExprResult(getPropNode).Srcref(parent)
	_ = definitionNode
	definitionNode.InsertAfter(parent)
	this_.compiler.ReportChangeToEnclosingScope(func() pkg_rhino.Node {
		if isCrDefinePropertyCall {
			return call
		}
		return getPropNode
	}())
}

// Java method: ChromePass.getTypeByCrPropertyKind(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node propertyKind)
func (this_ *ChromePass) priv_ChromePass_getTypeByCrPropertyKind_1133(propertyKind pkg_rhino.Node) pkg_rhino.Node {
	_ = propertyKind
	if (propertyKind == 0) || StaticGet_ChromePass_CR_PROPERTYKIND_JS().Matches(propertyKind) {
		return pkg_rhino.Node(0)
	}
	if StaticGet_ChromePass_CR_PROPERTYKIND_ATTR().Matches(propertyKind) {
		return pkg_rhino.IR_String("string")
	}
	if StaticGet_ChromePass_CR_PROPERTYKIND_BOOL_ATTR().Matches(propertyKind) {
		return pkg_rhino.IR_String("boolean")
	}
	this_.compiler.Report_1(JSError_Make_3(propertyKind, StaticGet_ChromePass_CR_DEFINE_PROPERTY_INVALID_PROPERTY_KIND(), []string{propertyKind.GetQualifiedName()}))
	return pkg_rhino.Node(0)
}

// Java static method: ChromePass.setJsDocWithType(com.google.javascript.rhino.Node target, com.google.javascript.rhino.Node type)
func ChromePass_SetJsDocWithType(target pkg_rhino.Node, type_ pkg_rhino.Node) {
	clinit_ChromePass()
	_ = target
	_ = type_
	var builder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_S_Builder()
	_ = builder
	builder.RecordType(pkg_rhino.NewJSTypeExpression(type_.SrcrefTree(StaticGet_ChromePass_VIRTUAL_NODE()), "<ChromePass.java>"))
	target.SetJSDocInfo(builder.Build_0())
}

// Java method: ChromePass.createAndInsertObjectsForQualifiedName(com.google.javascript.rhino.Node scriptChild, java.lang.String namespace)
func (this_ *ChromePass) priv_ChromePass_createAndInsertObjectsForQualifiedName_1134(scriptChild pkg_rhino.Node, namespace string) {
	_ = scriptChild
	_ = namespace
	var objectsForQualifiedName j2gort.Ext_java_util_List = this_.priv_ChromePass_createObjectsForQualifiedName_1136(namespace)
	_ = objectsForQualifiedName
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(objectsForQualifiedName))
		for cur_1.HasNext() {
			n := pkg_rhino.Node(cur_1.NextNode())
			_ = n
			n.SrcrefTree(scriptChild).InsertBefore(scriptChild)
		}
	}
	if !(objectsForQualifiedName.IsEmpty()) {
		this_.compiler.ReportChangeToEnclosingScope(scriptChild)
	}
}

// Java method: ChromePass.visitNamespaceDefinition(com.google.javascript.rhino.Node crDefineCallNode, com.google.javascript.rhino.Node parent)
func (this_ *ChromePass) priv_ChromePass_visitNamespaceDefinition_1135(crDefineCallNode pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = crDefineCallNode
	_ = parent
	if crDefineCallNode.GetChildCount() != int32(3) {
		this_.compiler.Report_1(JSError_Make_3(crDefineCallNode, StaticGet_ChromePass_CR_DEFINE_WRONG_NUMBER_OF_ARGUMENTS(), []string{}))
	}
	var namespaceArg pkg_rhino.Node = crDefineCallNode.GetSecondChild()
	_ = namespaceArg
	var function pkg_rhino.Node = crDefineCallNode.GetChildAtIndex(int32(2))
	_ = function
	if !(namespaceArg.IsStringLit()) {
		this_.compiler.Report_1(JSError_Make_3(namespaceArg, StaticGet_ChromePass_CR_DEFINE_INVALID_FIRST_ARGUMENT(), []string{}))
		return
	}
	var namespace string = namespaceArg.GetString()
	_ = namespace
	this_.priv_ChromePass_createAndInsertObjectsForQualifiedName_1134(parent, namespace)
	if !(function.IsFunction()) {
		this_.compiler.Report_1(JSError_Make_3(namespaceArg, StaticGet_ChromePass_CR_DEFINE_INVALID_SECOND_ARGUMENT(), []string{}))
		return
	}
	var returnNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = returnNode
	var objectLit pkg_rhino.Node = pkg_rhino.Node(0)
	_ = objectLit
	var functionBlock pkg_rhino.Node = function.GetLastChild()
	_ = functionBlock
	{
		var iifeCond_1 pkg_rhino.
			Node
		{

			var asgn_2 pkg_rhino.Node = functionBlock.GetLastChild()
			returnNode = asgn_2
			iifeCond_1 = asgn_2
		}
		if ((((iifeCond_1) == 0) || !(returnNode.IsReturn())) || ((func() pkg_rhino.Node {
			var asgn_3 pkg_rhino.Node = returnNode.GetFirstChild()
			objectLit = asgn_3
			return asgn_3
		}()) == 0)) || !(objectLit.IsObjectLit()) {
			this_.compiler.Report_1(JSError_Make_3(namespaceArg, StaticGet_ChromePass_CR_DEFINE_INVALID_RETURN_IN_FUNCTION(), []string{}))
			return
		}
	}
	var exports j2gort.Ext_java_util_Map = ChromePass_ObjectLitToMap(objectLit)
	_ = exports
	NodeTraversal_Traverse(this_.compiler, functionBlock, NewChromePass_RenameInternalsToExternalsCallback(this_, namespace, exports, functionBlock))
}

// Java static method: ChromePass.objectLitToMap(com.google.javascript.rhino.Node objectLit)
func ChromePass_ObjectLitToMap(objectLit pkg_rhino.Node) j2gort.Ext_java_util_Map {
	clinit_ChromePass()
	_ = objectLit
	var res j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = res
	{
		var keyNode pkg_rhino.Node = objectLit.GetFirstChild()
		_ = keyNode
		for ; keyNode != 0; keyNode = keyNode.GetNext() {
			var key string = keyNode.GetString()
			_ = key
			var valueNode pkg_rhino.Node = keyNode.GetFirstChild()
			_ = valueNode
			if valueNode.IsName() {
				var value string = keyNode.GetFirstChild().GetString()
				_ = value
				j2gort.MapPutStr(res, value, j2gort.BoxString(key))
			}
		}
	}
	return res
}

// Java method: ChromePass.createObjectsForQualifiedName(java.lang.String namespace)
func (this_ *ChromePass) priv_ChromePass_createObjectsForQualifiedName_1136(namespace string) j2gort.Ext_java_util_List {
	_ = namespace
	var objects j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = objects
	var parts []string = j2gort.JStrSplit(namespace, "\\.")
	_ = parts
	this_.priv_ChromePass_createObjectIfNew_1137(objects, parts[int32(0)], true)
	if int32(len(parts)) >= int32(2) {
		var currPrefix j2gort.LocalStringBuilder
		currPrefix.AppendStr(parts[int32(0)])

		{
			var i int32 = int32(1)
			_ = i
			for ; i < int32(len(parts)); i = i + 1 {
				currPrefix.AppendStr(".").AppendStr(parts[i])
				this_.priv_ChromePass_createObjectIfNew_1137(objects, currPrefix.String(), false)
			}
		}
	}
	return objects
}

// Java method: ChromePass.createObjectIfNew(java.util.List<com.google.javascript.rhino.Node> objects, java.lang.String name, boolean needVar)
func (this_ *ChromePass) priv_ChromePass_createObjectIfNew_1137(objects j2gort.Ext_java_util_List, name string, needVar bool) {
	_ = objects
	_ = name
	_ = needVar
	if this_.createdObjects.Contains(name) {
		return
	}
	this_.createdObjects.Add(name)
	var v pkg_rhino.Node = NodeUtil_NewQName_2_AbstractCompiler_string(this_.compiler, name)
	_ = v
	var lhs pkg_rhino.Node = pkg_rhino.IR_Or(v.CloneTree_0(), pkg_rhino.IR_Objectlit([]pkg_rhino.Node{}))
	_ = lhs
	objects.Add_1(j2gort.BoxNode(func() pkg_rhino.Node {
		if needVar {
			return pkg_rhino.IR_Var_2(v, lhs)
		}
		return pkg_rhino.IR_ExprResult(pkg_rhino.IR_Assign(v, lhs))
	}()))
}

// Java class: com.google.javascript.jscomp.ChromePass.RenameInternalsToExternalsCallback
type ChromePass_RenameInternalsToExternalsCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_         *ChromePass
	namespaceName  string                   `gson:"namespaceName"`
	exports        j2gort.Ext_java_util_Map `gson:"exports"`
	namespaceBlock pkg_rhino.Node           `gson:"namespaceBlock"`
}

func (this_ *ChromePass_RenameInternalsToExternalsCallback) BindSelf_(s any) {
	_ = s
}
func (this_ *ChromePass_RenameInternalsToExternalsCallback) IsJavaNil_() bool { return this_ == nil }
func (*ChromePass_RenameInternalsToExternalsCallback) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ChromePass$RenameInternalsToExternalsCallback", "RenameInternalsToExternalsCallback"
}
func (this_ *ChromePass_RenameInternalsToExternalsCallback) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*ChromePass_RenameInternalsToExternalsCallback) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ChromePass.RenameInternalsToExternalsCallback", "com.google.javascript.jscomp.ChromePass$RenameInternalsToExternalsCallback", func(v_ any) bool {
		return j2gort.InstanceOf[*ChromePass_RenameInternalsToExternalsCallback](v_)
	})
}
func (this_ *ChromePass_RenameInternalsToExternalsCallback) SetOuter_(o *ChromePass) {
	this_.outer_ = o
}
func (this_ *ChromePass_RenameInternalsToExternalsCallback) GetOuter_() *ChromePass {
	return this_.outer_
}
func (*ChromePass_RenameInternalsToExternalsCallback) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*ChromePass_RenameInternalsToExternalsCallback) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *ChromePass_RenameInternalsToExternalsCallback) initFields_ChromePass_RenameInternalsToExternalsCallback() {
	_ = this_
	this_.namespaceName = j2gort.NullString
}

// Java constructor: RenameInternalsToExternalsCallback(java.lang.String namespaceName, java.util.Map<java.lang.String,java.lang.String> exports, com.google.javascript.rhino.Node namespaceBlock)
func NewChromePass_RenameInternalsToExternalsCallback(outer_ *ChromePass, namespaceName string, exports j2gort.Ext_java_util_Map, namespaceBlock pkg_rhino.Node) *ChromePass_RenameInternalsToExternalsCallback {
	this_ := &ChromePass_RenameInternalsToExternalsCallback{}
	this_.outer_ = outer_
	InitChromePass_RenameInternalsToExternalsCallback(this_, namespaceName, exports, namespaceBlock)
	return this_
}

func InitChromePass_RenameInternalsToExternalsCallback(this_ *ChromePass_RenameInternalsToExternalsCallback, namespaceName string, exports j2gort.Ext_java_util_Map, namespaceBlock pkg_rhino.Node) {
	_ = namespaceName
	_ = exports
	_ = namespaceBlock
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_ChromePass_RenameInternalsToExternalsCallback()
	this_.namespaceName = namespaceName
	this_.exports = exports
	this_.namespaceBlock = namespaceBlock
}

// Java method: RenameInternalsToExternalsCallback.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ChromePass_RenameInternalsToExternalsCallback) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if ((n.IsFunction() || n.IsClass()) && (parent == this_.namespaceBlock)) && j2gort.MapContainsKeyStr(this_.exports, n.GetFirstChild().GetString()) {
		var clone pkg_rhino.Node = n.CloneTree_0()
		_ = clone
		if clone.IsClass() {
			var className pkg_rhino.Node = clone.GetFirstChild()
			_ = className
			className.ReplaceWith(pkg_rhino.IR_Empty().Srcref(className))
		}
		NodeUtil_MarkNewScopesChanged_2_Node_AbstractCompiler(clone, this_.outer_.compiler)
		var exprResult pkg_rhino.Node = pkg_rhino.IR_ExprResult(pkg_rhino.IR_Assign(this_.priv_RenameInternalsToExternalsCallback_buildQualifiedName_1138(n.GetFirstChild()), clone).Srcref(n)).Srcref(n)
		_ = exprResult
		if n.GetJSDocInfo() != nil {
			exprResult.GetFirstChild().SetJSDocInfo(n.GetJSDocInfo())
			clone.SetJSDocInfo(nil)
		}
		n.ReplaceWith(exprResult)
		NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(n, this_.outer_.compiler)
		this_.outer_.compiler.ReportChangeToEnclosingScope(exprResult)
	} else if ((n.IsName() && j2gort.MapContainsKeyStr(this_.exports, n.GetString())) && !(parent.IsFunction())) && !(parent.IsClass()) {
		if NodeUtil_IsNameDeclaration(parent) {
			if parent.GetParent() == this_.namespaceBlock {
				var varContent pkg_rhino.Node = n.RemoveFirstChild()
				_ = varContent
				var exprResult_2 pkg_rhino.Node = pkg_rhino.Node(0)
				_ = exprResult_2
				if varContent == 0 {
					exprResult_2 = pkg_rhino.IR_ExprResult(this_.priv_RenameInternalsToExternalsCallback_buildQualifiedName_1138(n)).Srcref(parent)
				} else {
					exprResult_2 = pkg_rhino.IR_ExprResult(pkg_rhino.IR_Assign(this_.priv_RenameInternalsToExternalsCallback_buildQualifiedName_1138(n), varContent).Srcref(parent)).Srcref(parent)
				}
				if parent.GetJSDocInfo() != nil {
					exprResult_2.GetFirstChild().SetJSDocInfo(parent.GetJSDocInfo().Clone_0())
				}
				parent.ReplaceWith(exprResult_2)
				this_.outer_.compiler.ReportChangeToEnclosingScope(exprResult_2)
			}
		} else {
			var newNode pkg_rhino.Node = this_.priv_RenameInternalsToExternalsCallback_buildQualifiedName_1138(n)
			_ = newNode
			if n.GetJSDocInfo() != nil {
				newNode.SetJSDocInfo(n.GetJSDocInfo().Clone_0())
			}
			if parent.IsCall() {
				parent.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), false)
			}
			n.ReplaceWith(newNode)
			this_.outer_.compiler.ReportChangeToEnclosingScope(newNode)
		}
	}
}

// Java method: RenameInternalsToExternalsCallback.buildQualifiedName(com.google.javascript.rhino.Node internalName)
func (this_ *ChromePass_RenameInternalsToExternalsCallback) priv_RenameInternalsToExternalsCallback_buildQualifiedName_1138(internalName pkg_rhino.Node) pkg_rhino.Node {
	_ = internalName
	var externalName string = j2gort.CastString(j2gort.MapGetStr(this_.exports, internalName.GetString()))
	_ = externalName
	return NodeUtil_NewQName_2_AbstractCompiler_string(this_.outer_.compiler, (j2gort.StrValString(this_.namespaceName) + "." + j2gort.StrValString(externalName))).SrcrefTree(internalName)
}
func (*ChromePass) AlwaysShouldTraverse_()                                         {}
func (*ChromePass) CanReleaseSyntacticScopes_()                                    {}
func (*ChromePass_RenameInternalsToExternalsCallback) AlwaysShouldTraverse_()      {}
func (*ChromePass_RenameInternalsToExternalsCallback) CanReleaseSyntacticScopes_() {}
