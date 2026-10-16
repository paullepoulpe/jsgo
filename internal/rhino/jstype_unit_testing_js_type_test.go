// Port of com/google/javascript/rhino/jstype/UnitTestingJSType.java @ closure-compiler v20260415

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
 *   Nick Santos
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino_test

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	. "github.com/paullepoulpe/jsgo/internal/rhino"
)

// unitTestingJSType is a minimal implementation of JSType for unit tests
// and nothing else.
//
// This class has no innate behaviour. It is intended as a stand-in for
// testing behaviours on JSType that require a concrete instance. Test cases
// are responsible for any configuration.
//
// Go-only: language-adapter — the tests' anonymous subclasses override
// isNoResolvedType() and getTypeClass(); those overrides are the optional
// fields.
type unitTestingJSType struct {
	JSTypeBase
	noResolvedType *bool
	typeClass      *JSTypeClass
}

// Java: UnitTestingJSType(JSTypeRegistry)
func newUnitTestingJSType(registry *JSTypeRegistry) *unitTestingJSType {
	u := &unitTestingJSType{}
	u.InitJSTypeBase(u, registry)
	return u
}

func (u *unitTestingJSType) IsNoResolvedType() bool {
	if u.noResolvedType != nil {
		return *u.noResolvedType
	}
	return u.JSTypeBase.IsNoResolvedType()
}

// Java: UnitTestingJSType.getTypeClass()
func (u *unitTestingJSType) GetTypeClass() JSTypeClass {
	if u.typeClass != nil {
		return *u.typeClass
	}
	panic(jutil.NewUnsupportedOperationException())
}

// Java: UnitTestingJSType.recursionUnsafeHashCode()
func (u *unitTestingJSType) RecursionUnsafeHashCode() int32 {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: UnitTestingJSType.getPossibleToBooleanOutcomes()
func (u *unitTestingJSType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: UnitTestingJSType.visit(Visitor<T>)
func (u *unitTestingJSType) Visit(visitor Visitor) any {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: UnitTestingJSType.visit(RelationshipVisitor<T>, JSType)
func (u *unitTestingJSType) VisitRelationship(visitor RelationshipVisitor, that JSType) any {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: UnitTestingJSType.resolveInternal(ErrorReporter)
func (u *unitTestingJSType) ResolveInternal(reporter ErrorReporter) JSType {
	return u
}

// Java: UnitTestingJSType.appendTo(TypeStringBuilder)
func (u *unitTestingJSType) AppendTo(sb *TypeStringBuilder) {
	panic(jutil.NewUnsupportedOperationException())
}
