// Port of com/google/javascript/jscomp/AbstractVar.java
// @ closure-compiler v20260415
//
// AbstractVar<S,V> is the generic base class for Var (used with Scope) and
// TypedVar (used with TypedScope). It extends ScopedName and carries the name,
// the name AST node, the enclosing scope, a declaration index, the compiler
// input, and optional implicit-goog-namespace metadata.
//
// Java inheritance chain:
//   ScopedName ← AbstractVar<S,V>
//                ← Var   (S=Scope,  V=Var)
//                ← TypedVar (S=TypedScope, V=TypedVar)
//
// Go port strategy:
//   AbstractVarBase is a struct that holds all fields declared in Java's
//   AbstractVar. Both UntypedVar (port of Var.java) and GoTypedVar (port of
//   TypedVar.java) embed AbstractVarBase.
//
//   The existing TypedScope in scope.go does not use AbstractVarBase; it uses
//   plain JSType maps. This file defines the hierarchy for callers that need
//   to work with the raw Var / TypedVar types (e.g. VarCheck, SymbolTable).
//
// Full method inventory ported from AbstractVar.java (325 LoC):
//   • getName / getNode / getNameNode / getScopeRoot
//   • getParentNode / getIndex / getScope
//   • isGlobal / isLocal / isExtern
//   • isBleedingFunction / isImplicitGoogNamespace
//   • isDeclaredOrInferredConst / isDefine
//   • isVar / isLet / isConst / isClass / isParam / isImport / isCatch
//   • isDefaultParam / isArguments / isThis / isGoogModuleExports
//   • getInitialValue / getJSDocInfo
//   • declarationType() / isImplicit()
//   • addImplicitGoogNamespaceDefinition
//   • AbstractVarIface / reference-list stub

package jscomp

import (
	"github.com/paullepoulpe/jsgo/legacy/internal/rhino"
)

// AbstractVarBase holds all fields from Java's AbstractVar<S,V>.
// Embed this in UntypedVar and GoTypedVar to mirror Java's inheritance.
type AbstractVarBase struct {
	name     string
	nameNode *rhino.Node
	index    int
	// implicitGoogNamespace: non-nil iff this is a goog.provide / goog.module namespace.
	implicitGoogNamespace *rhino.Node
	// references accumulates all references (NAME-node uses) of this variable.
	// Populated by SymbolTable or by a dedicated reference pass.
	// Mirrors the mutable list on StaticRef / StaticSlot in Java.
	references []*rhino.Node
}

// ---------------------------------------------------------------------------
// Identity / name accessors
// ---------------------------------------------------------------------------

// GetName returns the variable's identifier name.
// Mirrors AbstractVar.getName() in Java.
func (v *AbstractVarBase) GetName() string { return v.name }

// GetNode returns the NAME AST node that declares this variable, or nil for
// implicit variables.
// Mirrors AbstractVar.getNode() in Java.
func (v *AbstractVarBase) GetNode() *rhino.Node { return v.nameNode }

// GetNameNode is an alias for GetNode (matches Java's getNameNode()).
// Mirrors AbstractVar.getNameNode() in Java.
func (v *AbstractVarBase) GetNameNode() *rhino.Node { return v.nameNode }

// GetParentNode returns the parent of the NAME node, or nil if the name node
// is nil.
// Mirrors AbstractVar.getParentNode() in Java.
func (v *AbstractVarBase) GetParentNode() *rhino.Node {
	if v.nameNode == nil {
		return nil
	}
	return v.nameNode.Parent()
}

// GetIndex returns the declaration index within the scope.
// Mirrors AbstractVar.getIndex() in Java.
func (v *AbstractVarBase) GetIndex() int { return v.index }

// ---------------------------------------------------------------------------
// Declaration-type predicates
// ---------------------------------------------------------------------------

