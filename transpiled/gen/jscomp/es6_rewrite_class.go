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
// Port of com/google/javascript/jscomp/Es6RewriteClass.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"
	"sync/atomic"

	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_es6_rewrite_class_1 j2gort.Ext_com_google_common_base_Predicate

func init() {
	lambdaSingleton_es6_rewrite_class_1 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_1 uint32) bool {
		n := pkg_rhino.Node(rawNode_1)
		_ = n
		return !(n.IsFunction())
	})
}

// Java class: com.google.javascript.jscomp.Es6RewriteClass
type Es6RewriteClass struct {
	compiler                     AbstractCompiler                 `gson:"compiler"`
	astFactory                   *AstFactory                      `gson:"astFactory"`
	convertSuperConstructorCalls *Es6ConvertSuperConstructorCalls `gson:"convertSuperConstructorCalls"`
	transpilationNamespace       pkg_rhino.StaticScope            `gson:"transpilationNamespace"`
	jscompInherits               RuntimeJsLibManager_JsLibField   `gson:"jscompInherits"`
	idHash_                      atomic.Int32
}

func (this_ *Es6RewriteClass) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteClass) IsJavaNil_() bool { return this_ == nil }
func (*Es6RewriteClass) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteClass", "Es6RewriteClass"
}
func (this_ *Es6RewriteClass) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es6RewriteClass) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteClass", "com.google.javascript.jscomp.Es6RewriteClass", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteClass](v_)
	})
}
func (*Es6RewriteClass) Implements_jscomp_NodeTraversal_Callback() {}
func (*Es6RewriteClass) Implements_jscomp_CompilerPass()           {}

var static_Es6RewriteClass_features *pkg_parser.FeatureSet

func StaticGet_Es6RewriteClass_Features() *pkg_parser.FeatureSet {
	clinit_Es6RewriteClass()
	return static_Es6RewriteClass_features
}
func StaticSet_Es6RewriteClass_Features(v *pkg_parser.FeatureSet) {
	clinit_Es6RewriteClass()
	static_Es6RewriteClass_features = v
}

var static_Es6RewriteClass_CONSTRUCTOR_PRIMITIVE_RETURN *DiagnosticType

func StaticGet_Es6RewriteClass_CONSTRUCTOR_PRIMITIVE_RETURN() *DiagnosticType {
	clinit_Es6RewriteClass()
	return static_Es6RewriteClass_CONSTRUCTOR_PRIMITIVE_RETURN
}
func StaticSet_Es6RewriteClass_CONSTRUCTOR_PRIMITIVE_RETURN(v *DiagnosticType) {
	clinit_Es6RewriteClass()
	static_Es6RewriteClass_CONSTRUCTOR_PRIMITIVE_RETURN = v
}

var clinitDone_Es6RewriteClass uint32
var clinitState_Es6RewriteClass j2gort.ClinitState

func clinit_Es6RewriteClass() {
	if atomic.LoadUint32(&clinitDone_Es6RewriteClass) != 0 {
		return
	}
	clinitSlow_Es6RewriteClass()
}

func clinitSlow_Es6RewriteClass() {
	if !j2gort.EnterClinit(&clinitState_Es6RewriteClass, &clinitDone_Es6RewriteClass) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Es6RewriteClass, &clinitDone_Es6RewriteClass)
	static_Es6RewriteClass_features = pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature_7192([]*pkg_parser.FeatureSet_Feature{pkg_parser.StaticGet_FeatureSet_Feature_CLASSES(), pkg_parser.StaticGet_FeatureSet_Feature_CLASS_GETTER_SETTER(), pkg_parser.StaticGet_FeatureSet_Feature_NEW_TARGET(), pkg_parser.StaticGet_FeatureSet_Feature_SUPER()})
	static_Es6RewriteClass_CONSTRUCTOR_PRIMITIVE_RETURN = DiagnosticType_Error("JSC_CONSTRUCTOR_PRIMITIVE_RETURN", "Class constructors may only return an object or undefined")
}

func (this_ *Es6RewriteClass) initFields_Es6RewriteClass() {
	_ = this_
}

// Java constructor: Es6RewriteClass(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.CompilerOptions.Es6SubclassTranspilation es6SubclassTranspilation)
func NewEs6RewriteClass(compiler AbstractCompiler, es6SubclassTranspilation CompilerOptions_Es6SubclassTranspilation) *Es6RewriteClass {
	clinit_Es6RewriteClass()
	this_ := &Es6RewriteClass{}
	InitEs6RewriteClass(this_, compiler, es6SubclassTranspilation)
	return this_
}

func InitEs6RewriteClass(this_ *Es6RewriteClass, compiler AbstractCompiler, es6SubclassTranspilation CompilerOptions_Es6SubclassTranspilation) {
	_ = compiler
	_ = es6SubclassTranspilation
	this_.initFields_Es6RewriteClass()
	this_.compiler = compiler
	this_.astFactory = compiler.CreateAstFactory()
	this_.transpilationNamespace = compiler.GetTranspilationNamespace()
	this_.convertSuperConstructorCalls = NewEs6ConvertSuperConstructorCalls(compiler, es6SubclassTranspilation)
	this_.jscompInherits = compiler.GetRuntimeJsLibManager().GetJsLibField("$jscomp.inherits")
}

// Java method: Es6RewriteClass.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *Es6RewriteClass) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	TranspilationPasses_ProcessTranspile(this_.compiler, externs, StaticGet_Es6RewriteClass_Features(), []NodeTraversal_Callback{this_})
	TranspilationPasses_ProcessTranspile(this_.compiler, root, StaticGet_Es6RewriteClass_Features(), []NodeTraversal_Callback{this_})
	gn_ := NewGlobalNamespace_3_AbstractCompiler_Node_Node(this_.compiler, externs, root)
	this_.convertSuperConstructorCalls.SetGlobalNamespace(gn_)
	defer func() {
		this_.convertSuperConstructorCalls.SetGlobalNamespace(nil)
		gn_.ReleaseSlabs_()
	}()

	NodeTraversal_Traverse(this_.compiler, root, ToIface_NodeTraversal_Callback(this_.convertSuperConstructorCalls))
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(this_.compiler, root, StaticGet_Es6RewriteClass_Features())
}

// Java method: Es6RewriteClass.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteClass) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */ :
		if pkg_parser.StaticGet_FeatureSet_ES3().Contains_1_FeatureSet(this_.compiler.GetOptions().GetOutputFeatureSet()) {
			TranspilationUtil_CannotConvert(this_.compiler, n, "ES5 getters/setters (consider using --language_out=ES5)")
			return false
		}
	default:
	}
	return true
}

// Java method: Es6RewriteClass.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteClass) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsClass() {
		this_.priv_Es6RewriteClass_visitClass_2588(t, n, parent)
	}
}

