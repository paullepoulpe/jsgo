// Port of com/google/javascript/jscomp/AbstractVar.java @ closure-compiler v20261015

/*
 * Copyright 2018 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import (
	"sync/atomic"
	"unsafe"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// AbstractVar is used by Scope to store information about variables.
//
// Java: class AbstractVar<S extends AbstractScope<S, V>, V extends
// AbstractVar<S, V>> extends ScopedName implements StaticSlot, StaticRef.
//
// Go-only: language-adapter — as for AbstractScope (see the note at the top
// of abstract_scope.go): this interface is the wildcard AbstractVar<?, ?>,
// AbstractVarBase[S, V] holds the fields and implementations, and Var and
// TypedVar embed it. getScope() (S in Java), getDeclaration() and
// getSymbol() (V in Java) narrow rhino.StaticSlot / rhino.StaticRef methods,
// so they keep those Go signatures and never return a typed nil; callers
// that need the Java type assert: `v.GetScope().(*Scope)`.
type AbstractVar interface {
	ScopedName
	rhino.StaticSlot
	rhino.StaticRef

	// Java: AbstractVar.getInput()
	GetInput() *CompilerInput

	// Java: AbstractVar.getParentNode()
	GetParentNode() rhino.Node

	// Java: AbstractVar.isBleedingFunction()
	IsBleedingFunction() bool

	// Java: AbstractVar.getIndex()
	GetIndex() int

	// Java: AbstractVar.isGlobal()
	IsGlobal() bool

	// Java: AbstractVar.isLocal()
	IsLocal() bool

	// Java: AbstractVar.isExtern()
	IsExtern() bool

	// Java: AbstractVar.isDeclaredOrInferredConst()
	IsDeclaredOrInferredConst() bool

	// Java: AbstractVar.isDefine()
	IsDefine() bool

	// Java: AbstractVar.getInitialValue()
	GetInitialValue() rhino.Node

	// Java: AbstractVar.getNameNode()
	GetNameNode() rhino.Node

	// Java: AbstractVar.isVar()
	IsVar() bool

	// Java: AbstractVar.isCatch()
	IsCatch() bool

	// Java: AbstractVar.isLet()
	IsLet() bool

	// Java: AbstractVar.isConst()
	IsConst() bool

	// Java: AbstractVar.isClass()
	IsClass() bool

	// Java: AbstractVar.isParam()
	IsParam() bool

	// Java: AbstractVar.isDefaultParam()
	IsDefaultParam() bool

	// Java: AbstractVar.isImport()
	IsImport() bool

	// Java: AbstractVar.isArguments()
	IsArguments() bool

	// Java: AbstractVar.isGoogModuleExports()
	IsGoogModuleExports() bool

	// Java: AbstractVar.isThis()
	IsThis() bool

	// Java: AbstractVar.declarationType()
	DeclarationType() *rhino.Token

	// Java: AbstractVar.isImplicitGoogNamespace()
	IsImplicitGoogNamespace() bool

	// Java: AbstractVar.addImplicitGoogNamespaceDefinition(Node)
	AddImplicitGoogNamespaceDefinition(definition rhino.Node)

	// Java: AbstractVar.getImplicitGoogNamespaceStrength()
	GetImplicitGoogNamespaceStrength() rhino.StaticSourceFile_SourceKind

	// Java: Object.toString() (Var and TypedVar override it)
	String() string
}

// abstractVarBound is the bound `V extends AbstractVar<S, V>`: the methods
// AbstractScopeBase and AbstractVarBase call on a V, and the interfaces a V
// is converted to. getScopeS is promoted from AbstractVarBase.
//
// Go-only: language-adapter — the bound of a Java type parameter; see the
// note at the top of abstract_scope.go. AbstractVar is not embedded, for the
// reason given for abstractScopeBound (its GetInput() reaches CompilerInput,
// which holds a TypedScope).
type abstractVarBound[S any, V any] interface {
	comparable
	ScopedName
	rhino.StaticSlot
	rhino.StaticRef

	getScopeS() S
}

const (
	varIndexMask              uint32 = 0x00FFFFFF
	varImplicitStrengthShift         = 24
	varImplicitStrengthMask   uint32 = 0x03 << varImplicitStrengthShift
	varFlagMarkedEscaped      uint32 = 1 << 26
	varFlagMarkedAssignedOnce uint32 = 1 << 27
	varFlagTypeInferred       uint32 = 1 << 28
)

// AbstractVarBase holds the fields and the implementations of AbstractVar;
// Var and TypedVar embed it at offset 0 and call initAbstractVarBase.
//
// Go-only: performance — the name is not here: a Var holds it as an interned
// string ID (scopes look vars up by ID), a TypedVar as a string (most typed
// names are qualified names that nothing else interns).
type AbstractVarBase[S abstractScopeBound[S, V], V abstractVarBound[S, V]] struct {
	nameNode rhino.Node

	// Input source
	input *CompilerInput

	scope S

	// Lower 24 bits: uint32(index + 1); bits 24..25: implicit goog namespace
	// strength (0 = none, 1..3 = SourceKind + 1); bits 26..28: TypedVar flags.
	indexAndFlags uint32
	// Go-only: performance — immutable name hash, or per-pass scratch index.
	nameHash int32
}

var (
	_ [0]struct{} = [unsafe.Offsetof(Var{}.AbstractVarBase)]struct{}{}
	_ [0]struct{} = [unsafe.Offsetof(TypedVar{}.AbstractVarBase)]struct{}{}
)

// self is this var as a V.
func (b *AbstractVarBase[S, V]) self() V {
	p := unsafe.Pointer(b)
	return *(*V)(unsafe.Pointer(&p))
}

// initAbstractVarBase is the AbstractVar constructor, but for the name,
// which Var and TypedVar set.
//
// nameNode is the node representing this variables declaration or null.
// scope is the scope containing this var. index is the index at which the
// var is declared in the scope. Must be either positive, or -1 for implicit
// variables. input is the compiler input containing the given nameNode, if
// any. May be null to allow for declaring the native types in the type
// system.
//
// Java: AbstractVar(String, Node, S, int, CompilerInput, Node)
func (b *AbstractVarBase[S, V]) initAbstractVarBase(
	nameNode rhino.Node,
	scope S,
	index int,
	input *CompilerInput,
	implicitGoogNamespaceDefinition rhino.Node) {
	jutil.CheckArgumentMsg1(index >= -1, "%s", index)
	flags := uint32(index+1) & varIndexMask
	if !implicitGoogNamespaceDefinition.IsNil() {
		b.nameNode = rhino.NilNode
		strength := abstractVar_strengthOf(implicitGoogNamespaceDefinition)
		flags |= (uint32(strength) + 1) << varImplicitStrengthShift
	} else {
		b.nameNode = nameNode
	}
	b.scope = scope
	b.indexAndFlags = flags
	b.nameHash = 0
	b.input = input
}

// Java: AbstractVar.getScopeRoot()
func (b *AbstractVarBase[S, V]) GetScopeRoot() rhino.Node {
	return b.scope.GetRootNode()
}

// HashCode preserves ScopedName's final hash while avoiding repeated name
// hashing and virtual self calls. The scope root and name are immutable.
// Go-only: performance — indexAndFlags is 32 bits, leaving room for the cached
// hash without increasing the size of Var or TypedVar.
// Java: ScopedName.hashCode()
func (b *AbstractVarBase[S, V]) HashCode() int32 {
	if h := atomic.LoadInt32(&b.nameHash); h != 0 {
		return h
	}
	return b.computeHashCode()
}

//go:noinline
func (b *AbstractVarBase[S, V]) computeHashCode() int32 {
	root := b.scope.GetRootNode()
	if root.IsNil() {
		panic(jutil.NewNullPointerException())
	}
	h := 31*(31+jutil.StringHashCode(b.self().GetName())) + root.IdentityHashCode()
	atomic.StoreInt32(&b.nameHash, h)
	return h
}

// Java: AbstractVar.getNode()
func (b *AbstractVarBase[S, V]) GetNode() rhino.Node {
	return b.nameNode
}

// Java: AbstractVar.getInput()
func (b *AbstractVarBase[S, V]) GetInput() *CompilerInput {
	return b.input
}

// Java: AbstractVar.getSourceFile()
func (b *AbstractVarBase[S, V]) GetSourceFile() rhino.StaticSourceFile {
	n := b.GetNode()
	if n.IsNil() {
		n = b.scope.GetRootNode()
	}
	return n.GetStaticSourceFile()
}

// Java: AbstractVar.getParentNode()
func (b *AbstractVarBase[S, V]) GetParentNode() rhino.Node {
	if b.GetNode().IsNil() {
		return rhino.NilNode
	}
	return b.GetNode().GetParent()
}

// IsBleedingFunction reports whether this is a bleeding function (an
// anonymous named function that bleeds into the inner scope).
//
// Java: AbstractVar.isBleedingFunction()
func (b *AbstractVarBase[S, V]) IsBleedingFunction() bool {
	parent := b.GetParentNode()
	return !parent.IsNil() && NodeUtil_IsFunctionExpression(parent)
}

// GetScope is final in Java and returns S; see AbstractVar.
//
// Java: AbstractVar.getScope()
func (b *AbstractVarBase[S, V]) GetScope() rhino.StaticScope {
	var null S
	if b.scope == null {
		return nil
	}
	return b.scope
}

// getScopeS is getScope with its Java return type S.
//
// Go-only: language-adapter — see GetScope.
func (b *AbstractVarBase[S, V]) getScopeS() S {
	return b.scope
}

// GetIndex is the index at which the var is declared. e.g. if it's 0, it's
// the first declared variable in that scope
//
// Java: AbstractVar.getIndex()
func (b *AbstractVarBase[S, V]) GetIndex() int {
	return int(b.indexAndFlags&varIndexMask) - 1
}

// Java: AbstractVar.isGlobal()
func (b *AbstractVarBase[S, V]) IsGlobal() bool {
	return b.scope.IsGlobal()
}

// Java: AbstractVar.isLocal()
func (b *AbstractVarBase[S, V]) IsLocal() bool {
	return b.scope.IsLocal()
}

// Java: AbstractVar.isExtern()
func (b *AbstractVarBase[S, V]) IsExtern() bool {
	return b.input == nil || b.input.IsExtern()
}

// IsDeclaredOrInferredConst returns true if the variable is declared or
// inferred to be a constant.
//
// Java: AbstractVar.isDeclaredOrInferredConst()
func (b *AbstractVarBase[S, V]) IsDeclaredOrInferredConst() bool {
	declarationNode := b.GetNode()
	if declarationNode.IsNil() {
		return false
	}

	return declarationNode.IsDeclaredConstantVar() ||
		declarationNode.IsInferredConstantVar() ||
		declarationNode.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME)
}

// IsDefine returns true if the variable is declared as a define. A variable
// is a define if it is annotated by @define.
//
// Java: AbstractVar.isDefine()
func (b *AbstractVarBase[S, V]) IsDefine() bool {
	info := b.GetJSDocInfo()
	return info != nil && info.IsDefine()
}

// Java: AbstractVar.getInitialValue()
func (b *AbstractVarBase[S, V]) GetInitialValue() rhino.Node {
	if b.GetNode().IsNil() {
		return rhino.NilNode
	}
	return NodeUtil_GetRValueOfLValue(b.GetNode())
}

// Java: AbstractVar.getNameNode()
func (b *AbstractVarBase[S, V]) GetNameNode() rhino.Node {
	return b.GetNode()
}

// Java: AbstractVar.getJSDocInfo()
func (b *AbstractVarBase[S, V]) GetJSDocInfo() *rhino.JSDocInfo {
	if b.GetNode().IsNil() {
		return nil
	}
	return NodeUtil_GetBestJSDocInfo(b.GetNode())
}

// isDeclarationType is `declarationType() == token` (an enum reference
// compared with ==, so null never matches).
//
// Go-only: language-adapter — DeclarationType returns a nullable *Token.
func (b *AbstractVarBase[S, V]) isDeclarationType(token rhino.Token) bool {
	t := b.DeclarationType()
	return t != nil && *t == token
}

// Java: AbstractVar.isVar()
func (b *AbstractVarBase[S, V]) IsVar() bool {
	return b.isDeclarationType(rhino.Token_VAR)
}

// Java: AbstractVar.isCatch()
func (b *AbstractVarBase[S, V]) IsCatch() bool {
	return b.isDeclarationType(rhino.Token_CATCH)
}

// Java: AbstractVar.isLet()
func (b *AbstractVarBase[S, V]) IsLet() bool {
	return b.isDeclarationType(rhino.Token_LET)
}

// Java: AbstractVar.isConst()
func (b *AbstractVarBase[S, V]) IsConst() bool {
	return b.isDeclarationType(rhino.Token_CONST)
}

// Java: AbstractVar.isClass()
func (b *AbstractVarBase[S, V]) IsClass() bool {
	return b.isDeclarationType(rhino.Token_CLASS)
}

// Java: AbstractVar.isParam()
func (b *AbstractVarBase[S, V]) IsParam() bool {
	return b.isDeclarationType(rhino.Token_PARAM_LIST)
}

// Java: AbstractVar.isDefaultParam()
func (b *AbstractVarBase[S, V]) IsDefaultParam() bool {
	parent := b.GetNode().GetParent()
	return parent.GetParent().IsParamList() && parent.IsDefaultValue() &&
		parent.GetFirstChild() == b.nameNode
}

// Java: AbstractVar.isImport()
func (b *AbstractVarBase[S, V]) IsImport() bool {
	return b.isDeclarationType(rhino.Token_IMPORT)
}

// Java: AbstractVar.isArguments()
func (b *AbstractVarBase[S, V]) IsArguments() bool {
	return b.self().GetName() == Var_ARGUMENTS && b.scope.IsFunctionScope()
}

// Java: AbstractVar.isGoogModuleExports()
func (b *AbstractVarBase[S, V]) IsGoogModuleExports() bool {
	return b.self().GetName() == Var_EXPORTS && b.scope.IsModuleScope() && b.IsImplicit()
}

// Java: AbstractVar.isThis()
func (b *AbstractVarBase[S, V]) IsThis() bool {
	return b.self().GetName() == "this" && b.scope.IsFunctionScope()
}

// Java: AbstractVar.isImplicit()
func (b *AbstractVarBase[S, V]) IsImplicit() bool {
	if b.IsImplicitGoogNamespace() {
		return true
	}
	var_ := AbstractScope_ImplicitVar_Of(b.self().GetName())
	return var_ != nil && var_.IsMadeByScope(abstractScope_of(b.scope))
}

var abstractVar_DECLARATION_TYPES = jutil.NewEnumSet(
	rhino.Token_VAR,
	rhino.Token_LET,
	rhino.Token_CONST,
	rhino.Token_FUNCTION,
	rhino.Token_CLASS,
	rhino.Token_CATCH,
	rhino.Token_IMPORT,
	rhino.Token_PARAM_LIST)

// Go-only: language-adapter — a Java Token reference may be null; such
// values are *rhino.Token here (nil is null), always pointing into this
// table so that no call allocates.
var abstractVar_tokenPtrs = rhino.Token_values()

// Java: AbstractVar.declarationType()
func (b *AbstractVarBase[S, V]) DeclarationType() *rhino.Token {
	if b.IsImplicitGoogNamespace() {
		return nil
	}
	for current := b.GetNode(); !current.IsNil(); current = current.GetParent() {
		if abstractVar_DECLARATION_TYPES.Contains(current.GetToken()) {
			return &abstractVar_tokenPtrs[current.GetToken()]
		}
	}
	jutil.CheckStateMsg2(
		b.IsImplicit(),
		"The nameNode for %s must be a descendant of one of: %s", b.self(), abstractVar_DECLARATION_TYPES)
	return nil
}

// IsImplicitGoogNamespace reports whether this name comes from a
// goog.provide or legacy goog.module and does not have a normal syntactical
// definition
//
// Java: AbstractVar.isImplicitGoogNamespace()
func (b *AbstractVarBase[S, V]) IsImplicitGoogNamespace() bool {
	return b.indexAndFlags&varImplicitStrengthMask != 0
}

// AddImplicitGoogNamespaceDefinition indicates this namespace was provided
// at the given node; definition is the STRINGLIT in a goog.provide or
// goog.module call.
//
// Java: AbstractVar.addImplicitGoogNamespaceDefinition(Node)
func (b *AbstractVarBase[S, V]) AddImplicitGoogNamespaceDefinition(definition rhino.Node) {
	jutil.CheckStateMsg1(b.IsImplicitGoogNamespace(), "%s", b.self())
	cur := rhino.StaticSourceFile_SourceKind(((b.indexAndFlags & varImplicitStrengthMask) >> varImplicitStrengthShift) - 1)
	strength := abstractVar_strongerOf(cur, abstractVar_strengthOf(definition))
	b.indexAndFlags = (b.indexAndFlags &^ varImplicitStrengthMask) | ((uint32(strength) + 1) << varImplicitStrengthShift)
}

// GetImplicitGoogNamespaceStrength returns the strongest kind of all
// definitions of this namespace
//
// Java: AbstractVar.getImplicitGoogNamespaceStrength()
func (b *AbstractVarBase[S, V]) GetImplicitGoogNamespaceStrength() rhino.StaticSourceFile_SourceKind {
	jutil.CheckStateMsg1(b.IsImplicitGoogNamespace(), "%s", b.self())
	return rhino.StaticSourceFile_SourceKind(((b.indexAndFlags & varImplicitStrengthMask) >> varImplicitStrengthShift) - 1)
}

// Java: AbstractVar.strengthOf(Node)
func abstractVar_strengthOf(n rhino.Node) rhino.StaticSourceFile_SourceKind {
	source := n.GetStaticSourceFile()
	if source == nil {
		return rhino.StaticSourceFile_SourceKind_EXTERN
	}

	return source.GetKind()
}

// Java: AbstractVar.strongerOf(SourceKind, SourceKind)
func abstractVar_strongerOf(left, right rhino.StaticSourceFile_SourceKind) rhino.StaticSourceFile_SourceKind {
	if left == rhino.StaticSourceFile_SourceKind_STRONG || right == rhino.StaticSourceFile_SourceKind_STRONG {
		return rhino.StaticSourceFile_SourceKind_STRONG
	} else if left == rhino.StaticSourceFile_SourceKind_EXTERN || right == rhino.StaticSourceFile_SourceKind_EXTERN {
		// Externs are stronger because they aren't deleted.
		return rhino.StaticSourceFile_SourceKind_EXTERN
	}
	return rhino.StaticSourceFile_SourceKind_WEAK
}

// Equals is ScopedName's final equality.
//
// Java: ScopedName.equals(Object)
func (b *AbstractVarBase[S, V]) Equals(other any) bool {
	scoped, ok := other.(ScopedName)
	if !ok || jutil.IsNil(scoped) {
		return false
	}
	return b.self().GetName() == scoped.GetName() && scopedName_nodeEquals(b.GetScopeRoot(), scoped.GetScopeRoot())
}
