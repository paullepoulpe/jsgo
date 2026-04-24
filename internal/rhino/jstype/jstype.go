// Port of com/google/javascript/rhino/jstype/JSType.java @ closure-compiler v20260415
//
// The Java hierarchy is an abstract class `JSType` with many subclasses.
// In Go we expose an interface and concrete types; subtype logic is
// centralised in isSubtypeOf. Primitive type definitions live in per-type
// files matching the Java layout (all_type.go, unknown_type.go, …).

package jstype

type JSType interface {
	// TypeName matches Java JSType.toString() for primitive types (e.g.
	// "number"). Not used for conformance, only for diagnostics.
	TypeName() string

	// IsSubtypeOf ports JSType.isSubtypeOf(JSType). Semantics match Java:
	// every type is a subtype of AllType and Unknown; primitives are
	// subtypes only of themselves.
	IsSubtypeOf(other JSType) bool
}

// IsSubtype is the package-level entry point; callers should prefer this
// over t.IsSubtypeOf so union-type handling stays in one place. Mirrors
// how JSType.isSubtype in Java dispatches through the registry.
func IsSubtype(x, y JSType) bool {
	if x == nil || y == nil {
		return true
	}
	// Unresolved template vars accept any value — mirrors Java's
	// "open type" rule during unification. Once the substitution
	// resolves at a CALL site, the concrete type is used instead,
	// and this branch stops applying.
	if _, ok := y.(*TemplateType); ok {
		return true
	}
	if _, ok := x.(*TemplateType); ok {
		return true
	}
	// TemplatizedType ↔ underlying nominal: a `Foo<X>` is
	// assignment-compatible with a bare `Foo` (Java's invariant
	// generic rule degrades the args when one side has none).
	if xt, ok := x.(*TemplatizedType); ok {
		if y == xt.Base {
			return true
		}
	}
	if yt, ok := y.(*TemplatizedType); ok {
		if x == yt.Base {
			return true
		}
	}
	// When the target is a union, X matches if it's a subtype of any
	// member. Left-side unions stay in UnionType.IsSubtypeOf.
	if u, ok := y.(*UnionType); ok {
		if _, isUnion := x.(*UnionType); isUnion {
			return x.IsSubtypeOf(u)
		}
		for _, m := range u.members {
			if IsSubtype(x, m) {
				return true
			}
		}
		return false
	}
	return x.IsSubtypeOf(y)
}

func isAll(t JSType) bool     { _, ok := t.(allType); return ok }
func isUnknown(t JSType) bool { _, ok := t.(unknownType); return ok }

// IsSymbol reports whether t is the symbol type (SymbolType from symbol_type.go).
func IsSymbol(t JSType) bool { _, ok := t.(SymbolType); return ok }

// Canonical interned instances. Java JSTypeRegistry returns the same
// instance for each primitive; we mirror that with package-level vars.
var (
	All     JSType = allType{}
	Unknown JSType = unknownType{}
	Number  JSType = numberType{}
	String  JSType = stringType{}
	Boolean JSType = booleanType{}
	Null    JSType = nullType{}
	Void    JSType = voidType{}
	// Symbol is declared in symbol_type.go.
)

// ToMaybeFunctionType mirrors JSType.toMaybeFunctionType(JSType) — null-
// safe downcast to *FunctionType. Returns nil when t is nil or not a
// FunctionType.
func ToMaybeFunctionType(t JSType) *FunctionType {
	if t == nil {
		return nil
	}
	f, _ := t.(*FunctionType)
	return f
}

// ToMaybeEnumType mirrors JSType.toMaybeEnumType() — null-safe downcast
// to *EnumType. Returns nil when t is nil or not an EnumType.
func ToMaybeEnumType(t JSType) *EnumType {
	if t == nil {
		return nil
	}
	e, _ := t.(*EnumType)
	return e
}

// IsRecordType mirrors JSType.isRecordType() — returns true when t is a
// *RecordType (the Go equivalent of Java's RecordType subclass). Plain
// ObjectType instances with IsStructural=true are @record/@interface shapes
// but are not RecordType in Java's sense; only actual *RecordType values
// return true here.
func IsRecordType(t JSType) bool {
	if t == nil {
		return false
	}
	type recordTyper interface{ IsRecordType() bool }
	if rt, ok := t.(recordTyper); ok {
		return rt.IsRecordType()
	}
	return false
}

