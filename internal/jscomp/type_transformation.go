// Port of com/google/javascript/jscomp/TypeTransformation.java @ closure-compiler v20261021

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

package jscomp

import (
	"fmt"
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const typeTransformation_VIRTUAL_FILE = "<TypeTransformation.java>"

var TypeTransformation_UNKNOWN_TYPEVAR = DiagnosticType_Warning("TYPEVAR_UNDEFINED",
	"Reference to an unknown type variable {0}")
var TypeTransformation_UNKNOWN_STRVAR = DiagnosticType_Warning("UNKNOWN_STRVAR",
	"Reference to an unknown string variable {0}")
var TypeTransformation_UNKNOWN_TYPENAME = DiagnosticType_Warning("TYPENAME_UNDEFINED",
	"Reference to an unknown type name {0}")
var TypeTransformation_BASETYPE_INVALID = DiagnosticType_Warning("BASETYPE_INVALID",
	"The type {0} cannot be templatized")
var TypeTransformation_TEMPTYPE_INVALID = DiagnosticType_Warning("TEMPTYPE_INVALID",
	"Expected templatized type in {0} found {1}")
var TypeTransformation_INDEX_OUTOFBOUNDS = DiagnosticType_Warning("INDEX_OUTOFBOUNDS",
	"Index out of bounds in templateTypeOf: expected a number less than {0}, found {1}")
var TypeTransformation_DUPLICATE_VARIABLE = DiagnosticType_Warning("DUPLICATE_VARIABLE",
	"The variable {0} is already defined")

// TypeTransformation_UNKNOWN_NAMEVAR: This warning is never exercised.
var TypeTransformation_UNKNOWN_NAMEVAR = DiagnosticType_Warning("UNKNOWN_NAMEVAR",
	"Reference to an unknown name variable {0}")
var TypeTransformation_RECTYPE_INVALID = DiagnosticType_Warning("RECTYPE_INVALID",
	"The first parameter of a maprecord must be a record type, "+
		"found {0}")
var TypeTransformation_MAPRECORD_BODY_INVALID = DiagnosticType_Warning("MAPRECORD_BODY_INVALID",
	"The body of a maprecord function must evaluate to a record type or "+
		"a no type, found {0}")
var TypeTransformation_VAR_UNDEFINED = DiagnosticType_Warning("VAR_UNDEFINED",
	"Variable {0} is undefined in the scope")
var TypeTransformation_INVALID_CTOR = DiagnosticType_Warning("INVALID_CTOR",
	"Expected a constructor type, found {0}")
var TypeTransformation_RECPARAM_INVALID = DiagnosticType_Warning("RECPARAM_INVALID",
	"Expected a record type, found {0}")
var TypeTransformation_PROPTYPE_INVALID = DiagnosticType_Warning("PROPTYPE_INVALID",
	"Expected object type, found {0}")

// TypeTransformation is a class for processing type transformation
// expressions
type TypeTransformation struct {
	compiler AbstractCompiler
	registry *rhino.JSTypeRegistry
	typeEnv  rhino.StaticTypedScope
}

// typeTransformation_NameResolver is a helper class for holding the
// information about the type variables and the name variables in maprecord
// expressions
//
// Java: private static class TypeTransformation.NameResolver
type typeTransformation_NameResolver struct {
	typeVars *jutil.LinkedHashMap[string, rhino.JSType]
	nameVars *jutil.LinkedHashMap[string, string]
}

// Java: TypeTransformation.NameResolver(ImmutableMap<String, JSType>, ImmutableMap<String, String>)
func newTypeTransformation_NameResolver(
	typeVars *jutil.LinkedHashMap[string, rhino.JSType], nameVars *jutil.LinkedHashMap[string, string]) *typeTransformation_NameResolver {
	return &typeTransformation_NameResolver{typeVars: typeVars, nameVars: nameVars}
}

// Java: TypeTransformation(AbstractCompiler, StaticTypedScope)
func NewTypeTransformation(compiler AbstractCompiler, typeEnv rhino.StaticTypedScope) *TypeTransformation {
	return &TypeTransformation{
		compiler: compiler,
		registry: compiler.GetTypeRegistry(),
		typeEnv:  typeEnv,
	}
}

// Java: TypeTransformation.isTypeVar(Node)
func (t *TypeTransformation) isTypeVar(n *rhino.Node) bool {
	return n.IsName()
}

// Java: TypeTransformation.isTypeName(Node)
func (t *TypeTransformation) isTypeName(n *rhino.Node) bool {
	return n.IsStringLit()
}

// Java: TypeTransformation.isBooleanOperation(Node)
func (t *TypeTransformation) isBooleanOperation(n *rhino.Node) bool {
	return n.IsAnd() || n.IsOr() || n.IsNot()
}

// Java: TypeTransformation.nameToKeyword(String)
func (t *TypeTransformation) nameToKeyword(s string) parsing.TypeTransformationParser_Keywords {
	return parsing.TypeTransformationParser_Keywords_FromCaseInsensitiveASCII(s)
}

// Java: TypeTransformation.getType(String)
func (t *TypeTransformation) getType(typeName string) rhino.JSType {
	type_ := t.registry.GetType(t.typeEnv, typeName)
	if type_ != nil {
		return type_
	}
	slot, _ := t.typeEnv.GetSlot(typeName).(rhino.StaticTypedSlot)
	if slot != nil {
		type_ = slot.GetType()
	} else {
		type_ = nil
	}
	if type_ != nil {
		if type_.IsConstructor() || type_.IsInterface() {
			return type_.ToMaybeFunctionType().GetInstanceType().GetRawType()
		}
		if type_.IsEnumElementType() {
			return type_.GetEnumeratedTypeOfEnumElement()
		}
		return type_
	}
	var jsdoc *rhino.JSDocInfo
	if slot != nil {
		jsdoc = slot.GetJSDocInfo()
	}
	if jsdoc != nil && jsdoc.HasTypedefType() {
		return t.registry.EvaluateTypeExpression(jsdoc.GetTypedefType(), t.typeEnv)
	}
	return nil
}

// Java: TypeTransformation.getUnknownType()
func (t *TypeTransformation) getUnknownType() rhino.JSType {
	return t.registry.GetNativeObjectType(rhino.JSTypeNative_UNKNOWN_TYPE)
}

// Java: TypeTransformation.getNoType()
func (t *TypeTransformation) getNoType() rhino.JSType {
	return t.registry.GetNativeObjectType(rhino.JSTypeNative_NO_TYPE)
}

// Java: TypeTransformation.getAllType()
func (t *TypeTransformation) getAllType() rhino.JSType {
	return t.registry.GetNativeType(rhino.JSTypeNative_ALL_TYPE)
}

// Java: TypeTransformation.getObjectType()
func (t *TypeTransformation) getObjectType() rhino.JSType {
	return t.registry.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE)
}

