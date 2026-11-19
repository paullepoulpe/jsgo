/*
 * Copyright 2011 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/AstValidator.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.AstValidator
type AstValidator struct {
	compiler                         AbstractCompiler              `gson:"compiler"`
	violationHandler                 AstValidator_ViolationHandler `gson:"violationHandler"`
	currentScript                    pkg_rhino.Node                `gson:"currentScript"`
	idHash_                          atomic.Int32
	typeValidationMode               AstValidator_TypeInfoValidation `gson:"typeValidationMode"`
	isScriptFeatureValidationEnabled bool                            `gson:"isScriptFeatureValidationEnabled"`
	shouldValidateRequiredInlinings  bool                            `gson:"shouldValidateRequiredInlinings"`
}

func (this_ *AstValidator) BindSelf_(s any) {
	_ = s
}
func (this_ *AstValidator) IsJavaNil_() bool { return this_ == nil }
func (*AstValidator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AstValidator", "AstValidator"
}
func (this_ *AstValidator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AstValidator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AstValidator", "com.google.javascript.jscomp.AstValidator", func(v_ any) bool {
		return j2gort.InstanceOf[*AstValidator](v_)
	})
}
func (*AstValidator) Implements_jscomp_CompilerPass() {}

func StaticGet_AstValidator_AWAIT_NOT_WITHIN_ASYNC_FUNCTION() string {
	return "'await' expression is not within an async function"
}
func StaticGet_AstValidator_AWAIT_NOT_ALLOWED_IN_PARAMETER_LIST() string {
	return "'await' expression is not allowed in a parameter list"
}
func (this_ *AstValidator) initFields_AstValidator() {
	_ = this_
	this_.typeValidationMode = StaticGet_AstValidator_TypeInfoValidation_NONE()
}

// Java constructor: AstValidator(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.AstValidator.ViolationHandler handler, boolean validateScriptFeatures, boolean shouldValidateRequiredInlinings)
func NewAstValidator_4(compiler AbstractCompiler, handler AstValidator_ViolationHandler, validateScriptFeatures bool, shouldValidateRequiredInlinings bool) *AstValidator {
	this_ := &AstValidator{}
	InitAstValidator_4(this_, compiler, handler, validateScriptFeatures, shouldValidateRequiredInlinings)
	return this_
}

func InitAstValidator_4(this_ *AstValidator, compiler AbstractCompiler, handler AstValidator_ViolationHandler, validateScriptFeatures bool, shouldValidateRequiredInlinings bool) {
	_ = compiler
	_ = handler
	_ = validateScriptFeatures
	_ = shouldValidateRequiredInlinings
	this_.initFields_AstValidator()
	this_.compiler = compiler
	this_.violationHandler = handler
	this_.isScriptFeatureValidationEnabled = validateScriptFeatures
	this_.shouldValidateRequiredInlinings = shouldValidateRequiredInlinings
}

// Java constructor: AstValidator(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewAstValidator_1(compiler AbstractCompiler) *AstValidator {
	this_ := &AstValidator{}
	InitAstValidator_1(this_, compiler)
	return this_
}

func InitAstValidator_1(this_ *AstValidator, compiler AbstractCompiler) {
	_ = compiler
	InitAstValidator_2(this_, compiler, false)
}

// Java constructor: AstValidator(com.google.javascript.jscomp.AbstractCompiler compiler, boolean validateScriptFeatures)
func NewAstValidator_2(compiler AbstractCompiler, validateScriptFeatures bool) *AstValidator {
	this_ := &AstValidator{}
	InitAstValidator_2(this_, compiler, validateScriptFeatures)
	return this_
}

func InitAstValidator_2(this_ *AstValidator, compiler AbstractCompiler, validateScriptFeatures bool) {
	_ = compiler
	_ = validateScriptFeatures
	InitAstValidator_4(this_, compiler, NewAstValidator__Anon_2(this_), validateScriptFeatures, func() bool {
		re_1 := compiler.GetOptions().GetShouldValidateRequiredInlinings()
		ae_2 := pkg_base.StaticGet_Tri_TRUE()
		if re_1 == 0 {
			j2gort.ThrowNPE()
		}
		return re_1 == ae_2
	}())
}

// Java method: AstValidator.setTypeValidationMode(com.google.javascript.jscomp.AstValidator.TypeInfoValidation mode)
func (this_ *AstValidator) SetTypeValidationMode(mode AstValidator_TypeInfoValidation) *AstValidator {
	_ = mode
	this_.typeValidationMode = mode
	return this_
}

// Java method: AstValidator.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *AstValidator) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if externs != 0 {
		this_.ValidateCodeRoot(externs)
	}
	if root != 0 {
		this_.ValidateCodeRoot(root)
	}
}

// Java method: AstValidator.validateRoot(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) ValidateRoot(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_ROOT(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_890(n, int32(2))
	this_.ValidateCodeRoot(n.GetFirstChild())
	this_.ValidateCodeRoot(n.GetLastChild())
}

// Java method: AstValidator.validateCodeRoot(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) ValidateCodeRoot(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_ROOT(), n)
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.ValidateScript(c)
		}
	}
}

// Java method: AstValidator.validateScript(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) ValidateScript(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_SCRIPT(), n)
	this_.priv_AstValidator_validateHasSourceName_793(n)
	this_.priv_AstValidator_validateHasInputId_794(n)
	this_.currentScript = n
	if n.HasChildren() && n.GetFirstChild().IsModuleBody() {
		this_.priv_AstValidator_validateProperties_895(n)
		this_.priv_AstValidator_validateChildCount_890(n, int32(1))
		this_.ValidateModuleContents(n.GetFirstChild())
	} else {
		this_.ValidateStatements(n.GetFirstChild())
	}
	if this_.isScriptFeatureValidationEnabled {
		this_.priv_AstValidator_validateScriptFeatureSet_751(n)
	}
}

// Java method: AstValidator.validateScriptFeatureSet(com.google.javascript.rhino.Node scriptNode)
func (this_ *AstValidator) priv_AstValidator_validateScriptFeatureSet_751(scriptNode pkg_rhino.Node) {
	_ = scriptNode
	if !(scriptNode.IsScript()) {
		this_.priv_AstValidator_violation_886("Not a script node", scriptNode)
		return
	}
	var scriptFeatures *pkg_parser.FeatureSet = NodeUtil_GetFeatureSetOfScript(this_.currentScript)
	_ = scriptFeatures
	var allowableFeatures *pkg_parser.FeatureSet = this_.compiler.GetAllowableFeatures()
	_ = allowableFeatures
	if (scriptFeatures == nil) || (allowableFeatures == nil) {
		return
	}
	if !(allowableFeatures.Contains_1_FeatureSet(scriptFeatures)) {
		if !(scriptNode.IsFromExterns()) {
			var differentFeatures *pkg_parser.FeatureSet = scriptFeatures.Without_1(allowableFeatures)
			_ = differentFeatures
			this_.priv_AstValidator_violation_886(("SCRIPT node contains these unallowable features:" + j2gort.StrVal(differentFeatures.GetFeatures())), this_.currentScript)
		}
	}
}

// Java method: AstValidator.validateModuleContents(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) ValidateModuleContents(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_MODULE_BODY(), n)
	this_.ValidateStatements(n.GetFirstChild())
}

// Java method: AstValidator.validateStatements(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) ValidateStatements(n pkg_rhino.Node) {
	_ = n
	for n != 0 {
		this_.ValidateStatement_1(n)
		n = n.GetNext()
	}
}

// Java method: AstValidator.validateStatement(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) ValidateStatement_1(n pkg_rhino.Node) {
	_ = n
	this_.ValidateStatement_2(n, false)
}

// Java method: AstValidator.validateStatement(com.google.javascript.rhino.Node n, boolean isAmbient)
func (this_ *AstValidator) ValidateStatement_2(n pkg_rhino.Node, isAmbient bool) {
	_ = n
	_ = isAmbient
	switch n.TokenSwitchOrd_() {
	case int32(92) /* Token.LABEL */ :
		this_.priv_AstValidator_validateLabel_795(n)
	case int32(90) /* Token.BLOCK */ :
		this_.priv_AstValidator_validateBlock_792(n)
	case int32(73) /* Token.FUNCTION */ :
		if isAmbient {
			this_.priv_AstValidator_validateFunctionSignature_806(n)
		} else {
			this_.priv_AstValidator_validateFunctionStatement_804(n)
		}
	case int32(85) /* Token.WITH */ :
		this_.priv_AstValidator_validateWith_833(n)
	case int32(80) /* Token.FOR */ :
		this_.priv_AstValidator_validateFor_828(n)
	case int32(81) /* Token.FOR_IN */ :
		this_.priv_AstValidator_validateForIn_829(n)
	case int32(112) /* Token.FOR_OF */ :
		this_.priv_AstValidator_validateForOf_830(n)
	case int32(113) /* Token.FOR_AWAIT_OF */ :
		this_.priv_AstValidator_validateForAwaitOf_831(n)
	case int32(78) /* Token.WHILE */ :
		this_.priv_AstValidator_validateWhile_834(n)
	case int32(79) /* Token.DO */ :
		this_.priv_AstValidator_validateDo_835(n)
	case int32(75) /* Token.SWITCH */ :
		this_.priv_AstValidator_validateSwitch_845(n)
	case int32(74) /* Token.IF */ :
		this_.priv_AstValidator_validateIf_836(n)
	case int32(97) /* Token.CONST */, int32(84) /* Token.VAR */, int32(111) /* Token.LET */ :
		this_.priv_AstValidator_validateNameDeclarationHelper_822(n, n.GetToken(), n)
	case int32(93) /* Token.EXPR_RESULT */ :
		this_.priv_AstValidator_validateExprStmt_837(n)
	case int32(0) /* Token.RETURN */ :
		this_.priv_AstValidator_validateReturn_838(n)
	case int32(43) /* Token.THROW */ :
		this_.priv_AstValidator_validateThrow_839(n)
	case int32(48) /* Token.TRY */ :
		this_.priv_AstValidator_validateTry_842(n)
	case int32(82) /* Token.BREAK */ :
		this_.priv_AstValidator_validateBreak_840(n)
	case int32(83) /* Token.CONTINUE */ :
		this_.priv_AstValidator_validateContinue_841(n)
	case int32(88) /* Token.EMPTY */, int32(98) /* Token.DEBUGGER */ :
		this_.priv_AstValidator_validateProperties_895(n)
		this_.priv_AstValidator_validateChildless_849(n)
	case int32(105) /* Token.CLASS */ :
		this_.priv_AstValidator_validateClassDeclaration_784(n, isAmbient)
	case int32(116) /* Token.IMPORT */ :
		this_.priv_AstValidator_validateImport_767(n)
	case int32(120) /* Token.EXPORT */ :
		this_.priv_AstValidator_validateExport_770(n, isAmbient)
	case int32(164) /* Token.INTERFACE */ :
		this_.priv_AstValidator_validateInterface_777(n)
	case int32(167) /* Token.ENUM */ :
		this_.priv_AstValidator_validateEnum_781(n)
	case int32(170) /* Token.TYPE_ALIAS */ :
		this_.priv_AstValidator_validateTypeAlias_880(n)
	case int32(171) /* Token.DECLARE */ :
		this_.priv_AstValidator_validateAmbientDeclaration_881(n)
	case int32(175) /* Token.NAMESPACE */ :
		this_.priv_AstValidator_validateNamespace_883(n, isAmbient)
	case int32(123) /* Token.MODULE_BODY */ :
		if this_.compiler.GetOptions().GetSkipNonTranspilationPasses() {
			if !(this_.compiler.GetOptions().ShouldWrapGoogModulesForWhitespaceOnly()) {
				this_.ValidateModuleContents(n)
				return
			}
		}
		this_.priv_AstValidator_violation_886(("Expected statement but was " + j2gort.StrVal(n.GetToken()) + "."), n)
	default:
		if n.IsModuleBody() && this_.compiler.GetOptions().GetSkipNonTranspilationPasses() {
			if !(!(this_.compiler.GetOptions().ShouldWrapGoogModulesForWhitespaceOnly())) {
				j2gort.CheckStateStr(false, "Modules can exist in transpiler only if setWrapGoogModulesForWhitespaceOnly is false")
			}
			this_.ValidateModuleContents(n)
			return
		}
		this_.priv_AstValidator_violation_886(("Expected statement but was " + j2gort.StrVal(n.GetToken()) + "."), n)
	}
}

