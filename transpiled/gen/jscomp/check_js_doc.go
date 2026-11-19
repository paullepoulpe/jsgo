/*
 * Copyright 2015 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/CheckJSDoc.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.CheckJSDoc
type CheckJSDoc struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler               AbstractCompiler            `gson:"compiler"`
	checkJsDocTypesVisitor *CheckJSDoc_CheckJsdocTypes `gson:"checkJsDocTypesVisitor"`
	inExterns              bool                        `gson:"inExterns"`
}

func (this_ *CheckJSDoc) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckJSDoc) IsJavaNil_() bool { return this_ == nil }
func (*CheckJSDoc) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckJSDoc", "CheckJSDoc"
}
func (this_ *CheckJSDoc) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*CheckJSDoc) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckJSDoc", "com.google.javascript.jscomp.CheckJSDoc", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckJSDoc](v_)
	})
}
func (*CheckJSDoc) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*CheckJSDoc) Implements_jscomp_CompilerPass()                            {}
func (*CheckJSDoc) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_CheckJSDoc_MISPLACED_MSG_ANNOTATION *DiagnosticType

func StaticGet_CheckJSDoc_MISPLACED_MSG_ANNOTATION() *DiagnosticType {
	clinit_CheckJSDoc()
	return static_CheckJSDoc_MISPLACED_MSG_ANNOTATION
}
func StaticSet_CheckJSDoc_MISPLACED_MSG_ANNOTATION(v *DiagnosticType) {
	clinit_CheckJSDoc()
	static_CheckJSDoc_MISPLACED_MSG_ANNOTATION = v
}

var static_CheckJSDoc_MISPLACED_ANNOTATION *DiagnosticType

func StaticGet_CheckJSDoc_MISPLACED_ANNOTATION() *DiagnosticType {
	clinit_CheckJSDoc()
	return static_CheckJSDoc_MISPLACED_ANNOTATION
}
func StaticSet_CheckJSDoc_MISPLACED_ANNOTATION(v *DiagnosticType) {
	clinit_CheckJSDoc()
	static_CheckJSDoc_MISPLACED_ANNOTATION = v
}

var static_CheckJSDoc_MISPLACED_TS_TYPE_ANNOTATION *DiagnosticType

func StaticGet_CheckJSDoc_MISPLACED_TS_TYPE_ANNOTATION() *DiagnosticType {
	clinit_CheckJSDoc()
	return static_CheckJSDoc_MISPLACED_TS_TYPE_ANNOTATION
}
func StaticSet_CheckJSDoc_MISPLACED_TS_TYPE_ANNOTATION(v *DiagnosticType) {
	clinit_CheckJSDoc()
	static_CheckJSDoc_MISPLACED_TS_TYPE_ANNOTATION = v
}

var static_CheckJSDoc_ANNOTATION_DEPRECATED *DiagnosticType

func StaticGet_CheckJSDoc_ANNOTATION_DEPRECATED() *DiagnosticType {
	clinit_CheckJSDoc()
	return static_CheckJSDoc_ANNOTATION_DEPRECATED
}
func StaticSet_CheckJSDoc_ANNOTATION_DEPRECATED(v *DiagnosticType) {
	clinit_CheckJSDoc()
	static_CheckJSDoc_ANNOTATION_DEPRECATED = v
}

var static_CheckJSDoc_DISALLOWED_MEMBER_JSDOC *DiagnosticType

func StaticGet_CheckJSDoc_DISALLOWED_MEMBER_JSDOC() *DiagnosticType {
	clinit_CheckJSDoc()
	return static_CheckJSDoc_DISALLOWED_MEMBER_JSDOC
}
func StaticSet_CheckJSDoc_DISALLOWED_MEMBER_JSDOC(v *DiagnosticType) {
	clinit_CheckJSDoc()
	static_CheckJSDoc_DISALLOWED_MEMBER_JSDOC = v
}

var static_CheckJSDoc_ARROW_FUNCTION_AS_CONSTRUCTOR *DiagnosticType

func StaticGet_CheckJSDoc_ARROW_FUNCTION_AS_CONSTRUCTOR() *DiagnosticType {
	clinit_CheckJSDoc()
	return static_CheckJSDoc_ARROW_FUNCTION_AS_CONSTRUCTOR
}
func StaticSet_CheckJSDoc_ARROW_FUNCTION_AS_CONSTRUCTOR(v *DiagnosticType) {
	clinit_CheckJSDoc()
	static_CheckJSDoc_ARROW_FUNCTION_AS_CONSTRUCTOR = v
}

var static_CheckJSDoc_BAD_REST_PARAMETER_ANNOTATION *DiagnosticType

func StaticGet_CheckJSDoc_BAD_REST_PARAMETER_ANNOTATION() *DiagnosticType {
	clinit_CheckJSDoc()
	return static_CheckJSDoc_BAD_REST_PARAMETER_ANNOTATION
}
func StaticSet_CheckJSDoc_BAD_REST_PARAMETER_ANNOTATION(v *DiagnosticType) {
	clinit_CheckJSDoc()
	static_CheckJSDoc_BAD_REST_PARAMETER_ANNOTATION = v
}

var static_CheckJSDoc_DEFAULT_PARAM_MUST_BE_MARKED_OPTIONAL *DiagnosticType

func StaticGet_CheckJSDoc_DEFAULT_PARAM_MUST_BE_MARKED_OPTIONAL() *DiagnosticType {
	clinit_CheckJSDoc()
	return static_CheckJSDoc_DEFAULT_PARAM_MUST_BE_MARKED_OPTIONAL
}
func StaticSet_CheckJSDoc_DEFAULT_PARAM_MUST_BE_MARKED_OPTIONAL(v *DiagnosticType) {
	clinit_CheckJSDoc()
	static_CheckJSDoc_DEFAULT_PARAM_MUST_BE_MARKED_OPTIONAL = v
}

var static_CheckJSDoc_INVALID_NO_SIDE_EFFECT_ANNOTATION *DiagnosticType

func StaticGet_CheckJSDoc_INVALID_NO_SIDE_EFFECT_ANNOTATION() *DiagnosticType {
	clinit_CheckJSDoc()
	return static_CheckJSDoc_INVALID_NO_SIDE_EFFECT_ANNOTATION
}
func StaticSet_CheckJSDoc_INVALID_NO_SIDE_EFFECT_ANNOTATION(v *DiagnosticType) {
	clinit_CheckJSDoc()
	static_CheckJSDoc_INVALID_NO_SIDE_EFFECT_ANNOTATION = v
}

var static_CheckJSDoc_INVALID_MODIFIES_ANNOTATION *DiagnosticType

func StaticGet_CheckJSDoc_INVALID_MODIFIES_ANNOTATION() *DiagnosticType {
	clinit_CheckJSDoc()
	return static_CheckJSDoc_INVALID_MODIFIES_ANNOTATION
}
func StaticSet_CheckJSDoc_INVALID_MODIFIES_ANNOTATION(v *DiagnosticType) {
	clinit_CheckJSDoc()
	static_CheckJSDoc_INVALID_MODIFIES_ANNOTATION = v
}

var static_CheckJSDoc_INVALID_DEFINE_ON_LET *DiagnosticType

func StaticGet_CheckJSDoc_INVALID_DEFINE_ON_LET() *DiagnosticType {
	clinit_CheckJSDoc()
	return static_CheckJSDoc_INVALID_DEFINE_ON_LET
}
func StaticSet_CheckJSDoc_INVALID_DEFINE_ON_LET(v *DiagnosticType) {
	clinit_CheckJSDoc()
	static_CheckJSDoc_INVALID_DEFINE_ON_LET = v
}

var static_CheckJSDoc_MISPLACED_SUPPRESS *DiagnosticType

func StaticGet_CheckJSDoc_MISPLACED_SUPPRESS() *DiagnosticType {
	clinit_CheckJSDoc()
	return static_CheckJSDoc_MISPLACED_SUPPRESS
}
func StaticSet_CheckJSDoc_MISPLACED_SUPPRESS(v *DiagnosticType) {
	clinit_CheckJSDoc()
	static_CheckJSDoc_MISPLACED_SUPPRESS = v
}

var static_CheckJSDoc_JSDOC_IN_BLOCK_COMMENT *DiagnosticType

func StaticGet_CheckJSDoc_JSDOC_IN_BLOCK_COMMENT() *DiagnosticType {
	clinit_CheckJSDoc()
	return static_CheckJSDoc_JSDOC_IN_BLOCK_COMMENT
}
func StaticSet_CheckJSDoc_JSDOC_IN_BLOCK_COMMENT(v *DiagnosticType) {
	clinit_CheckJSDoc()
	static_CheckJSDoc_JSDOC_IN_BLOCK_COMMENT = v
}

var static_CheckJSDoc_JSDOC_ON_RETURN *DiagnosticType

func StaticGet_CheckJSDoc_JSDOC_ON_RETURN() *DiagnosticType {
	clinit_CheckJSDoc()
	return static_CheckJSDoc_JSDOC_ON_RETURN
}
func StaticSet_CheckJSDoc_JSDOC_ON_RETURN(v *DiagnosticType) {
	clinit_CheckJSDoc()
	static_CheckJSDoc_JSDOC_ON_RETURN = v
}

var static_CheckJSDoc_COMMENT_PATTERN *j2gort.Ext_java_util_regex_Pattern

func StaticGet_CheckJSDoc_COMMENT_PATTERN() *j2gort.Ext_java_util_regex_Pattern {
	clinit_CheckJSDoc()
	return static_CheckJSDoc_COMMENT_PATTERN
}
func StaticSet_CheckJSDoc_COMMENT_PATTERN(v *j2gort.Ext_java_util_regex_Pattern) {
	clinit_CheckJSDoc()
	static_CheckJSDoc_COMMENT_PATTERN = v
}

var clinitDone_CheckJSDoc uint32
var clinitState_CheckJSDoc j2gort.ClinitState

func clinit_CheckJSDoc() {
	if atomic.LoadUint32(&clinitDone_CheckJSDoc) != 0 {
		return
	}
	clinitSlow_CheckJSDoc()
}

func clinitSlow_CheckJSDoc() {
	if !j2gort.EnterClinit(&clinitState_CheckJSDoc, &clinitDone_CheckJSDoc) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckJSDoc, &clinitDone_CheckJSDoc)
	static_CheckJSDoc_MISPLACED_MSG_ANNOTATION = DiagnosticType_Disabled("JSC_MISPLACED_MSG_ANNOTATION", "Misplaced message annotation. @desc, @meaning, and @alternateMessageId annotations should be only on message nodes.\nMessage constants must be prefixed with 'MSG_'.")
	static_CheckJSDoc_MISPLACED_ANNOTATION = DiagnosticType_Warning("JSC_MISPLACED_ANNOTATION", "Misplaced {0} annotation. {1}")
	static_CheckJSDoc_MISPLACED_TS_TYPE_ANNOTATION = DiagnosticType_Warning("JSC_MISPLACED_TS_TYPE_ANNOTATION", "Misplaced @tsType annotation. {0}")
	static_CheckJSDoc_ANNOTATION_DEPRECATED = DiagnosticType_Warning("JSC_ANNOTATION_DEPRECATED", "The {0} annotation is deprecated. {1}")
	static_CheckJSDoc_DISALLOWED_MEMBER_JSDOC = DiagnosticType_Warning("JSC_DISALLOWED_MEMBER_JSDOC", "Class level JSDocs (@interface, @extends, etc.) are not allowed on class members")
	static_CheckJSDoc_ARROW_FUNCTION_AS_CONSTRUCTOR = DiagnosticType_Error("JSC_ARROW_FUNCTION_AS_CONSTRUCTOR", "Arrow functions cannot be used as constructors")
	static_CheckJSDoc_BAD_REST_PARAMETER_ANNOTATION = DiagnosticType_Warning("JSC_BAD_REST_PARAMETER_ANNOTATION", "Missing \"...\" in type annotation for rest parameter.")
	static_CheckJSDoc_DEFAULT_PARAM_MUST_BE_MARKED_OPTIONAL = DiagnosticType_Error("JSC_DEFAULT_PARAM_MUST_BE_MARKED_OPTIONAL", "Inline JSDoc on default parameters must be marked as optional")
	static_CheckJSDoc_INVALID_NO_SIDE_EFFECT_ANNOTATION = DiagnosticType_Error("JSC_INVALID_NO_SIDE_EFFECT_ANNOTATION", "@nosideeffects may only appear in externs files.")
	static_CheckJSDoc_INVALID_MODIFIES_ANNOTATION = DiagnosticType_Error("JSC_INVALID_MODIFIES_ANNOTATION", "@modifies may only appear in externs files.")
	static_CheckJSDoc_INVALID_DEFINE_ON_LET = DiagnosticType_Error("JSC_INVALID_DEFINE_ON_LET", "variables annotated with @define may only be declared with VARs, ASSIGNs, or CONSTs")
	static_CheckJSDoc_MISPLACED_SUPPRESS = DiagnosticType_Warning("JSC_MISPLACED_SUPPRESS", "@suppress annotation not allowed here. See https://github.com/google/closure-compiler/wiki/@suppress-annotations")
	static_CheckJSDoc_JSDOC_IN_BLOCK_COMMENT = DiagnosticType_Warning("JSC_JSDOC_IN_BLOCK_COMMENT", "Non-JSDoc comment has annotations. Did you mean to start it with '/**'?")
	static_CheckJSDoc_JSDOC_ON_RETURN = DiagnosticType_Warning("JSC_JSDOC_ON_RETURN", "JSDoc annotations are not supported on return.")
	static_CheckJSDoc_COMMENT_PATTERN = j2gort.Ext_java_util_regex_Pattern_Compile_1("(/|(\n[ \t]*))\\*[ \t]*@[a-zA-Z]+[ \t\n{]")
}

func (this_ *CheckJSDoc) initFields_CheckJSDoc() {
	_ = this_
	this_.checkJsDocTypesVisitor = NewCheckJSDoc_CheckJsdocTypes()
}

// Java constructor: CheckJSDoc(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckJSDoc(compiler AbstractCompiler) *CheckJSDoc {
	clinit_CheckJSDoc()
	this_ := &CheckJSDoc{}
	InitCheckJSDoc(this_, compiler)
	return this_
}

func InitCheckJSDoc(this_ *CheckJSDoc, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_CheckJSDoc()
	this_.compiler = compiler
}

// Java method: CheckJSDoc.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckJSDoc) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	this_.inExterns = true
	NodeTraversal_Traverse(this_.compiler, externs, this_)
	this_.inExterns = false
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: CheckJSDoc.checkJsDocInBlockComments(java.lang.String fileName)
func (this_ *CheckJSDoc) priv_CheckJSDoc_checkJsDocInBlockComments_1055(fileName string) {
	_ = fileName
	if !(this_.compiler.GetOptions().PreservesDetailedSourceInfo()) {
		return
	}
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.compiler.GetComments(fileName)))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			comment := j2gort.CastPtr[pkg_parser.Comment](rawElem_2)
			_ = comment
			if comment.FieldGet_type_() == pkg_parser.StaticGet_Comment_Type_BLOCK() {
				if StaticGet_CheckJSDoc_COMMENT_PATTERN().Matcher(j2gort.BoxString(comment.FieldGet_value())).Find_0() {
					this_.compiler.Report_1(JSError_Make_5(fileName, (comment.FieldGet_location().FieldGet_start().FieldGet_line() + int32(1)), comment.FieldGet_location().FieldGet_start().FieldGet_column(), StaticGet_CheckJSDoc_JSDOC_IN_BLOCK_COMMENT(), []string{}))
				}
			}
		}
	}
}

// Java method: CheckJSDoc.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckJSDoc) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsScript() {
		this_.priv_CheckJSDoc_checkJsDocInBlockComments_1055(n.GetSourceFileName())
	}
	var info *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = info
	this_.priv_CheckJSDoc_validateTypeAnnotations_1068(n, info)
	this_.priv_CheckJSDoc_validateFunctionJsDoc_1065(n, info)
	this_.priv_CheckJSDoc_validateMsgJsDoc_1067(n, info)
	this_.priv_CheckJSDoc_validateNoCollapse_1064(n, info)
	this_.priv_CheckJSDoc_validateClassLevelJsDoc_1062(n, info)
	this_.priv_CheckJSDoc_validateArrowFunction_1071(n)
	this_.priv_CheckJSDoc_validateRestParameter_1072(n)
	this_.priv_CheckJSDoc_validateDefaultValue_1073(n)
	this_.priv_CheckJSDoc_validateTemplates_1058(n, info)
	this_.priv_CheckJSDoc_validateTypedefs_1057(n, info)
	this_.priv_CheckJSDoc_validateNoSideEffects_1074(n, info)
	this_.priv_CheckJSDoc_validateAbstractJsDoc_1063(n, info)
	this_.priv_CheckJSDoc_validateDefinesDeclaration_1075(n, info)
	this_.priv_CheckJSDoc_validateSuppress_1056(n, info)
	this_.priv_CheckJSDoc_validateImplicitCast_1076(n, info)
	this_.priv_CheckJSDoc_validateClosurePrimitive_1077(n, info)
	this_.priv_CheckJSDoc_validateReturnJsDoc_1078(n, info)
	this_.priv_CheckJSDoc_validateTsType_1079(n, info)
	this_.priv_CheckJSDoc_validateJsDocTypeNames_1080(info)
	this_.priv_CheckJSDoc_validateIsUsedViaDotConstructor_1081(n, info)
	this_.priv_CheckJSDoc_validateLegacyNamedTsExport_1082(n, info)
}

// Java method: CheckJSDoc.validateSuppress(com.google.javascript.rhino.Node n, com.google.javascript.rhino.JSDocInfo info)
func (this_ *CheckJSDoc) priv_CheckJSDoc_validateSuppress_1056(n pkg_rhino.Node, info *pkg_rhino.JSDocInfo) {
	_ = n
	_ = info
	if (info == nil) || info.GetSuppressions().IsEmpty() {
		return
	}
	switch n.TokenSwitchOrd_() {
	case int32(73) /* Token.FUNCTION */, int32(105) /* Token.CLASS */, int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */, int32(94) /* Token.SCRIPT */, int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */, int32(108) /* Token.MEMBER_FIELD_DEF */, int32(109) /* Token.COMPUTED_FIELD_DEF */ :
		return
	case int32(129) /* Token.COMPUTED_PROP */ :
		if n.GetLastChild().IsFunction() {
			return
		}
	case int32(100) /* Token.STRING_KEY */ :
		if n.GetParent().IsObjectLit() {
			return
		}
	case int32(85) /* Token.WITH */ :
		if CheckJSDoc_ContainsOnlySuppressionFor(info, "with") {
			return
		}
	default:
	}
	if CheckJSDoc_ContainsOnlySuppressionFor(info, "missingRequire") {
		return
	}
	if n.GetParent().IsExprResult() {
		return
	}
	this_.compiler.Report_1(JSError_Make_3(n, StaticGet_CheckJSDoc_MISPLACED_SUPPRESS(), []string{}))
}

