// Clean-room modular Go AST + go/types post-pass framework.

package passes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"os"
	"strings"
)

// Pass is an independently toggleable Go-to-Go AST + go/types optimization pass.
type Pass interface {
	Name() string
	Description() string
	Apply(ctx *PassContext) (int, error)
}

// PassContext carries the parsed ASTs, token positions, and full go/types type information
// for a single Go package.
type PassContext struct {
	Fset          *token.FileSet
	Pkg           *types.Package
	Info          *types.Info
	Files         []*ast.File
	TargetFiles   map[*ast.File]bool // nil means all Files are transform targets
	ModifiedFiles map[*ast.File]bool // optional per-file dirty tracking populated by passes
	Importer      types.Importer
}

// ShouldTransform reports whether file is in scope for transformation.
func (c *PassContext) ShouldTransform(file *ast.File) bool {
	if c.TargetFiles == nil {
		return true
	}
	return c.TargetFiles[file]
}

// MarkFileModified records that a pass mutated file's AST so unchanged files in large
// packages can skip redundant import pruning and gofmt reformatting.
func (c *PassContext) MarkFileModified(file *ast.File) {
	if c.ModifiedFiles == nil {
		c.ModifiedFiles = make(map[*ast.File]bool)
	}
	c.ModifiedFiles[file] = true
}

// Options configures a pass pipeline run.
type Options struct {
	DryRun       bool // Compute transformations, formatting, and verification without writing disk files.
	IncludeTests bool // Also transform in-package *_test.go files.
}

// Stats records exact ground-truth transformation counts per pass and across files/packages.
type Stats struct {
	PassOrder     []string       `json:"pass_order"`
	ByPass        map[string]int `json:"by_pass"`
	TotalSites    int            `json:"total_sites"`
	FilesModified int            `json:"files_modified"`
	Packages      int            `json:"packages"`
}

func newStats(passList []Pass) *Stats {
	order := make([]string, 0, len(passList))
	byPass := make(map[string]int, len(passList))
	for _, p := range passList {
		if _, seen := byPass[p.Name()]; !seen {
			order = append(order, p.Name())
			byPass[p.Name()] = 0
		}
	}
	return &Stats{
		PassOrder: order,
		ByPass:    byPass,
	}
}

func (s *Stats) merge(other *Stats) {
	if other == nil {
		return
	}
	for _, name := range other.PassOrder {
		if _, exists := s.ByPass[name]; !exists {
			s.PassOrder = append(s.PassOrder, name)
		}
		s.ByPass[name] += other.ByPass[name]
	}
	s.TotalSites += other.TotalSites
	s.FilesModified += other.FilesModified
	s.Packages += other.Packages
}

// WriteText formats human-readable ground-truth pass statistics to w.
func (s *Stats) WriteText(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "=== gopasses statistics ===\n"); err != nil {
		return err
	}
	if len(s.PassOrder) == 0 {
		_, err := fmt.Fprintf(w, "passes enabled:  none\npackages:        %d\nfiles modified:  %d\ntotal sites:     %d\n",
			s.Packages, s.FilesModified, s.TotalSites)
		return err
	}
	for _, name := range s.PassOrder {
		if _, err := fmt.Fprintf(w, "  %-20s %d sites\n", name+":", s.ByPass[name]); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(w, "  %-20s %d sites (%d files modified across %d packages)\n",
		"total:", s.TotalSites, s.FilesModified, s.Packages)
	return err
}

// WriteJSON writes machine-readable JSON statistics to w.
func (s *Stats) WriteJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

// Run resolves patterns, type-checks each package, executes passList to a fixed point,
// prunes unused imports last, verifies go/types soundness, and writes modified files unless opts.DryRun is set.
func Run(patterns []string, passList []Pass, opts Options) (*Stats, error) {
	var hintDirs []string
	for _, p := range patterns {
		cleaned := strings.TrimSuffix(strings.TrimSuffix(p, "..."), "/")
		if cleaned != "" {
			if _, err := os.Stat(cleaned); err == nil {
				hintDirs = append(hintDirs, cleaned)
			}
		}
	}
	imp := NewModuleImporter(hintDirs...)
	targets, err := resolveTargets(patterns, imp.roots, opts.IncludeTests)
	if err != nil {
		return nil, err
	}

	totalStats := newStats(passList)
	for _, tgt := range targets {
		pkgStats, err := runOnTarget(tgt, passList, opts, imp)
		if err != nil {
			return totalStats, err
		}
		totalStats.merge(pkgStats)
	}
	return totalStats, nil
}

