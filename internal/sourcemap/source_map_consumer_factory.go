// Port of com/google/debugging/sourcemap/SourceMapConsumerFactory.java @ closure-compiler v20260415

package sourcemap

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// SourceMapConsumerFactory detects and parses the provided source map.
type SourceMapConsumerFactory struct{}

// SourceMapConsumerFactory_Parse parses the string representing the source
// map file contents. It panics with a *SourceMapParseException.
//
// Java: SourceMapConsumerFactory.parse(String)
func SourceMapConsumerFactory_Parse(contents string) SourceMapping {
	return SourceMapConsumerFactory_ParseWithSupplier(contents, nil)
}

// SourceMapConsumerFactory_ParseWithSupplier parses the string representing
// the source map file contents, using supplier (may be nil) for any
// referenced maps.
//
// Java: SourceMapConsumerFactory.parse(String, SourceMapSupplier)
func SourceMapConsumerFactory_ParseWithSupplier(contents string, supplier SourceMapSupplier) SourceMapping {
	// Version 1, starts with a magic string
	if strings.HasPrefix(contents, "/** Begin line maps. **/") {
		panic(NewSourceMapParseException(
			"This appears to be a V1 SourceMap, which is not supported."))
	} else if strings.HasPrefix(contents, "{") {
		sourceMapObject := SourceMapObjectParser_Parse(contents)

		// Check basic assertions about the format.
		switch sourceMapObject.GetVersion() {
		case 3:
			consumer := NewSourceMapConsumerV3()
			consumer.ParseSourceMapObject(sourceMapObject, supplier)
			return consumer
		default:
			panic(NewSourceMapParseException(
				"Unknown source map version:" + jutil.ToJavaString(sourceMapObject.GetVersion())))
		}
	}

	panic(NewSourceMapParseException("unable to detect source map format"))
}

// SourceMapConsumerFactory_ParseFast uses the fast parser to parse the
// source map. This only supports version 3 source maps.
//
// Java: SourceMapConsumerFactory.parseFast(String)
func SourceMapConsumerFactory_ParseFast(contents string) *SourceMapConsumerV3 {
	if strings.HasPrefix(contents, "{") {
		sourceMapObject := SourceMapObjectParser_ParseFast(contents)
		jutil.CheckArgumentMsg1(
			sourceMapObject.GetVersion() == 3,
			"Unknown source map version:%s",
			sourceMapObject.GetVersion())
		consumer := NewSourceMapConsumerV3()
		consumer.ParseSourceMapObject(sourceMapObject, nil)
		return consumer
	}
	panic(NewSourceMapParseException("unable to detect source map format"))
}
