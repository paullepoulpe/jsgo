// Port of com/google/javascript/jscomp/parsing/JsDocInfoParser.java @ closure-compiler v20260415
//
// Tier-0 and Tier-1 tags: `@type {T}`, `@param {T} name`, `@return {T}`.
// The Java parser is a proper tokeniser; here we use regex over the
// comment body, which is sufficient while type expressions stay simple
// (primitives, `!` non-null prefix). A real grammar port is queued for
// when Tier-2 templates/unions force it.

package rhino

import "regexp"

// Captures everything inside `{...}` so the registry can handle the
// full type expression (primitives, `!` non-null prefix, `?` nullable,
// `|` unions).
var typeExprRe = `([^}]+?)`

var (
	typeTagRe        = regexp.MustCompile(`@type\s*\{\s*` + typeExprRe + `\s*\}`)
	paramTagRe       = regexp.MustCompile(`@param\s*\{\s*` + typeExprRe + `\s*\}\s*([A-Za-z_][A-Za-z0-9_]*)`)
	returnTagRe      = regexp.MustCompile(`@return\s*\{\s*` + typeExprRe + `\s*\}`)
	extendsTagRe     = regexp.MustCompile(`@extends\s*\{\s*` + typeExprRe + `\s*\}`)
	implementsTagRe  = regexp.MustCompile(`@implements\s*\{\s*` + typeExprRe + `\s*\}`)
	constructorTagRe = regexp.MustCompile(`@constructor\b`)
	interfaceTagRe   = regexp.MustCompile(`@interface\b`)
	recordTagRe      = regexp.MustCompile(`@record\b`)
)

// ParseJSDoc returns a JSDocInfo for the given block-comment body (the
// text between /** and */, without the delimiters), or nil if none of
// the recognised tags appear.
func ParseJSDoc(body string) *JSDocInfo {
	info := &JSDocInfo{}
	any := false

	if m := typeTagRe.FindStringSubmatch(body); m != nil {
		info.SetType(NewJSTypeExpression(m[1]))
		any = true
	}
	for _, m := range paramTagRe.FindAllStringSubmatch(body, -1) {
		info.AddParam(m[2], NewJSTypeExpression(m[1]))
		any = true
	}
	if m := returnTagRe.FindStringSubmatch(body); m != nil {
		info.SetReturnType(NewJSTypeExpression(m[1]))
		any = true
	}
	if m := extendsTagRe.FindStringSubmatch(body); m != nil {
		info.SetBaseType(NewJSTypeExpression(m[1]))
		any = true
	}
	if constructorTagRe.MatchString(body) {
		info.SetConstructor(true)
		any = true
	}
	if interfaceTagRe.MatchString(body) {
		info.SetInterface(true)
		any = true
	}
	if recordTagRe.MatchString(body) {
		info.SetRecord(true)
		any = true
	}
	for _, m := range implementsTagRe.FindAllStringSubmatch(body, -1) {
		info.AddImplementation(NewJSTypeExpression(m[1]))
		any = true
	}

	if !any {
		return nil
	}
	return info
}
