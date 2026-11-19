/*
 * Copyright 2025 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/Es6NormalizeClasses.java (generated mechanically by j2go).

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

var lambdaSingleton_es6_normalize_classes_1 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_es6_normalize_classes_2 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_es6_normalize_classes_3 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_es6_normalize_classes_4 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_es6_normalize_classes_5 j2gort.Ext_com_google_common_base_Predicate

func init() {
	lambdaSingleton_es6_normalize_classes_1 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_10 uint32) bool {
		n := pkg_rhino.Node(rawNode_10)
		_ = n
		return (!(n.IsClass()) && (!(n.IsFunction()) || n.IsArrowFunction()))
	})
	lambdaSingleton_es6_normalize_classes_2 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_13 uint32) bool {
		n := pkg_rhino.Node(rawNode_13)
		_ = n
		return (!(n.IsClass()) && (!(n.IsFunction()) || n.IsArrowFunction()))
	})
	lambdaSingleton_es6_normalize_classes_3 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_14 uint32) bool {
		n := pkg_rhino.Node(rawNode_14)
		_ = n
		return !(n.IsClass())
	})
	lambdaSingleton_es6_normalize_classes_4 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_27 uint32) bool {
		n := pkg_rhino.Node(rawNode_27)
		_ = n
		return (n.IsSuper() || n.IsPrivateIdentifier())
	})
	lambdaSingleton_es6_normalize_classes_5 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_28 uint32) bool {
		n := pkg_rhino.Node(rawNode_28)
		_ = n
		return (!(n.IsClass()) && (!(n.IsFunction()) || n.IsArrowFunction()))
	})
}

// Java class: com.google.javascript.jscomp.Es6NormalizeClasses
type Es6NormalizeClasses struct {
	compiler                    AbstractCompiler                               `gson:"compiler"`
	astFactory                  *AstFactory                                    `gson:"astFactory"`
	expressionDecomposer        *ExpressionDecomposer                          `gson:"expressionDecomposer"`
	ctorCreator                 *SynthesizeExplicitConstructors                `gson:"ctorCreator"`
	jscompPrivateMap            RuntimeJsLibManager_JsLibField                 `gson:"jscompPrivateMap"`
	classStack                  j2gort.Deque[*Es6NormalizeClasses_ClassRecord] `gson:"classStack"`
	idHash_                     atomic.Int32
	transpileClassFields        bool `gson:"transpileClassFields"`
	transpilePrivateClassFields bool `gson:"transpilePrivateClassFields"`
	doStaticInheritanceRewrites bool `gson:"doStaticInheritanceRewrites"`
	doStaticSuperRewrites       bool `gson:"doStaticSuperRewrites"`
}

func (this_ *Es6NormalizeClasses) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6NormalizeClasses) IsJavaNil_() bool { return this_ == nil }
func (*Es6NormalizeClasses) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6NormalizeClasses", "Es6NormalizeClasses"
}
func (this_ *Es6NormalizeClasses) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es6NormalizeClasses) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6NormalizeClasses", "com.google.javascript.jscomp.Es6NormalizeClasses", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6NormalizeClasses](v_)
	})
}
func (*Es6NormalizeClasses) Implements_jscomp_NodeTraversal_ScopedCallback() {}
func (*Es6NormalizeClasses) Implements_jscomp_CompilerPass()                 {}
func (*Es6NormalizeClasses) Implements_jscomp_NodeTraversal_Callback()       {}

func StaticGet_Es6NormalizeClasses_CLASS_DECL_VAR() string          { return "$jscomp$classDecl$" }
func StaticGet_Es6NormalizeClasses_CLASS_EXTENDS_VAR() string       { return "$jscomp$classExtends$" }
func StaticGet_Es6NormalizeClasses_COMP_FIELD_VAR() string          { return "$jscomp$compField$" }
func StaticGet_Es6NormalizeClasses_STATIC_INIT_METHOD_NAME() string { return "$jscomp$staticInit$" }
func StaticGet_Es6NormalizeClasses_PRIVATE_MAP_VAR() string         { return "$jscomp$privateMap$" }
func StaticGet_Es6NormalizeClasses_STATIC_PRIVATE_MAP_VAR() string {
	return "$jscomp$staticPrivateMap$"
}
func StaticGet_Es6NormalizeClasses_PRIVATE_PROTO_VAR() string { return "$jscomp$priv$proto$" }
func StaticGet_Es6NormalizeClasses_PRIVATE_VAR() string       { return "$jscomp$priv$" }
func StaticGet_Es6NormalizeClasses_TMP_VAR() string           { return "$jscomp$tmp$" }

var static_Es6NormalizeClasses_ILLEGAL_PRIVATE_MEMBER_ASSIGNMENT *DiagnosticType

func StaticGet_Es6NormalizeClasses_ILLEGAL_PRIVATE_MEMBER_ASSIGNMENT() *DiagnosticType {
	clinit_Es6NormalizeClasses()
	return static_Es6NormalizeClasses_ILLEGAL_PRIVATE_MEMBER_ASSIGNMENT
}
func StaticSet_Es6NormalizeClasses_ILLEGAL_PRIVATE_MEMBER_ASSIGNMENT(v *DiagnosticType) {
	clinit_Es6NormalizeClasses()
	static_Es6NormalizeClasses_ILLEGAL_PRIVATE_MEMBER_ASSIGNMENT = v
}

var static_Es6NormalizeClasses_GENERIC_NAME_REPLACEMENTS j2gort.Ext_com_google_common_collect_ImmutableMap

func StaticGet_Es6NormalizeClasses_GENERIC_NAME_REPLACEMENTS() j2gort.Ext_com_google_common_collect_ImmutableMap {
	clinit_Es6NormalizeClasses()
	return static_Es6NormalizeClasses_GENERIC_NAME_REPLACEMENTS
}
func StaticSet_Es6NormalizeClasses_GENERIC_NAME_REPLACEMENTS(v j2gort.Ext_com_google_common_collect_ImmutableMap) {
	clinit_Es6NormalizeClasses()
	static_Es6NormalizeClasses_GENERIC_NAME_REPLACEMENTS = v
}

var clinitDone_Es6NormalizeClasses uint32
var clinitState_Es6NormalizeClasses j2gort.ClinitState

func clinit_Es6NormalizeClasses() {
	if atomic.LoadUint32(&clinitDone_Es6NormalizeClasses) != 0 {
		return
	}
	clinitSlow_Es6NormalizeClasses()
}

func clinitSlow_Es6NormalizeClasses() {
	if !j2gort.EnterClinit(&clinitState_Es6NormalizeClasses, &clinitDone_Es6NormalizeClasses) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Es6NormalizeClasses, &clinitDone_Es6NormalizeClasses)
	static_Es6NormalizeClasses_ILLEGAL_PRIVATE_MEMBER_ASSIGNMENT = DiagnosticType_Error("JSC_ILLEGAL_PRIVATE_MEMBER_ASSIGNMENT", "Invalid assignment to private method ''#{0}''. Private methods cannot be reassigned.")
	static_Es6NormalizeClasses_GENERIC_NAME_REPLACEMENTS = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder().Put_2("CLASS_DECL", "$jscomp$classDecl$")).Put_2("CLASS_EXTENDS", "$jscomp$classExtends$")).Put_2("COMP_FIELD", "$jscomp$compField$")).Put_2("STATIC_INIT", "$jscomp$staticInit$")).Put_2("STATIC_PRIVATE_MAP", "$jscomp$staticPrivateMap$")).Put_2("PRIVATE_MAP", "$jscomp$privateMap$")).Put_2("PRIVATE_PROTO", "$jscomp$priv$proto$")).Put_2("PRIVATE", "$jscomp$priv$")).Put_2("TMP", "$jscomp$tmp$")).BuildOrThrow())
}

func (this_ *Es6NormalizeClasses) initFields_Es6NormalizeClasses() {
	_ = this_
	this_.classStack = j2gort.Deque[*Es6NormalizeClasses_ClassRecord]{}
}

// Java constructor: Es6NormalizeClasses(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewEs6NormalizeClasses(compiler AbstractCompiler) *Es6NormalizeClasses {
	clinit_Es6NormalizeClasses()
	this_ := &Es6NormalizeClasses{}
	InitEs6NormalizeClasses(this_, compiler)
	return this_
}

func InitEs6NormalizeClasses(this_ *Es6NormalizeClasses, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_Es6NormalizeClasses()
	this_.compiler = compiler
	this_.astFactory = compiler.CreateAstFactory()
	this_.expressionDecomposer = compiler.CreateDefaultExpressionDecomposer()
	this_.ctorCreator = NewSynthesizeExplicitConstructors(compiler)
	var options *CompilerOptions = compiler.GetOptions()
	_ = options
	this_.transpileClassFields = options.NeedsTranspilationOf(pkg_parser.StaticGet_FeatureSet_Feature_PUBLIC_CLASS_FIELDS())
	this_.transpilePrivateClassFields = options.NeedsTranspilationOf(pkg_parser.StaticGet_FeatureSet_Feature_PRIVATE_ELEMENTS())
	this_.doStaticInheritanceRewrites = (options.GetAssumeStaticInheritanceIsNotUsed() && !(options.GetSkipNonTranspilationPasses()))
	this_.doStaticSuperRewrites = !(options.NeedsTranspilationOf(pkg_parser.StaticGet_FeatureSet_Feature_SUPER()))
	this_.jscompPrivateMap = compiler.GetRuntimeJsLibManager().GetJsLibField("$jscomp.PrivateMap")
}

// Java method: Es6NormalizeClasses.generateUniqueTmpVarName(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_generateUniqueTmpVarName_2494(t *NodeTraversal) string {
	_ = t
	return ("$jscomp$tmp$" + j2gort.StrValString(this_.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())))
}

// Java method: Es6NormalizeClasses.generateUniqueClassDeclVarName(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_generateUniqueClassDeclVarName_2495(t *NodeTraversal) string {
	_ = t
	return ("$jscomp$classDecl$" + j2gort.StrValString(this_.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())))
}

// Java method: Es6NormalizeClasses.generateUniqueClassExtendsVarName(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_generateUniqueClassExtendsVarName_2496(t *NodeTraversal) string {
	_ = t
	return ("$jscomp$classExtends$" + j2gort.StrValString(this_.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())))
}

// Java method: Es6NormalizeClasses.generateUniqueCompFieldVarName(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_generateUniqueCompFieldVarName_2497(t *NodeTraversal) string {
	_ = t
	return ("$jscomp$compField$" + j2gort.StrValString(this_.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())))
}

// Java method: Es6NormalizeClasses.generateUniqueStaticInitMethodName(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_generateUniqueStaticInitMethodName_2498(t *NodeTraversal) string {
	_ = t
	return ("$jscomp$staticInit$" + j2gort.StrValString(this_.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())))
}

// Java static method: Es6NormalizeClasses.generateUniquePrivateMapVarName(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.NodeTraversal t)
func Es6NormalizeClasses_GenerateUniquePrivateMapVarName(compiler AbstractCompiler, t *NodeTraversal) string {
	clinit_Es6NormalizeClasses()
	_ = compiler
	_ = t
	return ("$jscomp$privateMap$" + j2gort.StrValString(compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())))
}

// Java static method: Es6NormalizeClasses.generateUniqueStaticPrivateMapVarName(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.NodeTraversal t)
func Es6NormalizeClasses_GenerateUniqueStaticPrivateMapVarName(compiler AbstractCompiler, t *NodeTraversal) string {
	clinit_Es6NormalizeClasses()
	_ = compiler
	_ = t
	return ("$jscomp$staticPrivateMap$" + j2gort.StrValString(compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())))
}

// Java static method: Es6NormalizeClasses.generateUniquePrivateProtoVarName(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.NodeTraversal t)
func Es6NormalizeClasses_GenerateUniquePrivateProtoVarName(compiler AbstractCompiler, t *NodeTraversal) string {
	clinit_Es6NormalizeClasses()
	_ = compiler
	_ = t
	return ("$jscomp$priv$proto$" + j2gort.StrValString(compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())))
}

// Java static method: Es6NormalizeClasses.generateUniquePrivateVarName(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.NodeTraversal t)
func Es6NormalizeClasses_GenerateUniquePrivateVarName(compiler AbstractCompiler, t *NodeTraversal) string {
	clinit_Es6NormalizeClasses()
	_ = compiler
	_ = t
	return ("$jscomp$priv$" + j2gort.StrValString(compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())))
}

// Java method: Es6NormalizeClasses.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *Es6NormalizeClasses) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	sc_ := NewSyntacticScopeCreator_1(this_.compiler)
	defer func() {
		this_.classStack.Clear()
		sc_.ReleaseSlabs_()
	}()
	{
		clinit_NodeTraversal()
		var ntb_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_)
		(&ntb_).SetCompiler(this_.compiler)
		(&ntb_).SetCallback_1_NodeTraversal_Callback(this_)
		(&ntb_).SetScopeCreator(sc_)
		(&ntb_).TraverseRoots(externs, root)
	}

	var transpiledFeatures *pkg_parser.FeatureSet = pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_CLASS_STATIC_BLOCK())
	_ = transpiledFeatures
	if this_.transpileClassFields {
		transpiledFeatures = transpiledFeatures.With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_PUBLIC_CLASS_FIELDS())
	}
	if this_.transpilePrivateClassFields {
		transpiledFeatures = transpiledFeatures.With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_PRIVATE_ELEMENTS())
	}
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(this_.compiler, root, transpiledFeatures)
	GatherGetterAndSetterProperties_Update(this_.compiler, externs, root)
}

// Java method: Es6NormalizeClasses.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6NormalizeClasses) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(94) /* Token.SCRIPT */ :
		var scriptFeatures *pkg_parser.FeatureSet = NodeUtil_GetFeatureSetOfScript(n)
		_ = scriptFeatures
		return (((((scriptFeatures == nil) || scriptFeatures.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_CLASSES())) || scriptFeatures.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_CLASS_STATIC_BLOCK())) || scriptFeatures.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_PUBLIC_CLASS_FIELDS())) || scriptFeatures.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_PRIVATE_ELEMENTS()))
	case int32(105) /* Token.CLASS */ :
		var bindingIdentifier *j2gort.Ext_java_util_Optional
		_ = bindingIdentifier
		if n.GetFirstChild().IsName() {
			bindingIdentifier = j2gort.Ext_java_util_Optional_Of(j2gort.BoxNode(n.GetFirstChild()))
		} else {
			bindingIdentifier = j2gort.Ext_java_util_Optional_Empty()
		}
		var normalized *Es6NormalizeClasses_NormalizedNames = this_.priv_Es6NormalizeClasses_generateNormalizedClassNames_2536(t, n)
		_ = normalized
		var val_1 Es6NormalizeClasses_ClassRecord
		InitEs6NormalizeClasses_ClassRecord(&val_1, normalized, bindingIdentifier, n)
		var record *Es6NormalizeClasses_ClassRecord = &val_1
		_ = record
		var classMembers pkg_rhino.Node = n.GetLastChild()
		_ = classMembers
		if classMembers.IsClassMembers() {
			{
				var member pkg_rhino.Node = classMembers.GetFirstChild()
				_ = member
				for ; member != 0; member = member.GetNext() {
					if Es6NormalizeClasses_IsPrivateMember(member) {
						record.RecordPrivateMember(member)
					}
				}
			}
		}
		this_.classStack.AddFirst(record)
	case int32(109) /* Token.COMPUTED_FIELD_DEF */ :
		if !(!(this_.classStack.IsEmpty())) {
			j2gort.CheckState(false)
		}
		var classRecord *Es6NormalizeClasses_ClassRecord = this_.classStack.Peek()
		_ = classRecord
		classRecord.MaybeRecordComputedPropWithSideEffects(n)
		classRecord.RecordField(n)
	case int32(108) /* Token.MEMBER_FIELD_DEF */ :
		if !(!(this_.classStack.IsEmpty())) {
			j2gort.CheckState(false)
		}
		var classRecord_2 *Es6NormalizeClasses_ClassRecord = this_.classStack.Peek()
		_ = classRecord_2
		if n.IsPrivateIdentifier() {
			classRecord_2.RecordPrivateMember(n)
		}
		classRecord_2.RecordField(n)
	case int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */ :
		if n.IsPrivateIdentifier() {
			if !(!(this_.classStack.IsEmpty())) {
				j2gort.CheckState(false)
			}
			this_.classStack.Peek().RecordPrivateMember(n)
		}
	case int32(90) /* Token.BLOCK */ :
		if !(NodeUtil_IsClassStaticBlock(n)) {
			break
		}
		if !(!(this_.classStack.IsEmpty())) {
			j2gort.CheckState(false)
		}
		this_.classStack.Peek().RecordStaticBlock(n)
	case int32(129) /* Token.COMPUTED_PROP */ :
		if !(n.GetParent().IsClassMembers()) {
			break
		}
		if !(!(this_.classStack.IsEmpty())) {
			j2gort.CheckState(false)
		}
		this_.classStack.Peek().MaybeRecordComputedPropWithSideEffects(n)
	default:
	}
	return true
}

// Java method: Es6NormalizeClasses.shouldExtractExtends(com.google.javascript.rhino.Node classNode)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_shouldExtractExtends_2499(classNode pkg_rhino.Node) bool {
	_ = classNode
	if !(classNode.IsClass()) {
		j2gort.CheckArgument(false)
	}
	var superClassNode pkg_rhino.Node = classNode.GetSecondChild()
	_ = superClassNode
	return (!(superClassNode.IsEmpty()) && !(superClassNode.IsQualifiedName()))
}

