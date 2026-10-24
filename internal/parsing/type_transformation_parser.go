// Port of com/google/javascript/jscomp/parsing/TypeTransformationParser.java @ closure-compiler v20261021

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

package parsing

import (
	"math"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// TypeTransformationParser is a parser for the type transformation
// expressions (TTL-Exp) as in `@template T := TTL-Exp =:`.
//
// Java: final class TypeTransformationParser
type TypeTransformationParser struct {
	typeTransformationString string
	typeTransformationAst    rhino.Node
	sourceFile               rhino.StaticSourceFile
	errorReporter            rhino.ErrorReporter
	templateLineno           int
	templateCharno           int
}

const TypeTransformationParser_VAR_ARGS = math.MaxInt32

// Set the length of every TTL node as the length of "@template" (which is
// 9). This is used for error message logging where we underline the error
// location with ^^^^^^^^^.
var TypeTransformationParser_TTL_NODE_LENGTH = jutil.StringLength("@template")

// TypeTransformationParser_OperationKind is the classification of the
// keywords.
//
// Java: enum TypeTransformationParser.OperationKind
type TypeTransformationParser_OperationKind int32

const (
	TypeTransformationParser_OperationKind_TYPE_CONSTRUCTOR TypeTransformationParser_OperationKind = iota
	TypeTransformationParser_OperationKind_OPERATION
	TypeTransformationParser_OperationKind_STRING_PREDICATE
	TypeTransformationParser_OperationKind_TYPE_PREDICATE
	TypeTransformationParser_OperationKind_TYPEVAR_PREDICATE
)

var typeTransformationParser_OperationKind_names = [...]string{
	"TYPE_CONSTRUCTOR", "OPERATION", "STRING_PREDICATE", "TYPE_PREDICATE", "TYPEVAR_PREDICATE",
}

// Java: TypeTransformationParser.OperationKind.values()
func TypeTransformationParser_OperationKind_values() []TypeTransformationParser_OperationKind {
	return []TypeTransformationParser_OperationKind{
		TypeTransformationParser_OperationKind_TYPE_CONSTRUCTOR,
		TypeTransformationParser_OperationKind_OPERATION,
		TypeTransformationParser_OperationKind_STRING_PREDICATE,
		TypeTransformationParser_OperationKind_TYPE_PREDICATE,
		TypeTransformationParser_OperationKind_TYPEVAR_PREDICATE,
	}
}

// String is Enum.name().
func (k TypeTransformationParser_OperationKind) String() string {
	return typeTransformationParser_OperationKind_names[k]
}

// TypeTransformationParser_Keywords are the keywords of the type
// transformation language.
//
// Java: enum TypeTransformationParser.Keywords
type TypeTransformationParser_Keywords int32

const (
	TypeTransformationParser_Keywords_ALL TypeTransformationParser_Keywords = iota
	TypeTransformationParser_Keywords_COND
	TypeTransformationParser_Keywords_EQ
	TypeTransformationParser_Keywords_ISCTOR
	TypeTransformationParser_Keywords_ISDEFINED
	TypeTransformationParser_Keywords_ISRECORD
	TypeTransformationParser_Keywords_ISTEMPLATIZED
	TypeTransformationParser_Keywords_ISUNKNOWN
	TypeTransformationParser_Keywords_INSTANCEOF
	TypeTransformationParser_Keywords_MAPUNION
	TypeTransformationParser_Keywords_MAPRECORD
	TypeTransformationParser_Keywords_NONE
	TypeTransformationParser_Keywords_PRINTTYPE
	TypeTransformationParser_Keywords_PROPTYPE
	TypeTransformationParser_Keywords_RAWTYPEOF
	TypeTransformationParser_Keywords_SUB
	TypeTransformationParser_Keywords_STREQ
	TypeTransformationParser_Keywords_RECORD
	TypeTransformationParser_Keywords_TEMPLATETYPEOF
	TypeTransformationParser_Keywords_TYPE
	TypeTransformationParser_Keywords_TYPEEXPR
	TypeTransformationParser_Keywords_TYPEOFVAR
	TypeTransformationParser_Keywords_UNION
	TypeTransformationParser_Keywords_UNKNOWN
)

var typeTransformationParser_Keywords_constantNames = [...]string{
	"ALL", "COND", "EQ", "ISCTOR", "ISDEFINED", "ISRECORD", "ISTEMPLATIZED", "ISUNKNOWN", "INSTANCEOF",
	"MAPUNION", "MAPRECORD", "NONE", "PRINTTYPE", "PROPTYPE", "RAWTYPEOF", "SUB", "STREQ", "RECORD",
	"TEMPLATETYPEOF", "TYPE", "TYPEEXPR", "TYPEOFVAR", "UNION", "UNKNOWN",
}

// The fields of each Keywords constant, indexed by ordinal:
// Keywords(String name, int minParamCount, int maxParamCount, OperationKind kind).
var typeTransformationParser_Keywords_fields = [...]struct {
	name          string
	minParamCount int
	maxParamCount int
	kind          TypeTransformationParser_OperationKind
}{
	{"all", 0, 0, TypeTransformationParser_OperationKind_TYPE_CONSTRUCTOR},
	{"cond", 3, 3, TypeTransformationParser_OperationKind_OPERATION},
	{"eq", 2, 2, TypeTransformationParser_OperationKind_TYPE_PREDICATE},
	{"isCtor", 1, 1, TypeTransformationParser_OperationKind_TYPE_PREDICATE},
	{"isDefined", 1, 1, TypeTransformationParser_OperationKind_TYPEVAR_PREDICATE},
	{"isRecord", 1, 1, TypeTransformationParser_OperationKind_TYPE_PREDICATE},
	{"isTemplatized", 1, 1, TypeTransformationParser_OperationKind_TYPE_PREDICATE},
	{"isUnknown", 1, 1, TypeTransformationParser_OperationKind_TYPE_PREDICATE},
	{"instanceOf", 1, 1, TypeTransformationParser_OperationKind_OPERATION},
	{"mapunion", 2, 2, TypeTransformationParser_OperationKind_OPERATION},
	{"maprecord", 2, 2, TypeTransformationParser_OperationKind_OPERATION},
	{"none", 0, 0, TypeTransformationParser_OperationKind_TYPE_CONSTRUCTOR},
	{"printType", 2, 2, TypeTransformationParser_OperationKind_OPERATION},
	{"propType", 2, 2, TypeTransformationParser_OperationKind_OPERATION},
	{"rawTypeOf", 1, 1, TypeTransformationParser_OperationKind_TYPE_CONSTRUCTOR},
	{"sub", 2, 2, TypeTransformationParser_OperationKind_TYPE_PREDICATE},
	{"streq", 2, 2, TypeTransformationParser_OperationKind_STRING_PREDICATE},
	{"record", 1, TypeTransformationParser_VAR_ARGS, TypeTransformationParser_OperationKind_TYPE_CONSTRUCTOR},
	{"templateTypeOf", 2, 2, TypeTransformationParser_OperationKind_TYPE_CONSTRUCTOR},
	{"type", 2, TypeTransformationParser_VAR_ARGS, TypeTransformationParser_OperationKind_TYPE_CONSTRUCTOR},
	{"typeExpr", 1, 1, TypeTransformationParser_OperationKind_TYPE_CONSTRUCTOR},
	{"typeOfVar", 1, 1, TypeTransformationParser_OperationKind_OPERATION},
	{"union", 2, TypeTransformationParser_VAR_ARGS, TypeTransformationParser_OperationKind_TYPE_CONSTRUCTOR},
	{"unknown", 0, 0, TypeTransformationParser_OperationKind_TYPE_CONSTRUCTOR},
}

// Java: TypeTransformationParser.Keywords.values()
func TypeTransformationParser_Keywords_values() []TypeTransformationParser_Keywords {
	values := make([]TypeTransformationParser_Keywords, len(typeTransformationParser_Keywords_constantNames))
	for i := range values {
		values[i] = TypeTransformationParser_Keywords(i)
	}
	return values
}

// TypeTransformationParser_Keywords_ValueOf is Keywords.valueOf(String).
func TypeTransformationParser_Keywords_ValueOf(name string) TypeTransformationParser_Keywords {
	for i, n := range typeTransformationParser_Keywords_constantNames {
		if n == name {
			return TypeTransformationParser_Keywords(i)
		}
	}
	panic(jutil.NewIllegalArgumentExceptionMsg(
		"No enum constant com.google.javascript.jscomp.parsing.TypeTransformationParser.Keywords." + name))
}

func asciiEqualUpper(s string, upperConst string) bool {
	if len(s) != len(upperConst) {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		if c != upperConst[i] {
			return false
		}
	}
	return true
}

// TypeTransformationParser_Keywords_FromCaseInsensitiveASCII resolves
// Keywords.valueOf(Ascii.toUpperCase(name)) without allocating an uppercase copy.
func TypeTransformationParser_Keywords_FromCaseInsensitiveASCII(name string) TypeTransformationParser_Keywords {
	switch len(name) {
	case 2:
		if asciiEqualUpper(name, "EQ") {
			return TypeTransformationParser_Keywords_EQ
		}
	case 3:
		if asciiEqualUpper(name, "ALL") {
			return TypeTransformationParser_Keywords_ALL
		}
		if asciiEqualUpper(name, "SUB") {
			return TypeTransformationParser_Keywords_SUB
		}
	case 4:
		if asciiEqualUpper(name, "TYPE") {
			return TypeTransformationParser_Keywords_TYPE
		}
		if asciiEqualUpper(name, "COND") {
			return TypeTransformationParser_Keywords_COND
		}
		if asciiEqualUpper(name, "NONE") {
			return TypeTransformationParser_Keywords_NONE
		}
	case 5:
		if asciiEqualUpper(name, "UNION") {
			return TypeTransformationParser_Keywords_UNION
		}
		if asciiEqualUpper(name, "STREQ") {
			return TypeTransformationParser_Keywords_STREQ
		}
	case 6:
		if asciiEqualUpper(name, "RECORD") {
			return TypeTransformationParser_Keywords_RECORD
		}
		if asciiEqualUpper(name, "ISCTOR") {
			return TypeTransformationParser_Keywords_ISCTOR
		}
	case 7:
		if asciiEqualUpper(name, "UNKNOWN") {
			return TypeTransformationParser_Keywords_UNKNOWN
		}
	case 8:
		if asciiEqualUpper(name, "MAPUNION") {
			return TypeTransformationParser_Keywords_MAPUNION
		}
		if asciiEqualUpper(name, "ISRECORD") {
			return TypeTransformationParser_Keywords_ISRECORD
		}
		if asciiEqualUpper(name, "PROPTYPE") {
			return TypeTransformationParser_Keywords_PROPTYPE
		}
		if asciiEqualUpper(name, "TYPEEXPR") {
			return TypeTransformationParser_Keywords_TYPEEXPR
		}
	case 9:
		if asciiEqualUpper(name, "ISDEFINED") {
			return TypeTransformationParser_Keywords_ISDEFINED
		}
		if asciiEqualUpper(name, "ISUNKNOWN") {
			return TypeTransformationParser_Keywords_ISUNKNOWN
		}
		if asciiEqualUpper(name, "MAPRECORD") {
			return TypeTransformationParser_Keywords_MAPRECORD
		}
		if asciiEqualUpper(name, "PRINTTYPE") {
			return TypeTransformationParser_Keywords_PRINTTYPE
		}
		if asciiEqualUpper(name, "RAWTYPEOF") {
			return TypeTransformationParser_Keywords_RAWTYPEOF
		}
		if asciiEqualUpper(name, "TYPEOFVAR") {
			return TypeTransformationParser_Keywords_TYPEOFVAR
		}
	case 10:
		if asciiEqualUpper(name, "INSTANCEOF") {
			return TypeTransformationParser_Keywords_INSTANCEOF
		}
	case 13:
		if asciiEqualUpper(name, "ISTEMPLATIZED") {
			return TypeTransformationParser_Keywords_ISTEMPLATIZED
		}
	case 14:
		if asciiEqualUpper(name, "TEMPLATETYPEOF") {
			return TypeTransformationParser_Keywords_TEMPLATETYPEOF
		}
	}
	return TypeTransformationParser_Keywords_ValueOf(jutil.StringToUpperCaseASCII(name))
}

// String is Enum.name().
func (k TypeTransformationParser_Keywords) String() string {
	return typeTransformationParser_Keywords_constantNames[k]
}

// KeywordName is the field Keywords.name (the keyword as written, e.g.
// "rawTypeOf").
//
// Go-only: language-adapter — the field is named `name` in Java, which Go
// cannot tell apart from Enum.name() (String).
func (k TypeTransformationParser_Keywords) KeywordName() string {
	return typeTransformationParser_Keywords_fields[k].name
}

// MinParamCount is the field Keywords.minParamCount.
func (k TypeTransformationParser_Keywords) MinParamCount() int {
	return typeTransformationParser_Keywords_fields[k].minParamCount
}

// MaxParamCount is the field Keywords.maxParamCount.
func (k TypeTransformationParser_Keywords) MaxParamCount() int {
	return typeTransformationParser_Keywords_fields[k].maxParamCount
}

// Kind is the field Keywords.kind.
func (k TypeTransformationParser_Keywords) Kind() TypeTransformationParser_OperationKind {
	return typeTransformationParser_Keywords_fields[k].kind
}

// Java: TypeTransformationParser(String, StaticSourceFile, ErrorReporter, int, int)
func NewTypeTransformationParser(typeTransformationString string,
	sourceFile rhino.StaticSourceFile, errorReporter rhino.ErrorReporter,
	templateLineno int, templateCharno int) *TypeTransformationParser {
	p := &TypeTransformationParser{}
	p.typeTransformationString = typeTransformationString
	p.sourceFile = sourceFile
	p.errorReporter = errorReporter
	p.templateLineno = templateLineno
	p.templateCharno = templateCharno
	return p
}

// Java: TypeTransformationParser.getTypeTransformationAst()
func (p *TypeTransformationParser) GetTypeTransformationAst() rhino.Node {
	return p.typeTransformationAst
}

// Java: TypeTransformationParser.addNewWarning(Msg, String)
func (p *TypeTransformationParser) AddNewWarning(messageId rhino.Msg, messageArg string) {
	// TODO(lpino): Use the exact lineno and charno, it is currently using
	// the lineno and charno of the parent @template
	// TODO(lpino): Use only constants as parameters of this method
	p.errorReporter.Warning(
		"Bad type annotation. "+messageId.FormatArgs(messageArg),
		p.sourceFile.GetName(),
		p.templateLineno,
		p.templateCharno)
}

// Java: TypeTransformationParser.nameToKeyword(String)
func (p *TypeTransformationParser) NameToKeyword(s string) TypeTransformationParser_Keywords {
	return TypeTransformationParser_Keywords_ValueOf(jutil.StringToUpperCaseASCII(s))
}

// Java: TypeTransformationParser.isValidKeyword(String)
func (p *TypeTransformationParser) IsValidKeyword(name string) bool {
	for _, k := range TypeTransformationParser_Keywords_values() {
		if k.KeywordName() == name {
			return true
		}
	}
	return false
}

// Java: TypeTransformationParser.isOperationKind(String, OperationKind)
func (p *TypeTransformationParser) IsOperationKind(name string, kind TypeTransformationParser_OperationKind) bool {
	return p.IsValidKeyword(name) && p.NameToKeyword(name).Kind() == kind
}

// Java: TypeTransformationParser.isValidStringPredicate(String)
func (p *TypeTransformationParser) IsValidStringPredicate(name string) bool {
	return p.IsOperationKind(name, TypeTransformationParser_OperationKind_STRING_PREDICATE)
}

// Java: TypeTransformationParser.isValidTypePredicate(String)
func (p *TypeTransformationParser) IsValidTypePredicate(name string) bool {
	return p.IsOperationKind(name, TypeTransformationParser_OperationKind_TYPE_PREDICATE)
}

// Java: TypeTransformationParser.isValidTypevarPredicate(String)
func (p *TypeTransformationParser) IsValidTypevarPredicate(name string) bool {
	return p.IsOperationKind(name, TypeTransformationParser_OperationKind_TYPEVAR_PREDICATE)
}

// Java: TypeTransformationParser.isBooleanOperation(Node)
func (p *TypeTransformationParser) IsBooleanOperation(n rhino.Node) bool {
	return n.IsAnd() || n.IsOr() || n.IsNot()
}

// Java: TypeTransformationParser.isValidPredicate(String)
func (p *TypeTransformationParser) IsValidPredicate(name string) bool {
	return p.IsValidStringPredicate(name) ||
		p.IsValidTypePredicate(name) ||
		p.IsValidTypevarPredicate(name)
}

// Java: TypeTransformationParser.getFunctionParamCount(Node)
func (p *TypeTransformationParser) GetFunctionParamCount(n rhino.Node) int {
	jutil.CheckArgumentMsg1(n.IsFunction(),
		"Expected a function node, found %s", n)
	return n.GetSecondChild().GetChildCount()
}

// Java: TypeTransformationParser.getFunctionBody(Node)
func (p *TypeTransformationParser) GetFunctionBody(n rhino.Node) rhino.Node {
	jutil.CheckArgumentMsg1(n.IsFunction(),
		"Expected a function node, found %s", n)
	return n.GetChildAtIndex(2)
}

// Java: TypeTransformationParser.getCallName(Node)
func (p *TypeTransformationParser) GetCallName(n rhino.Node) string {
	jutil.CheckArgumentMsg1(n.IsCall(),
		"Expected a call node, found %s", n)
	return n.GetFirstChild().GetString()
}

// Java: TypeTransformationParser.getCallArgument(Node, int)
func (p *TypeTransformationParser) GetCallArgument(n rhino.Node, i int) rhino.Node {
	jutil.CheckArgumentMsg1(n.IsCall(),
		"Expected a call node, found %s", n)
	return n.GetChildAtIndex(i + 1)
}

// Java: TypeTransformationParser.getCallParamCount(Node)
func (p *TypeTransformationParser) GetCallParamCount(n rhino.Node) int {
	jutil.CheckArgumentMsg1(n.IsCall(),
		"Expected a call node, found %s", n)
	return n.GetChildCount() - 1
}

// Java: TypeTransformationParser.isTypeVar(Node)
func (p *TypeTransformationParser) IsTypeVar(n rhino.Node) bool {
	return n.IsName()
}

// Java: TypeTransformationParser.isTypeName(Node)
func (p *TypeTransformationParser) IsTypeName(n rhino.Node) bool {
	return n.IsStringLit()
}

// Java: TypeTransformationParser.isOperation(Node)
func (p *TypeTransformationParser) IsOperation(n rhino.Node) bool {
	return n.IsCall()
}

// IsValidExpression: a valid expression is either:
//   - NAME for a type variable
//   - STRING for a type name
//   - CALL for the other expressions
//
// Java: TypeTransformationParser.isValidExpression(Node)
func (p *TypeTransformationParser) IsValidExpression(e rhino.Node) bool {
	return p.IsTypeVar(e) || p.IsTypeName(e) || p.IsOperation(e)
}

// Java: TypeTransformationParser.warnInvalid(String)
func (p *TypeTransformationParser) WarnInvalid(msg string) {
	p.AddNewWarning(rhino.Msg_JSDOC_TYPETRANSFORMATION_INVALID, msg)
}

// Java: TypeTransformationParser.warnInvalidExpression(String)
func (p *TypeTransformationParser) WarnInvalidExpression(msg string) {
	p.AddNewWarning(rhino.Msg_JSDOC_TYPETRANSFORMATION_INVALID_EXPRESSION, msg)
}

// Java: TypeTransformationParser.warnMissingParam(String)
func (p *TypeTransformationParser) WarnMissingParam(msg string) {
	p.AddNewWarning(rhino.Msg_JSDOC_TYPETRANSFORMATION_MISSING_PARAM, msg)
}

// Java: TypeTransformationParser.warnExtraParam(String)
func (p *TypeTransformationParser) WarnExtraParam(msg string) {
	p.AddNewWarning(rhino.Msg_JSDOC_TYPETRANSFORMATION_EXTRA_PARAM, msg)
}

// Java: TypeTransformationParser.warnInvalidInside(String)
func (p *TypeTransformationParser) WarnInvalidInside(msg string) {
	p.AddNewWarning(rhino.Msg_JSDOC_TYPETRANSFORMATION_INVALID_INSIDE, msg)
}

// Java: TypeTransformationParser.checkParameterCount(Node, Keywords)
func (p *TypeTransformationParser) CheckParameterCount(expr rhino.Node, keyword TypeTransformationParser_Keywords) bool {
	paramCount := p.GetCallParamCount(expr)
	if paramCount < keyword.MinParamCount() {
		p.WarnMissingParam(keyword.KeywordName())
		return false
	}
	if paramCount > keyword.MaxParamCount() {
		p.WarnExtraParam(keyword.KeywordName())
		return false
	}
	return true
}

// ParseTypeTransformation takes a type transformation expression, transforms
// it to an AST using the ParserRunner of the JSCompiler and then verifies
// that it is a valid AST. Returns true if the parsing was successful
// otherwise it returns false and at least one warning is reported.
//
// Java: TypeTransformationParser.parseTypeTransformation()
func (p *TypeTransformationParser) ParseTypeTransformation() bool {
	config := Config_Builder_().
		SetLanguageMode(Config_LanguageMode_ES_NEXT).
		SetStrictMode(Config_StrictMode_SLOPPY).
		Build()
	// TODO(lpino): ParserRunner reports errors if the expression is not
	// ES6 valid. We need to abort the validation of the type transformation
	// whenever an error is reported.
	result := ParserRunner_Parse(
		p.sourceFile, p.typeTransformationString, config, p.errorReporter)
	ast := result.Ast
	// Check that the expression is a script with an expression result
	if ast.IsNil() ||
		!ast.IsScript() ||
		!ast.HasChildren() ||
		!ast.GetFirstChild().IsExprResult() {
		p.WarnInvalidExpression("type transformation")
		return false
	}

	expr := ast.GetFirstFirstChild()
	// The AST of the type transformation must correspond to a valid
	// expression
	if !p.ValidTypeTransformationExpression(expr) {
		// No need to add a new warning because the validation does it
		return false
	}
	p.FixTTLNodeLineNoCharNoAndLength(expr)
	// Store the result if the AST is valid
	p.typeTransformationAst = expr
	return true
}

// FixTTLNodeLineNoCharNoAndLength sets the lineno/charno of the TTL node to
// the lineno/charno of the "@" in the JSDoc's template annotation. This
// lineno/charno of the "@" is initialized by JsDocInfoParser when parsing
// TTLs.
//
// The TTL AST currently has lineno/charno relative to its position within
// the JSDoc, and `fixTTLNodeLineNoCharNoAndLength` changes this to be
// relative to its position in the source JS file.
//
// We are using the lineno/charno of "@" because adjusting the TTL AST node's
// lineno/charno to be more accurate is tricky, especially for handling
// multiline TTL expressions. If we set an inaccurate lineno/charno, the
// compiler may crash with an "index out of bounds" error when logging errors
// that occur in the TTL expression. Instead, we will use the lineno/charno
// of "@" to safely highlight which JSDoc annotation an error occurs in.
//
// Java: TypeTransformationParser.fixTTLNodeLineNoCharNoAndLength(Node)
func (p *TypeTransformationParser) FixTTLNodeLineNoCharNoAndLength(expr rhino.Node) {
	expr.SetLinenoCharno(p.templateLineno, p.templateCharno)
	expr.SetLength(TypeTransformationParser_TTL_NODE_LENGTH)
	for child := expr.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		p.FixTTLNodeLineNoCharNoAndLength(child)
	}
}