// Java: TypeTransformation.createUnionType(JSType[])
func (t *TypeTransformation) createUnionType(variants []rhino.JSType) rhino.JSType {
	return t.registry.CreateUnionType(variants...)
}

// Java: TypeTransformation.createRecordType(ImmutableMap<String, JSType>)
func (t *TypeTransformation) createRecordType(props *jutil.LinkedHashMap[string, rhino.JSType]) rhino.JSType {
	return t.registry.CreateRecordType(props)
}

// Java: TypeTransformation.reportWarning(Node, DiagnosticType, String...)
func (t *TypeTransformation) reportWarning(n *rhino.Node, msg *DiagnosticType, param ...string) {
	t.compiler.Report(JSError_MakeWithNode(n, msg, param...))
}

// Java: TypeTransformation.addNewEntry(ImmutableMap<String, T>, String, T)
func typeTransformation_addNewEntry[T any](
	map_ *jutil.LinkedHashMap[string, T], name string, type_ T) *jutil.LinkedHashMap[string, T] {
	builder := jutil.NewImmutableMapBuilder[string, T]()
	for k, v := range map_.All() {
		builder.Put(k, v)
	}
	return builder.Put(name, type_).BuildOrThrow()
}

// Java: TypeTransformation.getFunctionParameter(Node, int)
func (t *TypeTransformation) getFunctionParameter(n *rhino.Node, i int) string {
	jutil.CheckArgumentMsg1(n.IsFunction(), "Expected a function node, found %s", n)
	return n.GetSecondChild().GetChildAtIndex(i).GetString()
}

// Java: TypeTransformation.getCallName(Node)
func (t *TypeTransformation) getCallName(n *rhino.Node) string {
	jutil.CheckArgumentMsg1(n.IsCall(), "Expected a call node, found %s", n)
	return n.GetFirstChild().GetString()
}

// Java: TypeTransformation.getCallArgument(Node, int)
func (t *TypeTransformation) getCallArgument(n *rhino.Node, i int) *rhino.Node {
	jutil.CheckArgumentMsg1(n.IsCall(), "Expected a call node, found %s", n)
	return n.GetChildAtIndex(i + 1)
}

// Java: TypeTransformation.getCallParamCount(Node)
func (t *TypeTransformation) getCallParamCount(n *rhino.Node) int {
	jutil.CheckArgumentMsg1(n.IsCall(), "Expected a call node, found %s", n)
	return n.GetChildCount() - 1
}

// getCallParams: TODO(dimvar): rewrite the uses of this method to use
// siblings() and delete it. Copying is unnecessarily inefficient.
//
// Java: TypeTransformation.getCallParams(Node)
func (t *TypeTransformation) getCallParams(n *rhino.Node) []*rhino.Node {
	jutil.CheckArgumentMsg1(n.IsCall(), "Expected a call node, found %s", n)
	var builder []*rhino.Node
	for i := 0; i < t.getCallParamCount(n); i++ {
		builder = append(builder, t.getCallArgument(n, i))
	}
	return builder
}

// Java: TypeTransformation.getComputedPropValue(Node)
func (t *TypeTransformation) getComputedPropValue(n *rhino.Node) *rhino.Node {
	jutil.CheckArgumentMsg1(
		n.IsComputedProp(), "Expected a computed property node, found %s", n)
	return n.GetSecondChild()
}

// Java: TypeTransformation.getComputedPropName(Node)
func (t *TypeTransformation) getComputedPropName(n *rhino.Node) string {
	jutil.CheckArgumentMsg1(
		n.IsComputedProp(), "Expected a computed property node, found %s", n)
	return n.GetFirstChild().GetString()
}

