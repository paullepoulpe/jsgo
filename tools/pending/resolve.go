package main

// Mapping unresolved Go references back to the Java members they stand for
// (by PORTING.md §3, read backwards), and keeping the set of pending
// declarations.

import (
	"fmt"
	"go/types"
	"sort"
	"strings"
	"unicode"
)

// MemberKind is the kind of a pending declaration.
type MemberKind int

const (
	MCtor MemberKind = iota
	MStaticFunc
	MStaticVar
	MEnumConst
	MMethod
	MField
	MGoIfaceMethod // a method required by a ported Go interface the class implements
	MAsAdapter     // a Go-only As<Interface>() adapter returning the class as one of its interfaces
)

// PMember is one pending declaration of a class.
type PMember struct {
	kind   MemberKind
	goName string
	method *methodRef
	field  *fieldRef
	// For MGoIfaceMethod: the method's Go signature and the interface.
	goSig   *types.Signature
	goIface string
	// For MAsAdapter: the Java interface (internal name) it returns.
	asType string
}

// PClass is a pending Java class.
type PClass struct {
	c       *ClassInfo
	goName  string
	kind    Kind
	members map[string]*PMember // Go name → member (methods and fields share Go's selector namespace)
	statics map[string]*PMember // package-level names
	// typeNeeded: the Go type itself is referenced (a static-only class
	// needs no type).
	typeNeeded bool
	// funcSig is set for a functional interface, declared as a Go func type.
	funcSig *methodRef
}

// Gen holds the pending declarations being generated.
type Gen struct {
	model   *Model
	pkg     *goPackage
	res     *checkResult
	classes map[string]*PClass
	warned  map[string]bool
	warns   []string
	// otherByGoName: for the Go packages below jscomp, "rhino.Node" →
	// internal Java name; built lazily.
	otherByGoName map[string]string
}

func newGen(model *Model, pkg *goPackage) *Gen {
	return &Gen{model: model, pkg: pkg, classes: map[string]*PClass{}, warned: map[string]bool{}}
}

func (g *Gen) warn(msg string) {
	if !g.warned[msg] {
		g.warned[msg] = true
		g.warns = append(g.warns, msg)
	}
}

// realObject returns the package-level object of the real (non-generated)
// Go code named name, if any.
func (g *Gen) realObject(name string) types.Object {
	if g.res == nil || g.res.pkg == nil {
		return nil
	}
	obj := g.res.pkg.Scope().Lookup(name)
	if obj == nil {
		return nil
	}
	if g.res.pending != nil {
		pf := g.pkg.fset.File(g.res.pending.Pos())
		if pf != nil && pf.Base() <= int(obj.Pos()) && int(obj.Pos()) <= pf.Base()+pf.Size() {
			return nil
		}
	}
	return obj
}

// classPorted reports whether a Java class has a real Go port: its Java
// file is named by a "// Port of" header, or its Go type exists.
func (g *Gen) classPorted(c *ClassInfo) bool {
	if g.pkg.portedJava[g.model.javaFile(c)] {
		return true
	}
	name := g.model.goNames[c.Name]
	if obj := g.realObject(name); obj != nil {
		if _, ok := obj.(*types.TypeName); ok {
			return true
		}
	}
	return false
}

// ensureClass declares the pending class c (if not ported).
func (g *Gen) ensureClass(c *ClassInfo) *PClass {
	name := g.model.goNames[c.Name]
	if p, ok := g.classes[name]; ok {
		return p
	}
	if name == "" || g.classPorted(c) {
		return nil
	}
	p := &PClass{c: c, goName: name, kind: g.model.kind(c), members: map[string]*PMember{}, statics: map[string]*PMember{}}
	g.classes[name] = p
	if p.kind == KindIface {
		// An interface's type is its method set: declare its abstract
		// methods (and those of the supertypes that are not ported Go
		// interfaces, which are embedded instead).
		abstract := g.abstractMethods(c)
		if c.Access&accInterface != 0 && len(abstract) == 1 && len(g.realIfaceSupers(c)) == 0 {
			p.funcSig = &abstract[0]
		} else {
			for i := range abstract {
				ref := abstract[i]
				goName := capitalize(ref.m.Name)
				if _, has := p.members[goName]; !has {
					p.members[goName] = &PMember{kind: MMethod, goName: goName, method: &ref}
				}
			}
		}
	}
	return p
}

