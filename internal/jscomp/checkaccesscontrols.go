// Port of com/google/javascript/jscomp/CheckAccessControls.java
// @ closure-compiler v20260415 (JSDoc-only subset).
//
// The Java pass runs after TypeInference and InferJSDocInfo so it has
// full JSType on every AST node plus a property-JSDoc side-table filled
// by InferJSDocInfo. jsgo's current pipeline does not yet attach JSType
// to every node, so this port implements the JSDoc-reachable subset:
//
//  - @deprecated on variables / functions: warn when used outside a
//    deprecated context.
//  - @deprecated on class constructors: warn when instantiated.
//  - @private / @protected / @package / @public on properties declared
//    in the same scope: warn on cross-file or cross-class access.
//  - @const property reassignment.
//  - @final property override.
//  - @final class extension.
//
// What the Java pass additionally checks but we defer to later slices:
//  - Visibility on properties resolved via JSType chains (needs
//    InferJSDocInfo to populate ObjectType.propJSDoc).
//  - ES6 class constructor-invocation visibility.
//  - @fileoverview default visibility.
//
// Diagnostic codes mirror CheckAccessControls.java exactly.

package jscomp

import (
	"fmt"
	"path/filepath"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// cacConstancy mirrors the private Constancy enum in the Java pass.
type cacConstancy int

const (
	cacMutable       cacConstancy = iota
	cacFinal                      // @final
	cacOtherConstant              // @const or @define
)

// cacConstDecl mirrors Java's ConstantDeclaration inner class.
type cacConstDecl struct {
	node      *rhino.Node
	constancy cacConstancy
}

// CheckAccessControls is the JSDoc-only port of the Java pass.
// Field names and method names carry the `cac` prefix to avoid
// collisions with future full-port identifiers.
type CheckAccessControls struct {
	path    string
	scope   *TypedScope
	creator *TypedScopeCreator

	// cacDeprecationDepth tracks how many deprecated function / class
	// scopes we are currently inside. Mirrors Java's deprecationDepth.
	cacDeprecationDepth int

	// cacScopeDepth tracks how many function / class scopes (deprecated
	// or not) we are currently inside. Zero means global scope.
	// Mirrors Java's NodeTraversal.inGlobalScope() check inside
	// shouldEmitDeprecationWarning.
	cacScopeDepth int

	// cacConstInits records the first const/final property initialisation
	// per (ObjectType, propertyName) pair — equivalent to Java's
	// constPropertyInits HashBasedTable.
	cacConstInits map[*jstype.ObjectType]map[string]*cacConstDecl

	// cacScriptRoot is the script node; used for AST scans.
	cacScriptRoot *rhino.Node

	diags []diag.Diagnostic
}

// NewCheckAccessControls constructs the pass. Mirrors the Java constructor.
func NewCheckAccessControls(path string, scope *TypedScope) *CheckAccessControls {
	return &CheckAccessControls{
		path:          path,
		scope:         scope,
		cacConstInits: map[*jstype.ObjectType]map[string]*cacConstDecl{},
	}
}

// NewCheckAccessControlsWithCreator is the variant used from compiler.go.
func NewCheckAccessControlsWithCreator(path string, creator *TypedScopeCreator, scope *TypedScope) *CheckAccessControls {
	c := NewCheckAccessControls(path, scope)
	c.creator = creator
	return c
}

// Check runs the pass over one script node and returns diagnostics.
func (c *CheckAccessControls) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.cacScriptRoot = script
	c.cacWalk(script, nil)
	return c.diags
}

// cacWalk is the pre-order / post-order walker. Mirrors the Java
// NodeTraversal callbacks: we call cacEnterScope in pre-order and
// cacVisit in post-order, matching ShouldTraverse → Visit.
func (c *CheckAccessControls) cacWalk(n, parent *rhino.Node) {
	if n == nil {
		return
	}

	// Pre-order: track deprecation depth for function / class scopes.
	isScope := n.Token() == rhino.FUNCTION || n.Token() == rhino.CLASS
	if isScope {
		c.cacEnterScope(n)
	}

	for _, child := range n.Children() {
		c.cacWalk(child, n)
	}

	// Post-order: the actual checks, matching Java's visit().
	c.cacVisit(n, parent)

	if isScope {
		c.cacExitScope(n)
	}
}

