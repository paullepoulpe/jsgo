// Port of com/google/javascript/jscomp/AbstractVar.java @ closure-compiler v20261021

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
	"sync"
	"sync/atomic"
	"weak"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// AbstractVar is used by Scope to store information about variables.
//
// Java: class AbstractVar<S extends AbstractScope<S, V>, V extends
// AbstractVar<S, V>> extends ScopedName implements StaticSlot, StaticRef.
//
// Go-only: language-adapter — this interface is the wildcard AbstractVar<?, ?>;
// varCommon holds the shared representation embedded by concrete Var and
// TypedVar. GetScope, GetDeclaration and GetSymbol keep rhino.StaticSlot /
// rhino.StaticRef signatures and never return a typed nil; callers that need
// the concrete scope type assert (`v.GetScope().(*Scope)`).
type AbstractVar interface {
	ScopedName
	rhino.StaticSlot
	rhino.StaticRef
}

const (
	varIndexMask              uint32 = 0x00FFFFFF
	varImplicitStrengthShift         = 24
	varImplicitStrengthMask   uint32 = 0x03 << varImplicitStrengthShift
	varFlagMarkedEscaped      uint32 = 1 << 26
	varFlagMarkedAssignedOnce uint32 = 1 << 27
	varFlagTypeInferred       uint32 = 1 << 28
)

// varCommon holds the shared pointer-free fields of Var and TypedVar.
//
// Go-only: performance — Var stores its name as an interned string ID and
// TypedVar as a string; keeping the four 32-bit fields together (16 bytes)
// lets Var fit in 32 bytes and TypedVar in 48 bytes.
type varCommon struct {
	nameNode rhino.Node

	// Input source index in compilerInputTable (0: none), 4 bytes beside
	// nameNode instead of an 8-byte pointer.
	input uint32

	// Lower 24 bits: uint32(index + 1); bits 24..25: implicit goog namespace
	// strength (0 = none, 1..3 = SourceKind + 1); bits 26..28: TypedVar flags.
	indexAndFlags uint32
	// Cached ScopedName.hashCode() (0 = uncomputed).
	nameHash int32
}

// initVarCommon initializes the shared AbstractVar fields.
func (b *varCommon) initVarCommon(
	nameNode rhino.Node,
	index int,
	input *CompilerInput,
	implicitGoogNamespaceDefinition rhino.Node,
) {
	if index < -1 {
		failInvalidVarIndex(index)
	}
	flags := uint32(index+1) & varIndexMask
	if !implicitGoogNamespaceDefinition.IsNil() {
		b.nameNode = rhino.NilNode
		strength := abstractVar_strengthOf(implicitGoogNamespaceDefinition)
		flags |= (uint32(strength) + 1) << varImplicitStrengthShift
	} else {
		b.nameNode = nameNode
	}
	b.indexAndFlags = flags
	b.nameHash = 0
	b.input = compilerInputTable_indexOf(input)
}

//go:noinline
func failInvalidVarIndex(index int) {
	jutil.CheckArgumentMsg1(false, "%s", index)
}

//go:noinline
func (b *varCommon) computeHashCode(root rhino.Node, name string) int32 {
	if root.IsNil() {
		panic(jutil.NewNullPointerException())
	}
	h := 31*(31+jutil.StringHashCode(name)) + root.IdentityHashCode()
	atomic.StoreInt32(&b.nameHash, h)
	return h
}

// Java: AbstractVar.getNode()
func (b *varCommon) GetNode() rhino.Node {
	return b.nameNode
}

// Java: AbstractVar.getInput()
func (b *varCommon) GetInput() *CompilerInput {
	return compilerInputTable_at(b.input)
}

func (b *varCommon) sourceFile(scopeRoot rhino.Node) rhino.StaticSourceFile {
	n := b.nameNode
	if n.IsNil() {
		n = scopeRoot
	}
	return n.GetStaticSourceFile()
}

