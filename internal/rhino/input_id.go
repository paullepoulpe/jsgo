// Port of com/google/javascript/rhino/InputId.java @ closure-compiler v20260415

package rhino

import "github.com/paullepoulpe/jsgo/internal/jutil"

// InputId is an id used uniquely identify a CompilerInput.
type InputId struct {
	id string
}

// Java: InputId.serialVersionUID
const InputId_serialVersionUID int64 = 1

// Java: InputId(String)
func NewInputId(id string) *InputId {
	return &InputId{id: id}
}

// Java: InputId.getIdName()
func (i *InputId) GetIdName() string {
	return i.id
}

// Java: InputId.hashCode()
func (i *InputId) HashCode() int32 {
	return jutil.StringHashCode(i.id)
}

// Java: InputId.equals(Object)
func (i *InputId) Equals(obj any) bool {
	if other, ok := obj.(*InputId); ok && other == i {
		return true
	}
	inputId, ok := obj.(*InputId)
	if !ok || inputId == nil {
		return false
	}
	return i.id == inputId.id
}

// Java: InputId.toString()
func (i *InputId) String() string {
	return "InputId: " + i.GetIdName()
}
