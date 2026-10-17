package main

import (
	"fmt"
	"regexp"
	"sort"
)

// packageDir matches the directory of an npm package copied into test data:
// ".../testdata/.../src/<name>@<version>/", where <name> may be scoped
// ("@angular/core"). The match is the first such directory of a path, so
// a package's own src/ directories do not count.
var packageDir = regexp.MustCompile(`^(?:.*/)?testdata/(?:[^/]+/)*?src/(?:@[^/@]+/)?[^/@]+@[^/@]+/`)

// licenceFile matches the base name of a package's licence file, as npm
// packages name it: LICENSE, LICENCE, LICENSE.txt, LICENSE-MIT, COPYING,
// licenses.md and the like.
var licenceFile = regexp.MustCompile(`(?i)^(?:licen[cs]es?|copying)(?:[.-][^/]*)?$`)

// checkPackageLicences returns a problem for each npm package copied into
// test data (packageDir) whose directory holds no tracked licence file
// (licenceFile) directly. files are the tracked files.
func checkPackageLicences(files []string) []string {
	licensed := map[string]bool{}
	for _, f := range files {
		dir := packageDir.FindString(f)
		if dir == "" {
			continue
		}
		if licenceFile.MatchString(f[len(dir):]) {
			licensed[dir] = true
		} else if !licensed[dir] {
			licensed[dir] = false
		}
	}
	var problems []string
	for dir, ok := range licensed {
		if !ok {
			problems = append(problems, fmt.Sprintf("%s holds an npm package but no licence file: add the package's LICENSE from the registry tarball of that version, and list it in NOTICE", dir))
		}
	}
	sort.Strings(problems)
	return problems
}
