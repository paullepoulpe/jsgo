// Clean-room unused-import pruning helpers for Go AST passes.

package passes

import (
	"go/ast"
	"go/token"
	"go/types"
	"path"
	"strconv"
)

// pruneUnusedImports removes non-blank, non-dot imports from file that have zero package
// selector references in file.Decls (or whose positions appear in extraUnusedPos).
func pruneUnusedImports(file *ast.File, imp types.Importer, extraUnusedPos map[token.Pos]bool) bool {
	usedQualifiers := collectSelectorQualifiers(file)
	changed := false

	var newDecls []ast.Decl
	var newImports []*ast.ImportSpec

	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.IMPORT {
			newDecls = append(newDecls, decl)
			continue
		}
		var keptSpecs []ast.Spec
		for _, s := range gd.Specs {
			spec, ok := s.(*ast.ImportSpec)
			if !ok {
				keptSpecs = append(keptSpecs, s)
				continue
			}
			if spec.Name != nil && (spec.Name.Name == "_" || spec.Name.Name == ".") {
				keptSpecs = append(keptSpecs, spec)
				newImports = append(newImports, spec)
				continue
			}
			if extraUnusedPos != nil && (extraUnusedPos[spec.Pos()] || extraUnusedPos[spec.Path.Pos()]) {
				changed = true
				continue
			}
			localName := specLocalName(spec, imp)
			if localName != "" && !usedQualifiers[localName] {
				changed = true
				continue
			}
			keptSpecs = append(keptSpecs, spec)
			newImports = append(newImports, spec)
		}
		if len(keptSpecs) > 0 {
			gd.Specs = keptSpecs
			newDecls = append(newDecls, gd)
		} else {
			changed = true
		}
	}

	if changed {
		file.Decls = newDecls
		file.Imports = newImports
	}
	return changed
}

func specLocalName(spec *ast.ImportSpec, imp types.Importer) string {
	if spec.Name != nil {
		return spec.Name.Name
	}
	p, err := strconv.Unquote(spec.Path.Value)
	if err != nil {
		return ""
	}
	if imp != nil {
		if pkg, err := imp.Import(p); err == nil && pkg != nil && pkg.Name() != "" {
			return pkg.Name()
		}
	}
	return path.Base(p)
}

func collectSelectorQualifiers(file *ast.File) map[string]bool {
	qualifiers := make(map[string]bool)
	for _, decl := range file.Decls {
		if gd, ok := decl.(*ast.GenDecl); ok && gd.Tok == token.IMPORT {
			continue
		}
		ast.Inspect(decl, func(n ast.Node) bool {
			if sel, ok := n.(*ast.SelectorExpr); ok {
				if id, ok := sel.X.(*ast.Ident); ok {
					qualifiers[id.Name] = true
				}
			}
			return true
		})
	}
	return qualifiers
}
