// Port of test/com/google/javascript/jscomp/DiagnosticGroupTest.java @ closure-compiler v20260415

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