// typeOf is ensureClass for a class whose Go type is referenced.
func (g *Gen) typeOf(c *ClassInfo) *PClass {
	p := g.ensureClass(c)
	if p != nil {
		p.typeNeeded = true
	}
	return p
}

// abstractMethods lists the abstract methods of c and of its supertypes
// that are not ported Go interfaces (one per name and descriptor).
func (g *Gen) abstractMethods(c *ClassInfo) []methodRef {
	var out []methodRef
	seen := map[string]bool{}
	add := func(cls *ClassInfo) {
		for i, mi := range cls.Methods {
			if mi.Access&accAbstract == 0 || mi.Access&(accStatic|accSynthetic|accBridge) != 0 {
				continue
			}
			if seen[mi.Name+mi.Desc] {
				continue
			}
			seen[mi.Name+mi.Desc] = true
			out = append(out, methodRef{owner: cls, m: mi, index: i})
		}
	}
	add(c)
	for _, s := range g.model.supertypes(c) {
		if iface, _ := g.realInterface(s); iface != nil {
			continue
		}
		if sc, ok := g.model.classes[s]; ok {
			add(sc)
		}
	}
	return out
}

// realIfaceSupers lists the supertypes of c that are ported Go interfaces.
func (g *Gen) realIfaceSupers(c *ClassInfo) []string {
	var out []string
	for _, s := range g.model.supertypes(c) {
		if iface, _ := g.realInterface(s); iface != nil {
			out = append(out, s)
		}
	}
	return out
}

func hasSuffixBoundary(rest string) bool {
	if rest == "" {
		return true
	}
	r := []rune(rest)[0]
	return unicode.IsUpper(r) || unicode.IsDigit(r) || r == '_'
}