// Java method: Es6RewriteClass.visitClass(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node classNode, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteClass) priv_Es6RewriteClass_visitClass_2588(t *NodeTraversal, classNode pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = classNode
	_ = parent
	var metadata *Es6RewriteClass_ClassDeclarationMetadata = Es6RewriteClass_ClassDeclarationMetadata_Create(classNode, parent, this_.astFactory)
	_ = metadata
	if metadata == nil {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Can only convert classes that are declarations or the right hand side of a simple assignment: "+j2gort.StrVal(classNode))))
	}
	if metadata.HasSuperClass() {
		j2gort.CheckStateFmt(metadata.SuperClassNameNode().IsQualifiedName(), "Expected Es6NormalizeClasses to make all extends clauses into qualified  names, found %s", j2gort.BoxNode(metadata.SuperClassNameNode()))
	}
	j2gort.CheckStateFmt(NodeUtil_IsStatement(metadata.InsertionPoint().GetNode()), "insertion point must be a statement: %s", j2gort.BoxNode(metadata.InsertionPoint().GetNode()))
	var constructor pkg_rhino.Node = pkg_rhino.Node(0)
	_ = constructor
	var classMembers pkg_rhino.Node = classNode.GetLastChild()
	_ = classMembers
	{
		var member pkg_rhino.Node = classMembers.GetFirstChild()
		_ = member
		for member != 0 {
			var next pkg_rhino.Node = member.GetNext()
			_ = next
			if (member.IsComputedProp() && (member.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_GETTER()) || member.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_SETTER()))) || (member.IsGetterDef() || member.IsSetterDef()) {
				this_.priv_Es6RewriteClass_visitNonMethodMember_2597(member, metadata)
			} else if NodeUtil_IsEs6ConstructorMemberFunctionDef(member) {
				if !(classNode.IsFromExterns()) {
					this_.priv_Es6RewriteClass_reportPrimitiveReturns_2589(member.GetFirstChild())
				}
				constructor = member.RemoveFirstChild().SetColor(classNode.GetColor())
				if !(metadata.Anonymous()) {
					constructor.GetFirstChild().ReplaceWith(metadata.ClassNameNode().CloneNode_0())
				}
			} else if member.IsEmpty() {
			} else {
				if !(member.IsMemberFunctionDef() || member.IsComputedProp()) {
					j2gort.CheckStateFmt(false, "Unexpected class member: (%s)", j2gort.BoxNode(member))
				}
				if !(!(member.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_VARIABLE()))) {
					j2gort.CheckStateFmt(false, "Member variables should have been transpiled earlier: (%s)", j2gort.BoxNode(member))
				}
				this_.priv_Es6RewriteClass_visitMethod_2598(member, metadata)
			}
			member = next
		}
	}
	j2gort.CheckNotNullNode(constructor)
	this_.priv_Es6RewriteClass_flushDefinePropertiesForPrototype_2592(metadata, classNode)
	this_.priv_Es6RewriteClass_flushDefinePropertiesForClass_2593(metadata, classNode)
	var classJSDoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(classNode)
	_ = classJSDoc
	var newInfo *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(classJSDoc)
	_ = newInfo
	newInfo.RecordConstructor()
	var enclosingStatement pkg_rhino.Node = NodeUtil_GetEnclosingStatement(classNode)
	_ = enclosingStatement
	if metadata.HasSuperClass() {
		if !(classNode.IsFromExterns()) {
			var inheritsCall pkg_rhino.Node = pkg_rhino.IR_ExprResult(this_.astFactory.CreateCall(this_.astFactory.CreateQName_2_StaticScope_RuntimeJsLibManager_JsLibField(this_.transpilationNamespace, this_.jscompInherits), AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_NULL_OR_VOID()), []pkg_rhino.Node{metadata.FullClassNameNode().CloneTree_0(), metadata.SuperClassNameNode().CloneTree_0()})).SrcrefTreeIfMissing(metadata.SuperClassNameNode())
			_ = inheritsCall
			inheritsCall.InsertAfter(enclosingStatement)
		}
	}
	this_.priv_Es6RewriteClass_addTypeDeclarations_2599(metadata, enclosingStatement)
	if NodeUtil_IsStatement(classNode) {
		constructor.GetFirstChild().SetString("")
		var ctorVar pkg_rhino.Node = pkg_rhino.IR_Let_2(metadata.ClassNameNode().CloneNode_0(), constructor)
		_ = ctorVar
		ctorVar.SrcrefTreeIfMissing(classNode)
		classNode.ReplaceWith(ctorVar)
		NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_LET_DECLARATIONS(), this_.compiler)
	} else {
		classNode.ReplaceWith(constructor)
	}
	NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(classNode, this_.compiler)
	if NodeUtil_IsStatement(constructor) {
		constructor.SetJSDocInfo(newInfo.Build_0())
	} else if parent.IsName() {
		var var_ pkg_rhino.Node = parent.GetParent()
		_ = var_
		var_.SetJSDocInfo(newInfo.Build_0())
	} else if constructor.GetParent().IsName() {
		var var__2 pkg_rhino.Node = constructor.GetGrandparent()
		_ = var__2
		var__2.SetJSDocInfo(newInfo.Build_0())
	} else if parent.IsAssign() {
		parent.SetJSDocInfo(newInfo.Build_0())
	} else {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected parent node "+j2gort.StrVal(parent))))
	}
	t.ReportCodeChange_0()
}

// Java method: Es6RewriteClass.reportPrimitiveReturns(com.google.javascript.rhino.Node constructor)
func (this_ *Es6RewriteClass) priv_Es6RewriteClass_reportPrimitiveReturns_2589(constructor pkg_rhino.Node) {
	_ = constructor
	NodeUtil_VisitPostOrder_3(NodeUtil_GetFunctionBody(constructor), NewLambda_NodeUtil_Visitor(func(n pkg_rhino.Node) {
		_ = n
		if ((n.IsReturn() && n.HasChildren()) && NodeUtil_IsImmutableValue(n.GetFirstChild())) && (NodeUtil_GetKnownValueType(n.GetFirstChild()) != StaticGet_NodeUtil_ValueType_VOID()) {
			this_.compiler.Report_1(JSError_Make_3(n, StaticGet_Es6RewriteClass_CONSTRUCTOR_PRIMITIVE_RETURN(), []string{}))
		}
	}), lambdaSingleton_es6_rewrite_class_1)
}

// Java method: Es6RewriteClass.createObjectDotDefineProperties()
func (this_ *Es6RewriteClass) priv_Es6RewriteClass_createObjectDotDefineProperties_2590() pkg_rhino.Node {
	return this_.astFactory.CreateJSCompDotGlobalAccess(this_.transpilationNamespace, "Object.defineProperties")
}

// Java method: Es6RewriteClass.createObjectDotDefineProperty()
func (this_ *Es6RewriteClass) priv_Es6RewriteClass_createObjectDotDefineProperty_2591() pkg_rhino.Node {
	return this_.astFactory.CreateJSCompDotGlobalAccess(this_.transpilationNamespace, "Object.defineProperty")
}

// Java method: Es6RewriteClass.flushDefinePropertiesForPrototype(com.google.javascript.jscomp.Es6RewriteClass.ClassDeclarationMetadata metadata, com.google.javascript.rhino.Node srcrefNode)
func (this_ *Es6RewriteClass) priv_Es6RewriteClass_flushDefinePropertiesForPrototype_2592(metadata *Es6RewriteClass_ClassDeclarationMetadata, srcrefNode pkg_rhino.Node) {
	_ = metadata
	_ = srcrefNode
	if metadata.DefinePropertiesObjForPrototype().HasChildren() {
		var definePropsCall pkg_rhino.Node = pkg_rhino.IR_ExprResult(this_.astFactory.CreateCall(this_.priv_Es6RewriteClass_createObjectDotDefineProperties_2590(), AstFactory_S_Type_1_Node(metadata.ClassPrototypeNode()), []pkg_rhino.Node{metadata.ClassPrototypeNode().CloneTree_0(), metadata.DefinePropertiesObjForPrototype()}))
		_ = definePropsCall
		definePropsCall.SrcrefTreeIfMissing(srcrefNode)
		metadata.InsertNodeAndAdvance(definePropsCall)
		metadata.ResetDefinePropertiesObjForPrototype(this_.astFactory)
	}
}

// Java method: Es6RewriteClass.flushDefinePropertiesForClass(com.google.javascript.jscomp.Es6RewriteClass.ClassDeclarationMetadata metadata, com.google.javascript.rhino.Node srcrefNode)
func (this_ *Es6RewriteClass) priv_Es6RewriteClass_flushDefinePropertiesForClass_2593(metadata *Es6RewriteClass_ClassDeclarationMetadata, srcrefNode pkg_rhino.Node) {
	_ = metadata
	_ = srcrefNode
	if metadata.DefinePropertiesObjForClass().HasChildren() {
		var definePropsCall pkg_rhino.Node = pkg_rhino.IR_ExprResult(this_.astFactory.CreateCall(this_.priv_Es6RewriteClass_createObjectDotDefineProperties_2590(), AstFactory_S_Type_1_Node(metadata.FullClassNameNode()), []pkg_rhino.Node{metadata.FullClassNameNode().CloneTree_0(), metadata.DefinePropertiesObjForClass()}))
		_ = definePropsCall
		definePropsCall.SrcrefTreeIfMissing(srcrefNode)
		metadata.InsertNodeAndAdvance(definePropsCall)
		metadata.ResetDefinePropertiesObjForClass(this_.astFactory)
	}
}

// Java method: Es6RewriteClass.flushPendingDefineProperties(com.google.javascript.jscomp.Es6RewriteClass.ClassDeclarationMetadata metadata, boolean isStatic, com.google.javascript.rhino.Node srcrefNode)
func (this_ *Es6RewriteClass) priv_Es6RewriteClass_flushPendingDefineProperties_2594(metadata *Es6RewriteClass_ClassDeclarationMetadata, isStatic bool, srcrefNode pkg_rhino.Node) {
	_ = metadata
	_ = isStatic
	_ = srcrefNode
	if isStatic {
		this_.priv_Es6RewriteClass_flushDefinePropertiesForClass_2593(metadata, srcrefNode)
	} else {
		this_.priv_Es6RewriteClass_flushDefinePropertiesForPrototype_2592(metadata, srcrefNode)
	}
}

