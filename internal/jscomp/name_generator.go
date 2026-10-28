// Port of com/google/javascript/jscomp/NameGenerator.java @ closure-compiler v20261021

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
	"iter"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// stringMapSet adapts a plain Go map[string]struct{} to jutil.Set[string]
// without LinkedHashSet entry/order slice overhead.
type stringMapSet map[string]struct{}

func (s stringMapSet) Size() int              { return len(s) }
func (s stringMapSet) IsEmpty() bool          { return len(s) == 0 }
func (s stringMapSet) Contains(e string) bool { _, ok := s[e]; return ok }
func (s stringMapSet) All() iter.Seq[string] {
	return func(yield func(string) bool) {
		// order-independent: unordered set view used only for membership snapshots
		for k := range s {
			if !yield(k) {
				return
			}
		}
	}
}

// setNodeStringID updates a StringNode's payload to pre-interned id without
// re-interning, while preserving SourceInformationAnnotator's lazy OriginalName
// spill when kindFlagOrigNameIsString is set.
func setNodeStringID(n rhino.Node, id strpool.ID) {
	if id == 0 {
		panic(jutil.NewNullPointerException())
	}
	curID := n.GetStringID()
	if curID == uint32(id) {
		return
	}
	if n.GetOriginalNameID() == curID {
		n.SetString(strpool.String(id))
		return
	}
	n.SetStringID(id)
}

// NameGenerator is a class that generates unique JavaScript
// variable/property names.
//
// Go-only: language-adapter — the Java Set parameters are jutil.Set views:
// reservedNames is referenced, not copied (it may be a live view such as a
// map's key set).
type NameGenerator interface {
	// Reset reconfigures this NameGenerator, and resets it to the initial
	// state.
	//
	// reservedNames is the set of names that are reserved; generated names
	// will not include these names. This set is referenced rather than
	// copied, so changes to the set will be reflected in how names are
	// generated. All generated names begin with prefix. reservedCharacters
	// are characters that won't be used in generated names.
	//
	// Java: NameGenerator.reset(Set<String>, String, Set<Character>)
	Reset(reservedNames jutil.Set[string], prefix string, reservedCharacters jutil.Set[uint16])

	// ResetWithFirstAndNonFirstChars reconfigures this NameGenerator, and
	// resets it to the initial state.
	//
	// reservedNames is the set of names that are reserved; generated names
	// will not include these names. This set is referenced rather than
	// copied, so changes to the set will be reflected in how names are
	// generated. All generated names begin with prefix. If specified,
	// reservedFirstCharacters won't be used as the first character in
	// generated names, and reservedNonFirstCharacters won't be used for
	// characters (after the first) in generated names.
	//
	// Java: NameGenerator.reset(Set<String>, String, Set<Character>, Set<Character>)
	ResetWithFirstAndNonFirstChars(
		reservedNames jutil.Set[string],
		prefix string,
		reservedFirstCharacters jutil.Set[uint16],
		reservedNonFirstCharacters jutil.Set[uint16])

	// Clone returns a clone of this NameGenerator, reconfigured and reset.
	//
	// Java: NameGenerator.clone(Set<String>, String, Set<Character>)
	Clone(reservedNames jutil.Set[string], prefix string, reservedCharacters jutil.Set[uint16]) NameGenerator

	// GenerateNextName generates the next name.
	//
	// Java: NameGenerator.generateNextName()
	GenerateNextName() string
}
