// Port of com/google/javascript/rhino/jstype/AllType.java @ closure-compiler v20260415

package jstype

// allType represents the "All" type (*), the top of the JS type lattice.
// Every value is a member of All. It is only a subtype of itself and of
// Unknown (which accepts everything).
type allType struct{}

func (allType) TypeName() string            { return "*" }
func (allType) IsSubtypeOf(o JSType) bool   { _, ok := o.(allType); return ok || isUnknown(o) }
