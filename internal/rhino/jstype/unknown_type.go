// Port of com/google/javascript/rhino/jstype/UnknownType.java @ closure-compiler v20260415

package jstype

// unknownType represents the "?" (unknown) type.  In the Java type lattice
// every type is a subtype of Unknown, so IsSubtypeOf always returns true.
type unknownType struct{}

func (unknownType) TypeName() string          { return "?" }
func (unknownType) IsSubtypeOf(o JSType) bool { return true }