// RunOnPackageDir runs passList over a single Go package directory.
func RunOnPackageDir(dir string, passList []Pass, opts Options) (*Stats, error) {
	return Run([]string{dir}, passList, opts)
}

// TransformSource runs passList in-memory on a single-file Go package source string.
func TransformSource(filename, src string, passList []Pass) (string, *Stats, error) {
	if filename == "" {
		filename = "input.go"
	}
	imp := NewModuleImporter()
	memFiles := []memFile{{name: filename, src: []byte(src), target: true}}
	outFiles, stats, err := executePipeline("mem/pkg", memFiles, passList, imp)
	if err != nil {
		return "", nil, err
	}
	return string(outFiles[0].src), stats, nil
}

type memFile struct {
	name     string
	src      []byte
	target   bool
	modified bool
}

func runOnTarget(tgt packageTarget, passList []Pass, opts Options, imp *ModuleImporter) (*Stats, error) {
	filePaths, err := listPackageGoFiles(tgt.dir, opts.IncludeTests)
	if err != nil {
		return nil, err
	}
	if len(filePaths) == 0 {
		return newStats(passList), nil
	}

	files := make([]memFile, len(filePaths))
	for i, fp := range filePaths {
		raw, err := os.ReadFile(fp)
		if err != nil {
			return nil, err
		}
		files[i] = memFile{
			name:   fp,
			src:    raw,
			target: tgt.onlyFiles == nil || tgt.onlyFiles[fp],
		}
	}

	outFiles, stats, err := executePipeline(tgt.importPath, files, passList, imp)
	if err != nil {
		return nil, err
	}

	if !opts.DryRun {
		for _, f := range outFiles {
			if !f.modified {
				continue
			}
			mode := os.FileMode(0o644)
			if fi, statErr := os.Stat(f.name); statErr == nil {
				mode = fi.Mode().Perm()
			}
			if err := os.WriteFile(f.name, f.src, mode); err != nil {
				return stats, fmt.Errorf("writing %s: %w", f.name, err)
			}
		}
	}
	return stats, nil
}