// Java method: AstValidator.validateExpression(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) ValidateExpression(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateTypeInformation_758(n)
	this_.priv_AstValidator_validateRequiredInlinings_896(n)
	switch n.TokenSwitchOrd_() {
	case int32(135) /* Token.NEW_TARGET */ :
		this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_NEW_TARGET(), n)
		this_.priv_AstValidator_validateProperties_895(n)
		this_.priv_AstValidator_validateChildless_849(n)
	case int32(136) /* Token.IMPORT_META */ :
		this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_IMPORT_META(), n)
		this_.priv_AstValidator_validateProperties_895(n)
		this_.priv_AstValidator_validateChildless_849(n)
	case int32(38) /* Token.FALSE */, int32(36) /* Token.NULL */, int32(37) /* Token.THIS */, int32(39) /* Token.TRUE */ :
		this_.priv_AstValidator_validateProperties_895(n)
		this_.priv_AstValidator_validateChildless_849(n)
	case int32(24) /* Token.DELPROP */, int32(21) /* Token.POS */, int32(22) /* Token.NEG */, int32(19) /* Token.NOT */, int32(25) /* Token.TYPEOF */, int32(87) /* Token.VOID */, int32(20) /* Token.BITNOT */, int32(101) /* Token.CAST */ :
		this_.priv_AstValidator_validateUnaryOp_876(n)
	case int32(71) /* Token.INC */, int32(72) /* Token.DEC */ :
		this_.priv_AstValidator_validateIncDecOp_875(n)
	case int32(51) /* Token.ASSIGN */ :
		this_.priv_AstValidator_validateAssignmentExpression_850(n)
	case int32(63) /* Token.ASSIGN_EXPONENT */ :
		this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_EXPONENT_OP(), n)
		this_.priv_AstValidator_validateCompoundAssignmentExpression_851(n)
	case int32(52) /* Token.ASSIGN_BITOR */, int32(53) /* Token.ASSIGN_BITXOR */, int32(54) /* Token.ASSIGN_BITAND */, int32(55) /* Token.ASSIGN_LSH */, int32(56) /* Token.ASSIGN_RSH */, int32(57) /* Token.ASSIGN_URSH */, int32(58) /* Token.ASSIGN_ADD */, int32(59) /* Token.ASSIGN_SUB */, int32(60) /* Token.ASSIGN_MUL */, int32(61) /* Token.ASSIGN_DIV */, int32(62) /* Token.ASSIGN_MOD */ :
		this_.priv_AstValidator_validateCompoundAssignmentExpression_851(n)
	case int32(66) /* Token.ASSIGN_COALESCE */ :
		this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_NULL_COALESCE_OP(), n)
		this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_LOGICAL_ASSIGNMENT(), n)
		this_.priv_AstValidator_validateCompoundAssignmentExpression_851(n)
	case int32(64) /* Token.ASSIGN_OR */, int32(65) /* Token.ASSIGN_AND */ :
		this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_LOGICAL_ASSIGNMENT(), n)
		this_.priv_AstValidator_validateCompoundAssignmentExpression_851(n)
	case int32(67) /* Token.HOOK */ :
		this_.priv_AstValidator_validateTrinaryOp_878(n)
	case int32(35) /* Token.STRINGLIT */ :
		this_.priv_AstValidator_validateStringLit_859(n)
	case int32(33) /* Token.NUMBER */ :
		this_.priv_AstValidator_validateNumber_860(n)
	case int32(34) /* Token.BIGINT */ :
		this_.priv_AstValidator_validateBigInt_861(n)
	case int32(32) /* Token.NAME */ :
		this_.priv_AstValidator_validateName_800(n)
	case int32(18) /* Token.EXPONENT */ :
		this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_EXPONENT_OP(), n)
		this_.priv_AstValidator_validateBinaryOp_877(n)
	case int32(70) /* Token.COALESCE */ :
		this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_NULL_COALESCE_OP(), n)
		this_.priv_AstValidator_validateBinaryOp_877(n)
	case int32(50) /* Token.COMMA */, int32(68) /* Token.OR */, int32(69) /* Token.AND */, int32(1) /* Token.BITOR */, int32(2) /* Token.BITXOR */, int32(3) /* Token.BITAND */, int32(4) /* Token.EQ */, int32(5) /* Token.NE */, int32(40) /* Token.SHEQ */, int32(41) /* Token.SHNE */, int32(6) /* Token.LT */, int32(8) /* Token.GT */, int32(7) /* Token.LE */, int32(9) /* Token.GE */, int32(45) /* Token.INSTANCEOF */, int32(44) /* Token.IN */, int32(10) /* Token.LSH */, int32(11) /* Token.RSH */, int32(12) /* Token.URSH */, int32(14) /* Token.SUB */, int32(13) /* Token.ADD */, int32(15) /* Token.MUL */, int32(17) /* Token.MOD */, int32(16) /* Token.DIV */ :
		this_.priv_AstValidator_validateBinaryOp_877(n)
	case int32(27) /* Token.GETELEM */ :
		this_.priv_AstValidator_validateGetElem_853(n)
	case int32(30) /* Token.OPTCHAIN_GETELEM */ :
		this_.priv_AstValidator_validateOptChainGetElem_854(n)
	case int32(26) /* Token.GETPROP */ :
		this_.priv_AstValidator_validateGetProp_855(n)
	case int32(29) /* Token.OPTCHAIN_GETPROP */ :
		this_.priv_AstValidator_validateOptChainGetProp_856(n)
	case int32(46) /* Token.ARRAYLIT */ :
		this_.priv_AstValidator_validateArrayLit_862(n)
	case int32(47) /* Token.OBJECTLIT */ :
		this_.priv_AstValidator_validateObjectLit_863(n)
	case int32(42) /* Token.REGEXP */ :
		this_.priv_AstValidator_validateRegExpLit_858(n)
	case int32(28) /* Token.CALL */ :
		this_.priv_AstValidator_validateCall_812(n)
	case int32(31) /* Token.OPTCHAIN_CALL */ :
		this_.priv_AstValidator_validateOptChainCall_813(n)
	case int32(23) /* Token.NEW */ :
		this_.priv_AstValidator_validateNew_821(n)
	case int32(73) /* Token.FUNCTION */ :
		this_.priv_AstValidator_validateRequiredInlinings_896(n)
		this_.priv_AstValidator_validateFunctionExpression_805(n)
	case int32(105) /* Token.CLASS */ :
		this_.priv_AstValidator_validateClass_785(n)
	case int32(131) /* Token.TEMPLATELIT */ :
		this_.priv_AstValidator_validateTemplateLit_774(n)
	case int32(130) /* Token.TAGGED_TEMPLATELIT */ :
		this_.priv_AstValidator_validateTaggedTemplateLit_773(n)
	case int32(114) /* Token.YIELD */ :
		this_.priv_AstValidator_validateYield_762(n)
	case int32(115) /* Token.AWAIT */ :
		this_.priv_AstValidator_validateAwait_764(n)
	case int32(124) /* Token.DYNAMIC_IMPORT */ :
		this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_DYNAMIC_IMPORT(), n)
		this_.priv_AstValidator_validateUnaryOp_876(n)
	default:
		this_.priv_AstValidator_violation_886(("Expected expression but was " + j2gort.StrVal(n.GetToken())), n)
	}
}

// Java method: AstValidator.validatePseudoExpression(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Token[] allowedPseudoexpressions)
func (this_ *AstValidator) priv_AstValidator_validatePseudoExpression_757(n pkg_rhino.Node, allowedPseudoexpressions []pkg_rhino.Token) {
	_ = n
	_ = allowedPseudoexpressions
	var allowedTokensSet j2gort.Ext_com_google_common_collect_ImmutableSet = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_slice_any(func() []any {
		arr_3 := allowedPseudoexpressions
		if arr_3 == nil {
			return nil
		}
		out_4 := make([]any, len(arr_3))
		for i_, el_5 := range arr_3 {
			out_4[i_] = j2gort.BoxNullableVal(el_5)
		}
		return out_4
	}())
	_ = allowedTokensSet
	switch n.TokenSwitchOrd_() {
	case int32(88) /* Token.EMPTY */ :
		if !(allowedTokensSet.Contains(j2gort.BoxNullableVal(pkg_rhino.StaticGet_Token_EMPTY()))) {
			j2gort.CheckArgumentFmt(false, "Unexpected pseudoexpression %s", j2gort.BoxNode(n))
		}
		this_.priv_AstValidator_validateProperties_895(n)
		this_.priv_AstValidator_validateChildless_849(n)
	case int32(127) /* Token.ITER_SPREAD */ :
		if !(allowedTokensSet.Contains(j2gort.BoxNullableVal(pkg_rhino.StaticGet_Token_ITER_SPREAD()))) {
			j2gort.CheckArgumentFmt(false, "Unexpected pseudoexpression %s", j2gort.BoxNode(n))
		}
		this_.priv_AstValidator_validateProperties_895(n)
		this_.priv_AstValidator_validateChildCount_889(n)
		this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_SPREAD_EXPRESSIONS(), n)
		this_.ValidateExpression(n.GetFirstChild())
	default:
		this_.ValidateExpression(n)
	}
}

// Java method: AstValidator.validateTypeInformation(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateTypeInformation_758(n pkg_rhino.Node) {
	_ = n
	if n.GetIsInClosureUnawareSubtree() {
		return
	}
	{
		var iifeCond_1 bool
		{

			re_6 := this_.typeValidationMode
			ae_7 := StaticGet_AstValidator_TypeInfoValidation_NONE()
			if re_6 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_6 == ae_7
		}
		if iifeCond_1 {
			return
		}
	}
	{
		var iifeCond_1 bool
		{

			re_8 := this_.typeValidationMode
			ae_9 := StaticGet_AstValidator_TypeInfoValidation_JSTYPE()
			if re_8 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_8 == ae_9
		}
		if iifeCond_1 {
			var type_ pkg_rhino.JSType = n.GetJSType()
			_ = type_
			if (type_ != nil && !type_.IsJavaNil_()) && !(type_.IsResolved()) {
				this_.priv_AstValidator_violation_886(("Found unresolved type " + j2gort.StrVal(type_)), n)
			}
		}
	}
	switch n.TokenSwitchOrd_() {
	case int32(28) /* Token.CALL */ :
		if !(n.GetFirstChild().IsSuper()) {
			this_.priv_AstValidator_validateCallType_759(n)
		}
	default:
		this_.priv_AstValidator_expectSomeTypeInformation_760(n)
	}
}

// Java method: AstValidator.validateCallType(com.google.javascript.rhino.Node callNode)
func (this_ *AstValidator) priv_AstValidator_validateCallType_759(callNode pkg_rhino.Node) {
	_ = callNode
	switch (this_.typeValidationMode).Ordinal() {
	case int32(0) /* TypeInfoValidation.JSTYPE */ :
		var callee pkg_rhino.Node = callNode.GetFirstChild()
		_ = callee
		var calleeType pkg_rhino.JSType = j2gort.CheckNotNullIface(callee.GetJSType())
		_ = calleeType
		if calleeType.IsFunctionType() {
			var calleeFunctionType *pkg_rhino.FunctionType = calleeType.ToMaybeFunctionType()
			_ = calleeFunctionType
			var returnType pkg_rhino.JSType = calleeFunctionType.GetReturnType()
			_ = returnType
			if j2gort.IfaceEquals[pkg_rhino.JSType](callNode.GetJSTypeBeforeCast(), nil) && !(returnType.IsUnknownType()) {
				this_.priv_AstValidator_expectMatchingTypeInformation_761(callNode, returnType)
			}
		}
	case int32(1) /* TypeInfoValidation.COLOR */ :
		var callee_2 pkg_rhino.Node = callNode.GetFirstChild()
		_ = callee_2
		j2gort.CheckNotNullPtr(callee_2.GetColor())
	case int32(2) /* TypeInfoValidation.NONE */ :
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
	}
}

// Java method: AstValidator.expectSomeTypeInformation(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_expectSomeTypeInformation_760(n pkg_rhino.Node) {
	_ = n
	switch (this_.typeValidationMode).Ordinal() {
	case int32(0) /* TypeInfoValidation.JSTYPE */ :
		if j2gort.IfaceEquals[pkg_rhino.JSType](n.GetJSType(), nil) {
			this_.priv_AstValidator_violation_886(("Type information missing\n" + j2gort.StrValString(this_.compiler.ToSource_1_Node(NodeUtil_GetEnclosingStatement(n)))), n)
		}
	case int32(1) /* TypeInfoValidation.COLOR */ :
		if n.GetColor() == nil {
			this_.priv_AstValidator_violation_886(("Color information missing\n" + j2gort.StrValString(this_.compiler.ToSource_1_Node(NodeUtil_GetEnclosingStatement(n)))), n)
		}
	case int32(2) /* TypeInfoValidation.NONE */ :
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
	}
}

// Java method: AstValidator.expectMatchingTypeInformation(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType expectedTypeI)
func (this_ *AstValidator) priv_AstValidator_expectMatchingTypeInformation_761(n pkg_rhino.Node, expectedTypeI pkg_rhino.JSType) {
	_ = n
	_ = expectedTypeI
	var typeI pkg_rhino.JSType = n.GetJSType()
	_ = typeI
	if !(j2gort.ObjEqIface[pkg_rhino.JSType](expectedTypeI, typeI)) {
		this_.priv_AstValidator_violation_886(("Expected type: " + j2gort.StrValString(AstValidator_GetTypeAnnotationString(expectedTypeI)) + " Actual type: " + j2gort.StrValString(AstValidator_GetTypeAnnotationString(typeI))), n)
	}
}

// Java static method: AstValidator.getTypeAnnotationString(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType typeI)
func AstValidator_GetTypeAnnotationString(typeI pkg_rhino.JSType) string {
	_ = typeI
	if typeI == nil || typeI.IsJavaNil_() {
		return "NO TYPE INFORMATION"
	} else {
		return ("{" + j2gort.StrValString(typeI.ToAnnotationString(pkg_rhino.StaticGet_JSType_Nullability_EXPLICIT())) + "}")
	}
}

// Java method: AstValidator.validateYield(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateYield_762(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_GENERATORS(), n)
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_YIELD(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCountIn_891(n, int32(0), int32(1))
	if n.HasChildren() {
		this_.ValidateExpression(n.GetFirstChild())
	}
	this_.priv_AstValidator_validateYieldWithinGeneratorFunction_763(n)
}

