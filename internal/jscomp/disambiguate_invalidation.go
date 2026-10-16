// Port of com/google/javascript/jscomp/disambiguate/Invalidation.java @ closure-compiler v20260415

/*
 * Copyright 2020 The Closure Compiler Authors.
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

// Invalidation describes one way in which a property became invalidated.
//
// This information is only used diagnostically and doesn't affect the
// outcome of optimizations. Since diagnostic information isn't usually
// rendered, this class is designed to defer work as much as possible.
//
// Go-only: language-adapter — the private subclass WithRecieverType only
// adds a field, so it is folded into this struct (withRecieverType tells
// the two apart).
type Invalidation struct {
	reason Invalidation_Reason

	// WithRecieverType.recieverType
	recieverType     int
	withRecieverType bool
}

// Java: Invalidation.wellKnownProperty()
func Invalidation_WellKnownProperty() *Invalidation {
	return invalidation_WELL_KNOWN_PROPERTY
}

// Java: Invalidation.invalidatingType(int)
func Invalidation_InvalidatingType(recieverType int) *Invalidation {
	return newInvalidation_WithRecieverType(Invalidation_Reason_INVALIDATING_TYPE, recieverType)
}

// Java: Invalidation.undeclaredAccess(int)
func Invalidation_UndeclaredAccess(recieverType int) *Invalidation {
	return newInvalidation_WithRecieverType(Invalidation_Reason_UNDECLARED_ACCESS, recieverType)
}

// Invalidation_Reason is Invalidation.Reason.
type Invalidation_Reason int32

const (
	// Certain well-known properties like "prototype" are always
	// invalidated.
	Invalidation_Reason_WELL_KNOWN_PROPERTY Invalidation_Reason = iota
	// Properties accessed on invalidating types (like Object) are
	// invalidated.
	Invalidation_Reason_INVALIDATING_TYPE
	// Properties accessed on types that don't declare them are invalidated.
	Invalidation_Reason_UNDECLARED_ACCESS
)

var invalidation_Reason_names = [...]string{"WELL_KNOWN_PROPERTY", "INVALIDATING_TYPE", "UNDECLARED_ACCESS"}

// Invalidation_Reason_values is Reason.values().
func Invalidation_Reason_values() []Invalidation_Reason {
	return []Invalidation_Reason{
		Invalidation_Reason_WELL_KNOWN_PROPERTY,
		Invalidation_Reason_INVALIDATING_TYPE,
		Invalidation_Reason_UNDECLARED_ACCESS,
	}
}

// Name is Enum.name().
func (r Invalidation_Reason) Name() string { return invalidation_Reason_names[r] }

// String is Enum.toString().
func (r Invalidation_Reason) String() string { return r.Name() }

// Ordinal is Enum.ordinal().
func (r Invalidation_Reason) Ordinal() int { return int(r) }

// Java: Invalidation(Reason)
func newInvalidation(reason Invalidation_Reason) *Invalidation {
	return &Invalidation{reason: reason}
}

// newInvalidation_WithRecieverType is WithRecieverType(Reason, int).
// Anonymous classes are not serialized correctly by GSON.
//
// Java: Invalidation.WithRecieverType(Reason, int)
func newInvalidation_WithRecieverType(reason Invalidation_Reason, recieverType int) *Invalidation {
	i := newInvalidation(reason)
	i.recieverType = recieverType
	i.withRecieverType = true
	return i
}

var invalidation_WELL_KNOWN_PROPERTY = newInvalidation(Invalidation_Reason_WELL_KNOWN_PROPERTY)