// ValidTemplateTypeExpression: a template type expression must be of the
// form type(typename, TTLExp,...) or type(typevar, TTLExp...).
//
// Java: TypeTransformationParser.validTemplateTypeExpression(Node)
func (p *TypeTransformationParser) ValidTemplateTypeExpression(expr rhino.Node) bool {
	// The expression must have at least three children the type keyword,
	// a type name (or type variable) and a type expression
	if !p.CheckParameterCount(expr, TypeTransformationParser_Keywords_TYPE) {
		return false
	}
	paramCount := p.GetCallParamCount(expr)
	// The first parameter must be a type variable or a type name
	firstParam := p.GetCallArgument(expr, 0)
	if !p.IsTypeVar(firstParam) && !p.IsTypeName(firstParam) {
		p.WarnInvalid("type name or type variable")
		p.WarnInvalidInside("template type operation")
		return false
	}
	// The rest of the parameters must be valid type expressions
	for i := 1; i < paramCount; i++ {
		if !p.ValidTypeTransformationExpression(p.GetCallArgument(expr, i)) {
			p.WarnInvalidInside("template type operation")
			return false
		}
	}
	return true
}

// ValidUnionTypeExpression: a Union type expression must be a valid type
// variable or a union(TTLExp, TTLExp, ...).
//
// Java: TypeTransformationParser.validUnionTypeExpression(Node)
func (p *TypeTransformationParser) ValidUnionTypeExpression(expr rhino.Node) bool {
	// The expression must have at least three children: The union keyword
	// and two type expressions
	if !p.CheckParameterCount(expr, TypeTransformationParser_Keywords_UNION) {
		return false
	}
	paramCount := p.GetCallParamCount(expr)
	// Check if each of the members of the union is a valid type expression
	for i := 0; i < paramCount; i++ {
		if !p.ValidTypeTransformationExpression(p.GetCallArgument(expr, i)) {
			p.WarnInvalidInside("union type")
			return false
		}
	}
	return true
}

