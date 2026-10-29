package rhino

// Go-only: differential type-algebra oracle for the jstype port.
//
// This is the Go mirror of tools/oracle/java/TypeAlgebraDump.java. It builds
// the same type universe through the same JSTypeRegistry API calls, in the
// same order (many operations mutate types: lazy prototypes, call/apply/bind
// properties, caches), prints the same facts and compares the result with
// the Java output checked in at testdata/jstype_type_algebra.txt.gz.
//
// JSDoc type expressions need JsDocInfoParser, which lives in jscomp and is
// not ported here; the Java dump therefore prints the parsed type AST of
// every type string (the "ast" lines) and this test rebuilds those ASTs from
// the expected output before handing them to createTypeFromCommentNode.
//
// Regenerate the expected output with:
//
//	javac -cp closure-compiler.jar -d build tools/oracle/java/TypeAlgebraDump.java
//	java -cp closure-compiler.jar:build \
//	  com.google.javascript.rhino.jstype.TypeAlgebraDump --pairs=all | gzip -9 \
//	  > internal/rhino/testdata/jstype_type_algebra.txt.gz
//
// JSTYPE_ORACLE_FILE overrides the expected file (plain or .gz), e.g. for a
// dump made with a different --pairs setting.

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

type typeAlgebraDump struct {
	out      *strings.Builder
	labels   []string
	types    []JSType
	registry *JSTypeRegistry
	scope    *typeAlgebraTestScope
	seed     int64
	// Set by build(); used by the facts.
	tKey       *TemplateType
	partialMap *TemplateTypeMap
	totalMap   *TemplateTypeMap
	inferMap   *jutil.LinkedHashMap[*TemplateType, JSType]
	// asts maps a type string to the serialized AST the Java dump printed.
	asts map[string]string
}

// ---------------------------------------------------------------------
// Output helpers

func taEsc(s string) string {
	var sb strings.Builder
	for _, c := range jutil.StringToUTF16(s) {
		if c == '\\' {
			sb.WriteString("\\\\")
		} else if c == '\n' {
			sb.WriteString("\\n")
		} else if c < 0x20 || c >= 0x7f {
			fmt.Fprintf(&sb, "\\u%04x", int(c))
		} else {
			sb.WriteByte(byte(c))
		}
	}
	return sb.String()
}

var jsTypeReflectType = reflect.TypeOf((*JSType)(nil)).Elem()

func taStr(o any) string {
	if o == nil {
		return "null"
	}
	rv := reflect.ValueOf(o)
	switch rv.Kind() {
	case reflect.Slice:
		// Java collections returned by the API are non-null unless the fact
		// says otherwise (it then converts a nil slice to nil).
		var sb strings.Builder
		sb.WriteByte('[')
		for i := 0; i < rv.Len(); i++ {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(taStr(rv.Index(i).Interface()))
		}
		sb.WriteByte(']')
		return sb.String()
	case reflect.Pointer, reflect.Map, reflect.Interface:
		if rv.IsNil() {
			if rv.Type().Implements(jsTypeReflectType) {
				// A typed nil inside an interface would not compare equal to
				// nil in Go callers; flag it so that it shows as a diff.
				return "null(typed " + rv.Type().String() + ")"
			}
			return "null"
		}
	}
	switch v := o.(type) {
	case string:
		return v
	case *string:
		return *v
	case bool:
		return strconv.FormatBool(v)
	case int:
		return strconv.Itoa(v)
	case *FunctionType_Parameter:
		return "P(" + taStr(v.GetJSType()) + "," + strconv.FormatBool(v.IsOptional()) + "," +
			strconv.FormatBool(v.IsVariadic()) + ")"
	case *JSType_TypePair:
		return "<" + taStr(v.TypeA) + " ; " + taStr(v.TypeB) + ">"
	case fmt.Stringer:
		return v.String()
	}
	return jutil.ToJavaString(o)
}

func (d *typeAlgebraDump) line(s string) {
	d.out.WriteString(s)
	d.out.WriteByte('\n')
}

func (d *typeAlgebraDump) rec(label, value string) {
	d.line(label + "\t" + taEsc(value))
}

// taExceptionName is the Java simple class name of the exception that the
// panic r stands for.
func taExceptionName(r any) string {
	switch v := r.(type) {
	case jutil.Throwable:
		s := v.Error()
		if i := strings.Index(s, ": "); i >= 0 {
			s = s[:i]
		}
		if i := strings.LastIndexAny(s, ".$"); i >= 0 {
			s = s[i+1:]
		}
		return s
	case runtime.Error:
		msg := v.Error()
		switch {
		case strings.Contains(msg, "nil pointer dereference"):
			return "NullPointerException"
		case strings.Contains(msg, "interface conversion"):
			return "ClassCastException"
		case strings.Contains(msg, "index out of range"):
			return "IndexOutOfBoundsException"
		}
		return "GoRuntimeError(" + msg + ")"
	}
	return "GoPanic(" + fmt.Sprint(r) + ")"
}

// taPf runs s, returning its result or the class of the exception it threw.
func taPf(s func() any) (v string) {
	defer func() {
		if r := recover(); r != nil {
			v = "!" + taExceptionName(r)
		}
	}()
	return taStr(s())
}

func (d *typeAlgebraDump) fact(name string, s func() any) {
	d.rec(name, taPf(s))
}

// ---------------------------------------------------------------------
// Deterministic sequences

func (d *typeAlgebraDump) nextInt(n int) int {
	d.seed = d.seed*6364136223846793005 + 1442695040888963407
	return int(int64(uint64(d.seed)>>33) % int64(n))
}

// ---------------------------------------------------------------------
// A minimal scope and reporter

type typeAlgebraTestReporter struct {
	sb *strings.Builder
}

func (r *typeAlgebraTestReporter) Warning(message string, sourceName string, line int, lineOffset int) {
	r.sb.WriteString("WARNING\t" + taEsc(message+" @"+sourceName+":"+strconv.Itoa(line)+":"+strconv.Itoa(lineOffset)) + "\n")
}

func (r *typeAlgebraTestReporter) Error(message string, sourceName string, line int, lineOffset int) {
	r.sb.WriteString("ERROR\t" + taEsc(message+" @"+sourceName+":"+strconv.Itoa(line)+":"+strconv.Itoa(lineOffset)) + "\n")
}

type typeAlgebraTestScope struct {
	root  Node
	slots *jutil.LinkedHashMap[string, *SimpleSlotBase]
}

func (s *typeAlgebraTestScope) GetRootNode() Node { return s.root }

func (s *typeAlgebraTestScope) GetParentScope() StaticScope { return nil }

func (s *typeAlgebraTestScope) GetSlot(name string) StaticSlot {
	if slot := s.slots.Get(name); slot != nil {
		return slot
	}
	return nil
}

func (s *typeAlgebraTestScope) GetOwnSlot(name string) StaticSlot {
	if slot := s.slots.Get(name); slot != nil {
		return slot
	}
	return nil
}

func (s *typeAlgebraTestScope) GetTopmostScopeOfEventualDeclaration(name string) StaticScope {
	return StaticScope_GetTopmostScopeOfEventualDeclaration(s, name)
}

func (s *typeAlgebraTestScope) GetTypeOfThis() JSType { return nil }

func (s *typeAlgebraTestScope) LookupQualifiedName(qname QualifiedName) JSType {
	return StaticTypedScope_LookupQualifiedName(s, qname)
}

func newTypeAlgebraDump(asts map[string]string) *typeAlgebraDump {
	d := &typeAlgebraDump{
		out: &strings.Builder{}, seed: 0x5DEECE66D, asts: asts,
		inferMap: jutil.NewLinkedHashMap[*TemplateType, JSType](),
	}
	d.registry = NewJSTypeRegistryWithForwardDeclaredTypes(
		&typeAlgebraTestReporter{sb: d.out}, jutil.NewLinkedHashSet("Forwarded"))
	d.scope = &typeAlgebraTestScope{
		root:  NewNode(Token_ROOT),
		slots: jutil.NewLinkedHashMap[string, *SimpleSlotBase](),
	}
	return d
}

