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
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_ROOT(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_892(n, int32(2))
	this_.ValidateCodeRoot(n.GetFirstChild())
	this_.ValidateCodeRoot(n.GetLastChild())
}

// Java method: AstValidator.validateCodeRoot(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) ValidateCodeRoot(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_ROOT(), n)
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
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_SCRIPT(), n)
	this_.priv_AstValidator_validateHasSourceName_795(n)
	this_.priv_AstValidator_validateHasInputId_796(n)
	this_.currentScript = n
	if n.HasChildren() && n.GetFirstChild().IsModuleBody() {
		this_.priv_AstValidator_validateProperties_897(n)
		this_.priv_AstValidator_validateChildCount_892(n, int32(1))
		this_.ValidateModuleContents(n.GetFirstChild())
	} else {
		this_.ValidateStatements(n.GetFirstChild())
	}
	if this_.isScriptFeatureValidationEnabled {
		this_.priv_AstValidator_validateScriptFeatureSet_752(n)
	}
}

// Java method: AstValidator.validateScriptFeatureSet(com.google.javascript.rhino.Node scriptNode)
func (this_ *AstValidator) priv_AstValidator_validateScriptFeatureSet_752(scriptNode pkg_rhino.Node) {
	_ = scriptNode
	if !(scriptNode.IsScript()) {
		this_.priv_AstValidator_violation_888("Not a script node", scriptNode)
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
			this_.priv_AstValidator_violation_888(("SCRIPT node contains these unallowable features:" + j2gort.StrVal(differentFeatures.GetFeatures())), this_.currentScript)
		}
	}
}

// Java method: AstValidator.validateModuleContents(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) ValidateModuleContents(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_MODULE_BODY(), n)
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
		this_.priv_AstValidator_validateLabel_797(n)
	case int32(90) /* Token.BLOCK */ :
		this_.priv_AstValidator_validateBlock_794(n)
	case int32(73) /* Token.FUNCTION */ :
		if isAmbient {
			this_.priv_AstValidator_validateFunctionSignature_808(n)
		} else {
			this_.priv_AstValidator_validateFunctionStatement_806(n)
		}
	case int32(85) /* Token.WITH */ :
		this_.priv_AstValidator_validateWith_835(n)
	case int32(80) /* Token.FOR */ :
		this_.priv_AstValidator_validateFor_830(n)
	case int32(81) /* Token.FOR_IN */ :
		this_.priv_AstValidator_validateForIn_831(n)
	case int32(112) /* Token.FOR_OF */ :
		this_.priv_AstValidator_validateForOf_832(n)
	case int32(113) /* Token.FOR_AWAIT_OF */ :
		this_.priv_AstValidator_validateForAwaitOf_833(n)
	case int32(78) /* Token.WHILE */ :
		this_.priv_AstValidator_validateWhile_836(n)
	case int32(79) /* Token.DO */ :
		this_.priv_AstValidator_validateDo_837(n)
	case int32(75) /* Token.SWITCH */ :
		this_.priv_AstValidator_validateSwitch_847(n)
	case int32(74) /* Token.IF */ :
		this_.priv_AstValidator_validateIf_838(n)
	case int32(97) /* Token.CONST */, int32(84) /* Token.VAR */, int32(111) /* Token.LET */ :
		this_.priv_AstValidator_validateNameDeclarationHelper_824(n, n.GetToken(), n)
	case int32(93) /* Token.EXPR_RESULT */ :
		this_.priv_AstValidator_validateExprStmt_839(n)
	case int32(0) /* Token.RETURN */ :
		this_.priv_AstValidator_validateReturn_840(n)
	case int32(43) /* Token.THROW */ :
		this_.priv_AstValidator_validateThrow_841(n)
	case int32(48) /* Token.TRY */ :
		this_.priv_AstValidator_validateTry_844(n)
	case int32(82) /* Token.BREAK */ :
		this_.priv_AstValidator_validateBreak_842(n)
	case int32(83) /* Token.CONTINUE */ :
		this_.priv_AstValidator_validateContinue_843(n)
	case int32(88) /* Token.EMPTY */, int32(98) /* Token.DEBUGGER */ :
		this_.priv_AstValidator_validateProperties_897(n)
		this_.priv_AstValidator_validateChildless_851(n)
	case int32(105) /* Token.CLASS */ :
		this_.priv_AstValidator_validateClassDeclaration_786(n, isAmbient)
	case int32(116) /* Token.IMPORT */ :
		this_.priv_AstValidator_validateImport_769(n)
	case int32(120) /* Token.EXPORT */ :
		this_.priv_AstValidator_validateExport_772(n, isAmbient)
	case int32(164) /* Token.INTERFACE */ :
		this_.priv_AstValidator_validateInterface_779(n)
	case int32(167) /* Token.ENUM */ :
		this_.priv_AstValidator_validateEnum_783(n)
	case int32(170) /* Token.TYPE_ALIAS */ :
		this_.priv_AstValidator_validateTypeAlias_882(n)
	case int32(171) /* Token.DECLARE */ :
		this_.priv_AstValidator_validateAmbientDeclaration_883(n)
	case int32(175) /* Token.NAMESPACE */ :
		this_.priv_AstValidator_validateNamespace_885(n, isAmbient)
	case int32(123) /* Token.MODULE_BODY */ :
		if this_.compiler.GetOptions().GetSkipNonTranspilationPasses() {
			if !(this_.compiler.GetOptions().ShouldWrapGoogModulesForWhitespaceOnly()) {
				this_.ValidateModuleContents(n)
				return
			}
		}
		this_.priv_AstValidator_violation_888(("Expected statement but was " + j2gort.StrVal(n.GetToken()) + "."), n)
	default:
		if n.IsModuleBody() && this_.compiler.GetOptions().GetSkipNonTranspilationPasses() {
			if !(!(this_.compiler.GetOptions().ShouldWrapGoogModulesForWhitespaceOnly())) {
				j2gort.CheckStateStr(false, "Modules can exist in transpiler only if setWrapGoogModulesForWhitespaceOnly is false")
			}
			this_.ValidateModuleContents(n)
			return
		}
		this_.priv_AstValidator_violation_888(("Expected statement but was " + j2gort.StrVal(n.GetToken()) + "."), n)
	}
}

// Java method: AstValidator.validateExpression(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) ValidateExpression(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateTypeInformation_760(n)
	this_.priv_AstValidator_validateRequiredInlinings_898(n)
	switch n.TokenSwitchOrd_() {
	case int32(135) /* Token.NEW_TARGET */ :
		this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_NEW_TARGET(), n)
		this_.priv_AstValidator_validateProperties_897(n)
		this_.priv_AstValidator_validateChildless_851(n)
	case int32(136) /* Token.IMPORT_META */ :
		this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_IMPORT_META(), n)
		this_.priv_AstValidator_validateProperties_897(n)
		this_.priv_AstValidator_validateChildless_851(n)
	case int32(38) /* Token.FALSE */, int32(36) /* Token.NULL */, int32(37) /* Token.THIS */, int32(39) /* Token.TRUE */ :
		this_.priv_AstValidator_validateProperties_897(n)
		this_.priv_AstValidator_validateChildless_851(n)
	case int32(24) /* Token.DELPROP */, int32(21) /* Token.POS */, int32(22) /* Token.NEG */, int32(19) /* Token.NOT */, int32(25) /* Token.TYPEOF */, int32(87) /* Token.VOID */, int32(20) /* Token.BITNOT */, int32(101) /* Token.CAST */ :
		this_.priv_AstValidator_validateUnaryOp_878(n)
	case int32(71) /* Token.INC */, int32(72) /* Token.DEC */ :
		this_.priv_AstValidator_validateIncDecOp_877(n)
	case int32(51) /* Token.ASSIGN */ :
		this_.priv_AstValidator_validateAssignmentExpression_852(n)
	case int32(63) /* Token.ASSIGN_EXPONENT */ :
		this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_EXPONENT_OP(), n)
		this_.priv_AstValidator_validateCompoundAssignmentExpression_853(n)
	case int32(52) /* Token.ASSIGN_BITOR */, int32(53) /* Token.ASSIGN_BITXOR */, int32(54) /* Token.ASSIGN_BITAND */, int32(55) /* Token.ASSIGN_LSH */, int32(56) /* Token.ASSIGN_RSH */, int32(57) /* Token.ASSIGN_URSH */, int32(58) /* Token.ASSIGN_ADD */, int32(59) /* Token.ASSIGN_SUB */, int32(60) /* Token.ASSIGN_MUL */, int32(61) /* Token.ASSIGN_DIV */, int32(62) /* Token.ASSIGN_MOD */ :
		this_.priv_AstValidator_validateCompoundAssignmentExpression_853(n)
	case int32(66) /* Token.ASSIGN_COALESCE */ :
		this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_NULL_COALESCE_OP(), n)
		this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_LOGICAL_ASSIGNMENT(), n)
		this_.priv_AstValidator_validateCompoundAssignmentExpression_853(n)
	case int32(64) /* Token.ASSIGN_OR */, int32(65) /* Token.ASSIGN_AND */ :
		this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_LOGICAL_ASSIGNMENT(), n)
		this_.priv_AstValidator_validateCompoundAssignmentExpression_853(n)
	case int32(67) /* Token.HOOK */ :
		this_.priv_AstValidator_validateTrinaryOp_880(n)
	case int32(35) /* Token.STRINGLIT */ :
		this_.priv_AstValidator_validateStringLit_861(n)
	case int32(33) /* Token.NUMBER */ :
		this_.priv_AstValidator_validateNumber_862(n)
	case int32(34) /* Token.BIGINT */ :
		this_.priv_AstValidator_validateBigInt_863(n)
	case int32(32) /* Token.NAME */ :
		this_.priv_AstValidator_validateName_802(n)
	case int32(18) /* Token.EXPONENT */ :
		this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_EXPONENT_OP(), n)
		this_.priv_AstValidator_validateBinaryOp_879(n)
	case int32(70) /* Token.COALESCE */ :
		this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_NULL_COALESCE_OP(), n)
		this_.priv_AstValidator_validateBinaryOp_879(n)
	case int32(50) /* Token.COMMA */, int32(68) /* Token.OR */, int32(69) /* Token.AND */, int32(1) /* Token.BITOR */, int32(2) /* Token.BITXOR */, int32(3) /* Token.BITAND */, int32(4) /* Token.EQ */, int32(5) /* Token.NE */, int32(40) /* Token.SHEQ */, int32(41) /* Token.SHNE */, int32(6) /* Token.LT */, int32(8) /* Token.GT */, int32(7) /* Token.LE */, int32(9) /* Token.GE */, int32(45) /* Token.INSTANCEOF */, int32(44) /* Token.IN */, int32(10) /* Token.LSH */, int32(11) /* Token.RSH */, int32(12) /* Token.URSH */, int32(14) /* Token.SUB */, int32(13) /* Token.ADD */, int32(15) /* Token.MUL */, int32(17) /* Token.MOD */, int32(16) /* Token.DIV */ :
		this_.priv_AstValidator_validateBinaryOp_879(n)
	case int32(27) /* Token.GETELEM */ :
		this_.priv_AstValidator_validateGetElem_855(n)
	case int32(30) /* Token.OPTCHAIN_GETELEM */ :
		this_.priv_AstValidator_validateOptChainGetElem_856(n)
	case int32(26) /* Token.GETPROP */ :
		this_.priv_AstValidator_validateGetProp_857(n)
	case int32(29) /* Token.OPTCHAIN_GETPROP */ :
		this_.priv_AstValidator_validateOptChainGetProp_858(n)
	case int32(46) /* Token.ARRAYLIT */ :
		this_.priv_AstValidator_validateArrayLit_864(n)
	case int32(47) /* Token.OBJECTLIT */ :
		this_.priv_AstValidator_validateObjectLit_865(n)
	case int32(42) /* Token.REGEXP */ :
		this_.priv_AstValidator_validateRegExpLit_860(n)
	case int32(28) /* Token.CALL */ :
		this_.priv_AstValidator_validateCall_814(n)
	case int32(31) /* Token.OPTCHAIN_CALL */ :
		this_.priv_AstValidator_validateOptChainCall_815(n)
	case int32(23) /* Token.NEW */ :
		this_.priv_AstValidator_validateNew_823(n)
	case int32(73) /* Token.FUNCTION */ :
		this_.priv_AstValidator_validateRequiredInlinings_898(n)
		this_.priv_AstValidator_validateFunctionExpression_807(n)
	case int32(105) /* Token.CLASS */ :
		this_.priv_AstValidator_validateClass_787(n)
	case int32(131) /* Token.TEMPLATELIT */ :
		this_.priv_AstValidator_validateTemplateLit_776(n)
	case int32(130) /* Token.TAGGED_TEMPLATELIT */ :
		this_.priv_AstValidator_validateTaggedTemplateLit_775(n)
	case int32(114) /* Token.YIELD */ :
		this_.priv_AstValidator_validateYield_764(n)
	case int32(115) /* Token.AWAIT */ :
		this_.priv_AstValidator_validateAwait_766(n)
	case int32(124) /* Token.DYNAMIC_IMPORT */ :
		this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_DYNAMIC_IMPORT(), n)
		this_.priv_AstValidator_validateUnaryOp_878(n)
	default:
		this_.priv_AstValidator_violation_888(("Expected expression but was " + j2gort.StrVal(n.GetToken())), n)
	}
}

