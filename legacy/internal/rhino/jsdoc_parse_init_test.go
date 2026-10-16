// jsdoc_parse_init_test.go wires JSDocParserHook to parsing.ParseJSDocInfo
// so that rhino package tests can call the real parser without import cycles.
//
// This is in package rhino_test (external test package), which CAN import
// both rhino and parsing since it is a separate package. The init() here
// runs in the same test binary as package rhino tests, so those tests see
// JSDocParserHook set to the real implementation.

package rhino_test

import (
	"github.com/paullepoulpe/jsgo/legacy/internal/jscomp/parsing"
	"github.com/paullepoulpe/jsgo/legacy/internal/rhino"
)

func init() {
	rhino.JSDocParserHook = func(body string, startLine int) *rhino.JSDocInfo {
		return parsing.ParseJSDocInfo(body, startLine)
	}
}
