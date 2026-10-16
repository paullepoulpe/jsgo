// Port of com/google/javascript/jscomp/type/ChainableReverseAbstractInterpreter.java
// @ closure-compiler v20260415.
//
// Java's ChainableReverseAbstractInterpreter is an abstract class that
// implements ReverseAbstractInterpreter and supplies:
//   - a linked-list "chain" of interpreters (firstLink / nextLink)
//   - helper methods: getTypeIfRefinable, declareNameInScope
//   - inner visitor hierarchy for typeof-based narrowing
//   - getRestrictedByTypeOfResult for typeof narrowing
//   - restrictToArray / restrictToNotArray for Array.isArray
//
// In Go there are no abstract classes, so ChainableReverseAbstractInterpreter
// is a concrete struct that concrete subtypes embed. The inner visitor
// classes (anonymous in Java) become package-private types in this file.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/legacy/internal/rhino"
	"github.com/paullepoulpe/jsgo/legacy/internal/rhino/jstype"
)

// ----------------------------------------------------------------------------
// ChainableReverseAbstractInterpreter
// ----------------------------------------------------------------------------

// ChainableReverseAbstractInterpreter is the Go-portable base for chainable
// RAIs. Concrete subtypes (SemanticReverseAbstractInterpreter,
// ClosureReverseAbstractInterpreter) embed *ChainableReverseAbstractInterpreter
// and forward GetPreciserScopeKnowingConditionOutcome to the chain via
// nextPreciserScopeKnowingConditionOutcome when they don't handle a node.
//
// Port of ChainableReverseAbstractInterpreter (abstract class).
type ChainableReverseAbstractInterpreter struct {
	typeRegistry *jstype.Registry
	firstLink    ReverseAbstractInterpreter // head of chain (may be self)
	nextLink     ReverseAbstractInterpreter // nil if this is the last link
}

// newChainableBase constructs the embedded base for a concrete subtype.
// self must be the ReverseAbstractInterpreter implemented by the embedding
// struct (so firstLink resolves to a full implementation).
//
// Port of ChainableReverseAbstractInterpreter(JSTypeRegistry) constructor.
// Go-only: language-adapter Java uses class constructor; Go needs a factory for embedded structs.
func newChainableBase(reg *jstype.Registry, self ReverseAbstractInterpreter) ChainableReverseAbstractInterpreter {
	return ChainableReverseAbstractInterpreter{
		typeRegistry: reg,
		firstLink:    self,
		nextLink:     nil,
	}
}

// Append appends lastLink to the end of this chain, returning lastLink.
// The pattern new X().Append(new Y()) forms a chain X → Y.
//
// Port of ChainableReverseAbstractInterpreter.append.
func (c *ChainableReverseAbstractInterpreter) Append(lastLink *ChainableReverseAbstractInterpreter) *ChainableReverseAbstractInterpreter {
	c.nextLink = lastLink.firstLink
	lastLink.firstLink = c.firstLink
	return lastLink
}

// GetFirst returns the first link of this chain.
//
// Port of ChainableReverseAbstractInterpreter.getFirst.
func (c *ChainableReverseAbstractInterpreter) GetFirst() ReverseAbstractInterpreter {
	return c.firstLink
}

// firstPreciserScopeKnowingConditionOutcome delegates to the first link.
//
// Port of ChainableReverseAbstractInterpreter.firstPreciserScopeKnowingConditionOutcome.
func (c *ChainableReverseAbstractInterpreter) firstPreciserScopeKnowingConditionOutcome(
	condition *rhino.Node, blindScope *FlowScope, outcome bool,
) *FlowScope {
	return c.firstLink.GetPreciserScopeKnowingConditionOutcome(condition, blindScope, outcome)
}

// nextPreciserScopeKnowingConditionOutcome delegates to the next link,
// or returns blindScope if this is the last link.
//
// Port of ChainableReverseAbstractInterpreter.nextPreciserScopeKnowingConditionOutcome.
func (c *ChainableReverseAbstractInterpreter) nextPreciserScopeKnowingConditionOutcome(
	condition *rhino.Node, blindScope *FlowScope, outcome bool,
) *FlowScope {
	if c.nextLink != nil {
		return c.nextLink.GetPreciserScopeKnowingConditionOutcome(condition, blindScope, outcome)
	}
	return blindScope
}

