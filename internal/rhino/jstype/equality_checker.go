// Port of com/google/javascript/rhino/jstype/EqualityChecker.java
// @ closure-compiler v20260415
//
// Represents the computation of a single equality relationship.
// Each instance can only be used once. It acts as a scope for the computation
// and may accumulate state.
//
// NOTE: Types referenced by the Java original that are not yet ported to Go
// (RecordType, NamedType, TemplatizedType full API, FunctionType.getKind,
// FunctionType.getClosurePrimitive, FunctionType.getGoogModuleId) are stubbed
// with TODO markers. The public API (EqualityChecker, EqMethod) is complete
// and callable.

package jstype

// EqMethod mirrors EqualityChecker.EqMethod — the two modes of equality
// comparison.
type EqMethod int

const (
	// EqMethodIdentity indicates the two types should behave exactly the same
	// under all type operations. {string} != {?} and {Unresolved} != {?}.
	EqMethodIdentity EqMethod = iota
	// EqMethodDataFlow indicates the two types are compared for data-flow
	// termination. Unknown types are treated as equivalent.
	// {string} != {?} and {Unresolved} ~= {?}.
	EqMethodDataFlow
)

// eqCacheKey is the cache key for EqualityChecker's cycle-detection map.
// Mirrors EqualityChecker.CacheKey in Java.
// Note: Java uses identity (System.identityHashCode) + XOR for the key;
// we use the same pointer-identity semantics here.
type eqCacheKey struct {
	left  JSType
	right JSType
}

// potentiallyCyclicRecursionDepthEq mirrors
// EqualityChecker.POTENTIALLY_CYCLIC_RECURSION_DEPTH.
const potentiallyCyclicRecursionDepthEq = 20

// EqualityChecker represents the computation of a single equality relationship.
// Each instance can only be used once. Mirrors EqualityChecker in Java.
type EqualityChecker struct {
	eqMethod       EqMethod
	eqCache        map[eqCacheKey]matchStatus
	recursionDepth int
	hasRun         bool
}

// NewEqualityChecker allocates a fresh EqualityChecker.
// Mirrors EqualityChecker() default constructor in Java.
func NewEqualityChecker() *EqualityChecker {
	return &EqualityChecker{}
}

// SetEqMethod sets the equality method. Must be called before Check.
// Mirrors EqualityChecker.setEqMethod in Java.
func (ec *EqualityChecker) SetEqMethod(x EqMethod) *EqualityChecker {
	ec.checkHasNotRun()
	ec.eqMethod = x
	return ec
}

func (ec *EqualityChecker) checkHasNotRun() {
	if ec.hasRun {
		panic("EqualityChecker has already run")
	}
}

// Check runs the equality check for the given pair. Can only be called once.
// Mirrors EqualityChecker.check in Java.
func (ec *EqualityChecker) Check(left, right JSType) bool {
	ec.checkHasNotRun()
	ec.hasRun = true
	return ec.areEqualCaching(left, right)
}

// CheckParameters returns whether the parameters (ignoring any other aspects)
// of the two ArrowTypes are equal. Can only be called once.
// Mirrors EqualityChecker.checkParameters in Java.
func (ec *EqualityChecker) CheckParameters(left, right *ArrowType) bool {
	ec.checkHasNotRun()
	ec.hasRun = true
	return left == right || ec.areArrowParameterEqual(left, right)
}

// areEqualCaching is the main recursive entrypoint with cycle detection.
// Mirrors EqualityChecker.areEqualCaching in Java.
func (ec *EqualityChecker) areEqualCaching(left, right JSType) bool {
	// Wait to instantiate/use the cache until we have some hint of recursion.
	if ec.recursionDepth > potentiallyCyclicRecursionDepthEq {
		if ec.eqCache == nil {
			ec.eqCache = make(map[eqCacheKey]matchStatus)
		}
	}

	// Once the cache exists, use it consistently.
	if ec.eqCache == nil {
		ec.recursionDepth++
		result := ec.areEqualInternal(left, right)
		ec.recursionDepth--
		return result
	}

	key := eqCacheKey{left: left, right: right}
	cached, exists := ec.eqCache[key]

	if !exists {
		ec.eqCache[key] = matchStatusProcessing
		result := ec.areEqualInternal(left, right)
		ec.eqCache[key] = matchStatusValueOf(result)
		return result
	}

	if cached == matchStatusProcessing {
		ec.eqCache[key] = matchStatusMatch
		return true
	}

	return cached.subtypeValue()
}

