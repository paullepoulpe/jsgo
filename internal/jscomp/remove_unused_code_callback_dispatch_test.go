package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/diagnostic"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"testing"
)

type callbackDispatchRemovable struct {
	RemoveUnusedCode_RemovableBase
	calls *[]string
}

func (r *callbackDispatchRemovable) Remove(AbstractCompiler) {
	*r.calls = append(*r.calls, "remove:"+r.propertyName)
}
func (r *callbackDispatchRemovable) ApplyContinuations() {
	*r.calls = append(*r.calls, "apply:"+r.propertyName)
}
func (r *callbackDispatchRemovable) RemoveInternal(AbstractCompiler) {
	panic("custom Remove override was bypassed")
}

func callbackDispatchOwner() *RemoveUnusedCode {
	c := NewCompiler()
	c.InitOptions(NewCompilerOptions())
	r := NewRemoveUnusedCode(NewRemoveUnusedCode_Builder(c))
	r.removeUnusedPrototypeProperties = true
	r.removalLog = diagnostic.LogFile_CreateNoOp()
	r.pinnedPropertyNames = make(map[string]struct{})
	r.propRemovables = make(map[string]propRemovableList)
	return r
}
func callbackDispatchProperty(r *RemoveUnusedCode, name string, calls *[]string) *callbackDispatchRemovable {
	item := &callbackDispatchRemovable{calls: calls}
	builder := &RemoveUnusedCode_RemovableBuilder{propertyName: name, hasPropertyName: true, isPrototypeDotPropertyReference: true}
	item.InitRemoveUnusedCode_RemovableBase(item, r, rhino.IR_Name("target"), builder)
	r.ConsiderForIndependentRemoval(item)
	return item
}
func TestRemoveUnusedCodeCustomRemovableDispatch(t *testing.T) {
	r := callbackDispatchOwner()
	var calls []string
	callbackDispatchProperty(r, "unusedProperty", &calls)
	r.RemoveIndependentlyRemovableProperties()
	if len(calls) != 1 || calls[0] != "remove:unusedProperty" {
		t.Fatalf("Remove override lost: %v", calls)
	}
	calls = nil
	callbackDispatchProperty(r, "referencedProperty", &calls)
	r.MarkPropertyNameAsPinned("referencedProperty")
	if len(calls) != 1 || calls[0] != "apply:referencedProperty" {
		t.Fatalf("continuation override lost: %v", calls)
	}
}
func TestRemoveUnusedCodeImplicitPropertySnapshot(t *testing.T) {
	original := RemoveUnusedCode_IMPLICITLY_USED_PROPERTIES
	defer func() { RemoveUnusedCode_IMPLICITLY_USED_PROPERTIES = original }()
	RemoveUnusedCode_IMPLICITLY_USED_PROPERTIES = []string{"customImplicit"}
	custom := callbackDispatchOwner()
	RemoveUnusedCode_IMPLICITLY_USED_PROPERTIES = original
	stock := callbackDispatchOwner()
	if !custom.isPropertyPinned("customImplicit") || custom.isPropertyPinned("length") {
		t.Fatal("constructor did not preserve the replacement implicit-property list")
	}
	if !stock.isPropertyPinned("length") || stock.isPropertyPinned("customImplicit") {
		t.Fatal("constructor snapshot leaked into another pass")
	}
}
