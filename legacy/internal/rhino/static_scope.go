// Port of com/google/javascript/rhino/StaticScope.java
// @ closure-compiler v20260415

package rhino

// StaticScope must be implemented by any object that defines variables for the
// purposes of static analysis. It is distinguished from a run-time scope.
// Mirrors com.google.javascript.rhino.StaticScope.
type StaticScope interface {
	// GetRootNode returns the root node associated with this scope. May be nil.
	GetRootNode() *Node

	// GetParentScope returns the scope enclosing this one, or nil if none.
	GetParentScope() StaticScope

	// GetSlot returns any defined slot within this scope for the given name.
	// The call continues searching through parent scopes if a slot with this
	// name is not found in the current scope.
	GetSlot(name string) StaticSlot

	// GetOwnSlot is like GetSlot but does not recurse into parent scopes.
	GetOwnSlot(name string) StaticSlot
}

// GetTopmostScopeOfEventualDeclaration mirrors
// StaticScope.getTopmostScopeOfEventualDeclaration: it returns the topmost
// scope containing the given name, or nil if no scopes do.
//
// This is intended for use while scopes are still being built. Once scope
// building is complete the result must equal GetSlot(name).GetScope() or nil.
func GetTopmostScopeOfEventualDeclaration(s StaticScope, name string) StaticScope {
	if s == nil {
		return nil
	}
	slot := s.GetOwnSlot(name)
	if slot != nil {
		return slot.GetScope()
	}
	return GetTopmostScopeOfEventualDeclaration(s.GetParentScope(), name)
}
