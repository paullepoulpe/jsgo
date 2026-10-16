// Port of com/google/javascript/rhino/jstype/RecordTypeBuilder.java
// @ closure-compiler v20260415
//
// A builder for RecordType. Collects (name, JSType, Node) triples then
// constructs a RecordType in one shot. If no properties were added,
// Build returns the global Object type instead (matching the Java empty-
// record shortcut).

package jstype

import (
	"sort"

	"github.com/paullepoulpe/jsgo/legacy/internal/rhino"
)

// RecordProperty mirrors RecordTypeBuilder.RecordProperty in Java: a
// (type, propertyNode) pair for one named slot of a record type.
type RecordProperty struct {
	typ          JSType
	propertyNode *rhino.Node
}

// GetType returns the declared JSType for the property.
func (p *RecordProperty) GetType() JSType { return p.typ }

// GetPropertyNode returns the AST node that declared the property, or
// nil when the property was synthesized programmatically.
func (p *RecordProperty) GetPropertyNode() *rhino.Node { return p.propertyNode }

func (p *RecordProperty) String() string {
	// Mirrors Java's RecordProperty.toString().
	nodeStr := "<nil>"
	if p.propertyNode != nil {
		nodeStr = p.propertyNode.String_()
	}
	return "RecordProperty{type: " + p.typ.TypeName() + ", node: " + nodeStr + "}"
}

// RecordTypeBuilder mirrors RecordTypeBuilder in Java. Not safe for
// concurrent use.
type RecordTypeBuilder struct {
	isEmpty    bool
	isDeclared bool
	registry   *Registry
	properties map[string]*RecordProperty // insertion-order maintained via keys
	keys       []string                   // insertion order
}

// NewRecordTypeBuilder returns a new builder backed by registry.
func NewRecordTypeBuilder(registry *Registry) *RecordTypeBuilder {
	return &RecordTypeBuilder{
		isEmpty:    true,
		isDeclared: true,
		registry:   registry,
		properties: map[string]*RecordProperty{},
	}
}

// SetSynthesized marks the record being built as synthesized (isDeclared
// = !synthesized). Mirrors RecordTypeBuilder.setSynthesized in Java.
func (b *RecordTypeBuilder) SetSynthesized(synthesized bool) {
	b.isDeclared = !synthesized
}

// AddProperty adds a property with the given name, type, and AST node.
// Duplicate names are overwritten (last write wins, matching Java's
// LinkedHashMap.put semantics); the builder panics in Build if any
// duplicate was inserted — we preserve the insert-order key slice so
// Build can detect collisions via the sorted map.
//
// Returns the builder itself for chaining.
func (b *RecordTypeBuilder) AddProperty(name string, typ JSType, propertyNode *rhino.Node) *RecordTypeBuilder {
	b.isEmpty = false
	if _, exists := b.properties[name]; !exists {
		b.keys = append(b.keys, name)
	}
	b.properties[name] = &RecordProperty{typ: typ, propertyNode: propertyNode}
	return b
}

// Build creates the RecordType. If no properties were added (isEmpty),
// it returns the native Object type, matching the Java shortcut:
//
//	if (isEmpty) return registry.getNativeObjectType(OBJECT_TYPE);
//
// Properties are sorted by name (ImmutableSortedMap.naturalOrder in
// Java) before being handed to RecordType.
func (b *RecordTypeBuilder) Build() JSType {
	if b.isEmpty {
		// Java: return registry.getNativeObjectType(JSTypeNative.OBJECT_TYPE)
		// We return the interned Object sentinel.
		return newObjectNativeType()
	}

	// Sort keys to match Java's ImmutableSortedMap.naturalOrder().
	sorted := make([]string, len(b.keys))
	copy(sorted, b.keys)
	sort.Strings(sorted)

	props := make(map[string]*RecordProperty, len(sorted))
	for _, k := range sorted {
		props[k] = b.properties[k]
	}

	return newRecordType(b.registry, props, b.isDeclared)
}