// Java method: Es6RewriteClass.addToDefinePropertiesObject(com.google.javascript.jscomp.Es6RewriteClass.ClassDeclarationMetadata metadata, com.google.javascript.rhino.Node member)
func (this_ *Es6RewriteClass) priv_Es6RewriteClass_addToDefinePropertiesObject_2595(metadata *Es6RewriteClass_ClassDeclarationMetadata, member pkg_rhino.Node) {
	_ = metadata
	_ = member
	if !(!(member.IsComputedProp())) {
		j2gort.CheckArgument(false)
	}
	var obj pkg_rhino.Node
	_ = obj
	if member.IsStaticMember() {
		obj = metadata.DefinePropertiesObjForClass()
	} else {
		obj = metadata.DefinePropertiesObjForPrototype()
	}
	var prop pkg_rhino.Node = NodeUtil_GetFirstPropMatchingKey(obj, member.GetString())
	_ = prop
	var accessorKind string
	_ = accessorKind
	if member.IsGetterDef() {
		accessorKind = "get"
	} else {
		accessorKind = "set"
	}
	if (prop != 0) && (NodeUtil_GetFirstPropMatchingKey(prop, accessorKind) != 0) {
		this_.priv_Es6RewriteClass_flushPendingDefineProperties_2594(metadata, member.IsStaticMember(), member)
		if member.IsStaticMember() {
			obj = metadata.DefinePropertiesObjForClass()
		} else {
			obj = metadata.DefinePropertiesObjForPrototype()
		}
		prop = pkg_rhino.Node(0)
	}
	if prop == 0 {
		prop = this_.priv_Es6RewriteClass_createPropertyDescriptor_2601()
		if j2gort.CheckNotNullString(member.GetString()) == "__proto__" {
			var stringKey pkg_rhino.Node = this_.astFactory.CreateComputedProperty(this_.astFactory.CreateString("__proto__"), prop)
			_ = stringKey
			obj.AddChildToBack(stringKey)
			var scriptNode pkg_rhino.Node = NodeUtil_GetEnclosingScript(member)
			_ = scriptNode
			if scriptNode != 0 {
				NodeUtil_AddFeatureToScript(scriptNode, pkg_parser.StaticGet_FeatureSet_Feature_COMPUTED_PROPERTIES(), this_.compiler)
			}
		} else {
			var stringKey_2 pkg_rhino.Node = this_.astFactory.CreateStringKey(member.GetString(), prop)
			_ = stringKey_2
			if member.IsQuotedStringKey() {
				stringKey_2.PutBooleanProp(pkg_rhino.StaticGet_Node_QUOTED_PROP(), true)
			}
			obj.AddChildToBack(stringKey_2)
		}
	}
	var function pkg_rhino.Node = member.GetLastChild()
	_ = function
	var info *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(function)
	_ = info
	var stringKey_3 pkg_rhino.Node = this_.astFactory.CreateStringKey(accessorKind, function.Detach())
	_ = stringKey_3
	stringKey_3.SetJSDocInfo(info)
	prop.AddChildToBack(stringKey_3)
	prop.SrcrefTreeIfMissing(member)
}

// Java method: Es6RewriteClass.extractComputedProperty(com.google.javascript.rhino.Node computedMember, com.google.javascript.jscomp.Es6RewriteClass.ClassDeclarationMetadata metadata)
func (this_ *Es6RewriteClass) priv_Es6RewriteClass_extractComputedProperty_2596(computedMember pkg_rhino.Node, metadata *Es6RewriteClass_ClassDeclarationMetadata) {
	_ = computedMember
	_ = metadata
	this_.priv_Es6RewriteClass_flushPendingDefineProperties_2594(metadata, computedMember.IsStaticMember(), computedMember)
	var owner pkg_rhino.Node
	_ = owner
	if computedMember.IsStaticMember() {
		owner = metadata.FullClassNameNode()
	} else {
		owner = metadata.ClassPrototypeNode()
	}
	var property pkg_rhino.Node = computedMember.RemoveFirstChild()
	_ = property
	var propertyValue pkg_rhino.Node = computedMember.RemoveFirstChild()
	_ = propertyValue
	var propertyDescriptor pkg_rhino.Node = this_.priv_Es6RewriteClass_createPropertyDescriptor_2601()
	_ = propertyDescriptor
	var stringKey pkg_rhino.Node = this_.astFactory.CreateStringKey(func() string {
		if computedMember.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_GETTER()) {
			return "get"
		}
		return "set"
	}(), propertyValue)
	_ = stringKey
	propertyDescriptor.AddChildToBack(stringKey)
	var objectDefinePropertyCall pkg_rhino.Node = this_.astFactory.CreateCall(this_.priv_Es6RewriteClass_createObjectDotDefineProperty_2591(), AstFactory_S_Type_1_Node(owner), []pkg_rhino.Node{owner.CloneTree_0(), property, propertyDescriptor})
	_ = objectDefinePropertyCall
	metadata.InsertNodeAndAdvance(pkg_rhino.IR_ExprResult(objectDefinePropertyCall).SrcrefTreeIfMissing(computedMember))
}

// Java method: Es6RewriteClass.visitNonMethodMember(com.google.javascript.rhino.Node member, com.google.javascript.jscomp.Es6RewriteClass.ClassDeclarationMetadata metadata)
func (this_ *Es6RewriteClass) priv_Es6RewriteClass_visitNonMethodMember_2597(member pkg_rhino.Node, metadata *Es6RewriteClass_ClassDeclarationMetadata) {
	_ = member
	_ = metadata
	if member.IsComputedProp() {
		this_.priv_Es6RewriteClass_extractComputedProperty_2596(member.Detach(), metadata)
		return
	}
	this_.priv_Es6RewriteClass_addToDefinePropertiesObject_2595(metadata, member)
	if !(member.IsStaticMember()) {
		return
	}
	if !(!(member.IsComputedProp())) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(member))
	}
	var membersToDeclare j2gort.Ext_java_util_Map = metadata.ClassMembersToDeclare()
	_ = membersToDeclare
	var builder Es6RewriteClass_ClassProperty_Builder = Es6RewriteClass_ClassProperty_S_Builder()
	_ = builder
	var memberName string = member.GetString()
	_ = memberName
	if member.IsQuotedStringKey() {
		builder.Kind(StaticGet_Es6RewriteClass_ClassProperty_PropertyKind_QUOTED_PROPERTY())
	} else {
		builder.Kind(StaticGet_Es6RewriteClass_ClassProperty_PropertyKind_NORMAL_PROPERTY())
	}
	builder.PropertyKey(memberName).PropertyType(member.GetColor())
	var jsDoc *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_S_Builder()
	_ = jsDoc
	jsDoc.RecordNoCollapse()
	builder.JsDocInfo(jsDoc.Build_0())
	j2gort.MapPutStr(membersToDeclare, memberName, j2gort.PtrToAny(builder.Build()))
}

// Java method: Es6RewriteClass.visitMethod(com.google.javascript.rhino.Node member, com.google.javascript.jscomp.Es6RewriteClass.ClassDeclarationMetadata metadata)
func (this_ *Es6RewriteClass) priv_Es6RewriteClass_visitMethod_2598(member pkg_rhino.Node, metadata *Es6RewriteClass_ClassDeclarationMetadata) {
	_ = member
	_ = metadata
	this_.priv_Es6RewriteClass_flushPendingDefineProperties_2594(metadata, member.IsStaticMember(), member)
	var qualifiedMemberAccess pkg_rhino.Node = this_.priv_Es6RewriteClass_getQualifiedMemberAccess_2600(member, metadata)
	_ = qualifiedMemberAccess
	var method pkg_rhino.Node = member.GetLastChild().Detach()
	_ = method
	var assign pkg_rhino.Node = this_.astFactory.CreateAssign_2_Node_Node(qualifiedMemberAccess, method).SrcrefIfMissing(method)
	_ = assign
	var info *pkg_rhino.JSDocInfo = member.GetJSDocInfo()
	_ = info
	if member.IsStaticMember() && NodeUtil_ReferencesOwnReceiver(assign.GetLastChild()) {
		var memberDoc *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(info)
		_ = memberDoc
		memberDoc.RecordThisType(pkg_rhino.NewJSTypeExpression(pkg_rhino.NewNode_2(pkg_rhino.StaticGet_Token_BANG(), pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_QMARK())).SrcrefTree(member), member.GetSourceFileName()))
		info = memberDoc.Build_0()
	}
	if info != nil {
		assign.SetJSDocInfo(info)
	}
	var newNode pkg_rhino.Node = NodeUtil_NewExpr(assign)
	_ = newNode
	metadata.InsertNodeAndAdvance(newNode)
}