// resolveTop resolves an unresolved package-level Go name.
func (g *Gen) resolveTop(n *need) bool {
	name := n.name
	site := callSite{}
	if len(n.sites) > 0 {
		site = n.sites[0]
	}
	// A type?
	if c, ok := g.model.byGoName[name]; ok && (site.typePos || !site.called) {
		if g.classPorted(c) {
			g.warn(fmt.Sprintf("%s: %s is ported but not defined under that name", site.pos, name))
			return false
		}
		return g.typeOf(c) != nil
	}
	// A constructor?
	for _, prefix := range []string{"New", "new"} {
		if !strings.HasPrefix(name, prefix) {
			continue
		}
		rest := name[len(prefix):]
		best := ""
		for gname := range g.model.byGoName {
			if strings.HasPrefix(rest, gname) && hasSuffixBoundary(rest[len(gname):]) && len(gname) > len(best) {
				best = gname
			}
		}
		if best == "" {
			continue
		}
		c := g.model.byGoName[best]
		var ctors []methodRef
		for i, mi := range c.Methods {
			if mi.Name == "<init>" && mi.Access&accSynthetic == 0 {
				ctors = append(ctors, methodRef{owner: c, m: mi, index: i})
			}
		}
		if len(ctors) == 0 {
			continue
		}
		if g.classPorted(c) {
			g.warn(fmt.Sprintf("%s: %s: %s is ported but has no constructor named so", site.pos, name, best))
			return false
		}
		p := g.typeOf(c)
		if p == nil {
			return false
		}
		ref := g.chooseOverload(ctors, rest[len(best):], n.sites)
		p.statics[name] = &PMember{kind: MCtor, goName: name, method: &ref}
		return true
	}
	// A static member: <Class>_<member>, the longest class name first.
	var prefixes []string
	for gname := range g.model.byGoName {
		if strings.HasPrefix(name, gname+"_") {
			prefixes = append(prefixes, gname)
		}
	}
	sort.Slice(prefixes, func(i, j int) bool {
		if len(prefixes[i]) != len(prefixes[j]) {
			return len(prefixes[i]) > len(prefixes[j])
		}
		return prefixes[i] < prefixes[j]
	})
	for _, best := range prefixes {
		c := g.model.byGoName[best]
		rest := name[len(best)+1:]
		if rest == "" {
			continue
		}
		if f := g.model.fieldNamed(c, rest, true); f != nil {
			if g.classPorted(c) {
				g.warn(fmt.Sprintf("%s: %s: %s is ported but has no member named so", site.pos, name, best))
				return false
			}
			p := g.ensureClass(c)
			if p == nil {
				return false
			}
			kind := MStaticVar
			if f.f.Access&accEnum != 0 {
				kind = MEnumConst
			}
			p.statics[name] = &PMember{kind: kind, goName: name, field: f}
			return true
		}
		// Static methods; `values` keeps its Java spelling for enums
		// (Foo_values()); a trailing _ marks a static method whose plain
		// name is taken by a nested type (Foo_Builder_()).
		mname, suffix := "", ""
		for _, mi := range c.Methods {
			if mi.Access&accStatic == 0 || mi.Access&accSynthetic != 0 || strings.HasPrefix(mi.Name, "<") || strings.Contains(mi.Name, "$") {
				continue
			}
			goBase := capitalize(mi.Name)
			if c.Access&accEnum != 0 && mi.Name == "values" {
				goBase = "values"
			}
			for _, cand := range []string{goBase + "_", goBase} {
				if strings.HasPrefix(rest, cand) && (cand == rest || hasSuffixBoundary(rest[len(cand):])) && len(cand) > len(mname) {
					mname, suffix = cand, rest[len(cand):]
				}
			}
		}
		if mname == "" {
			continue
		}
		jname := strings.TrimSuffix(mname, "_")
		if jname != "values" {
			jname = decapitalize(jname)
		}
		cands := g.model.methodsNamed(c, jname, true)
		if len(cands) == 0 {
			continue
		}
		if g.classPorted(c) {
			g.warn(fmt.Sprintf("%s: %s: %s is ported but has no member named so", site.pos, name, best))
			return false
		}
		p := g.ensureClass(c)
		if p == nil {
			return false
		}
		ref := g.chooseOverload(cands, suffix, n.sites)
		ref.listResultArrayList = resultSelected(n.sites)
		p.statics[name] = &PMember{kind: MStaticFunc, goName: name, method: &ref}
		return true
	}
	if c, ok := g.model.byGoName[name]; ok {
		if g.classPorted(c) {
			g.warn(fmt.Sprintf("%s: %s is ported but not defined under that name", site.pos, name))
			return false
		}
		return g.typeOf(c) != nil
	}
	g.warn(fmt.Sprintf("%s: undefined %s matches no Java class or member", site.pos, name))
	return false
}

