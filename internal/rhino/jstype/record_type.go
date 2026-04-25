// Port of com/google/javascript/rhino/jstype/RecordType.java
// @ closure-compiler v20260415
//
// A record (structural) type.
//
// Subtyping: The subtyping of a record type is defined via structural
// comparison of a record type's properties. For example, a record type of
// the form { a : TYPE_1 } is a supertype of a record type of the form
// { b : TYPE_2, a : TYPE_1 } because B can be assigned to A and matches all
// constraints. Similarly, a defined type can be assigned to a record type so
// long as that defined type matches all property constraints of the record
// type. A record type of the form { a : A, b : B } can be assigned to a
// record of type { a : A }.
//
// In Go the Java class hierarchy (RecordType → PrototypeObjectType →
// ObjectType) is flattened: RecordType embeds an *ObjectType for the shared
// property-map behaviour and carries its own declared/isFrozen flags.

package jstype

import (
	"sort"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RecordType mirrors RecordType in Java. It is a structural (record)
// type whose assignability is determined by property-shape rather than
// nominal identity.
type RecordType struct {
	// ObjectType carries the shared property-map, TypeName, and
	// IsSubtypeOf logic inherited from PrototypeObjectType / ObjectType in
	// Java. RecordType overrides some methods below.
	*ObjectType

	// registry is the JSTypeRegistry under which this type lives.
	// Mirrors the `registry` inherited field in Java's PrototypeObjectType.
	// Required by GetGreatestSubtypeHelper to construct new RecordTypeBuilder
	// instances.
	registry *Registry

	// declared mirrors the `declared` field in Java — true for types
	// written explicitly in user code, false for synthesized types used
	// only for internal bookkeeping (e.g. greatest-subtype helpers).
	declared bool

	// isFrozen mirrors the `isFrozen` field in Java. Set to true at the
	// end of construction; subsequent calls to DefineProperty return false.
	isFrozen bool

	// propertyNodes stores the AST Node associated with each declared
	// property — parallels the Node stored inside ObjectType.Property in
	// Java's PrototypeObjectType.
	propertyNodes map[string]*rhino.Node
}

// newRecordType is the package-internal constructor called by
// RecordTypeBuilder.Build. It mirrors the package-private
// RecordType(JSTypeRegistry, Map<String,RecordProperty>, boolean)
// constructor in Java.
func newRecordType(registry *Registry, properties map[string]*RecordProperty, declared bool) *RecordType {
	rt := &RecordType{
		ObjectType:    NewObjectType(""),
		registry:      registry,
		declared:      declared,
		isFrozen:      false,
		propertyNodes: map[string]*rhino.Node{},
	}
	// Mark as structural so ObjectType.IsSubtypeOf applies structural
	// matching. Mirrors Java's RecordType.isStructuralType() = true and
	// the SubtypeChecker's use of isRecordType().
	rt.ObjectType.IsStructural = true

	for _, name := range sortedRecordKeys(properties) {
		prop := properties[name]
		if prop == nil {
			panic("RecordProperty associated with a property should not be nil!")
		}
		// Java calls defineDeclaredProperty / defineSynthesizedProperty —
		// both ultimately call PrototypeObjectType.defineProperty which
		// stores the type. We call DefineProperty here; inferred vs.
		// declared is tracked on the ObjectType via MarkInferred.
		rt.ObjectType.DefineProperty(name, prop.GetType())
		if !declared {
			rt.ObjectType.MarkInferred(name)
		}
		rt.propertyNodes[name] = prop.GetPropertyNode()
	}

	// Freeze the record type — no further property definitions allowed.
	rt.isFrozen = true

	return rt
}

// sortedRecordKeys returns the map keys in sorted order (matches Java's
// ImmutableSortedMap iteration which is always natural-order).
func sortedRecordKeys(m map[string]*RecordProperty) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// newObjectNativeType returns an ObjectType named "Object" as the
// "native Object type" sentinel used when RecordTypeBuilder.Build is
// called on an empty builder. Mirrors:
//
//	registry.getNativeObjectType(JSTypeNative.OBJECT_TYPE)
func newObjectNativeType() JSType {
	return NewObjectType("Object")
}

// IsSynthetic mirrors RecordType.isSynthetic() in Java. A synthesized
// record type is one built internally (e.g. during greatest-subtype
// computation) rather than declared in user code.
func (r *RecordType) IsSynthetic() bool {
	return !r.declared
}

// GetPropertyNode returns the AST Node associated with the named
// property, or nil if the property is unknown or was synthesized without
// a source node. Mirrors ObjectType.getPropertyNode(String) in Java.
func (r *RecordType) GetPropertyNode(name string) *rhino.Node {
	return r.propertyNodes[name]
}

// DefineProperty overrides ObjectType.DefineProperty: once the record
// is frozen (after construction), new properties are rejected (returns
// false). Mirrors RecordType.defineProperty(…) in Java.
func (r *RecordType) DefineProperty(name string, t JSType) bool {
	if r.isFrozen {
		return false
	}
	r.ObjectType.DefineProperty(name, t)
	return true
}

// IsRecordType returns true — this IS a RecordType. Mirrors
// JSType.isRecordType() / RecordType.toMaybeRecordType() != null in Java.
func (r *RecordType) IsRecordType() bool { return true }

// ToMaybeRecordType returns r itself, matching Java's
// RecordType.toMaybeRecordType() override.
func (r *RecordType) ToMaybeRecordType() *RecordType { return r }

// IsStructuralType returns true, matching Java's
// RecordType.isStructuralType() override.
func (r *RecordType) IsStructuralType() bool { return true }

// TypeName renders the record in `{k: T, …}` form, sorted by key —
// matches the Java RecordType's pretty-print output (setPrettyPrint(true)
// in the constructor routes through TypeStringBuilder which produces this
// form for record types).
func (r *RecordType) TypeName() string {
	names := r.ObjectType.MemberNames()
	sort.Strings(names)
	if len(names) == 0 {
		return "{}"
	}
	var b strings.Builder
	b.WriteString("{")
	for i, k := range names {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(k)
		b.WriteString(": ")
		t := r.ObjectType.GetPropertyType(k)
		if t != nil {
			b.WriteString(t.TypeName())
		} else {
			b.WriteString("?")
		}
	}
	b.WriteString("}")
	return b.String()
}

// IsSubtypeOf overrides ObjectType.IsSubtypeOf with the structural
// record-type rule from Java's SubtypeChecker: a type S is a subtype of
// record R iff S has every property R declares with a compatible type.
// The embedded ObjectType already handles structural checking because
// IsStructural is set to true; we delegate there.
func (r *RecordType) IsSubtypeOf(other JSType) bool {
	if other == nil {
		return true
	}
	if isAll(other) || isUnknown(other) {
		return true
	}
	return r.ObjectType.IsSubtypeOf(other)
}

// GetGreatestSubtypeHelper mirrors RecordType.getGreatestSubtypeHelper
// in Java. Called by the type system when the receiver is a RecordType.
//
// When `that` is also a RecordType, the result is a new synthetic
// RecordType whose properties are the union of both, with conflicting
// property types narrowed to their greatest subtype. If any conflict
// collapses to the bottom type (NO_TYPE — represented as nil here),
// nil is returned.
//
// When `that` is some other object type, the result is the join of all
// reference types that (a) are subtypes of `that` and (b) declare each
// of r's properties with a compatible type.
func (r *RecordType) GetGreatestSubtypeHelper(that JSType) JSType {
	if thatRecord, ok := that.(*RecordType); ok {
		builder := NewRecordTypeBuilder(r.registry)
		builder.SetSynthesized(true)

		// The greatest subtype consists of those *unique* properties of both
		// record types. If any property conflicts, NO_TYPE is returned.
		for _, prop := range r.ObjectType.MemberNames() {
			thisPropertyType := r.ObjectType.GetPropertyType(prop)
			var propType JSType
			if thatRecord.ObjectType.HasProperty(prop) {
				thatPropertyType := thatRecord.ObjectType.GetPropertyType(prop)
				propType = recordGetGreatestSubtype(thisPropertyType, thatPropertyType)
				if propType == nil {
					// Conflict — return bottom type (NO_TYPE).
					return nil
				}
			} else {
				propType = thisPropertyType
			}
			builder.AddProperty(prop, propType, r.GetPropertyNode(prop))
		}

		for _, prop := range thatRecord.ObjectType.MemberNames() {
			if !r.ObjectType.HasProperty(prop) {
				builder.AddProperty(prop, thatRecord.ObjectType.GetPropertyType(prop),
					thatRecord.GetPropertyNode(prop))
			}
		}
		return builder.Build()
	}

	// `that` is some other type. Java computes this as the join of all
	// object reference types that are subtypes of `that` and have
	// compatible property types (via registry.getEachReferenceTypeWithProperty).
	// That registry enumeration is not yet ported; return Unknown as a
	// conservative approximation.
	//
	// TODO: implement via registry enumeration when
	// JSTypeRegistry.getEachReferenceTypeWithProperty is ported.
	return Unknown
}

// recordGetGreatestSubtype returns the greatest subtype of two JSTypes
// for the record-vs-record case. Mirrors the call to
// thisPropertyType.getGreatestSubtype(thatPropertyType) in Java.
// Returns nil when the types have no common subtype (NO_TYPE).
func recordGetGreatestSubtype(a, b JSType) JSType {
	if a == nil || b == nil {
		return nil
	}
	// If either is Unknown, propagate Unknown (conservative).
	if isUnknown(a) || isUnknown(b) {
		return Unknown
	}
	// If they're the same type (by name), the greatest subtype is that type.
	if a.TypeName() == b.TypeName() {
		return a
	}
	// If one is a subtype of the other, the more specific wins.
	if IsSubtype(a, b) {
		return a
	}
	if IsSubtype(b, a) {
		return b
	}
	// No common subtype — bottom / NO_TYPE.
	return nil
}