// Eval evaluates the type transformation expression and returns the
// resulting type.
//
// ttlAst is the node representing the type transformation expression;
// typeVars is the environment containing the information about the type
// variables. Returns the resulting type after the transformation
//
// Java: TypeTransformation.eval(Node, ImmutableMap<String, JSType>)
func (t *TypeTransformation) Eval(ttlAst *rhino.Node, typeVars *jutil.LinkedHashMap[string, rhino.JSType]) rhino.JSType {
	return t.EvalWithNameVars(ttlAst, typeVars, jutil.NewLinkedHashMap[string, string]())
}

// EvalWithNameVars evaluates the type transformation expression and returns
// the resulting type.
//
// ttlAst is the node representing the type transformation expression;
// typeVars is the environment containing the information about the type
// variables; nameVars is the environment containing the information about
// the name variables. Returns the resulting type after the transformation
//
// Java: TypeTransformation.eval(Node, ImmutableMap<String, JSType>, ImmutableMap<String, String>)
func (t *TypeTransformation) EvalWithNameVars(
	ttlAst *rhino.Node, typeVars *jutil.LinkedHashMap[string, rhino.JSType], nameVars *jutil.LinkedHashMap[string, string]) rhino.JSType {
	result := t.evalInternal(ttlAst, newTypeTransformation_NameResolver(typeVars, nameVars))
	if result.IsEmptyType() {
		return t.getUnknownType()
	}
	return result
}

// Java: TypeTransformation.evalInternal(Node, NameResolver)
func (t *TypeTransformation) evalInternal(ttlAst *rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	if t.isTypeName(ttlAst) {
		return t.evalTypeName(ttlAst)
	}
	if t.isTypeVar(ttlAst) {
		return t.evalTypeVar(ttlAst, nameResolver)
	}
	name := t.getCallName(ttlAst)
	keyword := t.nameToKeyword(name)
	switch keyword.Kind() {
	case parsing.TypeTransformationParser_OperationKind_TYPE_CONSTRUCTOR:
		return t.evalTypeExpressionWithKeyword(ttlAst, keyword, nameResolver)
	case parsing.TypeTransformationParser_OperationKind_OPERATION:
		return t.evalOperationExpressionWithKeyword(ttlAst, keyword, nameResolver)
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Could not evaluate the type transformation expression"))
	}
}

// Java: TypeTransformation.evalOperationExpression(Node, NameResolver)
func (t *TypeTransformation) evalOperationExpression(ttlAst *rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	return t.evalOperationExpressionWithKeyword(ttlAst, t.nameToKeyword(t.getCallName(ttlAst)), nameResolver)
}

func (t *TypeTransformation) evalOperationExpressionWithKeyword(ttlAst *rhino.Node, keyword parsing.TypeTransformationParser_Keywords, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	switch keyword {
	case parsing.TypeTransformationParser_Keywords_COND:
		return t.evalConditional(ttlAst, nameResolver)
	case parsing.TypeTransformationParser_Keywords_MAPUNION:
		return t.evalMapunion(ttlAst, nameResolver)
	case parsing.TypeTransformationParser_Keywords_MAPRECORD:
		return t.evalMaprecord(ttlAst, nameResolver)
	case parsing.TypeTransformationParser_Keywords_TYPEOFVAR:
		return t.evalTypeOfVar(ttlAst)
	case parsing.TypeTransformationParser_Keywords_INSTANCEOF:
		return t.evalInstanceOf(ttlAst, nameResolver)
	case parsing.TypeTransformationParser_Keywords_PRINTTYPE:
		return t.evalPrintType(ttlAst, nameResolver)
	case parsing.TypeTransformationParser_Keywords_PROPTYPE:
		return t.evalPropType(ttlAst, nameResolver)
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Invalid type transformation operation"))
	}
}

// Java: TypeTransformation.evalTypeExpression(Node, NameResolver)
func (t *TypeTransformation) evalTypeExpression(ttlAst *rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	return t.evalTypeExpressionWithKeyword(ttlAst, t.nameToKeyword(t.getCallName(ttlAst)), nameResolver)
}

func (t *TypeTransformation) evalTypeExpressionWithKeyword(ttlAst *rhino.Node, keyword parsing.TypeTransformationParser_Keywords, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	switch keyword {
	case parsing.TypeTransformationParser_Keywords_TYPE:
		return t.evalTemplatizedType(ttlAst, nameResolver)
	case parsing.TypeTransformationParser_Keywords_UNION:
		return t.evalUnionType(ttlAst, nameResolver)
	case parsing.TypeTransformationParser_Keywords_NONE:
		return t.getNoType()
	case parsing.TypeTransformationParser_Keywords_ALL:
		return t.getAllType()
	case parsing.TypeTransformationParser_Keywords_UNKNOWN:
		return t.getUnknownType()
	case parsing.TypeTransformationParser_Keywords_RAWTYPEOF:
		return t.evalRawTypeOf(ttlAst, nameResolver)
	case parsing.TypeTransformationParser_Keywords_TEMPLATETYPEOF:
		return t.evalTemplateTypeOf(ttlAst, nameResolver)
	case parsing.TypeTransformationParser_Keywords_RECORD:
		return t.evalRecordType(ttlAst, nameResolver)
	case parsing.TypeTransformationParser_Keywords_TYPEEXPR:
		return t.evalNativeTypeExpr(ttlAst)
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Invalid type expression"))
	}
}