// Java static method: CheckJSDoc.containsOnlySuppressionFor(com.google.javascript.rhino.JSDocInfo jsdoc, java.lang.String allowedSuppression)
func CheckJSDoc_ContainsOnlySuppressionFor(jsdoc *pkg_rhino.JSDocInfo, allowedSuppression string) bool {
	clinit_CheckJSDoc()
	_ = jsdoc
	_ = allowedSuppression
	var suppressions j2gort.Ext_java_util_Set = jsdoc.GetSuppressions()
	_ = suppressions
	return ((suppressions.Size() == int32(1)) && (j2gort.CheckNotNullString(j2gort.CastString(j2gort.Ext_com_google_common_collect_Iterables_GetOnlyElement_1(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](suppressions)))) == allowedSuppression))
}

// Java method: CheckJSDoc.validateTypedefs(com.google.javascript.rhino.Node n, com.google.javascript.rhino.JSDocInfo info)
func (this_ *CheckJSDoc) priv_CheckJSDoc_validateTypedefs_1057(n pkg_rhino.Node, info *pkg_rhino.JSDocInfo) {
	_ = n
	_ = info
	if (info == nil) || !(info.HasTypedefType()) {
		return
	}
	if this_.priv_CheckJSDoc_isClassDecl_1059(n) {
		this_.priv_CheckJSDoc_reportMisplaced_1069(n, "typedef", "@typedef is not allowed on a class declaration.")
		return
	}
	var lvalue pkg_rhino.Node
	_ = lvalue
	if NodeUtil_IsNameDeclaration(n) || n.IsAssign() {
		lvalue = n.GetFirstChild()
	} else {
		lvalue = n
	}
	if !(lvalue.IsQualifiedName()) {
		this_.priv_CheckJSDoc_reportMisplaced_1069(n, "typedef", "@typedef is only allowed on qualified name declarations. Did you mean @type?")
	} else if CheckJSDoc_IsPrototypeOrInstanceDecl(lvalue) {
		this_.priv_CheckJSDoc_reportMisplaced_1069(n, "typedef", "@typedef is not allowed on instance or prototype properties. Did you mean @type?")
	}
}