// areEqualInternal is the core equality dispatch.
// Mirrors EqualityChecker.areEqualInternal in Java.
func (ec *EqualityChecker) areEqualInternal(left, right JSType) bool {
	if left == right {
		return true
	}
	if left == nil || right == nil {
		return false
	}

	// NoResolvedType equality by reference name.
	// TODO(EqualityChecker): implement when NoResolvedType is ported.
	// if left.isNoResolvedType() && right.isNoResolvedType() {
	//     return left.toObjectType().getReferenceName() == right.toObjectType().getReferenceName()
	// }

	leftUnknown := isUnknown(left)
	rightUnknown := isUnknown(right)
	if leftUnknown || rightUnknown {
		if ec.eqMethod == EqMethodDataFlow {
			// Both must be unknown for data-flow equality.
			return leftUnknown && rightUnknown
		} else if leftUnknown && rightUnknown {
			// If both unknown but one is nominal and other is not, fail.
			leftNominal := isNominalType(left)
			rightNominal := isNominalType(right)
			if leftNominal != rightNominal {
				return false
			}
		}
	}

	// Union equality.
	leftUnion, leftIsUnion := left.(*UnionType)
	rightUnion, rightIsUnion := right.(*UnionType)
	if leftIsUnion && rightIsUnion {
		return ec.areUnionEqual(leftUnion, rightUnion)
	} else if leftIsUnion {
		if len(leftUnion.members) == 1 {
			return ec.areEqualInternal(leftUnion.members[0], right)
		}
	} else if rightIsUnion {
		if len(rightUnion.members) == 1 {
			return ec.areEqualInternal(left, rightUnion.members[0])
		}
	}

	// FunctionType equality.
	leftFn, leftIsFn := left.(*FunctionType)
	rightFn, rightIsFn := right.(*FunctionType)
	if leftIsFn && rightIsFn {
		return ec.areFunctionEqual(leftFn, rightFn)
	}

	// ArrowType equality.
	leftArrow, leftIsArrow := left.(*ArrowType)
	rightArrow, rightIsArrow := right.(*ArrowType)
	if leftIsArrow && rightIsArrow {
		return ec.areArrowEqual(leftArrow, rightArrow)
	}

	// TemplateTypeMap equality.
	// TODO(EqualityChecker): implement when TemplateTypeMap full API is ported.
	// if !ec.areTypeMapEqual(left.getTemplateTypeMap(), right.getTemplateTypeMap()) {
	//     return false
	// }

	// Record type equality.
	// TODO(EqualityChecker): implement when RecordType is ported.
	// if left.isRecordType() && right.isRecordType() {
	//     return ec.areRecordEqual(left.toMaybeRecordType(), right.toMaybeRecordType())
	// }

	// Nominal type equality.
	if isNominalType(left) && isNominalType(right) {
		leftUnwrapped := unwrapNominalTypeProxies(left)
		rightUnwrapped := unwrapNominalTypeProxies(right)
		if isResolved(left) && isResolved(right) {
			return leftUnwrapped == rightUnwrapped
		}
		// TODO(EqualityChecker): pre-resolution cross-scope equality via reference names.
		// For now use pointer identity.
		return leftUnwrapped == rightUnwrapped
	}

	// Unwrap proxies (but not TemplateTypes, which have identity semantics).
	if lp, ok := left.(*ProxyObjectType); ok {
		return ec.areEqualCaching(lp.ReferencedType, right)
	}
	if rp, ok := right.(*ProxyObjectType); ok {
		return ec.areEqualCaching(left, rp.ReferencedType)
	}

	// Base case: only one instance of each sub-type will ever be created in a
	// given registry, so there is no need to verify members (pointer identity).
	return false
}

// areUnionEqual checks equality of two union types.
// Two union types are equal if, after flattening, they have the same number
// of alternates and all alternates are equal (order-independent).
// Mirrors EqualityChecker.areUnionEqual in Java.
func (ec *EqualityChecker) areUnionEqual(left, right *UnionType) bool {
	leftAlts := left.members
	rightAlts := right.members
	if ec.eqMethod == EqMethodIdentity && len(leftAlts) != len(rightAlts) {
		return false
	}

outer:
	for i := 0; i < len(rightAlts); i++ {
		rightAlt := rightAlts[i]
		for k := 0; k < len(leftAlts); k++ {
			if ec.areEqualCaching(leftAlts[k], rightAlt) {
				continue outer
			}
		}
		return false
	}

	return true
}