// Java method: AstValidator.validatePseudoExpression(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Token allowed)
func (this_ *AstValidator) priv_AstValidator_validatePseudoExpression_758(n pkg_rhino.Node, allowed pkg_rhino.Token) {
	_ = n
	_ = allowed
	this_.priv_AstValidator_validatePseudoExpression_759(n, allowed, pkg_rhino.Token(0))
}

// Java method: AstValidator.validatePseudoExpression(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Token allowed1, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Token allowed2)
func (this_ *AstValidator) priv_AstValidator_validatePseudoExpression_759(n pkg_rhino.Node, allowed1 pkg_rhino.Token, allowed2 pkg_rhino.Token) {
	_ = n
	_ = allowed1
	_ = allowed2
	switch n.TokenSwitchOrd_() {
	case int32(88) /* Token.EMPTY */ :
		if !((allowed1 == pkg_rhino.StaticGet_Token_EMPTY()) || (allowed2 == pkg_rhino.StaticGet_Token_EMPTY())) {
			j2gort.CheckArgumentFmt(false, "Unexpected pseudoexpression %s", j2gort.BoxNode(n))
		}
		this_.priv_AstValidator_validateProperties_897(n)
		this_.priv_AstValidator_validateChildless_851(n)
	case int32(127) /* Token.ITER_SPREAD */ :
		if !((allowed1 == pkg_rhino.StaticGet_Token_ITER_SPREAD()) || (allowed2 == pkg_rhino.StaticGet_Token_ITER_SPREAD())) {
			j2gort.CheckArgumentFmt(false, "Unexpected pseudoexpression %s", j2gort.BoxNode(n))
		}
		this_.priv_AstValidator_validateProperties_897(n)
		this_.priv_AstValidator_validateChildCount_891(n)
		this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_SPREAD_EXPRESSIONS(), n)
		this_.ValidateExpression(n.GetFirstChild())
	default:
		this_.ValidateExpression(n)
	}
}

// Java method: AstValidator.validateTypeInformation(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateTypeInformation_760(n pkg_rhino.Node) {
	_ = n
	if n.GetIsInClosureUnawareSubtree() {
		return
	}
	{
		var iifeCond_1 bool
		{

			re_3 := this_.typeValidationMode
			ae_4 := StaticGet_AstValidator_TypeInfoValidation_NONE()
			if re_3 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_3 == ae_4
		}
		if iifeCond_1 {
			return
		}
	}
	{
		var iifeCond_1 bool
		{

			re_5 := this_.typeValidationMode
			ae_6 := StaticGet_AstValidator_TypeInfoValidation_JSTYPE()
			if re_5 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_5 == ae_6
		}
		if iifeCond_1 {
			var type_ pkg_rhino.JSType = n.GetJSType()
			_ = type_
			if (type_ != nil && !type_.IsJavaNil_()) && !(type_.GetJSTypeBase().IsResolved()) {
				this_.priv_AstValidator_violation_888(("Found unresolved type " + j2gort.StrVal(type_)), n)
			}
		}
	}
	switch n.TokenSwitchOrd_() {
	case int32(28) /* Token.CALL */ :
		if !(n.GetFirstChild().IsSuper()) {
			this_.priv_AstValidator_validateCallType_761(n)
		}
	default:
		this_.priv_AstValidator_expectSomeTypeInformation_762(n)
	}
}

// Java method: AstValidator.validateCallType(com.google.javascript.rhino.Node callNode)
func (this_ *AstValidator) priv_AstValidator_validateCallType_761(callNode pkg_rhino.Node) {
	_ = callNode
	switch (this_.typeValidationMode).Ordinal() {
	case int32(0) /* TypeInfoValidation.JSTYPE */ :
		var callee pkg_rhino.Node = callNode.GetFirstChild()
		_ = callee
		var calleeType pkg_rhino.JSType = callee.GetJSType()
		_ = calleeType
		if calleeType == nil || calleeType.IsJavaNil_() {
			j2gort.CheckNotNullIface(calleeType)
		}
		if calleeType.GetJSTypeBase().IsFunctionType() {
			var calleeFunctionType *pkg_rhino.FunctionType = calleeType.ToMaybeFunctionType()
			_ = calleeFunctionType
			var returnType pkg_rhino.JSType = calleeFunctionType.GetReturnType()
			_ = returnType
			if j2gort.IfaceEquals[pkg_rhino.JSType](callNode.GetJSTypeBeforeCast(), nil) && !(returnType.GetJSTypeBase().IsUnknownType()) {
				this_.priv_AstValidator_expectMatchingTypeInformation_763(callNode, returnType)
			}
		}
	case int32(1) /* TypeInfoValidation.COLOR */ :
		var callee_2 pkg_rhino.Node = callNode.GetFirstChild()
		_ = callee_2
		if callee_2.GetColor() == nil {
			j2gort.CheckNotNullPtr(callee_2.GetColor())
		}
	case int32(2) /* TypeInfoValidation.NONE */ :
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
	}
}

// Java method: AstValidator.expectSomeTypeInformation(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_expectSomeTypeInformation_762(n pkg_rhino.Node) {
	_ = n
	switch (this_.typeValidationMode).Ordinal() {
	case int32(0) /* TypeInfoValidation.JSTYPE */ :
		if j2gort.IfaceEquals[pkg_rhino.JSType](n.GetJSType(), nil) {
			this_.priv_AstValidator_violation_888(("Type information missing\n" + j2gort.StrValString(this_.compiler.ToSource_1_Node(NodeUtil_GetEnclosingStatement(n)))), n)
		}
	case int32(1) /* TypeInfoValidation.COLOR */ :
		if n.GetColor() == nil {
			this_.priv_AstValidator_violation_888(("Color information missing\n" + j2gort.StrValString(this_.compiler.ToSource_1_Node(NodeUtil_GetEnclosingStatement(n)))), n)
		}
	case int32(2) /* TypeInfoValidation.NONE */ :
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
	}
}

// Java method: AstValidator.expectMatchingTypeInformation(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType expectedTypeI)
func (this_ *AstValidator) priv_AstValidator_expectMatchingTypeInformation_763(n pkg_rhino.Node, expectedTypeI pkg_rhino.JSType) {
	_ = n
	_ = expectedTypeI
	var typeI pkg_rhino.JSType = n.GetJSType()
	_ = typeI
	if !(j2gort.ObjEqIface[pkg_rhino.JSType](expectedTypeI, typeI)) {
		this_.priv_AstValidator_violation_888(("Expected type: " + j2gort.StrValString(AstValidator_GetTypeAnnotationString(expectedTypeI)) + " Actual type: " + j2gort.StrValString(AstValidator_GetTypeAnnotationString(typeI))), n)
	}
}

// Java static method: AstValidator.getTypeAnnotationString(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType typeI)
func AstValidator_GetTypeAnnotationString(typeI pkg_rhino.JSType) string {
	_ = typeI
	if typeI == nil || typeI.IsJavaNil_() {
		return "NO TYPE INFORMATION"
	} else {
		return ("{" + j2gort.StrValString(typeI.GetJSTypeBase().ToAnnotationString(pkg_rhino.StaticGet_JSType_Nullability_EXPLICIT())) + "}")
	}
}

// Java method: AstValidator.validateYield(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateYield_764(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_GENERATORS(), n)
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_YIELD(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCountIn_893(n, int32(0), int32(1))
	if n.HasChildren() {
		this_.ValidateExpression(n.GetFirstChild())
	}
	this_.priv_AstValidator_validateYieldWithinGeneratorFunction_765(n)
}

// Java method: AstValidator.validateYieldWithinGeneratorFunction(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateYieldWithinGeneratorFunction_765(n pkg_rhino.Node) {
	_ = n
	var parentFunction pkg_rhino.Node = NodeUtil_GetEnclosingFunction(n)
	_ = parentFunction
	if (parentFunction == 0) || !(parentFunction.IsGeneratorFunction()) {
		this_.priv_AstValidator_violation_888("'yield' expression is not within a generator function", n)
	} else if this_.priv_AstValidator_isInParameterListOfFunction_768(n, parentFunction) {
		this_.priv_AstValidator_violation_888("'yield' expression is not allowed in a parameter list", n)
	}
}

// Java method: AstValidator.validateAwait(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateAwait_766(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_ASYNC_FUNCTIONS(), n)
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_AWAIT(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	this_.ValidateExpression(n.GetFirstChild())
	this_.priv_AstValidator_validateAwaitWithinAsyncFunction_767(n)
}

// Java method: AstValidator.validateAwaitWithinAsyncFunction(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateAwaitWithinAsyncFunction_767(n pkg_rhino.Node) {
	_ = n
	var parentFunction pkg_rhino.Node = NodeUtil_GetEnclosingFunction(n)
	_ = parentFunction
	if parentFunction == 0 {
		if !(NodeUtil_GetEnclosingScript(n).GetBooleanProp(pkg_rhino.StaticGet_Node_ES6_MODULE())) {
			this_.priv_AstValidator_violation_888("'await' expression is not within an async function", n)
		}
		return
	}
	if !(parentFunction.IsAsyncFunction()) {
		this_.priv_AstValidator_violation_888("'await' expression is not within an async function", n)
	} else if this_.priv_AstValidator_isInParameterListOfFunction_768(n, parentFunction) {
		this_.priv_AstValidator_violation_888("'await' expression is not allowed in a parameter list", n)
	}
}

// Java method: AstValidator.isInParameterListOfFunction(com.google.javascript.rhino.Node child, com.google.javascript.rhino.Node functionNode)
func (this_ *AstValidator) priv_AstValidator_isInParameterListOfFunction_768(child pkg_rhino.Node, functionNode pkg_rhino.Node) bool {
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
func (this_ *AstValidator) priv_AstValidator_validateImport_769(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_MODULES(), n)
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_IMPORT(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	if n.GetFirstChild().IsName() {
		this_.priv_AstValidator_validateName_802(n.GetFirstChild())
	} else {
		this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_EMPTY(), n.GetFirstChild())
	}
	var secondChild pkg_rhino.Node = n.GetSecondChild()
	_ = secondChild
	switch secondChild.TokenSwitchOrd_() {
	case int32(117) /* Token.IMPORT_SPECS */ :
		this_.priv_AstValidator_validateImportSpecifiers_770(secondChild)
	case int32(119) /* Token.IMPORT_STAR */ :
		this_.priv_AstValidator_validateNonEmptyString_799(secondChild)
	default:
		this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_EMPTY(), secondChild)
	}
	this_.priv_AstValidator_validateStringLit_861(n.GetChildAtIndex(int32(2)))
}

