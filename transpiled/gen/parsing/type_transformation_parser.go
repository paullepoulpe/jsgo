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
// Port of com/google/javascript/jscomp/parsing/TypeTransformationParser.java (generated mechanically by j2go).

package parsing

import (
	"math"
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.parsing.TypeTransformationParser
type TypeTransformationParser struct {
	typeTransformationString string                     `gson:"typeTransformationString"`
	sourceFile               pkg_rhino.StaticSourceFile `gson:"sourceFile"`
	errorReporter            pkg_rhino.ErrorReporter    `gson:"errorReporter"`
	typeTransformationAst    pkg_rhino.Node             `gson:"typeTransformationAst"`
	templateLineno           int32                      `gson:"templateLineno"`
	templateCharno           int32                      `gson:"templateCharno"`
	idHash_                  atomic.Int32
}

func (this_ *TypeTransformationParser) BindSelf_(s any) {
	_ = s
}
func (this_ *TypeTransformationParser) IsJavaNil_() bool { return this_ == nil }
func (*TypeTransformationParser) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.TypeTransformationParser", "TypeTransformationParser"
}
func (this_ *TypeTransformationParser) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TypeTransformationParser) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.TypeTransformationParser", "com.google.javascript.jscomp.parsing.TypeTransformationParser", func(v_ any) bool {
		return j2gort.InstanceOf[*TypeTransformationParser](v_)
	})
}

func StaticGet_TypeTransformationParser_VAR_ARGS() int32 { return int32(2147483647) }

var static_TypeTransformationParser_TTL_NODE_LENGTH int32

func StaticGet_TypeTransformationParser_TTL_NODE_LENGTH() int32 {
	clinit_TypeTransformationParser()
	return static_TypeTransformationParser_TTL_NODE_LENGTH
}
func StaticSet_TypeTransformationParser_TTL_NODE_LENGTH(v int32) {
	clinit_TypeTransformationParser()
	static_TypeTransformationParser_TTL_NODE_LENGTH = v
}

var clinitDone_TypeTransformationParser uint32
var clinitState_TypeTransformationParser j2gort.ClinitState

func clinit_TypeTransformationParser() {
	if atomic.LoadUint32(&clinitDone_TypeTransformationParser) != 0 {
		return
	}
	clinitSlow_TypeTransformationParser()
}

func clinitSlow_TypeTransformationParser() {
	if !j2gort.EnterClinit(&clinitState_TypeTransformationParser, &clinitDone_TypeTransformationParser) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_TypeTransformationParser, &clinitDone_TypeTransformationParser)
	static_TypeTransformationParser_TTL_NODE_LENGTH = j2gort.JStrLen("@template")
}

func (this_ *TypeTransformationParser) initFields_TypeTransformationParser() {
	_ = this_
	this_.typeTransformationString = j2gort.NullString
}

// Java constructor: TypeTransformationParser(java.lang.String typeTransformationString, com.google.javascript.rhino.StaticSourceFile sourceFile, com.google.javascript.rhino.ErrorReporter errorReporter, int templateLineno, int templateCharno)
func NewTypeTransformationParser(typeTransformationString string, sourceFile pkg_rhino.StaticSourceFile, errorReporter pkg_rhino.ErrorReporter, templateLineno int32, templateCharno int32) *TypeTransformationParser {
	clinit_TypeTransformationParser()
	this_ := &TypeTransformationParser{}
	InitTypeTransformationParser(this_, typeTransformationString, sourceFile, errorReporter, templateLineno, templateCharno)
	return this_
}

func InitTypeTransformationParser(this_ *TypeTransformationParser, typeTransformationString string, sourceFile pkg_rhino.StaticSourceFile, errorReporter pkg_rhino.ErrorReporter, templateLineno int32, templateCharno int32) {
	_ = typeTransformationString
	_ = sourceFile
	_ = errorReporter
	_ = templateLineno
	_ = templateCharno
	this_.initFields_TypeTransformationParser()
	this_.typeTransformationString = typeTransformationString
	this_.sourceFile = sourceFile
	this_.errorReporter = errorReporter
	this_.templateLineno = templateLineno
	this_.templateCharno = templateCharno
}

// Java method: TypeTransformationParser.getTypeTransformationAst()
func (this_ *TypeTransformationParser) GetTypeTransformationAst() pkg_rhino.Node {
	return this_.typeTransformationAst
}

// Java method: TypeTransformationParser.addNewWarning(com.google.javascript.rhino.Msg messageId, java.lang.String messageArg)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_addNewWarning_7181(messageId pkg_rhino.Msg, messageArg string) {
	_ = messageId
	_ = messageArg
	this_.errorReporter.Warning(("Bad type annotation. " + j2gort.StrValString(messageId.Format_1([]any{j2gort.BoxString(messageArg)}))), this_.sourceFile.GetName(), this_.templateLineno, this_.templateCharno)
}

// Java method: TypeTransformationParser.nameToKeyword(java.lang.String s)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_nameToKeyword_7182(s string) *TypeTransformationParser_Keywords {
	_ = s
	return TypeTransformationParser_Keywords_ValueOf(j2gort.Ext_com_google_common_base_Ascii_ToUpperCase_1_string(s))
}

// Java method: TypeTransformationParser.isValidKeyword(java.lang.String name)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_isValidKeyword_7183(name string) bool {
	_ = name
	for _, el_1 := range j2gort.CheckNotNullSlice(TypeTransformationParser_Keywords_Values()) {
		k := el_1
		_ = k
		if j2gort.CheckNotNullString(k.name) == name {
			return true
		}
	}
	return false
}

// Java method: TypeTransformationParser.isOperationKind(java.lang.String name, com.google.javascript.jscomp.parsing.TypeTransformationParser.OperationKind kind)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_isOperationKind_7184(name string, kind TypeTransformationParser_OperationKind) bool {
	_ = name
	_ = kind
	return (this_.priv_TypeTransformationParser_isValidKeyword_7183(name) && (this_.priv_TypeTransformationParser_nameToKeyword_7182(name).kind == kind))
}

// Java method: TypeTransformationParser.isValidStringPredicate(java.lang.String name)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_isValidStringPredicate_7185(name string) bool {
	_ = name
	return this_.priv_TypeTransformationParser_isOperationKind_7184(name, StaticGet_TypeTransformationParser_OperationKind_STRING_PREDICATE())
}

// Java method: TypeTransformationParser.isValidTypePredicate(java.lang.String name)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_isValidTypePredicate_7186(name string) bool {
	_ = name
	return this_.priv_TypeTransformationParser_isOperationKind_7184(name, StaticGet_TypeTransformationParser_OperationKind_TYPE_PREDICATE())
}

// Java method: TypeTransformationParser.isValidTypevarPredicate(java.lang.String name)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_isValidTypevarPredicate_7187(name string) bool {
	_ = name
	return this_.priv_TypeTransformationParser_isOperationKind_7184(name, StaticGet_TypeTransformationParser_OperationKind_TYPEVAR_PREDICATE())
}

// Java method: TypeTransformationParser.isBooleanOperation(com.google.javascript.rhino.Node n)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_isBooleanOperation_7188(n pkg_rhino.Node) bool {
	_ = n
	return ((n.IsAnd() || n.IsOr()) || n.IsNot())
}

// Java method: TypeTransformationParser.isValidPredicate(java.lang.String name)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_isValidPredicate_7189(name string) bool {
	_ = name
	return ((this_.priv_TypeTransformationParser_isValidStringPredicate_7185(name) || this_.priv_TypeTransformationParser_isValidTypePredicate_7186(name)) || this_.priv_TypeTransformationParser_isValidTypevarPredicate_7187(name))
}

// Java method: TypeTransformationParser.getFunctionParamCount(com.google.javascript.rhino.Node n)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_getFunctionParamCount_7190(n pkg_rhino.Node) int32 {
	_ = n
	if !(n.IsFunction()) {
		j2gort.CheckArgumentFmt(false, "Expected a function node, found %s", j2gort.BoxNode(n))
	}
	return n.GetSecondChild().GetChildCount()
}

// Java method: TypeTransformationParser.getFunctionBody(com.google.javascript.rhino.Node n)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_getFunctionBody_7191(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if !(n.IsFunction()) {
		j2gort.CheckArgumentFmt(false, "Expected a function node, found %s", j2gort.BoxNode(n))
	}
	return n.GetChildAtIndex(int32(2))
}

// Java method: TypeTransformationParser.getCallName(com.google.javascript.rhino.Node n)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_getCallName_7192(n pkg_rhino.Node) string {
	_ = n
	if !(n.IsCall()) {
		j2gort.CheckArgumentFmt(false, "Expected a call node, found %s", j2gort.BoxNode(n))
	}
	return n.GetFirstChild().GetString()
}

// Java method: TypeTransformationParser.getCallArgument(com.google.javascript.rhino.Node n, int i)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_getCallArgument_7193(n pkg_rhino.Node, i int32) pkg_rhino.Node {
	_ = n
	_ = i
	if !(n.IsCall()) {
		j2gort.CheckArgumentFmt(false, "Expected a call node, found %s", j2gort.BoxNode(n))
	}
	return n.GetChildAtIndex((i + int32(1)))
}

