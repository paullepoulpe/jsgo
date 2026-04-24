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
	// Java's FunctionType.isStructuralInterface() returns true only for @record;
	// plain @interface uses nominal-compatible matching (requires explicit @implements).
	// When IsRecord is true, structural duck-typing applies; when IsRecord is false
	// but IsStructural is true, the interface requires nominal membership (super chain).
	IsRecord bool

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

	// constructor mirrors the back-pointer from an instance ObjectType to its
	// constructor FunctionType. Set by SetConstructor when the TypedScopeCreator
	// wires up the constructor/instance relationship. Mirrors the Java
	// relationship accessed via InstanceObjectType.getConstructor().
	constructor *FunctionType

	// ownerFunction mirrors the PrototypeObjectType.ownerFunction field in Java.
	// For prototype objects, this is the constructor FunctionType that owns the
	// prototype. Nil for regular instance objects. Set via SetOwnerFunction.
	ownerFunction *FunctionType
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
		// Template-argument invariance is not enforced here (matching Java's
		// default-mode behaviour where concrete extends Foo<string> can be
		// assigned to Foo<number> — TypeScript would reject but Closure
		// doesn't in default mode).
		if tt.Base != nil {
			for cur := o; cur != nil; cur = cur.super {
				if cur == tt.Base {
					return true
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
	}
	return false
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

// IsSealed mirrors the sealed-for-inference concept: once true, no new
// inferred properties should be added to this type by TypeInference.
func (o *ObjectType) IsSealed() bool { return o.sealedForInference }

// SealForInference sets the sealed flag. Called when TypeInference has
// finished processing a function body and wants to freeze the inferred
// shape of an object type.
func (o *ObjectType) SealForInference() { o.sealedForInference = true }

// IsRecordType mirrors JSType.isRecordType() for ObjectType: returns false
// for plain ObjectType; overridden by RecordType to return true.
func (o *ObjectType) IsRecordType() bool { return false }

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

// SetSuper wires the @extends chain. Called by TypedScopeCreator when
// it processes `@extends {Base}` on a constructor.
func (o *ObjectType) SetSuper(s *ObjectType) { o.super = s }

// Super returns the immediate superclass instance type (as wired by
// SetSuper) or nil at the root of the chain.
func (o *ObjectType) Super() *ObjectType { return o.super }

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

// GetTypeParameters returns the ordered @template keys declared on this
// type.  Mirrors ObjectType.getTypeParameters() in Java.
func (o *ObjectType) GetTypeParameters() []*TemplateType {
	return o.typeParameters
}

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

// GetConstructor mirrors ObjectType.getConstructor() in Java — returns the
// FunctionType whose instance this object is, or nil for literal objects.
// In the Go port this is stored via the constructor back-pointer field.
func (o *ObjectType) GetConstructor() *FunctionType {
	return o.constructor
}

// SetConstructor records the constructor FunctionType for this instance.
// Called by TypedScopeCreator when wiring up the instance/constructor
// relationship (mirrors Java's setImplicitPrototype → setOwnerFunction chain).
func (o *ObjectType) SetConstructor(ctor *FunctionType) {
	o.constructor = ctor
}

// GetSuperClassConstructor mirrors ObjectType.getSuperClassConstructor()
// in Java — climbs two implicit-prototype levels to find the superclass
// constructor. In the Go port: super.GetConstructor().
func (o *ObjectType) GetSuperClassConstructor() *FunctionType {
	if o.super == nil {
		return nil
	}
	return o.super.GetConstructor()
}

// GetImplicitPrototype mirrors ObjectType.getImplicitPrototype() in Java —
// returns the [[Prototype]] chain entry, which in the Go port is the super
// ObjectType. Returns nil at the top of the chain.
func (o *ObjectType) GetImplicitPrototype() *ObjectType {
	return o.super
}

// GetCtorImplementedInterfaces mirrors ObjectType.getCtorImplementedInterfaces()
// in Java — returns the @implements interfaces registered via the constructor.
// Delegates to the constructor's implementedInterfaces list.
func (o *ObjectType) GetCtorImplementedInterfaces() []*ObjectType {
	if o.constructor != nil {
		return o.constructor.GetImplementedInterfaces()
	}
	return nil
}

// GetCtorExtendedInterfaces mirrors ObjectType.getCtorExtendedInterfaces() in
// Java — returns the @extends interfaces for an interface's constructor.
// Delegates to the constructor's extendedInterfaces list.
func (o *ObjectType) GetCtorExtendedInterfaces() []*ObjectType {
	if o.constructor != nil {
		return o.constructor.GetExtendedInterfaces()
	}
	return nil
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
func (o *ObjectType) GetTypeOfThis() JSType {
	return nil
}

// GetOwnerFunction mirrors ObjectType.getOwnerFunction() in Java — returns
// nil for regular objects. For prototype objects (PrototypeObjectType in Java)
// this is set to the constructor.
func (o *ObjectType) GetOwnerFunction() *FunctionType {
	return o.ownerFunction
}

// SetOwnerFunction mirrors ObjectType.setOwnerFunction() in Java — records
// the constructor that owns this prototype.
func (o *ObjectType) SetOwnerFunction(fn *FunctionType) {
	o.ownerFunction = fn
}

// IsObject mirrors JSType.isObject() in Java — all ObjectTypes are objects.
func (o *ObjectType) IsObject() bool { return true }

// GetEnumeratedTypeOfEnumObject mirrors ObjectType.getEnumeratedTypeOfEnumObject().
// Returns nil for plain ObjectTypes; overridden by EnumType.
func (o *ObjectType) GetEnumeratedTypeOfEnumObject() JSType {
	return nil
}
