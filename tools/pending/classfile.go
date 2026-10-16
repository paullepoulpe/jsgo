package main

// A minimal reader for Java class files (JVMS §4): the constant pool, the
// class's flags, super types and inner classes, its fields (with
// ConstantValue and Signature) and methods (with Signature, parameter names
// from LocalVariableTable/MethodParameters, and the bytecode of the static
// initializer). That is all tools/pending needs to recover the signatures of
// not-yet-ported classes from the pinned closure-compiler jar.

import (
	"archive/zip"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strings"
	"unicode/utf16"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

const (
	accPublic     = 0x0001
	accPrivate    = 0x0002
	accStatic     = 0x0008
	accFinal      = 0x0010
	accVarargs    = 0x0080
	accBridge     = 0x0040
	accInterface  = 0x0200
	accAbstract   = 0x0400
	accSynthetic  = 0x1000
	accEnum       = 0x4000
	accAnnotation = 0x2000
)

type cpEntry struct {
	tag  byte
	a, b uint16
	str  string // Utf8 (decoded to a Go WTF-8 string)
	i64  int64  // Integer/Long
	f64  float64
}

// ClassInfo is one class file.
type ClassInfo struct {
	Name       string // internal name, e.g. com/google/javascript/jscomp/Foo$Bar
	Access     uint16
	Super      string
	Interfaces []string
	Signature  string // generic signature, or ""
	Fields     []*FieldInfo
	Methods    []*MethodInfo
	// Inner classes declared directly in this class: simple name → internal
	// name. OuterOf is the enclosing class of this class if it is a member
	// class.
	Inner   map[string]string
	OuterOf string
	// SimpleName is the member class's simple name (or the top-level name).
	SimpleName string
	Anonymous  bool
	// DiagnosticTypes: static field name → Go initializer recovered from
	// the static initializer (DiagnosticType.error/warning/disabled/make).
	StaticInits map[string]string
}

// FieldInfo is a field.
type FieldInfo struct {
	Access    uint16
	Name      string
	Desc      string
	Signature string
	Constant  *cpEntry
	Nullable  bool // @Nullable (a top-level type-use annotation)
}

// MethodInfo is a method or constructor.
type MethodInfo struct {
	Access     uint16
	Name       string
	Desc       string
	Signature  string
	ParamNames []string
	code       []byte
	// @Nullable type-use annotations on the return type and parameters.
	NullableReturn bool
	NullableParams map[int]bool
}

type classReader struct {
	b   []byte
	pos int
	err error
}

func (r *classReader) u1() byte {
	if r.pos+1 > len(r.b) {
		r.err = io.ErrUnexpectedEOF
		return 0
	}
	v := r.b[r.pos]
	r.pos++
	return v
}

func (r *classReader) u2() uint16 {
	if r.pos+2 > len(r.b) {
		r.err = io.ErrUnexpectedEOF
		return 0
	}
	v := binary.BigEndian.Uint16(r.b[r.pos:])
	r.pos += 2
	return v
}

func (r *classReader) u4() uint32 {
	if r.pos+4 > len(r.b) {
		r.err = io.ErrUnexpectedEOF
		return 0
	}
	v := binary.BigEndian.Uint32(r.b[r.pos:])
	r.pos += 4
	return v
}

func (r *classReader) bytes(n int) []byte {
	if r.pos+n > len(r.b) {
		r.err = io.ErrUnexpectedEOF
		return nil
	}
	v := r.b[r.pos : r.pos+n]
	r.pos += n
	return v
}

// decodeModifiedUTF8 decodes the JVM's modified UTF-8 into UTF-16 units and
// then into the port's WTF-8 strings.
func decodeModifiedUTF8(b []byte) string {
	units := make([]uint16, 0, len(b))
	for i := 0; i < len(b); {
		c := b[i]
		switch {
		case c < 0x80:
			units = append(units, uint16(c))
			i++
		case c&0xE0 == 0xC0 && i+1 < len(b):
			units = append(units, uint16(c&0x1F)<<6|uint16(b[i+1]&0x3F))
			i += 2
		case c&0xF0 == 0xE0 && i+2 < len(b):
			units = append(units, uint16(c&0x0F)<<12|uint16(b[i+1]&0x3F)<<6|uint16(b[i+2]&0x3F))
			i += 3
		default:
			units = append(units, uint16(c))
			i++
		}
	}
	_ = utf16.IsSurrogate
	return jutil.StringFromUTF16(units)
}

func parseClass(data []byte) (*ClassInfo, error) {
	r := &classReader{b: data}
	if r.u4() != 0xCAFEBABE {
		return nil, fmt.Errorf("bad magic")
	}
	r.u2()
	r.u2()
	n := int(r.u2())
	cp := make([]cpEntry, n)
	for i := 1; i < n; i++ {
		tag := r.u1()
		e := cpEntry{tag: tag}
		switch tag {
		case 1:
			l := int(r.u2())
			e.str = decodeModifiedUTF8(r.bytes(l))
		case 3:
			e.i64 = int64(int32(r.u4()))
		case 4:
			e.f64 = float64(math.Float32frombits(r.u4()))
		case 5:
			e.i64 = int64(uint64(r.u4())<<32 | uint64(r.u4()))
		case 6:
			e.f64 = math.Float64frombits(uint64(r.u4())<<32 | uint64(r.u4()))
		case 7, 8, 16, 19, 20:
			e.a = r.u2()
		case 9, 10, 11, 12, 17, 18:
			e.a, e.b = r.u2(), r.u2()
		case 15:
			e.a = uint16(r.u1())
			e.b = r.u2()
		default:
			return nil, fmt.Errorf("bad constant pool tag %d", tag)
		}
		cp[i] = e
		if tag == 5 || tag == 6 {
			i++
		}
	}
	utf := func(i uint16) string {
		if int(i) < len(cp) {
			return cp[i].str
		}
		return ""
	}
	className := func(i uint16) string {
		if i == 0 || int(i) >= len(cp) {
			return ""
		}
		return utf(cp[i].a)
	}
	c := &ClassInfo{Inner: map[string]string{}, StaticInits: map[string]string{}}
	c.Access = r.u2()
	c.Name = className(r.u2())
	c.Super = className(r.u2())
	for k := int(r.u2()); k > 0; k-- {
		c.Interfaces = append(c.Interfaces, className(r.u2()))
	}
	readAttrs := func(handle func(name string, body []byte)) {
		for k := int(r.u2()); k > 0; k-- {
			name := utf(r.u2())
			l := int(r.u4())
			body := r.bytes(l)
			handle(name, body)
		}
	}
	for k := int(r.u2()); k > 0; k-- {
		f := &FieldInfo{Access: r.u2(), Name: utf(r.u2()), Desc: utf(r.u2())}
		readAttrs(func(name string, body []byte) {
			switch name {
			case "Signature":
				f.Signature = utf(binary.BigEndian.Uint16(body))
			case "ConstantValue":
				e := cp[binary.BigEndian.Uint16(body)]
				if e.tag == 8 {
					e = cpEntry{tag: 8, str: utf(e.a)}
				}
				f.Constant = &e
			case "RuntimeVisibleTypeAnnotations", "RuntimeInvisibleTypeAnnotations":
				for _, a := range parseTypeAnnotations(body, utf) {
					if a.target == 0x13 && a.pathLen == 0 && isNullableAnnotation(a.typ) {
						f.Nullable = true
					}
				}
			}
		})
		c.Fields = append(c.Fields, f)
	}
	for k := int(r.u2()); k > 0; k-- {
		m := &MethodInfo{Access: r.u2(), Name: utf(r.u2()), Desc: utf(r.u2())}
		var lvt map[int]string
		var mp []string
		readAttrs(func(name string, body []byte) {
			switch name {
			case "Signature":
				m.Signature = utf(binary.BigEndian.Uint16(body))
			case "RuntimeVisibleTypeAnnotations", "RuntimeInvisibleTypeAnnotations":
				for _, a := range parseTypeAnnotations(body, utf) {
					if a.pathLen != 0 || !isNullableAnnotation(a.typ) {
						continue
					}
					switch a.target {
					case 0x14:
						m.NullableReturn = true
					case 0x16:
						if m.NullableParams == nil {
							m.NullableParams = map[int]bool{}
						}
						m.NullableParams[a.index] = true
					}
				}
			case "MethodParameters":
				cnt := int(body[0])
				for j := 0; j < cnt; j++ {
					mp = append(mp, utf(binary.BigEndian.Uint16(body[1+4*j:])))
				}
			case "Code":
				cr := &classReader{b: body}
				cr.u2()
				cr.u2()
				codeLen := int(cr.u4())
				m.code = cr.bytes(codeLen)
				for e := int(cr.u2()); e > 0; e-- {
					cr.bytes(8)
				}
				for a := int(cr.u2()); a > 0; a-- {
					aname := utf(cr.u2())
					abody := cr.bytes(int(cr.u4()))
					if aname == "LocalVariableTable" {
						if lvt == nil {
							lvt = map[int]string{}
						}
						cnt := int(binary.BigEndian.Uint16(abody))
						for j := 0; j < cnt; j++ {
							ent := abody[2+10*j:]
							start := binary.BigEndian.Uint16(ent)
							nm := utf(binary.BigEndian.Uint16(ent[4:]))
							idx := int(binary.BigEndian.Uint16(ent[8:]))
							if start == 0 {
								lvt[idx] = nm
							}
						}
					}
				}
			}
		})
		m.ParamNames = paramNames(m, mp, lvt)
		c.Methods = append(c.Methods, m)
	}
	readAttrs(func(name string, body []byte) {
		switch name {
		case "Signature":
			c.Signature = utf(binary.BigEndian.Uint16(body))
		case "InnerClasses":
			cnt := int(binary.BigEndian.Uint16(body))
			for j := 0; j < cnt; j++ {
				ent := body[2+8*j:]
				inner := className(binary.BigEndian.Uint16(ent))
				outer := className(binary.BigEndian.Uint16(ent[2:]))
				simple := utf(binary.BigEndian.Uint16(ent[4:]))
				if inner == c.Name {
					c.OuterOf = outer
					c.SimpleName = simple
					if simple == "" {
						c.Anonymous = true
					}
				}
				if outer == c.Name && simple != "" {
					c.Inner[simple] = inner
				}
			}
		}
	})
	if r.err != nil {
		return nil, r.err
	}
	if c.SimpleName == "" && !c.Anonymous {
		c.SimpleName = c.Name[strings.LastIndexByte(c.Name, '/')+1:]
		if i := strings.LastIndexByte(c.SimpleName, '$'); i >= 0 {
			// A local or anonymous class (no InnerClasses entry naming it).
			c.Anonymous = true
		}
	}
	for _, m := range c.Methods {
		if m.Name == "<clinit>" {
			scanStaticInit(c, cp, m.code)
		}
	}
	return c, nil
}

type typeAnnotation struct {
	target  byte
	index   int // formal parameter index (target 0x16)
	pathLen int
	typ     string
}

func isNullableAnnotation(desc string) bool {
	return strings.HasSuffix(desc, "/Nullable;")
}

// parseTypeAnnotations reads a Runtime(In)VisibleTypeAnnotations attribute
// of a field or method (JVMS §4.7.20).
func parseTypeAnnotations(body []byte, utf func(uint16) string) []typeAnnotation {
	r := &classReader{b: body}
	var out []typeAnnotation
	var skipElement func()
	var skipAnnotation func()
	skipElement = func() {
		switch r.u1() {
		case 'B', 'C', 'D', 'F', 'I', 'J', 'S', 'Z', 's', 'c':
			r.u2()
		case 'e':
			r.u2()
			r.u2()
		case '@':
			skipAnnotation()
		case '[':
			for n := int(r.u2()); n > 0; n-- {
				skipElement()
			}
		}
	}
	skipAnnotation = func() {
		r.u2()
		for n := int(r.u2()); n > 0; n-- {
			r.u2()
			skipElement()
		}
	}
	for n := int(r.u2()); n > 0 && r.err == nil; n-- {
		a := typeAnnotation{target: r.u1()}
		switch a.target {
		case 0x00, 0x01:
			r.u1()
		case 0x10:
			r.u2()
		case 0x11, 0x12:
			r.u1()
			r.u1()
		case 0x13, 0x14, 0x15:
		case 0x16:
			a.index = int(r.u1())
		case 0x17:
			r.u2()
		default:
			return out // code-attribute targets do not occur here
		}
		a.pathLen = int(r.u1())
		r.bytes(2 * a.pathLen)
		a.typ = utf(r.u2())
		for p := int(r.u2()); p > 0; p-- {
			r.u2()
			skipElement()
		}
		out = append(out, a)
	}
	return out
}

// paramNames recovers the parameter names: MethodParameters if present,
// else the LocalVariableTable entries live at pc 0 in the parameter slots.
func paramNames(m *MethodInfo, mp []string, lvt map[int]string) []string {
	params, _ := parseMethodDescriptor(m.Desc)
	names := make([]string, len(params))
	if len(mp) == len(params) {
		copy(names, mp)
	} else if lvt != nil {
		slot := 0
		if m.Access&accStatic == 0 {
			slot = 1
		}
		for i, p := range params {
			names[i] = lvt[slot]
			slot++
			if p == "J" || p == "D" {
				slot++
			}
		}
	}
	return names
}

// parseMethodDescriptor splits a method descriptor into parameter and
// return field descriptors.
func parseMethodDescriptor(desc string) ([]string, string) {
	var params []string
	i := 1
	for i < len(desc) && desc[i] != ')' {
		j := i
		for desc[j] == '[' {
			j++
		}
		if desc[j] == 'L' {
			j = strings.IndexByte(desc[j:], ';') + j
		}
		params = append(params, desc[i:j+1])
		i = j + 1
	}
	return params, desc[i+1:]
}

// opcodeLength is the length of the instruction at pc.
func opcodeLength(code []byte, pc int) int {
	op := code[pc]
	switch {
	case op == 16 || op == 18 || (op >= 21 && op <= 25) || (op >= 54 && op <= 58) || op == 169 || op == 188:
		return 2
	case op == 17 || op == 19 || op == 20 || op == 132 || (op >= 153 && op <= 168) ||
		(op >= 178 && op <= 184) || op == 187 || op == 189 || op == 192 || op == 193 || op == 198 || op == 199:
		return 3
	case op == 197:
		return 4
	case op == 185 || op == 186 || op == 200 || op == 201:
		return 5
	case op == 196:
		if code[pc+1] == 132 {
			return 6
		}
		return 4
	case op == 170: // tableswitch
		p := (pc + 4) &^ 3
		lo := int32(binary.BigEndian.Uint32(code[p+4:]))
		hi := int32(binary.BigEndian.Uint32(code[p+8:]))
		return p - pc + 12 + 4*int(hi-lo+1)
	case op == 171: // lookupswitch
		p := (pc + 4) &^ 3
		n := int(binary.BigEndian.Uint32(code[p+4:]))
		return p - pc + 8 + 8*n
	}
	return 1
}

// scanStaticInit recovers `static final DiagnosticType X =
// DiagnosticType.error|warning|disabled("KEY", "format")` (the format may be
// built with JSCompStrings.lines) and `DiagnosticType.make("KEY",
// CheckLevel.L, "format")` from <clinit>.
func scanStaticInit(c *ClassInfo, cp []cpEntry, code []byte) {
	// A tiny abstract interpreter over <clinit>: string and int constants,
	// String[] varargs arrays, CheckLevel constants, JSCompStrings.lines and
	// the DiagnosticType factories. Any other instruction clears the stack,
	// so only initializers built entirely from those are recovered.
	type value struct {
		kind  byte // 's' string, 'i' int, 'a' String[], 'l' CheckLevel, 'd' DiagnosticType expression
		s     string
		i     int
		elems []*string
	}
	var stack []*value
	push := func(v *value) { stack = append(stack, v) }
	pop := func() *value {
		if len(stack) == 0 {
			return nil
		}
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		return v
	}
	ref := func(idx uint16) (owner, name, desc string) {
		e := cp[idx]
		owner = cp[cp[e.a].a].str
		nt := cp[e.b]
		return owner, cp[nt.a].str, cp[nt.b].str
	}
	const dt = "com/google/javascript/jscomp/DiagnosticType"
	const strings_ = "com/google/javascript/jscomp/base/JSCompStrings"
	for pc := 0; pc < len(code); pc += opcodeLength(code, pc) {
		op := code[pc]
		switch {
		case op == 18 || op == 19: // ldc, ldc_w
			idx := uint16(code[pc+1])
			if op == 19 {
				idx = binary.BigEndian.Uint16(code[pc+1:])
			}
			if e := cp[idx]; e.tag == 8 {
				push(&value{kind: 's', s: cp[e.a].str})
			} else {
				push(nil)
			}
		case op >= 2 && op <= 8: // iconst_m1 .. iconst_5
			push(&value{kind: 'i', i: int(op) - 3})
		case op == 16: // bipush
			push(&value{kind: 'i', i: int(int8(code[pc+1]))})
		case op == 17: // sipush
			push(&value{kind: 'i', i: int(int16(binary.BigEndian.Uint16(code[pc+1:])))})
		case op == 189: // anewarray
			n := pop()
			if n != nil && n.kind == 'i' && n.i >= 0 && n.i < 1<<12 {
				push(&value{kind: 'a', elems: make([]*string, n.i)})
			} else {
				push(nil)
			}
		case op == 89: // dup
			if len(stack) == 0 {
				push(nil)
			}
			push(stack[len(stack)-1])
		case op == 83: // aastore
			v, idx, arr := pop(), pop(), pop()
			if arr != nil && arr.kind == 'a' && idx != nil && idx.kind == 'i' && idx.i < len(arr.elems) && v != nil && v.kind == 's' {
				s := v.s
				arr.elems[idx.i] = &s
			}
		case op == 178: // getstatic
			owner, name, _ := ref(binary.BigEndian.Uint16(code[pc+1:]))
			if owner == "com/google/javascript/jscomp/CheckLevel" {
				push(&value{kind: 'l', s: name})
			} else {
				push(nil)
			}
		case op == 184: // invokestatic
			owner, name, desc := ref(binary.BigEndian.Uint16(code[pc+1:]))
			switch {
			case owner == strings_ && name == "lines" && desc == "([Ljava/lang/String;)Ljava/lang/String;":
				arr := pop()
				var joined *value
				if arr != nil && arr.kind == 'a' {
					parts := make([]string, len(arr.elems))
					ok := true
					for i, e := range arr.elems {
						if e == nil {
							ok = false
							break
						}
						parts[i] = *e
					}
					if ok {
						joined = &value{kind: 's', s: strings.Join(parts, "\n")}
					}
				}
				push(joined)
			case owner == dt && desc == "(Ljava/lang/String;Ljava/lang/String;)L"+dt+";" &&
				(name == "error" || name == "warning" || name == "disabled"):
				s2, s1 := pop(), pop()
				if s1 != nil && s1.kind == 's' && s2 != nil && s2.kind == 's' {
					push(&value{kind: 'd', s: fmt.Sprintf("DiagnosticType_%s(%s, %s)", capitalize(name), goQuote(s1.s), goQuote(s2.s))})
				} else {
					push(nil)
				}
			case owner == dt && name == "make" && desc == "(Ljava/lang/String;Lcom/google/javascript/jscomp/CheckLevel;Ljava/lang/String;)L"+dt+";":
				s2, l, s1 := pop(), pop(), pop()
				if s1 != nil && s1.kind == 's' && l != nil && l.kind == 'l' && s2 != nil && s2.kind == 's' {
					push(&value{kind: 'd', s: fmt.Sprintf("DiagnosticType_Make(%s, CheckLevel_%s, %s)", goQuote(s1.s), l.s, goQuote(s2.s))})
				} else {
					push(nil)
				}
			default:
				stack = stack[:0]
			}
		case op == 179: // putstatic
			fowner, fname, _ := ref(binary.BigEndian.Uint16(code[pc+1:]))
			v := pop()
			if v != nil && v.kind == 'd' && fowner == c.Name {
				c.StaticInits[fname] = v.s
			}
		default:
			stack = stack[:0]
		}
	}
}

// readJar parses every class of the given package prefixes in the jar.
func readJar(path string, prefixes []string) (map[string]*ClassInfo, error) {
	zr, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	classes := map[string]*ClassInfo{}
	for _, f := range zr.File {
		if !strings.HasSuffix(f.Name, ".class") || strings.Contains(f.Name, "/jarjar/") {
			continue
		}
		keep := false
		for _, p := range prefixes {
			if strings.HasPrefix(f.Name, p) {
				keep = true
			}
		}
		if !keep {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return nil, err
		}
		c, err := parseClass(data)
		if err != nil {
			return nil, fmt.Errorf("%s: %v", f.Name, err)
		}
		classes[c.Name] = c
	}
	return classes, nil
}
