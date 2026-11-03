package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/diagnostic"
)

// Go-only: test-helper — Invalidation exists only to be logged with GSON
// (DisambiguateProperties' renaming_index log). The expected JSON is what
// the v20261015 jar's LogsGson prints for the same map (black-box, through
// a probe class in the disambiguate package).
func TestInvalidation_logsGsonMatchesJava(t *testing.T) {
	m := NewOrderedStringMap[disambiguateRenamingEntry]()
	m.Put("a", disambiguateRenamingEntry{invalidation: Invalidation_WellKnownProperty(), invalidated: true})
	m.Put("b", disambiguateRenamingEntry{invalidation: Invalidation_InvalidatingType(5), invalidated: true})
	m.Put("c", disambiguateRenamingEntry{invalidation: Invalidation_UndeclaredAccess(7), invalidated: true})
	const want = `{"a":{"reason":"WELL_KNOWN_PROPERTY"},"b":{"receiverType":5,"reason":"INVALIDATING_TYPE"},"c":{"receiverType":7,"reason":"UNDECLARED_ACCESS"}}`
	if got := diagnostic.LogsGson_ToJson(disambiguateProperties_renamingIndexJson{m: m}); got != want {
		t.Errorf("got  %s\nwant %s", got, want)
	}
	mZero := NewOrderedStringMap[disambiguateRenamingEntry]()
	mZero.Put("z", disambiguateRenamingEntry{invalidated: true})
	if got := diagnostic.LogsGson_ToJson(disambiguateProperties_renamingIndexJson{m: mZero}); got != `{"z":null}` {
		t.Errorf("zero invalidation: got %s, want {\"z\":null}", got)
	}
}
