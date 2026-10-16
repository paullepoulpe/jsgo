package main

// Java type descriptors and generic signatures (JVMS §4.3, §4.7.9.1).

import (
	"strings"
)

// JKind is the kind of a Java type.
type JKind int

const (
	JPrim JKind = iota
	JClass
	JArray
	JTypeVar
	JWildcard // '*' (Bound == nil), '+' (extends Bound) or '-' (super Bound)
)

// JType is a Java type.
type JType struct {
	Kind  JKind
	Prim  byte   // JPrim: Z B C S I J F D V
	Class string // JClass: internal name, nested classes joined with '$'
	Args  []*JType
	Elem  *JType // JArray
	Var   string // JTypeVar
	Bound *JType // JWildcard
	Super bool   // JWildcard: '-'
}

// JTypeParam is a type parameter with its first bound.
type JTypeParam struct {
	Name  string
	Bound *JType
}

// JMethodSig is a parsed method descriptor or signature.
type JMethodSig struct {
	TypeParams []JTypeParam
	Params     []*JType
	Result     *JType
}

type sigParser struct {
	s string
	i int
}

func (p *sigParser) peek() byte {
	if p.i < len(p.s) {
		return p.s[p.i]
	}
	return 0
}

func (p *sigParser) ident() string {
	j := p.i
	for p.i < len(p.s) && !strings.ContainsRune(".;[/<>:", rune(p.s[p.i])) {
		p.i++
	}
	return p.s[j:p.i]
}

func (p *sigParser) typeParams() []JTypeParam {
	if p.peek() != '<' {
		return nil
	}
	p.i++
	var out []JTypeParam
	for p.peek() != '>' {
		name := p.ident()
		var bound *JType
		// ClassBound: ':' [ReferenceTypeSignature]
		p.i++
		if c := p.peek(); c == 'L' || c == 'T' || c == '[' {
			bound = p.typeSig()
		}
		// InterfaceBounds
		for p.peek() == ':' {
			p.i++
			b := p.typeSig()
			if bound == nil {
				bound = b
			}
		}
		out = append(out, JTypeParam{Name: name, Bound: bound})
	}
	p.i++
	return out
}

func (p *sigParser) typeSig() *JType {
	c := p.peek()
	switch c {
	case 'Z', 'B', 'C', 'S', 'I', 'J', 'F', 'D', 'V':
		p.i++
		return &JType{Kind: JPrim, Prim: c}
	case '[':
		p.i++
		return &JType{Kind: JArray, Elem: p.typeSig()}
	case 'T':
		p.i++
		name := p.ident()
		p.i++ // ';'
		return &JType{Kind: JTypeVar, Var: name}
	case 'L':
		p.i++
		var sb strings.Builder
		var args []*JType
		for {
			seg := p.ident()
			sb.WriteString(seg)
			switch p.peek() {
			case '/':
				sb.WriteByte('/')
				p.i++
				continue
			case '<':
				args = p.typeArgs()
			}
			if p.peek() == '.' {
				// Inner class of a parameterized outer: the outer's type
				// arguments are dropped (static nesting in Go).
				sb.WriteByte('$')
				p.i++
				args = nil
				continue
			}
			break
		}
		p.i++ // ';'
		return &JType{Kind: JClass, Class: unshade(sb.String()), Args: args}
	}
	panic("bad signature " + p.s + " at " + p.s[p.i:])
}

func (p *sigParser) typeArgs() []*JType {
	p.i++ // '<'
	var out []*JType
	for p.peek() != '>' {
		switch p.peek() {
		case '*':
			p.i++
			out = append(out, &JType{Kind: JWildcard})
		case '+':
			p.i++
			out = append(out, &JType{Kind: JWildcard, Bound: p.typeSig()})
		case '-':
			p.i++
			out = append(out, &JType{Kind: JWildcard, Bound: p.typeSig(), Super: true})
		default:
			out = append(out, p.typeSig())
		}
	}
	p.i++
	return out
}

// unshade maps a class relocated into the jar's jarjar package back to its
// original name (the jar shades Guava, protobuf, ...).
func unshade(name string) string {
	return strings.TrimPrefix(name, "com/google/javascript/jscomp/jarjar/")
}

// parseFieldType parses a field descriptor or field signature.
func parseFieldType(s string) *JType {
	p := &sigParser{s: s}
	return p.typeSig()
}

// parseMethodSig parses a method descriptor or method signature.
func parseMethodSig(s string) *JMethodSig {
	p := &sigParser{s: s}
	m := &JMethodSig{TypeParams: p.typeParams()}
	p.i++ // '('
	for p.peek() != ')' {
		m.Params = append(m.Params, p.typeSig())
	}
	p.i++
	m.Result = p.typeSig()
	return m
}

// parseClassSig parses a class signature: type parameters, superclass and
// interfaces.
func parseClassSig(s string) ([]JTypeParam, []*JType) {
	p := &sigParser{s: s}
	tps := p.typeParams()
	var supers []*JType
	for p.i < len(p.s) {
		supers = append(supers, p.typeSig())
	}
	return tps, supers
}

// javaSimpleName renders a type the way a Java signature line shows it:
// simple class names, e.g. "List<Node>", "Node[]", "int".
func javaSimpleName(t *JType) string {
	switch t.Kind {
	case JPrim:
		return map[byte]string{'Z': "boolean", 'B': "byte", 'C': "char", 'S': "short", 'I': "int", 'J': "long", 'F': "float", 'D': "double", 'V': "void"}[t.Prim]
	case JArray:
		return javaSimpleName(t.Elem) + "[]"
	case JTypeVar:
		return t.Var
	case JWildcard:
		if t.Bound == nil {
			return "?"
		}
		if t.Super {
			return "? super " + javaSimpleName(t.Bound)
		}
		return "? extends " + javaSimpleName(t.Bound)
	}
	name := t.Class[strings.LastIndexByte(t.Class, '/')+1:]
	name = strings.ReplaceAll(name, "$", ".")
	if len(t.Args) > 0 {
		var args []string
		for _, a := range t.Args {
			args = append(args, javaSimpleName(a))
		}
		name += "<" + strings.Join(args, ", ") + ">"
	}
	return name
}

// javaErasedSimpleName is javaSimpleName without type arguments, as used in
// the PORTING.md `// Java:` lines (e.g. "List<Node>" stays generic there,
// so this is only used for varargs rendering).
func javaParamList(sig *JMethodSig, varargs bool) string {
	var parts []string
	for i, p := range sig.Params {
		s := javaSimpleName(p)
		if varargs && i == len(sig.Params)-1 && strings.HasSuffix(s, "[]") {
			s = strings.TrimSuffix(s, "[]") + "..."
		}
		parts = append(parts, s)
	}
	return strings.Join(parts, ", ")
}
