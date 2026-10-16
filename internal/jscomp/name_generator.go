// Port of com/google/javascript/jscomp/NameGenerator.java @ closure-compiler v20260415

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// NameGenerator is a class that generates unique JavaScript
// variable/property names.
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
	Reset(reservedNames *jutil.LinkedHashSet[string], prefix string, reservedCharacters *jutil.LinkedHashSet[uint16])

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
		reservedNames *jutil.LinkedHashSet[string],
		prefix string,
		reservedFirstCharacters *jutil.LinkedHashSet[uint16],
		reservedNonFirstCharacters *jutil.LinkedHashSet[uint16])

	// Clone returns a clone of this NameGenerator, reconfigured and reset.
	//
	// Java: NameGenerator.clone(Set<String>, String, Set<Character>)
	Clone(reservedNames *jutil.LinkedHashSet[string], prefix string, reservedCharacters *jutil.LinkedHashSet[uint16]) NameGenerator

	// GenerateNextName generates the next name.
	//
	// Java: NameGenerator.generateNextName()
	GenerateNextName() string
}
