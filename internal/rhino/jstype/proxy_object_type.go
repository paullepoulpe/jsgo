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

import "github.com/paullepoulpe/jsgo/internal/rhino"

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

// Go-only: NewProxyObjectType — Go constructor; Java uses package-private ProxyObjectType(registry, type).
// NewProxyObjectType constructs a ProxyObjectType wrapping referencedType.
// Mirrors the package-private Java constructor that takes only a registry
// and a JSType (no TemplateTypeMap).
func NewProxyObjectType(referencedType JSType) *ProxyObjectType {
	p := &ProxyObjectType{}
	p.SetReferencedType(referencedType)
	return p
}

// Go-only: NewProxyObjectTypeWithTemplates — Go constructor; Java uses ProxyObjectType(registry, type, templateTypeMap).
// NewProxyObjectTypeWithTemplates constructs a ProxyObjectType wrapping
// referencedType, with an associated TemplateTypeMap. The map is accepted
// for API parity with Java's two-arg constructor but is not stored on
// ProxyObjectType itself — subclasses that need it (e.g. TemplatizedType)
// track it separately.
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

// Go-only: TypeName — Go interface method; Java uses appendTo(TypeStringBuilder) instead.
// TypeName delegates to the referenced type. Mirrors Java's
// ProxyObjectType.appendTo which calls referencedType.appendTo.
func (p *ProxyObjectType) TypeName() string {
	return p.referencedType.TypeName()
}

// Go-only: IsSubtypeOf — Go JSType interface method; Java uses visit(RelationshipVisitor, JSType).
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

// GetImplicitPrototype mirrors ProxyObjectType.getImplicitPrototype() in Java —
// delegates to the referenced ObjectType's GetImplicitPrototype().
func (p *ProxyObjectType) GetImplicitPrototype() *ObjectType {
	if p.referencedObjType != nil {
		return p.referencedObjType.GetImplicitPrototype()
	}
	return nil
}

// Go-only: GetPropertyType — Go accessor; Java uses findPropertyTypeWithoutConsideringTemplateTypes.
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

// GetConstructor mirrors ProxyObjectType.getConstructor() in Java — delegates
// to the referenced ObjectType's GetConstructor().
func (p *ProxyObjectType) GetConstructor() *FunctionType {
	if p.referencedObjType != nil {
		return p.referencedObjType.GetConstructor()
	}
	return nil
}

// GetOwnerFunction mirrors ProxyObjectType.getOwnerFunction() in Java —
// delegates to the referenced ObjectType's GetOwnerFunction().
func (p *ProxyObjectType) GetOwnerFunction() *FunctionType {
	if p.referencedObjType != nil {
		return p.referencedObjType.GetOwnerFunction()
	}
	return nil
}

// GetCtorImplementedInterfaces mirrors ProxyObjectType.getCtorImplementedInterfaces()
// in Java — delegates to the referenced ObjectType.
func (p *ProxyObjectType) GetCtorImplementedInterfaces() []*ObjectType {
	if p.referencedObjType != nil {
		return p.referencedObjType.GetCtorImplementedInterfaces()
	}
	return nil
}

// GetCtorExtendedInterfaces mirrors ProxyObjectType.getCtorExtendedInterfaces()
// in Java — delegates to the referenced ObjectType.
func (p *ProxyObjectType) GetCtorExtendedInterfaces() []*ObjectType {
	if p.referencedObjType != nil {
		return p.referencedObjType.GetCtorExtendedInterfaces()
	}
	return nil
}

// GetTemplateTypes mirrors ProxyObjectType.getTemplateTypes() in Java —
// delegates to the referenced ObjectType, converting []*TemplateType to []JSType.
func (p *ProxyObjectType) GetTemplateTypes() []JSType {
	if p.referencedObjType == nil {
		return nil
	}
	params := p.referencedObjType.GetTypeParameters()
	if len(params) == 0 {
		return nil
	}
	out := make([]JSType, len(params))
	for i, tp := range params {
		out[i] = tp
	}
	return out
}

