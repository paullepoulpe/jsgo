// Port of test/com/google/javascript/jscomp/deps/DependencyInfoTest.java @ closure-compiler v20260415

package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func TestDependencyInfo_WriteAddDependency(t *testing.T) {
	sb := jutil.NewStringBuilder()
	flags := jutil.NewLinkedHashMap[string, string]()
	flags.Put("module", "goog")
	flags.Put("lang", "es6")
	DependencyInfo_Util_WriteAddDependency(
		sb,
		SimpleDependencyInfo_Builder_("../some/relative/path.js", "/unused/absolute/path.js").
			SetProvides([]string{"provided.symbol", "other.provide"}).
			SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("required.symbol"), DependencyInfo_Require_GoogRequireSymbol("other.require")).
			SetLoadFlags(flags).
			Build())
	jsaEq(t, sb.String(),
		"goog.addDependency('../some/relative/path.js', ['provided.symbol', 'other.provide'], "+
			"['required.symbol', 'other.require'], {'module': 'goog', 'lang': 'es6'});\n")
}

func TestDependencyInfo_WriteAddDependency_emptyArguments(t *testing.T) {
	sb := jutil.NewStringBuilder()
	DependencyInfo_Util_WriteAddDependency(sb, SimpleDependencyInfo_Builder_("path.js", "unused.js").Build())
	jsaEq(t, sb.String(), "goog.addDependency('path.js', [], []);\n")
}
