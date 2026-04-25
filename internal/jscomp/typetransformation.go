// Port of com/google/javascript/jscomp/TypeTransformation.java
// @ closure-compiler v20260415 (subset).
//
// TypeTransformation evaluates @template TTL expressions of the form
//   @template T := cond(U, T, R) =:
// at type-check time. The Java class is wired into
// TypedScopeCreator.processTemplateType and consumes a full
// JSTypeRegistry + StaticTypedScope + AbstractCompiler. This port
// uses the Go analogues (jstype.Registry, TypedScope, diag.Diagnostic)
// and is additive — it is not yet wired into the pipeline.
//
// Portions of the Java logic that depend on APIs not yet ported
// (evaluateTypeExpression, buildRecordTypeFromObject,
// createTemplatizedType) carry // TODO comments with the Java
// source location they correspond to.

package jscomp

import (
	"fmt"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// Diagnostic codes registered by TypeTransformation. Names match the
// Java DiagnosticType identifiers so corpus messages line up.
const (
	// TYPEVAR_UNDEFINED — "Reference to an unknown type variable {0}"
	ttlUnknownTypevar = "TYPEVAR_UNDEFINED"
	// UNKNOWN_STRVAR — "Reference to an unknown string variable {0}"
	ttlUnknownStrvar = "UNKNOWN_STRVAR"
	// TYPENAME_UNDEFINED — "Reference to an unknown type name {0}"
	ttlUnknownTypename = "TYPENAME_UNDEFINED"
	// BASETYPE_INVALID — "The type {0} cannot be templatized"
	ttlBasetypeInvalid = "BASETYPE_INVALID"
	// TEMPTYPE_INVALID — "Expected templatized type in {0} found {1}"
	ttlTemptypeInvalid = "TEMPTYPE_INVALID"
	// INDEX_OUTOFBOUNDS — "Index out of bounds in templateTypeOf…"
	ttlIndexOutofbounds = "INDEX_OUTOFBOUNDS"
	// DUPLICATE_VARIABLE — "The variable {0} is already defined"
	ttlDuplicateVariable = "DUPLICATE_VARIABLE"
	// UNKNOWN_NAMEVAR — "Reference to an unknown name variable {0}"
	ttlUnknownNamevar = "UNKNOWN_NAMEVAR"
	// RECTYPE_INVALID — "The first parameter of a maprecord must be a record type…"
	ttlRectypeInvalid = "RECTYPE_INVALID"
	// MAPRECORD_BODY_INVALID — "The body of a maprecord function must evaluate to a record type or a no type…"
	ttlMaprecordBodyInvalid = "MAPRECORD_BODY_INVALID"
	// VAR_UNDEFINED — "Variable {0} is undefined in the scope"
	ttlVarUndefined = "VAR_UNDEFINED"
	// INVALID_CTOR — "Expected a constructor type, found {0}"
	ttlInvalidCtor = "INVALID_CTOR"
	// RECPARAM_INVALID — "Expected a record type, found {0}"
	ttlRecparamInvalid = "RECPARAM_INVALID"
	// PROPTYPE_INVALID — "Expected object type, found {0}"
	ttlProptypeInvalid = "PROPTYPE_INVALID"
)

// OperationKind mirrors TypeTransformationParser.OperationKind.
type OperationKind int

const (
	OpKindTypeConstructor OperationKind = iota
	OpKindOperation
	OpKindStringPredicate
	OpKindTypePredicate
	OpKindTypevarPredicate
)

// TTLKeyword mirrors TypeTransformationParser.Keywords. Each keyword
// carries its source-level name, its OperationKind, and the allowed
// parameter count range (VAR_ARGS = math.MaxInt32).
type TTLKeyword struct {
	name          string
	minParamCount int
	maxParamCount int
	kind          OperationKind
}

const ttlVarArgs = 1<<31 - 1 // Integer.MAX_VALUE

// keywordAll / keywordCond / … mirror the Keywords enum values.
var (
	kwALL            = &TTLKeyword{"all", 0, 0, OpKindTypeConstructor}
	kwCOND           = &TTLKeyword{"cond", 3, 3, OpKindOperation}
	kwEQ             = &TTLKeyword{"eq", 2, 2, OpKindTypePredicate}
	kwISCTOR         = &TTLKeyword{"isCtor", 1, 1, OpKindTypePredicate}
	kwISDEFINED      = &TTLKeyword{"isDefined", 1, 1, OpKindTypevarPredicate}
	kwISRECORD       = &TTLKeyword{"isRecord", 1, 1, OpKindTypePredicate}
	kwISTEMPLATIZED  = &TTLKeyword{"isTemplatized", 1, 1, OpKindTypePredicate}
	kwISUNKNOWN      = &TTLKeyword{"isUnknown", 1, 1, OpKindTypePredicate}
	kwINSTANCEOF     = &TTLKeyword{"instanceOf", 1, 1, OpKindOperation}
	kwMAPUNION       = &TTLKeyword{"mapunion", 2, 2, OpKindOperation}
	kwMAPRECORD      = &TTLKeyword{"maprecord", 2, 2, OpKindOperation}
	kwNONE           = &TTLKeyword{"none", 0, 0, OpKindTypeConstructor}
	kwPRINTTYPE      = &TTLKeyword{"printType", 2, 2, OpKindOperation}
	kwPROPTYPE       = &TTLKeyword{"propType", 2, 2, OpKindOperation}
	kwRAWTYPEOF      = &TTLKeyword{"rawTypeOf", 1, 1, OpKindTypeConstructor}
	kwSUB            = &TTLKeyword{"sub", 2, 2, OpKindTypePredicate}
	kwSTREQ          = &TTLKeyword{"streq", 2, 2, OpKindStringPredicate}
	kwRECORD         = &TTLKeyword{"record", 1, ttlVarArgs, OpKindTypeConstructor}
	kwTEMPLATETYPEOF = &TTLKeyword{"templateTypeOf", 2, 2, OpKindTypeConstructor}
	kwTYPE           = &TTLKeyword{"type", 2, ttlVarArgs, OpKindTypeConstructor}
	kwTYPEEXPR       = &TTLKeyword{"typeExpr", 1, 1, OpKindTypeConstructor}
	kwTYPEOFVAR      = &TTLKeyword{"typeOfVar", 1, 1, OpKindOperation}
	kwUNION          = &TTLKeyword{"union", 2, ttlVarArgs, OpKindTypeConstructor}
	kwUNKNOWN        = &TTLKeyword{"unknown", 0, 0, OpKindTypeConstructor}

	// ttlKeywordByName maps the lower-case keyword name used in TTL
	// source to its descriptor. Matches Keywords.valueOf(Ascii.toUpperCase(s))
	// in Java, which resolves by enum name (uppercase).
	ttlKeywordByName = map[string]*TTLKeyword{
		"all":            kwALL,
		"cond":           kwCOND,
		"eq":             kwEQ,
		"isctor":         kwISCTOR,
		"isdefined":      kwISDEFINED,
		"isrecord":       kwISRECORD,
		"istemplatized":  kwISTEMPLATIZED,
		"isunknown":      kwISUNKNOWN,
		"instanceof":     kwINSTANCEOF,
		"mapunion":       kwMAPUNION,
		"maprecord":      kwMAPRECORD,
		"none":           kwNONE,
		"printtype":      kwPRINTTYPE,
		"proptype":       kwPROPTYPE,
		"rawtypeof":      kwRAWTYPEOF,
		"sub":            kwSUB,
		"streq":          kwSTREQ,
		"record":         kwRECORD,
		"templatetypeof": kwTEMPLATETYPEOF,
		"type":           kwTYPE,
		"typeexpr":       kwTYPEEXPR,
		"typeofvar":      kwTYPEOFVAR,
		"union":          kwUNION,
		"unknown":        kwUNKNOWN,
	}
)

// nameToKeyword mirrors TypeTransformation.nameToKeyword —
// Ascii.toUpperCase(s) → Keywords.valueOf lookup.
func nameToKeyword(s string) *TTLKeyword {
	kw := ttlKeywordByName[strings.ToLower(s)]
	return kw
}

// NameResolver holds the variable bindings in scope during TTL
// evaluation. Mirrors TypeTransformation.NameResolver.
type NameResolver struct {
	// typeVars maps @template type-variable names to their bound JSType.
	typeVars map[string]jstype.JSType
	// nameVars maps string-variable names (the key param of maprecord)
	// to their string values.
	nameVars map[string]string
}

// newNameResolver constructs a NameResolver from immutable maps,
// mirroring Java's `new NameResolver(typeVars, nameVars)`.
func newNameResolver(typeVars map[string]jstype.JSType, nameVars map[string]string) *NameResolver {
	return &NameResolver{typeVars: typeVars, nameVars: nameVars}
}

// addTypeVar returns a new NameResolver with name→typ added.
// Mirrors `addNewEntry(nameResolver.typeVars, paramName, unionType)`.
func (r *NameResolver) addTypeVar(name string, typ jstype.JSType) *NameResolver {
	next := make(map[string]jstype.JSType, len(r.typeVars)+1)
	for k, v := range r.typeVars {
		next[k] = v
	}
	next[name] = typ
	return &NameResolver{typeVars: next, nameVars: r.nameVars}
}

// addNameVar returns a new NameResolver with name→val added.
func (r *NameResolver) addNameVar(name, val string) *NameResolver {
	next := make(map[string]string, len(r.nameVars)+1)
	for k, v := range r.nameVars {
		next[k] = v
	}
	next[name] = val
	return &NameResolver{typeVars: r.typeVars, nameVars: next}
}

// TypeTransformation evaluates TTL type transformation expressions.
// Mirrors com/google/javascript/jscomp/TypeTransformation.java.
//
// Not yet wired into the pipeline — see compiler.go for the pass
// registry. The struct is exported so future callers (TypedScopeCreator,
// TypeCheck) can construct it directly.
type TypeTransformation struct {
	// path is the source file name used for diagnostic positions.
	path string
	// typeEnv is the global scope in which type names are resolved.
	typeEnv *TypedScope
	// diags accumulates warnings produced during evaluation.
	diags []diag.Diagnostic
}

// NewTypeTransformation mirrors TypeTransformation(AbstractCompiler, StaticTypedScope).
// path plays the role of the compiler/source-file for diagnostic positions.
func NewTypeTransformation(path string, typeEnv *TypedScope) *TypeTransformation {
	return &TypeTransformation{path: path, typeEnv: typeEnv}
}

// ---- node predicate helpers (mirror Java private helpers) -----------

// isTypeVar mirrors TypeTransformation.isTypeVar — a NAME node
// represents a type variable in TTL expressions.
func (t *TypeTransformation) isTypeVar(n *rhino.Node) bool {
	return n.Token() == rhino.NAME
}

// isTypeName mirrors TypeTransformation.isTypeName — a STRINGLIT node
// represents a type name reference.
func (t *TypeTransformation) isTypeName(n *rhino.Node) bool {
	return n.Token() == rhino.STRINGLIT
}

// isBooleanOperation mirrors TypeTransformation.isBooleanOperation.
func (t *TypeTransformation) isBooleanOperation(n *rhino.Node) bool {
	tok := n.Token()
	return tok == rhino.AND || tok == rhino.OR || tok == rhino.NOT
}

// ---- AST accessor helpers -----------------------------------------

// getFunctionParameter mirrors TypeTransformation.getFunctionParameter.
// Returns the name of the i-th parameter of FUNCTION node n.
func getFunctionParameter(n *rhino.Node, i int) string {
	// n is FUNCTION; children[1] is PARAM_LIST.
	paramList := n.Children()[1]
	return paramList.Children()[i].String_()
}

// getCallName mirrors TypeTransformation.getCallName.
func getCallName(n *rhino.Node) string {
	return n.FirstChild().String_()
}

// getCallArgument mirrors TypeTransformation.getCallArgument — 0-indexed
// past the callee child.
func getCallArgument(n *rhino.Node, i int) *rhino.Node {
	return n.Children()[i+1]
}

// getCallParamCount mirrors TypeTransformation.getCallParamCount.
func getCallParamCount(n *rhino.Node) int {
	return len(n.Children()) - 1
}

// getCallParams mirrors TypeTransformation.getCallParams — returns
// all argument nodes (children beyond the callee).
func getCallParams(n *rhino.Node) []*rhino.Node {
	kids := n.Children()
	if len(kids) <= 1 {
		return nil
	}
	out := make([]*rhino.Node, len(kids)-1)
	copy(out, kids[1:])
	return out
}

// getComputedPropValue mirrors TypeTransformation.getComputedPropValue.
// For a COMPUTED_PROP node, the second child is the value expression.
// TODO: COMPUTED_PROP token is not yet in rhino/token.go; uses GETELEM
// as a placeholder — adjust when the token is ported.
func getComputedPropValue(n *rhino.Node) *rhino.Node {
	return n.Children()[1]
}

// getComputedPropName mirrors TypeTransformation.getComputedPropName.
// For a COMPUTED_PROP node, the first child holds the computed name.
func getComputedPropName(n *rhino.Node) string {
	return n.FirstChild().String_()
}

// getFunctionBody mirrors NodeUtil.getFunctionBody — FUNCTION children
// are [name, PARAM_LIST, BLOCK]; the body is children[2].
func getFunctionBody(n *rhino.Node) *rhino.Node {
	return n.Children()[2]
}

// ---- type helpers -------------------------------------------------

// getUnknownType mirrors TypeTransformation.getUnknownType.
func getUnknownType() jstype.JSType { return jstype.Unknown }

// getNoType mirrors TypeTransformation.getNoType.
// Java's NO_TYPE is the bottom of the type lattice (the empty union).
// We use a sentinel noType value that isEmptyType() would return true for.
func getNoType() jstype.JSType { return noTypeSentinel }

// noTypeSentinel is a package-level singleton that acts as Java's
// JSTypeNative.NO_TYPE / registry.getNativeObjectType(NO_TYPE). It is
// treated as an empty type in the TTL evaluator: any expression that
// returns noTypeSentinel degrades to Unknown in Eval.
var noTypeSentinel jstype.JSType = noType{}

type noType struct{}

func (noType) TypeName() string          { return "NoType" }
func (noType) IsSubtypeOf(_ jstype.JSType) bool { return true }

// isNoType reports whether t is the no-type sentinel.
func isNoType(t jstype.JSType) bool { _, ok := t.(noType); return ok }

// isEmptyType mirrors JSType.isEmptyType() — true for NO_TYPE. Used
// in Eval to degrade bottom results to Unknown.
func isEmptyType(t jstype.JSType) bool { return isNoType(t) }

// getAllType mirrors TypeTransformation.getAllType.
func getAllType() jstype.JSType { return jstype.All }

// getObjectType mirrors TypeTransformation.getObjectType — the root
// nominal type `Object`. Used as the empty-record sentinel.
func getObjectType() jstype.JSType { return objectTypeSentinel }

var objectTypeSentinel jstype.JSType = jstype.NewObjectType("Object")

// isObjectType reports whether t is the `Object` sentinel used by
// maprecord / evalRecordType to skip empty-record entries.
func isObjectType(t jstype.JSType) bool {
	if obj, ok := t.(*jstype.ObjectType); ok {
		return obj.Name == "Object"
	}
	return false
}

// createUnionType mirrors TypeTransformation.createUnionType.
func createUnionType(variants []jstype.JSType) jstype.JSType {
	return jstype.NewUnion(variants...)
}

// createRecordType mirrors TypeTransformation.createRecordType —
// returns an anonymous ObjectType whose members match props.
func createRecordType(props map[string]jstype.JSType) jstype.JSType {
	obj := jstype.NewObjectType("")
	obj.IsStructural = true
	for k, v := range props {
		obj.DefineProperty(k, v)
	}
	return obj
}

// ---- diagnostic helper --------------------------------------------

// reportWarning mirrors TypeTransformation.reportWarning which calls
// compiler.report(JSError.make(n, msg, param…)).
func (t *TypeTransformation) reportWarning(n *rhino.Node, code string, params ...string) {
	msg := code
	for i, p := range params {
		msg = strings.ReplaceAll(msg, fmt.Sprintf("{%d}", i), p)
	}
	t.diags = append(t.diags, diag.Diagnostic{
		Path:    t.path,
		Line:    n.Lineno(),
		Col:     n.Charno(),
		Code:    code,
		Level:   diag.Warning,
		Message: msg,
	})
}

// ---- type resolution ----------------------------------------------

// getType mirrors TypeTransformation.getType. Looks up typeName in
// the type environment.
//
// Java consults:
//   1. registry.getType(typeEnv, typeName)
//   2. typeEnv.getSlot(typeName) — then unwraps ctor/enum/typedef
//
// In Go, TypedScope.GetDeclaredType covers both paths at the Tier-0
// depth we implement.
func (t *TypeTransformation) getType(typeName string) jstype.JSType {
	// Check declared types (includes scope chain).
	declared := t.typeEnv.GetDeclaredType(typeName)
	if declared != nil {
		// If it is a constructor, return the instance type —
		// mirrors type.isConstructor() → getInstanceType().getRawType().
		if fn, ok := declared.(*jstype.FunctionType); ok && fn.IsConstructor && fn.Instance != nil {
			return fn.Instance
		}
		// TODO: isEnumElementType() / getEnumeratedTypeOfEnumElement()
		// when enum types are ported (Java lines 144-146).
		return declared
	}
	// TODO: slot.getJSDocInfo().getTypedefType() evaluation
	// (Java lines 149-152) — needs evaluateTypeExpression.
	return nil
}

// ---- public entry points -----------------------------------------

// Eval mirrors TypeTransformation.eval(Node, ImmutableMap<String,JSType>).
// Returns the resulting type after evaluating the TTL expression.
// Any warnings produced are available via t.Diagnostics().
func (t *TypeTransformation) Eval(ttlAst *rhino.Node, typeVars map[string]jstype.JSType) jstype.JSType {
	return t.EvalWithNameVars(ttlAst, typeVars, nil)
}

// EvalWithNameVars mirrors the two-argument eval overload visible via
// @VisibleForTesting in Java.
func (t *TypeTransformation) EvalWithNameVars(
	ttlAst *rhino.Node,
	typeVars map[string]jstype.JSType,
	nameVars map[string]string,
) jstype.JSType {
	if typeVars == nil {
		typeVars = map[string]jstype.JSType{}
	}
	if nameVars == nil {
		nameVars = map[string]string{}
	}
	result := t.evalInternal(ttlAst, newNameResolver(typeVars, nameVars))
	if isEmptyType(result) {
		return getUnknownType()
	}
	return result
}

// Diagnostics returns warnings collected during the most-recent Eval call.
func (t *TypeTransformation) Diagnostics() []diag.Diagnostic { return t.diags }

// ---- internal evaluation -----------------------------------------

// evalInternal mirrors TypeTransformation.evalInternal.
func (t *TypeTransformation) evalInternal(ttlAst *rhino.Node, nr *NameResolver) jstype.JSType {
	if t.isTypeName(ttlAst) {
		return t.evalTypeName(ttlAst)
	}
	if t.isTypeVar(ttlAst) {
		return t.evalTypeVar(ttlAst, nr)
	}
	name := getCallName(ttlAst)
	kw := nameToKeyword(name)
	if kw == nil {
		panic(fmt.Sprintf("TypeTransformation: unknown keyword %q", name))
	}
	switch kw.kind {
	case OpKindTypeConstructor:
		return t.evalTypeExpression(ttlAst, nr)
	case OpKindOperation:
		return t.evalOperationExpression(ttlAst, nr)
	default:
		panic("TypeTransformation: could not evaluate the type transformation expression")
	}
}

// evalOperationExpression mirrors TypeTransformation.evalOperationExpression.
func (t *TypeTransformation) evalOperationExpression(ttlAst *rhino.Node, nr *NameResolver) jstype.JSType {
	name := getCallName(ttlAst)
	kw := nameToKeyword(name)
	switch kw {
	case kwCOND:
		return t.evalConditional(ttlAst, nr)
	case kwMAPUNION:
		return t.evalMapunion(ttlAst, nr)
	case kwMAPRECORD:
		return t.evalMaprecord(ttlAst, nr)
	case kwTYPEOFVAR:
		return t.evalTypeOfVar(ttlAst)
	case kwINSTANCEOF:
		return t.evalInstanceOf(ttlAst, nr)
	case kwPRINTTYPE:
		return t.evalPrintType(ttlAst, nr)
	case kwPROPTYPE:
		return t.evalPropType(ttlAst, nr)
	default:
		panic("TypeTransformation: invalid type transformation operation")
	}
}

// evalTypeExpression mirrors TypeTransformation.evalTypeExpression.
func (t *TypeTransformation) evalTypeExpression(ttlAst *rhino.Node, nr *NameResolver) jstype.JSType {
	name := getCallName(ttlAst)
	kw := nameToKeyword(name)
	switch kw {
	case kwTYPE:
		return t.evalTemplatizedType(ttlAst, nr)
	case kwUNION:
		return t.evalUnionType(ttlAst, nr)
	case kwNONE:
		return getNoType()
	case kwALL:
		return getAllType()
	case kwUNKNOWN:
		return getUnknownType()
	case kwRAWTYPEOF:
		return t.evalRawTypeOf(ttlAst, nr)
	case kwTEMPLATETYPEOF:
		return t.evalTemplateTypeOf(ttlAst, nr)
	case kwRECORD:
		return t.evalRecordType(ttlAst, nr)
	case kwTYPEEXPR:
		return t.evalNativeTypeExpr(ttlAst)
	default:
		panic("TypeTransformation: invalid type expression")
	}
}

// evalTypeName mirrors TypeTransformation.evalTypeName.
func (t *TypeTransformation) evalTypeName(ttlAst *rhino.Node) jstype.JSType {
	typeName := ttlAst.String_()
	resultingType := t.getType(typeName)
	if resultingType == nil {
		t.reportWarning(ttlAst, ttlUnknownTypename, typeName)
		return getUnknownType()
	}
	return resultingType
}

// evalTemplatizedType mirrors TypeTransformation.evalTemplatizedType.
// Java: registry.createTemplatizedType(baseType, templatizedTypes).
func (t *TypeTransformation) evalTemplatizedType(ttlAst *rhino.Node, nr *NameResolver) jstype.JSType {
	params := getCallParams(ttlAst)
	firstParam := t.evalInternal(params[0], nr)
	// isRawTypeOfTemplatizedType() — in Go we check that the type is
	// an ObjectType whose TemplateTypes are known (i.e. it is a generic
	// base that can be templatized).
	// TODO: full isRawTypeOfTemplatizedType() check when TemplateType
	// registry is wired (Java lines 321-323).
	baseObj, ok := firstParam.(*jstype.ObjectType)
	if !ok {
		t.reportWarning(ttlAst, ttlBasetypeInvalid, firstParam.TypeName())
		return getUnknownType()
	}
	templatizedTypes := make([]jstype.JSType, len(params)-1)
	for i := 1; i < len(params); i++ {
		templatizedTypes[i-1] = t.evalInternal(params[i], nr)
	}
	// Mirrors registry.createTemplatizedType(baseType, templatizedTypes).
	return jstype.NewTemplatizedType(baseObj, templatizedTypes)
}

// evalTypeVar mirrors TypeTransformation.evalTypeVar.
func (t *TypeTransformation) evalTypeVar(ttlAst *rhino.Node, nr *NameResolver) jstype.JSType {
	typeVar := ttlAst.String_()
	resultingType, ok := nr.typeVars[typeVar]
	if !ok {
		t.reportWarning(ttlAst, ttlUnknownTypevar, typeVar)
		return getUnknownType()
	}
	return resultingType
}

// evalUnionType mirrors TypeTransformation.evalUnionType.
func (t *TypeTransformation) evalUnionType(ttlAst *rhino.Node, nr *NameResolver) jstype.JSType {
	params := getCallParams(ttlAst)
	basicTypes := make([]jstype.JSType, len(params))
	for i, p := range params {
		basicTypes[i] = t.evalInternal(p, nr)
	}
	return createUnionType(basicTypes)
}

// evalTypeParams mirrors TypeTransformation.evalTypeParams.
func (t *TypeTransformation) evalTypeParams(ttlAst *rhino.Node, nr *NameResolver) []jstype.JSType {
	params := getCallParams(ttlAst)
	result := make([]jstype.JSType, len(params))
	for i, p := range params {
		result[i] = t.evalInternal(p, nr)
	}
	return result
}

// evalString mirrors TypeTransformation.evalString. Returns the string
// value of a NAME (resolved via nameVars) or STRINGLIT node.
func (t *TypeTransformation) evalString(ttlAst *rhino.Node, nr *NameResolver) string {
	if ttlAst.Token() == rhino.NAME {
		val, ok := nr.nameVars[ttlAst.String_()]
		if !ok {
			t.reportWarning(ttlAst, ttlUnknownStrvar, ttlAst.String_())
			return ""
		}
		return val
	}
	return ttlAst.String_()
}

// evalStringParams mirrors TypeTransformation.evalStringParams.
func (t *TypeTransformation) evalStringParams(ttlAst *rhino.Node, nr *NameResolver) []string {
	params := getCallParams(ttlAst)
	result := make([]string, len(params))
	for i, p := range params {
		result[i] = t.evalString(p, nr)
	}
	return result
}

// evalTypePredicate mirrors TypeTransformation.evalTypePredicate.
func (t *TypeTransformation) evalTypePredicate(ttlAst *rhino.Node, nr *NameResolver) bool {
	params := t.evalTypeParams(ttlAst, nr)
	name := getCallName(ttlAst)
	kw := nameToKeyword(name)
	typ := params[0]
	switch kw {
	case kwEQ:
		return typ.TypeName() == params[1].TypeName()
	case kwSUB:
		return jstype.IsSubtype(typ, params[1])
	case kwISCTOR:
		fn, ok := typ.(*jstype.FunctionType)
		return ok && fn.IsConstructor
	case kwISTEMPLATIZED:
		_, ok := typ.(*jstype.TemplatizedType)
		return ok
	case kwISRECORD:
		obj, ok := typ.(*jstype.ObjectType)
		return ok && obj.IsStructural
	case kwISUNKNOWN:
		return typ == jstype.Unknown
	default:
		panic("TypeTransformation: invalid type predicate in the type transformation")
	}
}

// evalStringPredicate mirrors TypeTransformation.evalStringPredicate.
func (t *TypeTransformation) evalStringPredicate(ttlAst *rhino.Node, nr *NameResolver) bool {
	params := t.evalStringParams(ttlAst, nr)
	// If any param resolved to empty (unresolved name var), always false.
	for _, p := range params {
		if p == "" {
			return false
		}
	}
	name := getCallName(ttlAst)
	kw := nameToKeyword(name)
	switch kw {
	case kwSTREQ:
		return params[0] == params[1]
	default:
		panic("TypeTransformation: invalid string predicate in the type transformation")
	}
}

// evalTypevarPredicate mirrors TypeTransformation.evalTypevarPredicate.
func (t *TypeTransformation) evalTypevarPredicate(ttlAst *rhino.Node, nr *NameResolver) bool {
	name := getCallName(ttlAst)
	kw := nameToKeyword(name)
	switch kw {
	case kwISDEFINED:
		_, ok := nr.typeVars[getCallArgument(ttlAst, 0).String_()]
		return ok
	default:
		panic("TypeTransformation: invalid typevar predicate in the type transformation")
	}
}

// evalBooleanOperation mirrors TypeTransformation.evalBooleanOperation.
func (t *TypeTransformation) evalBooleanOperation(ttlAst *rhino.Node, nr *NameResolver) bool {
	param0 := t.evalBoolean(ttlAst.FirstChild(), nr)
	switch ttlAst.Token() {
	case rhino.NOT:
		return !param0
	case rhino.AND:
		return param0 && t.evalBoolean(ttlAst.Children()[len(ttlAst.Children())-1], nr)
	case rhino.OR:
		return param0 || t.evalBoolean(ttlAst.Children()[len(ttlAst.Children())-1], nr)
	default:
		panic("TypeTransformation: invalid boolean predicate in the type transformation")
	}
}

// evalBoolean mirrors TypeTransformation.evalBoolean.
func (t *TypeTransformation) evalBoolean(ttlAst *rhino.Node, nr *NameResolver) bool {
	if t.isBooleanOperation(ttlAst) {
		return t.evalBooleanOperation(ttlAst, nr)
	}
	name := getCallName(ttlAst)
	kw := nameToKeyword(name)
	if kw == nil {
		panic(fmt.Sprintf("TypeTransformation: unknown keyword %q in boolean", name))
	}
	switch kw.kind {
	case OpKindStringPredicate:
		return t.evalStringPredicate(ttlAst, nr)
	case OpKindTypePredicate:
		return t.evalTypePredicate(ttlAst, nr)
	case OpKindTypevarPredicate:
		return t.evalTypevarPredicate(ttlAst, nr)
	default:
		panic("TypeTransformation: invalid boolean predicate in the type transformation")
	}
}

// evalConditional mirrors TypeTransformation.evalConditional.
func (t *TypeTransformation) evalConditional(ttlAst *rhino.Node, nr *NameResolver) jstype.JSType {
	params := getCallParams(ttlAst)
	if t.evalBoolean(params[0], nr) {
		return t.evalInternal(params[1], nr)
	}
	return t.evalInternal(params[2], nr)
}

// evalMapunion mirrors TypeTransformation.evalMapunion.
func (t *TypeTransformation) evalMapunion(ttlAst *rhino.Node, nr *NameResolver) jstype.JSType {
	params := getCallParams(ttlAst)
	unionParam := params[0]
	mapFunction := params[1]
	paramName := getFunctionParameter(mapFunction, 0)

	// The mapunion variable must not be defined in the environment.
	if _, exists := nr.typeVars[paramName]; exists {
		t.reportWarning(ttlAst, ttlDuplicateVariable, paramName)
		return getUnknownType()
	}

	mapFunctionBody := getFunctionBody(mapFunction)
	unionType := t.evalInternal(unionParam, nr)

	// If the union has only one member, evaluate with that single type.
	unionT, isUnion := unionType.(*jstype.UnionType)
	if !isUnion {
		newNR := nr.addTypeVar(paramName, unionType)
		return t.evalInternal(mapFunctionBody, newNR)
	}

	// Evaluate the map function body for each member.
	members := unionT.Members()
	newUnionElms := make([]jstype.JSType, len(members))
	for i, elm := range members {
		newNR := nr.addTypeVar(paramName, elm)
		newUnionElms[i] = t.evalInternal(mapFunctionBody, newNR)
	}
	return createUnionType(newUnionElms)
}

// evalRawTypeOf mirrors TypeTransformation.evalRawTypeOf.
func (t *TypeTransformation) evalRawTypeOf(ttlAst *rhino.Node, nr *NameResolver) jstype.JSType {
	params := getCallParams(ttlAst)
	typ := t.evalInternal(params[0], nr)
	templ, ok := typ.(*jstype.TemplatizedType)
	if !ok {
		t.reportWarning(ttlAst, ttlTemptypeInvalid, "rawTypeOf", typ.TypeName())
		return getUnknownType()
	}
	// Mirrors type.toMaybeObjectType().getRawType() — returns the base
	// ObjectType without template arguments.
	return templ.Base
}

// evalTemplateTypeOf mirrors TypeTransformation.evalTemplateTypeOf.
func (t *TypeTransformation) evalTemplateTypeOf(ttlAst *rhino.Node, nr *NameResolver) jstype.JSType {
	params := getCallParams(ttlAst)
	typ := t.evalInternal(params[0], nr)
	templ, ok := typ.(*jstype.TemplatizedType)
	if !ok {
		t.reportWarning(ttlAst, ttlTemptypeInvalid, "templateTypeOf", typ.TypeName())
		return getUnknownType()
	}
	index := int(params[1].Number())
	if index >= len(templ.Args) {
		t.reportWarning(ttlAst, ttlIndexOutofbounds,
			fmt.Sprintf("%d", len(templ.Args)),
			fmt.Sprintf("%d", index))
		return getUnknownType()
	}
	return templ.Args[index]
}

// evalRecord mirrors TypeTransformation.evalRecord. Walks an OBJECTLIT
// or a node-with-object-literal-shape and builds the property map.
//
// Java's evalRecord iterates child nodes; computed-prop children use
// nameVar resolution and regular prop children use their key string.
// COMPUTED_PROP is not yet in rhino/token.go — we handle the OBJECTLIT
// children using STRING_KEY which is what the parser emits for
// `{key: expr}` patterns. Computed props (GETELEM-based) are noted as
// a TODO.
func (t *TypeTransformation) evalRecord(record *rhino.Node, nr *NameResolver) jstype.JSType {
	props := map[string]jstype.JSType{}
	for _, propNode := range record.Children() {
		// TODO: handle COMPUTED_PROP (Java isComputedProp) when that
		// token is ported (Java lines 547-557).
		propName := propNode.String_()
		valueNode := propNode.FirstChild()
		if valueNode == nil {
			continue
		}
		resultingType := t.evalInternal(valueNode, nr)
		props[propName] = resultingType
	}
	return createRecordType(props)
}

// evalRecordParam mirrors TypeTransformation.evalRecordParam.
func (t *TypeTransformation) evalRecordParam(ttlAst *rhino.Node, nr *NameResolver) jstype.JSType {
	if ttlAst.Token() == rhino.OBJECTLIT {
		return t.evalRecord(ttlAst, nr)
	}
	return t.evalInternal(ttlAst, nr)
}

// evalRecordType mirrors TypeTransformation.evalRecordType.
func (t *TypeTransformation) evalRecordType(ttlAst *rhino.Node, nr *NameResolver) jstype.JSType {
	paramCount := getCallParamCount(ttlAst)
	var recTypes []*jstype.ObjectType
	for i := 0; i < paramCount; i++ {
		typ := t.evalRecordParam(getCallArgument(ttlAst, i), nr)
		obj, ok := typ.(*jstype.ObjectType)
		if !ok || typ == jstype.Unknown {
			t.reportWarning(ttlAst, ttlRecparamInvalid, typ.TypeName())
			return getUnknownType()
		}
		// Mirrors registry.buildRecordTypeFromObject(objType) + skip Object.
		// TODO: full buildRecordTypeFromObject when needed (Java line 587).
		if !isObjectType(typ) {
			recTypes = append(recTypes, obj)
		}
	}
	return t.joinRecordTypes(recTypes)
}

// putNewPropInPropertyMap mirrors TypeTransformation.putNewPropInPropertyMap.
func (t *TypeTransformation) putNewPropInPropertyMap(
	props map[string]jstype.JSType,
	newPropName string,
	newPropValue jstype.JSType,
) {
	existing, exists := props[newPropName]
	if !exists {
		props[newPropName] = newPropValue
		return
	}
	existingIsRecord := t.isRecordType(existing)
	newIsRecord := t.isRecordType(newPropValue)
	if !newIsRecord || !existingIsRecord {
		props[newPropName] = newPropValue
		return
	}
	// Both are records — join them.
	existObj := existing.(*jstype.ObjectType)
	newObj := newPropValue.(*jstype.ObjectType)
	props[newPropName] = t.joinRecordTypes([]*jstype.ObjectType{existObj, newObj})
}

// isRecordType mirrors JSType.isRecordType() — an ObjectType that is
// marked structural. Returns false for all other type kinds.
func (t *TypeTransformation) isRecordType(typ jstype.JSType) bool {
	obj, ok := typ.(*jstype.ObjectType)
	return ok && obj.IsStructural
}

// joinRecordTypes mirrors TypeTransformation.joinRecordTypes — merges
// all property maps from a list of record types into a single record.
func (t *TypeTransformation) joinRecordTypes(recTypes []*jstype.ObjectType) jstype.JSType {
	props := map[string]jstype.JSType{}
	for _, recType := range recTypes {
		for _, name := range recType.MemberNames() {
			val := recType.GetPropertyType(name)
			t.putNewPropInPropertyMap(props, name, val)
		}
	}
	return createRecordType(props)
}

// evalMaprecord mirrors TypeTransformation.evalMaprecord.
func (t *TypeTransformation) evalMaprecord(ttlAst *rhino.Node, nr *NameResolver) jstype.JSType {
	kids := ttlAst.Children()
	recordNode := kids[1]
	mapFunction := kids[2]
	typ := t.evalInternal(recordNode, nr)

	// If it is an empty record type (Object) then return.
	if isObjectType(typ) {
		return getObjectType()
	}

	// The parameter must be a valid record type.
	if !t.isRecordType(typ) {
		t.reportWarning(recordNode, ttlRectypeInvalid, typ.TypeName())
		return getUnknownType()
	}

	objtype := typ.(*jstype.ObjectType)

	// Fetch map function parameter names.
	paramKey := getFunctionParameter(mapFunction, 0)
	paramValue := getFunctionParameter(mapFunction, 1)

	// The maprecord variables must not already be defined.
	if _, exists := nr.nameVars[paramKey]; exists {
		t.reportWarning(ttlAst, ttlDuplicateVariable, paramKey)
		return getUnknownType()
	}
	if _, exists := nr.typeVars[paramValue]; exists {
		t.reportWarning(ttlAst, ttlDuplicateVariable, paramValue)
		return getUnknownType()
	}

	mapFnBody := getFunctionBody(mapFunction)
	newProps := map[string]jstype.JSType{}

	for _, propName := range objtype.MemberNames() {
		propValue := objtype.GetPropertyType(propName)

		newNR := nr.addTypeVar(paramValue, propValue).addNameVar(paramKey, propName)
		body := t.evalInternal(mapFnBody, newNR)

		// If body is unknown, the whole expression is unknown.
		if body == jstype.Unknown {
			return getUnknownType()
		}

		// Skip when body evaluates to NO_TYPE or empty Object record.
		if isEmptyType(body) || isObjectType(body) {
			continue
		}

		// Otherwise body must be a record type.
		if !t.isRecordType(body) {
			t.reportWarning(ttlAst, ttlMaprecordBodyInvalid, body.TypeName())
			return getUnknownType()
		}

		// Add properties from body record to newProps.
		bodyObj := body.(*jstype.ObjectType)
		for _, newPropName := range bodyObj.MemberNames() {
			newPropValue := bodyObj.GetPropertyType(newPropName)
			t.putNewPropInPropertyMap(newProps, newPropName, newPropValue)
		}
	}
	return createRecordType(newProps)
}

// evalTypeOfVar mirrors TypeTransformation.evalTypeOfVar.
func (t *TypeTransformation) evalTypeOfVar(ttlAst *rhino.Node) jstype.JSType {
	name := getCallArgument(ttlAst, 0).String_()
	typ := t.typeEnv.GetDeclaredType(name)
	if typ == nil {
		t.reportWarning(ttlAst, ttlVarUndefined, name)
		return getUnknownType()
	}
	return typ
}

// evalInstanceOf mirrors TypeTransformation.evalInstanceOf.
func (t *TypeTransformation) evalInstanceOf(ttlAst *rhino.Node, nr *NameResolver) jstype.JSType {
	typ := t.evalInternal(getCallArgument(ttlAst, 0), nr)
	fn, ok := typ.(*jstype.FunctionType)
	if typ == jstype.Unknown || !ok || !fn.IsConstructor {
		t.reportWarning(ttlAst, ttlInvalidCtor, typ.TypeName())
		return getUnknownType()
	}
	return fn.Instance
}

// evalNativeTypeExpr mirrors TypeTransformation.evalNativeTypeExpr.
// Java: registry.evaluateTypeExpression(new JSTypeExpression(arg, FILE), typeEnv).
// In Go, we delegate to the registry's ResolveExpr which handles the
// textual type-expression subset available today.
//
// TODO: full JSTypeExpression evaluation when rhino.JSTypeExpression and
// registry.evaluateTypeExpression are ported (Java lines 734-736).
func (t *TypeTransformation) evalNativeTypeExpr(ttlAst *rhino.Node) jstype.JSType {
	arg := getCallArgument(ttlAst, 0)
	return t.typeEnv.registry.ResolveExpr(arg.String_())
}

// evalPrintType mirrors TypeTransformation.evalPrintType — prints
// the type to stdout (Java uses System.out.println) and returns it.
func (t *TypeTransformation) evalPrintType(ttlAst *rhino.Node, nr *NameResolver) jstype.JSType {
	typ := t.evalInternal(getCallArgument(ttlAst, 1), nr)
	msg := getCallArgument(ttlAst, 0).String_() + typ.TypeName()
	fmt.Println(msg)
	return typ
}

// evalPropType mirrors TypeTransformation.evalPropType.
func (t *TypeTransformation) evalPropType(ttlAst *rhino.Node, nr *NameResolver) jstype.JSType {
	typ := t.evalInternal(getCallArgument(ttlAst, 1), nr)
	obj, ok := typ.(*jstype.ObjectType)
	if !ok {
		t.reportWarning(ttlAst, ttlProptypeInvalid, typ.TypeName())
		return getUnknownType()
	}
	propName := getCallArgument(ttlAst, 0).String_()
	propType := obj.GetPropertyType(propName)
	if propType == nil {
		return getUnknownType()
	}
	return propType
}

