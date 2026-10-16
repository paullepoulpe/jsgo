package jscomp

// API-sweep oracle for TypeValidator: mirrors
// tools/oracle/java/TypeValidatorDump.java step by step (same compiler
// configuration, same type universe - built by the raiDump of
// type_reverse_interpreter_oracle_test.go plus the same extra classes -,
// same calls in the same order, same seeded choices) and must reproduce its
// output byte for byte.
//
//   - TestTypeValidatorOracleGolden compares with the checked-in
//     testdata/type_validator.txt.gz (3000 seeded pairs):
//       java -cp .cache/closure-compiler.jar:.cache/oracle-build \
//         com.google.javascript.jscomp.TypeValidatorDump \
//         internal/jscomp/testdata/reverse_interpreter_types.txt \
//         internal/jscomp/testdata/type_validator.txt.gz 3000
//   - TestTypeValidatorOracle compares with any other dump:
//     TV_ORACLE=<dump.gz> TV_PAIRS=<the pairs it was made with>.
// TV_ORACLE_OUT=<path> also writes the Go output there.

import (
	"compress/gzip"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type tvDump struct {
	*raiDump
	compiler      AbstractCompiler
	reported      []string
	validator     *TypeValidator
	mismatchCount int
}

// tvErrorHandler is the ErrorHandler lambda of TypeValidatorDump.
type tvErrorHandler struct{ d *tvDump }

func (h *tvErrorHandler) Report(level *CheckLevel, error *JSError) {
	h.d.reported = append(h.d.reported,
		"  D "+jutil.ToJavaString(level)+
			" "+error.GetType().key+
			" "+raiEsc(error.GetDescription())+
			" @"+tvNodeLabel(error.GetNode()))
}

func newTvDump() *tvDump {
	d := &tvDump{}
	options := NewCompilerOptions()
	options.SetErrorHandler(&tvErrorHandler{d: d})
	d.compiler = newTestCompiler(options)
	d.raiDump = &raiDump{
		out:      &strings.Builder{},
		seed:     0x5DEECE66D,
		scope:    &raiTestScope{root: rhino.NewNode(rhino.Token_ROOT)},
		registry: d.compiler.GetTypeRegistry(),
	}
	return d
}

func tvNodeLabel(n *rhino.Node) string {
	if n == nil {
		return "null"
	}
	if n.IsName() || n.IsStringLit() {
		return n.GetToken().String() + "(" + raiEsc(n.GetString()) + ")"
	}
	return n.GetToken().String()
}

// build mirrors TypeValidatorDump.build.
func (d *tvDump) build(typeStrings []string) {
	d.raiDump.build(typeStrings)
	closer := d.registry.GetResolver().OpenForDefinition()
	defer closer.Close()
	d.buildExtraTypes()
}

// buildExtraTypes mirrors TypeValidatorDump.buildExtraTypes.
func (d *tvDump) buildExtraTypes() {
	r := d.registry
	number := d.nat(rhino.JSTypeNative_NUMBER_TYPE)
	str := d.nat(rhino.JSTypeNative_STRING_TYPE)
	numberOrUndefined := r.CreateUnionType(number, d.nat(rhino.JSTypeNative_VOID_TYPE))

	i2Ctor := rhino.FunctionType_Builder_(r).ForInterface().WithName(jutil.StringPtr("I2")).WithSourceNode(raiFn("I2")).Build()
	i2Ctor.GetPrototype().DefineDeclaredProperty("m", number, nil)
	i2Ctor.GetPrototype().DefineDeclaredProperty("opt", numberOrUndefined, nil)
	i2 := i2Ctor.GetInstanceType()
	i2.DefineDeclaredProperty("own", str, nil)

	badImplCtor := r.CreateConstructorType(
		jutil.StringPtr("BadImpl"), raiFn("BadImpl"), r.CreateParameters(), nil, nil, false)
	badImplCtor.SetImplementedInterfaces([]rhino.ObjectType{i2})

	wrongImplCtor := r.CreateConstructorType(
		jutil.StringPtr("WrongImpl"), raiFn("WrongImpl"), r.CreateParameters(), nil, nil, false)
	wrongImplCtor.SetImplementedInterfaces([]rhino.ObjectType{i2})
	wrongImplCtor.GetPrototype().DefineDeclaredProperty("m", str, nil)
	wrongImplCtor.GetInstanceType().DefineDeclaredProperty("own", number, nil)

	goodImplCtor := r.CreateConstructorType(
		jutil.StringPtr("GoodImpl"), raiFn("GoodImpl"), r.CreateParameters(), nil, nil, false)
	goodImplCtor.SetImplementedInterfaces([]rhino.ObjectType{i2})
	goodImplCtor.GetPrototype().DefineDeclaredProperty("m", number, nil)
	goodImplCtor.GetInstanceType().DefineDeclaredProperty("own", str, nil)

	i3Ctor := rhino.FunctionType_Builder_(r).ForInterface().WithName(jutil.StringPtr("I3")).WithSourceNode(raiFn("I3")).Build()
	i3Ctor.SetExtendedInterfaces([]rhino.ObjectType{i2})
	i3Ctor.GetPrototype().DefineDeclaredProperty("n", str, nil)

	abstractMethod := rhino.FunctionType_Builder_(r).WithReturnType(number).WithIsAbstract(true).Build()
	concreteMethod := rhino.FunctionType_Builder_(r).WithReturnType(number).Build()
	absCtor := rhino.FunctionType_Builder_(r).
		ForConstructor().
		WithName(jutil.StringPtr("Abs")).
		WithSourceNode(raiFn("Abs")).
		WithIsAbstract(true).
		Build()
	absCtor.GetPrototype().DefineDeclaredProperty("am", abstractMethod, nil)
	absCtor.GetPrototype().DefineDeclaredProperty("cm", concreteMethod, nil)
	abs2Ctor := rhino.FunctionType_Builder_(r).
		ForConstructor().
		WithName(jutil.StringPtr("Abs2")).
		WithSourceNode(raiFn("Abs2")).
		WithIsAbstract(true).
		Build()
	abs2Ctor.SetPrototypeBasedOn(absCtor.GetInstanceType())
	abs2Ctor.GetPrototype().DefineDeclaredProperty("am2", abstractMethod, nil)
	concreteCtor := r.CreateConstructorType(
		jutil.StringPtr("Concrete"), raiFn("Concrete"), r.CreateParameters(), nil, nil, false)
	concreteCtor.SetPrototypeBasedOn(abs2Ctor.GetInstanceType())
	concreteCtor.GetPrototype().DefineDeclaredProperty("am2", concreteMethod, nil)
	concrete2Ctor := r.CreateConstructorType(
		jutil.StringPtr("Concrete2"), raiFn("Concrete2"), r.CreateParameters(), nil, nil, false)
	concrete2Ctor.SetPrototypeBasedOn(absCtor.GetInstanceType())
	concrete2Ctor.GetPrototype().DefineDeclaredProperty("am", concreteMethod, nil)

	structCtor := r.CreateConstructorType(
		jutil.StringPtr("Struct"), raiFn("Struct"), r.CreateParameters(), nil, nil, false)
	structCtor.SetStruct()
	structCtor.GetPrototype().DefineDeclaredProperty("s", number, nil)

	d.add("I2", i2)
	d.add("ctor:I2", i2Ctor)
	d.add("ctor:BadImpl", badImplCtor)
	d.add("BadImpl", badImplCtor.GetInstanceType())
	d.add("ctor:WrongImpl", wrongImplCtor)
	d.add("WrongImpl", wrongImplCtor.GetInstanceType())
	d.add("ctor:GoodImpl", goodImplCtor)
	d.add("GoodImpl", goodImplCtor.GetInstanceType())
	d.add("ctor:I3", i3Ctor)
	d.add("I3", i3Ctor.GetInstanceType())
	d.add("abstractMethod", abstractMethod)
	d.add("ctor:Abs", absCtor)
	d.add("Abs", absCtor.GetInstanceType())
	d.add("ctor:Abs2", abs2Ctor)
	d.add("ctor:Concrete", concreteCtor)
	d.add("Concrete", concreteCtor.GetInstanceType())
	d.add("ctor:Concrete2", concrete2Ctor)
	d.add("ctor:Struct", structCtor)
	d.add("Struct", structCtor.GetInstanceType())
}

// tvExceptionResult is the harness's catch clauses: a Java exception is
// "X <class>: <message>" and a NullPointerException (a nil dereference in Go)
// is "X java.lang.NullPointerException".
func tvExceptionResult(r any) string {
	if re, ok := r.(runtime.Error); ok && strings.Contains(re.Error(), "nil pointer dereference") {
		return "X java.lang.NullPointerException"
	}
	t := jutil.AsThrowable(r)
	if t == nil {
		panic(r)
	}
	if _, ok := t.(*jutil.NullPointerException); ok {
		return "X java.lang.NullPointerException"
	}
	return "X " + cfaJavaClass(t) + ": " + raiEsc(cfaJavaMessage(t))
}

// call mirrors TypeValidatorDump.call.
func (d *tvDump) call(desc string, c func() string) {
	var result string
	func() {
		defer func() {
			if r := recover(); r != nil {
				result = tvExceptionResult(r)
			}
		}()
		result = c()
	}()
	d.line(desc + " -> " + result)
	for _, r := range d.reported {
		d.line(r)
	}
	d.reported = d.reported[:0]
	all := d.validator.GetMismatches()
	for i := d.mismatchCount; i < len(all); i++ {
		m := all[i]
		d.line("  M " + raiStr(m.Found()) + " | " + raiStr(m.Required()) + " @" + tvNodeLabel(m.Location()))
	}
	d.mismatchCount = len(all)
}

func tvBool(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

func tvTyped(n *rhino.Node, t rhino.JSType) *rhino.Node {
	if t != nil {
		n.SetJSType(t)
	}
	return n
}

var tvProps = []string{"p", "bar", "m", "own", "length", "value"}

// unary mirrors TypeValidatorDump.unary.
func (d *tvDump) unary(i int) {
	t := d.types[i]
	v := d.validator
	idx := " " + strconv.Itoa(i+1)
	name := func() *rhino.Node { return rhino.IR_Name("n") }
	void := func(f func()) func() string { return func() string { f(); return "" } }
	d.call("U expectObject"+idx, func() string { return tvBool(v.ExpectObject(name(), t, "msg")) })
	d.call("U expectActualObject"+idx, void(func() { v.ExpectActualObject(name(), t, "msg") }))
	d.call("U expectAnyObject"+idx, void(func() { v.ExpectAnyObject(name(), t, "msg") }))
	d.call("U expectAutoboxesToIterable"+idx, func() string { return tvBool(v.ExpectAutoboxesToIterable(name(), t, "msg")) })
	d.call("U expectAutoboxesToIterableOrAsyncIterable"+idx, func() string {
		r := v.ExpectAutoboxesToIterableOrAsyncIterable(name(), t, "msg")
		if r != nil {
			return raiStr(r)
		}
		return "empty"
	})
	d.call("U expectGeneratorSupertype"+idx, void(func() { v.ExpectGeneratorSupertype(name(), t, "msg") }))
	d.call("U expectAsyncGeneratorSupertype"+idx, void(func() { v.ExpectAsyncGeneratorSupertype(name(), t, "msg") }))
	d.call("U expectValidAsyncReturnType"+idx, void(func() { v.ExpectValidAsyncReturnType(name(), t) }))
	d.call("U expectITemplateArraySupertype"+idx, void(func() { v.ExpectITemplateArraySupertype(name(), t, "msg") }))
	d.call("U expectString"+idx, void(func() { v.ExpectString(name(), t, "msg") }))
	d.call("U expectNumber"+idx, void(func() { v.ExpectNumber(name(), t, "msg") }))
	d.call("U expectNumberStrict"+idx, void(func() { v.ExpectNumberStrict(name(), t, "msg") }))
	d.call("U expectBitwiseable"+idx, void(func() { v.ExpectBitwiseable(name(), t, "msg") }))
	d.call("U expectNumberOrSymbol"+idx, void(func() { v.ExpectNumberOrSymbol(name(), t, "msg") }))
	d.call("U expectStringOrSymbol"+idx, void(func() { v.ExpectStringOrSymbol(name(), t, "msg") }))
	d.call("U expectUnknownOrComparable"+idx, void(func() { v.ExpectUnknownOrComparable(name(), t, "msg") }))
	d.call("U expectStringOrNumberOrSymbol"+idx, void(func() { v.ExpectStringOrNumberOrSymbol(name(), t, "msg") }))
	d.call("U expectStringOrNumberOrSymbolStrict"+idx, void(func() { v.ExpectStringOrNumberOrSymbolStrict(name(), t, "msg") }))
	d.call("U expectBigIntOrNumber"+idx, void(func() { v.ExpectBigIntOrNumber(name(), t, "msg") }))
	d.call("U expectNotNullOrUndefined"+idx, func() string {
		return tvBool(v.ExpectNotNullOrUndefined(nil, name(), t, "msg", d.nat(rhino.JSTypeNative_OBJECT_TYPE)))
	})
	d.call("U expectWellFormedTemplatizedType"+idx, void(func() { v.ExpectWellFormedTemplatizedType(tvTyped(name(), t)) }))
	d.call("U getInstanceOfCtor"+idx, func() string {
		r := TypeValidator_GetInstanceOfCtor(t)
		if r == nil {
			return "null"
		}
		return raiStr(r)
	})
	fnType := t.ToMaybeFunctionType()
	if fnType != nil && (fnType.IsConstructor() || fnType.IsInterface()) {
		d.call("U expectAllInterfaceProperties"+idx, void(func() { v.ExpectAllInterfaceProperties(name(), fnType) }))
	}
	if fnType != nil && fnType.IsConstructor() {
		d.call("U expectAbstractMethodsImplemented"+idx, void(func() { v.ExpectAbstractMethodsImplemented(name(), fnType) }))
	}
}

// binary mirrors TypeValidatorDump.binary.
func (d *tvDump) binary() {
	n := len(d.types)
	i := d.nextInt(n)
	j := d.nextInt(n)
	k := d.nextInt(n)
	a := d.types[i]
	bt := d.types[j]
	c := d.types[k]
	v := d.validator
	ij := " " + strconv.Itoa(i+1) + " " + strconv.Itoa(j+1)
	ijk := ij + " " + strconv.Itoa(k+1)
	name := func(s string) *rhino.Node { return rhino.IR_Name(s) }
	void := func(f func()) func() string { return func() string { f(); return "" } }
	switch d.nextInt(11) {
	case 0:
		d.call("P expectMatchingTypesStrict"+ij, void(func() { v.ExpectMatchingTypesStrict(name("n"), a, bt, "msg") }))
	case 1:
		d.call("P expectSwitchMatchesCase"+ij, void(func() {
			v.ExpectSwitchMatchesCase(rhino.NewNodeWith2Children(rhino.Token_CASE, name("c"), rhino.IR_Block()), a, bt)
		}))
	case 2:
		d.call("P expectIndexMatch"+ij, void(func() {
			v.ExpectIndexMatch(rhino.IR_Getelem(name("o"), name("i")), a, bt)
		}))
	case 3:
		d.call("P expectCanAssignTo"+ij, func() string { return tvBool(v.ExpectCanAssignTo(name("n"), a, bt, "msg")) })
	case 4:
		prop := tvProps[d.nextInt(len(tvProps))]
		d.call("P expectCanAssignToPropertyOf"+ijk+" "+prop, func() string {
			return tvBool(v.ExpectCanAssignToPropertyOf(name("n"), a, bt, tvTyped(name("owner"), c), prop))
		})
	case 5:
		ordinal := 1 + d.nextInt(3)
		d.call("P expectArgumentMatchesParameter"+ijk+" "+strconv.Itoa(ordinal), void(func() {
			v.ExpectArgumentMatchesParameter(name("n"), a, bt, rhino.IR_Call(tvTyped(name("f"), c)), ordinal)
		}))
	case 6:
		d.call("P expectCanCast"+ij, void(func() { v.ExpectCanCast(name("n"), a, bt) }))
	case 7:
		d.call("P expectCanAssignToPrototype"+ij, void(func() { v.ExpectCanAssignToPrototype(a, name("n"), bt) }))
	case 8:
		sub := a.ToMaybeFunctionType()
		sup := bt.ToMaybeFunctionType()
		if sub != nil {
			d.call("P expectExtends"+ij, void(func() { v.ExpectExtends(name("n"), sub, sup) }))
		}
	case 9:
		parent := bt.ToMaybeObjectType()
		prop := tvProps[d.nextInt(len(tvProps))]
		if parent != nil {
			d.call("P checkPropertyType"+ijk+" "+prop, void(func() {
				v.CheckPropertyType(name("n"), a, parent, rhino.NewProperty_StringKey(prop), c)
			}))
		}
	default:
		d.call("P expectNotNullOrUndefined"+ij, func() string {
			return tvBool(v.ExpectNotNullOrUndefined(nil, name("n"), a, "msg", bt))
		})
	}
}

// sweep mirrors TypeValidatorDump.sweep.
func (d *tvDump) sweep(mode rhino.JSType_SubtypingMode, pairs int) {
	d.line("mode " + mode.String())
	d.validator = NewTypeValidator(d.compiler)
	d.validator.SetSubtypingMode(mode)
	d.mismatchCount = 0
	d.call("U expectValidTypeofName", func() string {
		d.validator.ExpectValidTypeofName(rhino.IR_Name("n"), "foo")
		return ""
	})
	for i := range d.types {
		d.unary(i)
	}
	for p := 0; p < pairs; p++ {
		if p%500 == 0 {
			// A fresh validator every 500 pairs keeps getMismatches() (a copy) short.
			d.line("reset")
			d.validator = NewTypeValidator(d.compiler)
			d.validator.SetSubtypingMode(mode)
			d.mismatchCount = 0
		}
		d.binary()
	}
}

func (d *tvDump) labelIndex(label string) int {
	for i, l := range d.labels {
		if l == label {
			return i
		}
	}
	return -1
}

// scenarios mirrors TypeValidatorDump.scenarios.
func (d *tvDump) scenarios() {
	d.line("scenarios")
	d.validator = NewTypeValidator(d.compiler)
	d.mismatchCount = 0
	v := d.validator
	r := d.registry
	number := d.nat(rhino.JSTypeNative_NUMBER_TYPE)
	str := d.nat(rhino.JSTypeNative_STRING_TYPE)
	func() {
		closer := r.GetResolver().OpenForDefinition()
		defer closer.Close()
		ctor := func(name string) rhino.FunctionType {
			return r.CreateConstructorType(jutil.StringPtr(name), raiFn(name), r.CreateParameters(), nil, nil, false)
		}
		aCtor := ctor("SA")
		bCtor := ctor("SB")
		cCtor := ctor("SC")
		cCtor.SetPrototypeBasedOn(aCtor.GetInstanceType())
		dCtor := ctor("SD")
		boxCtor := d.types[d.labelIndex("ctor:Box")].ToMaybeFunctionType()
		eCtor := ctor("SE")
		eCtor.SetPrototypeBasedOn(r.CreateTemplatizedType(boxCtor.GetInstanceType(), number))
		fCtor := ctor("SF")
		fCtor.SetPrototypeBasedOn(aCtor.GetInstanceType())
		superOf := func(c rhino.FunctionType) string {
			return raiStr(c.GetInstanceType().GetImplicitPrototype().GetImplicitPrototype())
		}
		d.call("S expectSuperType SB SC", func() string {
			v.ExpectSuperType(rhino.IR_Name("n"), bCtor.GetInstanceType(), cCtor.GetInstanceType())
			return superOf(cCtor)
		})
		d.call("S expectSuperType SA SC", func() string {
			v.ExpectSuperType(rhino.IR_Name("n"), aCtor.GetInstanceType(), cCtor.GetInstanceType())
			return superOf(cCtor)
		})
		d.call("S expectSuperType SA SD", func() string {
			v.ExpectSuperType(rhino.IR_Name("n"), aCtor.GetInstanceType(), dCtor.GetInstanceType())
			return superOf(dCtor)
		})
		d.call("S expectSuperType Box SE", func() string {
			v.ExpectSuperType(rhino.IR_Name("n"), boxCtor.GetInstanceType(), eCtor.GetInstanceType())
			return superOf(eCtor)
		})
		d.call("S expectSuperType ? SF", func() string {
			v.ExpectSuperType(rhino.IR_Name("n"), d.nat(rhino.JSTypeNative_UNKNOWN_TYPE).(rhino.ObjectType), fCtor.GetInstanceType())
			return superOf(fCtor)
		})
		fCtor.GetInstanceType().GetPropertyType("toString") // caches values
		d.call("S expectSuperType SB SF (cached)", func() string {
			v.ExpectSuperType(rhino.IR_Name("n"), bCtor.GetInstanceType(), fCtor.GetInstanceType())
			return superOf(fCtor)
		})
	}()

	// expectUndeclaredVariable
	script := rhino.IR_Script()
	root := rhino.IR_Root(script)
	global := TypedScope_CreateGlobalScope(root)
	input1 := NewCompilerInput(SourceFile_FromCode("first.js", ""))
	input2 := NewCompilerInput(SourceFile_FromCode("second.js", ""))
	names := []string{"a", "b", "c", "d", "e", "f", "g"}
	// The @suppress {duplicate} case (5) runs last: its outcome depends on SuppressDocWarningsGuard.
	for _, s := range []int{0, 1, 2, 3, 4, 6, 5} {
		if s == 5 {
			d.line(tvSuppressionMarker)
		}
		name := names[s]
		oldName := rhino.IR_Name(name)
		var oldDecl *rhino.Node
		if s == 3 {
			oldDecl = rhino.IR_ExprResult(oldName)
		} else {
			oldDecl = rhino.IR_Var(oldName)
		}
		if s == 5 {
			info := rhino.JSDocInfo_Builder_()
			info.RecordSuppressions(jutil.NewLinkedHashSet("duplicate"))
			oldDecl.SetJSDocInfo(info.Build())
		}
		script.AddChildToBack(oldDecl)
		oldName.SetLinenoCharno(10+s, 4)
		oldType := number
		if s == 4 {
			oldType = d.nat(rhino.JSTypeNative_UNKNOWN_TYPE)
		}
		var oldInput *CompilerInput
		if s != 0 {
			oldInput = input1
		}
		var_ := global.Declare(name, oldName, oldType, oldInput, false)
		newName := rhino.IR_Name(name)
		var newDecl *rhino.Node
		if s == 6 {
			newName = rhino.IR_Function(rhino.IR_Name(name), rhino.IR_ParamList(), rhino.IR_Block())
			newDecl = rhino.IR_ExprResult(newName)
		} else {
			newDecl = rhino.IR_Var(newName)
			newName.AddChildToBack(rhino.IR_Number(1))
		}
		script.AddChildToBack(newDecl)
		nn := newName
		newType := number
		if s == 2 || s == 5 {
			newType = str
		}
		d.call("S expectUndeclaredVariable "+name, func() string {
			parent := nn.GetParent()
			target := nn
			if s == 6 {
				parent = nn
				target = nn.GetFirstChild()
			}
			result := v.ExpectUndeclaredVariable("second.js", input2, target, parent, var_, name, newType)
			same := "new"
			if result == var_ {
				same = "same"
			}
			firstChild := "-"
			if target.HasChildren() {
				firstChild = raiStr(target.GetFirstChild().GetJSType())
			}
			return same + " " + raiEsc(result.GetInputName()) + " " + raiStr(result.GetType()) + " " +
				raiStr(target.GetJSType()) + " " + firstChild + " " + raiStr(parent.GetJSType())
		})
	}
}

// run mirrors TypeValidatorDump.main.
func (d *tvDump) run(typeStrings []string, pairs int) string {
	d.build(typeStrings)
	d.line("built")
	for _, r := range d.reported {
		d.line(r)
	}
	d.reported = d.reported[:0]
	d.sweep(rhino.JSType_SubtypingMode_NORMAL, pairs)
	d.sweep(rhino.JSType_SubtypingMode_IGNORE_NULL_UNDEFINED, pairs)
	d.scenarios()
	return d.out.String()
}

func tvReadGz(t *testing.T, path string) string {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	zr, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	b, err := io.ReadAll(zr)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func tvCompare(t *testing.T, want string, pairs int) {
	got := newTvDump().run(raiReadTypeStrings(t), pairs)
	if p := os.Getenv("TV_ORACLE_OUT"); p != "" {
		if err := os.WriteFile(p, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if got == want {
		return
	}
	// Everything before the suppression scenario must match; that scenario
	// needs SuppressDocWarningsGuard, which the test compiler does not have
	// yet (see test_compiler_test.go).
	wantHead, wantTail, _ := strings.Cut(want, tvSuppressionMarker+"\n")
	gotHead, gotTail, _ := strings.Cut(got, tvSuppressionMarker+"\n")
	if gotHead != wantHead {
		t.Fatalf("differs from Java at %s", cfaFirstDiff(want, got))
	}
	t.Skipf("only the suppression scenario differs (needs SuppressDocWarningsGuard): %s",
		cfaFirstDiff(wantTail, gotTail))
}

const tvSuppressionMarker = "suppression scenario (SuppressDocWarningsGuard)"

func TestTypeValidatorOracleGolden(t *testing.T) {
	tvCompare(t, tvReadGz(t, "testdata/type_validator.txt.gz"), 3000)
}

func TestTypeValidatorOracle(t *testing.T) {
	dump := os.Getenv("TV_ORACLE")
	if dump == "" {
		t.Skip("set TV_ORACLE=<TypeValidatorDump output .gz> and TV_PAIRS")
	}
	pairs, err := strconv.Atoi(os.Getenv("TV_PAIRS"))
	if err != nil {
		t.Fatal("TV_PAIRS: ", err)
	}
	tvCompare(t, tvReadGz(t, dump), pairs)
}
