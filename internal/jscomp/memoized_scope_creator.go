// Port of com/google/javascript/jscomp/MemoizedScopeCreator.java @ closure-compiler v20261015

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

type memoizedScopeEntry struct {
	node  *rhino.Node
	scope AbstractScope
}

// MemoizedScopeCreator memoizes a scope creator.
//
// This allows you to make multiple passes, without worrying about the
// expense of generating Scope objects over and over again.
//
// On the other hand, you also have to be more aware of what your passes are
// doing. Scopes are memoized stupidly, so if the underlying tree changes, the
// scope may be out of sync.
type MemoizedScopeCreator struct {
	scopesByScopeRoot map[*rhino.Node]AbstractScope
	small             [8]memoizedScopeEntry
	smallLen          int
	lastNode          *rhino.Node
	lastScope         AbstractScope
	delegate          ScopeCreator
}

// NewMemoizedScopeCreator: delegate is the real source of Scope objects.
//
// Java: MemoizedScopeCreator(ScopeCreator)
func NewMemoizedScopeCreator(delegate ScopeCreator) *MemoizedScopeCreator {
	return &MemoizedScopeCreator{
		delegate: delegate,
	}
}

// Java: MemoizedScopeCreator.createScope(Node, AbstractScope<?, ?>)
func (m *MemoizedScopeCreator) CreateScope(n *rhino.Node, parent AbstractScope) AbstractScope {
	if n == m.lastNode && m.lastScope != nil {
		jutil.CheckState(parent == m.lastScope.GetParentScope())
		return m.lastScope
	}
	var scope AbstractScope
	if m.scopesByScopeRoot != nil {
		scope = m.scopesByScopeRoot[n]
	} else {
		for i := 0; i < m.smallLen; i++ {
			if m.small[i].node == n {
				scope = m.small[i].scope
				break
			}
		}
	}
	if scope == nil {
		scope = m.delegate.CreateScope(n, parent)
		if m.scopesByScopeRoot != nil {
			m.scopesByScopeRoot[n] = scope
		} else if m.smallLen < len(m.small) {
			m.small[m.smallLen] = memoizedScopeEntry{node: n, scope: scope}
			m.smallLen++
		} else {
			mp := make(map[*rhino.Node]AbstractScope, 1024)
			for i := 0; i < len(m.small); i++ {
				mp[m.small[i].node] = m.small[i].scope
			}
			mp[n] = scope
			m.scopesByScopeRoot = mp
			m.small = [8]memoizedScopeEntry{}
		}
	} else {
		// Go-only: language-adapter — scope.getParent() on an AbstractScope<?, ?>
		// is the final getParentScope(), which returns getParent().
		jutil.CheckState(parent == scope.GetParentScope())
	}
	m.lastNode = n
	m.lastScope = scope
	return scope
}

func (m *MemoizedScopeCreator) ReleaseArena() {
	if m == nil {
		return
	}
	if m.scopesByScopeRoot != nil {
		if len(m.scopesByScopeRoot) <= 65536 {
			clear(m.scopesByScopeRoot)
		} else {
			m.scopesByScopeRoot = nil
		}
	}
	m.small = [8]memoizedScopeEntry{}
	m.smallLen = 0
	m.lastNode = nil
	m.lastScope = nil
	if sc, ok := m.delegate.(*SyntacticScopeCreator); ok && sc != nil {
		sc.ReleaseArena()
	}
}
