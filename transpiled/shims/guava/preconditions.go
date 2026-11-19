// Clean-room specification based on Google Guava Preconditions (Apache 2.0).

package guava

import (
	"fmt"
	"strings"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// CheckState ensures the truth of an expression involving the state of the calling instance.
func CheckState(expression bool) {
	if !expression {
		panic(javalang.NewIllegalStateException(""))
	}
}

// CheckArgument ensures the truth of an expression involving one or more arguments to the calling method.
func CheckArgument(expression bool) {
	if !expression {
		panic(javalang.NewIllegalArgumentException(""))
	}
}

// Fixed-arity fast paths: zero allocation on the passing path.

func CheckStateMsg1[A any](expression bool, template string, p1 A) {
	if !expression {
		panic(javalang.NewIllegalStateException(LenientFormat(template, p1)))
	}
}

func CheckStateMsg2[A, B any](expression bool, template string, p1 A, p2 B) {
	if !expression {
		panic(javalang.NewIllegalStateException(LenientFormat(template, p1, p2)))
	}
}

func CheckStateMsg3[A, B, C any](expression bool, template string, p1 A, p2 B, p3 C) {
	if !expression {
		panic(javalang.NewIllegalStateException(LenientFormat(template, p1, p2, p3)))
	}
}

func CheckStateMsg(expression bool, template string, args ...any) {
	if !expression {
		panic(javalang.NewIllegalStateException(LenientFormat(template, args...)))
	}
}

func CheckArgumentMsg1[A any](expression bool, template string, p1 A) {
	if !expression {
		panic(javalang.NewIllegalArgumentException(LenientFormat(template, p1)))
	}
}

func CheckArgumentMsg2[A, B any](expression bool, template string, p1 A, p2 B) {
	if !expression {
		panic(javalang.NewIllegalArgumentException(LenientFormat(template, p1, p2)))
	}
}

func CheckArgumentMsg3[A, B, C any](expression bool, template string, p1 A, p2 B, p3 C) {
	if !expression {
		panic(javalang.NewIllegalArgumentException(LenientFormat(template, p1, p2, p3)))
	}
}

func CheckArgumentMsg(expression bool, template string, args ...any) {
	if !expression {
		panic(javalang.NewIllegalArgumentException(LenientFormat(template, args...)))
	}
}

// CheckNotNull ensures that an object reference passed as a parameter to the calling method is not null.
func CheckNotNull[T any](reference T) T {
	if javalang.IsNil(reference) {
		panic(javalang.NewNullPointerException())
	}
	return reference
}

func CheckNotNullMsg1[T, A any](reference T, template string, p1 A) T {
	if javalang.IsNil(reference) {
		panic(javalang.NewNullPointerExceptionMsg(LenientFormat(template, p1)))
	}
	return reference
}

func CheckNotNullMsg2[T, A, B any](reference T, template string, p1 A, p2 B) T {
	if javalang.IsNil(reference) {
		panic(javalang.NewNullPointerExceptionMsg(LenientFormat(template, p1, p2)))
	}
	return reference
}

func CheckNotNullMsg[T any](reference T, template string, args ...any) T {
	if javalang.IsNil(reference) {
		panic(javalang.NewNullPointerExceptionMsg(LenientFormat(template, args...)))
	}
	return reference
}

// CheckElementIndex ensures that index specifies a valid element in an array/list of size size.
func CheckElementIndex(index, size int) int {
	if index < 0 || index >= size {
		panic(javalang.NewIndexOutOfBoundsException(fmt.Sprintf("index (%d) must be less than size (%d)", index, size)))
	}
	return index
}

// CheckPositionIndex ensures that index specifies a valid position in an array/list of size size.
func CheckPositionIndex(index, size int) int {
	if index < 0 || index > size {
		panic(javalang.NewIndexOutOfBoundsException(fmt.Sprintf("index (%d) must not be greater than size (%d)", index, size)))
	}
	return index
}

// CheckPositionIndexes ensures that [start, end) specifies a valid sub-range in an array/list of size size.
func CheckPositionIndexes(start, end, size int) {
	if start < 0 || end < start || end > size {
		panic(javalang.NewIndexOutOfBoundsException(fmt.Sprintf("invalid range [%d, %d) for size %d", start, end, size)))
	}
}

// LenientFormat substitutes %s in template with the string representations of args.
func LenientFormat(template string, args ...any) string {
	if len(args) == 0 {
		return template
	}
	var sb strings.Builder
	templateStart := 0
	i := 0
	for i < len(args) {
		placeholderStart := strings.Index(template[templateStart:], "%s")
		if placeholderStart == -1 {
			break
		}
		placeholderStart += templateStart
		sb.WriteString(template[templateStart:placeholderStart])
		sb.WriteString(javalang.ToJavaString(args[i]))
		templateStart = placeholderStart + 2
		i++
	}
	sb.WriteString(template[templateStart:])
	if i < len(args) {
		sb.WriteString(" [")
		sb.WriteString(javalang.ToJavaString(args[i]))
		for i++; i < len(args); i++ {
			sb.WriteString(", ")
			sb.WriteString(javalang.ToJavaString(args[i]))
		}
		sb.WriteString("]")
	}
	return sb.String()
}
