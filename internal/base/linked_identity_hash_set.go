// Port of com/google/javascript/jscomp/base/LinkedIdentityHashSet.java @ closure-compiler v20261015

/*
 * Copyright 2008 The Closure Compiler Authors.
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

package base

// LinkedIdentityHashSet is a deterministic drop-in replacement for
// Sets.newIdentityHashSet.
//
// This class is actually "set-like" in that it doesn't fully implement the
// Set interface. Only the methods needed to eliminate
// Sets.newIdentityHashSet were defined. The entire interface could be
// satisfied in the future if necessary.
type LinkedIdentityHashSet[E comparable] struct {
	innerMap LinkedIdentityHashMap[E, *LinkedIdentityHashSet_Sentinel]
}

// NewLinkedIdentityHashSet is the implicit default constructor.
//
// Java: LinkedIdentityHashSet()
func NewLinkedIdentityHashSet[E comparable]() *LinkedIdentityHashSet[E] {
	return &LinkedIdentityHashSet[E]{}
}

// Contains reports whether element is in the set.
//
// Java: LinkedIdentityHashSet.contains(E)
func (s *LinkedIdentityHashSet[E]) Contains(element E) bool {
	return s.innerMap.Get(element) == LinkedIdentityHashSet_SENTINEL
}

// Add adds element; reports whether it was absent.
//
// Java: LinkedIdentityHashSet.add(E)
func (s *LinkedIdentityHashSet[E]) Add(element E) bool {
	return s.innerMap.Put(element, LinkedIdentityHashSet_SENTINEL) == nil
}

// Remove removes element; reports whether it was present.
//
// Java: LinkedIdentityHashSet.remove(E)
func (s *LinkedIdentityHashSet[E]) Remove(element E) bool {
	return s.innerMap.Put(element, nil) == LinkedIdentityHashSet_SENTINEL
}

// LinkedIdentityHashSet_Sentinel marks membership.
type LinkedIdentityHashSet_Sentinel struct {
	// Go-only: language-adapter — a non-zero-size struct, so that distinct
	// allocations have distinct addresses like distinct Java objects.
	_ byte
}

// Java: LinkedIdentityHashSet.SENTINEL
var LinkedIdentityHashSet_SENTINEL = &LinkedIdentityHashSet_Sentinel{}
