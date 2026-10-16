// Port of com/google/javascript/rhino/InputId.java
// @ closure-compiler v20260415

package rhino

import "fmt"

// InputId uniquely identifies a CompilerInput.
// Mirrors com.google.javascript.rhino.InputId.
type InputId struct {
	id string
}

// NewInputId constructs an InputId with the given identifier string.
func NewInputId(id string) *InputId { return &InputId{id: id} }

// GetIdName returns the identifier string.
func (i *InputId) GetIdName() string { return i.id }

// Equals reports whether i equals other, mirroring InputId.equals in Java.
func (i *InputId) Equals(other *InputId) bool {
	if i == other {
		return true
	}
	if other == nil {
		return false
	}
	return i.id == other.id
}

func (i *InputId) String() string { return fmt.Sprintf("InputId: %s", i.id) }
