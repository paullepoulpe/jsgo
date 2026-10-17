// Port of com/google/javascript/jscomp/MemoizedScopeCreator.java @ closure-compiler v20260415

/*
 * Copyright 2017 The Closure Compiler Authors.
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

// MemoizedScopeCreator memoizes a scope creator.
//
// This allows you to make multiple passes, without worrying about the
// expense of generating Scope objects over and over again.
//
// On the other hand, you also have to be more aware of what your passes are
// doing. Scopes are memoized stupidly, so if the underlying tree changes, the
// scope may be out of sync.
type MemoizedScopeCreator struct {
	scopesByScopeRoot *jutil.LinkedHashMap[*rhino.Node, AbstractScope]
	delegate          ScopeCreator
}

// NewMemoizedScopeCreator: delegate is the real source of Scope objects.
//
// Java: MemoizedScopeCreator(ScopeCreator)
func NewMemoizedScopeCreator(delegate ScopeCreator) *MemoizedScopeCreator {
	return &MemoizedScopeCreator{
		scopesByScopeRoot: jutil.NewLinkedHashMap[*rhino.Node, AbstractScope](),
		delegate:          delegate,
	}
}

// Java: MemoizedScopeCreator.createScope(Node, AbstractScope<?, ?>)
func (m *MemoizedScopeCreator) CreateScope(n *rhino.Node, parent AbstractScope) AbstractScope {
	scope := m.scopesByScopeRoot.Get(n)
	if scope == nil {
		scope = m.delegate.CreateScope(n, parent)
		m.scopesByScopeRoot.Put(n, scope)
	} else {
		// Go-only: language-adapter — scope.getParent() on an AbstractScope<?, ?>
		// is the final getParentScope(), which returns getParent().
		jutil.CheckState(parent == scope.GetParentScope())
	}
	return scope
}