// Java method: TypeTransformationParser.getCallParamCount(com.google.javascript.rhino.Node n)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_getCallParamCount_7194(n pkg_rhino.Node) int32 {
	_ = n
	if !(n.IsCall()) {
		j2gort.CheckArgumentFmt(false, "Expected a call node, found %s", j2gort.BoxNode(n))
	}
	return (n.GetChildCount() - int32(1))
}

// Java method: TypeTransformationParser.isTypeVar(com.google.javascript.rhino.Node n)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_isTypeVar_7195(n pkg_rhino.Node) bool {
	_ = n
	return n.IsName()
}

// Java method: TypeTransformationParser.isTypeName(com.google.javascript.rhino.Node n)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_isTypeName_7196(n pkg_rhino.Node) bool {
	_ = n
	return n.IsStringLit()
}

// Java method: TypeTransformationParser.isOperation(com.google.javascript.rhino.Node n)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_isOperation_7197(n pkg_rhino.Node) bool {
	_ = n
	return n.IsCall()
}

// Java method: TypeTransformationParser.isValidExpression(com.google.javascript.rhino.Node e)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_isValidExpression_7198(e pkg_rhino.Node) bool {
	_ = e
	return ((this_.priv_TypeTransformationParser_isTypeVar_7195(e) || this_.priv_TypeTransformationParser_isTypeName_7196(e)) || this_.priv_TypeTransformationParser_isOperation_7197(e))
}

// Java method: TypeTransformationParser.warnInvalid(java.lang.String msg)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_warnInvalid_7199(msg string) {
	_ = msg
	this_.priv_TypeTransformationParser_addNewWarning_7181(pkg_rhino.StaticGet_Msg_JSDOC_TYPETRANSFORMATION_INVALID(), msg)
}

// Java method: TypeTransformationParser.warnInvalidExpression(java.lang.String msg)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_warnInvalidExpression_7200(msg string) {
	_ = msg
	this_.priv_TypeTransformationParser_addNewWarning_7181(pkg_rhino.StaticGet_Msg_JSDOC_TYPETRANSFORMATION_INVALID_EXPRESSION(), msg)
}

// Java method: TypeTransformationParser.warnMissingParam(java.lang.String msg)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_warnMissingParam_7201(msg string) {
	_ = msg
	this_.priv_TypeTransformationParser_addNewWarning_7181(pkg_rhino.StaticGet_Msg_JSDOC_TYPETRANSFORMATION_MISSING_PARAM(), msg)
}

// Java method: TypeTransformationParser.warnExtraParam(java.lang.String msg)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_warnExtraParam_7202(msg string) {
	_ = msg
	this_.priv_TypeTransformationParser_addNewWarning_7181(pkg_rhino.StaticGet_Msg_JSDOC_TYPETRANSFORMATION_EXTRA_PARAM(), msg)
}

// Java method: TypeTransformationParser.warnInvalidInside(java.lang.String msg)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_warnInvalidInside_7203(msg string) {
	_ = msg
	this_.priv_TypeTransformationParser_addNewWarning_7181(pkg_rhino.StaticGet_Msg_JSDOC_TYPETRANSFORMATION_INVALID_INSIDE(), msg)
}

// Java method: TypeTransformationParser.checkParameterCount(com.google.javascript.rhino.Node expr, com.google.javascript.jscomp.parsing.TypeTransformationParser.Keywords keyword)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_checkParameterCount_7204(expr pkg_rhino.Node, keyword *TypeTransformationParser_Keywords) bool {
	_ = expr
	_ = keyword
	var paramCount int32 = this_.priv_TypeTransformationParser_getCallParamCount_7194(expr)
	_ = paramCount
	if paramCount < keyword.minParamCount {
		this_.priv_TypeTransformationParser_warnMissingParam_7201(keyword.name)
		return false
	}
	if paramCount > keyword.maxParamCount {
		this_.priv_TypeTransformationParser_warnExtraParam_7202(keyword.name)
		return false
	}
	return true
}

// Java method: TypeTransformationParser.parseTypeTransformation()
func (this_ *TypeTransformationParser) ParseTypeTransformation() bool {
	var config Config = Config_S_Builder().SetLanguageMode(StaticGet_Config_LanguageMode_ES_NEXT()).SetStrictMode(StaticGet_Config_StrictMode_SLOPPY()).Build()
	_ = config
	var result *ParserRunner_ParseResult = ParserRunner_Parse(this_.sourceFile, this_.typeTransformationString, config, this_.errorReporter)
	_ = result
	var ast pkg_rhino.Node = result.ast
	_ = ast
	if (((ast == 0) || !(ast.IsScript())) || !(ast.HasChildren())) || !(ast.GetFirstChild().IsExprResult()) {
		this_.priv_TypeTransformationParser_warnInvalidExpression_7200("type transformation")
		return false
	}
	var expr pkg_rhino.Node = ast.GetFirstFirstChild()
	_ = expr
	if !(this_.priv_TypeTransformationParser_validTypeTransformationExpression_7233(expr)) {
		return false
	}
	this_.priv_TypeTransformationParser_fixTTLNodeLineNoCharNoAndLength_7206(expr)
	this_.typeTransformationAst = expr
	return true
}

// Java method: TypeTransformationParser.fixTTLNodeLineNoCharNoAndLength(com.google.javascript.rhino.Node expr)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_fixTTLNodeLineNoCharNoAndLength_7206(expr pkg_rhino.Node) {
	_ = expr
	expr.SetLinenoCharno(this_.templateLineno, this_.templateCharno)
	expr.SetLength(StaticGet_TypeTransformationParser_TTL_NODE_LENGTH())
	{
		var child pkg_rhino.Node = expr.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			this_.priv_TypeTransformationParser_fixTTLNodeLineNoCharNoAndLength_7206(child)
		}
	}
}

// Java method: TypeTransformationParser.validTemplateTypeExpression(com.google.javascript.rhino.Node expr)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_validTemplateTypeExpression_7207(expr pkg_rhino.Node) bool {
	_ = expr
	if !(this_.priv_TypeTransformationParser_checkParameterCount_7204(expr, StaticGet_TypeTransformationParser_Keywords_TYPE())) {
		return false
	}
	var paramCount int32 = this_.priv_TypeTransformationParser_getCallParamCount_7194(expr)
	_ = paramCount
	var firstParam pkg_rhino.Node = this_.priv_TypeTransformationParser_getCallArgument_7193(expr, int32(0))
	_ = firstParam
	if !(this_.priv_TypeTransformationParser_isTypeVar_7195(firstParam)) && !(this_.priv_TypeTransformationParser_isTypeName_7196(firstParam)) {
		this_.priv_TypeTransformationParser_warnInvalid_7199("type name or type variable")
		this_.priv_TypeTransformationParser_warnInvalidInside_7203("template type operation")
		return false
	}
	{
		var i int32 = int32(1)
		_ = i
		for ; i < paramCount; i = i + 1 {
			if !(this_.priv_TypeTransformationParser_validTypeTransformationExpression_7233(this_.priv_TypeTransformationParser_getCallArgument_7193(expr, i))) {
				this_.priv_TypeTransformationParser_warnInvalidInside_7203("template type operation")
				return false
			}
		}
	}
	return true
}

// Java method: TypeTransformationParser.validUnionTypeExpression(com.google.javascript.rhino.Node expr)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_validUnionTypeExpression_7208(expr pkg_rhino.Node) bool {
	_ = expr
	if !(this_.priv_TypeTransformationParser_checkParameterCount_7204(expr, StaticGet_TypeTransformationParser_Keywords_UNION())) {
		return false
	}
	var paramCount int32 = this_.priv_TypeTransformationParser_getCallParamCount_7194(expr)
	_ = paramCount
	{
		var i int32 = int32(0)
		_ = i
		for ; i < paramCount; i = i + 1 {
			if !(this_.priv_TypeTransformationParser_validTypeTransformationExpression_7233(this_.priv_TypeTransformationParser_getCallArgument_7193(expr, i))) {
				this_.priv_TypeTransformationParser_warnInvalidInside_7203("union type")
				return false
			}
		}
	}
	return true
}

// Java method: TypeTransformationParser.validNoneTypeExpression(com.google.javascript.rhino.Node expr)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_validNoneTypeExpression_7209(expr pkg_rhino.Node) bool {
	_ = expr
	return this_.priv_TypeTransformationParser_checkParameterCount_7204(expr, StaticGet_TypeTransformationParser_Keywords_NONE())
}

// Java method: TypeTransformationParser.validAllTypeExpression(com.google.javascript.rhino.Node expr)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_validAllTypeExpression_7210(expr pkg_rhino.Node) bool {
	_ = expr
	return this_.priv_TypeTransformationParser_checkParameterCount_7204(expr, StaticGet_TypeTransformationParser_Keywords_ALL())
}

// Java method: TypeTransformationParser.validUnknownTypeExpression(com.google.javascript.rhino.Node expr)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_validUnknownTypeExpression_7211(expr pkg_rhino.Node) bool {
	_ = expr
	return this_.priv_TypeTransformationParser_checkParameterCount_7204(expr, StaticGet_TypeTransformationParser_Keywords_UNKNOWN())
}

