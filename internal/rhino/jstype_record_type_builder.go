// Port of com/google/javascript/rhino/jstype/RecordTypeBuilder.java @ closure-compiler v20260415

package rhino

import "github.com/paullepoulpe/jsgo/internal/jutil"

// RecordTypeBuilder is a builder for record types.
type RecordTypeBuilder struct {
	isEmpty    bool
	isDeclared bool
	registry   *JSTypeRegistry
	properties *jutil.LinkedHashMap[string, *RecordTypeBuilder_RecordProperty]
}

// Java: RecordTypeBuilder(JSTypeRegistry)
func NewRecordTypeBuilder(registry *JSTypeRegistry) *RecordTypeBuilder {
	return &RecordTypeBuilder{
		isEmpty:    true,
		isDeclared: true,
		registry:   registry,
		properties: jutil.NewLinkedHashMap[string, *RecordTypeBuilder_RecordProperty](),
	}
}

// SetSynthesized: see the comments on RecordType about synthetic types.
//
// Java: RecordTypeBuilder.setSynthesized(boolean)
func (b *RecordTypeBuilder) SetSynthesized(synthesized bool) {
	b.isDeclared = !synthesized
}

// AddProperty adds a property with the given name and type to the record
// type. If you add a property that has already been added, then #build will
// fail.
//
// name: the name of the new property. type: the JSType of the new property.
// propertyNode: the node that holds this property definition. Returns the
// builder itself for chaining purposes.
//
// Java: RecordTypeBuilder.addProperty(String, JSType, Node)
func (b *RecordTypeBuilder) AddProperty(name string, type_ JSType, propertyNode *Node) *RecordTypeBuilder {
	b.isEmpty = false
	b.properties.Put(name, NewRecordTypeBuilder_RecordProperty(type_, propertyNode))
	return b
}

// Build creates a record. Fails if any duplicate property names were added.
//
// Java: RecordTypeBuilder.build()
func (b *RecordTypeBuilder) Build() JSType {
	// If we have an empty record, simply return the object type.
	if b.isEmpty {
		return b.registry.GetNativeObjectType(JSTypeNative_OBJECT_TYPE)
	}
	m := jutil.NewTreeMap[string, *RecordTypeBuilder_RecordProperty](jutil.StringCompareTo)
	for k, v := range b.properties.All() {
		m.Put(k, v)
	}
	return NewRecordType(b.registry, m, b.isDeclared)
}

// RecordTypeBuilder_RecordProperty is RecordTypeBuilder.RecordProperty.
type RecordTypeBuilder_RecordProperty struct {
	type_        JSType
	propertyNode *Node
}

// Java: RecordTypeBuilder.RecordProperty(JSType, Node)
func NewRecordTypeBuilder_RecordProperty(type_ JSType, propertyNode *Node) *RecordTypeBuilder_RecordProperty {
	return &RecordTypeBuilder_RecordProperty{type_: type_, propertyNode: propertyNode}
}

// Java: RecordTypeBuilder.RecordProperty.getType()
func (p *RecordTypeBuilder_RecordProperty) GetType() JSType {
	return p.type_
}

// Java: RecordTypeBuilder.RecordProperty.getPropertyNode()
func (p *RecordTypeBuilder_RecordProperty) GetPropertyNode() *Node {
	return p.propertyNode
}

// Java: RecordTypeBuilder.RecordProperty.toString()
func (p *RecordTypeBuilder_RecordProperty) String() string {
	return "RecordProperty{type: " + jutil.ToJavaString(p.type_) + ", node: " + jutil.ToJavaString(p.propertyNode) + "}"
}