// Java method: Es6NormalizeClasses.canExtractExtendsSimply(com.google.javascript.rhino.Node classNode)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_canExtractExtendsSimply_2500(classNode pkg_rhino.Node) bool {
	_ = classNode
	var enclosingStatement pkg_rhino.Node = j2gort.CheckNotNullNode(NodeUtil_GetEnclosingStatement(classNode))
	_ = enclosingStatement
	if enclosingStatement == classNode {
		return true
	}
	var classNodeParent pkg_rhino.Node = classNode.GetParent()
	_ = classNodeParent
	if (NodeUtil_IsNameDeclaration(enclosingStatement) && classNodeParent.IsName()) && classNodeParent.IsFirstChildOf(enclosingStatement) {
		return true
	}
	if ((enclosingStatement.IsExprResult() && classNodeParent.IsOnlyChildOf(enclosingStatement)) && classNodeParent.IsAssign()) && classNode.IsSecondChildOf(classNodeParent) {
		var lhsNode pkg_rhino.Node = classNodeParent.GetFirstChild()
		_ = lhsNode
		return !(this_.compiler.GetAstAnalyzer().MayHaveSideEffects(lhsNode))
	}
	return false
}

// Java method: Es6NormalizeClasses.extractExtends(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.Es6NormalizeClasses.ClassRecord classRecord)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_extractExtends_2501(t *NodeTraversal, classRecord *Es6NormalizeClasses_ClassRecord) {
	_ = t
	_ = classRecord
	var classNode pkg_rhino.Node = classRecord.classNode
	_ = classNode
	var statement pkg_rhino.Node = NodeUtil_GetEnclosingStatement(classNode)
	_ = statement
	var originalExtends pkg_rhino.Node = classNode.GetSecondChild()
	_ = originalExtends
	if originalExtends.IsClass() && originalExtends.GetFirstChild().IsName() {
		originalExtends.GetFirstChild().ReplaceWith(pkg_rhino.IR_Empty().Srcref(originalExtends.GetFirstChild()))
	}
	var nameNode pkg_rhino.Node = j2gort.CastNode[pkg_rhino.Node](classRecord.normalizedNames.superClassNode.Get())
	_ = nameNode
	originalExtends.ReplaceWith(nameNode)
	var extendsAlias pkg_rhino.Node = pkg_rhino.IR_ConstNode(nameNode.CloneNode_0(), originalExtends).Srcref(originalExtends)
	_ = extendsAlias
	extendsAlias.InsertBefore(statement)
	NodeUtil_AddFeatureToScript(NodeUtil_GetEnclosingScript(classNode), pkg_parser.StaticGet_FeatureSet_Feature_CONST_DECLARATIONS(), this_.compiler)
	t.ReportCodeChange_1(classNode)
}

// Java method: Es6NormalizeClasses.pushDownNodeIntoNewIife(com.google.javascript.rhino.Node n)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_pushDownNodeIntoNewIife_2502(n pkg_rhino.Node) {
	_ = n
	var functionBody pkg_rhino.Node = pkg_rhino.IR_Block_0()
	_ = functionBody
	var arrowFn pkg_rhino.Node = pkg_rhino.IR_ArrowFunction(pkg_rhino.IR_Name(""), pkg_rhino.IR_ParamList([]pkg_rhino.Node{}), functionBody).Srcref(n)
	_ = arrowFn
	arrowFn.SetColor(pkg_colors.StaticGet_StandardColors_UNKNOWN())
	var iife pkg_rhino.Node = this_.astFactory.CreateCallWithUnknownType(arrowFn, []pkg_rhino.Node{}).SrcrefTreeIfMissing(n)
	_ = iife
	n.ReplaceWith(iife)
	functionBody.AddChildToBack(this_.astFactory.CreateReturn(n).Srcref(n))
	NodeUtil_AddFeatureToScript(NodeUtil_GetEnclosingScript(iife), pkg_parser.StaticGet_FeatureSet_Feature_ARROW_FUNCTIONS(), this_.compiler)
	this_.compiler.ReportChangeToEnclosingScope(iife)
}

// Java method: Es6NormalizeClasses.shouldExtractClass(com.google.javascript.rhino.Node classNode)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_shouldExtractClass_2503(classNode pkg_rhino.Node) bool {
	_ = classNode
	var parent pkg_rhino.Node = classNode.GetParent()
	_ = parent
	if (NodeUtil_IsClassDeclaration(classNode) || parent.IsName()) || ((parent.IsAssign() && parent.GetFirstChild().IsQualifiedName()) && parent.GetParent().IsExprResult()) {
		return false
	}
	return true
}

// Java method: Es6NormalizeClasses.extractClass(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.Es6NormalizeClasses.ClassRecord classRecord)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_extractClass_2504(t *NodeTraversal, classRecord *Es6NormalizeClasses_ClassRecord) {
	_ = t
	_ = classRecord
	var classNode pkg_rhino.Node = classRecord.classNode
	_ = classNode
	var parent pkg_rhino.Node = classNode.GetParent()
	_ = parent
	var info *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(classNode)
	_ = info
	var statement pkg_rhino.Node = NodeUtil_GetEnclosingStatement(parent)
	_ = statement
	var classNameLhs pkg_rhino.Node = classRecord.CreateNewNameReferenceNode()
	_ = classNameLhs
	var classNameRhs pkg_rhino.Node = classNameLhs.CloneTree_0()
	_ = classNameRhs
	classNode.ReplaceWith(classNameRhs)
	var classDeclaration pkg_rhino.Node = pkg_rhino.IR_ConstNode(classNameLhs, classNode).SrcrefTreeIfMissing(classNode)
	_ = classDeclaration
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_CONST_DECLARATIONS(), this_.compiler)
	classDeclaration.SetJSDocInfo(pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(info).Build_0())
	classDeclaration.InsertBefore(statement)
	if (NodeUtil_IsNameDeclaration(statement) && statement.HasOneChild()) && (statement.GetOnlyChild() == parent) {
		this_.priv_Es6NormalizeClasses_addAtConstructor_2505(statement)
	} else if statement.IsExprResult() {
		var expr pkg_rhino.Node = statement.GetOnlyChild()
		_ = expr
		if (expr.IsAssign() && expr.GetFirstChild().IsQualifiedName()) && (expr.GetSecondChild() == classNameRhs) {
			this_.priv_Es6NormalizeClasses_addAtConstructor_2505(expr)
		}
	}
	this_.compiler.ReportChangeToEnclosingScope(classDeclaration)
}

// Java method: Es6NormalizeClasses.addAtConstructor(com.google.javascript.rhino.Node node)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_addAtConstructor_2505(node pkg_rhino.Node) {
	_ = node
	var builder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(node.GetJSDocInfo())
	_ = builder
	builder.RecordConstructor()
	node.SetJSDocInfo(builder.Build_0())
}

// Java method: Es6NormalizeClasses.getStatementDeclaringClass(com.google.javascript.rhino.Node classNode, com.google.javascript.rhino.Node classNameNode)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_getStatementDeclaringClass_2506(classNode pkg_rhino.Node, classNameNode pkg_rhino.Node) pkg_rhino.Node {
	_ = classNode
	_ = classNameNode
	if NodeUtil_IsClassDeclaration(classNode) {
		if !(NodeUtil_IsStatement(classNode)) {
			j2gort.CheckState(false)
		}
		return classNode
	}
	var parent pkg_rhino.Node = classNode.GetParent()
	_ = parent
	if parent.IsName() {
		if !(parent == classNameNode) {
			j2gort.CheckState(false)
		}
		if !(NodeUtil_IsStatement(classNameNode.GetParent())) {
			j2gort.CheckState(false)
		}
		return classNameNode.GetParent()
	}
	if (parent.IsAssign() && (parent.GetFirstChild() == classNameNode)) && parent.GetParent().IsExprResult() {
		if !(NodeUtil_IsStatement(classNameNode.GetGrandparent())) {
			j2gort.CheckState(false)
		}
		return classNameNode.GetGrandparent()
	}
	return pkg_rhino.Node(0)
}

// Java method: Es6NormalizeClasses.enterScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *Es6NormalizeClasses) EnterScope(t *NodeTraversal) {
	_ = t
	var scopeRoot pkg_rhino.Node = t.GetScopeRoot()
	_ = scopeRoot
	if NodeUtil_IsFunctionBlock(scopeRoot) && NodeUtil_IsEs6Constructor(scopeRoot.GetParent()) {
		this_.classStack.Peek().RecordConstructorScope(t.GetScope())
	}
}

// Java method: Es6NormalizeClasses.exitScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *Es6NormalizeClasses) ExitScope(t *NodeTraversal) {
	_ = t
}

// Java method: Es6NormalizeClasses.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6NormalizeClasses) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(44) /* Token.IN */ :
		if this_.transpilePrivateClassFields && n.GetFirstChild().IsPrivateIdentifier() {
			this_.priv_Es6NormalizeClasses_rewritePrivateIn_2509(t, n)
		}
	case int32(26) /* Token.GETPROP */, int32(29) /* Token.OPTCHAIN_GETPROP */ :
		if this_.transpilePrivateClassFields && n.IsPrivateIdentifier() {
			this_.priv_Es6NormalizeClasses_rewritePrivateGetProp_2510(t, n)
		}
	case int32(105) /* Token.CLASS */ :
		var currClassRecord *Es6NormalizeClasses_ClassRecord = this_.classStack.RemoveFirst()
		_ = currClassRecord
		if !(currClassRecord.classNode == n) {
			j2gort.CheckStateFmt(false, "unexpected node: %s", j2gort.BoxNode(n))
		}
		if currClassRecord.normalizedNames.shouldExtractExtends {
			if !(this_.priv_Es6NormalizeClasses_canExtractExtendsSimply_2500(n)) {
				this_.priv_Es6NormalizeClasses_pushDownNodeIntoNewIife_2502(n)
				parent = n.GetParent()
			}
			this_.priv_Es6NormalizeClasses_extractExtends_2501(t, currClassRecord)
		}
		if currClassRecord.normalizedNames.shouldExtractClass {
			if this_.expressionDecomposer.CanExposeExpression(n) != StaticGet_ExpressionDecomposer_DecompositionType_MOVABLE() {
				this_.priv_Es6NormalizeClasses_pushDownNodeIntoNewIife_2502(n)
			}
			this_.priv_Es6NormalizeClasses_extractClass_2504(t, currClassRecord)
			parent = n.GetParent()
		}
		var classNameNode pkg_rhino.Node = NodeUtil_GetNameNode(n)
		_ = classNameNode
		if !(classNameNode != 0) {
			j2gort.CheckStateFmt(false, "Class missing a name: %s", j2gort.BoxNode(n))
		}
		var classInsertionPoint pkg_rhino.Node = this_.priv_Es6NormalizeClasses_getStatementDeclaringClass_2506(n, classNameNode)
		_ = classInsertionPoint
		if !(classInsertionPoint != 0) {
			j2gort.CheckStateFmt(false, "Class was not extracted: %s", j2gort.BoxNode(n))
		}
		if currClassRecord.bindingIdentifier.IsPresent() && (parent.IsName() || parent.IsAssign()) {
			n.GetFirstChild().ReplaceWith(pkg_rhino.IR_Empty().Srcref(n.GetFirstChild()))
			this_.compiler.ReportChangeToEnclosingScope(n)
		}
		if this_.transpilePrivateClassFields && (currClassRecord.hasInstancePrivateMembers || currClassRecord.hasStaticPrivateMembers) {
			this_.priv_Es6NormalizeClasses_insertPrivateDeclarations_2507(t, currClassRecord, classInsertionPoint)
		}
		this_.priv_Es6NormalizeClasses_rewriteSideEffectedComputedProps_2518(t, currClassRecord, classInsertionPoint)
		if this_.transpileClassFields || (this_.transpilePrivateClassFields && currClassRecord.hasInstancePrivateMembers) {
			this_.priv_Es6NormalizeClasses_rewriteInstanceMembers_2520(t, currClassRecord)
		}
		this_.priv_Es6NormalizeClasses_rewriteStaticMembers_2532(t, currClassRecord, classInsertionPoint)
	case int32(32) /* Token.NAME */ :
		this_.priv_Es6NormalizeClasses_maybeUpdateClassSelfRef_2534(t, n)
	case int32(37) /* Token.THIS */, int32(110) /* Token.SUPER */ :
		this_.priv_Es6NormalizeClasses_visitThisAndSuper_2535(t, n)
	default:
	}
}

// Java method: Es6NormalizeClasses.insertPrivateDeclarations(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.Es6NormalizeClasses.ClassRecord record, com.google.javascript.rhino.Node classInsertionPoint)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_insertPrivateDeclarations_2507(t *NodeTraversal, record *Es6NormalizeClasses_ClassRecord, classInsertionPoint pkg_rhino.Node) {
	_ = t
	_ = record
	_ = classInsertionPoint
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_CONST_DECLARATIONS(), this_.compiler)
	if record.hasInstancePrivateMembers {
		this_.compiler.GetRuntimeJsLibManager().InjectLibForField("$jscomp.PrivateMap")
		var mapVar pkg_rhino.Node = this_.astFactory.CreateSingleConstNameDeclaration(record.GetPrivateMapVarName(this_.compiler, t), this_.astFactory.CreateNewNode(this_.astFactory.CreateQNameWithUnknownType_1_RuntimeJsLibManager_JsLibField(this_.jscompPrivateMap), []pkg_rhino.Node{})).SrcrefTreeIfMissing(record.classNode)
		_ = mapVar
		mapVar.InsertBefore(classInsertionPoint)
		if record.hasInstancePrivateMethodsOrAccessors {
			var descriptorsLit pkg_rhino.Node = this_.astFactory.CreateObjectLit_1([]pkg_rhino.Node{})
			_ = descriptorsLit
			for _, pm := range record.privateMembers {
				_ = pm
				if !(pm.isStatic) {
					if ((pm.kind == StaticGet_Es6NormalizeClasses_ClassRecord_PrivateMemberKind_METHOD()) && (pm.methodNode != 0)) && pm.methodNode.HasChildren() {
						this_.priv_Es6NormalizeClasses_rewriteSuperReferencesInPrivateMember_2517(pm.methodNode, record)
						var fnNode pkg_rhino.Node = pm.methodNode.GetFirstChild().Detach()
						_ = fnNode
						pm.methodNode.Detach()
						var valProp pkg_rhino.Node = this_.astFactory.CreateStringKey("value", fnNode)
						_ = valProp
						var descriptorLit pkg_rhino.Node = this_.astFactory.CreateObjectLit_1([]pkg_rhino.Node{valProp})
						_ = descriptorLit
						var memberKey pkg_rhino.Node = this_.astFactory.CreateStringKey(pm.PropertyName(), descriptorLit)
						_ = memberKey
						descriptorsLit.AddChildToBack(memberKey)
					} else if pm.kind == StaticGet_Es6NormalizeClasses_ClassRecord_PrivateMemberKind_ACCESSOR() {
						var descriptorLit_2 pkg_rhino.Node = this_.astFactory.CreateObjectLit_1([]pkg_rhino.Node{})
						_ = descriptorLit_2
						if (pm.getterNode != 0) && pm.getterNode.HasChildren() {
							this_.priv_Es6NormalizeClasses_rewriteSuperReferencesInPrivateMember_2517(pm.getterNode, record)
							var getFn pkg_rhino.Node = pm.getterNode.GetFirstChild().Detach()
							_ = getFn
							pm.getterNode.Detach()
							if this_.priv_Es6NormalizeClasses_rewriteThisToThisDotSelf_2515(getFn) {
								{
									r_4 := record
									r_4.needsInstanceSelf = true
								}
							}
							var getProp pkg_rhino.Node = this_.astFactory.CreateStringKey("get", getFn)
							_ = getProp
							descriptorLit_2.AddChildToBack(getProp)
						}
						if (pm.setterNode != 0) && pm.setterNode.HasChildren() {
							this_.priv_Es6NormalizeClasses_rewriteSuperReferencesInPrivateMember_2517(pm.setterNode, record)
							var setFn pkg_rhino.Node = pm.setterNode.GetFirstChild().Detach()
							_ = setFn
							pm.setterNode.Detach()
							if this_.priv_Es6NormalizeClasses_rewriteThisToThisDotSelf_2515(setFn) {
								{
									r_5 := record
									r_5.needsInstanceSelf = true
								}
							}
							var setProp pkg_rhino.Node = this_.astFactory.CreateStringKey("set", setFn)
							_ = setProp
							descriptorLit_2.AddChildToBack(setProp)
						}
						var memberKey_2 pkg_rhino.Node = this_.astFactory.CreateStringKey(pm.PropertyName(), descriptorLit_2)
						_ = memberKey_2
						descriptorsLit.AddChildToBack(memberKey_2)
					}
				}
			}

			var createProp pkg_rhino.Node = this_.astFactory.CreateGetPropWithUnknownType(this_.astFactory.CreateNameWithUnknownType("Object"), "create")
			_ = createProp
			var protoCall pkg_rhino.Node = this_.astFactory.CreateCallWithUnknownType(createProp, []pkg_rhino.Node{this_.astFactory.CreateNull(), descriptorsLit})
			_ = protoCall
			var protoVar pkg_rhino.Node = this_.astFactory.CreateSingleConstNameDeclaration(record.GetPrivateProtoVarName(this_.compiler, t), protoCall).SrcrefTreeIfMissing(record.classNode)
			_ = protoVar
			protoVar.InsertBefore(classInsertionPoint)
			this_.compiler.ReportChangeToEnclosingScope(protoVar)
			t.ReportCodeChange_0()
		}
	}
	if record.hasStaticPrivateMembers {
		this_.compiler.GetRuntimeJsLibManager().InjectLibForField("$jscomp.PrivateMap")
		var staticMapVar pkg_rhino.Node = this_.astFactory.CreateSingleConstNameDeclaration(record.GetStaticPrivateMapVarName(this_.compiler, t), this_.astFactory.CreateNewNode(this_.astFactory.CreateQNameWithUnknownType_1_RuntimeJsLibManager_JsLibField(this_.jscompPrivateMap), []pkg_rhino.Node{})).SrcrefTreeIfMissing(record.classNode)
		_ = staticMapVar
		staticMapVar.InsertBefore(classInsertionPoint)
	}
}