// Java method: CheckJSDoc.validateTemplates(com.google.javascript.rhino.Node n, com.google.javascript.rhino.JSDocInfo info)
func (this_ *CheckJSDoc) priv_CheckJSDoc_validateTemplates_1058(n pkg_rhino.Node, info *pkg_rhino.JSDocInfo) {
	_ = n
	_ = info
	if (((((info != nil) && !(info.GetTemplateTypeNames().IsEmpty())) && !(info.IsConstructorOrInterface())) && !(this_.priv_CheckJSDoc_isClassDecl_1059(n))) && !(info.ContainsFunctionDeclaration())) && (CheckJSDoc_GetFunctionDecl(n) == 0) {
		this_.priv_CheckJSDoc_reportMisplaced_1069(n, "template", "@template is only allowed in class, constructor, interface, function or method declarations")
	}
}

// Java static method: CheckJSDoc.getFunctionDecl(com.google.javascript.rhino.Node n)
func CheckJSDoc_GetFunctionDecl(n pkg_rhino.Node) pkg_rhino.Node {
	clinit_CheckJSDoc()
	_ = n
	if n.IsFunction() {
		return n
	}
	if n.IsMemberFunctionDef() {
		return n.GetFirstChild()
	}
	if (NodeUtil_IsNameDeclaration(n) && (n.GetFirstFirstChild() != 0)) && n.GetFirstFirstChild().IsFunction() {
		return n.GetFirstFirstChild()
	}
	if (n.IsAssign() && n.GetFirstChild().IsQualifiedName()) && n.GetLastChild().IsFunction() {
		return n.GetLastChild()
	}
	if n.IsGetterDef() || n.IsSetterDef() {
		return n.GetFirstChild()
	}
	if n.IsComputedProp() && n.GetLastChild().IsFunction() {
		return n.GetLastChild()
	}
	return pkg_rhino.Node(0)
}

