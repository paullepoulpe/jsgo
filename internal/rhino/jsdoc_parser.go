// Port of com/google/javascript/jscomp/parsing/JsDocInfoParser.java @ closure-compiler v20260415
//
// Tier-0 recognises only `@type {<name>}` inside a block comment. Later
// tiers add @param, @return, @template, union/record type expressions,
// etc. The Java parser is token-driven; this placeholder uses a small
// regex to keep the scaffold tight, and will be rewritten as the type
// grammar grows.

package rhino

import "regexp"

// Matches `@type { !? Name }` — the leading `!` / `?` nullability modifier
// is handled explicitly so we can distinguish non-null (safe to resolve
// to the base type) from nullable (too imprecise for Tier-0).
var typeTagRe = regexp.MustCompile(`@type\s*\{\s*([!?]?)\s*([A-Za-z_][A-Za-z0-9_]*)\s*\}`)

// ParseJSDoc returns a JSDocInfo for the given block-comment body (the
// text between /** and */, without the delimiters), or nil if the comment
// has no recognised tags.
func ParseJSDoc(body string) *JSDocInfo {
	m := typeTagRe.FindStringSubmatch(body)
	if m == nil {
		return nil
	}
	// Nullable-prefix (`?Foo`) is out of scope for Tier-0; leave the
	// type unrecorded so the checker is conservative instead of warning
	// on legitimate null assignments.
	if m[1] == "?" {
		return nil
	}
	info := &JSDocInfo{}
	info.SetType(NewJSTypeExpression(m[2]))
	return info
}