// Java method: Es6NormalizeClasses.findPrivateMember(java.lang.String name)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_findPrivateMember_2508(name string) *Es6NormalizeClasses_ResolvedPrivateMember {
	_ = name
	{
		cur_6 := this_.classStack.Iter()
		for cur_6.HasNext() {

			record := cur_6.Next()
			_ = record
			if record.HasPrivateMember(name) {
				return NewEs6NormalizeClasses_ResolvedPrivateMember(record, record.GetPrivateMember(name))
			}
		}
	}
	return nil
}

// Java method: Es6NormalizeClasses.rewritePrivateIn(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node inNode)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_rewritePrivateIn_2509(t *NodeTraversal, inNode pkg_rhino.Node) {
	_ = t
	_ = inNode
	if !(inNode.IsIn()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(inNode))
	}
	var left pkg_rhino.Node = inNode.GetFirstChild()
	_ = left
	var name string = left.GetString()
	_ = name
	var resolved *Es6NormalizeClasses_ResolvedPrivateMember = this_.priv_Es6NormalizeClasses_findPrivateMember_2508(name)
	_ = resolved
	if resolved == nil {
		this_.priv_Es6NormalizeClasses_maybeReportPrivatePropertiesCannotConvertYet_2525(left)
		return
	}
	var right pkg_rhino.Node = inNode.GetLastChild().Detach()
	_ = right
	var record *Es6NormalizeClasses_ClassRecord = resolved.ClassRecord()
	_ = record
	var isStatic bool = resolved.PrivateMember().isStatic
	_ = isStatic
	var mapVarName string
	_ = mapVarName
	if isStatic {
		mapVarName = record.GetStaticPrivateMapVarName(this_.compiler, t)
	} else {
		mapVarName = record.GetPrivateMapVarName(this_.compiler, t)
	}
	var mapNameNode pkg_rhino.Node = this_.astFactory.CreateNameWithUnknownType(mapVarName)
	_ = mapNameNode
	var hasCallee pkg_rhino.Node = this_.astFactory.CreateGetPropWithUnknownType(mapNameNode, "has")
	_ = hasCallee
	var hasCall pkg_rhino.Node = this_.astFactory.CreateCallWithUnknownType(hasCallee, []pkg_rhino.Node{right})
	_ = hasCall
	inNode.ReplaceWith(hasCall.SrcrefTree(inNode))
	t.ReportCodeChange_0()
}

// Java method: Es6NormalizeClasses.rewritePrivateGetProp(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node getPropNode)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_rewritePrivateGetProp_2510(t *NodeTraversal, getPropNode pkg_rhino.Node) {
	_ = t
	_ = getPropNode
	var name string = getPropNode.GetString()
	_ = name
	var resolved *Es6NormalizeClasses_ResolvedPrivateMember = this_.priv_Es6NormalizeClasses_findPrivateMember_2508(name)
	_ = resolved
	if resolved == nil {
		this_.priv_Es6NormalizeClasses_maybeReportPrivatePropertiesCannotConvertYet_2525(getPropNode)
		return
	}
	var record *Es6NormalizeClasses_ClassRecord = resolved.ClassRecord()
	_ = record
	var pm *Es6NormalizeClasses_ClassRecord_PrivateMember = resolved.PrivateMember()
	_ = pm
	if NodeUtil_IsLValue(getPropNode) && (pm.kind == StaticGet_Es6NormalizeClasses_ClassRecord_PrivateMemberKind_METHOD()) {
		this_.compiler.Report_1(JSError_Make_3(getPropNode, StaticGet_Es6NormalizeClasses_ILLEGAL_PRIVATE_MEMBER_ASSIGNMENT(), []string{pm.PropertyName()}))
		return
	}
	var parent pkg_rhino.Node = getPropNode.GetParent()
	_ = parent
	this_.priv_Es6NormalizeClasses_rewritePrivateMemberAccess_2511(t, getPropNode, record, pm, parent)
}

// Java method: Es6NormalizeClasses.rewritePrivateMemberAccess(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node getPropNode, com.google.javascript.jscomp.Es6NormalizeClasses.ClassRecord record, com.google.javascript.jscomp.Es6NormalizeClasses.ClassRecord.PrivateMember pm, com.google.javascript.rhino.Node parent)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_rewritePrivateMemberAccess_2511(t *NodeTraversal, getPropNode pkg_rhino.Node, record *Es6NormalizeClasses_ClassRecord, pm *Es6NormalizeClasses_ClassRecord_PrivateMember, parent pkg_rhino.Node) {
	_ = t
	_ = getPropNode
	_ = record
	_ = pm
	_ = parent
	var isStatic bool = pm.isStatic
	_ = isStatic
	var mapVarName string
	_ = mapVarName
	if isStatic {
		mapVarName = record.GetStaticPrivateMapVarName(this_.compiler, t)
	} else {
		mapVarName = record.GetPrivateMapVarName(this_.compiler, t)
	}
	var memberPropName string = pm.PropertyName()
	_ = memberPropName
	var receiver pkg_rhino.Node = getPropNode.GetFirstChild().Detach()
	_ = receiver
	var mapNameNode pkg_rhino.Node = this_.astFactory.CreateNameWithUnknownType(mapVarName)
	_ = mapNameNode
	var getCallee pkg_rhino.Node = this_.astFactory.CreateGetPropWithUnknownType(mapNameNode, "get")
	_ = getCallee
	var isOptChain bool = getPropNode.IsOptChainGetProp()
	_ = isOptChain
	var endOfSegment pkg_rhino.Node
	_ = endOfSegment
	if isOptChain {
		endOfSegment = NodeUtil_GetEndOfOptChainSegment(getPropNode)
	} else {
		endOfSegment = pkg_rhino.Node(0)
	}
	var isDirectCall bool = ((parent.IsCall() || parent.IsOptChainCall()) && (parent.GetFirstChild() == getPropNode))
	_ = isDirectCall
	if isDirectCall {
		var isOptionalCall bool = (parent.IsOptChainCall() && parent.IsOptionalChainStart())
		_ = isOptionalCall
		parent.SetIsOptionalChainStart(false)
		if this_.priv_Es6NormalizeClasses_receiverNeedsTempVar_2537(receiver) {
			var tempName string = this_.priv_Es6NormalizeClasses_generateUniqueTmpVarName_2494(t)
			_ = tempName
			var insertionBlock pkg_rhino.Node = this_.priv_Es6NormalizeClasses_getDeclarationInsertionPoint_2538(getPropNode, t)
			_ = insertionBlock
			var decl pkg_rhino.Node = this_.astFactory.CreateSingleLetNameDeclaration_1(tempName)
			_ = decl
			insertionBlock.AddChildToFront(decl.SrcrefTree(getPropNode))
			NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_LET_DECLARATIONS(), this_.compiler)
			var tempNameNode1 pkg_rhino.Node = this_.astFactory.CreateNameWithUnknownType(tempName).Srcref(getPropNode)
			_ = tempNameNode1
			var assign pkg_rhino.Node = this_.astFactory.CreateAssign_2_Node_Node(tempNameNode1, receiver).Srcref(getPropNode)
			_ = assign
			receiver = this_.astFactory.CreateNameWithUnknownType(tempName).Srcref(getPropNode)
			var tempNameNode2 pkg_rhino.Node = this_.astFactory.CreateNameWithUnknownType(tempName).Srcref(getPropNode)
			_ = tempNameNode2
			var tempNameNode3 pkg_rhino.Node = this_.astFactory.CreateNameWithUnknownType(tempName).Srcref(getPropNode)
			_ = tempNameNode3
			var getCall pkg_rhino.Node = this_.astFactory.CreateCallWithUnknownType(getCallee, []pkg_rhino.Node{tempNameNode2}).Srcref(getPropNode)
			_ = getCall
			var methodGetProp pkg_rhino.Node = this_.astFactory.CreateGetPropWithUnknownType(getCall, memberPropName)
			_ = methodGetProp
			var callCallee pkg_rhino.Node = this_.priv_Es6NormalizeClasses_createCallCallee_2514(methodGetProp, isOptionalCall, getPropNode)
			_ = callCallee
			parent.SetToken(func() pkg_rhino.Token {
				if callCallee.IsOptChainGetProp() {
					return pkg_rhino.StaticGet_Token_OPTCHAIN_CALL()
				}
				return pkg_rhino.StaticGet_Token_CALL()
			}())
			callCallee.SrcrefTreeIfMissing(getPropNode)
			tempNameNode3.SrcrefTree(getPropNode)
			getPropNode.Detach()
			parent.AddChildToFront(callCallee)
			tempNameNode3.InsertAfter(callCallee)
			var resultExpr pkg_rhino.Node = this_.priv_Es6NormalizeClasses_wrapOptionalChainStartIfNecessary_2513(receiver, parent, endOfSegment, isOptChain)
			_ = resultExpr
			var placeholder pkg_rhino.Node = pkg_rhino.IR_Empty()
			_ = placeholder
			resultExpr.ReplaceWith(placeholder)
			var comma pkg_rhino.Node = this_.astFactory.CreateComma(assign, resultExpr).Srcref(getPropNode)
			_ = comma
			placeholder.ReplaceWith(comma)
			t.ReportCodeChange_0()
		} else {
			var getCall_2 pkg_rhino.Node = this_.astFactory.CreateCallWithUnknownType(getCallee, []pkg_rhino.Node{receiver})
			_ = getCall_2
			var methodGetProp_2 pkg_rhino.Node = this_.astFactory.CreateGetPropWithUnknownType(getCall_2, memberPropName)
			_ = methodGetProp_2
			var receiverArg pkg_rhino.Node = receiver.CloneTree_0()
			_ = receiverArg
			receiverArg.SrcrefTree(getPropNode)
			getPropNode.Detach()
			var callCallee_2 pkg_rhino.Node = this_.priv_Es6NormalizeClasses_createCallCallee_2514(methodGetProp_2, isOptionalCall, getPropNode)
			_ = callCallee_2
			parent.SetToken(func() pkg_rhino.Token {
				if callCallee_2.IsOptChainGetProp() {
					return pkg_rhino.StaticGet_Token_OPTCHAIN_CALL()
				}
				return pkg_rhino.StaticGet_Token_CALL()
			}())
			callCallee_2.SrcrefTreeIfMissing(getPropNode)
			parent.AddChildToFront(callCallee_2)
			receiverArg.InsertAfter(callCallee_2)
			this_.priv_Es6NormalizeClasses_wrapOptionalChainStartIfNecessary_2513(receiver, parent, endOfSegment, isOptChain)
			t.ReportCodeChange_0()
		}
	} else {
		var assignComma pkg_rhino.Node = this_.priv_Es6NormalizeClasses_extractReceiverTmpIfNecessary_2512(t, receiver, getPropNode, isOptChain)
		_ = assignComma
		if assignComma != 0 {
			receiver = assignComma.GetFirstChild().CloneTree_0()
		}
		var getCall_3 pkg_rhino.Node = this_.astFactory.CreateCallWithUnknownType(getCallee, []pkg_rhino.Node{receiver})
		_ = getCall_3
		var methodGetProp_3 pkg_rhino.Node = this_.astFactory.CreateGetPropWithUnknownType(getCall_3, memberPropName)
		_ = methodGetProp_3
		getPropNode.ReplaceWith(methodGetProp_3.SrcrefTree(getPropNode))
		var result pkg_rhino.Node = this_.priv_Es6NormalizeClasses_wrapOptionalChainStartIfNecessary_2513(receiver, methodGetProp_3, endOfSegment, isOptChain)
		_ = result
		if assignComma != 0 {
			var placeholder_2 pkg_rhino.Node = pkg_rhino.IR_Empty()
			_ = placeholder_2
			result.ReplaceWith(placeholder_2)
			var comma_2 pkg_rhino.Node = this_.astFactory.CreateComma(assignComma, result).Srcref(getPropNode)
			_ = comma_2
			placeholder_2.ReplaceWith(comma_2)
		}
		t.ReportCodeChange_0()
	}
}

// Java method: Es6NormalizeClasses.extractReceiverTmpIfNecessary(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node receiver, com.google.javascript.rhino.Node getPropNode, boolean isOptChain)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_extractReceiverTmpIfNecessary_2512(t *NodeTraversal, receiver pkg_rhino.Node, getPropNode pkg_rhino.Node, isOptChain bool) pkg_rhino.Node {
	_ = t
	_ = receiver
	_ = getPropNode
	_ = isOptChain
	if !(isOptChain) || !(this_.priv_Es6NormalizeClasses_receiverNeedsTempVar_2537(receiver)) {
		return pkg_rhino.Node(0)
	}
	var tempName string = this_.priv_Es6NormalizeClasses_generateUniqueTmpVarName_2494(t)
	_ = tempName
	var insertionBlock pkg_rhino.Node = this_.priv_Es6NormalizeClasses_getDeclarationInsertionPoint_2538(getPropNode, t)
	_ = insertionBlock
	var decl pkg_rhino.Node = this_.astFactory.CreateSingleLetNameDeclaration_1(tempName)
	_ = decl
	insertionBlock.AddChildToFront(decl.SrcrefTree(getPropNode))
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_LET_DECLARATIONS(), this_.compiler)
	var tempNameNode pkg_rhino.Node = this_.astFactory.CreateNameWithUnknownType(tempName).Srcref(getPropNode)
	_ = tempNameNode
	return this_.astFactory.CreateAssign_2_Node_Node(tempNameNode, receiver).Srcref(getPropNode)
}

// Java method: Es6NormalizeClasses.wrapOptionalChainStartIfNecessary(com.google.javascript.rhino.Node receiver, com.google.javascript.rhino.Node accessExpr, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node endOfSegment, boolean isOptChain)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_wrapOptionalChainStartIfNecessary_2513(receiver pkg_rhino.Node, accessExpr pkg_rhino.Node, endOfSegment pkg_rhino.Node, isOptChain bool) pkg_rhino.Node {
	_ = receiver
	_ = accessExpr
	_ = endOfSegment
	_ = isOptChain
	if !(isOptChain) {
		return accessExpr
	}
	var segmentRoot pkg_rhino.Node
	_ = segmentRoot
	if (endOfSegment != 0) && (endOfSegment.GetParent() != 0) {
		segmentRoot = endOfSegment
	} else {
		segmentRoot = accessExpr
	}
	if !(segmentRoot == accessExpr) {
		NodeUtil_ConvertToNonOptionalChainSegmentDownTo(segmentRoot, accessExpr)
	}
	var placeholder pkg_rhino.Node = pkg_rhino.IR_Empty()
	_ = placeholder
	if segmentRoot.GetParent() != 0 {
		segmentRoot.ReplaceWith(placeholder)
	}
	var hook pkg_rhino.Node = this_.astFactory.CreateOptionalChainShortCircuit(receiver.CloneTree_0(), segmentRoot).SrcrefTreeIfMissing(segmentRoot)
	_ = hook
	if placeholder.GetParent() != 0 {
		placeholder.ReplaceWith(hook)
	}
	var hookParent pkg_rhino.Node = hook.GetParent()
	_ = hookParent
	if (hookParent != 0) && (hookParent.IsCall() || hookParent.IsOptChainCall()) {
		hook.SetIsParenthesized(true)
		if hookParent.GetFirstChild() == hook {
			hookParent.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
		}
	}
	return hook
}

// Java method: Es6NormalizeClasses.createCallCallee(com.google.javascript.rhino.Node methodGetProp, boolean isOptionalCall, com.google.javascript.rhino.Node getPropNode)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_createCallCallee_2514(methodGetProp pkg_rhino.Node, isOptionalCall bool, getPropNode pkg_rhino.Node) pkg_rhino.Node {
	_ = methodGetProp
	_ = isOptionalCall
	_ = getPropNode
	if isOptionalCall {
		return this_.astFactory.CreateStartOptChainGetprop(methodGetProp, "call", AstFactory_S_Type_1_Node(getPropNode))
	} else {
		return this_.astFactory.CreateGetPropWithUnknownType(methodGetProp, "call")
	}
}