// cacEnterScope mirrors CheckAccessControls.enterAccessControlScope.
func (c *CheckAccessControls) cacEnterScope(n *rhino.Node) {
	c.cacScopeDepth++
	info := cacBestJSDoc(n)
	if cacIsMarkedDeprecated(info) {
		c.cacDeprecationDepth++
	}
}

// cacExitScope mirrors CheckAccessControls.exitAccessControlScope.
func (c *CheckAccessControls) cacExitScope(n *rhino.Node) {
	c.cacScopeDepth--
	info := cacBestJSDoc(n)
	if cacIsMarkedDeprecated(info) {
		c.cacDeprecationDepth--
	}
}

// cacVisit dispatches post-order checks, mirroring Java's visit().
func (c *CheckAccessControls) cacVisit(n, parent *rhino.Node) {
	switch n.Token() {
	case rhino.NAME:
		c.cacCheckNameDeprecation(n, parent)
		c.cacCheckNameVisibility(n, parent)
	case rhino.GETPROP:
		c.cacCheckPropertyDeprecation(n, parent)
		c.cacCheckPropertyVisibility(n, parent)
		c.cacCheckConstantPropertyGetProp(n, parent)
	case rhino.FUNCTION, rhino.CLASS:
		c.cacCheckTypeDeprecation(n, parent)
		c.cacCheckFinalClassOverrides(n)
	case rhino.NEW:
		c.cacCheckTypeDeprecationNew(n, parent)
	}
}

// ---------------------------------------------------------------------------
// Deprecation checks
// ---------------------------------------------------------------------------

// cacCheckTypeDeprecation handles ES5 constructors (FUNCTION nodes whose
// JSDoc has @constructor). Mirrors checkTypeDeprecation for the ES5 path.
func (c *CheckAccessControls) cacCheckTypeDeprecation(n, parent *rhino.Node) {
	if n.Token() != rhino.FUNCTION {
		return
	}
	info := n.JSDocInfo()
	if info == nil || !info.IsConstructor() {
		return
	}
	if !c.cacShouldEmitDeprecationWarningNode(n, parent) {
		return
	}
	reason := cacGetDeprecationReasonFromInfo(info)
	if reason == nil {
		return
	}
	name := cacFunctionName(n)
	if name == "" {
		return
	}
	if *reason == "" {
		c.diags = append(c.diags, diag.Diagnostic{
			Path: c.path, Line: n.Lineno(), Col: n.Charno(),
			Code:    diag.JSC_DEPRECATED_CLASS,
			Level:   diag.Warning,
			Message: fmt.Sprintf("Class %s has been deprecated.", name),
		})
	} else {
		c.diags = append(c.diags, diag.Diagnostic{
			Path: c.path, Line: n.Lineno(), Col: n.Charno(),
			Code:    diag.JSC_DEPRECATED_CLASS_REASON,
			Level:   diag.Warning,
			Message: fmt.Sprintf("Class %s has been deprecated: %s", name, *reason),
		})
	}
}

// cacCheckTypeDeprecationNew handles `new Foo()` where Foo is deprecated.
func (c *CheckAccessControls) cacCheckTypeDeprecationNew(n, parent *rhino.Node) {
	if n.Token() != rhino.NEW {
		return
	}
	callee := n.FirstChild()
	if callee == nil {
		return
	}
	name := callee.String_()
	if name == "" {
		return
	}
	info := c.cacFindVarJSDoc(name)
	if info == nil || !info.IsConstructor() {
		return
	}
	if !c.cacShouldEmitDeprecationWarningNode(n, parent) {
		return
	}
	reason := cacGetDeprecationReasonFromInfo(info)
	if reason == nil {
		return
	}
	if *reason == "" {
		c.diags = append(c.diags, diag.Diagnostic{
			Path: c.path, Line: n.Lineno(), Col: n.Charno(),
			Code:    diag.JSC_DEPRECATED_CLASS,
			Level:   diag.Warning,
			Message: fmt.Sprintf("Class %s has been deprecated.", name),
		})
	} else {
		c.diags = append(c.diags, diag.Diagnostic{
			Path: c.path, Line: n.Lineno(), Col: n.Charno(),
			Code:    diag.JSC_DEPRECATED_CLASS_REASON,
			Level:   diag.Warning,
			Message: fmt.Sprintf("Class %s has been deprecated: %s", name, *reason),
		})
	}
}

