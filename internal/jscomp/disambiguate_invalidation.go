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

import "github.com/paullepoulpe/jsgo/internal/diagnostic"

// Invalidation describes one way in which a property became invalidated.
//
// This information is only used diagnostically and doesn't affect the
// outcome of optimizations. Since diagnostic information isn't usually
// rendered, this value struct avoids heap allocations during disambiguation
// and formats its GSON representation on demand via diagnostic.LogsGson_Able.
//
// Java: class Invalidation
type Invalidation struct {
	reason       Invalidation_Reason
	receiverType int32
	valid        bool
	hasReceiver  bool
}

// IsValid reports whether this Invalidation has been populated.
func (i Invalidation) IsValid() bool { return i.valid }

type invalidationJson struct {
	Reason Invalidation_Reason `gson:"reason"`
}

type invalidationWithReceiverJson struct {
	ReceiverType int                 `gson:"receiverType"`
	Reason       Invalidation_Reason `gson:"reason"`
}

// ToLogsGson formats Invalidation identically to Java's Invalidation and
// Invalidation.WithReceiverType ({"receiverType":5,"reason":"INVALIDATING_TYPE"}).
func (i Invalidation) ToLogsGson() any {
	if i.hasReceiver {
		return invalidationWithReceiverJson{
			ReceiverType: int(i.receiverType),
			Reason:       i.reason,
		}
	}
	return invalidationJson{Reason: i.reason}
}

var _ diagnostic.LogsGson_Able = Invalidation{}

// Java: Invalidation.wellKnownProperty()
func Invalidation_WellKnownProperty() Invalidation {
	return Invalidation{
		reason: Invalidation_Reason_WELL_KNOWN_PROPERTY,
		valid:  true,
	}
}

// Java: Invalidation.invalidatingType(int)
func Invalidation_InvalidatingType(receiverType int) Invalidation {
	return Invalidation{
		reason:       Invalidation_Reason_INVALIDATING_TYPE,
		receiverType: int32(receiverType),
		valid:        true,
		hasReceiver:  true,
	}
}

// Java: Invalidation.undeclaredAccess(int)
func Invalidation_UndeclaredAccess(receiverType int) Invalidation {
	return Invalidation{
		reason:       Invalidation_Reason_UNDECLARED_ACCESS,
		receiverType: int32(receiverType),
		valid:        true,
		hasReceiver:  true,
	}
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

// Name is Enum.name().
func (r Invalidation_Reason) Name() string { return invalidation_Reason_names[r] }

// String is Enum.toString().
func (r Invalidation_Reason) String() string { return r.Name() }

// Ordinal is Enum.ordinal().
func (r Invalidation_Reason) Ordinal() int { return int(r) }
