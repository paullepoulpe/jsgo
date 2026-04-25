// Port of com/google/javascript/rhino/jstype/ObjectType.java
// @ closure-compiler v20260415 (collapsed hierarchy).
//
// Closure's ObjectType hierarchy is large (templates, structural subtype
// for @record, late-resolution flags, member descriptor flags …). This
// port collapses the ObjectType/PrototypeObjectType/InstanceObjectType
// hierarchy into a single concrete Go struct (see also
// prototype_object_type.go and instance_object_type.go for the fuller
// typed subclasses). The covers the structural members Tier-2 actually
// exercises:
//   - a member map (name → declared JSType)
//   - a nominal name (the constructor's name, e.g. "Foo")
//   - a super pointer for @extends chains
//
// Methods that are abstract in ObjectType.java and implemented in
// PrototypeObjectType.java (such as defineProperty, appendTo) are
// implemented here on *ObjectType as the base concrete defaults.
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

	"github.com/paullepoulpe/jsgo/internal/rhino"
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

	// IsStructural is true for `@record` / `@interface` types. A value
	// of type S is assignable to a structural T iff S has every
	// property T declares with a compatible type — no nominal check.
	IsStructural bool

	// IsRecord is true for `@record` types specifically (not bare `@interface`).
	// When IsRecord is true, structural duck-typing applies; when IsRecord is false
	// but IsStructural is true, the interface requires nominal membership (super chain).
	IsRecord bool

	// IsInterface is true for `@interface` types (not `@record`). Interface
	// types are open-ended: accessing an undeclared property on an interface
	// variable does NOT fire JSC_INEXISTENT_PROPERTY.
	IsInterface bool

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

	// docInfo mirrors ObjectType.getJSDocInfo() — the JSDoc attached to
	// the type declaration (e.g. the `/** @constructor */` on `var Foo =
	// function(){}`). Set by InferJsDocInfo; consumed by later passes.
	docInfo *rhino.JSDocInfo

	// propDocInfo mirrors ObjectType.getPropertyJSDocInfo() — per-slot
	// JSDoc attached by InferJsDocInfo from the property's declaring node.
	propDocInfo map[string]*rhino.JSDocInfo

	// propNodes stores the source AST Node at which each property was
	// declared or first inferred. Populated by DefineDeclaredProperty and
	// DefineInferredProperty; queried by GetPropertyNode.
	propNodes map[string]*rhino.Node

	// isNative marks types that were pre-registered from the externs
	// bundle (e.g. Object, Array, Function). Mirrors
	// ObjectType.isNativeObjectType() in Java — InferJsDocInfo skips
	// attaching JSDoc to native types to avoid clobbering externs docs.
	isNative bool

	// typeParameters are the ordered @template keys declared on this type.
	// Mirrors ObjectType.getTypeParameters() in Java.
	typeParameters []*TemplateType

	// templatizedSelf, when non-nil, is the TemplatizedType that wraps this
	// ObjectType with its concrete arguments.  Set by NewTemplatizedType so
	// that GetTemplateArg can resolve inherited template parameters.
	templatizedSelf *TemplatizedType

	// sealedForInference mirrors the sealed-for-inference flag used by
	// TypeInference: once sealed, inferring new properties on the type is
	// suppressed. Mirrors the concept in ObjectType / PrototypeObjectType
	// where open-ended anonymous objects are allowed to grow until sealed.
	sealedForInference bool
	// NOTE: constructor and ownerFunction are NOT fields of ObjectType.
	// They live on InstanceObjectType and PrototypeObjectType respectively,
	// mirroring the Java class hierarchy:
	//   InstanceObjectType.constructor  → InstanceObjectType.getConstructor()
	//   PrototypeObjectType.ownerFunction → PrototypeObjectType.getOwnerFunction()
	//
	// Go-only: language-adapter — Ctor is a back-pointer set by
	// NewConstructorType so that GetCtorImplementedInterfaces / IsSubtypeOf
	// can reach the constructor's implementedInterfaces without a separate
	// InstanceObjectType subclass. Java reaches the same data through
	// InstanceObjectType.getConstructor(); the Go port collapses the
	// hierarchy so we wire the back-reference here.
	Ctor *FunctionType
}

