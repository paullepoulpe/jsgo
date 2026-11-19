// Clean-room modular Go AST + go/types post-pass loader.

package passes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// moduleInfo describes a Go module root discovered from a go.mod file.
type moduleInfo struct {
	modPath string
	rootDir string
}

// findModuleRoots discovers enclosing go.mod files starting from dirs and cwd.
func findModuleRoots(startDirs ...string) []moduleInfo {
	var roots []moduleInfo
	seen := make(map[string]bool)

	if wd, err := os.Getwd(); err == nil {
		startDirs = append(startDirs, wd)
	}

	for _, start := range startDirs {
		if start == "" {
			continue
		}
		dir, err := filepath.Abs(start)
		if err != nil {
			continue
		}
		if fi, err := os.Stat(dir); err == nil && !fi.IsDir() {
			dir = filepath.Dir(dir)
		}
		for {
			gomod := filepath.Join(dir, "go.mod")
			if data, err := os.ReadFile(gomod); err == nil {
				if !seen[dir] {
					if mp := parseModulePath(data); mp != "" {
						seen[dir] = true
						roots = append(roots, moduleInfo{modPath: mp, rootDir: dir})
					}
				}
				break
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}
	return roots
}

func parseModulePath(gomod []byte) string {
	for _, line := range strings.Split(string(gomod), "\n") {
		line = strings.TrimSpace(line)
		if idx := strings.Index(line, "//"); idx >= 0 {
			line = strings.TrimSpace(line[:idx])
		}
		if strings.HasPrefix(line, "module ") {
			rest := strings.TrimSpace(strings.TrimPrefix(line, "module "))
			rest = strings.Trim(rest, `"`)
			return rest
		}
	}
	return ""
}

// ModuleImporter implements types.Importer using stdlib export data for
// standard library packages and on-demand source type-checking for module
// and workspace packages.
type ModuleImporter struct {
	mu         sync.Mutex
	stdlib     types.Importer
	roots      []moduleInfo
	cache      map[string]*types.Package
	inProgress map[string]bool
	dirCache   map[string]string
}

// NewModuleImporter creates a cached module-aware types.Importer.
func NewModuleImporter(startDirs ...string) *ModuleImporter {
	return &ModuleImporter{
		stdlib:     importer.Default(),
		roots:      findModuleRoots(startDirs...),
		cache:      make(map[string]*types.Package),
		inProgress: make(map[string]bool),
		dirCache:   make(map[string]string),
	}
}

// HasPackage reports whether path is already present in the importer cache.
func (m *ModuleImporter) HasPackage(path string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.cache[path]
	return ok
}

// RegisterPackage pre-populates or updates a type-checked package in the cache,
// evicting any already-cached dependent packages that transitively referenced a
// previous *types.Package instance for path so go/types package identity stays consistent.
func (m *ModuleImporter) RegisterPackage(path string, pkg *types.Package) {
	m.mu.Lock()
	defer m.mu.Unlock()
	old, exists := m.cache[path]
	if exists && old == pkg {
		return
	}
	if exists && old != nil {
		// Collect keys to evict so no cached package holds a stale pointer to `old`.
		var evict []string
		// order-independent: map deletions do not depend on iteration order
		for k, cachedPkg := range m.cache {
			if k != path && cachedPkg != nil && pkgImportsTarget(cachedPkg, old, make(map[*types.Package]bool)) {
				evict = append(evict, k)
			}
		}
		for _, k := range evict {
			delete(m.cache, k)
		}
	}
	m.cache[path] = pkg
}

func pkgImportsTarget(cur, target *types.Package, visited map[*types.Package]bool) bool {
	if cur == nil || visited[cur] {
		return false
	}
	visited[cur] = true
	for _, imp := range cur.Imports() {
		if imp == target || pkgImportsTarget(imp, target, visited) {
			return true
		}
	}
	return false
}

// Import resolves and type-checks the package at path.
func (m *ModuleImporter) Import(path string) (*types.Package, error) {
	if path == "unsafe" {
		return types.Unsafe, nil
	}

	m.mu.Lock()
	if pkg, ok := m.cache[path]; ok {
		m.mu.Unlock()
		return pkg, nil
	}
	if m.inProgress[path] {
		m.mu.Unlock()
		return nil, fmt.Errorf("import cycle detected involving %q", path)
	}
	m.inProgress[path] = true
	m.mu.Unlock()

	defer func() {
		m.mu.Lock()
		delete(m.inProgress, path)
		m.mu.Unlock()
	}()

	// Standard library packages do not contain a '.' in their first path segment.
	firstSeg := path
	if slash := strings.IndexByte(path, '/'); slash >= 0 {
		firstSeg = path[:slash]
	}
	if !strings.Contains(firstSeg, ".") {
		pkg, err := m.stdlib.Import(path)
		if err == nil {
			m.mu.Lock()
			m.cache[path] = pkg
			m.mu.Unlock()
			return pkg, nil
		}
	}

	dir, err := m.resolvePackageDir(path)
	if err != nil {
		return nil, err
	}

	pkg, err := m.typecheckDir(path, dir)
	if err != nil {
		return nil, err
	}

	m.mu.Lock()
	m.cache[path] = pkg
	m.mu.Unlock()
	return pkg, nil
}

func (m *ModuleImporter) resolvePackageDir(importPath string) (string, error) {
	m.mu.Lock()
	if dir, ok := m.dirCache[importPath]; ok {
		m.mu.Unlock()
		return dir, nil
	}
	roots := m.roots
	m.mu.Unlock()

	for _, r := range roots {
		if importPath == r.modPath {
			if info, err := os.Stat(r.rootDir); err == nil && info.IsDir() {
				m.cacheDir(importPath, r.rootDir)
				return r.rootDir, nil
			}
		}
		prefix := r.modPath + "/"
		if strings.HasPrefix(importPath, prefix) {
			rel := filepath.FromSlash(strings.TrimPrefix(importPath, prefix))
			candidate := filepath.Join(r.rootDir, rel)
			if info, err := os.Stat(candidate); err == nil && info.IsDir() {
				m.cacheDir(importPath, candidate)
				return candidate, nil
			}
		}
	}

	// Fallback to `go list -f '{{.Dir}}' <importPath>` for external module dependencies.
	cmd := exec.Command("go", "list", "-f", "{{.Dir}}", importPath)
	if len(roots) > 0 {
		cmd.Dir = roots[0].rootDir
	}
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("cannot resolve package %q: %w", importPath, err)
	}
	dir := strings.TrimSpace(string(out))
	if dir == "" {
		return "", fmt.Errorf("empty directory for package %q", importPath)
	}
	m.cacheDir(importPath, dir)
	return dir, nil
}

func (m *ModuleImporter) cacheDir(importPath, dir string) {
	m.mu.Lock()
	m.dirCache[importPath] = dir
	m.mu.Unlock()
}

func (m *ModuleImporter) typecheckDir(importPath, dir string) (*types.Package, error) {
	goFiles, err := listPackageGoFiles(dir, false)
	if err != nil {
		return nil, err
	}
	if len(goFiles) == 0 {
		return nil, fmt.Errorf("no buildable Go source files in %s (%s)", dir, importPath)
	}

	fset := token.NewFileSet()
	astFiles := make([]*ast.File, 0, len(goFiles))
	for _, path := range goFiles {
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return nil, fmt.Errorf("parsing dependency %s: %w", path, err)
		}
		astFiles = append(astFiles, f)
	}

	conf := types.Config{
		Importer: m,
	}
	pkg, err := conf.Check(importPath, fset, astFiles, nil)
	if err != nil {
		return nil, fmt.Errorf("typechecking dependency %s: %w", importPath, err)
	}
	return pkg, nil
}