// Java method: TypeTransformationParser.validRawTypeOfTypeExpression(com.google.javascript.rhino.Node expr)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_validRawTypeOfTypeExpression_7212(expr pkg_rhino.Node) bool {
	_ = expr
	if !(this_.priv_TypeTransformationParser_checkParameterCount_7204(expr, StaticGet_TypeTransformationParser_Keywords_RAWTYPEOF())) {
		return false
	}
	if !(this_.priv_TypeTransformationParser_validTypeTransformationExpression_7233(this_.priv_TypeTransformationParser_getCallArgument_7193(expr, int32(0)))) {
		this_.priv_TypeTransformationParser_warnInvalidInside_7203(StaticGet_TypeTransformationParser_Keywords_RAWTYPEOF().name)
		return false
	}
	return true
}

// Java method: TypeTransformationParser.validTemplateTypeOfExpression(com.google.javascript.rhino.Node expr)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_validTemplateTypeOfExpression_7213(expr pkg_rhino.Node) bool {
	_ = expr
	if !(this_.priv_TypeTransformationParser_checkParameterCount_7204(expr, StaticGet_TypeTransformationParser_Keywords_TEMPLATETYPEOF())) {
		return false
	}
	if !(this_.priv_TypeTransformationParser_validTypeTransformationExpression_7233(this_.priv_TypeTransformationParser_getCallArgument_7193(expr, int32(0)))) {
		this_.priv_TypeTransformationParser_warnInvalidInside_7203(StaticGet_TypeTransformationParser_Keywords_TEMPLATETYPEOF().name)
		return false
	}
	if !(this_.priv_TypeTransformationParser_getCallArgument_7193(expr, int32(1)).IsNumber()) {
		this_.priv_TypeTransformationParser_warnInvalid_7199("index")
		this_.priv_TypeTransformationParser_warnInvalidInside_7203(StaticGet_TypeTransformationParser_Keywords_TEMPLATETYPEOF().name)
		return false
	}
	var index float64 = this_.priv_TypeTransformationParser_getCallArgument_7193(expr, int32(1)).GetDouble()
	_ = index
	if math.IsNaN(index) || !(pkg_base.JSCompDoubles_IsExactInt32(index)) {
		this_.priv_TypeTransformationParser_warnInvalid_7199("index")
		this_.priv_TypeTransformationParser_warnInvalidInside_7203(StaticGet_TypeTransformationParser_Keywords_TEMPLATETYPEOF().name)
		return false
	}
	return true
}

// Java method: TypeTransformationParser.validRecordParam(com.google.javascript.rhino.Node expr)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_validRecordParam_7214(expr pkg_rhino.Node) bool {
	_ = expr
	if expr.IsObjectLit() {
		{
			var prop pkg_rhino.Node = expr.GetFirstChild()
			_ = prop
			for ; prop != 0; prop = prop.GetNext() {
				if prop.IsShorthandProperty() {
					this_.priv_TypeTransformationParser_warnInvalid_7199("property, missing type")
					return false
				} else if !(this_.priv_TypeTransformationParser_validTypeTransformationExpression_7233(prop.GetFirstChild())) {
					return false
				}
			}
		}
	} else if !(this_.priv_TypeTransformationParser_validTypeTransformationExpression_7233(expr)) {
		return false
	}
	return true
}

// Java method: TypeTransformationParser.validRecordTypeExpression(com.google.javascript.rhino.Node expr)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_validRecordTypeExpression_7215(expr pkg_rhino.Node) bool {
	_ = expr
	if !(this_.priv_TypeTransformationParser_checkParameterCount_7204(expr, StaticGet_TypeTransformationParser_Keywords_RECORD())) {
		return false
	}
	{
		var i int32 = int32(0)
		_ = i
		for ; i < this_.priv_TypeTransformationParser_getCallParamCount_7194(expr); i = i + 1 {
			if !(this_.priv_TypeTransformationParser_validRecordParam_7214(this_.priv_TypeTransformationParser_getCallArgument_7193(expr, i))) {
				this_.priv_TypeTransformationParser_warnInvalidInside_7203(StaticGet_TypeTransformationParser_Keywords_RECORD().name)
				return false
			}
		}
	}
	return true
}

// Java method: TypeTransformationParser.validNativeTypeExpr(com.google.javascript.rhino.Node expr)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_validNativeTypeExpr_7216(expr pkg_rhino.Node) bool {
	_ = expr
	if !(this_.priv_TypeTransformationParser_checkParameterCount_7204(expr, StaticGet_TypeTransformationParser_Keywords_TYPEEXPR())) {
		return false
	}
	var typeString pkg_rhino.Node = this_.priv_TypeTransformationParser_getCallArgument_7193(expr, int32(0))
	_ = typeString
	if !(typeString.IsStringLit()) {
		this_.priv_TypeTransformationParser_warnInvalidExpression_7200("native type")
		this_.priv_TypeTransformationParser_warnInvalidInside_7203(StaticGet_TypeTransformationParser_Keywords_TYPEEXPR().name)
		return false
	}
	var typeExpr pkg_rhino.Node = JsDocInfoParser_ParseTypeString(typeString.GetString())
	_ = typeExpr
	typeString.Detach()
	expr.AddChildToBack(typeExpr)
	return true
}

// Java method: TypeTransformationParser.validTypeExpression(com.google.javascript.rhino.Node expr)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_validTypeExpression_7217(expr pkg_rhino.Node) bool {
	_ = expr
	var name string = this_.priv_TypeTransformationParser_getCallName_7192(expr)
	_ = name
	var keyword *TypeTransformationParser_Keywords = this_.priv_TypeTransformationParser_nameToKeyword_7182(name)
	_ = keyword
	switch j2gort.CheckNotNullPtr(keyword).Ordinal() {
	case int32(19) /* Keywords.TYPE */ :
		return this_.priv_TypeTransformationParser_validTemplateTypeExpression_7207(expr)
	case int32(22) /* Keywords.UNION */ :
		return this_.priv_TypeTransformationParser_validUnionTypeExpression_7208(expr)
	case int32(11) /* Keywords.NONE */ :
		return this_.priv_TypeTransformationParser_validNoneTypeExpression_7209(expr)
	case int32(0) /* Keywords.ALL */ :
		return this_.priv_TypeTransformationParser_validAllTypeExpression_7210(expr)
	case int32(23) /* Keywords.UNKNOWN */ :
		return this_.priv_TypeTransformationParser_validUnknownTypeExpression_7211(expr)
	case int32(14) /* Keywords.RAWTYPEOF */ :
		return this_.priv_TypeTransformationParser_validRawTypeOfTypeExpression_7212(expr)
	case int32(18) /* Keywords.TEMPLATETYPEOF */ :
		return this_.priv_TypeTransformationParser_validTemplateTypeOfExpression_7213(expr)
	case int32(17) /* Keywords.RECORD */ :
		return this_.priv_TypeTransformationParser_validRecordTypeExpression_7215(expr)
	case int32(20) /* Keywords.TYPEEXPR */ :
		return this_.priv_TypeTransformationParser_validNativeTypeExpr_7216(expr)
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Invalid type expression")))
	}
	panic("unreachable switch expression")
}

// Java method: TypeTransformationParser.validTypePredicate(com.google.javascript.rhino.Node expr, int paramCount)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_validTypePredicate_7218(expr pkg_rhino.Node, paramCount int32) bool {
	_ = expr
	_ = paramCount
	{
		var i int32 = int32(0)
		_ = i
		for ; i < paramCount; i = i + 1 {
			if !(this_.priv_TypeTransformationParser_validTypeTransformationExpression_7233(this_.priv_TypeTransformationParser_getCallArgument_7193(expr, i))) {
				this_.priv_TypeTransformationParser_warnInvalidInside_7203("boolean")
				return false
			}
		}
	}
	return true
}

// Java method: TypeTransformationParser.isValidStringParam(com.google.javascript.rhino.Node expr)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_isValidStringParam_7219(expr pkg_rhino.Node) bool {
	_ = expr
	if !(expr.IsName()) && !(expr.IsStringLit()) {
		this_.priv_TypeTransformationParser_warnInvalid_7199("string")
		return false
	}
	if j2gort.CheckNotNullString(expr.GetString()) == "" {
		this_.priv_TypeTransformationParser_warnInvalid_7199("string parameter")
		return false
	}
	return true
}

// Java method: TypeTransformationParser.validStringPredicate(com.google.javascript.rhino.Node expr, int paramCount)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_validStringPredicate_7220(expr pkg_rhino.Node, paramCount int32) bool {
	_ = expr
	_ = paramCount
	{
		var i int32 = int32(0)
		_ = i
		for ; i < paramCount; i = i + 1 {
			if !(this_.priv_TypeTransformationParser_isValidStringParam_7219(this_.priv_TypeTransformationParser_getCallArgument_7193(expr, i))) {
				this_.priv_TypeTransformationParser_warnInvalidInside_7203("boolean")
				return false
			}
		}
	}
	return true
}

// Java method: TypeTransformationParser.validTypevarParam(com.google.javascript.rhino.Node expr)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_validTypevarParam_7221(expr pkg_rhino.Node) bool {
	_ = expr
	if !(this_.priv_TypeTransformationParser_isTypeVar_7195(expr)) {
		this_.priv_TypeTransformationParser_warnInvalid_7199("name")
		return false
	}
	return true
}

// Java method: TypeTransformationParser.validTypevarPredicate(com.google.javascript.rhino.Node expr, int paramCount)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_validTypevarPredicate_7222(expr pkg_rhino.Node, paramCount int32) bool {
	_ = expr
	_ = paramCount
	{
		var i int32 = int32(0)
		_ = i
		for ; i < paramCount; i = i + 1 {
			if !(this_.priv_TypeTransformationParser_validTypevarParam_7221(this_.priv_TypeTransformationParser_getCallArgument_7193(expr, i))) {
				this_.priv_TypeTransformationParser_warnInvalidInside_7203("boolean")
				return false
			}
		}
	}
	return true
}

