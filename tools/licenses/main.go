// Command licenses keeps the licence notices of the ported sources in line
// with PORTING.md ("No GPL code, anywhere").
//
// Closure Compiler ships the Rhino-derived files under
// com/google/javascript/rhino (Node, Token, IR, JSDocInfo, the jstype
// package, DToA, and their tests) under a dual licence: MPL 1.1, or
// alternatively GPL 2.0 or later. A recipient may use such a file under
// either one. This repository uses them under the MPL 1.1 only.
//
// Without flags, every tracked Go file whose header names a dual-licensed
// Java file ("// Port of <path>.java ...") gets that file's MPL 1.1 notice,
// with the GPL alternative removed and a line saying the translation is
// distributed under the MPL 1.1 only. The notice replaces any licence
// block already in place, or is inserted after the "// Port of" line.
//
// Java sources are read from .cache/closure-src (the pinned sources jar).
// Test sources are fetched once from the pinned tag on GitHub into
// .cache/closure-test-src.
//
// With -check, it only verifies, without network or cache (so CI can run
// it), that:
//   - no tracked file mentions the GPL, other than the documents that
//     state this rule;
//   - every tracked Go file that names a Guava, Gson, protobuf-java or
//     args4j source in a "// Port of <path>" or "// Derived from <path>"
//     line carries that project's copyright and licence notice before its
//     package clause (thirdparty.go), and every such line before a
//     package clause names a Closure source or one of those projects;
//   - NOTICE and the licence texts under LICENSES/ exist.
//
// Usage:
//
//	go run ./tools/licenses          # rewrite the notices, then check
//	go run ./tools/licenses -check   # check only
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const tag = "v20260415"

// gplPattern is assembled from pieces so that this file does not match it.
var gplPattern = regexp.MustCompile(`(?i)(general` + ` public licen[cs]e|\b(?:l)?g` + `pl\b)`)

// exempt are the tracked files allowed to mention the GPL: they state the rule.
var exempt = map[string]bool{
	"PORTING.md":             true,
	"PLAN.md":                true,
	"tools/licenses/main.go": true,
}

var portOf = regexp.MustCompile(`^// Port of (com/google/[^ ]+\.java)`)

func main() {
	check := flag.Bool("check", false, "only check: no GPL mention, third-party notices, NOTICE and LICENSES/")
	flag.Parse()
	files, err := trackedFiles()
	if err != nil {
		fatal(err)
	}
	if !*check {
		n := 0
		for _, f := range files {
			if !strings.HasSuffix(f, ".go") {
				continue
			}
			changed, err := fixFile(f)
			if err != nil {
				fatal(fmt.Errorf("%s: %v", f, err))
			}
			if changed {
				n++
			}
		}
		fmt.Printf("licenses: %d files updated\n", n)
	}
	var bad, problems []string
	ported := map[string]int{}
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			continue // deleted in the working tree
		}
		if !exempt[f] && gplPattern.Match(b) {
			bad = append(bad, f)
		}
		if strings.HasSuffix(f, ".go") {
			p, us := checkNotices(f, b)
			problems = append(problems, p...)
			for _, u := range us {
				ported[u.project]++
			}
		}
	}
	problems = append(problems, checkRequiredFiles()...)
	if len(bad) > 0 {
		fmt.Fprintf(os.Stderr, "licenses: %d tracked files mention the GPL (PORTING.md: no GPL code, anywhere):\n", len(bad))
		for _, f := range bad {
			fmt.Fprintln(os.Stderr, "  "+f)
		}
	}
	if len(problems) > 0 {
		fmt.Fprintf(os.Stderr, "licenses: %d licence-notice problems (PORTING.md: licence notices):\n", len(problems))
		for _, p := range problems {
			fmt.Fprintln(os.Stderr, "  "+p)
		}
	}
	if len(bad) > 0 || len(problems) > 0 {
		os.Exit(1)
	}
	var counts []string
	for _, u := range thirdParty {
		if n := ported[u.project]; n > 0 {
			counts = append(counts, fmt.Sprintf("%s %d", u.project, n))
		}
	}
	fmt.Printf("licenses: no tracked file mentions the GPL; third-party notices present (%s); NOTICE and LICENSES/ present\n", strings.Join(counts, ", "))
}

func trackedFiles() ([]string, error) {
	out, err := exec.Command("git", "ls-files", "-z").Output()
	if err != nil {
		return nil, err
	}
	var files []string
	for _, f := range strings.Split(string(out), "\x00") {
		if f != "" {
			files = append(files, f)
		}
	}
	return files, nil
}