// Java method: AstValidator.validateImportSpecifiers(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateImportSpecifiers_770(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_IMPORT_SPECS(), n)
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validateImportSpecifier_771(c)
		}
	}
}

// Java method: AstValidator.validateImportSpecifier(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateImportSpecifier_771(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_IMPORT_SPEC(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_892(n, int32(2))
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validateName_802(c)
		}
	}
}

// Java method: AstValidator.validateExport(com.google.javascript.rhino.Node n, boolean isAmbient)
func (this_ *AstValidator) priv_AstValidator_validateExport_772(n pkg_rhino.Node, isAmbient bool) {
	_ = n
	_ = isAmbient
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_MODULES(), n)
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_EXPORT(), n)
	if n.GetBooleanProp(pkg_rhino.StaticGet_Node_EXPORT_ALL_FROM()) {
		this_.priv_AstValidator_validateProperties_897(n)
		this_.priv_AstValidator_validateChildCount_892(n, int32(2))
		this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_EMPTY(), n.GetFirstChild())
		this_.priv_AstValidator_validateStringLit_861(n.GetSecondChild())
	} else if n.GetBooleanProp(pkg_rhino.StaticGet_Node_EXPORT_DEFAULT()) {
		this_.priv_AstValidator_validateProperties_897(n)
		this_.priv_AstValidator_validateChildCount_892(n, int32(1))
		this_.ValidateExpression(n.GetFirstChild())
	} else {
		this_.priv_AstValidator_validateProperties_897(n)
		this_.priv_AstValidator_validateChildCountIn_893(n, int32(1), int32(2))
		if n.GetFirstChild().IsExportSpecs() {
			this_.priv_AstValidator_validateExportSpecifiers_773(n.GetFirstChild())
		} else {
			this_.ValidateStatement_2(n.GetFirstChild(), isAmbient)
		}
		if n.HasTwoChildren() {
			this_.priv_AstValidator_validateStringLit_861(n.GetSecondChild())
		}
	}
}

// Java method: AstValidator.validateExportSpecifiers(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateExportSpecifiers_773(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_EXPORT_SPECS(), n)
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validateExportSpecifier_774(c)
		}
	}
}

// Java method: AstValidator.validateExportSpecifier(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateExportSpecifier_774(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_EXPORT_SPEC(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_892(n, int32(2))
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validateName_802(c)
		}
	}
}

// Java method: AstValidator.validateTaggedTemplateLit(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateTaggedTemplateLit_775(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_TEMPLATE_LITERALS(), n)
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_TAGGED_TEMPLATELIT(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	this_.ValidateExpression(n.GetFirstChild())
	this_.priv_AstValidator_validateTemplateLit_776(n.GetLastChild())
}

// Java method: AstValidator.validateTemplateLit(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateTemplateLit_776(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_TEMPLATE_LITERALS(), n)
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_TEMPLATELIT(), n)
	{
		var child pkg_rhino.Node = n.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			if child.IsTemplateLitString() {
				this_.priv_AstValidator_validateTemplateLitString_777(child)
			} else {
				this_.priv_AstValidator_validateTemplateLitSub_778(child)
			}
		}
	}
}

// Java method: AstValidator.validateTemplateLitString(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateTemplateLitString_777(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_TEMPLATELIT_STRING(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	func() {
		defer func() {
			if rec_7 := recover(); rec_7 != nil {
				jex_8 := javalang.RecoverException(rec_7)
				_ = jex_8
				if j2gort.InstanceOf[javalang.IUnsupportedOperationException](jex_8) {
					e := j2gort.AsIface[javalang.IUnsupportedOperationException](jex_8)
					_ = e
					this_.priv_AstValidator_violation_888("Invalid TEMPLATELIT_STRING node.", n)
				} else {
					panic(rec_7)
				}
			}
		}()
		n.GetRawString()
	}()
}

// Java method: AstValidator.validateTemplateLitSub(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateTemplateLitSub_778(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_TEMPLATELIT_SUB(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	this_.ValidateExpression(n.GetFirstChild())
}

// Java method: AstValidator.validateInterface(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateInterface_779(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_INTERFACE(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	var name pkg_rhino.Node = n.GetFirstChild()
	_ = name
	this_.priv_AstValidator_validateName_802(name)
	var superTypes pkg_rhino.Node = name.GetNext()
	_ = superTypes
	if superTypes.IsEmpty() {
		this_.priv_AstValidator_validateProperties_897(superTypes)
		this_.priv_AstValidator_validateChildless_851(superTypes)
	} else {
		this_.priv_AstValidator_validateInterfaceExtends_780(superTypes)
	}
	this_.priv_AstValidator_validateInterfaceMembers_781(n.GetLastChild())
}

// Java method: AstValidator.validateInterfaceExtends(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateInterfaceExtends_780(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_INTERFACE_EXTENDS(), n)
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validateNamedType_881(c)
		}
	}
}

// Java method: AstValidator.validateInterfaceMembers(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateInterfaceMembers_781(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_INTERFACE_MEMBERS(), n)
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validateInterfaceMember_782(c)
		}
	}
}

// Java method: AstValidator.validateInterfaceMember(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateInterfaceMember_782(n pkg_rhino.Node) {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
		this_.priv_AstValidator_validateProperties_897(n)
		this_.priv_AstValidator_validateChildCount_891(n)
		this_.priv_AstValidator_validateFunctionSignature_808(n.GetFirstChild())
	case int32(172) /* Token.MEMBER_VARIABLE_DEF */ :
		this_.priv_AstValidator_validateProperties_897(n)
		this_.priv_AstValidator_validateChildless_851(n)
	case int32(173) /* Token.INDEX_SIGNATURE */ :
		this_.priv_AstValidator_validateProperties_897(n)
		this_.priv_AstValidator_validateChildCount_891(n)
		var child pkg_rhino.Node = n.GetFirstChild()
		_ = child
		this_.priv_AstValidator_validateProperties_897(child)
		this_.priv_AstValidator_validateChildless_851(child)
	case int32(174) /* Token.CALL_SIGNATURE */ :
		this_.priv_AstValidator_validateProperties_897(n)
		this_.priv_AstValidator_validateChildCount_891(n)
	default:
		this_.priv_AstValidator_violation_888(("Interface contained member of invalid type " + j2gort.StrVal(n.GetToken())), n)
	}
}

// Java method: AstValidator.validateEnum(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateEnum_783(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_ENUM(), n)
	this_.priv_AstValidator_validateName_802(n.GetFirstChild())
	this_.priv_AstValidator_validateEnumMembers_784(n.GetLastChild())
}

// Java method: AstValidator.validateEnumMembers(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateEnumMembers_784(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_ENUM_MEMBERS(), n)
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validateEnumStringKey_785(c)
		}
	}
}

// Java method: AstValidator.validateEnumStringKey(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateEnumStringKey_785(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_STRING_KEY(), n)
	this_.priv_AstValidator_validateObjectLiteralKeyName_876(n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_892(n, int32(0))
}

// Java method: AstValidator.validateClassDeclaration(com.google.javascript.rhino.Node n, boolean isAmbient)
func (this_ *AstValidator) priv_AstValidator_validateClassDeclaration_786(n pkg_rhino.Node, isAmbient bool) {
	_ = n
	_ = isAmbient
	this_.priv_AstValidator_validateClassHelper_788(n, isAmbient)
	this_.priv_AstValidator_validateName_802(n.GetFirstChild())
}

// Java method: AstValidator.validateClass(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateClass_787(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateClassHelper_788(n, false)
}

// Java method: AstValidator.validateClassHelper(com.google.javascript.rhino.Node n, boolean isAmbient)
func (this_ *AstValidator) priv_AstValidator_validateClassHelper_788(n pkg_rhino.Node, isAmbient bool) {
	_ = n
	_ = isAmbient
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_CLASSES(), n)
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_CLASS(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	var name pkg_rhino.Node = n.GetFirstChild()
	_ = name
	if name.IsEmpty() {
		this_.priv_AstValidator_validateProperties_897(name)
		this_.priv_AstValidator_validateChildless_851(name)
	} else {
		this_.priv_AstValidator_validateName_802(name)
	}
	var superClass pkg_rhino.Node = name.GetNext()
	_ = superClass
	if superClass.IsEmpty() {
		this_.priv_AstValidator_validateProperties_897(superClass)
		this_.priv_AstValidator_validateChildless_851(superClass)
	} else {
		this_.ValidateExpression(superClass)
	}
	this_.priv_AstValidator_validateClassMembers_789(n.GetLastChild(), isAmbient)
}

// Java method: AstValidator.validateClassMembers(com.google.javascript.rhino.Node n, boolean isAmbient)
func (this_ *AstValidator) priv_AstValidator_validateClassMembers_789(n pkg_rhino.Node, isAmbient bool) {
	_ = n
	_ = isAmbient
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_CLASS_MEMBERS(), n)
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validateClassMember_790(c, isAmbient)
		}
	}
}

// Java method: AstValidator.validateClassMember(com.google.javascript.rhino.Node n, boolean isAmbient)
func (this_ *AstValidator) priv_AstValidator_validateClassMember_790(n pkg_rhino.Node, isAmbient bool) {
	_ = n
	_ = isAmbient
	switch n.TokenSwitchOrd_() {
	case int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
		this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_MEMBER_DECLARATIONS(), n)
		if n.IsPrivateIdentifier() {
			this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_PRIVATE_ELEMENTS(), n)
		}
		this_.priv_AstValidator_validateObjectLiteralKeyName_876(n)
		this_.priv_AstValidator_validateProperties_897(n)
		this_.priv_AstValidator_validateChildCount_891(n)
		this_.priv_AstValidator_validateMemberFunction_791(n, isAmbient)
	case int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */ :
		this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_CLASS_GETTER_SETTER(), n)
		if n.IsPrivateIdentifier() {
			this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_PRIVATE_ELEMENTS(), n)
		}
		this_.priv_AstValidator_validateObjectLiteralKeyName_876(n)
		this_.priv_AstValidator_validateObjectLitKey_866(n)
		this_.priv_AstValidator_validateProperties_897(n)
		this_.priv_AstValidator_validateChildCount_891(n)
		this_.priv_AstValidator_validateMemberFunction_791(n, isAmbient)
	case int32(172) /* Token.MEMBER_VARIABLE_DEF */ :
		this_.priv_AstValidator_validateProperties_897(n)
		this_.priv_AstValidator_validateChildless_851(n)
	case int32(129) /* Token.COMPUTED_PROP */ :
		this_.priv_AstValidator_validateComputedPropClassMethod_873(n)
	case int32(108) /* Token.MEMBER_FIELD_DEF */ :
		this_.priv_AstValidator_validateClassField_792(n)
	case int32(109) /* Token.COMPUTED_FIELD_DEF */ :
		this_.priv_AstValidator_validateComputedPropClassField_793(n)
	case int32(173) /* Token.INDEX_SIGNATURE */ :
		this_.priv_AstValidator_validateProperties_897(n)
		this_.priv_AstValidator_validateChildCount_891(n)
		var child pkg_rhino.Node = n.GetFirstChild()
		_ = child
		this_.priv_AstValidator_validateProperties_897(child)
		this_.priv_AstValidator_validateChildless_851(child)
	case int32(174) /* Token.CALL_SIGNATURE */ :
		this_.priv_AstValidator_validateProperties_897(n)
		this_.priv_AstValidator_validateChildCount_891(n)
	case int32(90) /* Token.BLOCK */ :
		this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_CLASS_STATIC_BLOCK(), n)
		this_.priv_AstValidator_validateBlock_794(n)
	case int32(88) /* Token.EMPTY */ :
	default:
		this_.priv_AstValidator_violation_888(("Class contained member of invalid type " + j2gort.StrVal(n.GetToken())), n)
	}
}