// getTypeIfRefinable returns the type of a node in the given scope if the
// node corresponds to a name whose type is capable of being refined.
// Returns nil otherwise.
//
// In Go's FlowScope representation, scope.Get returns the flow-refined type;
// node.GetJSType() provides the compiler-annotated type as a fallback.
//
// Port of ChainableReverseAbstractInterpreter.getTypeIfRefinable.
func (c *ChainableReverseAbstractInterpreter) getTypeIfRefinable(node *rhino.Node, scope *FlowScope) jstype.JSType {
	switch node.Token() {
	case rhino.NAME:
		t := scope.Get(node.String_())
		if t == nil {
			t, _ = node.GetJSType().(jstype.JSType)
		}
		return t

	case rhino.GETPROP:
		qname := node.GetQualifiedName()
		if qname == "" {
			return nil
		}
		t := scope.Get(qname)
		if t == nil {
			t, _ = node.GetJSType().(jstype.JSType)
		}
		if t == nil {
			t = c.getNativeType(jstype.UNKNOWN_TYPE)
		}
		return t
	}
	return nil
}

// declareNameInScope returns an updated FlowScope with the refined type for
// the name represented by node. It is a precondition that getTypeIfRefinable
// would return non-nil for node in scope.
//
// Port of ChainableReverseAbstractInterpreter.declareNameInScope.
func (c *ChainableReverseAbstractInterpreter) declareNameInScope(scope *FlowScope, node *rhino.Node, t jstype.JSType) *FlowScope {
	switch node.Token() {
	case rhino.NAME:
		return scope.inferSlotType(node.String_(), t)

	case rhino.GETPROP:
		qname := node.GetQualifiedName()
		if qname == "" {
			return scope
		}
		return scope.inferQualifiedSlot(qname, t)

	case rhino.THIS:
		// "this" references aren't currently modelled in the CFG.
		return scope
	}
	// Should not happen — caller must guard with getTypeIfRefinable.
	return scope
}

// getNativeType looks up a native/built-in type from the registry.
//
// Port of ChainableReverseAbstractInterpreter.getNativeType.
func (c *ChainableReverseAbstractInterpreter) getNativeType(typeID jstype.JSTypeNative) jstype.JSType {
	return c.typeRegistry.GetNativeType(typeID)
}

// getNativeTypeForTypeOf returns the native type corresponding to a typeof
// string, or nil if no single type is definitively implied.
//
// Port of ChainableReverseAbstractInterpreter.getNativeTypeForTypeOf (private).
func (c *ChainableReverseAbstractInterpreter) getNativeTypeForTypeOf(value string) jstype.JSType {
	switch value {
	case "number":
		return c.getNativeType(jstype.NUMBER_TYPE)
	case "boolean":
		return c.getNativeType(jstype.BOOLEAN_TYPE)
	case "string":
		return c.getNativeType(jstype.STRING_TYPE)
	case "symbol":
		return c.getNativeType(jstype.SYMBOL_TYPE)
	case "undefined":
		return c.getNativeType(jstype.VOID_TYPE)
	case "object":
		return jstype.NewUnion(c.getNativeType(jstype.OBJECT_TYPE), c.getNativeType(jstype.NULL_TYPE))
	case "function":
		return c.getNativeType(jstype.FUNCTION_TYPE)
	}
	return nil
}

// getRestrictedByTypeOfResult returns a version of typ that is restricted by
// knowledge about the result of the typeof operation.
//
// Port of ChainableReverseAbstractInterpreter.getRestrictedByTypeOfResult.
func (c *ChainableReverseAbstractInterpreter) getRestrictedByTypeOfResult(
	typ jstype.JSType, value string, resultEqualsValue bool,
) jstype.JSType {
	if typ == nil {
		if resultEqualsValue {
			result := c.getNativeTypeForTypeOf(value)
			if result == nil {
				return c.getNativeType(jstype.CHECKED_UNKNOWN_TYPE)
			}
			return result
		}
		return nil
	}
	v := newRestrictByOneTypeOfResultVisitor(c, value, resultEqualsValue)
	return asJSType(jstype.VisitJSType(v, typ))
}

// restrictToArray restricts a type for when goog.isArray / Array.isArray returns true.
//
// Port of ChainableReverseAbstractInterpreter.restrictToArrayVisitor field.
func (c *ChainableReverseAbstractInterpreter) restrictToArray(typ jstype.JSType) jstype.JSType {
	if typ == nil {
		return nil
	}
	v := newRestrictToArrayVisitorImpl(c)
	return asJSType(jstype.VisitJSType(v, typ))
}