// listPackageGoFiles returns sorted absolute paths of buildable .go files in dir.
func listPackageGoFiles(dir string, includeTests bool) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	bctx := build.Default
	var files []string
	var testFiles []string
	var pkgName string

	for _, ent := range entries {
		if ent.IsDir() {
			continue
		}
		name := ent.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") {
			continue
		}
		isTest := strings.HasSuffix(name, "_test.go")
		if isTest && !includeTests {
			continue
		}
		match, err := bctx.MatchFile(dir, name)
		if err != nil || !match {
			continue
		}
		full := filepath.Join(dir, name)
		if isTest {
			testFiles = append(testFiles, full)
		} else {
			files = append(files, full)
			if pkgName == "" {
				if f, err := parser.ParseFile(token.NewFileSet(), full, nil, parser.PackageClauseOnly); err == nil {
					pkgName = f.Name.Name
				}
			}
		}
	}

	if includeTests && len(testFiles) > 0 {
		// Only include in-package test files (`package foo`), skipping external `package foo_test`
		// so a single types.Config.Check pass sees one coherent package name.
		for _, tf := range testFiles {
			f, err := parser.ParseFile(token.NewFileSet(), tf, nil, parser.PackageClauseOnly)
			if err != nil {
				continue
			}
			if pkgName == "" || f.Name.Name == pkgName {
				pkgName = f.Name.Name
				files = append(files, tf)
			}
		}
	}

	sort.Strings(files)
	return files, nil
}