// Java method: CheckJSDoc.isClassDecl(com.google.javascript.rhino.Node n)
func (this_ *CheckJSDoc) priv_CheckJSDoc_isClassDecl_1059(n pkg_rhino.Node) bool {
	_ = n
	return (((this_.priv_CheckJSDoc_isClass_1061(n) || (n.IsAssign() && this_.priv_CheckJSDoc_isClass_1061(n.GetLastChild()))) || (NodeUtil_IsNameDeclaration(n) && this_.priv_CheckJSDoc_isNameInitializeWithClass_1060(n.GetFirstChild()))) || this_.priv_CheckJSDoc_isNameInitializeWithClass_1060(n))
}

// Java method: CheckJSDoc.isNameInitializeWithClass(com.google.javascript.rhino.Node n)
func (this_ *CheckJSDoc) priv_CheckJSDoc_isNameInitializeWithClass_1060(n pkg_rhino.Node) bool {
	_ = n
	return ((((n != 0) && n.IsName()) && n.HasChildren()) && this_.priv_CheckJSDoc_isClassDecl_1059(n.GetFirstChild()))
}

// Java method: CheckJSDoc.isClass(com.google.javascript.rhino.Node n)
func (this_ *CheckJSDoc) priv_CheckJSDoc_isClass_1061(n pkg_rhino.Node) bool {
	_ = n
	return (n.IsClass() || (n.IsCall() && this_.compiler.GetCodingConvention().IsClassFactoryCall(n)))
}

// Java static method: CheckJSDoc.isPrototypeOrInstanceDecl(com.google.javascript.rhino.Node n)
func CheckJSDoc_IsPrototypeOrInstanceDecl(n pkg_rhino.Node) bool {
	clinit_CheckJSDoc()
	_ = n
	if n.IsStringKey() {
		return false
	}
	if NodeUtil_IsPrototypeProperty(n) {
		return true
	}
	var receiver pkg_rhino.Node = NodeUtil_GetRootOfQualifiedName_1_Node(n)
	_ = receiver
	return (receiver.IsThis() || receiver.IsSuper())
}

// Java method: CheckJSDoc.validateClassLevelJsDoc(com.google.javascript.rhino.Node n, com.google.javascript.rhino.JSDocInfo info)
func (this_ *CheckJSDoc) priv_CheckJSDoc_validateClassLevelJsDoc_1062(n pkg_rhino.Node, info *pkg_rhino.JSDocInfo) {
	_ = n
	_ = info
	if ((info != nil) && n.IsMemberFunctionDef()) && CheckJSDoc_HasClassLevelJsDoc(info) {
		this_.priv_CheckJSDoc_report_1070(n, StaticGet_CheckJSDoc_DISALLOWED_MEMBER_JSDOC(), []string{})
	}
}

