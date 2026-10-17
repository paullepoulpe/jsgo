package main

import (
	"slices"
	"strings"
	"testing"
)

func TestGoModRequires(t *testing.T) {
	gomod := "module m\n\ngo 1.24.0\n\nrequire golang.org/x/text v0.34.0 // indirect\n\nrequire (\n\texample.com/a v1.2.3\n\t// a comment\n\texample.com/b v0.0.0-2024 // indirect\n)\n"
	got := goModRequires(gomod)
	want := [][2]string{{"golang.org/x/text", "v0.34.0"}, {"example.com/a", "v1.2.3"}, {"example.com/b", "v0.0.0-2024"}}
	if !slices.Equal(got, want) {
		t.Errorf("goModRequires = %v, want %v", got, want)
	}
}

func TestCheckModules(t *testing.T) {
	const notice = "Go module dependencies\n  golang.org/x/text\n  v0.34.0 -- BSD 3-Clause License\n"
	for _, tc := range []struct {
		name, gomod, problem string // problem: substring of the only problem; "" for none
	}{
		{"no requirement", "module m\n\ngo 1.24.0\n", ""},
		{"listed and named", "module m\n\nrequire golang.org/x/text v0.34.0\n", ""},
		{"other version", "module m\n\nrequire golang.org/x/text v0.35.0\n", "NOTICE does not name golang.org/x/text v0.35.0"},
		{"unlisted module", "module m\n\nrequire (\n\texample.com/a v1.0.0\n)\n", "go.mod requires example.com/a, which tools/licenses/modules.go does not list"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			problems := checkModules(tc.gomod, notice)
			switch {
			case tc.problem == "" && len(problems) > 0:
				t.Errorf("unexpected problems: %q", problems)
			case tc.problem != "" && (len(problems) != 1 || !strings.Contains(problems[0], tc.problem)):
				t.Errorf("problems = %q, want one containing %q", problems, tc.problem)
			}
		})
	}
}