// areFunctionEqual checks equality of two function types.
// Two function types are equal if their signatures match. Since they don't
// have signatures, two interfaces are equal if their names match.
// Mirrors EqualityChecker.areFunctionEqual in Java.
func (ec *EqualityChecker) areFunctionEqual(left, right *FunctionType) bool {
	if left == right {
		return true
	}

	// TODO(EqualityChecker): left.getKind() / right.getKind() check.
	// Java distinguishes CONSTRUCTOR, INTERFACE, ORDINARY.
	// For now: constructors and interfaces use identity semantics (already
	// checked above); ordinary functions compare this-type + arrow + closure primitive.

	if left.IsConstructor != right.IsConstructor {
		return false
	}
	if left.IsConstructor {
		// Constructor/interface: identity semantics (already false from above check).
		return false
	}

	// Ordinary function: compare parameters and return type.
	return ec.areEqualCaching(left.Return, right.Return) &&
		ec.areFunctionParamsEqual(left, right)
	// TODO(EqualityChecker): getClosurePrimitive equality check.
}

// areFunctionParamsEqual compares parameter lists of two FunctionTypes.
func (ec *EqualityChecker) areFunctionParamsEqual(left, right *FunctionType) bool {
	if len(left.Params) != len(right.Params) {
		return false
	}
	for i := range left.Params {
		if !ec.areEqualCaching(left.Params[i], right.Params[i]) {
			return false
		}
	}
	return true
}

// areArrowEqual checks equality of two ArrowTypes.
// Mirrors EqualityChecker.areArrowEqual in Java.
func (ec *EqualityChecker) areArrowEqual(left, right *ArrowType) bool {
	return ec.areEqualCaching(left.ReturnType, right.ReturnType) &&
		ec.areArrowParameterEqual(left, right)
}

// areArrowParameterEqual checks equality of ArrowType parameter lists.
// Mirrors EqualityChecker.areArrowParameterEqual in Java.
func (ec *EqualityChecker) areArrowParameterEqual(left, right *ArrowType) bool {
	if len(left.Params) != len(right.Params) {
		return false
	}
	for i := range left.Params {
		lp := left.Params[i]
		rp := right.Params[i]

		lType := lp.JSType
		rType := rp.JSType
		if lType != nil {
			if rType != nil && !ec.areEqualCaching(lType, rType) {
				return false
			}
			// rType nil: left has type, right doesn't — not equal by Java's logic.
			// (Java: if leftParamType != null && rightParamType == null → return false, since
			// the else branch handles rightParamType != null only.)
		} else {
			if rType != nil {
				return false
			}
		}

		if lp.Optional != rp.Optional {
			return false
		}
		if lp.Variadic != rp.Variadic {
			return false
		}
	}
	return true
}

// areRecordEqual checks equality of two record types by their own properties.
// Mirrors EqualityChecker.areRecordEqual in Java.
// TODO(EqualityChecker): implement when RecordType is ported.
// func (ec *EqualityChecker) areRecordEqual(left, right *RecordType) bool { ... }

// areTypeMapEqual checks equality of two TemplateTypeMaps.
// Mirrors EqualityChecker.areTypeMapEqual in Java.
// TODO(EqualityChecker): implement when TemplateTypeMap full key/value API is ported.
// func (ec *EqualityChecker) areTypeMapEqual(left, right *TemplateTypeMap) bool { ... }

// ---- helper predicates ----

// isNominalType returns true for types that have a reference name and use
// nominal identity semantics. Mirrors JSType.isNominalType.
// TODO(EqualityChecker): expand when NamedType and full ObjectType hierarchy is ported.
func isNominalType(t JSType) bool {
	switch v := t.(type) {
	case *ObjectType:
		return v.Name != ""
	case *FunctionType:
		return v.IsConstructor
	}
	return false
}

// isResolved reports whether a type has been resolved (not a forward reference).
// Mirrors JSType.isResolved. For our port, all types are considered resolved
// unless they are ProxyObjectType wrappers.
// TODO(EqualityChecker): proper resolution tracking when NamedType is ported.
func isResolved(t JSType) bool {
	_, isProxy := t.(*ProxyObjectType)
	return !isProxy
}

// unwrapNominalTypeProxies unwraps NamedType and TemplatizedType proxies to
// reach the underlying nominal type. Mirrors EqualityChecker.unwrapNominalTypeProxies.
// TODO(EqualityChecker): handle NamedType when it is ported.
func unwrapNominalTypeProxies(t JSType) JSType {
	switch v := t.(type) {
	case *TemplatizedType:
		if v.Base != nil {
			return unwrapNominalTypeProxies(v.Base)
		}
	case *ProxyObjectType:
		return unwrapNominalTypeProxies(v.ReferencedType)
	}
	return t
}

// getGoogModuleId returns the goog.module ID for a nominal type, if any.
// Mirrors EqualityChecker.getGoogModuleId in Java.
// TODO(EqualityChecker): implement when FunctionType.getGoogModuleId is ported.
func getGoogModuleId(_ JSType) string { return "" }
