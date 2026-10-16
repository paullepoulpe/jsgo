// Port of com/google/javascript/rhino/jstype/EquivalenceMethod.java @ closure-compiler v20260415

package rhino

// EquivalenceMethod represents different ways for comparing equality among
// types.
type EquivalenceMethod int32

const (
	// EquivalenceMethod_IDENTITY indicates that the two types should behave
	// exactly the same under all type operations.
	//
	// Thus, {string} != {?} and {Unresolved} != {?}
	EquivalenceMethod_IDENTITY EquivalenceMethod = iota
	// EquivalenceMethod_DATA_FLOW indicates that the two types are almost
	// exactly the same, and that a data flow analysis algorithm comparing them
	// should consider them equal.
	//
	// In traditional type inference, the types form a finite lattice, and
	// this ensures that type inference will terminate.
	//
	// In our type system, the unknown types do not obey the lattice rules. So
	// if we continue to perform inference over the unknown types, we may never
	// terminate.
	//
	// By treating all unknown types as equivalent for the purposes of data
	// flow analysis, we ensure that the algorithm will terminate.
	//
	// Thus, {string} != {?} and {Unresolved} ~= {?}
	EquivalenceMethod_DATA_FLOW
)

var equivalenceMethod_names = [...]string{"IDENTITY", "DATA_FLOW"}

// EquivalenceMethod_values is EquivalenceMethod.values().
func EquivalenceMethod_values() []EquivalenceMethod {
	return []EquivalenceMethod{EquivalenceMethod_IDENTITY, EquivalenceMethod_DATA_FLOW}
}

// String is Enum.name().
func (m EquivalenceMethod) String() string { return equivalenceMethod_names[m] }
