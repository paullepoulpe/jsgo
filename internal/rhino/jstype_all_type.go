// Port of com/google/javascript/rhino/jstype/AllType.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// AllType is the all type, representing all values.
type AllType struct {
	JSTypeBase
}

// Java: AllType(JSTypeRegistry)
func NewAllType(registry *JSTypeRegistry) *AllType {
	t := &AllType{}
	t.objFlags = jsKindAll
	t.InitJSTypeBase(t, registry)
	t.EagerlyResolveToSelf()
	return t
}

// Java: AllType.getTypeClass()
func (t *AllType) GetTypeClass() JSTypeClass {
	return JSTypeClass_ALL
}

// Java: AllType.isAllType()
func (t *AllType) IsAllType() bool {
	return true
}

// Java: AllType.matchesStringContext()
func (t *AllType) MatchesStringContext() bool {
	// Be lenient.
	return true
}

// Java: AllType.matchesObjectContext()
func (t *AllType) MatchesObjectContext() bool {
	// Be lenient.
	return true
}

// Java: AllType.testForEquality(JSType)
func (t *AllType) TestForEquality(that JSType) base.Tri {
	return base.Tri_UNKNOWN
}

// Java: AllType.appendTo(TypeStringBuilder)
func (t *AllType) AppendTo(sb *TypeStringBuilder) {
	sb.Append("*")
}

// Java: AllType.getDisplayName()
func (t *AllType) GetDisplayName() *string {
	return &allType_displayName
}

var allType_displayName = "<Any Type>"

// Java: AllType.hasDisplayName()
func (t *AllType) HasDisplayName() bool {
	return true
}

// Java: AllType.getPossibleToBooleanOutcomes()
func (t *AllType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_BOTH
}

// Java: AllType.resolveInternal(ErrorReporter)
func (t *AllType) ResolveInternal(reporter ErrorReporter) JSType {
	panic(jutil.NewAssertionErrorNoMessage())
}

// Java: AllType.recursionUnsafeHashCode()
func (t *AllType) RecursionUnsafeHashCode() int32 {
	return jutil.IdentityHashCode(t)
}

// Java: AllType.isNullable()
func (t *AllType) IsNullable() bool {
	return true
}

// Java: AllType.isVoidable()
func (t *AllType) IsVoidable() bool {
	return true
}
