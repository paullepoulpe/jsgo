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

// TestUnsafeImportAllowlist guards against re-introducing "unsafe" outside the
// small, documented runtime/substrate files (Gson private-field reflection,
// pointer identity hashing, slab chunk sizing, RhinoStringPool byte pages, and
// the AST Node storage substrate).
func TestUnsafeImportAllowlist(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	allowed := map[string]bool{
		"internal/jutil/gson_from_json_reflective.go": true,
		"internal/jutil/gson_to_json.go":              true,
		"internal/jutil/objects.go":                   true,
		"internal/jutil/slab.go":                      true,
		"internal/strpool/strpool.go":                 true, // the string pool (was internal/rhino/rhino_string_pool.go)
	}

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
				if allowed[rel] ||
					strings.HasPrefix(rel, "internal/rhino/node") ||
					strings.HasPrefix(rel, "internal/rhino/flatspike/") {
					continue
				}
				t.Errorf("%s imports \"unsafe\" outside the documented substrate allowlist", rel)
			}
			return nil
		})
		if err != nil {
			t.Fatalf("walk %s: %v", top, err)
		}
	}
}
