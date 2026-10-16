// Port of com/google/javascript/rhino/jstype/ObjectType.java @ closure-compiler v20260415
//
// Additive extension to object_type.go — adds the property-map accessors and
// prototype-chain methods consumed by PropertyMap. These are split into a
// separate file so object_type.go (which covers the members map) is not
// rewritten. Only the methods required by Property.java / PropertyMap.java
// are added here; the full ObjectType surface is in object_type.go.

package jstype

// propertyMap is the Property-based property map for this ObjectType.
// nil until the first PropertyMap-based DefinePropertyKeyed call is made.
//
// NOTE: The existing object_type.go uses a plain map[string]JSType for
// members. PropertyMap is an additional parallel store used by the
// property.go / property_map.go layer. Both stores can coexist;
// callers using the PropertyMap API go through GetPropertyMap(), callers
// using the legacy API go through DefineProperty / GetPropertyType.
type objectTypePropertyMapField struct {
	propertyMap *PropertyMap
	// implicitPrototype points to the __proto__ object type. Mirrors
	// ObjectType.getImplicitPrototype() in Java.
	implicitPrototype *ObjectType
	// extendedInterfaces is the list of interfaces this type extends.
	// Used by PropertyMap.getSecondaryParentObjects().
	extendedInterfaces []*ObjectType
}

// We store the extra fields in a map keyed by *ObjectType pointer so we
// do not alter the ObjectType struct (which is declared in object_type.go
// and must not be rewritten).
//
// This is an alternative to embedding a field directly in the struct.
// The map approach is safe because ObjectType pointers are stable
// (never moved by the GC in Go).
var objExtFields = map[*ObjectType]*objectTypePropertyMapField{}

// Go-only: language-adapter — extFor — Go extension-field accessor; Java uses class fields directly.
func extFor(o *ObjectType) *objectTypePropertyMapField {
	if f, ok := objExtFields[o]; ok {
		return f
	}
	f := &objectTypePropertyMapField{}
	objExtFields[o] = f
	return f
}

// GetPropertyMap mirrors ObjectType.getPropertyMap(). Returns the
// PropertyMap for this object. Lazily initialised the first time.
func (o *ObjectType) GetPropertyMap() *PropertyMap {
	f := extFor(o)
	if f.propertyMap == nil {
		pm := NewPropertyMap()
		pm.SetParentSource(o)
		f.propertyMap = pm
	}
	return f.propertyMap
}

// getImplicitPrototype mirrors ObjectType.getImplicitPrototype(). Returns
// the prototype of this type for the purposes of property lookup.
//
// Backed by the super field in object_type.go when no explicit implicit-
// prototype has been set via SetImplicitPrototype.
func (o *ObjectType) getImplicitPrototype() *ObjectType {
	f := extFor(o)
	if f.implicitPrototype != nil {
		return f.implicitPrototype
	}
	// Fall back to the super pointer set by SetSuper (legacy API).
	return o.super
}

// Go-only: language-adapter — SetImplicitPrototype — Go setter; Java ObjectType.setImplicitPrototype is package-private.
// SetImplicitPrototype mirrors ObjectType.setImplicitPrototype(ObjectType).
// Sets the __proto__ for property-chain traversal.
func (o *ObjectType) SetImplicitPrototype(proto *ObjectType) {
	extFor(o).implicitPrototype = proto
}

// getCtorExtendedInterfaces mirrors ObjectType.getCtorExtendedInterfaces().
// Returns the interfaces extended by the interface associated with this type.
func (o *ObjectType) getCtorExtendedInterfaces() []*ObjectType {
	return extFor(o).extendedInterfaces
}

// Go-only: language-adapter — AddExtendedInterface — Go mutator; Java uses Builder.setExtendedInterfaces().
// AddExtendedInterface appends an interface to the list of extended
// interfaces for this type. Used by the scope creator when wiring
// `@extends` on interface declarations.
func (o *ObjectType) AddExtendedInterface(iface *ObjectType) {
	f := extFor(o)
	f.extendedInterfaces = append(f.extendedInterfaces, iface)
}

// Go-only: language-adapter — DefinePropertyKeyed — Go PropertyMap-layer helper; Java uses defineDeclaredProperty(Key, JSType, Node).
// DefinePropertyKeyed mirrors ObjectType.defineDeclaredProperty(Key, JSType, Node)
// at the PropertyMap layer. Unlike the legacy DefineProperty which stores
// into the members map[string]JSType, this stores a full Property into the
// PropertyMap, supporting both string and symbol keys.
//
// Returns true if the property was newly registered (false if it already
// existed and was overwritten).
func (o *ObjectType) DefinePropertyKeyed(key Key, typ JSType, inferred bool) bool {
	pm := o.GetPropertyMap()
	existing := pm.GetOwnProperty(key)
	prop := NewPropertyFromKey(key, typ, inferred, nil)
	pm.PutProperty(key, prop)
	return existing == nil
}