// packageTarget represents a resolved Go package directory and optional file subset.
type packageTarget struct {
	importPath string
	dir        string
	onlyFiles  map[string]bool // nil means all package files in dir
}

// resolveTargets resolves CLI patterns (file paths, directories, `<dir>/...`, or Go package patterns)
// into deterministic packageTargets.
func resolveTargets(patterns []string, roots []moduleInfo, includeTests bool) ([]packageTarget, error) {
	if len(patterns) == 0 {
		patterns = []string{"."}
	}

	byDir := make(map[string]*packageTarget)
	var dirOrder []string

	addDir := func(dir string, onlyFile string) error {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			return err
		}
		t, exists := byDir[absDir]
		if !exists {
			goFiles, err := listPackageGoFiles(absDir, includeTests)
			if err != nil || len(goFiles) == 0 {
				return nil
			}
			t = &packageTarget{
				importPath: inferImportPath(absDir, roots),
				dir:        absDir,
			}
			if onlyFile != "" {
				t.onlyFiles = make(map[string]bool)
			}
			byDir[absDir] = t
			dirOrder = append(dirOrder, absDir)
		}
		if onlyFile != "" && t.onlyFiles != nil {
			t.onlyFiles[onlyFile] = true
		} else if onlyFile == "" {
			t.onlyFiles = nil
		}
		return nil
	}

	for _, pat := range patterns {
		if pat == "" {
			continue
		}
		// Check recursive filesystem pattern `<dir>/...` or `...`
		if pat == "..." || strings.HasSuffix(pat, "/...") || strings.HasSuffix(pat, string(filepath.Separator)+"...") {
			base := "."
			if pat != "..." {
				base = strings.TrimSuffix(strings.TrimSuffix(pat, "..."), "/")
				base = strings.TrimSuffix(base, string(filepath.Separator))
				if base == "" {
					base = "/"
				}
			}
			if fi, err := os.Stat(base); err == nil && fi.IsDir() {
				err := filepath.WalkDir(base, func(path string, d fs.DirEntry, walkErr error) error {
					if walkErr != nil {
						return walkErr
					}
					if !d.IsDir() {
						return nil
					}
					name := d.Name()
					if path != base && (strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") || name == "vendor" || name == "testdata") {
						return filepath.SkipDir
					}
					return addDir(path, "")
				})
				if err != nil {
					return nil, err
				}
				continue
			}
		}

		// Check direct file or directory path on disk.
		if fi, err := os.Stat(pat); err == nil {
			if fi.IsDir() {
				if err := addDir(pat, ""); err != nil {
					return nil, err
				}
				continue
			}
			if strings.HasSuffix(pat, ".go") {
				absFile, err := filepath.Abs(pat)
				if err != nil {
					return nil, err
				}
				if err := addDir(filepath.Dir(absFile), absFile); err != nil {
					return nil, err
				}
				continue
			}
		}

		// Fallback to `go list -json <pat>` for Go module package patterns.
		dirs, err := goListDirs(pat, roots)
		if err != nil {
			return nil, err
		}
		for _, d := range dirs {
			if err := addDir(d, ""); err != nil {
				return nil, err
			}
		}
	}

	sort.Strings(dirOrder)
	out := make([]packageTarget, 0, len(dirOrder))
	for _, d := range dirOrder {
		out = append(out, *byDir[d])
	}
	return topoSortTargets(out, includeTests), nil
}