// Java method: AstValidator.validateYieldWithinGeneratorFunction(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateYieldWithinGeneratorFunction_763(n pkg_rhino.Node) {
	_ = n
	var parentFunction pkg_rhino.Node = NodeUtil_GetEnclosingFunction(n)
	_ = parentFunction
	if (parentFunction == 0) || !(parentFunction.IsGeneratorFunction()) {
		this_.priv_AstValidator_violation_886("'yield' expression is not within a generator function", n)
	} else if this_.priv_AstValidator_isInParameterListOfFunction_766(n, parentFunction) {
		this_.priv_AstValidator_violation_886("'yield' expression is not allowed in a parameter list", n)
	}
}

// Java method: AstValidator.validateAwait(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateAwait_764(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_ASYNC_FUNCTIONS(), n)
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_AWAIT(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	this_.ValidateExpression(n.GetFirstChild())
	this_.priv_AstValidator_validateAwaitWithinAsyncFunction_765(n)
}

// Java method: AstValidator.validateAwaitWithinAsyncFunction(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateAwaitWithinAsyncFunction_765(n pkg_rhino.Node) {
	_ = n
	var parentFunction pkg_rhino.Node = NodeUtil_GetEnclosingFunction(n)
	_ = parentFunction
	if parentFunction == 0 {
		if !(NodeUtil_GetEnclosingScript(n).GetBooleanProp(pkg_rhino.StaticGet_Node_ES6_MODULE())) {
			this_.priv_AstValidator_violation_886("'await' expression is not within an async function", n)
		}
		return
	}
	if !(parentFunction.IsAsyncFunction()) {
		this_.priv_AstValidator_violation_886("'await' expression is not within an async function", n)
	} else if this_.priv_AstValidator_isInParameterListOfFunction_766(n, parentFunction) {
		this_.priv_AstValidator_violation_886("'await' expression is not allowed in a parameter list", n)
	}
}

// Java method: AstValidator.isInParameterListOfFunction(com.google.javascript.rhino.Node child, com.google.javascript.rhino.Node functionNode)
func (this_ *AstValidator) priv_AstValidator_isInParameterListOfFunction_766(child pkg_rhino.Node, functionNode pkg_rhino.Node) bool {
	_ = child
	_ = functionNode
	var paramList pkg_rhino.Node = j2gort.CheckNotNullNode(functionNode.GetSecondChild())
	_ = paramList
	{
		var parent pkg_rhino.Node = child.GetParent()
		_ = parent
		for ; parent != functionNode; parent = parent.GetParent() {
			j2gort.CheckNotNullNode(parent)
			if parent == paramList {
				return true
			}
		}
	}
	return false
}

// Java method: AstValidator.validateImport(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateImport_767(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_MODULES(), n)
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_IMPORT(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	if n.GetFirstChild().IsName() {
		this_.priv_AstValidator_validateName_800(n.GetFirstChild())
	} else {
		this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_EMPTY(), n.GetFirstChild())
	}
	var secondChild pkg_rhino.Node = n.GetSecondChild()
	_ = secondChild
	switch secondChild.TokenSwitchOrd_() {
	case int32(117) /* Token.IMPORT_SPECS */ :
		this_.priv_AstValidator_validateImportSpecifiers_768(secondChild)
	case int32(119) /* Token.IMPORT_STAR */ :
		this_.priv_AstValidator_validateNonEmptyString_797(secondChild)
	default:
		this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_EMPTY(), secondChild)
	}
	this_.priv_AstValidator_validateStringLit_859(n.GetChildAtIndex(int32(2)))
}

// Java method: AstValidator.validateImportSpecifiers(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateImportSpecifiers_768(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_IMPORT_SPECS(), n)
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validateImportSpecifier_769(c)
		}
	}
}

// Java method: AstValidator.validateImportSpecifier(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateImportSpecifier_769(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_IMPORT_SPEC(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_890(n, int32(2))
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validateName_800(c)
		}
	}
}

// Java method: AstValidator.validateExport(com.google.javascript.rhino.Node n, boolean isAmbient)
func (this_ *AstValidator) priv_AstValidator_validateExport_770(n pkg_rhino.Node, isAmbient bool) {
	_ = n
	_ = isAmbient
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_MODULES(), n)
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_EXPORT(), n)
	if n.GetBooleanProp(pkg_rhino.StaticGet_Node_EXPORT_ALL_FROM()) {
		this_.priv_AstValidator_validateProperties_895(n)
		this_.priv_AstValidator_validateChildCount_890(n, int32(2))
		this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_EMPTY(), n.GetFirstChild())
		this_.priv_AstValidator_validateStringLit_859(n.GetSecondChild())
	} else if n.GetBooleanProp(pkg_rhino.StaticGet_Node_EXPORT_DEFAULT()) {
		this_.priv_AstValidator_validateProperties_895(n)
		this_.priv_AstValidator_validateChildCount_890(n, int32(1))
		this_.ValidateExpression(n.GetFirstChild())
	} else {
		this_.priv_AstValidator_validateProperties_895(n)
		this_.priv_AstValidator_validateChildCountIn_891(n, int32(1), int32(2))
		if n.GetFirstChild().IsExportSpecs() {
			this_.priv_AstValidator_validateExportSpecifiers_771(n.GetFirstChild())
		} else {
			this_.ValidateStatement_2(n.GetFirstChild(), isAmbient)
		}
		if n.HasTwoChildren() {
			this_.priv_AstValidator_validateStringLit_859(n.GetSecondChild())
		}
	}
}

// Java method: AstValidator.validateExportSpecifiers(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateExportSpecifiers_771(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_EXPORT_SPECS(), n)
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validateExportSpecifier_772(c)
		}
	}
}

// Java method: AstValidator.validateExportSpecifier(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateExportSpecifier_772(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_EXPORT_SPEC(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_890(n, int32(2))
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validateName_800(c)
		}
	}
}

// Java method: AstValidator.validateTaggedTemplateLit(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateTaggedTemplateLit_773(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_TEMPLATE_LITERALS(), n)
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_TAGGED_TEMPLATELIT(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	this_.ValidateExpression(n.GetFirstChild())
	this_.priv_AstValidator_validateTemplateLit_774(n.GetLastChild())
}

// Java method: AstValidator.validateTemplateLit(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateTemplateLit_774(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_TEMPLATE_LITERALS(), n)
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_TEMPLATELIT(), n)
	{
		var child pkg_rhino.Node = n.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			if child.IsTemplateLitString() {
				this_.priv_AstValidator_validateTemplateLitString_775(child)
			} else {
				this_.priv_AstValidator_validateTemplateLitSub_776(child)
			}
		}
	}
}

// Java method: AstValidator.validateTemplateLitString(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateTemplateLitString_775(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_TEMPLATELIT_STRING(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	func() {
		defer func() {
			if rec_10 := recover(); rec_10 != nil {
				jex_11 := javalang.RecoverException(rec_10)
				_ = jex_11
				if j2gort.InstanceOf[javalang.IUnsupportedOperationException](jex_11) {
					e := j2gort.AsIface[javalang.IUnsupportedOperationException](jex_11)
					_ = e
					this_.priv_AstValidator_violation_886("Invalid TEMPLATELIT_STRING node.", n)
				} else {
					panic(rec_10)
				}
			}
		}()
		n.GetRawString()
	}()
}

// Java method: AstValidator.validateTemplateLitSub(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateTemplateLitSub_776(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_TEMPLATELIT_SUB(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	this_.ValidateExpression(n.GetFirstChild())
}

// Java method: AstValidator.validateInterface(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateInterface_777(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_INTERFACE(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	var name pkg_rhino.Node = n.GetFirstChild()
	_ = name
	this_.priv_AstValidator_validateName_800(name)
	var superTypes pkg_rhino.Node = name.GetNext()
	_ = superTypes
	if superTypes.IsEmpty() {
		this_.priv_AstValidator_validateProperties_895(superTypes)
		this_.priv_AstValidator_validateChildless_849(superTypes)
	} else {
		this_.priv_AstValidator_validateInterfaceExtends_778(superTypes)
	}
	this_.priv_AstValidator_validateInterfaceMembers_779(n.GetLastChild())
}

// Java method: AstValidator.validateInterfaceExtends(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateInterfaceExtends_778(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_INTERFACE_EXTENDS(), n)
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validateNamedType_879(c)
		}
	}
}

// Java method: AstValidator.validateInterfaceMembers(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateInterfaceMembers_779(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_INTERFACE_MEMBERS(), n)
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validateInterfaceMember_780(c)
		}
	}
}

// Java method: AstValidator.validateInterfaceMember(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateInterfaceMember_780(n pkg_rhino.Node) {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
		this_.priv_AstValidator_validateProperties_895(n)
		this_.priv_AstValidator_validateChildCount_889(n)
		this_.priv_AstValidator_validateFunctionSignature_806(n.GetFirstChild())
	case int32(172) /* Token.MEMBER_VARIABLE_DEF */ :
		this_.priv_AstValidator_validateProperties_895(n)
		this_.priv_AstValidator_validateChildless_849(n)
	case int32(173) /* Token.INDEX_SIGNATURE */ :
		this_.priv_AstValidator_validateProperties_895(n)
		this_.priv_AstValidator_validateChildCount_889(n)
		var child pkg_rhino.Node = n.GetFirstChild()
		_ = child
		this_.priv_AstValidator_validateProperties_895(child)
		this_.priv_AstValidator_validateChildless_849(child)
	case int32(174) /* Token.CALL_SIGNATURE */ :
		this_.priv_AstValidator_validateProperties_895(n)
		this_.priv_AstValidator_validateChildCount_889(n)
	default:
		this_.priv_AstValidator_violation_886(("Interface contained member of invalid type " + j2gort.StrVal(n.GetToken())), n)
	}
}

// Java method: AstValidator.validateEnum(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateEnum_781(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_ENUM(), n)
	this_.priv_AstValidator_validateName_800(n.GetFirstChild())
	this_.priv_AstValidator_validateEnumMembers_782(n.GetLastChild())
}

// Java method: AstValidator.validateEnumMembers(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateEnumMembers_782(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_ENUM_MEMBERS(), n)
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validateEnumStringKey_783(c)
		}
	}
}

// Java method: AstValidator.validateEnumStringKey(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateEnumStringKey_783(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_STRING_KEY(), n)
	this_.priv_AstValidator_validateObjectLiteralKeyName_874(n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_890(n, int32(0))
}

// Java method: AstValidator.validateClassDeclaration(com.google.javascript.rhino.Node n, boolean isAmbient)
func (this_ *AstValidator) priv_AstValidator_validateClassDeclaration_784(n pkg_rhino.Node, isAmbient bool) {
	_ = n
	_ = isAmbient
	this_.priv_AstValidator_validateClassHelper_786(n, isAmbient)
	this_.priv_AstValidator_validateName_800(n.GetFirstChild())
}

// Java method: AstValidator.validateClass(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateClass_785(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateClassHelper_786(n, false)
}

// Java method: AstValidator.validateClassHelper(com.google.javascript.rhino.Node n, boolean isAmbient)
func (this_ *AstValidator) priv_AstValidator_validateClassHelper_786(n pkg_rhino.Node, isAmbient bool) {
	_ = n
	_ = isAmbient
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_CLASSES(), n)
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_CLASS(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	var name pkg_rhino.Node = n.GetFirstChild()
	_ = name
	if name.IsEmpty() {
		this_.priv_AstValidator_validateProperties_895(name)
		this_.priv_AstValidator_validateChildless_849(name)
	} else {
		this_.priv_AstValidator_validateName_800(name)
	}
	var superClass pkg_rhino.Node = name.GetNext()
	_ = superClass
	if superClass.IsEmpty() {
		this_.priv_AstValidator_validateProperties_895(superClass)
		this_.priv_AstValidator_validateChildless_849(superClass)
	} else {
		this_.ValidateExpression(superClass)
	}
	this_.priv_AstValidator_validateClassMembers_787(n.GetLastChild(), isAmbient)
}

// Java method: AstValidator.validateClassMembers(com.google.javascript.rhino.Node n, boolean isAmbient)
func (this_ *AstValidator) priv_AstValidator_validateClassMembers_787(n pkg_rhino.Node, isAmbient bool) {
	_ = n
	_ = isAmbient
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_CLASS_MEMBERS(), n)
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validateClassMember_788(c, isAmbient)
		}
	}
}

