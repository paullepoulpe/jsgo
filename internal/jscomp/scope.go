// Port of com/google/javascript/jscomp/Scope.java @ closure-compiler v20261015

/*
 * Copyright 2004 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Scope contains information about a variable scope in JavaScript. Scopes
// can be nested, a scope points back to its parent scope. A Scope contains
// information about variables defined in that scope.
//
// Most methods come from AbstractScopeBase[*Scope, *Var] and are typed:
// GetVar returns *Var, GetParent and GetClosestHoistScope return *Scope,
// GetVarIterable is an iter.Seq[*Var]. GetSlot, GetOwnSlot and
// GetParentScope keep rhino.StaticScope's signatures; see the note at the
// top of abstract_scope.go.
//
// See NodeTraversal.
type Scope struct {
	AbstractScopeBase[*Scope, *Var]

	parent *Scope
	depth  int
}

// Java: Scope.createGlobalScope(Node)
func Scope_CreateGlobalScope(rootNode *rhino.Node) *Scope {
	return newScope(rootNode)
}

// Java: Scope.createChildScope(Scope, Node)
func Scope_CreateChildScope(parent *Scope, rootNode *rhino.Node) *Scope {
	return newScopeWithParent(parent, rootNode)
}

// Java: Scope(Scope, Node)
func newScopeWithParent(parent *Scope, rootNode *rhino.Node) *Scope {
	s := &Scope{}
	s.InitAbstractScopeBase(s, rootNode)
	s.CheckChildScope(parent)
	s.parent = parent
	s.depth = parent.GetDepth() + 1
	return s
}

// Java: Scope(Node)
func newScope(rootNode *rhino.Node) *Scope {
	s := &Scope{}
	s.InitAbstractScopeBase(s, rootNode)
	s.CheckRootScope()
	s.parent = nil
	s.depth = 0
	return s
}

// Java: Scope.untyped()
func (s *Scope) Untyped() *Scope {
	return s
}

// Java: Scope.getDepth()
func (s *Scope) GetDepth() int {
	return s.depth
}

// Java: Scope.getParent()
func (s *Scope) GetParent() *Scope {
	return s.parent
}

// GetVar uses the concrete parent links of Java's final Scope class. The
// generic base must support TypedScope's virtual methods; untyped scopes can
// avoid those dictionary calls and classify implicit names once per lookup.
func (s *Scope) GetVar(name string) *Var {
	implicit := AbstractScope_ImplicitVar_Of(name)
	for {
		if v := s.vars.Get(name); v != nil {
			return v
		}
		if implicit != nil {
			if v := s.GetOwnImplicitSlot(implicit); v != nil {
				return v
			}
		}
		s = s.parent
		if s == nil {
			return nil
		}
	}
}

// Java: AbstractScope.getSlot(String) (final)
func (s *Scope) GetSlot(name string) rhino.StaticSlot {
	if v := s.GetVar(name); v != nil {
		return v
	}
	return nil
}

// Declare declares a variable.
//
// name is the name of the variable, nameNode the NAME node declaring the
// variable, input the input in which this variable is defined.
//
// Java: Scope.declare(String, Node, CompilerInput)
func (s *Scope) Declare(name string, nameNode *rhino.Node, input *CompilerInput) *Var {
	jutil.CheckArgument(name != "")
	// Make sure that it's declared only once
	jutil.CheckState(s.getOwnSlotV(name) == nil)
	var_ := NewVar(
		name,
		nameNode,
		s,
		s.GetVarCount(),
		input,
		/* implicitGoogNamespaceDefinition= */ nil)
	s.DeclareInternal(name, var_)
	return var_
}

// DeclareImplicitGoogNamespaceIfAbsent declares an implicit goog.provide or
// goog.module namespace in this scope.
//
// name is the name of the variable, definition the STRINGLIT node holding
// the full namespace.
//
// Java: Scope.declareImplicitGoogNamespaceIfAbsent(String, Node)
func (s *Scope) DeclareImplicitGoogNamespaceIfAbsent(name string, definition *rhino.Node) *Var {
	jutil.CheckArgument(name != "")
	jutil.CheckStateMsg1(s.IsGlobal(), "Cannot declare implicit goog namespace in local scope %s", s)
	// Allow redeclarations of provides, since they are implicit and don't have a single
	// declaration site.
	var_ := s.getOwnSlotV(name)
	if var_ == nil {
		var_ = Var_CreateImplicitGoogNamespace(name, s, definition)
		s.DeclareInternal(name, var_)
	} else if var_.IsImplicitGoogNamespace() {
		var_.AddImplicitGoogNamespaceDefinition(definition)
	}
	return var_
}

// Java: Scope.makeImplicitVar(ImplicitVar)
func (s *Scope) MakeImplicitVar(var_ AbstractScope_ImplicitVar) *Var {
	return NewVar(
		var_.name(),
		/* nameNode= */ nil,
		s,
		/* index= */ -1,
		/* input= */ nil,
		/* implicitGoogNamespaceDefinition= */ nil)
}