// ValidNoneTypeExpression: a none type expression must be of the form:
// none().
//
// Java: TypeTransformationParser.validNoneTypeExpression(Node)
func (p *TypeTransformationParser) ValidNoneTypeExpression(expr rhino.Node) bool {
	// The expression must have no children
	return p.CheckParameterCount(expr, TypeTransformationParser_Keywords_NONE)
}

// ValidAllTypeExpression: an all type expression must be of the form: all().
//
// Java: TypeTransformationParser.validAllTypeExpression(Node)
func (p *TypeTransformationParser) ValidAllTypeExpression(expr rhino.Node) bool {
	// The expression must have no children
	return p.CheckParameterCount(expr, TypeTransformationParser_Keywords_ALL)
}

// ValidUnknownTypeExpression: an unknown type expression must be of the
// form: unknown().
//
// Java: TypeTransformationParser.validUnknownTypeExpression(Node)
func (p *TypeTransformationParser) ValidUnknownTypeExpression(expr rhino.Node) bool {
	// The expression must have no children
	return p.CheckParameterCount(expr, TypeTransformationParser_Keywords_UNKNOWN)
}

// ValidRawTypeOfTypeExpression: a raw type expression must be of the form
// rawTypeOf(TTLExp).
//
// Java: TypeTransformationParser.validRawTypeOfTypeExpression(Node)
func (p *TypeTransformationParser) ValidRawTypeOfTypeExpression(expr rhino.Node) bool {
	// The expression must have two children. The rawTypeOf keyword and the
	// parameter
	if !p.CheckParameterCount(expr, TypeTransformationParser_Keywords_RAWTYPEOF) {
		return false
	}
	// The parameter must be a valid type expression
	if !p.ValidTypeTransformationExpression(p.GetCallArgument(expr, 0)) {
		p.WarnInvalidInside(TypeTransformationParser_Keywords_RAWTYPEOF.KeywordName())
		return false
	}
	return true
}

