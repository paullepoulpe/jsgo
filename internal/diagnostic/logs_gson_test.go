// Port of test/com/google/javascript/jscomp/diagnostic/LogsGsonTest.java @ closure-compiler v20260415

package diagnostic

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func TestAdapter_write_logsGsonable(t *testing.T) {
	assertLogsGson(t, []any{&TestLogsGsonable{}}, "[[0]]")
}

type TestLogsGsonable struct{}

func (*TestLogsGsonable) ToLogsGson() any {
	return []int32{0} // Java: ImmutableList.of(0)
}

func TestAdapter_write_multimap(t *testing.T) {
	// Given
	// Java: ImmutableMultimap.builder().put("foo", 0).put("foo", 1).put("bar", 2).build();
	// Go-only: language-adapter — jutil models Guava multimaps with
	// SetMultimap (no duplicate values here, so it is the same multimap).
	value := jutil.NewSetMultimap[any, any]()
	value.Put("foo", int32(0))
	value.Put("foo", int32(1))
	value.Put("bar", int32(2))

	// Then
	assertLogsGson(t, value, "{\"foo\":[0,1],\"bar\":[2]}")
}

func assertLogsGson(t *testing.T, value any, want string) {
	t.Helper()
	if got := LogsGson_ToJson(value); got != want {
		t.Errorf("LogsGson.toJson = %s, want %s", got, want)
	}
}
