// Port of test/com/google/javascript/jscomp/DiagnosticGroupTest.java @ closure-compiler v20260415

/*
 * Copyright 2009 The Closure Compiler Authors.
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

import "testing"

func TestDiagnosticGroup_Registration(t *testing.T) {
	dg := NewDiagnosticGroups()
	_ = dg // Java calls the static forName through this instance.
	jsaEq(t, DiagnosticGroups_ForName("deprecated"), DiagnosticGroups_DEPRECATED)
}

// Go-only: extra coverage — groups, composite groups and forType singletons.
func TestDiagnosticGroup_Basics(t *testing.T) {
	a := DiagnosticType_Warning("JSC_DG_A", "a")
	b := DiagnosticType_Warning("JSC_DG_B", "b")
	g1 := NewDiagnosticGroupWithName(jsaStr("g1"), a, b, a)
	jsaEq(t, g1.GetTypes().Size(), 2)
	jsaTrue(t, g1.MatchesType(DiagnosticType_Warning("JSC_DG_A", "other format")))
	jsaTrue(t, g1.Matches(JSError_Make(b)))
	jsaEq(t, g1.String(), "DiagnosticGroup<g1>")
	composite := NewDiagnosticGroupFromGroups(g1, NewDiagnosticGroup(DiagnosticType_Error("JSC_DG_C", "c")))
	jsaEq(t, composite.GetTypes().Size(), 3)
	jsaTrue(t, composite.GetName() == nil)
	jsaTrue(t, DiagnosticGroup_ForType(a) == DiagnosticGroup_ForType(a))
	jsaFalse(t, DiagnosticGroup_ForType(a).MatchesType(b))
}
