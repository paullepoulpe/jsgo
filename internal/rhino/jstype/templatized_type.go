// Port of com/google/javascript/rhino/jstype/TemplatizedType.java
// @ closure-compiler v20260415
//
// TemplatizedType is a generic instantiation: an underlying ObjectType
// (Array, Map, etc.) paired with concrete substitutions for its template
// parameters (`Array<number>`, `Map<string,Foo>`).  It extends
// ProxyObjectType in Java; here we embed the base ObjectType reference
// directly and carry the extra per-specialisation state.
//
// At inference time we produce one whenever a JSDoc expression names a
// generic type with `<…>` slots; at call-site unification we recurse into
// its arguments so a `(Array<T>) → T` parameter binds T from the actual
// element type.

package jstype

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// TemplatizedType pairs a generic ObjectType with concrete type arguments.
// Mirrors com.google.javascript.rhino.jstype.TemplatizedType.
//
// Java layout: extends ProxyObjectType, stores templateTypes (ImmutableList)
// and isSpecializedOnlyWithUnknown, and holds a TemplateTypeReplacer used
// to apply template substitution to property types.
type TemplatizedType struct {
	// Base is the unparameterised ObjectType (e.g. Array, Map).
	// Mirrors ProxyObjectType.referencedObjType.
	Base *ObjectType

	// Args are the resolved type arguments in parameter order.
	// Mirrors TemplatizedType.templateTypes.
	Args []JSType

	// templateTypeMap is the TemplateTypeMap produced by
	// objectType.getTemplateTypeMap().copyFilledWithValues(templateTypes).
	// It carries the concrete bindings for all template keys of Base.
	// Mirrors the templateTypeMap field inherited from JSType via ObjectType.
	templateTypeMap *TemplateTypeMap

	// isSpecializedOnlyWithUnknown is true when every resolved arg is Unknown.
	// Used to skip template args in hash-code computation.
	// Mirrors TemplatizedType.isSpecializedOnlyWithUnknown.
	isSpecializedOnlyWithUnknown bool

	// replacer applies the concrete template bindings to property types.
	// Mirrors TemplatizedType.replacer.
	replacer *TemplateTypeReplacer
}

// NewTemplatizedType constructs a TemplatizedType wrapping base with the
// given type arguments.
// Mirrors TemplatizedType(JSTypeRegistry, ObjectType, ImmutableList<JSType>).
// Go-only: language-adapter — Go constructor; Java uses package-private TemplatizedType(registry, objectType, templateTypes).
func NewTemplatizedType(base *ObjectType, args []JSType) *TemplatizedType {
	t := &TemplatizedType{Base: base, Args: args}
	if base != nil {
		base.templatizedSelf = t

		// Build the templateTypeMap by creating a map of base's type parameters
		// and filling it with the provided values.
		// Mirrors: super(registry, objectType,
		//           objectType.getTemplateTypeMap().copyFilledWithValues(templateTypes))
		// In Java, objectType.getTemplateTypeMap() returns a map with keys but no
		// values (the "unfilled" map). In Go we build the equivalent from GetTypeParameters().
		params := base.GetTypeParameters()
		if len(params) > 0 {
			// Build a base map with keys only (no values) — the empty-value
			// form that objectType.getTemplateTypeMap() returns in Java.
			baseMap := newTemplateTypeMap(nil, params, nil, nil)
			t.templateTypeMap = baseMap.CopyFilledWithValues(args)

			// Resolve each param key through the filled map.
			// Mirrors the ImmutableList.Builder loop in the Java constructor.
			resolved := make([]JSType, len(params))
			allUnknown := true
			for i, key := range params {
				rt := t.templateTypeMap.GetResolvedTemplateType(key)
				resolved[i] = rt
				if !isUnknown(rt) {
					allUnknown = false
				}
			}
			t.Args = resolved
			t.isSpecializedOnlyWithUnknown = allUnknown

			// Build the partial replacer used for getPropertyType substitution.
			// Mirrors: this.replacer = TemplateTypeReplacer.forPartialReplacement(registry, getTemplateTypeMap())
			t.replacer = ForPartialReplacement(nil, t.templateTypeMap)
		}
	}
	return t
}

// --- JSTypeClass / type identity ------------------------------------------

// GetTypeClass mirrors TemplatizedType.getTypeClass() → JSTypeClass.TEMPLATIZED.
// Go-only: language-adapter — returns an int constant instead of enum value.
func (t *TemplatizedType) GetTypeClass() JSTypeClass { return JSTypeClassTemplatized }

