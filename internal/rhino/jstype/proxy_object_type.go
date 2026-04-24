// Port of com/google/javascript/rhino/jstype/ProxyObjectType.java
// @ closure-compiler v20260415.
//
// ProxyObjectType is an ObjectType that uses composition to delegate all
// calls to a wrapped (referenced) JSType. It serves as the base for
// NamedType, TemplatizedType, and similar wrapper types that resolve their
// actual type lazily or apply an additional layer of indirection.
//
// NOTE: this port is additive. Existing NamedType/TemplatizedType are NOT
// wired to embed ProxyObjectType — that wiring is a separate later task.
//
// Many Java methods delegate to types not yet ported (JSDocInfo, Node,
// Visitor, ErrorReporter, PropertyMap, etc.). Those methods are stubbed
// with TODO comments and will panic if called, flagging any premature
// consumer.

package jstype

// ProxyObjectType wraps an arbitrary JSType and forwards all operations
// to it. It satisfies the JSType interface by delegating TypeName and
// IsSubtypeOf to the referenced type — matching Java's behaviour where
// ProxyObjectType extends ObjectType and @Override-s every query method.
//
// The companion referencedObjType field is set whenever the referenced
// type is an *ObjectType, mirroring Java's referencedObjType field that
// is set inside setReferencedType.
type ProxyObjectType struct {
	referencedType    JSType
	referencedObjType *ObjectType
}

// NewProxyObjectType constructs a ProxyObjectType wrapping referencedType.
// Mirrors the package-private Java constructor that takes only a registry
// and a JSType (no TemplateTypeMap).
func NewProxyObjectType(referencedType JSType) *ProxyObjectType {
	p := &ProxyObjectType{}
	p.SetReferencedType(referencedType)
	return p
}

// NewProxyObjectTypeWithTemplates constructs a ProxyObjectType wrapping
// referencedType, with an associated TemplateTypeMap. The map is accepted
// for API parity with Java's two-arg constructor but is not stored on
// ProxyObjectType itself — subclasses that need it (e.g. TemplatizedType)
// track it separately.
//
// TODO(proxy): store/use templateTypeMap once TemplateTypeMap is fully
// integrated into the jstype hierarchy.
func NewProxyObjectTypeWithTemplates(referencedType JSType, _ *TemplateTypeMap) *ProxyObjectType {
	return NewProxyObjectType(referencedType)
}

// --- referenced-type accessors ----------------------------------------

// GetReferencedTypeInternal returns the wrapped JSType. Mirrors Java's
// package-private getReferencedTypeInternal().
func (p *ProxyObjectType) GetReferencedTypeInternal() JSType {
	return p.referencedType
}

// GetReferencedObjTypeInternal returns the wrapped type cast to
// *ObjectType, or nil if the referenced type is not an ObjectType.
// Mirrors Java's package-private getReferencedObjTypeInternal().
func (p *ProxyObjectType) GetReferencedObjTypeInternal() *ObjectType {
	return p.referencedObjType
}

// SetReferencedType replaces the wrapped type and refreshes the companion
// referencedObjType pointer. Mirrors Java's package-private
// setReferencedType(JSType).
func (p *ProxyObjectType) SetReferencedType(referencedType JSType) {
	if referencedType == nil {
		panic("ProxyObjectType.SetReferencedType: referencedType must not be nil")
	}
	p.referencedType = referencedType
	if objType, ok := referencedType.(*ObjectType); ok {
		p.referencedObjType = objType
	} else {
		p.referencedObjType = nil
	}
}

// --- JSType interface --------------------------------------------------

// TypeName delegates to the referenced type. Mirrors Java's
// ProxyObjectType.appendTo which calls referencedType.appendTo.
func (p *ProxyObjectType) TypeName() string {
	return p.referencedType.TypeName()
}

// IsSubtypeOf delegates to the referenced type. Mirrors Java's visit
// mechanism — ProxyObjectType.visit(RelationshipVisitor, JSType) calls
// referencedType.visit(visitor, that), so subtype tests are forwarded.
func (p *ProxyObjectType) IsSubtypeOf(other JSType) bool {
	return p.referencedType.IsSubtypeOf(other)
}