// ValidTemplateTypeOfExpression: a template type of expression must be of
// the form templateTypeOf(TTLExp, index).
//
// Java: TypeTransformationParser.validTemplateTypeOfExpression(Node)
func (p *TypeTransformationParser) ValidTemplateTypeOfExpression(expr rhino.Node) bool {
	// The expression must have three children. The templateTypeOf keyword, a
	// templatized type and an index
	if !p.CheckParameterCount(expr, TypeTransformationParser_Keywords_TEMPLATETYPEOF) {
		return false
	}
	// The parameter must be a valid type expression
	if !p.ValidTypeTransformationExpression(p.GetCallArgument(expr, 0)) {
		p.WarnInvalidInside(TypeTransformationParser_Keywords_TEMPLATETYPEOF.KeywordName())
		return false
	}
	if !p.GetCallArgument(expr, 1).IsNumber() {
		p.WarnInvalid("index")
		p.WarnInvalidInside(TypeTransformationParser_Keywords_TEMPLATETYPEOF.KeywordName())
		return false
	}
	index := p.GetCallArgument(expr, 1).GetDouble()
	if math.IsNaN(index) || !base.JSCompDoubles_IsExactInt32(index) {
		p.WarnInvalid("index")
		p.WarnInvalidInside(TypeTransformationParser_Keywords_TEMPLATETYPEOF.KeywordName())
		return false
	}
	return true
}

