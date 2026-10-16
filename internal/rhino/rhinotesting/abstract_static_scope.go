// Port of com/google/javascript/rhino/testing/AbstractStaticScope.java @ closure-compiler v20260415

package rhinotesting

import "github.com/paullepoulpe/jsgo/internal/rhino"

// AbstractStaticScope is a basic implementation of StaticScope for testing.
type AbstractStaticScope interface {
	rhino.StaticTypedScope
}

// AbstractStaticScopeBase holds the methods shared by AbstractStaticScope
// implementations; the implementation provides GetSlot.
type AbstractStaticScopeBase struct {
	self AbstractStaticScope
}

// InitAbstractStaticScopeBase sets the most-derived receiver.
func (s *AbstractStaticScopeBase) InitAbstractStaticScopeBase(self AbstractStaticScope) {
	s.self = self
}

// Java: AbstractStaticScope.getRootNode()
func (s *AbstractStaticScopeBase) GetRootNode() *rhino.Node {
	return nil
}

// Java: AbstractStaticScope.getParentScope()
func (s *AbstractStaticScopeBase) GetParentScope() rhino.StaticScope {
	return nil
}

// Java: AbstractStaticScope.getOwnSlot(String)
func (s *AbstractStaticScopeBase) GetOwnSlot(name string) rhino.StaticSlot {
	return s.self.GetSlot(name)
}

// Java: AbstractStaticScope.getTypeOfThis()
func (s *AbstractStaticScopeBase) GetTypeOfThis() rhino.JSType {
	return nil
}

// Java: StaticScope.getTopmostScopeOfEventualDeclaration(String) (default method)
func (s *AbstractStaticScopeBase) GetTopmostScopeOfEventualDeclaration(name string) rhino.StaticScope {
	return rhino.StaticScope_GetTopmostScopeOfEventualDeclaration(s.self, name)
}

// Java: StaticTypedScope.lookupQualifiedName(QualifiedName) (default method)
func (s *AbstractStaticScopeBase) LookupQualifiedName(qname rhino.QualifiedName) rhino.JSType {
	return rhino.StaticTypedScope_LookupQualifiedName(s.self, qname)
}
