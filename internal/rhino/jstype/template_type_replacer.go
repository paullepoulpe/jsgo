// Port of com/google/javascript/rhino/jstype/TemplateTypeReplacer.java
// @ closure-compiler v20260415
//
// TemplateTypeReplacer specializes TemplatizedTypes according to a
// TemplateTypeMap of bindings.  It corresponds to the Java Visitor<JSType>
// that is dispatched via JSType.visit(this).
//
// In Go, the Visitor pattern is implemented as a type-switch in the
// package-level visitJSType function (replacing Java's JSType.visit),
// since the JSType interface does not carry a visit method.
//
// Subset ported: the three factory functions (forInference, forTotalReplacement,
// forPartialReplacement), all caseXxx methods, guardAgainstCycles, and
// setKeyType.  The caseObjectTypeUnguarded property-walking branch (which
// requires visitProperties=true, only set by forInference) is ported but
// ObjectType's GetOwnPropertyNames / GetPropertyNode APIs are stubs —
// in the current Go port, record-type property replacement falls through
// to a no-op, matching the narrow harness scope.
//
// TODO(object-property-replacer): implement record-type property walking
// when ObjectType.ownPropertyNames is fully ported (Tier-4+).
// TODO(proxy-object-type): caseProxyObjectType is a no-op stub until
// ProxyObjectType is ported (Tier-4+).

package jstype

// TemplateTypeReplacer specializes TemplatizedTypes according to a
// TemplateTypeMap of bindings.
// Mirrors com.google.javascript.rhino.jstype.TemplateTypeReplacer.
type TemplateTypeReplacer struct {
	registry *Registry
	bindings *TemplateTypeMap

	visitProperties          bool
	useUnknownForMissingKeys bool
	// useUnknownForMissingValues controls whether an unbound TemplateType
	// is replaced with Unknown (true) or left as-is (false).
	useUnknownForMissingValues bool

	hasMadeReplacement bool
	keyType            *TemplateType
	// ownSubMapBoundary is the exclusive upper index in bindings up to which
	// the current keyType's submap extends.  -1 means "no limit".
	ownSubMapBoundary int

	// seenTypes guards against cycles; nil until first use.
	seenTypes map[JSType]bool
	// visitedObjectTypes caches results for record-type property walks;
	// nil until first use.
	visitedObjectTypes map[JSType]JSType
}

// forInference creates a TemplateTypeReplacer for use during TypeInference.
// Mirrors TemplateTypeReplacer.forInference in Java.
func forInference(registry *Registry, bindings map[*TemplateType]JSType) *TemplateTypeReplacer {
	keys := make([]*TemplateType, 0, len(bindings))
	for k := range bindings {
		keys = append(keys, k)
	}
	values := make([]JSType, len(keys))
	for i, k := range keys {
		v := bindings[k]
		if v == nil {
			v = Unknown
		}
		values[i] = v
	}
	m := &TemplateTypeMap{}
	if registry != nil {
		m = createEmpty(registry)
	}
	m = m.CopyWithExtension(keys, values)
	return newTemplateTypeReplacer(registry, m, true, true, true)
}

// ForInference is the exported factory for external callers.
// Mirrors TemplateTypeReplacer.forInference in Java.
func ForInference(registry *Registry, bindings map[*TemplateType]JSType) *TemplateTypeReplacer {
	return forInference(registry, bindings)
}

// forTotalReplacement creates a replacer that always eliminates
// TemplateTypes; missing bindings resolve to Unknown.
// Mirrors TemplateTypeReplacer.forTotalReplacement in Java.
func forTotalReplacement(registry *Registry, bindings *TemplateTypeMap) *TemplateTypeReplacer {
	return newTemplateTypeReplacer(registry, bindings, false, false, true)
}

// ForTotalReplacement is the exported factory.
// Mirrors TemplateTypeReplacer.forTotalReplacement in Java.
func ForTotalReplacement(registry *Registry, bindings *TemplateTypeMap) *TemplateTypeReplacer {
	return forTotalReplacement(registry, bindings)
}