// Java method: AstValidator.validateClassMember(com.google.javascript.rhino.Node n, boolean isAmbient)
func (this_ *AstValidator) priv_AstValidator_validateClassMember_788(n pkg_rhino.Node, isAmbient bool) {
	_ = n
	_ = isAmbient
	switch n.TokenSwitchOrd_() {
	case int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
		this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_MEMBER_DECLARATIONS(), n)
		if n.IsPrivateIdentifier() {
			this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_PRIVATE_ELEMENTS(), n)
		}
		this_.priv_AstValidator_validateObjectLiteralKeyName_874(n)
		this_.priv_AstValidator_validateProperties_895(n)
		this_.priv_AstValidator_validateChildCount_889(n)
		this_.priv_AstValidator_validateMemberFunction_789(n, isAmbient)
	case int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */ :
		this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_CLASS_GETTER_SETTER(), n)
		if n.IsPrivateIdentifier() {
			this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_PRIVATE_ELEMENTS(), n)
		}
		this_.priv_AstValidator_validateObjectLiteralKeyName_874(n)
		this_.priv_AstValidator_validateObjectLitKey_864(n)
		this_.priv_AstValidator_validateProperties_895(n)
		this_.priv_AstValidator_validateChildCount_889(n)
		this_.priv_AstValidator_validateMemberFunction_789(n, isAmbient)
	case int32(172) /* Token.MEMBER_VARIABLE_DEF */ :
		this_.priv_AstValidator_validateProperties_895(n)
		this_.priv_AstValidator_validateChildless_849(n)
	case int32(129) /* Token.COMPUTED_PROP */ :
		this_.priv_AstValidator_validateComputedPropClassMethod_871(n)
	case int32(108) /* Token.MEMBER_FIELD_DEF */ :
		this_.priv_AstValidator_validateClassField_790(n)
	case int32(109) /* Token.COMPUTED_FIELD_DEF */ :
		this_.priv_AstValidator_validateComputedPropClassField_791(n)
	case int32(173) /* Token.INDEX_SIGNATURE */ :
		this_.priv_AstValidator_validateProperties_895(n)
		this_.priv_AstValidator_validateChildCount_889(n)
		var child pkg_rhino.Node = n.GetFirstChild()
		_ = child
		this_.priv_AstValidator_validateProperties_895(child)
		this_.priv_AstValidator_validateChildless_849(child)
	case int32(174) /* Token.CALL_SIGNATURE */ :
		this_.priv_AstValidator_validateProperties_895(n)
		this_.priv_AstValidator_validateChildCount_889(n)
	case int32(90) /* Token.BLOCK */ :
		this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_CLASS_STATIC_BLOCK(), n)
		this_.priv_AstValidator_validateBlock_792(n)
	case int32(88) /* Token.EMPTY */ :
	default:
		this_.priv_AstValidator_violation_886(("Class contained member of invalid type " + j2gort.StrVal(n.GetToken())), n)
	}
}

// Java method: AstValidator.validateMemberFunction(com.google.javascript.rhino.Node n, boolean isAmbient)
func (this_ *AstValidator) priv_AstValidator_validateMemberFunction_789(n pkg_rhino.Node, isAmbient bool) {
	_ = n
	_ = isAmbient
	var function pkg_rhino.Node = n.GetFirstChild()
	_ = function
	if isAmbient {
		this_.priv_AstValidator_validateFunctionSignature_806(function)
	} else {
		this_.priv_AstValidator_validateFunctionExpression_805(function)
	}
}

// Java method: AstValidator.validateClassField(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateClassField_790(n pkg_rhino.Node) {
	_ = n
	if n.IsPrivateIdentifier() {
		this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_PRIVATE_ELEMENTS(), n)
	} else {
		this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_PUBLIC_CLASS_FIELDS(), n)
	}
	this_.priv_AstValidator_validateNonEmptyString_797(n)
	if n.HasChildren() {
		this_.ValidateExpression(n.GetFirstChild())
	}
}

// Java method: AstValidator.validateComputedPropClassField(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateComputedPropClassField_791(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_PUBLIC_CLASS_FIELDS(), n)
	this_.ValidateExpression(n.GetFirstChild())
	if n.GetSecondChild() != 0 {
		this_.ValidateExpression(n.GetSecondChild())
	}
}

// Java method: AstValidator.validateBlock(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateBlock_792(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_BLOCK(), n)
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.ValidateStatement_1(c)
		}
	}
}

// Java method: AstValidator.validateHasSourceName(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateHasSourceName_793(n pkg_rhino.Node) {
	_ = n
	var sourceName string = n.GetSourceFileName()
	_ = sourceName
	if j2gort.Ext_com_google_common_base_Strings_IsNullOrEmpty(sourceName) {
		this_.priv_AstValidator_violation_886("Missing 'source name' annotation.", n)
	}
}

// Java method: AstValidator.validateHasInputId(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateHasInputId_794(n pkg_rhino.Node) {
	_ = n
	var inputId *pkg_rhino.InputId = n.GetInputId()
	_ = inputId
	if inputId == nil {
		this_.priv_AstValidator_violation_886("Missing 'input id' annotation.", n)
	}
}

// Java method: AstValidator.validateLabel(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateLabel_795(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_LABEL(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	this_.priv_AstValidator_validateLabelName_796(n.GetFirstChild())
	this_.ValidateStatement_1(n.GetLastChild())
}

// Java method: AstValidator.validateLabelName(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateLabelName_796(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_LABEL_NAME(), n)
	this_.priv_AstValidator_validateNonEmptyString_797(n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
}

// Java method: AstValidator.validateNonEmptyString(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateNonEmptyString_797(n pkg_rhino.Node) {
	_ = n
	if this_.priv_AstValidator_validateNonNullString_799(n) && (j2gort.CheckNotNullString(n.GetString()) == "") {
		this_.priv_AstValidator_violation_886("Expected non-empty string.", n)
	}
}

// Java method: AstValidator.validateEmptyString(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateEmptyString_798(n pkg_rhino.Node) {
	_ = n
	if this_.priv_AstValidator_validateNonNullString_799(n) && !(j2gort.CheckNotNullString(n.GetString()) == "") {
		this_.priv_AstValidator_violation_886("Expected empty string.", n)
	}
}

// Java method: AstValidator.validateNonNullString(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateNonNullString_799(n pkg_rhino.Node) bool {
	_ = n
	var didRet_12 bool
	var retVal_13 bool
	func() {
		defer func() {
			if rec_14 := recover(); rec_14 != nil {
				jex_15 := javalang.RecoverException(rec_14)
				_ = jex_15
				if j2gort.InstanceOf[javalang.IRuntimeException](jex_15) {
					e := j2gort.AsIface[javalang.IRuntimeException](jex_15)
					_ = e
					this_.priv_AstValidator_violation_886("Expected non-null string.", n)
					retVal_13 = false
					didRet_12 = true
					return
				} else {
					panic(rec_14)
				}
			}
		}()
		if n.GetString() == j2gort.NullString {
			this_.priv_AstValidator_violation_886("Expected non-null string.", n)
			retVal_13 = false
			didRet_12 = true
			return
		}
	}()
	if didRet_12 {
		return retVal_13
	}
	return true
}

// Java method: AstValidator.validateName(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateName_800(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_NAME(), n)
	this_.priv_AstValidator_validateNonEmptyString_797(n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	this_.priv_AstValidator_validateTypeInformation_758(n)
	this_.priv_AstValidator_validateShadowContentIfPresent_802(n)
}

// Java method: AstValidator.validateOptionalName(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateOptionalName_801(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_NAME(), n)
	this_.priv_AstValidator_validateNonNullString_799(n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	var isEmpty bool = ((n.GetString() != j2gort.NullString) && (j2gort.CheckNotNullString(n.GetString()) == ""))
	_ = isEmpty
	if !(isEmpty) {
		this_.priv_AstValidator_validateTypeInformation_758(n)
	}
}

// Java method: AstValidator.validateShadowContentIfPresent(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateShadowContentIfPresent_802(n pkg_rhino.Node) {
	_ = n
	var shadow pkg_rhino.Node = n.GetClosureUnawareShadow()
	_ = shadow
	if shadow == 0 {
		return
	}
	if !(shadow.IsRoot()) {
		this_.priv_AstValidator_violation_886("Shadow reference node is not a ROOT node", shadow)
		return
	}
	var shadowScript pkg_rhino.Node = shadow.GetFirstChild()
	_ = shadowScript
	if (shadowScript == 0) || !(shadowScript.IsScript()) {
		this_.priv_AstValidator_violation_886("Shadow root node's child is not a script node", shadowScript)
		return
	}
	if shadowScript.GetChildCount() != int32(1) {
		this_.priv_AstValidator_violation_886("Shadow SCRIPT node child has more than one child", shadowScript)
		return
	}
	var exprResult pkg_rhino.Node = shadowScript.GetFirstChild()
	_ = exprResult
	if (exprResult == 0) || !(exprResult.IsExprResult()) {
		this_.priv_AstValidator_violation_886("Shadow SCRIPT node child is not an expr result node", exprResult)
		return
	}
	if exprResult.GetChildCount() != int32(1) {
		this_.priv_AstValidator_violation_886("Shadow EXPR_RESULT node should have exactly one child", exprResult)
		return
	}
	var shadowJsCall pkg_rhino.Node = exprResult.GetOnlyChild()
	_ = shadowJsCall
	if !(shadowJsCall.IsCall()) {
		this_.priv_AstValidator_violation_886("Shadow node EXPR_RESULT child is not a call", shadowJsCall)
		return
	}
	var shadowJsFunction pkg_rhino.Node = shadowJsCall.GetLastChild()
	_ = shadowJsFunction
	if !(shadowJsFunction.IsFunction()) {
		this_.priv_AstValidator_violation_886("Shadow node CALL child is not a function", shadowJsFunction)
		return
	}
	this_.priv_AstValidator_validateFunctionExpression_805(shadowJsFunction)
}

// Java method: AstValidator.validateEmptyName(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateEmptyName_803(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_NAME(), n)
	this_.priv_AstValidator_validateEmptyString_798(n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
}

// Java method: AstValidator.validateFunctionStatement(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateFunctionStatement_804(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_FUNCTION(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	this_.priv_AstValidator_validateName_800(n.GetFirstChild())
	this_.priv_AstValidator_validateParameters_810(n.GetSecondChild())
	this_.priv_AstValidator_validateFunctionBody_809(n.GetLastChild(), false)
	this_.priv_AstValidator_validateFunctionFeatures_808(n)
	if n.GetParent().IsBlock() && !(n.GetGrandparent().IsFunction()) {
		this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION(), n)
	}
}

// Java method: AstValidator.validateFunctionExpression(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateFunctionExpression_805(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFunctionExpressionHelper_807(n, false)
}

// Java method: AstValidator.validateFunctionSignature(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateFunctionSignature_806(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFunctionExpressionHelper_807(n, true)
}

// Java method: AstValidator.validateFunctionExpressionHelper(com.google.javascript.rhino.Node n, boolean isAmbient)
func (this_ *AstValidator) priv_AstValidator_validateFunctionExpressionHelper_807(n pkg_rhino.Node, isAmbient bool) {
	_ = n
	_ = isAmbient
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_FUNCTION(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	this_.priv_AstValidator_validateParameters_810(n.GetSecondChild())
	var name pkg_rhino.Node = n.GetFirstChild()
	_ = name
	var body pkg_rhino.Node = n.GetLastChild()
	_ = body
	if n.IsArrowFunction() {
		this_.priv_AstValidator_validateEmptyName_803(name)
		if body.IsBlock() {
			this_.priv_AstValidator_validateBlock_792(body)
		} else {
			this_.ValidateExpression(body)
		}
	} else {
		this_.priv_AstValidator_validateOptionalName_801(name)
		this_.priv_AstValidator_validateFunctionBody_809(body, isAmbient)
	}
	this_.priv_AstValidator_validateFunctionFeatures_808(n)
}

// Java method: AstValidator.validateFunctionFeatures(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateFunctionFeatures_808(n pkg_rhino.Node) {
	_ = n
	if n.IsArrowFunction() {
		this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_ARROW_FUNCTIONS(), n)
	}
	if n.IsGeneratorFunction() {
		this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_GENERATORS(), n)
	}
	if n.IsAsyncFunction() {
		this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_ASYNC_FUNCTIONS(), n)
	}
	if n.IsAsyncFunction() && n.IsGeneratorFunction() {
		this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_ASYNC_GENERATORS(), n)
	}
}

// Java method: AstValidator.validateFunctionBody(com.google.javascript.rhino.Node n, boolean noBlock)
func (this_ *AstValidator) priv_AstValidator_validateFunctionBody_809(n pkg_rhino.Node, noBlock bool) {
	_ = n
	_ = noBlock
	if noBlock {
		this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_EMPTY(), n)
	} else {
		this_.priv_AstValidator_validateBlock_792(n)
	}
}

// Java method: AstValidator.validateParameters(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateParameters_810(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_PARAM_LIST(), n)
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			if c.IsRest() {
				this_.priv_AstValidator_validateRestParameters_816(pkg_rhino.StaticGet_Token_PARAM_LIST(), c)
			} else if c.IsDefaultValue() {
				this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_DEFAULT_PARAMETERS(), c)
				this_.priv_AstValidator_validateDefaultValue_811(pkg_rhino.StaticGet_Token_PARAM_LIST(), c)
			} else {
				if c.IsName() {
					this_.priv_AstValidator_validateName_800(c)
				} else if c.IsArrayPattern() {
					this_.priv_AstValidator_validateArrayPattern_826(pkg_rhino.StaticGet_Token_PARAM_LIST(), c)
				} else {
					this_.priv_AstValidator_validateObjectPattern_827(pkg_rhino.StaticGet_Token_PARAM_LIST(), c)
				}
			}
		}
	}
}

// Java method: AstValidator.validateDefaultValue(com.google.javascript.rhino.Token contextType, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateDefaultValue_811(contextType pkg_rhino.Token, n pkg_rhino.Node) {
	_ = contextType
	_ = n
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	this_.priv_AstValidator_validateLHS_824(contextType, n.GetFirstChild())
	this_.ValidateExpression(n.GetLastChild())
}

// Java method: AstValidator.validateCall(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateCall_812(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_CALL(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateMinimumChildCount_892(n, int32(1))
	var callee pkg_rhino.Node = n.GetFirstChild()
	_ = callee
	if callee.IsSuper() {
		this_.priv_AstValidator_validateSuper_814(callee)
	} else {
		this_.ValidateExpression(callee)
	}
	{
		var c pkg_rhino.Node = callee.GetNext()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validatePseudoExpression_757(c, []pkg_rhino.Token{pkg_rhino.StaticGet_Token_ITER_SPREAD()})
		}
	}
}