// Java method: Es6NormalizeClasses.rewriteThisToThisDotSelf(com.google.javascript.rhino.Node fn)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_rewriteThisToThisDotSelf_2515(fn pkg_rhino.Node) bool {
	_ = fn
	var fnBody pkg_rhino.Node = NodeUtil_GetFunctionBody(fn)
	_ = fnBody
	if fnBody == 0 {
		return false
	}
	var hasThis []bool = make([]bool, int32(1))
	_ = hasThis
	NodeUtil_VisitPostOrder_3(fnBody, NewLambda_NodeUtil_Visitor(func(n pkg_rhino.Node) {
		_ = n
		if n.IsThis() {
			{
				a_8 := hasThis
				i_9 := int32(0)
				a_8[i_9] = true
			}
			var target pkg_rhino.Node = this_.astFactory.CreateGetPropWithUnknownType(this_.astFactory.CreateThis(AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN())), "$self").Srcref(n)
			_ = target
			n.ReplaceWith(target)
			this_.compiler.ReportChangeToEnclosingScope(target)
		}
	}), lambdaSingleton_es6_normalize_classes_1)
	return hasThis[int32(0)]
}

// Java method: Es6NormalizeClasses.hasThisReference(com.google.javascript.rhino.Node fn)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_hasThisReference_2516(fn pkg_rhino.Node) bool {
	_ = fn
	var fnBody pkg_rhino.Node = NodeUtil_GetFunctionBody(fn)
	_ = fnBody
	if fnBody == 0 {
		return false
	}
	var hasThis []bool = make([]bool, int32(1))
	_ = hasThis
	NodeUtil_VisitPostOrder_3(fnBody, NewLambda_NodeUtil_Visitor(func(n pkg_rhino.Node) {
		_ = n
		if n.IsThis() {
			{
				a_11 := hasThis
				i_12 := int32(0)
				a_11[i_12] = true
			}
		}
	}), lambdaSingleton_es6_normalize_classes_2)
	return hasThis[int32(0)]
}

// Java method: Es6NormalizeClasses.rewriteSuperReferencesInPrivateMember(com.google.javascript.rhino.Node memberNode, com.google.javascript.jscomp.Es6NormalizeClasses.ClassRecord record)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_rewriteSuperReferencesInPrivateMember_2517(memberNode pkg_rhino.Node, record *Es6NormalizeClasses_ClassRecord) {
	_ = memberNode
	_ = record
	if record.normalizedNames.SuperClassNode().IsEmpty() {
		return
	}
	var fn pkg_rhino.Node = memberNode.GetFirstChild()
	_ = fn
	var fnBody pkg_rhino.Node = NodeUtil_GetFunctionBody(fn)
	_ = fnBody
	if (fnBody == 0) || !(NodeUtil_ReferencesSuper(fnBody)) {
		return
	}
	var isStatic bool = memberNode.IsStaticMember()
	_ = isStatic
	NodeUtil_VisitPostOrder_3(fnBody, NewLambda_NodeUtil_Visitor(func(n pkg_rhino.Node) {
		_ = n
		if n.IsSuper() {
			var parent pkg_rhino.Node = n.GetParent()
			_ = parent
			if parent.IsGetProp() || parent.IsGetElem() {
				var target pkg_rhino.Node = pkg_rhino.Node(0)
				_ = target
				if isStatic {
					target = record.CreateNewSuperclassReferenceNode().SrcrefTree(n)
				} else {
					target = this_.astFactory.CreatePrototypeAccess(record.CreateNewSuperclassReferenceNode()).SrcrefTree(n)
				}
				var grandParent pkg_rhino.Node = parent.GetParent()
				_ = grandParent
				if ((grandParent != 0) && grandParent.IsCall()) && (grandParent.GetFirstChild() == parent) {
					n.ReplaceWith(target)
					if !(isStatic) || !(this_.compiler.GetOptions().GetAssumeStaticInheritanceIsNotUsed()) {
						var callCallee pkg_rhino.Node = this_.astFactory.CreateGetPropWithUnknownType(parent.Detach(), "call")
						_ = callCallee
						grandParent.AddChildToFront(callCallee.SrcrefTreeIfMissing(parent))
						var thisArg pkg_rhino.Node = this_.astFactory.CreateThis(AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN())).Srcref(n)
						_ = thisArg
						thisArg.InsertAfter(callCallee)
					}
				} else if ((grandParent != 0) && grandParent.IsAssign()) && (grandParent.GetFirstChild() == parent) {
					var rhs pkg_rhino.Node = grandParent.GetSecondChild().Detach()
					_ = rhs
					var propKey pkg_rhino.Node
					_ = propKey
					if parent.IsGetProp() {
						propKey = this_.astFactory.CreateCall(this_.astFactory.CreateName_2_string_AstFactory_Type("JSCompiler_renameProperty", AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN())), AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_STRING()), []pkg_rhino.Node{this_.astFactory.CreateString(parent.GetString()), record.CreateNewSuperclassReferenceNode().SrcrefTree(n)})
					} else {
						propKey = parent.GetSecondChild().Detach()
					}
					var reflectSet pkg_rhino.Node = this_.astFactory.CreateCall(this_.astFactory.CreateGetPropWithUnknownType(this_.astFactory.CreateName_2_string_AstFactory_Type("Reflect", AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN())), "set"), AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN()), []pkg_rhino.Node{target, propKey, rhs, this_.astFactory.CreateThis(AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN())).Srcref(n)})
					_ = reflectSet
					grandParent.ReplaceWith(reflectSet.SrcrefTreeIfMissing(grandParent))
				} else {
					var propKey_2 pkg_rhino.Node
					_ = propKey_2
					if parent.IsGetProp() {
						propKey_2 = this_.astFactory.CreateCall(this_.astFactory.CreateName_2_string_AstFactory_Type("JSCompiler_renameProperty", AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN())), AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_STRING()), []pkg_rhino.Node{this_.astFactory.CreateString(parent.GetString()), record.CreateNewSuperclassReferenceNode().SrcrefTree(n)})
					} else {
						propKey_2 = parent.GetSecondChild().Detach()
					}
					var reflectGet pkg_rhino.Node = this_.astFactory.CreateCall(this_.astFactory.CreateGetPropWithUnknownType(this_.astFactory.CreateName_2_string_AstFactory_Type("Reflect", AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN())), "get"), AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN()), []pkg_rhino.Node{target, propKey_2, this_.astFactory.CreateThis(AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN())).Srcref(n)})
					_ = reflectGet
					parent.ReplaceWith(reflectGet.SrcrefTreeIfMissing(parent))
				}
				this_.compiler.ReportChangeToEnclosingScope(fnBody)
			}
		}
	}), lambdaSingleton_es6_normalize_classes_3)
}

// Java method: Es6NormalizeClasses.rewriteSideEffectedComputedProps(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.Es6NormalizeClasses.ClassRecord record, com.google.javascript.rhino.Node insertionPoint)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_rewriteSideEffectedComputedProps_2518(t *NodeTraversal, record *Es6NormalizeClasses_ClassRecord, insertionPoint pkg_rhino.Node) {
	_ = t
	_ = record
	_ = insertionPoint
	var computedPropsWithSideEffects *j2gort.Deque[pkg_rhino.Node] = record.computedPropsWithSideEffects
	_ = computedPropsWithSideEffects
	if computedPropsWithSideEffects.IsEmpty() {
		return
	}
	for !(computedPropsWithSideEffects.IsEmpty()) {
		var computedPropMember pkg_rhino.Node = computedPropsWithSideEffects.Pop()
		_ = computedPropMember
		this_.priv_Es6NormalizeClasses_extractExpressionFromCompField_2519(t, record, computedPropMember, insertionPoint)
	}
	t.ReportCodeChange_0()
}

// Java method: Es6NormalizeClasses.extractExpressionFromCompField(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.Es6NormalizeClasses.ClassRecord record, com.google.javascript.rhino.Node memberField, com.google.javascript.rhino.Node insertionPoint)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_extractExpressionFromCompField_2519(t *NodeTraversal, record *Es6NormalizeClasses_ClassRecord, memberField pkg_rhino.Node, insertionPoint pkg_rhino.Node) {
	_ = t
	_ = record
	_ = memberField
	_ = insertionPoint
	if !(memberField.IsComputedFieldDef() || memberField.IsComputedProp()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(memberField))
	}
	var compExpression pkg_rhino.Node = memberField.RemoveFirstChild()
	_ = compExpression
	var compFieldVar pkg_rhino.Node = this_.astFactory.CreateSingleVarNameDeclaration_2(this_.priv_Es6NormalizeClasses_generateUniqueCompFieldVarName_2497(t), compExpression).SrcrefTreeIfMissing(record.classNode)
	_ = compFieldVar
	var compFieldName pkg_rhino.Node = compFieldVar.GetFirstChild()
	_ = compFieldName
	memberField.AddChildToFront(compFieldName.CloneNode_0())
	compFieldVar.InsertBefore(insertionPoint)
	compFieldVar.SrcrefTreeIfMissing(record.classNode)
}

// Java method: Es6NormalizeClasses.rewriteInstanceMembers(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.Es6NormalizeClasses.ClassRecord record)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_rewriteInstanceMembers_2520(t *NodeTraversal, record *Es6NormalizeClasses_ClassRecord) {
	_ = t
	_ = record
	var instanceMembers *j2gort.Deque[pkg_rhino.Node] = record.instanceMembers
	_ = instanceMembers
	var hasPriv bool = (this_.transpilePrivateClassFields && record.hasInstancePrivateMembers)
	_ = hasPriv
	if instanceMembers.IsEmpty() && !(hasPriv) {
		return
	}
	this_.ctorCreator.SynthesizeClassConstructorIfMissing(t, record.classNode)
	var ctor pkg_rhino.Node = NodeUtil_GetEs6ClassConstructorMemberFunctionDef(record.classNode)
	_ = ctor
	var ctorFn pkg_rhino.Node = ctor.GetFirstChild()
	_ = ctorFn
	var ctorBlock pkg_rhino.Node = ctorFn.GetLastChild()
	_ = ctorBlock
	var isDerived bool = !(record.classNode.GetSecondChild().IsEmpty())
	_ = isDerived
	var insertionPoint pkg_rhino.Node = this_.priv_Es6NormalizeClasses_addTemporaryInsertionPoint_2523(ctorBlock, record)
	_ = insertionPoint
	if hasPriv {
		var mapName string = record.GetPrivateMapVarName(this_.compiler, t)
		_ = mapName
		var privName string = record.GetPrivateVarName(this_.compiler, t)
		_ = privName
		var privInit pkg_rhino.Node = pkg_rhino.Node(0)
		_ = privInit
		if record.hasInstancePrivateMethodsOrAccessors {
			var protoName pkg_rhino.Node = this_.astFactory.CreateNameWithUnknownType(record.GetPrivateProtoVarName(this_.compiler, t))
			_ = protoName
			var createProp pkg_rhino.Node = this_.astFactory.CreateGetPropWithUnknownType(this_.astFactory.CreateNameWithUnknownType("Object"), "create")
			_ = createProp
			privInit = this_.astFactory.CreateCallWithUnknownType(createProp, []pkg_rhino.Node{protoName})
		} else {
			var createProp_2 pkg_rhino.Node = this_.astFactory.CreateGetPropWithUnknownType(this_.astFactory.CreateNameWithUnknownType("Object"), "create")
			_ = createProp_2
			privInit = this_.astFactory.CreateCallWithUnknownType(createProp_2, []pkg_rhino.Node{this_.astFactory.CreateNull()})
		}
		var privDecl pkg_rhino.Node = this_.astFactory.CreateSingleConstNameDeclaration(privName, privInit).SrcrefTreeIfMissing(record.classNode)
		_ = privDecl
		privDecl.InsertBefore(insertionPoint)
		if record.needsInstanceSelf {
			var privSelf pkg_rhino.Node = this_.astFactory.CreateGetPropWithUnknownType(this_.astFactory.CreateConstantName(privName, AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN())), "$self")
			_ = privSelf
			var assignSelf pkg_rhino.Node = this_.astFactory.CreateAssignStatement(privSelf, this_.astFactory.CreateThisForEs6Class(record.classNode)).SrcrefTreeIfMissing(record.classNode)
			_ = assignSelf
			assignSelf.InsertBefore(insertionPoint)
		}
		var mapNameNode pkg_rhino.Node = this_.astFactory.CreateNameWithUnknownType(mapName)
		_ = mapNameNode
		var mapSetCallee pkg_rhino.Node = this_.astFactory.CreateGetPropWithUnknownType(mapNameNode, "set")
		_ = mapSetCallee
		var mapSetCall pkg_rhino.Node = this_.astFactory.CreateCallWithUnknownType(mapSetCallee, []pkg_rhino.Node{this_.astFactory.CreateThisForEs6Class(record.classNode), this_.astFactory.CreateConstantName(privName, AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN()))})
		_ = mapSetCall
		var mapSetStmt pkg_rhino.Node = this_.astFactory.ExprResult(mapSetCall).SrcrefTreeIfMissing(record.classNode)
		_ = mapSetStmt
		mapSetStmt.InsertBefore(insertionPoint)
	}
	for !(instanceMembers.IsEmpty()) {
		var instanceMember pkg_rhino.Node = instanceMembers.Pop()
		_ = instanceMember
		if !(instanceMember.IsMemberFieldDef() || instanceMember.IsComputedFieldDef()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(instanceMember))
		}
		if instanceMember.IsPrivateIdentifier() {
			if this_.transpilePrivateClassFields {
				var privName_2 string = record.GetPrivateVarName(this_.compiler, t)
				_ = privName_2
				var propName string = j2gort.JStrSubFrom(instanceMember.GetString(), int32(1))
				_ = propName
				var lhs pkg_rhino.Node = this_.astFactory.CreateGetPropWithUnknownType(this_.astFactory.CreateConstantName(privName_2, AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN())), propName)
				_ = lhs
				var fieldValue pkg_rhino.Node
				_ = fieldValue
				if instanceMember.HasChildren() {
					fieldValue = instanceMember.RemoveFirstChild()
				} else {
					fieldValue = this_.astFactory.CreateUndefinedValue()
				}
				var assign pkg_rhino.Node = this_.astFactory.CreateAssignStatement(lhs, fieldValue)
				_ = assign
				assign.SrcrefTreeIfMissing(instanceMember)
				assign.InsertBefore(insertionPoint)
				instanceMember.Detach()
			}
		} else {
			var thisNode pkg_rhino.Node = this_.astFactory.CreateThisForEs6ClassMember(instanceMember)
			_ = thisNode
			var transpiledNode pkg_rhino.Node
			_ = transpiledNode
			if instanceMember.IsMemberFieldDef() {
				transpiledNode = this_.priv_Es6NormalizeClasses_convNonCompFieldToGetProp_2527(thisNode, instanceMember, false)
			} else {
				transpiledNode = this_.priv_Es6NormalizeClasses_convCompFieldToGetElem_2528(thisNode, instanceMember, false)
			}
			if transpiledNode != 0 {
				transpiledNode.InsertBefore(insertionPoint)
			}
		}
	}
	if this_.transpilePrivateClassFields {
		for _, pm := range record.privateMembers {
			_ = pm
			if !(pm.isStatic) {
				{
					cur_17 := j2gort.NewIterCursor(j2gort.CheckNotNull(pm.declarationNodes))
					for cur_17.HasNext() {
						declNode := pkg_rhino.Node(cur_17.NextNode())
						_ = declNode
						if declNode.GetParent() != 0 {
							NodeUtil_DeleteNode(declNode, this_.compiler)
						}
					}
				}
			}
		}

	}
	if !(isDerived) && NodeUtil_ReferencesEnclosingReceiver(ctorFn.GetSecondChild()) {
		this_.priv_Es6NormalizeClasses_moveConstructorParamDefaultsAfterFields_2521(t, ctorFn, insertionPoint)
	}
	insertionPoint.Detach()
	t.ReportCodeChange_0()
	t.ReportCodeChange_1(ctorBlock)
}

