// Command licenses keeps the licence notices of the ported sources in line
// with PORTING.md ("No GPL code, anywhere").
//
// Closure Compiler ships the Rhino-derived files under
// com/google/javascript/rhino (Node, Token, IR, JSDocInfo, the jstype
// package, DToA, and their tests) under a dual licence: MPL 1.1, or
// alternatively GPL 2.0 or later. A recipient may use such a file under
// either one. This repository uses them under the MPL 1.1 only.
//
// Without flags, every tracked Go file whose header names Closure Java
// files ("// Port of <path>.java ...", main or test sources) gets the
// licence comment of each right after its "// Port of" lines (closure.go):
// the file's own Apache 2.0 notice, or for a dual-licensed one its MPL 1.1
// notice, with the GPL alternative removed and a line saying the
// translation is distributed under the MPL 1.1 only. Licence comments
// already there that are not those are replaced. A Java file with no
// licence comment must be listed in noNotice (closure.go), with the reason.
//
// Java sources are read from srcDir, .cache/closure-src (the pinned sources
// jar, which tools/oracle/fetch.sh unpacks there); javaSource refuses them
// if they are another release than tag. Test sources are fetched once from
// the pinned tag on GitHub into testSrcDir, .cache/closure-test-src/<tag>.
//
// With -check, it only verifies, without network or cache (so CI can run
// it), that:
//   - no tracked file mentions the GPL, other than the documents that
//     state this rule. A gzip file is checked decompressed, a zip or jar
//     archive member by member, a file in another compressed format is an
//     error, and other binary files are skipped (content.go);
//   - every tracked Go file that names a Guava, Gson, protobuf-java or
//     args4j source in a "// Port of <path>" or "// Derived from <path>"
//     line carries that project's copyright and licence notice before its
//     package clause (thirdparty.go), and every such line before a
//     package clause names a Closure source or one of those projects;
//   - every tracked Go file whose header names a Closure source that
//     noNotice does not list carries a Closure licence comment before its
//     package clause, and every source noNotice lists is ported
//     (closure.go). Without the Java sources it cannot tell one Closure
//     notice from another; the fixer puts in each source's own;
//   - no tracked file (decompressed, as above) mentions a corpus file of
//     tools/oracle/gpl-excluded-corpus.txt (tools/oracle/corpus): it
//     contains the file's path inside its npm package;
//   - every npm package copied into test data, a tracked
//     ".../testdata/.../src/<name>@<version>/" directory, holds the
//     package's licence file (LICENSE, LICENSE.txt, licenses.md, ...)
//     directly (packages.go);
//   - every module that go.mod requires is listed in modules.go, and NOTICE
//     names it at the required version (modules.go);
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
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/paullepoulpe/jsgo/tools/oracle/corpus"
)

const tag = "v20261021"

// srcDir holds the main Java sources of the pinned release (the sources
// jar), and testSrcDir caches its test sources, fetched from GitHub. Both
// are relative to the repository root.
const (
	srcDir     = ".cache/closure-src"
	testSrcDir = ".cache/closure-test-src/" + tag
)

// gplPattern is assembled from pieces so that this file does not match it.
var gplPattern = regexp.MustCompile(`(?i)(general` + ` public licen[cs]e|\b(?:l)?g` + `pl\b)`)

// exempt are the tracked files allowed to mention the GPL: they state the rule.
var exempt = map[string]bool{
	"PORTING.md":                true,
	"PLAN.md":                   true,
	"tools/licenses/main.go":    true,
	"tools/licenses/content.go": true,
	corpus.ListPath:             true,
}