// Java: TypeTransformation.evalTypeName(Node)
func (t *TypeTransformation) evalTypeName(ttlAst *rhino.Node) rhino.JSType {
	typeName := ttlAst.GetString()
	resultingType := t.getType(typeName)
	// If the type name is not defined then return UNKNOWN and report a warning
	if resultingType == nil {
		t.reportWarning(ttlAst, TypeTransformation_UNKNOWN_TYPENAME, typeName)
		return t.getUnknownType()
	}
	return resultingType
}

// Java: TypeTransformation.evalTemplatizedType(Node, NameResolver)
func (t *TypeTransformation) evalTemplatizedType(ttlAst *rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	jutil.CheckArgumentMsg1(ttlAst.IsCall(), "Expected a call node, found %s", ttlAst)
	firstArgNode := ttlAst.GetFirstChild().GetNext()
	firstParam := t.evalInternal(firstArgNode, nameResolver)
	if !firstParam.IsRawTypeOfTemplatizedType() {
		t.reportWarning(ttlAst, TypeTransformation_BASETYPE_INVALID, firstParam.String())
		return t.getUnknownType()
	}
	// TODO(lpino): Check that the number of parameters correspond with the
	// number of template types that the base type can take when creating
	// a templatized type. For instance, if the base type is Array then there
	// must be just one parameter.
	restCount := 0
	for arg := firstArgNode.GetNext(); arg != nil; arg = arg.GetNext() {
		restCount++
	}
	templatizedTypes := make([]rhino.JSType, restCount)
	i := 0
	for arg := firstArgNode.GetNext(); arg != nil; arg = arg.GetNext() {
		templatizedTypes[i] = t.evalInternal(arg, nameResolver)
		i++
	}
	baseType := firstParam.ToMaybeObjectType()
	return t.registry.CreateTemplatizedType(baseType, templatizedTypes...)
}

// Java: TypeTransformation.evalTypeVar(Node, NameResolver)
func (t *TypeTransformation) evalTypeVar(ttlAst *rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	typeVar := ttlAst.GetString()
	resultingType := nameResolver.typeVars.Get(typeVar)
	// If the type variable is not defined then return UNKNOWN and report a warning
	if resultingType == nil {
		t.reportWarning(ttlAst, TypeTransformation_UNKNOWN_TYPEVAR, typeVar)
		return t.getUnknownType()
	}
	return resultingType
}

// Java: TypeTransformation.evalUnionType(Node, NameResolver)
func (t *TypeTransformation) evalUnionType(ttlAst *rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	jutil.CheckArgumentMsg1(ttlAst.IsCall(), "Expected a call node, found %s", ttlAst)
	firstArg := ttlAst.GetFirstChild().GetNext()
	paramCount := 0
	for arg := firstArg; arg != nil; arg = arg.GetNext() {
		paramCount++
	}
	// Create an array of types after evaluating each parameter
	basicTypes := make([]rhino.JSType, paramCount)
	i := 0
	for arg := firstArg; arg != nil; arg = arg.GetNext() {
		basicTypes[i] = t.evalInternal(arg, nameResolver)
		i++
	}
	return t.createUnionType(basicTypes)
}

// Java: TypeTransformation.evalTypeParams(Node, NameResolver)
func (t *TypeTransformation) evalTypeParams(ttlAst *rhino.Node, nameResolver *typeTransformation_NameResolver) []rhino.JSType {
	jutil.CheckArgumentMsg1(ttlAst.IsCall(), "Expected a call node, found %s", ttlAst)
	firstArg := ttlAst.GetFirstChild().GetNext()
	paramCount := 0
	for arg := firstArg; arg != nil; arg = arg.GetNext() {
		paramCount++
	}
	result := make([]rhino.JSType, paramCount)
	i := 0
	for arg := firstArg; arg != nil; arg = arg.GetNext() {
		result[i] = t.evalInternal(arg, nameResolver)
		i++
	}
	return result
}

// Java: TypeTransformation.evalString(Node, NameResolver)
func (t *TypeTransformation) evalString(ttlAst *rhino.Node, nameResolver *typeTransformation_NameResolver) string {
	if ttlAst.IsName() {
		// Return the empty string if the name variable cannot be resolved
		if !nameResolver.nameVars.ContainsKey(ttlAst.GetString()) {
			t.reportWarning(ttlAst, TypeTransformation_UNKNOWN_STRVAR, ttlAst.GetString())
			return ""
		}
		return nameResolver.nameVars.Get(ttlAst.GetString())
	}
	return ttlAst.GetString()
}

// Java: TypeTransformation.evalStringParams(Node, NameResolver)
func (t *TypeTransformation) evalStringParams(ttlAst *rhino.Node, nameResolver *typeTransformation_NameResolver) []string {
	params := t.getCallParams(ttlAst)
	paramCount := len(params)
	result := make([]string, paramCount)
	for i := 0; i < paramCount; i++ {
		result[i] = t.evalString(params[i], nameResolver)
	}
	return result
}