// Dereference mirrors JSType.dereference() — returns the ObjectType when
// t is an *ObjectType or *FunctionType (which is also an object in the
// Java hierarchy). Returns nil for primitives, null, undefined, and
// unknown. Mirrors Java's autobox().toObjectType() at a simplified level.
func Dereference(t JSType) *ObjectType {
	if t == nil {
		return nil
	}
	if o, ok := t.(*ObjectType); ok {
		return o
	}
	// FunctionType is an ObjectType in the Java hierarchy — InferJsDocInfo
	// calls dereference() on the inferred type and expects to get an
	// ObjectType back for function expressions. Return nil here because our
	// FunctionType is not an ObjectType subtype; the caller must handle the
	// FunctionType case separately if needed.
	return nil
}

// Outcome represents the truthiness / nullishness of a branch outcome.
// Port of com.google.javascript.rhino.Outcome @ closure-compiler v20260415.
type Outcome int

const (
	// OutcomeTrue — truthy values (e.g. {}, true, 1).
	OutcomeTrue Outcome = iota
	// OutcomeFalse — all falsy values (includes nullish and non-null falsy).
	OutcomeFalse
	// OutcomeFalseNotNull — falsy, but definitely not null/undefined (e.g. 0, "").
	OutcomeFalseNotNull
	// OutcomeNullish — null or undefined.
	OutcomeNullish
)

// IsTruthy mirrors Outcome.isTruthy().
func (o Outcome) IsTruthy() bool { return o == OutcomeTrue }

// IsNullish mirrors Outcome.isNullish() → true only for OutcomeNullish, false
// for OutcomeFalseNotNull, and unknown (represented as -1 sentinel) for
// OutcomeFalse (which covers both nullish and non-null-falsy).
// Callers compare with OutcomeNullish / OutcomeFalseNotNull directly.
func (o Outcome) IsNullishBool() bool { return o == OutcomeNullish }

// TypePair holds a pair of restricted JSTypes produced by equality/inequality
// narrowing. Port of JSType.TypePair @ closure-compiler v20260415.
type TypePair struct {
	TypeA JSType // may be nil when the branch is impossible
	TypeB JSType // may be nil when the branch is impossible
}

// Meet returns the greatest common subtype (lattice meet / infimum) of x and y.
// Package-level wrapper over UnionGetGreatestSubtype / JSType.getGreatestSubtype.
// Port of JSType.getGreatestSubtype(JSType) @ closure-compiler v20260415.
func Meet(x, y JSType) JSType {
	if x == nil || y == nil {
		return None
	}
	if x == y {
		return x
	}
	if isUnknown(x) {
		return y
	}
	if isUnknown(y) {
		return x
	}
	if u, ok := x.(*UnionType); ok {
		return UnionGetGreatestSubtype(u, y)
	}
	if u, ok := y.(*UnionType); ok {
		return UnionGetGreatestSubtype(u, x)
	}
	if IsSubtype(x, y) {
		return x
	}
	if IsSubtype(y, x) {
		return y
	}
	// No common subtype — return the appropriate bottom.
	if _, xIsObj := x.(*ObjectType); xIsObj {
		if _, yIsObj := y.(*ObjectType); yIsObj {
			return NoObject
		}
	}
	return None
}

// IsStrictSubtypeOf reports x <: y with structural subtyping disabled.
// Port of JSType.isSubtypeWithoutStructuralTyping() @ closure-compiler v20260415.
// Implemented as: IsSubtype(x, y) && x != y (approximation without registry).
func IsStrictSubtypeOf(x, y JSType) bool {
	if x == nil || y == nil {
		return false
	}
	return IsSubtype(x, y) && x != y && !isUnknown(x) && !isUnknown(y)
}

// matchesObjectContexter is the optional interface for types that override
// the default matchesObjectContext() behaviour.
type matchesObjectContexter interface{ MatchesObjectContext() bool }

// MatchesObjectContext reports whether t can appear in an Object context.
// Port of JSType.matchesObjectContext() @ closure-compiler v20260415.
// Default: false for JSType; overridden to true by object types, boxable
// primitives (string, number, boolean, bigint, symbol), Unknown, AllType.
func MatchesObjectContext(t JSType) bool {
	if t == nil {
		return false
	}
	if m, ok := t.(matchesObjectContexter); ok {
		return m.MatchesObjectContext()
	}
	// Primitives that can be auto-boxed return true.
	switch t.(type) {
	case stringType, numberType, booleanType, bigIntType, SymbolType:
		return true
	case unknownType, allType:
		return true
	case *UnionType:
		u := t.(*UnionType)
		return AnyTypeMatchesUnion(MatchesObjectContext, u)
	}
	return false
}