// Java method: AstValidator.validateMemberFunction(com.google.javascript.rhino.Node n, boolean isAmbient)
func (this_ *AstValidator) priv_AstValidator_validateMemberFunction_791(n pkg_rhino.Node, isAmbient bool) {
	_ = n
	_ = isAmbient
	var function pkg_rhino.Node = n.GetFirstChild()
	_ = function
	if isAmbient {
		this_.priv_AstValidator_validateFunctionSignature_808(function)
	} else {
		this_.priv_AstValidator_validateFunctionExpression_807(function)
	}
}

// Java method: AstValidator.validateClassField(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateClassField_792(n pkg_rhino.Node) {
	_ = n
	if n.IsPrivateIdentifier() {
		this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_PRIVATE_ELEMENTS(), n)
	} else {
		this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_PUBLIC_CLASS_FIELDS(), n)
	}
	this_.priv_AstValidator_validateNonEmptyString_799(n)
	if n.HasChildren() {
		this_.ValidateExpression(n.GetFirstChild())
	}
}

// Java method: AstValidator.validateComputedPropClassField(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateComputedPropClassField_793(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_PUBLIC_CLASS_FIELDS(), n)
	this_.ValidateExpression(n.GetFirstChild())
	if n.GetSecondChild() != 0 {
		this_.ValidateExpression(n.GetSecondChild())
	}
}

// Java method: AstValidator.validateBlock(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateBlock_794(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_BLOCK(), n)
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.ValidateStatement_1(c)
		}
	}
}

// Java method: AstValidator.validateHasSourceName(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateHasSourceName_795(n pkg_rhino.Node) {
	_ = n
	var sourceName string = n.GetSourceFileName()
	_ = sourceName
	if j2gort.Ext_com_google_common_base_Strings_IsNullOrEmpty(sourceName) {
		this_.priv_AstValidator_violation_888("Missing 'source name' annotation.", n)
	}
}

// Java method: AstValidator.validateHasInputId(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateHasInputId_796(n pkg_rhino.Node) {
	_ = n
	var inputId *pkg_rhino.InputId = n.GetInputId()
	_ = inputId
	if inputId == nil {
		this_.priv_AstValidator_violation_888("Missing 'input id' annotation.", n)
	}
}

// Java method: AstValidator.validateLabel(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateLabel_797(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_LABEL(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	this_.priv_AstValidator_validateLabelName_798(n.GetFirstChild())
	this_.ValidateStatement_1(n.GetLastChild())
}

// Java method: AstValidator.validateLabelName(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateLabelName_798(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_LABEL_NAME(), n)
	this_.priv_AstValidator_validateNonEmptyString_799(n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
}

// Java method: AstValidator.validateNonEmptyString(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateNonEmptyString_799(n pkg_rhino.Node) {
	_ = n
	if this_.priv_AstValidator_validateNonNullString_801(n) && (j2gort.CheckNotNullString(n.GetString()) == "") {
		this_.priv_AstValidator_violation_888("Expected non-empty string.", n)
	}
}

// Java method: AstValidator.validateEmptyString(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateEmptyString_800(n pkg_rhino.Node) {
	_ = n
	if this_.priv_AstValidator_validateNonNullString_801(n) && !(j2gort.CheckNotNullString(n.GetString()) == "") {
		this_.priv_AstValidator_violation_888("Expected empty string.", n)
	}
}

// Java method: AstValidator.validateNonNullString(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateNonNullString_801(n pkg_rhino.Node) bool {
	_ = n
	var didRet_9 bool
	var retVal_10 bool
	func() {
		defer func() {
			if rec_11 := recover(); rec_11 != nil {
				jex_12 := javalang.RecoverException(rec_11)
				_ = jex_12
				if j2gort.InstanceOf[javalang.IRuntimeException](jex_12) {
					e := j2gort.AsIface[javalang.IRuntimeException](jex_12)
					_ = e
					this_.priv_AstValidator_violation_888("Expected non-null string.", n)
					retVal_10 = false
					didRet_9 = true
					return
				} else {
					panic(rec_11)
				}
			}
		}()
		if n.GetString() == j2gort.NullString {
			this_.priv_AstValidator_violation_888("Expected non-null string.", n)
			retVal_10 = false
			didRet_9 = true
			return
		}
	}()
	if didRet_9 {
		return retVal_10
	}
	return true
}

// Java method: AstValidator.validateName(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateName_802(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_NAME(), n)
	this_.priv_AstValidator_validateNonEmptyString_799(n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	this_.priv_AstValidator_validateTypeInformation_760(n)
	this_.priv_AstValidator_validateShadowContentIfPresent_804(n)
}

// Java method: AstValidator.validateOptionalName(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateOptionalName_803(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_NAME(), n)
	this_.priv_AstValidator_validateNonNullString_801(n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	var isEmpty bool = ((n.GetString() != j2gort.NullString) && (j2gort.CheckNotNullString(n.GetString()) == ""))
	_ = isEmpty
	if !(isEmpty) {
		this_.priv_AstValidator_validateTypeInformation_760(n)
	}
}

// Java method: AstValidator.validateShadowContentIfPresent(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateShadowContentIfPresent_804(n pkg_rhino.Node) {
	_ = n
	var shadow pkg_rhino.Node = n.GetClosureUnawareShadow()
	_ = shadow
	if shadow == 0 {
		return
	}
	if !(shadow.IsRoot()) {
		this_.priv_AstValidator_violation_888("Shadow reference node is not a ROOT node", shadow)
		return
	}
	var shadowScript pkg_rhino.Node = shadow.GetFirstChild()
	_ = shadowScript
	if (shadowScript == 0) || !(shadowScript.IsScript()) {
		this_.priv_AstValidator_violation_888("Shadow root node's child is not a script node", shadowScript)
		return
	}
	if shadowScript.GetChildCount() != int32(1) {
		this_.priv_AstValidator_violation_888("Shadow SCRIPT node child has more than one child", shadowScript)
		return
	}
	var exprResult pkg_rhino.Node = shadowScript.GetFirstChild()
	_ = exprResult
	if (exprResult == 0) || !(exprResult.IsExprResult()) {
		this_.priv_AstValidator_violation_888("Shadow SCRIPT node child is not an expr result node", exprResult)
		return
	}
	if exprResult.GetChildCount() != int32(1) {
		this_.priv_AstValidator_violation_888("Shadow EXPR_RESULT node should have exactly one child", exprResult)
		return
	}
	var shadowJsCall pkg_rhino.Node = exprResult.GetOnlyChild()
	_ = shadowJsCall
	if !(shadowJsCall.IsCall()) {
		this_.priv_AstValidator_violation_888("Shadow node EXPR_RESULT child is not a call", shadowJsCall)
		return
	}
	var shadowJsFunction pkg_rhino.Node = shadowJsCall.GetLastChild()
	_ = shadowJsFunction
	if !(shadowJsFunction.IsFunction()) {
		this_.priv_AstValidator_violation_888("Shadow node CALL child is not a function", shadowJsFunction)
		return
	}
	this_.priv_AstValidator_validateFunctionExpression_807(shadowJsFunction)
}

// Java method: AstValidator.validateEmptyName(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateEmptyName_805(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_NAME(), n)
	this_.priv_AstValidator_validateEmptyString_800(n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
}

// Java method: AstValidator.validateFunctionStatement(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateFunctionStatement_806(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_FUNCTION(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	this_.priv_AstValidator_validateName_802(n.GetFirstChild())
	this_.priv_AstValidator_validateParameters_812(n.GetSecondChild())
	this_.priv_AstValidator_validateFunctionBody_811(n.GetLastChild(), false)
	this_.priv_AstValidator_validateFunctionFeatures_810(n)
	if n.GetParent().IsBlock() && !(n.GetGrandparent().IsFunction()) {
		this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION(), n)
	}
}

// Java method: AstValidator.validateFunctionExpression(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateFunctionExpression_807(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFunctionExpressionHelper_809(n, false)
}

// Java method: AstValidator.validateFunctionSignature(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateFunctionSignature_808(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFunctionExpressionHelper_809(n, true)
}

// Java method: AstValidator.validateFunctionExpressionHelper(com.google.javascript.rhino.Node n, boolean isAmbient)
func (this_ *AstValidator) priv_AstValidator_validateFunctionExpressionHelper_809(n pkg_rhino.Node, isAmbient bool) {
	_ = n
	_ = isAmbient
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_FUNCTION(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	this_.priv_AstValidator_validateParameters_812(n.GetSecondChild())
	var name pkg_rhino.Node = n.GetFirstChild()
	_ = name
	var body pkg_rhino.Node = n.GetLastChild()
	_ = body
	if n.IsArrowFunction() {
		this_.priv_AstValidator_validateEmptyName_805(name)
		if body.IsBlock() {
			this_.priv_AstValidator_validateBlock_794(body)
		} else {
			this_.ValidateExpression(body)
		}
	} else {
		this_.priv_AstValidator_validateOptionalName_803(name)
		this_.priv_AstValidator_validateFunctionBody_811(body, isAmbient)
	}
	this_.priv_AstValidator_validateFunctionFeatures_810(n)
}

// Java method: AstValidator.validateFunctionFeatures(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateFunctionFeatures_810(n pkg_rhino.Node) {
	_ = n
	if n.IsArrowFunction() {
		this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_ARROW_FUNCTIONS(), n)
	}
	if n.IsGeneratorFunction() {
		this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_GENERATORS(), n)
	}
	if n.IsAsyncFunction() {
		this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_ASYNC_FUNCTIONS(), n)
	}
	if n.IsAsyncFunction() && n.IsGeneratorFunction() {
		this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_ASYNC_GENERATORS(), n)
	}
}

// Java method: AstValidator.validateFunctionBody(com.google.javascript.rhino.Node n, boolean noBlock)
func (this_ *AstValidator) priv_AstValidator_validateFunctionBody_811(n pkg_rhino.Node, noBlock bool) {
	_ = n
	_ = noBlock
	if noBlock {
		this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_EMPTY(), n)
	} else {
		this_.priv_AstValidator_validateBlock_794(n)
	}
}

// Java method: AstValidator.validateParameters(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateParameters_812(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_PARAM_LIST(), n)
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			if c.IsRest() {
				this_.priv_AstValidator_validateRestParameters_818(pkg_rhino.StaticGet_Token_PARAM_LIST(), c)
			} else if c.IsDefaultValue() {
				this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_DEFAULT_PARAMETERS(), c)
				this_.priv_AstValidator_validateDefaultValue_813(pkg_rhino.StaticGet_Token_PARAM_LIST(), c)
			} else {
				if c.IsName() {
					this_.priv_AstValidator_validateName_802(c)
				} else if c.IsArrayPattern() {
					this_.priv_AstValidator_validateArrayPattern_828(pkg_rhino.StaticGet_Token_PARAM_LIST(), c)
				} else {
					this_.priv_AstValidator_validateObjectPattern_829(pkg_rhino.StaticGet_Token_PARAM_LIST(), c)
				}
			}
		}
	}
}

// Java method: AstValidator.validateDefaultValue(com.google.javascript.rhino.Token contextType, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateDefaultValue_813(contextType pkg_rhino.Token, n pkg_rhino.Node) {
	_ = contextType
	_ = n
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	this_.priv_AstValidator_validateLHS_826(contextType, n.GetFirstChild())
	this_.ValidateExpression(n.GetLastChild())
}