// Java: TypeTransformation.evalTypePredicate(Node, NameResolver)
func (t *TypeTransformation) evalTypePredicate(ttlAst *rhino.Node, nameResolver *typeTransformation_NameResolver) bool {
	return t.evalTypePredicateWithKeyword(ttlAst, t.nameToKeyword(t.getCallName(ttlAst)), nameResolver)
}

func (t *TypeTransformation) evalTypePredicateWithKeyword(ttlAst *rhino.Node, keyword parsing.TypeTransformationParser_Keywords, nameResolver *typeTransformation_NameResolver) bool {
	jutil.CheckArgumentMsg1(ttlAst.IsCall(), "Expected a call node, found %s", ttlAst)
	arg0 := ttlAst.GetFirstChild().GetNext()
	type_ := t.evalInternal(arg0, nameResolver)
	var type1 rhino.JSType
	if arg1 := arg0.GetNext(); arg1 != nil {
		type1 = t.evalInternal(arg1, nameResolver)
		for extra := arg1.GetNext(); extra != nil; extra = extra.GetNext() {
			t.evalInternal(extra, nameResolver)
		}
	}
	switch keyword {
	case parsing.TypeTransformationParser_Keywords_EQ:
		return type_.Equals(type1)
	case parsing.TypeTransformationParser_Keywords_SUB:
		return type_.IsSubtypeOf(type1)
	case parsing.TypeTransformationParser_Keywords_ISCTOR:
		return type_.IsConstructor()
	case parsing.TypeTransformationParser_Keywords_ISTEMPLATIZED:
		return type_.IsTemplatizedType()
	case parsing.TypeTransformationParser_Keywords_ISRECORD:
		return type_.IsRecordType()
	case parsing.TypeTransformationParser_Keywords_ISUNKNOWN:
		return type_.IsUnknownType()
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Invalid type predicate in the type transformation"))
	}
}

// Java: TypeTransformation.evalStringPredicate(Node, NameResolver)
func (t *TypeTransformation) evalStringPredicate(ttlAst *rhino.Node,
	nameResolver *typeTransformation_NameResolver) bool {
	return t.evalStringPredicateWithKeyword(ttlAst, t.nameToKeyword(t.getCallName(ttlAst)), nameResolver)
}

func (t *TypeTransformation) evalStringPredicateWithKeyword(ttlAst *rhino.Node, keyword parsing.TypeTransformationParser_Keywords, nameResolver *typeTransformation_NameResolver) bool {
	jutil.CheckArgumentMsg1(ttlAst.IsCall(), "Expected a call node, found %s", ttlAst)
	arg0 := ttlAst.GetFirstChild().GetNext()
	arg1 := arg0.GetNext()
	s0 := t.evalString(arg0, nameResolver)
	s1 := t.evalString(arg1, nameResolver)
	hasEmpty := s0 == "" || s1 == ""
	for extra := arg1.GetNext(); extra != nil; extra = extra.GetNext() {
		if t.evalString(extra, nameResolver) == "" {
			hasEmpty = true
		}
	}
	if hasEmpty {
		return false
	}
	switch keyword {
	case parsing.TypeTransformationParser_Keywords_STREQ:
		return s0 == s1
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Invalid string predicate in the type transformation"))
	}
}

// Java: TypeTransformation.evalTypevarPredicate(Node, NameResolver)
func (t *TypeTransformation) evalTypevarPredicate(ttlAst *rhino.Node, nameResolver *typeTransformation_NameResolver) bool {
	return t.evalTypevarPredicateWithKeyword(ttlAst, t.nameToKeyword(t.getCallName(ttlAst)), nameResolver)
}

func (t *TypeTransformation) evalTypevarPredicateWithKeyword(ttlAst *rhino.Node, keyword parsing.TypeTransformationParser_Keywords, nameResolver *typeTransformation_NameResolver) bool {
	switch keyword {
	case parsing.TypeTransformationParser_Keywords_ISDEFINED:
		return nameResolver.typeVars.ContainsKey(t.getCallArgument(ttlAst, 0).GetString())
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Invalid typevar predicate in the type transformation"))
	}
}

// Java: TypeTransformation.evalBooleanOperation(Node, NameResolver)
func (t *TypeTransformation) evalBooleanOperation(ttlAst *rhino.Node, nameResolver *typeTransformation_NameResolver) bool {
	param0 := t.evalBoolean(ttlAst.GetFirstChild(), nameResolver)
	if ttlAst.IsNot() {
		return !param0
	}
	if ttlAst.IsAnd() {
		return param0 && t.evalBoolean(ttlAst.GetLastChild(), nameResolver)
	}
	if ttlAst.IsOr() {
		return param0 || t.evalBoolean(ttlAst.GetLastChild(), nameResolver)
	}
	panic(jutil.NewIllegalStateExceptionMsg(
		"Invalid boolean predicate in the type transformation"))
}

// Java: TypeTransformation.evalBoolean(Node, NameResolver)
func (t *TypeTransformation) evalBoolean(ttlAst *rhino.Node, nameResolver *typeTransformation_NameResolver) bool {
	if t.isBooleanOperation(ttlAst) {
		return t.evalBooleanOperation(ttlAst, nameResolver)
	}
	name := t.getCallName(ttlAst)
	keyword := t.nameToKeyword(name)
	switch keyword.Kind() {
	case parsing.TypeTransformationParser_OperationKind_STRING_PREDICATE:
		return t.evalStringPredicateWithKeyword(ttlAst, keyword, nameResolver)
	case parsing.TypeTransformationParser_OperationKind_TYPE_PREDICATE:
		return t.evalTypePredicateWithKeyword(ttlAst, keyword, nameResolver)
	case parsing.TypeTransformationParser_OperationKind_TYPEVAR_PREDICATE:
		return t.evalTypevarPredicateWithKeyword(ttlAst, keyword, nameResolver)
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Invalid boolean predicate in the type transformation"))
	}
}