// ValidRecordParam: a record must be a valid type transformation expression
// or a node of the form: {prop:TTLExp, prop:TTLExp, ...}. Notice that the
// values are mandatory and they must be valid type transformation
// expressions.
//
// Java: TypeTransformationParser.validRecordParam(Node)
func (p *TypeTransformationParser) ValidRecordParam(expr rhino.Node) bool {
	if expr.IsObjectLit() {
		// Each value of a property must be a valid expression
		for prop := expr.GetFirstChild(); !prop.IsNil(); prop = prop.GetNext() {
			if prop.IsShorthandProperty() {
				p.WarnInvalid("property, missing type")
				return false
			} else if !p.ValidTypeTransformationExpression(prop.GetFirstChild()) {
				return false
			}
		}
	} else if !p.ValidTypeTransformationExpression(expr) {
		return false
	}
	return true
}

// ValidRecordTypeExpression: a record type expression must be of the form:
// record(RecordExp, RecordExp, ...).
//
// Java: TypeTransformationParser.validRecordTypeExpression(Node)
func (p *TypeTransformationParser) ValidRecordTypeExpression(expr rhino.Node) bool {
	// The expression must have at least two children. The record keyword and
	// a record expression
	if !p.CheckParameterCount(expr, TypeTransformationParser_Keywords_RECORD) {
		return false
	}
	// Each child must be a valid record
	for i := 0; i < p.GetCallParamCount(expr); i++ {
		if !p.ValidRecordParam(p.GetCallArgument(expr, i)) {
			p.WarnInvalidInside(TypeTransformationParser_Keywords_RECORD.KeywordName())
			return false
		}
	}
	return true
}