// Java method: AstValidator.validateOptChainCall(com.google.javascript.rhino.Node node)
func (this_ *AstValidator) priv_AstValidator_validateOptChainCall_813(node pkg_rhino.Node) {
	_ = node
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_OPTIONAL_CHAINING(), node)
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_OPTCHAIN_CALL(), node)
	this_.priv_AstValidator_validateProperties_895(node)
	this_.priv_AstValidator_validateMinimumChildCount_892(node, int32(1))
	var callee pkg_rhino.Node = node.GetFirstChild()
	_ = callee
	this_.ValidateExpression(callee)
	{
		var argument pkg_rhino.Node = callee.GetNext()
		_ = argument
		for ; argument != 0; argument = argument.GetNext() {
			this_.priv_AstValidator_validatePseudoExpression_757(argument, []pkg_rhino.Token{pkg_rhino.StaticGet_Token_ITER_SPREAD()})
		}
	}
	this_.priv_AstValidator_validateFirstNodeOfOptChain_887(node)
}

// Java method: AstValidator.validateSuper(com.google.javascript.rhino.Node superNode)
func (this_ *AstValidator) priv_AstValidator_validateSuper_814(superNode pkg_rhino.Node) {
	_ = superNode
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_SUPER(), superNode)
	this_.priv_AstValidator_validateProperties_895(superNode)
	this_.priv_AstValidator_validateChildless_849(superNode)
	this_.priv_AstValidator_validateTypeInformation_758(superNode)
	var superParent pkg_rhino.Node = superNode.GetParent()
	_ = superParent
	var methodNode pkg_rhino.Node = NodeUtil_GetEnclosingNonArrowFunction(superParent)
	_ = methodNode
	if NodeUtil_IsNormalGet(superParent) && superNode.IsFirstChildOf(superParent) {
		if !(this_.priv_AstValidator_allowsSuperPropertyReference_815(superParent)) {
			this_.priv_AstValidator_violation_886("super property references are only allowed in methods, class static blocks and class fields.", superNode)
		}
	} else if superParent.IsCall() && superNode.IsFirstChildOf(superParent) {
		if (methodNode == 0) || !(NodeUtil_IsEs6Constructor(methodNode)) {
			this_.priv_AstValidator_violation_886("super constructor call is only allowed in a constructor method", superNode)
		} else {
			var extendsNode pkg_rhino.Node = methodNode.GetParent().GetParent().GetParent().GetSecondChild()
			_ = extendsNode
			if extendsNode.IsEmpty() {
				this_.priv_AstValidator_violation_886("super constructor call in a class that extends nothing", superNode)
			}
		}
	} else {
		this_.priv_AstValidator_violation_886("`super` is a syntax error here", superNode)
	}
}

// Java method: AstValidator.allowsSuperPropertyReference(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_allowsSuperPropertyReference_815(n pkg_rhino.Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(94) /* Token.SCRIPT */ :
		return false
	case int32(108) /* Token.MEMBER_FIELD_DEF */, int32(109) /* Token.COMPUTED_FIELD_DEF */ :
		return true
	case int32(73) /* Token.FUNCTION */ :
		if NodeUtil_IsMethodDeclaration(n) {
			return true
		} else if !(n.IsArrowFunction()) {
			return false
		}
		return this_.priv_AstValidator_allowsSuperPropertyReference_815(n.GetParent())
	case int32(90) /* Token.BLOCK */ :
		if NodeUtil_IsClassStaticBlock(n) {
			return true
		}
		return this_.priv_AstValidator_allowsSuperPropertyReference_815(n.GetParent())
	default:
		return this_.priv_AstValidator_allowsSuperPropertyReference_815(n.GetParent())
	}
	panic("unreachable switch expression")
}

// Java method: AstValidator.validateRestParameters(com.google.javascript.rhino.Token contextType, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateRestParameters_816(contextType pkg_rhino.Token, n pkg_rhino.Node) {
	_ = contextType
	_ = n
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_REST_PARAMETERS(), n)
	this_.priv_AstValidator_validateRest_819(contextType, n)
}

// Java method: AstValidator.validateArrayPatternRest(com.google.javascript.rhino.Token contextType, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateArrayPatternRest_817(contextType pkg_rhino.Token, n pkg_rhino.Node) {
	_ = contextType
	_ = n
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_ARRAY_PATTERN_REST(), n)
	this_.priv_AstValidator_validateRest_819(contextType, n)
}

// Java method: AstValidator.validateObjectPatternRest(com.google.javascript.rhino.Token contextType, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateObjectPatternRest_818(contextType pkg_rhino.Token, n pkg_rhino.Node) {
	_ = contextType
	_ = n
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_OBJECT_PATTERN_REST(), n)
	this_.priv_AstValidator_validateRest_819(contextType, n)
}

// Java method: AstValidator.validateRest(com.google.javascript.rhino.Token contextType, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateRest_819(contextType pkg_rhino.Token, n pkg_rhino.Node) {
	_ = contextType
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(125) /* Token.ITER_REST */, int32(126) /* Token.OBJECT_REST */ :
	default:
		this_.priv_AstValidator_violation_886("Unexpected node type.", n)
		return
	}
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	this_.priv_AstValidator_validateLHS_824(contextType, n.GetFirstChild())
	if n.GetNext() != 0 {
		this_.priv_AstValidator_violation_886("Rest parameters must come after all other parameters.", n)
	}
}

// Java method: AstValidator.validateObjectSpread(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateObjectSpread_820(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD(), n)
	this_.ValidateExpression(n.GetFirstChild())
}

// Java method: AstValidator.validateNew(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateNew_821(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_NEW(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateMinimumChildCount_892(n, int32(1))
	this_.ValidateExpression(n.GetFirstChild())
	{
		var c pkg_rhino.Node = n.GetSecondChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validatePseudoExpression_757(c, []pkg_rhino.Token{pkg_rhino.StaticGet_Token_ITER_SPREAD()})
		}
	}
}

// Java method: AstValidator.validateNameDeclarationHelper(com.google.javascript.rhino.Node statement, com.google.javascript.rhino.Token declaration, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateNameDeclarationHelper_822(statement pkg_rhino.Node, declaration pkg_rhino.Token, n pkg_rhino.Node) {
	_ = statement
	_ = declaration
	_ = n
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateMinimumChildCount_892(n, int32(1))
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validateNameDeclarationChild_823(statement, declaration, c)
		}
	}
	{
		var iifeCond_2 bool
		{

			re_16 := declaration
			ae_17 := pkg_rhino.StaticGet_Token_LET()
			if re_16 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_2 = re_16 == ae_17
		}
		if iifeCond_2 {
			this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_LET_DECLARATIONS(), n)
		} else {
			var iifeCond_1 bool
			{

				re_18 := declaration
				ae_19 := pkg_rhino.StaticGet_Token_CONST()
				if re_18 == 0 {
					j2gort.ThrowNPE()
				}
				iifeCond_1 = re_18 == ae_19
			}
			if iifeCond_1 {
				this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_CONST_DECLARATIONS(), n)
			}
		}
	}
}

// Java method: AstValidator.validateNameDeclarationChild(com.google.javascript.rhino.Node statement, com.google.javascript.rhino.Token declaration, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateNameDeclarationChild_823(statement pkg_rhino.Node, declaration pkg_rhino.Token, n pkg_rhino.Node) {
	_ = statement
	_ = declaration
	_ = n
	var inEnhancedFor bool = NodeUtil_IsEnhancedFor(statement)
	_ = inEnhancedFor
	var inForIn bool = statement.IsForIn()
	_ = inForIn
	var minValues int32 = int32(0)
	_ = minValues
	var maxValues int32 = int32(0)
	_ = maxValues
	if inForIn && func() bool {
		re_20 := declaration
		ae_21 := pkg_rhino.StaticGet_Token_VAR()
		if re_20 == 0 {
			j2gort.ThrowNPE()
		}
		return re_20 == ae_21
	}() {
		minValues = int32(0)
		maxValues = int32(1)
	} else if inEnhancedFor {
		minValues = int32(0)
		maxValues = int32(0)
	} else if n.IsDestructuringLhs() || func() bool {
		re_22 := declaration
		ae_23 := pkg_rhino.StaticGet_Token_CONST()
		if re_22 == 0 {
			j2gort.ThrowNPE()
		}
		return re_22 == ae_23
	}() {
		minValues = int32(1)
		maxValues = int32(1)
	} else {
		minValues = int32(0)
		maxValues = int32(1)
	}
	if n.IsName() {
		this_.priv_AstValidator_validateNonEmptyString_797(n)
		this_.priv_AstValidator_validateProperties_895(n)
		this_.priv_AstValidator_validateChildCountIn_891(n, minValues, maxValues)
		if n.HasChildren() {
			this_.ValidateExpression(n.GetFirstChild())
		}
	} else if n.IsDestructuringLhs() {
		this_.priv_AstValidator_validateProperties_895(n)
		this_.priv_AstValidator_validateChildCountIn_891(n, (int32(1) + minValues), (int32(1) + maxValues))
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		switch c.TokenSwitchOrd_() {
		case int32(102) /* Token.ARRAY_PATTERN */ :
			this_.priv_AstValidator_validateArrayPattern_826(declaration, c)
		case int32(103) /* Token.OBJECT_PATTERN */ :
			this_.priv_AstValidator_validateObjectPattern_827(declaration, c)
		default:
			this_.priv_AstValidator_violation_886(("Invalid destructuring lhs first child for " + j2gort.StrVal(declaration) + " node"), n)
		}
		if n.HasTwoChildren() {
			this_.ValidateExpression(n.GetSecondChild())
		}
	} else {
		this_.priv_AstValidator_violation_886(("Invalid child for " + j2gort.StrVal(declaration) + " node"), n)
	}
}

// Java method: AstValidator.validateLHS(com.google.javascript.rhino.Token contextType, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateLHS_824(contextType pkg_rhino.Token, n pkg_rhino.Node) {
	_ = contextType
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		this_.priv_AstValidator_validateName_800(n)
	case int32(102) /* Token.ARRAY_PATTERN */ :
		this_.priv_AstValidator_validateArrayPattern_826(contextType, n)
	case int32(103) /* Token.OBJECT_PATTERN */ :
		this_.priv_AstValidator_validateObjectPattern_827(contextType, n)
	case int32(26) /* Token.GETPROP */, int32(27) /* Token.GETELEM */ :
		this_.priv_AstValidator_validateGetPropGetElemInLHS_825(contextType, n)
	case int32(101) /* Token.CAST */ :
		this_.priv_AstValidator_validateLHS_824(contextType, n.GetOnlyChild())
	default:
		this_.priv_AstValidator_violation_886(("Invalid child for " + j2gort.StrVal(contextType) + " node"), n)
	}
}

// Java method: AstValidator.validateGetPropGetElemInLHS(com.google.javascript.rhino.Token contextType, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateGetPropGetElemInLHS_825(contextType pkg_rhino.Token, n pkg_rhino.Node) {
	_ = contextType
	_ = n
	if (((contextType == pkg_rhino.StaticGet_Token_CONST()) || (contextType == pkg_rhino.StaticGet_Token_LET())) || (contextType == pkg_rhino.StaticGet_Token_VAR())) || (contextType == pkg_rhino.StaticGet_Token_PARAM_LIST()) {
		this_.priv_AstValidator_violation_886(("Invalid child for " + j2gort.StrVal(contextType) + " node"), n)
		return
	}
	switch n.TokenSwitchOrd_() {
	case int32(26) /* Token.GETPROP */ :
		this_.priv_AstValidator_validateGetProp_855(n)
	case int32(27) /* Token.GETELEM */ :
		this_.priv_AstValidator_validateGetElem_853(n)
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Expected GETPROP or GETELEM but instead got node "+j2gort.StrVal(n.GetToken()))))
	}
}

// Java method: AstValidator.validateArrayPattern(com.google.javascript.rhino.Token type, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateArrayPattern_826(type_ pkg_rhino.Token, n pkg_rhino.Node) {
	_ = type_
	_ = n
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_ARRAY_DESTRUCTURING(), n)
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_ARRAY_PATTERN(), n)
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			switch c.TokenSwitchOrd_() {
			case int32(134) /* Token.DEFAULT_VALUE */ :
				this_.priv_AstValidator_validateDefaultValue_811(type_, c)
			case int32(125) /* Token.ITER_REST */ :
				this_.priv_AstValidator_validateArrayPatternRest_817(type_, c)
			case int32(88) /* Token.EMPTY */ :
				this_.priv_AstValidator_validateProperties_895(c)
				this_.priv_AstValidator_validateChildless_849(c)
			default:
				this_.priv_AstValidator_validateLHS_824(type_, c)
			}
		}
	}
}

// Java method: AstValidator.validateObjectPattern(com.google.javascript.rhino.Token type, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateObjectPattern_827(type_ pkg_rhino.Token, n pkg_rhino.Node) {
	_ = type_
	_ = n
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_OBJECT_DESTRUCTURING(), n)
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_OBJECT_PATTERN(), n)
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			switch c.TokenSwitchOrd_() {
			case int32(100) /* Token.STRING_KEY */ :
				this_.priv_AstValidator_validateObjectPatternStringKey_868(type_, c)
			case int32(126) /* Token.OBJECT_REST */ :
				this_.priv_AstValidator_validateObjectPatternRest_818(type_, c)
			case int32(129) /* Token.COMPUTED_PROP */ :
				this_.priv_AstValidator_validateObjectPatternComputedPropKey_870(type_, c)
			default:
				this_.priv_AstValidator_violation_886(("Invalid object pattern child for " + j2gort.StrVal(type_) + " node"), n)
			}
		}
	}
}