// Java: TypeTransformation.evalConditional(Node, NameResolver)
func (t *TypeTransformation) evalConditional(ttlAst *rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	jutil.CheckArgumentMsg1(ttlAst.IsCall(), "Expected a call node, found %s", ttlAst)
	p0 := ttlAst.GetFirstChild().GetNext()
	p1 := p0.GetNext()
	p2 := p1.GetNext()
	if t.evalBoolean(p0, nameResolver) {
		return t.evalInternal(p1, nameResolver)
	} else {
		return t.evalInternal(p2, nameResolver)
	}
}

// Java: TypeTransformation.evalMapunion(Node, NameResolver)
func (t *TypeTransformation) evalMapunion(ttlAst *rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	jutil.CheckArgumentMsg1(ttlAst.IsCall(), "Expected a call node, found %s", ttlAst)
	unionParam := ttlAst.GetFirstChild().GetNext()
	mapFunction := unionParam.GetNext()
	paramName := t.getFunctionParameter(mapFunction, 0)

	// The mapunion variable must not be defined in the environment
	if nameResolver.typeVars.ContainsKey(paramName) {
		t.reportWarning(ttlAst, TypeTransformation_DUPLICATE_VARIABLE, paramName)
		return t.getUnknownType()
	}

	mapFunctionBody := NodeUtil_GetFunctionBody(mapFunction)
	unionType := t.evalInternal(unionParam, nameResolver)
	// If the first parameter does not correspond to a union type then
	// consider it as a union with a single type and evaluate
	if !unionType.IsUnionType() {
		newNameResolver := newTypeTransformation_NameResolver(
			typeTransformation_addNewEntry(nameResolver.typeVars, paramName, unionType),
			nameResolver.nameVars)
		return t.evalInternal(mapFunctionBody, newNameResolver)
	}

	// Otherwise obtain the elements in the union type. Note that the block
	// above guarantees the casting to be safe
	unionElms := append([]rhino.JSType(nil), unionType.GetUnionMembers()...)
	// Evaluate the map function body using each element in the union type
	unionSize := len(unionElms)
	newUnionElms := make([]rhino.JSType, unionSize)
	i := 0
	for _, elm := range unionElms {
		newNameResolver := newTypeTransformation_NameResolver(
			typeTransformation_addNewEntry(nameResolver.typeVars, paramName, elm),
			nameResolver.nameVars)
		newUnionElms[i] = t.evalInternal(mapFunctionBody, newNameResolver)
		i++
	}

	return t.createUnionType(newUnionElms)
}

// Java: TypeTransformation.evalRawTypeOf(Node, NameResolver)
func (t *TypeTransformation) evalRawTypeOf(ttlAst *rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	jutil.CheckArgumentMsg1(ttlAst.IsCall(), "Expected a call node, found %s", ttlAst)
	type_ := t.evalInternal(ttlAst.GetFirstChild().GetNext(), nameResolver)
	if !type_.IsTemplatizedType() {
		t.reportWarning(ttlAst, TypeTransformation_TEMPTYPE_INVALID, "rawTypeOf", type_.String())
		return t.getUnknownType()
	}
	return type_.ToMaybeObjectType().GetRawType()
}

// Java: TypeTransformation.evalTemplateTypeOf(Node, NameResolver)
func (t *TypeTransformation) evalTemplateTypeOf(ttlAst *rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	jutil.CheckArgumentMsg1(ttlAst.IsCall(), "Expected a call node, found %s", ttlAst)
	p0 := ttlAst.GetFirstChild().GetNext()
	p1 := p0.GetNext()
	type_ := t.evalInternal(p0, nameResolver)
	if !type_.IsTemplatizedType() {
		t.reportWarning(ttlAst, TypeTransformation_TEMPTYPE_INVALID, "templateTypeOf", type_.String())
		return t.getUnknownType()
	}
	index := int(jutil.D2I(p1.GetDouble()))
	templateTypes := type_.ToMaybeObjectType().GetTemplateTypes()
	if index >= len(templateTypes) {
		t.reportWarning(ttlAst, TypeTransformation_INDEX_OUTOFBOUNDS,
			strconv.Itoa(len(templateTypes)), strconv.Itoa(index))
		return t.getUnknownType()
	}
	return templateTypes[index]
}

