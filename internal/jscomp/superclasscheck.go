// Port of com/google/javascript/jscomp/TypeCheck.java's
// checkPropertyVisibility / validateMethodOverride branches
// @ closure-compiler v20260415 (narrow; no access-modifier check,
// no `@override`-required diagnostic).
//
// Fires JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH when an override of
// a superclass property has an incompatible type, and
// JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH when the incompatible type
// comes from an @implements interface. Three override shapes are
// recognised:
//
//   X.prototype.method = function(...) {...};   // if X extends Y
//   /** @override */
//   Y.prototype.method = …
//
//   class Derived extends Base { method() { ... } }
//   /** @implements {IFoo} */ class Bar { method() { ... } }
//
// The check is scope-driven: it walks class / prototype-assign AST
// nodes, resolves the override's type via TypeInference, and
// compares against the superclass's slot via ObjectType's super
// chain.

package jscomp

import (
	"fmt"
	"sort"
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
		s.checkClassImplements(n)
		s.checkExtendsNonObject(n)
	case rhino.FUNCTION:
		s.checkFunctionImplements(n)
		s.checkInterfaceExtendsConflicts(n)
		s.checkExtendsKindMismatch(n)
		s.checkExtendsNonObject(n)
		s.checkStructDictRequireConstructor(n)
		s.checkFunctionDictExtendsStruct(n)
	case rhino.VAR, rhino.LET, rhino.CONST:
		s.checkDeclaredConstructorImplements(n)
	case rhino.EXPR_RESULT:
		if a := n.FirstChild(); a != nil && a.Token() == rhino.ASSIGN {
			s.checkPrototypeAssign(n, a)
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

// checkExtendsNonObject fires JSC_EXTENDS_NON_OBJECT when @extends
// references a primitive or union type. Mirrors TypeCheck.java's
// visitClass / checkConstructor guard. Skips CLASS nodes whose JS
// `extends` clause is populated — Java gates on the inferred type of
// that expression, not the JSDoc annotation.
func (s *SuperclassCheck) checkExtendsNonObject(n *rhino.Node) {
	info := n.JSDocInfo()
	if info == nil {
		return
	}
	kids := n.Children()
	if len(kids) < 1 {
		return
	}
	name := kids[0].String_()
	if name == "" {
		return
	}
	if n.Token() == rhino.CLASS && len(kids) >= 2 && kids[1].String_() != "" {
		return
	}
	bases := info.BaseTypes()
	if len(bases) == 0 {
		if b := info.BaseType(); b != nil {
			bases = []*rhino.JSTypeExpression{b}
		}
	}
	for _, b := range bases {
		if s.scope.LookupNominalType(b.TypeName) != nil {
			continue // resolved to a nominal object type — fine
		}
		resolved := s.scope.ResolveExpr(b.TypeName)
		if resolved == nil || resolved == jstype.Unknown {
			// The tag resolved to nothing. Fire JSC_RESOLVED_TAG_EMPTY when
			// the type name IS a locally-declared identifier (param/var — so
			// JSC_UNRECOGNIZED_TYPE_ERROR is suppressed by allNames) but is
			// NOT in a template-parameterized context (which Java handles
			// specially). Skip wildcards and names that checkjsdoctypes
			// already reports as UNRECOGNIZED_TYPE_ERROR.
			tn := b.TypeName
			if tn != "?" && tn != "" && s.isDeclaredInAnyFunctionScope(tn, n) &&
				!enclosingFunctionHasTemplate(n) {
				s.diags = append(s.diags, diag.Diagnostic{
					Path: s.path, Line: n.Lineno(), Col: n.Charno(),
					Code:    diag.JSC_RESOLVED_TAG_EMPTY,
					Level:   diag.Error,
					Message: fmt.Sprintf("Could not resolve type in @extends tag of %s", name),
				})
			}
			break
		}
		if isObjectLikeType(resolved) {
			continue // object / templatized / function — fine
		}
		s.diags = append(s.diags, diag.Diagnostic{
			Path: s.path, Line: n.Lineno(), Col: n.Charno(),
			Code:    diag.JSC_EXTENDS_NON_OBJECT,
			Level:   diag.Error,
			Message: fmt.Sprintf("%s @extends non-object type %s", name, resolved.TypeName()),
		})
		break // one diagnostic per declaration
	}
}

// isDeclaredInAnyFunctionScope returns true when `name` appears as a declared
// binding in any function-local scope tracked by the creator. Used to
// distinguish `@extends {clazz}` (clazz is a local param → RESOLVED_TAG_EMPTY)
// from `@extends {not_base}` (not_base not declared → UNRECOGNIZED_TYPE_ERROR).
func (s *SuperclassCheck) isDeclaredInAnyFunctionScope(name string, near *rhino.Node) bool {
	// Walk the parent chain looking for the nearest enclosing FUNCTION
	// whose scope might have the name.
	for cur := near.Parent(); cur != nil; cur = cur.Parent() {
		if cur.Token() != rhino.FUNCTION {
			continue
		}
		if fnScope, ok := s.creator.Functions[cur]; ok && fnScope != nil {
			if fnScope.GetRawDeclaredType(name) != nil {
				return true
			}
		}
	}
	return false
}

// enclosingFunctionHasTemplate returns true when any ancestor FUNCTION node
// (up to the script root) carries a @template annotation. Java handles
// template-parameterized extends specially and does not fire RESOLVED_TAG_EMPTY.
func enclosingFunctionHasTemplate(n *rhino.Node) bool {
	for cur := n.Parent(); cur != nil; cur = cur.Parent() {
		if cur.Token() != rhino.FUNCTION {
			continue
		}
		if info := cur.JSDocInfo(); info != nil && len(info.TemplateNames()) > 0 {
			return true
		}
	}
	return false
}

// checkClassImplements fires JSC_IMPLEMENTS_NON_INTERFACE for ES6 class
// nodes, delegating to validateImplements like checkFunctionImplements.
func (s *SuperclassCheck) checkClassImplements(cls *rhino.Node) {
	info := cls.JSDocInfo()
	if info == nil {
		return
	}
	kids := cls.Children()
	if len(kids) < 1 {
		return
	}
	name := kids[0].String_()
	if name == "" {
		return
	}
	s.validateImplements(cls, name, info)
}

// ifacePropEntry holds an interface instance type plus the display
// name (e.g. "Int0<number>") that should appear in diagnostics.
// Needed because ObjectType.TypeName() loses template arguments.
type ifacePropEntry struct {
	obj     *jstype.ObjectType
	display string
}

// checkInterfaceExtendsConflicts fires
// JSC_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE when an @interface
// @extends multiple other interfaces and two of those parents
// declare the same property with incompatible types. Mirrors
// TypeCheck.java's checkConflictingExtendedInterfaces loop and
// checkInterfaceConflictProperties recursive helper — uses an
// accumulator map so each property is compared against the
// interface that last declared it, matching Java's pair naming.
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
	// properties accumulates propName → (ObjectType, displayName) of
	// the interface that last declared the property. Mirrors Java's
	// LinkedHashMap<Property.Key, ObjectType> properties.
	properties := map[string]ifacePropEntry{}
	currentProperties := map[string]ifacePropEntry{}
	for _, b := range bases {
		// Strip template arguments for lookup; keep original for display.
		baseName := stripTemplateArgs(b.TypeName)
		interfaceType := s.scope.LookupNominalType(baseName)
		if interfaceType == nil {
			continue
		}
		currentProperties = map[string]ifacePropEntry{}
		s.checkInterfaceConflictProperties(fn, name, properties, currentProperties, interfaceType, b.TypeName)
		for k, v := range currentProperties {
			properties[k] = v
		}
	}
}

// checkInterfaceConflictProperties mirrors TypeCheck.java's
// checkInterfaceConflictProperties. For each own property of
// interfaceType, it checks whether `properties` already holds a
// conflicting declaration from an earlier sibling. It then recurses
// into each interface that interfaceType extends (via JSDocInfo
// BaseTypes, to capture multi-inheritance), accumulating into
// currentProperties.
func (s *SuperclassCheck) checkInterfaceConflictProperties(
	fn *rhino.Node,
	functionName string,
	properties map[string]ifacePropEntry,
	currentProperties map[string]ifacePropEntry,
	interfaceType *jstype.ObjectType,
	displayName string,
) {
	// Collect own property names, sorted for deterministic output.
	ownNames := interfaceType.MemberNames()
	sort.Strings(ownNames)
	for _, pname := range ownNames {
		currentProperties[pname] = ifacePropEntry{obj: interfaceType, display: displayName}
		old := properties[pname]
		if old.obj == nil {
			continue
		}
		thisPropType := interfaceType.GetPropertyType(pname)
		oPropType := old.obj.GetPropertyType(pname)
		if thisPropType == nil || oPropType == nil {
			continue
		}
		if jstype.IsSubtype(thisPropType, oPropType) || jstype.IsSubtype(oPropType, thisPropType) {
			continue
		}
		s.diags = append(s.diags, diag.Diagnostic{
			Path: s.path, Line: fn.Lineno(), Col: fn.Charno(),
			Code: diag.JSC_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE, Level: diag.Error,
			Message: fmt.Sprintf(
				"Interface %s has a property %s with incompatible types in its super interfaces %s and %s",
				functionName, pname, old.display, displayName),
		})
	}
	// Recurse into the interfaces that interfaceType extends. Mirrors
	// Java's:
	//   for (ObjectType iType : interfaceType.getCtorExtendedInterfaces())
	// Our ObjectType.Super() only stores the first parent (single-
	// inheritance), so we scan the AST for the interface declaration's
	// JSDocInfo to enumerate all @extends declarations.
	docInfo := s.findConstructorJSDoc(interfaceType.TypeName())
	if docInfo == nil {
		// Fall back to single-inheritance Super chain.
		if sup := interfaceType.Super(); sup != nil {
			s.checkInterfaceConflictProperties(fn, functionName, properties, currentProperties, sup, sup.TypeName())
		}
		return
	}
	for _, b := range docInfo.BaseTypes() {
		baseName := stripTemplateArgs(b.TypeName)
		iType := s.scope.LookupNominalType(baseName)
		if iType == nil {
			continue
		}
		s.checkInterfaceConflictProperties(fn, functionName, properties, currentProperties, iType, b.TypeName)
	}
}

// stripTemplateArgs removes any `<…>` suffix from a type name so
// that "Int0<number>" → "Int0" for scope lookup purposes.
func stripTemplateArgs(name string) string {
	if i := strings.IndexByte(name, '<'); i >= 0 {
		return name[:i]
	}
	return name
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
	// Use the inner FUNCTION node as the report position, matching
	// Java's behaviour of pointing at the constructor expression.
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
		iface := s.scope.LookupNominalType(stripTemplateArgs(expr.TypeName))
		if iface == nil {
			// @implements a primitive or union (e.g. {string}, {Foo|Bar}).
			// Java's TypeCheck fires JSC_IMPLEMENTS_NON_INTERFACE in this
			// case since only object types can be interfaces.
			resolved := s.scope.ResolveExpr(expr.TypeName)
			if resolved != nil && resolved != jstype.Unknown && !isObjectLikeType(resolved) {
				s.diags = append(s.diags, diag.Diagnostic{
					Path: s.path, Line: reportAt.Lineno(), Col: reportAt.Charno(),
					Code: diag.JSC_IMPLEMENTS_NON_INTERFACE, Level: diag.Error,
					Message: "can only implement interfaces",
				})
			}
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
		// plus any inherited via @extends chains). Only fire for
		// function-typed slots — data properties and Unknown-typed
		// stubs (e.g. goog.abstractMethod assignments) are skipped,
		// matching Java's validateImplements which gates on the slot
		// being a function type before emitting the diagnostic.
		// The diagnostic names the interface that originally declared
		// the property (which may be a parent of `iface`), matching
		// Java's output for multi-extends chains.
		for propName, entry := range s.allInterfacePropEntries(iface) {
			if !isFunctionTypeSlot(entry.typ) {
				continue
			}
			if inst.GetPropertyType(propName) != nil {
				continue
			}
			s.diags = append(s.diags, diag.Diagnostic{
				Path: s.path, Line: reportAt.Lineno(), Col: reportAt.Charno(),
				Code: diag.JSC_INTERFACE_METHOD_NOT_IMPLEMENTED, Level: diag.Error,
				Message: fmt.Sprintf(
					"property %s on interface %s is not implemented by type %s",
					propName, entry.ifaceName, ctorName),
			})
		}
	}
}

// isObjectLikeType reports whether t is an object, templatized, or
// function type — i.e. NOT a raw primitive, multi-member union, or
// unknown. A nullable object (Foo|null) is treated as object-like to
// match Java's behaviour: only unions with multiple non-null members
// trigger JSC_EXTENDS_NON_OBJECT / JSC_IMPLEMENTS_NON_INTERFACE.
func isObjectLikeType(t jstype.JSType) bool {
	switch v := t.(type) {
	case *jstype.ObjectType, *jstype.TemplatizedType, *jstype.FunctionType:
		return true
	case *jstype.UnionType:
		// Nullable object (Foo|null or Foo|undefined) — fine.
		// Union of two-or-more non-null types — not object-like.
		nonNull := nonNullUnionMembers(v)
		return len(nonNull) == 1 && isObjectLikeType(nonNull[0])
	}
	return false
}

// nonNullUnionMembers returns the members of u that are not null or void.
func nonNullUnionMembers(u *jstype.UnionType) []jstype.JSType {
	var out []jstype.JSType
	for _, m := range u.Members() {
		if m == jstype.Null || m == jstype.Void {
			continue
		}
		out = append(out, m)
	}
	return out
}

// ifacePropInfo holds the declared type and the name of the interface
// that originally declared a property (which may differ from the
// @implements target when properties are inherited via @extends).
type ifacePropInfo struct {
	typ       jstype.JSType
	ifaceName string
}

// allInterfacePropEntries returns every property name → {type, declaring interface}
// pair reachable from `iface` through its full @extends graph.
// Mirrors TypeCheck.java's validateImplements which calls
// getCtorExtendedInterfaces recursively to collect all properties
// from multi-extends chains (e.g. Int2 @extends {Int0} @extends {Int1}).
// We use the SuperclassCheck AST lookup to enumerate all @extends
// bases rather than only the single Super() pointer (which stores
// only the first parent). The ifaceName in each entry names the
// interface that first declared the property, so error messages point
// to the right interface (e.g. "interface Int1" not "interface Int2").
func (s *SuperclassCheck) allInterfacePropEntries(iface *jstype.ObjectType) map[string]ifacePropInfo {
	out := map[string]ifacePropInfo{}
	s.collectIfacePropEntries(iface, out)
	return out
}

// collectIfacePropEntries is the recursive helper for allInterfacePropEntries.
// It registers each own member of `cur` (first-wins dedup) then recurses
// into every @extends base found in the interface's JSDoc.
func (s *SuperclassCheck) collectIfacePropEntries(cur *jstype.ObjectType, out map[string]ifacePropInfo) {
	if cur == nil {
		return
	}
	for _, name := range cur.MemberNames() {
		if _, exists := out[name]; !exists {
			out[name] = ifacePropInfo{typ: cur.GetPropertyType(name), ifaceName: cur.TypeName()}
		}
	}
	// Use the AST-based JSDoc lookup to enumerate ALL @extends bases,
	// not just the single Super() pointer (which holds only one parent).
	docInfo := s.findConstructorJSDoc(cur.TypeName())
	if docInfo != nil {
		for _, b := range docInfo.BaseTypes() {
			baseName := stripTemplateArgs(b.TypeName)
			base := s.scope.LookupNominalType(baseName)
			if base != nil {
				s.collectIfacePropEntries(base, out)
			}
		}
	} else if cur.Super() != nil {
		// Fall back to single Super() chain when no AST doc is found.
		s.collectIfacePropEntries(cur.Super(), out)
	}
}

// checkImplementedConflictsForStub detects type conflicts for an
// untyped prototype stub `X.prototype.name;` when two @implements
// interfaces declare the property with incompatible types. Java
// propagates the first interface's type to the stub; we replicate
// by comparing all @implements slots against the first.
func (s *SuperclassCheck) checkImplementedConflictsForStub(propName string, propNode *rhino.Node, ownerName string, fn *jstype.FunctionType) {
	own := fn.Instance.GetPropertyType(propName)
	if own == nil || own != jstype.Unknown {
		return // only for un-annotated stubs
	}
	ctorDoc := s.findConstructorJSDoc(ownerName)
	if ctorDoc == nil {
		return
	}
	var firstSlot jstype.JSType
	for _, expr := range ctorDoc.Implementations() {
		iface := s.scope.LookupNominalType(stripTemplateArgs(expr.TypeName))
		if iface == nil || !iface.IsStructural {
			continue
		}
		slot := iface.GetPropertyType(propName)
		if slot == nil {
			continue
		}
		if firstSlot == nil {
			firstSlot = slot
			continue
		}
		if !overrideCompatible(firstSlot, slot) {
			s.diags = append(s.diags, diag.Diagnostic{
				Path: s.path, Line: propNode.Lineno(), Col: propNode.Charno(),
				Code: diag.JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH, Level: diag.Error,
				Message: fmt.Sprintf("mismatch of the %s property on type %s and the type of the property it overrides from interface %s\noriginal: %s\noverride: %s",
					propName, ownerName, iface.TypeName(), slot.TypeName(), firstSlot.TypeName()),
			})
			return
		}
	}
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
	if len(innerKids) < 2 || innerKids[1].String_() != "prototype" {
		return
	}
	var fn *jstype.FunctionType
	var ownerName string
	switch innerKids[0].Token() {
	case rhino.NAME:
		ownerName = innerKids[0].String_()
		fn, _ = s.scope.GetRawDeclaredType(ownerName).(*jstype.FunctionType)
	case rhino.GETPROP:
		if gpK := innerKids[0].Children(); len(gpK) >= 2 && gpK[0].Token() == rhino.NAME {
			ownerName = gpK[0].String_() + "." + gpK[1].String_()
			fn, _ = s.scope.LookupNamespaceProp(gpK[0].String_(), gpK[1].String_()).(*jstype.FunctionType)
		}
	}
	if fn == nil || fn.Instance == nil {
		return
	}
	info := stmt.JSDocInfo()
	if info == nil {
		info = gp.JSDocInfo()
	}
	if info == nil || !info.HasType() {
		// No explicit type. Fire if two @implements interfaces conflict
		// on this property (untyped stub — Java propagates the first
		// interface's type; we detect conflict by comparing them).
		s.checkImplementedConflictsForStub(propName, propNode, ownerName, fn)
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
		// Detect X.prototype = {…} — object-literal prototype override.
		if propName == "prototype" && inner.Token() == rhino.NAME &&
			rhs.Token() == rhino.OBJECTLIT {
			s.checkPrototypeObjectLitAssign(inner.String_(), rhs)
		}
		return
	}
	innerKids := inner.Children()
	if len(innerKids) < 2 {
		return
	}
	if innerKids[1].String_() != "prototype" {
		return
	}
	// Resolve constructor for both simple `Foo.prototype.bar` and namespace
	// `goog.Sub.prototype.bar` patterns. Mirrors getLValueRootScope in Java.
	var fn *jstype.FunctionType
	var ownerName string
	switch innerKids[0].Token() {
	case rhino.NAME:
		ownerName = innerKids[0].String_()
		fn, _ = s.scope.GetRawDeclaredType(ownerName).(*jstype.FunctionType)
	case rhino.GETPROP:
		if gpK := innerKids[0].Children(); len(gpK) >= 2 && gpK[0].Token() == rhino.NAME {
			ownerName = gpK[0].String_() + "." + gpK[1].String_()
			fn, _ = s.scope.LookupNamespaceProp(gpK[0].String_(), gpK[1].String_()).(*jstype.FunctionType)
		}
	}
	if fn == nil || fn.Instance == nil {
		return
	}
	overrideType := s.inferAssignValue(stmt, assign, rhs, fn.Instance)
	if overrideType == nil {
		return
	}
	// Read JSDoc to detect @override — mirrors Java's declaresOverride(info).
	jsDocInfo := assign.JSDocInfo()
	if jsDocInfo == nil {
		jsDocInfo = stmt.JSDocInfo()
	}
	declaredOverride := jsDocInfo.IsOverride()
	// Mirrors Java's checkDeclaredPropertyAgainstNominalInheritance:
	// interfaces are checked elsewhere; skip them here.
	if fn.Instance.IsStructural {
		if !declaredOverride {
			return
		}
	}
	foundProperty := false
	super := fn.Instance.Super()
	impls := s.implementedInstances(ownerName)
	// Superclass slot (via @extends chain).
	if super != nil {
		if super.IsPropertyInferred(propName) {
			return
		}
		if slot := super.GetPropertyType(propName); slot != nil {
			foundProperty = true
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
	// @implements interfaces (with template instantiation) + @extends
	// parents for interface owners (multi-extends like Sub @extends {A} @extends {B}).
	// Mirrors TypeCheck.java's checkDeclaredPropertyAgainstInterfaceDeclaredInheritance,
	// which walks each @implements and applies the target's template args.
	ctorDoc := s.findConstructorJSDoc(ownerName)
	if ctorDoc != nil {
		for _, expr := range ctorDoc.Implementations() {
			iface := s.scope.LookupNominalType(stripTemplateArgs(expr.TypeName))
			if iface == nil {
				continue
			}
			slot := iface.GetPropertyType(propName)
			if slot == nil {
				continue
			}
			foundProperty = true
			slot = s.instantiateSlot(expr.TypeName, iface, slot)
			if overrideCompatible(overrideType, slot) {
				continue
			}
			s.diags = append(s.diags, diag.Diagnostic{
				Path: s.path, Line: propNode.Lineno(), Col: propNode.Charno(),
				Code: diag.JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH, Level: diag.Error,
				Message: fmt.Sprintf("mismatch of the %s property on type %s and the type of the property it overrides from interface %s\noriginal: %s\noverride: %s",
					propName, ownerName, iface.TypeName(), slot.TypeName(), overrideType.TypeName()),
			})
			return
		}
		if fn.Instance.IsStructural {
			for _, b := range ctorDoc.BaseTypes() {
				iface := s.scope.LookupNominalType(stripTemplateArgs(b.TypeName))
				if iface == nil {
					continue
				}
				slot := iface.GetPropertyType(propName)
				if slot == nil {
					continue
				}
				foundProperty = true
				if overrideCompatible(overrideType, slot) {
					continue
				}
				s.diags = append(s.diags, diag.Diagnostic{
					Path: s.path, Line: propNode.Lineno(), Col: propNode.Charno(),
					Code: diag.JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH, Level: diag.Error,
					Message: fmt.Sprintf("mismatch of the %s property on type %s and the type of the property it overrides from interface %s\noriginal: %s\noverride: %s",
						propName, ownerName, iface.TypeName(), slot.TypeName(), overrideType.TypeName()),
				})
				return
			}
		}
	}
	// Mirrors Java's checkDeclaredPropertyAgainstNominalInheritance:
	// if @override was declared but the property exists on neither the
	// resolved superclass nor any @implements interface, fire
	// JSC_UNKNOWN_OVERRIDE. Guard: fire when there is at least one
	// resolved ancestor OR when the JSDoc declares a plain (non-template)
	// @implements whose base type couldn't be resolved (i.e. truly
	// undefined, not just unresolvable due to template args).
	// Unresolved @extends (e.g. template @extends {Foo<T>}) leaves
	// super==nil with no plain @implements, so we skip — matching Java's
	// hasUnknownOrEmptySupertype guard.
	hasAncestors := super != nil || len(impls) > 0
	if !hasAncestors {
		ownerInfo := s.findConstructorJSDoc(ownerName)
		if ownerInfo != nil {
			for _, e := range ownerInfo.Implementations() {
				// Only count plain (non-template) @implements that failed
				// to resolve — that indicates a truly undefined type rather
				// than a type we can't resolve due to template arguments.
				if !strings.ContainsRune(e.TypeName, '<') {
					hasAncestors = true
					break
				}
			}
		}
	}
	if !foundProperty && declaredOverride && hasAncestors {
		s.diags = append(s.diags, diag.Diagnostic{
			Path: s.path, Line: stmt.Lineno(), Col: stmt.Charno(),
			Code: diag.JSC_UNKNOWN_OVERRIDE, Level: diag.Error,
			Message: fmt.Sprintf("property %s not defined on any superclass of %s",
				propName, fn.Instance.TypeName()),
		})
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
	of, okF := override.(*jstype.FunctionType)
	sf, okS := slot.(*jstype.FunctionType)
	if okF && okS {
		// Varargs sentinel (mirrors Java SubtypeChecker.isArrowTypeSubtype
		// lines 516-529): if slot is variadic with a specific non-Unknown
		// element type and the override has exactly the same number of
		// declared params (all explicitly typed non-Unknown), the override
		// is NOT compatible — it has a fixed arity that can't handle 0
		// args the way the slot's varargs can. When the vararg type is
		// Unknown (the Closure "top-function" sentinel, e.g.
		// `function(...?)`) any override is accepted; likewise when the
		// override has more params it likely also uses varargs.
		if sf.MaxArity == jstype.MaxArgsUnlimited && of.MaxArity != jstype.MaxArgsUnlimited &&
			len(of.Params) == len(sf.Params) &&
			len(of.Params) > 0 && len(sf.Params) > 0 {
			varargType := sf.Params[len(sf.Params)-1]
			overrideParamType := of.Params[len(of.Params)-1]
			if varargType != jstype.Unknown && overrideParamType != jstype.Unknown {
				return false
			}
		}
	}
	if jstype.IsSubtype(override, slot) {
		return true
	}
	if !okF || !okS {
		return false
	}
	// Override with fewer params is acceptable to Java so long
	// as every overlapping param has a compatible type. When the
	// slot is varargs, the override may also have MORE params (they
	// are all covered by the variadic element). Return type must
	// still be a subtype.
	if len(of.Params) > len(sf.Params) && sf.MaxArity != jstype.MaxArgsUnlimited {
		return false
	}
	limit := len(sf.Params)
	if len(of.Params) < limit {
		limit = len(of.Params)
	}
	for i := 0; i < limit; i++ {
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
		if inst := s.scope.LookupNominalType(stripTemplateArgs(e.TypeName)); inst != nil {
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
		if c.Token() == rhino.CLASS {
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
// declaration against B's slot and against each @implements interface.
// Mirrors TypeCheck.java's visitObjectOrClassLiteralKey →
// checkPropertyInheritanceOnClassMember →
// checkDeclaredPropertyAgainstNominalInheritance path.
func (s *SuperclassCheck) checkClass(cls *rhino.Node) {
	kids := cls.Children()
	if len(kids) < 3 {
		return
	}
	nameNode := kids[0]
	className := nameNode.String_()
	body := kids[2]
	preStub, ok := s.scope.GetRawDeclaredType(className).(*jstype.FunctionType)
	if !ok || preStub.Instance == nil {
		return
	}
	super := preStub.Instance.Super() // nil when there is no @extends
	ifaces := s.implementedInstances(className)
	// Nothing to check when neither superclass nor @implements are present.
	if super == nil && len(ifaces) == 0 {
		return
	}
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
			info := m.JSDocInfo()
			ownSlot := functionTypeFromJSDoc(info, m, s.scope)
			if ownSlot == nil {
				continue
			}
			// Superclass slot (via @extends chain).
			if super != nil {
				if superSlot := super.GetPropertyType(methodName); superSlot != nil {
					if !overrideCompatible(ownSlot, superSlot) {
						s.diags = append(s.diags, diag.Diagnostic{
							Path:  s.path,
							Line:  m.Lineno(),
							Col:   m.Charno(),
							Code:  diag.JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH,
							Level: diag.Error,
							Message: fmt.Sprintf("mismatch of the %s property type and the type of the property it overrides from superclass %s\noriginal: %s\noverride: %s",
								methodName, super.TypeName(), superSlot.TypeName(), ownSlot.TypeName()),
						})
						continue // one diagnostic per member
					}
				}
			}
			// @implements interface slots.
			for _, iface := range ifaces {
				ifaceSlot := iface.GetPropertyType(methodName)
				if ifaceSlot == nil {
					continue
				}
				if overrideCompatible(ownSlot, ifaceSlot) {
					continue
				}
				s.diags = append(s.diags, diag.Diagnostic{
					Path:  s.path,
					Line:  m.Lineno(),
					Col:   m.Charno(),
					Code:  diag.JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH,
					Level: diag.Error,
					Message: fmt.Sprintf("mismatch of the %s property on type %s and the type of the property it overrides from interface %s\noriginal: %s\noverride: %s",
						methodName, className, iface.TypeName(), ifaceSlot.TypeName(), ownSlot.TypeName()),
				})
				break // one diagnostic per member
			}
		case rhino.CLASS_FIELD_DEF:
			if len(m.Children()) < 1 {
				continue
			}
			fieldName := m.FirstChild().String_()
			if fieldName == "" {
				continue
			}
			info := m.JSDocInfo()
			if info == nil || !info.HasType() {
				continue
			}
			ownType := s.scope.ResolveExpr(info.Type().TypeName)
			// Superclass slot.
			if super != nil {
				if superSlot := super.GetPropertyType(fieldName); superSlot != nil {
					if !jstype.IsSubtype(ownType, superSlot) {
						s.diags = append(s.diags, diag.Diagnostic{
							Path:  s.path,
							Line:  m.Lineno(),
							Col:   m.Charno(),
							Code:  diag.JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH,
							Level: diag.Error,
							Message: fmt.Sprintf("mismatch of the %s property type and the type of the property it overrides from superclass %s\noriginal: %s\noverride: %s",
								fieldName, super.TypeName(), superSlot.TypeName(), ownType.TypeName()),
						})
						continue
					}
				}
			}
			// @implements interface slots.
			for _, iface := range ifaces {
				ifaceSlot := iface.GetPropertyType(fieldName)
				if ifaceSlot == nil {
					continue
				}
				if jstype.IsSubtype(ownType, ifaceSlot) {
					continue
				}
				s.diags = append(s.diags, diag.Diagnostic{
					Path:  s.path,
					Line:  m.Lineno(),
					Col:   m.Charno(),
					Code:  diag.JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH,
					Level: diag.Error,
					Message: fmt.Sprintf("mismatch of the %s property on type %s and the type of the property it overrides from interface %s\noriginal: %s\noverride: %s",
						fieldName, className, iface.TypeName(), ifaceSlot.TypeName(), ownType.TypeName()),
				})
				break
			}
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

// instantiateSlot substitutes template args from typeExpr (e.g. "Int<string>")
// into slot, mapping each of iface's TemplateTypes to the corresponding arg.
// Returns slot unchanged when no template args are present.
func (s *SuperclassCheck) instantiateSlot(typeExpr string, iface *jstype.ObjectType, slot jstype.JSType) jstype.JSType {
	i := strings.IndexByte(typeExpr, '<')
	if i < 0 || slot == nil {
		return slot
	}
	fn, ok := s.scope.GetRawDeclaredType(iface.TypeName()).(*jstype.FunctionType)
	if !ok || len(fn.TemplateTypes) == 0 {
		return slot
	}
	inner := strings.TrimSpace(typeExpr[i+1 : len(typeExpr)-1])
	argStrs := strings.Split(inner, ",")
	bindings := map[*jstype.TemplateType]jstype.JSType{}
	for j, tt := range fn.TemplateTypes {
		if j >= len(argStrs) {
			break
		}
		bindings[tt] = s.scope.ResolveExpr(strings.TrimSpace(argStrs[j]))
	}
	if len(bindings) == 0 {
		return slot
	}
	r := jstype.ForInference(s.scope.Registry(), bindings)
	return jstype.ApplyTo(r, slot)
}

// checkPrototypeObjectLitAssign handles `X.prototype = {key: fn, …}` —
// fires JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH when a key's type
// conflicts with a slot inherited from a super-interface.
func (s *SuperclassCheck) checkPrototypeObjectLitAssign(ownerName string, lit *rhino.Node) {
	fn, ok := s.scope.GetRawDeclaredType(ownerName).(*jstype.FunctionType)
	if !ok || fn.Instance == nil {
		return
	}
	super := fn.Instance.Super()
	if super == nil {
		return
	}
	for _, kv := range lit.Children() {
		if kv.Token() != rhino.STRING_KEY {
			continue
		}
		child := kv.FirstChild()
		if child == nil || child.Token() != rhino.FUNCTION {
			continue
		}
		methodName := kv.String_()
		ownSlot := functionTypeFromJSDoc(kv.JSDocInfo(), child, s.scope)
		if ownSlot == nil {
			continue
		}
		superSlot := super.GetPropertyType(methodName)
		if superSlot == nil || overrideCompatible(ownSlot, superSlot) {
			continue
		}
		if fn.Instance.IsStructural || super.IsStructural {
			s.diags = append(s.diags, diag.Diagnostic{
				Path: s.path, Line: kv.Lineno(), Col: kv.Charno(),
				Code: diag.JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH, Level: diag.Error,
				Message: fmt.Sprintf("mismatch of the %s property on type %s and the type of the property it overrides from interface %s\noriginal: %s\noverride: %s",
					methodName, ownerName, super.TypeName(), superSlot.TypeName(), ownSlot.TypeName()),
			})
		} else {
			s.diags = append(s.diags, diag.Diagnostic{
				Path: s.path, Line: kv.Lineno(), Col: kv.Charno(),
				Code: diag.JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH, Level: diag.Error,
				Message: fmt.Sprintf("mismatch of the %s property type and the type of the property it overrides from superclass %s\noriginal: %s\noverride: %s",
					methodName, super.TypeName(), superSlot.TypeName(), ownSlot.TypeName()),
			})
		}
	}
}