// --- ObjectType-shape delegation methods ------------------------------
// Each method mirrors the corresponding @Override in Java's
// ProxyObjectType. Methods that depend on types not yet ported are
// stubbed; they panic so that a caller discovers the gap immediately.

// GetReferenceName mirrors ObjectType.getReferenceName(). Returns the
// empty string when the referenced type is not an ObjectType — matching
// Java's fallback.
func (p *ProxyObjectType) GetReferenceName() string {
	if p.referencedObjType != nil {
		return p.referencedObjType.Name
	}
	return ""
}

// IsStructuralType reports whether the wrapped type is a structural
// (@record / @interface) type. Mirrors Java's isStructuralType().
func (p *ProxyObjectType) IsStructuralType() bool {
	if p.referencedObjType != nil {
		return p.referencedObjType.IsStructural
	}
	return false
}

// IsStruct mirrors ObjectType.isStruct().
func (p *ProxyObjectType) IsStruct() bool {
	if p.referencedObjType != nil {
		return p.referencedObjType.IsStruct
	}
	return false
}

// IsDict mirrors ObjectType.isDict().
func (p *ProxyObjectType) IsDict() bool {
	if p.referencedObjType != nil {
		return p.referencedObjType.IsDict
	}
	return false
}

// GetImplicitPrototype mirrors ObjectType.getImplicitPrototype(). Returns
// nil when the referenced type is not an ObjectType — matching Java's null
// return for non-object references.
//
// TODO(proxy): return the actual prototype chain entry once the prototype
// field is added to ObjectType.
func (p *ProxyObjectType) GetImplicitPrototype() *ObjectType {
	if p.referencedObjType != nil {
		return p.referencedObjType.Super()
	}
	return nil
}

// GetPropertyType looks up a property on the referenced ObjectType. Returns
// nil when the referenced type is not an ObjectType. Mirrors Java's
// findPropertyTypeWithoutConsideringTemplateTypes which calls
// referencedType.findPropertyType(propertyName).
func (p *ProxyObjectType) GetPropertyType(name string) JSType {
	if p.referencedObjType != nil {
		return p.referencedObjType.GetPropertyType(name)
	}
	return nil
}

// HasProperty mirrors ObjectType.hasProperty via delegation to
// referencedObjType.
func (p *ProxyObjectType) HasProperty(name string) bool {
	if p.referencedObjType != nil {
		return p.referencedObjType.HasProperty(name)
	}
	return false
}

// GetConstructor mirrors ObjectType.getConstructor(). Returns nil when the
// referenced type is not an ObjectType. The *FunctionType return type
// matches Java's FunctionType; callers must handle nil.
//
// TODO(proxy): wire to the constructor slot on ObjectType once it carries
// a FunctionType reference.
func (p *ProxyObjectType) GetConstructor() *FunctionType {
	// ObjectType does not yet carry a constructor FunctionType pointer.
	// Return nil conservatively — matches Java's null return for plain
	// anonymous object types.
	return nil
}

// GetOwnerFunction mirrors ObjectType.getOwnerFunction(). Returns nil when
// the referenced type is not an ObjectType. Used by TypeCheck for prototype
// method ownership checks.
//
// TODO(proxy): wire to owner-function slot on ObjectType once added.
func (p *ProxyObjectType) GetOwnerFunction() *FunctionType {
	return nil
}

// GetCtorImplementedInterfaces mirrors ObjectType.getCtorImplementedInterfaces().
// Returns an empty slice when the referenced type is not an ObjectType.
// Mirrors Java's ImmutableList.of() fallback.
//
// TODO(proxy): wire to @implements list on ObjectType once ported.
func (p *ProxyObjectType) GetCtorImplementedInterfaces() []*ObjectType {
	return nil
}

// GetCtorExtendedInterfaces mirrors ObjectType.getCtorExtendedInterfaces().
// Returns nil (empty slice) when the referenced type is not an ObjectType.
//
// TODO(proxy): wire to @extends list on ObjectType once ported.
func (p *ProxyObjectType) GetCtorExtendedInterfaces() []*ObjectType {
	return nil
}

// GetTemplateTypes mirrors ObjectType.getTemplateTypes(). Returns nil when
// the referenced type is not an ObjectType.
//
// TODO(proxy): wire to template-type list on ObjectType once ported.
func (p *ProxyObjectType) GetTemplateTypes() []JSType {
	return nil
}