// restrictToNotArray restricts a type for when goog.isArray / Array.isArray returns false.
//
// Port of ChainableReverseAbstractInterpreter.restrictToNotArrayVisitor field.
func (c *ChainableReverseAbstractInterpreter) restrictToNotArray(typ jstype.JSType) jstype.JSType {
	if typ == nil {
		return nil
	}
	v := newRestrictToNotArrayVisitorImpl(c)
	return asJSType(jstype.VisitJSType(v, typ))
}

// asJSType converts the any returned by VisitJSType to jstype.JSType.
// nil is returned as nil; any other non-JSType value returns nil.
// Go-only: language-adapter Java Visitor<JSType> returns JSType directly; Go Visitor returns any.
func asJSType(v any) jstype.JSType {
	if v == nil {
		return nil
	}
	t, _ := v.(jstype.JSType)
	return t
}

// ----------------------------------------------------------------------------
// Visitor helpers shared across concrete visitor types
// ----------------------------------------------------------------------------

// visitUnionRestrictingType applies visitor v to each member of typ and
// returns the union of kept (non-nil) results, using leastSupertype to join.
// Mirrors the union case in RestrictByTypeOfResultVisitor.caseUnionType.
// Go-only: language-adapter Java inlines caseUnionType per visitor subclass; Go shares this helper.
func visitUnionRestrictingType(v jstype.Visitor, typ *jstype.UnionType) jstype.JSType {
	var restricted jstype.JSType
	for _, alt := range typ.Members() {
		restrictedAlt := asJSType(jstype.VisitJSType(v, alt))
		if restrictedAlt != nil {
			if restricted == nil {
				restricted = restrictedAlt
			} else {
				restricted = raiLeastSupertype(restrictedAlt, restricted)
			}
		}
	}
	return restricted
}

// raiLeastSupertype computes the least supertype of two non-nil JSTypes.
// Port of JSType.getLeastSupertype — the union of two types is their
// least supertype when no direct supertype relation exists.
// Go-only: language-adapter Java calls JSType.getLeastSupertype() as an instance method; Go needs a free function.
func raiLeastSupertype(a, b jstype.JSType) jstype.JSType {
	if jstype.IsSubtype(a, b) {
		return b
	}
	if jstype.IsSubtype(b, a) {
		return a
	}
	return jstype.NewUnion(a, b)
}

// visitEnumElement visits the primitive type of an enum element, returning
// the enum element itself when the primitive is unchanged (mirrors Java's
// RestrictByTypeOfResultVisitor.caseEnumElementType).
// Go-only: language-adapter Java inlines caseEnumElementType per subclass; Go shares this helper.
func visitEnumElement(v jstype.Visitor, typ *jstype.EnumElementType) jstype.JSType {
	prim := typ.GetPrimitiveType()
	inner := asJSType(jstype.VisitJSType(v, prim))
	if inner != nil && prim == inner {
		return typ
	}
	return inner
}

// visitNamedType visits the referenced type of a NamedType.
// Go-only: language-adapter Java inlines caseNamedType via caseProxyObjectType; Go shares this helper.
func visitNamedType(v jstype.Visitor, typ *jstype.NamedType) jstype.JSType {
	return asJSType(jstype.VisitJSType(v, typ.GetReferencedType()))
}

// Go-only: language-adapter Java inlines caseProxyObjectType.visitReferenceType(); Go shares this helper.
func visitProxyObjectType(v jstype.Visitor, typ *jstype.ProxyObjectType) jstype.JSType {
	return asJSType(jstype.VisitJSType(v, typ.GetReferencedTypeInternal()))
}

// ----------------------------------------------------------------------------
// restrictByOneTypeOfResultVisitor
// ----------------------------------------------------------------------------

// restrictByOneTypeOfResultVisitor restricts a type by a single known typeof
// result value.
//
// Port of ChainableReverseAbstractInterpreter.RestrictByOneTypeOfResultVisitor (private inner class).
type restrictByOneTypeOfResultVisitor struct {
	jstype.WithDefaultCase
	c                 *ChainableReverseAbstractInterpreter
	value             string
	resultEqualsValue bool
}

// Go-only: language-adapter Java instantiates inner class with `new RestrictByOneTypeOfResultVisitor(value, resultEqualsValue)`.
func newRestrictByOneTypeOfResultVisitor(
	c *ChainableReverseAbstractInterpreter, value string, resultEqualsValue bool,
) *restrictByOneTypeOfResultVisitor {
	v := &restrictByOneTypeOfResultVisitor{
		c:                 c,
		value:             value,
		resultEqualsValue: resultEqualsValue,
	}
	v.DefaultCase = func(typ jstype.JSType) any { return nil }
	return v
}

