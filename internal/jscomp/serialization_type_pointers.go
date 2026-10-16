// Port of com/google/javascript/jscomp/serialization/TypePointers.java @ closure-compiler v20260415

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// TypePointers holds TypePointer utilities.
type TypePointers struct{}

// TypePointers_OFFSET_TO_AXIOMATIC_COLOR lists the axiomatic colors that the
// first N TypePointer pool offsets correspond to, in this order.
//
// These colors are never serialized because all their information is
// constant.
var TypePointers_OFFSET_TO_AXIOMATIC_COLOR = []*colors.Color{
	colors.StandardColors_UNKNOWN,
	colors.StandardColors_BOOLEAN,
	colors.StandardColors_STRING,
	colors.StandardColors_NUMBER,
	colors.StandardColors_NULL_OR_VOID,
	colors.StandardColors_SYMBOL,
	colors.StandardColors_BIGINT,
	colors.StandardColors_TOP_OBJECT,
}

// Java: TypePointers.AXIOMATIC_COLOR_COUNT
var TypePointers_AXIOMATIC_COLOR_COUNT = len(TypePointers_OFFSET_TO_AXIOMATIC_COLOR)

// Java: TypePointers.trimOffset(int)
func TypePointers_TrimOffset(x int) int {
	jutil.CheckStateMsg1(x >= TypePointers_AXIOMATIC_COLOR_COUNT, "%s", x)
	return x - TypePointers_AXIOMATIC_COLOR_COUNT
}

// Java: TypePointers.untrimOffset(int)
func TypePointers_UntrimOffset(x int) int {
	jutil.CheckStateMsg1(x >= 0, "%s", x)
	return x + TypePointers_AXIOMATIC_COLOR_COUNT
}

// Java: TypePointers.isAxiomatic(int)
func TypePointers_IsAxiomatic(x int) bool {
	return x < TypePointers_AXIOMATIC_COLOR_COUNT
}

// Java: TypePointers()
func newTypePointers() *TypePointers {
	panic(jutil.NewAssertionErrorNoMessage())
}
