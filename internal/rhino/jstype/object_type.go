// Port of com/google/javascript/rhino/jstype/PrototypeObjectType.java +
// InstanceObjectType.java + ObjectType.java @ closure-compiler v20260415
//
// Closure's ObjectType hierarchy is large (templates, structural subtype
// for @record, late-resolution flags, member descriptor flags …). This
// port covers the structural members Tier-2 actually exercises:
//   - a member map (name → declared JSType)
//   - a nominal name (the constructor's name, e.g. "Foo")
//   - a super pointer for @extends chains (added when Tier-2 lands
//     inheritance)
//
// Properties are nominal: two ObjectType instances compare equal only
// when they're literally the same Go value. That mirrors Closure's
// reference-equality on interned object types and keeps subtype checks
// O(1) for the common case.

package jstype

import (
	"fmt"
	"sort"
	"strings"
)

type ObjectType struct {
	// Name is the source-level identifier (e.g. "Foo" for `function
	// Foo(){}`). Empty for anonymous object literals.
	Name string

	// members is the declared property map, populated by
	// TypedScopeCreator as it walks `this.x = …`, `Foo.prototype.x = …`,
	// and instance class field declarations.
	members map[string]JSType

	// staticMembers is the separate slot for `static` class members.
	// Accessed via the constructor name (`C.foo`) rather than an
	// instance (`(new C).foo`). Mirrors FunctionType.getOwnerFunction's
	// own-property map in Java.
	staticMembers map[string]JSType

	// super is the @extends superclass instance type. nil at top of
	// hierarchy. (Inheritance lookup walks this chain.)
	super *ObjectType

	// extendedInterfaces carries every @extends expression for a
	// multi-inherited interface (e.g. `@extends {A} @extends {B}`).
	// Mirrors FunctionType.getExtendedInterfaces() in Java. Empty for
	// concrete classes (which use `super` for their single base) and
	// for interfaces with at most one @extends. Callers that need the
	// full supertype graph (INCOMPATIBLE_EXTENDED_PROPERTY_TYPE,
	// INTERFACE_METHOD_NOT_IMPLEMENTED) walk this list.
	extendedInterfaces []*ObjectType

	// templatizedSuper is the @extends supertype expression when it
	// carried template arguments (e.g. `@extends {Bar<!Array<T>>}`).
	// Mirrors ObjectType.getImplicitPrototype() returning a
	// TemplatizedType in Java. `super` still points at the underlying
	// Bar instance so nominal subtype walks work; this slot adds the
	// substitution map used when resolving inherited property types.
	templatizedSuper *TemplatizedType

	// TemplateParams are the class-level @template parameters in the
	// order declared. Used when a TemplatizedType pairs this ObjectType
	// with concrete Args to build the substitution map for inherited
	// property lookup. Empty for non-generic classes. Mirrors Java's
	// ObjectType.getTemplateTypeMap keys.
	TemplateParams []*TemplateType

	// IsStructural is true for `@record` / `@interface` types. A value
	// of type S is assignable to a structural T iff S has every
	// property T declares with a compatible type — no nominal check.
	IsStructural bool

	// IsStruct / IsDict mirror ObjectType.isStruct() / isDict() in Java —
	// flags set from `@struct` / `@dict` on the constructor JSDoc.
	// A @struct instance forbids `x['k']` computed access; a @dict
	// instance forbids `x.k` dotted access. Both flags surface via
	// JSC_ILLEGAL_PROPERTY_ACCESS at the TypeCheck visitor.
	IsStruct bool
	IsDict   bool

	// inferredProps carries property names whose type came from an
	// RHS literal inference (e.g. `this.x = true`) rather than a
	// declared @type. TypeCheck's assignment-type mismatch guard
	// skips inferred properties — Java's default-mode behaviour
	// allows open reassignment when the slot wasn't @type'd.
	inferredProps map[string]bool
}

func NewObjectType(name string) *ObjectType {
	return &ObjectType{
		Name:          name,
		members:       map[string]JSType{},
		staticMembers: map[string]JSType{},
		inferredProps: map[string]bool{},
	}
}

// MarkInferred flags property `name` as literal-inferred. Callers
// that register a property type from an RHS literal (rather than
// an explicit @type annotation) mark it so visitAssign can skip
// the subtype enforcement on later writes.
func (o *ObjectType) MarkInferred(name string) {
	if o.inferredProps == nil {
		o.inferredProps = map[string]bool{}
	}
	o.inferredProps[name] = true
}

