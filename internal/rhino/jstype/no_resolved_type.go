// Port of com/google/javascript/rhino/jstype/NoResolvedType.java
// @ closure-compiler v20260415
//
// NoResolvedType represents a forward-declared type that was never
// pulled into the binary. In most cases it behaves like the top object
// type. Clutz relies on it to emit .d.ts typings without seeing all
// dependencies of a JS file — it maps unresolved types back to the
// corresponding goog.require import.
//
// In the Java hierarchy NoResolvedType extends ObjectType and carries a
// referenceName (the original forward-declared identifier) plus optional
// template arguments. The type equality check returns UNKNOWN (Tri) and
// getPossibleToBooleanOutcomes returns BOTH — matching the Java source.

package jstype

// NoResolvedType is an unresolved forward-declared type. Mirrors
// com.google.javascript.rhino.jstype.NoResolvedType.
type NoResolvedType struct {
	// ReferenceName is the original identifier used to forward-declare
	// this type (e.g. "goog.module.Foo"). Mirrors
	// NoResolvedType.getReferenceName().
	ReferenceName string

	// TemplateTypes carries any template arguments supplied at the use
	// site (e.g. `Foo<Bar>`). nil when none were provided. This field is
	// not used by JSCompiler's type checking; it is only needed by Clutz.
	// Mirrors NoResolvedType.getTemplateTypes().
	TemplateTypes []JSType
}

// NewNoResolvedType allocates a fresh NoResolvedType. templateTypes may
// be nil. Mirrors the package-private constructor in Java.
func NewNoResolvedType(referenceName string, templateTypes []JSType) *NoResolvedType {
	return &NoResolvedType{
		ReferenceName: referenceName,
		TemplateTypes: templateTypes,
	}
}

// TypeName returns the internal string representation. Java's appendTo
// emits "?" for annotation contexts and "NoResolvedType<name>" otherwise;
// we use the non-annotation form for diagnostics.
func (n *NoResolvedType) TypeName() string {
	return "NoResolvedType<" + n.ReferenceName + ">"
}

// IsNoResolvedType mirrors NoResolvedType.isNoResolvedType() → true.
func (n *NoResolvedType) IsNoResolvedType() bool { return true }

// IsSubtypeOf implements the permissive rule for unresolved types: a
// forward-declared type is treated as compatible with any target.
// This mirrors how Java's NoResolvedType behaves — it is a sibling of
// the top object type and therefore assignment-compatible everywhere an
// object is accepted. Being conservative (accepting all targets) avoids
// false positives when dependencies are absent.
func (n *NoResolvedType) IsSubtypeOf(other JSType) bool { return true }