// Java: TypeTransformationParser.validNativeTypeExpr(Node)
func (p *TypeTransformationParser) ValidNativeTypeExpr(expr rhino.Node) bool {
	// The expression must have two children:
	// - The typeExpr keyword
	// - A string
	if !p.CheckParameterCount(expr, TypeTransformationParser_Keywords_TYPEEXPR) {
		return false
	}
	typeString := p.GetCallArgument(expr, 0)
	if !typeString.IsStringLit() {
		p.WarnInvalidExpression("native type")
		p.WarnInvalidInside(TypeTransformationParser_Keywords_TYPEEXPR.KeywordName())
		return false
	}
	typeExpr := JsDocInfoParser_ParseTypeString(typeString.GetString())
	typeString.Detach()
	expr.AddChildToBack(typeExpr)
	return true
}

// ValidTypeExpression: a TTL type expression must be a union type, a
// template type, a record type or any of the type predicates (none,
// rawTypeOf, templateTypeOf).
//
// Java: TypeTransformationParser.validTypeExpression(Node)
func (p *TypeTransformationParser) ValidTypeExpression(expr rhino.Node) bool {
	name := p.GetCallName(expr)
	keyword := p.NameToKeyword(name)
	switch keyword {
	case TypeTransformationParser_Keywords_TYPE:
		return p.ValidTemplateTypeExpression(expr)
	case TypeTransformationParser_Keywords_UNION:
		return p.ValidUnionTypeExpression(expr)
	case TypeTransformationParser_Keywords_NONE:
		return p.ValidNoneTypeExpression(expr)
	case TypeTransformationParser_Keywords_ALL:
		return p.ValidAllTypeExpression(expr)
	case TypeTransformationParser_Keywords_UNKNOWN:
		return p.ValidUnknownTypeExpression(expr)
	case TypeTransformationParser_Keywords_RAWTYPEOF:
		return p.ValidRawTypeOfTypeExpression(expr)
	case TypeTransformationParser_Keywords_TEMPLATETYPEOF:
		return p.ValidTemplateTypeOfExpression(expr)
	case TypeTransformationParser_Keywords_RECORD:
		return p.ValidRecordTypeExpression(expr)
	case TypeTransformationParser_Keywords_TYPEEXPR:
		return p.ValidNativeTypeExpr(expr)
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Invalid type expression"))
	}
}

// Java: TypeTransformationParser.validTypePredicate(Node, int)
func (p *TypeTransformationParser) ValidTypePredicate(expr rhino.Node, paramCount int) bool {
	// All the types must be valid type expressions
	for i := 0; i < paramCount; i++ {
		if !p.ValidTypeTransformationExpression(p.GetCallArgument(expr, i)) {
			p.WarnInvalidInside("boolean")
			return false
		}
	}
	return true
}

// Java: TypeTransformationParser.isValidStringParam(Node)
func (p *TypeTransformationParser) IsValidStringParam(expr rhino.Node) bool {
	if !expr.IsName() && !expr.IsStringLit() {
		p.WarnInvalid("string")
		return false
	}
	if expr.GetString() == "" {
		p.WarnInvalid("string parameter")
		return false
	}
	return true
}

// Java: TypeTransformationParser.validStringPredicate(Node, int)
func (p *TypeTransformationParser) ValidStringPredicate(expr rhino.Node, paramCount int) bool {
	// Each parameter must be valid string parameter
	for i := 0; i < paramCount; i++ {
		if !p.IsValidStringParam(p.GetCallArgument(expr, i)) {
			p.WarnInvalidInside("boolean")
			return false
		}
	}
	return true
}

// Java: TypeTransformationParser.validTypevarParam(Node)
func (p *TypeTransformationParser) ValidTypevarParam(expr rhino.Node) bool {
	if !p.IsTypeVar(expr) {
		p.WarnInvalid("name")
		return false
	}
	return true
}

// Java: TypeTransformationParser.validTypevarPredicate(Node, int)
func (p *TypeTransformationParser) ValidTypevarPredicate(expr rhino.Node, paramCount int) bool {
	// Each parameter must be valid string parameter
	for i := 0; i < paramCount; i++ {
		if !p.ValidTypevarParam(p.GetCallArgument(expr, i)) {
			p.WarnInvalidInside("boolean")
			return false
		}
	}
	return true
}

