// Command maprangecheck flags iteration over Go maps whose order is not
// accounted for.
//
// Go randomises map iteration order, while the compiler's output must be
// byte-identical with the Java jar's, whose collections iterate in a fixed
// order (insertion order for the Linked* collections). A `for ... range m`
// over a map, or a call of maps.Keys, maps.Values or maps.All, in a
// non-test file of the checked packages must therefore carry a comment
// that says why its order cannot reach anything order-sensitive (output,
// diagnostics, generated names or ids, list order):
//
//	// order-independent: only counts the entries
//	for _, v := range m {
//
// The comment is on the line of the statement or on the line just above
// it, starts with "order-independent:" and gives a reason. Everything
// else is reported, and the command exits with status 1.
//
//	go run ./tools/maprangecheck [-all] [packages]   (default: ./internal/...)
//
// -all also lists the annotated sites with their reasons. The command uses
// only the standard library: it type-checks each package from source with
// go/types against the export data that `go list -export` produces for its
// dependencies.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// DefaultPackages are the packages the guard covers.
var DefaultPackages = []string{"./internal/..."}

// Marker starts the annotation of an order-independent site.
const Marker = "order-independent:"

// Site is one map iteration.
type Site struct {
	Pos    token.Position
	What   string // "range over m" or "maps.Keys(m)"
	Reason string // the annotation's reason; "" when unannotated
}

func main() {
	all := flag.Bool("all", false, "also list annotated sites")
	flag.Parse()
	pkgs := flag.Args()
	if len(pkgs) == 0 {
		pkgs = DefaultPackages
	}
	sites, err := Check(".", pkgs)
	if err != nil {
		fmt.Fprintln(os.Stderr, "maprangecheck:", err)
		os.Exit(2)
	}
	bad := 0
	for _, s := range sites {
		if s.Reason == "" {
			bad++
			fmt.Printf("%s: %s without an %q comment\n", s.Pos, s.What, "// "+Marker+" <reason>")
		} else if *all {
			fmt.Printf("%s: %s: %s\n", s.Pos, s.What, s.Reason)
		}
	}
	if bad > 0 {
		fmt.Fprintf(os.Stderr, "maprangecheck: %d unannotated map iteration(s); keep Java's order (an insertion-ordered slice, or sorted keys) or annotate why the order cannot matter\n", bad)
		os.Exit(1)
	}
}

type listedPackage struct {
	ImportPath string
	Dir        string
	Export     string
	GoFiles    []string
	CgoFiles   []string
	DepOnly    bool
	Standard   bool
	Error      *struct{ Err string }
}

// Check lists the map iterations in the non-test files of the packages
// matched by patterns, resolved in dir.
func Check(dir string, patterns []string) ([]Site, error) {
	args := append([]string{"list", "-e", "-export", "-deps", "-json"}, patterns...)
	cmd := exec.Command("go", args...)
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("go list: %v\n%s", err, stderr.Bytes())
	}
	exports := map[string]string{}
	var targets []*listedPackage
	dec := json.NewDecoder(bytes.NewReader(out))
	for {
		p := new(listedPackage)
		if err := dec.Decode(p); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if p.Error != nil {
			return nil, fmt.Errorf("%s: %s", p.ImportPath, p.Error.Err)
		}
		exports[p.ImportPath] = p.Export
		if !p.DepOnly && !p.Standard {
			targets = append(targets, p)
		}
	}
	fset := token.NewFileSet()
	imp := importer.ForCompiler(fset, "gc", func(path string) (io.ReadCloser, error) {
		f, ok := exports[path]
		if !ok || f == "" {
			return nil, fmt.Errorf("no export data for %s", path)
		}
		return os.Open(f)
	})
	var sites []Site
	for _, p := range targets {
		var files []*ast.File
		for _, name := range append(p.GoFiles, p.CgoFiles...) {
			f, err := parser.ParseFile(fset, filepath.Join(p.Dir, name), nil, parser.ParseComments|parser.SkipObjectResolution)
			if err != nil {
				return nil, err
			}
			files = append(files, f)
		}
		info := &types.Info{Types: map[ast.Expr]types.TypeAndValue{}, Uses: map[*ast.Ident]types.Object{}}
		conf := types.Config{Importer: imp}
		if _, err := conf.Check(p.ImportPath, fset, files, info); err != nil {
			return nil, fmt.Errorf("type-check %s: %v", p.ImportPath, err)
		}
		for _, f := range files {
			sites = append(sites, fileSites(fset, f, info)...)
		}
	}
	sort.SliceStable(sites, func(i, j int) bool {
		a, b := sites[i].Pos, sites[j].Pos
		if a.Filename != b.Filename {
			return a.Filename < b.Filename
		}
		return a.Offset < b.Offset
	})
	return sites, nil
}

func fileSites(fset *token.FileSet, f *ast.File, info *types.Info) []Site {
	// reasons[line] is the annotation that covers that line.
	reasons := map[int]string{}
	for _, cg := range f.Comments {
		for _, c := range cg.List {
			text := strings.TrimSpace(strings.TrimPrefix(c.Text, "//"))
			if !strings.HasPrefix(text, Marker) {
				continue
			}
			reason := strings.TrimSpace(strings.TrimPrefix(text, Marker))
			if reason == "" {
				continue
			}
			line := fset.Position(c.Pos()).Line
			reasons[line] = reason
			reasons[line+1] = reason
		}
	}
	var sites []Site
	add := func(n ast.Node, what string) {
		pos := fset.Position(n.Pos())
		sites = append(sites, Site{Pos: pos, What: what, Reason: reasons[pos.Line]})
	}
	ast.Inspect(f, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.RangeStmt:
			if isMap(info.TypeOf(n.X)) {
				add(n, "range over map "+types.ExprString(n.X))
			}
		case *ast.CallExpr:
			fun := n.Fun
			if ix, ok := fun.(*ast.IndexExpr); ok {
				fun = ix.X
			} else if ix, ok := fun.(*ast.IndexListExpr); ok {
				fun = ix.X
			}
			sel, ok := fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			fn, ok := info.Uses[sel.Sel].(*types.Func)
			if !ok || fn.Pkg() == nil || fn.Pkg().Path() != "maps" {
				return true
			}
			switch fn.Name() {
			case "Keys", "Values", "All":
				add(n, types.ExprString(n))
			}
		}
		return true
	})
	return sites
}

// isMap reports whether t is a map, or a type parameter whose type set
// holds only maps.
func isMap(t types.Type) bool {
	if t == nil {
		return false
	}
	if tp, ok := types.Unalias(t).(*types.TypeParam); ok {
		iface, _ := tp.Underlying().(*types.Interface)
		if iface == nil || iface.NumEmbeddeds() == 0 {
			return false
		}
		for i := 0; i < iface.NumEmbeddeds(); i++ {
			if u, ok := iface.EmbeddedType(i).(*types.Union); ok {
				for j := 0; j < u.Len(); j++ {
					if !isMap(u.Term(j).Type()) {
						return false
					}
				}
			} else if !isMap(iface.EmbeddedType(i)) {
				return false
			}
		}
		return true
	}
	_, ok := t.Underlying().(*types.Map)
	return ok
}