// GetTemplateParamCount mirrors JSType.getTemplateParamCount(). Delegates to
// the referenced ObjectType's type parameter count.
// Mirrors Java's ProxyObjectType.getTemplateParamCount() → referencedType.getTemplateParamCount().
func (p *ProxyObjectType) GetTemplateParamCount() int {
	if p.referencedObjType != nil {
		return len(p.referencedObjType.GetTypeParameters())
	}
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
// referencedType.isInterface().
func (p *ProxyObjectType) IsInterface() bool {
	if f, ok := p.referencedType.(*FunctionType); ok {
		return f.kind == FunctionKindInterface
	}
	return false
}

// IsNominalType mirrors JSType.isNominalType(). Delegates to the
// referenced type's Go representation — true for named ObjectTypes and
// constructor FunctionTypes.
// Mirrors Java's ProxyObjectType.isNominalType() → referencedType.isNominalType().
func (p *ProxyObjectType) IsNominalType() bool {
	if p.referencedObjType != nil && p.referencedObjType.Name != "" {
		return true
	}
	if _, ok := p.referencedType.(*FunctionType); ok {
		return true
	}
	return false
}

// IsInstanceType mirrors JSType.isInstanceType(). Returns true iff the
// referenced ObjectType was created as the instance type of a constructor.
// Mirrors Java's ProxyObjectType.isInstanceType() → referencedType.isInstanceType().
func (p *ProxyObjectType) IsInstanceType() bool {
	if p.referencedObjType != nil {
		return p.referencedObjType.constructor != nil
	}
	return false
}

// IsOrdinaryFunction mirrors JSType.isOrdinaryFunction(). Returns true iff
// the referenced type is a non-constructor, non-interface FunctionType.
// Mirrors Java's ProxyObjectType.isOrdinaryFunction() → referencedType.isOrdinaryFunction().
func (p *ProxyObjectType) IsOrdinaryFunction() bool {
	if f, ok := p.referencedType.(*FunctionType); ok {
		return f.kind == FunctionKindOrdinary
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
// does not distinguish checked-unknown from plain unknown; delegates to
// the referenced type's unknown-type check.
// Mirrors Java's ProxyObjectType.isCheckedUnknownType() → referencedType.isCheckedUnknownType().
func (p *ProxyObjectType) IsCheckedUnknownType() bool {
	return isUnknown(p.referencedType)
}

// IsNoType mirrors JSType.isNoType() — delegates to the referenced type.
// Mirrors Java's ProxyObjectType: final boolean isNoType() { return referencedType.isNoType(); }
func (p *ProxyObjectType) IsNoType() bool {
	if t, ok := p.referencedType.(*NoType); ok {
		return t.IsNoType()
	}
	return false
}

// IsNoObjectType mirrors JSType.isNoObjectType() — delegates to the referenced type.
// Mirrors Java's ProxyObjectType: final boolean isNoObjectType() { return referencedType.isNoObjectType(); }
func (p *ProxyObjectType) IsNoObjectType() bool {
	if t, ok := p.referencedType.(*NoObjectType); ok {
		return t.IsNoObjectType()
	}
	return false
}

// IsNoResolvedType mirrors JSType.isNoResolvedType() — delegates to the referenced type.
// Mirrors Java's ProxyObjectType: final boolean isNoResolvedType() { return referencedType.isNoResolvedType(); }
func (p *ProxyObjectType) IsNoResolvedType() bool {
	if t, ok := p.referencedType.(*NoResolvedType); ok {
		return t.IsNoResolvedType()
	}
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

// IsNativeObjectType mirrors ObjectType.isNativeObjectType() — delegates to
// the referenced ObjectType's native flag.
// Mirrors Java's ProxyObjectType: final boolean isNativeObjectType() {
//   return referencedObjType != null && referencedObjType.isNativeObjectType(); }
func (p *ProxyObjectType) IsNativeObjectType() bool {
	if p.referencedObjType != nil {
		return p.referencedObjType.IsNativeObjectType()
	}
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
// Mirrors Java's ProxyObjectType.hasAnyTemplateTypesInternal() → referencedType check.
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

// VisitReferenceType mirrors ProxyObjectType.visitReferenceType(Visitor<T>) —
// returns the referenced type for direct access without the proxy layer.
// In Go, since we don't have the Java Visitor interface, this returns the
// referenced JSType directly. Mirrors Java's visitReferenceType.
func (p *ProxyObjectType) VisitReferenceType() JSType {
	return p.referencedType
}

// --- Additional Java-parity methods ----------------------------------------

// GetJSDocInfo mirrors ProxyObjectType.getJSDocInfo() in Java — delegates to
// the referenced type's JSDocInfo.
func (p *ProxyObjectType) GetJSDocInfo() *rhino.JSDocInfo {
	if p.referencedObjType != nil {
		return p.referencedObjType.GetJSDocInfo()
	}
	return nil
}

// SetJSDocInfo mirrors ProxyObjectType.setJSDocInfo(JSDocInfo) in Java —
// delegates to the referenced ObjectType.
func (p *ProxyObjectType) SetJSDocInfo(info *rhino.JSDocInfo) {
	if p.referencedObjType != nil {
		p.referencedObjType.SetJSDocInfo(info)
	}
}

// GetTypeOfThis mirrors ProxyObjectType.getTypeOfThis() in Java — delegates
// to the referenced ObjectType.
func (p *ProxyObjectType) GetTypeOfThis() JSType {
	if p.referencedObjType != nil {
		return p.referencedObjType.GetTypeOfThis()
	}
	return nil
}

// CanBeCalled mirrors ProxyObjectType.canBeCalled() in Java — true if the
// referenced type is a FunctionType.
func (p *ProxyObjectType) CanBeCalled() bool {
	if _, ok := p.referencedType.(*FunctionType); ok {
		return true
	}
	return false
}