// Java method: Es6NormalizeClasses.moveConstructorParamDefaultsAfterFields(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node ctorFn, com.google.javascript.rhino.Node insertionPoint)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_moveConstructorParamDefaultsAfterFields_2521(t *NodeTraversal, ctorFn pkg_rhino.Node, insertionPoint pkg_rhino.Node) {
	_ = t
	_ = ctorFn
	_ = insertionPoint
	var paramList pkg_rhino.Node = ctorFn.GetSecondChild()
	_ = paramList
	{
		var param pkg_rhino.Node = paramList.GetFirstChild()
		_ = param
		for param != 0 {
			var next pkg_rhino.Node = param.GetNext()
			_ = next
			if param.IsDefaultValue() {
				var lhs pkg_rhino.Node = param.RemoveFirstChild()
				_ = lhs
				var jsDoc *pkg_rhino.JSDocInfo = lhs.GetJSDocInfo()
				_ = jsDoc
				lhs.SetJSDocInfo(nil)
				var defaultValue pkg_rhino.Node = param.RemoveFirstChild()
				_ = defaultValue
				if lhs.IsName() {
					var rhs pkg_rhino.Node = this_.priv_Es6NormalizeClasses_createDefaultValueHook_2522(lhs, defaultValue)
					_ = rhs
					var assignStmt pkg_rhino.Node = this_.astFactory.CreateAssignStatement(lhs.CloneNode_0(), rhs).SrcrefTreeIfMissing(param)
					_ = assignStmt
					assignStmt.InsertBefore(insertionPoint)
					param.ReplaceWith(lhs)
					lhs.SetJSDocInfo(jsDoc)
				} else {
					if !(lhs.IsDestructuringPattern()) {
						j2gort.CheckStateMsg(false, j2gort.BoxNode(lhs))
					}
					var tempParam pkg_rhino.Node = this_.astFactory.CreateName_2_string_AstFactory_Type(this_.priv_Es6NormalizeClasses_generateUniqueTmpVarName_2494(t), AstFactory_S_Type_1_Node(lhs))
					_ = tempParam
					var rhs_2 pkg_rhino.Node = this_.priv_Es6NormalizeClasses_createDefaultValueHook_2522(tempParam, defaultValue)
					_ = rhs_2
					var letStmt pkg_rhino.Node = pkg_rhino.IR_Let_2(lhs, rhs_2).SrcrefTreeIfMissing(param)
					_ = letStmt
					letStmt.InsertBefore(insertionPoint)
					NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_LET_DECLARATIONS(), this_.compiler)
					param.ReplaceWith(tempParam)
					tempParam.SetJSDocInfo(jsDoc)
				}
			} else if param.IsDestructuringPattern() || (param.IsRest() && param.GetOnlyChild().IsDestructuringPattern()) {
				var pattern pkg_rhino.Node
				_ = pattern
				if param.IsRest() {
					pattern = param.GetOnlyChild()
				} else {
					pattern = param
				}
				var jsDoc_2 *pkg_rhino.JSDocInfo = pattern.GetJSDocInfo()
				_ = jsDoc_2
				pattern.SetJSDocInfo(nil)
				var tempParam_2 pkg_rhino.Node = this_.astFactory.CreateName_2_string_AstFactory_Type(this_.priv_Es6NormalizeClasses_generateUniqueTmpVarName_2494(t), AstFactory_S_Type_1_Node(pattern))
				_ = tempParam_2
				pattern.ReplaceWith(tempParam_2)
				tempParam_2.SetJSDocInfo(jsDoc_2)
				var letStmt_2 pkg_rhino.Node = pkg_rhino.IR_Let_2(pattern, tempParam_2.CloneNode_0()).SrcrefTreeIfMissing(param)
				_ = letStmt_2
				letStmt_2.InsertBefore(insertionPoint)
				NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_LET_DECLARATIONS(), this_.compiler)
			}
			param = next
		}
	}
}

// Java method: Es6NormalizeClasses.createDefaultValueHook(com.google.javascript.rhino.Node paramName, com.google.javascript.rhino.Node defaultValue)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_createDefaultValueHook_2522(paramName pkg_rhino.Node, defaultValue pkg_rhino.Node) pkg_rhino.Node {
	_ = paramName
	_ = defaultValue
	var undefined pkg_rhino.Node = this_.astFactory.CreateUndefinedValue()
	_ = undefined
	undefined.MakeNonIndexable()
	return this_.astFactory.CreateHook(this_.astFactory.CreateSheq(paramName.CloneNode_0(), undefined), defaultValue, paramName.CloneNode_0())
}

// Java method: Es6NormalizeClasses.addTemporaryInsertionPoint(com.google.javascript.rhino.Node ctorBlock, com.google.javascript.jscomp.Es6NormalizeClasses.ClassRecord classRecord)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_addTemporaryInsertionPoint_2523(ctorBlock pkg_rhino.Node, classRecord *Es6NormalizeClasses_ClassRecord) pkg_rhino.Node {
	_ = ctorBlock
	_ = classRecord
	if !(int32(len(classRecord.superCalls)) <= int32(1)) {
		j2gort.CheckStateFmt(false, "classes with public fields must have only one super() call at the constructor root: %s", j2gort.BoxNode(classRecord.classNode))
	}
	var tempNode pkg_rhino.Node = pkg_rhino.IR_Empty()
	_ = tempNode
	if int32(len(classRecord.superCalls)) == int32(1) {
		var stmt pkg_rhino.Node = this_.priv_Es6NormalizeClasses_ensureSuperCallIsOwnStatement_2524(classRecord.superCalls[int32(0)])
		_ = stmt
		tempNode.InsertAfter(stmt)
		return tempNode
	}
	{
		var stmt_2 pkg_rhino.Node = ctorBlock.GetFirstChild()
		_ = stmt_2
		for ; stmt_2 != 0; stmt_2 = stmt_2.GetNext() {
			if NodeUtil_IsExprCall(stmt_2) && stmt_2.GetFirstFirstChild().IsSuper() {
				tempNode.InsertAfter(stmt_2)
				return tempNode
			}
		}
	}
	var insertionPoint pkg_rhino.Node = NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(ctorBlock)
	_ = insertionPoint
	if insertionPoint != 0 {
		tempNode.InsertBefore(insertionPoint)
	} else {
		ctorBlock.AddChildToBack(tempNode)
	}
	return tempNode
}

// Java method: Es6NormalizeClasses.ensureSuperCallIsOwnStatement(com.google.javascript.rhino.Node superCall)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_ensureSuperCallIsOwnStatement_2524(superCall pkg_rhino.Node) pkg_rhino.Node {
	_ = superCall
	var parent pkg_rhino.Node = superCall.GetParent()
	_ = parent
	if parent.IsExprResult() {
		return parent
	}
	if !(this_.expressionDecomposer.CanExposeExpression(superCall) != StaticGet_ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE()) {
		j2gort.CheckStateFmt(false, "Cannot decompose super() call in a class with class fields. Move super() call to the root of the constructor. %s", j2gort.BoxNode(superCall))
	}
	this_.expressionDecomposer.MaybeExposeExpression(superCall)
	this_.expressionDecomposer.MoveExpression(superCall)
	return NodeUtil_GetEnclosingStatement(superCall)
}

// Java method: Es6NormalizeClasses.maybeReportPrivatePropertiesCannotConvertYet(com.google.javascript.rhino.Node n)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_maybeReportPrivatePropertiesCannotConvertYet_2525(n pkg_rhino.Node) {
	_ = n
	if this_.transpilePrivateClassFields && n.IsPrivateIdentifier() {
		TranspilationUtil_CannotConvertYet(this_.compiler, n, "Private elements")
	}
}

// Java method: Es6NormalizeClasses.createObjectDotDefineProperty()
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_createObjectDotDefineProperty_2526() pkg_rhino.Node {
	this_.compiler.GetRuntimeJsLibManager().InjectLibForField("$jscomp.global")
	var global RuntimeJsLibManager_JsLibField = this_.compiler.GetRuntimeJsLibManager().GetJsLibField("$jscomp.global")
	_ = global
	var jscompDotGlobal pkg_rhino.Node = this_.astFactory.CreateQNameWithUnknownType_1_RuntimeJsLibManager_JsLibField(global)
	_ = jscompDotGlobal
	var objectNode pkg_rhino.Node = this_.astFactory.CreateGetPropWithUnknownType(jscompDotGlobal, "Object")
	_ = objectNode
	return this_.astFactory.CreateGetPropWithUnknownType(objectNode, "defineProperty")
}

// Java method: Es6NormalizeClasses.convNonCompFieldToGetProp(com.google.javascript.rhino.Node receiver, com.google.javascript.rhino.Node noncomputedField, boolean isStatic)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_convNonCompFieldToGetProp_2527(receiver pkg_rhino.Node, noncomputedField pkg_rhino.Node, isStatic bool) pkg_rhino.Node {
	_ = receiver
	_ = noncomputedField
	_ = isStatic
	if !(noncomputedField.IsMemberFieldDef()) {
		j2gort.CheckArgument(false)
	}
	if !(receiver.GetParent() == 0) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(receiver))
	}
	var fieldValue pkg_rhino.Node = pkg_rhino.Node(0)
	_ = fieldValue
	if isStatic || this_.transpileClassFields {
		fieldValue = noncomputedField.GetFirstChild()
	}
	var propName string = noncomputedField.GetString()
	_ = propName
	if propName == "__proto__" {
		return this_.priv_Es6NormalizeClasses_createFieldDefineProperty_2529(noncomputedField, fieldValue, receiver, this_.astFactory.CreateString(propName))
	}
	return this_.priv_Es6NormalizeClasses_createFieldAssignment_2530(noncomputedField, fieldValue, j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return j2gort.BoxNode(this_.astFactory.CreateGetProp(receiver, propName, AstFactory_S_Type_1_Node(noncomputedField)))
	}))
}

// Java method: Es6NormalizeClasses.convCompFieldToGetElem(com.google.javascript.rhino.Node receiver, com.google.javascript.rhino.Node computedField, boolean isStatic)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_convCompFieldToGetElem_2528(receiver pkg_rhino.Node, computedField pkg_rhino.Node, isStatic bool) pkg_rhino.Node {
	_ = receiver
	_ = computedField
	_ = isStatic
	if !(computedField.IsComputedFieldDef()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(computedField))
	}
	if !(receiver.GetParent() == 0) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(receiver))
	}
	var fieldValue pkg_rhino.Node = pkg_rhino.Node(0)
	_ = fieldValue
	if (computedField.GetChildCount() == int32(2)) && (isStatic || this_.transpileClassFields) {
		fieldValue = computedField.GetLastChild()
	}
	var compFieldNameExpr pkg_rhino.Node = pkg_rhino.Node(0)
	_ = compFieldNameExpr
	if this_.transpileClassFields {
		compFieldNameExpr = computedField.RemoveFirstChild()
	} else {
		compFieldNameExpr = computedField.GetFirstChild().CloneTree_0()
	}
	if NodeUtil_GetStringValue(compFieldNameExpr) == "__proto__" {
		return this_.priv_Es6NormalizeClasses_createFieldDefineProperty_2529(computedField, fieldValue, receiver, compFieldNameExpr)
	}
	return this_.priv_Es6NormalizeClasses_createFieldAssignment_2530(computedField, fieldValue, j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return j2gort.BoxNode(this_.astFactory.CreateGetElem(receiver, compFieldNameExpr))
	}))
}

// Java method: Es6NormalizeClasses.createFieldDefineProperty(com.google.javascript.rhino.Node field, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node fieldValue, com.google.javascript.rhino.Node receiver, com.google.javascript.rhino.Node keyNode)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_createFieldDefineProperty_2529(field pkg_rhino.Node, fieldValue pkg_rhino.Node, receiver pkg_rhino.Node, keyNode pkg_rhino.Node) pkg_rhino.Node {
	_ = field
	_ = fieldValue
	_ = receiver
	_ = keyNode
	return this_.priv_Es6NormalizeClasses_createFieldInitializationStatement_2531(field, fieldValue, j2gort.NewLambda_Ext_java_util_function_Function(func(lp_18 any) any {
		val := j2gort.CastNode[pkg_rhino.Node](lp_18)
		_ = val
		var descriptor pkg_rhino.Node = this_.astFactory.CreateObjectLit_1([]pkg_rhino.Node{this_.astFactory.CreateStringKey("configurable", this_.astFactory.CreateBoolean(true)), this_.astFactory.CreateStringKey("enumerable", this_.astFactory.CreateBoolean(true)), this_.astFactory.CreateStringKey("writable", this_.astFactory.CreateBoolean(true)), this_.astFactory.CreateStringKey("value", val)})
		_ = descriptor
		var definePropertyCall pkg_rhino.Node = this_.astFactory.CreateCallWithUnknownType(this_.priv_Es6NormalizeClasses_createObjectDotDefineProperty_2526(), []pkg_rhino.Node{receiver, keyNode, descriptor})
		_ = definePropertyCall
		return j2gort.BoxNode(this_.astFactory.ExprResult(definePropertyCall))
	}))
}

// Java method: Es6NormalizeClasses.createFieldAssignment(com.google.javascript.rhino.Node field, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node fieldValue, java.util.function.Supplier<com.google.javascript.rhino.Node> createAssignLhs)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_createFieldAssignment_2530(field pkg_rhino.Node, fieldValue pkg_rhino.Node, createAssignLhs j2gort.Ext_java_util_function_Supplier) pkg_rhino.Node {
	_ = field
	_ = fieldValue
	_ = createAssignLhs
	return this_.priv_Es6NormalizeClasses_createFieldInitializationStatement_2531(field, fieldValue, j2gort.NewLambda_Ext_java_util_function_Function(func(lp_19 any) any {
		val := j2gort.CastNode[pkg_rhino.Node](lp_19)
		_ = val
		return j2gort.BoxNode(this_.astFactory.CreateAssignStatement(j2gort.CastNode[pkg_rhino.Node](createAssignLhs.Get()), val))
	}))
}

// Java method: Es6NormalizeClasses.createFieldInitializationStatement(com.google.javascript.rhino.Node field, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node fieldValue, java.util.function.Function<com.google.javascript.rhino.Node,com.google.javascript.rhino.Node> createStatement)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_createFieldInitializationStatement_2531(field pkg_rhino.Node, fieldValue pkg_rhino.Node, createStatement j2gort.Ext_java_util_function_Function) pkg_rhino.Node {
	_ = field
	_ = fieldValue
	_ = createStatement
	if this_.transpileClassFields && !(field.IsPrivateIdentifier()) {
		field.Detach()
	}
	if fieldValue == 0 {
		if !(this_.transpileClassFields) {
			return pkg_rhino.Node(0)
		}
		fieldValue = this_.astFactory.CreateUndefinedValue()
	} else {
		fieldValue.Detach()
	}
	var result pkg_rhino.Node = j2gort.CastNode[pkg_rhino.Node](createStatement.Apply(j2gort.BoxNode(fieldValue)))
	_ = result
	if this_.transpileClassFields {
		result.GetFirstChild().SetJSDocInfo(field.GetJSDocInfo())
	}
	result.SrcrefTreeIfMissing(field)
	return result
}