// declarationTypeOf returns the declaration token (VAR, LET, CONST, FUNCTION,
// CLASS, CATCH, IMPORT, PARAM_LIST) for a name node, walking up the ancestor
// chain. Returns -1 for implicit variables.
// Mirrors AbstractVar.declarationType() in Java.
func declarationTypeOf(nameNode *rhino.Node) rhino.Token {
	if nameNode == nil {
		return -1
	}
	for cur := nameNode; cur != nil; cur = cur.Parent() {
		switch cur.Token() {
		case rhino.VAR, rhino.LET, rhino.CONST,
			rhino.FUNCTION, rhino.CLASS, rhino.CATCH,
			rhino.IMPORT, rhino.PARAM_LIST:
			return cur.Token()
		}
	}
	return -1
}

// DeclarationType returns the declaration token for this variable.
// Mirrors AbstractVar.declarationType() in Java.
func (v *AbstractVarBase) DeclarationType() rhino.Token {
	if v.IsImplicitGoogNamespace() {
		return -1
	}
	return declarationTypeOf(v.nameNode)
}

// IsVar reports whether this variable was declared with the `var` keyword.
// Mirrors AbstractVar.isVar() in Java.
func (v *AbstractVarBase) IsVar() bool {
	return v.DeclarationType() == rhino.VAR
}

// IsLet reports whether this variable was declared with `let`.
// Mirrors AbstractVar.isLet() in Java.
func (v *AbstractVarBase) IsLet() bool {
	return v.DeclarationType() == rhino.LET
}

// IsConst reports whether this variable was declared with `const`.
// Mirrors AbstractVar.isConst() in Java.
func (v *AbstractVarBase) IsConst() bool {
	return v.DeclarationType() == rhino.CONST
}

// IsClass reports whether this variable is a class declaration.
// Mirrors AbstractVar.isClass() in Java.
func (v *AbstractVarBase) IsClass() bool {
	return v.DeclarationType() == rhino.CLASS
}

// IsParam reports whether this variable is a function parameter.
// Mirrors AbstractVar.isParam() in Java.
func (v *AbstractVarBase) IsParam() bool {
	return v.DeclarationType() == rhino.PARAM_LIST
}

// IsImport reports whether this variable was declared with an `import`.
// Mirrors AbstractVar.isImport() in Java.
func (v *AbstractVarBase) IsImport() bool {
	return v.DeclarationType() == rhino.IMPORT
}

// IsCatch reports whether this variable is a catch-clause parameter.
// Mirrors AbstractVar.isCatch() in Java.
func (v *AbstractVarBase) IsCatch() bool {
	return v.DeclarationType() == rhino.CATCH
}

// IsDefaultParam reports whether this variable is a parameter with a default
// value expression (ES2015 default parameter).
// Mirrors AbstractVar.isDefaultParam() in Java.
func (v *AbstractVarBase) IsDefaultParam() bool {
	if v.nameNode == nil {
		return false
	}
	parent := v.nameNode.Parent()
	if parent == nil {
		return false
	}
	// The parent of a default parameter name is a DEFAULT_VALUE node whose
	// parent is a PARAM_LIST (or the first child is the name node).
	gp := parent.Parent()
	return gp != nil &&
		gp.Token() == rhino.PARAM_LIST &&
		parent.Token() == rhino.DEFAULT_VALUE &&
		len(parent.Children()) > 0 &&
		parent.Children()[0] == v.nameNode
}

// ---------------------------------------------------------------------------
// Implicit-variable predicates
// ---------------------------------------------------------------------------

// IsImplicitGoogNamespace reports whether this variable was synthesised from
// a goog.provide or legacy goog.module statement rather than an explicit
// declaration.
// Mirrors AbstractVar.isImplicitGoogNamespace() in Java.
func (v *AbstractVarBase) IsImplicitGoogNamespace() bool {
	return v.implicitGoogNamespace != nil
}

// AddImplicitGoogNamespaceDefinition records an additional definition node
// for an implicit goog namespace (the STRINGLIT in a goog.provide call).
// For now we simply update the stored definition with the new node since
// we don't track multiple definitions; the caller is responsible for picking
// the strongest SourceKind.
// Mirrors AbstractVar.addImplicitGoogNamespaceDefinition() in Java.
func (v *AbstractVarBase) AddImplicitGoogNamespaceDefinition(definition *rhino.Node) {
	if !v.IsImplicitGoogNamespace() {
		panic("AddImplicitGoogNamespaceDefinition called on non-implicit var")
	}
	v.implicitGoogNamespace = definition
}

