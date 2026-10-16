package main

// The Go side: type-check internal/jscomp (without pending_gen.go, with the
// declarations generated so far) and collect, from the root files, every
// reference that does not resolve.

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// goPackage is the parsed Go package directory.
type goPackage struct {
	dir   string
	fset  *token.FileSet
	files map[string]*ast.File // base name → AST, without the generated file
	// portedJava: Java source paths (com/google/.../Foo.java) named in a
	// "// Port of" header of a non-generated file.
	portedJava map[string]bool
	importer   types.Importer
}

func loadGoPackage(dir, generated string) (*goPackage, error) {
	p := &goPackage{dir: dir, fset: token.NewFileSet(), files: map[string]*ast.File{}, portedJava: map[string]bool{}}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || name == generated {
			continue
		}
		f, err := parser.ParseFile(p.fset, filepath.Join(dir, name), nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		if f.Name.Name != "jscomp" {
			continue // an external test package
		}
		p.files[name] = f
		// Header: the Java files this Go file ports.
		for i, cg := range f.Comments {
			if i > 1 || cg.Pos() > f.Package {
				break
			}
			for _, c := range cg.List {
				for _, w := range strings.FieldsFunc(c.Text, func(r rune) bool {
					return r == ' ' || r == '(' || r == ')' || r == ','
				}) {
					if strings.HasPrefix(w, "com/google/") && strings.HasSuffix(w, ".java") {
						p.portedJava[w] = true
					}
				}
			}
		}
	}
	if build.Default.GOROOT == "" {
		// A binary built by `go run` may not know its GOROOT; the source
		// importer needs it for the standard library.
		out, err := exec.Command("go", "env", "GOROOT").Output()
		if err != nil {
			return nil, fmt.Errorf("go env GOROOT: %v", err)
		}
		build.Default.GOROOT = strings.TrimSpace(string(out))
	}
	p.importer = importer.ForCompiler(p.fset, "source", nil)
	return p, nil
}

const jscompPath = "github.com/paullepoulpe/jsgo/internal/jscomp"

// checkResult is one type-checking run.
type checkResult struct {
	pkg     *types.Package
	info    *types.Info
	files   []*ast.File
	pending *ast.File
	errs    []types.Error
}

func (p *goPackage) typecheck(pendingSrc []byte) (*checkResult, error) {
	var files []*ast.File
	var names []string
	for n := range p.files {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		files = append(files, p.files[n])
	}
	res := &checkResult{}
	if pendingSrc != nil {
		f, err := parser.ParseFile(p.fset, filepath.Join(p.dir, "pending_gen.go"), pendingSrc, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("generated code does not parse: %v\n%s", err, pendingSrc)
		}
		files = append(files, f)
		res.pending = f
	}
	res.files = files
	res.info = &types.Info{
		Types:      map[ast.Expr]types.TypeAndValue{},
		Defs:       map[*ast.Ident]types.Object{},
		Uses:       map[*ast.Ident]types.Object{},
		Selections: map[*ast.SelectorExpr]*types.Selection{},
	}
	conf := types.Config{
		Importer: p.importer,
		Error: func(err error) {
			if te, ok := err.(types.Error); ok {
				res.errs = append(res.errs, te)
			}
		},
	}
	res.pkg, _ = conf.Check(jscompPath, p.fset, files, res.info)
	return res, nil
}

// callSite describes how an unresolved name is used.
type callSite struct {
	called   bool
	arity    int
	ellipsis bool
	argTypes []types.Type // nil entries for untyped/invalid arguments
	typePos  bool         // used where a type is expected
	pos      token.Position
}

// need is an unresolved reference from a root file.
type need struct {
	// owner is "" for a package-level name, else the Go type name whose
	// member is referenced.
	owner string
	name  string
	sites []callSite
}

func (n *need) key() string { return n.owner + "." + n.name }