// Java method: Es6NormalizeClasses.rewriteStaticMembers(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.Es6NormalizeClasses.ClassRecord record, com.google.javascript.rhino.Node insertionPoint)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_rewriteStaticMembers_2532(t *NodeTraversal, record *Es6NormalizeClasses_ClassRecord, insertionPoint pkg_rhino.Node) {
	_ = t
	_ = record
	_ = insertionPoint
	var staticMembers *j2gort.Deque[pkg_rhino.Node] = record.staticMembers
	_ = staticMembers
	var hasPriv bool = (this_.transpilePrivateClassFields && record.hasStaticPrivateMembers)
	_ = hasPriv
	if staticMembers.IsEmpty() && !(hasPriv) {
		return
	}
	var staticInitTempParent pkg_rhino.Node = pkg_rhino.IR_Block_0()
	_ = staticInitTempParent
	if hasPriv {
		var staticMapName string = record.GetStaticPrivateMapVarName(this_.compiler, t)
		_ = staticMapName
		var privName string = record.GetStaticPrivateVarName(this_.compiler, t)
		_ = privName
		for _, pm := range record.privateMembers {
			_ = pm
			if pm.isStatic && (pm.kind == StaticGet_Es6NormalizeClasses_ClassRecord_PrivateMemberKind_ACCESSOR()) {
				if ((pm.getterNode != 0) && pm.getterNode.HasChildren()) && this_.priv_Es6NormalizeClasses_hasThisReference_2516(pm.getterNode.GetFirstChild()) {
					{
						r_22 := record
						r_22.needsStaticSelf = true
					}
				}
				if ((pm.setterNode != 0) && pm.setterNode.HasChildren()) && this_.priv_Es6NormalizeClasses_hasThisReference_2516(pm.setterNode.GetFirstChild()) {
					{
						r_23 := record
						r_23.needsStaticSelf = true
					}
				}
			}
		}

		var createProp pkg_rhino.Node = this_.astFactory.CreateGetPropWithUnknownType(this_.astFactory.CreateNameWithUnknownType("Object"), "create")
		_ = createProp
		var staticPrivInit pkg_rhino.Node = this_.astFactory.CreateCallWithUnknownType(createProp, []pkg_rhino.Node{this_.astFactory.CreateNull()})
		_ = staticPrivInit
		var privDecl pkg_rhino.Node = this_.astFactory.CreateSingleConstNameDeclaration(privName, staticPrivInit).SrcrefTreeIfMissing(record.classNode)
		_ = privDecl
		staticInitTempParent.AddChildToBack(privDecl)
		if record.needsStaticSelf {
			var privSelf pkg_rhino.Node = this_.astFactory.CreateGetPropWithUnknownType(this_.astFactory.CreateConstantName(privName, AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN())), "$self")
			_ = privSelf
			var assignSelf pkg_rhino.Node = this_.astFactory.CreateAssignStatement(privSelf, record.CreateNewNameReferenceNode().SrcrefTreeIfMissing(record.classNode)).SrcrefTreeIfMissing(record.classNode)
			_ = assignSelf
			staticInitTempParent.AddChildToBack(assignSelf)
		}
		var mapNameNode pkg_rhino.Node = this_.astFactory.CreateNameWithUnknownType(staticMapName)
		_ = mapNameNode
		var mapSetCallee pkg_rhino.Node = this_.astFactory.CreateGetPropWithUnknownType(mapNameNode, "set")
		_ = mapSetCallee
		var mapSetCall pkg_rhino.Node = this_.astFactory.CreateCallWithUnknownType(mapSetCallee, []pkg_rhino.Node{record.CreateNewNameReferenceNode().SrcrefTreeIfMissing(record.classNode), this_.astFactory.CreateConstantName(privName, AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN()))})
		_ = mapSetCall
		var mapSetStmt pkg_rhino.Node = this_.astFactory.ExprResult(mapSetCall).SrcrefTreeIfMissing(record.classNode)
		_ = mapSetStmt
		staticInitTempParent.AddChildToBack(mapSetStmt)
		this_.priv_Es6NormalizeClasses_addStaticPrivateMembers_2533(record, privName, staticInitTempParent)
	}
	for !(staticMembers.IsEmpty()) {
		var staticMember pkg_rhino.Node = staticMembers.Pop()
		_ = staticMember
		if staticMember.IsPrivateIdentifier() {
			if this_.transpilePrivateClassFields {
				var privName_2 string = record.GetStaticPrivateVarName(this_.compiler, t)
				_ = privName_2
				var propName string = j2gort.JStrSubFrom(staticMember.GetString(), int32(1))
				_ = propName
				var lhs pkg_rhino.Node = this_.astFactory.CreateGetPropWithUnknownType(this_.astFactory.CreateConstantName(privName_2, AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN())), propName)
				_ = lhs
				var fieldValue pkg_rhino.Node
				_ = fieldValue
				if staticMember.HasChildren() {
					fieldValue = staticMember.RemoveFirstChild()
				} else {
					fieldValue = this_.astFactory.CreateUndefinedValue()
				}
				var assign pkg_rhino.Node = this_.astFactory.CreateAssignStatement(lhs, fieldValue)
				_ = assign
				assign.SrcrefTreeIfMissing(staticMember)
				staticInitTempParent.AddChildToBack(assign)
				staticMember.Detach()
			}
		} else {
			var nameToUse pkg_rhino.Node = record.CreateNewNameReferenceNode().SrcrefTree(staticMember)
			_ = nameToUse
			var transpiledNode pkg_rhino.Node
		iifeSwitch_1:
			switch {
			default:
				switch staticMember.TokenSwitchOrd_() {
				case int32(90) /* Token.BLOCK */ :
					var block pkg_rhino.Node = staticMember.Detach()
					_ = block
					t.ReportCodeChange_0()
					{

						if block.HasChildren() {
							transpiledNode = block
							break iifeSwitch_1
						}
						transpiledNode = pkg_rhino.Node(0)
						break iifeSwitch_1
					}

				case int32(108) /* Token.MEMBER_FIELD_DEF */ :
					transpiledNode = this_.priv_Es6NormalizeClasses_convNonCompFieldToGetProp_2527(nameToUse, staticMember, true)
					break iifeSwitch_1
				case int32(109) /* Token.COMPUTED_FIELD_DEF */ :
					transpiledNode = this_.priv_Es6NormalizeClasses_convCompFieldToGetElem_2528(nameToUse, staticMember, true)
					break iifeSwitch_1
				default:
					panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.StrVal(staticMember))))
				}
				panic("unreachable switch expression")
			}

			_ = transpiledNode
			if transpiledNode != 0 {
				staticInitTempParent.AddChildToBack(transpiledNode)
			}
		}
	}
	if this_.transpilePrivateClassFields {
		for _, pm_2 := range record.privateMembers {
			_ = pm_2
			if pm_2.isStatic {
				{
					cur_26 := j2gort.NewIterCursor(j2gort.CheckNotNull(pm_2.declarationNodes))
					for cur_26.HasNext() {
						declNode := pkg_rhino.Node(cur_26.NextNode())
						_ = declNode
						if declNode.GetParent() != 0 {
							NodeUtil_DeleteNode(declNode, this_.compiler)
						}
					}
				}
			}
		}

	}
	if !(staticInitTempParent.HasChildren()) {
		return
	}
	if !(hasPriv) && !(NodeUtil_Has(staticInitTempParent, lambdaSingleton_es6_normalize_classes_4, lambdaSingleton_es6_normalize_classes_5)) {
		for staticInitTempParent.HasChildren() {
			staticInitTempParent.GetLastChild().Detach().InsertAfter(insertionPoint)
		}
		t.ReportCodeChange_1(insertionPoint.GetParent())
		return
	}
	var staticInitMethodName string = this_.priv_Es6NormalizeClasses_generateUniqueStaticInitMethodName_2498(t)
	_ = staticInitMethodName
	var staticInitMethodFunc pkg_rhino.Node = this_.astFactory.CreateEmptyFunction(AstFactory_S_Type_2_JSTypeNative_Color(pkg_rhino.StaticGet_JSTypeNative_FUNCTION_TYPE(), pkg_colors.StaticGet_StandardColors_TOP_OBJECT()))
	_ = staticInitMethodFunc
	var staticInitMethod pkg_rhino.Node = this_.astFactory.CreateMemberFunctionDef(staticInitMethodName, staticInitMethodFunc)
	_ = staticInitMethod
	staticInitMethod.SetStaticMember(true)
	var staticInitBlock pkg_rhino.Node = staticInitMethod.GetFirstChild().GetLastChild()
	_ = staticInitBlock
	staticInitBlock.AddChildrenToFront(staticInitTempParent.RemoveChildren())
	var classMembers pkg_rhino.Node = NodeUtil_GetClassMembers(record.classNode)
	_ = classMembers
	staticInitMethod.SrcrefTreeIfMissing(classMembers)
	classMembers.AddChildToBack(staticInitMethod)
	this_.compiler.ReportChangeToChangeScope(staticInitMethodFunc)
	t.ReportCodeChange_1(staticInitMethod)
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_MEMBER_DECLARATIONS(), this_.compiler)
	var staticInitCall pkg_rhino.Node = this_.astFactory.CreateCall(this_.astFactory.CreateGetProp(record.CreateNewNameReferenceNode(), staticInitMethodName, AstFactory_S_Type_1_Node(staticInitMethodFunc)), AstFactory_S_Type_2_JSTypeNative_Color(pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE(), pkg_colors.StaticGet_StandardColors_NULL_OR_VOID()), []pkg_rhino.Node{}).SrcrefTreeIfMissing(classMembers)
	_ = staticInitCall
	var staticInitCallStmt pkg_rhino.Node = this_.astFactory.ExprResult(staticInitCall)
	_ = staticInitCallStmt
	staticInitCallStmt.InsertAfter(insertionPoint)
}

// Java method: Es6NormalizeClasses.addStaticPrivateMembers(com.google.javascript.jscomp.Es6NormalizeClasses.ClassRecord record, java.lang.String privName, com.google.javascript.rhino.Node staticInitTempParent)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_addStaticPrivateMembers_2533(record *Es6NormalizeClasses_ClassRecord, privName string, staticInitTempParent pkg_rhino.Node) {
	_ = record
	_ = privName
	_ = staticInitTempParent
	for _, pm := range record.privateMembers {
		_ = pm
		if pm.isStatic {
			if ((pm.kind == StaticGet_Es6NormalizeClasses_ClassRecord_PrivateMemberKind_METHOD()) && (pm.methodNode != 0)) && pm.methodNode.HasChildren() {
				this_.priv_Es6NormalizeClasses_rewriteSuperReferencesInPrivateMember_2517(pm.methodNode, record)
				var fnNode pkg_rhino.Node = pm.methodNode.GetFirstChild().Detach()
				_ = fnNode
				pm.methodNode.Detach()
				var lhs pkg_rhino.Node = this_.astFactory.CreateGetPropWithUnknownType(this_.astFactory.CreateConstantName(privName, AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN())), pm.PropertyName())
				_ = lhs
				var assign pkg_rhino.Node = this_.astFactory.CreateAssignStatement(lhs, fnNode)
				_ = assign
				assign.SrcrefTreeIfMissing(pm.methodNode)
				staticInitTempParent.AddChildToBack(assign)
			} else if pm.kind == StaticGet_Es6NormalizeClasses_ClassRecord_PrivateMemberKind_ACCESSOR() {
				var descriptorLit pkg_rhino.Node = this_.astFactory.CreateObjectLit_1([]pkg_rhino.Node{})
				_ = descriptorLit
				if (pm.getterNode != 0) && pm.getterNode.HasChildren() {
					this_.priv_Es6NormalizeClasses_rewriteSuperReferencesInPrivateMember_2517(pm.getterNode, record)
					var getFn pkg_rhino.Node = pm.getterNode.GetFirstChild().Detach()
					_ = getFn
					pm.getterNode.Detach()
					this_.priv_Es6NormalizeClasses_rewriteThisToThisDotSelf_2515(getFn)
					descriptorLit.AddChildToBack(this_.astFactory.CreateStringKey("get", getFn))
				}
				if (pm.setterNode != 0) && pm.setterNode.HasChildren() {
					this_.priv_Es6NormalizeClasses_rewriteSuperReferencesInPrivateMember_2517(pm.setterNode, record)
					var setFn pkg_rhino.Node = pm.setterNode.GetFirstChild().Detach()
					_ = setFn
					pm.setterNode.Detach()
					this_.priv_Es6NormalizeClasses_rewriteThisToThisDotSelf_2515(setFn)
					descriptorLit.AddChildToBack(this_.astFactory.CreateStringKey("set", setFn))
				}
				var defineProp pkg_rhino.Node = this_.astFactory.CreateGetPropWithUnknownType(this_.astFactory.CreateNameWithUnknownType("Object"), "defineProperty")
				_ = defineProp
				var call pkg_rhino.Node = this_.astFactory.CreateCallWithUnknownType(defineProp, []pkg_rhino.Node{this_.astFactory.CreateConstantName(privName, AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN())), this_.astFactory.CreateString(pm.PropertyName()), descriptorLit})
				_ = call
				var declNode pkg_rhino.Node = pm.GetPrimaryDeclarationNode()
				_ = declNode
				var exprResult pkg_rhino.Node = this_.astFactory.ExprResult(call)
				_ = exprResult
				if declNode != 0 {
					exprResult.SrcrefTreeIfMissing(declNode)
				}
				staticInitTempParent.AddChildToBack(exprResult)
			}
		}
	}

}

// Java method: Es6NormalizeClasses.maybeUpdateClassSelfRef(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node nameNode)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_maybeUpdateClassSelfRef_2534(t *NodeTraversal, nameNode pkg_rhino.Node) {
	_ = t
	_ = nameNode
	{
		cur_31 := this_.classStack.Iter()
		for cur_31.HasNext() {

			klass := cur_31.Next()
			_ = klass
			if klass.bindingIdentifier.IsEmpty() {
				continue
			}
			var bindingIdentifier pkg_rhino.Node = j2gort.CastNode[pkg_rhino.Node](klass.bindingIdentifier.Get())
			_ = bindingIdentifier
			if (bindingIdentifier.MatchesName_1_Node(klass.normalizedNames.classNameNode) || (nameNode == bindingIdentifier)) || !(nameNode.MatchesQualifiedName_1_Node(bindingIdentifier)) {
				continue
			}
			var var_ *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(nameNode.GetString()))
			_ = var_
			if (var_ == nil) || (var_.GetNameNode() != bindingIdentifier) {
				continue
			}
			var newNameNode pkg_rhino.Node = klass.CreateNewNameReferenceNode().Srcref(nameNode)
			_ = newNameNode
			nameNode.ReplaceWith(newNameNode)
			this_.compiler.ReportChangeToEnclosingScope(newNameNode)
			return
		}
	}
}

// Java method: Es6NormalizeClasses.visitThisAndSuper(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_visitThisAndSuper_2535(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	if Es6NormalizeClasses_IsSuperCall(n.GetParent()) {
		var classRecord *Es6NormalizeClasses_ClassRecord = this_.classStack.Peek()
		_ = classRecord
		classRecord.superCalls = append(classRecord.superCalls, n.GetParent())
		return
	}
	var rootNode pkg_rhino.Node = t.GetClosestScopeRootNodeBindingThisOrSuper()
	_ = rootNode
	var rootParent pkg_rhino.Node = rootNode.GetParent()
	_ = rootParent
	if rootNode.IsFunction() && (((rootParent.IsMemberFunctionDef() || rootParent.IsComputedProp()) || rootParent.IsGetterDef()) || rootParent.IsSetterDef()) {
		rootNode = rootParent
	}
	var inStaticBlock bool = (rootNode.IsBlock() && rootNode.GetParent().IsClassMembers())
	_ = inStaticBlock
	if !(rootNode.IsStaticMember()) && !(inStaticBlock) {
		return
	}
	var forceRewriteStaticThis bool = (n.IsThis() && Es6NormalizeClasses_IsInStaticInitializationContext(rootNode))
	_ = forceRewriteStaticThis
	if !(this_.doStaticInheritanceRewrites) && !(forceRewriteStaticThis) {
		return
	}
	var classNode pkg_rhino.Node = rootNode.GetGrandparent()
	_ = classNode
	var classRecord_2 *Es6NormalizeClasses_ClassRecord = this_.classStack.Peek()
	_ = classRecord_2
	if !(classRecord_2.classNode == classNode) {
		j2gort.CheckStateFmt(false, "wrong class node: %s != %s", j2gort.BoxNode(classRecord_2.classNode), j2gort.BoxNode(classNode))
	}
	var newNameNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = newNameNode
	if forceRewriteStaticThis {
		newNameNode = classRecord_2.CreateNewNameReferenceNode().SrcrefTree(n)
	} else if this_.doStaticSuperRewrites && n.IsSuper() {
		newNameNode = classRecord_2.CreateNewSuperclassReferenceNode()
	}
	if newNameNode != 0 {
		var originalName string = j2gort.CheckNotNullString(n.GetQualifiedName())
		_ = originalName
		newNameNode.SrcrefTree(n).SetOriginalName(originalName)
		n.ReplaceWith(newNameNode)
		t.ReportCodeChange_1(newNameNode)
	}
}

// Java static method: Es6NormalizeClasses.isSuperCall(com.google.javascript.rhino.Node n)
func Es6NormalizeClasses_IsSuperCall(n pkg_rhino.Node) bool {
	clinit_Es6NormalizeClasses()
	_ = n
	return (n.IsCall() && n.GetFirstChild().IsSuper())
}

// Java static method: Es6NormalizeClasses.isInStaticInitializationContext(com.google.javascript.rhino.Node rootNode)
func Es6NormalizeClasses_IsInStaticInitializationContext(rootNode pkg_rhino.Node) bool {
	clinit_Es6NormalizeClasses()
	_ = rootNode
	switch rootNode.TokenSwitchOrd_() {
	case int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(129) /* Token.COMPUTED_PROP */, int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */ :
		if !(rootNode.IsStaticMember()) {
			j2gort.CheckState(false)
		}
		return false
	case int32(108) /* Token.MEMBER_FIELD_DEF */, int32(109) /* Token.COMPUTED_FIELD_DEF */ :
		if !(rootNode.IsStaticMember()) {
			j2gort.CheckState(false)
		}
		return true
	case int32(90) /* Token.BLOCK */ :
		if !(rootNode.GetParent().IsClassMembers()) {
			j2gort.CheckState(false)
		}
		return true
	default:
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "invalid root "+j2gort.StrVal(rootNode))))
	}
	panic("unreachable")
}

// Java method: Es6NormalizeClasses.generateNormalizedClassNames(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node classNode)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_generateNormalizedClassNames_2536(t *NodeTraversal, classNode pkg_rhino.Node) *Es6NormalizeClasses_NormalizedNames {
	_ = t
	_ = classNode
	var extendsClause pkg_rhino.Node = classNode.GetSecondChild()
	_ = extendsClause
	var shouldExtractExtends bool = this_.priv_Es6NormalizeClasses_shouldExtractExtends_2499(classNode)
	_ = shouldExtractExtends
	var superClass *j2gort.Ext_java_util_Optional = nil
	_ = superClass
	if shouldExtractExtends {
		var name string = this_.priv_Es6NormalizeClasses_generateUniqueClassExtendsVarName_2496(t)
		_ = name
		superClass = j2gort.Ext_java_util_Optional_Of(j2gort.BoxNode(this_.astFactory.CreateName_2_string_AstFactory_Type(name, AstFactory_S_Type_1_Node(extendsClause))))
	} else {
		if extendsClause.IsEmpty() {
			superClass = j2gort.Ext_java_util_Optional_Empty()
		} else {
			superClass = j2gort.Ext_java_util_Optional_Of(j2gort.BoxNode(extendsClause))
		}
	}
	var shouldExtractClass bool = (this_.priv_Es6NormalizeClasses_shouldExtractClass_2503(classNode) || (shouldExtractExtends && !(this_.priv_Es6NormalizeClasses_canExtractExtendsSimply_2500(classNode))))
	_ = shouldExtractClass
	var classNameNodeTemplate pkg_rhino.Node = pkg_rhino.Node(0)
	_ = classNameNodeTemplate
	if shouldExtractClass {
		var name_2 string = this_.priv_Es6NormalizeClasses_generateUniqueClassDeclVarName_2495(t)
		_ = name_2
		classNameNodeTemplate = this_.astFactory.CreateConstantName(name_2, AstFactory_S_Type_1_Node(classNode))
	} else {
		classNameNodeTemplate = NodeUtil_GetNameNode(classNode)
	}
	return NewEs6NormalizeClasses_NormalizedNames(classNameNodeTemplate, superClass, shouldExtractClass, shouldExtractExtends)
}

// Java method: Es6NormalizeClasses.receiverNeedsTempVar(com.google.javascript.rhino.Node receiver)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_receiverNeedsTempVar_2537(receiver pkg_rhino.Node) bool {
	_ = receiver
	return (!(receiver.IsThis()) && !(receiver.IsName()))
}