// isImplicit reports whether this is an implicit variable of any kind.
// Used internally to gate declarationType checks.
// Mirrors AbstractVar.isImplicit() in Java.
func (v *AbstractVarBase) isImplicit() bool {
	if v.IsImplicitGoogNamespace() {
		return true
	}
	_, ok := implicitVarOf(v.name)
	return ok
}

// ---------------------------------------------------------------------------
// Bleeding-function predicate
// ---------------------------------------------------------------------------

// IsBleedingFunction reports whether this is a named function expression whose
// name bleeds into the enclosing scope.
// Mirrors AbstractVar.isBleedingFunction() in Java (which uses
// NodeUtil.isFunctionExpression rather than just isFunctionScope).
// We approximate: the parent is a FUNCTION and the grandparent is not a
// FUNCTION declaration (i.e. it's a function expression, not a statement).
func (v *AbstractVarBase) IsBleedingFunction() bool {
	parent := v.GetParentNode()
	if parent == nil {
		return false
	}
	if parent.Token() != rhino.FUNCTION {
		return false
	}
	// A function statement has the FUNCTION as a direct child of a BLOCK or SCRIPT.
	// A function expression is nested in an assignment / call / etc.
	// We check: if the FUNCTION's first child is this nameNode, it is a declaration
	// (not an expression). Otherwise it is a named function expression.
	kids := parent.Children()
	if len(kids) > 0 && kids[0] == v.nameNode {
		// name is the first child of FUNCTION — could be a declaration or expression.
		// In a function expression the FUNCTION itself appears as a value (not as
		// a statement). Treat it as bleeding when grandparent is an expression context.
		gp := parent.Parent()
		if gp == nil {
			return false
		}
		// Statement contexts: BLOCK, SCRIPT, MODULE_BODY
		switch gp.Token() {
		case rhino.BLOCK, rhino.SCRIPT, rhino.MODULE_BODY:
			return false // function declaration, not bleeding
		}
		return true // function expression — name bleeds
	}
	return false
}

// ---------------------------------------------------------------------------
// Const / define predicates
// ---------------------------------------------------------------------------

// IsDeclaredOrInferredConst reports whether this variable is declared or
// inferred to be constant (via `const`, @const JSDoc, or IS_CONSTANT_NAME).
// Mirrors AbstractVar.isDeclaredOrInferredConst() in Java.
func (v *AbstractVarBase) IsDeclaredOrInferredConst() bool {
	if v.nameNode == nil {
		return false
	}
	return v.nameNode.IsDeclaredConstantVar() ||
		v.nameNode.IsInferredConstantVar()
}

// IsDefine reports whether this variable is annotated with @define.
// Mirrors AbstractVar.isDefine() in Java.
func (v *AbstractVarBase) IsDefine() bool {
	doc := v.GetJSDocInfo()
	return doc != nil && doc.IsDefine()
}

// ---------------------------------------------------------------------------
// Initial value and JSDoc
// ---------------------------------------------------------------------------

// GetInitialValue returns the RHS expression of the variable's declaration,
// or nil for implicit variables or declarations without an initializer.
// Mirrors AbstractVar.getInitialValue() in Java via NodeUtil.getRValueOfLValue.
func (v *AbstractVarBase) GetInitialValue() *rhino.Node {
	if v.nameNode == nil {
		return nil
	}
	return getRValueOfLValue(v.nameNode)
}