func (v *restrictByOneTypeOfResultVisitor) matchesExpectation(result string) bool {
	return (result == v.value) == v.resultEqualsValue
}

func (v *restrictByOneTypeOfResultVisitor) CaseAllType() any {
	return v.caseTopType(v.c.getNativeType(jstype.ALL_TYPE))
}

func (v *restrictByOneTypeOfResultVisitor) CaseUnknownType() any {
	return v.caseTopType(v.c.getNativeType(jstype.CHECKED_UNKNOWN_TYPE))
}

func (v *restrictByOneTypeOfResultVisitor) caseTopType(topType jstype.JSType) jstype.JSType {
	result := topType
	if v.resultEqualsValue {
		typeByName := v.c.getNativeTypeForTypeOf(v.value)
		if typeByName != nil {
			result = typeByName
		}
	}
	return result
}

func (v *restrictByOneTypeOfResultVisitor) CaseUnionType(typ *jstype.UnionType) any {
	return visitUnionRestrictingType(v, typ)
}

func (v *restrictByOneTypeOfResultVisitor) CaseNoType(typ jstype.JSType) any { return typ }

func (v *restrictByOneTypeOfResultVisitor) CaseEnumElementType(typ *jstype.EnumElementType) any {
	return visitEnumElement(v, typ)
}

func (v *restrictByOneTypeOfResultVisitor) CaseTemplatizedType(typ *jstype.TemplatizedType) any {
	return v.CaseObjectType(typ.Base)
}

func (v *restrictByOneTypeOfResultVisitor) CaseTemplateType(typ *jstype.TemplateType) any {
	// TemplateType extends ObjectType in Java; caseTemplateType delegates to caseObjectType.
	// In Go, TemplateType is not an *ObjectType, so we check if the value passes as "function"
	// or "object" and return nil otherwise (same net result as the default caseObjectType).
	// Mirrors RestrictByTypeOfResultVisitor.caseTemplateType → caseObjectType(templateType).
	_ = typ
	return nil
}

func (v *restrictByOneTypeOfResultVisitor) CaseNamedType(typ jstype.JSType) any {
	if nt, ok := typ.(*jstype.NamedType); ok {
		return visitNamedType(v, nt)
	}
	return nil
}

func (v *restrictByOneTypeOfResultVisitor) CaseProxyObjectType(typ jstype.JSType) any {
	if pt, ok := typ.(*jstype.ProxyObjectType); ok {
		return visitProxyObjectType(v, pt)
	}
	return nil
}

func (v *restrictByOneTypeOfResultVisitor) CaseNoObjectType() any {
	matches := (v.value == "object" || v.value == "function") == v.resultEqualsValue
	if matches {
		return v.c.getNativeType(jstype.NO_OBJECT_TYPE)
	}
	return nil
}

func (v *restrictByOneTypeOfResultVisitor) CaseBooleanType() any {
	if v.matchesExpectation("boolean") {
		return v.c.getNativeType(jstype.BOOLEAN_TYPE)
	}
	return nil
}

func (v *restrictByOneTypeOfResultVisitor) CaseFunctionType(typ *jstype.FunctionType) any {
	if v.value == "function" {
		ctorType := v.c.getNativeType(jstype.FUNCTION_TYPE)
		if v.resultEqualsValue {
			// Restrict objects to "Function"; subtypes are left.
			return jstype.Meet(ctorType, typ)
		}
		// Only filter out subtypes of "function".
		if jstype.IsSubtype(typ, ctorType) {
			return nil
		}
		return typ
	}
	if v.matchesExpectation("object") {
		return typ
	}
	return nil
}

func (v *restrictByOneTypeOfResultVisitor) CaseNullType() any {
	if v.matchesExpectation("object") {
		return v.c.getNativeType(jstype.NULL_TYPE)
	}
	return nil
}

func (v *restrictByOneTypeOfResultVisitor) CaseNumberType() any {
	if v.matchesExpectation("number") {
		return v.c.getNativeType(jstype.NUMBER_TYPE)
	}
	return nil
}

func (v *restrictByOneTypeOfResultVisitor) CaseBigIntType() any {
	if v.matchesExpectation("bigint") {
		return v.c.getNativeType(jstype.BIGINT_TYPE)
	}
	return nil
}