// Java method: TypeTransformationParser.validBooleanOperation(com.google.javascript.rhino.Node expr)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_validBooleanOperation_7223(expr pkg_rhino.Node) bool {
	_ = expr
	var valid bool = false
	_ = valid
	if expr.IsNot() {
		valid = this_.priv_TypeTransformationParser_validBooleanExpression_7224(expr.GetFirstChild())
	} else {
		valid = (this_.priv_TypeTransformationParser_validBooleanExpression_7224(expr.GetFirstChild()) && this_.priv_TypeTransformationParser_validBooleanExpression_7224(expr.GetSecondChild()))
	}
	if !(valid) {
		this_.priv_TypeTransformationParser_warnInvalidInside_7203("boolean")
		return false
	}
	return true
}

// Java method: TypeTransformationParser.validBooleanExpression(com.google.javascript.rhino.Node expr)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_validBooleanExpression_7224(expr pkg_rhino.Node) bool {
	_ = expr
	if this_.priv_TypeTransformationParser_isBooleanOperation_7188(expr) {
		return this_.priv_TypeTransformationParser_validBooleanOperation_7223(expr)
	}
	if !(this_.priv_TypeTransformationParser_isOperation_7197(expr)) {
		this_.priv_TypeTransformationParser_warnInvalidExpression_7200("boolean")
		return false
	}
	if !(this_.priv_TypeTransformationParser_isValidPredicate_7189(this_.priv_TypeTransformationParser_getCallName_7192(expr))) {
		this_.priv_TypeTransformationParser_warnInvalid_7199("boolean predicate")
		return false
	}
	var keyword *TypeTransformationParser_Keywords = this_.priv_TypeTransformationParser_nameToKeyword_7182(this_.priv_TypeTransformationParser_getCallName_7192(expr))
	_ = keyword
	if !(this_.priv_TypeTransformationParser_checkParameterCount_7204(expr, keyword)) {
		return false
	}
	switch (keyword.kind).Ordinal() {
	case int32(3) /* OperationKind.TYPE_PREDICATE */ :
		return this_.priv_TypeTransformationParser_validTypePredicate_7218(expr, this_.priv_TypeTransformationParser_getCallParamCount_7194(expr))
	case int32(2) /* OperationKind.STRING_PREDICATE */ :
		return this_.priv_TypeTransformationParser_validStringPredicate_7220(expr, this_.priv_TypeTransformationParser_getCallParamCount_7194(expr))
	case int32(4) /* OperationKind.TYPEVAR_PREDICATE */ :
		return this_.priv_TypeTransformationParser_validTypevarPredicate_7222(expr, this_.priv_TypeTransformationParser_getCallParamCount_7194(expr))
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Invalid boolean expression")))
	}
	panic("unreachable switch expression")
}

// Java method: TypeTransformationParser.validConditionalExpression(com.google.javascript.rhino.Node expr)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_validConditionalExpression_7225(expr pkg_rhino.Node) bool {
	_ = expr
	if !(this_.priv_TypeTransformationParser_checkParameterCount_7204(expr, StaticGet_TypeTransformationParser_Keywords_COND())) {
		return false
	}
	if !(this_.priv_TypeTransformationParser_validBooleanExpression_7224(this_.priv_TypeTransformationParser_getCallArgument_7193(expr, int32(0)))) {
		this_.priv_TypeTransformationParser_warnInvalidInside_7203("conditional")
		return false
	}
	if !(this_.priv_TypeTransformationParser_validTypeTransformationExpression_7233(this_.priv_TypeTransformationParser_getCallArgument_7193(expr, int32(1)))) {
		this_.priv_TypeTransformationParser_warnInvalidInside_7203("conditional")
		return false
	}
	if !(this_.priv_TypeTransformationParser_validTypeTransformationExpression_7233(this_.priv_TypeTransformationParser_getCallArgument_7193(expr, int32(2)))) {
		this_.priv_TypeTransformationParser_warnInvalidInside_7203("conditional")
		return false
	}
	return true
}

// Java method: TypeTransformationParser.validMapunionExpression(com.google.javascript.rhino.Node expr)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_validMapunionExpression_7226(expr pkg_rhino.Node) bool {
	_ = expr
	if !(this_.priv_TypeTransformationParser_checkParameterCount_7204(expr, StaticGet_TypeTransformationParser_Keywords_MAPUNION())) {
		return false
	}
	if !(this_.priv_TypeTransformationParser_validTypeTransformationExpression_7233(this_.priv_TypeTransformationParser_getCallArgument_7193(expr, int32(0)))) {
		this_.priv_TypeTransformationParser_warnInvalidInside_7203(StaticGet_TypeTransformationParser_Keywords_MAPUNION().name)
		return false
	}
	if !(this_.priv_TypeTransformationParser_getCallArgument_7193(expr, int32(1)).IsFunction()) {
		this_.priv_TypeTransformationParser_warnInvalid_7199("map function")
		this_.priv_TypeTransformationParser_warnInvalidInside_7203(StaticGet_TypeTransformationParser_Keywords_MAPUNION().name)
		return false
	}
	var mapFn pkg_rhino.Node = this_.priv_TypeTransformationParser_getCallArgument_7193(expr, int32(1))
	_ = mapFn
	var mapFnParamCount int32 = this_.priv_TypeTransformationParser_getFunctionParamCount_7190(mapFn)
	_ = mapFnParamCount
	if mapFnParamCount < int32(1) {
		this_.priv_TypeTransformationParser_warnMissingParam_7201("map function")
		this_.priv_TypeTransformationParser_warnInvalidInside_7203(StaticGet_TypeTransformationParser_Keywords_MAPUNION().name)
		return false
	}
	if mapFnParamCount > int32(1) {
		this_.priv_TypeTransformationParser_warnExtraParam_7202("map function")
		this_.priv_TypeTransformationParser_warnInvalidInside_7203(StaticGet_TypeTransformationParser_Keywords_MAPUNION().name)
		return false
	}
	var mapFnBody pkg_rhino.Node = this_.priv_TypeTransformationParser_getFunctionBody_7191(mapFn)
	_ = mapFnBody
	if !(this_.priv_TypeTransformationParser_validTypeTransformationExpression_7233(mapFnBody)) {
		this_.priv_TypeTransformationParser_warnInvalidInside_7203("map function body")
		return false
	}
	return true
}

// Java method: TypeTransformationParser.validMaprecordExpression(com.google.javascript.rhino.Node expr)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_validMaprecordExpression_7227(expr pkg_rhino.Node) bool {
	_ = expr
	if !(this_.priv_TypeTransformationParser_checkParameterCount_7204(expr, StaticGet_TypeTransformationParser_Keywords_MAPRECORD())) {
		return false
	}
	if !(this_.priv_TypeTransformationParser_validTypeTransformationExpression_7233(this_.priv_TypeTransformationParser_getCallArgument_7193(expr, int32(0)))) {
		this_.priv_TypeTransformationParser_warnInvalidInside_7203(StaticGet_TypeTransformationParser_Keywords_MAPRECORD().name)
		return false
	}
	if !(this_.priv_TypeTransformationParser_getCallArgument_7193(expr, int32(1)).IsFunction()) {
		this_.priv_TypeTransformationParser_warnInvalid_7199("map function")
		this_.priv_TypeTransformationParser_warnInvalidInside_7203(StaticGet_TypeTransformationParser_Keywords_MAPRECORD().name)
		return false
	}
	var mapFn pkg_rhino.Node = this_.priv_TypeTransformationParser_getCallArgument_7193(expr, int32(1))
	_ = mapFn
	var mapFnParamCount int32 = this_.priv_TypeTransformationParser_getFunctionParamCount_7190(mapFn)
	_ = mapFnParamCount
	if mapFnParamCount < int32(2) {
		this_.priv_TypeTransformationParser_warnMissingParam_7201("map function")
		this_.priv_TypeTransformationParser_warnInvalidInside_7203(StaticGet_TypeTransformationParser_Keywords_MAPRECORD().name)
		return false
	}
	if mapFnParamCount > int32(2) {
		this_.priv_TypeTransformationParser_warnExtraParam_7202("map function")
		this_.priv_TypeTransformationParser_warnInvalidInside_7203(StaticGet_TypeTransformationParser_Keywords_MAPRECORD().name)
		return false
	}
	var mapFnBody pkg_rhino.Node = this_.priv_TypeTransformationParser_getFunctionBody_7191(mapFn)
	_ = mapFnBody
	if !(this_.priv_TypeTransformationParser_validTypeTransformationExpression_7233(mapFnBody)) {
		this_.priv_TypeTransformationParser_warnInvalidInside_7203("map function body")
		return false
	}
	return true
}

// Java method: TypeTransformationParser.validTypeOfVarExpression(com.google.javascript.rhino.Node expr)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_validTypeOfVarExpression_7228(expr pkg_rhino.Node) bool {
	_ = expr
	if !(this_.priv_TypeTransformationParser_checkParameterCount_7204(expr, StaticGet_TypeTransformationParser_Keywords_TYPEOFVAR())) {
		return false
	}
	if !(this_.priv_TypeTransformationParser_getCallArgument_7193(expr, int32(0)).IsStringLit()) {
		this_.priv_TypeTransformationParser_warnInvalid_7199("name")
		this_.priv_TypeTransformationParser_warnInvalidInside_7203(StaticGet_TypeTransformationParser_Keywords_TYPEOFVAR().name)
		return false
	}
	return true
}