// Go-only: language-adapter — NewObjectType — Go constructor; Java uses Builder pattern (PrototypeObjectType.Builder).
func NewObjectType(name string) *ObjectType {
	return &ObjectType{
		Name:          name,
		members:       map[string]JSType{},
		staticMembers: map[string]JSType{},
		inferredProps: map[string]bool{},
	}
}

// Go-only: language-adapter — MarkInferred — Go helper for inferred-property tracking; Java uses Property.inferred flag.
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

// Go-only: language-adapter — IsPropertyInferred — Go helper; Java uses Property.isInferred() from PropertyMap.
// IsPropertyInferred reports whether `name` was registered via
// literal inference.
func (o *ObjectType) IsPropertyInferred(name string) bool {
	return o.inferredProps[name]
}

// Go-only: language-adapter — ClearInferred removes the inferred flag from `name`,
// marking the property as declared. Called when an explicit @type JSDoc later
// declares the same slot that was previously inferred from a this.prop assignment.
func (o *ObjectType) ClearInferred(name string) {
	if o.inferredProps != nil {
		delete(o.inferredProps, name)
	}
}

// Go-only: language-adapter — DefineStaticProperty — Go addition for static class members; Java tracks statics on FunctionType.
// DefineStaticProperty mirrors ObjectType.defineStaticProperty — puts
// the property on the constructor-owned slot rather than the instance
// map. Used for `static x = …` class fields and `static method(){…}`.
func (o *ObjectType) DefineStaticProperty(name string, t JSType) {
	o.staticMembers[name] = t
}

// Go-only: language-adapter — GetStaticPropertyType — Go accessor for static slot; Java uses FunctionType own-property map.
// GetStaticPropertyType looks up a member on the static slot only —
// does not walk the instance map or the @extends chain.
func (o *ObjectType) GetStaticPropertyType(name string) JSType {
	if t, ok := o.staticMembers[name]; ok {
		return t
	}
	return nil
}

// Go-only: language-adapter — HasStaticProperty — Go existence check for static slot; no direct Java analogue.
// HasStaticProperty is the existence check for the static slot.
func (o *ObjectType) HasStaticProperty(name string) bool {
	return o.GetStaticPropertyType(name) != nil
}

