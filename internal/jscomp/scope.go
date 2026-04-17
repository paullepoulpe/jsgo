// Port of com/google/javascript/jscomp/TypedScopeCreator.java @ closure-compiler v20260415
//
// TypedScopeCreator in Java walks the AST and builds a TypedScope tree,
// resolving JSDoc type annotations into JSType values via the
// JSTypeRegistry. Tier-0 needs only the script-level scope and
// declared-type lookup for a single `var` binding.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// TypedScope maps identifier names to their declared JSType.
type TypedScope struct {
	declared map[string]jstype.JSType
}

func (s *TypedScope) GetDeclaredType(name string) jstype.JSType {
	if s == nil {
		return nil
	}
	return s.declared[name]
}

type TypedScopeCreator struct {
	registry *jstype.Registry
}

func NewTypedScopeCreator(r *jstype.Registry) *TypedScopeCreator {
	return &TypedScopeCreator{registry: r}
}

// CreateScope mirrors TypedScopeCreator.createScope for the SCRIPT root.
// Walks top-level VAR declarations, reads attached JSDocInfo's @type, and
// records (name → declared JSType) in the returned scope.
func (c *TypedScopeCreator) CreateScope(script *rhino.Node) *TypedScope {
	scope := &TypedScope{declared: make(map[string]jstype.JSType)}
	for _, child := range script.Children() {
		if !isDeclToken(child.Token()) {
			continue
		}
		info := child.JSDocInfo()
		name := child.FirstChild()
		if name == nil || name.Token() != rhino.NAME {
			continue
		}
		if !info.HasType() {
			continue
		}
		scope.declared[name.String_()] = c.registry.ResolveName(info.Type().TypeName)
	}
	return scope
}

func isDeclToken(t rhino.Token) bool {
	return t == rhino.VAR || t == rhino.LET || t == rhino.CONST
}