// Java: TypeTransformationParser.validBooleanOperation(Node)
func (p *TypeTransformationParser) ValidBooleanOperation(expr rhino.Node) bool {
	var valid bool
	if expr.IsNot() {
		valid = p.ValidBooleanExpression(expr.GetFirstChild())
	} else {
		valid = p.ValidBooleanExpression(expr.GetFirstChild()) &&
			p.ValidBooleanExpression(expr.GetSecondChild())
	}
	if !valid {
		p.WarnInvalidInside("boolean")
		return false
	}
	return true
}

// ValidBooleanExpression: a boolean expression must be a boolean predicate
// or a boolean type predicate.
//
// Java: TypeTransformationParser.validBooleanExpression(Node)
func (p *TypeTransformationParser) ValidBooleanExpression(expr rhino.Node) bool {
	if p.IsBooleanOperation(expr) {
		return p.ValidBooleanOperation(expr)
	}

	if !p.IsOperation(expr) {
		p.WarnInvalidExpression("boolean")
		return false
	}
	if !p.IsValidPredicate(p.GetCallName(expr)) {
		p.WarnInvalid("boolean predicate")
		return false
	}
	keyword := p.NameToKeyword(p.GetCallName(expr))
	if !p.CheckParameterCount(expr, keyword) {
		return false
	}
	switch keyword.Kind() {
	case TypeTransformationParser_OperationKind_TYPE_PREDICATE:
		return p.ValidTypePredicate(expr, p.GetCallParamCount(expr))
	case TypeTransformationParser_OperationKind_STRING_PREDICATE:
		return p.ValidStringPredicate(expr, p.GetCallParamCount(expr))
	case TypeTransformationParser_OperationKind_TYPEVAR_PREDICATE:
		return p.ValidTypevarPredicate(expr, p.GetCallParamCount(expr))
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Invalid boolean expression"))
	}
}

// ValidConditionalExpression: a conditional type transformation expression
// must be of the form cond(BoolExp, TTLExp, TTLExp).
//
// Java: TypeTransformationParser.validConditionalExpression(Node)
func (p *TypeTransformationParser) ValidConditionalExpression(expr rhino.Node) bool {
	// The expression must have four children:
	// - The cond keyword
	// - A boolean expression
	// - A type transformation expression with the 'if' branch
	// - A type transformation expression with the 'else' branch
	if !p.CheckParameterCount(expr, TypeTransformationParser_Keywords_COND) {
		return false
	}
	// Check for the validity of the boolean and the expressions
	if !p.ValidBooleanExpression(p.GetCallArgument(expr, 0)) {
		p.WarnInvalidInside("conditional")
		return false
	}
	if !p.ValidTypeTransformationExpression(p.GetCallArgument(expr, 1)) {
		p.WarnInvalidInside("conditional")
		return false
	}
	if !p.ValidTypeTransformationExpression(p.GetCallArgument(expr, 2)) {
		p.WarnInvalidInside("conditional")
		return false
	}
	return true
}

// ValidMapunionExpression: a mapunion type transformation expression must
// be of the form mapunion(TTLExp, (typevar) => TTLExp).
//
// Java: TypeTransformationParser.validMapunionExpression(Node)
func (p *TypeTransformationParser) ValidMapunionExpression(expr rhino.Node) bool {
	// The expression must have four children:
	// - The mapunion keyword
	// - A union type expression
	// - A map function
	if !p.CheckParameterCount(expr, TypeTransformationParser_Keywords_MAPUNION) {
		return false
	}
	// The second child must be a valid union type expression
	if !p.ValidTypeTransformationExpression(p.GetCallArgument(expr, 0)) {
		p.WarnInvalidInside(TypeTransformationParser_Keywords_MAPUNION.KeywordName())
		return false
	}
	// The third child must be a function
	if !p.GetCallArgument(expr, 1).IsFunction() {
		p.WarnInvalid("map function")
		p.WarnInvalidInside(TypeTransformationParser_Keywords_MAPUNION.KeywordName())
		return false
	}
	mapFn := p.GetCallArgument(expr, 1)
	// The map function must have only one parameter
	mapFnParamCount := p.GetFunctionParamCount(mapFn)
	if mapFnParamCount < 1 {
		p.WarnMissingParam("map function")
		p.WarnInvalidInside(TypeTransformationParser_Keywords_MAPUNION.KeywordName())
		return false
	}
	if mapFnParamCount > 1 {
		p.WarnExtraParam("map function")
		p.WarnInvalidInside(TypeTransformationParser_Keywords_MAPUNION.KeywordName())
		return false
	}
	// The body must be a valid type transformation expression
	mapFnBody := p.GetFunctionBody(mapFn)
	if !p.ValidTypeTransformationExpression(mapFnBody) {
		p.WarnInvalidInside("map function body")
		return false
	}
	return true
}

// ValidMaprecordExpression: a maprecord type transformation expression must
// be of the form maprecord(TTLExp, (typevar, typevar) => TTLExp).
//
// Java: TypeTransformationParser.validMaprecordExpression(Node)
func (p *TypeTransformationParser) ValidMaprecordExpression(expr rhino.Node) bool {
	// The expression must have four children:
	// - The maprecord keyword
	// - A type expression
	// - A map function
	if !p.CheckParameterCount(expr, TypeTransformationParser_Keywords_MAPRECORD) {
		return false
	}
	// The second child must be a valid expression
	if !p.ValidTypeTransformationExpression(p.GetCallArgument(expr, 0)) {
		p.WarnInvalidInside(TypeTransformationParser_Keywords_MAPRECORD.KeywordName())
		return false
	}
	// The third child must be a function
	if !p.GetCallArgument(expr, 1).IsFunction() {
		p.WarnInvalid("map function")
		p.WarnInvalidInside(TypeTransformationParser_Keywords_MAPRECORD.KeywordName())
		return false
	}
	mapFn := p.GetCallArgument(expr, 1)
	// The map function must have exactly two parameters
	mapFnParamCount := p.GetFunctionParamCount(mapFn)
	if mapFnParamCount < 2 {
		p.WarnMissingParam("map function")
		p.WarnInvalidInside(TypeTransformationParser_Keywords_MAPRECORD.KeywordName())
		return false
	}
	if mapFnParamCount > 2 {
		p.WarnExtraParam("map function")
		p.WarnInvalidInside(TypeTransformationParser_Keywords_MAPRECORD.KeywordName())
		return false
	}
	// The body must be a valid type transformation expression
	mapFnBody := p.GetFunctionBody(mapFn)
	if !p.ValidTypeTransformationExpression(mapFnBody) {
		p.WarnInvalidInside("map function body")
		return false
	}
	return true
}