// Java method: Es6RewriteClass.addTypeDeclarations(com.google.javascript.jscomp.Es6RewriteClass.ClassDeclarationMetadata metadata, com.google.javascript.rhino.Node insertionPoint)
func (this_ *Es6RewriteClass) priv_Es6RewriteClass_addTypeDeclarations_2599(metadata *Es6RewriteClass_ClassDeclarationMetadata, insertionPoint pkg_rhino.Node) {
	_ = metadata
	_ = insertionPoint
	{
		cur_2 := j2gort.NewMapValCursor(metadata.ClassMembersToDeclare())
		for cur_2.HasNext() {
			rawElem_3 := cur_2.Next()
			property := j2gort.CastPtr[Es6RewriteClass_ClassProperty](rawElem_3)
			_ = property
			var declaration pkg_rhino.Node = property.GetDeclaration(this_.astFactory, metadata.FullClassNameNode().CloneTree_0())
			_ = declaration
			declaration.SrcrefTreeIfMissing(metadata.ClassNameNode())
			declaration.InsertAfter(insertionPoint)
			insertionPoint = declaration
		}
	}
}

// Java method: Es6RewriteClass.getQualifiedMemberAccess(com.google.javascript.rhino.Node member, com.google.javascript.jscomp.Es6RewriteClass.ClassDeclarationMetadata metadata)
func (this_ *Es6RewriteClass) priv_Es6RewriteClass_getQualifiedMemberAccess_2600(member pkg_rhino.Node, metadata *Es6RewriteClass_ClassDeclarationMetadata) pkg_rhino.Node {
	_ = member
	_ = metadata
	var context pkg_rhino.Node
	_ = context
	if member.IsStaticMember() {
		context = metadata.FullClassNameNode().CloneTree_0()
	} else {
		context = metadata.ClassPrototypeNode().CloneTree_0()
	}
	context.MakeNonIndexableRecursive()
	if member.IsComputedProp() {
		return this_.astFactory.CreateGetElem(context, member.RemoveFirstChild()).SrcrefTreeIfMissing(member)
	} else {
		var methodName pkg_rhino.Node = member.GetFirstFirstChild()
		_ = methodName
		return this_.astFactory.CreateGetProp(context, member.GetString(), AstFactory_S_Type_1_Node(member)).SrcrefTree(methodName)
	}
}

// Java method: Es6RewriteClass.createPropertyDescriptor()
func (this_ *Es6RewriteClass) priv_Es6RewriteClass_createPropertyDescriptor_2601() pkg_rhino.Node {
	return this_.astFactory.CreateObjectLit_1([]pkg_rhino.Node{this_.astFactory.CreateStringKey("configurable", this_.astFactory.CreateBoolean(true)), this_.astFactory.CreateStringKey("enumerable", this_.astFactory.CreateBoolean(true))})
}

// Java record: com.google.javascript.jscomp.Es6RewriteClass.ClassProperty
type Es6RewriteClass_ClassProperty struct {
	propertyKey  string                                     `gson:"propertyKey"`
	jsDocInfo    *pkg_rhino.JSDocInfo                       `gson:"jsDocInfo"`
	propertyType *pkg_colors.Color                          `gson:"propertyType"`
	kind         Es6RewriteClass_ClassProperty_PropertyKind `gson:"kind"`
}

func (this_ *Es6RewriteClass_ClassProperty) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteClass_ClassProperty) IsJavaNil_() bool { return this_ == nil }
func (*Es6RewriteClass_ClassProperty) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteClass$ClassProperty", "ClassProperty"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteClass.ClassProperty", "com.google.javascript.jscomp.Es6RewriteClass$ClassProperty", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteClass_ClassProperty](v_)
	})
}
func (*Es6RewriteClass_ClassProperty) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *Es6RewriteClass_ClassProperty) initFields_Es6RewriteClass_ClassProperty() {
	_ = this_
}

// Java constructor: ClassProperty(java.lang.String propertyKey, com.google.javascript.jscomp.Es6RewriteClass.ClassProperty.PropertyKind kind, com.google.javascript.rhino.JSDocInfo jsDocInfo, com.google.javascript.jscomp.colors.@org.jspecify.annotations.Nullable Color propertyType)
func NewEs6RewriteClass_ClassProperty(propertyKey string, kind Es6RewriteClass_ClassProperty_PropertyKind, jsDocInfo *pkg_rhino.JSDocInfo, propertyType *pkg_colors.Color) *Es6RewriteClass_ClassProperty {
	this_ := &Es6RewriteClass_ClassProperty{}
	InitEs6RewriteClass_ClassProperty(this_, propertyKey, kind, jsDocInfo, propertyType)
	return this_
}

func InitEs6RewriteClass_ClassProperty(this_ *Es6RewriteClass_ClassProperty, propertyKey string, kind Es6RewriteClass_ClassProperty_PropertyKind, jsDocInfo *pkg_rhino.JSDocInfo, propertyType *pkg_colors.Color) {
	_ = propertyKey
	_ = kind
	_ = jsDocInfo
	_ = propertyType
	this_.initFields_Es6RewriteClass_ClassProperty()
	j2gort.CheckNotNullString(propertyKey)
	j2gort.CheckNotNullEnum(kind)
	j2gort.CheckNotNullPtr(jsDocInfo)
	this_.propertyKey = propertyKey
	this_.kind = kind
	this_.jsDocInfo = jsDocInfo
	this_.propertyType = propertyType
}

// Java method: ClassProperty.getDeclaration(com.google.javascript.jscomp.AstFactory astFactory, com.google.javascript.rhino.Node toDeclareOn)
func (this_ *Es6RewriteClass_ClassProperty) GetDeclaration(astFactory *AstFactory, toDeclareOn pkg_rhino.Node) pkg_rhino.Node {
	_ = astFactory
	_ = toDeclareOn
	var decl pkg_rhino.Node = pkg_rhino.Node(0)
	_ = decl
	switch (this_.Kind()).Ordinal() {
	case int32(0) /* PropertyKind.QUOTED_PROPERTY */ :
		decl = astFactory.CreateGetElem(toDeclareOn, astFactory.CreateString(this_.PropertyKey()))
	case int32(1) /* PropertyKind.COMPUTED_PROPERTY */ :
		panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", j2gort.BoxString(this_.String()))))
	case int32(2) /* PropertyKind.NORMAL_PROPERTY */ :
		decl = astFactory.CreateGetProp(toDeclareOn, this_.PropertyKey(), AstFactory_S_Type_1_Color(this_.PropertyType()))
	}
	decl.SetJSDocInfo(this_.JsDocInfo())
	decl = astFactory.ExprResult(decl)
	return decl
}

// Java static method: ClassProperty.builder()
func Es6RewriteClass_ClassProperty_S_Builder() Es6RewriteClass_ClassProperty_Builder {
	return NewAutoBuilder_Es6RewriteClass_ClassProperty_Builder()
}