// Java method: Es6NormalizeClasses.getDeclarationInsertionPoint(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.NodeTraversal t)
func (this_ *Es6NormalizeClasses) priv_Es6NormalizeClasses_getDeclarationInsertionPoint_2538(n pkg_rhino.Node, t *NodeTraversal) pkg_rhino.Node {
	_ = n
	_ = t
	var enclosingStmt pkg_rhino.Node = NodeUtil_GetEnclosingStatement(n)
	_ = enclosingStmt
	if (enclosingStmt != 0) && (enclosingStmt.GetParent() != 0) {
		return enclosingStmt.GetParent()
	}
	return t.GetScopeRoot()
}

// Java static method: Es6NormalizeClasses.isPrivateMember(com.google.javascript.rhino.Node n)
func Es6NormalizeClasses_IsPrivateMember(n pkg_rhino.Node) bool {
	clinit_Es6NormalizeClasses()
	_ = n
	return (n.IsPrivateIdentifier() || (n.HasChildren() && n.GetFirstChild().IsPrivateIdentifier()))
}

// Java record: com.google.javascript.jscomp.Es6NormalizeClasses.ResolvedPrivateMember
type Es6NormalizeClasses_ResolvedPrivateMember struct {
	classRecord   *Es6NormalizeClasses_ClassRecord               `gson:"classRecord"`
	privateMember *Es6NormalizeClasses_ClassRecord_PrivateMember `gson:"privateMember"`
}

func (this_ *Es6NormalizeClasses_ResolvedPrivateMember) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6NormalizeClasses_ResolvedPrivateMember) IsJavaNil_() bool { return this_ == nil }
func (*Es6NormalizeClasses_ResolvedPrivateMember) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6NormalizeClasses$ResolvedPrivateMember", "ResolvedPrivateMember"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6NormalizeClasses.ResolvedPrivateMember", "com.google.javascript.jscomp.Es6NormalizeClasses$ResolvedPrivateMember", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6NormalizeClasses_ResolvedPrivateMember](v_)
	})
}
func (*Es6NormalizeClasses_ResolvedPrivateMember) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *Es6NormalizeClasses_ResolvedPrivateMember) initFields_Es6NormalizeClasses_ResolvedPrivateMember() {
	_ = this_
}

// Java constructor: ResolvedPrivateMember(com.google.javascript.jscomp.Es6NormalizeClasses.ClassRecord classRecord, com.google.javascript.jscomp.Es6NormalizeClasses.ClassRecord.PrivateMember privateMember)
func NewEs6NormalizeClasses_ResolvedPrivateMember(classRecord *Es6NormalizeClasses_ClassRecord, privateMember *Es6NormalizeClasses_ClassRecord_PrivateMember) *Es6NormalizeClasses_ResolvedPrivateMember {
	this_ := &Es6NormalizeClasses_ResolvedPrivateMember{}
	InitEs6NormalizeClasses_ResolvedPrivateMember(this_, classRecord, privateMember)
	return this_
}

func InitEs6NormalizeClasses_ResolvedPrivateMember(this_ *Es6NormalizeClasses_ResolvedPrivateMember, classRecord *Es6NormalizeClasses_ClassRecord, privateMember *Es6NormalizeClasses_ClassRecord_PrivateMember) {
	_ = classRecord
	_ = privateMember
	this_.initFields_Es6NormalizeClasses_ResolvedPrivateMember()
	this_.classRecord = classRecord
	this_.privateMember = privateMember
}

func (this_ *Es6NormalizeClasses_ResolvedPrivateMember) ClassRecord() *Es6NormalizeClasses_ClassRecord {
	return this_.classRecord
}
func (this_ *Es6NormalizeClasses_ResolvedPrivateMember) PrivateMember() *Es6NormalizeClasses_ClassRecord_PrivateMember {
	return this_.privateMember
}
func (this_ *Es6NormalizeClasses_ResolvedPrivateMember) Equals(other_ any) bool {
	o_, ok_ := other_.(*Es6NormalizeClasses_ResolvedPrivateMember)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.classRecord != o_.classRecord {
		return false
	}
	if this_.privateMember != o_.privateMember {
		return false
	}
	return true
}
func (this_ *Es6NormalizeClasses_ResolvedPrivateMember) HashCode() int32 {
	var h_ int32 = 0
	{
		hp_33 := this_.classRecord
		h_ = 31 * h_
		if hp_33 != nil {
			h_ += j2gort.IdentityHash(hp_33)
		}
	}
	{
		hp_34 := this_.privateMember
		h_ = 31 * h_
		if hp_34 != nil {
			h_ += j2gort.IdentityHash(hp_34)
		}
	}
	return h_
}
func (this_ *Es6NormalizeClasses_ResolvedPrivateMember) String() string {
	_ = fmt.Sprint
	return "ResolvedPrivateMember[" + "classRecord=" + j2gort.StrVal(j2gort.PtrToAny(this_.classRecord)) + ", " + "privateMember=" + j2gort.StrVal(j2gort.PtrToAny(this_.privateMember)) + "]"
}

// Java class: com.google.javascript.jscomp.Es6NormalizeClasses.ClassRecord
type Es6NormalizeClasses_ClassRecord struct {
	privateMembers                       []*Es6NormalizeClasses_ClassRecord_PrivateMember                          `gson:"privateMembers"`
	privateMembersByName                 j2gort.OrderedMap[string, *Es6NormalizeClasses_ClassRecord_PrivateMember] `gson:"privateMembersByName"`
	privateMapVarName                    string                                                                    `gson:"privateMapVarName"`
	staticPrivateMapVarName              string                                                                    `gson:"staticPrivateMapVarName"`
	privateVarName                       string                                                                    `gson:"privateVarName"`
	staticPrivateVarName                 string                                                                    `gson:"staticPrivateVarName"`
	privateProtoVarName                  string                                                                    `gson:"privateProtoVarName"`
	bindingIdentifier                    *j2gort.Ext_java_util_Optional                                            `gson:"bindingIdentifier"`
	normalizedNames                      *Es6NormalizeClasses_NormalizedNames                                      `gson:"normalizedNames"`
	instanceMembers                      *j2gort.Deque[pkg_rhino.Node]                                             `gson:"instanceMembers"`
	staticMembers                        *j2gort.Deque[pkg_rhino.Node]                                             `gson:"staticMembers"`
	computedPropsWithSideEffects         *j2gort.Deque[pkg_rhino.Node]                                             `gson:"computedPropsWithSideEffects"`
	constructorVars                      j2gort.Ext_com_google_common_collect_ImmutableSet                         `gson:"constructorVars"`
	superCalls                           []pkg_rhino.Node                                                          `gson:"superCalls"`
	classNode                            pkg_rhino.Node                                                            `gson:"classNode"`
	idHash_                              atomic.Int32
	hasInstancePrivateMembers            bool `gson:"hasInstancePrivateMembers"`
	hasInstancePrivateMethodsOrAccessors bool `gson:"hasInstancePrivateMethodsOrAccessors"`
	hasStaticPrivateMembers              bool `gson:"hasStaticPrivateMembers"`
	needsInstanceSelf                    bool `gson:"needsInstanceSelf"`
	needsStaticSelf                      bool `gson:"needsStaticSelf"`
}

func (this_ *Es6NormalizeClasses_ClassRecord) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6NormalizeClasses_ClassRecord) IsJavaNil_() bool { return this_ == nil }
func (*Es6NormalizeClasses_ClassRecord) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6NormalizeClasses$ClassRecord", "ClassRecord"
}
func (this_ *Es6NormalizeClasses_ClassRecord) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es6NormalizeClasses_ClassRecord) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6NormalizeClasses.ClassRecord", "com.google.javascript.jscomp.Es6NormalizeClasses$ClassRecord", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6NormalizeClasses_ClassRecord](v_)
	})
}

func (this_ *Es6NormalizeClasses_ClassRecord) initFields_Es6NormalizeClasses_ClassRecord() {
	_ = this_
	this_.privateMembers = []*Es6NormalizeClasses_ClassRecord_PrivateMember{}
	this_.privateMembersByName = j2gort.NewOrderedMap[string, *Es6NormalizeClasses_ClassRecord_PrivateMember]()
	this_.privateMapVarName = j2gort.NullString
	this_.staticPrivateMapVarName = j2gort.NullString
	this_.hasInstancePrivateMembers = false
	this_.hasInstancePrivateMethodsOrAccessors = false
	this_.hasStaticPrivateMembers = false
	this_.needsInstanceSelf = false
	this_.needsStaticSelf = false
	this_.privateVarName = j2gort.NullString
	this_.staticPrivateVarName = j2gort.NullString
	this_.privateProtoVarName = j2gort.NullString
	this_.instanceMembers = &j2gort.Deque[pkg_rhino.Node]{}
	this_.staticMembers = &j2gort.Deque[pkg_rhino.Node]{}
	this_.computedPropsWithSideEffects = &j2gort.Deque[pkg_rhino.Node]{}
	this_.constructorVars = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
	this_.superCalls = []pkg_rhino.Node{}
}

// Java constructor: ClassRecord(com.google.javascript.jscomp.Es6NormalizeClasses.NormalizedNames normalizedNames, java.util.Optional<com.google.javascript.rhino.Node> bindingIdentifier, com.google.javascript.rhino.Node classNode)
func NewEs6NormalizeClasses_ClassRecord(normalizedNames *Es6NormalizeClasses_NormalizedNames, bindingIdentifier *j2gort.Ext_java_util_Optional, classNode pkg_rhino.Node) *Es6NormalizeClasses_ClassRecord {
	this_ := &Es6NormalizeClasses_ClassRecord{}
	InitEs6NormalizeClasses_ClassRecord(this_, normalizedNames, bindingIdentifier, classNode)
	return this_
}

func InitEs6NormalizeClasses_ClassRecord(this_ *Es6NormalizeClasses_ClassRecord, normalizedNames *Es6NormalizeClasses_NormalizedNames, bindingIdentifier *j2gort.Ext_java_util_Optional, classNode pkg_rhino.Node) {
	_ = normalizedNames
	_ = bindingIdentifier
	_ = classNode
	this_.initFields_Es6NormalizeClasses_ClassRecord()
	this_.normalizedNames = normalizedNames
	this_.bindingIdentifier = bindingIdentifier
	this_.classNode = classNode
}

// Java method: ClassRecord.recordPrivateMember(com.google.javascript.rhino.Node n)
func (this_ *Es6NormalizeClasses_ClassRecord) RecordPrivateMember(n pkg_rhino.Node) {
	_ = n
	if !(n.IsPrivateIdentifier()) {
		j2gort.CheckArgument(false)
	}
	var name string = n.GetString()
	_ = name
	if !(javalang.StringStartsWith(j2gort.CheckNotNullString(name), "#")) {
		j2gort.CheckArgumentStr(false, name)
	}
	var isStatic bool = n.IsStaticMember()
	_ = isStatic
	var kind Es6NormalizeClasses_ClassRecord_PrivateMemberKind = Es6NormalizeClasses_ClassRecord_PrivateMemberKind(0)
	_ = kind
	if n.IsMemberFieldDef() {
		kind = StaticGet_Es6NormalizeClasses_ClassRecord_PrivateMemberKind_FIELD()
	} else if n.IsMemberFunctionDef() {
		kind = StaticGet_Es6NormalizeClasses_ClassRecord_PrivateMemberKind_METHOD()
	} else if n.IsGetterDef() || n.IsSetterDef() {
		kind = StaticGet_Es6NormalizeClasses_ClassRecord_PrivateMemberKind_ACCESSOR()
	} else {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected node for private member: "+j2gort.StrVal(n))))
	}
	var pm *Es6NormalizeClasses_ClassRecord_PrivateMember = nil
	_ = pm
	if this_.privateMembersByName.ContainsKey(name) {
		pm = this_.privateMembersByName.Get(name)
		pm.declarationNodes.Add_1(j2gort.BoxNode(n))
		if n.IsGetterDef() {
			{
				r_33 := pm
				r_33.getterNode = n
			}
		} else if n.IsSetterDef() {
			{
				r_34 := pm
				r_34.setterNode = n
			}
		}
	} else {
		pm = NewEs6NormalizeClasses_ClassRecord_PrivateMember(name, kind, isStatic, n)
		this_.privateMembers = append(this_.privateMembers, pm)
		this_.privateMembersByName.Put(name, pm)
	}
	if ((kind == StaticGet_Es6NormalizeClasses_ClassRecord_PrivateMemberKind_METHOD()) || (kind == StaticGet_Es6NormalizeClasses_ClassRecord_PrivateMemberKind_ACCESSOR())) && !(isStatic) {
		this_.hasInstancePrivateMethodsOrAccessors = true
	}
	if isStatic {
		this_.hasStaticPrivateMembers = true
	} else {
		this_.hasInstancePrivateMembers = true
	}
}

// Java method: ClassRecord.hasPrivateMember(java.lang.String name)
func (this_ *Es6NormalizeClasses_ClassRecord) HasPrivateMember(name string) bool {
	_ = name
	return this_.privateMembersByName.ContainsKey(name)
}

// Java method: ClassRecord.getPrivateMember(java.lang.String name)
func (this_ *Es6NormalizeClasses_ClassRecord) GetPrivateMember(name string) *Es6NormalizeClasses_ClassRecord_PrivateMember {
	_ = name
	return this_.privateMembersByName.Get(name)
}

// Java method: ClassRecord.getPrivateMapVarName(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.NodeTraversal t)
func (this_ *Es6NormalizeClasses_ClassRecord) GetPrivateMapVarName(compiler AbstractCompiler, t *NodeTraversal) string {
	_ = compiler
	_ = t
	if this_.privateMapVarName == j2gort.NullString {
		this_.privateMapVarName = Es6NormalizeClasses_GenerateUniquePrivateMapVarName(compiler, t)
	}
	return this_.privateMapVarName
}

// Java method: ClassRecord.getPrivateProtoVarName(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.NodeTraversal t)
func (this_ *Es6NormalizeClasses_ClassRecord) GetPrivateProtoVarName(compiler AbstractCompiler, t *NodeTraversal) string {
	_ = compiler
	_ = t
	if this_.privateProtoVarName == j2gort.NullString {
		this_.privateProtoVarName = Es6NormalizeClasses_GenerateUniquePrivateProtoVarName(compiler, t)
	}
	return this_.privateProtoVarName
}

// Java method: ClassRecord.getStaticPrivateMapVarName(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.NodeTraversal t)
func (this_ *Es6NormalizeClasses_ClassRecord) GetStaticPrivateMapVarName(compiler AbstractCompiler, t *NodeTraversal) string {
	_ = compiler
	_ = t
	if this_.staticPrivateMapVarName == j2gort.NullString {
		this_.staticPrivateMapVarName = Es6NormalizeClasses_GenerateUniqueStaticPrivateMapVarName(compiler, t)
	}
	return this_.staticPrivateMapVarName
}

// Java method: ClassRecord.getPrivateVarName(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.NodeTraversal t)
func (this_ *Es6NormalizeClasses_ClassRecord) GetPrivateVarName(compiler AbstractCompiler, t *NodeTraversal) string {
	_ = compiler
	_ = t
	if this_.privateVarName == j2gort.NullString {
		this_.privateVarName = Es6NormalizeClasses_GenerateUniquePrivateVarName(compiler, t)
	}
	return this_.privateVarName
}

// Java method: ClassRecord.getStaticPrivateVarName(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.NodeTraversal t)
func (this_ *Es6NormalizeClasses_ClassRecord) GetStaticPrivateVarName(compiler AbstractCompiler, t *NodeTraversal) string {
	_ = compiler
	_ = t
	if this_.staticPrivateVarName == j2gort.NullString {
		this_.staticPrivateVarName = Es6NormalizeClasses_GenerateUniquePrivateVarName(compiler, t)
	}
	return this_.staticPrivateVarName
}

// Java method: ClassRecord.recordField(com.google.javascript.rhino.Node field)
func (this_ *Es6NormalizeClasses_ClassRecord) RecordField(field pkg_rhino.Node) {
	_ = field
	if !(field.IsComputedFieldDef() || field.IsMemberFieldDef()) {
		j2gort.CheckArgument(false)
	}
	if field.IsStaticMember() {
		this_.staticMembers.AddLast(j2gort.CheckNotNullNode(field))
	} else {
		this_.instanceMembers.AddLast(j2gort.CheckNotNullNode(field))
	}
}

// Java method: ClassRecord.maybeRecordComputedPropWithSideEffects(com.google.javascript.rhino.Node computedProp)
func (this_ *Es6NormalizeClasses_ClassRecord) MaybeRecordComputedPropWithSideEffects(computedProp pkg_rhino.Node) {
	_ = computedProp
	if !(computedProp.IsComputedProp() || computedProp.IsComputedFieldDef()) {
		j2gort.CheckArgument(false)
	}
	if NodeUtil_CanBeSideEffected_1(computedProp.GetFirstChild()) {
		this_.computedPropsWithSideEffects.AddLast(j2gort.CheckNotNullNode(computedProp))
	}
}

// Java method: ClassRecord.recordStaticBlock(com.google.javascript.rhino.Node block)
func (this_ *Es6NormalizeClasses_ClassRecord) RecordStaticBlock(block pkg_rhino.Node) {
	_ = block
	if !(NodeUtil_IsClassStaticBlock(block)) {
		j2gort.CheckArgument(false)
	}
	this_.staticMembers.AddLast(j2gort.CheckNotNullNode(block))
}

