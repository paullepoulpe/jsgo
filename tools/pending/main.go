// Command pending generates internal/jscomp/pending_gen.go: Go declarations
// for the not-yet-ported Java classes that the files of package
// internal/jscomp (tests included) reference (PORTING.md §1, the one
// sanctioned exception), so that the package builds, vets and tests while
// the port proceeds.
//
// It type-checks package internal/jscomp without pending_gen.go, collects
// from the root files (-roots, by default every non-generated file) every
// identifier that does not resolve and
// every selector on a pending type, maps each back to the Java member it
// names by PORTING.md's naming rules (overloads are told apart by the call
// site's arity and argument types), reads that member's signature from the
// pinned jar, and declares it with the Go types the ports use. Bodies panic
// with &PendingPort{"<Java signature>"} ("pending port: ..."), which
// RunJUnit4 reports as a skip. It repeats until nothing new
// resolves (a declared return type lets the next selector resolve), so the
// result holds exactly the members the root files use. Classes that have a
// real port (a "// Port of" header names their Java file, or their Go type
// exists) are never declared: their declarations vanish as soon as the port
// lands.
//
// Usage:
//
//	go run ./tools/pending           # regenerate internal/jscomp/pending_gen.go
//	go run ./tools/pending -check    # fail if the file is stale or duplicates a real declaration
//	go run ./tools/pending -v        # also list the references it could not resolve
//	go run ./tools/pending -roots compiler.go,'*_test.go'   # only those files' references
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const generatedName = "pending_gen.go"

// allRoots makes every non-generated file of the package (tests included)
// a root: main must build, vet and test while files refer to classes that
// are not ported yet.
const allRoots = "all"

func main() {
	dir := flag.String("dir", "internal/jscomp", "the Go package directory")
	jar := flag.String("jar", "", "the pinned closure-compiler jar (default: .cache/closure-compiler.jar of this checkout or of the main worktree)")
	roots := flag.String("roots", allRoots, `comma-separated root file names or globs (in -dir), or "all" for every non-generated file`)
	check := flag.Bool("check", false, "only check: fail if "+generatedName+" is stale or declares something a real file declares")
	verbose := flag.Bool("v", false, "list unresolved references")
	flag.Parse()

	if *jar == "" {
		*jar = defaultJar()
	}
	if err := run(*dir, *jar, strings.Split(*roots, ","), *check, *verbose); err != nil {
		fmt.Fprintln(os.Stderr, "pending:", err)
		os.Exit(1)
	}
}

func run(dir, jar string, rootPatterns []string, check, verbose bool) error {
	pkg, err := loadGoPackage(dir, generatedName)
	if err != nil {
		return err
	}
	outPath := filepath.Join(dir, generatedName)
	if check {
		if err := checkDuplicates(pkg, outPath); err != nil {
			return err
		}
		if !fileExists(jar) {
			// Without the pinned jar (tools/oracle/fetch.sh) only the
			// duplicate check can run.
			fmt.Fprintf(os.Stderr, "pending: %s not found; checked for duplicate declarations only\n", jar)
			return nil
		}
	}
	classes, err := readJar(jar, []string{"com/google/javascript/", "com/google/debugging/"})
	if err != nil {
		return err
	}
	model := newModel(classes)
	roots := map[string]bool{}
	for name := range pkg.files {
		for _, pat := range rootPatterns {
			pat = strings.TrimSpace(pat)
			if pat == allRoots {
				roots[name] = true
			} else if ok, _ := filepath.Match(pat, name); ok {
				roots[name] = true
			}
		}
	}

	g := newGen(model, pkg)
	handled := map[string]bool{}
	var src []byte
	var res *checkResult
	for iter := 0; ; iter++ {
		if iter > 60 {
			return fmt.Errorf("no fixed point after %d iterations", iter)
		}
		res, err = pkg.typecheck(src)
		if err != nil {
			return err
		}
		g.res = res
		pendingTypes := map[string]bool{}
		for n := range g.classes {
			pendingTypes[n] = true
		}
		needs := collectNeeds(res, roots, pendingTypes, pkg.fset, func(string) {})
		progress := false
		for _, n := range needs {
			if handled[n.key()] {
				continue
			}
			handled[n.key()] = true
			progress = true
			if n.owner == "" {
				g.resolveTop(n)
			} else {
				g.resolveMember(n)
			}
		}
		for g.complete() {
			progress = true
		}
		// Rendering can declare more classes (types used in signatures):
		// render until stable.
		var out []byte
		for {
			before := len(g.classes)
			out = g.render()
			for g.complete() {
			}
			if len(g.classes) == before {
				out = g.render()
				break
			}
		}
		if bytes.Equal(out, src) && !progress {
			break
		}
		src = out
	}
	// Report what is left: unresolved references and type errors in the
	// root files.
	var rootErrs []string
	for _, e := range res.errs {
		pos := pkg.fset.Position(e.Pos)
		if roots[filepath.Base(pos.Filename)] || filepath.Base(pos.Filename) == generatedName {
			rootErrs = append(rootErrs, fmt.Sprintf("%s: %s", pos, e.Msg))
		}
	}
	sort.Strings(g.warns)
	if verbose {
		for _, w := range g.warns {
			fmt.Fprintln(os.Stderr, "unresolved:", w)
		}
		for _, e := range rootErrs {
			fmt.Fprintln(os.Stderr, "error:", e)
		}
	}
	nDecls := 0
	for _, p := range g.classes {
		nDecls += 1 + len(p.members) + len(p.statics)
	}
	fmt.Fprintf(os.Stderr, "pending: %d classes, %d declarations, %d unresolved references, %d type errors in root files\n",
		len(g.classes), nDecls, len(g.warns), len(rootErrs))

	if len(g.classes) == 0 {
		src = nil
	}
	if check {
		existing, err := os.ReadFile(outPath)
		if os.IsNotExist(err) {
			existing = nil
		} else if err != nil {
			return err
		}
		if !bytes.Equal(existing, src) {
			return fmt.Errorf("%s is stale: run `go run ./tools/pending`", outPath)
		}
		return nil
	}
	if src == nil {
		if err := os.Remove(outPath); err != nil && !os.IsNotExist(err) {
			return err
		}
		return nil
	}
	return os.WriteFile(outPath, src, 0o644)
}

