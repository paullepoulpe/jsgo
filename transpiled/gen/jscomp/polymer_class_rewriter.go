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
// Port of com/google/javascript/jscomp/PolymerClassRewriter.java (generated mechanically by j2go).

package jscomp

import (
	"strconv"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_polymer_class_rewriter_1 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_polymer_class_rewriter_2 j2gort.Ext_com_google_common_base_Predicate

func init() {
	lambdaSingleton_polymer_class_rewriter_1 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_1 uint32) bool {
		node := pkg_rhino.Node(rawNode_1)
		_ = node
		return (node.IsScript() || node.IsModuleBody())
	})
	lambdaSingleton_polymer_class_rewriter_2 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_4 uint32) bool {
		node := pkg_rhino.Node(rawNode_4)
		_ = node
		return (((node.IsScript() || node.IsModuleBody()) || PolymerClassRewriter_IsIIFE(node)) || PolymerClassRewriter_IsFunctionArgInGoogLoadModule(node))
	})
}

// Java class: com.google.javascript.jscomp.PolymerClassRewriter
type PolymerClassRewriter struct {
	compiler                   AbstractCompiler `gson:"compiler"`
	externsInsertionRef        pkg_rhino.Node   `gson:"externsInsertionRef"`
	idHash_                    atomic.Int32
	propertySinkExternInjected bool `gson:"propertySinkExternInjected"`
}

func (this_ *PolymerClassRewriter) BindSelf_(s any) {
	_ = s
}
func (this_ *PolymerClassRewriter) IsJavaNil_() bool { return this_ == nil }
func (*PolymerClassRewriter) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PolymerClassRewriter", "PolymerClassRewriter"
}
func (this_ *PolymerClassRewriter) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PolymerClassRewriter) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PolymerClassRewriter", "com.google.javascript.jscomp.PolymerClassRewriter", func(v_ any) bool {
		return j2gort.InstanceOf[*PolymerClassRewriter](v_)
	})
}

func StaticGet_PolymerClassRewriter_VIRTUAL_FILE() string { return "<PolymerClassRewriter.java>" }
func StaticGet_PolymerClassRewriter_POLYMER_ELEMENT_PROP_CONFIG() string {
	return "PolymerElementProperties"
}

var static_PolymerClassRewriter_IMPLICIT_GLOBAL_CONFLICT *DiagnosticType

func StaticGet_PolymerClassRewriter_IMPLICIT_GLOBAL_CONFLICT() *DiagnosticType {
	clinit_PolymerClassRewriter()
	return static_PolymerClassRewriter_IMPLICIT_GLOBAL_CONFLICT
}
func StaticSet_PolymerClassRewriter_IMPLICIT_GLOBAL_CONFLICT(v *DiagnosticType) {
	clinit_PolymerClassRewriter()
	static_PolymerClassRewriter_IMPLICIT_GLOBAL_CONFLICT = v
}

var static_PolymerClassRewriter_POLYMER_ELEMENT_CONFLICT *DiagnosticType

func StaticGet_PolymerClassRewriter_POLYMER_ELEMENT_CONFLICT() *DiagnosticType {
	clinit_PolymerClassRewriter()
	return static_PolymerClassRewriter_POLYMER_ELEMENT_CONFLICT
}
func StaticSet_PolymerClassRewriter_POLYMER_ELEMENT_CONFLICT(v *DiagnosticType) {
	clinit_PolymerClassRewriter()
	static_PolymerClassRewriter_POLYMER_ELEMENT_CONFLICT = v
}

var clinitDone_PolymerClassRewriter uint32
var clinitState_PolymerClassRewriter j2gort.ClinitState

func clinit_PolymerClassRewriter() {
	if atomic.LoadUint32(&clinitDone_PolymerClassRewriter) != 0 {
		return
	}
	clinitSlow_PolymerClassRewriter()
}

func clinitSlow_PolymerClassRewriter() {
	if !j2gort.EnterClinit(&clinitState_PolymerClassRewriter, &clinitDone_PolymerClassRewriter) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_PolymerClassRewriter, &clinitDone_PolymerClassRewriter)
	static_PolymerClassRewriter_IMPLICIT_GLOBAL_CONFLICT = DiagnosticType_Error("JSC_POLYMER_IMPLICIT_GLOBAL_CONFLICT", "Implicit global name for Polymer element conflicts with existing var {0}. Either give the element a lhs or rename {0}. (Or move to class-based Polymer 2 elements)")
	static_PolymerClassRewriter_POLYMER_ELEMENT_CONFLICT = DiagnosticType_Error("JSC_POLYMER_ELEMENT_CONFLICT", "Cannot generate correct types for Polymer call due to PolymerElement definition at {0}:{2}:{1}.\nRename the local PolymerElement to avoid shadowing the PolymerElement externs.")
}

func (this_ *PolymerClassRewriter) initFields_PolymerClassRewriter() {
	_ = this_

}

// Java constructor: PolymerClassRewriter(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewPolymerClassRewriter(compiler AbstractCompiler) *PolymerClassRewriter {
	clinit_PolymerClassRewriter()
	this_ := &PolymerClassRewriter{}
	InitPolymerClassRewriter(this_, compiler)
	return this_
}

func InitPolymerClassRewriter(this_ *PolymerClassRewriter, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_PolymerClassRewriter()
	this_.compiler = compiler
	this_.externsInsertionRef = compiler.GetSynthesizedExternsInput().GetAstRoot(compiler)
}

// Java static method: PolymerClassRewriter.isIIFE(com.google.javascript.rhino.Node n)
func PolymerClassRewriter_IsIIFE(n pkg_rhino.Node) bool {
	clinit_PolymerClassRewriter()
	_ = n
	return (n.IsCall() && n.GetFirstChild().IsFunction())
}