// resolveMember resolves a selector on a pending type.
func (g *Gen) resolveMember(n *need) bool {
	p := g.classes[n.owner]
	if p == nil {
		return false
	}
	site := callSite{}
	if len(n.sites) > 0 {
		site = n.sites[0]
	}
	c := p.c
	// Java Object/Enum methods under their Go names.
	special := map[string]string{"String": "toString", "Ordinal": "ordinal", "Equals": "equals", "HashCode": "hashCode", "CompareTo": "compareTo", "Name": "name"}
	if p.kind == KindEnum && n.name == "String" {
		special["String"] = "name"
	}
	if jn, ok := special[n.name]; ok {
		cands := g.model.methodsNamed(c, jn, false)
		if len(cands) == 0 {
			// Inherited from java.lang.Object/Enum: not in the jar index.
			p.members[n.name] = &PMember{kind: MMethod, goName: n.name, method: &methodRef{owner: c, m: &MethodInfo{Name: jn, Desc: specialDesc(jn)}, index: 1 << 20}}
			return true
		}
		ref := g.chooseOverload(cands, "", n.sites)
		p.members[n.name] = &PMember{kind: MMethod, goName: n.name, method: &ref}
		return true
	}
	if !site.called {
		if f := g.model.fieldNamed(c, n.name, false); f != nil {
			p.members[n.name] = &PMember{kind: MField, goName: n.name, field: f}
			return true
		}
	}
	// A Go-only As<Interface>() adapter (PORTING.md: a Java class whose
	// methods for two interfaces collide in Go returns itself as one of them
	// through such a method).
	if want, ok := strings.CutPrefix(n.name, "As"); ok && site.called && want != "" {
		for _, s := range g.model.supertypes(c) {
			sc := g.model.classes[s]
			if sc == nil || sc.Access&accInterface == 0 {
				continue
			}
			if goName := g.model.goNames[s]; goName == want || strings.HasSuffix(goName, "_"+want) {
				p.members[n.name] = &PMember{kind: MAsAdapter, goName: n.name, asType: s}
				return true
			}
		}
	}
	// Methods of the class and its supertypes.
	mname := ""
	visit := func(cls *ClassInfo) {
		for _, mi := range cls.Methods {
			if mi.Access&(accStatic|accSynthetic|accBridge) != 0 || strings.HasPrefix(mi.Name, "<") || strings.Contains(mi.Name, "$") {
				continue
			}
			base := capitalize(mi.Name)
			if strings.HasPrefix(n.name, base) && hasSuffixBoundary(n.name[len(base):]) && len(base) > len(mname) {
				mname = base
			}
		}
	}
	visit(c)
	for _, s := range g.model.supertypes(c) {
		if sc, ok := g.model.classes[s]; ok {
			visit(sc)
		}
	}
	if mname == "" {
		if f := g.model.fieldNamed(c, n.name, false); f != nil {
			p.members[n.name] = &PMember{kind: MField, goName: n.name, field: f}
			return true
		}
		g.warn(fmt.Sprintf("%s: %s.%s matches no Java member of %s", site.pos, n.owner, n.name, c.Name))
		return false
	}
	cands := g.model.methodsNamed(c, decapitalize(mname), false)
	if len(cands) == 0 {
		g.warn(fmt.Sprintf("%s: %s.%s: no overloads", site.pos, n.owner, n.name))
		return false
	}
	ref := g.chooseOverload(cands, n.name[len(mname):], n.sites)
	ref.listResultArrayList = resultSelected(n.sites)
	p.members[n.name] = &PMember{kind: MMethod, goName: n.name, method: &ref}
	return true
}

// resultSelected reports whether a call site selects a field or method on
// the call's result.
func resultSelected(sites []callSite) bool {
	for _, s := range sites {
		if s.called && s.resultSelected {
			return true
		}
	}
	return false
}

func specialDesc(name string) string {
	switch name {
	case "toString", "name":
		return "()Ljava/lang/String;"
	case "ordinal", "hashCode":
		return "()I"
	case "equals":
		return "(Ljava/lang/Object;)Z"
	case "compareTo":
		return "(Ljava/lang/Object;)I"
	}
	return "()V"
}