// Java: TypeTransformation.evalRecord(Node, NameResolver)
func (t *TypeTransformation) evalRecord(record *rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	props := jutil.NewLinkedHashMap[string, rhino.JSType]()
	for propNode := record.GetFirstChild(); propNode != nil; propNode = propNode.GetNext() {
		// If it is a computed property then find the property name using the resolver
		if propNode.IsComputedProp() {
			compPropName := t.getComputedPropName(propNode)
			// If the name does not exist then report a warning
			if !nameResolver.nameVars.ContainsKey(compPropName) {
				t.reportWarning(record, TypeTransformation_UNKNOWN_NAMEVAR, compPropName)
				return t.getUnknownType()
			}
			// Otherwise add the property
			propValue := t.getComputedPropValue(propNode)
			resolvedName := nameResolver.nameVars.Get(compPropName)
			resultingType := t.evalInternal(propValue, nameResolver)
			props.Put(resolvedName, resultingType)
		} else {
			propName := propNode.GetString()
			resultingType := t.evalInternal(propNode.GetFirstChild(),
				nameResolver)
			props.Put(propName, resultingType)
		}
	}
	return t.registry.CreateRecordType(props)
}

// Java: TypeTransformation.evalRecordParam(Node, NameResolver)
func (t *TypeTransformation) evalRecordParam(ttlAst *rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	if ttlAst.IsObjectLit() {
		return t.evalRecord(ttlAst, nameResolver)
	}
	// The parameter of record can be a type transformation expression
	return t.evalInternal(ttlAst, nameResolver)
}

// Java: TypeTransformation.evalRecordType(Node, NameResolver)
func (t *TypeTransformation) evalRecordType(ttlAst *rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	paramCount := t.getCallParamCount(ttlAst)
	var recTypesBuilder []rhino.ObjectType
	for i := 0; i < paramCount; i++ {
		type_ := t.evalRecordParam(t.getCallArgument(ttlAst, i), nameResolver)
		// Check that each parameter evaluates to an object
		objType := type_.ToMaybeObjectType()
		if objType == nil || objType.IsUnknownType() {
			t.reportWarning(ttlAst, TypeTransformation_RECPARAM_INVALID, type_.String())
			return t.getUnknownType()
		}
		recType := t.registry.BuildRecordTypeFromObject(objType)
		if !recType.Equals(t.getObjectType()) {
			recTypesBuilder = append(recTypesBuilder, jutil.CheckNotNull(recType.ToMaybeObjectType()))
		}
	}
	return t.joinRecordTypes(recTypesBuilder)
}

// Java: TypeTransformation.putNewPropInPropertyMap(Map<String, JSType>, String, JSType)
func (t *TypeTransformation) putNewPropInPropertyMap(props *jutil.LinkedHashMap[string, rhino.JSType],
	newPropName string, newPropValue rhino.JSType) {
	// TODO(lpino): Decide if the best strategy is to collapse the properties
	// to a union type or not. So far, new values replace the old ones except
	// if they are two record types in which case the properties are joined
	// together

	// Three cases:
	// (i) If the key does not exist then add it to the map with the new value
	// (ii) If the key to be added already exists in the map and the new value
	// is not a record type then the current value is replaced with the new one
	// (iii) If the new value is a record type and the current is not then
	// the current value is replaced with the new one
	if !props.ContainsKey(newPropName) ||
		!newPropValue.IsRecordType() ||
		!props.Get(newPropName).IsRecordType() {
		props.Put(newPropName, newPropValue)
		return
	}
	// Otherwise join the current value with the new one since both are records
	props.Put(newPropName,
		t.joinRecordTypes([]rhino.ObjectType{
			jutil.CheckNotNull(typeTransformation_castToObjectType(props.Get(newPropName))),
			jutil.CheckNotNull(typeTransformation_castToObjectType(newPropValue))}))
}

// typeTransformation_castToObjectType is the Java cast (ObjectType) of a
// JSType: it throws ClassCastException when the type is not an object type.
//
// Go-only: language-adapter — Java casts the JSType reference; Go asserts
// the interface.
func typeTransformation_castToObjectType(t rhino.JSType) rhino.ObjectType {
	o, ok := t.(rhino.ObjectType)
	if !ok {
		panic(jutil.NewClassCastException(fmt.Sprintf("%T", t), "com.google.javascript.rhino.jstype.ObjectType"))
	}
	return o
}

// joinRecordTypes merges a list of record types. Example
// {r:{s:string, n:number}} and {a:boolean} is transformed into
// {r:{s:string, n:number}, a:boolean}
//
// Java: TypeTransformation.joinRecordTypes(ImmutableList<ObjectType>)
func (t *TypeTransformation) joinRecordTypes(recTypes []rhino.ObjectType) rhino.JSType {
	props := jutil.NewLinkedHashMap[string, rhino.JSType]()
	for _, recType := range recTypes {
		for _, newPropName := range recType.GetOwnPropertyNames() {
			newPropValue := recType.GetPropertyType(newPropName)
			// Put the new property depending if it already exists in the map
			t.putNewPropInPropertyMap(props, newPropName, newPropValue)
		}
	}
	return t.createRecordType(jutil.NewLinkedHashMapFrom(props))
}