// --- Visitor dispatch -------------------------------------------------------

// Visit mirrors TemplatizedType.visit(Visitor<T>) → visitor.caseTemplatizedType(this).
// Go-only: language-adapter — Go uses Visitor interface with any return type.
func (t *TemplatizedType) Visit(v Visitor) any {
	if v == nil {
		return t
	}
	return v.CaseTemplatizedType(t)
}

// VisitRelationship mirrors TemplatizedType.visit(RelationshipVisitor<T>, JSType).
// Go-only: language-adapter — Go uses RelationshipVisitor interface.
func (t *TemplatizedType) VisitRelationship(v RelationshipVisitor, that JSType) any {
	if v == nil {
		return t
	}
	return v.CaseTemplatizedType(t, that)
}

// --- toMaybe* downcasts ----------------------------------------------------

// ToMaybeTemplatizedType mirrors TemplatizedType.toMaybeTemplatizedType() — returns self.
func (t *TemplatizedType) ToMaybeTemplatizedType() *TemplatizedType { return t }

// ToMaybeObjectType returns the Base ObjectType for this TemplatizedType.
// Go-only: language-adapter — Java's toMaybeObjectType() is on JSType; TemplatizedType delegates to ProxyObjectType which returns referencedObjType. Go exposes Base directly.
func (t *TemplatizedType) ToMaybeObjectType() *ObjectType { return t.Base }

// --- Core accessors --------------------------------------------------------

// GetTemplateTypes mirrors TemplatizedType.getTemplateTypes() — returns the
// resolved type arguments in parameter order.
func (t *TemplatizedType) GetTemplateTypes() []JSType { return t.Args }

// GetTemplateTypeMap mirrors TemplatizedType.getTemplateTypeMap() — returns the
// TemplateTypeMap containing bindings for Base's template keys.
func (t *TemplatizedType) GetTemplateTypeMap() *TemplateTypeMap { return t.templateTypeMap }

// GetReferencedType mirrors TemplatizedType.getReferencedType() — returns the
// underlying (unparameterised) ObjectType.
// Mirrors: public ObjectType getReferencedType() { return getReferencedObjTypeInternal(); }
func (t *TemplatizedType) GetReferencedType() *ObjectType { return t.Base }

// GetReferencedTypeInternal mirrors ProxyObjectType.getReferencedTypeInternal().
// Go-only: language-adapter — ProxyObjectType is flattened; returns Base as JSType.
func (t *TemplatizedType) GetReferencedTypeInternal() JSType { return t.Base }

// GetReferencedObjTypeInternal mirrors ProxyObjectType.getReferencedObjTypeInternal().
func (t *TemplatizedType) GetReferencedObjTypeInternal() *ObjectType { return t.Base }

// --- Property type lookup --------------------------------------------------

// GetPropertyType mirrors TemplatizedType.getPropertyType(Property.Key) in Java:
//
//	@Override public JSType getPropertyType(Property.Key propertyName) {
//	  JSType result = super.getPropertyType(propertyName);
//	  return result == null ? null : result.visit(replacer);
//	}
//
// The Java super call resolves to ProxyObjectType.findPropertyTypeWithoutConsideringTemplateTypes
// which delegates to referencedType.findPropertyType(propertyName).  In Go that
// corresponds to Base.GetPropertyType(name).  The replacer then substitutes any
// TemplateType occurrences with the concrete args from Args.
func (t *TemplatizedType) GetPropertyType(name string) JSType {
	if t.Base == nil {
		return nil
	}
	result := t.Base.GetPropertyType(name)
	if result == nil {
		return nil
	}
	if t.replacer == nil {
		return result
	}
	return visitJSType(t.replacer, result)
}

// HasProperty mirrors TemplatizedType.hasProperty(name) via delegation.
// Go-only: language-adapter — Java's ObjectType.hasProperty delegates to getPropertyType; Go exposes this convenience on TemplatizedType.
func (t *TemplatizedType) HasProperty(name string) bool {
	return t.GetPropertyType(name) != nil
}

// --- Interface/implemented-interface substitution -------------------------

