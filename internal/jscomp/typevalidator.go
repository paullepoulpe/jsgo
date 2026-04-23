// Port of com/google/javascript/jscomp/TypeValidator.java
// @ closure-compiler v20260415 (subset).
//
// TypeValidator is a central reporter for type violations. It accumulates
// TypeMismatch entries and exposes expect* methods that are called from
// TypeCheck (and eventually other passes). The methods mirror their Java
// counterparts 1:1 in name and signature shape; bodies are faithful
// subsets adapted to jsgo's narrower JSType hierarchy.
//
// The subset ported here covers the public API surface Java exposes:
//   - expectCanAssignTo / expectCanAssignToPropertyOf
//   - expectArgumentMatchesParameter
//   - expectCanCast
//   - expectIndexMatch
//   - expectSwitchMatchesCase
//   - expectValidTypeofName
//   - expectObject / expectActualObject / expectAnyObject
//   - expectString / expectNumber / expectNumberStrict
//   - expectStringOrNumberOrSymbol / expectStringOrSymbol / expectNumberOrSymbol
//   - expectBitwiseable / expectBigIntOrNumber / expectUnknownOrComparable
//   - expectNotNullOrUndefined
//   - getMismatches / setSubtypingMode
//
// Methods that depend on Java types not yet ported (FunctionType-builder
// chains, TemplateType walking, CompilerInput, TypedVar/undeclare) are
// marked TODO and left as stubs.

package jscomp