// topoSortTargets orders resolved targets in deterministic topological dependency order
// (dependencies before dependents, breaking ties by the input lexicographical dirOrder)
// so cross-package method synthesis (e.g. Slice_<Getter>() in rhino) is visible to
// downstream packages (e.g. jscomp) in a single gopasses pass.
func topoSortTargets(targets []packageTarget, includeTests bool) []packageTarget {
	n := len(targets)
	if n <= 1 {
		return targets
	}
	idxByImport := make(map[string]int, n)
	for i, t := range targets {
		if t.importPath != "" {
			idxByImport[t.importPath] = i
		}
	}
	if len(idxByImport) <= 1 {
		return targets
	}

	inDegree := make([]int, n)
	dependents := make([][]int, n)
	fset := token.NewFileSet()

	for i, t := range targets {
		goFiles, err := listPackageGoFiles(t.dir, includeTests)
		if err != nil {
			continue
		}
		seenDep := make(map[int]bool)
		var deps []int
		for _, fp := range goFiles {
			af, err := parser.ParseFile(fset, fp, nil, parser.ImportsOnly)
			if err != nil {
				continue
			}
			for _, imp := range af.Imports {
				if imp.Path == nil {
					continue
				}
				impPath, err := strconv.Unquote(imp.Path.Value)
				if err != nil {
					continue
				}
				if depIdx, ok := idxByImport[impPath]; ok && depIdx != i && !seenDep[depIdx] {
					seenDep[depIdx] = true
					deps = append(deps, depIdx)
				}
			}
		}
		sort.Ints(deps)
		inDegree[i] = len(deps)
		for _, depIdx := range deps {
			dependents[depIdx] = append(dependents[depIdx], i)
		}
	}

	var ready []int
	for i := 0; i < n; i++ {
		if inDegree[i] == 0 {
			ready = append(ready, i)
		}
	}

	sorted := make([]packageTarget, 0, n)
	emitted := make([]bool, n)
	for len(ready) > 0 {
		u := ready[0]
		ready = ready[1:]
		emitted[u] = true
		sorted = append(sorted, targets[u])
		for _, v := range dependents[u] {
			inDegree[v]--
			if inDegree[v] == 0 {
				insertAt := sort.SearchInts(ready, v)
				ready = append(ready, 0)
				copy(ready[insertAt+1:], ready[insertAt:])
				ready[insertAt] = v
			}
		}
	}

	if len(sorted) < n {
		for i := 0; i < n; i++ {
			if !emitted[i] {
				sorted = append(sorted, targets[i])
			}
		}
	}
	return sorted
}

func inferImportPath(absDir string, roots []moduleInfo) string {
	for _, r := range roots {
		if absDir == r.rootDir {
			return r.modPath
		}
		prefix := r.rootDir + string(filepath.Separator)
		if strings.HasPrefix(absDir, prefix) {
			rel := strings.TrimPrefix(absDir, prefix)
			return r.modPath + "/" + filepath.ToSlash(rel)
		}
	}
	return filepath.ToSlash(absDir)
}

func goListDirs(pattern string, roots []moduleInfo) ([]string, error) {
	cmd := exec.Command("go", "list", "-json", pattern)
	if len(roots) > 0 {
		cmd.Dir = roots[0].rootDir
	}
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("cannot resolve pattern %q: %w", pattern, err)
	}
	dec := json.NewDecoder(bytes.NewReader(out))
	var dirs []string
	for dec.More() {
		var item struct {
			Dir string `json:"Dir"`
		}
		if err := dec.Decode(&item); err != nil {
			return nil, err
		}
		if item.Dir != "" {
			dirs = append(dirs, item.Dir)
		}
	}
	return dirs, nil
}