// Java static method: PolymerClassRewriter.isFunctionArgInGoogLoadModule(com.google.javascript.rhino.Node n)
func PolymerClassRewriter_IsFunctionArgInGoogLoadModule(n pkg_rhino.Node) bool {
	clinit_PolymerClassRewriter()
	_ = n
	if !(n.IsFunction()) {
		return false
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	return ((parent != 0) && NodeUtil_IsBundledGoogModuleCall(parent))
}

// Java method: PolymerClassRewriter.insertGeneratedDeclarationCodeToGlobalScope(com.google.javascript.rhino.Node enclosingNode, com.google.javascript.rhino.Node declarationCode)
func (this_ *PolymerClassRewriter) priv_PolymerClassRewriter_insertGeneratedDeclarationCodeToGlobalScope_4373(enclosingNode pkg_rhino.Node, declarationCode pkg_rhino.Node) {
	_ = enclosingNode
	_ = declarationCode
	switch enclosingNode.TokenSwitchOrd_() {
	case int32(123) /* Token.MODULE_BODY */ :
		{
			var insertionPoint pkg_rhino.Node = this_.priv_PolymerClassRewriter_getNodeForInsertion_4374(enclosingNode.GetParent())
			_ = insertionPoint
			insertionPoint.AddChildToFront(declarationCode)
			this_.compiler.ReportChangeToChangeScope(NodeUtil_GetEnclosingScript(insertionPoint))
		}
	case int32(94) /* Token.SCRIPT */ :
		{
			enclosingNode.AddChildToFront(declarationCode)
			this_.compiler.ReportChangeToChangeScope(NodeUtil_GetEnclosingScript(enclosingNode))
		}
	case int32(28) /* Token.CALL */ :
		{
			if !(PolymerClassRewriter_IsIIFE(enclosingNode)) {
				j2gort.CheckState(false)
			}
			var enclosingNodeForIIFE pkg_rhino.Node = NodeUtil_GetEnclosingNode(enclosingNode.GetParent(), lambdaSingleton_polymer_class_rewriter_1)
			_ = enclosingNodeForIIFE
			if enclosingNodeForIIFE.IsScript() {
				enclosingNodeForIIFE.AddChildToFront(declarationCode)
				this_.compiler.ReportChangeToChangeScope(NodeUtil_GetEnclosingScript(enclosingNodeForIIFE))
			} else {
				if !(enclosingNodeForIIFE.IsModuleBody()) {
					j2gort.CheckState(false)
				}
				var insertionPoint_2 pkg_rhino.Node = this_.priv_PolymerClassRewriter_getNodeForInsertion_4374(enclosingNodeForIIFE.GetParent())
				_ = insertionPoint_2
				insertionPoint_2.AddChildToFront(declarationCode)
				this_.compiler.ReportChangeToChangeScope(NodeUtil_GetEnclosingScript(insertionPoint_2))
			}
		}
	case int32(73) /* Token.FUNCTION */ :
		{
			if !(PolymerClassRewriter_IsFunctionArgInGoogLoadModule(enclosingNode)) {
				j2gort.CheckState(false)
			}
			var enclosingScript pkg_rhino.Node = NodeUtil_GetEnclosingScript(enclosingNode)
			_ = enclosingScript
			var insertionPoint_3 pkg_rhino.Node = this_.priv_PolymerClassRewriter_getNodeForInsertion_4374(enclosingScript)
			_ = insertionPoint_3
			insertionPoint_3.AddChildToFront(declarationCode)
			this_.compiler.ReportChangeToChangeScope(insertionPoint_3)
		}
	default:
		panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", "Enclosing node for Polymer is incorrect")))
	}
}

// Java method: PolymerClassRewriter.getNodeForInsertion(com.google.javascript.rhino.Node enclosingScript)
func (this_ *PolymerClassRewriter) priv_PolymerClassRewriter_getNodeForInsertion_4374(enclosingScript pkg_rhino.Node) pkg_rhino.Node {
	_ = enclosingScript
	if NodeUtil_IsFromTypeSummary(enclosingScript) {
		return this_.compiler.GetSynthesizedTypeSummaryInput().GetAstRoot(this_.compiler)
	} else {
		return this_.compiler.GetNodeForCodeInsertion(nil)
	}
}

// Java method: PolymerClassRewriter.insertGeneratedPropsAndBehaviorCode(com.google.javascript.rhino.Node enclosingNode, com.google.javascript.rhino.Node statements)
func (this_ *PolymerClassRewriter) priv_PolymerClassRewriter_insertGeneratedPropsAndBehaviorCode_4375(enclosingNode pkg_rhino.Node, statements pkg_rhino.Node) {
	_ = enclosingNode
	_ = statements
	switch enclosingNode.TokenSwitchOrd_() {
	case int32(123) /* Token.MODULE_BODY */ :
		{
			if enclosingNode.GetParent().GetBooleanProp(pkg_rhino.StaticGet_Node_GOOG_MODULE()) {
				var insertionPoint pkg_rhino.Node = PolymerClassRewriter_GetInsertionPointForGoogModule(enclosingNode)
				_ = insertionPoint
				enclosingNode.AddChildrenAfter(statements, insertionPoint)
			} else {
				enclosingNode.AddChildrenToFront(statements)
			}
		}
	case int32(94) /* Token.SCRIPT */ :
		enclosingNode.AddChildrenToBack(statements)
		this_.compiler.ReportChangeToChangeScope(NodeUtil_GetEnclosingScript(enclosingNode))
	case int32(28) /* Token.CALL */ :
		{
			if !(PolymerClassRewriter_IsIIFE(enclosingNode)) {
				j2gort.CheckState(false)
			}
			var functionNode pkg_rhino.Node = enclosingNode.GetFirstChild()
			_ = functionNode
			var functionBlock pkg_rhino.Node = functionNode.GetLastChild()
			_ = functionBlock
			functionBlock.AddChildrenToFront(statements)
		}
	case int32(73) /* Token.FUNCTION */ :
		if !(PolymerClassRewriter_IsFunctionArgInGoogLoadModule(enclosingNode)) {
			j2gort.CheckState(false)
		}
		var functionBlock_2 pkg_rhino.Node = enclosingNode.GetLastChild()
		_ = functionBlock_2
		var insertionPoint_2 pkg_rhino.Node = PolymerClassRewriter_GetInsertionPointForGoogModule(functionBlock_2)
		_ = insertionPoint_2
		if insertionPoint_2 != 0 {
			functionBlock_2.AddChildrenAfter(statements, insertionPoint_2)
		}
	default:
	}
}

// Java method: PolymerClassRewriter.rewritePolymerCall(com.google.javascript.jscomp.PolymerClassDefinition cls, com.google.javascript.jscomp.NodeTraversal traversal)
func (this_ *PolymerClassRewriter) RewritePolymerCall(cls *PolymerClassDefinition, traversal *NodeTraversal) {
	_ = cls
	_ = traversal
	var callParent pkg_rhino.Node = cls.definition.GetParent()
	_ = callParent
	var exprRoot pkg_rhino.Node
	_ = exprRoot
	if callParent.IsExprResult() {
		exprRoot = callParent
	} else {
		exprRoot = callParent.GetParent()
	}
	if !(NodeUtil_IsStatementParent(exprRoot.GetParent())) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(exprRoot.GetParent()))
	}
	var objLit pkg_rhino.Node = j2gort.CheckNotNullNode(cls.descriptor)
	_ = objLit
	var objLitDoc *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_S_Builder().ParseDocumentation()
	_ = objLitDoc
	var jsTypeExpression pkg_rhino.JSTypeExpression = pkg_rhino.NewJSTypeExpression(pkg_rhino.IR_String((j2gort.StrValString(cls.target.GetQualifiedName()) + ".prototype")).Srcref(exprRoot), exprRoot.GetSourceFileName())
	_ = jsTypeExpression
	objLitDoc.RecordLends(jsTypeExpression)
	objLit.SetJSDocInfo(objLitDoc.Build_0())
	this_.priv_PolymerClassRewriter_addTypesToFunctions_4379(objLit, cls.target.GetQualifiedName(), cls.defType)
	PolymerPassStaticUtils_SwitchDollarSignPropsToBrackets(objLit, this_.compiler)
	PolymerPassStaticUtils_QuoteListenerAndHostAttributeKeys(objLit, this_.compiler)
	PolymerPassStaticUtils_ProtectObserverAndPropertyFunctionKeys(objLit)
	{
		cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(cls.props))
		for cur_2.HasNext() {
			rawElem_3 := cur_2.Next()
			prop := j2gort.CastPtr[PolymerPass_MemberDefinition](rawElem_3)
			_ = prop
			if prop.value.IsObjectLit() {
				PolymerPassStaticUtils_SwitchDollarSignPropsToBrackets(prop.value, this_.compiler)
			}
		}
	}
	var propsAndBehaviorBlock pkg_rhino.Node = pkg_rhino.IR_Block_0()
	_ = propsAndBehaviorBlock
	var constructorDoc *pkg_rhino.JSDocInfo_Builder = this_.priv_PolymerClassRewriter_getConstructorDoc_4381(cls)
	_ = constructorDoc
	var ctorKey pkg_rhino.Node = cls.constructor.value.GetParent()
	_ = ctorKey
	if ctorKey != 0 {
		ctorKey.SetJSDocInfo(nil)
	}
	if !(traversal.InGlobalScope()) {
		var polymerElement *Var = j2gort.CastPtr[Var](traversal.GetScope().GetVar("PolymerElement"))
		_ = polymerElement
		if (polymerElement != nil) && !(j2gort.CastPtr[Scope](polymerElement.GetScope()).IsGlobal()) {
			var nameNode pkg_rhino.Node = polymerElement.GetNameNode()
			_ = nameNode
			this_.compiler.Report_1(JSError_Make_3(cls.constructor.value, StaticGet_PolymerClassRewriter_POLYMER_ELEMENT_CONFLICT(), []string{nameNode.GetSourceFileName(), strconv.FormatInt(int64(nameNode.GetLineno()), 10), strconv.FormatInt(int64(nameNode.GetCharno()), 10)}))
		}
	}
	var declarationCode pkg_rhino.Node = this_.priv_PolymerClassRewriter_generateDeclarationCode_4382(exprRoot, cls, constructorDoc, traversal)
	_ = declarationCode
	var basePath string = (j2gort.StrValString(cls.target.GetQualifiedName()) + ".prototype.")
	_ = basePath
	this_.priv_PolymerClassRewriter_appendBehaviorPropertiesToBlock_4385(cls, propsAndBehaviorBlock, basePath, false)
	this_.priv_PolymerClassRewriter_appendPropertiesToBlock_4386(cls.props, propsAndBehaviorBlock, basePath, false)
	this_.priv_PolymerClassRewriter_appendBehaviorMembersToBlock_4390(cls, propsAndBehaviorBlock)
	var readOnlyPropsAll j2gort.Ext_com_google_common_collect_ImmutableList = this_.priv_PolymerClassRewriter_parseReadOnlyProperties_4380(cls, propsAndBehaviorBlock)
	_ = readOnlyPropsAll
	this_.priv_PolymerClassRewriter_createExportsAndExterns_4392(cls, j2gort.AsIface[j2gort.Ext_java_util_List](readOnlyPropsAll))
	this_.priv_PolymerClassRewriter_removePropertyDocs_4388(objLit, cls.defType)
	var propsAndBehaviorCode pkg_rhino.Node = propsAndBehaviorBlock.RemoveChildren()
	_ = propsAndBehaviorCode
	var parent pkg_rhino.Node = exprRoot.GetParent()
	_ = parent
	if (!(traversal.InGlobalScope()) && cls.hasGeneratedLhs) && !(cls.target.IsGetProp()) {
		var enclosingNode pkg_rhino.Node = NodeUtil_GetEnclosingNode(parent, lambdaSingleton_polymer_class_rewriter_2)
		_ = enclosingNode
		this_.priv_PolymerClassRewriter_insertGeneratedDeclarationCodeToGlobalScope_4373(enclosingNode, declarationCode)
		if propsAndBehaviorCode != 0 {
			this_.priv_PolymerClassRewriter_insertGeneratedPropsAndBehaviorCode_4375(enclosingNode, propsAndBehaviorCode)
		}
	} else {
		var beforeRoot pkg_rhino.Node = exprRoot.GetPrevious()
		_ = beforeRoot
		if beforeRoot == 0 {
			if propsAndBehaviorCode != 0 {
				parent.AddChildrenToFront(propsAndBehaviorCode)
			}
			parent.AddChildToFront(declarationCode)
		} else {
			if propsAndBehaviorCode != 0 {
				parent.AddChildrenAfter(propsAndBehaviorCode, beforeRoot)
			}
			declarationCode.InsertAfter(beforeRoot)
		}
		this_.compiler.ReportChangeToEnclosingScope(parent)
	}
	if propsAndBehaviorCode != 0 {
		this_.compiler.ReportChangeToEnclosingScope(propsAndBehaviorCode)
	}
	if cls.features != nil {
		var scriptNode pkg_rhino.Node = NodeUtil_GetEnclosingScript(parent)
		_ = scriptNode
		var oldFeatures *pkg_parser.FeatureSet = j2gort.CastPtr[pkg_parser.FeatureSet](scriptNode.GetFeatureSet())
		_ = oldFeatures
		var newFeatures *pkg_parser.FeatureSet = oldFeatures.Union(cls.features)
		_ = newFeatures
		if !(newFeatures.Equals(j2gort.PtrToAny(oldFeatures))) {
			scriptNode.SetFeatureSet(j2gort.PtrToAny(newFeatures))
			this_.compiler.ReportChangeToChangeScope(scriptNode)
		}
	}
	if NodeUtil_IsNameDeclaration(exprRoot) {
		var assignExpr pkg_rhino.Node = PolymerClassRewriter_VarToAssign(exprRoot)
		_ = assignExpr
		exprRoot.ReplaceWith(assignExpr)
		this_.compiler.ReportChangeToEnclosingScope(assignExpr)
	}
	if cls.descriptor != 0 {
		var props pkg_rhino.Node = NodeUtil_GetFirstPropMatchingKey(cls.descriptor, "properties")
		_ = props
		if (props != 0) && props.IsObjectLit() {
			this_.priv_PolymerClassRewriter_addPropertiesConfigObjectReflection_4378(cls, props)
		}
	}
}