func (d *typeAlgebraDump) nat(n JSTypeNative) JSType {
	return d.registry.GetNativeType(n)
}

func (d *typeAlgebraDump) natObj(n JSTypeNative) ObjectType {
	return d.registry.GetNativeObjectType(n)
}

func (d *typeAlgebraDump) add(label string, t JSType) {
	d.labels = append(d.labels, label)
	d.types = append(d.types, t)
	d.line("add\t" + strconv.Itoa(len(d.types)) + "\t" + label)
}

func taAdd[T JSType](d *typeAlgebraDump, label string, t T) T {
	d.add(label, t)
	return t
}

func taFn(name string) Node {
	return IR_Function(IR_Name(name), IR_ParamList(), IR_Block())
}

// taSexp serializes a type AST: (TOKEN[ "string"] line:char child...).
func taSexp(n Node) string {
	var sb strings.Builder
	taSexpTo(n, &sb)
	return sb.String()
}

func taSexpTo(n Node, sb *strings.Builder) {
	sb.WriteByte('(')
	sb.WriteString(n.GetToken().Name())
	if n.nodeKind() == node_kindString {
		sb.WriteString(" \"")
		sb.WriteString(strings.ReplaceAll(taEsc(n.GetString()), "\"", "\\\""))
		sb.WriteByte('"')
	}
	sb.WriteByte(' ')
	sb.WriteString(strconv.Itoa(n.GetLineno()))
	sb.WriteByte(':')
	sb.WriteString(strconv.Itoa(n.GetCharno()))
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		sb.WriteByte(' ')
		taSexpTo(c, sb)
	}
	sb.WriteByte(')')
}

// taUnesc reverses taEsc.
func taUnesc(s string) string {
	var chars []uint16
	src := jutil.StringToUTF16(s)
	for i := 0; i < len(src); i++ {
		c := src[i]
		if c == '\\' && i+1 < len(src) {
			switch src[i+1] {
			case '\\':
				chars = append(chars, '\\')
				i++
				continue
			case 'n':
				chars = append(chars, '\n')
				i++
				continue
			case 'u':
				v, err := strconv.ParseUint(jutil.StringFromUTF16(src[i+2:i+6]), 16, 16)
				if err != nil {
					panic(err)
				}
				chars = append(chars, uint16(v))
				i += 5
				continue
			}
		}
		chars = append(chars, c)
	}
	return jutil.StringFromUTF16(chars)
}

var taTokensByName = func() map[string]Token {
	m := map[string]Token{}
	for _, t := range Token_values() {
		m[t.Name()] = t
	}
	return m
}()

// taParseSexp rebuilds the AST printed by taSexp / TypeAlgebraDump.sexp.
func taParseSexp(s string) Node {
	p := 0
	var parse func() Node
	expect := func(c byte) {
		if s[p] != c {
			panic(fmt.Sprintf("sexp: expected %q at %d in %s", c, p, s))
		}
		p++
	}
	parse = func() Node {
		expect('(')
		start := p
		for s[p] != ' ' {
			p++
		}
		tok, ok := taTokensByName[s[start:p]]
		if !ok {
			panic("sexp: unknown token " + s[start:p])
		}
		expect(' ')
		var n Node
		if s[p] == '"' {
			p++
			var sb strings.Builder
			for s[p] != '"' {
				if s[p] == '\\' && s[p+1] == '"' {
					sb.WriteByte('"')
					p += 2
					continue
				}
				if s[p] == '\\' {
					sb.WriteByte(s[p])
					sb.WriteByte(s[p+1])
					p += 2
					continue
				}
				sb.WriteByte(s[p])
				p++
			}
			p++
			expect(' ')
			n = Node_NewStringWithToken(tok, taUnesc(sb.String()))
		} else {
			n = NewNode(tok)
		}
		start = p
		for s[p] != ' ' && s[p] != ')' {
			p++
		}
		lc := strings.SplitN(s[start:p], ":", 2)
		line, _ := strconv.Atoi(lc[0])
		char, _ := strconv.Atoi(lc[1])
		n.SetLinenoCharno(line, char)
		for s[p] == ' ' {
			p++
			n.AddChildToBack(parse())
		}
		expect(')')
		return n
	}
	return parse()
}

func (d *typeAlgebraDump) parsed(typeString string, s StaticTypedScope) JSType {
	sexp, ok := d.asts[typeString]
	if !ok {
		panic("no AST for type string " + typeString)
	}
	n := taParseSexp(sexp)
	d.line("ast\t" + taEsc(typeString) + "\t" + taEsc(taSexp(n)))
	return d.registry.CreateTypeFromCommentNodeWithScope(n, "src.js", s)
}

func taPropertyTypeMap(o ObjectType) *jutil.LinkedHashMap[string, JSType] {
	propTypeMap := jutil.NewLinkedHashMap[string, JSType]()
	if fn, ok := o.(FunctionType); ok {
		var cacheBuf [8]FunctionType
		taUpdateFnPropertyTypeMap(fn, propTypeMap, cacheBuf[:0])
		return propTypeMap
	}
	for _, name := range o.GetPropertyNames() {
		if propTypeMap.ContainsKey(name) {
			panic(jutil.NewIllegalArgumentExceptionMsg("Multiple entries with same key: " + name))
		}
		propTypeMap.Put(name, o.GetPropertyType(name))
	}
	return propTypeMap
}

func taUpdateFnPropertyTypeMap(type_ FunctionType, propTypeMap *jutil.LinkedHashMap[string, JSType], cache []FunctionType) {
	if type_ == nil {
		return
	}
	if prototype := type_.GetPrototype(); prototype != nil {
		for _, name := range prototype.GetOwnPropertyNames() {
			if !propTypeMap.ContainsKey(name) {
				propTypeMap.Put(name, prototype.GetPropertyType(name))
			}
		}
	}
	for _, interfaceType := range type_.GetExtendedInterfaces() {
		superConstructor := interfaceType.GetConstructor()
		if superConstructor == nil || fnSliceContains(cache, superConstructor) {
			continue
		}
		taUpdateFnPropertyTypeMap(superConstructor, propTypeMap, append(cache, superConstructor))
	}
}

func taRecord(keysAndValues ...any) *jutil.LinkedHashMap[string, JSType] {
	m := jutil.NewLinkedHashMap[string, JSType]()
	for i := 0; i < len(keysAndValues); i += 2 {
		m.Put(keysAndValues[i].(string), keysAndValues[i+1].(JSType))
	}
	return m
}

// ---------------------------------------------------------------------
// The type universe

func (d *typeAlgebraDump) build() {
	// Like BaseJSTypeTestCase: define everything with the resolver open, so
	// that types are resolved together when it closes.
	closer := d.registry.GetResolver().OpenForDefinition()
	d.buildTypes()
	closer.Close()
}