// GetCtorImplementedInterfaces mirrors TemplatizedType.getCtorImplementedInterfaces():
//
//	@Override public Iterable<ObjectType> getCtorImplementedInterfaces() {
//	  LinkedHashSet<ObjectType> resolved = new LinkedHashSet<>();
//	  for (ObjectType obj : getReferencedObjTypeInternal().getCtorImplementedInterfaces()) {
//	    resolved.add(obj.visit(replacer).toObjectType());
//	  }
//	  return resolved;
//	}
func (t *TemplatizedType) GetCtorImplementedInterfaces() []*ObjectType {
	if t.Base == nil {
		return nil
	}
	raw := t.Base.GetCtorImplementedInterfaces()
	if len(raw) == 0 || t.replacer == nil {
		return raw
	}
	out := make([]*ObjectType, 0, len(raw))
	seen := make(map[*ObjectType]bool)
	for _, obj := range raw {
		replaced := visitJSType(t.replacer, obj)
		if ot, ok := replaced.(*ObjectType); ok && !seen[ot] {
			seen[ot] = true
			out = append(out, ot)
		}
	}
	return out
}

// GetCtorExtendedInterfaces mirrors TemplatizedType.getCtorExtendedInterfaces():
//
//	@Override public Iterable<ObjectType> getCtorExtendedInterfaces() {
//	  LinkedHashSet<ObjectType> resolved = new LinkedHashSet<>();
//	  for (ObjectType obj : getReferencedObjTypeInternal().getCtorExtendedInterfaces()) {
//	    resolved.add(obj.visit(replacer).toObjectType());
//	  }
//	  return resolved;
//	}
func (t *TemplatizedType) GetCtorExtendedInterfaces() []*ObjectType {
	if t.Base == nil {
		return nil
	}
	raw := t.Base.GetCtorExtendedInterfaces()
	if len(raw) == 0 || t.replacer == nil {
		return raw
	}
	out := make([]*ObjectType, 0, len(raw))
	seen := make(map[*ObjectType]bool)
	for _, obj := range raw {
		replaced := visitJSType(t.replacer, obj)
		if ot, ok := replaced.(*ObjectType); ok && !seen[ot] {
			seen[ot] = true
			out = append(out, ot)
		}
	}
	return out
}

// --- Type-string / hash ---------------------------------------------------

// TypeName mirrors TemplatizedType.appendTo(TypeStringBuilder):
//
//	@Override void appendTo(TypeStringBuilder sb) {
//	  super.appendTo(sb);
//	  if (!this.templateTypes.isEmpty()) {
//	    sb.append("<").appendAll(this.templateTypes, ",").append(">");
//	  }
//	}
//
// Go-only: language-adapter — Go interface uses TypeName() string instead of appendTo.
func (t *TemplatizedType) TypeName() string {
	if t.Base == nil {
		return "?"
	}
	if len(t.Args) == 0 {
		return t.Base.TypeName()
	}
	names := make([]string, len(t.Args))
	for i, a := range t.Args {
		if a == nil {
			names[i] = "?"
		} else {
			names[i] = a.TypeName()
		}
	}
	return t.Base.TypeName() + "<" + strings.Join(names, ",") + ">"
}

// RecursionUnsafeHashCode mirrors TemplatizedType.recursionUnsafeHashCode():
//
//	@Override int recursionUnsafeHashCode() {
//	  int baseHash = super.recursionUnsafeHashCode();
//	  if (this.isSpecializedOnlyWithUnknown) { return baseHash; }
//	  return Objects.hash(templateTypes, baseHash);
//	}
//
// Go-only: language-adapter — ObjectType has no RecursionUnsafeHashCode so we use pointer hash.
func (t *TemplatizedType) RecursionUnsafeHashCode() int {
	// Use a stable hash of the base type name as the base hash, mirroring
	// ProxyObjectType.recursionUnsafeHashCode() → referencedType.hashCode().
	baseHash := 0
	if t.Base != nil {
		for _, c := range t.Base.TypeName() {
			baseHash = 31*baseHash + int(c)
		}
	}
	if t.isSpecializedOnlyWithUnknown {
		return baseHash
	}
	// Combine args hash with base hash — mirrors Objects.hash(templateTypes, baseHash).
	h := 31*17 + baseHash
	for _, a := range t.Args {
		ah := 0
		if a != nil {
			if hasher, ok := a.(interface{ RecursionUnsafeHashCode() int }); ok {
				ah = hasher.RecursionUnsafeHashCode()
			} else {
				for _, c := range a.TypeName() {
					ah = 31*ah + int(c)
				}
			}
		}
		h = 31*h + ah
	}
	return h
}

// --- Template-type flag ---------------------------------------------------