// Java method: PolymerClassRewriter.rewritePolymerClassDeclaration(com.google.javascript.rhino.Node clazz, com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.jscomp.PolymerClassDefinition cls)
func (this_ *PolymerClassRewriter) RewritePolymerClassDeclaration(clazz pkg_rhino.Node, traversal *NodeTraversal, cls *PolymerClassDefinition) {
	_ = clazz
	_ = traversal
	_ = cls
	if cls.descriptor != 0 {
		this_.priv_PolymerClassRewriter_addTypesToFunctions_4379(cls.descriptor, cls.target.GetQualifiedName(), cls.defType)
	}
	PolymerPassStaticUtils_SwitchDollarSignPropsToBrackets(NodeUtil_GetClassMembers(clazz), this_.compiler)
	{
		cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(cls.props))
		for cur_5.HasNext() {
			rawElem_6 := cur_5.Next()
			prop := j2gort.CastPtr[PolymerPass_MemberDefinition](rawElem_6)
			_ = prop
			if prop.value.IsObjectLit() {
				PolymerPassStaticUtils_SwitchDollarSignPropsToBrackets(prop.value, this_.compiler)
			}
		}
	}
	var block pkg_rhino.Node = pkg_rhino.IR_Block_0()
	_ = block
	this_.priv_PolymerClassRewriter_appendBehaviorPropertiesToBlock_4385(cls, block, (j2gort.StrValString(cls.target.GetQualifiedName()) + ".prototype."), false)
	this_.priv_PolymerClassRewriter_appendPropertiesToBlock_4386(cls.props, block, (j2gort.StrValString(cls.target.GetQualifiedName()) + ".prototype."), false)
	var readOnlyProps j2gort.Ext_com_google_common_collect_ImmutableList = this_.priv_PolymerClassRewriter_parseReadOnlyProperties_4380(cls, block)
	_ = readOnlyProps
	this_.priv_PolymerClassRewriter_createExportsAndExterns_4392(cls, j2gort.AsIface[j2gort.Ext_java_util_List](readOnlyProps))
	var jsDocInfoNode pkg_rhino.Node = NodeUtil_GetBestJSDocInfoNode(clazz)
	_ = jsDocInfoNode
	var classInfo *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(jsDocInfoNode.GetJSDocInfo())
	_ = classInfo
	var interfaceName string = cls.GetInterfaceName(this_.compiler.GetUniqueIdSupplier())
	_ = interfaceName
	var interfaceType pkg_rhino.JSTypeExpression = pkg_rhino.NewJSTypeExpression(pkg_rhino.NewNode_2(pkg_rhino.StaticGet_Token_BANG(), pkg_rhino.IR_String(interfaceName)).SrcrefTree(jsDocInfoNode), "<PolymerClassRewriter.java>")
	_ = interfaceType
	classInfo.RecordImplementedInterface(interfaceType)
	jsDocInfoNode.SetJSDocInfo(classInfo.Build_0())
	var insertAfterReference pkg_rhino.Node = NodeUtil_GetEnclosingStatement(clazz)
	_ = insertAfterReference
	if block.HasChildren() {
		this_.priv_PolymerClassRewriter_removePropertyDocs_4388(cls.descriptor, cls.defType)
		var newInsertAfterReference pkg_rhino.Node = block.GetLastChild()
		_ = newInsertAfterReference
		insertAfterReference.GetParent().AddChildrenAfter(block.RemoveChildren(), insertAfterReference)
		this_.compiler.ReportChangeToEnclosingScope(insertAfterReference)
		insertAfterReference = newInsertAfterReference
	}
	PolymerClassRewriter_AddReturnTypeIfMissing(cls, "is", pkg_rhino.NewJSTypeExpression(pkg_rhino.IR_String("string"), "<PolymerClassRewriter.java>"))
	var type_ pkg_rhino.Node = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_BANG())
	_ = type_
	var array pkg_rhino.Node = pkg_rhino.IR_String("Array")
	_ = array
	type_.AddChildToBack(array)
	var arrayTemplateType pkg_rhino.Node = pkg_rhino.NewNode_2(pkg_rhino.StaticGet_Token_BLOCK(), pkg_rhino.IR_String("string"))
	_ = arrayTemplateType
	array.AddChildToBack(arrayTemplateType)
	PolymerClassRewriter_AddReturnTypeIfMissing(cls, "observers", pkg_rhino.NewJSTypeExpression(type_, "<PolymerClassRewriter.java>"))
	PolymerClassRewriter_AddReturnTypeIfMissing(cls, "properties", pkg_rhino.NewJSTypeExpression(pkg_rhino.IR_String("PolymerElementProperties"), "<PolymerClassRewriter.java>"))
	if cls.descriptor != 0 {
		this_.priv_PolymerClassRewriter_convertSimpleObserverStringsToReferences_4394(cls)
		this_.priv_PolymerClassRewriter_addPropertiesConfigObjectReflection_4378(cls, cls.descriptor)
	}
	if cls.descriptor != 0 {
		this_.priv_PolymerClassRewriter_validateComputedPropertiesReflectionCalls_4395(cls)
		this_.priv_PolymerClassRewriter_validateComplexObserverReflectionCalls_4396(cls)
	}
}

// Java static method: PolymerClassRewriter.addReturnTypeIfMissing(com.google.javascript.jscomp.PolymerClassDefinition cls, java.lang.String getterPropName, com.google.javascript.rhino.JSTypeExpression jsType)
func PolymerClassRewriter_AddReturnTypeIfMissing(cls *PolymerClassDefinition, getterPropName string, jsType pkg_rhino.JSTypeExpression) {
	clinit_PolymerClassRewriter()
	_ = cls
	_ = getterPropName
	_ = jsType
	var classMembers pkg_rhino.Node = NodeUtil_GetClassMembers(cls.definition)
	_ = classMembers
	var getter pkg_rhino.Node = NodeUtil_GetFirstGetterMatchingKey(classMembers, getterPropName)
	_ = getter
	if getter != 0 {
		var info *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(getter)
		_ = info
		if (info == nil) || !(info.HasReturnType()) {
			var builder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(info)
			_ = builder
			builder.RecordReturnType(jsType)
			jsType.GetRoot().SrcrefTreeIfMissing(getter)
			getter.SetJSDocInfo(builder.Build_0())
		}
	}
}

// Java method: PolymerClassRewriter.addPropertiesConfigObjectReflection(com.google.javascript.jscomp.PolymerClassDefinition cls, com.google.javascript.rhino.Node propertiesLiteral)
func (this_ *PolymerClassRewriter) priv_PolymerClassRewriter_addPropertiesConfigObjectReflection_4378(cls *PolymerClassDefinition, propertiesLiteral pkg_rhino.Node) {
	_ = cls
	_ = propertiesLiteral
	j2gort.CheckNotNullNode(propertiesLiteral)
	if !(propertiesLiteral.IsObjectLit()) {
		j2gort.CheckState(false)
	}
	var parent pkg_rhino.Node = propertiesLiteral.GetParent()
	_ = parent
	var objReflectCall pkg_rhino.Node = pkg_rhino.IR_Call(NodeUtil_NewQName_2_AbstractCompiler_string(this_.compiler, "$jscomp.reflectObject"), []pkg_rhino.Node{cls.target.CloneTree_0(), propertiesLiteral.Detach()}).SrcrefTreeIfMissing(propertiesLiteral)
	_ = objReflectCall
	parent.AddChildToFront(objReflectCall)
	this_.compiler.ReportChangeToEnclosingScope(parent)
	this_.compiler.GetRuntimeJsLibManager().EnsureLibraryInjected("util/reflectobject", false)
}