func (this_ *Es6RewriteClass_ClassProperty) PropertyKey() string { return this_.propertyKey }
func (this_ *Es6RewriteClass_ClassProperty) Kind() Es6RewriteClass_ClassProperty_PropertyKind {
	return this_.kind
}
func (this_ *Es6RewriteClass_ClassProperty) JsDocInfo() *pkg_rhino.JSDocInfo { return this_.jsDocInfo }
func (this_ *Es6RewriteClass_ClassProperty) PropertyType() *pkg_colors.Color {
	return this_.propertyType
}
func (this_ *Es6RewriteClass_ClassProperty) Equals(other_ any) bool {
	o_, ok_ := other_.(*Es6RewriteClass_ClassProperty)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.propertyKey != o_.propertyKey {
		return false
	}
	if this_.kind != o_.kind {
		return false
	}
	if this_.jsDocInfo != o_.jsDocInfo {
		return false
	}
	if !javalang.ObjectsEquals(j2gort.PtrToAny(this_.propertyType), j2gort.PtrToAny(o_.propertyType)) {
		return false
	}
	return true
}
func (this_ *Es6RewriteClass_ClassProperty) HashCode() int32 {
	var h_ int32 = 0
	{
		hs_4 := this_.propertyKey
		h_ = 31 * h_
		if hs_4 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_4)
		}
	}
	h_ = 31*h_ + (this_.kind).HashCode()
	{
		hp_5 := this_.jsDocInfo
		h_ = 31 * h_
		if hp_5 != nil {
			h_ += j2gort.IdentityHash(hp_5)
		}
	}
	{
		hp_6 := this_.propertyType
		h_ = 31 * h_
		if hp_6 != nil {
			h_ += hp_6.HashCode()
		}
	}
	return h_
}
func (this_ *Es6RewriteClass_ClassProperty) String() string {
	_ = fmt.Sprint
	return "ClassProperty[" + "propertyKey=" + j2gort.StrVal(j2gort.BoxString(this_.propertyKey)) + ", " + "kind=" + j2gort.StrVal(j2gort.BoxNullableVal(this_.kind)) + ", " + "jsDocInfo=" + j2gort.StrVal(j2gort.PtrToAny(this_.jsDocInfo)) + ", " + "propertyType=" + j2gort.StrVal(j2gort.PtrToAny(this_.propertyType)) + "]"
}

// Java enum: com.google.javascript.jscomp.Es6RewriteClass.ClassProperty.PropertyKind
type Es6RewriteClass_ClassProperty_PropertyKind uint8

const (
	enum_Es6RewriteClass_ClassProperty_PropertyKind_QUOTED_PROPERTY   Es6RewriteClass_ClassProperty_PropertyKind = 1
	enum_Es6RewriteClass_ClassProperty_PropertyKind_COMPUTED_PROPERTY Es6RewriteClass_ClassProperty_PropertyKind = 2
	enum_Es6RewriteClass_ClassProperty_PropertyKind_NORMAL_PROPERTY   Es6RewriteClass_ClassProperty_PropertyKind = 3
)

var enumBaseTable_Es6RewriteClass_ClassProperty_PropertyKind = [...]j2gort.EnumBase{
	{Name_: "QUOTED_PROPERTY", Ordinal_: 0},
	{Name_: "COMPUTED_PROPERTY", Ordinal_: 1},
	{Name_: "NORMAL_PROPERTY", Ordinal_: 2},
}
var enumValues_Es6RewriteClass_ClassProperty_PropertyKind = []Es6RewriteClass_ClassProperty_PropertyKind{1, 2, 3}

func StaticGet_Es6RewriteClass_ClassProperty_PropertyKind_QUOTED_PROPERTY() Es6RewriteClass_ClassProperty_PropertyKind {
	return enum_Es6RewriteClass_ClassProperty_PropertyKind_QUOTED_PROPERTY
}
func StaticGet_Es6RewriteClass_ClassProperty_PropertyKind_COMPUTED_PROPERTY() Es6RewriteClass_ClassProperty_PropertyKind {
	return enum_Es6RewriteClass_ClassProperty_PropertyKind_COMPUTED_PROPERTY
}
func StaticGet_Es6RewriteClass_ClassProperty_PropertyKind_NORMAL_PROPERTY() Es6RewriteClass_ClassProperty_PropertyKind {
	return enum_Es6RewriteClass_ClassProperty_PropertyKind_NORMAL_PROPERTY
}

func Es6RewriteClass_ClassProperty_PropertyKind_Values() []Es6RewriteClass_ClassProperty_PropertyKind {
	return j2gort.CloneSlice(enumValues_Es6RewriteClass_ClassProperty_PropertyKind)
}

func Es6RewriteClass_ClassProperty_PropertyKind_ValueOf(name string) Es6RewriteClass_ClassProperty_PropertyKind {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_Es6RewriteClass_ClassProperty_PropertyKind {
		if enumBaseTable_Es6RewriteClass_ClassProperty_PropertyKind[i_].Name_ == name {
			return Es6RewriteClass_ClassProperty_PropertyKind(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.Es6RewriteClass.ClassProperty.PropertyKind." + name))
}

func (this_ Es6RewriteClass_ClassProperty_PropertyKind) IsJavaNil_() bool                { return this_ == 0 }
func (Es6RewriteClass_ClassProperty_PropertyKind) HasIdentityEquals_()                   {}
func (Es6RewriteClass_ClassProperty_PropertyKind) Implements_j2gort_Ext_java_lang_Enum() {}
func (Es6RewriteClass_ClassProperty_PropertyKind) Implements_j2gort_Ext_java_lang_constant_Constable() {
}
func (this_ Es6RewriteClass_ClassProperty_PropertyKind) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (Es6RewriteClass_ClassProperty_PropertyKind) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteClass$ClassProperty$PropertyKind", "PropertyKind"
}
func (this_ Es6RewriteClass_ClassProperty_PropertyKind) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_Es6RewriteClass_ClassProperty_PropertyKind[this_-1]
}
func (this_ Es6RewriteClass_ClassProperty_PropertyKind) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ Es6RewriteClass_ClassProperty_PropertyKind) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ Es6RewriteClass_ClassProperty_PropertyKind) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Es6RewriteClass_ClassProperty_PropertyKind[this_-1].Name_
}
func (this_ Es6RewriteClass_ClassProperty_PropertyKind) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_Es6RewriteClass_ClassProperty_PropertyKind[this_-1].Name_
}
func (this_ Es6RewriteClass_ClassProperty_PropertyKind) Equals(other_ any) bool {
	o_, ok_ := other_.(Es6RewriteClass_ClassProperty_PropertyKind)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ Es6RewriteClass_ClassProperty_PropertyKind) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_Es6RewriteClass_ClassProperty_PropertyKind[this_-1].HashCode()
}
func (this_ Es6RewriteClass_ClassProperty_PropertyKind) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Es6RewriteClass_ClassProperty_PropertyKind[this_-1].CompareTo(other_)
}
func (this_ Es6RewriteClass_ClassProperty_PropertyKind) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_Es6RewriteClass_ClassProperty_PropertyKind[this_-1].Name_)
}
func init() {
	vals_ := []any{Es6RewriteClass_ClassProperty_PropertyKind(1), Es6RewriteClass_ClassProperty_PropertyKind(2), Es6RewriteClass_ClassProperty_PropertyKind(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.Es6RewriteClass.ClassProperty.PropertyKind", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.Es6RewriteClass$ClassProperty$PropertyKind", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteClass.ClassProperty.PropertyKind", "com.google.javascript.jscomp.Es6RewriteClass$ClassProperty$PropertyKind", func(v_ any) bool {
		return j2gort.InstanceOf[Es6RewriteClass_ClassProperty_PropertyKind](v_)
	})
}

// Java polymorphic class interface: com.google.javascript.jscomp.Es6RewriteClass.ClassProperty.Builder
type Es6RewriteClass_ClassProperty_Builder interface {
	Implements_jscomp_Es6RewriteClass_ClassProperty_Builder()
	IsJavaNil_() bool
	GetEs6RewriteClass_ClassProperty_BuilderBase() *Es6RewriteClass_ClassProperty_BuilderBase
	PropertyKey(string) Es6RewriteClass_ClassProperty_Builder
	Kind(Es6RewriteClass_ClassProperty_PropertyKind) Es6RewriteClass_ClassProperty_Builder
	JsDocInfo(*pkg_rhino.JSDocInfo) Es6RewriteClass_ClassProperty_Builder
	PropertyType(*pkg_colors.Color) Es6RewriteClass_ClassProperty_Builder
	Build() *Es6RewriteClass_ClassProperty
}

