// Port of com/google/javascript/rhino/jstype/NullType.java @ closure-compiler v20260415

package jstype

// nullType is the JS null type (the singleton null value).  It is a subtype
// of itself, All, and Unknown.
type nullType struct{}

func (nullType) TypeName() string { return "null" }
func (nullType) IsSubtypeOf(o JSType) bool {
	if _, ok := o.(nullType); ok {
		return true
	}
	return isAll(o) || isUnknown(o)
}