// fixFile gives a Go port of a dual-licensed Java file the MPL-only notice.
func fixFile(path string) (bool, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	lines := strings.SplitAfter(string(b), "\n")
	portLine := -1
	var javaPath string
	for i := 0; i < len(lines) && i < 5; i++ {
		if m := portOf.FindStringSubmatch(lines[i]); m != nil {
			portLine, javaPath = i, m[1]
			break
		}
	}
	// Only the Rhino-derived files are dual-licensed; everything else
	// Closure ships is Apache 2.0.
	if portLine < 0 || !strings.HasPrefix(javaPath, "com/google/javascript/rhino/") {
		return false, nil
	}
	src, err := javaSource(javaPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "licenses: warning: %s: %v\n", path, err)
		return false, nil
	}
	notice, dual, err := mplOnlyNotice(src)
	if err != nil {
		return false, fmt.Errorf("%s: %v", javaPath, err)
	}
	if !dual {
		return false, nil
	}
	// The licence block already in the Go file, if any: a /* ... */ comment
	// right after the "// Port of" line (blank lines allowed) that mentions a
	// licence.
	head := strings.Join(lines[:portLine+1], "")
	rest := strings.Join(lines[portLine+1:], "")
	trimmed := strings.TrimLeft(rest, "\n")
	if strings.HasPrefix(trimmed, "/*") {
		if end := strings.Index(trimmed, "*/"); end >= 0 {
			block := trimmed[:end+2]
			if strings.Contains(block, "LICENSE BLOCK") || strings.Contains(block, "Mozilla Public License") {
				trimmed = strings.TrimLeft(trimmed[end+2:], "\n")
			}
		}
	}
	out := head + "\n" + notice + "\n\n" + trimmed
	// A file that ports several Java files carries a notice for each; give
	// the later ones the same treatment.
	out, err = mplOnlyBlocks(out)
	if err != nil {
		return false, err
	}
	if out == string(b) {
		return false, nil
	}
	return true, os.WriteFile(path, []byte(out), 0o644)
}

// mplOnlyBlocks applies mplOnlyNotice to every dual-licence comment block of s.
func mplOnlyBlocks(s string) (string, error) {
	var sb strings.Builder
	for {
		start := strings.Index(s, "/*")
		if start < 0 {
			break
		}
		end := strings.Index(s[start:], "*/")
		if end < 0 {
			break
		}
		block := s[start : start+end+2]
		sb.WriteString(s[:start])
		notice, dual, err := mplOnlyNotice([]byte(block))
		if err != nil {
			return "", err
		}
		if dual {
			sb.WriteString(notice)
		} else {
			sb.WriteString(block)
		}
		s = s[start+end+2:]
	}
	sb.WriteString(s)
	return sb.String(), nil
}

// mplOnlyNotice returns the leading licence comment of a Java source with the
// GPL alternative removed, and whether the source was dual-licensed at all.
func mplOnlyNotice(src []byte) (string, bool, error) {
	s := string(src)
	start := strings.Index(s, "/*")
	if start < 0 {
		return "", false, nil
	}
	end := strings.Index(s[start:], "*/")
	if end < 0 {
		return "", false, nil
	}
	block := s[start : start+end+2]
	if !strings.Contains(block, "Mozilla Public License") || !gplPattern.MatchString(block) {
		return "", false, nil
	}
	var out []string
	skipping := false
	for _, l := range strings.Split(block, "\n") {
		switch {
		case strings.Contains(l, "Alternatively, the contents of this file may be used"):
			skipping = true
			continue
		case skipping:
			if strings.Contains(l, "under either the MPL or the GPL") || strings.Contains(l, "the terms of any one of the MPL") {
				skipping = false
			}
			continue
		}
		if strings.Contains(l, "***** END LICENSE BLOCK") {
			out = append(out,
				" * This Go translation is distributed under the MPL 1.1 only, as the",
				" * original's dual licence allows; the other licence it offers is not",
				" * used.",
				" *")
		}
		l = strings.Replace(l, "Version: MPL 1.1/GPL 2.0", "Version: MPL 1.1", 1)
		out = append(out, l)
	}
	// Collapse the blank comment lines left where the paragraph was.
	notice := strings.Join(out, "\n")
	for strings.Contains(notice, " *\n *\n *\n") {
		notice = strings.ReplaceAll(notice, " *\n *\n *\n", " *\n *\n")
	}
	notice = strings.ReplaceAll(notice, " *\n *\n * This Go translation", " *\n * This Go translation")
	if gplPattern.MatchString(notice) {
		return "", true, fmt.Errorf("could not remove the GPL alternative from the notice:\n%s", notice)
	}
	return notice, true, nil
}

// javaSource reads a Java file of the pinned Closure release: main sources
// from .cache/closure-src, test sources from .cache/closure-test-src
// (fetched from GitHub on first use).
func javaSource(p string) ([]byte, error) {
	if b, err := os.ReadFile(filepath.Join(".cache", "closure-src", p)); err == nil {
		return b, nil
	}
	cached := filepath.Join(".cache", "closure-test-src", p)
	if b, err := os.ReadFile(cached); err == nil {
		return b, nil
	}
	// Test sources live under test/; test-only helpers such as
	// rhino/testing live under src/ but are left out of the sources jar.
	var errs []string
	for _, dir := range []string{"test", "src"} {
		url := "https://raw.githubusercontent.com/google/closure-compiler/" + tag + "/" + dir + "/" + p
		b, err := fetch(url)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		if err := os.MkdirAll(filepath.Dir(cached), 0o755); err != nil {
			return nil, err
		}
		return b, os.WriteFile(cached, b, 0o644)
	}
	return nil, fmt.Errorf("%s not in .cache/closure-src; %s", p, strings.Join(errs, "; "))
}

func fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: %s", url, resp.Status)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if !bytes.Contains(b, []byte("package ")) {
		return nil, fmt.Errorf("GET %s: not a Java source", url)
	}
	return b, nil
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "licenses:", err)
	os.Exit(1)
}