// Java method: AstValidator.validateCall(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateCall_814(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_CALL(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateMinimumChildCount_894(n, int32(1))
	var callee pkg_rhino.Node = n.GetFirstChild()
	_ = callee
	if callee.IsSuper() {
		this_.priv_AstValidator_validateSuper_816(callee)
	} else {
		this_.ValidateExpression(callee)
	}
	{
		var c pkg_rhino.Node = callee.GetNext()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validatePseudoExpression_758(c, pkg_rhino.StaticGet_Token_ITER_SPREAD())
		}
	}
}

// Java method: AstValidator.validateOptChainCall(com.google.javascript.rhino.Node node)
func (this_ *AstValidator) priv_AstValidator_validateOptChainCall_815(node pkg_rhino.Node) {
	_ = node
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_OPTIONAL_CHAINING(), node)
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_OPTCHAIN_CALL(), node)
	this_.priv_AstValidator_validateProperties_897(node)
	this_.priv_AstValidator_validateMinimumChildCount_894(node, int32(1))
	var callee pkg_rhino.Node = node.GetFirstChild()
	_ = callee
	this_.ValidateExpression(callee)
	{
		var argument pkg_rhino.Node = callee.GetNext()
		_ = argument
		for ; argument != 0; argument = argument.GetNext() {
			this_.priv_AstValidator_validatePseudoExpression_758(argument, pkg_rhino.StaticGet_Token_ITER_SPREAD())
		}
	}
	this_.priv_AstValidator_validateFirstNodeOfOptChain_889(node)
}

// Java method: AstValidator.validateSuper(com.google.javascript.rhino.Node superNode)
func (this_ *AstValidator) priv_AstValidator_validateSuper_816(superNode pkg_rhino.Node) {
	_ = superNode
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_SUPER(), superNode)
	this_.priv_AstValidator_validateProperties_897(superNode)
	this_.priv_AstValidator_validateChildless_851(superNode)
	this_.priv_AstValidator_validateTypeInformation_760(superNode)
	var superParent pkg_rhino.Node = superNode.GetParent()
	_ = superParent
	var methodNode pkg_rhino.Node = NodeUtil_GetEnclosingNonArrowFunction(superParent)
	_ = methodNode
	if NodeUtil_IsNormalGet(superParent) && superNode.IsFirstChildOf(superParent) {
		if !(this_.priv_AstValidator_allowsSuperPropertyReference_817(superParent)) {
			this_.priv_AstValidator_violation_888("super property references are only allowed in methods, class static blocks and class fields.", superNode)
		}
	} else if superParent.IsCall() && superNode.IsFirstChildOf(superParent) {
		if (methodNode == 0) || !(NodeUtil_IsEs6Constructor(methodNode)) {
			this_.priv_AstValidator_violation_888("super constructor call is only allowed in a constructor method", superNode)
		} else {
			var extendsNode pkg_rhino.Node = methodNode.GetParent().GetParent().GetParent().GetSecondChild()
			_ = extendsNode
			if extendsNode.IsEmpty() {
				this_.priv_AstValidator_violation_888("super constructor call in a class that extends nothing", superNode)
			}
		}
	} else {
		this_.priv_AstValidator_violation_888("`super` is a syntax error here", superNode)
	}
}

// Java method: AstValidator.allowsSuperPropertyReference(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_allowsSuperPropertyReference_817(n pkg_rhino.Node) bool {
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
		return this_.priv_AstValidator_allowsSuperPropertyReference_817(n.GetParent())
	case int32(90) /* Token.BLOCK */ :
		if NodeUtil_IsClassStaticBlock(n) {
			return true
		}
		return this_.priv_AstValidator_allowsSuperPropertyReference_817(n.GetParent())
	default:
		return this_.priv_AstValidator_allowsSuperPropertyReference_817(n.GetParent())
	}
	panic("unreachable switch expression")
}

// Java method: AstValidator.validateRestParameters(com.google.javascript.rhino.Token contextType, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateRestParameters_818(contextType pkg_rhino.Token, n pkg_rhino.Node) {
	_ = contextType
	_ = n
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_REST_PARAMETERS(), n)
	this_.priv_AstValidator_validateRest_821(contextType, n)
}

// Java method: AstValidator.validateArrayPatternRest(com.google.javascript.rhino.Token contextType, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateArrayPatternRest_819(contextType pkg_rhino.Token, n pkg_rhino.Node) {
	_ = contextType
	_ = n
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_ARRAY_PATTERN_REST(), n)
	this_.priv_AstValidator_validateRest_821(contextType, n)
}

// Java method: AstValidator.validateObjectPatternRest(com.google.javascript.rhino.Token contextType, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateObjectPatternRest_820(contextType pkg_rhino.Token, n pkg_rhino.Node) {
	_ = contextType
	_ = n
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_OBJECT_PATTERN_REST(), n)
	this_.priv_AstValidator_validateRest_821(contextType, n)
}

// Java method: AstValidator.validateRest(com.google.javascript.rhino.Token contextType, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateRest_821(contextType pkg_rhino.Token, n pkg_rhino.Node) {
	_ = contextType
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(125) /* Token.ITER_REST */, int32(126) /* Token.OBJECT_REST */ :
	default:
		this_.priv_AstValidator_violation_888("Unexpected node type.", n)
		return
	}
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	this_.priv_AstValidator_validateLHS_826(contextType, n.GetFirstChild())
	if n.GetNext() != 0 {
		this_.priv_AstValidator_violation_888("Rest parameters must come after all other parameters.", n)
	}
}

// Java method: AstValidator.validateObjectSpread(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateObjectSpread_822(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD(), n)
	this_.ValidateExpression(n.GetFirstChild())
}

// Java method: AstValidator.validateNew(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateNew_823(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_NEW(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateMinimumChildCount_894(n, int32(1))
	this_.ValidateExpression(n.GetFirstChild())
	{
		var c pkg_rhino.Node = n.GetSecondChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validatePseudoExpression_758(c, pkg_rhino.StaticGet_Token_ITER_SPREAD())
		}
	}
}

// Java method: AstValidator.validateNameDeclarationHelper(com.google.javascript.rhino.Node statement, com.google.javascript.rhino.Token declaration, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateNameDeclarationHelper_824(statement pkg_rhino.Node, declaration pkg_rhino.Token, n pkg_rhino.Node) {
	_ = statement
	_ = declaration
	_ = n
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateMinimumChildCount_894(n, int32(1))
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validateNameDeclarationChild_825(statement, declaration, c)
		}
	}
	{
		var iifeCond_2 bool
		{

			re_13 := declaration
			ae_14 := pkg_rhino.StaticGet_Token_LET()
			if re_13 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_2 = re_13 == ae_14
		}
		if iifeCond_2 {
			this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_LET_DECLARATIONS(), n)
		} else {
			var iifeCond_1 bool
			{

				re_15 := declaration
				ae_16 := pkg_rhino.StaticGet_Token_CONST()
				if re_15 == 0 {
					j2gort.ThrowNPE()
				}
				iifeCond_1 = re_15 == ae_16
			}
			if iifeCond_1 {
				this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_CONST_DECLARATIONS(), n)
			}
		}
	}
}

// Java method: AstValidator.validateNameDeclarationChild(com.google.javascript.rhino.Node statement, com.google.javascript.rhino.Token declaration, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateNameDeclarationChild_825(statement pkg_rhino.Node, declaration pkg_rhino.Token, n pkg_rhino.Node) {
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
		re_17 := declaration
		ae_18 := pkg_rhino.StaticGet_Token_VAR()
		if re_17 == 0 {
			j2gort.ThrowNPE()
		}
		return re_17 == ae_18
	}() {
		minValues = int32(0)
		maxValues = int32(1)
	} else if inEnhancedFor {
		minValues = int32(0)
		maxValues = int32(0)
	} else if n.IsDestructuringLhs() || func() bool {
		re_19 := declaration
		ae_20 := pkg_rhino.StaticGet_Token_CONST()
		if re_19 == 0 {
			j2gort.ThrowNPE()
		}
		return re_19 == ae_20
	}() {
		minValues = int32(1)
		maxValues = int32(1)
	} else {
		minValues = int32(0)
		maxValues = int32(1)
	}
	if n.IsName() {
		this_.priv_AstValidator_validateNonEmptyString_799(n)
		this_.priv_AstValidator_validateProperties_897(n)
		this_.priv_AstValidator_validateChildCountIn_893(n, minValues, maxValues)
		if n.HasChildren() {
			this_.ValidateExpression(n.GetFirstChild())
		}
	} else if n.IsDestructuringLhs() {
		this_.priv_AstValidator_validateProperties_897(n)
		this_.priv_AstValidator_validateChildCountIn_893(n, (int32(1) + minValues), (int32(1) + maxValues))
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		switch c.TokenSwitchOrd_() {
		case int32(102) /* Token.ARRAY_PATTERN */ :
			this_.priv_AstValidator_validateArrayPattern_828(declaration, c)
		case int32(103) /* Token.OBJECT_PATTERN */ :
			this_.priv_AstValidator_validateObjectPattern_829(declaration, c)
		default:
			this_.priv_AstValidator_violation_888(("Invalid destructuring lhs first child for " + j2gort.StrVal(declaration) + " node"), n)
		}
		if n.HasTwoChildren() {
			this_.ValidateExpression(n.GetSecondChild())
		}
	} else {
		this_.priv_AstValidator_violation_888(("Invalid child for " + j2gort.StrVal(declaration) + " node"), n)
	}
}

// Java method: AstValidator.validateLHS(com.google.javascript.rhino.Token contextType, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateLHS_826(contextType pkg_rhino.Token, n pkg_rhino.Node) {
	_ = contextType
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		this_.priv_AstValidator_validateName_802(n)
	case int32(102) /* Token.ARRAY_PATTERN */ :
		this_.priv_AstValidator_validateArrayPattern_828(contextType, n)
	case int32(103) /* Token.OBJECT_PATTERN */ :
		this_.priv_AstValidator_validateObjectPattern_829(contextType, n)
	case int32(26) /* Token.GETPROP */, int32(27) /* Token.GETELEM */ :
		this_.priv_AstValidator_validateGetPropGetElemInLHS_827(contextType, n)
	case int32(101) /* Token.CAST */ :
		this_.priv_AstValidator_validateLHS_826(contextType, n.GetOnlyChild())
	default:
		this_.priv_AstValidator_violation_888(("Invalid child for " + j2gort.StrVal(contextType) + " node"), n)
	}
}

// Java method: AstValidator.validateGetPropGetElemInLHS(com.google.javascript.rhino.Token contextType, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateGetPropGetElemInLHS_827(contextType pkg_rhino.Token, n pkg_rhino.Node) {
	_ = contextType
	_ = n
	if (((contextType == pkg_rhino.StaticGet_Token_CONST()) || (contextType == pkg_rhino.StaticGet_Token_LET())) || (contextType == pkg_rhino.StaticGet_Token_VAR())) || (contextType == pkg_rhino.StaticGet_Token_PARAM_LIST()) {
		this_.priv_AstValidator_violation_888(("Invalid child for " + j2gort.StrVal(contextType) + " node"), n)
		return
	}
	switch n.TokenSwitchOrd_() {
	case int32(26) /* Token.GETPROP */ :
		this_.priv_AstValidator_validateGetProp_857(n)
	case int32(27) /* Token.GETELEM */ :
		this_.priv_AstValidator_validateGetElem_855(n)
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Expected GETPROP or GETELEM but instead got node "+j2gort.StrVal(n.GetToken()))))
	}
}