// Java method: PolymerClassRewriter.addTypesToFunctions(com.google.javascript.rhino.Node objLit, java.lang.String thisType, com.google.javascript.jscomp.PolymerClassDefinition.DefinitionType defType)
func (this_ *PolymerClassRewriter) priv_PolymerClassRewriter_addTypesToFunctions_4379(objLit pkg_rhino.Node, thisType string, defType PolymerClassDefinition_DefinitionType) {
	_ = objLit
	_ = thisType
	_ = defType
	if !(objLit.IsObjectLit()) {
		j2gort.CheckState(false)
	}
	{
		var keyNode pkg_rhino.Node = objLit.GetFirstChild()
		_ = keyNode
		for ; keyNode != 0; keyNode = keyNode.GetNext() {
			var value pkg_rhino.Node = keyNode.GetLastChild()
			_ = value
			if (value != 0) && value.IsFunction() {
				var fnDoc *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(keyNode.GetJSDocInfo())
				_ = fnDoc
				fnDoc.RecordThisType(pkg_rhino.NewJSTypeExpression(pkg_rhino.NewNode_2(pkg_rhino.StaticGet_Token_BANG(), pkg_rhino.IR_String(thisType)).SrcrefTree(keyNode), "<PolymerClassRewriter.java>"))
				keyNode.SetJSDocInfo(fnDoc.Build_0())
			}
		}
	}
	{
		cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(PolymerPassStaticUtils_ExtractProperties(objLit, defType, this_.compiler, pkg_rhino.Node(0))))
		for cur_7.HasNext() {
			rawElem_8 := cur_7.Next()
			property := j2gort.CastPtr[PolymerPass_MemberDefinition](rawElem_8)
			_ = property
			if !(property.value.IsObjectLit()) {
				continue
			}
			var defaultValue pkg_rhino.Node = NodeUtil_GetFirstPropMatchingKey(property.value, "value")
			_ = defaultValue
			if (defaultValue == 0) || !(defaultValue.IsFunction()) {
				continue
			}
			var defaultValueKey pkg_rhino.Node = defaultValue.GetParent()
			_ = defaultValueKey
			var fnDoc_2 *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(defaultValueKey.GetJSDocInfo())
			_ = fnDoc_2
			fnDoc_2.RecordThisType(pkg_rhino.NewJSTypeExpression(pkg_rhino.NewNode_2(pkg_rhino.StaticGet_Token_BANG(), pkg_rhino.IR_String(thisType)).SrcrefTree(defaultValueKey), "<PolymerClassRewriter.java>"))
			fnDoc_2.RecordReturnType(PolymerPassStaticUtils_GetTypeFromProperty(property, this_.compiler))
			defaultValueKey.SetJSDocInfo(fnDoc_2.Build_0())
		}
	}
}

// Java method: PolymerClassRewriter.parseReadOnlyProperties(com.google.javascript.jscomp.PolymerClassDefinition cls, com.google.javascript.rhino.Node block)
func (this_ *PolymerClassRewriter) priv_PolymerClassRewriter_parseReadOnlyProperties_4380(cls *PolymerClassDefinition, block pkg_rhino.Node) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = cls
	_ = block
	var qualifiedPath string = (j2gort.StrValString(cls.target.GetQualifiedName()) + ".prototype.")
	_ = qualifiedPath
	var readOnlyProps []*PolymerPass_MemberDefinition = []*PolymerPass_MemberDefinition{}
	_ = readOnlyProps
	{
		cur_9 := j2gort.NewIterCursor(j2gort.CheckNotNull(cls.props))
		for cur_9.HasNext() {
			rawElem_10 := cur_9.Next()
			prop := j2gort.CastPtr[PolymerPass_MemberDefinition](rawElem_10)
			_ = prop
			if prop.value.IsObjectLit() {
				var readOnlyValue pkg_rhino.Node = NodeUtil_GetFirstPropMatchingKey(prop.value, "readOnly")
				_ = readOnlyValue
				if (readOnlyValue != 0) && readOnlyValue.IsTrue() {
					var setter pkg_rhino.Node = this_.priv_PolymerClassRewriter_makeReadOnlySetter_4391(prop, qualifiedPath)
					_ = setter
					setter.SrcrefTreeIfMissing(prop.name)
					block.AddChildToBack(setter)
					readOnlyProps = append(readOnlyProps, j2gort.CheckNotNullElem(prop))
				}
			}
		}
	}
	if cls.behaviorProps != nil {
		{
			cur_11 := j2gort.NewMapKeyCursor(cls.behaviorProps)
			for cur_11.HasNext() {
				rawElem_12 := cur_11.Next()
				prop_2 := j2gort.CastPtr[PolymerPass_MemberDefinition](rawElem_12)
				_ = prop_2
				if prop_2.value.IsObjectLit() {
					var readOnlyValue_2 pkg_rhino.Node = NodeUtil_GetFirstPropMatchingKey(prop_2.value, "readOnly")
					_ = readOnlyValue_2
					if (readOnlyValue_2 != 0) && readOnlyValue_2.IsTrue() {
						var setter_2 pkg_rhino.Node = this_.priv_PolymerClassRewriter_makeReadOnlySetter_4391(prop_2, qualifiedPath)
						_ = setter_2
						setter_2.SrcrefTreeIfMissing(prop_2.name)
						block.AddChildToBack(setter_2)
						readOnlyProps = append(readOnlyProps, j2gort.CheckNotNullElem(prop_2))
					}
				}
			}
		}
	}
	return j2gort.ImmListFromTypedSlice(readOnlyProps)
}

// Java method: PolymerClassRewriter.getConstructorDoc(com.google.javascript.jscomp.PolymerClassDefinition cls)
func (this_ *PolymerClassRewriter) priv_PolymerClassRewriter_getConstructorDoc_4381(cls *PolymerClassDefinition) *pkg_rhino.JSDocInfo_Builder {
	_ = cls
	var constructorDoc *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(cls.constructor.info)
	_ = constructorDoc
	constructorDoc.RecordConstructor()
	var baseType pkg_rhino.JSTypeExpression = pkg_rhino.NewJSTypeExpression(pkg_rhino.NewNode_2(pkg_rhino.StaticGet_Token_BANG(), pkg_rhino.IR_String(PolymerPassStaticUtils_GetPolymerElementType(cls))).SrcrefTree(cls.definition), "<PolymerClassRewriter.java>")
	_ = baseType
	constructorDoc.RecordBaseType(baseType)
	var interfaceName string = cls.GetInterfaceName(this_.compiler.GetUniqueIdSupplier())
	_ = interfaceName
	var interfaceType pkg_rhino.JSTypeExpression = pkg_rhino.NewJSTypeExpression(pkg_rhino.NewNode_2(pkg_rhino.StaticGet_Token_BANG(), pkg_rhino.IR_String(interfaceName)).SrcrefTree(cls.definition), "<PolymerClassRewriter.java>")
	_ = interfaceType
	constructorDoc.RecordImplementedInterface(interfaceType)
	return constructorDoc
}

// Java method: PolymerClassRewriter.generateDeclarationCode(com.google.javascript.rhino.Node exprRoot, com.google.javascript.jscomp.PolymerClassDefinition cls, com.google.javascript.rhino.JSDocInfo.Builder constructorDoc, com.google.javascript.jscomp.NodeTraversal traversal)
func (this_ *PolymerClassRewriter) priv_PolymerClassRewriter_generateDeclarationCode_4382(exprRoot pkg_rhino.Node, cls *PolymerClassDefinition, constructorDoc *pkg_rhino.JSDocInfo_Builder, traversal *NodeTraversal) pkg_rhino.Node {
	_ = exprRoot
	_ = cls
	_ = constructorDoc
	_ = traversal
	if cls.target.IsGetProp() {
		var assign pkg_rhino.Node = pkg_rhino.IR_Assign(cls.target.CloneTree_0(), cls.constructor.value.CloneTree_0())
		_ = assign
		NodeUtil_MarkNewScopesChanged_2_Node_AbstractCompiler(assign, this_.compiler)
		assign.SetJSDocInfo(constructorDoc.Build_0())
		var exprResult pkg_rhino.Node = pkg_rhino.IR_ExprResult(assign)
		_ = exprResult
		exprResult.SrcrefTreeIfMissing(cls.target)
		return exprResult
	} else {
		var var_ pkg_rhino.Node = pkg_rhino.IR_Var_2(cls.target.CloneTree_0(), cls.constructor.value.CloneTree_0())
		_ = var_
		NodeUtil_MarkNewScopesChanged_2_Node_AbstractCompiler(var_, this_.compiler)
		var_.SrcrefTreeIfMissing(exprRoot)
		var_.SetJSDocInfo(constructorDoc.Build_0())
		var name string = cls.target.GetString()
		_ = name
		var existingVar *Var = j2gort.CastPtr[Var](traversal.GetScope().GetSlot(name))
		_ = existingVar
		if (existingVar != nil) && cls.hasGeneratedLhs {
			this_.compiler.Report_1(JSError_Make_3(cls.constructor.value, StaticGet_PolymerClassRewriter_IMPLICIT_GLOBAL_CONFLICT(), []string{name}))
		}
		return var_
	}
}