func main() {
	check := flag.Bool("check", false, "only check: no GPL mention, Closure and third-party notices, NOTICE and LICENSES/")
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
	excluded, err := corpus.Excluded()
	if err != nil {
		fatal(err)
	}
	index := newExcludedIndex(excluded)
	var bad, mentions, problems []string
	ported := map[string]int{}
	named := map[string]bool{}
	closurePorts := 0
	skipped := 0
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			continue // deleted in the working tree
		}
		text, skip, err := fileText(f, b)
		switch {
		case err != nil:
			problems = append(problems, err.Error())
		case skip:
			skipped++
		default:
			if !exempt[f] && mentionsGPL(text) {
				bad = append(bad, f)
			}
			if f != corpus.ListPath {
				for _, e := range index.mentioned(text) {
					mentions = append(mentions, f+": "+e)
				}
			}
		}
		if strings.HasSuffix(f, ".go") {
			p, us := checkNotices(f, b)
			problems = append(problems, p...)
			for _, u := range us {
				ported[u.project]++
			}
			p, port := checkClosureNotices(f, b, named)
			problems = append(problems, p...)
			if port {
				closurePorts++
			}
		}
	}
	problems = append(problems, checkNoNotice(named)...)
	problems = append(problems, checkPackageLicences(files)...)
	problems = append(problems, checkModuleNotices()...)
	problems = append(problems, checkRequiredFiles()...)
	if len(bad) > 0 {
		fmt.Fprintf(os.Stderr, "licenses: %d tracked files mention the GPL (PORTING.md: no GPL code, anywhere):\n", len(bad))
		for _, f := range bad {
			fmt.Fprintln(os.Stderr, "  "+f)
		}
	}
	if len(mentions) > 0 {
		fmt.Fprintf(os.Stderr, "licenses: %d mentions of corpus files that %s excludes (tools/oracle/README.md):\n", len(mentions), corpus.ListPath)
		for _, m := range mentions {
			fmt.Fprintln(os.Stderr, "  "+m)
		}
	}
	if len(problems) > 0 {
		fmt.Fprintf(os.Stderr, "licenses: %d licence-notice problems (PORTING.md: licence notices):\n", len(problems))
		for _, p := range problems {
			fmt.Fprintln(os.Stderr, "  "+p)
		}
	}
	if len(bad) > 0 || len(mentions) > 0 || len(problems) > 0 {
		os.Exit(1)
	}
	var counts []string
	for _, u := range thirdParty {
		if n := ported[u.project]; n > 0 {
			counts = append(counts, fmt.Sprintf("%s %d", u.project, n))
		}
	}
	fmt.Printf("licenses: no tracked file mentions the GPL or one of the %d excluded corpus files (gzip and zip files decompressed, %d other binary files skipped); %d Closure ports carry their sources' licence comments (%d sources have none); third-party notices present (%s); test-data packages licensed; %d module dependencies in NOTICE; NOTICE and LICENSES/ present\n", len(excluded), skipped, closurePorts, len(noNotice), strings.Join(counts, ", "), len(modules))
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

// fixFile gives a Go port of Closure Java sources the licence comment of
// each (portNotices, withNotices), and turns any dual-licence comment left
// in the file into its MPL-only form. A source it cannot read is a warning.
func fixFile(path string) (bool, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	lines := strings.Split(string(b), "\n")
	paths := closureSources(lines[:packageClause(lines)])
	if len(paths) == 0 {
		return false, nil
	}
	var readErr error
	notices, err := portNotices(paths, func(p string) ([]byte, error) {
		src, err := javaSource(p)
		readErr = err
		return src, err
	})
	switch {
	case readErr != nil:
		fmt.Fprintf(os.Stderr, "licenses: warning: %s: %v\n", path, readErr)
		return false, nil
	case err != nil:
		return false, err
	}
	out := string(b)
	if len(notices) > 0 {
		var replaced []string
		out, replaced = withNotices(out, notices)
		for _, r := range replaced {
			fmt.Fprintf(os.Stderr, "licenses: %s: replaced a licence comment that is not its source's:\n%s\n", path, r)
		}
	}
	// A file that ports several Java files carries a notice for each; give
	// any dual-licence one the MPL-only treatment.
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
// from srcDir, test sources from testSrcDir (fetched from GitHub on first
// use). A "test/" or "src/" prefix of p says which directory of the
// repository holds it.
func javaSource(p string) ([]byte, error) {
	dirs := []string{"test", "src"}
	switch {
	case strings.HasPrefix(p, "test/"):
		dirs = dirs[:1]
	case strings.HasPrefix(p, "src/"):
		dirs = dirs[1:]
	}
	p = canonical(p)
	if err := checkSrcVersion(); err != nil {
		return nil, err
	}
	if b, err := os.ReadFile(filepath.Join(filepath.FromSlash(srcDir), filepath.FromSlash(p))); err == nil {
		return b, nil
	}
	// The sources jar holds CommandLineRunnerVersion.java, which the build
	// generates, at its root.
	if b, err := os.ReadFile(filepath.Join(filepath.FromSlash(srcDir), path.Base(p))); err == nil {
		if pkg := strings.ReplaceAll(path.Dir(p), "/", "."); bytes.Contains(b, []byte("package "+pkg+";")) {
			return b, nil
		}
	}
	cached := filepath.Join(filepath.FromSlash(testSrcDir), filepath.FromSlash(p))
	if b, err := os.ReadFile(cached); err == nil {
		return b, nil
	}
	// Test sources live under test/; test-only helpers such as
	// rhino/testing live under src/ but are left out of the sources jar.
	var errs []string
	for _, dir := range dirs {
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
	return nil, fmt.Errorf("%s not in %s; %s", p, srcDir, strings.Join(errs, "; "))
}

// checkSrcVersion fails if srcDir holds the sources of another release than
// tag: the sources jar's generated CommandLineRunnerVersion.java, at its
// root, names the release. A missing srcDir is not an error here (every
// main source then falls back to GitHub).
func checkSrcVersion() error {
	b, err := os.ReadFile(filepath.Join(filepath.FromSlash(srcDir), "CommandLineRunnerVersion.java"))
	if err != nil {
		return nil
	}
	if !bytes.Contains(b, []byte(`"`+tag+`"`)) {
		return fmt.Errorf("%s holds the sources of another release than %s (tools/oracle/pin.env): run tools/oracle/fetch.sh", srcDir, tag)
	}
	return nil
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