func ToIface_Es6RewriteClass_ClassProperty_Builder[T interface {
	Es6RewriteClass_ClassProperty_Builder
	comparable
}](v T) Es6RewriteClass_ClassProperty_Builder {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_Es6RewriteClass_ClassProperty_Builder[T Es6RewriteClass_ClassProperty_Builder](v T) Es6RewriteClass_ClassProperty_Builder {
	var vi Es6RewriteClass_ClassProperty_Builder = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.Es6RewriteClass.ClassProperty.Builder
type Es6RewriteClass_ClassProperty_BuilderBase struct {
	idHash_ atomic.Int32
}

func (this_ *Es6RewriteClass_ClassProperty_BuilderBase) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteClass_ClassProperty_BuilderBase) IsJavaNil_() bool { return this_ == nil }
func (*Es6RewriteClass_ClassProperty_BuilderBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteClass$ClassProperty$Builder", "Builder"
}
func (this_ *Es6RewriteClass_ClassProperty_BuilderBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteClass.ClassProperty.Builder", "com.google.javascript.jscomp.Es6RewriteClass$ClassProperty$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[Es6RewriteClass_ClassProperty_Builder](v_)
	})
}
func (b *Es6RewriteClass_ClassProperty_BuilderBase) GetEs6RewriteClass_ClassProperty_BuilderBase() *Es6RewriteClass_ClassProperty_BuilderBase {
	return b
}
func (*Es6RewriteClass_ClassProperty_BuilderBase) Implements_jscomp_Es6RewriteClass_ClassProperty_Builder() {
}

func (this_ *Es6RewriteClass_ClassProperty_BuilderBase) initFields_Es6RewriteClass_ClassProperty_Builder() {
	_ = this_
}

func InitEs6RewriteClass_ClassProperty_Builder(this_ *Es6RewriteClass_ClassProperty_BuilderBase) {
	this_.initFields_Es6RewriteClass_ClassProperty_Builder()
}

// Java method: Builder.propertyKey(java.lang.String value)
func (this_ *Es6RewriteClass_ClassProperty_BuilderBase) PropertyKey(value string) Es6RewriteClass_ClassProperty_Builder {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Builder.kind(com.google.javascript.jscomp.Es6RewriteClass.ClassProperty.PropertyKind value)
func (this_ *Es6RewriteClass_ClassProperty_BuilderBase) Kind(value Es6RewriteClass_ClassProperty_PropertyKind) Es6RewriteClass_ClassProperty_Builder {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Builder.jsDocInfo(com.google.javascript.rhino.JSDocInfo value)
func (this_ *Es6RewriteClass_ClassProperty_BuilderBase) JsDocInfo(value *pkg_rhino.JSDocInfo) Es6RewriteClass_ClassProperty_Builder {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Builder.propertyType(com.google.javascript.jscomp.colors.@org.jspecify.annotations.Nullable Color type)
func (this_ *Es6RewriteClass_ClassProperty_BuilderBase) PropertyType(type_ *pkg_colors.Color) Es6RewriteClass_ClassProperty_Builder {
	_ = type_
	panic("abstract or native method invoked")
}

// Java method: Builder.build()
func (this_ *Es6RewriteClass_ClassProperty_BuilderBase) Build() *Es6RewriteClass_ClassProperty {
	panic("abstract or native method invoked")
}

// Java class: com.google.javascript.jscomp.Es6RewriteClass.ClassDeclarationMetadata
type Es6RewriteClass_ClassDeclarationMetadata struct {
	insertionPoint                  *Es6RewriteClass_InsertionPoint `gson:"insertionPoint"`
	classMembersToDeclare           j2gort.Ext_java_util_Map        `gson:"classMembersToDeclare"`
	definePropertiesObjForPrototype pkg_rhino.Node                  `gson:"definePropertiesObjForPrototype"`
	definePropertiesObjForClass     pkg_rhino.Node                  `gson:"definePropertiesObjForClass"`
	fullClassNameNode               pkg_rhino.Node                  `gson:"fullClassNameNode"`
	classPrototypeNode              pkg_rhino.Node                  `gson:"classPrototypeNode"`
	classNameNode                   pkg_rhino.Node                  `gson:"classNameNode"`
	superClassNameNode              pkg_rhino.Node                  `gson:"superClassNameNode"`
	idHash_                         atomic.Int32
	anonymous                       bool `gson:"anonymous"`
}

func (this_ *Es6RewriteClass_ClassDeclarationMetadata) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteClass_ClassDeclarationMetadata) IsJavaNil_() bool { return this_ == nil }
func (*Es6RewriteClass_ClassDeclarationMetadata) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteClass$ClassDeclarationMetadata", "ClassDeclarationMetadata"
}
func (this_ *Es6RewriteClass_ClassDeclarationMetadata) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es6RewriteClass_ClassDeclarationMetadata) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteClass.ClassDeclarationMetadata", "com.google.javascript.jscomp.Es6RewriteClass$ClassDeclarationMetadata", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteClass_ClassDeclarationMetadata](v_)
	})
}

func (this_ *Es6RewriteClass_ClassDeclarationMetadata) initFields_Es6RewriteClass_ClassDeclarationMetadata() {
	_ = this_
}

// Java constructor: ClassDeclarationMetadata(com.google.javascript.jscomp.Es6RewriteClass.InsertionPoint insertionPoint, com.google.javascript.rhino.Node definePropertiesObjForPrototype, com.google.javascript.rhino.Node definePropertiesObjForClass, java.util.Map<java.lang.String,com.google.javascript.jscomp.Es6RewriteClass.ClassProperty> classMembersToDeclare, com.google.javascript.rhino.Node fullClassNameNode, com.google.javascript.rhino.Node classPrototypeNode, boolean anonymous, com.google.javascript.rhino.Node classNameNode, com.google.javascript.rhino.Node superClassNameNode)
func NewEs6RewriteClass_ClassDeclarationMetadata(insertionPoint *Es6RewriteClass_InsertionPoint, definePropertiesObjForPrototype pkg_rhino.Node, definePropertiesObjForClass pkg_rhino.Node, classMembersToDeclare j2gort.Ext_java_util_Map, fullClassNameNode pkg_rhino.Node, classPrototypeNode pkg_rhino.Node, anonymous bool, classNameNode pkg_rhino.Node, superClassNameNode pkg_rhino.Node) *Es6RewriteClass_ClassDeclarationMetadata {
	this_ := &Es6RewriteClass_ClassDeclarationMetadata{}
	InitEs6RewriteClass_ClassDeclarationMetadata(this_, insertionPoint, definePropertiesObjForPrototype, definePropertiesObjForClass, classMembersToDeclare, fullClassNameNode, classPrototypeNode, anonymous, classNameNode, superClassNameNode)
	return this_
}

func InitEs6RewriteClass_ClassDeclarationMetadata(this_ *Es6RewriteClass_ClassDeclarationMetadata, insertionPoint *Es6RewriteClass_InsertionPoint, definePropertiesObjForPrototype pkg_rhino.Node, definePropertiesObjForClass pkg_rhino.Node, classMembersToDeclare j2gort.Ext_java_util_Map, fullClassNameNode pkg_rhino.Node, classPrototypeNode pkg_rhino.Node, anonymous bool, classNameNode pkg_rhino.Node, superClassNameNode pkg_rhino.Node) {
	_ = insertionPoint
	_ = definePropertiesObjForPrototype
	_ = definePropertiesObjForClass
	_ = classMembersToDeclare
	_ = fullClassNameNode
	_ = classPrototypeNode
	_ = anonymous
	_ = classNameNode
	_ = superClassNameNode
	this_.initFields_Es6RewriteClass_ClassDeclarationMetadata()
	this_.insertionPoint = j2gort.CheckNotNullPtr(insertionPoint)
	this_.definePropertiesObjForPrototype = j2gort.CheckNotNullNode(definePropertiesObjForPrototype)
	this_.definePropertiesObjForClass = j2gort.CheckNotNullNode(definePropertiesObjForClass)
	this_.classMembersToDeclare = j2gort.CheckNotNull(classMembersToDeclare)
	this_.fullClassNameNode = j2gort.CheckNotNullNode(fullClassNameNode)
	this_.classPrototypeNode = j2gort.CheckNotNullNode(classPrototypeNode)
	this_.classNameNode = j2gort.CheckNotNullNode(classNameNode)
	this_.superClassNameNode = j2gort.CheckNotNullNode(superClassNameNode)
	this_.anonymous = anonymous
}

// Java method: ClassDeclarationMetadata.insertionPoint()
func (this_ *Es6RewriteClass_ClassDeclarationMetadata) InsertionPoint() *Es6RewriteClass_InsertionPoint {
	return this_.insertionPoint
}

// Java method: ClassDeclarationMetadata.definePropertiesObjForPrototype()
func (this_ *Es6RewriteClass_ClassDeclarationMetadata) DefinePropertiesObjForPrototype() pkg_rhino.Node {
	return this_.definePropertiesObjForPrototype
}