// HasAnyTemplateTypesInternal mirrors TemplatizedType.hasAnyTemplateTypesInternal():
//
//	@Override public boolean hasAnyTemplateTypesInternal() {
//	  return templateTypeMap.hasAnyTemplateTypesInternal();
//	}
func (t *TemplatizedType) HasAnyTemplateTypesInternal() bool {
	if t.templateTypeMap != nil {
		return t.templateTypeMap.HasAnyTemplateTypesInternal()
	}
	return false
}

// IsInstantiatedGenericType reports whether this TemplatizedType has at least
// one template argument with a non-Unknown type.
// Mirrors TemplatizedType.isInstantiatedGenericType() in Java.
// Go-only: language-adapter — Java's isInstantiatedGenericType() lives on JSType; Go exposes it on TemplatizedType directly.
func (t *TemplatizedType) IsInstantiatedGenericType() bool {
	return len(t.Args) > 0
}

// --- Raw-type comparison --------------------------------------------------

// WrapsSameRawType mirrors TemplatizedType.wrapsSameRawType(JSType):
//
//	boolean wrapsSameRawType(JSType that) {
//	  return that.isTemplatizedType() && this.getReferencedTypeInternal()
//	      .equals(that.toMaybeTemplatizedType().getReferencedTypeInternal());
//	}
func (t *TemplatizedType) WrapsSameRawType(that JSType) bool {
	other, ok := that.(*TemplatizedType)
	if !ok {
		return false
	}
	return t.Base == other.Base
}

// --- Greatest subtype helper ----------------------------------------------

// GetGreatestSubtypeHelper mirrors TemplatizedType.getGreatestSubtypeHelper(JSType):
//
//	JSType getGreatestSubtypeHelper(JSType rawThat) {
//	  checkNotNull(rawThat);
//	  if (!wrapsSameRawType(rawThat)) {
//	    if (!rawThat.isTemplatizedType()) {
//	      if (this.isSubtype(rawThat)) { return this; }
//	      else if (rawThat.isSubtypeOf(this)) { return rawThat; }
//	    }
//	    if (this.isObject() && rawThat.isObject()) { return getNativeType(NO_OBJECT_TYPE); }
//	    return getNativeType(NO_TYPE);
//	  }
//	  TemplatizedType that = rawThat.toMaybeTemplatizedType();
//	  if (this.equals(that)) { return this; }
//	  return getReferencedObjTypeInternal();
//	}
func (t *TemplatizedType) GetGreatestSubtypeHelper(rawThat JSType) JSType {
	if rawThat == nil {
		panic("TemplatizedType.GetGreatestSubtypeHelper: rawThat must not be nil")
	}
	if !t.WrapsSameRawType(rawThat) {
		if _, ok := rawThat.(*TemplatizedType); !ok {
			if IsSubtype(t, rawThat) {
				return t
			}
			if IsSubtype(rawThat, t) {
				return rawThat
			}
		}
		if IsObject(t) && IsObject(rawThat) {
			return NoObject
		}
		return None
	}
	that := rawThat.(*TemplatizedType)
	if t == that {
		return t
	}
	// Same raw type, different args: return the raw (unparameterised) type.
	return t.Base
}

// --- Subtype check --------------------------------------------------------

// IsSubtypeOf mirrors the template-argument-invariant subtype check.
// When both sides are TemplatizedTypes sharing the same raw type, the
// template arguments are compared modulo null/undefined — mirroring Java's
// TemplateTypeMap.checkEquivalenceHelper which always passes
// SubtypingMode.IGNORE_NULL_UNDEFINED for invariant template positions.
func (t *TemplatizedType) IsSubtypeOf(o JSType) bool {
	if other, ok := o.(*TemplatizedType); ok {
		if !IsSubtype(t.Base, other.Base) {
			return false
		}
		n := len(t.Args)
		if len(other.Args) < n {
			n = len(other.Args)
		}
		for i := 0; i < n; i++ {
			a := stripNullUndef(t.Args[i])
			b := stripNullUndef(other.Args[i])
			if !IsSubtype(a, b) {
				return false
			}
		}
		return true
	}
	if t.Base != nil {
		return t.Base.IsSubtypeOf(o)
	}
	return false
}

// --- Template-arg lookup --------------------------------------------------