// Java method: TypeTransformationParser.validInstanceOfExpression(com.google.javascript.rhino.Node expr)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_validInstanceOfExpression_7229(expr pkg_rhino.Node) bool {
	_ = expr
	if !(this_.priv_TypeTransformationParser_checkParameterCount_7204(expr, StaticGet_TypeTransformationParser_Keywords_INSTANCEOF())) {
		return false
	}
	if !(this_.priv_TypeTransformationParser_validTypeTransformationExpression_7233(this_.priv_TypeTransformationParser_getCallArgument_7193(expr, int32(0)))) {
		this_.priv_TypeTransformationParser_warnInvalidInside_7203(StaticGet_TypeTransformationParser_Keywords_INSTANCEOF().name)
		return false
	}
	return true
}

// Java method: TypeTransformationParser.validPrintTypeExpression(com.google.javascript.rhino.Node expr)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_validPrintTypeExpression_7230(expr pkg_rhino.Node) bool {
	_ = expr
	if !(this_.priv_TypeTransformationParser_checkParameterCount_7204(expr, StaticGet_TypeTransformationParser_Keywords_PRINTTYPE())) {
		return false
	}
	if !(this_.priv_TypeTransformationParser_getCallArgument_7193(expr, int32(0)).IsStringLit()) {
		this_.priv_TypeTransformationParser_warnInvalid_7199("message")
		this_.priv_TypeTransformationParser_warnInvalidInside_7203(StaticGet_TypeTransformationParser_Keywords_PRINTTYPE().name)
		return false
	}
	if !(this_.priv_TypeTransformationParser_validTypeTransformationExpression_7233(this_.priv_TypeTransformationParser_getCallArgument_7193(expr, int32(1)))) {
		this_.priv_TypeTransformationParser_warnInvalidInside_7203(StaticGet_TypeTransformationParser_Keywords_PRINTTYPE().name)
		return false
	}
	return true
}

// Java method: TypeTransformationParser.validPropTypeExpression(com.google.javascript.rhino.Node expr)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_validPropTypeExpression_7231(expr pkg_rhino.Node) bool {
	_ = expr
	if !(this_.priv_TypeTransformationParser_checkParameterCount_7204(expr, StaticGet_TypeTransformationParser_Keywords_PROPTYPE())) {
		return false
	}
	if !(this_.priv_TypeTransformationParser_getCallArgument_7193(expr, int32(0)).IsStringLit()) {
		this_.priv_TypeTransformationParser_warnInvalid_7199("property name")
		this_.priv_TypeTransformationParser_warnInvalidInside_7203(StaticGet_TypeTransformationParser_Keywords_PROPTYPE().name)
		return false
	}
	if !(this_.priv_TypeTransformationParser_validTypeTransformationExpression_7233(this_.priv_TypeTransformationParser_getCallArgument_7193(expr, int32(1)))) {
		this_.priv_TypeTransformationParser_warnInvalidInside_7203(StaticGet_TypeTransformationParser_Keywords_PROPTYPE().name)
		return false
	}
	return true
}

// Java method: TypeTransformationParser.validOperationExpression(com.google.javascript.rhino.Node expr)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_validOperationExpression_7232(expr pkg_rhino.Node) bool {
	_ = expr
	var name string = this_.priv_TypeTransformationParser_getCallName_7192(expr)
	_ = name
	var keyword *TypeTransformationParser_Keywords = this_.priv_TypeTransformationParser_nameToKeyword_7182(name)
	_ = keyword
	switch j2gort.CheckNotNullPtr(keyword).Ordinal() {
	case int32(1) /* Keywords.COND */ :
		return this_.priv_TypeTransformationParser_validConditionalExpression_7225(expr)
	case int32(9) /* Keywords.MAPUNION */ :
		return this_.priv_TypeTransformationParser_validMapunionExpression_7226(expr)
	case int32(10) /* Keywords.MAPRECORD */ :
		return this_.priv_TypeTransformationParser_validMaprecordExpression_7227(expr)
	case int32(21) /* Keywords.TYPEOFVAR */ :
		return this_.priv_TypeTransformationParser_validTypeOfVarExpression_7228(expr)
	case int32(8) /* Keywords.INSTANCEOF */ :
		return this_.priv_TypeTransformationParser_validInstanceOfExpression_7229(expr)
	case int32(12) /* Keywords.PRINTTYPE */ :
		return this_.priv_TypeTransformationParser_validPrintTypeExpression_7230(expr)
	case int32(13) /* Keywords.PROPTYPE */ :
		return this_.priv_TypeTransformationParser_validPropTypeExpression_7231(expr)
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Invalid type transformation operation")))
	}
	panic("unreachable switch expression")
}

// Java method: TypeTransformationParser.validTypeTransformationExpression(com.google.javascript.rhino.Node expr)
func (this_ *TypeTransformationParser) priv_TypeTransformationParser_validTypeTransformationExpression_7233(expr pkg_rhino.Node) bool {
	_ = expr
	if !(this_.priv_TypeTransformationParser_isValidExpression_7198(expr)) {
		this_.priv_TypeTransformationParser_warnInvalidExpression_7200("type transformation")
		return false
	}
	if this_.priv_TypeTransformationParser_isTypeVar_7195(expr) || this_.priv_TypeTransformationParser_isTypeName_7196(expr) {
		return true
	}
	var name string = this_.priv_TypeTransformationParser_getCallName_7192(expr)
	_ = name
	if !(this_.priv_TypeTransformationParser_isValidKeyword_7183(name)) {
		this_.priv_TypeTransformationParser_warnInvalidExpression_7200("type transformation")
		return false
	}
	var keyword *TypeTransformationParser_Keywords = this_.priv_TypeTransformationParser_nameToKeyword_7182(name)
	_ = keyword
	switch (keyword.kind).Ordinal() {
	case int32(0) /* OperationKind.TYPE_CONSTRUCTOR */ :
		return this_.priv_TypeTransformationParser_validTypeExpression_7217(expr)
	case int32(1) /* OperationKind.OPERATION */ :
		return this_.priv_TypeTransformationParser_validOperationExpression_7232(expr)
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Invalid type transformation expression")))
	}
	panic("unreachable switch expression")
}

// Java enum: com.google.javascript.jscomp.parsing.TypeTransformationParser.OperationKind
type TypeTransformationParser_OperationKind uint8

const (
	enum_TypeTransformationParser_OperationKind_TYPE_CONSTRUCTOR  TypeTransformationParser_OperationKind = 1
	enum_TypeTransformationParser_OperationKind_OPERATION         TypeTransformationParser_OperationKind = 2
	enum_TypeTransformationParser_OperationKind_STRING_PREDICATE  TypeTransformationParser_OperationKind = 3
	enum_TypeTransformationParser_OperationKind_TYPE_PREDICATE    TypeTransformationParser_OperationKind = 4
	enum_TypeTransformationParser_OperationKind_TYPEVAR_PREDICATE TypeTransformationParser_OperationKind = 5
)

var enumBaseTable_TypeTransformationParser_OperationKind = [...]j2gort.EnumBase{
	{Name_: "TYPE_CONSTRUCTOR", Ordinal_: 0},
	{Name_: "OPERATION", Ordinal_: 1},
	{Name_: "STRING_PREDICATE", Ordinal_: 2},
	{Name_: "TYPE_PREDICATE", Ordinal_: 3},
	{Name_: "TYPEVAR_PREDICATE", Ordinal_: 4},
}
var enumValues_TypeTransformationParser_OperationKind = []TypeTransformationParser_OperationKind{1, 2, 3, 4, 5}

func StaticGet_TypeTransformationParser_OperationKind_TYPE_CONSTRUCTOR() TypeTransformationParser_OperationKind {
	return enum_TypeTransformationParser_OperationKind_TYPE_CONSTRUCTOR
}
func StaticGet_TypeTransformationParser_OperationKind_OPERATION() TypeTransformationParser_OperationKind {
	return enum_TypeTransformationParser_OperationKind_OPERATION
}
func StaticGet_TypeTransformationParser_OperationKind_STRING_PREDICATE() TypeTransformationParser_OperationKind {
	return enum_TypeTransformationParser_OperationKind_STRING_PREDICATE
}
func StaticGet_TypeTransformationParser_OperationKind_TYPE_PREDICATE() TypeTransformationParser_OperationKind {
	return enum_TypeTransformationParser_OperationKind_TYPE_PREDICATE
}
func StaticGet_TypeTransformationParser_OperationKind_TYPEVAR_PREDICATE() TypeTransformationParser_OperationKind {
	return enum_TypeTransformationParser_OperationKind_TYPEVAR_PREDICATE
}

func TypeTransformationParser_OperationKind_Values() []TypeTransformationParser_OperationKind {
	return j2gort.CloneSlice(enumValues_TypeTransformationParser_OperationKind)
}