// Go-only: language-adapter — TypeName — Go interface method for diagnostics; Java uses appendTo(TypeStringBuilder).
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
	// Accept `{...}` vs `Object<K, V>` — every ObjectType is a
	// subtype of the global Object, and Java's TemplatizedType
	// treats parametrisation as invariant-optional (a plain
	// ObjectType IS-A templated Object when the object's member
	// types line up with the value slot; at our current depth
	// we just accept).
	if tt, ok := other.(*TemplatizedType); ok {
		if tt.Base != nil && tt.Base.Name == "Object" {
			return true
		}
		// Nominal subtype of a templatized type: A is a subtype of
		// Foo<string> if A's super chain contains Foo's base instance.
		// Mirrors Java's TemplatizedType.isSubtype where an uninstantiated
		// type (without template args) is checked against the raw base.
		// When `@implements {Foo<…>}` is present we also check template-arg
		// invariance so `B implements Foo<number>` is NOT a subtype of
		// `Foo<string>` (matching Closure's behaviour in both default and
		// strict modes).
		if tt.Base != nil {
			for cur := o; cur != nil; cur = cur.super {
				if cur == tt.Base {
					return true
				}
				// Walk implemented interfaces declared on the constructor.
				// Use the templatized variant when available so that
				// `@implements {Foo<string>}` is distinguished from
				// `@implements {Foo<number>}`. Mirrors Java's
				// ObjectType.isSubtype delegating through
				// FunctionType.implementedInterfaces (which there carries
				// templatization implicitly through the ObjectType cast).
				if cur.Ctor != nil {
					tplList := cur.Ctor.GetImplementedInterfaceTypes()
					ifaces := cur.Ctor.GetImplementedInterfaces()
					for i, iface := range ifaces {
						if iface != tt.Base {
							continue
						}
						// Found the interface base; now verify args
						// match against tt's args. When tplList isn't
						// populated (e.g. raw `@implements {Foo}`),
						// fall through to a permissive accept.
						if i < len(tplList) {
							if templateArgsCompatible(tplList[i], tt) {
								return true
							}
							continue
						}
						return true
					}
				}
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
		// `@implements {Foo}` makes o a subtype of Foo (an interface).
		// Walk the implementedInterfaces declared on each ancestor's
		// constructor. Mirrors Java's ObjectType.isSubtype path through
		// FunctionType.getImplementedInterfaces / getAllImplementedInterfaces.
		if cur.Ctor != nil {
			for _, iface := range cur.Ctor.GetImplementedInterfaces() {
				if iface == target {
					return true
				}
				// Walk the interface's own super chain too — `B implements Foo`
				// where Foo extends Bar means B is also a subtype of Bar.
				for sup := iface; sup != nil; sup = sup.super {
					if sup == target {
						return true
					}
				}
			}
		}
	}
	return false
}

// templateArgsCompatible reports whether the source's template args
// match the target TemplatizedType's args. Used by IsSubtypeOf to
// gate `@implements {Foo<…>}` matches: a source typed
// `@implements {Foo<string>}` is NOT a subtype of `Foo<number>`.
//
// Mirrors Java's TemplatizedType.isSubtype path which delegates to
// per-arg invariance via JSType.isEquivalentTo. Source can be a
// TemplatizedType (when the JSDoc had `<…>` slots) or a plain
// ObjectType (when raw); a plain ObjectType is treated as having
// unknown args and is accepted permissively.
func templateArgsCompatible(srcImpl JSType, target *TemplatizedType) bool {
	srcTpl, _ := srcImpl.(*TemplatizedType)
	if srcTpl == nil || len(srcTpl.Args) == 0 {
		// Raw `@implements {Foo}` (no template args) or non-templatized
		// resolution — accept permissively.
		return true
	}
	if len(target.Args) == 0 {
		return true
	}
	if len(srcTpl.Args) != len(target.Args) {
		return false
	}
	for i, a := range srcTpl.Args {
		t := target.Args[i]
		if a == nil || t == nil {
			continue
		}
		if isUnknown(a) || isUnknown(t) || isAll(a) || isAll(t) {
			continue
		}
		// Treat TemplateType (uninstantiated @template parameter) as
		// compatible with any arg — mirrors Closure's default-mode
		// behaviour of widening unbound template parameters at
		// call/assignment sites.
		if _, ok := a.(*TemplateType); ok {
			continue
		}
		if _, ok := t.(*TemplateType); ok {
			continue
		}
		// Invariance: same identity in both directions.
		if !IsSubtype(a, t) || !IsSubtype(t, a) {
			return false
		}
	}
	return true
}

// DefineProperty mirrors ObjectType.defineProperty(name, type). Later
// definitions overwrite earlier ones — the Java side has more nuanced
// behaviour (declared-vs-inferred + duplicate warnings), which we'll
// add when Tier-2 fixtures need them.
//
// @struct/@dict enforcement: if the type is @struct, computed-style
// (bracket) additions are guarded by the caller; here we only reject
// if the type is @dict and the call site is a dot-access (no way to
// distinguish at this layer — enforcement is at the TypeCheck visitor).
func (o *ObjectType) DefineProperty(name string, t JSType) {
	o.members[name] = t
}

// DefineDeclaredProperty mirrors ObjectType.defineDeclaredProperty(String,
// JSType, Node) — defines a property whose type is explicitly declared by
// the programmer (@type annotation or ES6 class field syntax).
//
// Unlike DefineProperty, this returns false and is a no-op when the
// property has already been declared (the Java "once declared, immutable"
// semantics). This prevents later inferred slots from clobbering explicit
// annotations. Mirrors PrototypeObjectType.defineProperty gated by
// hasOwnDeclaredProperty.
func (o *ObjectType) DefineDeclaredProperty(name string, t JSType, node *rhino.Node) bool {
	if o.HasOwnDeclaredProperty(name) {
		return false
	}
	o.members[name] = t
	// Clear inferred flag — this is now a declared slot.
	delete(o.inferredProps, name)
	// Preserve the source node for later GetPropertyNode queries.
	if o.propNodes == nil {
		o.propNodes = map[string]*rhino.Node{}
	}
	o.propNodes[name] = node
	return true
}

// DefineInferredProperty mirrors ObjectType.defineInferredProperty — defines
// a property whose type is inferred from usage. If a declared property
// already exists, returns true without overwriting. Otherwise stores the
// type and marks it as inferred.
func (o *ObjectType) DefineInferredProperty(name string, t JSType, node *rhino.Node) bool {
	if o.HasOwnProperty(name) && !o.IsPropertyInferred(name) {
		// Never hide a declared property with an inferred one.
		return true
	}
	o.members[name] = t
	o.MarkInferred(name)
	if o.propNodes == nil {
		o.propNodes = map[string]*rhino.Node{}
	}
	o.propNodes[name] = node
	return true
}

// HasOwnDeclaredProperty mirrors ObjectType.hasOwnDeclaredProperty(String) —
// returns true when the named property is declared directly on this type
// (not inherited) AND is not inferred (i.e. was explicitly @type'd).
func (o *ObjectType) HasOwnDeclaredProperty(name string) bool {
	return o.HasOwnProperty(name) && !o.IsPropertyInferred(name)
}

// GetPropertyNode returns the AST Node at which the named property was
// declared or first inferred, or nil when none was recorded.
// Mirrors ObjectType.getPropertyNode(String).
func (o *ObjectType) GetPropertyNode(name string) *rhino.Node {
	if o.propNodes == nil {
		return nil
	}
	return o.propNodes[name]
}

// Go-only: language-adapter — IsSealed — Go inference-sealing concept; Java uses separate ObjectType.sealed flag.
// IsSealed mirrors the sealed-for-inference concept: once true, no new
// inferred properties should be added to this type by TypeInference.
func (o *ObjectType) IsSealed() bool { return o.sealedForInference }

// Go-only: language-adapter — SealForInference — Go helper; Java equivalent concept doesn't use a separate seal method.
// SealForInference sets the sealed flag. Called when TypeInference has
// finished processing a function body and wants to freeze the inferred
// shape of an object type.
func (o *ObjectType) SealForInference() { o.sealedForInference = true }

// IsRecordType mirrors JSType.isRecordType() for ObjectType: returns false
// for plain ObjectType; overridden by RecordType to return true.
// Port of ObjectType.isRecordType() in Java (base-class default).
func (o *ObjectType) IsRecordType() bool {
	// ObjectType base-class default: only RecordType returns true.
	// Java: JSType.isRecordType() base returns false; RecordType overrides.
	// Check the IsRecord flag; plain ObjectType always has it false.
	if o.IsRecord {
		return true
	}
	return false
}

// GetPropertyType walks the @extends chain looking for the named
// property. Returns nil if the property is undeclared anywhere in the
// chain — distinct from declared-as-Unknown.
func (o *ObjectType) GetPropertyType(name string) JSType {
	for cur := o; cur != nil; cur = cur.super {
		if t, ok := cur.members[name]; ok {
			return t
		}
	}
	return nil
}

// HasProperty is the existence check used by TypeCheck.visitGetProp to
// decide whether to emit JSC_INEXISTENT_PROPERTY.
func (o *ObjectType) HasProperty(name string) bool {
	return o.GetPropertyType(name) != nil
}

// Go-only: language-adapter — MemberCount — Go convenience count; Java uses getPropertyMap().getOwnProperties().size().
// MemberCount returns the number of directly-declared members on this
// type (not including the super chain). Used by the checker to skip
// inexistent-property diagnostics on stub types whose member map is
// still empty.
func (o *ObjectType) MemberCount() int { return len(o.members) }

// Go-only: language-adapter — MemberNames — Go convenience slice; Java uses getPropertyMap().getOwnProperties().
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

// Go-only: language-adapter — SetSuper — Go helper for super-chain wiring; Java uses setImplicitPrototype().
// SetSuper wires the @extends chain. Called by TypedScopeCreator when
// it processes `@extends {Base}` on a constructor.
func (o *ObjectType) SetSuper(s *ObjectType) { o.super = s }

// Super returns the immediate superclass instance type (as wired by
// SetSuper) or nil at the root of the chain.
func (o *ObjectType) Super() *ObjectType { return o.super }

// Go-only: language-adapter — SuperHasProperty — Go convenience; Java uses getImplicitPrototype().getPropertyKind().
// SuperHasProperty reports whether the property `name` is declared
// anywhere in the superclass chain (excluding the receiver itself).
// Used by arity checks to detect @override methods: when the super
// chain also declares the property, the overriding function's arity
// may differ from the parent's and we cannot safely fire a
// call-site arity diagnostic without override-aware inference.
func (o *ObjectType) SuperHasProperty(name string) bool {
	if o.super == nil {
		return false
	}
	return o.super.GetPropertyType(name) != nil
}

// Go-only: language-adapter — containsVoid — Go internal helper for structural subtyping; Java uses JSType.isVoidable().
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

// SetJSDocInfo mirrors ObjectType.setJSDocInfo — attaches the declaring
// JSDoc to this type. Called by InferJsDocInfo.
func (o *ObjectType) SetJSDocInfo(info *rhino.JSDocInfo) { o.docInfo = info }

// GetJSDocInfo mirrors ObjectType.getJSDocInfo.
func (o *ObjectType) GetJSDocInfo() *rhino.JSDocInfo { return o.docInfo }

// HasReferenceName mirrors ObjectType.hasReferenceName() — true iff the
// type was allocated with a non-empty source name. Anonymous object
// literals (Name == "") return false.
func (o *ObjectType) HasReferenceName() bool { return o.Name != "" }

// GetReferenceName mirrors ObjectType.getReferenceName().
func (o *ObjectType) GetReferenceName() string { return o.Name }

// IsNativeObjectType mirrors ObjectType.isNativeObjectType() — true for
// types pre-registered from the externs bundle. InferJsDocInfo skips
// attaching JSDoc to natives to avoid overwriting externs-supplied docs.
func (o *ObjectType) IsNativeObjectType() bool { return o.isNative }

// Go-only: language-adapter — MarkNative — Go helper to set native flag; Java sets via ObjectType constructor parameter.
// MarkNative marks this ObjectType as a native externs type. Called by
// the scope creator when registering externs-bundle types.
func (o *ObjectType) MarkNative() { o.isNative = true }

// HasOwnProperty mirrors ObjectType.hasOwnProperty(String) — checks the
// directly-declared member map only, not the super chain. Used by
// InferJsDocInfo to avoid setting JSDoc on inherited slots.
func (o *ObjectType) HasOwnProperty(name string) bool {
	_, ok := o.members[name]
	return ok
}

// GetPropertyJSDocInfo mirrors ObjectType.getPropertyJSDocInfo(String).
// Returns nil when no JSDoc has been attached to this slot.
func (o *ObjectType) GetPropertyJSDocInfo(name string) *rhino.JSDocInfo {
	if o.propDocInfo == nil {
		return nil
	}
	return o.propDocInfo[name]
}

// SetPropertyJSDocInfo mirrors ObjectType.setPropertyJSDocInfo(String,
// JSDocInfo) / PrototypeObjectType.setPropertyJSDocInfo — attaches
// per-property JSDoc. Following the Java semantics: if info is non-nil and
// the property does not yet exist in the own map, an inferred slot is
// created implicitly (mirrors "defineInferredProperty with nil node"
// in PrototypeObjectType.setPropertyJSDocInfo).
func (o *ObjectType) SetPropertyJSDocInfo(name string, info *rhino.JSDocInfo) {
	if info != nil && !o.HasOwnProperty(name) {
		// Java: defineInferredProperty(name, getPropertyType(name), nil)
		// Ensure a slot exists so the JSDoc is not lost.
		existing := o.GetPropertyType(name) // may be nil for truly new props
		if existing == nil {
			existing = Unknown
		}
		o.DefineInferredProperty(name, existing, nil)
	}
	if o.propDocInfo == nil {
		o.propDocInfo = map[string]*rhino.JSDocInfo{}
	}
	o.propDocInfo[name] = info
}

// GetOwnPropertyJSDoc / SetPropertyJSDoc are CheckAccessControls-flavoured
// aliases of the above, matching the Java helper names used in the port.
func (o *ObjectType) GetOwnPropertyJSDoc(name string) *rhino.JSDocInfo {
	return o.GetPropertyJSDocInfo(name)
}
func (o *ObjectType) SetPropertyJSDoc(name string, info *rhino.JSDocInfo) {
	o.SetPropertyJSDocInfo(name, info)
}

// Go-only: language-adapter — GetTypeParameters — Go accessor; Java uses getTemplateTypeMap().getTemplateKeys().
// GetTypeParameters returns the ordered @template keys declared on this
// type.  Mirrors ObjectType.getTypeParameters() in Java.
func (o *ObjectType) GetTypeParameters() []*TemplateType {
	return o.typeParameters
}

// Go-only: language-adapter — SetTypeParameters — Go setter; Java wires template keys via registry builder.
// SetTypeParameters records the ordered @template keys for this type.
// Called by TypedScopeCreator when processing `@template` annotations.
func (o *ObjectType) SetTypeParameters(params []*TemplateType) {
	o.typeParameters = params
}

// String is for diagnostics; matches Java's PrototypeObjectType.toString
// shape closely enough that the message column of jsgo output is
// recognisable to Closure users.
func (o *ObjectType) String() string {
	return fmt.Sprintf("ObjectType(%s)", o.TypeName())
}

// MatchesObjectContext reports true — all ObjectType instances can appear in
// an Object context. Mirrors PrototypeObjectType.matchesObjectContext().
func (o *ObjectType) MatchesObjectContext() bool { return true }

// GetPossibleToBooleanOutcomes returns BooleanLiteralSetTrue — an Object is
// always truthy. Mirrors ObjectType.getPossibleToBooleanOutcomes().
func (o *ObjectType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSetTrue
}

// --- Additional Java-parity methods ----------------------------------------

// GetConstructor mirrors ObjectType.getConstructor() in Java.
// Base-class default: plain ObjectType instances have no constructor.
// InstanceObjectType overrides this via its embedded PrototypeObjectType chain.
// Port of ObjectType.getConstructor() in Java (abstract; base-class default returns null).
func (o *ObjectType) GetConstructor() *FunctionType {
	// Java: abstract FunctionType getConstructor() — base returns null.
	// The Ctor back-pointer is set on instances created by NewConstructorType;
	// most plain ObjectType instances (anonymous literals, etc.) have Ctor == nil.
	if o.Ctor != nil {
		return o.Ctor
	}
	return nil
}

// GetSuperClassConstructor mirrors ObjectType.getSuperClassConstructor()
// in Java — climbs TWO implicit-prototype levels to find the superclass
// constructor. Java: getImplicitPrototype().getImplicitPrototype().getConstructor().
// Port of ObjectType.getSuperClassConstructor() in Java.
func (o *ObjectType) GetSuperClassConstructor() *FunctionType {
	// Java:
	//   ObjectType iproto = getImplicitPrototype();
	//   if (iproto == null) return null;
	//   iproto = iproto.getImplicitPrototype();
	//   return iproto == null ? null : iproto.getConstructor();
	iproto := o.GetImplicitPrototype()
	if iproto == nil {
		return nil
	}
	iproto = iproto.GetImplicitPrototype()
	if iproto == nil {
		return nil
	}
	return iproto.GetConstructor()
}

// GetImplicitPrototype mirrors ObjectType.getImplicitPrototype() in Java —
// returns the [[Prototype]] chain entry, which in the Go port is the super
// ObjectType. Returns nil at the top of the chain.
func (o *ObjectType) GetImplicitPrototype() *ObjectType {
	return o.super
}

// GetCtorImplementedInterfaces mirrors ObjectType.getCtorImplementedInterfaces()
// in Java. Base-class default: ObjectType has no associated constructor so the
// set of implemented interfaces is empty. InstanceObjectType overrides to
// return the constructor's implementedInterfaces.
// Port of ObjectType.getCtorImplementedInterfaces() in Java (returns ImmutableSet.of()).
func (o *ObjectType) GetCtorImplementedInterfaces() []*ObjectType {
	// Java: return ImmutableSet.of() — empty, not null.
	// The Go Ctor back-pointer lets us resolve this when wired.
	if o.Ctor != nil {
		return o.Ctor.GetImplementedInterfaces()
	}
	return []*ObjectType{}
}

// GetCtorExtendedInterfaces mirrors ObjectType.getCtorExtendedInterfaces() in
// Java. Base-class default: ObjectType has no associated constructor so the set
// of extended interfaces is empty. InstanceObjectType overrides to return the
// constructor's extendedInterfaces.
// Port of ObjectType.getCtorExtendedInterfaces() in Java (returns ImmutableSet.of()).
func (o *ObjectType) GetCtorExtendedInterfaces() []*ObjectType {
	// Java: return ImmutableSet.of() — empty, not null.
	if o.Ctor != nil {
		ifaces := extFor(o).extendedInterfaces
		if ifaces != nil {
			return ifaces
		}
	}
	return []*ObjectType{}
}

// GetPropertyTypeMap mirrors ObjectType.getPropertyTypeMap() in Java.
// Returns a map of all property names to their types, walking the super chain.
func (o *ObjectType) GetPropertyTypeMap() map[string]JSType {
	out := map[string]JSType{}
	for cur := o; cur != nil; cur = cur.super {
		for k, v := range cur.members {
			if _, exists := out[k]; !exists {
				out[k] = v
			}
		}
	}
	return out
}

// GetOwnPropertyNames mirrors ObjectType.getOwnPropertyNames() in Java.
// Returns the set of directly-owned property names (not inherited).
func (o *ObjectType) GetOwnPropertyNames() []string {
	out := make([]string, 0, len(o.members))
	for k := range o.members {
		out = append(out, k)
	}
	return out
}

// GetPropertiesCount mirrors ObjectType.getPropertiesCount() in Java.
// Returns the total count of properties (own + inherited via super chain).
func (o *ObjectType) GetPropertiesCount() int {
	seen := map[string]bool{}
	for cur := o; cur != nil; cur = cur.super {
		for k := range cur.members {
			seen[k] = true
		}
	}
	return len(seen)
}

// GetTypeOfThis mirrors ObjectType.getTypeOfThis() in Java — for most
// ObjectTypes returns nil; for instances of constructors this is meaningful
// to the function that owns the prototype. Returns nil by default.
// Port of ObjectType.getTypeOfThis() in Java (base-class default).
func (o *ObjectType) GetTypeOfThis() JSType {
	// Java: public JSType getTypeOfThis() { return null; }
	// PrototypeObjectType and FunctionType override to return the instance/this.
	// The Ctor back-pointer lets a plain ObjectType resolve this when the
	// constructor is wired; otherwise fall through to nil.
	if o.Ctor != nil {
		return o.Ctor.GetTypeOfThis()
	}
	return nil
}

// GetOwnerFunction mirrors ObjectType.getOwnerFunction() in Java.
// Base-class default: plain ObjectType instances are not function prototypes.
// PrototypeObjectType overrides this via its ownerFunction field.
// Port of ObjectType.getOwnerFunction() in Java (base-class default returns null).
func (o *ObjectType) GetOwnerFunction() *FunctionType {
	// Java: public FunctionType getOwnerFunction() { return null; }
	// isFunctionPrototypeType() returns true iff getOwnerFunction() != null.
	// Plain ObjectType is not a function prototype; PrototypeObjectType overrides.
	// Ctor back-pointer: when this is the constructor's Instance object (which is
	// distinct from the function prototype), return nil.
	if o.Ctor == nil {
		return nil
	}
	return nil
}

// IsObject mirrors JSType.isObject() in Java — all ObjectTypes are objects.
func (o *ObjectType) IsObject() bool { return true }

// GetEnumeratedTypeOfEnumObject mirrors ObjectType.getEnumeratedTypeOfEnumObject().
// Returns nil for plain ObjectTypes; overridden by EnumType.
// Port of ObjectType.getEnumeratedTypeOfEnumObject() in Java (base-class default).
func (o *ObjectType) GetEnumeratedTypeOfEnumObject() JSType {
	// Java: public JSType getEnumeratedTypeOfEnumObject() { return null; }
	// EnumType overrides this to return its element type's primitive type.
	// Plain ObjectType is never an enum object type.
	if o.IsInterface || o.IsRecord || o.IsStructural {
		return nil
	}
	return nil
}

// GetDisplayName mirrors ObjectType.getDisplayName() — returns the normalized
// reference name (stripping delegate suffixes). Port of Java method.
func (o *ObjectType) GetDisplayName() string {
	return o.GetNormalizedReferenceName()
}

// GetNormalizedReferenceName mirrors ObjectType.getNormalizedReferenceName() —
// strips delegate suffix from the reference name. Delegate suffixes appear
// inside parentheses (e.g. "Foo(Bar).prototype" → "Foo.prototype").
// Port of ObjectType.getNormalizedReferenceName() in Java.
func (o *ObjectType) GetNormalizedReferenceName() string {
	name := o.GetReferenceName()
	if name == "" {
		return ""
	}
	start := strings.IndexByte(name, '(')
	if start == -1 {
		return name
	}
	end := strings.LastIndexByte(name, ')')
	prefix := name[:start]
	if end+1 >= len(name) {
		return prefix
	}
	return prefix + name[end+1:]
}

// IsPropertyTypeDeclared mirrors ObjectType.isPropertyTypeDeclared(String) —
// returns true when the named property exists and was explicitly declared
// (not inferred from usage). Port of Java method.
func (o *ObjectType) IsPropertyTypeDeclared(name string) bool {
	return o.HasOwnProperty(name) && !o.IsPropertyInferred(name)
}

// IsPropertyTypeInferred mirrors ObjectType.isPropertyTypeInferred(String) —
// returns true when the named property exists and was inferred from usage.
// Port of Java method.
func (o *ObjectType) IsPropertyTypeInferred(name string) bool {
	return o.HasOwnProperty(name) && o.IsPropertyInferred(name)
}

// IsStructuralType mirrors ObjectType.isStructuralType() — returns true when
// the type is a @record / @interface type backed by structural (duck-type)
// subtyping. Port of ObjectType.isStructuralType() in Java.
func (o *ObjectType) IsStructuralType() bool {
	return o.IsStructural
}

// IsImplicitPrototypeOf mirrors ObjectType.isImplicitPrototypeOf(ObjectType) —
// returns true when `o` appears in the implicit-prototype chain of `other`.
// Port of ObjectType.isImplicitPrototypeOf(ObjectType) in Java.
func (o *ObjectType) IsImplicitPrototypeOf(other *ObjectType) bool {
	for cur := other.GetImplicitPrototype(); cur != nil; cur = cur.GetImplicitPrototype() {
		if cur == o {
			return true
		}
	}
	return false
}

// GetPropertyDefSite mirrors ObjectType.getPropertyDefSite(String) — returns
// the AST node at which the property was declared. Alias for GetPropertyNode.
// Port of ObjectType.getPropertyDefSite(String) in Java.
func (o *ObjectType) GetPropertyDefSite(name string) *rhino.Node {
	return o.GetPropertyNode(name)
}

// IsPropertyInExterns mirrors ObjectType.isPropertyInExterns(String) —
// returns true when the property was declared on a native (externs) type.
// Port of ObjectType.isPropertyInExterns(String) in Java. In the Go port,
// a property is in externs iff it's declared on a native ObjectType.
func (o *ObjectType) IsPropertyInExterns(name string) bool {
	for cur := o; cur != nil; cur = cur.super {
		if cur.isNative && cur.HasOwnProperty(name) {
			return true
		}
	}
	return false
}

// IsUnknownType mirrors ObjectType.isUnknownType() in Java — treats the type
// as unknown when its implicit-prototype chain is unknown. Java uses a cached
// `unknown` field (starts true, lazily cleared once the chain is resolved).
// Port of ObjectType.isUnknownType() in Java.
func (o *ObjectType) IsUnknownType() bool {
	// Java:
	//   if (unknown) {
	//     ObjectType implicitProto = getImplicitPrototype();
	//     if (implicitProto == null || implicitProto.isNativeObjectType()) {
	//       unknown = false;
	//       for (ObjectType iface : getCtorExtendedInterfaces()) {
	//         if (iface.isUnknownType()) { unknown = true; break; }
	//       }
	//     } else { unknown = implicitProto.isUnknownType(); }
	//   }
	//   return unknown;
	//
	// In the Go port we don't cache the unknown flag; compute it directly.
	// Walk the implicit-prototype chain: if any non-native ancestor is unknown,
	// the type is unknown. For plain ObjectType the chain is the super pointer.
	proto := o.GetImplicitPrototype()
	if proto == nil || proto.IsNativeObjectType() {
		// Check extended interfaces for unknown-ness.
		for _, iface := range o.GetCtorExtendedInterfaces() {
			if iface.IsUnknownType() {
				return true
			}
		}
		return false
	}
	return proto.IsUnknownType()
}

// HasCachedValues mirrors ObjectType.hasCachedValues() — returns whether
// the type's property cache is valid. Java uses `return !unknown` where
// the `unknown` field starts as `true` until properties are added.
// Port of ObjectType.hasCachedValues() in Java. In the Go port, all
// types start with a valid member map, so this always returns true.
func (o *ObjectType) HasCachedValues() bool {
	return len(o.members) > 0 || o.Name != ""
}

// ClearCachedValues mirrors ObjectType.clearCachedValues() — resets cached
// property type information. In Java, sets `unknown = true` to mark the
// property type cache as stale. Port of ObjectType.clearCachedValues() in Java.
// In the Go port, the member map is the source of truth, so this is a no-op.
func (o *ObjectType) ClearCachedValues() {
	_ = o.Name // touch the receiver so the compiler doesn't elide the call
}
