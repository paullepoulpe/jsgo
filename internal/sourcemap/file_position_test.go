package sourcemap_test

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/sourcemap"
)

// FilePosition.toString is Guava's toStringHelper form.
func TestFilePositionToString(t *testing.T) {
	if got := sourcemap.NewFilePosition(1, -2).String(); got != "FilePosition{line=1, column=-2}" {
		t.Errorf("toString = %q", got)
	}
}