func TypeTransformationParser_OperationKind_ValueOf(name string) TypeTransformationParser_OperationKind {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_TypeTransformationParser_OperationKind {
		if enumBaseTable_TypeTransformationParser_OperationKind[i_].Name_ == name {
			return TypeTransformationParser_OperationKind(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.parsing.TypeTransformationParser.OperationKind." + name))
}

func (this_ TypeTransformationParser_OperationKind) IsJavaNil_() bool                              { return this_ == 0 }
func (TypeTransformationParser_OperationKind) HasIdentityEquals_()                                 {}
func (TypeTransformationParser_OperationKind) Implements_j2gort_Ext_java_lang_Enum()               {}
func (TypeTransformationParser_OperationKind) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ TypeTransformationParser_OperationKind) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (TypeTransformationParser_OperationKind) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.TypeTransformationParser$OperationKind", "OperationKind"
}
func (this_ TypeTransformationParser_OperationKind) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_TypeTransformationParser_OperationKind[this_-1]
}
func (this_ TypeTransformationParser_OperationKind) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ TypeTransformationParser_OperationKind) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ TypeTransformationParser_OperationKind) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_TypeTransformationParser_OperationKind[this_-1].Name_
}
func (this_ TypeTransformationParser_OperationKind) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_TypeTransformationParser_OperationKind[this_-1].Name_
}
func (this_ TypeTransformationParser_OperationKind) Equals(other_ any) bool {
	o_, ok_ := other_.(TypeTransformationParser_OperationKind)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ TypeTransformationParser_OperationKind) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_TypeTransformationParser_OperationKind[this_-1].HashCode()
}
func (this_ TypeTransformationParser_OperationKind) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_TypeTransformationParser_OperationKind[this_-1].CompareTo(other_)
}
func (this_ TypeTransformationParser_OperationKind) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_TypeTransformationParser_OperationKind[this_-1].Name_)
}
func init() {
	vals_ := []any{TypeTransformationParser_OperationKind(1), TypeTransformationParser_OperationKind(2), TypeTransformationParser_OperationKind(3), TypeTransformationParser_OperationKind(4), TypeTransformationParser_OperationKind(5)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.TypeTransformationParser.OperationKind", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.TypeTransformationParser$OperationKind", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.TypeTransformationParser.OperationKind", "com.google.javascript.jscomp.parsing.TypeTransformationParser$OperationKind", func(v_ any) bool {
		return j2gort.InstanceOf[TypeTransformationParser_OperationKind](v_)
	})
}

// Java enum: com.google.javascript.jscomp.parsing.TypeTransformationParser.Keywords
type TypeTransformationParser_Keywords struct {
	j2gort.EnumBase
	name          string
	minParamCount int32
	maxParamCount int32
	kind          TypeTransformationParser_OperationKind
}

func (this_ *TypeTransformationParser_Keywords) GetEnumBase() *j2gort.EnumBase {
	if this_ == nil {
		return nil
	}
	return &this_.EnumBase
}
func (this_ *TypeTransformationParser_Keywords) EnumOrdinal_() int32 {
	if this_ == nil {
		return -1
	}
	return this_.EnumBase.Ordinal_
}
func (*TypeTransformationParser_Keywords) Implements_j2gort_Ext_java_lang_Enum()               {}
func (*TypeTransformationParser_Keywords) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ *TypeTransformationParser_Keywords) FieldGet_kind() TypeTransformationParser_OperationKind {
	return this_.kind
}
func (this_ *TypeTransformationParser_Keywords) initFields_TypeTransformationParser_Keywords() {
	_ = this_
	this_.name = j2gort.NullString
}

func initEnumTypeTransformationParser_Keywords(this_ *TypeTransformationParser_Keywords, enumName_ string, enumOrd_ int32, name string, minParamCount int32, maxParamCount int32, kind TypeTransformationParser_OperationKind) *TypeTransformationParser_Keywords {
	this_.EnumBase.Name_ = enumName_
	this_.EnumBase.Ordinal_ = enumOrd_
	_ = name
	_ = minParamCount
	_ = maxParamCount
	_ = kind
	this_.initFields_TypeTransformationParser_Keywords()
	this_.name = name
	this_.minParamCount = minParamCount
	this_.maxParamCount = maxParamCount
	this_.kind = kind
	return this_
}

var enum_TypeTransformationParser_Keywords_ALL *TypeTransformationParser_Keywords

func StaticGet_TypeTransformationParser_Keywords_ALL() *TypeTransformationParser_Keywords {
	clinit_TypeTransformationParser_Keywords()
	return enum_TypeTransformationParser_Keywords_ALL
}

var enum_TypeTransformationParser_Keywords_COND *TypeTransformationParser_Keywords

func StaticGet_TypeTransformationParser_Keywords_COND() *TypeTransformationParser_Keywords {
	clinit_TypeTransformationParser_Keywords()
	return enum_TypeTransformationParser_Keywords_COND
}

var enum_TypeTransformationParser_Keywords_EQ *TypeTransformationParser_Keywords

func StaticGet_TypeTransformationParser_Keywords_EQ() *TypeTransformationParser_Keywords {
	clinit_TypeTransformationParser_Keywords()
	return enum_TypeTransformationParser_Keywords_EQ
}

var enum_TypeTransformationParser_Keywords_ISCTOR *TypeTransformationParser_Keywords

func StaticGet_TypeTransformationParser_Keywords_ISCTOR() *TypeTransformationParser_Keywords {
	clinit_TypeTransformationParser_Keywords()
	return enum_TypeTransformationParser_Keywords_ISCTOR
}

var enum_TypeTransformationParser_Keywords_ISDEFINED *TypeTransformationParser_Keywords

func StaticGet_TypeTransformationParser_Keywords_ISDEFINED() *TypeTransformationParser_Keywords {
	clinit_TypeTransformationParser_Keywords()
	return enum_TypeTransformationParser_Keywords_ISDEFINED
}

var enum_TypeTransformationParser_Keywords_ISRECORD *TypeTransformationParser_Keywords

func StaticGet_TypeTransformationParser_Keywords_ISRECORD() *TypeTransformationParser_Keywords {
	clinit_TypeTransformationParser_Keywords()
	return enum_TypeTransformationParser_Keywords_ISRECORD
}

var enum_TypeTransformationParser_Keywords_ISTEMPLATIZED *TypeTransformationParser_Keywords

func StaticGet_TypeTransformationParser_Keywords_ISTEMPLATIZED() *TypeTransformationParser_Keywords {
	clinit_TypeTransformationParser_Keywords()
	return enum_TypeTransformationParser_Keywords_ISTEMPLATIZED
}

var enum_TypeTransformationParser_Keywords_ISUNKNOWN *TypeTransformationParser_Keywords

func StaticGet_TypeTransformationParser_Keywords_ISUNKNOWN() *TypeTransformationParser_Keywords {
	clinit_TypeTransformationParser_Keywords()
	return enum_TypeTransformationParser_Keywords_ISUNKNOWN
}

var enum_TypeTransformationParser_Keywords_INSTANCEOF *TypeTransformationParser_Keywords

func StaticGet_TypeTransformationParser_Keywords_INSTANCEOF() *TypeTransformationParser_Keywords {
	clinit_TypeTransformationParser_Keywords()
	return enum_TypeTransformationParser_Keywords_INSTANCEOF
}

var enum_TypeTransformationParser_Keywords_MAPUNION *TypeTransformationParser_Keywords

func StaticGet_TypeTransformationParser_Keywords_MAPUNION() *TypeTransformationParser_Keywords {
	clinit_TypeTransformationParser_Keywords()
	return enum_TypeTransformationParser_Keywords_MAPUNION
}

var enum_TypeTransformationParser_Keywords_MAPRECORD *TypeTransformationParser_Keywords

func StaticGet_TypeTransformationParser_Keywords_MAPRECORD() *TypeTransformationParser_Keywords {
	clinit_TypeTransformationParser_Keywords()
	return enum_TypeTransformationParser_Keywords_MAPRECORD
}

var enum_TypeTransformationParser_Keywords_NONE *TypeTransformationParser_Keywords

func StaticGet_TypeTransformationParser_Keywords_NONE() *TypeTransformationParser_Keywords {
	clinit_TypeTransformationParser_Keywords()
	return enum_TypeTransformationParser_Keywords_NONE
}

var enum_TypeTransformationParser_Keywords_PRINTTYPE *TypeTransformationParser_Keywords

func StaticGet_TypeTransformationParser_Keywords_PRINTTYPE() *TypeTransformationParser_Keywords {
	clinit_TypeTransformationParser_Keywords()
	return enum_TypeTransformationParser_Keywords_PRINTTYPE
}

var enum_TypeTransformationParser_Keywords_PROPTYPE *TypeTransformationParser_Keywords

func StaticGet_TypeTransformationParser_Keywords_PROPTYPE() *TypeTransformationParser_Keywords {
	clinit_TypeTransformationParser_Keywords()
	return enum_TypeTransformationParser_Keywords_PROPTYPE
}

var enum_TypeTransformationParser_Keywords_RAWTYPEOF *TypeTransformationParser_Keywords

func StaticGet_TypeTransformationParser_Keywords_RAWTYPEOF() *TypeTransformationParser_Keywords {
	clinit_TypeTransformationParser_Keywords()
	return enum_TypeTransformationParser_Keywords_RAWTYPEOF
}

var enum_TypeTransformationParser_Keywords_SUB *TypeTransformationParser_Keywords

func StaticGet_TypeTransformationParser_Keywords_SUB() *TypeTransformationParser_Keywords {
	clinit_TypeTransformationParser_Keywords()
	return enum_TypeTransformationParser_Keywords_SUB
}

