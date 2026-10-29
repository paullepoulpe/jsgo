package jutil

import "testing"

func TestCaseFormatLowerHyphenToUpperCamelMatchesGuava(t *testing.T) {
	for _, r := range readTSV(t, "caseformat_hyphen.tsv") {
		in := StringFromUTF16(decodeUnits(t, r[0]))
		if got := encodeUnits(StringToUTF16(CaseFormat_LowerHyphenToUpperCamel(in))); got != r[1] {
			t.Errorf("LOWER_HYPHEN.to(UPPER_CAMEL, %s) = %s, Guava %s", r[0], got, r[1])
		}
	}

	cases := [][2]string{
		{"", ""},
		{"-", ""},
		{"--", ""},
		{"foo", "Foo"},
		{"Foo", "Foo"},
		{"FOO", "Foo"},
		{"foo-bar-baz", "FooBarBaz"},
		{"foo-BAR-baz", "FooBarBaz"},
		{"--a--B--", "AB"},
		{"é-ñ", "éñ"},
	}
	for _, tc := range cases {
		if got := CaseFormat_LowerHyphenToUpperCamel(tc[0]); got != tc[1] {
			t.Errorf("CaseFormat_LowerHyphenToUpperCamel(%q) = %q, want %q", tc[0], got, tc[1])
		}
	}
}

func TestCaseFormatLowerCamelToUpperUnderscore(t *testing.T) {
	cases := [][2]string{
		{"", ""},
		{"a", "A"},
		{"A", "A"},
		{"123", "123"},
		{"foo", "FOO"},
		{"Foo", "FOO"},
		{"FOO", "F_O_O"},
		{"fooBar", "FOO_BAR"},
		{"FooBar", "FOO_BAR"},
		{"fooBarBaz", "FOO_BAR_BAZ"},
		{"htmlParser", "HTML_PARSER"},
		{"xmlHTTPRequest", "XML_H_T_T_P_REQUEST"},
		{"éBar", "é_BAR"},
		{StringFromUTF16([]uint16{0xd800, 'a', 'B'}), StringFromUTF16([]uint16{0xd800, 'A', '_', 'B'})},
	}
	for _, tc := range cases {
		if got := CaseFormat_LowerCamelToUpperUnderscore(tc[0]); got != tc[1] {
			t.Errorf("CaseFormat_LowerCamelToUpperUnderscore(%q) = %q, want %q", tc[0], got, tc[1])
		}
	}
}