// Java method: PolymerClassRewriter.replaceJSDocAndAddNewVars(com.google.javascript.jscomp.PolymerPass.MemberDefinition prop, com.google.javascript.rhino.JSTypeExpression propType, com.google.javascript.rhino.Node block)
func (this_ *PolymerClassRewriter) priv_PolymerClassRewriter_replaceJSDocAndAddNewVars_4383(prop *PolymerPass_MemberDefinition, propType pkg_rhino.JSTypeExpression, block pkg_rhino.Node) *pkg_rhino.JSDocInfo {
	_ = prop
	_ = propType
	_ = block
	var infoBuilder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(prop.info)
	_ = infoBuilder
	infoBuilder.RecordType(propType)
	var origInfo *pkg_rhino.JSDocInfo = infoBuilder.Build_0()
	_ = origInfo
	var propertyNames j2gort.Ext_com_google_common_collect_ImmutableSet = propType.GetRecordPropertyNames()
	_ = propertyNames
	this_.priv_PolymerClassRewriter_createVarsInExternsBlock_4387(block, propertyNames, propType, prop)
	var unknown pkg_rhino.JSTypeExpression = pkg_rhino.NewJSTypeExpression(pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_QMARK()), propType.GetSourceName())
	_ = unknown
	var newInfoBuilder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFromWithNewType(origInfo, unknown)
	_ = newInfoBuilder
	return newInfoBuilder.Build_0()
}

// Java method: PolymerClassRewriter.getPropertyNode(com.google.javascript.jscomp.PolymerPass.MemberDefinition prop, java.lang.String basePath)
func (this_ *PolymerClassRewriter) priv_PolymerClassRewriter_getPropertyNode_4384(prop *PolymerPass_MemberDefinition, basePath string) pkg_rhino.Node {
	_ = prop
	_ = basePath
	if prop.name.IsQuotedStringKey() {
		return pkg_rhino.Node(0)
	}
	var propertyNode pkg_rhino.Node = pkg_rhino.IR_ExprResult(NodeUtil_NewQName_2_AbstractCompiler_string(this_.compiler, (j2gort.StrValString(basePath) + j2gort.StrValString(prop.name.GetString()))))
	_ = propertyNode
	propertyNode.SrcrefTreeIfMissing(prop.name)
	return propertyNode
}

// Java method: PolymerClassRewriter.appendBehaviorPropertiesToBlock(com.google.javascript.jscomp.PolymerClassDefinition cls, com.google.javascript.rhino.Node block, java.lang.String basePath, boolean isExternsBlock)
func (this_ *PolymerClassRewriter) priv_PolymerClassRewriter_appendBehaviorPropertiesToBlock_4385(cls *PolymerClassDefinition, block pkg_rhino.Node, basePath string, isExternsBlock bool) {
	_ = cls
	_ = block
	_ = basePath
	_ = isExternsBlock
	if ((cls.behaviors == nil) || len(cls.behaviors) == 0) || (cls.behaviorProps == nil) {
		return
	}
	{
		cur_13 := j2gort.NewMapEntryCursor(cls.behaviorProps)
		for cur_13.NextEntry() {
			var behavior *PolymerBehaviorExtractor_BehaviorDefinition = j2gort.CastPtr[PolymerBehaviorExtractor_BehaviorDefinition](cur_13.GetValue())
			_ = behavior
			var prop *PolymerPass_MemberDefinition = j2gort.CastPtr[PolymerPass_MemberDefinition](cur_13.GetKey())
			_ = prop
			var propertyNode pkg_rhino.Node = this_.priv_PolymerClassRewriter_getPropertyNode_4384(prop, basePath)
			_ = propertyNode
			if propertyNode == 0 {
				continue
			}
			var propType pkg_rhino.JSTypeExpression = PolymerPassStaticUtils_GetTypeFromProperty(prop, this_.compiler)
			_ = propType
			if propType == (pkg_rhino.JSTypeExpression{}) {
				continue
			}
			var info *pkg_rhino.JSDocInfo = nil
			_ = info
			if isExternsBlock {
				info = this_.priv_PolymerClassRewriter_replaceJSDocAndAddNewVars_4383(prop, propType, block)
			} else {
				var infoBuilder *pkg_rhino.JSDocInfo_Builder = this_.priv_PolymerClassRewriter_getJSDocInfoBuilderForBehavior_4389(behavior, prop)
				_ = infoBuilder
				infoBuilder.RecordType(propType)
				info = infoBuilder.Build_0()
			}
			propertyNode.GetFirstChild().SetJSDocInfo(info)
			block.AddChildToBack(propertyNode)
		}
	}
}

// Java method: PolymerClassRewriter.appendPropertiesToBlock(java.util.List<com.google.javascript.jscomp.PolymerPass.MemberDefinition> props, com.google.javascript.rhino.Node block, java.lang.String basePath, boolean isExternsBlock)
func (this_ *PolymerClassRewriter) priv_PolymerClassRewriter_appendPropertiesToBlock_4386(props j2gort.Ext_java_util_List, block pkg_rhino.Node, basePath string, isExternsBlock bool) {
	_ = props
	_ = block
	_ = basePath
	_ = isExternsBlock
	{
		cur_14 := j2gort.NewIterCursor(j2gort.CheckNotNull(props))
		for cur_14.HasNext() {
			rawElem_15 := cur_14.Next()
			prop := j2gort.CastPtr[PolymerPass_MemberDefinition](rawElem_15)
			_ = prop
			var propertyNode pkg_rhino.Node = this_.priv_PolymerClassRewriter_getPropertyNode_4384(prop, basePath)
			_ = propertyNode
			if propertyNode == 0 {
				continue
			}
			var propType pkg_rhino.JSTypeExpression = PolymerPassStaticUtils_GetTypeFromProperty(prop, this_.compiler)
			_ = propType
			if propType == (pkg_rhino.JSTypeExpression{}) {
				continue
			}
			var info *pkg_rhino.JSDocInfo = nil
			_ = info
			if isExternsBlock {
				info = this_.priv_PolymerClassRewriter_replaceJSDocAndAddNewVars_4383(prop, propType, block)
			} else {
				var infoBuilder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(prop.info)
				_ = infoBuilder
				infoBuilder.RecordType(propType)
				info = infoBuilder.Build_0()
			}
			propertyNode.GetFirstChild().SetJSDocInfo(info)
			block.AddChildToBack(propertyNode)
		}
	}
}

// Java method: PolymerClassRewriter.createVarsInExternsBlock(com.google.javascript.rhino.Node block, com.google.common.collect.ImmutableSet<java.lang.String> propertyNames, com.google.javascript.rhino.JSTypeExpression propType, com.google.javascript.jscomp.PolymerPass.MemberDefinition prop)
func (this_ *PolymerClassRewriter) priv_PolymerClassRewriter_createVarsInExternsBlock_4387(block pkg_rhino.Node, propertyNames j2gort.Ext_com_google_common_collect_ImmutableSet, propType pkg_rhino.JSTypeExpression, prop *PolymerPass_MemberDefinition) {
	_ = block
	_ = propertyNames
	_ = propType
	_ = prop
	{
		cur_16 := j2gort.NewIterCursor(j2gort.CheckNotNull(propertyNames))
		for cur_16.HasNext() {
			propName := cur_16.NextStr()
			_ = propName
			var varName string = ("PolymerDummyVar" + j2gort.StrValString(j2gort.CastString(j2gort.AsIface[j2gort.Ext_java_util_function_Supplier](this_.compiler.GetUniqueNameIdSupplier()).Get())))
			_ = varName
			var n pkg_rhino.Node = pkg_rhino.Node_NewString_2(pkg_rhino.StaticGet_Token_NAME(), varName)
			_ = n
			var var_ pkg_rhino.Node = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_VAR())
			_ = var_
			var_.AddChildToBack(n)
			var newType pkg_rhino.JSTypeExpression = PolymerClassRewriter_CreateNewTypeExpressionForExtern(propName, propType.GetSourceName(), prop)
			_ = newType
			var oldInfoBuilder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(prop.info)
			_ = oldInfoBuilder
			var info *pkg_rhino.JSDocInfo = oldInfoBuilder.Build_0()
			_ = info
			var newInfo *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_CopyFromWithNewType(info, newType)
			_ = newInfo
			var_.SetJSDocInfo(newInfo.Build_0())
			block.AddChildToBack(var_)
		}
	}
}