import (
	"fmt"
	"sort"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// TypeMismatch is defined in typemismatch.go (canonical port). TypeValidator
// uses the same struct; the accumulator below wraps []TypeMismatch.

// typeMismatchAccumulator holds the running list of mismatches seen
// during one pass. Equivalent to TypeMismatch.Accumulator in Java.
type typeMismatchAccumulator struct {
	mismatches []TypeMismatch
}

func (a *typeMismatchAccumulator) registerMismatch(n *rhino.Node, found, required jstype.JSType) {
	a.mismatches = append(a.mismatches, TypeMismatch{Found: found, Required: required, Location: n})
}

func (a *typeMismatchAccumulator) getMismatches() []TypeMismatch {
	return a.mismatches
}

// ---- TypeValidator --------------------------------------------------- //

// TypeValidator is a central reporter for all type violations. Holds the
// type registry, an error reporter (path + diag accumulator), and a
// running accumulator of TypeMismatch entries. Mirrors Java's
// TypeValidator class.
type TypeValidator struct {
	// path is the source-file path used when constructing diagnostics.
	path string
	// registry is the JSType registry — used to resolve native types
	// and build union/templatized types. Mirrors typeRegistry field.
	registry *jstype.Registry
	// diags accumulates diagnostics emitted during the pass.
	diags []diag.Diagnostic
	// mismatches accumulates TypeMismatch pairs.
	mismatches typeMismatchAccumulator
	// subtypingMode mirrors Java's SubtypingMode field. When set to
	// tvSubtypingIgnoreNullUndefined, null/undefined are ignored in
	// subtype checks (used for .java.js files in Java).
	subtypingMode tvSubtypingMode
	// allBitwisableValueTypes mirrors Java's field of the same name:
	// string|number|boolean|null|undefined — valid bitwise operands.
	allBitwisableValueTypes jstype.JSType
	// nullOrUndefined mirrors Java's nullOrUndefined field.
	nullOrUndefined jstype.JSType
}

// tvSubtypingMode mirrors Java's SubtypingMode enum.
type tvSubtypingMode int

const (
	tvSubtypingNormal                tvSubtypingMode = iota
	tvSubtypingIgnoreNullUndefined   tvSubtypingMode = iota
)

// NewTypeValidator constructs a TypeValidator. Mirrors Java's
// TypeValidator(AbstractCompiler) constructor.
func NewTypeValidator(path string, registry *jstype.Registry) *TypeValidator {
	if registry == nil {
		registry = jstype.NewRegistry()
	}
	return &TypeValidator{
		path:     path,
		registry: registry,
		// string|number|boolean|null|undefined
		allBitwisableValueTypes: jstype.NewUnion(
			jstype.String, jstype.Number, jstype.Boolean, jstype.Null, jstype.Void),
		// null|undefined
		nullOrUndefined: jstype.NewUnion(jstype.Null, jstype.Void),
	}
}

// GetMismatches returns the list of type violations accumulated so far.
// Mirrors Java's TypeValidator.getMismatches().
func (tv *TypeValidator) GetMismatches() []TypeMismatch {
	return tv.mismatches.getMismatches()
}

// Diagnostics returns the list of diagnostics emitted by this validator.
func (tv *TypeValidator) Diagnostics() []diag.Diagnostic {
	return tv.diags
}

// SetSubtypingMode mirrors Java's TypeValidator.setSubtypingMode(SubtypingMode).
func (tv *TypeValidator) SetSubtypingMode(mode tvSubtypingMode) {
	tv.subtypingMode = mode
}

// ---- public expect* methods ----------------------------------------- //

// expectValidTypeofName fires JSC_UNKNOWN_TYPEOF_VALUE when the operand
// of a typeof expression is not one of the known type strings. Mirrors
// Java's TypeValidator.expectValidTypeofName(Node, String).
func (tv *TypeValidator) expectValidTypeofName(n *rhino.Node, found string) {
	tv.emit(diag.Diagnostic{
		Path:    tv.path,
		Line:    n.Lineno(),
		Col:     n.Charno(),
		Code:    diag.JSC_UNRECOGNIZED_TYPEOF_ERROR,
		Level:   diag.Warning,
		Message: fmt.Sprintf("unknown type: %s", found),
	})
}

// expectObject fires JSC_TYPE_MISMATCH when `typ` cannot be used in an
// object context. Returns true when no warning was issued.
// Mirrors Java's TypeValidator.expectObject(Node, JSType, String).
func (tv *TypeValidator) expectObject(n *rhino.Node, typ jstype.JSType, msg string) bool {
	if tvMatchesObjectContext(typ) {
		return true
	}
	tv.tvMismatch(n, msg, typ, jstype.NewObjectType("Object"))
	return false
}

// expectActualObject fires JSC_TYPE_MISMATCH when `typ` is not an object
// (types merely convertible to objects are not accepted).
// Mirrors Java's TypeValidator.expectActualObject(Node, JSType, String).
func (tv *TypeValidator) expectActualObject(n *rhino.Node, typ jstype.JSType, msg string) {
	if !tvIsObject(typ) {
		tv.tvMismatch(n, msg, typ, jstype.NewObjectType("Object"))
	}
}

// expectAnyObject fires JSC_TYPE_MISMATCH when `typ` can never contain
// an object (i.e. is definitely a pure primitive).
// Mirrors Java's TypeValidator.expectAnyObject(Node, JSType, String).
func (tv *TypeValidator) expectAnyObject(n *rhino.Node, typ jstype.JSType, msg string) {
	anyObject := jstype.NewObjectType("Object")
	if !jstype.IsSubtype(anyObject, typ) && !tvIsEmptyType(typ) {
		tv.tvMismatch(n, msg, typ, anyObject)
	}
}

// expectString fires JSC_TYPE_MISMATCH when `typ` is not convertible
// to a string. Mirrors Java's TypeValidator.expectString(Node, JSType, String).
func (tv *TypeValidator) expectString(n *rhino.Node, typ jstype.JSType, msg string) {
	if !tvMatchesStringContext(typ) {
		tv.tvMismatch(n, msg, typ, jstype.String)
	}
}

// expectNumber fires JSC_TYPE_MISMATCH when `typ` is not convertible
// to a number. When it is number-compatible but not a strict subtype,
// it also fires the strict (INVALID_OPERAND_TYPE) path.
// Mirrors Java's TypeValidator.expectNumber(Node, JSType, String).
func (tv *TypeValidator) expectNumber(n *rhino.Node, typ jstype.JSType, msg string) {
	if !tvMatchesNumberContext(typ) {
		tv.tvMismatch(n, msg, typ, jstype.Number)
	} else {
		tv.expectNumberStrict(n, typ, msg)
	}
}

// expectNumberStrict fires JSC_INVALID_OPERAND_TYPE when `typ` is not a
// strict subtype of number. Mirrors Java's TypeValidator.expectNumberStrict.
func (tv *TypeValidator) expectNumberStrict(n *rhino.Node, typ jstype.JSType, msg string) {
	if !jstype.IsSubtype(typ, jstype.Number) {
		tv.tvRegisterMismatchAndReportWithDiag(n, diag.JSC_INVALID_OPERAND_TYPE, msg, typ, jstype.Number, nil, nil)
	}
}

// expectMatchingTypesStrict fires JSC_INVALID_OPERAND_TYPE when neither
// left is subtype of right nor right is subtype of left.
// Mirrors Java's TypeValidator.expectMatchingTypesStrict.
func (tv *TypeValidator) expectMatchingTypesStrict(n *rhino.Node, left, right jstype.JSType, msg string) {
	if !jstype.IsSubtype(left, right) && !jstype.IsSubtype(right, left) {
		tv.tvRegisterMismatchAndReportWithDiag(n, diag.JSC_INVALID_OPERAND_TYPE, msg, right, left, nil, nil)
	}
}

// expectBitwiseable fires JSC_TYPE_MISMATCH when `typ` is not a valid
// bitwise operand (number-context or primitive), then also fires the strict
// number check. Mirrors Java's TypeValidator.expectBitwiseable.
func (tv *TypeValidator) expectBitwiseable(n *rhino.Node, typ jstype.JSType, msg string) {
	if !tvMatchesNumberContext(typ) && !jstype.IsSubtype(typ, tv.allBitwisableValueTypes) {
		tv.tvMismatch(n, msg, typ, tv.allBitwisableValueTypes)
	} else {
		tv.expectNumberStrict(n, typ, msg)
	}
}

// expectNumberOrSymbol fires when `typ` is neither number-context nor
// symbol-context. Mirrors Java's TypeValidator.expectNumberOrSymbol.
func (tv *TypeValidator) expectNumberOrSymbol(n *rhino.Node, typ jstype.JSType, msg string) {
	if !tvMatchesNumberContext(typ) && !tvMatchesSymbolContext(typ) {
		// Report mismatch against number (symbol not in our type lattice)
		tv.tvMismatch(n, msg, typ, jstype.Number)
	}
}

// expectStringOrSymbol fires when `typ` is neither string-context nor
// symbol-context. Mirrors Java's TypeValidator.expectStringOrSymbol.
func (tv *TypeValidator) expectStringOrSymbol(n *rhino.Node, typ jstype.JSType, msg string) {
	if !tvMatchesStringContext(typ) && !tvMatchesSymbolContext(typ) {
		tv.tvMismatch(n, msg, typ, jstype.String)
	}
}

// expectStringOrNumberOrSymbol fires when `typ` is not string, number,
// or symbol context. Mirrors Java's TypeValidator.expectStringOrNumberOrSymbol.
func (tv *TypeValidator) expectStringOrNumberOrSymbol(n *rhino.Node, typ jstype.JSType, msg string) {
	if !tvMatchesNumberContext(typ) && !tvMatchesStringContext(typ) && !tvMatchesSymbolContext(typ) {
		tv.tvMismatch(n, msg, typ, jstype.NewUnion(jstype.String, jstype.Number))
	} else {
		tv.expectStringOrNumberOrSymbolStrict(n, typ, msg)
	}
}

// expectStringOrNumberOrSymbolStrict fires JSC_INVALID_OPERAND_TYPE when
// `typ` is not a strict subtype of string|number.
// Mirrors Java's TypeValidator.expectStringOrNumberOrSymbolStrict.
func (tv *TypeValidator) expectStringOrNumberOrSymbolStrict(n *rhino.Node, typ jstype.JSType, msg string) {
	strNum := jstype.NewUnion(jstype.String, jstype.Number)
	if !jstype.IsSubtype(typ, strNum) {
		tv.tvRegisterMismatchAndReportWithDiag(n, diag.JSC_INVALID_OPERAND_TYPE, msg, typ, strNum, nil, nil)
	}
}

// expectUnknownOrComparable fires when `typ` is not a comparable type
// (bigint/number/string). In jsgo, bigint is not modeled, so we use
// number|string as the comparable set.
// Mirrors Java's TypeValidator.expectUnknownOrComparable.
func (tv *TypeValidator) expectUnknownOrComparable(n *rhino.Node, typ jstype.JSType, msg string) {
	strNum := jstype.NewUnion(jstype.String, jstype.Number)
	if jstype.IsSubtype(typ, strNum) {
		return
	}
	if tvMatchesNumberContext(typ) || tvMatchesStringContext(typ) {
		tv.tvRegisterMismatchAndReportWithDiag(n, diag.JSC_INVALID_OPERAND_TYPE, msg, typ, strNum, nil, nil)
	} else {
		tv.tvMismatch(n, msg, typ, strNum)
	}
}

// expectBigIntOrNumber fires when `typ` is not a number/bigint type.
// In jsgo bigint is not modeled; we use number as the target.
// Mirrors Java's TypeValidator.expectBigIntOrNumber.
func (tv *TypeValidator) expectBigIntOrNumber(n *rhino.Node, typ jstype.JSType, msg string) {
	if jstype.IsSubtype(typ, jstype.Number) {
		return
	}
	if tvMatchesNumberContext(typ) {
		tv.tvRegisterMismatchAndReportWithDiag(n, diag.JSC_INVALID_OPERAND_TYPE, msg, typ, jstype.Number, nil, nil)
	} else {
		tv.tvMismatch(n, msg, typ, jstype.Number)
	}
}

// expectNotNullOrUndefined fires JSC_TYPE_MISMATCH when `typ` is
// definitely null or undefined. Returns false when a warning was issued.
// Mirrors Java's TypeValidator.expectNotNullOrUndefined.
func (tv *TypeValidator) expectNotNullOrUndefined(n *rhino.Node, typ jstype.JSType, msg string, expectedType jstype.JSType) bool {
	if tvIsNoType(typ) || tvIsUnknownType(typ) {
		return true
	}
	if !jstype.IsSubtype(typ, tv.nullOrUndefined) {
		return true
	}
	// containsForwardDeclaredUnresolvedName equivalent — jsgo does not
	// track unresolved forward declarations, so skip that guard.
	//
	// Java's edge case: if (this.x == null) { this.x.foo() } — suppress
	// for GETPROP inside non-global scope when the type is null. Omitted
	// since jsgo does not track global-scope flag here.
	tv.tvMismatch(n, msg, typ, expectedType)
	return false
}

// expectSwitchMatchesCase fires JSC_TYPE_MISMATCH when `switchType` and
// `caseType` can never be equal (canTestForShallowEqualityWith returns false).
// Mirrors Java's TypeValidator.expectSwitchMatchesCase.
func (tv *TypeValidator) expectSwitchMatchesCase(n *rhino.Node, switchType, caseType jstype.JSType) {
	first := n.FirstChild()
	if first == nil {
		first = n
	}
	if !tvCanTestForShallowEqualityWith(switchType, caseType) {
		tv.tvMismatch(first, "case expression doesn't match switch", caseType, switchType)
	}
}

// expectIndexMatch fires diagnostics when `indexType` is not a valid
// index for `objType`. Mirrors Java's TypeValidator.expectIndexMatch.
func (tv *TypeValidator) expectIndexMatch(n *rhino.Node, objType, indexType jstype.JSType) {
	// Determine the index node: for GETELEM the last child (the index
	// expression); our simplified AST exposes children as a slice, so
	// we use the last child for GETELEM and first for COMPUTED_PROP.
	// Fall back to n itself when children are unavailable.
	kids := n.Children()
	var indexNode *rhino.Node
	if n.Token() == rhino.GETELEM && len(kids) >= 2 {
		indexNode = kids[len(kids)-1]
	} else if len(kids) >= 1 {
		indexNode = kids[0]
	}
	if indexNode == nil {
		indexNode = n
	}

	// Symbol indexes are always allowed.
	if tvIsSymbolValueType(indexType) {
		return
	}

	if tvIsUnknownType(objType) {
		tv.expectStringOrNumberOrSymbol(indexNode, indexType, "property access")
		return
	}

	obj := tvDereferencedObject(objType)
	if obj != nil && obj.Name == "Array" {
		tv.expectNumberOrSymbol(indexNode, indexType, "array access")
		return
	}
	if obj != nil && obj.IsStruct {
		tv.emit(diag.Diagnostic{
			Path:    tv.path,
			Line:    indexNode.Lineno(),
			Col:     indexNode.Charno(),
			Code:    diag.JSC_ILLEGAL_PROPERTY_ACCESS,
			Level:   diag.Error,
			Message: "Cannot do '[]' access on a struct",
		})
		return
	}
	if tvMatchesObjectContext(objType) {
		tv.expectStringOrSymbol(indexNode, indexType, "property access")
		return
	}
	// Not an object at all.
	tv.tvMismatch(n, "only arrays or objects can be accessed", objType,
		jstype.NewUnion(jstype.NewObjectType("Array"), jstype.NewObjectType("Object")))
}

// expectCanAssignTo checks that `rightType` is assignable to `leftType`.
// Returns false and emits JSC_TYPE_MISMATCH when not.
// Mirrors Java's TypeValidator.expectCanAssignTo(Node, JSType, JSType, String).
func (tv *TypeValidator) expectCanAssignTo(n *rhino.Node, rightType, leftType jstype.JSType, msg string) bool {
	// Template types: always accept (Java's open-type rule).
	if _, ok := leftType.(*jstype.TemplateType); ok {
		return true
	}
	if !jstype.IsSubtype(rightType, leftType) {
		tv.tvMismatch(n, msg, rightType, leftType)
		return false
	}
	return true
}

// expectCanAssignToPropertyOf checks that `rightType` can be assigned to
// a property `propName` owned by `ownerName` of type `leftType`.
// Mirrors Java's TypeValidator.expectCanAssignToPropertyOf(Node, JSType, JSType, JSType, Supplier<String>, String).
func (tv *TypeValidator) expectCanAssignToPropertyOf(n *rhino.Node, rightType, leftType jstype.JSType, ownerName string, propName string) bool {
	// Template types: always accept.
	if _, ok := leftType.(*jstype.TemplateType); ok {
		return true
	}
	if !jstype.IsSubtype(rightType, leftType) {
		// Java skips interface methods: if ownerFn.isInterface() &&
		// both types are functions. We approximate: if both are FunctionType,
		// skip. (Interfaces carry no explicit flag in jsgo yet.)
		if _, rf := rightType.(*jstype.FunctionType); rf {
			if _, lf := leftType.(*jstype.FunctionType); lf {
				// Possibly an interface method; stay quiet.
				return true
			}
		}
		tv.tvMismatch(n,
			fmt.Sprintf("assignment to property %s of %s", propName, ownerName),
			rightType, leftType)
		return false
	}
	return true
}

// expectArgumentMatchesParameter fires JSC_TYPE_MISMATCH when `argType`
// is not a subtype of `paramType` at a call site.
// Mirrors Java's TypeValidator.expectArgumentMatchesParameter.
func (tv *TypeValidator) expectArgumentMatchesParameter(n *rhino.Node, argType, paramType jstype.JSType, callNode *rhino.Node, ordinal int) {
	if jstype.IsSubtype(argType, paramType) {
		return
	}
	callName := tvCallNodeName(callNode)
	msg := fmt.Sprintf("actual parameter %d of %s does not match formal parameter", ordinal, callName)
	tv.tvMismatch(n, msg, argType, paramType)
}

// expectCanCast fires JSC_INVALID_CAST when `sourceType` cannot be cast
// to `targetType`. Mirrors Java's TypeValidator.expectCanCast.
func (tv *TypeValidator) expectCanCast(n *rhino.Node, targetType, sourceType jstype.JSType) {
	if !tvCanCastTo(sourceType, targetType) {
		tv.emit(diag.Diagnostic{
			Path:  tv.path,
			Line:  n.Lineno(),
			Col:   n.Charno(),
			Code:  diag.JSC_INVALID_CAST,
			Level: diag.Warning,
			Message: fmt.Sprintf(
				"invalid cast - must be a subtype or supertype\nfrom: %s\nto  : %s",
				tvTypeName(sourceType), tvTypeName(targetType)),
		})
		tv.mismatches.registerMismatch(n, sourceType, targetType)
		return
	}
	// Valid cast but still potentially informative for disambiguation.
	if !jstype.IsSubtype(sourceType, targetType) {
		tv.mismatches.registerMismatch(n, sourceType, targetType)
	}
}

// ---- private helpers ------------------------------------------------- //

// tvMismatch is the internal mismatch dispatcher. Mirrors Java's private
// mismatch(Node, String, JSType, JSType) method. Computes structural
// property details and delegates to tvRegisterMismatchAndReportWithDiag.
func (tv *TypeValidator) tvMismatch(n *rhino.Node, msg string, found, required jstype.JSType) {
	if tvIsSubtypeRespectingMode(found, required, tv.subtypingMode) {
		return
	}
	// Structural mismatch analysis: collect missing / mismatched properties
	// when the required type is structural (a @record/@interface).
	var missing, mismatch []string
	if reqObj, ok := required.(*jstype.ObjectType); ok && reqObj.IsStructural {
		if foundObj, ok := found.(*jstype.ObjectType); ok {
			for _, pname := range reqObj.MemberNames() {
				reqProp := reqObj.GetPropertyType(pname)
				if foundObj.HasProperty(pname) {
					foundProp := foundObj.GetPropertyType(pname)
					if !tvIsSubtypeRespectingMode(foundProp, reqProp, tv.subtypingMode) {
						mismatch = append(mismatch, pname)
					}
				} else if !containsVoidTV(reqProp) {
					missing = append(missing, pname)
				}
			}
			sort.Strings(missing)
			sort.Strings(mismatch)
		}
	}
	tv.tvRegisterMismatchAndReportWithDiag(n, diag.JSC_TYPE_MISMATCH, msg, found, required, missing, mismatch)
}

// tvRegisterMismatchAndReportWithDiag builds and emits a diagnostic,
// then registers the mismatch in the accumulator.
// Mirrors Java's registerMismatchAndReport(Node, DiagnosticType, String, …).
func (tv *TypeValidator) tvRegisterMismatchAndReportWithDiag(n *rhino.Node, code string, msg string, found, required jstype.JSType, missing, mismatch []string) {
	formatted := tvFormatFoundRequired(msg, found, required, missing, mismatch)
	tv.emit(diag.Diagnostic{
		Path:    tv.path,
		Line:    n.Lineno(),
		Col:     n.Charno(),
		Code:    code,
		Level:   diag.Error,
		Message: formatted,
	})
	tv.mismatches.registerMismatch(n, found, required)
}

// tvFormatFoundRequired mirrors Java's formatFoundRequired: builds the
// "found : T\nrequired: U" style message, plus optional missing/mismatch
// property lists when present.
func tvFormatFoundRequired(description string, found, required jstype.JSType, missing, mismatch []string) string {
	foundStr := tvTypeName(found)
	requiredStr := tvTypeName(required)
	missingStr := strings.Join(missing, ",")
	mismatchStr := strings.Join(mismatch, ",")
	if missingStr != "" || mismatchStr != "" {
		return fmt.Sprintf("%s\nfound   : %s\nrequired: %s\nmissing : [%s]\nmismatch: [%s]",
			description, foundStr, requiredStr, missingStr, mismatchStr)
	}
	return fmt.Sprintf("%s\nfound   : %s\nrequired: %s",
		description, foundStr, requiredStr)
}

// emit appends a diagnostic to the validator's list.
func (tv *TypeValidator) emit(d diag.Diagnostic) {
	tv.diags = append(tv.diags, d)
}

// ---- local helpers prefixed tv* ------------------------------------- //
// These mirror internal Java methods but are local to this file.

// tvTypeName returns the diagnostic string for a type.
func tvTypeName(t jstype.JSType) string {
	if t == nil {
		return "?"
	}
	return t.TypeName()
}

// tvIsUnknownType mirrors JSType.isUnknownType(): true for Unknown and All.
func tvIsUnknownType(t jstype.JSType) bool {
	return t == jstype.Unknown || t == jstype.All
}

// tvIsNoType mirrors JSType.isNoType(): true for nil (the "no type" sentinel).
func tvIsNoType(t jstype.JSType) bool {
	return t == nil
}

// tvIsEmptyType mirrors JSType.isEmptyType() — no type / bottom type.
func tvIsEmptyType(t jstype.JSType) bool {
	return t == nil
}

// tvIsObject mirrors JSType.isObject(): true for ObjectType/FunctionType.
func tvIsObject(t jstype.JSType) bool {
	if t == nil {
		return false
	}
	switch t.(type) {
	case *jstype.ObjectType, *jstype.FunctionType, *jstype.TemplatizedType:
		return true
	}
	return false
}

// tvMatchesObjectContext mirrors JSType.matchesObjectContext(): primitives
// can box to Object, so string/number/boolean qualify, as do object types.
func tvMatchesObjectContext(t jstype.JSType) bool {
	if t == nil {
		return false
	}
	if tvIsUnknownType(t) {
		return true
	}
	switch t {
	case jstype.String, jstype.Number, jstype.Boolean:
		return true // can box
	}
	return tvIsObject(t)
}

// tvMatchesStringContext mirrors JSType.matchesStringContext(): types that
// can be used where a string is expected.
func tvMatchesStringContext(t jstype.JSType) bool {
	if t == nil {
		return false
	}
	if tvIsUnknownType(t) {
		return true
	}
	if t == jstype.String || t == jstype.Number || t == jstype.Boolean ||
		t == jstype.Null || t == jstype.Void {
		return true
	}
	if u, ok := t.(*jstype.UnionType); ok {
		for _, m := range u.Members() {
			if !tvMatchesStringContext(m) {
				return false
			}
		}
		return true
	}
	return false
}

// tvMatchesNumberContext mirrors JSType.matchesNumberContext(): types that
// can be used where a number is expected.
func tvMatchesNumberContext(t jstype.JSType) bool {
	if t == nil {
		return false
	}
	if tvIsUnknownType(t) {
		return true
	}
	if t == jstype.Number || t == jstype.Boolean ||
		t == jstype.Null || t == jstype.Void {
		return true
	}
	if u, ok := t.(*jstype.UnionType); ok {
		for _, m := range u.Members() {
			if !tvMatchesNumberContext(m) {
				return false
			}
		}
		return true
	}
	if ee, ok := t.(*jstype.EnumElementType); ok {
		return tvMatchesNumberContext(ee.Enum.Underlying)
	}
	return false
}

// tvMatchesSymbolContext mirrors JSType.matchesSymbolContext(): in jsgo we
// have no symbol type, so this always returns false (conservative).
func tvMatchesSymbolContext(_ jstype.JSType) bool {
	return false
}

// tvIsSymbolValueType mirrors JSType.isSymbolValueType(): always false in
// jsgo since symbols are not modeled.
func tvIsSymbolValueType(_ jstype.JSType) bool {
	return false
}

// tvDereferencedObject mirrors JSType.dereference(): returns the underlying
// ObjectType for object and templatized types.
func tvDereferencedObject(t jstype.JSType) *jstype.ObjectType {
	if t == nil {
		return nil
	}
	if tt, ok := t.(*jstype.TemplatizedType); ok {
		t = tt.Base
	}
	if obj, ok := t.(*jstype.ObjectType); ok {
		return obj
	}
	return nil
}

// tvCanTestForShallowEqualityWith mirrors JSType.canTestForShallowEqualityWith:
// two types are testable for equality when they share at least one possible
// value — i.e. they are not disjoint. In our simplified lattice, two types
// are disjoint when both are concrete and different primitive singletons.
func tvCanTestForShallowEqualityWith(a, b jstype.JSType) bool {
	if tvIsUnknownType(a) || tvIsUnknownType(b) {
		return true
	}
	if tvIsNoType(a) || tvIsNoType(b) {
		return true
	}
	// Collect concrete members of each side, then check for overlap.
	aLeaves := tvLeaves(a)
	bLeaves := tvLeaves(b)
	for _, x := range aLeaves {
		for _, y := range bLeaves {
			if tvSameConcreteType(x, y) {
				return true
			}
		}
	}
	return false
}

// tvLeaves returns the flat list of non-union members in `t`.
func tvLeaves(t jstype.JSType) []jstype.JSType {
	if u, ok := t.(*jstype.UnionType); ok {
		var out []jstype.JSType
		for _, m := range u.Members() {
			out = append(out, tvLeaves(m)...)
		}
		return out
	}
	return []jstype.JSType{t}
}

// tvSameConcreteType reports whether two non-union types represent the
// same leaf type.
func tvSameConcreteType(a, b jstype.JSType) bool {
	if a == b {
		return true
	}
	// Objects of any kind can share values (structural compatibility).
	if tvIsObject(a) && tvIsObject(b) {
		return true
	}
	// One side is an object/function: it can hold any value at runtime.
	if tvIsObject(a) || tvIsObject(b) {
		return true
	}
	// EnumElementType shares values with its Underlying.
	if ea, ok := a.(*jstype.EnumElementType); ok {
		return tvSameConcreteType(ea.Enum.Underlying, b)
	}
	if eb, ok := b.(*jstype.EnumElementType); ok {
		return tvSameConcreteType(a, eb.Enum.Underlying)
	}
	return false
}

// tvCanCastTo mirrors JSType.canCastTo(JSType): a cast is valid when the
// source is a subtype of target, or target is a subtype of source (i.e.
// one side can be the other at runtime).
func tvCanCastTo(source, target jstype.JSType) bool {
	if tvIsUnknownType(source) || tvIsUnknownType(target) {
		return true
	}
	return jstype.IsSubtype(source, target) || jstype.IsSubtype(target, source)
}

// tvIsSubtypeRespectingMode mirrors JSType.isSubtypeOf(JSType, SubtypingMode):
// in IGNORE_NULL_UNDEFINED mode, strips null/undefined before checking.
func tvIsSubtypeRespectingMode(a, b jstype.JSType, mode tvSubtypingMode) bool {
	if mode == tvSubtypingIgnoreNullUndefined {
		a = tvStripNullUndefined(a)
		b = tvStripNullUndefined(b)
	}
	return jstype.IsSubtype(a, b)
}

// tvStripNullUndefined removes Null and Void from a union; returns `t`
// unchanged for non-union types.
func tvStripNullUndefined(t jstype.JSType) jstype.JSType {
	if t == jstype.Null || t == jstype.Void {
		return jstype.Unknown
	}
	u, ok := t.(*jstype.UnionType)
	if !ok {
		return t
	}
	var members []jstype.JSType
	for _, m := range u.Members() {
		if m != jstype.Null && m != jstype.Void {
			members = append(members, m)
		}
	}
	if len(members) == 0 {
		return jstype.Unknown
	}
	return jstype.NewUnion(members...)
}

// tvCallNodeName returns a human-readable name for the callee of a CALL node.
// Mirrors how Java's typeRegistry.getReadableTypeNameNoDeref is used in
// expectArgumentMatchesParameter.
func tvCallNodeName(callNode *rhino.Node) string {
	if callNode == nil {
		return "<unknown>"
	}
	first := callNode.FirstChild()
	if first == nil {
		return "<unknown>"
	}
	// NAME node → just use the string
	if first.Token() == rhino.NAME {
		return first.String_()
	}
	// GETPROP → "obj.prop"
	if first.Token() == rhino.GETPROP {
		kids := first.Children()
		if len(kids) >= 2 {
			recv := kids[0].String_()
			prop := kids[1].String_()
			if recv != "" && prop != "" {
				return recv + "." + prop
			}
		}
	}
	return "<unknown>"
}

// containsVoidTV mirrors ObjectType.containsVoid — checks whether a type
// includes Void (undefined), indicating an optional property.
func containsVoidTV(t jstype.JSType) bool {
	if t == jstype.Void {
		return true
	}
	if u, ok := t.(*jstype.UnionType); ok {
		for _, m := range u.Members() {
			if containsVoidTV(m) {
				return true
			}
		}
	}
	return false
}
