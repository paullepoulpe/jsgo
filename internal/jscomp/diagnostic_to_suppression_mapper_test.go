// Port of test/com/google/javascript/jscomp/DiagnosticToSuppressionMapperTest.java @ closure-compiler v20260415

package jscomp

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
)

func jsaStr(s string) *string { return &s }

func jsaTreeMapString(m *jutil.TreeMap[string, string]) string { return m.String() }

func TestDiagnosticToSuppressionMapper_CreateSuppressionMap_asJson_skipsNonParseableSuppressions(t *testing.T) {
	testError1 := DiagnosticType_Error("JSC_TEST_ERROR_1", "testing")
	testError2 := DiagnosticType_Error("JSC_TEST_ERROR_2", "testing")
	diagnosticGroups := jutil.NewLinkedHashMap[string, *DiagnosticGroup]()
	diagnosticGroups.Put("validSuppress", DiagnosticGroup_ForType(testError1))
	diagnosticGroups.Put("invalidSuppress", DiagnosticGroup_ForType(testError2))

	suppressions := NewDiagnosticToSuppressionMapper(jutil.NewLinkedHashSet("validSuppress"), diagnosticGroups).
		CreateSuppressionMap()

	jsaEq(t, jsaTreeMapString(suppressions), "{"+testError1.key+"=validSuppress}")
}

func TestDiagnosticToSuppressionMapper_CreateSuppressionMap_choosesArbitrarySuppressionGivenMultiple(t *testing.T) {
	testError := DiagnosticType_Error("JSC_TEST_ERROR", "testing")
	diagnosticGroups := jutil.NewLinkedHashMap[string, *DiagnosticGroup]()
	diagnosticGroups.Put("suppress0", DiagnosticGroup_ForType(testError))
	diagnosticGroups.Put("suppress1", DiagnosticGroup_ForType(testError))

	suppressions := NewDiagnosticToSuppressionMapper(
		jutil.NewLinkedHashSet("suppress0", "suppress1"), diagnosticGroups).
		CreateSuppressionMap()

	jsaEq(t, jsaTreeMapString(suppressions), "{"+testError.key+"=suppress1}")
}

func TestDiagnosticToSuppressionMapper_PrintSuppressionMap_asJson_printsValidJson(t *testing.T) {
	var myOut bytes.Buffer
	old := jutil.System_Out()
	jutil.System_SetOut(jutil.NewPrintStream(&myOut))
	defer jutil.System_SetOut(old)

	mapper := NewDiagnosticToSuppressionMapper(
		parsing.ParserRunner_GetSuppressionNames(), DiagnosticGroups_GetRegisteredGroups())
	mapper.PrintSuppressions(DiagnosticToSuppressionMapper_OutputFormat_JSON)

	var suppressions map[string]string
	if err := json.Unmarshal(myOut.Bytes(), &suppressions); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	// Java: assertThat(suppressions).isEqualTo(mapper.createSuppressionMap()) (Map equality)
	expected := mapper.CreateSuppressionMap()
	jsaEq(t, len(suppressions), expected.Size())
	for k, v := range expected.All() {
		got, ok := suppressions[k]
		if !ok || got != v {
			t.Errorf("suppressions[%q] = %q (present %v), want %q", k, got, ok, v)
		}
	}
}

// Go-only: extra coverage — the JSON output for a small map.
func TestDiagnosticToSuppressionMapper_PrintSuppressionMap_asJson_small(t *testing.T) {
	testError1 := DiagnosticType_Error("JSC_TEST_ERROR_1", "testing")
	testError2 := DiagnosticType_Error("JSC_A<B", "testing")
	diagnosticGroups := jutil.NewLinkedHashMap[string, *DiagnosticGroup]()
	diagnosticGroups.Put("s1", DiagnosticGroup_ForType(testError1))
	diagnosticGroups.Put("s2", DiagnosticGroup_ForType(testError2))
	var myOut bytes.Buffer
	old := jutil.System_Out()
	jutil.System_SetOut(jutil.NewPrintStream(&myOut))
	defer jutil.System_SetOut(old)
	NewDiagnosticToSuppressionMapper(jutil.NewLinkedHashSet("s1", "s2"), diagnosticGroups).
		PrintSuppressions(DiagnosticToSuppressionMapper_OutputFormat_JSON)
	jsaEq(t, myOut.String(), `{"JSC_A\u003cB":"s2","JSC_TEST_ERROR_1":"s1"}`+"\n")
}

func TestDiagnosticToSuppressionMapper_PrintSuppressionMap_asMarkdown_printsSortedTable(t *testing.T) {
	testError1 := DiagnosticType_Error("JSC_TEST_ERROR_1", "testing")
	testError2 := DiagnosticType_Error("JSC_TEST_ERROR_2", "testing")
	testError3 := DiagnosticType_Error("JSC_OTHER_ERROR", "testing")
	diagnosticGroups := jutil.NewLinkedHashMap[string, *DiagnosticGroup]()
	diagnosticGroups.Put("suppressionA", NewDiagnosticGroupWithName(jsaStr("suppressionA"), testError2, testError3))
	diagnosticGroups.Put("suppressionB", DiagnosticGroup_ForType(testError1))

	var myOut bytes.Buffer
	old := jutil.System_Out()
	jutil.System_SetOut(jutil.NewPrintStream(&myOut))
	defer jutil.System_SetOut(old)

	mapper := NewDiagnosticToSuppressionMapper(
		jutil.NewLinkedHashSet("suppressionA", "suppressionB"), diagnosticGroups)
	mapper.PrintSuppressions(DiagnosticToSuppressionMapper_OutputFormat_MD)

	jsaEq(t, myOut.String(),
		"| Error | Suppression tag |\n"+
			"|---|---|\n"+
			"|JSC_OTHER_ERROR|suppressionA|\n"+
			"|JSC_TEST_ERROR_1|suppressionB|\n"+
			"|JSC_TEST_ERROR_2|suppressionA|\n")
}