func executePipeline(importPath string, files []memFile, passList []Pass, imp *ModuleImporter) ([]memFile, *Stats, error) {
	stats := newStats(passList)
	stats.Packages = 1

	fset, astFiles, pkg, info, err := parseAndTypecheck(importPath, files, imp, nil)
	if err != nil {
		return nil, nil, err
	}

	if len(passList) == 0 {
		if !imp.HasPackage(importPath) {
			imp.RegisterPackage(importPath, pkg)
		}
		return files, stats, nil
	}

	anyChanged := false
	changedByFile := make([]bool, len(files))

	recheckWithImportCleanup := func(stepLabel string) error {
		unusedPos := make(map[token.Pos]bool)
		checkFset, checkASTs, nextPkg, nextInfo, chkErr := parseAndTypecheck(importPath, files, imp, unusedPos)
		if chkErr != nil && len(unusedPos) > 0 {
			for i, af := range checkASTs {
				if !files[i].target {
					continue
				}
				if pruneUnusedImports(af, imp, unusedPos) {
					formatted, fmtErr := formatASTFile(checkFset, af)
					if fmtErr != nil {
						return fmtErr
					}
					files[i].src = formatted
					changedByFile[i] = true
				}
			}
			checkFset, checkASTs, nextPkg, nextInfo, chkErr = parseAndTypecheck(importPath, files, imp, nil)
		}
		if chkErr != nil {
			return fmt.Errorf("post-pass verification failed after %s on %s: %w", stepLabel, importPath, chkErr)
		}
		fset = checkFset
		astFiles = checkASTs
		pkg = nextPkg
		info = nextInfo
		return nil
	}

	runSinglePass := func(p Pass) (int, error) {
		targetMap := make(map[*ast.File]bool, len(astFiles))
		for i, af := range astFiles {
			if files[i].target {
				targetMap[af] = true
			}
		}

		ctx := &PassContext{
			Fset:          fset,
			Pkg:           pkg,
			Info:          info,
			Files:         astFiles,
			TargetFiles:   targetMap,
			ModifiedFiles: make(map[*ast.File]bool),
			Importer:      imp,
		}

		count, err := p.Apply(ctx)
		if err != nil {
			return 0, fmt.Errorf("pass %s failed on %s: %w", p.Name(), importPath, err)
		}
		stats.ByPass[p.Name()] += count
		stats.TotalSites += count
		if count == 0 {
			return 0, nil
		}
		anyChanged = true

		trackDirty := len(ctx.ModifiedFiles) > 0
		for i, af := range astFiles {
			if !files[i].target || (trackDirty && !ctx.ModifiedFiles[af]) {
				continue
			}
			pruneUnusedImports(af, imp, nil)
			formatted, fmtErr := formatASTFile(fset, af)
			if fmtErr != nil {
				return 0, fmt.Errorf("formatting %s after pass %s: %w", files[i].name, p.Name(), fmtErr)
			}
			if !bytes.Equal(files[i].src, formatted) {
				files[i].src = formatted
				changedByFile[i] = true
			}
		}

		if err := recheckWithImportCleanup(p.Name()); err != nil {
			return 0, err
		}
		return count, nil
	}

	// Iterate passes to fixed point (bounded to 4 rounds) so unwrapping an outer CheckNotNull*
	// never leaves behind an exposed AsIface[T](x) or T(x) for a second gopasses invocation.
	const maxRounds = 4
	for round := 0; round < maxRounds; round++ {
		roundCount := 0
		for _, p := range passList {
			c, err := runSinglePass(p)
			if err != nil {
				return nil, nil, err
			}
			roundCount += c
		}
		if roundCount == 0 {
			break
		}
	}

	// Final import_cleanup pass LAST across any modified target files.
	if anyChanged {
		finalPruned := false
		for i, af := range astFiles {
			if !files[i].target || !changedByFile[i] {
				continue
			}
			if pruneUnusedImports(af, imp, nil) {
				formatted, fmtErr := formatASTFile(fset, af)
				if fmtErr != nil {
					return nil, nil, fmtErr
				}
				if !bytes.Equal(files[i].src, formatted) {
					files[i].src = formatted
					finalPruned = true
				}
			}
		}
		if finalPruned {
			if err := recheckWithImportCleanup("import_cleanup"); err != nil {
				return nil, nil, err
			}
		}
		for i := range files {
			if changedByFile[i] {
				files[i].modified = true
				stats.FilesModified++
			}
		}
	}

	if importPath != "" {
		imp.RegisterPackage(importPath, pkg)
	}
	return files, stats, nil
}

func parseAndTypecheck(
	importPath string,
	files []memFile,
	imp types.Importer,
	recordUnusedImportPos map[token.Pos]bool,
) (*token.FileSet, []*ast.File, *types.Package, *types.Info, error) {
	fset := token.NewFileSet()
	astFiles := make([]*ast.File, len(files))
	for i, f := range files {
		af, err := parser.ParseFile(fset, f.name, f.src, parser.ParseComments)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("parse error in %s: %w", f.name, err)
		}
		astFiles[i] = af
	}

	info := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
	}

	var firstErr error
	conf := types.Config{
		Importer: imp,
		Error: func(err error) {
			if firstErr == nil {
				firstErr = err
			}
			if recordUnusedImportPos != nil {
				if te, ok := err.(types.Error); ok && strings.Contains(te.Msg, "imported and not used") {
					recordUnusedImportPos[te.Pos] = true
				}
			}
		},
	}

	pkg, err := conf.Check(importPath, fset, astFiles, info)
	if err != nil {
		if firstErr != nil {
			return fset, astFiles, nil, nil, firstErr
		}
		return fset, astFiles, nil, nil, err
	}
	return fset, astFiles, pkg, info, nil
}

func formatASTFile(fset *token.FileSet, file *ast.File) ([]byte, error) {
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, file); err != nil {
		return nil, err
	}
	return format.Source(buf.Bytes())
}