var enum_TypeTransformationParser_Keywords_STREQ *TypeTransformationParser_Keywords

func StaticGet_TypeTransformationParser_Keywords_STREQ() *TypeTransformationParser_Keywords {
	clinit_TypeTransformationParser_Keywords()
	return enum_TypeTransformationParser_Keywords_STREQ
}

var enum_TypeTransformationParser_Keywords_RECORD *TypeTransformationParser_Keywords

func StaticGet_TypeTransformationParser_Keywords_RECORD() *TypeTransformationParser_Keywords {
	clinit_TypeTransformationParser_Keywords()
	return enum_TypeTransformationParser_Keywords_RECORD
}

var enum_TypeTransformationParser_Keywords_TEMPLATETYPEOF *TypeTransformationParser_Keywords

func StaticGet_TypeTransformationParser_Keywords_TEMPLATETYPEOF() *TypeTransformationParser_Keywords {
	clinit_TypeTransformationParser_Keywords()
	return enum_TypeTransformationParser_Keywords_TEMPLATETYPEOF
}

var enum_TypeTransformationParser_Keywords_TYPE *TypeTransformationParser_Keywords

func StaticGet_TypeTransformationParser_Keywords_TYPE() *TypeTransformationParser_Keywords {
	clinit_TypeTransformationParser_Keywords()
	return enum_TypeTransformationParser_Keywords_TYPE
}

var enum_TypeTransformationParser_Keywords_TYPEEXPR *TypeTransformationParser_Keywords

func StaticGet_TypeTransformationParser_Keywords_TYPEEXPR() *TypeTransformationParser_Keywords {
	clinit_TypeTransformationParser_Keywords()
	return enum_TypeTransformationParser_Keywords_TYPEEXPR
}

var enum_TypeTransformationParser_Keywords_TYPEOFVAR *TypeTransformationParser_Keywords

func StaticGet_TypeTransformationParser_Keywords_TYPEOFVAR() *TypeTransformationParser_Keywords {
	clinit_TypeTransformationParser_Keywords()
	return enum_TypeTransformationParser_Keywords_TYPEOFVAR
}

var enum_TypeTransformationParser_Keywords_UNION *TypeTransformationParser_Keywords

func StaticGet_TypeTransformationParser_Keywords_UNION() *TypeTransformationParser_Keywords {
	clinit_TypeTransformationParser_Keywords()
	return enum_TypeTransformationParser_Keywords_UNION
}

var enum_TypeTransformationParser_Keywords_UNKNOWN *TypeTransformationParser_Keywords

func StaticGet_TypeTransformationParser_Keywords_UNKNOWN() *TypeTransformationParser_Keywords {
	clinit_TypeTransformationParser_Keywords()
	return enum_TypeTransformationParser_Keywords_UNKNOWN
}

var clinitDone_TypeTransformationParser_Keywords uint32
var clinitState_TypeTransformationParser_Keywords j2gort.ClinitState

func clinit_TypeTransformationParser_Keywords() {
	if atomic.LoadUint32(&clinitDone_TypeTransformationParser_Keywords) != 0 {
		return
	}
	clinitSlow_TypeTransformationParser_Keywords()
}