// Java: TypeTransformation.evalMaprecord(Node, NameResolver)
func (t *TypeTransformation) evalMaprecord(ttlAst *rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	recordNode := ttlAst.GetSecondChild()
	mapFunction := ttlAst.GetChildAtIndex(2)
	type_ := t.evalInternal(recordNode, nameResolver)

	// If it is an empty record type (Object) then return
	if type_.Equals(t.getObjectType()) {
		return t.getObjectType()
	}

	// The parameter must be a valid record type
	if !type_.IsRecordType() {
		// TODO(lpino): Decide how to handle non-record types
		t.reportWarning(recordNode, TypeTransformation_RECTYPE_INVALID, type_.String())
		return t.getUnknownType()
	}

	objtype := type_.ToMaybeObjectType()
	// Fetch the information of the map function
	paramKey := t.getFunctionParameter(mapFunction, 0)
	paramValue := t.getFunctionParameter(mapFunction, 1)

	// The maprecord variables must not be defined in the environment
	if nameResolver.nameVars.ContainsKey(paramKey) {
		t.reportWarning(ttlAst, TypeTransformation_DUPLICATE_VARIABLE, paramKey)
		return t.getUnknownType()
	}
	if nameResolver.typeVars.ContainsKey(paramValue) {
		t.reportWarning(ttlAst, TypeTransformation_DUPLICATE_VARIABLE, paramValue)
		return t.getUnknownType()
	}

	// Compute the new properties using the map function
	mapFnBody := NodeUtil_GetFunctionBody(mapFunction)
	newProps := jutil.NewLinkedHashMap[string, rhino.JSType]()
	for _, propName := range objtype.GetOwnPropertyNames() {
		// The value of the current property
		propValue := objtype.GetPropertyType(propName)

		// Evaluate the map function body with paramValue and paramKey replaced
		// by the values of the current property
		newNameResolver := newTypeTransformation_NameResolver(
			typeTransformation_addNewEntry(nameResolver.typeVars, paramValue, propValue),
			typeTransformation_addNewEntry(nameResolver.nameVars, paramKey, propName))
		body := t.evalInternal(mapFnBody, newNameResolver)

		// If the body returns unknown then the whole expression returns unknown
		if body.IsUnknownType() {
			return t.getUnknownType()
		}

		// Skip the property when the body evaluates to NO_TYPE
		// or the empty record (Object)
		if body.IsEmptyType() || body.Equals(t.getObjectType()) {
			continue
		}

		// Otherwise the body must evaluate to a record type
		if !body.IsRecordType() {
			t.reportWarning(ttlAst, TypeTransformation_MAPRECORD_BODY_INVALID, body.String())
			return t.getUnknownType()
		}

		// Add the properties of the resulting record type to the original one
		bodyAsObj := body.ToMaybeObjectType()
		for _, newPropName := range bodyAsObj.GetOwnPropertyNames() {
			newPropValue := bodyAsObj.GetPropertyType(newPropName)
			// If the key already exists then we have to mix it with the current property value
			t.putNewPropInPropertyMap(newProps, newPropName, newPropValue)
		}
	}
	return t.createRecordType(jutil.NewLinkedHashMapFrom(newProps))
}

// Java: TypeTransformation.evalTypeOfVar(Node)
func (t *TypeTransformation) evalTypeOfVar(ttlAst *rhino.Node) rhino.JSType {
	name := t.getCallArgument(ttlAst, 0).GetString()
	slot, _ := t.typeEnv.GetSlot(name).(rhino.StaticTypedSlot)
	var type_ rhino.JSType
	if slot != nil {
		type_ = slot.GetType()
	}
	if type_ == nil {
		t.reportWarning(ttlAst, TypeTransformation_VAR_UNDEFINED, name)
		return t.getUnknownType()
	}
	return type_
}

// Java: TypeTransformation.evalInstanceOf(Node, NameResolver)
func (t *TypeTransformation) evalInstanceOf(ttlAst *rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	type_ := t.evalInternal(t.getCallArgument(ttlAst, 0), nameResolver)
	if type_.IsUnknownType() || !type_.IsConstructor() {
		t.reportWarning(ttlAst, TypeTransformation_INVALID_CTOR, jutil.ToJavaString(type_.GetDisplayName()))
		return t.getUnknownType()
	}
	return type_.ToMaybeFunctionType().GetInstanceType()
}

// Java: TypeTransformation.evalNativeTypeExpr(Node)
func (t *TypeTransformation) evalNativeTypeExpr(ttlAst *rhino.Node) rhino.JSType {
	expr :=
		rhino.NewJSTypeExpression(t.getCallArgument(ttlAst, 0).CloneTree(), typeTransformation_VIRTUAL_FILE)
	return t.registry.EvaluateTypeExpression(expr, t.typeEnv)
}

// Java: TypeTransformation.evalPrintType(Node, NameResolver)
func (t *TypeTransformation) evalPrintType(ttlAst *rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	type_ := t.evalInternal(t.getCallArgument(ttlAst, 1), nameResolver)
	msg := t.getCallArgument(ttlAst, 0).GetString() + type_.String()
	jutil.System_Out().Println(msg)
	return type_
}

// Java: TypeTransformation.evalPropType(Node, NameResolver)
func (t *TypeTransformation) evalPropType(ttlAst *rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	type_ := t.evalInternal(t.getCallArgument(ttlAst, 1), nameResolver)
	objType := type_.ToMaybeObjectType()
	if objType == nil {
		t.reportWarning(ttlAst, TypeTransformation_PROPTYPE_INVALID, type_.String())
		return t.getUnknownType()
	}
	propType := objType.GetPropertyType(t.getCallArgument(ttlAst, 0).GetString())
	// firstNonNull(propType, getUnknownType())
	if propType != nil {
		return propType
	}
	return t.getUnknownType()
}