// Java method: AstValidator.validateArrayPattern(com.google.javascript.rhino.Token type, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateArrayPattern_828(type_ pkg_rhino.Token, n pkg_rhino.Node) {
	_ = type_
	_ = n
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_ARRAY_DESTRUCTURING(), n)
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_ARRAY_PATTERN(), n)
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			switch c.TokenSwitchOrd_() {
			case int32(134) /* Token.DEFAULT_VALUE */ :
				this_.priv_AstValidator_validateDefaultValue_813(type_, c)
			case int32(125) /* Token.ITER_REST */ :
				this_.priv_AstValidator_validateArrayPatternRest_819(type_, c)
			case int32(88) /* Token.EMPTY */ :
				this_.priv_AstValidator_validateProperties_897(c)
				this_.priv_AstValidator_validateChildless_851(c)
			default:
				this_.priv_AstValidator_validateLHS_826(type_, c)
			}
		}
	}
}

// Java method: AstValidator.validateObjectPattern(com.google.javascript.rhino.Token type, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateObjectPattern_829(type_ pkg_rhino.Token, n pkg_rhino.Node) {
	_ = type_
	_ = n
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_OBJECT_DESTRUCTURING(), n)
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_OBJECT_PATTERN(), n)
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			switch c.TokenSwitchOrd_() {
			case int32(100) /* Token.STRING_KEY */ :
				this_.priv_AstValidator_validateObjectPatternStringKey_870(type_, c)
			case int32(126) /* Token.OBJECT_REST */ :
				this_.priv_AstValidator_validateObjectPatternRest_820(type_, c)
			case int32(129) /* Token.COMPUTED_PROP */ :
				this_.priv_AstValidator_validateObjectPatternComputedPropKey_872(type_, c)
			default:
				this_.priv_AstValidator_violation_888(("Invalid object pattern child for " + j2gort.StrVal(type_) + " node"), n)
			}
		}
	}
}

// Java method: AstValidator.validateFor(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateFor_830(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_FOR(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_892(n, int32(4))
	var target pkg_rhino.Node = n.GetFirstChild()
	_ = target
	if NodeUtil_IsNameDeclaration(target) {
		this_.priv_AstValidator_validateNameDeclarationHelper_824(n, target.GetToken(), target)
	} else {
		this_.priv_AstValidator_validatePseudoExpression_758(target, pkg_rhino.StaticGet_Token_EMPTY())
	}
	this_.priv_AstValidator_validatePseudoExpression_758(n.GetSecondChild(), pkg_rhino.StaticGet_Token_EMPTY())
	this_.priv_AstValidator_validatePseudoExpression_758(n.GetChildAtIndex(int32(2)), pkg_rhino.StaticGet_Token_EMPTY())
	this_.priv_AstValidator_validateBlock_794(n.GetLastChild())
}

// Java method: AstValidator.validateForIn(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateForIn_831(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_FOR_IN(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	this_.priv_AstValidator_validateEnhancedForVarOrAssignmentTarget_834(n, n.GetFirstChild())
	this_.ValidateExpression(n.GetSecondChild())
	this_.priv_AstValidator_validateBlock_794(n.GetLastChild())
}

// Java method: AstValidator.validateForOf(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateForOf_832(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_FOR_OF(), n)
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_FOR_OF(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	this_.priv_AstValidator_validateEnhancedForVarOrAssignmentTarget_834(n, n.GetFirstChild())
	this_.ValidateExpression(n.GetSecondChild())
	this_.priv_AstValidator_validateBlock_794(n.GetLastChild())
}

// Java method: AstValidator.validateForAwaitOf(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateForAwaitOf_833(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_FOR_AWAIT_OF(), n)
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_FOR_AWAIT_OF(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	this_.priv_AstValidator_validateEnhancedForVarOrAssignmentTarget_834(n, n.GetFirstChild())
	this_.ValidateExpression(n.GetSecondChild())
	this_.priv_AstValidator_validateBlock_794(n.GetLastChild())
	this_.priv_AstValidator_validateAwaitWithinAsyncFunction_767(n)
}

// Java method: AstValidator.validateEnhancedForVarOrAssignmentTarget(com.google.javascript.rhino.Node forNode, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateEnhancedForVarOrAssignmentTarget_834(forNode pkg_rhino.Node, n pkg_rhino.Node) {
	_ = forNode
	_ = n
	if NodeUtil_IsNameDeclaration(n) {
		this_.priv_AstValidator_validateProperties_897(n)
		this_.priv_AstValidator_validateChildCount_892(n, int32(1))
		this_.priv_AstValidator_validateNameDeclarationHelper_824(forNode, n.GetToken(), n)
	} else {
		this_.priv_AstValidator_validateLHS_826(n.GetParent().GetToken(), n)
	}
}

// Java method: AstValidator.validateWith(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateWith_835(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_WITH(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	this_.ValidateExpression(n.GetFirstChild())
	this_.priv_AstValidator_validateBlock_794(n.GetLastChild())
}

// Java method: AstValidator.validateWhile(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateWhile_836(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_WHILE(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	this_.ValidateExpression(n.GetFirstChild())
	this_.priv_AstValidator_validateBlock_794(n.GetLastChild())
}

// Java method: AstValidator.validateDo(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateDo_837(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_DO(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	this_.priv_AstValidator_validateBlock_794(n.GetFirstChild())
	this_.ValidateExpression(n.GetLastChild())
}

// Java method: AstValidator.validateIf(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateIf_838(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_IF(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCountIn_893(n, int32(2), int32(3))
	this_.ValidateExpression(n.GetFirstChild())
	this_.priv_AstValidator_validateBlock_794(n.GetSecondChild())
	if n.HasXChildren(int32(3)) {
		this_.priv_AstValidator_validateBlock_794(n.GetLastChild())
	}
}

// Java method: AstValidator.validateExprStmt(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateExprStmt_839(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_EXPR_RESULT(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	this_.ValidateExpression(n.GetFirstChild())
}

// Java method: AstValidator.validateReturn(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateReturn_840(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_RETURN(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateMaximumChildCount_895(n, int32(1))
	if n.HasChildren() {
		this_.ValidateExpression(n.GetFirstChild())
	}
}

// Java method: AstValidator.validateThrow(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateThrow_841(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_THROW(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	this_.ValidateExpression(n.GetFirstChild())
}

// Java method: AstValidator.validateBreak(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateBreak_842(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_BREAK(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateMaximumChildCount_895(n, int32(1))
	if n.HasChildren() {
		this_.priv_AstValidator_validateLabelName_798(n.GetFirstChild())
	}
}

// Java method: AstValidator.validateContinue(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateContinue_843(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_CONTINUE(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateMaximumChildCount_895(n, int32(1))
	if n.HasChildren() {
		this_.priv_AstValidator_validateLabelName_798(n.GetFirstChild())
	}
}

// Java method: AstValidator.validateTry(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateTry_844(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_TRY(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCountIn_893(n, int32(2), int32(3))
	this_.priv_AstValidator_validateBlock_794(n.GetFirstChild())
	var seenCatchOrFinally bool = false
	_ = seenCatchOrFinally
	var catches pkg_rhino.Node = n.GetSecondChild()
	_ = catches
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_BLOCK(), catches)
	this_.priv_AstValidator_validateProperties_897(catches)
	this_.priv_AstValidator_validateMaximumChildCount_895(catches, int32(1))
	if catches.HasChildren() {
		this_.priv_AstValidator_validateCatch_845(catches.GetFirstChild())
		seenCatchOrFinally = true
	}
	if n.HasXChildren(int32(3)) {
		this_.priv_AstValidator_validateBlock_794(n.GetLastChild())
		seenCatchOrFinally = true
	}
	if !(seenCatchOrFinally) {
		this_.priv_AstValidator_violation_888("Missing catch or finally for try statement.", n)
	}
}

// Java method: AstValidator.validateCatch(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateCatch_845(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_CATCH(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	var caught pkg_rhino.Node = n.GetFirstChild()
	_ = caught
	if caught.IsName() {
		this_.priv_AstValidator_validateName_802(caught)
	} else if caught.IsArrayPattern() {
		this_.priv_AstValidator_validateArrayPattern_828(pkg_rhino.StaticGet_Token_CATCH(), caught)
	} else if caught.IsObjectPattern() {
		this_.priv_AstValidator_validateObjectPattern_829(pkg_rhino.StaticGet_Token_CATCH(), caught)
	} else if caught.IsEmpty() {
		this_.priv_AstValidator_validateNoCatchBinding_846(caught)
	} else {
		this_.priv_AstValidator_violation_888(("Unexpected catch binding: " + j2gort.StrVal(caught)), n)
	}
	this_.priv_AstValidator_validateBlock_794(n.GetLastChild())
}

// Java method: AstValidator.validateNoCatchBinding(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateNoCatchBinding_846(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_OPTIONAL_CATCH_BINDING(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
}

// Java method: AstValidator.validateSwitch(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateSwitch_847(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_SWITCH(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_892(n, int32(2))
	this_.ValidateExpression(n.GetFirstChild())
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_SWITCH_BODY(), n.GetSecondChild())
	var cases pkg_rhino.Node = n.GetSecondChild()
	_ = cases
	var defaults int32 = int32(0)
	_ = defaults
	{
		var c pkg_rhino.Node = cases.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validateSwitchMember_848(c)
			if c.IsDefaultCase() {
				defaults = defaults + 1
			}
		}
	}
	if defaults > int32(1) {
		this_.priv_AstValidator_violation_888(("Expected at most 1 'default' in switch but was " + j2gort.StrValInt32(defaults)), n)
	}
}

// Java method: AstValidator.validateSwitchMember(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateSwitchMember_848(n pkg_rhino.Node) {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(76) /* Token.CASE */ :
		this_.priv_AstValidator_validateCase_850(n)
	case int32(77) /* Token.DEFAULT_CASE */ :
		this_.priv_AstValidator_validateDefaultCase_849(n)
	default:
		this_.priv_AstValidator_violation_888(("Expected switch member but was " + j2gort.StrVal(n.GetToken())), n)
	}
}

// Java method: AstValidator.validateDefaultCase(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateDefaultCase_849(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_DEFAULT_CASE(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	this_.priv_AstValidator_validateBlock_794(n.GetLastChild())
}

// Java method: AstValidator.validateCase(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateCase_850(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_CASE(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	this_.ValidateExpression(n.GetFirstChild())
	this_.priv_AstValidator_validateBlock_794(n.GetLastChild())
}

// Java method: AstValidator.validateChildless(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateChildless_851(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateChildCount_892(n, int32(0))
}

// Java method: AstValidator.validateAssignmentExpression(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateAssignmentExpression_852(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	this_.priv_AstValidator_validateLHS_826(n.GetToken(), n.GetFirstChild())
	this_.ValidateExpression(n.GetLastChild())
}

// Java method: AstValidator.validateCompoundAssignmentExpression(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateCompoundAssignmentExpression_853(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	var contextType pkg_rhino.Token = n.GetToken()
	_ = contextType
	var lhs pkg_rhino.Node = n.GetFirstChild()
	_ = lhs
	this_.priv_AstValidator_validateAssignmentOpTarget_854(lhs, contextType)
	this_.ValidateExpression(n.GetLastChild())
}

// Java method: AstValidator.validateAssignmentOpTarget(com.google.javascript.rhino.Node lhs, com.google.javascript.rhino.Token contextType)
func (this_ *AstValidator) priv_AstValidator_validateAssignmentOpTarget_854(lhs pkg_rhino.Node, contextType pkg_rhino.Token) {
	_ = lhs
	_ = contextType
	switch lhs.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		this_.priv_AstValidator_validateName_802(lhs)
	case int32(26) /* Token.GETPROP */, int32(27) /* Token.GETELEM */ :
		this_.priv_AstValidator_validateGetPropGetElemInLHS_827(contextType, lhs)
	case int32(101) /* Token.CAST */ :
		this_.priv_AstValidator_validateProperties_897(lhs)
		this_.priv_AstValidator_validateChildCount_892(lhs, int32(1))
		this_.priv_AstValidator_validateAssignmentOpTarget_854(lhs.GetFirstChild(), contextType)
	default:
		this_.priv_AstValidator_violation_888(("Invalid child for " + j2gort.StrVal(contextType) + " node"), lhs)
	}
}