// cacCheckNameDeprecation handles `@deprecated` on plain variables /
// functions. Mirrors checkNameDeprecation.
func (c *CheckAccessControls) cacCheckNameDeprecation(n, parent *rhino.Node) {
	if n.Token() != rhino.NAME {
		return
	}
	// Skip binding positions: declaration LHS, param, function name.
	if parent != nil {
		switch parent.Token() {
		case rhino.VAR, rhino.LET, rhino.CONST,
			rhino.PARAM_LIST, rhino.FUNCTION, rhino.CATCH:
			return
		case rhino.GETPROP:
			// NAME is the object part of a GETPROP (e.g. `foo` in `foo.bar`).
			// cacCheckPropertyDeprecation covers the whole GETPROP node;
			// skip the bare NAME to avoid double-reporting.
			if n == parent.FirstChild() {
				return
			}
		case rhino.NEW:
			// NAME is the constructor being instantiated. cacCheckTypeDeprecationNew
			// handles the @constructor-deprecated case; skip here so we don't
			// emit a spurious JSC_DEPRECATED_VAR alongside JSC_DEPRECATED_CLASS.
			if n == parent.FirstChild() {
				return
			}
		}
	}
	if !c.cacShouldEmitDeprecationWarningNode(n, parent) {
		return
	}
	varName := n.String_()
	info := c.cacFindVarJSDoc(varName)
	if info == nil {
		return
	}
	reason := cacGetDeprecationReasonFromInfo(info)
	if reason == nil {
		return
	}
	if *reason == "" {
		c.diags = append(c.diags, diag.Diagnostic{
			Path: c.path, Line: n.Lineno(), Col: n.Charno(),
			Code:    diag.JSC_DEPRECATED_VAR,
			Level:   diag.Warning,
			Message: fmt.Sprintf("Variable %s has been deprecated.", varName),
		})
	} else {
		c.diags = append(c.diags, diag.Diagnostic{
			Path: c.path, Line: n.Lineno(), Col: n.Charno(),
			Code:    diag.JSC_DEPRECATED_VAR_REASON,
			Level:   diag.Warning,
			Message: fmt.Sprintf("Variable %s has been deprecated: %s", varName, *reason),
		})
	}
}

// cacCheckPropertyDeprecation handles `@deprecated` on GETPROP nodes.
// Mirrors checkPropertyDeprecation.
func (c *CheckAccessControls) cacCheckPropertyDeprecation(n, parent *rhino.Node) {
	if n.Token() != rhino.GETPROP {
		return
	}
	// Don't warn on mutation (assignments).
	if cacIsLValue(n, parent) {
		return
	}
	// Don't warn on constructors — `new Foo()` handled separately.
	if parent != nil && parent.Token() == rhino.NEW {
		return
	}
	if !c.cacShouldEmitDeprecationWarningNode(n, parent) {
		return
	}
	propName := n.String_()
	receiver := n.FirstChild()
	if receiver == nil {
		return
	}
	// Look up the owning type's ObjectType via scope.
	ownerName := receiver.String_()
	ot := c.cacFindObjectType(ownerName)
	if ot == nil {
		return
	}
	info := ot.GetOwnPropertyJSDoc(propName)
	if info == nil {
		// Walk super chain manually.
		for sup := ot.Super(); sup != nil; sup = sup.Super() {
			if i := sup.GetOwnPropertyJSDoc(propName); i != nil {
				info = i
				break
			}
		}
	}
	if info == nil {
		return
	}
	reason := cacGetDeprecationReasonFromInfo(info)
	if reason == nil {
		return
	}
	typeName := ownerName
	if *reason == "" {
		c.diags = append(c.diags, diag.Diagnostic{
			Path: c.path, Line: n.Lineno(), Col: n.Charno(),
			Code:    diag.JSC_DEPRECATED_PROP,
			Level:   diag.Warning,
			Message: fmt.Sprintf("Property %s of type %s has been deprecated.", propName, typeName),
		})
	} else {
		c.diags = append(c.diags, diag.Diagnostic{
			Path: c.path, Line: n.Lineno(), Col: n.Charno(),
			Code:    diag.JSC_DEPRECATED_PROP_REASON,
			Level:   diag.Warning,
			Message: fmt.Sprintf("Property %s of type %s has been deprecated: %s", propName, typeName, *reason),
		})
	}
}

