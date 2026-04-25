package jscomp

import (
	"fmt"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

func TestDebugStripNullUndef(t *testing.T) {
	// Simulate F.Instance and UnionType{F.Instance|null}
	fInst := jstype.NewObjectType("F")
	fInst.DefineProperty("z", jstype.Number)
	
	union := jstype.NewUnion(fInst, jstype.Null)
	fmt.Printf("union: %T\n", union)
	
	stripped := stripNullUndef(union)
	fmt.Printf("stripped: %T / %v\n", stripped, stripped)
	
	if obj, ok := stripped.(*jstype.ObjectType); ok {
		fmt.Printf("obj.GetPropertyType(z): %v\n", obj.GetPropertyType("z"))
		fmt.Printf("obj.IsPropertyInferred(z): %v\n", obj.IsPropertyInferred("z"))
	} else {
		fmt.Printf("stripped is not ObjectType: %T\n", stripped)
	}
}