// defaultJar finds .cache/closure-compiler.jar in this checkout or, for a
// git worktree, in the main checkout (tools/oracle/fetch.sh puts it there).
func defaultJar() string {
	const rel = ".cache/closure-compiler.jar"
	if _, err := os.Stat(rel); err == nil {
		return rel
	}
	out, err := exec.Command("git", "rev-parse", "--path-format=absolute", "--git-common-dir").Output()
	if err == nil {
		main := filepath.Dir(strings.TrimSpace(string(out)))
		if p := filepath.Join(main, rel); fileExists(p) {
			return p
		}
	}
	return rel
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

// checkDuplicates fails if the generated file declares a name (or a method)
// that a non-generated file of the package also declares.
func checkDuplicates(pkg *goPackage, outPath string) error {
	data, err := os.ReadFile(outPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	fset := token.NewFileSet()
	gen, err := parser.ParseFile(fset, outPath, data, 0)
	if err != nil {
		return err
	}
	real := map[string]string{}
	for name, f := range pkg.files {
		for k := range declKeys(f) {
			real[k] = name
		}
	}
	var dups []string
	for k := range declKeys(gen) {
		if where, ok := real[k]; ok {
			dups = append(dups, fmt.Sprintf("%s (also declared in %s)", k, where))
		}
	}
	if len(dups) > 0 {
		sort.Strings(dups)
		return fmt.Errorf("%s duplicates real declarations; run `go run ./tools/pending`:\n  %s", outPath, strings.Join(dups, "\n  "))
	}
	return nil
}

// declKeys lists a file's package-level names and methods (Type.Method).
func declKeys(f *ast.File) map[string]bool {
	keys := map[string]bool{}
	for _, d := range f.Decls {
		switch d := d.(type) {
		case *ast.FuncDecl:
			if d.Recv == nil {
				keys[d.Name.Name] = true
				continue
			}
			t := d.Recv.List[0].Type
			if s, ok := t.(*ast.StarExpr); ok {
				t = s.X
			}
			switch x := t.(type) {
			case *ast.IndexExpr:
				t = x.X
			case *ast.IndexListExpr:
				t = x.X
			}
			if id, ok := t.(*ast.Ident); ok {
				keys[id.Name+"."+d.Name.Name] = true
			}
		case *ast.GenDecl:
			for _, s := range d.Specs {
				switch s := s.(type) {
				case *ast.TypeSpec:
					keys[s.Name.Name] = true
				case *ast.ValueSpec:
					for _, n := range s.Names {
						if n.Name != "_" {
							keys[n.Name] = true
						}
					}
				}
			}
		}
	}
	return keys
}