// Java method: CheckJSDoc.validateAbstractJsDoc(com.google.javascript.rhino.Node n, com.google.javascript.rhino.JSDocInfo info)
func (this_ *CheckJSDoc) priv_CheckJSDoc_validateAbstractJsDoc_1063(n pkg_rhino.Node, info *pkg_rhino.JSDocInfo) {
	_ = n
	_ = info
	if (info == nil) || !(info.IsAbstract()) {
		return
	}
	if this_.priv_CheckJSDoc_isClassDecl_1059(n) {
		return
	}
	if n.IsGetProp() {
		if !(n.GetFirstChild().IsThis() && (info.GetReturnType() != (pkg_rhino.JSTypeExpression{}))) {
			this_.priv_CheckJSDoc_report_1070(n, StaticGet_CheckJSDoc_MISPLACED_ANNOTATION(), []string{"@abstract", "abstract undeclared methods can only be written as class fields"})
		} else if !((n.GetParent().IsExprResult() && n.GetGrandparent().IsBlock()) && n.GetParent().GetGrandparent().IsFunction()) {
			this_.priv_CheckJSDoc_report_1070(n, StaticGet_CheckJSDoc_MISPLACED_ANNOTATION(), []string{"@abstract", "abstract methods without an initializer must be declared in a constructor function."})
			return
		}
		return
	}
	var functionNode pkg_rhino.Node = CheckJSDoc_GetFunctionDecl(n)
	_ = functionNode
	if functionNode == 0 {
		this_.priv_CheckJSDoc_report_1070(n, StaticGet_CheckJSDoc_MISPLACED_ANNOTATION(), []string{"@abstract", "only functions or non-static methods can be abstract"})
		return
	}
	if !(info.IsConstructor()) && NodeUtil_GetFunctionBody(functionNode).HasChildren() {
		this_.priv_CheckJSDoc_report_1070(n, StaticGet_CheckJSDoc_MISPLACED_ANNOTATION(), []string{"@abstract", "function with a non-empty body cannot be abstract"})
		return
	}
	if NodeUtil_IsEs6ConstructorMemberFunctionDef(n) {
		this_.priv_CheckJSDoc_report_1070(n, StaticGet_CheckJSDoc_MISPLACED_ANNOTATION(), []string{"@abstract", "constructors cannot be abstract"})
		return
	}
	if (((((!(info.IsConstructor()) && !(n.IsMemberFunctionDef())) && !(n.IsStringKey())) && !(n.IsComputedProp())) && !(n.IsGetterDef())) && !(n.IsSetterDef())) && !(NodeUtil_IsPrototypeMethod(functionNode)) {
		this_.priv_CheckJSDoc_report_1070(n, StaticGet_CheckJSDoc_MISPLACED_ANNOTATION(), []string{"@abstract", "only functions or non-static methods can be abstract"})
		return
	}
	if n.IsStaticMember() {
		this_.priv_CheckJSDoc_report_1070(n, StaticGet_CheckJSDoc_MISPLACED_ANNOTATION(), []string{"@abstract", "static methods cannot be abstract"})
		return
	}
}

// Java static method: CheckJSDoc.hasClassLevelJsDoc(com.google.javascript.rhino.JSDocInfo info)
func CheckJSDoc_HasClassLevelJsDoc(info *pkg_rhino.JSDocInfo) bool {
	clinit_CheckJSDoc()
	_ = info
	return (((info.IsConstructorOrInterface() || info.HasBaseType()) || (info.GetImplementedInterfaceCount() != int32(0))) || (info.GetExtendedInterfacesCount() != int32(0)))
}

// Java method: CheckJSDoc.validateNoCollapse(com.google.javascript.rhino.Node n, com.google.javascript.rhino.JSDocInfo info)
func (this_ *CheckJSDoc) priv_CheckJSDoc_validateNoCollapse_1064(n pkg_rhino.Node, info *pkg_rhino.JSDocInfo) {
	_ = n
	_ = info
	if (info == nil) || !(info.IsNoCollapse()) {
		return
	}
	if n.IsFromExterns() {
		this_.priv_CheckJSDoc_reportMisplaced_1069(n, "nocollapse", "This JSDoc has no effect in externs.")
		return
	}
	if NodeUtil_IsPrototypePropertyDeclaration(n.GetParent()) || (n.GetParent().IsClassMembers() && !(n.IsStaticMember())) {
		this_.priv_CheckJSDoc_reportMisplaced_1069(n, "nocollapse", "This JSDoc has no effect on prototype properties and non-static fields.")
	}
	if n.IsAssign() {
		var assignee pkg_rhino.Node = n.GetFirstChild()
		_ = assignee
		if assignee.IsQualifiedName() {
			var rootOfQname pkg_rhino.Node = NodeUtil_GetRootOfQualifiedName_1_Node(assignee)
			_ = rootOfQname
			if !(rootOfQname.IsName()) {
				this_.priv_CheckJSDoc_reportMisplaced_1069(n, "nocollapse", "This JSDoc has no effect.")
			}
		}
	}
}

// Java method: CheckJSDoc.validateFunctionJsDoc(com.google.javascript.rhino.Node n, com.google.javascript.rhino.JSDocInfo info)
func (this_ *CheckJSDoc) priv_CheckJSDoc_validateFunctionJsDoc_1065(n pkg_rhino.Node, info *pkg_rhino.JSDocInfo) {
	_ = n
	_ = info
	if info == nil {
		return
	}
	if (info.ContainsFunctionDeclaration() && !(info.HasType())) && !(this_.priv_CheckJSDoc_isJSDocOnFunctionNode_1066(n, info)) {
		this_.priv_CheckJSDoc_reportMisplaced_1069(n, "function", "This JSDoc is not attached to a function node. Are you missing parentheses?")
	}
}

