package jutil

import "testing"

func TestCaseFormatLowerHyphenToUpperCamelMatchesGuava(t *testing.T) {
	for _, r := range readTSV(t, "caseformat_hyphen.tsv") {
		in := StringFromUTF16(decodeUnits(t, r[0]))
		if got := encodeUnits(StringToUTF16(CaseFormat_LowerHyphenToUpperCamel(in))); got != r[1] {
			t.Errorf("LOWER_HYPHEN.to(UPPER_CAMEL, %s) = %s, Guava %s", r[0], got, r[1])
		}
	}
}