// Java: AbstractVar.getParentNode()
func (b *varCommon) GetParentNode() rhino.Node {
	if b.nameNode.IsNil() {
		return rhino.NilNode
	}
	return b.nameNode.GetParent()
}

// IsBleedingFunction reports whether this is a bleeding function (an
// anonymous named function that bleeds into the inner scope).
//
// Java: AbstractVar.isBleedingFunction()
func (b *varCommon) IsBleedingFunction() bool {
	parent := b.GetParentNode()
	return !parent.IsNil() && NodeUtil_IsFunctionExpression(parent)
}

// GetIndex is the index at which the var is declared (e.g. 0 for the first
// declared variable in that scope, or -1 for implicit variables).
//
// Java: AbstractVar.getIndex()
func (b *varCommon) GetIndex() int {
	return int(b.indexAndFlags&varIndexMask) - 1
}

// Java: AbstractVar.isExtern()
func (b *varCommon) IsExtern() bool {
	input := b.GetInput()
	return input == nil || input.IsExtern()
}

// IsDeclaredOrInferredConst returns true if the variable is declared or
// inferred to be a constant.
//
// Java: AbstractVar.isDeclaredOrInferredConst()
func (b *varCommon) IsDeclaredOrInferredConst() bool {
	declarationNode := b.nameNode
	if declarationNode.IsNil() {
		return false
	}
	return declarationNode.IsDeclaredConstantVar() ||
		declarationNode.IsInferredConstantVar() ||
		declarationNode.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME)
}

// IsDefine returns true if the variable is declared as a define (@define).
//
// Java: AbstractVar.isDefine()
func (b *varCommon) IsDefine() bool {
	info := b.GetJSDocInfo()
	return info != nil && info.IsDefine()
}

// Java: AbstractVar.getInitialValue()
func (b *varCommon) GetInitialValue() rhino.Node {
	if b.nameNode.IsNil() {
		return rhino.NilNode
	}
	return NodeUtil_GetRValueOfLValue(b.nameNode)
}

// Java: AbstractVar.getNameNode()
func (b *varCommon) GetNameNode() rhino.Node {
	return b.nameNode
}

// Java: AbstractVar.getJSDocInfo()
func (b *varCommon) GetJSDocInfo() *rhino.JSDocInfo {
	if b.nameNode.IsNil() {
		return nil
	}
	return NodeUtil_GetBestJSDocInfo(b.nameNode)
}

// Java: AbstractVar.isDefaultParam()
func (b *varCommon) IsDefaultParam() bool {
	parent := b.nameNode.GetParent()
	return parent.GetParent().IsParamList() && parent.IsDefaultValue() &&
		parent.GetFirstChild() == b.nameNode
}

const abstractVar_DECLARATION_TYPES_STR = "[PARAM_LIST, FUNCTION, VAR, CATCH, CONST, CLASS, LET, IMPORT]"

func abstractVar_isDeclarationType(tok rhino.Token) bool {
	switch tok {
	case rhino.Token_VAR,
		rhino.Token_LET,
		rhino.Token_CONST,
		rhino.Token_FUNCTION,
		rhino.Token_CLASS,
		rhino.Token_CATCH,
		rhino.Token_IMPORT,
		rhino.Token_PARAM_LIST:
		return true
	default:
		return false
	}
}

func (b *varCommon) findDeclarationToken() (rhino.Token, bool) {
	if b.IsImplicitGoogNamespace() {
		return 0, false
	}
	for current := b.nameNode; !current.IsNil(); current = current.GetParent() {
		if tok := current.GetToken(); abstractVar_isDeclarationType(tok) {
			return tok, true
		}
	}
	return 0, false
}

//go:noinline
func failInvalidDeclarationType() {
	jutil.CheckStateMsg1(
		false,
		"The nameNode must be a descendant of one of: %s", abstractVar_DECLARATION_TYPES_STR)
}