// ForPartialReplacement creates a replacer that may leave TemplateTypes
// in place when no binding is available.
// Mirrors TemplateTypeReplacer.forPartialReplacement in Java.
func ForPartialReplacement(registry *Registry, bindings *TemplateTypeMap) *TemplateTypeReplacer {
	return newTemplateTypeReplacer(registry, bindings, false, false, false)
}

// newTemplateTypeReplacer is the private constructor.
func newTemplateTypeReplacer(
	registry *Registry,
	bindings *TemplateTypeMap,
	visitProperties bool,
	useUnknownForMissingKeys bool,
	useUnknownForMissingValues bool,
) *TemplateTypeReplacer {
	return &TemplateTypeReplacer{
		registry:                   registry,
		bindings:                   bindings,
		visitProperties:            visitProperties,
		useUnknownForMissingKeys:   useUnknownForMissingKeys,
		useUnknownForMissingValues: useUnknownForMissingValues,
		ownSubMapBoundary:          -1,
	}
}

// HasMadeReplacement reports whether any substitution was performed.
// Mirrors TemplateTypeReplacer.hasMadeReplacement in Java.
func (r *TemplateTypeReplacer) HasMadeReplacement() bool { return r.hasMadeReplacement }

// setKeyType sets the current key being processed and its submap boundary.
// Called from newTemplateTypeMap during resolution.
// Mirrors TemplateTypeReplacer.setKeyType in Java.
func (r *TemplateTypeReplacer) setKeyType(keyType *TemplateType, ownSubMapBoundary int) {
	r.keyType = keyType
	r.ownSubMapBoundary = ownSubMapBoundary
}

// initSeenTypes lazily initialises the cycle-detection set.
func (r *TemplateTypeReplacer) initSeenTypes() {
	if r.seenTypes == nil {
		r.seenTypes = make(map[JSType]bool)
	}
}

// visitJSType dispatches t through the TemplateTypeReplacer, implementing
// the Java Visitor<JSType> pattern via type switch.
// Corresponds to JSType.visit(replacer) in Java.
func visitJSType(r *TemplateTypeReplacer, t JSType) JSType {
	if t == nil {
		return nil
	}
	switch v := t.(type) {
	case *TemplateType:
		return r.caseTemplateType(v)
	case *UnionType:
		return r.caseUnionType(v)
	case *TemplatizedType:
		return r.caseTemplatizedType(v)
	case *FunctionType:
		return r.caseFunctionType(v)
	case *ObjectType:
		return r.caseObjectType(v)
	case *EnumType:
		return r.caseEnumElementTypeOrEnum(t)
	case *EnumElementType:
		return r.caseEnumElementTypeOrEnum(v)
	case unknownType:
		return Unknown
	case allType:
		return All
	case numberType:
		return Number
	case stringType:
		return String
	case booleanType:
		return Boolean
	case nullType:
		return Null
	case voidType:
		return Void
	default:
		// Unknown / unrecognised type: return as-is.
		return t
	}
}

// caseNoType corresponds to Visitor.caseNoType.  NoType is not ported
// yet; return t unchanged.
func (r *TemplateTypeReplacer) caseNoType(t JSType) JSType { return t }

// caseEnumElementTypeOrEnum corresponds to Visitor.caseEnumElementType.
// EnumElementType and EnumType are returned unchanged.
func (r *TemplateTypeReplacer) caseEnumElementTypeOrEnum(t JSType) JSType { return t }

// caseFunctionType guards against cycles and delegates to
// caseFunctionTypeUnguarded.
// Mirrors TemplateTypeReplacer.caseFunctionType in Java.
func (r *TemplateTypeReplacer) caseFunctionType(t *FunctionType) JSType {
	return r.guardAgainstCycles(t, func(u JSType) JSType {
		return r.caseFunctionTypeUnguarded(u.(*FunctionType))
	})
}