func (d *typeAlgebraDump) buildTypes() {
	registry := d.registry
	scope := d.scope
	s := jutil.StringPtr
	for _, n := range JSTypeNative_values() {
		d.add("native:"+n.String(), d.nat(n))
	}
	number := d.nat(JSTypeNative_NUMBER_TYPE)
	string_ := d.nat(JSTypeNative_STRING_TYPE)
	bool_ := d.nat(JSTypeNative_BOOLEAN_TYPE)
	nul := d.nat(JSTypeNative_NULL_TYPE)
	undef := d.nat(JSTypeNative_VOID_TYPE)
	unknown := d.nat(JSTypeNative_UNKNOWN_TYPE)
	all := d.nat(JSTypeNative_ALL_TYPE)
	symbol := d.nat(JSTypeNative_SYMBOL_TYPE)
	array := d.natObj(JSTypeNative_ARRAY_TYPE)
	object := d.natObj(JSTypeNative_OBJECT_TYPE)
	readonlyArray := d.natObj(JSTypeNative_READONLY_ARRAY_TYPE)
	iterable := d.natObj(JSTypeNative_ITERABLE_TYPE)
	promise := d.natObj(JSTypeNative_PROMISE_TYPE)
	ithenable := d.natObj(JSTypeNative_I_THENABLE_TYPE)
	iobject := d.natObj(JSTypeNative_I_OBJECT_TYPE)
	map_ := d.natObj(JSTypeNative_MAP_TYPE)

	// Classes, interfaces and inheritance.
	fooCtor := registry.CreateConstructorType(
		s("Foo"), taFn("Foo"), registry.CreateParameters(number), nil, nil, false)
	foo := fooCtor.GetInstanceType()
	fooCtor.GetPrototype().DefineDeclaredProperty("bar", number, NilNode)
	foo.DefineDeclaredProperty("baz", string_, NilNode)
	fooCtor.DefineDeclaredProperty("staticProp", bool_, NilNode)
	registry.DeclareType(scope, "Foo", foo)
	d.add("Foo", foo)
	d.add("ctor:Foo", fooCtor)
	d.add("Foo.prototype", fooCtor.GetPrototype())

	subFooCtor := registry.CreateConstructorType(
		s("SubFoo"), taFn("SubFoo"), registry.CreateParameters(), nil, nil, false)
	subFooCtor.SetPrototypeBasedOn(foo)
	subFoo := subFooCtor.GetInstanceType()
	subFoo.DefineDeclaredProperty("qux", bool_, NilNode)
	registry.DeclareType(scope, "SubFoo", subFoo)
	d.add("SubFoo", subFoo)
	d.add("ctor:SubFoo", subFooCtor)
	d.add("SubFoo.prototype", subFooCtor.GetPrototype())

	iCtor := FunctionType_Builder_(registry).ForInterface().WithName(s("I")).WithSourceNode(taFn("I")).Build()
	iCtor.GetPrototype().DefineDeclaredProperty("bar", number, NilNode)
	iface := iCtor.GetInstanceType()
	registry.DeclareType(scope, "I", iface)
	d.add("I", iface)
	d.add("ctor:I", iCtor)

	jCtor := FunctionType_Builder_(registry).ForInterface().WithName(s("J")).WithSourceNode(taFn("J")).Build()
	jCtor.SetExtendedInterfaces([]ObjectType{iface})
	jCtor.GetPrototype().DefineDeclaredProperty("jay", string_, NilNode)
	jface := jCtor.GetInstanceType()
	d.add("J", jface)
	d.add("ctor:J", jCtor)

	implCtor := registry.CreateConstructorType(
		s("Impl"), taFn("Impl"), registry.CreateParameters(), nil, nil, false)
	implCtor.SetImplementedInterfaces([]ObjectType{jface})
	implCtor.GetPrototype().DefineDeclaredProperty("bar", number, NilNode)
	impl := implCtor.GetInstanceType()
	registry.DeclareType(scope, "Impl", impl)
	d.add("Impl", impl)

	recCtor := FunctionType_Builder_(registry).ForInterface().WithName(s("Rec")).WithSourceNode(taFn("Rec")).Build()
	recCtor.SetImplicitMatch(true)
	recCtor.GetPrototype().DefineDeclaredProperty("bar", number, NilNode)
	recCtor.GetPrototype().DefineDeclaredProperty("opt", registry.CreateOptionalType(string_), NilNode)
	rec := recCtor.GetInstanceType()
	registry.DeclareType(scope, "Rec", rec)
	d.add("Rec", rec)
	d.add("ctor:Rec", recCtor)

	tKey := registry.CreateTemplateType("T")
	d.tKey = tKey
	boxCtor := FunctionType_Builder_(registry).
		ForConstructor().
		WithName(s("Box")).
		WithSourceNode(taFn("Box")).
		WithParametersOfList(registry.CreateParameters(tKey)).
		WithTemplateKeys(tKey).
		Build()
	boxCtor.GetPrototype().DefineDeclaredProperty("value", tKey, NilNode)
	box := boxCtor.GetInstanceType()
	registry.DeclareType(scope, "Box", box)
	d.add("Box", box)
	d.add("ctor:Box", boxCtor)
	d.add("T", tKey)
	boxString := taAdd(d, "Box<string>", registry.CreateTemplatizedType(box, string_))
	boxNumber := taAdd(d, "Box<number>", registry.CreateTemplatizedType(box, number))
	d.add("Box<?>", registry.CreateTemplatizedType(box, unknown))
	d.add("Box<string>.value", boxString.GetPropertyType("value"))

	structCtor := registry.CreateConstructorType(
		s("Struct"), taFn("Struct"), registry.CreateParameters(), nil, nil, false)
	structCtor.SetStruct()
	structCtor.GetPrototype().DefineDeclaredProperty("s", number, NilNode)
	d.add("Struct", structCtor.GetInstanceType())
	dictCtor := registry.CreateConstructorType(
		s("Dict"), taFn("Dict"), registry.CreateParameters(), nil, nil, false)
	dictCtor.SetDict()
	d.add("Dict", dictCtor.GetInstanceType())
	abstractCtor := registry.CreateConstructorType(
		s("Abs"), taFn("Abs"), registry.CreateParameters(), nil, nil, true)
	abstractCtor.SetImplementedInterfaces([]ObjectType{iface})
	d.add("Abs", abstractCtor.GetInstanceType())
	subStructCtor := registry.CreateConstructorType(
		s("SubStruct"), taFn("SubStruct"), registry.CreateParameters(), nil, nil, false)
	subStructCtor.SetPrototypeBasedOn(structCtor.GetInstanceType())
	d.add("SubStruct", subStructCtor.GetInstanceType())
	anonCtor := registry.CreateConstructorType(s(""), NilNode, registry.CreateParameters(), nil, nil, false)
	d.add("anonCtor", anonCtor)
	d.add("anonInstance", anonCtor.GetInstanceType())
	structuralCtor := FunctionType_Builder_(registry).ForConstructor().WithTypeOfThis(foo).Build()
	d.add("structuralCtor", structuralCtor)

	lit := registry.CreateAnonymousObjectType(nil)
	lit.DefineInferredProperty("a", number, NilNode)
	lit.DefineInferredProperty("b", string_, NilNode)
	d.add("objlit", lit)
	emptyLit := registry.CreateAnonymousObjectType(nil)
	d.add("emptyObjlit", emptyLit)
	named := registry.CreateObjectType(s("Named"), object)
	named.DefineDeclaredProperty("n", number, NilNode)
	d.add("Named", named)

	// Records.
	recA := taAdd(d, "{a:number}", registry.CreateRecordType(taRecord("a", number)))
	recAB := taAdd(d, "{a:number,b:string}", registry.CreateRecordType(taRecord("b", string_, "a", number)))
	d.add("{a:number=}", registry.CreateRecordType(taRecord("a", registry.CreateOptionalType(number))))
	recAS := taAdd(d, "{a:string}", registry.CreateRecordType(taRecord("a", string_)))
	d.add("{bar:number}", registry.CreateRecordType(taRecord("bar", number)))
	d.add("{r:{a:number},f:...}", registry.CreateRecordType(
		taRecord("r", recA, "f", JSType(registry.CreateFunctionType(number, string_)))))
	rMany := jutil.NewLinkedHashMap[string, JSType]()
	for i := 0; i < 12; i++ {
		if i%2 == 0 {
			rMany.Put("p"+string(rune('a'+i)), number)
		} else {
			rMany.Put("p"+string(rune('a'+i)), string_)
		}
	}
	d.add("{many}", registry.CreateRecordType(rMany))
	synth := NewRecordTypeBuilder(registry)
	synth.SetSynthesized(true)
	synth.AddProperty("x", bool_, NilNode)
	d.add("{x:boolean}(synthetic)", synth.Build())
	d.add("emptyRecord", NewRecordTypeBuilder(registry).Build())
	d.add("buildRecordTypeFromObject(objlit)", registry.BuildRecordTypeFromObject(lit))

	// Unions.
	d.add("number|string", registry.CreateUnionType(number, string_))
	d.add("number|null", registry.CreateNullableType(number))
	d.add("Foo|null", registry.CreateNullableType(foo))
	d.add("Foo|SubFoo", registry.CreateUnionType(foo, subFoo))
	d.add("SubFoo|Foo", registry.CreateUnionType(subFoo, foo))
	d.add(
		"(number|string)|boolean",
		registry.CreateUnionType(registry.CreateUnionType(number, string_), bool_))
	d.add("?|number", registry.CreateUnionType(unknown, number))
	d.add("*|number", registry.CreateUnionType(all, number))
	d.add("?|undefined", registry.CreateUnionType(unknown, undef))
	d.add("number|undefined", registry.CreateOptionalType(number))
	d.add("Foo|null|undefined", registry.CreateOptionalNullableType(foo))
	d.add(
		"Array<string>|Array<number>",
		registry.CreateUnionType(
			registry.CreateTemplatizedType(array, string_),
			registry.CreateTemplatizedType(array, number)))
	d.add(
		"ReadonlyArray<string>|ReadonlyArray<number>",
		registry.CreateUnionType(
			registry.CreateTemplatizedType(readonlyArray, string_),
			registry.CreateTemplatizedType(readonlyArray, number)))
	d.add("Array|Array<string>", registry.CreateUnionType(array, registry.CreateTemplatizedType(array, string_)))
	d.add("Object|Array<string>", registry.CreateUnionType(object, registry.CreateTemplatizedType(array, string_)))
	d.add("{a:number}|{a:number,b:string}", registry.CreateUnionType(recA, recAB))
	d.add("{a:number}|{a:string}", registry.CreateUnionType(recA, recAS))
	d.add(
		"fn|fn",
		registry.CreateUnionType(
			registry.CreateFunctionType(number, string_), registry.CreateFunctionType(string_, string_)))
	d.add(
		"fn(number)|fn(string)",
		registry.CreateUnionType(
			registry.CreateFunctionType(number, number), registry.CreateFunctionType(number, string_)))
	d.add("null|undefined|number|string|boolean", registry.CreateUnionType(nul, undef, number, string_, bool_))
	d.add("symbol|number", registry.CreateUnionType(symbol, number))
	d.add("Foo|Impl|I", registry.CreateUnionType(foo, impl, iface))
	d.add("T|number", registry.CreateUnionType(tKey, number))
	d.add("Box<string>|Box<number>", registry.CreateUnionType(boxString, boxNumber))
	d.add("emptyUnion", UnionType_Builder_(registry).Build())
	d.add("checked?|?", registry.CreateUnionType(d.nat(JSTypeNative_CHECKED_UNKNOWN_TYPE), unknown))
	d.add("checked?|number", registry.CreateUnionType(d.nat(JSTypeNative_CHECKED_UNKNOWN_TYPE), number))
	big := UnionType_Builder_(registry)
	for i := 0; i < 35; i++ {
		big.AddAlternate(registry.CreateRecordType(taRecord("k"+strconv.Itoa(i), number)))
	}
	d.add("bigUnion", big.Build())

	// Functions.
	d.add("function(number):string", registry.CreateFunctionType(string_, number))
	d.add(
		"function(number=):string",
		registry.CreateFunctionTypeWithParameters(string_, registry.CreateOptionalParameters(number)))
	d.add("function(...number):undefined", registry.CreateFunctionTypeWithVarArgs(undef, number))
	d.add("function():?", registry.CreateFunctionType(unknown))
	pb := NewFunctionParamBuilder(registry)
	pb.AddRequiredParams(number)
	pb.AddOptionalParams(string_)
	pb.AddVarArgs(bool_)
	d.add(
		"function(number,string=,...boolean):Foo",
		FunctionType_Builder_(registry).WithParametersOfList(pb.Build()).WithReturnType(foo).Build())
	d.add(
		"function(this:Foo):undefined",
		FunctionType_Builder_(registry).WithParameters().WithReturnType(undef).WithTypeOfThis(foo).Build())
	d.add(
		"function(new:Foo):Foo",
		FunctionType_Builder_(registry).
			ForConstructor().
			WithParameters().
			WithReturnType(foo).
			WithTypeOfThis(foo).
			Build())
	d.add("function(?):number", registry.CreateFunctionType(number, unknown))
	d.add("function(T):T", FunctionType_Builder_(registry).
		WithParametersOfList(registry.CreateParameters(tKey)).
		WithReturnType(tKey).
		WithTemplateKeys(tKey).
		Build())
	d.add(
		"function(number):string(inferred)",
		FunctionType_Builder_(registry).
			WithParametersOfList(registry.CreateParameters(number)).
			WithInferredReturnType(string_).
			Build())
	d.add(
		"function(this:Foo, number):number",
		registry.CreateFunctionTypeWithInstanceType(foo, number, []JSType{number}))
	withProps := registry.CreateFunctionType(number, string_)
	withProps.DefineDeclaredProperty("extra", bool_, NilNode)
	d.add("fnWithProps", withProps)
	d.add("bind(-1)", fooCtor.GetBindReturnType(-1))
	d.add("bind(1)", registry.CreateFunctionType(number, string_, bool_).GetBindReturnType(1))
	d.add("bind(2)", registry.CreateFunctionType(number, string_, bool_).GetBindReturnType(2))
	d.add("forget(Foo)", fooCtor.ForgetParameterAndReturnTypes())
	d.add(
		"newReturnType",
		registry.CreateFunctionTypeWithNewReturnType(registry.CreateFunctionType(number, string_), bool_))
	d.add("toBuilder(SubFoo)", subFooCtor.ToBuilder().Build())

	// Enums.
	eNum := EnumType_Builder_(registry).SetName(s("E")).SetElementType(number).SetSource(NilNode).Build()
	eNum.DefineElement("A", NilNode)
	eNum.DefineElement("B", NilNode)
	registry.DeclareType(scope, "E", eNum.GetElementsType())
	registry.IdentifyNonNullableName(scope, "E")
	d.add("enum E", eNum)
	d.add("E", eNum.GetElementsType())
	eStr := EnumType_Builder_(registry).SetName(s("S")).SetElementType(string_).Build()
	eStr.DefineElement("X", NilNode)
	d.add("enum S", eStr)
	d.add("S", eStr.GetElementsType())
	eFoo := EnumType_Builder_(registry).SetName(s("EF")).SetElementType(foo).Build()
	eFoo.DefineElement("F", NilNode)
	d.add("enum EF", eFoo)
	d.add("EF", eFoo.GetElementsType())
	eUnion := EnumType_Builder_(registry).
		SetName(s("EU")).
		SetElementType(registry.CreateUnionType(number, string_)).
		Build()
	d.add("EU", eUnion.GetElementsType())
	d.add("E|S", registry.CreateUnionType(eNum.GetElementsType(), eStr.GetElementsType()))

	// Templatized and template types.
	d.add("Array<string>", registry.CreateTemplatizedType(array, string_))
	d.add("Array<number>", registry.CreateTemplatizedType(array, number))
	d.add("Array<?>", registry.CreateTemplatizedType(array, unknown))
	d.add("Array<Foo>", registry.CreateTemplatizedType(array, foo))
	d.add("Array<SubFoo>", registry.CreateTemplatizedType(array, subFoo))
	d.add("Object<string,number>", registry.CreateTemplatizedType(object, string_, number))
	d.add("IObject<number,string>", registry.CreateTemplatizedType(iobject, number, string_))
	d.add("Iterable<string>", registry.CreateTemplatizedType(iterable, string_))
	d.add("ReadonlyArray<string>", registry.CreateTemplatizedType(readonlyArray, string_))
	d.add("ReadonlyArray<Foo>", registry.CreateTemplatizedType(readonlyArray, foo))
	d.add("Promise<number>", registry.CreateTemplatizedType(promise, number))
	d.add("IThenable<string>", registry.CreateTemplatizedType(ithenable, string_))
	d.add("Map<string,number>", registry.CreateTemplatizedType(map_, string_, number))
	d.add("Array<Array<string>>", registry.CreateTemplatizedType(array, registry.CreateTemplatizedType(array, string_)))
	uKey := taAdd(d, "U extends number", registry.CreateTemplateTypeWithBound("U", number))
	d.add("Box<T>", registry.CreateTemplatizedType(box, tKey))
	tmap := jutil.NewLinkedHashMap[*TemplateType, JSType]()
	tmap.Put(tKey, bool_)
	d.add("Box<boolean>(map)", registry.CreateTemplatizedTypeOfMap(box, tmap))
	replacer := TemplateTypeReplacer_ForInference(registry, tmap)
	d.add("replace(function(T):T)", replacer.VisitJSType(d.types[d.indexOfLabel("function(T):T")]))
	d.add("replace(Box<T>)", replacer.VisitJSType(d.types[d.indexOfLabel("Box<T>")]))
	d.add("replace(T|number)", replacer.VisitJSType(d.types[d.indexOfLabel("T|number")]))
	d.add("U|string", registry.CreateUnionType(uKey, string_))

	// Proxies and named types.
	d.add("proxy(number)", NewProxyObjectType(registry, number))
	d.add("proxy(Foo)", NewProxyObjectType(registry, foo))
	d.add("proxy(number|string)", NewProxyObjectType(registry, registry.CreateUnionType(number, string_)))
	d.add("named:Foo", registry.CreateNamedType(scope, "Foo", s("src.js"), 1, 2))
	d.add("named:Missing", registry.CreateNamedType(scope, "Missing", s("src.js"), 3, 4))
	d.add("named:Forwarded", registry.CreateNamedType(scope, "Forwarded", s("src.js"), 5, 6))
	d.add("named:Box", registry.CreateNamedType(scope, "Box", s("src.js"), 7, 8))
	d.add("noResolved", NewNoResolvedType(registry, "NR", []JSType{}))
	iter := taAdd(d, "knownSymbol", NewKnownSymbolType(registry, "Symbol.iterator"))
	withSym := registry.CreateAnonymousObjectType(nil)
	withSym.DefineDeclaredPropertyByKey(NewProperty_SymbolKey(iter), number, NilNode)
	d.add("objlitWithSymbol", withSym)

	// More generics, structural interfaces and forward references.
	kKey := registry.CreateTemplateType("K")
	vKey := registry.CreateTemplateType("V")
	pairCtor := FunctionType_Builder_(registry).
		ForInterface().
		WithName(s("Pair")).
		WithSourceNode(taFn("Pair")).
		WithTemplateKeys(kKey, vKey).
		Build()
	pairCtor.SetImplicitMatch(true)
	pairCtor.GetPrototype().DefineDeclaredProperty("k", kKey, NilNode)
	pairCtor.GetPrototype().DefineDeclaredProperty("v", vKey, NilNode)
	pairType := pairCtor.GetInstanceType()
	registry.DeclareType(scope, "Pair", pairType)
	d.add("Pair", pairType)
	pairSN := taAdd(d, "Pair<string,number>", registry.CreateTemplatizedType(pairType, string_, number))
	d.add("Pair<number,number>", registry.CreateTemplatizedType(pairType, number, number))
	d.add("Pair<K,V>", registry.CreateTemplatizedType(pairType, kKey, vKey))
	d.add("{k:string,v:number}", registry.CreateRecordType(taRecord("k", string_, "v", number)))
	d.add("{k:string}", registry.CreateRecordType(taRecord("k", string_)))

	strBoxCtor := registry.CreateConstructorType(
		s("StrBox"), taFn("StrBox"), registry.CreateParameters(), nil, nil, false)
	strBoxCtor.SetPrototypeBasedOn(boxString)
	registry.DeclareType(scope, "StrBox", strBoxCtor.GetInstanceType())
	d.add("StrBox", strBoxCtor.GetInstanceType())
	d.add("ctor:StrBox", strBoxCtor)

	elKey := registry.CreateTemplateType("EL")
	listCtor := FunctionType_Builder_(registry).
		ForInterface().
		WithName(s("List")).
		WithSourceNode(taFn("List")).
		WithTemplateKeys(elKey).
		Build()
	listCtor.SetExtendedInterfaces([]ObjectType{registry.CreateTemplatizedType(iterable, elKey)})
	listCtor.GetPrototype().DefineDeclaredProperty("get", registry.CreateFunctionType(elKey, number), NilNode)
	registry.DeclareType(scope, "List", listCtor.GetInstanceType())
	d.add("List", listCtor.GetInstanceType())
	d.add("ctor:List", listCtor)
	d.add("List<string>", registry.CreateTemplatizedType(listCtor.GetInstanceType(), string_))
	d.add("Box<Box<string>>", registry.CreateTemplatizedType(box, boxString))
	d.add("Box<T>|null", registry.CreateNullableType(registry.CreateTemplatizedType(box, tKey)))
	d.add(
		"function(K,V):Pair<K,V>",
		FunctionType_Builder_(registry).
			WithParametersOfList(registry.CreateParameters(kKey, vKey)).
			WithReturnType(registry.CreateTemplatizedType(pairType, kKey, vKey)).
			WithTemplateKeys(kKey, vKey).
			Build())
	d.add(
		"function(this:T):T",
		FunctionType_Builder_(registry).
			WithParameters().
			WithReturnType(tKey).
			WithTypeOfThis(tKey).
			WithTemplateKeys(tKey).
			Build())
	d.add(
		"function(function(T):T, ...T):Array<T>",
		FunctionType_Builder_(registry).
			WithParametersOfList(
				registry.CreateParametersWithVarArgs(
					d.types[d.indexOfLabel("function(T):T")], tKey)).
			WithReturnType(registry.CreateTemplatizedType(array, tKey)).
			WithTemplateKeys(tKey).
			Build())
	d.add(
		"function(new:Box<T>, T)",
		FunctionType_Builder_(registry).
			ForConstructor().
			WithParametersOfList(registry.CreateParameters(tKey)).
			WithTypeOfThis(registry.CreateTemplatizedType(box, tKey)).
			WithTemplateKeys(tKey).
			Build())
	d.add("Pair<string,number>|{k:string}", registry.CreateUnionType(pairSN, d.types[d.indexOfLabel("{k:string}")]))
	d.add("T|K|V", registry.CreateUnionType(tKey, kKey, vKey))

	// Forward references, resolved when the resolver closes.
	late := registry.CreateNamedType(scope, "Late", s("src.js"), 9, 10)
	late.DefineDeclaredProperty("lateProp", number, NilNode)
	d.add("named:Late", late)
	lateMissing := registry.CreateNamedType(scope, "LateMissing", s("src.js"), 11, 12)
	lateMissing.DefineDeclaredProperty("p", number, NilNode)
	d.add("named:LateMissing", lateMissing)
	d.add("Box<Late>", registry.CreateTemplatizedType(box, late))
	d.add("Late|null", registry.CreateNullableType(late))
	lateCtor := registry.CreateConstructorType(
		s("Late"), taFn("Late"), registry.CreateParameters(), nil, nil, false)
	registry.DeclareType(scope, "Late", lateCtor.GetInstanceType())
	d.add("Late", lateCtor.GetInstanceType())

	// The property index behind canPropertyBeDefined.
	registry.RegisterPropertyOnType("bar", foo)
	registry.RegisterPropertyOnType("bar", rec)
	registry.RegisterPropertyOnType("bar", registry.CreateUnionType(impl, nul))
	registry.RegisterPropertyOnType("a", recA)
	registry.RegisterPropertyOnType("a", lit)
	registry.RegisterPropertyOnType("k", pairType)
	registry.RegisterPropertyOnType("zzz", d.types[d.indexOfLabel("{x:boolean}(synthetic)")])
	registry.RegisterPropertyOnType("zzz", lateCtor.GetInstanceType())

	// Bindings for the replacer facts.
	d.partialMap = boxString.GetTemplateTypeMap()
	d.totalMap = pairSN.GetTemplateTypeMap()
	d.inferMap.Put(tKey, number)
	d.inferMap.Put(kKey, string_)

	// Types from JSDoc type expressions.
	scope.slots.Put("ns", NewSimpleSlot("ns", lit, false))
	typeStrings := []string{
		"number",
		"!Array<string>",
		"Array<string>",
		"?Foo",
		"Foo",
		"!Foo",
		"{a: number, b: (string|undefined)}",
		"{a: number, 'quoted': string}",
		"function(this:Foo, number=, ...string): boolean",
		"function(new:Foo, number)",
		"function(new:number)",
		"function(...?)",
		"function(number=, string)",
		"Object<string, number>",
		"Object<number>",
		"(number|string|null)",
		"*",
		"?",
		"undefined",
		"void",
		"Box<string>",
		"Box<string, number>",
		"!Box",
		"Missing",
		"Forwarded",
		"E",
		"!E",
		"?E",
		"Rec",
		"SubFoo",
		"typeof ns",
		"!Map<string, !Array<?>>",
		"function(): undefined",
		"string=",
		"Pair<string, number>",
		"StrBox",
		"List<number>",
		"Late",
		"Box<Box<string>>",
		"{k: string, v: number}",
		"function(this:StrBox): string",
		"LateMissing",
	}
	for _, ts := range typeStrings {
		d.add("jsdoc:"+ts, d.parsed(ts, scope))
	}
	tScope := registry.CreateScopeWithTemplates(scope, []*TemplateType{tKey})
	templateTypeStrings := []string{
		"T", "Box<T>", "function(T): T", "(Array<T>|T)", "Pair<T, string>", "K",
	}
	for _, ts := range templateTypeStrings {
		d.add("jsdocT:"+ts, d.parsed(ts, tScope))
	}
}

