package main

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

const portFoo = "// Port of com/google/javascript/jscomp/Foo.java @ closure-compiler v20261021\n"

// closureNoticeOf returns a Closure Apache 2.0 licence comment of the year.
func closureNoticeOf(year string) string {
	return `/*
 * Copyright ` + year + ` The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 */`
}

var apacheNotice = closureNoticeOf("2018")

// gayNotice stands for the second licence comment of DToA.java.
const gayNotice = `/****************************************************************
 *
 * The author of this software is David M. Gay.
 *
 * Copyright (c) 1991, 2000, 2001 by Lucent Technologies.
 *
 ***************************************************************/`

// dualNotice is a Rhino licence comment, assembled so that this file does
// not mention the other licence.
var dualNotice = strings.NewReplacer("{G}", gplWord, "{GPLong}", "General"+" Public License").Replace(`/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1/{G} 2.0
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License.
 *
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Alternatively, the contents of this file may be used under the terms of
 * the GNU {GPLong} Version 2 or later (the "{G}"), in which
 * case the provisions of the {G} are applicable instead of those above. If
 * not delete the provisions above, a recipient may use your version of this
 * file under either the MPL or the {G}.
 *
 * ***** END LICENSE BLOCK ***** */`)

func mplOnlyOf(t *testing.T, notice string) string {
	t.Helper()
	mpl, dual, err := mplOnlyNotice([]byte(notice))
	if err != nil || !dual {
		t.Fatalf("mplOnlyNotice: dual=%v, err=%v", dual, err)
	}
	return mpl
}

func TestClosureSources(t *testing.T) {
	header := []string{
		portFoo,
		"// Port of test/com/google/javascript/jscomp/FooTest.java @ closure-compiler v20261021",
		"// Port of com/google/javascript/jscomp/FooTest.java @ closure-compiler v20261021", // same source
		"// Port of com/google/javascript/jscomp/Foo.java @ closure-compiler v20261021",     // repeated
		"// (test/com/google/javascript/jscomp/BarTest.java)",                               // not a Port of line
		"// Port of com/google/common/base/Ascii.java (Guava 33.4.6-jre)",                   // third party
		"// Port of com/google/javascript/jscomp/js/es6_runtime.js",                         // not Java
		"// Port of src/com/google/debugging/sourcemap/Util.java, with a note",
	}
	got := closureSources(header)
	want := []string{
		"com/google/javascript/jscomp/Foo.java",
		"test/com/google/javascript/jscomp/FooTest.java",
		"src/com/google/debugging/sourcemap/Util.java",
	}
	if !slices.Equal(got, want) {
		t.Errorf("closureSources = %q, want %q", got, want)
	}
}

