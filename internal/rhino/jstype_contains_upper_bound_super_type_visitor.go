// Port of com/google/javascript/rhino/jstype/ContainsUpperBoundSuperTypeVisitor.java @ closure-compiler v20261015

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

import "github.com/paullepoulpe/jsgo/internal/base"

// ContainsUpperBoundSuperTypeVisitor is a type visitor that traverse through
// the referenced types of "forwaring types" to search for a target type.
//
// The result of the visitation indicates whether the target type was
// present, absent, or undcidable due to a cycle.
//
// This is used by TemplateType as a helper in cycle detection and
// subtyping. Cycles are checked for once after TemplateTypes are constructed
// in FunctionTypeBuilder and an error is emitted if they are found.
//
// Java: extends Visitor.WithDefaultCase<Result>; the case methods return
// Results boxed in `any` (see Visitor).
type ContainsUpperBoundSuperTypeVisitor struct {
	Visitor_WithDefaultCaseBase

	target JSType
	seen   *base.LinkedIdentityHashSet[JSType]
}

// Java: ContainsUpperBoundSuperTypeVisitor(JSType)
func NewContainsUpperBoundSuperTypeVisitor(target JSType) *ContainsUpperBoundSuperTypeVisitor {
	v := &ContainsUpperBoundSuperTypeVisitor{}
	v.InitVisitor_WithDefaultCaseBase(v)
	v.seen = base.NewLinkedIdentityHashSet[JSType]()
	v.target = target
	return v
}

// Java: ContainsUpperBoundSuperTypeVisitor.caseDefault(JSType)
func (v *ContainsUpperBoundSuperTypeVisitor) CaseDefault(type_ JSType) any {
	if v.target == nil {
		return ContainsUpperBoundSuperTypeVisitor_Result_ABSENT
	}

	if type_ == v.target {
		return ContainsUpperBoundSuperTypeVisitor_Result_PRESENT
	}
	return ContainsUpperBoundSuperTypeVisitor_Result_ABSENT
}

// Java: ContainsUpperBoundSuperTypeVisitor.caseTemplateType(TemplateType)
func (v *ContainsUpperBoundSuperTypeVisitor) CaseTemplateType(type_ *TemplateType) any {
	return v.CaseForwardingType(type_, type_.GetBound())
}

// Java: ContainsUpperBoundSuperTypeVisitor.caseNamedType(NamedType)
func (v *ContainsUpperBoundSuperTypeVisitor) CaseNamedType(type_ *NamedType) any {
	return v.CaseForwardingType(type_, type_.GetReferencedType())
}

// Java: ContainsUpperBoundSuperTypeVisitor.caseTemplatizedType(TemplatizedType)
func (v *ContainsUpperBoundSuperTypeVisitor) CaseTemplatizedType(type_ *TemplatizedType) any {
	return v.CaseForwardingType(type_, type_.GetReferencedType())
}

// Java: ContainsUpperBoundSuperTypeVisitor.caseUnionType(UnionType)
func (v *ContainsUpperBoundSuperTypeVisitor) CaseUnionType(type_ *UnionType) any {
	if JSType(type_) == v.target {
		return ContainsUpperBoundSuperTypeVisitor_Result_PRESENT
	}

	for _, alt := range type_.GetAlternates() {
		foundInAlt := alt.Visit(v).(ContainsUpperBoundSuperTypeVisitor_Result)
		if foundInAlt != ContainsUpperBoundSuperTypeVisitor_Result_ABSENT {
			return foundInAlt
		}
	}

	return ContainsUpperBoundSuperTypeVisitor_Result_ABSENT
}

// Java: ContainsUpperBoundSuperTypeVisitor.caseForwardingType(JSType, JSType)
func (v *ContainsUpperBoundSuperTypeVisitor) CaseForwardingType(type_ JSType, reference JSType) any {
	if type_ == v.target {
		return ContainsUpperBoundSuperTypeVisitor_Result_PRESENT
	} else if v.seen.Contains(type_) {
		return ContainsUpperBoundSuperTypeVisitor_Result_CYCLE
	} else {
		v.seen.Add(type_)
		return reference.Visit(v)
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