// Java static method: PolymerClassRewriter.createNewTypeExpressionForExtern(java.lang.String propName, java.lang.String sourceName, com.google.javascript.jscomp.PolymerPass.MemberDefinition prop)
func PolymerClassRewriter_CreateNewTypeExpressionForExtern(propName string, sourceName string, prop *PolymerPass_MemberDefinition) pkg_rhino.JSTypeExpression {
	clinit_PolymerClassRewriter()
	_ = propName
	_ = sourceName
	_ = prop
	var leftCurly pkg_rhino.Node = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_LC())
	_ = leftCurly
	var leftBracket pkg_rhino.Node = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_LB())
	_ = leftBracket
	var colon pkg_rhino.Node = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_COLON())
	_ = colon
	var propertyName pkg_rhino.Node = pkg_rhino.Node_NewString_1(propName)
	_ = propertyName
	propertyName.SetToken(pkg_rhino.StaticGet_Token_STRING_KEY())
	var unknown pkg_rhino.Node = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_QMARK())
	_ = unknown
	colon.AddChildToBack(propertyName)
	colon.AddChildToBack(unknown)
	leftBracket.AddChildToBack(colon)
	leftCurly.AddChildToBack(leftBracket)
	leftCurly.SrcrefTreeIfMissing(prop.name)
	return pkg_rhino.NewJSTypeExpression(leftCurly, sourceName)
}

// Java method: PolymerClassRewriter.removePropertyDocs(com.google.javascript.rhino.Node objLit, com.google.javascript.jscomp.PolymerClassDefinition.DefinitionType defType)
func (this_ *PolymerClassRewriter) priv_PolymerClassRewriter_removePropertyDocs_4388(objLit pkg_rhino.Node, defType PolymerClassDefinition_DefinitionType) {
	_ = objLit
	_ = defType
	{
		cur_17 := j2gort.NewIterCursor(j2gort.CheckNotNull(PolymerPassStaticUtils_ExtractProperties(objLit, defType, this_.compiler, pkg_rhino.Node(0))))
		for cur_17.HasNext() {
			rawElem_18 := cur_17.Next()
			prop := j2gort.CastPtr[PolymerPass_MemberDefinition](rawElem_18)
			_ = prop
			prop.name.SetJSDocInfo(nil)
		}
	}
}

// Java method: PolymerClassRewriter.getJSDocInfoBuilderForBehavior(com.google.javascript.jscomp.PolymerBehaviorExtractor.BehaviorDefinition behavior, com.google.javascript.jscomp.PolymerPass.MemberDefinition behaviorFunctionOrProp)
func (this_ *PolymerClassRewriter) priv_PolymerClassRewriter_getJSDocInfoBuilderForBehavior_4389(behavior *PolymerBehaviorExtractor_BehaviorDefinition, behaviorFunctionOrProp *PolymerPass_MemberDefinition) *pkg_rhino.JSDocInfo_Builder {
	_ = behavior
	_ = behaviorFunctionOrProp
	var info *pkg_rhino.JSDocInfo_Builder = nil
	_ = info
	if (!(behavior.isGlobalDeclaration) && (behaviorFunctionOrProp.info != nil)) && behaviorFunctionOrProp.info.ContainsTypeDeclaration() {
		if behavior.behaviorModule != 0 {
			info = pkg_rhino.JSDocInfo_Builder_MaybeCopyFromAndReplaceNames(behaviorFunctionOrProp.info, behavior.GetModuleLocalNames(this_.compiler))
		} else {
			info = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(behaviorFunctionOrProp.info)
		}
	} else {
		info = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(behaviorFunctionOrProp.info)
	}
	return info
}

// Java method: PolymerClassRewriter.appendBehaviorMembersToBlock(com.google.javascript.jscomp.PolymerClassDefinition cls, com.google.javascript.rhino.Node block)
func (this_ *PolymerClassRewriter) priv_PolymerClassRewriter_appendBehaviorMembersToBlock_4390(cls *PolymerClassDefinition, block pkg_rhino.Node) {
	_ = cls
	_ = block
	var qualifiedPath string = (j2gort.StrValString(cls.target.GetQualifiedName()) + ".prototype.")
	_ = qualifiedPath
	var nameToExprResult j2gort.OrderedMap[string, pkg_rhino.Node] = j2gort.NewOrderedMap[string, pkg_rhino.Node]()
	_ = nameToExprResult
	{

		for _, behavior := range cls.behaviors {

			_ = behavior
			{
				cur_21 := j2gort.NewIterCursor(j2gort.CheckNotNull(behavior.functionsToCopy))
				for cur_21.HasNext() {
					rawElem_22 := cur_21.Next()
					behaviorFunction := j2gort.CastPtr[PolymerPass_MemberDefinition](rawElem_22)
					_ = behaviorFunction
					var fnName string = behaviorFunction.name.GetString()
					_ = fnName
					if NodeUtil_GetFirstPropMatchingKey(cls.descriptor, fnName) != 0 {
						continue
					}
					if nameToExprResult.ContainsKey(fnName) {
						nameToExprResult.Get(fnName).Detach()
					}
					var fnValue pkg_rhino.Node = behaviorFunction.value.CloneTree_0()
					_ = fnValue
					NodeUtil_MarkNewScopesChanged_2_Node_AbstractCompiler(fnValue, this_.compiler)
					var exprResult pkg_rhino.Node = pkg_rhino.IR_ExprResult(pkg_rhino.IR_Assign(NodeUtil_NewQName_2_AbstractCompiler_string(this_.compiler, (j2gort.StrValString(qualifiedPath)+j2gort.StrValString(fnName))), fnValue))
					_ = exprResult
					exprResult.SrcrefTreeIfMissing(behaviorFunction.name)
					var info *pkg_rhino.JSDocInfo_Builder = this_.priv_PolymerClassRewriter_getJSDocInfoBuilderForBehavior_4389(behavior, behaviorFunction)
					_ = info
					if (behaviorFunction.info != nil) && (behaviorFunction.info.GetVisibility() == pkg_rhino.StaticGet_JSDocInfo_Visibility_PROTECTED()) {
						info.OverwriteVisibility(pkg_rhino.StaticGet_JSDocInfo_Visibility_PUBLIC())
					}
					if !(behavior.isGlobalDeclaration) {
						var body pkg_rhino.Node = NodeUtil_GetFunctionBody(fnValue)
						_ = body
						if fnValue.IsArrowFunction() && !(NodeUtil_GetFunctionBody(fnValue).IsBlock()) {
							body.ReplaceWith(NodeUtil_NewUndefinedNode(body))
						} else {
							body.RemoveChildren()
						}
						var paramIndex int32 = int32(0)
						_ = paramIndex
						{
							var param pkg_rhino.Node = NodeUtil_GetFunctionParameters(fnValue).GetFirstChild()
							_ = param
							for param != 0 {
								var next pkg_rhino.Node = param.GetNext()
								_ = next
								PolymerClassRewriter_MakeParamSafe(param, func() int32 { old_23 := paramIndex; nw_24 := old_23 + 1; paramIndex = nw_24; return old_23 }())
								param = next
							}
						}
						if (behaviorFunction.info == nil) || (behaviorFunction.info.GetReturnType() == (pkg_rhino.JSTypeExpression{})) {
							info.RecordReturnType(pkg_rhino.NewJSTypeExpression(pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_QMARK()).Srcref(fnValue), fnValue.GetSourceFileName()))
						}
					}
					exprResult.GetFirstChild().SetJSDocInfo(info.Build_0())
					block.AddChildToBack(exprResult)
					nameToExprResult.Put(fnName, exprResult)
				}
			}
			{
				cur_25 := j2gort.NewIterCursor(j2gort.CheckNotNull(behavior.nonPropertyMembersToCopy))
				for cur_25.HasNext() {
					rawElem_26 := cur_25.Next()
					behaviorProp := j2gort.CastPtr[PolymerPass_MemberDefinition](rawElem_26)
					_ = behaviorProp
					var propName string = behaviorProp.name.GetString()
					_ = propName
					if nameToExprResult.ContainsKey(propName) {
						nameToExprResult.Get(propName).Detach()
					}
					var exprResult_2 pkg_rhino.Node = pkg_rhino.IR_ExprResult(NodeUtil_NewQName_2_AbstractCompiler_string(this_.compiler, (j2gort.StrValString(qualifiedPath) + j2gort.StrValString(propName))))
					_ = exprResult_2
					exprResult_2.SrcrefTree(behaviorProp.name)
					var info_2 *pkg_rhino.JSDocInfo_Builder = this_.priv_PolymerClassRewriter_getJSDocInfoBuilderForBehavior_4389(behavior, behaviorProp)
					_ = info_2
					if behaviorProp.name.IsGetterDef() {
						info_2 = pkg_rhino.JSDocInfo_S_Builder().ParseDocumentation()
						if (behaviorProp.info != nil) && (behaviorProp.info.GetReturnType() != (pkg_rhino.JSTypeExpression{})) {
							info_2.RecordType(behaviorProp.info.GetReturnType())
							if behaviorProp.enclosingModule != 0 {
								info_2 = pkg_rhino.JSDocInfo_Builder_MaybeCopyFromAndReplaceNames(info_2.Build_0(), behavior.GetModuleLocalNames(this_.compiler))
							}
						}
					}
					exprResult_2.GetFirstChild().SetJSDocInfo(info_2.Build_0())
					block.AddChildToBack(exprResult_2)
					nameToExprResult.Put(propName, exprResult_2)
				}
			}
		}
	}
}