func (d *typeAlgebraDump) indexOfLabel(label string) int {
	for i, l := range d.labels {
		if l == label {
			return i
		}
	}
	return -1
}

// ---------------------------------------------------------------------
// Facts

var taProps = []string{
	"toString", "valueOf", "length", "prototype", "call", "apply", "bind", "bar", "baz", "value", "a",
	"constructor", "then", "missing",
}

func taBits(bs ...bool) string {
	var sb strings.Builder
	for _, b := range bs {
		if b {
			sb.WriteByte('1')
		} else {
			sb.WriteByte('0')
		}
	}
	return sb.String()
}

// taPtr maps a nil pointer (Java null) to nil.
func taPtr[T any](p *T) any {
	if p == nil {
		return nil
	}
	return p
}

// taNullableList maps a nil slice (Java null) to nil.
func taNullableList[T any](l []T) any {
	if l == nil {
		return nil
	}
	return l
}

func (d *typeAlgebraDump) unaryFacts(i int) {
	t := d.types[i]
	registry := d.registry
	d.line("=== " + strconv.Itoa(i+1) + " " + d.labels[i])
	d.fact("toString", func() any { return t.String() })
	d.fact("annot", func() any { return t.ToAnnotationString(JSType_Nullability_EXPLICIT) })
	d.fact("annotImplicit", func() any { return t.ToAnnotationString(JSType_Nullability_IMPLICIT) })
	d.fact("typeClass", func() any { return t.GetTypeClass() })
	d.fact("preds", func() any {
		return taBits(
			t.IsNoType(), t.IsNoResolvedType(), t.IsNoObjectType(), t.IsEmptyType(),
			t.IsNumberObjectType(), t.IsNumberValueType(), t.IsBigIntObjectType(),
			t.IsBigIntValueType(), t.IsFunctionPrototypeType(), t.IsStringObjectType(),
			t.IsSymbolObjectType(), t.IsTheObjectType(), t.IsStringValueType(),
			t.IsSymbolValueType(), t.IsKnownSymbolValueType(), t.IsString(), t.IsNumber(),
			t.IsSymbol(), t.IsBigIntOrNumber(), t.IsOnlyBigInt(), t.IsArrayType(),
			t.IsReadonlyArrayType(), t.IsBooleanObjectType(), t.IsBooleanValueType(),
			t.IsRegexpType(), t.IsDateType(), t.IsNullType(), t.IsVoidType(), t.IsAllType(),
			t.IsUnknownType(), t.IsSomeUnknownType(), t.IsCheckedUnknownType(),
			t.IsUnionType(), t.IsRawTypeOfTemplatizedType(), t.IsLiteralObject(),
			t.IsGlobalThisType(), t.IsFunctionType(), t.IsEnumElementType(), t.IsEnumType(),
			t.IsNamedType(), t.IsRecordType(), t.IsStructuralInterface(),
			t.IsStructuralType(), t.IsTemplatizedType(), t.IsTemplateType(),
			t.HasAnyTemplateTypes(), t.IsObject(), t.IsObjectType(), t.IsConstructor(),
			t.IsNominalType(), t.IsNominalConstructorOrInterface(), t.IsNativeObjectType(),
			t.LoosenTypecheckingDueToForwardReferencedSupertype(), t.IsInstanceType(),
			t.IsInterface(), t.IsOrdinaryFunction())
	})
	d.fact("preds2", func() any {
		return taBits(
			t.MatchesNumberContext(), t.MatchesStringContext(), t.MatchesSymbolContext(),
			t.MatchesObjectContext(), t.CanBeCalled(), t.IsBoxableScalar(), t.IsNullable(),
			t.IsVoidable(), t.IsExplicitlyVoidable(), t.IsResolved(),
			t.IsSuccessfullyResolved(), t.IsUnsuccessfullyResolved(), t.HasDisplayName())
	})
	d.fact("isStruct", func() any { return t.IsStruct() })
	d.fact("isDict", func() any { return t.IsDict() })
	d.fact("displayName", func() any { return t.GetDisplayName() })
	d.fact("autoboxesTo", func() any { return t.AutoboxesTo() })
	d.fact("autobox", func() any { return t.Autobox() })
	d.fact("dereference", func() any { return t.Dereference() })
	d.fact("restrictByNotNullOrUndefined", func() any { return t.RestrictByNotNullOrUndefined() })
	d.fact("restrictByNotUndefined", func() any { return t.RestrictByNotUndefined() })
	d.fact("restrictByNotNull", func() any { return t.RestrictByNotNull() })
	d.fact("collapseUnion", func() any { return t.CollapseUnion() })
	d.fact("toBoolean", func() any { return t.GetPossibleToBooleanOutcomes() })
	for _, o := range Outcome_values() {
		d.fact("givenOutcome:"+o.String(), func() any { return t.GetRestrictedTypeGivenOutcome(o) })
	}
	d.fact("templateParamCount", func() any { return t.GetTemplateParamCount() })
	d.fact("typeParameters", func() any { return t.GetTypeParameters() })
	d.fact("templateTypeMap", func() any { return t.GetTemplateTypeMap().String() })
	d.fact("enumeratedTypeOfEnumElement", func() any { return t.GetEnumeratedTypeOfEnumElement() })
	d.fact("unionMembers", func() any { return taNullableList(t.GetUnionMembers()) })
	for _, p := range taProps {
		d.fact("findPropertyType:"+p, func() any { return t.FindPropertyType(p) })
		d.fact("propertyKind:"+p, func() any {
			return t.GetPropertyKind(p).String() + "/" + t.GetPropertyKindWithAutobox(p, false).String()
		})
	}
	for _, p := range []string{"bar", "a", "k", "zzz"} {
		d.fact("canPropertyBeDefined:"+p, func() any { return registry.CanPropertyBeDefined(t, p) })
	}
	d.fact("replacePartial", func() any {
		r := TemplateTypeReplacer_ForPartialReplacement(registry, d.partialMap)
		return taStr(r.VisitJSType(t)) + " " + strconv.FormatBool(r.HasMadeReplacement())
	})
	d.fact("replaceTotal", func() any {
		r := TemplateTypeReplacer_ForTotalReplacement(registry, d.totalMap)
		return taStr(r.VisitJSType(t)) + " " + strconv.FormatBool(r.HasMadeReplacement())
	})
	d.fact("replaceInfer", func() any {
		r := TemplateTypeReplacer_ForInference(registry, d.inferMap)
		return taStr(r.VisitJSType(t)) + " " + strconv.FormatBool(r.HasMadeReplacement())
	})
	d.fact("ttm.keys", func() any { return t.GetTemplateTypeMap().GetTemplateKeys() })
	d.fact("ttm.values", func() any { return t.GetTemplateTypeMap().GetTemplateValues() })
	d.fact("ttm.resolved", func() any {
		m := t.GetTemplateTypeMap()
		l := []JSType{}
		for _, k := range m.GetTemplateKeys() {
			l = append(l, m.GetResolvedTemplateType(k))
		}
		return l
	})
	d.fact("ttm.hasT", func() any {
		m := t.GetTemplateTypeMap()
		return strconv.FormatBool(m.HasTemplateKey(d.tKey)) + "/" + strconv.FormatBool(m.HasTemplateType(d.tKey)) +
			"/" + strconv.Itoa(m.Size()) + "/" + strconv.FormatBool(m.IsEmpty())
	})
	d.fact("ttm.lastByName(T)", func() any { return taPtr(t.GetTemplateTypeMap().GetLastTemplateTypeKeyByName("T")) })
	d.fact("readableTypeName", func() any {
		n := IR_Name("x")
		n.SetJSType(t)
		return registry.GetReadableTypeName(n) + "/" + registry.GetReadableTypeNameNoDeref(n)
	})
	if o := t.ToMaybeObjectType(); o != nil {
		d.fact("o.referenceName", func() any { return o.GetReferenceName() })
		d.fact("o.normalizedReferenceName", func() any {
			if !o.HasReferenceName() {
				return nil
			}
			return o.GetNormalizedReferenceName()
		})
		d.fact("o.hasReferenceName", func() any { return o.HasReferenceName() })
		d.fact("o.constructor", func() any { return o.GetConstructor() })
		d.fact("o.implicitPrototype", func() any { return o.GetImplicitPrototype() })
		d.fact("o.ownPropertyNames", func() any { return o.GetOwnPropertyNames() })
		d.fact("o.propertyNames", func() any { return o.GetPropertyNames() })
		d.fact("o.propertiesCount", func() any { return o.GetPropertiesCount() })
		d.fact("o.templateTypes", func() any { return taNullableList(o.GetTemplateTypes()) })
		d.fact("o.ownerFunction", func() any { return o.GetOwnerFunction() })
		d.fact("o.ctorImplementedInterfaces", func() any { return o.GetCtorImplementedInterfaces() })
		d.fact("o.ctorExtendedInterfaces", func() any { return o.GetCtorExtendedInterfaces() })
		d.fact("o.typeOfThis", func() any { return o.GetTypeOfThis() })
		d.fact("o.superClassConstructor", func() any { return o.GetSuperClassConstructor() })
		d.fact("o.rawType", func() any { return o.GetRawType() })
		d.fact("o.propertyTypeMap", func() any { return taPropertyTypeMap(o).String() })
		d.fact("o.ownPropertyKeys", func() any { return o.GetOwnPropertyKeys() })
		d.fact("o.hasCachedValues", func() any { return o.HasCachedValues() })
		d.fact("o.enumeratedTypeOfEnumObject", func() any { return o.GetEnumeratedTypeOfEnumObject() })
		d.fact("o.detectInheritanceCycle", func() any { return o.DetectInheritanceCycle() })
		d.fact("o.slot:bar", func() any { return o.GetSlot("bar") })
		d.fact("o.ownSlot:bar", func() any { return o.GetOwnSlot("bar") })
		d.fact("o.slot:k", func() any { return o.GetSlot("k") })
		for _, p := range taProps {
			d.fact("o.propertyType:"+p, func() any { return o.GetPropertyType(p) })
			d.fact("o.hasOwnProperty:"+p, func() any {
				return strconv.FormatBool(o.HasOwnProperty(p)) + "/" +
					strconv.FormatBool(o.IsPropertyTypeDeclared(p)) + "/" +
					strconv.FormatBool(o.IsPropertyTypeInferred(p))
			})
			d.fact("o.closestDefiningType:"+p, func() any { return o.GetClosestDefiningType(p) })
		}
	}
	if f := t.ToMaybeFunctionType(); f != nil {
		number := d.nat(JSTypeNative_NUMBER_TYPE)
		d.fact("f.kind", func() any { return f.GetKind() })
		d.fact("f.parameters", func() any { return f.GetParameters() })
		d.fact("f.minArity", func() any { return f.GetMinArity() })
		d.fact("f.maxArity", func() any { return f.GetMaxArity() })
		d.fact("f.returnType", func() any { return f.GetReturnType() })
		d.fact("f.isReturnTypeInferred", func() any { return f.IsReturnTypeInferred() })
		d.fact("f.typeOfThis", func() any { return f.GetTypeOfThis() })
		d.fact("f.hasInstanceType", func() any { return f.HasInstanceType() })
		d.fact("f.instanceType", func() any { return f.GetInstanceType() })
		d.fact("f.prototype", func() any { return f.GetPrototype() })
		d.fact("f.implementedInterfaces", func() any { return f.GetImplementedInterfaces() })
		d.fact("f.extendedInterfaces", func() any { return f.GetExtendedInterfaces() })
		d.fact("f.allImplementedInterfaces", func() any { return f.GetAllImplementedInterfaces() })
		d.fact("f.ancestorInterfaces", func() any { return f.GetAncestorInterfaces() })
		d.fact("f.isAbstract", func() any { return f.IsAbstract() })
		d.fact("f.makesStructs", func() any { return f.MakesStructs() })
		d.fact("f.makesDicts", func() any { return f.MakesDicts() })
		d.fact("f.hasProperties", func() any { return f.HasProperties() })
		d.fact("f.isAmbiguousConstructor", func() any { return f.IsAmbiguousConstructor() })
		d.fact("f.checkExtendsLoop", func() any { return taNullableList(f.CheckExtendsLoop()) })
		d.fact("f.constructorOnlyTemplateParameters", func() any { return f.GetConstructorOnlyTemplateParameters() })
		d.fact("f.bindReturnType(0)", func() any { return f.GetBindReturnType(0) })
		d.fact("f.bindReturnType(1)", func() any { return f.GetBindReturnType(1) })
		d.fact("f.acceptsArguments()", func() any { return f.AcceptsArguments([]JSType{}) })
		d.fact("f.acceptsArguments(number)", func() any { return f.AcceptsArguments([]JSType{number}) })
		d.fact(
			"f.acceptsArguments(number,string)",
			func() any { return f.AcceptsArguments([]JSType{number, d.nat(JSTypeNative_STRING_TYPE)}) })
		d.fact("f.propertyTypeMap", func() any { return taPropertyTypeMap(f).String() })
	}
	if u := t.ToMaybeUnionType(); u != nil {
		d.fact("u.alternates", func() any { return u.GetAlternates() })
		d.fact("u.restrictedUnion(null)", func() any { return u.GetRestrictedUnion(d.nat(JSTypeNative_NULL_TYPE)) })
		d.fact("u.restrictedUnion(Object)", func() any { return u.GetRestrictedUnion(d.nat(JSTypeNative_OBJECT_TYPE)) })
		d.fact("u.contains(number)", func() any { return u.Contains(d.nat(JSTypeNative_NUMBER_TYPE)) })
	}
	if e := t.ToMaybeEnumType(); e != nil {
		d.fact("e.elements", func() any { return e.GetElements() })
		d.fact("e.elementsType", func() any { return e.GetElementsType() })
	}
	if tt := t.ToMaybeTemplatizedType(); tt != nil {
		d.fact("tt.referencedType", func() any { return tt.GetReferencedType() })
		d.fact("tt.templateTypes", func() any { return taNullableList(tt.GetTemplateTypes()) })
	}
	if tk := t.ToMaybeTemplateType(); tk != nil {
		d.fact("tk.bound", func() any { return tk.GetBound() })
		d.fact("tk.containsCycle", func() any { return tk.ContainsCycle() })
	}
	if nt := t.ToMaybeNamedType(); nt != nil {
		d.fact("nt.referencedType", func() any { return nt.GetReferencedType() })
	}
}

