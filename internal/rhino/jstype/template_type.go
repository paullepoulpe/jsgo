// Port of com/google/javascript/rhino/jstype/TemplateType.java
// @ closure-compiler v20260415 (narrow; no bounds, no late-binding).
//
// The Java TemplateType carries a bound (like Java generics'
// `extends T`), interned-per-registry identity, and a full
// TemplateTypeMap for resolving substitutions during visitCall. This
// port tracks just the name — enough for call-site unification to
// thread arg types through @param {T} → @return {T}. Bounds,
// `@template T extends Base` syntax, and the implicit-template
// handling live in the Tier-4+ slice (see TODO.md P0-4).
//
// The full Java-faithful TemplateTypeMap is in template_type_map.go;
// this file retains CallSiteMap, the narrow call-site substitution
// struct used by TypeInference.

package jstype

// TemplateType represents an unresolved generic parameter (e.g. `T`
// in `/** @template T */`). Two TemplateTypes compare equal only by
// pointer — the registry interns one per `@template` name per
// function, mirroring Java's TemplateType reference semantics.
type TemplateType struct {
	Name string
}

// NewTemplateType returns a fresh TemplateType carrying the named
// `@template` parameter. Not interned: the caller (TypedScopeCreator)
// is expected to create one per function and reuse it across that
// function's @param / @return slots.
func NewTemplateType(name string) *TemplateType {
	return &TemplateType{Name: name}
}

func (t *TemplateType) TypeName() string { return t.Name }

// IsSubtypeOf is permissive — a TemplateType is treated like Unknown
// so unresolved generics neither false-positive nor false-negative.
// Once unification refines the type at a call site, the refined
// JSType (not the TemplateType) participates in the subtype check.
func (t *TemplateType) IsSubtypeOf(o JSType) bool { return true }

// CallSiteMap is the substitution `TemplateType → JSType` built
// at a CALL site. Narrow port of the call-site unification portion
// of rhino/jstype/TemplateTypeMap.java; used by TypeInference.
// The full Java-faithful TemplateTypeMap is in template_type_map.go.
type CallSiteMap struct {
	subs map[*TemplateType]JSType
}

func NewCallSiteMap() *CallSiteMap {
	return &CallSiteMap{subs: map[*TemplateType]JSType{}}
}

// Bind records tmpl → t. If a prior binding exists and the new type
// differs, the two are unioned — matches Java's "pick the join" rule
// when the same template parameter is filled from multiple argument
// positions.
func (m *CallSiteMap) Bind(tmpl *TemplateType, t JSType) {
	if tmpl == nil || t == nil {
		return
	}
	if prev, ok := m.subs[tmpl]; ok {
		if prev.TypeName() == t.TypeName() {
			return
		}
		m.subs[tmpl] = NewUnion(prev, t)
		return
	}
	m.subs[tmpl] = t
}

// Lookup returns the substitution for tmpl (or nil if the template is
// unresolved, e.g. because no argument positions mention it).
func (m *CallSiteMap) Lookup(tmpl *TemplateType) JSType {
	if m == nil {
		return nil
	}
	return m.subs[tmpl]
}

// Substitute walks `t` replacing every TemplateType with its mapped
// concrete type. ObjectType / FunctionType / UnionType / Templatized
// are handled structurally. Unmapped template vars stay as
// TemplateType so a downstream caller can still decide how to
// degrade them (our convention is to treat them as Unknown).
func (m *CallSiteMap) Substitute(t JSType) JSType {
	if m == nil || t == nil {
		return t
	}
	switch v := t.(type) {
	case *TemplateType:
		if sub := m.Lookup(v); sub != nil {
			return sub
		}
		return t
	case *UnionType:
		var out []JSType
		changed := false
		for _, m2 := range v.Members() {
			s := m.Substitute(m2)
			out = append(out, s)
			if s != m2 {
				changed = true
			}
		}
		if !changed {
			return t
		}
		return NewUnion(out...)
	case *TemplatizedType:
		var args []JSType
		changed := false
		for _, a := range v.Args {
			s := m.Substitute(a)
			args = append(args, s)
			if s != a {
				changed = true
			}
		}
		if !changed {
			return t
		}
		return NewTemplatizedType(v.Base, args)
	case *FunctionType:
		var params []JSType
		changed := false
		for _, p := range v.Params {
			s := m.Substitute(p)
			params = append(params, s)
			if s != p {
				changed = true
			}
		}
		ret := m.Substitute(v.Return)
		if ret != v.Return {
			changed = true
		}
		if !changed {
			return t
		}
		out := &FunctionType{
			Params:        params,
			Return:        ret,
			IsConstructor: v.IsConstructor,
			Instance:      v.Instance,
			TemplateTypes: v.TemplateTypes,
		}
		return out
	}
	return t
}
