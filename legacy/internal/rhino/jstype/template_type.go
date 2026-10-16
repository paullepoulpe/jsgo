// Port of com/google/javascript/rhino/jstype/TemplateType.java @ closure-compiler v20260415
//
// TemplateType is a placeholder type used as keys in TemplateTypeMaps.
// In Java it extends ProxyObjectType, wrapping its bound (or Unknown when
// unbound) as the referenced type. The Go port keeps a flat struct to
// avoid restructuring all callers; ProxyObjectType delegation is expressed
// as explicit method bodies mirroring Java's @Override chain.
//
// The full Java-faithful TemplateTypeMap is in template_type_map.go;
// this file retains CallSiteMap, the narrow call-site substitution
// struct used by TypeInference.

package jstype

import (
	"unsafe"

	"github.com/paullepoulpe/jsgo/legacy/internal/rhino"
)

// TemplateType is a placeholder type, used as keys in TemplateTypeMaps.
// Port of TemplateType extends ProxyObjectType.
type TemplateType struct {
	// name is the template parameter name (e.g. "T").
	// Port of: private final String name.
	// Exported as Name for direct field access by callers in this package.
	Name string

	// bound is the upper-bound constraint from `@template T extends Base`.
	// nil means unbounded (equivalent to Unknown in Java).
	// Port of: private JSType bound.
	bound JSType

	// typeTransformation holds the AST node for TTL (type-transformation
	// language) @template annotations.
	// Port of: private final Node typeTransformation.
	typeTransformation *rhino.Node
}

// GetTypeClass returns JSTypeClassTemplate.
// Port of rhino/jstype/TemplateType.java @ closure-compiler v20260415
// Port of: JSTypeClass getTypeClass() { return TYPE_CLASS; }
func (t *TemplateType) GetTypeClass() JSTypeClass {
	return JSTypeClassTemplate
}

// GetReferenceName returns the template parameter name.
// Port of rhino/jstype/TemplateType.java @ closure-compiler v20260415
// Port of: public String getReferenceName() { return name; }
func (t *TemplateType) GetReferenceName() string {
	return t.Name
}

// AppendTo renders the template type into sb.
// When the bound is unknown (nil or Unknown), only the name is written.
// When a non-unknown bound is set, the name and " extends <bound>" are written.
// Port of rhino/jstype/TemplateType.java @ closure-compiler v20260415
// Port of: void appendTo(TypeStringBuilder sb) { ... }
func (t *TemplateType) AppendTo(sb *TypeStringBuilder) {
	bound := t.getBound()
	if isUnknown(bound) {
		sb.Append(t.Name)
	} else {
		sb.Append(t.Name).Append(" extends ").Append(bound.TypeName())
	}
}

// ToMaybeTemplateType returns this TemplateType.
// Port of rhino/jstype/TemplateType.java @ closure-compiler v20260415
// Port of: public TemplateType toMaybeTemplateType() { return this; }
func (t *TemplateType) ToMaybeTemplateType() *TemplateType {
	return t
}

// HasAnyTemplateTypesInternal reports whether this type contains any
// unresolved template type parameters. TemplateType itself always returns true.
// Port of rhino/jstype/TemplateType.java @ closure-compiler v20260415
// Port of: public boolean hasAnyTemplateTypesInternal() { return true; }
func (t *TemplateType) HasAnyTemplateTypesInternal() bool {
	return true
}

// Visit dispatches this TemplateType to the given Visitor.
// Port of rhino/jstype/TemplateType.java @ closure-compiler v20260415
// Port of: public <T> T visit(Visitor<T> visitor) { return visitor.caseTemplateType(this); }
func (t *TemplateType) Visit(visitor Visitor) any {
	return visitor.CaseTemplateType(t)
}

// VisitRelationship dispatches this TemplateType to the given RelationshipVisitor.
// Port of rhino/jstype/TemplateType.java @ closure-compiler v20260415
// Port of: <T> T visit(RelationshipVisitor<T> visitor, JSType that) { return visitor.caseTemplateType(this, that); }
func (t *TemplateType) VisitRelationship(visitor RelationshipVisitor, that JSType) any {
	return visitor.CaseTemplateType(t, that)
}

