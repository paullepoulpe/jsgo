package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/diagnostic"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Go-only: test-helper — Invalidation exists only to be logged with GSON
// (DisambiguateProperties' renaming_index log). The expected JSON is what
// the v20261015 jar's LogsGson prints for the same map (black-box, through
// a probe class in the disambiguate package).
func TestInvalidation_logsGsonMatchesJava(t *testing.T) {
	m := jutil.NewLinkedHashMap[string, any]()
	m.Put("a", Invalidation_WellKnownProperty())
	m.Put("b", Invalidation_InvalidatingType(5))
	m.Put("c", Invalidation_UndeclaredAccess(7))
	const want = `{"a":{"reason":"WELL_KNOWN_PROPERTY"},"b":{"receiverType":5,"reason":"INVALIDATING_TYPE"},"c":{"receiverType":7,"reason":"UNDECLARED_ACCESS"}}`
	if got := diagnostic.LogsGson_ToJson(m); got != want {
		t.Errorf("got  %s\nwant %s", got, want)
	}
}
