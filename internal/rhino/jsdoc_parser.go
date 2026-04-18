// Port of com/google/javascript/jscomp/parsing/JsDocInfoParser.java @ closure-compiler v20260415
//
// Recognises `@type`, `@param`, `@return`, `@extends`, `@implements`,
// `@typedef`, `@template`, plus the marker tags `@constructor`,
// `@interface`, `@record`. The Java implementation is a full
// tokeniser + type-expression parser; we use a small hand-rolled
// scanner that walks the comment body and extracts balanced `{...}`
// bodies so inline record types like `{foo: number}` come through
// intact.

package rhino

import (
	"strings"
	"unicode"
)

// ParseJSDoc returns a JSDocInfo for the given block-comment body (the
// text between /** and */, without the delimiters), or nil if none of
// the recognised tags appear.
func ParseJSDoc(body string) *JSDocInfo {
	info := &JSDocInfo{}
	any := false
	i := 0
	for i < len(body) {
		if body[i] != '@' {
			i++
			continue
		}
		i++ // consume '@'
		tag, next := readTag(body, i)
		i = next
		switch tag {
		case "type":
			if expr, j, ok := readTypeExpr(body, i); ok {
				info.SetType(NewJSTypeExpression(expr))
				any = true
				i = j
			}
		case "param":
			if expr, j, ok := readTypeExpr(body, i); ok {
				ident, k := readIdentAfter(body, j)
				if ident != "" {
					info.AddParam(ident, NewJSTypeExpression(expr))
					any = true
					i = k
					continue
				}
				i = j
			}
		case "return", "returns":
			if expr, j, ok := readTypeExpr(body, i); ok {
				info.SetReturnType(NewJSTypeExpression(expr))
				any = true
				i = j
			}
		case "extends":
			if expr, j, ok := readTypeExpr(body, i); ok {
				info.SetBaseType(NewJSTypeExpression(expr))
				any = true
				i = j
			}
		case "implements":
			if expr, j, ok := readTypeExpr(body, i); ok {
				info.AddImplementation(NewJSTypeExpression(expr))
				any = true
				i = j
			}
		case "typedef":
			if expr, j, ok := readTypeExpr(body, i); ok {
				info.SetTypedefType(NewJSTypeExpression(expr))
				any = true
				i = j
			}
		case "template":
			names, j := readTemplateNames(body, i)
			for _, tn := range names {
				info.AddTemplateName(tn)
			}
			if len(names) > 0 {
				any = true
			}
			i = j
		case "constructor":
			info.SetConstructor(true)
			any = true
		case "interface":
			info.SetInterface(true)
			any = true
		case "record":
			info.SetRecord(true)
			any = true
		}
	}
	if !any {
		return nil
	}
	return info
}

// readTag reads a tag name (letters only) and returns (tag, nextPos).
func readTag(body string, i int) (string, int) {
	start := i
	for i < len(body) {
		c := body[i]
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z') {
			break
		}
		i++
	}
	return strings.ToLower(body[start:i]), i
}

// readTypeExpr skips whitespace then a `{`, consumes balanced `{...}`
// content, and returns the inner body + position past the closing `}`.
// Handles nested `{}` so inline record types come through intact.
func readTypeExpr(body string, i int) (string, int, bool) {
	for i < len(body) && (body[i] == ' ' || body[i] == '\t' || body[i] == '\n' || body[i] == '\r') {
		i++
	}
	if i >= len(body) || body[i] != '{' {
		return "", i, false
	}
	i++
	depth := 1
	start := i
	for i < len(body) && depth > 0 {
		switch body[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return body[start:i], i + 1, true
			}
		}
		i++
	}
	return "", i, false
}

// readIdentAfter skips whitespace + optional `=` (for `@param {T} [x=default]`)
// and reads a JS-identifier. Returns ("", pos) if none.
func readIdentAfter(body string, i int) (string, int) {
	for i < len(body) && (body[i] == ' ' || body[i] == '\t' || body[i] == '\n' || body[i] == '\r' || body[i] == '*' || body[i] == '[') {
		i++
	}
	start := i
	for i < len(body) {
		c := body[i]
		if !(unicode.IsLetter(rune(c)) || unicode.IsDigit(rune(c)) || c == '_' || c == '$') {
			break
		}
		i++
	}
	return body[start:i], i
}

// readTemplateNames reads `@template T, U, V` style lists, stopping at
// newline or start of another @-tag.
func readTemplateNames(body string, i int) ([]string, int) {
	var out []string
	for i < len(body) && (body[i] == ' ' || body[i] == '\t') {
		i++
	}
	start := i
	for i < len(body) && body[i] != '\n' && body[i] != '@' {
		i++
	}
	raw := body[start:i]
	for _, p := range strings.Split(raw, ",") {
		p = strings.TrimSpace(p)
		// Also trim trailing `*` / whitespace from JSDoc line starts.
		p = strings.TrimRight(p, " \t*")
		if p != "" && isTemplateIdent(p) {
			out = append(out, p)
		}
	}
	return out, i
}

func isTemplateIdent(s string) bool {
	for i := 0; i < len(s); i++ {
		c := rune(s[i])
		if !(unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_' || c == '$') {
			return false
		}
	}
	return true
}
