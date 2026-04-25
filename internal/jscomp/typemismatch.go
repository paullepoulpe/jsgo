// Port of com/google/javascript/jscomp/TypeMismatch.java @ closure-compiler v20260415
//
// TypeMismatch is a data record that signals the found type and the
// required type were used interchangeably at a specific source location.
// Type-based optimisations consult the accumulated set so they don't
// silently break code that already has type warnings.
//
// The Java side is a Java 16 record with an inner Accumulator class that
// builds the set during type validation; both are ported here.
// Accumulator's full logic is wired when TypeValidator lands; for now it
// is present as a faithful structural port so TypeValidator can import it
// without changes.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// TypeMismatch signals that found and required were used interchangeably
// at location.
//
//   - Found    — the RHS type; the type of the value being assigned.
//   - Required — the LHS type; the declared / expected type.
//   - Location — the AST node where the mismatch was detected.
//
// Mirrors the Java record TypeMismatch(JSType found, JSType required,
// Node location).
type TypeMismatch struct {
	Found    jstype.JSType
	Required jstype.JSType
	Location *rhino.Node
}

// GetFound ports the @InlineMe accessor `this.found()`.
func (m *TypeMismatch) GetFound() jstype.JSType { return m.Found }

// GetRequired ports the @InlineMe accessor `this.required()`.
func (m *TypeMismatch) GetRequired() jstype.JSType { return m.Required }

// GetLocation ports the @InlineMe accessor `this.location()`.
func (m *TypeMismatch) GetLocation() *rhino.Node { return m.Location }

// testLocation is the sentinel Node used by CreateForTesting, mirroring
// Java's `private static final Node TEST_LOCATION = IR.empty()`.
var testLocation = rhino.New(rhino.ERROR)

// create is the internal factory, mirroring the private static
// TypeMismatch.create(JSType, JSType, Node).
func newTypeMismatch(found, required jstype.JSType, location *rhino.Node) *TypeMismatch {
	if found == nil {
		panic("TypeMismatch: found must not be nil")
	}
	if required == nil {
		panic("TypeMismatch: required must not be nil")
	}
	if location == nil {
		panic("TypeMismatch: location must not be nil")
	}
	return &TypeMismatch{Found: found, Required: required, Location: location}
}

// CreateForTesting is the @VisibleForTesting factory used in tests when
// no real source location is available. Mirrors Java's
// TypeMismatch.createForTesting(JSType, JSType).
func CreateForTesting(found, required jstype.JSType) *TypeMismatch {
	return newTypeMismatch(found, required, testLocation)
}

// Accumulator collects a set of related TypeMismatches. It is the port of
// the inner class TypeMismatch.Accumulator in Java. TypeValidator
// (not yet ported) calls RegisterMismatch for every detected mismatch so
// that type-based optimisations can suppress unsafe transformations.
type Accumulator struct {
	mismatches []*TypeMismatch
}

// RegisterMismatch records a mismatch at location between found and
// required, expanding into per-parameter / return mismatches when both
// types are function types. Mirrors
// Accumulator.registerMismatch(Node, JSType, JSType).
//
// Per Java semantics the method skips:
//   - mismatches that disappear once null/undefined and template wrappers
//     are stripped (restrictByNotNullOrUndefined + raw templatized type),
//   - pairs where one is already a subtype of the other after stripping,
//   - pairs where a required-union member is a subtype of found (downcast
//     guard), and
//   - pairs where either side is a raw TemplateType (open type variable).
func (a *Accumulator) RegisterMismatch(location *rhino.Node, found, required jstype.JSType) {
	found = removeNullUndefinedAndTemplates(found)
	required = removeNullUndefinedAndTemplates(required)

	if jstype.IsSubtype(found, required) || jstype.IsSubtype(required, found) {
		return
	}

	// If required is a union and any alternate is a subtype of found,
	// this is a safe downcast — skip it.
	if u, ok := required.(*jstype.UnionType); ok {
		for _, alt := range u.Members() {
			if jstype.IsSubtype(alt, found) {
				return
			}
		}
	}

	if bothAreNotTemplateTypes(found, required) {
		a.mismatches = append(a.mismatches, newTypeMismatch(found, required, location))
	}

	// Recurse into paired function parameters and return types.
	fnFound, foundIsFn := found.(*jstype.FunctionType)
	fnRequired, requiredIsFn := required.(*jstype.FunctionType)
	if foundIsFn && requiredIsFn {
		minLen := len(fnFound.Params)
		if len(fnRequired.Params) < minLen {
			minLen = len(fnRequired.Params)
		}
		for i := 0; i < minLen; i++ {
			a.registerIfMismatch(location, fnFound.Params[i], fnRequired.Params[i])
		}
		a.registerIfMismatch(location, fnFound.Return, fnRequired.Return)
	}
}