// autoboxeser is the optional interface for types that override autoboxesTo().
type autoboxeser interface{ AutoboxesTo() JSType }

// AutoboxesTo returns the boxed Object type for scalar primitives, or nil
// for types that are not auto-boxed. Port of JSType.autoboxesTo().
func AutoboxesTo(t JSType) JSType {
	if t == nil {
		return nil
	}
	if a, ok := t.(autoboxeser); ok {
		return a.AutoboxesTo()
	}
	return nil
}

// Autobox dereferences t for property access: strips null/undefined and
// returns the auto-boxed type (if any), otherwise the type itself.
// Port of JSType.autobox() @ closure-compiler v20260415.
func Autobox(t JSType) JSType {
	if t == nil {
		return nil
	}
	restricted := RestrictByNotNullOrUndefined(t)
	if boxed := AutoboxesTo(restricted); boxed != nil {
		return boxed
	}
	return restricted
}

// RestrictByNotNullOrUndefined strips null and undefined from t.
// For non-union types that are neither null nor undefined, returns t unchanged.
// Port of JSType.restrictByNotNullOrUndefined() @ closure-compiler v20260415.
func RestrictByNotNullOrUndefined(t JSType) JSType {
	if t == nil {
		return nil
	}
	if u, ok := t.(*UnionType); ok {
		return u.RestrictByNotNullOrUndefined()
	}
	if isNullType(t) || isVoidType(t) {
		return None
	}
	return t
}

// collapseUnioner is the optional interface for union-like types.
type collapseUnioner interface{ CollapseUnion() JSType }

// CollapseUnion returns the least non-union supertype.
// For non-union types returns t unchanged.
// Port of JSType.collapseUnion() @ closure-compiler v20260415.
func CollapseUnion(t JSType) JSType {
	if t == nil {
		return nil
	}
	if c, ok := t.(collapseUnioner); ok {
		return c.CollapseUnion()
	}
	return t
}

// ownerFunctioner is the optional interface for prototype / proxy object types.
type ownerFunctioner interface{ GetOwnerFunction() *FunctionType }

// GetOwnerFunction returns the FunctionType whose prototype this object is,
// or nil when t is not a prototype. Port of ObjectType.getOwnerFunction().
func GetOwnerFunction(t JSType) *FunctionType {
	if t == nil {
		return nil
	}
	if o, ok := t.(ownerFunctioner); ok {
		return o.GetOwnerFunction()
	}
	return nil
}

// possibleToBooleaner is the optional interface for types that provide
// getPossibleToBooleanOutcomes.
type possibleToBooleaner interface{ GetPossibleToBooleanOutcomes() BooleanLiteralSet }

// GetPossibleToBooleanOutcomes returns the set of boolean outcomes for t.
// Port of JSType.getPossibleToBooleanOutcomes() @ closure-compiler v20260415.
func GetPossibleToBooleanOutcomes(t JSType) BooleanLiteralSet {
	if t == nil {
		return BooleanLiteralSetBoth
	}
	if p, ok := t.(possibleToBooleaner); ok {
		return p.GetPossibleToBooleanOutcomes()
	}
	// Primitives: null and void are always false; objects always true.
	switch t.(type) {
	case nullType, voidType:
		return BooleanLiteralSetFalse
	case *ObjectType, *FunctionType, *PrototypeObjectType:
		return BooleanLiteralSetTrue
	}
	// Unknown / All / other: both outcomes possible.
	return BooleanLiteralSetBoth
}

// GetRestrictedTypeGivenOutcome narrows t based on a boolean outcome.
// Port of JSType.getRestrictedTypeGivenOutcome(Outcome) @ closure-compiler v20260415.
func GetRestrictedTypeGivenOutcome(t JSType, outcome Outcome) JSType {
	if t == nil {
		return None
	}
	// Union: map across alternates.
	if u, ok := t.(*UnionType); ok {
		mapped := MapTypes(func(m JSType) JSType {
			return GetRestrictedTypeGivenOutcome(m, outcome)
		}, u.Members())
		return NewUnion(mapped...)
	}
	// Truthy unknown → checked unknown (approximated as Unknown here).
	if outcome.IsTruthy() && isUnknown(t) {
		return Unknown
	}
	// Nullish branch: only null/undefined survive.
	if outcome.IsNullishBool() {
		if isNullType(t) || isVoidType(t) {
			return t
		}
		return None
	}
	// General: keep t iff it can produce the requested boolean.
	literals := GetPossibleToBooleanOutcomes(t)
	if literals.Contains(outcome.IsTruthy()) {
		return t
	}
	return None
}