// caseFunctionTypeUnguarded performs the actual function-type replacement.
// Mirrors TemplateTypeReplacer.caseFunctionTypeUnguarded in Java.
func (r *TemplateTypeReplacer) caseFunctionTypeUnguarded(t *FunctionType) JSType {
	// Native function types are left alone.
	if t.IsConstructor && t.Instance == nil {
		return t
	}
	if !t.IsConstructor && len(t.Params) == 0 && t.Return == nil {
		return t
	}

	changed := false

	beforeReturn := t.Return
	afterReturn := visitJSType(r, beforeReturn)
	if afterReturn != beforeReturn {
		changed = true
	}

	paramsChanged := false
	afterParams := make([]JSType, len(t.Params))
	for i, p := range t.Params {
		afterP := visitJSType(r, p)
		afterParams[i] = afterP
		if afterP != p {
			changed = true
			paramsChanged = true
		}
	}

	if !changed {
		return t
	}

	params := t.Params
	if paramsChanged {
		params = afterParams
	}
	if afterReturn == nil {
		afterReturn = Unknown
	}
	return &FunctionType{
		Params:        params,
		Return:        afterReturn,
		MinArity:      t.MinArity,
		MaxArity:      t.MaxArity,
		IsConstructor: t.IsConstructor,
		IsES6Class:    t.IsES6Class,
		Instance:      t.Instance,
		TemplateTypes: t.TemplateTypes,
	}
}

// caseObjectType guards against cycles and delegates to
// caseObjectTypeUnguarded.
// Mirrors TemplateTypeReplacer.caseObjectType in Java.
func (r *TemplateTypeReplacer) caseObjectType(t *ObjectType) JSType {
	return r.guardAgainstCycles(t, func(u JSType) JSType {
		return r.caseObjectTypeUnguarded(u.(*ObjectType))
	})
}

// caseObjectTypeUnguarded handles record-type property walking when
// visitProperties is true.  ObjectType's own property APIs are not
// fully ported; this returns the input unchanged for non-record types.
// Mirrors TemplateTypeReplacer.caseObjectTypeUnguarded in Java.
func (r *TemplateTypeReplacer) caseObjectTypeUnguarded(t *ObjectType) JSType {
	if !r.visitProperties || t.IsStructural == false || t.Name != "" {
		// Nominal type or visitProperties disabled — return as-is.
		// TODO(object-property-replacer): walk record properties.
		return t
	}
	// For structural (record) types with visitProperties=true, check cache.
	if r.visitedObjectTypes != nil {
		if cached, ok := r.visitedObjectTypes[t]; ok {
			return cached
		}
	} else {
		r.visitedObjectTypes = make(map[JSType]JSType)
	}

	changed := false
	newObj := NewObjectType("")
	newObj.IsStructural = true
	for name, beforeType := range t.members {
		afterType := visitJSType(r, beforeType)
		if afterType != beforeType {
			changed = true
		}
		newObj.DefineProperty(name, afterType)
	}

	var result JSType = t
	if changed {
		result = newObj
	}
	r.visitedObjectTypes[t] = result
	return result
}

// caseTemplatizedType guards against cycles and delegates.
// Mirrors TemplateTypeReplacer.caseTemplatizedType in Java.
func (r *TemplateTypeReplacer) caseTemplatizedType(t *TemplatizedType) JSType {
	return r.guardAgainstCycles(t, func(u JSType) JSType {
		return r.caseTemplatizedTypeUnguarded(u.(*TemplatizedType))
	})
}

