package javaast

import (
	"regexp"
	"strings"
)

// TryTrivial attempts to recognise one of the "trivial" body patterns
// (#161) and returns the translated Go body text (without outer braces)
// plus true on success. The caller should write:
//
//	"func (...) Name(...) RetType {\n" + body + "}\n"
//
// Parameters:
//
//	javaBody   – raw text inside the Java method braces
//	recvName   – the Go receiver variable name (e.g. "c")
//	className  – the Go struct name (e.g. "CheckJSDoc")
//	typeMapper – Java→Go type converter
func TryTrivial(javaBody, recvName, className string, typeMapper func(string) string) (string, bool) {
	body := strings.TrimSpace(javaBody)

	// Getter: `return foo;` or `return this.foo;`
	if s, ok := tryGetter(body, recvName); ok {
		return s, true
	}
	// Setter: `this.foo = foo;`
	if s, ok := trySetter(body, recvName); ok {
		return s, true
	}
	// Delegation: `return other.method(args);` or `other.method(args);`
	if s, ok := tryDelegation(body, recvName); ok {
		return s, true
	}
	// Single-expression return: `return expr;`
	if s, ok := trySingleReturn(body, recvName); ok {
		return s, true
	}
	// Null-check early return: `if (x == null) return y;`
	if s, ok := tryNullCheck(body, recvName); ok {
		return s, true
	}
	// Field init: `this.foo = bar; this.baz = qux;` (multiple assignments only)
	if s, ok := tryFieldInit(body, recvName); ok {
		return s, true
	}
	return "", false
}

// ---- getter ------------------------------------------------------------
// Patterns:
//   `return foo;`
//   `return this.foo;`
//   `return foo != null ? foo : "";`  (guarded getter — still trivial)

var getterRe = regexp.MustCompile(`^return\s+(this\.)?(\w+)\s*;$`)

func tryGetter(body, recv string) (string, bool) {
	m := getterRe.FindStringSubmatch(body)
	if m == nil {
		return "", false
	}
	field := m[2]
	return "\treturn " + recv + "." + field + "\n", true
}

// ---- setter ------------------------------------------------------------
// Pattern: `this.foo = foo;`

var setterRe = regexp.MustCompile(`^this\.(\w+)\s*=\s*(\w+)\s*;$`)

func trySetter(body, recv string) (string, bool) {
	m := setterRe.FindStringSubmatch(body)
	if m == nil {
		return "", false
	}
	field := m[1]
	param := m[2]
	return "\t" + recv + "." + field + " = " + param + "\n", true
}

// ---- delegation --------------------------------------------------------
// Pattern: `[return] receiver.method(args...);`
// where receiver may be `this`, a field, or another local.

var delegationRe = regexp.MustCompile(`^(return\s+)?(\w+)\.(\w+)\(([^)]*)\)\s*;$`)

func tryDelegation(body, recv string) (string, bool) {
	m := delegationRe.FindStringSubmatch(body)
	if m == nil {
		return "", false
	}
	hasReturn := m[1] != ""
	obj := m[2]
	method := goMethodName(m[3])
	args := strings.TrimSpace(m[4])

	goObj := obj
	if obj == "this" {
		goObj = recv
	} else {
		goObj = recv + "." + obj
	}

	call := goObj + "." + method + "(" + args + ")"
	if hasReturn {
		return "\treturn " + call + "\n", true
	}
	return "\t" + call + "\n", true
}

// ---- single-expression return -----------------------------------------
// Pattern: `return <expr>;`  (anything simple)

var singleReturnRe = regexp.MustCompile(`^return\s+(.+)\s*;$`)

func trySingleReturn(body, recv string) (string, bool) {
	m := singleReturnRe.FindStringSubmatch(body)
	if m == nil {
		return "", false
	}
	exprText := strings.TrimSpace(m[1])
	// Don't handle multi-line or block expressions here.
	if strings.ContainsAny(exprText, "\n{}") {
		return "", false
	}
	e := parseExpr(exprText)
	tr := &translator{ctx: MethodContext{RecvName: recv}, indent: 1}
	return "\treturn " + tr.emitExpr(e) + "\n", true
}

// ---- null check --------------------------------------------------------
// Pattern: `if (x == null) return y;`
//          `if (x != null) return y;`

var nullCheckRe = regexp.MustCompile(`^if\s*\(\s*(\w+)\s*(==|!=)\s*null\s*\)\s*return\s+(.+?)\s*;$`)

func tryNullCheck(body, recv string) (string, bool) {
	m := nullCheckRe.FindStringSubmatch(body)
	if m == nil {
		return "", false
	}
	varName := m[1]
	op := m[2]
	retExpr := strings.TrimSpace(m[3])

	goOp := op // == or !=
	e := parseExpr(retExpr)
	tr := &translator{ctx: MethodContext{RecvName: recv}, indent: 2}
	ret := tr.emitExpr(e)

	var sb strings.Builder
	sb.WriteString("\tif " + varName + " " + goOp + " nil {\n")
	sb.WriteString("\t\treturn " + ret + "\n")
	sb.WriteString("\t}\n")
	return sb.String(), true
}

// ---- field init --------------------------------------------------------
// Pattern: `this.foo = bar; this.baz = qux; ...` (all this.X = Y assignments)

var fieldAssignRe = regexp.MustCompile(`this\.(\w+)\s*=\s*([^;]+?)\s*;`)

func tryFieldInit(body, recv string) (string, bool) {
	// Every statement must match `this.X = Y;`
	stmts := splitSemicolon(body)
	if len(stmts) == 0 {
		return "", false
	}
	var lines []string
	for _, stmt := range stmts {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		m := fieldAssignRe.FindStringSubmatch(stmt + ";")
		if m == nil {
			return "", false
		}
		field := m[1]
		val := strings.TrimSpace(m[2])
		lines = append(lines, "\t"+recv+"."+field+" = "+val)
	}
	if len(lines) == 0 {
		return "", false
	}
	return strings.Join(lines, "\n") + "\n", true
}

func splitSemicolon(s string) []string {
	var parts []string
	depth := 0
	start := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			if depth > 0 {
				depth--
			}
		case ';':
			if depth == 0 {
				parts = append(parts, s[start:i])
				start = i + 1
			}
		}
	}
	tail := strings.TrimSpace(s[start:])
	if tail != "" {
		parts = append(parts, tail)
	}
	return parts
}
