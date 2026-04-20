// Port of com/google/javascript/jscomp/TypeCheck.java's
// checkPropertyVisibility / validateMethodOverride branches
// @ closure-compiler v20260415 (narrow; no access-modifier check,
// no `@override`-required diagnostic).
//
// Fires JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH when an override of
// a superclass property has an incompatible type. Two override
// shapes are recognised:
//
//   X.prototype.method = function(...) {...};   // if X extends Y
//   /** @override */
//   Y.prototype.method = …
//
//   class Derived extends Base { method() { ... } }
//
// The check is scope-driven: it walks class / prototype-assign AST
// nodes, resolves the override's type via TypeInference, and
// compares against the superclass's slot via ObjectType's super
// chain.

package jscomp

import (
	"fmt"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

type SuperclassCheck struct {
	path       string
	creator    *TypedScopeCreator
	scope      *TypedScope
	scriptRoot *rhino.Node
	diags      []diag.Diagnostic
	// SuppressedNames lists variable names whose declaration already
	// triggered JSC_VAR_MULTIPLY_DECLARED_ERROR. @implements /
	// @interface-derived diagnostics keyed on those names are
	// suppressed — Java's TypeCheck skips them when the symbol slot
	// was rejected by an earlier VarCheck pass.
	SuppressedNames map[string]bool
}

func NewSuperclassCheck(path string, creator *TypedScopeCreator, scope *TypedScope) *SuperclassCheck {
	return &SuperclassCheck{path: path, creator: creator, scope: scope}
}

func (s *SuperclassCheck) Check(script *rhino.Node) []diag.Diagnostic {
	s.diags = nil
	s.scriptRoot = script
	s.walk(script)
	return s.diags
}

func (s *SuperclassCheck) walk(n *rhino.Node) {
	if n == nil {
		return
	}
	switch n.Token() {
	case rhino.CLASS:
		s.checkClass(n)
		s.checkClassDictExtendsStruct(n)
	case rhino.FUNCTION:
		s.checkFunctionImplements(n)
		s.checkInterfaceExtendsConflicts(n)
		s.checkExtendsKindMismatch(n)
		s.checkStructDictRequireConstructor(n)
		s.checkFunctionDictExtendsStruct(n)
	case rhino.VAR, rhino.LET, rhino.CONST:
		s.checkDeclaredConstructorImplements(n)
	case rhino.EXPR_RESULT:
		if a := n.FirstChild(); a != nil && a.Token() == rhino.ASSIGN {
			s.checkPrototypeAssign(n, a)
			s.checkPrototypeObjectLitAssign(n, a)
		}
		if gp := n.FirstChild(); gp != nil && gp.Token() == rhino.GETPROP {
			s.checkPrototypeStub(n, gp)
		}
	}
	for _, c := range n.Children() {
		s.walk(c)
	}
}

// checkFunctionDictExtendsStruct fires JSC_DICT_EXTEND_STRUCT_TYPE
// when a @constructor @dict function extends a @struct parent.
// Mirrors TypeCheck.java's struct/dict hierarchy guard.
func (s *SuperclassCheck) checkFunctionDictExtendsStruct(fn *rhino.Node) {
	info := fn.JSDocInfo()
	if info == nil || !info.IsDict() || !info.IsConstructor() {
		return
	}
	base := info.BaseType()
	if base == nil {
		return
	}
	super := s.scope.LookupNominalType(base.TypeName)
	if super == nil || !super.IsStruct {
		return
	}
	kids := fn.Children()
	if len(kids) < 1 {
		return
	}
	name := kids[0].String_()
	if name == "" {
		return
	}
	s.diags = append(s.diags, diag.Diagnostic{
		Path: s.path, Line: fn.Lineno(), Col: fn.Charno(),
		Code: diag.JSC_DICT_EXTEND_STRUCT_TYPE, Level: diag.Error,
		Message: fmt.Sprintf("@dict class %s cannot extend @struct class %s", name, base.TypeName),
	})
}

// checkClassDictExtendsStruct is the ES6-class variant of
// checkFunctionDictExtendsStruct. ES6 classes are @struct by
// default — a @dict-tagged class extending any normal class
// (struct-implicit) triggers the same diagnostic.
func (s *SuperclassCheck) checkClassDictExtendsStruct(cls *rhino.Node) {
	info := cls.JSDocInfo()
	if info == nil || !info.IsDict() {
		return
	}
	kids := cls.Children()
	if len(kids) < 2 {
		return
	}
	nameNode, superNode := kids[0], kids[1]
	name := nameNode.String_()
	if name == "" {
		return
	}
	if superNode == nil || superNode.String_() == "" {
		return
	}
	// ES6 classes are @struct by default (unless explicitly @dict).
	superName := superNode.String_()
	super := s.scope.LookupNominalType(superName)
	isSuperStruct := super != nil && super.IsStruct
	// Treat a plain ES6-class parent (no JSDoc) as struct-by-default.
	if super != nil && !super.IsStruct && !super.IsDict && !super.IsStructural {
		isSuperStruct = true
	}
	if !isSuperStruct {
		return
	}
	s.diags = append(s.diags, diag.Diagnostic{
		Path: s.path, Line: cls.Lineno(), Col: cls.Charno(),
		Code: diag.JSC_DICT_EXTEND_STRUCT_TYPE, Level: diag.Error,
		Message: fmt.Sprintf("@dict class %s cannot extend @struct class %s", name, superName),
	})
}

// checkStructDictRequireConstructor fires JSC_CONSTRUCTOR_REQUIRED
// when a FUNCTION carries @struct or @dict but no @constructor /
// @interface / @record. The struct/dict flag only has meaning on
// a class-shaped declaration — Java rejects it on plain functions.
func (s *SuperclassCheck) checkStructDictRequireConstructor(fn *rhino.Node) {
	info := fn.JSDocInfo()
	if info == nil {
		return
	}
	if !info.IsStruct() && !info.IsDict() {
		return
	}
	if info.IsConstructor() || info.IsInterface() || info.IsRecord() {
		return
	}
	kids := fn.Children()
	if len(kids) < 1 {
		return
	}
	name := kids[0].String_()
	if name == "" {
		return
	}
	tag := "@struct"
	if info.IsDict() {
		tag = "@dict"
	}
	s.diags = append(s.diags, diag.Diagnostic{
		Path: s.path, Line: fn.Lineno(), Col: fn.Charno(),
		Code: diag.JSC_CONSTRUCTOR_REQUIRED, Level: diag.Error,
		Message: fmt.Sprintf("%s used without @constructor for %s", tag, name),
	})
}

// checkExtendsKindMismatch fires JSC_CONFLICTING_EXTENDED_TYPE
// when an @interface extends a non-interface, or a @constructor
// extends a non-constructor. Mirrors TypeCheck.java's
// checkConstructor / checkInterface-extends kind guard.
func (s *SuperclassCheck) checkExtendsKindMismatch(fn *rhino.Node) {
	info := fn.JSDocInfo()
	if info == nil {
		return
	}
	isIface := info.IsInterface() || info.IsRecord()
	isCtor := info.IsConstructor()
	if !isIface && !isCtor {
		return
	}
	kids := fn.Children()
	if len(kids) < 1 {
		return
	}
	name := kids[0].String_()
	if name == "" {
		return
	}
	// Walk every @extends — fires once for the first parent whose
	// kind conflicts with the child's. Matches Java's behaviour of
	// reporting one conflicting-extends diagnostic per declaration.
	bases := info.BaseTypes()
	if len(bases) == 0 {
		if b := info.BaseType(); b != nil {
			bases = []*rhino.JSTypeExpression{b}
		}
	}
	reported := false
	for _, b := range bases {
		if reported {
			break
		}
		super := s.scope.LookupNominalType(b.TypeName)
		if super == nil {
			continue
		}
		if isIface && !super.IsStructural {
			s.diags = append(s.diags, diag.Diagnostic{
				Path: s.path, Line: fn.Lineno(), Col: fn.Charno(),
				Code: diag.JSC_CONFLICTING_EXTENDED_TYPE, Level: diag.Error,
				Message: fmt.Sprintf(
					"%s cannot extend this type; interfaces can only extend interfaces",
					name),
			})
			reported = true
		}
		if isCtor && super.IsStructural {
			s.diags = append(s.diags, diag.Diagnostic{
				Path: s.path, Line: fn.Lineno(), Col: fn.Charno(),
				Code: diag.JSC_CONFLICTING_EXTENDED_TYPE, Level: diag.Error,
				Message: fmt.Sprintf(
					"%s cannot extend this type; constructors can only extend constructors",
					name),
			})
			reported = true
		}
	}
}

// checkInterfaceExtendsConflicts fires
// JSC_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE when an @interface
// @extends multiple other interfaces and two of those parents
// declare the same property with incompatible types. One
// diagnostic per pair of conflicting parents, reported at the
// interface declaration's position. Mirrors TypeCheck.java's
// checkConflictingExtendedInterfaces loop.
func (s *SuperclassCheck) checkInterfaceExtendsConflicts(fn *rhino.Node) {
	info := fn.JSDocInfo()
	if info == nil || !info.IsInterface() {
		return
	}
	bases := info.BaseTypes()
	if len(bases) < 2 {
		return
	}
	kids := fn.Children()
	if len(kids) < 1 {
		return
	}
	name := kids[0].String_()
	if name == "" {
		return
	}
	// Resolve each direct parent interface to its InstanceType,
	// tracking the declared name (including template args, e.g.
	// `Int0<number>`) for diagnostic messages.
	type parentEnt struct {
		inst *jstype.ObjectType
		name string
	}
	var directParents []parentEnt
	for _, b := range bases {
		inst := s.scope.LookupNominalType(stripTemplateArgs(b.TypeName))
		if inst == nil {
			continue
		}
		directParents = append(directParents, parentEnt{inst: inst, name: b.TypeName})
	}
	if len(directParents) < 2 {
		return
	}
	// Port of TypeCheck.java's checkInterfaceConflictProperties loop.
	// `properties` carries property → owning-interface-name for every
	// parent-subtree processed so far; `current` carries the same
	// data for the parent currently being walked. A conflict is
	// reported each time a parent declares a property whose type
	// disagrees with an earlier-recorded parent's declared type.
	properties := map[string]string{}
	reported := map[string]bool{}
	var walk func(ownerName string, inst *jstype.ObjectType, current map[string]string)
	walk = func(ownerName string, inst *jstype.ObjectType, current map[string]string) {
		if inst == nil {
			return
		}
		for _, pname := range inst.MemberNames() {
			tj := inst.GetPropertyType(pname)
			if tj == nil {
				continue
			}
			otherOwner, ok := properties[pname]
			current[pname] = ownerName
			if !ok {
				continue
			}
			otherInst := s.scope.LookupNominalType(stripTemplateArgs(otherOwner))
			var ti jstype.JSType
			if otherInst != nil {
				ti = otherInst.GetPropertyType(pname)
			}
			if ti == nil {
				continue
			}
			if jstype.IsSubtype(ti, tj) || jstype.IsSubtype(tj, ti) {
				continue
			}
			key := pname + "\x00" + otherOwner + "\x00" + ownerName
			if reported[key] {
				continue
			}
			reported[key] = true
			s.diags = append(s.diags, diag.Diagnostic{
				Path: s.path, Line: fn.Lineno(), Col: fn.Charno(),
				Code: diag.JSC_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE, Level: diag.Error,
				Message: fmt.Sprintf(
					"Interface %s has a property %s with incompatible types in its super interfaces %s and %s",
					name, pname, otherOwner, ownerName),
			})
		}
		if exts := inst.ExtendedInterfaces(); len(exts) > 0 {
			for _, sup := range exts {
				walk(sup.TypeName(), sup, current)
			}
		} else if sup := inst.Super(); sup != nil {
			walk(sup.TypeName(), sup, current)
		}
	}
	for _, pe := range directParents {
		current := map[string]string{}
		walk(pe.name, pe.inst, current)
		for k, v := range current {
			properties[k] = v
		}
	}
}

// checkFunctionImplements walks an AST FUNCTION declaration and,
// when it carries @constructor + @implements {Foo, Bar, …}, emits
// JSC_INTERFACE_METHOD_NOT_IMPLEMENTED for each interface property
// the constructor's instance doesn't implement (directly or
// through its @extends chain). Mirrors TypeCheck's
// validateImplements loop.
func (s *SuperclassCheck) checkFunctionImplements(fn *rhino.Node) {
	info := fn.JSDocInfo()
	if !info.IsConstructor() {
		return
	}
	kids := fn.Children()
	if len(kids) < 1 {
		return
	}
	nameNode := kids[0]
	if nameNode.String_() == "" {
		return
	}
	s.validateImplements(fn, nameNode.String_(), info)
}

// checkDeclaredConstructorImplements handles the
// `/** @constructor @implements {Foo} */ var Bar = function(){}`
// form — JSDoc is on the VAR, not on the inner FUNCTION.
func (s *SuperclassCheck) checkDeclaredConstructorImplements(v *rhino.Node) {
	info := v.JSDocInfo()
	if !info.IsConstructor() {
		return
	}
	nameNode := v.FirstChild()
	if nameNode == nil || nameNode.Token() != rhino.NAME || nameNode.String_() == "" {
		return
	}
	// Report at the FUNCTION initialiser's position when present —
	// Java's TypeCheck fires INTERFACE_METHOD_NOT_IMPLEMENTED at the
	// function-expression keyword column, not the VAR's (which is
	// 0:0 for synthetic VARs emitted by the parser).
	reportAt := v
	if init := nameNode.FirstChild(); init != nil && init.Token() == rhino.FUNCTION {
		reportAt = init
	}
	s.validateImplements(reportAt, nameNode.String_(), info)
}

// validateImplements verifies `ctorName` implements every property
// each of its @implements interfaces declares. Emits one
// JSC_INTERFACE_METHOD_NOT_IMPLEMENTED per missing property. The
// diagnostic position matches Java's — at the constructor
// declaration node's starting column.
func (s *SuperclassCheck) validateImplements(reportAt *rhino.Node, ctorName string, info *rhino.JSDocInfo) {
	impls := info.Implementations()
	if len(impls) == 0 {
		return
	}
	ctor, ok := s.scope.GetRawDeclaredType(ctorName).(*jstype.FunctionType)
	if !ok || ctor.Instance == nil {
		return
	}
	inst := ctor.Instance
	// Detect duplicate @implements — Java fires
	// JSC_SAME_INTERFACE_MULTIPLE_IMPLEMENTS when the same
	// interface name appears twice (ignoring template args).
	seenIface := map[string]bool{}
	for _, expr := range impls {
		base := expr.TypeName
		// Strip template-arg suffix when present.
		if lt := strings.Index(base, "<"); lt > 0 {
			base = base[:lt]
		}
		// Strip `!` / `?` nullability markers.
		base = strings.TrimLeft(base, "!?")
		if seenIface[base] {
			s.diags = append(s.diags, diag.Diagnostic{
				Path: s.path, Line: reportAt.Lineno(), Col: reportAt.Charno(),
				Code: diag.JSC_SAME_INTERFACE_MULTIPLE_IMPLEMENTS, Level: diag.Error,
				Message: "Cannot @implement the same interface more than once\nRepeated interface: " + base,
			})
			continue
		}
		seenIface[base] = true
	}
	for _, expr := range impls {
		iface := s.scope.LookupNominalType(expr.TypeName)
		if iface == nil {
			continue
		}
		// Skip when the interface name was flagged by a prior
		// VAR_MULTIPLY_DECLARED_ERROR — the script's @interface
		// declaration was rejected, so its property set should not
		// drive INTERFACE_METHOD_NOT_IMPLEMENTED diagnostics.
		ifaceName := stripTemplateArgs(expr.TypeName)
		if s.SuppressedNames[ifaceName] {
			continue
		}
		// Java's TypeCheck fires JSC_IMPLEMENTS_NON_INTERFACE
		// when @implements references a non-interface constructor
		// (e.g. `@implements {Boolean}` or a plain user-defined
		// @constructor). Gated on the target not being
		// @record / @interface.
		if !iface.IsStructural {
			s.diags = append(s.diags, diag.Diagnostic{
				Path: s.path, Line: reportAt.Lineno(), Col: reportAt.Charno(),
				Code: diag.JSC_IMPLEMENTS_NON_INTERFACE, Level: diag.Error,
				Message: "can only implement interfaces",
			})
			continue
		}
		// Walk the interface's full property set (its own members
		// plus any inherited via @extends {Other}). The reported
		// owner is the interface that actually declares the property,
		// not the @implements'd root — matches Java's
		// TypeCheck.validateInterfaceMethodImplemented.
		for _, pe := range allInterfacePropertyOwners(iface) {
			if inst.GetPropertyType(pe.name) != nil {
				continue
			}
			s.diags = append(s.diags, diag.Diagnostic{
				Path: s.path, Line: reportAt.Lineno(), Col: reportAt.Charno(),
				Code: diag.JSC_INTERFACE_METHOD_NOT_IMPLEMENTED, Level: diag.Error,
				Message: fmt.Sprintf(
					"property %s on interface %s is not implemented by type %s",
					pe.name, pe.owner, ctorName),
			})
		}
	}
}

// allInterfaceProperties returns every property name declared on
// `iface` or reachable through its super chain — either the single
// @extends (via Super()) or the multi-extends list
// (ExtendedInterfaces, used by interfaces that declare more than
// one parent). Deduplicated. Mirrors Java's
// FunctionType.getAllImplementedInterfaces walk.
func allInterfaceProperties(iface *jstype.ObjectType) []string {
	var out []string
	for _, pe := range allInterfacePropertyOwners(iface) {
		out = append(out, pe.name)
	}
	return out
}

// propOwner pairs a property name with the interface in the super
// chain that actually declares it.
type propOwner struct {
	name  string
	owner string
}

// allInterfacePropertyOwners returns every property on `iface`'s
// transitive supertype graph alongside the name of the interface
// that directly declares each property. First-seen wins (matching
// Java's interfaceToImplementingClasses walk where a child's
// declaration shadows the grandparent's).
func allInterfacePropertyOwners(iface *jstype.ObjectType) []propOwner {
	seen := map[string]bool{}
	var out []propOwner
	visited := map[*jstype.ObjectType]bool{}
	var walk func(cur *jstype.ObjectType)
	walk = func(cur *jstype.ObjectType) {
		if cur == nil || visited[cur] {
			return
		}
		visited[cur] = true
		for _, name := range cur.MemberNames() {
			if seen[name] {
				continue
			}
			seen[name] = true
			out = append(out, propOwner{name: name, owner: cur.TypeName()})
		}
		if exts := cur.ExtendedInterfaces(); len(exts) > 0 {
			for _, sup := range exts {
				walk(sup)
			}
		} else {
			walk(cur.Super())
		}
	}
	walk(iface)
	return out
}

// checkPrototypeStub handles the stub form `X.prototype.name;` with
// a `/** @type {T} */` JSDoc. Fires JSC_HIDDEN_SUPERCLASS /
// INTERFACE_PROPERTY_MISMATCH when the stub's declared type
// conflicts with a super / @implements slot of the same name.
func (s *SuperclassCheck) checkPrototypeStub(stmt *rhino.Node, gp *rhino.Node) {
	kids := gp.Children()
	if len(kids) < 2 {
		return
	}
	inner := kids[0]
	propName := kids[1].String_()
	propNode := kids[1]
	if inner.Token() != rhino.GETPROP {
		return
	}
	innerKids := inner.Children()
	if len(innerKids) < 2 {
		return
	}
	if innerKids[0].Token() != rhino.NAME || innerKids[1].String_() != "prototype" {
		return
	}
	ownerName := innerKids[0].String_()
	fn, ok := s.scope.GetRawDeclaredType(ownerName).(*jstype.FunctionType)
	if !ok || fn.Instance == nil {
		return
	}
	info := stmt.JSDocInfo()
	if info == nil {
		info = gp.JSDocInfo()
	}
	if info == nil || !info.HasType() {
		return
	}
	overrideType := s.scope.ResolveExpr(info.Type().TypeName)
	if overrideType == nil {
		return
	}
	// Superclass slot. Use HIDDEN_INTERFACE when the owner itself is
	// an interface (so the super chain is all interfaces) or when
	// the super is structural. Data-field overrides on concrete
	// classes pass through without complaint — Java defers type-
	// mismatch to use-site checks for non-function slots.
	if super := fn.Instance.Super(); super != nil {
		slot := super.GetPropertyType(propName)
		if slot != nil && !jstype.IsSubtype(overrideType, slot) && isFunctionTypeSlot(slot) {
			if fn.Instance.IsStructural || super.IsStructural {
				s.diags = append(s.diags, diag.Diagnostic{
					Path: s.path, Line: propNode.Lineno(), Col: propNode.Charno(),
					Code: diag.JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH, Level: diag.Error,
					Message: fmt.Sprintf(
						"mismatch of the %s property on type %s and the type of the property it overrides from interface %s\noriginal: %s\noverride: %s",
						propName, ownerName, super.TypeName(), slot.TypeName(), overrideType.TypeName()),
				})
			} else {
				s.diags = append(s.diags, diag.Diagnostic{
					Path: s.path, Line: propNode.Lineno(), Col: propNode.Charno(),
					Code: diag.JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH, Level: diag.Error,
					Message: fmt.Sprintf(
						"mismatch of the %s property type and the type of the property it overrides from superclass %s\noriginal: %s\noverride: %s",
						propName, super.TypeName(), slot.TypeName(), overrideType.TypeName()),
				})
			}
			return
		}
	}
	// Transitive base interfaces — only when the owner is a class /
	// @constructor @implements some iface, OR the owner is an
	// interface with additional @extends parents beyond the direct
	// Super chain.
	bases := baseInterfaceInstances(fn.Instance)
	for _, iface := range bases {
		slot := iface.GetPropertyType(propName)
		if slot == nil || jstype.IsSubtype(overrideType, slot) {
			continue
		}
		s.diags = append(s.diags, diag.Diagnostic{
			Path: s.path, Line: propNode.Lineno(), Col: propNode.Charno(),
			Code: diag.JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH, Level: diag.Error,
			Message: fmt.Sprintf(
				"mismatch of the %s property on type %s and the type of the property it overrides from interface %s\noriginal: %s\noverride: %s",
				propName, ownerName, iface.TypeName(), slot.TypeName(), overrideType.TypeName()),
		})
		return
	}
}

// isFunctionTypeSlot reports whether the slot's declared type is
// itself a function. Used to narrow the override-type-mismatch
// check to method slots; data-field overrides are accepted by
// Java in default mode.
func isFunctionTypeSlot(t jstype.JSType) bool {
	if t == nil {
		return false
	}
	if _, ok := t.(*jstype.FunctionType); ok {
		return true
	}
	return false
}

// baseInterfaceInstances returns every super InstanceType reachable
// through `obj`'s @extends chain (for interfaces, a single chain is
// typical; for multi-extends interfaces we return each parent).
func baseInterfaceInstances(obj *jstype.ObjectType) []*jstype.ObjectType {
	var out []*jstype.ObjectType
	seen := map[*jstype.ObjectType]bool{}
	var walk func(o *jstype.ObjectType)
	walk = func(o *jstype.ObjectType) {
		if o == nil || seen[o] {
			return
		}
		seen[o] = true
		if s := o.Super(); s != nil && s.IsStructural {
			out = append(out, s)
			walk(s)
		}
	}
	walk(obj)
	return out
}

// checkPrototypeAssign handles `X.prototype.name = rhs;`. When X has
// a super (via @extends) and super declares `name`, verify the rhs
// type matches — else fire JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH.
// Same shape runs against each @implemented interface's slots,
// firing JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH instead.
func (s *SuperclassCheck) checkPrototypeAssign(stmt *rhino.Node, assign *rhino.Node) {
	kids := assign.Children()
	if len(kids) < 2 {
		return
	}
	lhs, rhs := kids[0], kids[1]
	if lhs.Token() != rhino.GETPROP {
		return
	}
	lhsKids := lhs.Children()
	if len(lhsKids) < 2 {
		return
	}
	inner := lhsKids[0]
	propName := lhsKids[1].String_()
	propNode := lhsKids[1]
	if inner.Token() != rhino.GETPROP {
		return
	}
	innerKids := inner.Children()
	if len(innerKids) < 2 {
		return
	}
	if innerKids[0].Token() != rhino.NAME || innerKids[1].String_() != "prototype" {
		return
	}
	ownerName := innerKids[0].String_()
	fn, ok := s.scope.GetRawDeclaredType(ownerName).(*jstype.FunctionType)
	if !ok || fn.Instance == nil {
		return
	}
	overrideType := s.inferAssignValue(stmt, assign, rhs, fn.Instance)
	if overrideType == nil {
		return
	}
	// Superclass slot (via @extends chain).
	if super := fn.Instance.Super(); super != nil {
		if super.IsPropertyInferred(propName) {
			return
		}
		if slot := super.GetPropertyType(propName); slot != nil {
			if !overrideCompatible(overrideType, slot) {
				s.diags = append(s.diags, diag.Diagnostic{
					Path:  s.path,
					Line:  stmt.Lineno(),
					Col:   stmt.Charno(),
					Code:  diag.JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH,
					Level: diag.Error,
					Message: fmt.Sprintf("mismatch of the %s property type and the type of the property it overrides from superclass %s\noriginal: %s\noverride: %s",
						propName, super.TypeName(), slot.TypeName(), overrideType.TypeName()),
				})
				return
			}
		}
	}
	// @implements interfaces.
	for _, iface := range s.implementedInstances(ownerName) {
		slot := iface.GetPropertyType(propName)
		if slot == nil {
			continue
		}
		if overrideCompatible(overrideType, slot) {
			continue
		}
		s.diags = append(s.diags, diag.Diagnostic{
			Path:  s.path,
			Line:  propNode.Lineno(),
			Col:   propNode.Charno(),
			Code:  diag.JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH,
			Level: diag.Error,
			Message: fmt.Sprintf("mismatch of the %s property on type %s and the type of the property it overrides from interface %s\noriginal: %s\noverride: %s",
				propName, ownerName, iface.TypeName(), slot.TypeName(), overrideType.TypeName()),
		})
		return
	}
}

// checkPrototypeObjectLitAssign handles the wholesale-prototype
// form `Foo.prototype = { k: v, … }`. Each STRING_KEY's value is
// compared against the superclass's slot of the same name, firing
// JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH on type mismatch. Mirrors
// TypeCheck.java's validateMethodOverride when the override is
// attached via an OBJECTLIT rather than a dotted ASSIGN.
func (s *SuperclassCheck) checkPrototypeObjectLitAssign(stmt *rhino.Node, assign *rhino.Node) {
	kids := assign.Children()
	if len(kids) < 2 {
		return
	}
	lhs, rhs := kids[0], kids[1]
	if lhs.Token() != rhino.GETPROP || rhs.Token() != rhino.OBJECTLIT {
		return
	}
	lhsKids := lhs.Children()
	if len(lhsKids) < 2 {
		return
	}
	if lhsKids[0].Token() != rhino.NAME || lhsKids[1].String_() != "prototype" {
		return
	}
	ownerName := lhsKids[0].String_()
	fn, ok := s.scope.GetRawDeclaredType(ownerName).(*jstype.FunctionType)
	if !ok || fn.Instance == nil {
		return
	}
	super := fn.Instance.Super()
	ifaces := s.implementedInstances(ownerName)
	if super == nil && len(ifaces) == 0 {
		return
	}
	for _, kv := range rhs.Children() {
		if kv.Token() != rhino.STRING_KEY {
			continue
		}
		propName := kv.String_()
		if propName == "" {
			continue
		}
		valKids := kv.Children()
		if len(valKids) == 0 {
			continue
		}
		val := valKids[0]
		info := kv.JSDocInfo()
		var overrideType jstype.JSType
		if info != nil && info.HasType() {
			overrideType = s.scope.ResolveExpr(info.Type().TypeName)
		} else if val.Token() == rhino.FUNCTION {
			overrideType = functionTypeFromJSDoc(info, val, s.scope)
		} else {
			ti := NewTypeInference(s.scope)
			overrideType = ti.InferType(val)
		}
		if overrideType == nil {
			continue
		}
		// Superclass slot.
		if super != nil {
			if super.IsPropertyInferred(propName) {
				continue
			}
			if slot := super.GetPropertyType(propName); slot != nil {
				if !overrideCompatible(overrideType, slot) {
					s.diags = append(s.diags, diag.Diagnostic{
						Path:  s.path,
						Line:  kv.Lineno(),
						Col:   kv.Charno(),
						Code:  diag.JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH,
						Level: diag.Error,
						Message: fmt.Sprintf("mismatch of the %s property type and the type of the property it overrides from superclass %s\noriginal: %s\noverride: %s",
							propName, super.TypeName(), slot.TypeName(), overrideType.TypeName()),
					})
					continue
				}
			}
		}
		// @implements interfaces.
		for _, iface := range ifaces {
			slot := iface.GetPropertyType(propName)
			if slot == nil {
				continue
			}
			if overrideCompatible(overrideType, slot) {
				continue
			}
			s.diags = append(s.diags, diag.Diagnostic{
				Path:  s.path,
				Line:  kv.Lineno(),
				Col:   kv.Charno(),
				Code:  diag.JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH,
				Level: diag.Error,
				Message: fmt.Sprintf("mismatch of the %s property on type %s and the type of the property it overrides from interface %s\noriginal: %s\noverride: %s",
					propName, ownerName, iface.TypeName(), slot.TypeName(), overrideType.TypeName()),
			})
			break
		}
	}
}

// overrideCompatible reports whether `override` is a valid
// Java-style override of `slot`. Delegates to IsSubtype, but
// relaxes the check for FunctionTypes so an override with
// *fewer* parameters (Java's "extra params are OK to drop")
// still matches. Java's actual compatibility logic is more
// nuanced — we approximate by accepting any overlapping
// param-prefix whose types match.
func overrideCompatible(override, slot jstype.JSType) bool {
	if jstype.IsSubtype(override, slot) {
		return true
	}
	of, ok := override.(*jstype.FunctionType)
	if !ok {
		return false
	}
	sf, ok := slot.(*jstype.FunctionType)
	if !ok {
		return false
	}
	// Override with fewer params is acceptable to Java so long
	// as every overlapping param has a compatible type. Return
	// type must still be a subtype.
	if len(of.Params) > len(sf.Params) {
		return false
	}
	for i := range of.Params {
		if !jstype.IsSubtype(of.Params[i], sf.Params[i]) {
			return false
		}
	}
	if sf.Return != nil && of.Return != nil {
		if !jstype.IsSubtype(of.Return, sf.Return) {
			return false
		}
	}
	return true
}

// implementedInstances returns the interface instance types declared
// via `@implements {Foo}` on the constructor named by ownerName.
func (s *SuperclassCheck) implementedInstances(ownerName string) []*jstype.ObjectType {
	// The scope doesn't currently expose the JSDocInfo that produced
	// the FunctionType, so we walk the AST once to collect
	// @implements declarations. Caching would be worthwhile if this
	// became hot.
	info := s.findConstructorJSDoc(ownerName)
	if info == nil {
		return nil
	}
	var out []*jstype.ObjectType
	for _, e := range info.Implementations() {
		if inst := s.scope.LookupNominalType(e.TypeName); inst != nil {
			out = append(out, inst)
		}
	}
	return out
}

// findConstructorJSDoc scans the AST for a `function Foo(…)` or
// `var Foo = function(…)` constructor named `name` and returns its
// JSDocInfo. Returns nil when no matching declaration exists.
func (s *SuperclassCheck) findConstructorJSDoc(name string) *rhino.JSDocInfo {
	return s.scanConstructorJSDoc(s.scriptRoot, name)
}

func (s *SuperclassCheck) scanConstructorJSDoc(n *rhino.Node, name string) *rhino.JSDocInfo {
	if n == nil {
		return nil
	}
	for _, c := range n.Children() {
		if c.Token() == rhino.FUNCTION {
			if first := c.FirstChild(); first != nil && first.String_() == name {
				return c.JSDocInfo()
			}
		}
		if c.Token() == rhino.VAR || c.Token() == rhino.LET || c.Token() == rhino.CONST {
			nm := c.FirstChild()
			if nm != nil && nm.String_() == name {
				if init := nm.FirstChild(); init != nil && init.Token() == rhino.FUNCTION {
					if info := c.JSDocInfo(); info != nil {
						return info
					}
					return init.JSDocInfo()
				}
			}
		}
		if found := s.scanConstructorJSDoc(c, name); found != nil {
			return found
		}
	}
	return nil
}

// inferAssignValue resolves the type attached to `X.prototype.m = …`.
// Prefers an explicit @type / @param-based FunctionType from the
// enclosing ASSIGN's JSDoc; falls back to inferring the RHS via
// TypeInference.
func (s *SuperclassCheck) inferAssignValue(stmt, assign, rhs *rhino.Node, inst *jstype.ObjectType) jstype.JSType {
	info := assign.JSDocInfo()
	if info == nil {
		info = stmt.JSDocInfo()
	}
	if info != nil && info.HasType() {
		return s.scope.ResolveExpr(info.Type().TypeName)
	}
	if rhs.Token() == rhino.FUNCTION {
		return functionTypeFromJSDoc(info, rhs, s.scope)
	}
	ti := NewTypeInference(s.scope)
	return ti.InferType(rhs)
}

// checkClass verifies each method / field of a `class C extends B`
// declaration against B's slot.
func (s *SuperclassCheck) checkClass(cls *rhino.Node) {
	kids := cls.Children()
	if len(kids) < 3 {
		return
	}
	nameNode := kids[0]
	body := kids[2]
	preStub, ok := s.scope.GetRawDeclaredType(nameNode.String_()).(*jstype.FunctionType)
	if !ok || preStub.Instance == nil || preStub.Instance.Super() == nil {
		return
	}
	super := preStub.Instance.Super()
	for _, m := range body.Children() {
		switch m.Token() {
		case rhino.FUNCTION:
			if len(m.Children()) < 1 {
				continue
			}
			methodName := m.FirstChild().String_()
			if methodName == "" || methodName == "constructor" {
				continue
			}
			superSlot := super.GetPropertyType(methodName)
			if superSlot == nil {
				continue
			}
			info := m.JSDocInfo()
			ownSlot := functionTypeFromJSDoc(info, m, s.scope)
			if ownSlot == nil {
				continue
			}
			if jstype.IsSubtype(ownSlot, superSlot) {
				continue
			}
			s.diags = append(s.diags, diag.Diagnostic{
				Path:  s.path,
				Line:  m.Lineno(),
				Col:   m.Charno(),
				Code:  diag.JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH,
				Level: diag.Error,
				Message: fmt.Sprintf("mismatch of the %s property type and the type of the property it overrides from superclass %s\noriginal: %s\noverride: %s",
					methodName, super.TypeName(), superSlot.TypeName(), ownSlot.TypeName()),
			})
		case rhino.CLASS_FIELD_DEF:
			if len(m.Children()) < 1 {
				continue
			}
			fieldName := m.FirstChild().String_()
			if fieldName == "" {
				continue
			}
			superSlot := super.GetPropertyType(fieldName)
			if superSlot == nil {
				continue
			}
			info := m.JSDocInfo()
			if info == nil || !info.HasType() {
				continue
			}
			ownType := s.scope.ResolveExpr(info.Type().TypeName)
			if jstype.IsSubtype(ownType, superSlot) {
				continue
			}
			s.diags = append(s.diags, diag.Diagnostic{
				Path:  s.path,
				Line:  m.Lineno(),
				Col:   m.Charno(),
				Code:  diag.JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH,
				Level: diag.Error,
				Message: fmt.Sprintf("mismatch of the %s property type and the type of the property it overrides from superclass %s\noriginal: %s\noverride: %s",
					fieldName, super.TypeName(), superSlot.TypeName(), ownType.TypeName()),
			})
		}
	}
}

// functionTypeFromJSDoc builds a FunctionType from a FUNCTION node's
// JSDoc (@param / @return). Returns nil when the node doesn't carry
// enough JSDoc to infer a full signature.
func functionTypeFromJSDoc(info *rhino.JSDocInfo, fn *rhino.Node, scope *TypedScope) *jstype.FunctionType {
	kids := fn.Children()
	if len(kids) < 2 {
		return nil
	}
	params := kids[1]
	var paramTypes []jstype.JSType
	for _, p := range params.Children() {
		var pt jstype.JSType = jstype.Unknown
		if info != nil {
			if expr := info.ParameterType(p.String_()); expr != nil {
				pt = scope.ResolveExpr(expr.TypeName)
			}
		}
		paramTypes = append(paramTypes, pt)
	}
	var ret jstype.JSType = jstype.Unknown
	if info != nil && info.ReturnType() != nil {
		ret = scope.ResolveExpr(info.ReturnType().TypeName)
	}
	return jstype.NewFunctionType(paramTypes, ret)
}
