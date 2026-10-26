package jutil

import (
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// TestUnsafeImportAllowlist guards against re-introducing "unsafe" anywhere
// in cmd/ or internal/.
func TestUnsafeImportAllowlist(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join("..", ".."))

	fset := token.NewFileSet()
	for _, top := range []string{"cmd", "internal"} {
		err := filepath.WalkDir(filepath.Join(repoRoot, top), func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			rel, err := filepath.Rel(repoRoot, path)
			if err != nil {
				return err
			}
			rel = filepath.ToSlash(rel)

			f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
			if err != nil {
				t.Fatalf("parse imports of %s: %v", rel, err)
				return nil
			}
			for _, imp := range f.Imports {
				pkg, err := strconv.Unquote(imp.Path.Value)
				if err != nil || pkg != "unsafe" {
					continue
				}
				t.Errorf("%s imports \"unsafe\" (zero unsafe imports allowed)", rel)
			}
			return nil
		})
		if err != nil {
			t.Fatalf("walk %s: %v", top, err)
		}
	}
}
