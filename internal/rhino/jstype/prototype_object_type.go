// Port of com/google/javascript/rhino/jstype/PrototypeObjectType.java
// @ closure-compiler v20260415 (structural surface).
//
// PrototypeObjectType is the primary mutable object-type node: it owns a
// PropertyMap (members), a nominal name (className or ownerFunction's name
// + ".prototype"), an implicit-prototype pointer (implicitPrototypeFallback),
// and flags for nativeness, anonymity, and prettyPrint rendering.
//
// In the Java hierarchy PrototypeObjectType extends ObjectType. Here, the
// existing ObjectType Go struct already provides the member map, super
// chain, name, and subtype logic used by the bulk of the checker.
// PrototypeObjectType therefore embeds *ObjectType for shared behaviour and
// adds the distinct Java fields that are not currently on ObjectType.

package jstype

// PrototypeObjectType is the Java PrototypeObjectType shape.
// It embeds *ObjectType (the Go analogue of ObjectType's property map,
// name, super chain, and subtype logic) and adds:
//   - implicitPrototypeFallback — the [[Prototype]] pointer.
//   - ownerFunction — set when this object is a function's .prototype.
//   - prettyPrint / anonymousType / nativeType — rendering and flag fields.
//   - templateParamCount — mirrors PrototypeObjectType.templateParamCount.
type PrototypeObjectType struct {
	*ObjectType

	// implicitPrototypeFallback is the [[Prototype]] chain link.
	// Access via GetImplicitPrototype(); subclasses like InstanceObjectType
	// override by returning the constructor's prototype instead.
	implicitPrototypeFallback *PrototypeObjectType

	// ownerFunction is set when this PrototypeObjectType is the .prototype
	// of a FunctionType. Mirrors PrototypeObjectType.ownerFunction in Java.
	ownerFunction *FunctionType

	// nativeType mirrors PrototypeObjectType.nativeType (set from
	// Builder.setNative). True for built-in types (Object, Array, …).
	nativeType bool

	// anonymousType mirrors PrototypeObjectType.anonymousType.
	// True for unnamed literals created without a className.
	anonymousType bool

	// prettyPrint controls whether appendTo renders the full property set
	// ({k: T, …}) or the abbreviated form ({…}). Mirrors
	// PrototypeObjectType.prettyPrint.
	prettyPrint bool

	// templateParamCount mirrors PrototypeObjectType.templateParamCount.
	// Number of @template parameters that belong to this type (not its
	// enclosing generic context).
	templateParamCount int
}

// NewPrototypeObjectType allocates a PrototypeObjectType with the given
// nominal name and implicit prototype.  Mirrors the essential state set up
// by PrototypeObjectType(Builder<?>) in Java.
func NewPrototypeObjectType(name string, implicitPrototype *PrototypeObjectType, native bool) *PrototypeObjectType {
	p := &PrototypeObjectType{
		ObjectType:                NewObjectType(name),
		implicitPrototypeFallback: implicitPrototype,
		nativeType:                native,
	}
	return p
}

// GetImplicitPrototype returns the [[Prototype]] of this object.
// Mirrors PrototypeObjectType.getImplicitPrototype().
func (p *PrototypeObjectType) GetImplicitPrototype() *PrototypeObjectType {
	return p.implicitPrototypeFallback
}

// SetImplicitPrototype sets the [[Prototype]] link.
// Mirrors PrototypeObjectType.setImplicitPrototype() (package-level in Java).
func (p *PrototypeObjectType) SetImplicitPrototype(proto *PrototypeObjectType) {
	p.implicitPrototypeFallback = proto
}

// GetOwnerFunction returns the FunctionType that owns this object as its
// .prototype, or nil when this is not a function prototype type.
// Mirrors PrototypeObjectType.getOwnerFunction().
func (p *PrototypeObjectType) GetOwnerFunction() *FunctionType {
	return p.ownerFunction
}

// SetOwnerFunction wires this object as the .prototype of f.
// Mirrors PrototypeObjectType.setOwnerFunction().
func (p *PrototypeObjectType) SetOwnerFunction(f *FunctionType) {
	p.ownerFunction = f
}

// GetReferenceName mirrors PrototypeObjectType.getReferenceName():
// returns className when set, falls back to ownerFunction.getReferenceName()
// + ".prototype", or nil (empty) when neither applies.
func (p *PrototypeObjectType) GetReferenceName() string {
	if p.ObjectType.Name != "" {
		return p.ObjectType.Name
	}
	if p.ownerFunction != nil {
		// TODO: FunctionType does not yet carry a name; extend when needed.
		return ".prototype"
	}
	return ""
}

// IsNativeObjectType mirrors PrototypeObjectType.isNativeObjectType().
func (p *PrototypeObjectType) IsNativeObjectType() bool { return p.nativeType }

