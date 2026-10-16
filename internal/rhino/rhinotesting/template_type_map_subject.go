// Port of com/google/javascript/rhino/testing/TemplateTypeMapSubject.java @ closure-compiler v20260415

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

package rhinotesting

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// TemplateTypeMapSubject is a Truth Subject for TemplateTypeMap.
type TemplateTypeMapSubject struct {
	t      testing.TB
	actual *rhino.TemplateTypeMap
}

// Java: TemplateTypeMapSubject.assertThat(TemplateTypeMap)
func TemplateTypeMapSubject_AssertThat(t testing.TB, actual *rhino.TemplateTypeMap) *TemplateTypeMapSubject {
	return &TemplateTypeMapSubject{t: t, actual: actual}
}

// Java: TemplateTypeMapSubject.hasKeysAndValues(ImmutableList<TemplateType>, ImmutableList<? extends JSType>)
func (s *TemplateTypeMapSubject) HasKeysAndValues(keys []*rhino.TemplateType, values []rhino.JSType) {
	s.t.Helper()
	tr := Truth{T: s.t}
	tr.AssertWithMessage("value of: templateTypeMap.getTemplateKeys()").That(s.actual.GetTemplateKeys()).IsEqualTo(keys)
	tr.AssertWithMessage("value of: templateTypeMap.getTemplateValues()").That(s.actual.GetTemplateValues()).IsEqualTo(values)
}

// Java: TemplateTypeMapSubject.hasSubmapsAt(int...)
func (s *TemplateTypeMapSubject) HasSubmapsAt(indices ...int) {
	s.t.Helper()
	Truth{T: s.t}.AssertWithMessage("value of: templateTypeMap.getIndicesOfSubmapsForTesting()").
		That(s.actual.GetIndicesOfSubmapsForTesting()).
		ContainsExactlyElementsIn(indices)
}
