package jscomp

import (
	"os"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Var's equality is final ScopedName equality: immutable name plus declaration
// scope-root identity. The key type here is *Var, so an embedding Go type cannot
// introduce an overriding Equals implementation into this map. Distinct Var or
// Scope instances with the same name/root still address the same collection.
type referenceCollectorScopedKey struct {
	root *rhino.Node
	name string
}

func referenceCollectorKey(v *Var) referenceCollectorScopedKey {
	if v == nil {
		// The original object map accepts a null key. A non-null Var with a
		// null scope root, by contrast, throws from its HashCode method.
		return referenceCollectorScopedKey{}
	}
	root := v.GetScopeRoot()
	if root == nil {
		panic(jutil.NewNullPointerException())
	}
	return referenceCollectorScopedKey{root: root, name: v.GetName()}
}

func (c *ReferenceCollector) initNativeReferenceLookup() {
	if os.Getenv("JSGO_EXPERIMENT_NATIVE_REFERENCES") == "1" {
		c.nativeReferences = make(map[referenceCollectorScopedKey]*ReferenceCollection)
	}
}

func (c *ReferenceCollector) referenceMapWrapper() *ReferenceCollector_ReferenceMapWrapper {
	wrapper := NewReferenceCollector_ReferenceMapWrapper(c.referenceMap)
	if c.nativeReferences != nil {
		wrapper.owner = c
	}
	return wrapper
}
