// Port of com/google/javascript/jscomp/parsing/JsDocInfoParser.java @ closure-compiler v20260415
//
// Tier-0 and Tier-1 tags: `@type {T}`, `@param {T} name`, `@return {T}`.
// The Java parser is a proper tokeniser; here we use regex over the
// comment body, which is sufficient while type expressions stay simple
// (primitives, `!` non-null prefix). A real grammar port is queued for
// when Tier-2 templates/unions force it.

package rhino

import "regexp"

// Captures `!?` modifier + primitive/ident name inside `{...}`.
var typeExprRe = `([!?]?)\s*([A-Za-z_][A-Za-z0-9_]*)`

var (
	typeTagRe   = regexp.MustCompile(`@type\s*\{\s*` + typeExprRe + `\s*\}`)
	paramTagRe  = regexp.MustCompile(`@param\s*\{\s*` + typeExprRe + `\s*\}\s*([A-Za-z_][A-Za-z0-9_]*)`)
	returnTagRe = regexp.MustCompile(`@return\s*\{\s*` + typeExprRe + `\s*\}`)
)

// ParseJSDoc returns a JSDocInfo for the given block-comment body (the
// text between /** and */, without the delimiters), or nil if none of
// the recognised tags appear.
func ParseJSDoc(body string) *JSDocInfo {
	info := &JSDocInfo{}
	any := false

	if m := typeTagRe.FindStringSubmatch(body); m != nil {
		// Nullable `?Foo` is out of scope for Tier-0/1; skip to stay
		// conservative.
		if m[1] != "?" {
			info.SetType(NewJSTypeExpression(m[2]))
			any = true
		}
	}
	for _, m := range paramTagRe.FindAllStringSubmatch(body, -1) {
		if m[1] == "?" {
			// Nullable param type → leave unrecorded (conservative).
			continue
		}
		info.AddParam(m[3], NewJSTypeExpression(m[2]))
		any = true
	}
	if m := returnTagRe.FindStringSubmatch(body); m != nil {
		if m[1] != "?" {
			info.SetReturnType(NewJSTypeExpression(m[2]))
			any = true
		}
	}

	if !any {
		return nil
	}
	return info
}