func taIdent(r, a, b JSType) string {
	if r == a && r == b {
		return "=AB"
	} else if r == a {
		return "=A"
	} else if r == b {
		return "=B"
	}
	return taStr(r)
}

func (d *typeAlgebraDump) pair(i, j int) {
	a := d.types[i]
	b := d.types[j]
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(i + 1))
	sb.WriteByte(' ')
	sb.WriteString(strconv.Itoa(j + 1))
	sb.WriteByte('\t')
	sb.WriteString(taPf(func() any {
		return taBits(
			a.IsSubtypeOf(b),
			a.IsSubtypeOfWithMode(b, JSType_SubtypingMode_IGNORE_NULL_UNDEFINED),
			a.IsSubtypeWithoutStructuralTyping(b),
			a.Equals(b),
			a.DiffersFrom(b),
			a.CanTestForEqualityWith(b),
			a.CanTestForShallowEqualityWith(b),
			a.CanCastTo(b),
			JSType_AreSimilar(a, b))
	}))
	sb.WriteByte('\t')
	sb.WriteString(taPf(func() any { return a.TestForEquality(b) }))
	sb.WriteByte('\t')
	sb.WriteString(taPf(func() any { return taIdent(a.GetLeastSupertype(b), a, b) }))
	sb.WriteByte('\t')
	sb.WriteString(taPf(func() any { return taIdent(a.GetGreatestSubtype(b), a, b) }))
	sb.WriteByte('\t')
	sb.WriteString(taPf(func() any { return a.GetTypesUnderEquality(b) }))
	sb.WriteByte('\t')
	sb.WriteString(taPf(func() any { return a.GetTypesUnderInequality(b) }))
	sb.WriteByte('\t')
	sb.WriteString(taPf(func() any { return a.GetTypesUnderShallowInequality(b) }))
	sb.WriteByte('\t')
	sb.WriteString(taPf(func() any { return a.GetTypesUnderShallowEquality(b) }))
	sb.WriteByte('\t')
	sb.WriteString(taPf(func() any {
		fa, aIsFn := a.(FunctionType)
		fb, bIsFn := b.(FunctionType)
		if aIsFn && bIsFn {
			return "fn:" + strconv.FormatBool(fa.HasEqualCallType(fb))
		}
		oa, aIsObj := a.(ObjectType)
		ob, bIsObj := b.(ObjectType)
		if aIsObj && bIsObj {
			return "obj:" + taStr(d.registry.FindCommonSuperObject(oa, ob))
		}
		return "-"
	}))
	d.line(taEsc(sb.String()))
}