// Java method: AstValidator.validateFor(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateFor_828(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_FOR(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_890(n, int32(4))
	var target pkg_rhino.Node = n.GetFirstChild()
	_ = target
	if NodeUtil_IsNameDeclaration(target) {
		this_.priv_AstValidator_validateNameDeclarationHelper_822(n, target.GetToken(), target)
	} else {
		this_.priv_AstValidator_validatePseudoExpression_757(target, []pkg_rhino.Token{pkg_rhino.StaticGet_Token_EMPTY()})
	}
	this_.priv_AstValidator_validatePseudoExpression_757(n.GetSecondChild(), []pkg_rhino.Token{pkg_rhino.StaticGet_Token_EMPTY()})
	this_.priv_AstValidator_validatePseudoExpression_757(n.GetChildAtIndex(int32(2)), []pkg_rhino.Token{pkg_rhino.StaticGet_Token_EMPTY()})
	this_.priv_AstValidator_validateBlock_792(n.GetLastChild())
}

// Java method: AstValidator.validateForIn(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateForIn_829(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_FOR_IN(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	this_.priv_AstValidator_validateEnhancedForVarOrAssignmentTarget_832(n, n.GetFirstChild())
	this_.ValidateExpression(n.GetSecondChild())
	this_.priv_AstValidator_validateBlock_792(n.GetLastChild())
}

// Java method: AstValidator.validateForOf(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateForOf_830(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_FOR_OF(), n)
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_FOR_OF(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	this_.priv_AstValidator_validateEnhancedForVarOrAssignmentTarget_832(n, n.GetFirstChild())
	this_.ValidateExpression(n.GetSecondChild())
	this_.priv_AstValidator_validateBlock_792(n.GetLastChild())
}

// Java method: AstValidator.validateForAwaitOf(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateForAwaitOf_831(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_FOR_AWAIT_OF(), n)
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_FOR_AWAIT_OF(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	this_.priv_AstValidator_validateEnhancedForVarOrAssignmentTarget_832(n, n.GetFirstChild())
	this_.ValidateExpression(n.GetSecondChild())
	this_.priv_AstValidator_validateBlock_792(n.GetLastChild())
	this_.priv_AstValidator_validateAwaitWithinAsyncFunction_765(n)
}

// Java method: AstValidator.validateEnhancedForVarOrAssignmentTarget(com.google.javascript.rhino.Node forNode, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateEnhancedForVarOrAssignmentTarget_832(forNode pkg_rhino.Node, n pkg_rhino.Node) {
	_ = forNode
	_ = n
	if NodeUtil_IsNameDeclaration(n) {
		this_.priv_AstValidator_validateProperties_895(n)
		this_.priv_AstValidator_validateChildCount_890(n, int32(1))
		this_.priv_AstValidator_validateNameDeclarationHelper_822(forNode, n.GetToken(), n)
	} else {
		this_.priv_AstValidator_validateLHS_824(n.GetParent().GetToken(), n)
	}
}

// Java method: AstValidator.validateWith(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateWith_833(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_WITH(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	this_.ValidateExpression(n.GetFirstChild())
	this_.priv_AstValidator_validateBlock_792(n.GetLastChild())
}

// Java method: AstValidator.validateWhile(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateWhile_834(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_WHILE(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	this_.ValidateExpression(n.GetFirstChild())
	this_.priv_AstValidator_validateBlock_792(n.GetLastChild())
}

// Java method: AstValidator.validateDo(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateDo_835(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_DO(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	this_.priv_AstValidator_validateBlock_792(n.GetFirstChild())
	this_.ValidateExpression(n.GetLastChild())
}

// Java method: AstValidator.validateIf(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateIf_836(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_IF(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCountIn_891(n, int32(2), int32(3))
	this_.ValidateExpression(n.GetFirstChild())
	this_.priv_AstValidator_validateBlock_792(n.GetSecondChild())
	if n.HasXChildren(int32(3)) {
		this_.priv_AstValidator_validateBlock_792(n.GetLastChild())
	}
}

// Java method: AstValidator.validateExprStmt(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateExprStmt_837(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_EXPR_RESULT(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	this_.ValidateExpression(n.GetFirstChild())
}

// Java method: AstValidator.validateReturn(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateReturn_838(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_RETURN(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateMaximumChildCount_893(n, int32(1))
	if n.HasChildren() {
		this_.ValidateExpression(n.GetFirstChild())
	}
}

// Java method: AstValidator.validateThrow(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateThrow_839(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_THROW(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	this_.ValidateExpression(n.GetFirstChild())
}

// Java method: AstValidator.validateBreak(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateBreak_840(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_BREAK(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateMaximumChildCount_893(n, int32(1))
	if n.HasChildren() {
		this_.priv_AstValidator_validateLabelName_796(n.GetFirstChild())
	}
}

// Java method: AstValidator.validateContinue(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateContinue_841(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_CONTINUE(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateMaximumChildCount_893(n, int32(1))
	if n.HasChildren() {
		this_.priv_AstValidator_validateLabelName_796(n.GetFirstChild())
	}
}

// Java method: AstValidator.validateTry(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateTry_842(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_TRY(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCountIn_891(n, int32(2), int32(3))
	this_.priv_AstValidator_validateBlock_792(n.GetFirstChild())
	var seenCatchOrFinally bool = false
	_ = seenCatchOrFinally
	var catches pkg_rhino.Node = n.GetSecondChild()
	_ = catches
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_BLOCK(), catches)
	this_.priv_AstValidator_validateProperties_895(catches)
	this_.priv_AstValidator_validateMaximumChildCount_893(catches, int32(1))
	if catches.HasChildren() {
		this_.priv_AstValidator_validateCatch_843(catches.GetFirstChild())
		seenCatchOrFinally = true
	}
	if n.HasXChildren(int32(3)) {
		this_.priv_AstValidator_validateBlock_792(n.GetLastChild())
		seenCatchOrFinally = true
	}
	if !(seenCatchOrFinally) {
		this_.priv_AstValidator_violation_886("Missing catch or finally for try statement.", n)
	}
}

// Java method: AstValidator.validateCatch(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateCatch_843(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_CATCH(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	var caught pkg_rhino.Node = n.GetFirstChild()
	_ = caught
	if caught.IsName() {
		this_.priv_AstValidator_validateName_800(caught)
	} else if caught.IsArrayPattern() {
		this_.priv_AstValidator_validateArrayPattern_826(pkg_rhino.StaticGet_Token_CATCH(), caught)
	} else if caught.IsObjectPattern() {
		this_.priv_AstValidator_validateObjectPattern_827(pkg_rhino.StaticGet_Token_CATCH(), caught)
	} else if caught.IsEmpty() {
		this_.priv_AstValidator_validateNoCatchBinding_844(caught)
	} else {
		this_.priv_AstValidator_violation_886(("Unexpected catch binding: " + j2gort.StrVal(caught)), n)
	}
	this_.priv_AstValidator_validateBlock_792(n.GetLastChild())
}

// Java method: AstValidator.validateNoCatchBinding(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateNoCatchBinding_844(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_OPTIONAL_CATCH_BINDING(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
}

// Java method: AstValidator.validateSwitch(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateSwitch_845(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_SWITCH(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_890(n, int32(2))
	this_.ValidateExpression(n.GetFirstChild())
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_SWITCH_BODY(), n.GetSecondChild())
	var cases pkg_rhino.Node = n.GetSecondChild()
	_ = cases
	var defaults int32 = int32(0)
	_ = defaults
	{
		var c pkg_rhino.Node = cases.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validateSwitchMember_846(c)
			if c.IsDefaultCase() {
				defaults = defaults + 1
			}
		}
	}
	if defaults > int32(1) {
		this_.priv_AstValidator_violation_886(("Expected at most 1 'default' in switch but was " + j2gort.StrValInt32(defaults)), n)
	}
}

// Java method: AstValidator.validateSwitchMember(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateSwitchMember_846(n pkg_rhino.Node) {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(76) /* Token.CASE */ :
		this_.priv_AstValidator_validateCase_848(n)
	case int32(77) /* Token.DEFAULT_CASE */ :
		this_.priv_AstValidator_validateDefaultCase_847(n)
	default:
		this_.priv_AstValidator_violation_886(("Expected switch member but was " + j2gort.StrVal(n.GetToken())), n)
	}
}

// Java method: AstValidator.validateDefaultCase(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateDefaultCase_847(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_DEFAULT_CASE(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	this_.priv_AstValidator_validateBlock_792(n.GetLastChild())
}

// Java method: AstValidator.validateCase(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateCase_848(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_CASE(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	this_.ValidateExpression(n.GetFirstChild())
	this_.priv_AstValidator_validateBlock_792(n.GetLastChild())
}

// Java method: AstValidator.validateChildless(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateChildless_849(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateChildCount_890(n, int32(0))
}

// Java method: AstValidator.validateAssignmentExpression(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateAssignmentExpression_850(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	this_.priv_AstValidator_validateLHS_824(n.GetToken(), n.GetFirstChild())
	this_.ValidateExpression(n.GetLastChild())
}

// Java method: AstValidator.validateCompoundAssignmentExpression(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateCompoundAssignmentExpression_851(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	var contextType pkg_rhino.Token = n.GetToken()
	_ = contextType
	var lhs pkg_rhino.Node = n.GetFirstChild()
	_ = lhs
	this_.priv_AstValidator_validateAssignmentOpTarget_852(lhs, contextType)
	this_.ValidateExpression(n.GetLastChild())
}

// Java method: AstValidator.validateAssignmentOpTarget(com.google.javascript.rhino.Node lhs, com.google.javascript.rhino.Token contextType)
func (this_ *AstValidator) priv_AstValidator_validateAssignmentOpTarget_852(lhs pkg_rhino.Node, contextType pkg_rhino.Token) {
	_ = lhs
	_ = contextType
	switch lhs.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		this_.priv_AstValidator_validateName_800(lhs)
	case int32(26) /* Token.GETPROP */, int32(27) /* Token.GETELEM */ :
		this_.priv_AstValidator_validateGetPropGetElemInLHS_825(contextType, lhs)
	case int32(101) /* Token.CAST */ :
		this_.priv_AstValidator_validateProperties_895(lhs)
		this_.priv_AstValidator_validateChildCount_890(lhs, int32(1))
		this_.priv_AstValidator_validateAssignmentOpTarget_852(lhs.GetFirstChild(), contextType)
	default:
		this_.priv_AstValidator_violation_886(("Invalid child for " + j2gort.StrVal(contextType) + " node"), lhs)
	}
}

// Java method: AstValidator.validateGetElem(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateGetElem_853(n pkg_rhino.Node) {
	_ = n
	if !(n.IsGetElem()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(n))
	}
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_890(n, int32(2))
	this_.priv_AstValidator_validatePropertyReferenceTarget_857(n.GetFirstChild())
	this_.ValidateExpression(n.GetLastChild())
}

// Java method: AstValidator.validateOptChainGetElem(com.google.javascript.rhino.Node node)
func (this_ *AstValidator) priv_AstValidator_validateOptChainGetElem_854(node pkg_rhino.Node) {
	_ = node
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_OPTIONAL_CHAINING(), node)
	if !(node.IsOptChainGetElem()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(node))
	}
	this_.priv_AstValidator_validateProperties_895(node)
	this_.priv_AstValidator_validateChildCount_890(node, int32(2))
	this_.ValidateExpression(node.GetFirstChild())
	this_.ValidateExpression(node.GetLastChild())
	this_.priv_AstValidator_validateFirstNodeOfOptChain_887(node)
}

// Java method: AstValidator.validateGetProp(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateGetProp_855(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_GETPROP(), n)
	this_.priv_AstValidator_validatePropertyReferenceTarget_857(n.GetFirstChild())
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	this_.priv_AstValidator_validateNonEmptyString_797(n)
}

// Java method: AstValidator.validateOptChainGetProp(com.google.javascript.rhino.Node node)
func (this_ *AstValidator) priv_AstValidator_validateOptChainGetProp_856(node pkg_rhino.Node) {
	_ = node
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_OPTIONAL_CHAINING(), node)
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_OPTCHAIN_GETPROP(), node)
	this_.ValidateExpression(node.GetFirstChild())
	this_.priv_AstValidator_validateFirstNodeOfOptChain_887(node)
	this_.priv_AstValidator_validateProperties_895(node)
	this_.priv_AstValidator_validateChildCount_889(node)
	this_.priv_AstValidator_validateNonEmptyString_797(node)
}

// Java method: AstValidator.validatePropertyReferenceTarget(com.google.javascript.rhino.Node objectNode)
func (this_ *AstValidator) priv_AstValidator_validatePropertyReferenceTarget_857(objectNode pkg_rhino.Node) {
	_ = objectNode
	if objectNode.IsSuper() {
		this_.priv_AstValidator_validateSuper_814(objectNode)
	} else {
		this_.ValidateExpression(objectNode)
	}
}

// Java method: AstValidator.validateRegExpLit(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateRegExpLit_858(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_SYNTAX(), n)
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_REGEXP(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCountIn_891(n, int32(1), int32(2))
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validateStringLit_859(c)
		}
	}
}

