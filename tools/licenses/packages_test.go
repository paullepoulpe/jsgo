package main

import (
	"slices"
	"strings"
	"testing"
)

func TestPackageDir(t *testing.T) {
	for _, tc := range []struct{ path, want string }{
		{"internal/jscomp/testdata/typecheck/src/lodash@4.17.21/package/lodash.js", "internal/jscomp/testdata/typecheck/src/lodash@4.17.21/"},
		{"internal/jscomp/testdata/typecheck/src/@angular/core@19.2.14/package/x.cjs", "internal/jscomp/testdata/typecheck/src/@angular/core@19.2.14/"},
		{"internal/jscomp/testdata/es_modules/src/three@0.186.1/package/src/a@b/c.js", "internal/jscomp/testdata/es_modules/src/three@0.186.1/"},
		{"testdata/src/x@1/LICENSE", "testdata/src/x@1/"},
		{"internal/jscomp/testdata/typecheck/expected/lodash@4.17.21/package/lodash.js.typecheck.gz", ""},
		{"internal/jscomp/testdata/checks/src/vars.js", ""},
		{"internal/jscomp/testdata/closure_rewrite/src/custom/a.js", ""},
		{"src/x@1/a.js", ""},
	} {
		if got := packageDir.FindString(tc.path); got != tc.want {
			t.Errorf("packageDir(%q) = %q, want %q", tc.path, got, tc.want)
		}
	}
}

func TestCheckPackageLicences(t *testing.T) {
	const td = "internal/jscomp/testdata/"
	files := []string{
		td + "typecheck/src/a@1.0.0/package/a.js",
		td + "typecheck/src/a@1.0.0/LICENSE",
		td + "typecheck/src/@s/b@2.0.0/package/b.js",
		td + "typecheck/src/@s/b@2.0.0/LICENSE.txt",
		td + "typecheck/src/c@3.0.0/licenses.md",
		td + "typecheck/src/c@3.0.0/package/c.js",
		td + "es_modules/src/a@1.0.0/package/a.js",    // no licence
		td + "es_modules/src/d@1.0.0/package/LICENSE", // not next to the package
		td + "es_modules/src/d@1.0.0/package/d.js",
		td + "typecheck/src/plain.js",
	}
	got := checkPackageLicences(files)
	want := []string{td + "es_modules/src/a@1.0.0/", td + "es_modules/src/d@1.0.0/"}
	if len(got) != len(want) {
		t.Fatalf("problems = %q, want %d", got, len(want))
	}
	for i, p := range got {
		if !strings.HasPrefix(p, want[i]+" ") {
			t.Errorf("problem %d = %q, want one about %s", i, p, want[i])
		}
	}
	for _, name := range []string{"LICENSE", "LICENCE", "License.md", "LICENSE-MIT", "licenses.md", "COPYING", "LICENSE.txt"} {
		if !licenceFile.MatchString(name) {
			t.Errorf("licenceFile does not match %q", name)
		}
	}
	for _, name := range []string{"package", "licensed.js", "README.md", "LICENSE/x"} {
		if licenceFile.MatchString(name) {
			t.Errorf("licenceFile matches %q", name)
		}
	}
	if got := checkPackageLicences(slices.Clip(files[:6])); len(got) != 0 {
		t.Errorf("unexpected problems: %q", got)
	}
}