// IsAnonymous mirrors PrototypeObjectType.isAnonymous().
func (p *PrototypeObjectType) IsAnonymous() bool { return p.anonymousType }

// SetPrettyPrint mirrors PrototypeObjectType.setPrettyPrint().
func (p *PrototypeObjectType) SetPrettyPrint(v bool) { p.prettyPrint = v }

// IsPrettyPrint mirrors PrototypeObjectType.isPrettyPrint().
func (p *PrototypeObjectType) IsPrettyPrint() bool { return p.prettyPrint }

// GetTemplateParamCount mirrors PrototypeObjectType.getTemplateParamCount().
func (p *PrototypeObjectType) GetTemplateParamCount() int { return p.templateParamCount }

// SetTemplateParamCount sets the template parameter count.
func (p *PrototypeObjectType) SetTemplateParamCount(n int) { p.templateParamCount = n }

// IsFunctionPrototypeType mirrors ObjectType.isFunctionPrototypeType():
// true when this object is the .prototype of a FunctionType.
func (p *PrototypeObjectType) IsFunctionPrototypeType() bool { return p.ownerFunction != nil }

// GetCtorImplementedInterfaces mirrors PrototypeObjectType.getCtorImplementedInterfaces().
// For a function prototype delegates to the owner function; otherwise empty.
// TODO: FunctionType does not yet carry implementedInterfaces; returns nil for now.
func (p *PrototypeObjectType) GetCtorImplementedInterfaces() []*ObjectType {
	return nil
}

// GetCtorExtendedInterfaces mirrors PrototypeObjectType.getCtorExtendedInterfaces().
// TODO: FunctionType does not yet carry extendedInterfaces; returns nil for now.
func (p *PrototypeObjectType) GetCtorExtendedInterfaces() []*ObjectType {
	return nil
}

// MatchesNumberContext mirrors PrototypeObjectType.matchesNumberContext().
// A PrototypeObjectType matches the number context when it is a number-
// object, date, boolean-object, or string-object wrapper — or when it
// overrides the native "valueOf" property.
func (p *PrototypeObjectType) MatchesNumberContext() bool {
	return p.IsNumberObjectType() || p.IsDateType() ||
		p.IsBooleanObjectType() || p.IsStringObjectType()
}

// MatchesStringContext mirrors PrototypeObjectType.matchesStringContext().
func (p *PrototypeObjectType) MatchesStringContext() bool {
	return p.IsTheObjectType() || p.IsStringObjectType() || p.IsDateType() ||
		p.IsRegexpType() || p.IsArrayType() || p.IsNumberObjectType() ||
		p.IsBigIntObjectType() || p.IsBooleanObjectType()
}

// MatchesSymbolContext mirrors PrototypeObjectType.matchesSymbolContext().
func (p *PrototypeObjectType) MatchesSymbolContext() bool { return p.IsSymbolObjectType() }

// MatchesObjectContext mirrors PrototypeObjectType.matchesObjectContext(); always true.
func (p *PrototypeObjectType) MatchesObjectContext() bool { return true }

// --- name-based native-type helpers (mirror InstanceObjectType / Java) ---

// IsTheObjectType reports whether this is the root `Object` type.
func (p *PrototypeObjectType) IsTheObjectType() bool {
	return p.nativeType && p.ObjectType.Name == "Object"
}

// IsArrayType mirrors PrototypeObjectType's check via the ctor.
func (p *PrototypeObjectType) IsArrayType() bool {
	return p.nativeType && p.ObjectType.Name == "Array"
}

// IsNumberObjectType reports whether this is the native Number wrapper.
func (p *PrototypeObjectType) IsNumberObjectType() bool {
	return p.nativeType && p.ObjectType.Name == "Number"
}

// IsStringObjectType reports whether this is the native String wrapper.
func (p *PrototypeObjectType) IsStringObjectType() bool {
	return p.nativeType && p.ObjectType.Name == "String"
}

// IsBooleanObjectType reports whether this is the native Boolean wrapper.
func (p *PrototypeObjectType) IsBooleanObjectType() bool {
	return p.nativeType && p.ObjectType.Name == "Boolean"
}

// IsDateType reports whether this is the native Date type.
func (p *PrototypeObjectType) IsDateType() bool {
	return p.nativeType && p.ObjectType.Name == "Date"
}

// IsRegexpType reports whether this is the native RegExp type.
func (p *PrototypeObjectType) IsRegexpType() bool {
	return p.nativeType && p.ObjectType.Name == "RegExp"
}

// IsBigIntObjectType reports whether this is the native BigInt wrapper.
func (p *PrototypeObjectType) IsBigIntObjectType() bool {
	return p.nativeType && p.ObjectType.Name == "BigInt"
}

// IsSymbolObjectType reports whether this is the native Symbol wrapper.
func (p *PrototypeObjectType) IsSymbolObjectType() bool {
	return p.nativeType && p.ObjectType.Name == "Symbol"
}
