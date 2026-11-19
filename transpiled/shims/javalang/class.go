// Clean-room specification based on java.lang.Class Javadoc.

package javalang

// Class is a lightweight reflection-free class descriptor used as a token
// and map key by translated Java code.
type Class struct {
	name       string
	simpleName string
}

// NewClass creates a new Class descriptor.
func NewClass(name, simpleName string) *Class {
	return &Class{name: name, simpleName: simpleName}
}

// GetName returns the fully qualified name of the class.
func (c *Class) GetName() string {
	if c == nil {
		return ""
	}
	return c.name
}

// GetSimpleName returns the simple unqualified name of the underlying class.
func (c *Class) GetSimpleName() string {
	if c == nil {
		return ""
	}
	return c.simpleName
}

// String returns "class <name>".
func (c *Class) String() string {
	return "class " + c.GetName()
}

// Equals reports whether other is a Class descriptor with the same name.
func (c *Class) Equals(other any) bool {
	o, ok := other.(*Class)
	return ok && (c == o || (c != nil && o != nil && c.name == o.name))
}

// HashCode returns the polynomial hash code of the class name.
func (c *Class) HashCode() int32 {
	if c == nil {
		return 0
	}
	return StringHashCode(c.name)
}