// Java method: AstValidator.validateGetElem(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateGetElem_855(n pkg_rhino.Node) {
	_ = n
	if !(n.IsGetElem()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(n))
	}
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_892(n, int32(2))
	this_.priv_AstValidator_validatePropertyReferenceTarget_859(n.GetFirstChild())
	this_.ValidateExpression(n.GetLastChild())
}

// Java method: AstValidator.validateOptChainGetElem(com.google.javascript.rhino.Node node)
func (this_ *AstValidator) priv_AstValidator_validateOptChainGetElem_856(node pkg_rhino.Node) {
	_ = node
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_OPTIONAL_CHAINING(), node)
	if !(node.IsOptChainGetElem()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(node))
	}
	this_.priv_AstValidator_validateProperties_897(node)
	this_.priv_AstValidator_validateChildCount_892(node, int32(2))
	this_.ValidateExpression(node.GetFirstChild())
	this_.ValidateExpression(node.GetLastChild())
	this_.priv_AstValidator_validateFirstNodeOfOptChain_889(node)
}

// Java method: AstValidator.validateGetProp(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateGetProp_857(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_GETPROP(), n)
	this_.priv_AstValidator_validatePropertyReferenceTarget_859(n.GetFirstChild())
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	this_.priv_AstValidator_validateNonEmptyString_799(n)
}

// Java method: AstValidator.validateOptChainGetProp(com.google.javascript.rhino.Node node)
func (this_ *AstValidator) priv_AstValidator_validateOptChainGetProp_858(node pkg_rhino.Node) {
	_ = node
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_OPTIONAL_CHAINING(), node)
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_OPTCHAIN_GETPROP(), node)
	this_.ValidateExpression(node.GetFirstChild())
	this_.priv_AstValidator_validateFirstNodeOfOptChain_889(node)
	this_.priv_AstValidator_validateProperties_897(node)
	this_.priv_AstValidator_validateChildCount_891(node)
	this_.priv_AstValidator_validateNonEmptyString_799(node)
}

// Java method: AstValidator.validatePropertyReferenceTarget(com.google.javascript.rhino.Node objectNode)
func (this_ *AstValidator) priv_AstValidator_validatePropertyReferenceTarget_859(objectNode pkg_rhino.Node) {
	_ = objectNode
	if objectNode.IsSuper() {
		this_.priv_AstValidator_validateSuper_816(objectNode)
	} else {
		this_.ValidateExpression(objectNode)
	}
}

// Java method: AstValidator.validateRegExpLit(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateRegExpLit_860(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_SYNTAX(), n)
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_REGEXP(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCountIn_893(n, int32(1), int32(2))
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validateStringLit_861(c)
		}
	}
}

// Java method: AstValidator.validateStringLit(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateStringLit_861(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_STRINGLIT(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	func() {
		defer func() {
			if rec_21 := recover(); rec_21 != nil {
				jex_22 := javalang.RecoverException(rec_21)
				_ = jex_22
				if j2gort.InstanceOf[javalang.IUnsupportedOperationException](jex_22) {
					e := j2gort.AsIface[javalang.IUnsupportedOperationException](jex_22)
					_ = e
					this_.priv_AstValidator_violation_888("Invalid STRING node.", n)
				} else {
					panic(rec_21)
				}
			}
		}()
		n.GetString()
	}()
}

// Java method: AstValidator.validateNumber(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateNumber_862(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_NUMBER(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	func() {
		defer func() {
			if rec_23 := recover(); rec_23 != nil {
				jex_24 := javalang.RecoverException(rec_23)
				_ = jex_24
				if j2gort.InstanceOf[javalang.IUnsupportedOperationException](jex_24) {
					e := j2gort.AsIface[javalang.IUnsupportedOperationException](jex_24)
					_ = e
					this_.priv_AstValidator_violation_888("Invalid NUMBER node.", n)
				} else {
					panic(rec_23)
				}
			}
		}()
		n.GetDouble()
	}()
}

// Java method: AstValidator.validateBigInt(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateBigInt_863(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_BIGINT(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	func() {
		defer func() {
			if rec_25 := recover(); rec_25 != nil {
				jex_26 := javalang.RecoverException(rec_25)
				_ = jex_26
				if j2gort.InstanceOf[javalang.IUnsupportedOperationException](jex_26) {
					e := j2gort.AsIface[javalang.IUnsupportedOperationException](jex_26)
					_ = e
					this_.priv_AstValidator_violation_888("Invalid BIGINT node.", n)
				} else {
					panic(rec_25)
				}
			}
		}()
		n.GetBigInt()
	}()
}

// Java method: AstValidator.validateArrayLit(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateArrayLit_864(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_ARRAYLIT(), n)
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validatePseudoExpression_759(c, pkg_rhino.StaticGet_Token_EMPTY(), pkg_rhino.StaticGet_Token_ITER_SPREAD())
			break
		}
	}
}

// Java method: AstValidator.validateObjectLit(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateObjectLit_865(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_OBJECTLIT(), n)
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_AstValidator_validateObjectLitKey_866(c)
		}
	}
}

// Java method: AstValidator.validateObjectLitKey(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateObjectLitKey_866(n pkg_rhino.Node) {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(95) /* Token.GETTER_DEF */ :
		this_.priv_AstValidator_validateObjectLitGetKey_867(n)
	case int32(96) /* Token.SETTER_DEF */ :
		this_.priv_AstValidator_validateObjectLitSetKey_868(n)
	case int32(100) /* Token.STRING_KEY */ :
		this_.priv_AstValidator_validateObjectLitStringKey_869(n)
	case int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
		this_.priv_AstValidator_validateClassMember_790(n, false)
		if n.IsStaticMember() {
			this_.priv_AstValidator_violation_888("Keys in an object literal should not be static.", n)
		}
	case int32(129) /* Token.COMPUTED_PROP */ :
		this_.priv_AstValidator_validateObjectLitComputedPropKey_871(n)
	case int32(128) /* Token.OBJECT_SPREAD */ :
		this_.priv_AstValidator_validateObjectSpread_822(n)
	default:
		this_.priv_AstValidator_violation_888(("Expected object literal key expression but was " + j2gort.StrVal(n.GetToken())), n)
	}
}

// Java method: AstValidator.validateObjectLitGetKey(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateObjectLitGetKey_867(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_GETTER(), n)
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_GETTER_DEF(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	this_.priv_AstValidator_validateObjectLiteralKeyName_876(n)
	var function pkg_rhino.Node = n.GetFirstChild()
	_ = function
	this_.priv_AstValidator_validateFunctionExpression_807(function)
	if !(j2gort.CheckNotNullString(function.GetFirstChild().GetString()) == "") {
		this_.priv_AstValidator_violation_888("Expected unnamed function expression.", n)
	}
	var functionParams pkg_rhino.Node = function.GetSecondChild()
	_ = functionParams
	if functionParams.HasChildren() {
		this_.priv_AstValidator_violation_888("get methods must not have parameters.", n)
	}
}

// Java method: AstValidator.validateObjectLitSetKey(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateObjectLitSetKey_868(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_SETTER(), n)
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_SETTER_DEF(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	this_.priv_AstValidator_validateObjectLiteralKeyName_876(n)
	var function pkg_rhino.Node = n.GetFirstChild()
	_ = function
	this_.priv_AstValidator_validateFunctionExpression_807(function)
	if !(j2gort.CheckNotNullString(function.GetFirstChild().GetString()) == "") {
		this_.priv_AstValidator_violation_888("Expected unnamed function expression.", n)
	}
	var functionParams pkg_rhino.Node = function.GetSecondChild()
	_ = functionParams
	if !(functionParams.HasOneChild()) {
		this_.priv_AstValidator_violation_888("set methods must have exactly one parameter.", n)
	}
}

// Java method: AstValidator.validateObjectLitStringKey(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateObjectLitStringKey_869(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_STRING_KEY(), n)
	this_.priv_AstValidator_validateObjectLiteralKeyName_876(n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_892(n, int32(1))
	this_.ValidateExpression(n.GetFirstChild())
	if n.GetBooleanProp(pkg_rhino.StaticGet_Node_IS_SHORTHAND_PROPERTY()) {
		this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES(), n)
	}
}

// Java method: AstValidator.validateObjectPatternStringKey(com.google.javascript.rhino.Token type, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateObjectPatternStringKey_870(type_ pkg_rhino.Token, n pkg_rhino.Node) {
	_ = type_
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_STRING_KEY(), n)
	this_.priv_AstValidator_validateObjectLiteralKeyName_876(n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_892(n, int32(1))
	var c pkg_rhino.Node = n.GetFirstChild()
	_ = c
	switch c.TokenSwitchOrd_() {
	case int32(134) /* Token.DEFAULT_VALUE */ :
		this_.priv_AstValidator_validateDefaultValue_813(type_, c)
	default:
		this_.priv_AstValidator_validateLHS_826(type_, c)
	}
}

// Java method: AstValidator.validateObjectLitComputedPropKey(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateObjectLitComputedPropKey_871(n pkg_rhino.Node) {
	_ = n
	if !(n.GetIsInClosureUnawareSubtree() && NodeUtil_IsGetOrSetKey(n)) {
		this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_COMPUTED_PROPERTIES(), n)
	}
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_COMPUTED_PROP(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	this_.ValidateExpression(n.GetFirstChild())
	this_.ValidateExpression(n.GetLastChild())
}

// Java method: AstValidator.validateObjectPatternComputedPropKey(com.google.javascript.rhino.Token type, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateObjectPatternComputedPropKey_872(type_ pkg_rhino.Token, n pkg_rhino.Node) {
	_ = type_
	_ = n
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_COMPUTED_PROPERTIES(), n)
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_COMPUTED_PROP(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	this_.ValidateExpression(n.GetFirstChild())
	if n.GetLastChild().IsDefaultValue() {
		this_.priv_AstValidator_validateDefaultValue_813(type_, n.GetLastChild())
	} else {
		this_.priv_AstValidator_validateLHS_826(n.GetLastChild().GetToken(), n.GetLastChild())
	}
}

// Java method: AstValidator.validateComputedPropClassMethod(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateComputedPropClassMethod_873(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_COMPUTED_PROPERTIES(), n)
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_COMPUTED_PROP(), n)
	this_.ValidateExpression(n.GetFirstChild())
	if n.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_VARIABLE()) {
		this_.priv_AstValidator_validateProperties_897(n)
		this_.priv_AstValidator_validateChildCount_892(n, int32(1))
	} else {
		this_.priv_AstValidator_validateProperties_897(n)
		this_.priv_AstValidator_validateChildCount_892(n, int32(2))
		this_.priv_AstValidator_validateFunctionExpression_807(n.GetLastChild())
		if n.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_GETTER()) {
			this_.priv_AstValidator_validateObjectLitComputedPropGetKey_874(n)
			this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_CLASS_GETTER_SETTER(), n)
		} else if n.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_SETTER()) {
			this_.priv_AstValidator_validateObjectLitComputedPropSetKey_875(n)
			this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_CLASS_GETTER_SETTER(), n)
		}
	}
}

// Java method: AstValidator.validateObjectLitComputedPropGetKey(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateObjectLitComputedPropGetKey_874(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_COMPUTED_PROPERTIES(), n)
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_GETTER(), n)
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_COMPUTED_PROP(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	var function pkg_rhino.Node = n.GetLastChild()
	_ = function
	this_.priv_AstValidator_validateFunctionExpression_807(function)
	if !(j2gort.CheckNotNullString(function.GetFirstChild().GetString()) == "") {
		this_.priv_AstValidator_violation_888("Expected unnamed function expression.", n)
	}
	var functionParams pkg_rhino.Node = function.GetSecondChild()
	_ = functionParams
	if functionParams.HasChildren() {
		this_.priv_AstValidator_violation_888("get methods must not have parameters.", n)
	}
}