// IsTypeTransformation reports whether this template has a type-transformation expression.
// Port of rhino/jstype/TemplateType.java @ closure-compiler v20260415
// Port of: public boolean isTypeTransformation() { return typeTransformation != null; }
func (t *TemplateType) IsTypeTransformation() bool {
	return t.typeTransformation != nil
}

// GetTypeTransformation returns the type-transformation AST node, or nil.
// Port of rhino/jstype/TemplateType.java @ closure-compiler v20260415
// Port of: public Node getTypeTransformation() { return typeTransformation; }
func (t *TemplateType) GetTypeTransformation() *rhino.Node {
	return t.typeTransformation
}

// GetBound returns the upper-bound constraint for this template parameter.
// Returns Unknown when unbound.
// Port of rhino/jstype/TemplateType.java @ closure-compiler v20260415
// Port of: public JSType getBound() { return bound; }
func (t *TemplateType) GetBound() JSType {
	return t.getBound()
}

// SetBound updates the bound and (mirroring Java's ProxyObjectType) also
// updates the referenced type to match the new bound.
// Port of rhino/jstype/TemplateType.java @ closure-compiler v20260415
// Port of: public void setBound(JSType bound) { this.bound = bound; this.setReferencedType(bound); }
func (t *TemplateType) SetBound(b JSType) {
	t.bound = b
	// setReferencedType(bound) — in Java this updates the ProxyObjectType's
	// referencedType. In Go, the bound IS the referenced type for TemplateType.
}

// SetValidator applies validator to this TemplateType and returns the result.
// ProxyObjectType would delegate to the unknown type's setValidator, which
// we don't want — some validators special-case template types.
// Port of rhino/jstype/TemplateType.java @ closure-compiler v20260415
// Port of: public boolean setValidator(Predicate<JSType> validator) { return validator.apply(this); }
func (t *TemplateType) SetValidator(validator func(JSType) bool) bool {
	return validator(t)
}

// ContainsCycle reports whether there is a cycle in the reference chain of
// this type. The visitor is passed an unreachable target (nil) so it only
// returns when a cycle is found or a terminating type is reached.
// Port of rhino/jstype/TemplateType.java @ closure-compiler v20260415
// Port of: public boolean containsCycle() { ... }
func (t *TemplateType) ContainsCycle() bool {
	typeVisitor := NewContainsUpperBoundSuperTypeVisitor(nil)
	result := t.Visit(typeVisitor).(UpperBoundResult)
	return result == UpperBoundCycle
}

// RecursionUnsafeHashCode returns the identity hash code for this type.
// TemplateType uses identity semantics for equality (see EqualityChecker),
// so identity hash code is used to avoid spurious hash conflicts.
// Port of rhino/jstype/TemplateType.java @ closure-compiler v20260415
// Port of: int recursionUnsafeHashCode() { return System.identityHashCode(this); }
func (t *TemplateType) RecursionUnsafeHashCode() uintptr {
	return uintptr(unsafe.Pointer(t))
}

// --- Go constructors (no Java analog — Java uses package-private constructors with registry) ---

// NewTemplateType returns a fresh TemplateType carrying the named
// `@template` parameter. Not interned: the caller (TypedScopeCreator)
// is expected to create one per function and reuse it across that
// function's @param / @return slots.
// Go-only: language-adapter — Go constructor; Java uses package-private TemplateType(registry, name).
func NewTemplateType(name string) *TemplateType {
	return &TemplateType{Name: name}
}

// NewBoundedTemplateType returns a fresh TemplateType with an upper-bound
// constraint. Mirrors TemplateType(registry, name, bound) in Java.
// Go-only: language-adapter — Go constructor; Java uses package-private TemplateType(registry, name, bound).
func NewBoundedTemplateType(name string, bound JSType) *TemplateType {
	return &TemplateType{Name: name, bound: bound}
}

// NewTemplateTypeWithTransformation returns a fresh TemplateType with a
// type-transformation expression. Mirrors TemplateType(registry, name, typeTransformation) in Java.
// Go-only: language-adapter — Go constructor; Java uses package-private TemplateType(registry, name, typeTransformation).
func NewTemplateTypeWithTransformation(name string, typeTransformation *rhino.Node) *TemplateType {
	return &TemplateType{Name: name, typeTransformation: typeTransformation}
}

// --- Go-only interface adapters (satisfy the JSType interface) ---

