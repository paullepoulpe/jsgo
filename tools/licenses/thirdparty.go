package main

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
)

// upstream is a third-party project whose code some Go files translate.
// Such a file names its source in a "// Port of <path> ..." or
// "// Derived from <path> ..." line and, before its package clause,
// carries the source's copyright and licence notice.
type upstream struct {
	project  string
	prefixes []string // Java source paths of the project
	// copyright and licence must both match the file's header, with
	// comment markers removed and white space collapsed.
	copyright *regexp.Regexp
	licence   *regexp.Regexp
	text      string // the licence text under LICENSES/
	// pointer: the header must also name text. protobuf's own header refers
	// to "the LICENSE file", which is text in this repository.
	pointer bool
}

var apache2 = regexp.MustCompile(`Licensed under the Apache License, Version 2\.0`)

var thirdParty = []*upstream{
	{
		project:   "Guava",
		prefixes:  []string{"com/google/common/"},
		copyright: regexp.MustCompile(`Copyright \(C\) \d{4}(-\d{4})? The Guava Authors`),
		licence:   apache2,
		text:      "LICENSES/Apache-2.0.txt",
	},
	{
		project:   "Gson",
		prefixes:  []string{"com/google/gson/"},
		copyright: regexp.MustCompile(`Copyright \(C\) \d{4}(-\d{4})? (Google Inc\.|Google LLC|The Gson [Aa]uthors|The Android Open Source Project)`),
		licence:   apache2,
		text:      "LICENSES/Apache-2.0.txt",
	},
	{
		project:   "protobuf-java",
		prefixes:  []string{"com/google/protobuf/"},
		copyright: regexp.MustCompile(`Copyright \d{4} Google Inc\.`),
		licence:   regexp.MustCompile(`governed by a BSD-style license|Redistribution and use in source and binary forms`),
		text:      "LICENSES/BSD-3-Clause-protobuf.txt",
		pointer:   true,
	},
	{
		project:   "args4j",
		prefixes:  []string{"org/kohsuke/args4j/"},
		copyright: regexp.MustCompile(`Copyright \(c\) \d{4} Kohsuke Kawaguchi`),
		licence:   regexp.MustCompile(`Permission is hereby granted, free of charge`),
		text:      "LICENSES/MIT-args4j.txt",
	},
}

// closurePrefixes are the Java source paths of Closure Compiler itself
// (Apache 2.0, or MPL 1.1 for the Rhino-derived files fixFile handles).
var closurePrefixes = []string{"com/google/javascript/", "com/google/debugging/"}

// requiredFiles must exist (and be non-empty) for -check to pass: NOTICE
// and the licence texts it points to.
func requiredFiles() []string {
	files := []string{"NOTICE", "LICENSES/Apache-2.0.txt", "LICENSES/MPL-1.1.txt"}
	for _, u := range thirdParty {
		if !slices.Contains(files, u.text) {
			files = append(files, u.text)
		}
	}
	for _, m := range modules {
		if !slices.Contains(files, m.text) {
			files = append(files, m.text)
		}
	}
	return files
}

// sourceMarker matches a line naming the source a file translates.
var sourceMarker = regexp.MustCompile(`^\s*// (?:Port of|Derived from) (\S+)`)

// classify returns the third-party project of a Java source path, or
// closure=true for a Closure source. Both are zero for an unknown path.
func classify(path string) (u *upstream, closure bool) {
	p := strings.TrimPrefix(strings.TrimPrefix(path, "test/"), "src/")
	for _, pre := range closurePrefixes {
		if strings.HasPrefix(p, pre) {
			return nil, true
		}
	}
	for _, u := range thirdParty {
		for _, pre := range u.prefixes {
			if strings.HasPrefix(p, pre) {
				return u, false
			}
		}
	}
	return nil, false
}

// checkNotices returns the problems with the third-party notices of the Go
// source b (path f):
//   - a "// Port of" or "// Derived from" line anywhere in the file that
//     names a Guava, Gson, protobuf or args4j source, without that
//     project's copyright and licence notice before the package clause;
//   - a "// Port of" or "// Derived from" line before the package clause
//     that names a source of no known project (add the project above).
//
// It also returns the projects whose notice the file carries.
func checkNotices(f string, b []byte) (problems []string, ported []*upstream) {
	lines := strings.Split(string(b), "\n")
	pkg := len(lines)
	for i, l := range lines {
		if strings.HasPrefix(l, "package ") {
			pkg = i
			break
		}
	}
	var needed []*upstream
	for i, l := range lines {
		m := sourceMarker.FindStringSubmatch(l)
		if m == nil {
			continue
		}
		path := strings.TrimRight(m[1], ",:;.)")
		u, closure := classify(path)
		switch {
		case u != nil:
			if !slices.Contains(needed, u) {
				needed = append(needed, u)
			}
		case !closure && i < pkg && strings.Contains(path, "/") && strings.HasSuffix(path, ".java"):
			problems = append(problems, fmt.Sprintf("%s:%d: %s is not a Closure source nor one of a known third-party project (tools/licenses/thirdparty.go)", f, i+1, path))
		}
	}
	if len(needed) == 0 {
		return problems, nil
	}
	header := normalizeComment(lines[:pkg])
	for _, u := range needed {
		var missing []string
		if !u.copyright.MatchString(header) {
			missing = append(missing, u.project+"'s copyright line")
		}
		if !u.licence.MatchString(header) {
			missing = append(missing, u.project+"'s licence")
		}
		if u.pointer && !strings.Contains(header, u.text) {
			missing = append(missing, "a pointer to "+u.text)
		}
		if len(missing) > 0 {
			problems = append(problems, fmt.Sprintf("%s: translates %s code but its header (before the package clause) lacks %s", f, u.project, strings.Join(missing, " and ")))
			continue
		}
		ported = append(ported, u)
	}
	return problems, ported
}

// normalizeComment joins lines of Go and Java-style comments into one line
// of text: comment markers are dropped and runs of white space collapse.
func normalizeComment(lines []string) string {
	var words []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		for _, p := range []string{"//", "/*", "*/", "*"} {
			l = strings.TrimPrefix(l, p)
		}
		l = strings.TrimSuffix(l, "*/")
		words = append(words, strings.Fields(l)...)
	}
	return strings.Join(words, " ")
}

// checkRequiredFiles returns the problems with NOTICE and LICENSES/.
func checkRequiredFiles() []string {
	var problems []string
	for _, f := range requiredFiles() {
		fi, err := os.Stat(f)
		switch {
		case err != nil:
			problems = append(problems, fmt.Sprintf("%s is missing (PORTING.md, licence notices)", f))
		case fi.Size() == 0:
			problems = append(problems, fmt.Sprintf("%s is empty", f))
		}
	}
	return problems
}
