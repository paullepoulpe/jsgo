// jsdoc_parse_hook.go provides a function variable that the parsing package
// can set (from its test init) to wire up ParseJSDocInfo for use in rhino
// package tests, avoiding the import cycle rhino→parsing.
//
// Production code never sets this variable; it is nil by default.

package rhino

// JSDocParserHook is optionally set by the parsing package's test init to
// provide a real implementation of JSDoc parsing to the rhino package tests.
// When nil, jsdocInfoFromBody returns nil (stub behaviour).
// The signature matches parsing.ParseJSDocInfo(body string, startLine int).
var JSDocParserHook func(body string, startLine int) *JSDocInfo