// Java static method: PolymerClassRewriter.makeParamSafe(com.google.javascript.rhino.Node param, int index)
func PolymerClassRewriter_MakeParamSafe(param pkg_rhino.Node, index int32) {
	clinit_PolymerClassRewriter()
	_ = param
	_ = index
	if param.IsRest() {
		param = param.GetOnlyChild()
	} else if param.IsDefaultValue() {
		var value pkg_rhino.Node = param.GetSecondChild()
		_ = value
		value.ReplaceWith(NodeUtil_NewUndefinedNode(param))
		param = param.GetFirstChild()
	}
	if param.IsDestructuringPattern() {
		param.ReplaceWith(pkg_rhino.IR_Name(("param$polymer$" + j2gort.StrValInt32(index))).Srcref(param))
	}
}

// Java method: PolymerClassRewriter.makeReadOnlySetter(com.google.javascript.jscomp.PolymerPass.MemberDefinition prop, java.lang.String qualifiedPath)
func (this_ *PolymerClassRewriter) priv_PolymerClassRewriter_makeReadOnlySetter_4391(prop *PolymerPass_MemberDefinition, qualifiedPath string) pkg_rhino.Node {
	_ = prop
	_ = qualifiedPath
	var propName string = prop.name.GetString()
	_ = propName
	var setterName string = ("_set" + j2gort.StrValString(javalang.StringToUpperCaseASCII(j2gort.JStrSub(propName, int32(0), int32(1)))) + j2gort.StrValString(j2gort.JStrSubFrom(propName, int32(1))))
	_ = setterName
	var fnNode pkg_rhino.Node = pkg_rhino.IR_Function(pkg_rhino.IR_Name(""), pkg_rhino.IR_ParamList([]pkg_rhino.Node{pkg_rhino.IR_Name(propName)}), pkg_rhino.IR_Block_0())
	_ = fnNode
	this_.compiler.ReportChangeToChangeScope(fnNode)
	var exprResNode pkg_rhino.Node = pkg_rhino.IR_ExprResult(pkg_rhino.IR_Assign(NodeUtil_NewQName_2_AbstractCompiler_string(this_.compiler, (j2gort.StrValString(qualifiedPath)+j2gort.StrValString(setterName))), fnNode))
	_ = exprResNode
	var info *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_S_Builder().ParseDocumentation()
	_ = info
	info.RecordOverride()
	var propType pkg_rhino.JSTypeExpression = PolymerPassStaticUtils_GetTypeFromProperty(prop, this_.compiler)
	_ = propType
	info.RecordParameter(propName, propType)
	exprResNode.GetFirstChild().SetJSDocInfo(info.Build_0())
	return exprResNode
}

// Java method: PolymerClassRewriter.createExportsAndExterns(com.google.javascript.jscomp.PolymerClassDefinition cls, java.util.List<com.google.javascript.jscomp.PolymerPass.MemberDefinition> readOnlyProps)
func (this_ *PolymerClassRewriter) priv_PolymerClassRewriter_createExportsAndExterns_4392(cls *PolymerClassDefinition, readOnlyProps j2gort.Ext_java_util_List) {
	_ = cls
	_ = readOnlyProps
	var block pkg_rhino.Node = pkg_rhino.IR_Block_0()
	_ = block
	var interfaceName string = cls.GetInterfaceName(this_.compiler.GetUniqueIdSupplier())
	_ = interfaceName
	var fnNode pkg_rhino.Node = NodeUtil_EmptyFunction()
	_ = fnNode
	this_.compiler.ReportChangeToChangeScope(fnNode)
	var varNode pkg_rhino.Node = pkg_rhino.IR_Var_2(NodeUtil_NewQName_2_AbstractCompiler_string(this_.compiler, interfaceName), fnNode)
	_ = varNode
	var info *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_S_Builder().ParseDocumentation()
	_ = info
	info.RecordInterface()
	varNode.SetJSDocInfo(info.Build_0())
	block.AddChildToBack(varNode)
	var interfaceBasePath string = (j2gort.StrValString(interfaceName) + ".prototype.")
	_ = interfaceBasePath
	this_.priv_PolymerClassRewriter_appendBehaviorPropertiesToBlock_4385(cls, block, interfaceBasePath, true)
	this_.priv_PolymerClassRewriter_appendPropertiesToBlock_4386(cls.props, block, interfaceBasePath, true)
	var uniqueMethods j2gort.OrderedMap[string, *PolymerPass_MemberDefinition] = j2gort.NewOrderedMap[string, *PolymerPass_MemberDefinition]()
	_ = uniqueMethods
	if cls.behaviors != nil {
		{

			for _, behavior := range cls.behaviors {

				_ = behavior
				{
					cur_29 := j2gort.NewIterCursor(j2gort.CheckNotNull(behavior.functionsToCopy))
					for cur_29.HasNext() {
						rawElem_30 := cur_29.Next()
						method := j2gort.CastPtr[PolymerPass_MemberDefinition](rawElem_30)
						_ = method
						uniqueMethods.Put(method.name.GetString(), method)
					}
				}
			}
		}
	}
	{
		cur_31 := j2gort.NewIterCursor(j2gort.CheckNotNull(cls.methods))
		for cur_31.HasNext() {
			rawElem_32 := cur_31.Next()
			method_2 := j2gort.CastPtr[PolymerPass_MemberDefinition](rawElem_32)
			_ = method_2
			uniqueMethods.Put(method_2.name.GetString(), method_2)
		}
	}
	{
		cur_33 := uniqueMethods.IterValues()
		for cur_33.HasNext() {

			method_3 := cur_33.
				Next()
			_ = method_3
			this_.priv_PolymerClassRewriter_addMethodToObjectExternsUsingExportAnnotation_4393(cls, method_3)
		}
	}
	{
		cur_35 := j2gort.NewIterCursor(j2gort.CheckNotNull(readOnlyProps))
		for cur_35.HasNext() {
			rawElem_36 := cur_35.Next()
			prop := j2gort.CastPtr[PolymerPass_MemberDefinition](rawElem_36)
			_ = prop
			var propName string = prop.name.GetString()
			_ = propName
			var setterName string = ("_set" + j2gort.StrValString(javalang.StringToUpperCaseASCII(j2gort.JStrSub(propName, int32(0), int32(1)))) + j2gort.StrValString(j2gort.JStrSubFrom(propName, int32(1))))
			_ = setterName
			var setterExprNode pkg_rhino.Node = pkg_rhino.IR_ExprResult(NodeUtil_NewQName_2_AbstractCompiler_string(this_.compiler, (j2gort.StrValString(interfaceBasePath) + j2gort.StrValString(setterName))))
			_ = setterExprNode
			var setterInfo *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_S_Builder().ParseDocumentation()
			_ = setterInfo
			var propType pkg_rhino.JSTypeExpression = PolymerPassStaticUtils_GetTypeFromProperty(prop, this_.compiler)
			_ = propType
			var unknown pkg_rhino.JSTypeExpression = pkg_rhino.NewJSTypeExpression(pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_QMARK()), propType.GetSourceName())
			_ = unknown
			setterInfo.RecordParameter(propName, unknown)
			setterExprNode.GetFirstChild().SetJSDocInfo(setterInfo.Build_0())
			block.AddChildToBack(setterExprNode)
		}
	}
	var stmts pkg_rhino.Node = pkg_rhino.Node(0)
	_ = stmts
	if NodeUtil_IsFromTypeSummary(cls.input.GetAstRoot(this_.compiler)) {
		var typeSummaryInsertionRef pkg_rhino.Node = this_.compiler.GetSynthesizedTypeSummaryInput().GetAstRoot(this_.compiler)
		_ = typeSummaryInsertionRef
		block.SrcrefTreeIfMissing(typeSummaryInsertionRef)
		stmts = block.RemoveChildren()
		typeSummaryInsertionRef.AddChildrenToBack(stmts)
	} else {
		block.SrcrefTreeIfMissing(this_.externsInsertionRef)
		stmts = block.RemoveChildren()
		this_.externsInsertionRef.AddChildrenToBack(stmts)
	}
	this_.compiler.ReportChangeToEnclosingScope(stmts)
}

// Java method: PolymerClassRewriter.addMethodToObjectExternsUsingExportAnnotation(com.google.javascript.jscomp.PolymerClassDefinition cls, com.google.javascript.jscomp.PolymerPass.MemberDefinition method)
func (this_ *PolymerClassRewriter) priv_PolymerClassRewriter_addMethodToObjectExternsUsingExportAnnotation_4393(cls *PolymerClassDefinition, method *PolymerPass_MemberDefinition) {
	_ = cls
	_ = method
	var getprop pkg_rhino.Node = NodeUtil_NewQName_2_AbstractCompiler_string(this_.compiler, (j2gort.StrValString(cls.target.GetQualifiedName()) + ".prototype." + j2gort.StrValString(method.name.GetString())))
	_ = getprop
	var info *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_S_Builder().ParseDocumentation()
	_ = info
	if method.info != nil {
		info.RecordVisibility(method.info.GetVisibility())
	}
	info.RecordExport()
	getprop.SetJSDocInfo(info.Build_0())
	var expression pkg_rhino.Node = pkg_rhino.IR_ExprResult(getprop).SrcrefTreeIfMissing(method.name)
	_ = expression
	var insertAfter pkg_rhino.Node = cls.definition
	_ = insertAfter
	for !(NodeUtil_IsStatementBlock(insertAfter.GetParent())) {
		insertAfter = insertAfter.GetParent()
	}
	expression.InsertAfter(insertAfter)
	this_.compiler.ReportChangeToEnclosingScope(expression)
}

