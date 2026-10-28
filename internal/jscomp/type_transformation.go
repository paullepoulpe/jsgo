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

	"github.com/paullepoulpe/jsgo/internal/ast"
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

type ttlTypeVar struct {
	name string
	typ  rhino.JSType
}

type ttlNameVar struct {
	name string
	val  string
}

type ttlRecordProp struct {
	name string
	typ  rhino.JSType
}

// typeTransformation_NameResolver holds the type variable and name variable
// bindings during TTL evaluation.
//
// Java: private static class TypeTransformation.NameResolver
type typeTransformation_NameResolver struct {
	typeVars []ttlTypeVar
	nameVars []ttlNameVar
}

func (r *typeTransformation_NameResolver) lookupTypeVar(name string) (rhino.JSType, bool) {
	for i := len(r.typeVars) - 1; i >= 0; i-- {
		if r.typeVars[i].name == name {
			return r.typeVars[i].typ, true
		}
	}
	return nil, false
}

func (r *typeTransformation_NameResolver) lookupNameVar(name string) (string, bool) {
	for i := len(r.nameVars) - 1; i >= 0; i-- {
		if r.nameVars[i].name == name {
			return r.nameVars[i].val, true
		}
	}
	return "", false
}

// Java: TypeTransformation(AbstractCompiler, StaticTypedScope)
func NewTypeTransformation(compiler AbstractCompiler, typeEnv rhino.StaticTypedScope) *TypeTransformation {
	return &TypeTransformation{
		compiler: compiler,
		registry: compiler.GetTypeRegistry(),
		typeEnv:  typeEnv,
	}
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

// Java: TypeTransformation.reportWarning(Node, DiagnosticType, String...)
func (t *TypeTransformation) reportWarning(n rhino.Node, msg *DiagnosticType, param ...string) {
	t.compiler.Report(JSError_MakeWithNode(n, msg, param...))
}

// Java: TypeTransformation.getFunctionParameter(Node, int)
func (t *TypeTransformation) getFunctionParameter(n rhino.Node, i int) string {
	jutil.CheckArgumentMsg1(n.IsFunction(), "Expected a function node, found %s", n)
	param, _ := ast.MustFunction(n).Params().ParamAt(i)
	return param.Node().GetString()
}

// Java: TypeTransformation.getCallName(Node)
func (t *TypeTransformation) getCallName(n rhino.Node) string {
	jutil.CheckArgumentMsg1(n.IsCall(), "Expected a call node, found %s", n)
	return ast.MustCall(n).Callee().Node().GetString()
}

// Java: TypeTransformation.getCallArgument(Node, int)
func (t *TypeTransformation) getCallArgument(n rhino.Node, i int) rhino.Node {
	jutil.CheckArgumentMsg1(n.IsCall(), "Expected a call node, found %s", n)
	arg, _ := ast.MustCall(n).AsCallOrNew().ArgAt(i)
	return arg.Node()
}

// Java: TypeTransformation.getComputedPropValue(Node)
func (t *TypeTransformation) getComputedPropValue(n rhino.Node) rhino.Node {
	jutil.CheckArgumentMsg1(
		n.IsComputedProp(), "Expected a computed property node, found %s", n)
	val, _ := ast.MustComputedProp(n).Value()
	return val.Node()
}

// Java: TypeTransformation.getComputedPropName(Node)
func (t *TypeTransformation) getComputedPropName(n rhino.Node) string {
	jutil.CheckArgumentMsg1(
		n.IsComputedProp(), "Expected a computed property node, found %s", n)
	return ast.MustComputedProp(n).Key().Node().GetString()
}

// Eval evaluates the type transformation expression and returns the
// resulting type.
//
// ttlAst is the node representing the type transformation expression;
// typeVars is the environment containing the information about the type
// variables. Returns the resulting type after the transformation
//
// Java: TypeTransformation.eval(Node, ImmutableMap<String, JSType>)
func (t *TypeTransformation) Eval(ttlAst rhino.Node, typeVars *jutil.LinkedHashMap[string, rhino.JSType]) rhino.JSType {
	return t.EvalWithNameVars(ttlAst, typeVars, nil)
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
	ttlAst rhino.Node, typeVars *jutil.LinkedHashMap[string, rhino.JSType], nameVars *jutil.LinkedHashMap[string, string]) rhino.JSType {
	var tvBuf [8]ttlTypeVar
	tv := tvBuf[:0]
	if typeVars != nil {
		for k, v := range typeVars.All() {
			tv = append(tv, ttlTypeVar{name: k, typ: v})
		}
	}
	var nvBuf [4]ttlNameVar
	nv := nvBuf[:0]
	if nameVars != nil {
		for k, v := range nameVars.All() {
			nv = append(nv, ttlNameVar{name: k, val: v})
		}
	}
	return t.evalWithSlices(ttlAst, tv, nv)
}

func (t *TypeTransformation) evalWithSlices(ttlAst rhino.Node, typeVars []ttlTypeVar, nameVars []ttlNameVar) rhino.JSType {
	resolver := typeTransformation_NameResolver{typeVars: typeVars, nameVars: nameVars}
	result := t.evalInternal(ttlAst, &resolver)
	if result.IsEmptyType() {
		return t.getUnknownType()
	}
	return result
}

// Java: TypeTransformation.evalInternal(Node, NameResolver)
func (t *TypeTransformation) evalInternal(ttlAst rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	if ttlAst.IsStringLit() {
		return t.evalTypeName(ttlAst)
	}
	if ttlAst.IsName() {
		return t.evalTypeVar(ttlAst, nameResolver)
	}
	keyword := parsing.TypeTransformationParser_Keywords_FromCaseInsensitiveASCII(t.getCallName(ttlAst))
	switch keyword.Kind() {
	case parsing.TypeTransformationParser_OperationKind_TYPE_CONSTRUCTOR:
		return t.evalTypeExpressionWithKeyword(ttlAst, keyword, nameResolver)
	case parsing.TypeTransformationParser_OperationKind_OPERATION:
		return t.evalOperationExpressionWithKeyword(ttlAst, keyword, nameResolver)
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Could not evaluate the type transformation expression"))
	}
}

func (t *TypeTransformation) evalOperationExpressionWithKeyword(ttlAst rhino.Node, keyword parsing.TypeTransformationParser_Keywords, nameResolver *typeTransformation_NameResolver) rhino.JSType {
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

func (t *TypeTransformation) evalTypeExpressionWithKeyword(ttlAst rhino.Node, keyword parsing.TypeTransformationParser_Keywords, nameResolver *typeTransformation_NameResolver) rhino.JSType {
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
func (t *TypeTransformation) evalTypeName(ttlAst rhino.Node) rhino.JSType {
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
func (t *TypeTransformation) evalTemplatizedType(ttlAst rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	call := ast.MustCall(ttlAst)
	firstArg, _ := call.FirstArg()
	firstParam := t.evalInternal(firstArg.Node(), nameResolver)
	if !firstParam.IsRawTypeOfTemplatizedType() {
		t.reportWarning(ttlAst, TypeTransformation_BASETYPE_INVALID, firstParam.String())
		return t.getUnknownType()
	}
	// TODO(lpino): Check that the number of parameters correspond with the
	// number of template types that the base type can take when creating
	// a templatized type. For instance, if the base type is Array then there
	// must be just one parameter.
	restCount := call.NumArgs() - 1
	templatizedTypes := make([]rhino.JSType, 0, restCount)
	for arg := firstArg.Node().GetNext(); !arg.IsNil(); arg = arg.GetNext() {
		templatizedTypes = append(templatizedTypes, t.evalInternal(arg, nameResolver))
	}
	baseType := firstParam.ToMaybeObjectType()
	return t.registry.CreateTemplatizedType(baseType, templatizedTypes...)
}

// Java: TypeTransformation.evalTypeVar(Node, NameResolver)
func (t *TypeTransformation) evalTypeVar(ttlAst rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	typeVar := ttlAst.GetString()
	resultingType, ok := nameResolver.lookupTypeVar(typeVar)
	// If the type variable is not defined then return UNKNOWN and report a warning
	if !ok || resultingType == nil {
		t.reportWarning(ttlAst, TypeTransformation_UNKNOWN_TYPEVAR, typeVar)
		return t.getUnknownType()
	}
	return resultingType
}

// Java: TypeTransformation.evalUnionType(Node, NameResolver)
func (t *TypeTransformation) evalUnionType(ttlAst rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	jutil.CheckArgumentMsg1(ttlAst.IsCall(), "Expected a call node, found %s", ttlAst)
	var buf [8]rhino.JSType
	variants := buf[:0]
	for arg := ttlAst.GetSecondChild(); !arg.IsNil(); arg = arg.GetNext() {
		variants = append(variants, t.evalInternal(arg, nameResolver))
	}
	return t.registry.CreateUnionType(variants...)
}

// Java: TypeTransformation.evalString(Node, NameResolver)
func (t *TypeTransformation) evalString(ttlAst rhino.Node, nameResolver *typeTransformation_NameResolver) string {
	s := ttlAst.GetString()
	if ttlAst.IsName() {
		// Return the empty string if the name variable cannot be resolved
		val, ok := nameResolver.lookupNameVar(s)
		if !ok {
			t.reportWarning(ttlAst, TypeTransformation_UNKNOWN_STRVAR, s)
			return ""
		}
		return val
	}
	return s
}

func (t *TypeTransformation) evalTypePredicateWithKeyword(ttlAst rhino.Node, keyword parsing.TypeTransformationParser_Keywords, nameResolver *typeTransformation_NameResolver) bool {
	arg0 := ttlAst.GetSecondChild()
	p0 := t.evalInternal(arg0, nameResolver)
	switch keyword {
	case parsing.TypeTransformationParser_Keywords_EQ:
		return p0.Equals(t.evalInternal(arg0.GetNext(), nameResolver))
	case parsing.TypeTransformationParser_Keywords_SUB:
		return p0.IsSubtypeOf(t.evalInternal(arg0.GetNext(), nameResolver))
	case parsing.TypeTransformationParser_Keywords_ISCTOR:
		return p0.IsConstructor()
	case parsing.TypeTransformationParser_Keywords_ISTEMPLATIZED:
		return p0.IsTemplatizedType()
	case parsing.TypeTransformationParser_Keywords_ISRECORD:
		return p0.IsRecordType()
	case parsing.TypeTransformationParser_Keywords_ISUNKNOWN:
		return p0.IsUnknownType()
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Invalid type predicate in the type transformation"))
	}
}

func (t *TypeTransformation) evalStringPredicateWithKeyword(ttlAst rhino.Node, keyword parsing.TypeTransformationParser_Keywords, nameResolver *typeTransformation_NameResolver) bool {
	arg0 := ttlAst.GetSecondChild()
	p0 := t.evalString(arg0, nameResolver)
	p1 := t.evalString(arg0.GetNext(), nameResolver)
	if p0 == "" || p1 == "" {
		return false
	}
	switch keyword {
	case parsing.TypeTransformationParser_Keywords_STREQ:
		return p0 == p1
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Invalid string predicate in the type transformation"))
	}
}

func (t *TypeTransformation) evalTypevarPredicateWithKeyword(ttlAst rhino.Node, keyword parsing.TypeTransformationParser_Keywords, nameResolver *typeTransformation_NameResolver) bool {
	switch keyword {
	case parsing.TypeTransformationParser_Keywords_ISDEFINED:
		_, ok := nameResolver.lookupTypeVar(t.getCallArgument(ttlAst, 0).GetString())
		return ok
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Invalid typevar predicate in the type transformation"))
	}
}

// Java: TypeTransformation.evalBooleanOperation(Node, NameResolver)
func (t *TypeTransformation) evalBooleanOperation(ttlAst rhino.Node, nameResolver *typeTransformation_NameResolver) bool {
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
func (t *TypeTransformation) evalBoolean(ttlAst rhino.Node, nameResolver *typeTransformation_NameResolver) bool {
	if ttlAst.IsAnd() || ttlAst.IsOr() || ttlAst.IsNot() {
		return t.evalBooleanOperation(ttlAst, nameResolver)
	}
	keyword := parsing.TypeTransformationParser_Keywords_FromCaseInsensitiveASCII(t.getCallName(ttlAst))
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
func (t *TypeTransformation) evalConditional(ttlAst rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	jutil.CheckArgumentMsg1(ttlAst.IsCall(), "Expected a call node, found %s", ttlAst)
	p0 := ttlAst.GetSecondChild()
	p1 := p0.GetNext()
	if t.evalBoolean(p0, nameResolver) {
		return t.evalInternal(p1, nameResolver)
	}
	return t.evalInternal(p1.GetNext(), nameResolver)
}

// Java: TypeTransformation.evalMapunion(Node, NameResolver)
func (t *TypeTransformation) evalMapunion(ttlAst rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	jutil.CheckArgumentMsg1(ttlAst.IsCall(), "Expected a call node, found %s", ttlAst)
	unionParam := ttlAst.GetSecondChild()
	mapFunction := unionParam.GetNext()
	paramName := t.getFunctionParameter(mapFunction, 0)

	// The mapunion variable must not be defined in the environment
	if _, exists := nameResolver.lookupTypeVar(paramName); exists {
		t.reportWarning(ttlAst, TypeTransformation_DUPLICATE_VARIABLE, paramName)
		return t.getUnknownType()
	}

	mapFunctionBody := NodeUtil_GetFunctionBody(mapFunction)
	unionType := t.evalInternal(unionParam, nameResolver)

	var tvBuf [8]ttlTypeVar
	var extTypeVars []ttlTypeVar
	if len(nameResolver.typeVars)+1 <= len(tvBuf) {
		extTypeVars = append(tvBuf[:0], nameResolver.typeVars...)
	} else {
		extTypeVars = make([]ttlTypeVar, len(nameResolver.typeVars), len(nameResolver.typeVars)+1)
		copy(extTypeVars, nameResolver.typeVars)
	}
	extTypeVars = append(extTypeVars, ttlTypeVar{name: paramName})
	lastIdx := len(extTypeVars) - 1
	childResolver := typeTransformation_NameResolver{
		typeVars: extTypeVars,
		nameVars: nameResolver.nameVars,
	}

	// If the first parameter does not correspond to a union type then
	// consider it as a union with a single type and evaluate
	if !unionType.IsUnionType() {
		extTypeVars[lastIdx].typ = unionType
		return t.evalInternal(mapFunctionBody, &childResolver)
	}

	// Otherwise obtain the elements in the union type. Note that the block
	// above guarantees the casting to be safe
	unionElms := unionType.GetUnionMembers()
	var outBuf [8]rhino.JSType
	var newUnionElms []rhino.JSType
	if len(unionElms) <= len(outBuf) {
		newUnionElms = outBuf[:len(unionElms)]
	} else {
		newUnionElms = make([]rhino.JSType, len(unionElms))
	}
	for i, elm := range unionElms {
		extTypeVars[lastIdx].typ = elm
		newUnionElms[i] = t.evalInternal(mapFunctionBody, &childResolver)
	}

	return t.registry.CreateUnionType(newUnionElms...)
}

// Java: TypeTransformation.evalRawTypeOf(Node, NameResolver)
func (t *TypeTransformation) evalRawTypeOf(ttlAst rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	type_ := t.evalInternal(t.getCallArgument(ttlAst, 0), nameResolver)
	if !type_.IsTemplatizedType() {
		t.reportWarning(ttlAst, TypeTransformation_TEMPTYPE_INVALID, "rawTypeOf", type_.String())
		return t.getUnknownType()
	}
	return type_.ToMaybeObjectType().GetRawType()
}

// Java: TypeTransformation.evalTemplateTypeOf(Node, NameResolver)
func (t *TypeTransformation) evalTemplateTypeOf(ttlAst rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	jutil.CheckArgumentMsg1(ttlAst.IsCall(), "Expected a call node, found %s", ttlAst)
	p0 := ttlAst.GetSecondChild()
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
func (t *TypeTransformation) evalRecord(record rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	builder := rhino.NewRecordTypeBuilder(t.registry)
	for member := range ast.MustObjectLit(record).Members() {
		propNode := member.Node()
		// If it is a computed property then find the property name using the resolver
		if propNode.IsComputedProp() {
			compPropName := t.getComputedPropName(propNode)
			resolvedName, ok := nameResolver.lookupNameVar(compPropName)
			if !ok {
				t.reportWarning(record, TypeTransformation_UNKNOWN_NAMEVAR, compPropName)
				return t.getUnknownType()
			}
			propValue := t.getComputedPropValue(propNode)
			resultingType := t.evalInternal(propValue, nameResolver)
			builder.AddProperty(resolvedName, resultingType, rhino.NilNode)
		} else {
			key := ast.MustStringKey(propNode)
			resultingType := t.evalInternal(key.Value().Node(), nameResolver)
			builder.AddProperty(key.Key(), resultingType, rhino.NilNode)
		}
	}
	return builder.Build()
}

// Java: TypeTransformation.evalRecordParam(Node, NameResolver)
func (t *TypeTransformation) evalRecordParam(ttlAst rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	if ttlAst.IsObjectLit() {
		return t.evalRecord(ttlAst, nameResolver)
	}
	// The parameter of record can be a type transformation expression
	return t.evalInternal(ttlAst, nameResolver)
}

// Java: TypeTransformation.evalRecordType(Node, NameResolver)
func (t *TypeTransformation) evalRecordType(ttlAst rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	jutil.CheckArgumentMsg1(ttlAst.IsCall(), "Expected a call node, found %s", ttlAst)
	var recBuf [4]rhino.ObjectType
	recTypesBuilder := recBuf[:0]
	for arg := ttlAst.GetSecondChild(); !arg.IsNil(); arg = arg.GetNext() {
		type_ := t.evalRecordParam(arg, nameResolver)
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

func (t *TypeTransformation) putNewPropInPropertyMap(
	props []ttlRecordProp, newPropName string, newPropValue rhino.JSType) []ttlRecordProp {
	for i := range props {
		if props[i].name == newPropName {
			if !newPropValue.IsRecordType() || !props[i].typ.IsRecordType() {
				props[i].typ = newPropValue
				return props
			}
			pair := [2]rhino.ObjectType{
				jutil.CheckNotNull(typeTransformation_castToObjectType(props[i].typ)),
				jutil.CheckNotNull(typeTransformation_castToObjectType(newPropValue)),
			}
			props[i].typ = t.joinRecordTypes(pair[:])
			return props
		}
	}
	return append(props, ttlRecordProp{name: newPropName, typ: newPropValue})
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

func (t *TypeTransformation) buildRecordFromProps(props []ttlRecordProp) rhino.JSType {
	builder := rhino.NewRecordTypeBuilder(t.registry)
	for _, p := range props {
		builder.AddProperty(p.name, p.typ, rhino.NilNode)
	}
	return builder.Build()
}

// joinRecordTypes merges a list of record types. Example
// {r:{s:string, n:number}} and {a:boolean} is transformed into
// {r:{s:string, n:number}, a:boolean}
//
// Java: TypeTransformation.joinRecordTypes(ImmutableList<ObjectType>)
func (t *TypeTransformation) joinRecordTypes(recTypes []rhino.ObjectType) rhino.JSType {
	var propBuf [8]ttlRecordProp
	props := propBuf[:0]
	for _, recType := range recTypes {
		for _, newPropName := range recType.GetOwnPropertyNames() {
			newPropValue := recType.GetPropertyType(newPropName)
			props = t.putNewPropInPropertyMap(props, newPropName, newPropValue)
		}
	}
	return t.buildRecordFromProps(props)
}

// Java: TypeTransformation.evalMaprecord(Node, NameResolver)
func (t *TypeTransformation) evalMaprecord(ttlAst rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	jutil.CheckArgumentMsg1(ttlAst.IsCall(), "Expected a call node, found %s", ttlAst)
	recordNode := ttlAst.GetSecondChild()
	mapFunction := recordNode.GetNext()
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
	if _, exists := nameResolver.lookupNameVar(paramKey); exists {
		t.reportWarning(ttlAst, TypeTransformation_DUPLICATE_VARIABLE, paramKey)
		return t.getUnknownType()
	}
	if _, exists := nameResolver.lookupTypeVar(paramValue); exists {
		t.reportWarning(ttlAst, TypeTransformation_DUPLICATE_VARIABLE, paramValue)
		return t.getUnknownType()
	}

	var tvBuf [8]ttlTypeVar
	var extTypeVars []ttlTypeVar
	if len(nameResolver.typeVars)+1 <= len(tvBuf) {
		extTypeVars = append(tvBuf[:0], nameResolver.typeVars...)
	} else {
		extTypeVars = make([]ttlTypeVar, len(nameResolver.typeVars), len(nameResolver.typeVars)+1)
		copy(extTypeVars, nameResolver.typeVars)
	}
	extTypeVars = append(extTypeVars, ttlTypeVar{name: paramValue})
	tvLast := len(extTypeVars) - 1

	var nvBuf [4]ttlNameVar
	var extNameVars []ttlNameVar
	if len(nameResolver.nameVars)+1 <= len(nvBuf) {
		extNameVars = append(nvBuf[:0], nameResolver.nameVars...)
	} else {
		extNameVars = make([]ttlNameVar, len(nameResolver.nameVars), len(nameResolver.nameVars)+1)
		copy(extNameVars, nameResolver.nameVars)
	}
	extNameVars = append(extNameVars, ttlNameVar{name: paramKey})
	nvLast := len(extNameVars) - 1

	childResolver := typeTransformation_NameResolver{
		typeVars: extTypeVars,
		nameVars: extNameVars,
	}

	// Compute the new properties using the map function
	mapFnBody := NodeUtil_GetFunctionBody(mapFunction)
	var propBuf [8]ttlRecordProp
	newProps := propBuf[:0]
	for _, propName := range objtype.GetOwnPropertyNames() {
		// The value of the current property
		propValue := objtype.GetPropertyType(propName)

		// Evaluate the map function body with paramValue and paramKey replaced
		// by the values of the current property
		extTypeVars[tvLast].typ = propValue
		extNameVars[nvLast].val = propName
		body := t.evalInternal(mapFnBody, &childResolver)

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
			newProps = t.putNewPropInPropertyMap(newProps, newPropName, newPropValue)
		}
	}
	return t.buildRecordFromProps(newProps)
}

// Java: TypeTransformation.evalTypeOfVar(Node)
func (t *TypeTransformation) evalTypeOfVar(ttlAst rhino.Node) rhino.JSType {
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
func (t *TypeTransformation) evalInstanceOf(ttlAst rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	type_ := t.evalInternal(t.getCallArgument(ttlAst, 0), nameResolver)
	if type_.IsUnknownType() || !type_.IsConstructor() {
		t.reportWarning(ttlAst, TypeTransformation_INVALID_CTOR, jutil.ToJavaString(type_.GetDisplayName()))
		return t.getUnknownType()
	}
	return type_.ToMaybeFunctionType().GetInstanceType()
}

// Java: TypeTransformation.evalNativeTypeExpr(Node)
func (t *TypeTransformation) evalNativeTypeExpr(ttlAst rhino.Node) rhino.JSType {
	expr :=
		rhino.NewJSTypeExpression(t.getCallArgument(ttlAst, 0).CloneTree(), typeTransformation_VIRTUAL_FILE)
	return t.registry.EvaluateTypeExpression(expr, t.typeEnv)
}

// Java: TypeTransformation.evalPrintType(Node, NameResolver)
func (t *TypeTransformation) evalPrintType(ttlAst rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
	type_ := t.evalInternal(t.getCallArgument(ttlAst, 1), nameResolver)
	msg := t.getCallArgument(ttlAst, 0).GetString() + type_.String()
	jutil.System_Out().Println(msg)
	return type_
}

// Java: TypeTransformation.evalPropType(Node, NameResolver)
func (t *TypeTransformation) evalPropType(ttlAst rhino.Node, nameResolver *typeTransformation_NameResolver) rhino.JSType {
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
