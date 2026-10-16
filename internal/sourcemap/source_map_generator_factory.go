// Port of com/google/debugging/sourcemap/SourceMapGeneratorFactory.java @ closure-compiler v20260415

package sourcemap

import "github.com/paullepoulpe/jsgo/internal/jutil"

// SourceMapGeneratorFactory creates source map generators.
type SourceMapGeneratorFactory struct{}

// SourceMapGeneratorFactory_GetInstance returns the appropriate source map
// object for the given source map format.
//
// Java: SourceMapGeneratorFactory.getInstance(SourceMapFormat)
func SourceMapGeneratorFactory_GetInstance(format SourceMapFormat) SourceMapGenerator {
	switch format {
	case SourceMapFormat_DEFAULT, SourceMapFormat_V3:
		return NewSourceMapGeneratorV3()
	}
	// Go-only: language-adapter — unreachable for a valid enum value: the
	// Java switch expression is exhaustive.
	panic(jutil.NewIllegalStateException())
}