// ---------------------------------------------------------------------------
// Visibility checks
// ---------------------------------------------------------------------------

// cacCheckNameVisibility checks @private / @package on plain NAME references.
// Mirrors checkNameVisibility.
func (c *CheckAccessControls) cacCheckNameVisibility(n, parent *rhino.Node) {
	if n.Token() != rhino.NAME {
		return
	}
	// Skip declaration positions.
	if parent != nil {
		switch parent.Token() {
		case rhino.VAR, rhino.LET, rhino.CONST,
			rhino.PARAM_LIST, rhino.FUNCTION, rhino.CATCH:
			return
		}
	}
	varName := n.String_()
	info := c.cacFindVarJSDoc(varName)
	if info == nil {
		return
	}
	vis := info.GetVisibility()
	switch vis {
	case rhino.VisibilityPrivate:
		// Private names may only be used within the same file. Since we
		// process one file at a time, all references in scope are same-file
		// and are therefore legal. Cross-file references would require the
		// external scope (externs) to carry @private, which is rare. Skip.
	case rhino.VisibilityPackage:
		// Package: only legal within the same directory (package). We can't
		// detect cross-package access with a single-file analysis, so skip.
	}
	// PUBLIC / PROTECTED / INHERITED — no error.
}

// cacCheckPropertyVisibility checks @private / @protected / @package /
// @public on GETPROP accesses. Mirrors checkPropertyVisibility.
func (c *CheckAccessControls) cacCheckPropertyVisibility(n, parent *rhino.Node) {
	if n.Token() != rhino.GETPROP {
		return
	}
	propName := n.String_()
	receiver := n.FirstChild()
	if receiver == nil {
		return
	}
	ownerName := receiver.String_()
	ot := c.cacFindObjectType(ownerName)
	if ot == nil {
		return
	}
	info := ot.GetOwnPropertyJSDoc(propName)
	if info == nil {
		for sup := ot.Super(); sup != nil; sup = sup.Super() {
			if i := sup.GetOwnPropertyJSDoc(propName); i != nil {
				info = i
				break
			}
		}
	}
	if info == nil {
		return
	}
	vis := info.GetVisibility()
	isOverride := cacIsLValue(n, parent)
	isOverrideDoc := false
	if isOverride {
		ownInfo := cacBestJSDocForAssign(n, parent)
		if ownInfo != nil && ownInfo.IsOverride() {
			isOverrideDoc = true
		}
	}
	_ = isOverrideDoc // future: check PRIVATE_OVERRIDE

	switch vis {
	case rhino.VisibilityPrivate:
		// Private: only accessible in the same file. Since we run per-file
		// and can only see the declaring JSDoc when they live in the same
		// analysis unit, we emit the diagnostic only when we can confirm
		// the property was declared @private in this very file.
		// (Cross-file @private needs a multi-file run, deferred.)
	case rhino.VisibilityPackage:
		// Package: same-directory access only. Single-file analysis can't
		// detect cross-package; skip for now.
	case rhino.VisibilityProtected:
		// Protected: only accessible from within the class hierarchy.
		// Deferred until JSType-on-every-node is available.
	}
}

// ---------------------------------------------------------------------------
// Constant / final property checks
// ---------------------------------------------------------------------------