// Java method: CheckJSDoc.isJSDocOnFunctionNode(com.google.javascript.rhino.Node n, com.google.javascript.rhino.JSDocInfo info)
func (this_ *CheckJSDoc) priv_CheckJSDoc_isJSDocOnFunctionNode_1066(n pkg_rhino.Node, info *pkg_rhino.JSDocInfo) bool {
	_ = n
	_ = info
	switch n.TokenSwitchOrd_() {
	case int32(73) /* Token.FUNCTION */, int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */, int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(100) /* Token.STRING_KEY */, int32(129) /* Token.COMPUTED_PROP */, int32(120) /* Token.EXPORT */, int32(108) /* Token.MEMBER_FIELD_DEF */, int32(109) /* Token.COMPUTED_FIELD_DEF */ :
		return true
	case int32(27) /* Token.GETELEM */, int32(26) /* Token.GETPROP */ :
		if n.GetFirstChild().IsQualifiedName() {
			return true
		}
		return false
	case int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */, int32(51) /* Token.ASSIGN */ :
		var lhs pkg_rhino.Node = n.GetFirstChild()
		_ = lhs
		var rhs pkg_rhino.Node = NodeUtil_GetRValueOfLValue(lhs)
		_ = rhs
		if ((rhs != 0) && this_.priv_CheckJSDoc_isClass_1061(rhs)) && !(info.IsConstructor()) {
			return false
		}
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: CheckJSDoc.validateMsgJsDoc(com.google.javascript.rhino.Node n, com.google.javascript.rhino.JSDocInfo info)
func (this_ *CheckJSDoc) priv_CheckJSDoc_validateMsgJsDoc_1067(n pkg_rhino.Node, info *pkg_rhino.JSDocInfo) {
	_ = n
	_ = info
	if info == nil {
		return
	}
	var hasNonDescMsgTag bool = ((info.GetMeaning() != j2gort.NullString) || (info.GetAlternateMessageId() != j2gort.NullString))
	_ = hasNonDescMsgTag
	if hasNonDescMsgTag || ((info.GetDescription() != j2gort.NullString) && !(CheckJSDoc_IsFromTs(n))) {
		var descOkay bool = false
		_ = descOkay
		switch n.TokenSwitchOrd_() {
		case int32(51) /* Token.ASSIGN */, int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */ :
			descOkay = CheckJSDoc_IsValidMsgName(n.GetFirstChild())
		case int32(100) /* Token.STRING_KEY */ :
			descOkay = CheckJSDoc_IsValidMsgName(n)
		case int32(26) /* Token.GETPROP */ :
			if n.IsFromExterns() && n.IsQualifiedName() {
				descOkay = CheckJSDoc_IsValidMsgName(n)
			}
		default:
		}
		if !(descOkay) {
			this_.priv_CheckJSDoc_report_1070(n, StaticGet_CheckJSDoc_MISPLACED_MSG_ANNOTATION(), []string{})
		}
	}
}

// Java static method: CheckJSDoc.isValidMsgName(com.google.javascript.rhino.Node nameNode)
func CheckJSDoc_IsValidMsgName(nameNode pkg_rhino.Node) bool {
	clinit_CheckJSDoc()
	_ = nameNode
	if nameNode.IsName() || nameNode.IsStringKey() {
		return javalang.StringStartsWith(j2gort.CheckNotNullString(nameNode.GetString()), "MSG_")
	} else if nameNode.IsQualifiedName() {
		return javalang.StringStartsWith(j2gort.CheckNotNullString(nameNode.GetString()), "MSG_")
	} else {
		return false
	}
}

// Java method: CheckJSDoc.validateTypeAnnotations(com.google.javascript.rhino.Node n, com.google.javascript.rhino.JSDocInfo info)
func (this_ *CheckJSDoc) priv_CheckJSDoc_validateTypeAnnotations_1068(n pkg_rhino.Node, info *pkg_rhino.JSDocInfo) {
	_ = n
	_ = info
	if (info != nil) && info.HasType() {
		var valid bool = false
		_ = valid
		switch n.TokenSwitchOrd_() {
		case int32(73) /* Token.FUNCTION */ :
			valid = NodeUtil_IsFunctionDeclaration(n)
		case int32(32) /* Token.NAME */ :
			valid = CheckJSDoc_IsTypeAnnotationAllowedForName(n)
		case int32(102) /* Token.ARRAY_PATTERN */, int32(103) /* Token.OBJECT_PATTERN */ :
			valid = n.GetParent().IsParamList()
		case int32(101) /* Token.CAST */, int32(120) /* Token.EXPORT */, int32(100) /* Token.STRING_KEY */, int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */, int32(108) /* Token.MEMBER_FIELD_DEF */, int32(109) /* Token.COMPUTED_FIELD_DEF */ :
			valid = true
		case int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */ :
			valid = !(NodeUtil_IsDestructuringDeclaration(n))
		case int32(51) /* Token.ASSIGN */ :
			var lvalue pkg_rhino.Node = n.GetFirstChild()
			_ = lvalue
			valid = (n.GetParent().IsExprResult() && ((lvalue.IsGetProp() || lvalue.IsGetElem()) || lvalue.MatchesName_1_string("exports")))
		case int32(26) /* Token.GETPROP */ :
			valid = (n.GetParent().IsExprResult() && n.IsQualifiedName())
		case int32(28) /* Token.CALL */ :
			valid = info.IsDefine()
		default:
		}
		if !(valid) {
			this_.priv_CheckJSDoc_reportMisplaced_1069(n, "type", "Type annotations are not allowed here. Are you missing parentheses?")
		}
	}
}

// Java static method: CheckJSDoc.isTypeAnnotationAllowedForName(com.google.javascript.rhino.Node n)
func CheckJSDoc_IsTypeAnnotationAllowedForName(n pkg_rhino.Node) bool {
	clinit_CheckJSDoc()
	_ = n
	if !(n.IsName()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	if !(NodeUtil_IsLValue(n)) {
		return false
	}
	var rootTarget pkg_rhino.Node = NodeUtil_GetRootTarget(n)
	_ = rootTarget
	return !(NodeUtil_IsLhsOfAssign(rootTarget))
}

// Java method: CheckJSDoc.reportMisplaced(com.google.javascript.rhino.Node n, java.lang.String annotationName, java.lang.String note)
func (this_ *CheckJSDoc) priv_CheckJSDoc_reportMisplaced_1069(n pkg_rhino.Node, annotationName string, note string) {
	_ = n
	_ = annotationName
	_ = note
	this_.compiler.Report_1(JSError_Make_3(n, StaticGet_CheckJSDoc_MISPLACED_ANNOTATION(), []string{annotationName, note}))
}

// Java method: CheckJSDoc.report(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.DiagnosticType type, java.lang.String[] arguments)
func (this_ *CheckJSDoc) priv_CheckJSDoc_report_1070(n pkg_rhino.Node, type_ *DiagnosticType, arguments []string) {
	_ = n
	_ = type_
	_ = arguments
	this_.compiler.Report_1(JSError_Make_3(n, type_, arguments))
}

// Java method: CheckJSDoc.validateArrowFunction(com.google.javascript.rhino.Node n)
func (this_ *CheckJSDoc) priv_CheckJSDoc_validateArrowFunction_1071(n pkg_rhino.Node) {
	_ = n
	if n.IsArrowFunction() {
		var info *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(n)
		_ = info
		if (info != nil) && info.IsConstructorOrInterface() {
			this_.priv_CheckJSDoc_report_1070(n, StaticGet_CheckJSDoc_ARROW_FUNCTION_AS_CONSTRUCTOR(), []string{})
		}
	}
}

// Java method: CheckJSDoc.validateRestParameter(com.google.javascript.rhino.Node restParam)
func (this_ *CheckJSDoc) priv_CheckJSDoc_validateRestParameter_1072(restParam pkg_rhino.Node) {
	_ = restParam
	if !(restParam.IsRest()) || !(restParam.GetParent().IsParamList()) {
		return
	}
	var paramList pkg_rhino.Node = restParam.GetParent()
	_ = paramList
	var inlineInfo *pkg_rhino.JSDocInfo = restParam.GetFirstChild().GetJSDocInfo()
	_ = inlineInfo
	var functionInfo *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(paramList.GetParent())
	_ = functionInfo
	var paramTypeAnnotation pkg_rhino.JSTypeExpression = (pkg_rhino.JSTypeExpression{})
	_ = paramTypeAnnotation
	if inlineInfo != nil {
		paramTypeAnnotation = inlineInfo.GetType()
	} else if functionInfo != nil {
		if restParam.GetFirstChild().IsName() {
			var paramName string = restParam.GetFirstChild().GetString()
			_ = paramName
			paramTypeAnnotation = functionInfo.GetParameterType(paramName)
		} else {
			var indexOfRest int32 = paramList.GetIndexOfChild(restParam)
			_ = indexOfRest
			if functionInfo.GetParameterCount() >= indexOfRest {
				paramTypeAnnotation = functionInfo.GetParameterType(functionInfo.GetParameterNameAt(indexOfRest))
			} else {
				paramTypeAnnotation = (pkg_rhino.JSTypeExpression{})
			}
		}
	} else {
		paramTypeAnnotation = (pkg_rhino.JSTypeExpression{})
	}
	if (paramTypeAnnotation != (pkg_rhino.JSTypeExpression{})) && (paramTypeAnnotation.GetRoot().RawToken_() != uint8(125)) {
		this_.compiler.Report_1(JSError_Make_3(restParam, StaticGet_CheckJSDoc_BAD_REST_PARAMETER_ANNOTATION(), []string{}))
	}
}

// Java method: CheckJSDoc.validateDefaultValue(com.google.javascript.rhino.Node n)
func (this_ *CheckJSDoc) priv_CheckJSDoc_validateDefaultValue_1073(n pkg_rhino.Node) {
	_ = n
	if n.IsDefaultValue() && n.GetParent().IsParamList() {
		var targetNode pkg_rhino.Node = n.GetFirstChild()
		_ = targetNode
		var info *pkg_rhino.JSDocInfo = targetNode.GetJSDocInfo()
		_ = info
		if info == nil {
			return
		}
		var typeExpr pkg_rhino.JSTypeExpression = info.GetType()
		_ = typeExpr
		if typeExpr == (pkg_rhino.JSTypeExpression{}) {
			return
		}
		var typeNode pkg_rhino.Node = typeExpr.GetRoot()
		_ = typeNode
		if typeNode.RawToken_() != uint8(160) {
			this_.priv_CheckJSDoc_report_1070(typeNode, StaticGet_CheckJSDoc_DEFAULT_PARAM_MUST_BE_MARKED_OPTIONAL(), []string{})
		}
	}
}

// Java method: CheckJSDoc.validateNoSideEffects(com.google.javascript.rhino.Node n, com.google.javascript.rhino.JSDocInfo info)
func (this_ *CheckJSDoc) priv_CheckJSDoc_validateNoSideEffects_1074(n pkg_rhino.Node, info *pkg_rhino.JSDocInfo) {
	_ = n
	_ = info
	if info == nil {
		return
	}
	if n.IsFromExterns() {
		return
	}
	if info.HasSideEffectsArgumentsAnnotation() || info.ModifiesThis() {
		this_.priv_CheckJSDoc_report_1070(n, StaticGet_CheckJSDoc_INVALID_MODIFIES_ANNOTATION(), []string{})
	}
}

// Java method: CheckJSDoc.validateDefinesDeclaration(com.google.javascript.rhino.Node n, com.google.javascript.rhino.JSDocInfo info)
func (this_ *CheckJSDoc) priv_CheckJSDoc_validateDefinesDeclaration_1075(n pkg_rhino.Node, info *pkg_rhino.JSDocInfo) {
	_ = n
	_ = info
	if ((info != nil) && info.IsDefine()) && n.IsLet() {
		this_.priv_CheckJSDoc_report_1070(n, StaticGet_CheckJSDoc_INVALID_DEFINE_ON_LET(), []string{})
	}
}

// Java method: CheckJSDoc.validateImplicitCast(com.google.javascript.rhino.Node n, com.google.javascript.rhino.JSDocInfo info)
func (this_ *CheckJSDoc) priv_CheckJSDoc_validateImplicitCast_1076(n pkg_rhino.Node, info *pkg_rhino.JSDocInfo) {
	_ = n
	_ = info
	if (!(this_.inExterns) && (info != nil)) && info.IsImplicitCast() {
		this_.priv_CheckJSDoc_report_1070(n, StaticGet_TypeCheck_ILLEGAL_IMPLICIT_CAST(), []string{})
	}
}

// Java method: CheckJSDoc.validateClosurePrimitive(com.google.javascript.rhino.Node n, com.google.javascript.rhino.JSDocInfo info)
func (this_ *CheckJSDoc) priv_CheckJSDoc_validateClosurePrimitive_1077(n pkg_rhino.Node, info *pkg_rhino.JSDocInfo) {
	_ = n
	_ = info
	if (info == nil) || !(info.HasClosurePrimitiveId()) {
		return
	}
	if !(this_.priv_CheckJSDoc_isJSDocOnFunctionNode_1066(n, info)) {
		this_.priv_CheckJSDoc_report_1070(n, StaticGet_CheckJSDoc_MISPLACED_ANNOTATION(), []string{"closurePrimitive", "must be on a function node"})
	}
}

// Java method: CheckJSDoc.validateReturnJsDoc(com.google.javascript.rhino.Node n, com.google.javascript.rhino.JSDocInfo info)
func (this_ *CheckJSDoc) priv_CheckJSDoc_validateReturnJsDoc_1078(n pkg_rhino.Node, info *pkg_rhino.JSDocInfo) {
	_ = n
	_ = info
	if !(n.IsReturn()) || (info == nil) {
		return
	}
	if (info.ContainsDeclaration() && !(info.HasType())) && !(info.HasTypedefType()) {
		this_.priv_CheckJSDoc_report_1070(n, StaticGet_CheckJSDoc_JSDOC_ON_RETURN(), []string{})
	}
}

// Java method: CheckJSDoc.validateTsType(com.google.javascript.rhino.Node n, com.google.javascript.rhino.JSDocInfo info)
func (this_ *CheckJSDoc) priv_CheckJSDoc_validateTsType_1079(n pkg_rhino.Node, info *pkg_rhino.JSDocInfo) {
	_ = n
	_ = info
	if (info == nil) || info.GetTsTypes().IsEmpty() {
		return
	}
	if !(this_.priv_CheckJSDoc_isJSDocOnFunctionNode_1066(n, info)) {
		this_.priv_CheckJSDoc_report_1070(n, StaticGet_CheckJSDoc_MISPLACED_TS_TYPE_ANNOTATION(), []string{"must be on a function node"})
	}
}

// Java static method: CheckJSDoc.isFromTs(com.google.javascript.rhino.Node n)
func CheckJSDoc_IsFromTs(n pkg_rhino.Node) bool {
	clinit_CheckJSDoc()
	_ = n
	return n.GetStaticSourceFile().IsTypeScriptSource()
}

// Java method: CheckJSDoc.validateJsDocTypeNames(com.google.javascript.rhino.JSDocInfo info)
func (this_ *CheckJSDoc) priv_CheckJSDoc_validateJsDocTypeNames_1080(info *pkg_rhino.JSDocInfo) {
	_ = info
	if info == nil {
		return
	}
	{
		var i int32 = int32(0)
		_ = i
		var count int32 = info.GetTypeNodeCount()
		_ = count
		for ; i < count; i = i + 1 {
			NodeUtil_VisitPreOrder_2(info.GetTypeNodeAt(i), ToIface_NodeUtil_Visitor(this_.checkJsDocTypesVisitor))
		}
	}
}

// Java method: CheckJSDoc.validateIsUsedViaDotConstructor(com.google.javascript.rhino.Node n, com.google.javascript.rhino.JSDocInfo info)
func (this_ *CheckJSDoc) priv_CheckJSDoc_validateIsUsedViaDotConstructor_1081(n pkg_rhino.Node, info *pkg_rhino.JSDocInfo) {
	_ = n
	_ = info
	if (info == nil) || !(info.IsUsedViaDotConstructor()) {
		return
	}
	if !(n.IsFunction() && info.IsConstructor()) && !(NodeUtil_IsEs6ConstructorMemberFunctionDef(n)) {
		this_.priv_CheckJSDoc_report_1070(n, StaticGet_CheckJSDoc_MISPLACED_ANNOTATION(), []string{"usedViaDotConstructor", "must be on a constructor"})
	}
}

// Java method: CheckJSDoc.validateLegacyNamedTsExport(com.google.javascript.rhino.Node n, com.google.javascript.rhino.JSDocInfo info)
func (this_ *CheckJSDoc) priv_CheckJSDoc_validateLegacyNamedTsExport_1082(n pkg_rhino.Node, info *pkg_rhino.JSDocInfo) {
	_ = n
	_ = info
	if (info == nil) || !(info.IsLegacyNamedTsExport()) {
		return
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if !(n.IsCall()) || !(NodeUtil_IsGoogModuleCall(parent) || NodeUtil_IsGoogProvideCall(parent)) {
		this_.priv_CheckJSDoc_report_1070(n, StaticGet_CheckJSDoc_MISPLACED_ANNOTATION(), []string{"@legacyNamedTsExport", "must be on a goog.module or goog.provide call"})
	}
}

// Java class: com.google.javascript.jscomp.CheckJSDoc.CheckJsdocTypes
type CheckJSDoc_CheckJsdocTypes struct {
	idHash_ atomic.Int32
}

func (this_ *CheckJSDoc_CheckJsdocTypes) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckJSDoc_CheckJsdocTypes) IsJavaNil_() bool { return this_ == nil }
func (*CheckJSDoc_CheckJsdocTypes) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckJSDoc$CheckJsdocTypes", "CheckJsdocTypes"
}
func (this_ *CheckJSDoc_CheckJsdocTypes) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CheckJSDoc_CheckJsdocTypes) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckJSDoc.CheckJsdocTypes", "com.google.javascript.jscomp.CheckJSDoc$CheckJsdocTypes", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckJSDoc_CheckJsdocTypes](v_)
	})
}
func (*CheckJSDoc_CheckJsdocTypes) Implements_jscomp_NodeUtil_Visitor() {}