// GetTypesUnderEquality returns the pair of types that survive an == check
// between t and that. Port of JSType.getTypesUnderEquality(JSType).
func GetTypesUnderEquality(t, that JSType) TypePair {
	if that == nil {
		return TypePair{t, that}
	}
	// Union on the right: flip.
	if u, ok := that.(*UnionType); ok {
		p := unionTypesUnderEquality(u, t)
		return TypePair{p.TypeB, p.TypeA}
	}
	// Union on the left: map across alternates.
	if u, ok := t.(*UnionType); ok {
		return unionTypesUnderEquality(u, that)
	}
	// Non-union: see if equality can hold.
	switch testForEquality(t, that) {
	case eqFalse:
		return TypePair{nil, nil}
	default:
		return TypePair{t, that}
	}
}

// GetTypesUnderInequality returns the pair of types that survive a != check.
// Port of JSType.getTypesUnderInequality(JSType).
func GetTypesUnderInequality(t, that JSType) TypePair {
	if that == nil {
		return TypePair{t, that}
	}
	// Union on the right: flip.
	if u, ok := that.(*UnionType); ok {
		p := unionTypesUnderInequality(u, t)
		return TypePair{p.TypeB, p.TypeA}
	}
	// Union on the left: map across alternates.
	if u, ok := t.(*UnionType); ok {
		return unionTypesUnderInequality(u, that)
	}
	// Non-union: null/void pairs resolve to No.
	if (isNullType(t) || isVoidType(t)) && (isNullType(that) || isVoidType(that)) {
		return TypePair{None, None}
	}
	return TypePair{t, that}
}

// eqResult is a simplified tri-state for equality.
type eqResult int

const (
	eqUnknown eqResult = iota
	eqTrue
	eqFalse
)

// testForEquality returns whether t == that is definitely false, definitely
// true, or unknown. Mirrors JSType.testForEquality's non-union base case.
func testForEquality(t, that JSType) eqResult {
	if t == nil || that == nil {
		return eqUnknown
	}
	if isAll(t) || isAll(that) || isUnknown(t) || isUnknown(that) {
		return eqUnknown
	}
	// null == undefined is true; null == anything_else is false.
	if isNullType(t) && isVoidType(that) {
		return eqTrue
	}
	if isVoidType(t) && isNullType(that) {
		return eqTrue
	}
	// Same singleton primitives.
	if t == that {
		return eqTrue
	}
	return eqUnknown
}

// unionTypesUnderEquality implements GetTypesUnderEquality for union on left.
func unionTypesUnderEquality(u *UnionType, that JSType) TypePair {
	b := NewUnionTypeBuilder()
	thatB := NewUnionTypeBuilder()
	for _, m := range u.Members() {
		p := GetTypesUnderEquality(m, that)
		if p.TypeA != nil {
			b.AddAlternate(p.TypeA)
		}
		if p.TypeB != nil {
			thatB.AddAlternate(p.TypeB)
		}
	}
	return TypePair{b.Build(), thatB.Build()}
}

// unionTypesUnderInequality implements GetTypesUnderInequality for union on left.
func unionTypesUnderInequality(u *UnionType, that JSType) TypePair {
	b := NewUnionTypeBuilder()
	thatB := NewUnionTypeBuilder()
	for _, m := range u.Members() {
		p := GetTypesUnderInequality(m, that)
		if p.TypeA != nil {
			b.AddAlternate(p.TypeA)
		}
		if p.TypeB != nil {
			thatB.AddAlternate(p.TypeB)
		}
	}
	return TypePair{b.Build(), thatB.Build()}
}

// RestrictByTypeof narrows t to the subset compatible with a typeof result.
// Port of the typeof-narrowing logic in ClosureReverseAbstractInterpreter.
// Returns None when the typeof value is inconsistent with t.
func RestrictByTypeof(t JSType, typeofValue string) JSType {
	if t == nil || isUnknown(t) || isAll(t) {
		return t
	}
	if u, ok := t.(*UnionType); ok {
		mapped := MapTypes(func(m JSType) JSType {
			return RestrictByTypeof(m, typeofValue)
		}, u.Members())
		return NewUnion(mapped...)
	}
	if typeofMatchesType(typeofValue, t) {
		return t
	}
	return None
}