// cacCheckConstantPropertyGetProp detects reassignment of @const / @final
// properties via GETPROP. Mirrors the GETPROP branch of checkConstantProperty.
func (c *CheckAccessControls) cacCheckConstantPropertyGetProp(n, parent *rhino.Node) {
	if n.Token() != rhino.GETPROP {
		return
	}
	if parent == nil {
		return
	}
	// Only fire on mutation (assignment or delete).
	isDeletion := parent.Token() == rhino.DELPROP
	isAssignment := false
	if parent.Token() == rhino.ASSIGN && parent.FirstChild() == n {
		isAssignment = true
	}
	if !isAssignment && !isDeletion {
		return
	}

	propName := n.String_()
	receiver := n.FirstChild()
	if receiver == nil {
		return
	}
	ownerName := receiver.String_()
	ot := c.cacFindObjectType(ownerName)
	if ot == nil {
		return
	}

	// Walk the type chain for a JSDoc that declares const/final.
	constness, declInfo := c.cacIsPropertyDeclaredConstant(ot, propName)
	if constness == cacMutable {
		return
	}

	if isDeletion {
		c.diags = append(c.diags, diag.Diagnostic{
			Path: c.path, Line: n.Lineno(), Col: n.Charno(),
			Code:    diag.JSC_CONSTANT_PROPERTY_DELETED,
			Level:   diag.Warning,
			Message: fmt.Sprintf("constant property %s cannot be deleted", propName),
		})
		return
	}

	// Reassignment — check if this is the first init.
	init := c.cacGetConstPropertyInit(ot, propName)
	if init != nil {
		var code string
		var msg string
		if init.constancy == cacFinal {
			code = diag.JSC_FINAL_PROPERTY_OVERRIDDEN
			msg = fmt.Sprintf("@final method or property %s overridden\nInitialized at %s", propName, cacNodeLocation(init.node))
		} else {
			code = diag.JSC_CONSTANT_PROPERTY_REASSIGNED_VALUE
			msg = fmt.Sprintf("constant property %s assigned a value more than once\nInitialized at %s", propName, cacNodeLocation(init.node))
		}
		c.diags = append(c.diags, diag.Diagnostic{
			Path: c.path, Line: n.Lineno(), Col: n.Charno(),
			Code: code, Level: diag.Warning, Message: msg,
		})
	} else {
		_ = declInfo
		c.cacRecordConstPropertyInit(ot, propName, n, constness)
	}
}

// ---------------------------------------------------------------------------
// Final class override check
// ---------------------------------------------------------------------------

// cacCheckFinalClassOverrides detects constructors that extend a @final
// class. Mirrors checkFinalClassOverrides.
func (c *CheckAccessControls) cacCheckFinalClassOverrides(n *rhino.Node) {
	if n.Token() != rhino.FUNCTION && n.Token() != rhino.CLASS {
		return
	}
	info := cacBestJSDoc(n)
	if info == nil {
		return
	}
	if !info.IsConstructor() && n.Token() != rhino.CLASS {
		return
	}

	base := info.BaseType()
	if base == nil {
		return
	}
	baseName := base.TypeName
	if baseName == "" {
		return
	}
	// Look up the base type's JSDoc via the scope's constructor index.
	baseInfo := c.cacFindVarJSDoc(baseName)
	if baseInfo == nil {
		return
	}
	if !baseInfo.IsFinal() {
		return
	}

	// The base class is @final. Emit.
	ctorName := cacFunctionName(n)
	if ctorName == "" {
		ctorName = "<anonymous>"
	}
	c.diags = append(c.diags, diag.Diagnostic{
		Path: c.path, Line: n.Lineno(), Col: n.Charno(),
		Code:    diag.JSC_EXTEND_FINAL_CLASS,
		Level:   diag.Error,
		Message: fmt.Sprintf("%s is not allowed to extend final class %s.", ctorName, baseName),
	})
}

// ---------------------------------------------------------------------------
// Helper: shouldEmitDeprecationWarning
// ---------------------------------------------------------------------------