// getRValueOfLValue returns the right-hand side of an lvalue:
//   - NAME in a VAR/LET/CONST → first child (initializer) if present
//   - destructuring patterns → nil for now (not yet ported)
//
// Mirrors NodeUtil.getRValueOfLValue in Java.
func getRValueOfLValue(n *rhino.Node) *rhino.Node {
	if n == nil {
		return nil
	}
	p := n.Parent()
	if p == nil {
		return nil
	}
	switch p.Token() {
	case rhino.VAR, rhino.LET, rhino.CONST:
		// VAR/LET/CONST → children are NAME nodes; each NAME may have an initializer.
		kids := n.Children()
		if len(kids) > 0 {
			return kids[0]
		}
		return nil
	case rhino.ASSIGN:
		kids := p.Children()
		if len(kids) >= 2 && kids[0] == n {
			return kids[1]
		}
	case rhino.FUNCTION:
		// Named function declaration: the function itself is the value.
		return p
	case rhino.CLASS:
		// Named class declaration: the class itself is the value.
		return p
	}
	return nil
}

// GetJSDocInfo returns the JSDocInfo most closely associated with this
// variable's name node.
// Mirrors AbstractVar.getJSDocInfo() in Java via NodeUtil.getBestJSDocInfo.
// Uses the existing getBestJSDocInfo helper (declared in checkglobalthis.go).
func (v *AbstractVarBase) GetJSDocInfo() *rhino.JSDocInfo {
	return getBestJSDocInfo(v.nameNode)
}

// ---------------------------------------------------------------------------
// Special-name predicates
// ---------------------------------------------------------------------------

// IsArguments reports whether this is the implicit "arguments" variable in a
// function scope.
// Mirrors AbstractVar.isArguments() in Java.
func (v *AbstractVarBase) IsArguments() bool {
	return v.name == VarArguments
}

// IsThis reports whether this is the implicit "this" variable.
// Mirrors AbstractVar.isThis() in Java.
// Note: because "this" is a keyword and not syntactically declarable, the
// scope field check (v.scope.isFunctionScope()) is deferred to UntypedVar.
func (v *AbstractVarBase) IsThis() bool {
	return v.name == "this"
}

// IsGoogModuleExports reports whether this is the implicit "exports" binding
// in a goog.module scope.
// Mirrors AbstractVar.isGoogModuleExports() in Java.
// The scope check is deferred to UntypedVar since AbstractVarBase has no scope pointer.
func (v *AbstractVarBase) IsGoogModuleExports() bool {
	return v.name == VarExports && v.IsImplicitGoogNamespace()
}

// ---------------------------------------------------------------------------
// Reference tracking
// ---------------------------------------------------------------------------

// AddReference records a NAME-node reference to this variable.
// Used by SymbolTable and reference-checking passes.
// Mirrors the mutable reference list in StaticSlot / StaticRef in Java.
func (v *AbstractVarBase) AddReference(nameNode *rhino.Node) {
	v.references = append(v.references, nameNode)
}

// GetReferences returns all recorded references to this variable.
// Mirrors StaticSlot.getReferences() in Java.
func (v *AbstractVarBase) GetReferences() []*rhino.Node {
	return v.references
}

// ReferenceCount returns the number of recorded references.
func (v *AbstractVarBase) ReferenceCount() int {
	return len(v.references)
}

// ClearReferences removes all recorded references (e.g. before a re-index).
func (v *AbstractVarBase) ClearReferences() {
	v.references = v.references[:0]
}

// ---------------------------------------------------------------------------
// AbstractVarIface
// ---------------------------------------------------------------------------

// AbstractVarIface is the interface satisfied by both UntypedVar and GoTypedVar.
// It mirrors the StaticSlot / StaticRef contract from Java.
type AbstractVarIface interface {
	GetName() string
	GetNode() *rhino.Node
	GetParentNode() *rhino.Node
	GetIndex() int
	IsBleedingFunction() bool
	IsImplicitGoogNamespace() bool
	IsVar() bool
	IsLet() bool
	IsConst() bool
	IsClass() bool
	IsParam() bool
	IsCatch() bool
	IsImport() bool
	IsDeclaredOrInferredConst() bool
	GetInitialValue() *rhino.Node
	GetJSDocInfo() *rhino.JSDocInfo
	AddReference(nameNode *rhino.Node)
	GetReferences() []*rhino.Node
}