func (this_ *CheckJSDoc_CheckJsdocTypes) initFields_CheckJSDoc_CheckJsdocTypes() {
	_ = this_
}

// Java constructor: CheckJsdocTypes()
func NewCheckJSDoc_CheckJsdocTypes() *CheckJSDoc_CheckJsdocTypes {
	this_ := &CheckJSDoc_CheckJsdocTypes{}
	InitCheckJSDoc_CheckJsdocTypes(this_)
	return this_
}

func InitCheckJSDoc_CheckJsdocTypes(this_ *CheckJSDoc_CheckJsdocTypes) {
	this_.initFields_CheckJSDoc_CheckJsdocTypes()
}

// Java method: CheckJsdocTypes.visit(com.google.javascript.rhino.Node typeRefNode)
func (this_ *CheckJSDoc_CheckJsdocTypes) Visit(typeRefNode pkg_rhino.Node) {
	_ = typeRefNode
	if !(typeRefNode.IsStringLit()) {
		return
	}
	var typeName string = typeRefNode.GetString()
	_ = typeName
	var dot int32 = j2gort.JStrIndexOfCp(typeName, int32(uint16(46)))
	_ = dot
	var rootOfType string
	_ = rootOfType
	if dot == int32(-1) {
		rootOfType = typeName
	} else {
		rootOfType = j2gort.JStrSub(typeName, int32(0), dot)
	}
	if ClosureRewriteModule_IsModuleExport(rootOfType) || ClosureRewriteModule_IsModuleContent(rootOfType) {
		typeRefNode.SetString(("UnrecognizedType_" + j2gort.StrValString(typeName)))
	}
}
func (*CheckJSDoc) AlwaysShouldTraverse_()      {}
func (*CheckJSDoc) CanReleaseSyntacticScopes_() {}
