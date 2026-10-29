// Port of com/google/javascript/jscomp/Var.java @ closure-compiler v20261021

/*
 * Copyright 2015 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Var is used by Scope to store information about variables.
//
// GetScope, GetDeclaration and GetSymbol keep rhino.StaticSlot /
// rhino.StaticRef's signatures; see the note at the top of abstract_scope.go
// (`v.GetScope().(*Scope)`).
//
// Equality is ScopedName's (name and scope root), so Var-keyed maps and sets
// use the jutil.Object* collections.
type Var struct {
	varCommon

	scope  *Scope
	nameID uint32
}

var _ AbstractVar = (*Var)(nil)

const (
	Var_ARGUMENTS = "arguments"
	Var_EXPORTS   = "exports"
)

// Java: Var(String, Node, Scope, int, CompilerInput, Node)
func NewVar(
	name string,
	nameNode rhino.Node,
	scope *Scope,
	index int,
	input *CompilerInput,
	implicitGoogNamespaceDefinition rhino.Node) *Var {
	return newVarWithID(rhino.RhinoStringPool_InternID(name), nameNode, scope, index, input, implicitGoogNamespaceDefinition)
}

// newVarWithID is NewVar for an interned name.
func newVarWithID(
	nameID uint32,
	nameNode rhino.Node,
	scope *Scope,
	index int,
	input *CompilerInput,
	implicitGoogNamespaceDefinition rhino.Node) *Var {
	if !nameNode.IsNil() {
		switch nameNode.GetToken() {
		case rhino.Token_MODULE_BODY, rhino.Token_NAME, rhino.Token_IMPORT_STAR:
		default:
			failInvalidVarNameNode(nameNode)
		}
	}
	return initVar(&Var{}, nameID, nameNode, scope, index, input, implicitGoogNamespaceDefinition)
}

func initVar(
	v *Var,
	nameID uint32,
	nameNode rhino.Node,
	scope *Scope,
	index int,
	input *CompilerInput,
	implicitGoogNamespaceDefinition rhino.Node,
) *Var {
	v.scope = scope
	v.nameID = nameID
	v.initVarCommon(nameNode, index, input, implicitGoogNamespaceDefinition)
	return v
}

//go:noinline
func failInvalidVarNameNode(nameNode rhino.Node) {
	panic(jutil.NewIllegalArgumentExceptionMsg("Invalid name node " + jutil.ToJavaString(nameNode)))
}

// Java: AbstractVar.getName()
func (v *Var) GetName() string {
	return rhino.RhinoStringPool_StringFromID(v.nameID)
}

// Java: AbstractVar.getScope()
func (v *Var) GetScope() rhino.StaticScope {
	if v.scope == nil {
		return nil
	}
	return v.scope
}

func (v *Var) getScopeS() *Scope {
	return v.scope
}

// Java: AbstractVar.getScopeRoot()
func (v *Var) GetScopeRoot() rhino.Node {
	return v.scope.rootNode
}

// Java: AbstractVar.getSourceFile()
func (v *Var) GetSourceFile() rhino.StaticSourceFile {
	return v.sourceFile(v.scope.rootNode)
}

// Java: AbstractVar.isGlobal()
func (v *Var) IsGlobal() bool {
	return v.scope.parent == nil
}

// Java: AbstractVar.isLocal()
func (v *Var) IsLocal() bool {
	return v.scope.parent != nil
}

// Java: ScopedName.hashCode()
func (v *Var) HashCode() int32 {
	if h := atomic.LoadInt32(&v.nameHash); h != 0 {
		return h
	}
	return v.computeHashCode(v.scope.rootNode, v.GetName())
}

// Java: ScopedName.equals(Object)
func (v *Var) Equals(other any) bool {
	if o, ok := other.(*Var); ok {
		return o != nil && v.nameID == o.nameID && scopedName_nodeEquals(v.scope.rootNode, o.scope.rootNode)
	}
	return equalsScopedName(v.GetName(), v.scope.rootNode, other)
}

// DeclarationToken returns the enclosing declaration token of v, or (0, false)
// when v is an implicit variable.
func (v *Var) DeclarationToken() (rhino.Token, bool) {
	if tok, ok := v.findDeclarationToken(); ok {
		return tok, true
	}
	if !v.IsImplicit() {
		failInvalidDeclarationType()
	}
	return 0, false
}

func (v *Var) isDeclarationToken(want rhino.Token) bool {
	tok, ok := v.DeclarationToken()
	return ok && tok == want
}

func (v *Var) IsVar() bool       { return v.isDeclarationToken(rhino.Token_VAR) }
func (v *Var) IsCatch() bool     { return v.isDeclarationToken(rhino.Token_CATCH) }
func (v *Var) IsLet() bool       { return v.isDeclarationToken(rhino.Token_LET) }
func (v *Var) IsConst() bool     { return v.isDeclarationToken(rhino.Token_CONST) }
func (v *Var) IsClass() bool     { return v.isDeclarationToken(rhino.Token_CLASS) }
func (v *Var) IsParam() bool     { return v.isDeclarationToken(rhino.Token_PARAM_LIST) }
func (v *Var) IsImport() bool    { return v.isDeclarationToken(rhino.Token_IMPORT) }
func (v *Var) IsArguments() bool { return v.nameID == implicitArgumentsID && v.scope.IsFunctionScope() }
func (v *Var) IsGoogModuleExports() bool {
	return v.nameID == implicitExportsID && v.scope.IsModuleScope() && v.IsImplicit()
}
func (v *Var) IsThis() bool { return v.nameID == implicitThisID && v.scope.IsFunctionScope() }

func (v *Var) IsImplicit() bool {
	if v.IsImplicitGoogNamespace() {
		return true
	}
	imp, ok := AbstractScope_ImplicitVar_LookupID(v.nameID)
	return ok && imp.isMadeByRootNode(v.scope.rootNode)
}

func (v *Var) AddImplicitGoogNamespaceDefinition(definition rhino.Node) {
	v.addImplicitGoogNamespaceDefinition(v, definition)
}

func (v *Var) GetImplicitGoogNamespaceStrength() rhino.StaticSourceFile_SourceKind {
	return v.getImplicitGoogNamespaceStrength(v)
}

// Java: AbstractVar.getSymbol()
func (v *Var) GetSymbol() rhino.StaticSlot {
	return v
}

// Java: AbstractVar.getDeclaration()
func (v *Var) GetDeclaration() rhino.StaticRef {
	if v.nameNode.IsNil() {
		return nil
	}
	return v
}

// Java: Var.createImplicitGoogNamespace(String, Scope, Node)
func Var_CreateImplicitGoogNamespace(name string, scope *Scope, definition rhino.Node) *Var {
	return NewVar(name, rhino.NilNode /* nameNode */, scope, -1, nil /* input */, definition)
}

// Java: Var.toString()
func (v *Var) String() string {
	return "Var " + v.GetName() + " @ " + jutil.ToJavaString(v.nameNode)
}