// chooseOverload picks the Java overload a Go call site stands for: by
// arity, then by argument types, then (PORTING.md §3) the plain name is
// the overload with the fewest parameters.
func (g *Gen) chooseOverload(cands []methodRef, suffix string, sites []callSite) methodRef {
	if len(cands) == 1 {
		return cands[0]
	}
	var site *callSite
	for i := range sites {
		if sites[i].called {
			site = &sites[i]
			break
		}
	}
	nparams := func(r methodRef) int { return len(methodSig(r.m).Params) }
	filtered := cands
	if site != nil {
		var byArity []methodRef
		for _, c := range cands {
			n := nparams(c)
			varargs := c.m.Access&accVarargs != 0
			if n == site.arity || (varargs && !site.ellipsis && site.arity >= n-1) {
				byArity = append(byArity, c)
			}
		}
		if len(byArity) > 0 {
			filtered = byArity
		}
		if len(filtered) > 1 {
			type scored struct {
				r     methodRef
				score int
			}
			var ok []scored
			for _, c := range filtered {
				sig := methodSig(c.m)
				score, compatible := 0, true
				for i, at := range site.argTypes {
					if at == nil || i >= len(sig.Params) {
						continue
					}
					pt := sig.Params[i]
					if c.m.Access&accVarargs != 0 && i >= len(sig.Params)-1 && pt.Kind == JArray && !site.ellipsis {
						pt = pt.Elem
					}
					switch g.compatible(at, pt) {
					case 1:
						score++
					case -1:
						compatible = false
					}
				}
				if compatible {
					ok = append(ok, scored{c, score})
				}
			}
			if len(ok) > 0 {
				sort.SliceStable(ok, func(i, j int) bool { return ok[i].score > ok[j].score })
				var best []methodRef
				for _, s := range ok {
					if s.score == ok[0].score {
						best = append(best, s.r)
					}
				}
				filtered = best
			}
		}
	}
	sort.SliceStable(filtered, func(i, j int) bool { return nparams(filtered[i]) < nparams(filtered[j]) })
	if suffix != "" && len(filtered) > 1 && nparams(filtered[0]) < nparams(filtered[1]) {
		// A suffixed name is not the fewest-parameter overload.
		return filtered[1]
	}
	return filtered[0]
}

// compatible compares a Go argument type with a Java parameter type: 1 if
// they correspond, -1 if they cannot, 0 if unknown.
func (g *Gen) compatible(at types.Type, pt *JType) int {
	switch pt.Kind {
	case JPrim:
		b, ok := at.Underlying().(*types.Basic)
		if !ok {
			return -1
		}
		switch pt.Prim {
		case 'Z':
			if b.Info()&types.IsBoolean != 0 {
				return 1
			}
		case 'D', 'F':
			if b.Info()&types.IsNumeric != 0 {
				return 1
			}
		default:
			if b.Info()&types.IsInteger != 0 {
				return 1
			}
		}
		return -1
	case JArray:
		if _, ok := at.Underlying().(*types.Slice); ok {
			return 1
		}
		return -1
	case JTypeVar, JWildcard:
		return 0
	}
	switch pt.Class {
	case "java/lang/String", "java/lang/CharSequence":
		if b, ok := at.Underlying().(*types.Basic); ok && b.Info()&types.IsString != 0 {
			return 1
		}
		if p, ok := at.(*types.Pointer); ok {
			if b, ok := p.Elem().Underlying().(*types.Basic); ok && b.Info()&types.IsString != 0 {
				return 1
			}
		}
		return -1
	case "java/lang/Object":
		return 0
	case "java/lang/Integer", "java/lang/Long", "java/lang/Short", "java/lang/Byte":
		if b, ok := at.Underlying().(*types.Basic); ok && b.Info()&types.IsInteger != 0 {
			return 1
		}
		return 0
	case "java/lang/Boolean":
		if b, ok := at.Underlying().(*types.Basic); ok && b.Info()&types.IsBoolean != 0 {
			return 1
		}
		return 0
	}
	// A Go slice stands for the Java collections that map to []T (emit.go).
	if _, isSlice := at.Underlying().(*types.Slice); isSlice {
		switch pt.Class {
		case "java/util/Collection", "java/util/List", "java/util/Queue", "java/util/Deque",
			"com/google/common/collect/ImmutableList", "com/google/common/collect/ImmutableCollection",
			"java/util/stream/Stream", "com/google/common/collect/FluentIterable", "java/util/SortedSet",
			"com/google/common/collect/ImmutableSortedSet":
			return 1
		}
		return -1
	}
	// A class: compare with the Java class of the Go type.
	jc := g.javaClassOfGoType(at)
	if jc == "" {
		if _, isBasic := at.Underlying().(*types.Basic); isBasic {
			return -1
		}
		return 0
	}
	if jc == pt.Class {
		return 1
	}
	if c, ok := g.model.classes[jc]; ok {
		for _, s := range g.model.supertypes(c) {
			if s == pt.Class {
				return 1
			}
		}
		if _, known := g.model.classes[pt.Class]; known {
			return -1
		}
	}
	return 0
}