func clinitSlow_TypeTransformationParser_Keywords() {
	if !j2gort.EnterClinit(&clinitState_TypeTransformationParser_Keywords, &clinitDone_TypeTransformationParser_Keywords) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_TypeTransformationParser_Keywords, &clinitDone_TypeTransformationParser_Keywords)
	enum_TypeTransformationParser_Keywords_ALL = &TypeTransformationParser_Keywords{}
	initEnumTypeTransformationParser_Keywords(enum_TypeTransformationParser_Keywords_ALL, "ALL", int32(0), "all", int32(0), int32(0), StaticGet_TypeTransformationParser_OperationKind_TYPE_CONSTRUCTOR())
	enum_TypeTransformationParser_Keywords_COND = &TypeTransformationParser_Keywords{}
	initEnumTypeTransformationParser_Keywords(enum_TypeTransformationParser_Keywords_COND, "COND", int32(1), "cond", int32(3), int32(3), StaticGet_TypeTransformationParser_OperationKind_OPERATION())
	enum_TypeTransformationParser_Keywords_EQ = &TypeTransformationParser_Keywords{}
	initEnumTypeTransformationParser_Keywords(enum_TypeTransformationParser_Keywords_EQ, "EQ", int32(2), "eq", int32(2), int32(2), StaticGet_TypeTransformationParser_OperationKind_TYPE_PREDICATE())
	enum_TypeTransformationParser_Keywords_ISCTOR = &TypeTransformationParser_Keywords{}
	initEnumTypeTransformationParser_Keywords(enum_TypeTransformationParser_Keywords_ISCTOR, "ISCTOR", int32(3), "isCtor", int32(1), int32(1), StaticGet_TypeTransformationParser_OperationKind_TYPE_PREDICATE())
	enum_TypeTransformationParser_Keywords_ISDEFINED = &TypeTransformationParser_Keywords{}
	initEnumTypeTransformationParser_Keywords(enum_TypeTransformationParser_Keywords_ISDEFINED, "ISDEFINED", int32(4), "isDefined", int32(1), int32(1), StaticGet_TypeTransformationParser_OperationKind_TYPEVAR_PREDICATE())
	enum_TypeTransformationParser_Keywords_ISRECORD = &TypeTransformationParser_Keywords{}
	initEnumTypeTransformationParser_Keywords(enum_TypeTransformationParser_Keywords_ISRECORD, "ISRECORD", int32(5), "isRecord", int32(1), int32(1), StaticGet_TypeTransformationParser_OperationKind_TYPE_PREDICATE())
	enum_TypeTransformationParser_Keywords_ISTEMPLATIZED = &TypeTransformationParser_Keywords{}
	initEnumTypeTransformationParser_Keywords(enum_TypeTransformationParser_Keywords_ISTEMPLATIZED, "ISTEMPLATIZED", int32(6), "isTemplatized", int32(1), int32(1), StaticGet_TypeTransformationParser_OperationKind_TYPE_PREDICATE())
	enum_TypeTransformationParser_Keywords_ISUNKNOWN = &TypeTransformationParser_Keywords{}
	initEnumTypeTransformationParser_Keywords(enum_TypeTransformationParser_Keywords_ISUNKNOWN, "ISUNKNOWN", int32(7), "isUnknown", int32(1), int32(1), StaticGet_TypeTransformationParser_OperationKind_TYPE_PREDICATE())
	enum_TypeTransformationParser_Keywords_INSTANCEOF = &TypeTransformationParser_Keywords{}
	initEnumTypeTransformationParser_Keywords(enum_TypeTransformationParser_Keywords_INSTANCEOF, "INSTANCEOF", int32(8), "instanceOf", int32(1), int32(1), StaticGet_TypeTransformationParser_OperationKind_OPERATION())
	enum_TypeTransformationParser_Keywords_MAPUNION = &TypeTransformationParser_Keywords{}
	initEnumTypeTransformationParser_Keywords(enum_TypeTransformationParser_Keywords_MAPUNION, "MAPUNION", int32(9), "mapunion", int32(2), int32(2), StaticGet_TypeTransformationParser_OperationKind_OPERATION())
	enum_TypeTransformationParser_Keywords_MAPRECORD = &TypeTransformationParser_Keywords{}
	initEnumTypeTransformationParser_Keywords(enum_TypeTransformationParser_Keywords_MAPRECORD, "MAPRECORD", int32(10), "maprecord", int32(2), int32(2), StaticGet_TypeTransformationParser_OperationKind_OPERATION())
	enum_TypeTransformationParser_Keywords_NONE = &TypeTransformationParser_Keywords{}
	initEnumTypeTransformationParser_Keywords(enum_TypeTransformationParser_Keywords_NONE, "NONE", int32(11), "none", int32(0), int32(0), StaticGet_TypeTransformationParser_OperationKind_TYPE_CONSTRUCTOR())
	enum_TypeTransformationParser_Keywords_PRINTTYPE = &TypeTransformationParser_Keywords{}
	initEnumTypeTransformationParser_Keywords(enum_TypeTransformationParser_Keywords_PRINTTYPE, "PRINTTYPE", int32(12), "printType", int32(2), int32(2), StaticGet_TypeTransformationParser_OperationKind_OPERATION())
	enum_TypeTransformationParser_Keywords_PROPTYPE = &TypeTransformationParser_Keywords{}
	initEnumTypeTransformationParser_Keywords(enum_TypeTransformationParser_Keywords_PROPTYPE, "PROPTYPE", int32(13), "propType", int32(2), int32(2), StaticGet_TypeTransformationParser_OperationKind_OPERATION())
	enum_TypeTransformationParser_Keywords_RAWTYPEOF = &TypeTransformationParser_Keywords{}
	initEnumTypeTransformationParser_Keywords(enum_TypeTransformationParser_Keywords_RAWTYPEOF, "RAWTYPEOF", int32(14), "rawTypeOf", int32(1), int32(1), StaticGet_TypeTransformationParser_OperationKind_TYPE_CONSTRUCTOR())
	enum_TypeTransformationParser_Keywords_SUB = &TypeTransformationParser_Keywords{}
	initEnumTypeTransformationParser_Keywords(enum_TypeTransformationParser_Keywords_SUB, "SUB", int32(15), "sub", int32(2), int32(2), StaticGet_TypeTransformationParser_OperationKind_TYPE_PREDICATE())
	enum_TypeTransformationParser_Keywords_STREQ = &TypeTransformationParser_Keywords{}
	initEnumTypeTransformationParser_Keywords(enum_TypeTransformationParser_Keywords_STREQ, "STREQ", int32(16), "streq", int32(2), int32(2), StaticGet_TypeTransformationParser_OperationKind_STRING_PREDICATE())
	enum_TypeTransformationParser_Keywords_RECORD = &TypeTransformationParser_Keywords{}
	initEnumTypeTransformationParser_Keywords(enum_TypeTransformationParser_Keywords_RECORD, "RECORD", int32(17), "record", int32(1), int32(2147483647), StaticGet_TypeTransformationParser_OperationKind_TYPE_CONSTRUCTOR())
	enum_TypeTransformationParser_Keywords_TEMPLATETYPEOF = &TypeTransformationParser_Keywords{}
	initEnumTypeTransformationParser_Keywords(enum_TypeTransformationParser_Keywords_TEMPLATETYPEOF, "TEMPLATETYPEOF", int32(18), "templateTypeOf", int32(2), int32(2), StaticGet_TypeTransformationParser_OperationKind_TYPE_CONSTRUCTOR())
	enum_TypeTransformationParser_Keywords_TYPE = &TypeTransformationParser_Keywords{}
	initEnumTypeTransformationParser_Keywords(enum_TypeTransformationParser_Keywords_TYPE, "TYPE", int32(19), "type", int32(2), int32(2147483647), StaticGet_TypeTransformationParser_OperationKind_TYPE_CONSTRUCTOR())
	enum_TypeTransformationParser_Keywords_TYPEEXPR = &TypeTransformationParser_Keywords{}
	initEnumTypeTransformationParser_Keywords(enum_TypeTransformationParser_Keywords_TYPEEXPR, "TYPEEXPR", int32(20), "typeExpr", int32(1), int32(1), StaticGet_TypeTransformationParser_OperationKind_TYPE_CONSTRUCTOR())
	enum_TypeTransformationParser_Keywords_TYPEOFVAR = &TypeTransformationParser_Keywords{}
	initEnumTypeTransformationParser_Keywords(enum_TypeTransformationParser_Keywords_TYPEOFVAR, "TYPEOFVAR", int32(21), "typeOfVar", int32(1), int32(1), StaticGet_TypeTransformationParser_OperationKind_OPERATION())
	enum_TypeTransformationParser_Keywords_UNION = &TypeTransformationParser_Keywords{}
	initEnumTypeTransformationParser_Keywords(enum_TypeTransformationParser_Keywords_UNION, "UNION", int32(22), "union", int32(2), int32(2147483647), StaticGet_TypeTransformationParser_OperationKind_TYPE_CONSTRUCTOR())
	enum_TypeTransformationParser_Keywords_UNKNOWN = &TypeTransformationParser_Keywords{}
	initEnumTypeTransformationParser_Keywords(enum_TypeTransformationParser_Keywords_UNKNOWN, "UNKNOWN", int32(23), "unknown", int32(0), int32(0), StaticGet_TypeTransformationParser_OperationKind_TYPE_CONSTRUCTOR())
	enumValues_TypeTransformationParser_Keywords = []*TypeTransformationParser_Keywords{enum_TypeTransformationParser_Keywords_ALL, enum_TypeTransformationParser_Keywords_COND, enum_TypeTransformationParser_Keywords_EQ, enum_TypeTransformationParser_Keywords_ISCTOR, enum_TypeTransformationParser_Keywords_ISDEFINED, enum_TypeTransformationParser_Keywords_ISRECORD, enum_TypeTransformationParser_Keywords_ISTEMPLATIZED, enum_TypeTransformationParser_Keywords_ISUNKNOWN, enum_TypeTransformationParser_Keywords_INSTANCEOF, enum_TypeTransformationParser_Keywords_MAPUNION, enum_TypeTransformationParser_Keywords_MAPRECORD, enum_TypeTransformationParser_Keywords_NONE, enum_TypeTransformationParser_Keywords_PRINTTYPE, enum_TypeTransformationParser_Keywords_PROPTYPE, enum_TypeTransformationParser_Keywords_RAWTYPEOF, enum_TypeTransformationParser_Keywords_SUB, enum_TypeTransformationParser_Keywords_STREQ, enum_TypeTransformationParser_Keywords_RECORD, enum_TypeTransformationParser_Keywords_TEMPLATETYPEOF, enum_TypeTransformationParser_Keywords_TYPE, enum_TypeTransformationParser_Keywords_TYPEEXPR, enum_TypeTransformationParser_Keywords_TYPEOFVAR, enum_TypeTransformationParser_Keywords_UNION, enum_TypeTransformationParser_Keywords_UNKNOWN}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.TypeTransformationParser.Keywords", []any{enum_TypeTransformationParser_Keywords_ALL, enum_TypeTransformationParser_Keywords_COND, enum_TypeTransformationParser_Keywords_EQ, enum_TypeTransformationParser_Keywords_ISCTOR, enum_TypeTransformationParser_Keywords_ISDEFINED, enum_TypeTransformationParser_Keywords_ISRECORD, enum_TypeTransformationParser_Keywords_ISTEMPLATIZED, enum_TypeTransformationParser_Keywords_ISUNKNOWN, enum_TypeTransformationParser_Keywords_INSTANCEOF, enum_TypeTransformationParser_Keywords_MAPUNION, enum_TypeTransformationParser_Keywords_MAPRECORD, enum_TypeTransformationParser_Keywords_NONE, enum_TypeTransformationParser_Keywords_PRINTTYPE, enum_TypeTransformationParser_Keywords_PROPTYPE, enum_TypeTransformationParser_Keywords_RAWTYPEOF, enum_TypeTransformationParser_Keywords_SUB, enum_TypeTransformationParser_Keywords_STREQ, enum_TypeTransformationParser_Keywords_RECORD, enum_TypeTransformationParser_Keywords_TEMPLATETYPEOF, enum_TypeTransformationParser_Keywords_TYPE, enum_TypeTransformationParser_Keywords_TYPEEXPR, enum_TypeTransformationParser_Keywords_TYPEOFVAR, enum_TypeTransformationParser_Keywords_UNION, enum_TypeTransformationParser_Keywords_UNKNOWN})
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.TypeTransformationParser$Keywords", []any{enum_TypeTransformationParser_Keywords_ALL, enum_TypeTransformationParser_Keywords_COND, enum_TypeTransformationParser_Keywords_EQ, enum_TypeTransformationParser_Keywords_ISCTOR, enum_TypeTransformationParser_Keywords_ISDEFINED, enum_TypeTransformationParser_Keywords_ISRECORD, enum_TypeTransformationParser_Keywords_ISTEMPLATIZED, enum_TypeTransformationParser_Keywords_ISUNKNOWN, enum_TypeTransformationParser_Keywords_INSTANCEOF, enum_TypeTransformationParser_Keywords_MAPUNION, enum_TypeTransformationParser_Keywords_MAPRECORD, enum_TypeTransformationParser_Keywords_NONE, enum_TypeTransformationParser_Keywords_PRINTTYPE, enum_TypeTransformationParser_Keywords_PROPTYPE, enum_TypeTransformationParser_Keywords_RAWTYPEOF, enum_TypeTransformationParser_Keywords_SUB, enum_TypeTransformationParser_Keywords_STREQ, enum_TypeTransformationParser_Keywords_RECORD, enum_TypeTransformationParser_Keywords_TEMPLATETYPEOF, enum_TypeTransformationParser_Keywords_TYPE, enum_TypeTransformationParser_Keywords_TYPEEXPR, enum_TypeTransformationParser_Keywords_TYPEOFVAR, enum_TypeTransformationParser_Keywords_UNION, enum_TypeTransformationParser_Keywords_UNKNOWN})
}

func init() {
	j2gort.RegisterEnumConstantsLoader("com.google.javascript.jscomp.parsing.TypeTransformationParser.Keywords", "com.google.javascript.jscomp.parsing.TypeTransformationParser$Keywords", func() {
		clinit_TypeTransformationParser_Keywords()
	})
}

var enumValues_TypeTransformationParser_Keywords []*TypeTransformationParser_Keywords

func TypeTransformationParser_Keywords_Values() []*TypeTransformationParser_Keywords {
	clinit_TypeTransformationParser_Keywords()
	return j2gort.CloneSlice(enumValues_TypeTransformationParser_Keywords)
}

func TypeTransformationParser_Keywords_ValueOf(name string) *TypeTransformationParser_Keywords {
	clinit_TypeTransformationParser_Keywords()
	for _, c := range enumValues_TypeTransformationParser_Keywords {
		if c.Name() == name {
			return c
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.parsing.TypeTransformationParser.Keywords." + name))
}

func (this_ *TypeTransformationParser_Keywords) IsJavaNil_() bool { return this_ == nil }
func (*TypeTransformationParser_Keywords) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.TypeTransformationParser$Keywords", "Keywords"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.TypeTransformationParser.Keywords", "com.google.javascript.jscomp.parsing.TypeTransformationParser$Keywords", func(v_ any) bool {
		return j2gort.InstanceOf[*TypeTransformationParser_Keywords](v_)
	})
}
