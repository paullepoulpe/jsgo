package transpile_test

import (
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/transpile"
)

func TestClosureBundlerWithInputMinifier(t *testing.T) {
	var seenPaths []string
	bundler := transpile.NewClosureBundler().
		WithPath("pkg/mod.js").
		WithInputMinifier(func(path, code string) string {
			seenPaths = append(seenPaths, path)
			return "/*min*/" + strings.TrimSpace(code)
		})

	sb := jutil.NewStringWriter()
	info := jscomp.SimpleDependencyInfo_Builder_("pkg/mod.js", "pkg/mod.js").
		SetGoogModule(true).
		Build()
	bundler.AppendTo(sb, info, "exports.x = 1 + 2;\n")

	got := sb.String()
	if len(seenPaths) != 1 || seenPaths[0] != "pkg/mod.js" {
		t.Fatalf("seenPaths = %v, want [pkg/mod.js]", seenPaths)
	}
	if !strings.Contains(got, "/*min*/exports.x = 1 + 2;") {
		t.Fatalf("unexpected bundled output: %q", got)
	}
}