// Java method: ClassDeclarationMetadata.resetDefinePropertiesObjForPrototype(com.google.javascript.jscomp.AstFactory astFactory)
func (this_ *Es6RewriteClass_ClassDeclarationMetadata) ResetDefinePropertiesObjForPrototype(astFactory *AstFactory) {
	_ = astFactory
	this_.definePropertiesObjForPrototype = astFactory.CreateObjectLit_2(AstFactory_S_Type_1_Node(this_.fullClassNameNode), []pkg_rhino.Node{})
}

// Java method: ClassDeclarationMetadata.definePropertiesObjForClass()
func (this_ *Es6RewriteClass_ClassDeclarationMetadata) DefinePropertiesObjForClass() pkg_rhino.Node {
	return this_.definePropertiesObjForClass
}

// Java method: ClassDeclarationMetadata.resetDefinePropertiesObjForClass(com.google.javascript.jscomp.AstFactory astFactory)
func (this_ *Es6RewriteClass_ClassDeclarationMetadata) ResetDefinePropertiesObjForClass(astFactory *AstFactory) {
	_ = astFactory
	this_.definePropertiesObjForClass = astFactory.CreateObjectLit_2(AstFactory_S_Type_1_Node(this_.fullClassNameNode), []pkg_rhino.Node{})
}

// Java method: ClassDeclarationMetadata.classMembersToDeclare()
func (this_ *Es6RewriteClass_ClassDeclarationMetadata) ClassMembersToDeclare() j2gort.Ext_java_util_Map {
	return this_.classMembersToDeclare
}

// Java method: ClassDeclarationMetadata.fullClassNameNode()
func (this_ *Es6RewriteClass_ClassDeclarationMetadata) FullClassNameNode() pkg_rhino.Node {
	return this_.fullClassNameNode
}

// Java method: ClassDeclarationMetadata.classPrototypeNode()
func (this_ *Es6RewriteClass_ClassDeclarationMetadata) ClassPrototypeNode() pkg_rhino.Node {
	return this_.classPrototypeNode
}

// Java method: ClassDeclarationMetadata.anonymous()
func (this_ *Es6RewriteClass_ClassDeclarationMetadata) Anonymous() bool {
	return this_.anonymous
}

// Java method: ClassDeclarationMetadata.classNameNode()
func (this_ *Es6RewriteClass_ClassDeclarationMetadata) ClassNameNode() pkg_rhino.Node {
	return this_.classNameNode
}

// Java method: ClassDeclarationMetadata.superClassNameNode()
func (this_ *Es6RewriteClass_ClassDeclarationMetadata) SuperClassNameNode() pkg_rhino.Node {
	return this_.superClassNameNode
}

// Java static method: ClassDeclarationMetadata.builder()
func Es6RewriteClass_ClassDeclarationMetadata_S_Builder() Es6RewriteClass_ClassDeclarationMetadata_Builder {
	return NewAutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder().SetClassMembersToDeclare(j2gort.NewExt_java_util_LinkedHashMap())
}

// Java static method: ClassDeclarationMetadata.create(com.google.javascript.rhino.Node classNode, com.google.javascript.rhino.Node parent, com.google.javascript.jscomp.AstFactory astFactory)
func Es6RewriteClass_ClassDeclarationMetadata_Create(classNode pkg_rhino.Node, parent pkg_rhino.Node, astFactory *AstFactory) *Es6RewriteClass_ClassDeclarationMetadata {
	_ = classNode
	_ = parent
	_ = astFactory
	var classNameNode pkg_rhino.Node = classNode.GetFirstChild()
	_ = classNameNode
	var superClassNameNode pkg_rhino.Node = classNameNode.GetNext()
	_ = superClassNameNode
	var builder Es6RewriteClass_ClassDeclarationMetadata_Builder = Es6RewriteClass_ClassDeclarationMetadata_S_Builder().SetSuperClassNameNode(superClassNameNode).SetClassNameNode(classNameNode)
	_ = builder
	if NodeUtil_IsClassDeclaration(classNode) {
		builder.SetInsertionPoint(Es6RewriteClass_InsertionPoint_From(classNode)).SetFullClassNameNode(classNameNode).SetAnonymous(false)
	} else if parent.IsAssign() && parent.GetParent().IsExprResult() {
		var fullClassNameNode pkg_rhino.Node = parent.GetFirstChild()
		_ = fullClassNameNode
		if !(fullClassNameNode.IsQualifiedName()) {
			return nil
		}
		builder.SetInsertionPoint(Es6RewriteClass_InsertionPoint_From(parent.GetParent())).SetFullClassNameNode(fullClassNameNode).SetAnonymous(true)
	} else if parent.IsExport() {
		builder.SetInsertionPoint(Es6RewriteClass_InsertionPoint_From(classNode)).SetFullClassNameNode(classNameNode).SetAnonymous(false)
	} else if parent.IsName() {
		builder.SetInsertionPoint(Es6RewriteClass_InsertionPoint_From(parent.GetParent())).SetFullClassNameNode(parent.CloneNode_0()).SetAnonymous(true)
	} else {
		return nil
	}
	var classType AstFactory_Type = AstFactory_S_Type_1_Node(builder.FullClassNameNode())
	_ = classType
	builder.SetClassPrototypeNode(astFactory.CreatePrototypeAccess(builder.FullClassNameNode().CloneTree_0()))
	builder.SetDefinePropertiesObjForClass(astFactory.CreateObjectLit_2(classType, []pkg_rhino.Node{}))
	builder.SetDefinePropertiesObjForPrototype(astFactory.CreateObjectLit_2(classType, []pkg_rhino.Node{}))
	return builder.Build()
}

// Java method: ClassDeclarationMetadata.insertNodeAndAdvance(com.google.javascript.rhino.Node newNode)
func (this_ *Es6RewriteClass_ClassDeclarationMetadata) InsertNodeAndAdvance(newNode pkg_rhino.Node) {
	_ = newNode
	this_.InsertionPoint().InsertNodeAndAdvance(newNode)
}

// Java method: ClassDeclarationMetadata.hasSuperClass()
func (this_ *Es6RewriteClass_ClassDeclarationMetadata) HasSuperClass() bool {
	return !(this_.SuperClassNameNode().IsEmpty())
}

// Java polymorphic class interface: com.google.javascript.jscomp.Es6RewriteClass.ClassDeclarationMetadata.Builder
type Es6RewriteClass_ClassDeclarationMetadata_Builder interface {
	Implements_jscomp_Es6RewriteClass_ClassDeclarationMetadata_Builder()
	IsJavaNil_() bool
	GetEs6RewriteClass_ClassDeclarationMetadata_BuilderBase() *Es6RewriteClass_ClassDeclarationMetadata_BuilderBase
	SetInsertionPoint(*Es6RewriteClass_InsertionPoint) Es6RewriteClass_ClassDeclarationMetadata_Builder
	SetFullClassNameNode(pkg_rhino.Node) Es6RewriteClass_ClassDeclarationMetadata_Builder
	FullClassNameNode() pkg_rhino.Node
	SetClassMembersToDeclare(j2gort.Ext_java_util_Map) Es6RewriteClass_ClassDeclarationMetadata_Builder
	SetAnonymous(bool) Es6RewriteClass_ClassDeclarationMetadata_Builder
	SetClassNameNode(pkg_rhino.Node) Es6RewriteClass_ClassDeclarationMetadata_Builder
	SetSuperClassNameNode(pkg_rhino.Node) Es6RewriteClass_ClassDeclarationMetadata_Builder
	SetClassPrototypeNode(pkg_rhino.Node) Es6RewriteClass_ClassDeclarationMetadata_Builder
	SetDefinePropertiesObjForClass(pkg_rhino.Node) Es6RewriteClass_ClassDeclarationMetadata_Builder
	SetDefinePropertiesObjForPrototype(pkg_rhino.Node) Es6RewriteClass_ClassDeclarationMetadata_Builder
	Build() *Es6RewriteClass_ClassDeclarationMetadata
}

