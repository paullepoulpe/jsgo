// Port of com/google/javascript/rhino/jstype/ContainsUpperBoundSuperTypeVisitor.java @ closure-compiler v20261021

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import (
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// containsUpperBoundSuperTypePresent reports whether target appears in the
// reference chain of forwarding types reachable from type_.
//
// Java: new ContainsUpperBoundSuperTypeVisitor(target) == Result.PRESENT
func containsUpperBoundSuperTypePresent(type_ JSType, target JSType) bool {
	return containsUpperBoundSuperType(type_, target) == ContainsUpperBoundSuperTypeVisitor_Result_PRESENT
}

// containsUpperBoundSuperTypeCycle reports whether the reference chain of
// forwarding types reachable from type_ contains a cycle.
//
// Java: new ContainsUpperBoundSuperTypeVisitor(null) == Result.CYCLE
func containsUpperBoundSuperTypeCycle(type_ JSType) bool {
	return containsUpperBoundSuperType(type_, nil) == ContainsUpperBoundSuperTypeVisitor_Result_CYCLE
}

func containsUpperBoundSuperType(type_ JSType, target JSType) ContainsUpperBoundSuperTypeVisitor_Result {
	var seenBuf [8]JSType
	seen := seenBuf[:0]
	return containsUpperBoundSuperTypeRec(type_, target, &seen)
}

func containsUpperBoundSuperTypeRec(type_ JSType, target JSType, seen *[]JSType) ContainsUpperBoundSuperTypeVisitor_Result {
	switch x := type_.(type) {
	case *TemplateType:
		return containsUpperBoundForwarding(x, x.GetBound(), target, seen)
	case *NamedType:
		return containsUpperBoundForwarding(x, x.GetReferencedType(), target, seen)
	case *TemplatizedType:
		return containsUpperBoundForwarding(x, x.GetReferencedType(), target, seen)
	case *UnionType:
		if JSType(x) == target {
			return ContainsUpperBoundSuperTypeVisitor_Result_PRESENT
		}
		for _, alt := range x.GetAlternates() {
			foundInAlt := containsUpperBoundSuperTypeRec(alt, target, seen)
			if foundInAlt != ContainsUpperBoundSuperTypeVisitor_Result_ABSENT {
				return foundInAlt
			}
		}
		return ContainsUpperBoundSuperTypeVisitor_Result_ABSENT
	case *AllType, *BooleanType, *NoObjectTypeBase, *UnknownType, *NullType,
		*NumberType, *BigIntType, *StringType, *VoidType:
		// Visitor.WithDefaultCase passes nil for zero-arg cases, so CaseDefault
		// always returns ABSENT for these types.
		return ContainsUpperBoundSuperTypeVisitor_Result_ABSENT
	case *NoType, *EnumElementType, *FunctionTypeBase,
		*InstanceObjectType, *PrototypeObjectTypeBase, *RecordType, *EnumType, *NoResolvedType,
		*ProxyObjectTypeBase, *SymbolTypeBase, *KnownSymbolType:
		if target != nil && type_ == target {
			return ContainsUpperBoundSuperTypeVisitor_Result_PRESENT
		}
		return ContainsUpperBoundSuperTypeVisitor_Result_ABSENT
	default:
		if _, ok := type_.(NoObjectType); ok {
			return ContainsUpperBoundSuperTypeVisitor_Result_ABSENT
		}
		switch type_.(type) {
		case FunctionType, ProxyObjectType, ObjectType, SymbolType:
			if target != nil && type_ == target {
				return ContainsUpperBoundSuperTypeVisitor_Result_PRESENT
			}
			return ContainsUpperBoundSuperTypeVisitor_Result_ABSENT
		}
		panic(jutil.NewUnsupportedOperationException())
	}
}

// Java: ContainsUpperBoundSuperTypeVisitor.caseForwardingType(JSType, JSType)
func containsUpperBoundForwarding(type_ JSType, reference JSType, target JSType, seen *[]JSType) ContainsUpperBoundSuperTypeVisitor_Result {
	if type_ == target {
		return ContainsUpperBoundSuperTypeVisitor_Result_PRESENT
	} else if slices.Contains(*seen, type_) {
		return ContainsUpperBoundSuperTypeVisitor_Result_CYCLE
	} else {
		*seen = append(*seen, type_)
		return containsUpperBoundSuperTypeRec(reference, target, seen)
	}
}

// ContainsUpperBoundSuperTypeVisitor_Result represents the outcome of a
// visitation of the ContainsUpperBoundSuperTypeVisitor.
type ContainsUpperBoundSuperTypeVisitor_Result int32

const (
	ContainsUpperBoundSuperTypeVisitor_Result_PRESENT ContainsUpperBoundSuperTypeVisitor_Result = iota
	ContainsUpperBoundSuperTypeVisitor_Result_ABSENT
	// Containment is undecidable due to a reference cycle.
	ContainsUpperBoundSuperTypeVisitor_Result_CYCLE
)

var containsUpperBoundSuperTypeVisitor_Result_names = [...]string{"PRESENT", "ABSENT", "CYCLE"}

// String is Enum.name().
func (r ContainsUpperBoundSuperTypeVisitor_Result) String() string {
	return containsUpperBoundSuperTypeVisitor_Result_names[r]
}