// javaClassOfGoType maps a named Go type (or pointer to one) to its Java
// class.
func (g *Gen) javaClassOfGoType(t types.Type) string {
	if p, ok := t.(*types.Pointer); ok {
		t = p.Elem()
	}
	named, ok := t.(*types.Named)
	if !ok {
		return ""
	}
	obj := named.Obj()
	if obj.Pkg() == nil {
		return ""
	}
	if g.res != nil && obj.Pkg() == g.res.pkg {
		if c, ok := g.model.byGoName[obj.Name()]; ok {
			return c.Name
		}
		return ""
	}
	if g.otherByGoName == nil {
		g.otherByGoName = map[string]string{}
		for name, c := range g.model.classes {
			path, pkgName, ok := goPackageOf(name)
			if !ok || path == "" || c.Anonymous {
				continue
			}
			g.otherByGoName[pkgName+"."+g.otherGoName(c)] = name
		}
	}
	return g.otherByGoName[obj.Pkg().Name()+"."+obj.Name()]
}

// otherGoName is the Go type name of a class of a package below jscomp.
func (g *Gen) otherGoName(c *ClassInfo) string {
	if c.OuterOf != "" {
		if o, ok := g.model.classes[c.OuterOf]; ok && !o.Anonymous {
			return g.otherGoName(o) + "_" + c.SimpleName
		}
		return ""
	}
	return c.SimpleName
}

// complete adds, to every pending struct/enum class, the methods of the
// interfaces its Java class implements: the method sets of ported Go
// interfaces, and the declared methods of pending interfaces.
func (g *Gen) complete() bool {
	changed := false
	var names []string
	for n := range g.classes {
		names = append(names, n)
	}
	sort.Strings(names)
	// The static fields that recovered initializers refer to.
	for _, name := range names {
		p := g.classes[name]
		for _, m := range sortedMembers(p.statics) {
			if m.kind != MStaticVar {
				continue
			}
			init := m.field.owner.StaticInits[m.field.f.Name]
			if init == nil {
				continue
			}
			for _, ref := range init.Refs {
				if g.declareStaticField(ref[0], ref[1]) {
					changed = true
				}
			}
		}
	}
	for _, name := range names {
		p := g.classes[name]
		if p.kind == KindIface {
			continue
		}
		for _, s := range g.model.supertypes(p.c) {
			if sp, ok := g.classes[g.model.goNames[s]]; ok && sp.kind == KindIface && g.model.goNames[s] != "" {
				for _, m := range sortedMembers(sp.members) {
					if m.kind != MMethod && m.kind != MGoIfaceMethod {
						continue
					}
					if _, has := p.members[m.goName]; !has {
						cp := *m
						p.members[m.goName] = &cp
						changed = true
					}
				}
				continue
			}
			iface, ifaceName := g.realInterface(s)
			if iface == nil {
				continue
			}
			ms := types.NewMethodSet(iface)
			for i := 0; i < ms.Len(); i++ {
				fn := ms.At(i).Obj().(*types.Func)
				if _, has := p.members[fn.Name()]; has {
					continue
				}
				sig := fn.Type().(*types.Signature)
				if strings.Contains(types.TypeString(sig, nil), "invalid type") {
					// The interface itself refers to undefined types: declare
					// the Java method instead (its types get declared too).
					if ref := g.javaMethodForGo(p.c, fn.Name(), sig.Params().Len()); ref != nil {
						p.members[fn.Name()] = &PMember{kind: MMethod, goName: fn.Name(), method: ref}
						changed = true
					}
					continue
				}
				p.members[fn.Name()] = &PMember{kind: MGoIfaceMethod, goName: fn.Name(), goSig: sig, goIface: ifaceName}
				changed = true
			}
		}
	}
	return changed
}

