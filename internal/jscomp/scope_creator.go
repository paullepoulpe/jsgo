// Port of com/google/javascript/jscomp/ScopeCreator.java @ closure-compiler v20261021

/*
 * Copyright 2006 The Closure Compiler Authors.
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

// ScopeCreator defines how objects capable of creating scopes from the parse
// tree behave.
type ScopeCreator interface {
	// CreateScope creates a Scope object.
	//
	// n is the root node (either a FUNCTION node, a SCRIPT node, or a
	// synthetic block node whose children are all SCRIPT nodes); parent is
	// the parent Scope object (may be null).
	//
	// Go-only: language-adapter — implementations whose Java override narrows
	// the result (SyntacticScopeCreator returns Scope, TypedScopeCreator
	// TypedScope) keep this signature; callers assert (`.(*Scope)`).
	//
	// Java: ScopeCreator.createScope(Node, AbstractScope<?, ?>)
	CreateScope(n rhino.Node, parent AbstractScope) AbstractScope
}

// ScopeCreator_ASSERT_NO_SCOPES_CREATED is a scope creator that can be
// provided to NodeTraversal to ensure that no scopes are actually ever
// created. This is in contrast to the default SyntacticScopeCreator, which
// will create scopes on demand.
var ScopeCreator_ASSERT_NO_SCOPES_CREATED ScopeCreator = &scopeCreator_assertNoScopesCreated{}

// scopeCreator_assertNoScopesCreated is the anonymous class of
// ASSERT_NO_SCOPES_CREATED.
type scopeCreator_assertNoScopesCreated struct {
	_ byte // Go-only: language-adapter — distinct identity
}

// Java: ScopeCreator.ASSERT_NO_SCOPES_CREATED.createScope(Node, AbstractScope<?, ?>)
func (*scopeCreator_assertNoScopesCreated) CreateScope(n rhino.Node, parent AbstractScope) AbstractScope {
	panic(jutil.NewIllegalStateExceptionMsg(
		"ScopeCreator.ASSERT_NO_SCOPES_CREATED cannot create child scopes."))
}
