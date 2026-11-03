package jutil

import "testing"

func TestIsValidPattern(t *testing.T) {
	cases := []struct {
		pattern string
		want    bool
	}{
		{"", true},
		{`.*\.closure\.js`, true},
		{`^foo/(bar|baz)[0-9]+\.js$`, true},
		{`\Q([unclosed\E+`, true},
		{`\Qunclosed_quote([`, true},
		{`[\Q]-\E]`, true},
		{`\Eabc`, true},
		{`[]a]`, true},
		{`[^]a]`, true},
		{`[\d\w\s\v\h]`, true},
		{`(foo)\1`, true},
		{`\bword\b\z`, true},
		{`a{2}{3}`, true},
		{`a{2000}`, true},
		{`\`, false},
		{`\q`, false},
		{`\M`, false},
		{`[abc`, false},
		{`[]`, false},
		{`[^]`, false},
		{`(abc`, false},
		{`abc)`, false},
		{`*abc`, false},
		{`[z-a]`, false},
	}
	for _, tc := range cases {
		if got := IsValidPattern(tc.pattern); got != tc.want {
			t.Errorf("IsValidPattern(%q) = %v, want %v", tc.pattern, got, tc.want)
		}
	}
}