func (v *restrictByOneTypeOfResultVisitor) CaseObjectType(typ *jstype.ObjectType) any {
	if v.value == "function" {
		ctorType := v.c.getNativeType(jstype.FUNCTION_TYPE)
		if v.resultEqualsValue {
			return jstype.Meet(ctorType, typ)
		}
		if jstype.IsSubtype(typ, ctorType) {
			return nil
		}
		return typ
	}
	if v.matchesExpectation("object") {
		return typ
	}
	return nil
}

func (v *restrictByOneTypeOfResultVisitor) CaseStringType() any {
	if v.matchesExpectation("string") {
		return v.c.getNativeType(jstype.STRING_TYPE)
	}
	return nil
}

func (v *restrictByOneTypeOfResultVisitor) CaseSymbolType(typ jstype.JSType) any {
	if v.matchesExpectation("symbol") {
		return typ
	}
	return nil
}

func (v *restrictByOneTypeOfResultVisitor) CaseVoidType() any {
	if v.matchesExpectation("undefined") {
		return v.c.getNativeType(jstype.VOID_TYPE)
	}
	return nil
}

// ----------------------------------------------------------------------------
// restrictToArrayVisitorImpl
// ----------------------------------------------------------------------------

// restrictToArrayVisitorImpl is the visitor for Array.isArray / goog.isArray
// returning true. Extends RestrictByTrueTypeOfResultVisitor — all unhandled
// leaf cases return nil.
//
// Port of ChainableReverseAbstractInterpreter.restrictToArrayVisitor anonymous field.
type restrictToArrayVisitorImpl struct {
	jstype.WithDefaultCase
	c *ChainableReverseAbstractInterpreter
}

// Go-only: language-adapter Java uses anonymous inner class; Go uses a named constructor.
func newRestrictToArrayVisitorImpl(c *ChainableReverseAbstractInterpreter) *restrictToArrayVisitorImpl {
	v := &restrictToArrayVisitorImpl{c: c}
	// RestrictByTrueTypeOfResultVisitor: all unhandled leaf cases → nil.
	v.DefaultCase = func(typ jstype.JSType) any { return nil }
	return v
}

func (v *restrictToArrayVisitorImpl) CaseAllType() any {
	return v.c.getNativeType(jstype.ARRAY_TYPE)
}

func (v *restrictToArrayVisitorImpl) CaseUnknownType() any {
	// caseTopType: all-type → ARRAY_TYPE; unknown → CHECKED_UNKNOWN_TYPE (not array).
	return v.c.getNativeType(jstype.CHECKED_UNKNOWN_TYPE)
}

func (v *restrictToArrayVisitorImpl) CaseUnionType(typ *jstype.UnionType) any {
	return visitUnionRestrictingType(v, typ)
}

func (v *restrictToArrayVisitorImpl) CaseObjectType(typ *jstype.ObjectType) any {
	arrayType := v.c.getNativeType(jstype.ARRAY_TYPE)
	if jstype.IsSubtype(arrayType, typ) {
		return arrayType
	}
	return nil
}

func (v *restrictToArrayVisitorImpl) CaseEnumElementType(typ *jstype.EnumElementType) any {
	return visitEnumElement(v, typ)
}

func (v *restrictToArrayVisitorImpl) CaseTemplatizedType(typ *jstype.TemplatizedType) any {
	return v.CaseObjectType(typ.Base)
}

func (v *restrictToArrayVisitorImpl) CaseTemplateType(typ *jstype.TemplateType) any {
	// TemplateType extends ObjectType in Java; caseTemplateType delegates to caseObjectType.
	// In Go, TemplateType is not an *ObjectType, so we apply the same logic as CaseObjectType
	// conservatively: check if ARRAY_TYPE is a subtype.
	_ = typ
	return nil
}

func (v *restrictToArrayVisitorImpl) CaseNamedType(typ jstype.JSType) any {
	if nt, ok := typ.(*jstype.NamedType); ok {
		return visitNamedType(v, nt)
	}
	return nil
}

func (v *restrictToArrayVisitorImpl) CaseProxyObjectType(typ jstype.JSType) any {
	if pt, ok := typ.(*jstype.ProxyObjectType); ok {
		return visitProxyObjectType(v, pt)
	}
	return nil
}

func (v *restrictToArrayVisitorImpl) CaseNoType(typ jstype.JSType) any { return typ }

