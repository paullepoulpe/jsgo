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

type memoizedScopeCache struct {
	escaped           bool
	scopesByScopeRoot map[*rhino.Node]AbstractScope
	small             [8]memoizedScopeEntry
	smallLen          int
	lastNode          *rhino.Node
	lastScope         AbstractScope
}

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
	// Keep the entire mutable cache behind one pointer: initialized Go value
	// copies shared their memoized entries before the small-cache optimization.
	cache    *memoizedScopeCache
	delegate ScopeCreator
}

// NewMemoizedScopeCreator: delegate is the real source of Scope objects.
//
// Java: MemoizedScopeCreator(ScopeCreator)
func NewMemoizedScopeCreator(delegate ScopeCreator) *MemoizedScopeCreator {
	// Coallocate the shared cache with its initial creator without publishing a
	// self pointer. Copies keep the cache owner alive and observe later entries.
	owner := new(struct {
		creator MemoizedScopeCreator
		cache   memoizedScopeCache
	})
	owner.creator.cache = &owner.cache
	owner.creator.delegate = delegate
	return &owner.creator
}

// Java: MemoizedScopeCreator.createScope(Node, AbstractScope<?, ?>)
func (m *MemoizedScopeCreator) CreateScope(n *rhino.Node, parent AbstractScope) AbstractScope {
	if n == m.cache.lastNode && m.cache.lastScope != nil {
		jutil.CheckState(parent == m.cache.lastScope.GetParentScope())
		return m.cache.lastScope
	}
	var scope AbstractScope
	if m.cache.scopesByScopeRoot != nil {
		scope = m.cache.scopesByScopeRoot[n]
	} else {
		for i := 0; i < m.cache.smallLen; i++ {
			if m.cache.small[i].node == n {
				scope = m.cache.small[i].scope
				break
			}
		}
	}
	if scope == nil {
		scope = m.delegate.CreateScope(n, parent)
		if m.cache.scopesByScopeRoot != nil {
			m.cache.scopesByScopeRoot[n] = scope
		} else if m.cache.smallLen < len(m.cache.small) {
			m.cache.small[m.cache.smallLen] = memoizedScopeEntry{node: n, scope: scope}
			m.cache.smallLen++
		} else {
			mp := make(map[*rhino.Node]AbstractScope, 1024)
			for i := 0; i < len(m.cache.small); i++ {
				mp[m.cache.small[i].node] = m.cache.small[i].scope
			}
			mp[n] = scope
			m.cache.scopesByScopeRoot = mp
			m.cache.small = [8]memoizedScopeEntry{}
		}
	} else {
		// Go-only: language-adapter — scope.getParent() on an AbstractScope<?, ?>
		// is the final getParentScope(), which returns getParent().
		jutil.CheckState(parent == scope.GetParentScope())
	}
	m.cache.lastNode = n
	m.cache.lastScope = scope
	return scope
}

func (m *MemoizedScopeCreator) markArenaEscaped() {
	if m == nil {
		return
	}
	if m.cache != nil {
		m.cache.escaped = true
	}
	switch delegate := m.delegate.(type) {
	case *SyntacticScopeCreator:
		delegate.markArenaEscaped()
	case *MemoizedScopeCreator:
		delegate.markArenaEscaped()
	}
}

// ReleaseArena clears a private memoized cache and releases its delegate's
// scratch storage after use. Published creators retain their cache identity.
func (m *MemoizedScopeCreator) ReleaseArena() {
	if m == nil {
		return
	}
	if m.cache == nil || m.cache.escaped {
		return
	}
	if m.cache.scopesByScopeRoot != nil {
		if len(m.cache.scopesByScopeRoot) <= 65536 {
			clear(m.cache.scopesByScopeRoot)
		} else {
			m.cache.scopesByScopeRoot = nil
		}
	}
	m.cache.small = [8]memoizedScopeEntry{}
	m.cache.smallLen = 0
	m.cache.lastNode = nil
	m.cache.lastScope = nil
	if sc, ok := m.delegate.(*SyntacticScopeCreator); ok && sc != nil {
		sc.ReleaseArena()
	}
}