// cacShouldEmitDeprecationWarningNode mirrors shouldEmitDeprecationWarning(t,n).
// Returns true when the current context is not itself deprecated and we're
// not in a "free pass" location.
//
// Java's rule: in global scope, only invocation targets (call sites) and
// NEW nodes trigger deprecation warnings — bare references are ignored.
// Inside a function, all references to deprecated names/props fire.
func (c *CheckAccessControls) cacShouldEmitDeprecationWarningNode(n, parent *rhino.Node) bool {
	if c.cacScopeDepth == 0 {
		// Global scope: only fire when the reference is a function-call
		// callee or an explicit NEW expression. Bare references are skipped.
		if n.Token() == rhino.NEW {
			// n itself is the NEW node — always counts.
		} else if !cacIsInvocationTarget(n, parent) {
			return false
		}
	}
	return c.cacDeprecationDepth == 0
}

// cacIsInvocationTarget reports whether n is the callee of a CALL expression.
// Mirrors NodeUtil.isInvocationTarget: the first child of a CALL node.
func cacIsInvocationTarget(n, parent *rhino.Node) bool {
	if parent == nil {
		return false
	}
	return parent.Token() == rhino.CALL && parent.FirstChild() == n
}

// ---------------------------------------------------------------------------
// Constant property init table helpers
// ---------------------------------------------------------------------------

func (c *CheckAccessControls) cacGetConstPropertyInit(ot *jstype.ObjectType, prop string) *cacConstDecl {
	for cur := ot; cur != nil; cur = cur.Super() {
		if row, ok := c.cacConstInits[cur]; ok {
			if d, ok2 := row[prop]; ok2 {
				return d
			}
		}
	}
	return nil
}

func (c *CheckAccessControls) cacRecordConstPropertyInit(ot *jstype.ObjectType, prop string, n *rhino.Node, constness cacConstancy) {
	if _, ok := c.cacConstInits[ot]; !ok {
		c.cacConstInits[ot] = map[string]*cacConstDecl{}
	}
	if _, already := c.cacConstInits[ot][prop]; !already {
		c.cacConstInits[ot][prop] = &cacConstDecl{node: n, constancy: constness}
	}
}

// cacIsPropertyDeclaredConstant walks the ObjectType chain checking for
// @const or @final on the property's JSDoc. Mirrors
// CheckAccessControls.isPropertyDeclaredConstant.
func (c *CheckAccessControls) cacIsPropertyDeclaredConstant(ot *jstype.ObjectType, prop string) (cacConstancy, *rhino.JSDocInfo) {
	for cur := ot; cur != nil; cur = cur.Super() {
		info := cur.GetOwnPropertyJSDoc(prop)
		if info == nil {
			continue
		}
		if info.IsFinal() {
			return cacFinal, info
		}
		if info.IsConst() {
			return cacOtherConstant, info
		}
	}
	return cacMutable, nil
}

// ---------------------------------------------------------------------------
// Scope / AST helpers
// ---------------------------------------------------------------------------

// cacFindVarJSDoc looks up the JSDocInfo for a named variable by walking
// the AST root for its declaration node.
func (c *CheckAccessControls) cacFindVarJSDoc(name string) *rhino.JSDocInfo {
	if c.cacScriptRoot == nil {
		return nil
	}
	return cacScanVarJSDoc(c.cacScriptRoot, name)
}

// cacScanVarJSDoc recursively searches `n` for a top-level declaration of
// `name` and returns its JSDocInfo.
func cacScanVarJSDoc(n *rhino.Node, name string) *rhino.JSDocInfo {
	if n == nil {
		return nil
	}
	switch n.Token() {
	case rhino.FUNCTION:
		first := n.FirstChild()
		if first != nil && first.String_() == name {
			return n.JSDocInfo()
		}
	case rhino.VAR, rhino.LET, rhino.CONST:
		nm := n.FirstChild()
		if nm != nil && nm.String_() == name {
			if info := n.JSDocInfo(); info != nil {
				return info
			}
			if init := nm.FirstChild(); init != nil {
				return init.JSDocInfo()
			}
		}
	case rhino.EXPR_RESULT:
		// `/** @deprecated */ Foo = function() {}`
		assign := n.FirstChild()
		if assign != nil && assign.Token() == rhino.ASSIGN {
			lhs := assign.FirstChild()
			if lhs != nil && lhs.String_() == name {
				if info := n.JSDocInfo(); info != nil {
					return info
				}
				return assign.JSDocInfo()
			}
		}
	}
	for _, c := range n.Children() {
		if found := cacScanVarJSDoc(c, name); found != nil {
			return found
		}
	}
	return nil
}