func (d *typeAlgebraDump) registryFacts() {
	d.line("=== registry")
	for _, p := range []string{"bar", "a", "k", "zzz"} {
		d.fact("eachReferenceTypeWithProperty:"+p, func() any { return d.registry.GetEachReferenceTypeWithProperty(p) })
	}
	for _, n := range []string{"Foo", "Box", "Late", "Missing", "E", "Pair", "Forwarded", "List", "ns"} {
		d.fact("getType:"+n, func() any { return d.registry.GetType(d.scope, n) })
		d.fact("getGlobalType:"+n, func() any { return d.registry.GetGlobalType(n) })
		d.fact("isNonNullableName:"+n, func() any { return d.registry.IsNonNullableName(d.scope, n) })
	}
}

func (d *typeAlgebraDump) run(pairsArg string) {
	d.build()
	d.registryFacts()
	for i := range d.types {
		d.unaryFacts(i)
	}
	d.line("=== pairs")
	n := len(d.types)
	if pairsArg == "all" {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				d.pair(i, j)
			}
		}
	} else {
		count, err := strconv.Atoi(pairsArg)
		if err != nil {
			panic(err)
		}
		for k := 0; k < count; k++ {
			d.pair(d.nextInt(n), d.nextInt(n))
		}
	}
	d.line("=== end")
}

