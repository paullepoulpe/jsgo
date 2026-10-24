// Port of com/google/javascript/jscomp/ScopedName.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"sync/atomic"
)

// ScopedName is a name together with a scope root node. This is suitable
// for use as a map key and is the base class for Var. Note that Equals and
// HashCode as defined at this level, even though subclasses may add
// additional mutable state. The idea is that if two var objects have the
// same name and scope root, then they really are "the same var" and should
// be treated as interchangeable. This does mean that we need to be careful
// to not depend on any other details in situations where vars are stored in
// sets.
type ScopedName interface {
	// GetName returns the name as a string.
	//
	// Java: ScopedName.getName()
	GetName() string

	// GetScopeRoot returns the root node of the scope in which this name is
	// defined.
	//
	// Java: ScopedName.getScopeRoot()
	GetScopeRoot() rhino.Node

	// Java: ScopedName.equals(Object)
	Equals(other any) bool

	// Java: ScopedName.hashCode()
	HashCode() int32
}

// ScopedNameBase holds the final equals/hashCode of ScopedName; subclasses
// embed it and call InitScopedNameBase with themselves.
type ScopedNameBase struct {
	self ScopedName
}

// InitScopedNameBase is the ScopedName() constructor.
func (b *ScopedNameBase) InitScopedNameBase(self ScopedName) { b.self = self }

// GetScopedNameBase gives access to the base through the interface.
func (b *ScopedNameBase) GetScopedNameBase() *ScopedNameBase { return b }

// Java: ScopedName.of(String, Node)
func ScopedName_Of(name string, scopeRoot rhino.Node) ScopedName {
	return newScopedName_Simple(name, scopeRoot)
}

// ScopedName_Simple is the simple implementation with no additional data or
// semantics.
type ScopedName_Simple struct {
	name      string
	scopeRoot rhino.Node
	hash      int32
}

// Java: ScopedName.Simple(String, Node)
func newScopedName_Simple(name string, scopeRoot rhino.Node) *ScopedName_Simple {
	return &ScopedName_Simple{
		name:      name,
		scopeRoot: scopeRoot,
	}
}

// Java: ScopedName.Simple.getName()
func (s *ScopedName_Simple) GetName() string {
	return s.name
}

// Java: ScopedName.Simple.getScopeRoot()
func (s *ScopedName_Simple) GetScopeRoot() rhino.Node {
	return s.scopeRoot
}

// Java: ScopedName.equals(Object)
func (s *ScopedName_Simple) Equals(other any) bool {
	scopedName, ok := other.(ScopedName)
	if !ok || jutil.IsNil(scopedName) {
		return false
	}
	return s.name == scopedName.GetName() &&
		scopedName_nodeEquals(s.scopeRoot, scopedName.GetScopeRoot())
}

// Java: ScopedName.hashCode()
func (s *ScopedName_Simple) HashCode() int32 {
	if h := atomic.LoadInt32(&s.hash); h != 0 {
		return h
	}
	return s.computeHashCode()
}

//go:noinline
func (s *ScopedName_Simple) computeHashCode() int32 {
	if s.scopeRoot.IsNil() {
		panic(jutil.NewNullPointerException())
	}
	h := 31*(31+jutil.StringHashCode(s.name)) + s.scopeRoot.IdentityHashCode()
	atomic.StoreInt32(&s.hash, h)
	return h
}

// Equals is final in Java.
//
// Java: ScopedName.equals(Object)
func (b *ScopedNameBase) Equals(other any) bool {
	scopedName, ok := other.(ScopedName)
	if !ok || jutil.IsNil(scopedName) {
		return false
	}
	return b.self.GetName() == scopedName.GetName() &&
		scopedName_nodeEquals(b.self.GetScopeRoot(), scopedName.GetScopeRoot())
}

// scopedName_nodeEquals is a.equals(b) for Nodes (identity; null a throws).
func scopedName_nodeEquals(a, b rhino.Node) bool {
	if a.IsNil() {
		scopedName_failNull()
	}
	return a == b
}

//go:noinline
func scopedName_failNull() {
	panic(jutil.NewNullPointerException())
}

// HashCode is final in Java.
//
// Java: ScopedName.hashCode()
func (b *ScopedNameBase) HashCode() int32 {
	// NOTE: here we avoid Object.hash because of this function is called
	// frequently enough that the overhead of implicit array creation of a
	// varargs function is enough that we want to avoid it.
	name := b.self.GetName()
	root := b.self.GetScopeRoot()
	if root.IsNil() {
		panic(jutil.NewNullPointerException())
	}
	var result int32 = 1
	result = 31*result + jutil.StringHashCode(name)
	result = 31*result + root.IdentityHashCode() // Go-only: performance — IdentityHashCode(root), unboxed
	return result
}