// caseTemplatizedTypeUnguarded replaces the base and type arguments.
// Mirrors TemplateTypeReplacer.caseTemplatizedTypeUnguarded in Java.
func (r *TemplateTypeReplacer) caseTemplatizedTypeUnguarded(t *TemplatizedType) JSType {
	changed := false

	beforeBase := t.Base
	afterBaseJSType := visitJSType(r, beforeBase)
	afterBase, _ := afterBaseJSType.(*ObjectType)
	if afterBase == nil {
		afterBase = beforeBase
	}
	if afterBase != beforeBase {
		changed = true
	}

	afterArgs := make([]JSType, len(t.Args))
	for i, before := range t.Args {
		after := visitJSType(r, before)
		afterArgs[i] = after
		if after != before {
			changed = true
		}
	}

	if !changed {
		return t
	}
	return NewTemplatizedType(afterBase, afterArgs)
}

// caseUnionType guards against cycles and delegates.
// Mirrors TemplateTypeReplacer.caseUnionType in Java.
func (r *TemplateTypeReplacer) caseUnionType(t *UnionType) JSType {
	return r.guardAgainstCycles(t, func(u JSType) JSType {
		return r.caseUnionTypeUnguarded(u.(*UnionType))
	})
}

// caseUnionTypeUnguarded replaces each union member.
// Mirrors TemplateTypeReplacer.caseUnionTypeUnguarded in Java.
func (r *TemplateTypeReplacer) caseUnionTypeUnguarded(t *UnionType) JSType {
	changed := false
	results := make([]JSType, len(t.members))
	for i, alt := range t.members {
		repl := visitJSType(r, alt)
		results[i] = repl
		if repl != alt {
			changed = true
		}
	}
	if !changed {
		return t
	}
	return NewUnion(results...)
}

// caseTemplateType handles TemplateType replacement.
// Mirrors TemplateTypeReplacer.caseTemplateType in Java.
func (r *TemplateTypeReplacer) caseTemplateType(t *TemplateType) JSType {
	r.hasMadeReplacement = true

	if t != r.keyType && !r.bindings.hasTemplateKeyFrom(t, r.ownSubMapBoundary) {
		if r.useUnknownForMissingKeys {
			return Unknown
		}
		return t
	}

	r.initSeenTypes()
	if r.seenTypes[t] {
		// Cycle detected — return the TemplateType itself.
		return t
	}
	if !r.bindings.HasTemplateType(t) {
		if r.useUnknownForMissingValues {
			return Unknown
		}
		return t
	}

	replacement := r.bindings.GetUnresolvedOriginalTemplateType(t)

	// Guard against recursive templatized types such as T => Foo<T>.
	if replacement == r.keyType || r.seenTypes[replacement] {
		return t
	}

	r.seenTypes[t] = true
	visited := visitJSType(r, replacement)
	delete(r.seenTypes, t)

	if visited == r.keyType {
		panic("TemplateTypeReplacer: trying to replace key with the same value")
	}
	return visited
}

// caseNamedType returns the named type unchanged — its internals are not
// interesting to template replacement.
// Mirrors TemplateTypeReplacer.caseNamedType in Java.
// TODO(named-type): add *NamedType case to visitJSType when named_type.go
// is present in this package (it exists in the main repo but not this
// worktree snapshot).

// caseProxyObjectType is a stub; ProxyObjectType is not yet ported.
// Mirrors TemplateTypeReplacer.caseProxyObjectType in Java (no-op for now).
// TODO(proxy-object-type): unwrap and re-wrap when ProxyObjectType is ported.
func (r *TemplateTypeReplacer) caseProxyObjectType(t JSType) JSType { return t }

// guardAgainstCycles adds t to seenTypes, calls mapper(t), then removes t.
// Returns t immediately if t is already in seenTypes.
// Mirrors TemplateTypeReplacer.guardAgainstCycles in Java.
func (r *TemplateTypeReplacer) guardAgainstCycles(t JSType, mapper func(JSType) JSType) JSType {
	r.initSeenTypes()
	if r.seenTypes[t] {
		return t
	}
	r.seenTypes[t] = true
	result := mapper(t)
	delete(r.seenTypes, t)
	return result
}