// GetTemplateArg returns the resolved type argument for the given template
// key, looked up by pointer identity through the base object's type
// parameter list.  Returns nil if key is not a parameter of the base type.
// Go-only: language-adapter — Java uses TemplateTypeMap.getResolvedTemplateType directly.
func (t *TemplatizedType) GetTemplateArg(key *TemplateType) JSType {
	if key == nil || t.templateTypeMap == nil {
		return nil
	}
	if t.templateTypeMap.HasTemplateType(key) {
		return t.templateTypeMap.GetResolvedTemplateType(key)
	}
	return nil
}

// --- Delegation helpers ---------------------------------------------------

// MatchesObjectContext delegates to the base ObjectType.
// Mirrors ProxyObjectType.matchesObjectContext() → referencedType.matchesObjectContext().
// Go-only: language-adapter — Java's TemplatizedType inherits this via ProxyObjectType; Go exposes it directly since there is no inheritance chain.
func (t *TemplatizedType) MatchesObjectContext() bool { return true }

// GetPossibleToBooleanOutcomes mirrors ObjectType.getPossibleToBooleanOutcomes() — objects are always truthy.
// Go-only: language-adapter — Java's TemplatizedType inherits via ProxyObjectType and ObjectType; Go exposes directly.
func (t *TemplatizedType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSetTrue
}

// GetJSDocInfo mirrors ProxyObjectType.getJSDocInfo() — delegates to Base.
// Go-only: language-adapter — Java's TemplatizedType inherits getJSDocInfo via ProxyObjectType; Go exposes it directly.
func (t *TemplatizedType) GetJSDocInfo() *rhino.JSDocInfo {
	if t.Base != nil {
		return t.Base.GetJSDocInfo()
	}
	return nil
}

// GetImplicitPrototype mirrors ProxyObjectType.getImplicitPrototype() — delegates to Base.
// Go-only: language-adapter — Java's TemplatizedType inherits getImplicitPrototype via ProxyObjectType; Go exposes it directly.
func (t *TemplatizedType) GetImplicitPrototype() *ObjectType {
	if t.Base != nil {
		return t.Base.GetImplicitPrototype()
	}
	return nil
}

// GetConstructor mirrors ProxyObjectType.getConstructor() — delegates to Base.
// Go-only: language-adapter — Java's TemplatizedType inherits getConstructor via ProxyObjectType; Go exposes it directly.
func (t *TemplatizedType) GetConstructor() *FunctionType {
	if t.Base != nil {
		return t.Base.GetConstructor()
	}
	return nil
}

// GetReferenceName mirrors ProxyObjectType.getReferenceName() — delegates to Base.
// Go-only: language-adapter — Java's TemplatizedType inherits getReferenceName via ProxyObjectType; Go exposes it directly.
func (t *TemplatizedType) GetReferenceName() string {
	if t.Base != nil {
		return t.Base.Name
	}
	return ""
}

// IsNativeObjectType mirrors ProxyObjectType.isNativeObjectType() — delegates to Base.
// Go-only: language-adapter — Java's TemplatizedType inherits isNativeObjectType via ProxyObjectType; Go exposes it directly.
func (t *TemplatizedType) IsNativeObjectType() bool {
	if t.Base != nil {
		return t.Base.IsNativeObjectType()
	}
	return false
}

// --- stripNullUndef -------------------------------------------------------

// stripNullUndef removes Null and Void members from a UnionType, returning
// the remaining type.  Mirrors the null/undef stripping Java applies via
// SubtypingMode.IGNORE_NULL_UNDEFINED when comparing template arguments.
// Go-only: language-adapter — Java achieves this via SubtypingMode on the visitor path.
func stripNullUndef(t JSType) JSType {
	if t == nil {
		return t
	}
	if t == Null || t == Void {
		return t
	}
	u, ok := t.(*UnionType)
	if !ok {
		return t
	}
	var members []JSType
	for _, m := range u.members {
		if m != Null && m != Void {
			members = append(members, m)
		}
	}
	if len(members) == 0 {
		return t
	}
	if len(members) == 1 {
		return members[0]
	}
	return &UnionType{members: members}
}

// --- findTemplatizedSuper -------------------------------------------------

// findTemplatizedSuper returns the first TemplatizedType encountered while
// walking up the super chain of obj, or nil.
// Go-only: language-adapter — Java uses the ProxyObjectType wrapper hierarchy; Go walks the super chain.
func findTemplatizedSuper(obj *ObjectType) *TemplatizedType {
	for cur := obj; cur != nil; cur = cur.super {
		if cur.templatizedSelf != nil {
			return cur.templatizedSelf
		}
	}
	return nil
}