// javaMethodForGo finds the Java instance method a Go method name with the
// given arity stands for (the longest Java name that prefixes it).
func (g *Gen) javaMethodForGo(c *ClassInfo, goName string, arity int) *methodRef {
	var best *methodRef
	bestLen := 0
	consider := func(cls *ClassInfo) {
		for i, mi := range cls.Methods {
			if mi.Access&(accStatic|accSynthetic|accBridge) != 0 || strings.HasPrefix(mi.Name, "<") {
				continue
			}
			base := capitalize(mi.Name)
			if !strings.HasPrefix(goName, base) || !hasSuffixBoundary(goName[len(base):]) {
				continue
			}
			if len(methodSig(mi).Params) != arity || len(base) < bestLen {
				continue
			}
			if len(base) > bestLen || best == nil {
				best = &methodRef{owner: cls, m: mi, index: i}
				bestLen = len(base)
			}
		}
	}
	consider(c)
	for _, s := range g.model.supertypes(c) {
		if sc, ok := g.model.classes[s]; ok {
			consider(sc)
		}
	}
	return best
}

// realInterface returns the ported Go interface type for a Java class, if
// it is one.
func (g *Gen) realInterface(internalName string) (types.Type, string) {
	path, pkgName, ok := goPackageOf(internalName)
	if !ok {
		return nil, ""
	}
	var obj types.Object
	var qname string
	if path == "" {
		name := g.model.goNames[internalName]
		if name == "" {
			return nil, ""
		}
		obj = g.realObject(name)
		qname = name
	} else {
		c, ok := g.model.classes[internalName]
		if !ok {
			return nil, ""
		}
		name := g.otherGoName(c)
		pkg := g.importedPackage(path)
		if pkg == nil {
			return nil, ""
		}
		obj = pkg.Scope().Lookup(name)
		qname = pkgName + "." + name
	}
	tn, ok := obj.(*types.TypeName)
	if !ok {
		return nil, ""
	}
	if _, isIface := tn.Type().Underlying().(*types.Interface); !isIface {
		return nil, ""
	}
	if named, ok := tn.Type().(*types.Named); ok && named.TypeParams().Len() > 0 {
		return nil, "" // generic interfaces need instantiation; leave to the compiler
	}
	return tn.Type(), qname
}

// importedPackage returns a dependency package of jscomp by import path.
func (g *Gen) importedPackage(path string) *types.Package {
	if g.res == nil || g.res.pkg == nil {
		return nil
	}
	for _, p := range g.res.pkg.Imports() {
		if p.Path() == path {
			return p
		}
	}
	p, err := g.pkg.importer.Import(path)
	if err != nil {
		return nil
	}
	return p
}

func sortedMembers(m map[string]*PMember) []*PMember {
	var out []*PMember
	for _, v := range m {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		oi, oj := memberOrder(out[i]), memberOrder(out[j])
		if oi != oj {
			return oi < oj
		}
		return out[i].goName < out[j].goName
	})
	return out
}

// memberOrder orders members as Java declares them.
func memberOrder(m *PMember) int {
	switch {
	case m.field != nil:
		return m.field.index
	case m.method != nil:
		return 100000 + m.method.index
	}
	return 1 << 30
}

// staticFieldGoName is the Go name of a Java static field: Class_FIELD.
func (g *Gen) staticFieldGoName(owner, name string) string {
	return g.model.goNames[owner] + "_" + name
}

// declareStaticField declares a static field of a pending class (one that
// a recovered initializer refers to); it reports whether it added one.
func (g *Gen) declareStaticField(owner, name string) bool {
	c := g.model.classes[owner]
	if c == nil || g.classPorted(c) {
		return false
	}
	f := g.model.fieldNamed(c, name, true)
	if f == nil {
		return false
	}
	goName := g.staticFieldGoName(owner, name)
	p := g.ensureClass(c)
	if p == nil {
		return false
	}
	if _, has := p.statics[goName]; has {
		return false
	}
	p.statics[goName] = &PMember{kind: MStaticVar, goName: goName, field: f}
	return true
}