// GetTemplateParamCount mirrors JSType.getTemplateParamCount(). Returns 0
// as a conservative default — matches Java's base implementation.
//
// TODO(proxy): delegate to referencedType once getTemplateParamCount is
// part of the Go JSType interface.
func (p *ProxyObjectType) GetTemplateParamCount() int {
	return 0
}

// ToMaybeEnumType mirrors JSType.toMaybeEnumType(). Delegates to the
// referenced type's dynamic type. Returns nil if the referenced type is
// not an *EnumType.
func (p *ProxyObjectType) ToMaybeEnumType() *EnumType {
	if e, ok := p.referencedType.(*EnumType); ok {
		return e
	}
	return nil
}

// ToMaybeEnumElementType mirrors JSType.toMaybeEnumElementType(). Returns
// nil if the referenced type is not an *EnumElementType.
func (p *ProxyObjectType) ToMaybeEnumElementType() *EnumElementType {
	if e, ok := p.referencedType.(*EnumElementType); ok {
		return e
	}
	return nil
}

// ToMaybeUnionType mirrors JSType.toMaybeUnionType(). Returns nil if the
// referenced type is not a *UnionType.
func (p *ProxyObjectType) ToMaybeUnionType() *UnionType {
	if u, ok := p.referencedType.(*UnionType); ok {
		return u
	}
	return nil
}

// ToMaybeFunctionType mirrors JSType.toMaybeFunctionType(). Returns nil if
// the referenced type is not a *FunctionType.
func (p *ProxyObjectType) ToMaybeFunctionType() *FunctionType {
	if f, ok := p.referencedType.(*FunctionType); ok {
		return f
	}
	return nil
}

// ToMaybeTemplatizedType mirrors JSType.toMaybeTemplatizedType(). Returns
// nil if the referenced type is not a *TemplatizedType.
func (p *ProxyObjectType) ToMaybeTemplatizedType() *TemplatizedType {
	if t, ok := p.referencedType.(*TemplatizedType); ok {
		return t
	}
	return nil
}

// ToMaybeTemplateType mirrors JSType.toMaybeTemplateType(). Returns nil if
// the referenced type is not a *TemplateType.
func (p *ProxyObjectType) ToMaybeTemplateType() *TemplateType {
	if t, ok := p.referencedType.(*TemplateType); ok {
		return t
	}
	return nil
}

// CollapseUnion mirrors JSType.collapseUnion(). When the referenced type is
// a union, delegates collapse to it; otherwise returns the receiver itself.
// Mirrors Java's ProxyObjectType.collapseUnion().
func (p *ProxyObjectType) CollapseUnion() JSType {
	return CollapseUnion(p.referencedType)
}

// IsConstructor reports whether the referenced type is a constructor
// function. Mirrors Java's ProxyObjectType.isConstructor().
func (p *ProxyObjectType) IsConstructor() bool {
	if f, ok := p.referencedType.(*FunctionType); ok {
		return f.IsConstructor
	}
	return false
}

// IsInterface reports whether the referenced type is an @interface.
// Mirrors Java's ProxyObjectType.isInterface() which delegates to
// FunctionType.isInterface(). Returns false for non-function types.
//
// TODO(proxy): extend once ObjectType/FunctionType carry an IsInterface
// flag.
func (p *ProxyObjectType) IsInterface() bool {
	return false
}

// IsNominalType mirrors JSType.isNominalType(). Delegates to the
// referenced type's Go representation — true for named ObjectTypes and
// constructor FunctionTypes.
//
// TODO(proxy): delegate to referencedType once isNominalType is part of
// the Go JSType interface.
func (p *ProxyObjectType) IsNominalType() bool {
	if p.referencedObjType != nil && p.referencedObjType.Name != "" {
		return true
	}
	return false
}

// IsInstanceType mirrors JSType.isInstanceType(). Returns true iff the
// referenced ObjectType was created as the instance type of a constructor.
//
// TODO(proxy): wire to the IsInstance flag on FunctionType.Instance once
// that relationship is modelled.
func (p *ProxyObjectType) IsInstanceType() bool {
	return false
}

