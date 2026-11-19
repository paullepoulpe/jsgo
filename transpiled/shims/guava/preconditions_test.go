package guava

import (
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

func TestPreconditions(t *testing.T) {
	// CheckState passing
	CheckState(true)
	CheckStateMsg1(true, "err %s", 1)
	CheckStateMsg2(true, "err %s %s", 1, 2)
	CheckStateMsg3(true, "err %s %s %s", 1, 2, 3)
	CheckStateMsg(true, "err %s", 1, 2, 3, 4)

	// CheckState failing
	var stateCaught javalang.Throwable
	func() {
		defer func() { stateCaught = javalang.RecoverException(recover()) }()
		CheckStateMsg1(false, "state failure with arg %s", "XYZ")
	}()
	if stateCaught == nil || stateCaught.GetMessage() != "state failure with arg XYZ" {
		t.Errorf("CheckStateMsg1 failure message got %q", stateCaught.GetMessage())
	}

	// CheckArgument failing
	var argCaught javalang.Throwable
	func() {
		defer func() { argCaught = javalang.RecoverException(recover()) }()
		CheckArgumentMsg2(false, "expected %s but got %s", "A", "B")
	}()
	if argCaught == nil || argCaught.GetMessage() != "expected A but got B" {
		t.Errorf("CheckArgumentMsg2 failure message got %q", argCaught.GetMessage())
	}

	// CheckNotNull passing and failing
	val := CheckNotNull("hello")
	if val != "hello" {
		t.Errorf("CheckNotNull returned %q, want \"hello\"", val)
	}

	var nullCaught javalang.Throwable
	func() {
		defer func() { nullCaught = javalang.RecoverException(recover()) }()
		var nilPtr *string
		CheckNotNullMsg1(nilPtr, "variable %s must not be null", "foo")
	}()
	if nullCaught == nil || nullCaught.GetMessage() != "variable foo must not be null" {
		t.Errorf("CheckNotNullMsg1 failure message got %q", nullCaught.GetMessage())
	}

	// Boundary checks
	if CheckElementIndex(2, 5) != 2 {
		t.Errorf("CheckElementIndex failed")
	}
	if CheckPositionIndex(5, 5) != 5 {
		t.Errorf("CheckPositionIndex failed")
	}
	CheckPositionIndexes(1, 3, 5)
}

func TestLenientFormat(t *testing.T) {
	res1 := LenientFormat("foo %s bar %s", 1, 2)
	if res1 != "foo 1 bar 2" {
		t.Errorf("LenientFormat got %q, want \"foo 1 bar 2\"", res1)
	}

	// Extra args appended in brackets
	res2 := LenientFormat("value: %s", "A", "B", "C")
	if res2 != "value: A [B, C]" {
		t.Errorf("LenientFormat with extra args got %q, want \"value: A [B, C]\"", res2)
	}
}