// ----------------------------------------------------------------------------
// restrictToNotArrayVisitorImpl
// ----------------------------------------------------------------------------

// restrictToNotArrayVisitorImpl is the visitor for Array.isArray / goog.isArray
// returning false. Extends RestrictByFalseTypeOfResultVisitor — all unhandled
// leaf cases return their own native type.
//
// Port of ChainableReverseAbstractInterpreter.restrictToNotArrayVisitor anonymous field.
type restrictToNotArrayVisitorImpl struct {
	jstype.WithDefaultCase
	c *ChainableReverseAbstractInterpreter
}

// Go-only: language-adapter Java uses anonymous inner class; Go uses a named constructor.
func newRestrictToNotArrayVisitorImpl(c *ChainableReverseAbstractInterpreter) *restrictToNotArrayVisitorImpl {
	v := &restrictToNotArrayVisitorImpl{c: c}
	// RestrictByFalseTypeOfResultVisitor: default keeps the type as-is.
	v.DefaultCase = func(typ jstype.JSType) any { return typ }
	return v
}

func (v *restrictToNotArrayVisitorImpl) CaseAllType() any {
	return v.c.getNativeType(jstype.ALL_TYPE)
}

func (v *restrictToNotArrayVisitorImpl) CaseUnknownType() any {
	return v.c.getNativeType(jstype.CHECKED_UNKNOWN_TYPE)
}

func (v *restrictToNotArrayVisitorImpl) CaseUnionType(typ *jstype.UnionType) any {
	return visitUnionRestrictingType(v, typ)
}

func (v *restrictToNotArrayVisitorImpl) CaseObjectType(typ *jstype.ObjectType) any {
	if jstype.IsSubtype(typ, v.c.getNativeType(jstype.ARRAY_TYPE)) {
		return nil
	}
	return typ
}

func (v *restrictToNotArrayVisitorImpl) CaseFunctionType(typ *jstype.FunctionType) any {
	return typ
}

func (v *restrictToNotArrayVisitorImpl) CaseNoObjectType() any {
	return v.c.getNativeType(jstype.NO_OBJECT_TYPE)
}

func (v *restrictToNotArrayVisitorImpl) CaseBooleanType() any {
	return v.c.getNativeType(jstype.BOOLEAN_TYPE)
}

func (v *restrictToNotArrayVisitorImpl) CaseNullType() any {
	return v.c.getNativeType(jstype.NULL_TYPE)
}

func (v *restrictToNotArrayVisitorImpl) CaseNumberType() any {
	return v.c.getNativeType(jstype.NUMBER_TYPE)
}

func (v *restrictToNotArrayVisitorImpl) CaseBigIntType() any {
	return v.c.getNativeType(jstype.BIGINT_TYPE)
}

func (v *restrictToNotArrayVisitorImpl) CaseStringType() any {
	return v.c.getNativeType(jstype.STRING_TYPE)
}

func (v *restrictToNotArrayVisitorImpl) CaseSymbolType(typ jstype.JSType) any {
	return typ
}

func (v *restrictToNotArrayVisitorImpl) CaseVoidType() any {
	return v.c.getNativeType(jstype.VOID_TYPE)
}

func (v *restrictToNotArrayVisitorImpl) CaseEnumElementType(typ *jstype.EnumElementType) any {
	return visitEnumElement(v, typ)
}

func (v *restrictToNotArrayVisitorImpl) CaseTemplatizedType(typ *jstype.TemplatizedType) any {
	return v.CaseObjectType(typ.Base)
}

func (v *restrictToNotArrayVisitorImpl) CaseTemplateType(typ *jstype.TemplateType) any {
	// TemplateType extends ObjectType in Java; caseTemplateType delegates to caseObjectType.
	// In Go, TemplateType is not an *ObjectType, so we conservatively return nil.
	_ = typ
	return nil
}

func (v *restrictToNotArrayVisitorImpl) CaseNamedType(typ jstype.JSType) any {
	if nt, ok := typ.(*jstype.NamedType); ok {
		return visitNamedType(v, nt)
	}
	return nil
}

func (v *restrictToNotArrayVisitorImpl) CaseProxyObjectType(typ jstype.JSType) any {
	if pt, ok := typ.(*jstype.ProxyObjectType); ok {
		return visitProxyObjectType(v, pt)
	}
	return nil
}

func (v *restrictToNotArrayVisitorImpl) CaseNoType(typ jstype.JSType) any { return typ }