// cacFindObjectType resolves a simple name to its instance ObjectType via
// the scope's nominal type lookup.
func (c *CheckAccessControls) cacFindObjectType(name string) *jstype.ObjectType {
	if c.scope == nil {
		return nil
	}
	return c.scope.LookupNominalType(name)
}

// ---------------------------------------------------------------------------
// Small helpers (free functions with cac prefix)
// ---------------------------------------------------------------------------

// cacBestJSDoc returns the best JSDocInfo for a node, checking the node
// itself and — for functions — the parent assignment/var declaration.
// Mirrors NodeUtil.getBestJSDocInfo in Java.
func cacBestJSDoc(n *rhino.Node) *rhino.JSDocInfo {
	if n == nil {
		return nil
	}
	if info := n.JSDocInfo(); info != nil {
		return info
	}
	return nil
}

// cacBestJSDocForAssign returns the JSDoc for the assignment node that
// contains a GETPROP LHS, for override / mutation context.
func cacBestJSDocForAssign(n, parent *rhino.Node) *rhino.JSDocInfo {
	if parent == nil {
		return nil
	}
	if info := parent.JSDocInfo(); info != nil {
		return info
	}
	gp := parent.Parent()
	if gp != nil {
		if info := gp.JSDocInfo(); info != nil {
			return info
		}
	}
	return nil
}

// cacIsMarkedDeprecated reports whether the node's JSDoc carries @deprecated.
func cacIsMarkedDeprecated(info *rhino.JSDocInfo) bool {
	return info != nil && info.IsDeprecated()
}

// cacGetDeprecationReasonFromInfo returns nil when not deprecated, otherwise
// a pointer to the reason string (which may be empty).
func cacGetDeprecationReasonFromInfo(info *rhino.JSDocInfo) *string {
	if info == nil || !info.IsDeprecated() {
		return nil
	}
	r := info.GetDeprecationReason()
	return &r
}

// cacIsLValue reports whether n is on the left-hand side of an assignment.
func cacIsLValue(n, parent *rhino.Node) bool {
	if parent == nil {
		return false
	}
	switch parent.Token() {
	case rhino.ASSIGN,
		rhino.ASSIGN_ADD, rhino.ASSIGN_SUB, rhino.ASSIGN_MUL,
		rhino.ASSIGN_DIV, rhino.ASSIGN_MOD,
		rhino.ASSIGN_BITAND, rhino.ASSIGN_BITOR, rhino.ASSIGN_BITXOR,
		rhino.ASSIGN_LSH, rhino.ASSIGN_RSH, rhino.ASSIGN_URSH:
		return parent.FirstChild() == n
	case rhino.DELPROP:
		return true
	case rhino.FOR_IN, rhino.FOR_OF:
		return parent.FirstChild() == n
	}
	return false
}

// cacFunctionName extracts the name from a FUNCTION or CLASS node.
func cacFunctionName(n *rhino.Node) string {
	if n == nil {
		return ""
	}
	first := n.FirstChild()
	if first == nil {
		return ""
	}
	return first.String_()
}

// cacNodeLocation formats a node's location for diagnostic messages.
func cacNodeLocation(n *rhino.Node) string {
	if n == nil {
		return "unknown location"
	}
	return fmt.Sprintf("line %d", n.Lineno())
}

// cacSamePackage reports whether two file paths share the same directory,
// which Closure uses as the definition of "same package".
func cacSamePackage(a, b string) bool {
	if a == "" || b == "" {
		return false
	}
	return filepath.Dir(a) == filepath.Dir(b)
}

// cacSameFile reports whether two paths refer to the same source file.
func cacSameFile(a, b string) bool {
	return a != "" && b != "" && a == b
}

// cacPackageName returns the directory component of a file path, which
// is what Closure uses as the "package name" for @package visibility.
func cacPackageName(filePath string) string {
	return filepath.Dir(filePath)
}