// IsPropertyInferred reports whether `name` was registered via
// literal inference.
func (o *ObjectType) IsPropertyInferred(name string) bool {
	return o.inferredProps[name]
}

// DefineStaticProperty mirrors ObjectType.defineStaticProperty — puts
// the property on the constructor-owned slot rather than the instance
// map. Used for `static x = …` class fields and `static method(){…}`.
func (o *ObjectType) DefineStaticProperty(name string, t JSType) {
	o.staticMembers[name] = t
}

// GetStaticPropertyType looks up a member on the static slot only —
// does not walk the instance map or the @extends chain.
func (o *ObjectType) GetStaticPropertyType(name string) JSType {
	if t, ok := o.staticMembers[name]; ok {
		return t
	}
	return nil
}

// HasStaticProperty is the existence check for the static slot.
func (o *ObjectType) HasStaticProperty(name string) bool {
	return o.GetStaticPropertyType(name) != nil
}

func (o *ObjectType) TypeName() string {
	if o.Name != "" {
		return o.Name
	}
	// Anonymous structural type: render `{k: T, …}` so diagnostics
	// carry useful surface. Sorted keys give deterministic output.
	if len(o.members) == 0 {
		return "{}"
	}
	keys := make([]string, 0, len(o.members))
	for k := range o.members {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	b.WriteString("{")
	for i, k := range keys {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(k)
		b.WriteString(": ")
		b.WriteString(o.members[k].TypeName())
	}
	b.WriteString("}")
	return b.String()
}

// IsSubtypeOf handles three cases: nominal subtyping via the @extends
// chain (default), structural subtyping when the target is
// `@record`/`@interface`, and the two top-of-lattice sinks (AllType,
// Unknown). Mirrors ObjectType.isSubtype in Java.
func (o *ObjectType) IsSubtypeOf(other JSType) bool {
	if other == nil {
		return true
	}
	if isAll(other) || isUnknown(other) {
		return true
	}
	if tt, ok := other.(*TemplatizedType); ok && tt.Base != nil {
		// Accept `{...}` vs `Object<K, V>` — every ObjectType is a
		// subtype of the global Object.
		if tt.Base.Name == "Object" {
			return true
		}
		// Array IS-A Iterable<T> / IArrayLike<T> / Iterator<T> in
		// Closure's built-in type graph. Mirrors NativeType
		// pre-registration in JSTypeRegistry.
		if o.Name == "Array" {
			switch tt.Base.Name {
			case "Iterable", "IArrayLike", "Iterator", "IterableIterator", "ReadonlyArray":
				return true
			}
		}
		// Nominal @extends of a templatized target: walk our nominal super
		// chain looking for the target's Base, then check the matching
		// templatized-super link's args against the target's args.
		for cur := o; cur != nil; cur = cur.super {
			if cur.templatizedSuper != nil && cur.templatizedSuper.Base == tt.Base {
				return cur.templatizedSuper.IsSubtypeOf(tt)
			}
			if cur == tt.Base {
				return true
			}
		}
		return false
	}
	target, ok := other.(*ObjectType)
	if !ok {
		return false
	}
	if target.IsStructural {
		// Structural: self must have each of target's declared
		// properties with a compatible type. Walks the super chain so
		// inherited members count. A target property whose type
		// includes undefined is optional — the source is allowed to
		// omit it entirely.
		for pname, ptype := range target.members {
			srcType := o.GetPropertyType(pname)
			if srcType == nil {
				if containsVoid(ptype) {
					continue
				}
				return false
			}
			if !IsSubtype(srcType, ptype) {
				return false
			}
		}
		return true
	}
	// Every ObjectType is a subtype of the global `Object` constructor's
	// instance type — Java's ObjectType treats Object as the root of the
	// nominal hierarchy. Without explicit super-wiring for every anon
	// literal, compare by name: the externs-registered `Object` instance
	// has Name == "Object".
	if target.Name == "Object" {
		return true
	}
	// Nominal: same type identity, or target is an ancestor.
	for cur := o; cur != nil; cur = cur.super {
		if cur == target {
			return true
		}
	}
	return false
}

// DefineProperty mirrors ObjectType.defineProperty(name, type). Later
// definitions overwrite earlier ones — the Java side has more nuanced
// behaviour (declared-vs-inferred + duplicate warnings), which we'll
// add when Tier-2 fixtures need them.
func (o *ObjectType) DefineProperty(name string, t JSType) {
	o.members[name] = t
}

// GetPropertyType walks the @extends chain looking for the named
// property. Returns nil if the property is undeclared anywhere in the
// chain — distinct from declared-as-Unknown. When a link in the chain
// was wired as a TemplatizedType (`@extends {Bar<T>}`), the property
// type is rewritten via that link's substitution map — mirroring
// Java's ObjectType.getPropertyType which threads the substitution
// through getImplicitPrototype() at each step.
func (o *ObjectType) GetPropertyType(name string) JSType {
	if t, ok := o.members[name]; ok {
		return t
	}
	// Walk the templatized super first when present — the substitution
	// must be applied at every step for nested generic inheritance
	// (`Foo<T> extends Bar<Array<T>>` resolving a Bar-declared property
	// that mentions Bar's template parameter).
	if o.templatizedSuper != nil {
		if t := o.templatizedSuper.GetPropertyType(name); t != nil {
			return t
		}
		return nil
	}
	if o.super != nil {
		return o.super.GetPropertyType(name)
	}
	return nil
}

// HasProperty is the existence check used by TypeCheck.visitGetProp to
// decide whether to emit JSC_INEXISTENT_PROPERTY.
func (o *ObjectType) HasProperty(name string) bool {
	return o.GetPropertyType(name) != nil
}

// MemberCount returns the number of directly-declared members on this
// type (not including the super chain). Used by the checker to skip
// inexistent-property diagnostics on stub types whose member map is
// still empty.
func (o *ObjectType) MemberCount() int { return len(o.members) }

// MemberNames returns the set of directly-declared member names on
// this type. Caller owns the slice; iteration order is undefined.
// Used by passes that need to fan out over an interface's properties
// (e.g. JSC_INTERFACE_METHOD_NOT_IMPLEMENTED).
func (o *ObjectType) MemberNames() []string {
	out := make([]string, 0, len(o.members))
	for k := range o.members {
		out = append(out, k)
	}
	return out
}

// OwnMembers returns the directly-declared member map (name → type),
// excluding anything inherited via the super chain. Mirrors
// ObjectType.getOwnPropertyNames iterated with getOwnPropertyType.
// The returned map aliases the internal slot — callers must treat it
// as read-only.
func (o *ObjectType) OwnMembers() map[string]JSType { return o.members }

// SetSuper wires the @extends chain. Called by TypedScopeCreator when
// it processes `@extends {Base}` on a constructor.
func (o *ObjectType) SetSuper(s *ObjectType) { o.super = s }

// Super returns the immediate superclass instance type (as wired by
// SetSuper) or nil at the root of the chain.
func (o *ObjectType) Super() *ObjectType { return o.super }

// AddExtendedInterface appends an entry to the multi-extends list used
// by interfaces that declare more than one `@extends`. Matches Java's
// FunctionType.extendInterface.
func (o *ObjectType) AddExtendedInterface(s *ObjectType) {
	if s == nil {
		return
	}
	o.extendedInterfaces = append(o.extendedInterfaces, s)
}

// ExtendedInterfaces returns the slice of multi-inherited interface
// parents registered via AddExtendedInterface. Matches Java's
// FunctionType.getExtendedInterfaces().
func (o *ObjectType) ExtendedInterfaces() []*ObjectType {
	return o.extendedInterfaces
}

// containsVoid reports whether `t` is (or is a union containing)
// the Void primitive. Used by the structural-subtype check to
// treat record-type properties with `{…: (T|undefined)}` as
// optional on the source side. Java's JSType.isExplicitlyVoidable.
func containsVoid(t JSType) bool {
	if t == Void {
		return true
	}
	if u, ok := t.(*UnionType); ok {
		for _, m := range u.members {
			if containsVoid(m) {
				return true
			}
		}
	}
	return false
}
// SetTemplatizedSuper records the templatized form of the @extends
// link — e.g. `@extends {Bar<Array<T>>}` on a `@template T` class.
// The underlying nominal super is still expected to be wired through
// SetSuper; this slot carries only the extra substitution bindings.
// Mirrors ObjectType.setImplicitPrototype when Java sees a
// TemplatizedType on the @extends position.
func (o *ObjectType) SetTemplatizedSuper(t *TemplatizedType) { o.templatizedSuper = t }

// TemplatizedSuper returns the templatized @extends link (or nil if
// the super was wired as a bare ObjectType). Callers that need the
// substitution-aware property walk should prefer this over Super().
func (o *ObjectType) TemplatizedSuper() *TemplatizedType { return o.templatizedSuper }

// String is for diagnostics; matches Java's PrototypeObjectType.toString
// shape closely enough that the message column of jsgo output is
// recognisable to Closure users.
func (o *ObjectType) String() string {
	return fmt.Sprintf("ObjectType(%s)", o.TypeName())
}