// ---------------------------------------------------------------------
// The test

func readTypeAlgebraExpected(t *testing.T) []string {
	path := os.Getenv("JSTYPE_ORACLE_FILE")
	if path == "" {
		path = "testdata/jstype_type_algebra.txt.gz"
	}
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var r io.Reader = f
	if strings.HasSuffix(path, ".gz") {
		zr, err := gzip.NewReader(f)
		if err != nil {
			t.Fatal(err)
		}
		defer zr.Close()
		r = zr
	}
	var lines []string
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 1<<20), 1<<24)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	return lines
}

func TestJSType_TypeAlgebraOracle(t *testing.T) {
	expected := readTypeAlgebraExpected(t)
	asts := map[string]string{}
	pairsArg := "all"
	pairLines := 0
	inPairs := false
	for _, l := range expected {
		if rest, ok := strings.CutPrefix(l, "ast\t"); ok {
			parts := strings.SplitN(rest, "\t", 2)
			asts[taUnesc(parts[0])] = taUnesc(parts[1])
		}
		if l == "=== pairs" {
			inPairs = true
		} else if l == "=== end" {
			inPairs = false
		} else if inPairs {
			pairLines++
		}
	}
	d := newTypeAlgebraDump(asts)
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panic while building the dump: %v\noutput so far ends with:\n%s",
					r, lastLines(d.out.String(), 20))
				panic(r)
			}
		}()
		// The pair section is either all pairs or a sample; recover which.
		nTypes := 0
		for _, l := range expected {
			if strings.HasPrefix(l, "add\t") {
				nTypes++
			}
		}
		if pairLines != nTypes*nTypes {
			pairsArg = strconv.Itoa(pairLines)
		}
		d.run(pairsArg)
	}()
	actual := strings.Split(strings.TrimSuffix(d.out.String(), "\n"), "\n")

	if dir := os.Getenv("JSTYPE_ORACLE_OUT"); dir != "" {
		_ = os.WriteFile(dir, []byte(d.out.String()), 0o644)
	}

	diffs := 0
	n := max(len(expected), len(actual))
	for i := 0; i < n; i++ {
		var e, a string
		if i < len(expected) {
			e = expected[i]
		}
		if i < len(actual) {
			a = actual[i]
		}
		if e != a {
			diffs++
			if diffs <= 40 {
				t.Errorf("line %d (%s):\n  java: %s\n  go:   %s", i+1, typeAlgebraContext(expected, i), e, a)
			}
		}
	}
	if diffs > 0 {
		t.Errorf("%d of %d lines differ", diffs, len(expected))
	} else {
		t.Logf("%d types, %d pair lines, %d lines identical", len(d.types), pairLines, len(expected))
	}
}

// typeAlgebraContext names the section a line of the dump belongs to.
func typeAlgebraContext(lines []string, i int) string {
	for j := min(i, len(lines)-1); j >= 0; j-- {
		if strings.HasPrefix(lines[j], "=== ") {
			return lines[j]
		}
	}
	return "build"
}

func lastLines(s string, n int) string {
	ls := strings.Split(s, "\n")
	if len(ls) > n {
		ls = ls[len(ls)-n:]
	}
	return strings.Join(ls, "\n")
}