// Java static method: PolymerClassRewriter.varToAssign(com.google.javascript.rhino.Node var)
func PolymerClassRewriter_VarToAssign(var_ pkg_rhino.Node) pkg_rhino.Node {
	clinit_PolymerClassRewriter()
	_ = var_
	var assign pkg_rhino.Node = pkg_rhino.IR_Assign(var_.GetFirstChild().CloneNode_0(), var_.GetFirstChild().RemoveFirstChild())
	_ = assign
	return pkg_rhino.IR_ExprResult(assign).SrcrefTreeIfMissing(var_)
}

// Java method: PolymerClassRewriter.convertSimpleObserverStringsToReferences(com.google.javascript.jscomp.PolymerClassDefinition cls)
func (this_ *PolymerClassRewriter) priv_PolymerClassRewriter_convertSimpleObserverStringsToReferences_4394(cls *PolymerClassDefinition) {
	_ = cls
	{
		cur_37 := j2gort.NewIterCursor(j2gort.CheckNotNull(cls.props))
		for cur_37.HasNext() {
			rawElem_38 := cur_37.Next()
			prop := j2gort.CastPtr[PolymerPass_MemberDefinition](rawElem_38)
			_ = prop
			if prop.value.IsObjectLit() {
				var observer pkg_rhino.Node = NodeUtil_GetFirstPropMatchingKey(prop.value, "observer")
				_ = observer
				if (observer != 0) && observer.IsStringLit() {
					var observerDirectReference pkg_rhino.Node = pkg_rhino.IR_Getprop_3(cls.target.CloneTree_0(), "prototype", []string{observer.GetString()}).Srcref(observer)
					_ = observerDirectReference
					observer.ReplaceWith(observerDirectReference)
					this_.compiler.ReportChangeToEnclosingScope(observerDirectReference)
				}
			}
		}
	}
}

// Java method: PolymerClassRewriter.validateComputedPropertiesReflectionCalls(com.google.javascript.jscomp.PolymerClassDefinition cls)
func (this_ *PolymerClassRewriter) priv_PolymerClassRewriter_validateComputedPropertiesReflectionCalls_4395(cls *PolymerClassDefinition) {
	_ = cls
	{
		cur_39 := j2gort.NewIterCursor(j2gort.CheckNotNull(cls.props))
		for cur_39.HasNext() {
			rawElem_40 := cur_39.Next()
			prop := j2gort.CastPtr[PolymerPass_MemberDefinition](rawElem_40)
			_ = prop
			if prop.value.IsObjectLit() {
				var computed pkg_rhino.Node = NodeUtil_GetFirstPropMatchingKey(prop.value, "computed")
				_ = computed
				if (computed != 0) && computed.IsStringLit() {
					this_.priv_PolymerClassRewriter_validateReflectedString_4397(computed)
				}
			}
		}
	}
}

// Java method: PolymerClassRewriter.validateComplexObserverReflectionCalls(com.google.javascript.jscomp.PolymerClassDefinition cls)
func (this_ *PolymerClassRewriter) priv_PolymerClassRewriter_validateComplexObserverReflectionCalls_4396(cls *PolymerClassDefinition) {
	_ = cls
	var classMembers pkg_rhino.Node = NodeUtil_GetClassMembers(cls.definition)
	_ = classMembers
	var getter pkg_rhino.Node = NodeUtil_GetFirstGetterMatchingKey(classMembers, "observers")
	_ = getter
	if getter != 0 {
		var complexObservers pkg_rhino.Node = pkg_rhino.Node(0)
		_ = complexObservers
		{
			var child pkg_rhino.Node = NodeUtil_GetFunctionBody(getter.GetFirstChild()).GetFirstChild()
			_ = child
			for ; child != 0; child = child.GetNext() {
				if child.IsReturn() {
					if child.HasChildren() && child.GetFirstChild().IsArrayLit() {
						complexObservers = child.GetFirstChild()
						break
					}
				}
			}
		}
		if complexObservers != 0 {
			{
				var complexObserver pkg_rhino.Node = complexObservers.GetFirstChild()
				_ = complexObserver
				for complexObserver != 0 {
					var next pkg_rhino.Node = complexObserver.GetNext()
					_ = next
					if complexObserver.IsStringLit() {
						this_.priv_PolymerClassRewriter_validateReflectedString_4397(complexObserver)
					}
					complexObserver = next
				}
			}
		}
	}
}

// Java method: PolymerClassRewriter.validateReflectedString(com.google.javascript.rhino.Node methodSignature)
func (this_ *PolymerClassRewriter) priv_PolymerClassRewriter_validateReflectedString_4397(methodSignature pkg_rhino.Node) {
	_ = methodSignature
	if !(methodSignature.IsStringLit()) {
		j2gort.CheckArgument(false)
	}
	var methodSignatureString string = javalang.StringTrim(j2gort.CheckNotNullString(methodSignature.GetString()))
	_ = methodSignatureString
	var openParenIndex int32 = j2gort.JStrIndexOfCp(methodSignatureString, int32(uint16(40)))
	_ = openParenIndex
	if (int32(j2gort.JStrCharAt(methodSignatureString, (j2gort.JStrLen(methodSignatureString)-int32(1)))) != int32(uint16(41))) || (openParenIndex < int32(1)) {
		this_.compiler.Report_1(JSError_Make_3(methodSignature, StaticGet_PolymerPassErrors_POLYMER_UNPARSABLE_STRING(), []string{}))
		return
	}
	if openParenIndex < (j2gort.JStrLen(methodSignatureString) - int32(2)) {
		var methodParamsString string = javalang.StringTrim(j2gort.CheckNotNullString(j2gort.JStrSub(methodSignatureString, (openParenIndex + int32(1)), (j2gort.JStrLen(methodSignatureString) - int32(1)))))
		_ = methodParamsString
		this_.priv_PolymerClassRewriter_parseMethodParams_4398(methodParamsString, methodSignature)
	}
}

// Java method: PolymerClassRewriter.parseMethodParams(java.lang.String methodParameters, com.google.javascript.rhino.Node methodSignature)
func (this_ *PolymerClassRewriter) priv_PolymerClassRewriter_parseMethodParams_4398(methodParameters string, methodSignature pkg_rhino.Node) {
	_ = methodParameters
	_ = methodSignature
	var nextDelimeter uint16 = uint16(44)
	_ = nextDelimeter
	{
		var i int32 = int32(0)
		_ = i
		for ; i < j2gort.JStrLen(methodParameters); i = i + 1 {
			if int32(j2gort.JStrCharAt(methodParameters, i)) == int32(nextDelimeter) {
				nextDelimeter = uint16(44)
			} else {
				if (int32(j2gort.JStrCharAt(methodParameters, i)) == int32(uint16(34))) || (int32(j2gort.JStrCharAt(methodParameters, i)) == int32(uint16(39))) {
					nextDelimeter = j2gort.JStrCharAt(methodParameters, i)
				}
			}
		}
	}
	if int32(nextDelimeter) != int32(uint16(44)) {
		this_.compiler.Report_1(JSError_Make_3(methodSignature, StaticGet_PolymerPassErrors_POLYMER_UNPARSABLE_STRING(), []string{}))
	}
}

// Java static method: PolymerClassRewriter.getInsertionPointForGoogModule(com.google.javascript.rhino.Node moduleBody)
func PolymerClassRewriter_GetInsertionPointForGoogModule(moduleBody pkg_rhino.Node) pkg_rhino.Node {
	clinit_PolymerClassRewriter()
	_ = moduleBody
	var insertionPoint pkg_rhino.Node = moduleBody.GetFirstChild()
	_ = insertionPoint
	var next pkg_rhino.Node = insertionPoint.GetNext()
	_ = next
	for (PolymerClassRewriter_IsGoogRequireExpr(next) || NodeUtil_IsGoogModuleDeclareLegacyNamespaceCall(next)) || NodeUtil_IsGoogSetTestOnlyCall(next) {
		insertionPoint = next
		next = next.GetNext()
	}
	return insertionPoint
}

// Java static method: PolymerClassRewriter.isGoogRequireExpr(com.google.javascript.rhino.Node statement)
func PolymerClassRewriter_IsGoogRequireExpr(statement pkg_rhino.Node) bool {
	clinit_PolymerClassRewriter()
	_ = statement
	if NodeUtil_IsExprCall(statement) && ModuleImportResolver_IsGoogModuleDependencyCall(statement.GetOnlyChild()) {
		return true
	}
	if !(NodeUtil_IsNameDeclaration(statement)) {
		return false
	}
	var rhs pkg_rhino.Node
	_ = rhs
	if statement.GetFirstChild().IsName() {
		rhs = statement.GetFirstFirstChild()
	} else {
		rhs = statement.GetFirstChild().GetSecondChild()
	}
	return ModuleImportResolver_IsGoogModuleDependencyCall(rhs)
}