func TestJavaNotices(t *testing.T) {
	for _, tc := range []struct {
		name, src string
		want      []string
		err       string
	}{
		{"apache", apacheNotice + "\npackage com.google.javascript.jscomp;\n\n/* Copyright in the body. */\nclass Foo {}\n", []string{apacheNotice}, ""},
		{"generated, none", "package com.google.javascript.jscomp;\n\nimport javax.annotation.processing.Generated;\n", nil, ""},
		{"two comments", dualNotice + "\n\n" + gayNotice + "\npackage com.google.javascript.rhino.dtoa;\n", []string{dualNotice, gayNotice}, ""},
		{"javadoc only", "/** Helpers. */\npackage com.google.javascript.jscomp;\n", nil, ""},
		{"line comments", "// Copyright 2020 Someone.\npackage com.google.javascript.jscomp;\n", nil, "does not recognise"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := javaNotices([]byte(tc.src))
			switch {
			case tc.err != "" && (err == nil || !strings.Contains(err.Error(), tc.err)):
				t.Errorf("err = %v, want one containing %q", err, tc.err)
			case tc.err == "" && err != nil:
				t.Errorf("unexpected error: %v", err)
			case !slices.Equal(got, tc.want):
				t.Errorf("javaNotices = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestPortNotices(t *testing.T) {
	const pkg = "\npackage com.google.javascript.jscomp;\n"
	sources := map[string]string{
		"com/google/javascript/jscomp/Foo.java":                                apacheNotice + pkg,
		"com/google/javascript/jscomp/Bar.java":                                apacheNotice + pkg,
		"com/google/javascript/jscomp/Baz.java":                                closureNoticeOf("2009") + pkg,
		"com/google/javascript/jscomp/FooTest.java":                            closureNoticeOf("2010") + pkg,
		"com/google/javascript/rhino/dtoa/DToA.java":                           dualNotice + "\n\n" + gayNotice + pkg,
		"com/google/javascript/jscomp/AutoValue_PassFactory.java":              pkg,
		"com/google/javascript/jscomp/Gen.java":                                pkg,
		"com/google/javascript/jscomp/modules/AutoBuilder_Export_Builder.java": apacheNotice + pkg,
		"com/google/javascript/jscomp/Odd.java":                                guavaHeader + pkg,
	}
	source := func(p string) ([]byte, error) {
		if s, ok := sources[canonical(p)]; ok {
			return []byte(s), nil
		}
		return nil, errors.New("no such source")
	}
	const jscomp = "com/google/javascript/jscomp/"
	for _, tc := range []struct {
		name  string
		paths []string
		want  []string
		err   string
	}{
		{"one", []string{jscomp + "Foo.java"}, []string{apacheNotice}, ""},
		{"test source", []string{"test/" + jscomp + "FooTest.java"}, []string{closureNoticeOf("2010")}, ""},
		{"same notice once", []string{jscomp + "Foo.java", jscomp + "Bar.java"}, []string{apacheNotice}, ""},
		{"each one's", []string{jscomp + "Foo.java", jscomp + "Baz.java"}, []string{apacheNotice, closureNoticeOf("2009")}, ""},
		{"with a generated source", []string{jscomp + "Foo.java", jscomp + "AutoValue_PassFactory.java"}, []string{apacheNotice}, ""},
		{"generated only", []string{jscomp + "AutoValue_PassFactory.java"}, nil, ""},
		{"dual licence", []string{"com/google/javascript/rhino/dtoa/DToA.java"}, []string{mplOnlyOf(t, dualNotice), gayNotice}, ""},
		{"no notice, not listed", []string{jscomp + "Gen.java"}, nil, "add it to noNotice"},
		{"listed, has a notice", []string{"com/google/javascript/jscomp/modules/AutoBuilder_Export_Builder.java"}, nil, "remove it there"},
		{"not Closure's notice", []string{jscomp + "Odd.java"}, nil, "would not recognise"},
		{"unreadable", []string{jscomp + "Missing.java"}, nil, "no such source"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := portNotices(tc.paths, source)
			switch {
			case tc.err != "" && (err == nil || !strings.Contains(err.Error(), tc.err)):
				t.Errorf("err = %v, want one containing %q", err, tc.err)
			case tc.err == "" && err != nil:
				t.Errorf("unexpected error: %v", err)
			case !slices.Equal(got, tc.want):
				t.Errorf("portNotices = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestWithNotices(t *testing.T) {
	const body = "package jscomp\n\nfunc f() {}\n"
	const cont = "// (test/com/google/javascript/jscomp/FooTest.java)\n"
	old := closureNoticeOf("2007")
	for _, tc := range []struct {
		name, src string
		notices   []string
		want      string
		replaced  []string
	}{
		{"added", portFoo + "\n" + body, []string{apacheNotice}, portFoo + "\n" + apacheNotice + "\n\n" + body, nil},
		{"after a continuation line", portFoo + cont + "\n" + body, []string{apacheNotice}, portFoo + cont + "\n" + apacheNotice + "\n\n" + body, nil},
		{"package comment kept", portFoo + "\n// Package jscomp does it.\n" + body, []string{apacheNotice}, portFoo + "\n" + apacheNotice + "\n\n// Package jscomp does it.\n" + body, nil},
		{"present", portFoo + "\n" + apacheNotice + "\n\n" + body, []string{apacheNotice}, portFoo + "\n" + apacheNotice + "\n\n" + body, nil},
		{"present, other spacing", portFoo + "\n\n" + strings.ReplaceAll(apacheNotice, "License\");", "License\");  ") + "\n" + body, []string{apacheNotice}, portFoo + "\n\n" + strings.ReplaceAll(apacheNotice, "License\");", "License\");  ") + "\n" + body, nil},
		{"wrong one replaced", portFoo + "\n" + old + "\n\n" + body, []string{apacheNotice}, portFoo + "\n" + apacheNotice + "\n\n" + body, []string{old}},
		{"one of two missing", portFoo + "\n" + apacheNotice + "\n\n" + body, []string{apacheNotice, gayNotice}, portFoo + "\n" + apacheNotice + "\n\n" + gayNotice + "\n\n" + body, nil},
		{"later notice of another source kept", portFoo + "\n" + apacheNotice + "\n\n// (with rhino/testing/X.java)\n\n" + old + "\n\n" + body, []string{apacheNotice}, portFoo + "\n" + apacheNotice + "\n\n// (with rhino/testing/X.java)\n\n" + old + "\n\n" + body, nil},
		{"no notices", portFoo + "\n" + body, nil, portFoo + "\n" + body, nil},
		{"not a port", "// Go-only.\n\n" + body, []string{apacheNotice}, "// Go-only.\n\n" + body, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, replaced := withNotices(tc.src, tc.notices)
			if got != tc.want {
				t.Errorf("withNotices =\n%s\nwant\n%s", got, tc.want)
			}
			if !slices.Equal(replaced, tc.replaced) {
				t.Errorf("replaced = %q, want %q", replaced, tc.replaced)
			}
			if again, _ := withNotices(got, tc.notices); again != got {
				t.Errorf("not idempotent: second run gives\n%s", again)
			}
		})
	}
}

func TestCheckClosureNotices(t *testing.T) {
	const body = "\npackage x\n"
	const generatedPort = "// Port of com/google/javascript/jscomp/AutoValue_PassFactory.java @ closure-compiler v20261021\n"
	mplOnly := mplOnlyOf(t, dualNotice)
	mplNoOnlyLine := strings.Replace(mplOnly, "This Go translation", "This translation", 1)
	for _, tc := range []struct {
		name, src string
		problem   string // substring of the only problem; "" for none
		port      bool
	}{
		{"not a port", "// Go-only.\n" + body, "", false},
		{"third-party port", "// Port of com/google/common/base/Ascii.java (Guava 33.4.6-jre)\n\n" + guavaHeader + body, "", false},
		{"notice present", portFoo + "\n" + apacheNotice + "\n" + body, "", true},
		{"notice missing", portFoo + body, "ports com/google/javascript/jscomp/Foo.java but lacks its licence comment", true},
		{"test source, missing", "// Port of test/com/google/javascript/jscomp/FooTest.java @ closure-compiler v20261021\n" + body, "ports test/com/google/javascript/jscomp/FooTest.java but lacks", true},
		{"generated source only", generatedPort + body, "", true},
		{"generated source and another", portFoo + generatedPort + body, "ports com/google/javascript/jscomp/Foo.java but", true},
		{"MPL-only notice", "// Port of com/google/javascript/rhino/Node.java @ closure-compiler v20261021\n\n" + mplOnly + "\n" + body, "", true},
		{"MPL notice without the MPL-only line", "// Port of com/google/javascript/rhino/Node.java @ closure-compiler v20261021\n\n" + mplNoOnlyLine + "\n" + body, "lacks its licence comment", true},
		{"another project's notice", portFoo + "\n" + guavaHeader + body, "lacks its licence comment", true},
		{"notice after the package clause", portFoo + body + "\n" + apacheNotice + "\n", "lacks its licence comment", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			problems, port := checkClosureNotices("x.go", []byte(tc.src), map[string]bool{})
			switch {
			case port != tc.port:
				t.Errorf("port = %v, want %v", port, tc.port)
			case tc.problem == "" && len(problems) > 0:
				t.Errorf("unexpected problems: %q", problems)
			case tc.problem != "" && (len(problems) != 1 || !strings.Contains(problems[0], tc.problem)):
				t.Errorf("problems = %q, want one containing %q", problems, tc.problem)
			}
		})
	}
}

func TestCheckNoNotice(t *testing.T) {
	named := map[string]bool{}
	checkClosureNotices("x.go", []byte("// Port of test/com/google/javascript/jscomp/FooTest.java\n\n"+apacheNotice+"\n\npackage x\n"), named)
	if !named["com/google/javascript/jscomp/FooTest.java"] {
		t.Errorf("named = %v, want the source without its test/ prefix", named)
	}
	for p := range noNotice {
		named[p] = true
	}
	if problems := checkNoNotice(named); len(problems) > 0 {
		t.Errorf("unexpected problems: %q", problems)
	}
	const stale = "com/google/javascript/jscomp/AutoValue_PassFactory.java"
	delete(named, stale)
	if problems := checkNoNotice(named); len(problems) != 1 || !strings.Contains(problems[0], stale) {
		t.Errorf("problems = %q, want one naming %s", problems, stale)
	}
	for p, reason := range noNotice {
		if reason == "" {
			t.Errorf("noNotice[%s] gives no reason", p)
		}
	}
}

func TestJavaSourceFromCache(t *testing.T) {
	t.Chdir(t.TempDir())
	write := func(p, s string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(s), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	version := "package com.google.javascript.jscomp;\nclass CommandLineRunnerVersion {\n  static final String COMPILER_VERSION = \"" + tag + "\";\n}\n"
	write(srcDir+"/CommandLineRunnerVersion.java", version)
	write(srcDir+"/com/google/javascript/jscomp/Foo.java", "package com.google.javascript.jscomp;\nclass Foo {}\n")
	write(testSrcDir+"/com/google/javascript/jscomp/FooTest.java", "package com.google.javascript.jscomp;\nclass FooTest {}\n")
	for p, want := range map[string]string{
		"com/google/javascript/jscomp/Foo.java":                      "class Foo {}",
		"test/com/google/javascript/jscomp/FooTest.java":             "class FooTest {}",
		"com/google/javascript/jscomp/CommandLineRunnerVersion.java": "class CommandLineRunnerVersion {",
	} {
		b, err := javaSource(p)
		if err != nil || !strings.Contains(string(b), want) {
			t.Errorf("javaSource(%s) = %q, %v; want the source with %q", p, b, err, want)
		}
	}
	// Sources of another release are refused, not silently used.
	write(srcDir+"/CommandLineRunnerVersion.java", strings.ReplaceAll(version, tag, "v20000101"))
	if _, err := javaSource("com/google/javascript/jscomp/Foo.java"); err == nil || !strings.Contains(err.Error(), "another release") {
		t.Errorf("javaSource with stale sources: err = %v, want another release", err)
	}
}