// Java method: AstValidator.validateStringLit(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateStringLit_859(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_STRINGLIT(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	func() {
		defer func() {
			if rec_24 := recover(); rec_24 != nil {
				jex_25 := javalang.RecoverException(rec_24)
				_ = jex_25
				if j2gort.InstanceOf[javalang.IUnsupportedOperationException](jex_25) {
					e := j2gort.AsIface[javalang.IUnsupportedOperationException](jex_25)
					_ = e
					this_.priv_AstValidator_violation_886("Invalid STRING node.", n)
				} else {
					panic(rec_24)
				}
			}
		}()
		n.GetString()
	}()
}

// Java method: AstValidator.validateNumber(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateNumber_860(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_NUMBER(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	func() {
		defer func() {
			if rec_26 := recover(); rec_26 != nil {
				jex_27 := javalang.RecoverException(rec_26)
				_ = jex_27
				if j2gort.InstanceOf[javalang.IUnsupportedOperationException](jex_27) {
					e := j2gort.AsIface[javalang.IUnsupportedOperationException](jex_27)
					_ = e
					this_.priv_AstValidator_violation_886("Invalid NUMBER node.", n)
				} else {
					panic(rec_26)
				}
			}
		}()
		n.GetDouble()
	}()
}

// Java method: AstValidator.validateBigInt(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateBigInt_861(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_BIGINT(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	func() {
		defer func() {
			if rec_28 := recover(); rec_28 != nil {
				jex_29 := javalang.RecoverException(rec_28)
				_ = jex_29
				if j2gort.InstanceOf[javalang.IUnsupportedOperationException](jex_29) {
					e := j2gort.AsIface[javalang.IUnsupportedOperationException](jex_29)
					_ = e
					this_.priv_AstValidator_violation_886("Invalid BIGINT node.", n)
				} else {
					panic(rec_28)
				}
			}
		}()
		n.GetBigInt()
	}()
}

// Java method: AstValidator.validateArrayLit(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateArrayLit_862(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_ARRAYLIT(), n)
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validatePseudoExpression_757(c, []pkg_rhino.Token{pkg_rhino.StaticGet_Token_EMPTY(), pkg_rhino.StaticGet_Token_ITER_SPREAD()})
			break
		}
	}
}

// Java method: AstValidator.validateObjectLit(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateObjectLit_863(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_OBJECTLIT(), n)
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validateObjectLitKey_864(c)
		}
	}
}

// Java method: AstValidator.validateObjectLitKey(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateObjectLitKey_864(n pkg_rhino.Node) {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(95) /* Token.GETTER_DEF */ :
		this_.priv_AstValidator_validateObjectLitGetKey_865(n)
	case int32(96) /* Token.SETTER_DEF */ :
		this_.priv_AstValidator_validateObjectLitSetKey_866(n)
	case int32(100) /* Token.STRING_KEY */ :
		this_.priv_AstValidator_validateObjectLitStringKey_867(n)
	case int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
		this_.priv_AstValidator_validateClassMember_788(n, false)
		if n.IsStaticMember() {
			this_.priv_AstValidator_violation_886("Keys in an object literal should not be static.", n)
		}
	case int32(129) /* Token.COMPUTED_PROP */ :
		this_.priv_AstValidator_validateObjectLitComputedPropKey_869(n)
	case int32(128) /* Token.OBJECT_SPREAD */ :
		this_.priv_AstValidator_validateObjectSpread_820(n)
	default:
		this_.priv_AstValidator_violation_886(("Expected object literal key expression but was " + j2gort.StrVal(n.GetToken())), n)
	}
}

// Java method: AstValidator.validateObjectLitGetKey(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateObjectLitGetKey_865(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_GETTER(), n)
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_GETTER_DEF(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	this_.priv_AstValidator_validateObjectLiteralKeyName_874(n)
	var function pkg_rhino.Node = n.GetFirstChild()
	_ = function
	this_.priv_AstValidator_validateFunctionExpression_805(function)
	if !(j2gort.CheckNotNullString(function.GetFirstChild().GetString()) == "") {
		this_.priv_AstValidator_violation_886("Expected unnamed function expression.", n)
	}
	var functionParams pkg_rhino.Node = function.GetSecondChild()
	_ = functionParams
	if functionParams.HasChildren() {
		this_.priv_AstValidator_violation_886("get methods must not have parameters.", n)
	}
}

// Java method: AstValidator.validateObjectLitSetKey(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateObjectLitSetKey_866(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_SETTER(), n)
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_SETTER_DEF(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	this_.priv_AstValidator_validateObjectLiteralKeyName_874(n)
	var function pkg_rhino.Node = n.GetFirstChild()
	_ = function
	this_.priv_AstValidator_validateFunctionExpression_805(function)
	if !(j2gort.CheckNotNullString(function.GetFirstChild().GetString()) == "") {
		this_.priv_AstValidator_violation_886("Expected unnamed function expression.", n)
	}
	var functionParams pkg_rhino.Node = function.GetSecondChild()
	_ = functionParams
	if !(functionParams.HasOneChild()) {
		this_.priv_AstValidator_violation_886("set methods must have exactly one parameter.", n)
	}
}

// Java method: AstValidator.validateObjectLitStringKey(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateObjectLitStringKey_867(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_STRING_KEY(), n)
	this_.priv_AstValidator_validateObjectLiteralKeyName_874(n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_890(n, int32(1))
	this_.ValidateExpression(n.GetFirstChild())
	if n.GetBooleanProp(pkg_rhino.StaticGet_Node_IS_SHORTHAND_PROPERTY()) {
		this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES(), n)
	}
}

// Java method: AstValidator.validateObjectPatternStringKey(com.google.javascript.rhino.Token type, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateObjectPatternStringKey_868(type_ pkg_rhino.Token, n pkg_rhino.Node) {
	_ = type_
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_STRING_KEY(), n)
	this_.priv_AstValidator_validateObjectLiteralKeyName_874(n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_890(n, int32(1))
	var c pkg_rhino.Node = n.GetFirstChild()
	_ = c
	switch c.TokenSwitchOrd_() {
	case int32(134) /* Token.DEFAULT_VALUE */ :
		this_.priv_AstValidator_validateDefaultValue_811(type_, c)
	default:
		this_.priv_AstValidator_validateLHS_824(type_, c)
	}
}

// Java method: AstValidator.validateObjectLitComputedPropKey(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateObjectLitComputedPropKey_869(n pkg_rhino.Node) {
	_ = n
	if !(n.GetIsInClosureUnawareSubtree() && NodeUtil_IsGetOrSetKey(n)) {
		this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_COMPUTED_PROPERTIES(), n)
	}
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_COMPUTED_PROP(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	this_.ValidateExpression(n.GetFirstChild())
	this_.ValidateExpression(n.GetLastChild())
}

// Java method: AstValidator.validateObjectPatternComputedPropKey(com.google.javascript.rhino.Token type, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateObjectPatternComputedPropKey_870(type_ pkg_rhino.Token, n pkg_rhino.Node) {
	_ = type_
	_ = n
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_COMPUTED_PROPERTIES(), n)
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_COMPUTED_PROP(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	this_.ValidateExpression(n.GetFirstChild())
	if n.GetLastChild().IsDefaultValue() {
		this_.priv_AstValidator_validateDefaultValue_811(type_, n.GetLastChild())
	} else {
		this_.priv_AstValidator_validateLHS_824(n.GetLastChild().GetToken(), n.GetLastChild())
	}
}

// Java method: AstValidator.validateComputedPropClassMethod(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateComputedPropClassMethod_871(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_COMPUTED_PROPERTIES(), n)
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_COMPUTED_PROP(), n)
	this_.ValidateExpression(n.GetFirstChild())
	if n.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_VARIABLE()) {
		this_.priv_AstValidator_validateProperties_895(n)
		this_.priv_AstValidator_validateChildCount_890(n, int32(1))
	} else {
		this_.priv_AstValidator_validateProperties_895(n)
		this_.priv_AstValidator_validateChildCount_890(n, int32(2))
		this_.priv_AstValidator_validateFunctionExpression_805(n.GetLastChild())
		if n.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_GETTER()) {
			this_.priv_AstValidator_validateObjectLitComputedPropGetKey_872(n)
			this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_CLASS_GETTER_SETTER(), n)
		} else if n.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_SETTER()) {
			this_.priv_AstValidator_validateObjectLitComputedPropSetKey_873(n)
			this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_CLASS_GETTER_SETTER(), n)
		}
	}
}

// Java method: AstValidator.validateObjectLitComputedPropGetKey(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateObjectLitComputedPropGetKey_872(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_COMPUTED_PROPERTIES(), n)
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_GETTER(), n)
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_COMPUTED_PROP(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	var function pkg_rhino.Node = n.GetLastChild()
	_ = function
	this_.priv_AstValidator_validateFunctionExpression_805(function)
	if !(j2gort.CheckNotNullString(function.GetFirstChild().GetString()) == "") {
		this_.priv_AstValidator_violation_886("Expected unnamed function expression.", n)
	}
	var functionParams pkg_rhino.Node = function.GetSecondChild()
	_ = functionParams
	if functionParams.HasChildren() {
		this_.priv_AstValidator_violation_886("get methods must not have parameters.", n)
	}
}

// Java method: AstValidator.validateObjectLitComputedPropSetKey(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateObjectLitComputedPropSetKey_873(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_COMPUTED_PROPERTIES(), n)
	this_.priv_AstValidator_validateFeature_894(pkg_parser.StaticGet_FeatureSet_Feature_SETTER(), n)
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_COMPUTED_PROP(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	var function pkg_rhino.Node = n.GetLastChild()
	_ = function
	this_.priv_AstValidator_validateFunctionExpression_805(function)
	if !(j2gort.CheckNotNullString(function.GetFirstChild().GetString()) == "") {
		this_.priv_AstValidator_violation_886("Expected unnamed function expression.", n)
	}
	var functionParams pkg_rhino.Node = function.GetSecondChild()
	_ = functionParams
	if !(functionParams.HasOneChild()) {
		this_.priv_AstValidator_violation_886("set methods must have exactly one parameter.", n)
	}
}

// Java method: AstValidator.validateObjectLiteralKeyName(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateObjectLiteralKeyName_874(n pkg_rhino.Node) {
	_ = n
	if n.IsQuotedStringKey() {
		func() {
			defer func() {
				if rec_30 := recover(); rec_30 != nil {
					jex_31 := javalang.RecoverException(rec_30)
					_ = jex_31
					if j2gort.InstanceOf[javalang.IUnsupportedOperationException](jex_31) {
						e := j2gort.AsIface[javalang.IUnsupportedOperationException](jex_31)
						_ = e
						this_.priv_AstValidator_violation_886(("getString failed for" + j2gort.StrVal(n.GetToken())), n)
					} else {
						panic(rec_30)
					}
				}
			}()
			n.GetString()
		}()
	} else {
		this_.priv_AstValidator_validateNonEmptyString_797(n)
	}
}

// Java method: AstValidator.validateIncDecOp(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateIncDecOp_875(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_890(n, int32(1))
	this_.priv_AstValidator_validateAssignmentOpTarget_852(n.GetFirstChild(), n.GetToken())
}

// Java method: AstValidator.validateUnaryOp(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateUnaryOp_876(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_890(n, int32(1))
	this_.ValidateExpression(n.GetFirstChild())
}

// Java method: AstValidator.validateBinaryOp(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateBinaryOp_877(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_890(n, int32(2))
	this_.ValidateExpression(n.GetFirstChild())
	this_.ValidateExpression(n.GetLastChild())
}

// Java method: AstValidator.validateTrinaryOp(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateTrinaryOp_878(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_890(n, int32(3))
	var first pkg_rhino.Node = n.GetFirstChild()
	_ = first
	this_.ValidateExpression(first)
	this_.ValidateExpression(first.GetNext())
	this_.ValidateExpression(n.GetLastChild())
}

// Java method: AstValidator.validateNamedType(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateNamedType_879(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_NAMED_TYPE(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	this_.priv_AstValidator_validateName_800(n.GetFirstChild())
}

// Java method: AstValidator.validateTypeAlias(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateTypeAlias_880(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_TYPE_ALIAS(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
}

// Java method: AstValidator.validateAmbientDeclaration(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateAmbientDeclaration_881(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_DECLARE(), n)
	this_.priv_AstValidator_validateAmbientDeclarationHelper_882(n.GetFirstChild())
}

// Java method: AstValidator.validateAmbientDeclarationHelper(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateAmbientDeclarationHelper_882(n pkg_rhino.Node) {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */ :
		this_.priv_AstValidator_validateNameDeclarationHelper_822(n.GetParent(), n.GetToken(), n)
	case int32(73) /* Token.FUNCTION */ :
		this_.priv_AstValidator_validateFunctionSignature_806(n)
	case int32(105) /* Token.CLASS */ :
		this_.priv_AstValidator_validateClassDeclaration_784(n, true)
	case int32(167) /* Token.ENUM */ :
		this_.priv_AstValidator_validateEnum_781(n)
	case int32(175) /* Token.NAMESPACE */ :
		this_.priv_AstValidator_validateNamespace_883(n, true)
	case int32(170) /* Token.TYPE_ALIAS */ :
		this_.priv_AstValidator_validateTypeAlias_880(n)
	case int32(120) /* Token.EXPORT */ :
		this_.priv_AstValidator_validateExport_770(n, true)
	default:
	}
}

// Java method: AstValidator.validateNamespace(com.google.javascript.rhino.Node n, boolean isAmbient)
func (this_ *AstValidator) priv_AstValidator_validateNamespace_883(n pkg_rhino.Node, isAmbient bool) {
	_ = n
	_ = isAmbient
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_NAMESPACE(), n)
	this_.priv_AstValidator_validateProperties_895(n)
	this_.priv_AstValidator_validateChildCount_889(n)
	this_.priv_AstValidator_validateNamespaceName_884(n.GetFirstChild())
	this_.priv_AstValidator_validateNamespaceElements_885(n.GetLastChild(), isAmbient)
}