// Java method: ClassRecord.recordConstructorScope(com.google.javascript.jscomp.Scope s)
func (this_ *Es6NormalizeClasses_ClassRecord) RecordConstructorScope(s *Scope) {
	_ = s
	if !(s.IsFunctionBlockScope()) {
		j2gort.CheckArgumentMsg(false, j2gort.PtrToAny(s))
	}
	if !(this_.constructorVars.IsEmpty()) {
		j2gort.CheckStateMsg(false, this_.constructorVars)
	}
	var argsScope *Scope = j2gort.CastPtr[Scope](s.GetParent())
	_ = argsScope
	this_.constructorVars = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder().AddAll_1_Ext_java_lang_Iterable(s.GetAllSymbols())).AddAll_1_Ext_java_lang_Iterable(argsScope.GetAllSymbols())).Build())
}

// Java method: ClassRecord.createNewNameReferenceNode()
func (this_ *Es6NormalizeClasses_ClassRecord) CreateNewNameReferenceNode() pkg_rhino.Node {
	if this_.normalizedNames.ClassNameNode().IsName() {
		return this_.normalizedNames.ClassNameNode().CloneNode_0()
	}
	return this_.normalizedNames.ClassNameNode().CloneTree_0()
}

// Java method: ClassRecord.createNewSuperclassReferenceNode()
func (this_ *Es6NormalizeClasses_ClassRecord) CreateNewSuperclassReferenceNode() pkg_rhino.Node {
	if this_.normalizedNames.SuperClassNode().IsEmpty() {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "no superclass")))
	}
	return j2gort.CastNode[pkg_rhino.Node](this_.normalizedNames.SuperClassNode().Get()).CloneTree_0()
}

// Java enum: com.google.javascript.jscomp.Es6NormalizeClasses.ClassRecord.PrivateMemberKind
type Es6NormalizeClasses_ClassRecord_PrivateMemberKind uint8

const (
	enum_Es6NormalizeClasses_ClassRecord_PrivateMemberKind_FIELD    Es6NormalizeClasses_ClassRecord_PrivateMemberKind = 1
	enum_Es6NormalizeClasses_ClassRecord_PrivateMemberKind_METHOD   Es6NormalizeClasses_ClassRecord_PrivateMemberKind = 2
	enum_Es6NormalizeClasses_ClassRecord_PrivateMemberKind_ACCESSOR Es6NormalizeClasses_ClassRecord_PrivateMemberKind = 3
)

var enumBaseTable_Es6NormalizeClasses_ClassRecord_PrivateMemberKind = [...]j2gort.EnumBase{
	{Name_: "FIELD", Ordinal_: 0},
	{Name_: "METHOD", Ordinal_: 1},
	{Name_: "ACCESSOR", Ordinal_: 2},
}
var enumValues_Es6NormalizeClasses_ClassRecord_PrivateMemberKind = []Es6NormalizeClasses_ClassRecord_PrivateMemberKind{1, 2, 3}

func StaticGet_Es6NormalizeClasses_ClassRecord_PrivateMemberKind_FIELD() Es6NormalizeClasses_ClassRecord_PrivateMemberKind {
	return enum_Es6NormalizeClasses_ClassRecord_PrivateMemberKind_FIELD
}
func StaticGet_Es6NormalizeClasses_ClassRecord_PrivateMemberKind_METHOD() Es6NormalizeClasses_ClassRecord_PrivateMemberKind {
	return enum_Es6NormalizeClasses_ClassRecord_PrivateMemberKind_METHOD
}
func StaticGet_Es6NormalizeClasses_ClassRecord_PrivateMemberKind_ACCESSOR() Es6NormalizeClasses_ClassRecord_PrivateMemberKind {
	return enum_Es6NormalizeClasses_ClassRecord_PrivateMemberKind_ACCESSOR
}

func Es6NormalizeClasses_ClassRecord_PrivateMemberKind_Values() []Es6NormalizeClasses_ClassRecord_PrivateMemberKind {
	return j2gort.CloneSlice(enumValues_Es6NormalizeClasses_ClassRecord_PrivateMemberKind)
}

func Es6NormalizeClasses_ClassRecord_PrivateMemberKind_ValueOf(name string) Es6NormalizeClasses_ClassRecord_PrivateMemberKind {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_Es6NormalizeClasses_ClassRecord_PrivateMemberKind {
		if enumBaseTable_Es6NormalizeClasses_ClassRecord_PrivateMemberKind[i_].Name_ == name {
			return Es6NormalizeClasses_ClassRecord_PrivateMemberKind(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.Es6NormalizeClasses.ClassRecord.PrivateMemberKind." + name))
}

func (this_ Es6NormalizeClasses_ClassRecord_PrivateMemberKind) IsJavaNil_() bool                { return this_ == 0 }
func (Es6NormalizeClasses_ClassRecord_PrivateMemberKind) HasIdentityEquals_()                   {}
func (Es6NormalizeClasses_ClassRecord_PrivateMemberKind) Implements_j2gort_Ext_java_lang_Enum() {}
func (Es6NormalizeClasses_ClassRecord_PrivateMemberKind) Implements_j2gort_Ext_java_lang_constant_Constable() {
}
func (this_ Es6NormalizeClasses_ClassRecord_PrivateMemberKind) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (Es6NormalizeClasses_ClassRecord_PrivateMemberKind) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6NormalizeClasses$ClassRecord$PrivateMemberKind", "PrivateMemberKind"
}
func (this_ Es6NormalizeClasses_ClassRecord_PrivateMemberKind) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_Es6NormalizeClasses_ClassRecord_PrivateMemberKind[this_-1]
}
func (this_ Es6NormalizeClasses_ClassRecord_PrivateMemberKind) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ Es6NormalizeClasses_ClassRecord_PrivateMemberKind) EnumOrdinal_() int32 {
	return int32(this_) - 1
}
func (this_ Es6NormalizeClasses_ClassRecord_PrivateMemberKind) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Es6NormalizeClasses_ClassRecord_PrivateMemberKind[this_-1].Name_
}
func (this_ Es6NormalizeClasses_ClassRecord_PrivateMemberKind) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_Es6NormalizeClasses_ClassRecord_PrivateMemberKind[this_-1].Name_
}
func (this_ Es6NormalizeClasses_ClassRecord_PrivateMemberKind) Equals(other_ any) bool {
	o_, ok_ := other_.(Es6NormalizeClasses_ClassRecord_PrivateMemberKind)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ Es6NormalizeClasses_ClassRecord_PrivateMemberKind) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_Es6NormalizeClasses_ClassRecord_PrivateMemberKind[this_-1].HashCode()
}
func (this_ Es6NormalizeClasses_ClassRecord_PrivateMemberKind) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Es6NormalizeClasses_ClassRecord_PrivateMemberKind[this_-1].CompareTo(other_)
}
func (this_ Es6NormalizeClasses_ClassRecord_PrivateMemberKind) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_Es6NormalizeClasses_ClassRecord_PrivateMemberKind[this_-1].Name_)
}
func init() {
	vals_ := []any{Es6NormalizeClasses_ClassRecord_PrivateMemberKind(1), Es6NormalizeClasses_ClassRecord_PrivateMemberKind(2), Es6NormalizeClasses_ClassRecord_PrivateMemberKind(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.Es6NormalizeClasses.ClassRecord.PrivateMemberKind", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.Es6NormalizeClasses$ClassRecord$PrivateMemberKind", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6NormalizeClasses.ClassRecord.PrivateMemberKind", "com.google.javascript.jscomp.Es6NormalizeClasses$ClassRecord$PrivateMemberKind", func(v_ any) bool {
		return j2gort.InstanceOf[Es6NormalizeClasses_ClassRecord_PrivateMemberKind](v_)
	})
}

// Java class: com.google.javascript.jscomp.Es6NormalizeClasses.ClassRecord.PrivateMember
type Es6NormalizeClasses_ClassRecord_PrivateMember struct {
	name             string                    `gson:"name"`
	declarationNodes j2gort.Ext_java_util_List `gson:"declarationNodes"`
	methodNode       pkg_rhino.Node            `gson:"methodNode"`
	getterNode       pkg_rhino.Node            `gson:"getterNode"`
	setterNode       pkg_rhino.Node            `gson:"setterNode"`
	idHash_          atomic.Int32
	kind             Es6NormalizeClasses_ClassRecord_PrivateMemberKind `gson:"kind"`
	isStatic         bool                                              `gson:"isStatic"`
}

func (this_ *Es6NormalizeClasses_ClassRecord_PrivateMember) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6NormalizeClasses_ClassRecord_PrivateMember) IsJavaNil_() bool { return this_ == nil }
func (*Es6NormalizeClasses_ClassRecord_PrivateMember) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6NormalizeClasses$ClassRecord$PrivateMember", "PrivateMember"
}
func (this_ *Es6NormalizeClasses_ClassRecord_PrivateMember) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es6NormalizeClasses_ClassRecord_PrivateMember) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6NormalizeClasses.ClassRecord.PrivateMember", "com.google.javascript.jscomp.Es6NormalizeClasses$ClassRecord$PrivateMember", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6NormalizeClasses_ClassRecord_PrivateMember](v_)
	})
}

func (this_ *Es6NormalizeClasses_ClassRecord_PrivateMember) initFields_Es6NormalizeClasses_ClassRecord_PrivateMember() {
	_ = this_
	this_.name = j2gort.NullString
	this_.declarationNodes = j2gort.NewExt_java_util_ArrayList()
}

// Java constructor: PrivateMember(java.lang.String name, com.google.javascript.jscomp.Es6NormalizeClasses.ClassRecord.PrivateMemberKind kind, boolean isStatic, com.google.javascript.rhino.Node initialNode)
func NewEs6NormalizeClasses_ClassRecord_PrivateMember(name string, kind Es6NormalizeClasses_ClassRecord_PrivateMemberKind, isStatic bool, initialNode pkg_rhino.Node) *Es6NormalizeClasses_ClassRecord_PrivateMember {
	this_ := &Es6NormalizeClasses_ClassRecord_PrivateMember{}
	InitEs6NormalizeClasses_ClassRecord_PrivateMember(this_, name, kind, isStatic, initialNode)
	return this_
}

func InitEs6NormalizeClasses_ClassRecord_PrivateMember(this_ *Es6NormalizeClasses_ClassRecord_PrivateMember, name string, kind Es6NormalizeClasses_ClassRecord_PrivateMemberKind, isStatic bool, initialNode pkg_rhino.Node) {
	_ = name
	_ = kind
	_ = isStatic
	_ = initialNode
	this_.initFields_Es6NormalizeClasses_ClassRecord_PrivateMember()
	this_.name = name
	this_.kind = kind
	this_.isStatic = isStatic
	this_.declarationNodes.Add_1(j2gort.BoxNode(initialNode))
	if kind == StaticGet_Es6NormalizeClasses_ClassRecord_PrivateMemberKind_METHOD() {
		this_.methodNode = initialNode
	} else if initialNode.IsGetterDef() {
		this_.getterNode = initialNode
	} else if initialNode.IsSetterDef() {
		this_.setterNode = initialNode
	}
}

// Java method: PrivateMember.propertyName()
func (this_ *Es6NormalizeClasses_ClassRecord_PrivateMember) PropertyName() string {
	if javalang.StringStartsWith(j2gort.CheckNotNullString(this_.name), "#") {
		return j2gort.JStrSubFrom(this_.name, int32(1))
	} else {
		return this_.name
	}
}

// Java method: PrivateMember.getPrimaryDeclarationNode()
func (this_ *Es6NormalizeClasses_ClassRecord_PrivateMember) GetPrimaryDeclarationNode() pkg_rhino.Node {
	if this_.getterNode != 0 {
		return this_.getterNode
	}
	if this_.setterNode != 0 {
		return this_.setterNode
	}
	return j2gort.CastNode[pkg_rhino.Node](j2gort.Ext_com_google_common_collect_Iterables_GetFirst(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.declarationNodes), nil))
}

// Java record: com.google.javascript.jscomp.Es6NormalizeClasses.NormalizedNames
type Es6NormalizeClasses_NormalizedNames struct {
	superClassNode       *j2gort.Ext_java_util_Optional `gson:"superClassNode"`
	classNameNode        pkg_rhino.Node                 `gson:"classNameNode"`
	shouldExtractClass   bool                           `gson:"shouldExtractClass"`
	shouldExtractExtends bool                           `gson:"shouldExtractExtends"`
}

func (this_ *Es6NormalizeClasses_NormalizedNames) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6NormalizeClasses_NormalizedNames) IsJavaNil_() bool { return this_ == nil }
func (*Es6NormalizeClasses_NormalizedNames) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6NormalizeClasses$NormalizedNames", "NormalizedNames"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6NormalizeClasses.NormalizedNames", "com.google.javascript.jscomp.Es6NormalizeClasses$NormalizedNames", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6NormalizeClasses_NormalizedNames](v_)
	})
}
func (*Es6NormalizeClasses_NormalizedNames) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *Es6NormalizeClasses_NormalizedNames) initFields_Es6NormalizeClasses_NormalizedNames() {
	_ = this_
}

// Java constructor: NormalizedNames(com.google.javascript.rhino.Node classNameNode, java.util.Optional<com.google.javascript.rhino.Node> superClassNode, boolean shouldExtractClass, boolean shouldExtractExtends)
func NewEs6NormalizeClasses_NormalizedNames(classNameNode pkg_rhino.Node, superClassNode *j2gort.Ext_java_util_Optional, shouldExtractClass bool, shouldExtractExtends bool) *Es6NormalizeClasses_NormalizedNames {
	this_ := &Es6NormalizeClasses_NormalizedNames{}
	InitEs6NormalizeClasses_NormalizedNames(this_, classNameNode, superClassNode, shouldExtractClass, shouldExtractExtends)
	return this_
}

func InitEs6NormalizeClasses_NormalizedNames(this_ *Es6NormalizeClasses_NormalizedNames, classNameNode pkg_rhino.Node, superClassNode *j2gort.Ext_java_util_Optional, shouldExtractClass bool, shouldExtractExtends bool) {
	_ = classNameNode
	_ = superClassNode
	_ = shouldExtractClass
	_ = shouldExtractExtends
	this_.initFields_Es6NormalizeClasses_NormalizedNames()
	if !(classNameNode.IsQualifiedName()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(classNameNode))
	}
	if !(superClassNode.IsEmpty() || j2gort.CastNode[pkg_rhino.Node](superClassNode.Get()).IsQualifiedName()) {
		j2gort.CheckArgumentMsg(false, j2gort.PtrToAny(superClassNode))
	}
	this_.classNameNode = classNameNode
	this_.superClassNode = superClassNode
	this_.shouldExtractClass = shouldExtractClass
	this_.shouldExtractExtends = shouldExtractExtends
}

func (this_ *Es6NormalizeClasses_NormalizedNames) ClassNameNode() pkg_rhino.Node {
	return this_.classNameNode
}
func (this_ *Es6NormalizeClasses_NormalizedNames) SuperClassNode() *j2gort.Ext_java_util_Optional {
	return this_.superClassNode
}
func (this_ *Es6NormalizeClasses_NormalizedNames) ShouldExtractClass() bool {
	return this_.shouldExtractClass
}
func (this_ *Es6NormalizeClasses_NormalizedNames) ShouldExtractExtends() bool {
	return this_.shouldExtractExtends
}
func (this_ *Es6NormalizeClasses_NormalizedNames) Equals(other_ any) bool {
	o_, ok_ := other_.(*Es6NormalizeClasses_NormalizedNames)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.classNameNode != o_.classNameNode {
		return false
	}
	if !javalang.ObjectsEquals(j2gort.PtrToAny(this_.superClassNode), j2gort.PtrToAny(o_.superClassNode)) {
		return false
	}
	if this_.shouldExtractClass != o_.shouldExtractClass {
		return false
	}
	if this_.shouldExtractExtends != o_.shouldExtractExtends {
		return false
	}
	return true
}
func (this_ *Es6NormalizeClasses_NormalizedNames) HashCode() int32 {
	var h_ int32 = 0
	h_ = 31*h_ + int32(this_.classNameNode)
	h_ = 31*h_ + j2gort.ObjectHash(j2gort.PtrToAny(this_.superClassNode))
	if this_.shouldExtractClass {
		h_ = 31*h_ + 1231
	} else {
		h_ = 31*h_ + 1237
	}
	if this_.shouldExtractExtends {
		h_ = 31*h_ + 1231
	} else {
		h_ = 31*h_ + 1237
	}
	return h_
}
func (this_ *Es6NormalizeClasses_NormalizedNames) String() string {
	_ = fmt.Sprint
	return "NormalizedNames[" + "classNameNode=" + j2gort.StrVal(j2gort.BoxNode(this_.classNameNode)) + ", " + "superClassNode=" + j2gort.StrVal(j2gort.PtrToAny(this_.superClassNode)) + ", " + "shouldExtractClass=" + j2gort.StrVal(j2gort.ValueOfBoolean(this_.shouldExtractClass)) + ", " + "shouldExtractExtends=" + j2gort.StrVal(j2gort.ValueOfBoolean(this_.shouldExtractExtends)) + "]"
}