// --- Additional Java-parity helpers ----------------------------------------

// IsNullable reports whether t includes null (is nullable). For union types,
// true if any member is null. Mirrors JSType.isNullable() in Java.
func IsNullable(t JSType) bool {
	if t == nil {
		return false
	}
	if isNullType(t) {
		return true
	}
	if u, ok := t.(*UnionType); ok {
		return u.IsNullable()
	}
	return false
}

// IsVoidable reports whether t includes undefined/void. For union types, true
// if any member is void. Mirrors JSType.isVoidable() in Java.
func IsVoidable(t JSType) bool {
	if t == nil {
		return false
	}
	if isVoidType(t) {
		return true
	}
	if u, ok := t.(*UnionType); ok {
		return u.IsVoidable()
	}
	return false
}

// IsObject reports whether t is an object type (ObjectType, FunctionType, or
// TemplatizedType). Mirrors JSType.isObject() in Java.
func IsObject(t JSType) bool {
	if t == nil {
		return false
	}
	switch t.(type) {
	case *ObjectType, *FunctionType, *TemplatizedType:
		return true
	case *UnionType:
		return t.(*UnionType).IsObject()
	}
	return false
}

// ToMaybeObjectType mirrors JSType.toObjectType() / toMaybeObjectType() — null-safe
// downcast to *ObjectType. Returns nil when t is nil or not an ObjectType.
func ToMaybeObjectType(t JSType) *ObjectType {
	if t == nil {
		return nil
	}
	o, _ := t.(*ObjectType)
	return o
}

// ToMaybeUnionType mirrors JSType.toMaybeUnionType() — null-safe downcast to
// *UnionType. Returns nil when t is not a UnionType.
func ToMaybeUnionType(t JSType) *UnionType {
	if t == nil {
		return nil
	}
	u, _ := t.(*UnionType)
	return u
}

// IsUnionType reports whether t is a UnionType.
// Mirrors JSType.isUnionType() in Java.
func IsUnionTypeValue(t JSType) bool {
	return ToMaybeUnionType(t) != nil
}

// GetLeastSupertype returns the least supertype (join / supremum) of x and y.
// Mirrors JSType.getLeastSupertype(JSType) in Java. For non-union types with
// no common supertype, returns a union of the two.
func GetLeastSupertype(x, y JSType) JSType {
	if x == nil {
		return y
	}
	if y == nil {
		return x
	}
	if x == y {
		return x
	}
	if isUnknown(x) || isAll(x) {
		return x
	}
	if isUnknown(y) || isAll(y) {
		return y
	}
	if IsSubtype(x, y) {
		return y
	}
	if IsSubtype(y, x) {
		return x
	}
	// No common supertype — create a union.
	return NewUnion(x, y)
}

// IsConstructorType reports whether t is a @constructor FunctionType.
// Mirrors JSType.isConstructor() in Java.
func IsConstructorType(t JSType) bool {
	if f, ok := t.(*FunctionType); ok {
		return f.kind == FunctionKindConstructor
	}
	return false
}

// IsInterfaceType reports whether t is a @interface FunctionType.
// Mirrors JSType.isInterface() in Java.
func IsInterfaceType(t JSType) bool {
	if f, ok := t.(*FunctionType); ok {
		return f.kind == FunctionKindInterface
	}
	return false
}

// IsTemplatizedType reports whether t is a *TemplatizedType.
// Mirrors JSType.isTemplatizedType() in Java.
func IsTemplatizedType(t JSType) bool {
	_, ok := t.(*TemplatizedType)
	return ok
}

// IsNamedType reports whether t is a *NamedType.
// Mirrors JSType.isNamedType() in Java.
func IsNamedType(t JSType) bool {
	_, ok := t.(*NamedType)
	return ok
}

// typeofMatchesType reports whether a typeof string is consistent with t.
func typeofMatchesType(typeofValue string, t JSType) bool {
	switch typeofValue {
	case "number":
		_, ok := t.(numberType)
		return ok
	case "string":
		_, ok := t.(stringType)
		return ok
	case "boolean":
		_, ok := t.(booleanType)
		return ok
	case "bigint":
		_, ok := t.(bigIntType)
		return ok
	case "symbol":
		return IsSymbol(t)
	case "undefined":
		return isVoidType(t)
	case "function":
		_, ok := t.(*FunctionType)
		return ok
	case "object":
		if isNullType(t) {
			return true // typeof null === "object"
		}
		_, ok := t.(*ObjectType)
		return ok
	}
	return false
}