// Java method: AstValidator.validateNamespaceName(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateNamespaceName_884(n pkg_rhino.Node) {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		this_.priv_AstValidator_validateName_800(n)
	case int32(26) /* Token.GETPROP */ :
		this_.priv_AstValidator_validateGetProp_855(n)
	default:
	}
}

// Java method: AstValidator.validateNamespaceElements(com.google.javascript.rhino.Node n, boolean isAmbient)
func (this_ *AstValidator) priv_AstValidator_validateNamespaceElements_885(n pkg_rhino.Node, isAmbient bool) {
	_ = n
	_ = isAmbient
	this_.priv_AstValidator_validateNodeType_888(pkg_rhino.StaticGet_Token_NAMESPACE_ELEMENTS(), n)
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			if isAmbient {
				this_.priv_AstValidator_validateAmbientDeclarationHelper_882(c)
			} else {
				this_.ValidateStatement_1(c)
			}
		}
	}
}

// Java method: AstValidator.violation(java.lang.String message, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_violation_886(message string, n pkg_rhino.Node) {
	_ = message
	_ = n
	this_.violationHandler.HandleViolation(message, n)
}

// Java method: AstValidator.validateFirstNodeOfOptChain(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateFirstNodeOfOptChain_887(n pkg_rhino.Node) {
	_ = n
	if !(NodeUtil_IsOptChainNode(n.GetFirstChild())) {
		if !(n.IsOptionalChainStart()) {
			this_.priv_AstValidator_violation_886(("Start of optional chain node " + j2gort.StrVal(n.GetToken()) + " is not marked as the start."), n)
		}
	}
}

// Java method: AstValidator.validateNodeType(com.google.javascript.rhino.Token type, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateNodeType_888(type_ pkg_rhino.Token, n pkg_rhino.Node) {
	_ = type_
	_ = n
	if n.GetToken() != type_ {
		this_.priv_AstValidator_violation_886(("Expected " + j2gort.StrVal(type_) + " but was " + j2gort.StrVal(n.GetToken())), n)
	}
}

// Java method: AstValidator.validateChildCount(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateChildCount_889(n pkg_rhino.Node) {
	_ = n
	var expectedArity int32 = pkg_rhino.Token_Arity(n.GetToken())
	_ = expectedArity
	if expectedArity != int32(-1) {
		this_.priv_AstValidator_validateChildCount_890(n, expectedArity)
	}
}

// Java method: AstValidator.validateChildCount(com.google.javascript.rhino.Node n, int expected)
func (this_ *AstValidator) priv_AstValidator_validateChildCount_890(n pkg_rhino.Node, expected int32) {
	_ = n
	_ = expected
	var count int32 = n.GetChildCount()
	_ = count
	if expected != count {
		this_.priv_AstValidator_violation_886(("Expected " + j2gort.StrValInt32(expected) + " children, but was " + j2gort.StrValInt32(count)), n)
	}
}

// Java method: AstValidator.validateChildCountIn(com.google.javascript.rhino.Node n, int min, int max)
func (this_ *AstValidator) priv_AstValidator_validateChildCountIn_891(n pkg_rhino.Node, min_ int32, max_ int32) {
	_ = n
	_ = min_
	_ = max_
	if max_ == min_ {
		this_.priv_AstValidator_validateChildCount_890(n, min_)
		return
	}
	var count int32 = n.GetChildCount()
	_ = count
	if (count < min_) || (count > max_) {
		this_.priv_AstValidator_violation_886(("Expected child count in [" + j2gort.StrValInt32(min_) + ", " + j2gort.StrValInt32(max_) + "], but was " + j2gort.StrValInt32(count)), n)
	}
}

// Java method: AstValidator.validateMinimumChildCount(com.google.javascript.rhino.Node n, int i)
func (this_ *AstValidator) priv_AstValidator_validateMinimumChildCount_892(n pkg_rhino.Node, i int32) {
	_ = n
	_ = i
	var valid bool = false
	_ = valid
	if i == int32(1) {
		valid = n.HasChildren()
	} else if i == int32(2) {
		valid = n.HasMoreThanOneChild()
	} else {
		valid = (n.GetChildCount() >= i)
	}
	if !(valid) {
		this_.priv_AstValidator_violation_886(("Expected at least " + j2gort.StrValInt32(i) + " children, but was " + j2gort.StrValInt32(n.GetChildCount())), n)
	}
}

// Java method: AstValidator.validateMaximumChildCount(com.google.javascript.rhino.Node n, int i)
func (this_ *AstValidator) priv_AstValidator_validateMaximumChildCount_893(n pkg_rhino.Node, i int32) {
	_ = n
	_ = i
	var valid bool = false
	_ = valid
	if i == int32(1) {
		valid = !(n.HasMoreThanOneChild())
	} else if i == int32(-1) {
		valid = true
	} else {
		valid = (n.GetChildCount() <= i)
	}
	if !(valid) {
		this_.priv_AstValidator_violation_886(("Expected no more than " + j2gort.StrValInt32(i) + " children, but was " + j2gort.StrValInt32(n.GetChildCount())), n)
	}
}

// Java method: AstValidator.validateFeature(com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature feature, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateFeature_894(feature *pkg_parser.FeatureSet_Feature, n pkg_rhino.Node) {
	_ = feature
	_ = n
	var allowbleFeatures *pkg_parser.FeatureSet = this_.compiler.GetAllowableFeatures()
	_ = allowbleFeatures
	if !(n.IsFromExterns()) && !(allowbleFeatures.Has(feature)) {
		this_.priv_AstValidator_violation_886(("AST should not contain " + j2gort.StrVal(feature)), n)
	}
	if !(this_.isScriptFeatureValidationEnabled) || (this_.currentScript == 0) {
		return
	}
	var scriptFeatures *pkg_parser.FeatureSet = NodeUtil_GetFeatureSetOfScript(this_.currentScript)
	_ = scriptFeatures
	if (scriptFeatures == nil) || !(NodeUtil_GetFeatureSetOfScript(this_.currentScript).Has(feature)) {
		this_.priv_AstValidator_violation_886(("SCRIPT node should be marked as containing feature " + j2gort.StrVal(feature)), this_.currentScript)
	}
}

// Java method: AstValidator.validateProperties(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateProperties_895(n pkg_rhino.Node) {
	_ = n
	if n.HasPropertyViolations() {
		n.ValidateProperties(j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_32 any) {
			errorMessage := j2gort.CastString(lp_32)
			_ = errorMessage
			this_.priv_AstValidator_violation_886(errorMessage, n)
		}))
	}
}

// Java method: AstValidator.validateRequiredInlinings(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateRequiredInlinings_896(n pkg_rhino.Node) {
	_ = n
	if !(this_.shouldValidateRequiredInlinings) {
		return
	}
	var jsdocInfo *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = jsdocInfo
	if (jsdocInfo != nil) && jsdocInfo.IsRequireInlining() {
		this_.priv_AstValidator_violation_886("@requireInlining node failed to be inlined.", n)
	}
}

// Java interface: com.google.javascript.jscomp.AstValidator.ViolationHandler
type AstValidator_ViolationHandler interface {
	Implements_jscomp_AstValidator_ViolationHandler()
	IsJavaNil_() bool
	HandleViolation(string, pkg_rhino.Node)
}

func ToIface_AstValidator_ViolationHandler[T interface {
	AstValidator_ViolationHandler
	comparable
}](v T) AstValidator_ViolationHandler {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AstValidator.ViolationHandler", "com.google.javascript.jscomp.AstValidator$ViolationHandler", func(v_ any) bool {
		return j2gort.InstanceOf[AstValidator_ViolationHandler](v_)
	})
}

// Java enum: com.google.javascript.jscomp.AstValidator.TypeInfoValidation
type AstValidator_TypeInfoValidation uint8

const (
	enum_AstValidator_TypeInfoValidation_JSTYPE AstValidator_TypeInfoValidation = 1
	enum_AstValidator_TypeInfoValidation_COLOR  AstValidator_TypeInfoValidation = 2
	enum_AstValidator_TypeInfoValidation_NONE   AstValidator_TypeInfoValidation = 3
)

var enumBaseTable_AstValidator_TypeInfoValidation = [...]j2gort.EnumBase{
	{Name_: "JSTYPE", Ordinal_: 0},
	{Name_: "COLOR", Ordinal_: 1},
	{Name_: "NONE", Ordinal_: 2},
}
var enumValues_AstValidator_TypeInfoValidation = []AstValidator_TypeInfoValidation{1, 2, 3}

func StaticGet_AstValidator_TypeInfoValidation_JSTYPE() AstValidator_TypeInfoValidation {
	return enum_AstValidator_TypeInfoValidation_JSTYPE
}
func StaticGet_AstValidator_TypeInfoValidation_COLOR() AstValidator_TypeInfoValidation {
	return enum_AstValidator_TypeInfoValidation_COLOR
}
func StaticGet_AstValidator_TypeInfoValidation_NONE() AstValidator_TypeInfoValidation {
	return enum_AstValidator_TypeInfoValidation_NONE
}

func AstValidator_TypeInfoValidation_Values() []AstValidator_TypeInfoValidation {
	return j2gort.CloneSlice(enumValues_AstValidator_TypeInfoValidation)
}

func AstValidator_TypeInfoValidation_ValueOf(name string) AstValidator_TypeInfoValidation {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_AstValidator_TypeInfoValidation {
		if enumBaseTable_AstValidator_TypeInfoValidation[i_].Name_ == name {
			return AstValidator_TypeInfoValidation(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.AstValidator.TypeInfoValidation." + name))
}

func (this_ AstValidator_TypeInfoValidation) IsJavaNil_() bool                              { return this_ == 0 }
func (AstValidator_TypeInfoValidation) HasIdentityEquals_()                                 {}
func (AstValidator_TypeInfoValidation) Implements_j2gort_Ext_java_lang_Enum()               {}
func (AstValidator_TypeInfoValidation) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ AstValidator_TypeInfoValidation) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (AstValidator_TypeInfoValidation) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AstValidator$TypeInfoValidation", "TypeInfoValidation"
}
func (this_ AstValidator_TypeInfoValidation) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_AstValidator_TypeInfoValidation[this_-1]
}
func (this_ AstValidator_TypeInfoValidation) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ AstValidator_TypeInfoValidation) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ AstValidator_TypeInfoValidation) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_AstValidator_TypeInfoValidation[this_-1].Name_
}
func (this_ AstValidator_TypeInfoValidation) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_AstValidator_TypeInfoValidation[this_-1].Name_
}
func (this_ AstValidator_TypeInfoValidation) Equals(other_ any) bool {
	o_, ok_ := other_.(AstValidator_TypeInfoValidation)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ AstValidator_TypeInfoValidation) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_AstValidator_TypeInfoValidation[this_-1].HashCode()
}
func (this_ AstValidator_TypeInfoValidation) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_AstValidator_TypeInfoValidation[this_-1].CompareTo(other_)
}
func (this_ AstValidator_TypeInfoValidation) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_AstValidator_TypeInfoValidation[this_-1].Name_)
}
func init() {
	vals_ := []any{AstValidator_TypeInfoValidation(1), AstValidator_TypeInfoValidation(2), AstValidator_TypeInfoValidation(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.AstValidator.TypeInfoValidation", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.AstValidator$TypeInfoValidation", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AstValidator.TypeInfoValidation", "com.google.javascript.jscomp.AstValidator$TypeInfoValidation", func(v_ any) bool {
		return j2gort.InstanceOf[AstValidator_TypeInfoValidation](v_)
	})
}

// Java class:
type AstValidator__Anon_2 struct {
	outer_  *AstValidator
	idHash_ atomic.Int32
}

func (this_ *AstValidator__Anon_2) BindSelf_(s any) {
	_ = s
}
func (this_ *AstValidator__Anon_2) IsJavaNil_() bool { return this_ == nil }
func (*AstValidator__Anon_2) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AstValidator$1", ""
}
func (this_ *AstValidator__Anon_2) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AstValidator__Anon_2) HasIdentityEquals_()                              {}
func (this_ *AstValidator__Anon_2) SetOuter_(o *AstValidator)                  { this_.outer_ = o }
func (this_ *AstValidator__Anon_2) GetOuter_() *AstValidator                   { return this_.outer_ }
func (*AstValidator__Anon_2) Implements_jscomp_AstValidator_ViolationHandler() {}

func (this_ *AstValidator__Anon_2) initFields_AstValidator__Anon_2() {
	_ = this_
}

func NewAstValidator__Anon_2(outer_ *AstValidator) *AstValidator__Anon_2 {
	this_ := &AstValidator__Anon_2{}
	this_.outer_ = outer_
	this_.initFields_AstValidator__Anon_2()
	return this_
}

// Java method: .handleViolation(java.lang.String message, com.google.javascript.rhino.Node n)
func (this_ *AstValidator__Anon_2) HandleViolation(message string, n pkg_rhino.Node) {
	_ = message
	_ = n
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString((j2gort.StrValString(message) + ". Reference node:\n" + j2gort.StrValString(n.ToStringTree()) + "\n Parent node:\n" + j2gort.StrValString(func() string {
		if n.HasParent() {
			return n.GetParent().ToStringTree()
		}
		return " no parent "
	}()))))))
}