// Java method: AstValidator.validateObjectLitComputedPropSetKey(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateObjectLitComputedPropSetKey_875(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_COMPUTED_PROPERTIES(), n)
	this_.priv_AstValidator_validateFeature_896(pkg_parser.StaticGet_FeatureSet_Feature_SETTER(), n)
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_COMPUTED_PROP(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	var function pkg_rhino.Node = n.GetLastChild()
	_ = function
	this_.priv_AstValidator_validateFunctionExpression_807(function)
	if !(j2gort.CheckNotNullString(function.GetFirstChild().GetString()) == "") {
		this_.priv_AstValidator_violation_888("Expected unnamed function expression.", n)
	}
	var functionParams pkg_rhino.Node = function.GetSecondChild()
	_ = functionParams
	if !(functionParams.HasOneChild()) {
		this_.priv_AstValidator_violation_888("set methods must have exactly one parameter.", n)
	}
}

// Java method: AstValidator.validateObjectLiteralKeyName(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateObjectLiteralKeyName_876(n pkg_rhino.Node) {
	_ = n
	if n.IsQuotedStringKey() {
		func() {
			defer func() {
				if rec_27 := recover(); rec_27 != nil {
					jex_28 := javalang.RecoverException(rec_27)
					_ = jex_28
					if j2gort.InstanceOf[javalang.IUnsupportedOperationException](jex_28) {
						e := j2gort.AsIface[javalang.IUnsupportedOperationException](jex_28)
						_ = e
						this_.priv_AstValidator_violation_888(("getString failed for" + j2gort.StrVal(n.GetToken())), n)
					} else {
						panic(rec_27)
					}
				}
			}()
			n.GetString()
		}()
	} else {
		this_.priv_AstValidator_validateNonEmptyString_799(n)
	}
}

// Java method: AstValidator.validateIncDecOp(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateIncDecOp_877(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_892(n, int32(1))
	this_.priv_AstValidator_validateAssignmentOpTarget_854(n.GetFirstChild(), n.GetToken())
}

// Java method: AstValidator.validateUnaryOp(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateUnaryOp_878(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_892(n, int32(1))
	this_.ValidateExpression(n.GetFirstChild())
}

// Java method: AstValidator.validateBinaryOp(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateBinaryOp_879(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_892(n, int32(2))
	this_.ValidateExpression(n.GetFirstChild())
	this_.ValidateExpression(n.GetLastChild())
}

// Java method: AstValidator.validateTrinaryOp(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateTrinaryOp_880(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_892(n, int32(3))
	var first pkg_rhino.Node = n.GetFirstChild()
	_ = first
	this_.ValidateExpression(first)
	this_.ValidateExpression(first.GetNext())
	this_.ValidateExpression(n.GetLastChild())
}

// Java method: AstValidator.validateNamedType(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateNamedType_881(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_NAMED_TYPE(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	this_.priv_AstValidator_validateName_802(n.GetFirstChild())
}

// Java method: AstValidator.validateTypeAlias(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateTypeAlias_882(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_TYPE_ALIAS(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
}

// Java method: AstValidator.validateAmbientDeclaration(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateAmbientDeclaration_883(n pkg_rhino.Node) {
	_ = n
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_DECLARE(), n)
	this_.priv_AstValidator_validateAmbientDeclarationHelper_884(n.GetFirstChild())
}

// Java method: AstValidator.validateAmbientDeclarationHelper(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateAmbientDeclarationHelper_884(n pkg_rhino.Node) {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */ :
		this_.priv_AstValidator_validateNameDeclarationHelper_824(n.GetParent(), n.GetToken(), n)
	case int32(73) /* Token.FUNCTION */ :
		this_.priv_AstValidator_validateFunctionSignature_808(n)
	case int32(105) /* Token.CLASS */ :
		this_.priv_AstValidator_validateClassDeclaration_786(n, true)
	case int32(167) /* Token.ENUM */ :
		this_.priv_AstValidator_validateEnum_783(n)
	case int32(175) /* Token.NAMESPACE */ :
		this_.priv_AstValidator_validateNamespace_885(n, true)
	case int32(170) /* Token.TYPE_ALIAS */ :
		this_.priv_AstValidator_validateTypeAlias_882(n)
	case int32(120) /* Token.EXPORT */ :
		this_.priv_AstValidator_validateExport_772(n, true)
	default:
	}
}

// Java method: AstValidator.validateNamespace(com.google.javascript.rhino.Node n, boolean isAmbient)
func (this_ *AstValidator) priv_AstValidator_validateNamespace_885(n pkg_rhino.Node, isAmbient bool) {
	_ = n
	_ = isAmbient
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_NAMESPACE(), n)
	this_.priv_AstValidator_validateProperties_897(n)
	this_.priv_AstValidator_validateChildCount_891(n)
	this_.priv_AstValidator_validateNamespaceName_886(n.GetFirstChild())
	this_.priv_AstValidator_validateNamespaceElements_887(n.GetLastChild(), isAmbient)
}

// Java method: AstValidator.validateNamespaceName(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateNamespaceName_886(n pkg_rhino.Node) {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		this_.priv_AstValidator_validateName_802(n)
	case int32(26) /* Token.GETPROP */ :
		this_.priv_AstValidator_validateGetProp_857(n)
	default:
	}
}

// Java method: AstValidator.validateNamespaceElements(com.google.javascript.rhino.Node n, boolean isAmbient)
func (this_ *AstValidator) priv_AstValidator_validateNamespaceElements_887(n pkg_rhino.Node, isAmbient bool) {
	_ = n
	_ = isAmbient
	this_.priv_AstValidator_validateNodeType_890(pkg_rhino.StaticGet_Token_NAMESPACE_ELEMENTS(), n)
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			if isAmbient {
				this_.priv_AstValidator_validateAmbientDeclarationHelper_884(c)
			} else {
				this_.ValidateStatement_1(c)
			}
		}
	}
}

// Java method: AstValidator.violation(java.lang.String message, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_violation_888(message string, n pkg_rhino.Node) {
	_ = message
	_ = n
	this_.violationHandler.HandleViolation(message, n)
}

// Java method: AstValidator.validateFirstNodeOfOptChain(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateFirstNodeOfOptChain_889(n pkg_rhino.Node) {
	_ = n
	if !(NodeUtil_IsOptChainNode(n.GetFirstChild())) {
		if !(n.IsOptionalChainStart()) {
			this_.priv_AstValidator_violation_888(("Start of optional chain node " + j2gort.StrVal(n.GetToken()) + " is not marked as the start."), n)
		}
	}
}

// Java method: AstValidator.validateNodeType(com.google.javascript.rhino.Token type, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateNodeType_890(type_ pkg_rhino.Token, n pkg_rhino.Node) {
	_ = type_
	_ = n
	if n.GetToken() != type_ {
		this_.priv_AstValidator_violation_888(("Expected " + j2gort.StrVal(type_) + " but was " + j2gort.StrVal(n.GetToken())), n)
	}
}

// Java method: AstValidator.validateChildCount(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateChildCount_891(n pkg_rhino.Node) {
	_ = n
	var expectedArity int32 = pkg_rhino.Token_Arity(n.GetToken())
	_ = expectedArity
	if expectedArity != int32(-1) {
		this_.priv_AstValidator_validateChildCount_892(n, expectedArity)
	}
}

// Java method: AstValidator.validateChildCount(com.google.javascript.rhino.Node n, int expected)
func (this_ *AstValidator) priv_AstValidator_validateChildCount_892(n pkg_rhino.Node, expected int32) {
	_ = n
	_ = expected
	var count int32 = n.GetChildCount()
	_ = count
	if expected != count {
		this_.priv_AstValidator_violation_888(("Expected " + j2gort.StrValInt32(expected) + " children, but was " + j2gort.StrValInt32(count)), n)
	}
}

// Java method: AstValidator.validateChildCountIn(com.google.javascript.rhino.Node n, int min, int max)
func (this_ *AstValidator) priv_AstValidator_validateChildCountIn_893(n pkg_rhino.Node, min_ int32, max_ int32) {
	_ = n
	_ = min_
	_ = max_
	if max_ == min_ {
		this_.priv_AstValidator_validateChildCount_892(n, min_)
		return
	}
	var count int32 = n.GetChildCount()
	_ = count
	if (count < min_) || (count > max_) {
		this_.priv_AstValidator_violation_888(("Expected child count in [" + j2gort.StrValInt32(min_) + ", " + j2gort.StrValInt32(max_) + "], but was " + j2gort.StrValInt32(count)), n)
	}
}

// Java method: AstValidator.validateMinimumChildCount(com.google.javascript.rhino.Node n, int i)
func (this_ *AstValidator) priv_AstValidator_validateMinimumChildCount_894(n pkg_rhino.Node, i int32) {
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
		this_.priv_AstValidator_violation_888(("Expected at least " + j2gort.StrValInt32(i) + " children, but was " + j2gort.StrValInt32(n.GetChildCount())), n)
	}
}

// Java method: AstValidator.validateMaximumChildCount(com.google.javascript.rhino.Node n, int i)
func (this_ *AstValidator) priv_AstValidator_validateMaximumChildCount_895(n pkg_rhino.Node, i int32) {
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
		this_.priv_AstValidator_violation_888(("Expected no more than " + j2gort.StrValInt32(i) + " children, but was " + j2gort.StrValInt32(n.GetChildCount())), n)
	}
}

// Java method: AstValidator.validateFeature(com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature feature, com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateFeature_896(feature *pkg_parser.FeatureSet_Feature, n pkg_rhino.Node) {
	_ = feature
	_ = n
	var allowbleFeatures *pkg_parser.FeatureSet = this_.compiler.GetAllowableFeatures()
	_ = allowbleFeatures
	if !(n.IsFromExterns()) && !(allowbleFeatures.Has(feature)) {
		this_.priv_AstValidator_violation_888(("AST should not contain " + j2gort.StrVal(feature)), n)
	}
	if !(this_.isScriptFeatureValidationEnabled) || (this_.currentScript == 0) {
		return
	}
	var scriptFeatures *pkg_parser.FeatureSet = NodeUtil_GetFeatureSetOfScript(this_.currentScript)
	_ = scriptFeatures
	if (scriptFeatures == nil) || !(scriptFeatures.Has(feature)) {
		this_.priv_AstValidator_violation_888(("SCRIPT node should be marked as containing feature " + j2gort.StrVal(feature)), this_.currentScript)
	}
}

// Java method: AstValidator.validateProperties(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateProperties_897(n pkg_rhino.Node) {
	_ = n
	if n.HasPropertyViolations() {
		n.ValidateProperties(j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_29 any) {
			errorMessage := j2gort.CastString(lp_29)
			_ = errorMessage
			this_.priv_AstValidator_violation_888(errorMessage, n)
		}))
	}
}

// Java method: AstValidator.validateRequiredInlinings(com.google.javascript.rhino.Node n)
func (this_ *AstValidator) priv_AstValidator_validateRequiredInlinings_898(n pkg_rhino.Node) {
	_ = n
	if !(this_.shouldValidateRequiredInlinings) {
		return
	}
	var jsdocInfo *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = jsdocInfo
	if (jsdocInfo != nil) && jsdocInfo.IsRequireInlining() {
		this_.priv_AstValidator_violation_888("@requireInlining node failed to be inlined.", n)
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

func Upcast_AstValidator_ViolationHandler[T AstValidator_ViolationHandler](v T) AstValidator_ViolationHandler {
	var vi AstValidator_ViolationHandler = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
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
