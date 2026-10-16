// Port of com/google/javascript/rhino/jstype/NoResolvedType.java @ closure-compiler v20260415

package rhino

import (
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// NoResolvedType is an unresolved type that was forward declared. So we
// know it exists, but that it wasn't pulled into this binary.
//
// In most cases, it behaves like a sibling of the top object type.
//
// Why have this class at all, and not just use the top object type or `*`?
//
// Clutz relies on this class to emit .d.ts typings for JS files without
// having to see all their dependencies: Clutz expects unresolved types to be
// available when it runs, so that it can map them back to the corresponding
// goog.required import.
//
// Clutz also sets a special JSCompiler flag so that all unresolved types are
// treated as forward declared. (In most non-Clutz cases, this class is only
// rarely used, and only if there's actually an explicit
// `goog.forwardDeclare` call.)
//
// To address this, we would need to significantly refactor Clutz.
type NoResolvedType struct {
	ObjectTypeBase

	// The name originally used to reference this type.
	referenceName string

	// Any template arguments to this type, or null if none. This field is not
	// used for JSCompiler's type checking; it is only needed by Clutz.
	templateTypes []JSType
}

// Java: NoResolvedType(JSTypeRegistry, String, ImmutableList<JSType>)
func NewNoResolvedType(registry *JSTypeRegistry, referenceName string, templateTypes []JSType) *NoResolvedType {
	t := &NoResolvedType{}
	t.InitObjectTypeBase(t, registry)
	t.referenceName = referenceName
	t.templateTypes = templateTypes
	t.EagerlyResolveToSelf()
	return t
}

// Java: NoResolvedType.getTypeClass()
func (t *NoResolvedType) GetTypeClass() JSTypeClass {
	return JSTypeClass_NO_RESOLVED
}

// Java: NoResolvedType.getReferenceName()
func (t *NoResolvedType) GetReferenceName() *string {
	return &t.referenceName
}

// Java: NoResolvedType.getTemplateTypes()
func (t *NoResolvedType) GetTemplateTypes() []JSType {
	return t.templateTypes
}

// Java: NoResolvedType.isNoResolvedType()
func (t *NoResolvedType) IsNoResolvedType() bool {
	return true
}

// Java: NoResolvedType.getPossibleToBooleanOutcomes()
func (t *NoResolvedType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_BOTH
}

// Java: NoResolvedType.testForEquality(JSType)
func (t *NoResolvedType) TestForEquality(that JSType) base.Tri {
	return base.Tri_UNKNOWN
}

// Java: NoResolvedType.resolveInternal(ErrorReporter)
func (t *NoResolvedType) ResolveInternal(reporter ErrorReporter) JSType {
	panic(jutil.NewAssertionErrorNoMessage())
}

// Java: NoResolvedType.getConstructor()
func (t *NoResolvedType) GetConstructor() FunctionType {
	return nil
}

// Java: NoResolvedType.getImplicitPrototype()
func (t *NoResolvedType) GetImplicitPrototype() ObjectType {
	return nil
}

// Java: NoResolvedType.recursionUnsafeHashCode()
func (t *NoResolvedType) RecursionUnsafeHashCode() int32 {
	// NoResolvedType instances are unique within a JSTypeRegistry.
	return jutil.IdentityHashCode(t)
}

// Java: NoResolvedType.defineProperty(Property.Key, JSType, boolean, Node)
func (t *NoResolvedType) DefinePropertyByKey(propertyName Property_Key, type_ JSType, inferred bool, propertyNode *Node) bool {
	// nothing, all properties are defined
	return true
}

// Java: NoResolvedType.appendTo(TypeStringBuilder)
func (t *NoResolvedType) AppendTo(sb *TypeStringBuilder) {
	if sb.IsForAnnotations() {
		sb.Append("?")
	} else {
		sb.Append("NoResolvedType<" + t.referenceName + ">")
	}
}
