package main

import (
	"fmt"
	"os"
	"strings"
)

// module is a Go module that go.mod may require. Its code is linked into
// the programs that import it, so its licence notice must ship with them:
// the licence text under LICENSES/ and an entry in NOTICE.
type module struct {
	path string
	text string // the licence text under LICENSES/
}

var modules = []module{
	{"golang.org/x/text", "LICENSES/BSD-3-Clause-golang-x-text.txt"},
}

// goModRequires returns the module paths and versions that a go.mod file
// requires, in file order: single-line require directives and require
// blocks, with comments removed.
func goModRequires(gomod string) [][2]string {
	var reqs [][2]string
	inBlock := false
	for _, line := range strings.Split(gomod, "\n") {
		if i := strings.Index(line, "//"); i >= 0 {
			line = line[:i]
		}
		f := strings.Fields(line)
		switch {
		case len(f) == 0:
		case inBlock && f[0] == ")":
			inBlock = false
		case inBlock && len(f) >= 2:
			reqs = append(reqs, [2]string{f[0], f[1]})
		case f[0] == "require" && len(f) == 2 && f[1] == "(":
			inBlock = true
		case f[0] == "require" && len(f) >= 3:
			reqs = append(reqs, [2]string{f[1], f[2]})
		}
	}
	return reqs
}

// checkModules returns the problems with the notices of the modules that
// gomod requires: a module missing from modules, or one that notice does
// not name at the required version ("<path> <version>").
func checkModules(gomod, notice string) []string {
	var problems []string
	notice = strings.Join(strings.Fields(notice), " ")
	for _, r := range goModRequires(gomod) {
		var m *module
		for i := range modules {
			if modules[i].path == r[0] {
				m = &modules[i]
			}
		}
		switch {
		case m == nil:
			problems = append(problems, fmt.Sprintf("go.mod requires %s, which tools/licenses/modules.go does not list: add its licence text under LICENSES/, an entry in NOTICE, and the module to the list", r[0]))
		case !strings.Contains(notice, r[0]+" "+r[1]):
			problems = append(problems, fmt.Sprintf("NOTICE does not name %s %s, the version go.mod requires (its entry must give the module path and version)", r[0], r[1]))
		}
	}
	return problems
}

// checkModuleNotices runs checkModules on the repository's go.mod and NOTICE.
func checkModuleNotices() []string {
	gomod, err := os.ReadFile("go.mod")
	if err != nil {
		return []string{"go.mod: " + err.Error()}
	}
	notice, err := os.ReadFile("NOTICE")
	if err != nil {
		return nil // reported by checkRequiredFiles
	}
	return checkModules(string(gomod), string(notice))
}
