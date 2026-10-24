// Port of com/google/javascript/jscomp/disambiguate/Invalidation.java @ closure-compiler v20261021

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
// The class has a subclass (WithReceiverType), so it is an interface over
// InvalidationBase (PORTING.md §4). Its only use is GSON serialization of
// the renaming index, which writes a WithReceiverType's own field before
// the base class's: {"receiverType":5,"reason":"INVALIDATING_TYPE"}.
//
// Java: class Invalidation
type Invalidation interface {
	GetInvalidationBase() *InvalidationBase
}

// InvalidationBase holds the fields of Invalidation.
type InvalidationBase struct {
	reason Invalidation_Reason
}

// GetInvalidationBase gives access to the fields of an Invalidation.
func (b *InvalidationBase) GetInvalidationBase() *InvalidationBase { return b }

// Java: Invalidation.wellKnownProperty()
func Invalidation_WellKnownProperty() Invalidation {
	return invalidation_WELL_KNOWN_PROPERTY
}

// Java: Invalidation.invalidatingType(int)
func Invalidation_InvalidatingType(receiverType int) Invalidation {
	return newInvalidation_WithReceiverType(Invalidation_Reason_INVALIDATING_TYPE, receiverType)
}

// Java: Invalidation.undeclaredAccess(int)
func Invalidation_UndeclaredAccess(receiverType int) Invalidation {
	return newInvalidation_WithReceiverType(Invalidation_Reason_UNDECLARED_ACCESS, receiverType)
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
func newInvalidation(reason Invalidation_Reason) *InvalidationBase {
	b := &InvalidationBase{}
	b.InitInvalidationBase(reason)
	return b
}

// InitInvalidationBase is the constructor body, for subclasses.
func (b *InvalidationBase) InitInvalidationBase(reason Invalidation_Reason) {
	b.reason = reason
}

// Invalidation_WithReceiverType is the private subclass WithReceiverType.
// Anonymous classes are not serialized correctly by GSON.
//
// Java: private static final class Invalidation.WithReceiverType extends Invalidation
type Invalidation_WithReceiverType struct {
	receiverType int
	InvalidationBase
}

// Java: Invalidation.WithReceiverType(Reason, int)
func newInvalidation_WithReceiverType(reason Invalidation_Reason, receiverType int) *Invalidation_WithReceiverType {
	w := &Invalidation_WithReceiverType{}
	w.InitInvalidationBase(reason)
	w.receiverType = receiverType
	return w
}

var invalidation_WELL_KNOWN_PROPERTY Invalidation = newInvalidation(Invalidation_Reason_WELL_KNOWN_PROPERTY)