func ToIface_Es6RewriteClass_ClassDeclarationMetadata_Builder[T interface {
	Es6RewriteClass_ClassDeclarationMetadata_Builder
	comparable
}](v T) Es6RewriteClass_ClassDeclarationMetadata_Builder {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_Es6RewriteClass_ClassDeclarationMetadata_Builder[T Es6RewriteClass_ClassDeclarationMetadata_Builder](v T) Es6RewriteClass_ClassDeclarationMetadata_Builder {
	var vi Es6RewriteClass_ClassDeclarationMetadata_Builder = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.Es6RewriteClass.ClassDeclarationMetadata.Builder
type Es6RewriteClass_ClassDeclarationMetadata_BuilderBase struct {
	idHash_ atomic.Int32
}

func (this_ *Es6RewriteClass_ClassDeclarationMetadata_BuilderBase) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteClass_ClassDeclarationMetadata_BuilderBase) IsJavaNil_() bool {
	return this_ == nil
}
func (*Es6RewriteClass_ClassDeclarationMetadata_BuilderBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteClass$ClassDeclarationMetadata$Builder", "Builder"
}
func (this_ *Es6RewriteClass_ClassDeclarationMetadata_BuilderBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteClass.ClassDeclarationMetadata.Builder", "com.google.javascript.jscomp.Es6RewriteClass$ClassDeclarationMetadata$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[Es6RewriteClass_ClassDeclarationMetadata_Builder](v_)
	})
}
func (b *Es6RewriteClass_ClassDeclarationMetadata_BuilderBase) GetEs6RewriteClass_ClassDeclarationMetadata_BuilderBase() *Es6RewriteClass_ClassDeclarationMetadata_BuilderBase {
	return b
}
func (*Es6RewriteClass_ClassDeclarationMetadata_BuilderBase) Implements_jscomp_Es6RewriteClass_ClassDeclarationMetadata_Builder() {
}

func (this_ *Es6RewriteClass_ClassDeclarationMetadata_BuilderBase) initFields_Es6RewriteClass_ClassDeclarationMetadata_Builder() {
	_ = this_
}

func InitEs6RewriteClass_ClassDeclarationMetadata_Builder(this_ *Es6RewriteClass_ClassDeclarationMetadata_BuilderBase) {
	this_.initFields_Es6RewriteClass_ClassDeclarationMetadata_Builder()
}

// Java method: Builder.setInsertionPoint(com.google.javascript.jscomp.Es6RewriteClass.InsertionPoint insertionPoint)
func (this_ *Es6RewriteClass_ClassDeclarationMetadata_BuilderBase) SetInsertionPoint(insertionPoint *Es6RewriteClass_InsertionPoint) Es6RewriteClass_ClassDeclarationMetadata_Builder {
	_ = insertionPoint
	panic("abstract or native method invoked")
}

// Java method: Builder.setFullClassNameNode(com.google.javascript.rhino.Node fullClassNameNode)
func (this_ *Es6RewriteClass_ClassDeclarationMetadata_BuilderBase) SetFullClassNameNode(fullClassNameNode pkg_rhino.Node) Es6RewriteClass_ClassDeclarationMetadata_Builder {
	_ = fullClassNameNode
	panic("abstract or native method invoked")
}

// Java method: Builder.fullClassNameNode()
func (this_ *Es6RewriteClass_ClassDeclarationMetadata_BuilderBase) FullClassNameNode() pkg_rhino.Node {
	panic("abstract or native method invoked")
}

// Java method: Builder.setClassMembersToDeclare(java.util.Map<java.lang.String,com.google.javascript.jscomp.Es6RewriteClass.ClassProperty> classMembersToDeclare)
func (this_ *Es6RewriteClass_ClassDeclarationMetadata_BuilderBase) SetClassMembersToDeclare(classMembersToDeclare j2gort.Ext_java_util_Map) Es6RewriteClass_ClassDeclarationMetadata_Builder {
	_ = classMembersToDeclare
	panic("abstract or native method invoked")
}

// Java method: Builder.setAnonymous(boolean anonymous)
func (this_ *Es6RewriteClass_ClassDeclarationMetadata_BuilderBase) SetAnonymous(anonymous bool) Es6RewriteClass_ClassDeclarationMetadata_Builder {
	_ = anonymous
	panic("abstract or native method invoked")
}

// Java method: Builder.setClassNameNode(com.google.javascript.rhino.Node classNameNode)
func (this_ *Es6RewriteClass_ClassDeclarationMetadata_BuilderBase) SetClassNameNode(classNameNode pkg_rhino.Node) Es6RewriteClass_ClassDeclarationMetadata_Builder {
	_ = classNameNode
	panic("abstract or native method invoked")
}

// Java method: Builder.setSuperClassNameNode(com.google.javascript.rhino.Node superClassNameNode)
func (this_ *Es6RewriteClass_ClassDeclarationMetadata_BuilderBase) SetSuperClassNameNode(superClassNameNode pkg_rhino.Node) Es6RewriteClass_ClassDeclarationMetadata_Builder {
	_ = superClassNameNode
	panic("abstract or native method invoked")
}

// Java method: Builder.setClassPrototypeNode(com.google.javascript.rhino.Node node)
func (this_ *Es6RewriteClass_ClassDeclarationMetadata_BuilderBase) SetClassPrototypeNode(node pkg_rhino.Node) Es6RewriteClass_ClassDeclarationMetadata_Builder {
	_ = node
	panic("abstract or native method invoked")
}

// Java method: Builder.setDefinePropertiesObjForClass(com.google.javascript.rhino.Node node)
func (this_ *Es6RewriteClass_ClassDeclarationMetadata_BuilderBase) SetDefinePropertiesObjForClass(node pkg_rhino.Node) Es6RewriteClass_ClassDeclarationMetadata_Builder {
	_ = node
	panic("abstract or native method invoked")
}

// Java method: Builder.setDefinePropertiesObjForPrototype(com.google.javascript.rhino.Node node)
func (this_ *Es6RewriteClass_ClassDeclarationMetadata_BuilderBase) SetDefinePropertiesObjForPrototype(node pkg_rhino.Node) Es6RewriteClass_ClassDeclarationMetadata_Builder {
	_ = node
	panic("abstract or native method invoked")
}

// Java method: Builder.build()
func (this_ *Es6RewriteClass_ClassDeclarationMetadata_BuilderBase) Build() *Es6RewriteClass_ClassDeclarationMetadata {
	panic("abstract or native method invoked")
}

// Java class: com.google.javascript.jscomp.Es6RewriteClass.InsertionPoint
type Es6RewriteClass_InsertionPoint struct {
	insertionPoint pkg_rhino.Node `gson:"insertionPoint"`
	idHash_        atomic.Int32
}

func (this_ *Es6RewriteClass_InsertionPoint) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteClass_InsertionPoint) IsJavaNil_() bool { return this_ == nil }
func (*Es6RewriteClass_InsertionPoint) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteClass$InsertionPoint", "InsertionPoint"
}
func (this_ *Es6RewriteClass_InsertionPoint) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es6RewriteClass_InsertionPoint) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteClass.InsertionPoint", "com.google.javascript.jscomp.Es6RewriteClass$InsertionPoint", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteClass_InsertionPoint](v_)
	})
}

func (this_ *Es6RewriteClass_InsertionPoint) initFields_Es6RewriteClass_InsertionPoint() {
	_ = this_
}

// Java constructor: InsertionPoint(com.google.javascript.rhino.Node insertionPoint)
func NewEs6RewriteClass_InsertionPoint(insertionPoint pkg_rhino.Node) *Es6RewriteClass_InsertionPoint {
	this_ := &Es6RewriteClass_InsertionPoint{}
	InitEs6RewriteClass_InsertionPoint(this_, insertionPoint)
	return this_
}

func InitEs6RewriteClass_InsertionPoint(this_ *Es6RewriteClass_InsertionPoint, insertionPoint pkg_rhino.Node) {
	_ = insertionPoint
	this_.initFields_Es6RewriteClass_InsertionPoint()
	this_.insertionPoint = insertionPoint
}

// Java method: InsertionPoint.insertNodeAndAdvance(com.google.javascript.rhino.Node newNode)
func (this_ *Es6RewriteClass_InsertionPoint) InsertNodeAndAdvance(newNode pkg_rhino.Node) {
	_ = newNode
	newNode.InsertAfter(this_.insertionPoint)
	this_.insertionPoint = newNode
}

// Java static method: InsertionPoint.from(com.google.javascript.rhino.Node start)
func Es6RewriteClass_InsertionPoint_From(start pkg_rhino.Node) *Es6RewriteClass_InsertionPoint {
	_ = start
	return NewEs6RewriteClass_InsertionPoint(start)
}

// Java method: InsertionPoint.getNode()
func (this_ *Es6RewriteClass_InsertionPoint) GetNode() pkg_rhino.Node {
	return this_.insertionPoint
}
func (*Es6RewriteClass) CanReleaseSyntacticScopes_() {}