// ValidTypeOfVarExpression: a typeOfVar expression must be of the form
// typeOfVar('name').
//
// Java: TypeTransformationParser.validTypeOfVarExpression(Node)
func (p *TypeTransformationParser) ValidTypeOfVarExpression(expr rhino.Node) bool {
	// The expression must have two children:
	// - The typeOfVar keyword
	// - A string
	if !p.CheckParameterCount(expr, TypeTransformationParser_Keywords_TYPEOFVAR) {
		return false
	}
	if !p.GetCallArgument(expr, 0).IsStringLit() {
		p.WarnInvalid("name")
		p.WarnInvalidInside(TypeTransformationParser_Keywords_TYPEOFVAR.KeywordName())
		return false
	}
	return true
}

// ValidInstanceOfExpression: a typeOfVar expression must be of the form
// instanceOf('name').
//
// Java: TypeTransformationParser.validInstanceOfExpression(Node)
func (p *TypeTransformationParser) ValidInstanceOfExpression(expr rhino.Node) bool {
	// The expression must have two children:
	// - The instanceOf keyword
	// - A string
	if !p.CheckParameterCount(expr, TypeTransformationParser_Keywords_INSTANCEOF) {
		return false
	}
	if !p.ValidTypeTransformationExpression(p.GetCallArgument(expr, 0)) {
		p.WarnInvalidInside(TypeTransformationParser_Keywords_INSTANCEOF.KeywordName())
		return false
	}
	return true
}

// Java: TypeTransformationParser.validPrintTypeExpression(Node)
func (p *TypeTransformationParser) ValidPrintTypeExpression(expr rhino.Node) bool {
	// The expression must have three children. The printType keyword, a
	// message and a type transformation expression
	if !p.CheckParameterCount(expr, TypeTransformationParser_Keywords_PRINTTYPE) {
		return false
	}
	if !p.GetCallArgument(expr, 0).IsStringLit() {
		p.WarnInvalid("message")
		p.WarnInvalidInside(TypeTransformationParser_Keywords_PRINTTYPE.KeywordName())
		return false
	}
	if !p.ValidTypeTransformationExpression(p.GetCallArgument(expr, 1)) {
		p.WarnInvalidInside(TypeTransformationParser_Keywords_PRINTTYPE.KeywordName())
		return false
	}
	return true
}

// Java: TypeTransformationParser.validPropTypeExpression(Node)
func (p *TypeTransformationParser) ValidPropTypeExpression(expr rhino.Node) bool {
	// The expression must have three children. The propType keyword, a
	// a string and a type transformation expression
	if !p.CheckParameterCount(expr, TypeTransformationParser_Keywords_PROPTYPE) {
		return false
	}
	if !p.GetCallArgument(expr, 0).IsStringLit() {
		p.WarnInvalid("property name")
		p.WarnInvalidInside(TypeTransformationParser_Keywords_PROPTYPE.KeywordName())
		return false
	}
	if !p.ValidTypeTransformationExpression(p.GetCallArgument(expr, 1)) {
		p.WarnInvalidInside(TypeTransformationParser_Keywords_PROPTYPE.KeywordName())
		return false
	}
	return true
}

// ValidOperationExpression: an operation expression is a cond or a
// mapunion.
//
// Java: TypeTransformationParser.validOperationExpression(Node)
func (p *TypeTransformationParser) ValidOperationExpression(expr rhino.Node) bool {
	name := p.GetCallName(expr)
	keyword := p.NameToKeyword(name)
	switch keyword {
	case TypeTransformationParser_Keywords_COND:
		return p.ValidConditionalExpression(expr)
	case TypeTransformationParser_Keywords_MAPUNION:
		return p.ValidMapunionExpression(expr)
	case TypeTransformationParser_Keywords_MAPRECORD:
		return p.ValidMaprecordExpression(expr)
	case TypeTransformationParser_Keywords_TYPEOFVAR:
		return p.ValidTypeOfVarExpression(expr)
	case TypeTransformationParser_Keywords_INSTANCEOF:
		return p.ValidInstanceOfExpression(expr)
	case TypeTransformationParser_Keywords_PRINTTYPE:
		return p.ValidPrintTypeExpression(expr)
	case TypeTransformationParser_Keywords_PROPTYPE:
		return p.ValidPropTypeExpression(expr)
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Invalid type transformation operation"))
	}
}

// ValidTypeTransformationExpression checks the structure of the AST of a
// type transformation expression in @template T := TTLExp =:.
//
// Java: TypeTransformationParser.validTypeTransformationExpression(Node)
func (p *TypeTransformationParser) ValidTypeTransformationExpression(expr rhino.Node) bool {
	if !p.IsValidExpression(expr) {
		p.WarnInvalidExpression("type transformation")
		return false
	}
	if p.IsTypeVar(expr) || p.IsTypeName(expr) {
		return true
	}
	// Check for valid keyword
	name := p.GetCallName(expr)
	if !p.IsValidKeyword(name) {
		p.WarnInvalidExpression("type transformation")
		return false
	}
	keyword := p.NameToKeyword(name)
	// Check the rest of the expression depending on the kind
	switch keyword.Kind() {
	case TypeTransformationParser_OperationKind_TYPE_CONSTRUCTOR:
		return p.ValidTypeExpression(expr)
	case TypeTransformationParser_OperationKind_OPERATION:
		return p.ValidOperationExpression(expr)
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Invalid type transformation expression"))
	}
}