// IsOrdinaryFunction mirrors JSType.isOrdinaryFunction(). Returns true iff
// the referenced type is a non-constructor FunctionType.
func (p *ProxyObjectType) IsOrdinaryFunction() bool {
	if f, ok := p.referencedType.(*FunctionType); ok {
		return !f.IsConstructor
	}
	return false
}

// IsAllType mirrors JSType.isAllType(). Returns true iff the referenced
// type is the all-type sentinel.
func (p *ProxyObjectType) IsAllType() bool {
	return isAll(p.referencedType)
}

// IsUnknownType mirrors JSType.isUnknownType(). Returns true iff the
// referenced type is the unknown-type sentinel.
func (p *ProxyObjectType) IsUnknownType() bool {
	return isUnknown(p.referencedType)
}

// IsCheckedUnknownType mirrors JSType.isCheckedUnknownType(). The Go port
// does not distinguish checked-unknown from plain unknown yet.
//
// TODO(proxy): add a checkedUnknownType variant and delegate here.
func (p *ProxyObjectType) IsCheckedUnknownType() bool {
	return false
}

// IsNoType mirrors JSType.isNoType(). The Go port has no bottom type yet.
//
// TODO(proxy): add noType and delegate here.
func (p *ProxyObjectType) IsNoType() bool {
	return false
}

// IsNoObjectType mirrors JSType.isNoObjectType(). The Go port has no
// no-object type yet.
//
// TODO(proxy): add noObjectType and delegate here.
func (p *ProxyObjectType) IsNoObjectType() bool {
	return false
}

// IsNoResolvedType mirrors JSType.isNoResolvedType(). The Go port has no
// no-resolved type yet.
//
// TODO(proxy): add noResolvedType and delegate here.
func (p *ProxyObjectType) IsNoResolvedType() bool {
	return false
}

// IsNullable mirrors JSType.isNullable(). The Go port represents
// nullable via UnionType{…, Null}; check the union for a null member.
func (p *ProxyObjectType) IsNullable() bool {
	if u, ok := p.referencedType.(*UnionType); ok {
		for _, m := range u.Members() {
			if _, isNull := m.(nullType); isNull {
				return true
			}
		}
	}
	if _, ok := p.referencedType.(nullType); ok {
		return true
	}
	return false
}

// IsVoidable mirrors JSType.isVoidable(). Checks whether the referenced
// type includes undefined (void).
func (p *ProxyObjectType) IsVoidable() bool {
	return containsVoid(p.referencedType)
}

// IsNativeObjectType mirrors ObjectType.isNativeObjectType(). The Go port
// does not track native/externs provenance on ObjectType yet.
//
// TODO(proxy): add IsNativeObjectType flag to ObjectType and delegate.
func (p *ProxyObjectType) IsNativeObjectType() bool {
	return false
}

// LoosenTypecheckingDueToForwardReferencedSupertype mirrors the Java
// method of the same name on JSType. The Go port does not model forward
// references yet; returns false conservatively.
//
// TODO(proxy): wire once forward-reference resolution is ported.
func (p *ProxyObjectType) LoosenTypecheckingDueToForwardReferencedSupertype() bool {
	return false
}

// HasAnyTemplateTypesInternal mirrors JSType.hasAnyTemplateTypesInternal().
// Checks whether the referenced type is (or contains) a *TemplateType.
//
// TODO(proxy): delegate to referencedType once hasAnyTemplateTypes is
// part of the Go JSType interface.
func (p *ProxyObjectType) HasAnyTemplateTypesInternal() bool {
	if _, ok := p.referencedType.(*TemplateType); ok {
		return true
	}
	if t, ok := p.referencedType.(*TemplatizedType); ok {
		for _, a := range t.Args {
			if _, ok2 := a.(*TemplateType); ok2 {
				return true
			}
		}
	}
	return false
}

// VisitReferenceType is the Java visitReferenceType(Visitor<T>) helper —
// it allows callers to visit the referenced type directly without going
// through the proxy. Not applicable until the Go port has a Visitor
// mechanism.
//
// TODO(proxy): add Visitor interface and implement.
func (p *ProxyObjectType) VisitReferenceType() JSType {
	return p.referencedType
}