// GetMismatches returns the accumulated slice. Mirrors
// Accumulator.getMismatches() → ImmutableCollection<TypeMismatch>.
func (a *Accumulator) GetMismatches() []*TypeMismatch {
	out := make([]*TypeMismatch, len(a.mismatches))
	copy(out, a.mismatches)
	return out
}

// registerIfMismatch calls RegisterMismatch only when neither type is nil
// and found is not a subtype of required. Mirrors the private helper
// Accumulator.registerIfMismatch(Node, JSType, JSType).
func (a *Accumulator) registerIfMismatch(location *rhino.Node, found, required jstype.JSType) {
	if found != nil && required != nil && !jstype.IsSubtype(found, required) {
		a.RegisterMismatch(location, found, required)
	}
}

// bothAreNotTemplateTypes returns true iff neither type is a raw
// TemplateType. Mirrors the private static helper in Java.
//
// A type variable is treated as unknown so any concrete type is
// accepted in its position without warning; recording a mismatch
// against an open type variable would create false positives.
func bothAreNotTemplateTypes(found, required jstype.JSType) bool {
	_, foundIsTemplate := found.(*jstype.TemplateType)
	_, requiredIsTemplate := required.(*jstype.TemplateType)
	return !foundIsTemplate && !requiredIsTemplate
}

// removeNullUndefinedAndTemplates strips null/undefined alternatives from
// a union and, for templatised types, returns the raw base type. Mirrors
// Java's Accumulator.removeNullUndefinedAndTemplates(JSType).
//
// Java calls t.restrictByNotNullOrUndefined() which removes null and
// undefined members from a union. We replicate that logic directly: for
// union types, rebuild the union without null/undefined members; for
// non-union types, return as-is (null and undefined are leaves, not
// wrapping types). Then, if the result is a TemplatizedType, return its
// raw ObjectType base — all specialisations of a generic class share the
// same JS code so the template arguments are irrelevant for mismatch
// tracking.
func removeNullUndefinedAndTemplates(t jstype.JSType) jstype.JSType {
	// Strip null / undefined from unions.
	result := restrictByNotNullOrUndefined(t)

	// If the result is a templatised type, return the raw base.
	if tt, ok := result.(*jstype.TemplatizedType); ok {
		return tt.Base
	}
	return result
}

// restrictByNotNullOrUndefined removes null and undefined from a union
// type, returning the collapsed result. For non-union types it is a
// no-op (null and undefined are singletons, not wrappers). Mirrors
// JSType.restrictByNotNullOrUndefined() in Java.
func restrictByNotNullOrUndefined(t jstype.JSType) jstype.JSType {
	u, ok := t.(*jstype.UnionType)
	if !ok {
		return t
	}
	var keep []jstype.JSType
	for _, m := range u.Members() {
		if m == jstype.Null || m == jstype.Void {
			continue
		}
		keep = append(keep, m)
	}
	if len(keep) == 0 {
		return t // degenerate: all members were null/undefined
	}
	return jstype.NewUnion(keep...)
}