// collectNeeds walks the root files and returns the unresolved references:
// package-level identifiers, and selectors on generated (pending) types.
func collectNeeds(res *checkResult, roots map[string]bool, pendingTypes map[string]bool, fset *token.FileSet, warn func(string)) []*need {
	needs := map[string]*need{}
	add := func(owner, name string, site callSite) {
		k := owner + "." + name
		n, ok := needs[k]
		if !ok {
			n = &need{owner: owner, name: name}
			needs[k] = n
		}
		n.sites = append(n.sites, site)
	}
	info := res.info
	for _, f := range res.files {
		fname := filepath.Base(fset.Position(f.Pos()).Filename)
		if !roots[fname] {
			continue
		}
		// Parent links for context.
		parents := map[ast.Node]ast.Node{}
		var stack []ast.Node
		ast.Inspect(f, func(n ast.Node) bool {
			if n == nil {
				stack = stack[:len(stack)-1]
				return true
			}
			if len(stack) > 0 {
				parents[n] = stack[len(stack)-1]
			}
			stack = append(stack, n)
			return true
		})
		siteOf := func(e ast.Expr) callSite {
			s := callSite{pos: fset.Position(e.Pos())}
			p := parents[e]
			if call, ok := p.(*ast.CallExpr); ok && call.Fun == e {
				s.called = true
				s.arity = len(call.Args)
				s.ellipsis = call.Ellipsis.IsValid()
				for _, a := range call.Args {
					var t types.Type
					if tv, ok := info.Types[a]; ok && tv.Type != nil {
						if b, isBasic := tv.Type.(*types.Basic); !isBasic || b.Kind() != types.Invalid {
							if tv.IsNil() {
								t = nil
							} else {
								t = tv.Type
							}
						}
					}
					s.argTypes = append(s.argTypes, t)
				}
			}
			s.typePos = isTypePosition(e, parents)
			return s
		}
		skip := map[*ast.Ident]bool{f.Name: true}
		ast.Inspect(f, func(n ast.Node) bool {
			if ts, ok := n.(*ast.TypeSwitchStmt); ok {
				if as, ok := ts.Assign.(*ast.AssignStmt); ok {
					for _, l := range as.Lhs {
						if id, ok := l.(*ast.Ident); ok {
							skip[id] = true
						}
					}
				}
			}
			return true
		})
		ast.Inspect(f, func(n ast.Node) bool {
			switch n := n.(type) {
			case *ast.SelectorExpr:
				if _, ok := info.Selections[n]; ok {
					return true
				}
				if _, ok := info.Uses[n.Sel]; ok {
					return true
				}
				if id, ok := n.X.(*ast.Ident); ok {
					if _, isPkg := info.Uses[id].(*types.PkgName); isPkg {
						warn(fmt.Sprintf("%s: unresolved %s.%s (another package; fix by hand)", fset.Position(n.Pos()), id.Name, n.Sel.Name))
						return false
					}
				}
				tv, ok := info.Types[n.X]
				if !ok || tv.Type == nil {
					return true
				}
				t := tv.Type
				if p, ok := t.(*types.Pointer); ok {
					t = p.Elem()
				}
				named, ok := t.(*types.Named)
				if !ok {
					return true
				}
				obj := named.Obj()
				if obj.Pkg() == nil || obj.Pkg() != res.pkg {
					return true
				}
				if pendingTypes[obj.Name()] {
					add(obj.Name(), n.Sel.Name, siteOf(n))
				} else if obj.Parent() == res.pkg.Scope() {
					warn(fmt.Sprintf("%s: %s has no member %s (ported type; fix the call or the port)", fset.Position(n.Pos()), obj.Name(), n.Sel.Name))
				}
				return true
			case *ast.Ident:
				if n.Name == "_" || skip[n] || info.Uses[n] != nil || info.Defs[n] != nil {
					return true
				}
				p := parents[n]
				if sel, ok := p.(*ast.SelectorExpr); ok && sel.Sel == n {
					return true
				}
				if kv, ok := p.(*ast.KeyValueExpr); ok && kv.Key == n {
					if lit, inLit := parents[kv].(*ast.CompositeLit); inLit {
						tv, ok := info.Types[lit]
						if !ok || tv.Type == nil {
							return true
						}
						if _, isStruct := tv.Type.Underlying().(*types.Struct); isStruct {
							return true
						}
					}
				}
				if _, ok := p.(*ast.LabeledStmt); ok {
					return true
				}
				if br, ok := p.(*ast.BranchStmt); ok && br.Label == n {
					return true
				}
				if universe := types.Universe.Lookup(n.Name); universe != nil {
					return true
				}
				add("", n.Name, siteOf(n))
			}
			return true
		})
	}
	var out []*need
	for _, n := range needs {
		out = append(out, n)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].key() < out[j].key() })
	return out
}

// isTypePosition reports whether e is used where Go expects a type.
func isTypePosition(e ast.Expr, parents map[ast.Node]ast.Node) bool {
	var child ast.Node = e
	for p := parents[child]; p != nil; child, p = p, parents[p] {
		switch p := p.(type) {
		case *ast.StarExpr, *ast.ArrayType, *ast.MapType, *ast.ChanType, *ast.Ellipsis:
			continue
		case *ast.Field:
			return p.Type == child
		case *ast.ValueSpec:
			return p.Type == child
		case *ast.TypeSpec:
			return p.Type == child
		case *ast.CompositeLit:
			return p.Type == child
		case *ast.TypeAssertExpr:
			return p.Type == child
		case *ast.CaseClause:
			if blk, ok := parents[p].(*ast.BlockStmt); ok {
				_, isTypeSwitch := parents[blk].(*ast.TypeSwitchStmt)
				return isTypeSwitch
			}
			return false
		case *ast.IndexExpr, *ast.IndexListExpr:
			return true
		case *ast.FuncType, *ast.FieldList, *ast.InterfaceType, *ast.StructType:
			return true
		case *ast.CallExpr:
			// new(T) / make(T)
			if id, ok := p.Fun.(*ast.Ident); ok && (id.Name == "new" || id.Name == "make") && len(p.Args) > 0 && p.Args[0] == child {
				return true
			}
			return false
		default:
			return false
		}
	}
	return false
}