// TypeName returns the string representation of the template type for
// diagnostics. In Java, appendTo(TypeStringBuilder) serves this role.
// Go-only: language-adapter — Go JSType interface method; Java uses appendTo(TypeStringBuilder) instead.
func (t *TemplateType) TypeName() string { return t.Name }

// IsSubtypeOf is permissive — a TemplateType is treated like Unknown
// so unresolved generics neither false-positive nor false-negative.
// Once unification refines the type at a call site, the refined
// JSType (not the TemplateType) participates in the subtype check.
// Go-only: language-adapter — Go JSType interface method; Java uses visit(RelationshipVisitor, JSType).
func (t *TemplateType) IsSubtypeOf(o JSType) bool { return true }

// --- CallSiteMap — Go-only narrow call-site substitution map ---
//
// CallSiteMap is the substitution `TemplateType → JSType` built
// at a CALL site. Narrow port of the call-site unification portion
// of rhino/jstype/TemplateTypeMap.java; used by TypeInference.
// The full Java-faithful TemplateTypeMap is in template_type_map.go.
//
// Go-only: language-adapter — CallSiteMap has no direct Java analog in TemplateType.java; it is a Go helper for call-site unification.
type CallSiteMap struct {
	subs map[*TemplateType]JSType
}

// NewCallSiteMap returns a new empty CallSiteMap.
// Go-only: language-adapter — Go constructor for CallSiteMap.
func NewCallSiteMap() *CallSiteMap {
	return &CallSiteMap{subs: map[*TemplateType]JSType{}}
}

// Bind records tmpl → t. If a prior binding exists and the new type
// differs, the two are unioned — matches Java's "pick the join" rule
// when the same template parameter is filled from multiple argument
// positions.
// Go-only: language-adapter — call-site unification helper with no Java analog in TemplateType.java.
func (m *CallSiteMap) Bind(tmpl *TemplateType, t JSType) {
	if tmpl == nil || t == nil {
		return
	}
	if prev, ok := m.subs[tmpl]; ok {
		if prev.TypeName() == t.TypeName() {
			return
		}
		m.subs[tmpl] = NewUnion(prev, t)
		return
	}
	m.subs[tmpl] = t
}

// Lookup returns the substitution for tmpl (or nil if the template is
// unresolved, e.g. because no argument positions mention it).
// Go-only: language-adapter — call-site unification helper with no Java analog in TemplateType.java.
func (m *CallSiteMap) Lookup(tmpl *TemplateType) JSType {
	if m == nil {
		return nil
	}
	return m.subs[tmpl]
}

// Substitute walks `t` replacing every TemplateType with its mapped
// concrete type. ObjectType / FunctionType / UnionType / Templatized
// are handled structurally. Unmapped template vars stay as
// TemplateType so a downstream caller can still decide how to
// degrade them (our convention is to treat them as Unknown).
// Go-only: language-adapter — call-site substitution helper with no Java analog in TemplateType.java.
func (m *CallSiteMap) Substitute(t JSType) JSType {
	if m == nil || t == nil {
		return t
	}
	switch v := t.(type) {
	case *TemplateType:
		if sub := m.Lookup(v); sub != nil {
			return sub
		}
		return t
	case *UnionType:
		var out []JSType
		changed := false
		for _, m2 := range v.Members() {
			s := m.Substitute(m2)
			out = append(out, s)
			if s != m2 {
				changed = true
			}
		}
		if !changed {
			return t
		}
		return NewUnion(out...)
	case *TemplatizedType:
		var args []JSType
		changed := false
		for _, a := range v.Args {
			s := m.Substitute(a)
			args = append(args, s)
			if s != a {
				changed = true
			}
		}
		if !changed {
			return t
		}
		return NewTemplatizedType(v.Base, args)
	case *FunctionType:
		var params []JSType
		changed := false
		for _, p := range v.Params {
			s := m.Substitute(p)
			params = append(params, s)
			if s != p {
				changed = true
			}
		}
		ret := m.Substitute(v.Return)
		if ret != v.Return {
			changed = true
		}
		if !changed {
			return t
		}
		out := &FunctionType{
			Params:        params,
			Return:        ret,
			IsConstructor: v.IsConstructor,
			Instance:      v.Instance,
			TemplateTypes: v.TemplateTypes,
		}
		return out
	}
	return t
}