//go:noinline
func failNotImplicitGoogNamespace(self any) {
	jutil.CheckStateMsg1(false, "%s", self)
}

// IsImplicitGoogNamespace reports whether this name comes from a
// goog.provide or legacy goog.module and does not have a normal syntactical
// definition.
//
// Java: AbstractVar.isImplicitGoogNamespace()
func (b *varCommon) IsImplicitGoogNamespace() bool {
	return b.indexAndFlags&varImplicitStrengthMask != 0
}

func (b *varCommon) addImplicitGoogNamespaceDefinition(definition rhino.Node) bool {
	if !b.IsImplicitGoogNamespace() {
		return false
	}
	cur := rhino.StaticSourceFile_SourceKind(((b.indexAndFlags & varImplicitStrengthMask) >> varImplicitStrengthShift) - 1)
	strength := abstractVar_strongerOf(cur, abstractVar_strengthOf(definition))
	b.indexAndFlags = (b.indexAndFlags &^ varImplicitStrengthMask) | ((uint32(strength) + 1) << varImplicitStrengthShift)
	return true
}

func (b *varCommon) getImplicitGoogNamespaceStrength() (rhino.StaticSourceFile_SourceKind, bool) {
	if !b.IsImplicitGoogNamespace() {
		return 0, false
	}
	return rhino.StaticSourceFile_SourceKind(((b.indexAndFlags & varImplicitStrengthMask) >> varImplicitStrengthShift) - 1), true
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

func equalsScopedName(name string, scopeRoot rhino.Node, other any) bool {
	scoped, ok := other.(ScopedName)
	if !ok || isScopedNameNil(scoped) {
		return false
	}
	return name == scoped.GetName() && scopedName_nodeEquals(scopeRoot, scoped.GetScopeRoot())
}

// compilerInputTable numbers the CompilerInputs that variables refer to, so
// that a variable keeps a 32-bit index instead of a pointer.
//
// Go-only: performance — with the index beside the name node, a Var takes
// 32 bytes instead of 48 and a TypedVar 48 instead of 64, for the hundreds
// of thousands of variables a large compilation declares. The table holds
// the inputs weakly, so it never keeps a compilation alive; a variable is
// used while its compiler, which holds its inputs, is.
var compilerInputTable struct {
	mu       sync.Mutex
	entries  []weak.Pointer[CompilerInput]                 // entry i has index i+1
	snapshot atomic.Pointer[[]weak.Pointer[CompilerInput]] // immutable prefix snapshot for lock-free reads
}

// compilerInputTable_indexOf returns the index of input (0 for nil),
// adding it to the table on first use.
func compilerInputTable_indexOf(input *CompilerInput) uint32 {
	if input == nil {
		return 0
	}
	if i := input.varTableIndex.Load(); i != 0 {
		return i
	}
	compilerInputTable.mu.Lock()
	defer compilerInputTable.mu.Unlock()
	if i := input.varTableIndex.Load(); i != 0 {
		return i
	}
	cur := compilerInputTable.entries
	var next []weak.Pointer[CompilerInput]
	if len(cur) < cap(cur) {
		next = cur[:len(cur)+1]
		next[len(cur)] = weak.Make(input)
	} else {
		next = make([]weak.Pointer[CompilerInput], len(cur)+1, max(cap(cur)*2, 64))
		copy(next, cur)
		next[len(cur)] = weak.Make(input)
	}
	compilerInputTable.entries = next
	snap := new([]weak.Pointer[CompilerInput])
	*snap = next
	compilerInputTable.snapshot.Store(snap)
	i := uint32(len(next))
	input.varTableIndex.Store(i)
	return i
}

// compilerInputTable_at returns the input with index i (nil for 0).
func compilerInputTable_at(i uint32) *CompilerInput {
	if i == 0 {
		return nil
	}
	return (*compilerInputTable.snapshot.Load())[i-1].Value()
}
