package jscomp

// TypedScope and TypedVar have no Java unit tests; this replays the scenario
// of tools/oracle/java/TypedScopeVectors.java (TypedScopes built by hand over
// a parsed file, every query printed) and compares with its output,
// testdata/typed_scope.txt. Regenerate with
//
//	java -cp .cache/$CLOSURE_JAR:<build> com.google.javascript.jscomp.TypedScopeVectors \
//	  internal/jscomp/testdata/typed_scope.txt

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const typedScopeVectors_SRC = "var a = 1;\n" +
	"a.b = 2;\n" +
	"function f(x) { var a; var u; let w; this.q; }\n" +
	"class C { m() {} static { } fld = 1; static sfld = 2; }\n" +
	"var u;\n" +
	"var arguments;\n"

var typedScopeVectors_NAMES = []string{
	"a", "a.b", "a.c", "b", "x", "u", "z", "w", "this", "super", "arguments", "exports", "f", "C",
}

// typedScopeVectors holds TypedScopeVectors' labels and output.
type typedScopeVectors struct {
	labels map[any]string
	out    strings.Builder
}

func (v *typedScopeVectors) println(s string) {
	v.out.WriteString(s)
	v.out.WriteByte('\n')
}

// label is TypedScopeVectors.label.
func (v *typedScopeVectors) label(scope any) string {
	if jutil.IsNil(scope) {
		return "null"
	}
	if l, ok := v.labels[scope]; ok {
		return l
	}
	return "?"
}

// typedScopeVectorsStr is TypedScopeVectors.str.
func typedScopeVectorsStr(s func() string) (result string) {
	jutil.Try(func() { result = s() }, func(e jutil.Throwable) { result = "!" + scopeDumpJavaClass(e) })
	return result
}

// describe is TypedScopeVectors.describe.
func (v *typedScopeVectors) describe(tv *TypedVar) string {
	if tv == nil {
		return "null"
	}
	node := "null"
	if !tv.GetNode().IsNil() {
		node = tv.GetNode().GetToken().String()
	}
	scope, _ := tv.GetScope().(*TypedScope)
	return tv.String() +
		"|scope=" + v.label(scope) +
		"|type=" + jutil.ToJavaString(tv.GetType()) +
		"|inferred=" + scopeDumpBool(tv.IsTypeInferred()) +
		"|index=" + strconv.Itoa(tv.GetIndex()) +
		"|node=" + node +
		"|input=" + tv.GetInputName() +
		"|extern=" + scopeDumpBool(tv.IsExtern())
}

func (v *typedScopeVectors) run(t *testing.T) string {
	options := NewCompilerOptions()
	options.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	compiler := newTestCompiler(options)
	config := compiler.GetParserConfig(AbstractCompiler_ConfigContext_DEFAULT)
	file := SourceFile_FromCode("typed.js", typedScopeVectors_SRC)
	script := parsing.ParserRunner_Parse(file, typedScopeVectors_SRC, config, scopeDumpReporter{}).Ast
	script.SetInputId(rhino.NewInputId("typed.js"))
	root := rhino.IR_Root(script)
	input := NewCompilerInput(file)
	registry := compiler.GetTypeRegistry()

	aName := script.GetFirstChild().GetFirstChild()
	abGetprop := script.GetSecondChild().GetFirstChild().GetFirstChild()
	fn := script.GetChildAtIndex(2)
	xName := fn.GetSecondChild().GetFirstChild()
	fBody := fn.GetLastChild()
	innerA := fBody.GetFirstChild().GetFirstChild()
	innerU := fBody.GetSecondChild().GetFirstChild()
	cls := script.GetChildAtIndex(3)
	members := cls.GetLastChild()
	mFn := members.GetFirstChild().GetFirstChild()
	staticBlock := members.GetSecondChild()
	fld := members.GetChildAtIndex(2)
	sfld := members.GetChildAtIndex(3)
	uName := script.GetChildAtIndex(4).GetFirstChild()
	argsName := script.GetChildAtIndex(5).GetFirstChild()

	num := registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)
	string_ := registry.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
	bool_ := registry.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE)
	obj := registry.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE)
	ctor := registry.GetNativeType(rhino.JSTypeNative_DATE_FUNCTION_TYPE)
	root.SetJSType(obj)
	fn.SetJSType(ctor)
	cls.SetJSType(ctor)

	global := TypedScope_CreateGlobalScope(root)
	v.labels[global] = "global"
	global.Declare("a", aName, num, input, false)
	global.Declare("a.b", abGetprop, string_, input, true)
	global.Declare("u", uName, nil, input, false)
	global.Declare("arguments", argsName, obj, nil, false)
	global.Declare("f", fn.GetFirstChild(), ctor, input, false)
	global.Declare("C", cls.GetFirstChild(), ctor, input, false)
	v.println("declare-again " + typedScopeVectorsStr(func() string {
		return v.describe(global.Declare("a", aName, num, input, false))
	}))

	fScope := NewTypedScope(global, fn)
	v.labels[fScope] = "fScope"
	fScope.Declare("x", xName, nil, input, true)
	fBlock := newTypedScopeOwnedReservedNames(fScope, fBody, []string{"a", "z", "u"}, nil)
	v.labels[fBlock] = "fBlock"
	fBlock.Declare("a", innerA, bool_, input, false)
	fBlock.Declare("u", innerU, nil, input, false)
	clsScope := NewTypedScope(global, cls)
	v.labels[clsScope] = "clsScope"
	mScope := NewTypedScope(clsScope, mFn)
	v.labels[mScope] = "mScope"
	sbScope := NewTypedScope(clsScope, staticBlock)
	v.labels[sbScope] = "sbScope"
	fldScope := NewTypedScope(clsScope, fld)
	v.labels[fldScope] = "fldScope"
	sfldScope := NewTypedScope(clsScope, sfld)
	v.labels[sfldScope] = "sfldScope"
	bottom := TypedScope_CreateLatticeBottom(rhino.NewNode(rhino.Token_ROOT))
	v.labels[bottom] = "bottom"
	v.println("bad-child " + typedScopeVectorsStr(func() string {
		return v.label(NewTypedScope(global, script.GetChildAtIndex(4)))
	}))

	for _, s := range []*TypedScope{global, fScope, fBlock, clsScope, mScope, sbScope, fldScope, sfldScope, bottom} {
		module := jutil.ToJavaString(s.GetModule())
		v.println("scope " + v.label(s) +
			" depth=" + strconv.Itoa(s.GetDepth()) +
			" parent=" + v.label(s.GetParent()) +
			" bottom=" + scopeDumpBool(s.IsBottom()) +
			" module=" + module +
			" this=" + typedScopeVectorsStr(func() string { return jutil.ToJavaString(s.GetTypeOfThis()) }) +
			" count=" + strconv.Itoa(s.GetVarCount()) +
			" global=" + v.label(s.GetGlobalScope()))
		v.println("  vars " + jutil.CollectionToString(scopeTestVarNames(s.GetVarIterable())))
		v.println("  unbound " + jutil.CollectionToString(scopeTestVarNames(s.GetDeclarativelyUnboundVarsWithoutTypes())))
		v.println("  args " + typedScopeVectorsStr(func() string { return v.describe(s.GetArgumentsVar()) }))
		for _, name := range typedScopeVectors_NAMES {
			v.println("  " + name +
				" get=" + typedScopeVectorsStr(func() string { return v.describe(s.GetVar(name)) }) +
				" slot=" + typedScopeVectorsStr(func() string {
				slot, _ := s.GetSlot(name).(*TypedVar)
				return v.describe(slot)
			}) +
				" own=" + typedScopeVectorsStr(func() string {
				own, _ := s.GetOwnSlot(name).(*TypedVar)
				return v.describe(own)
			}) +
				" hasOwn=" + typedScopeVectorsStr(func() string { return scopeDumpBool(s.HasOwnSlot(name)) }) +
				" has=" + typedScopeVectorsStr(func() string { return scopeDumpBool(s.HasSlot(name)) }) +
				" topmost=" + typedScopeVectorsStr(func() string {
				return v.label(s.GetTopmostScopeOfEventualDeclaration(name))
			}))
		}
		for _, ns := range []string{"a", "a.b", "a.b.c", "f", "nope"} {
			v.println("  ns " + ns + " " +
				typedScopeVectorsStr(func() string { return jutil.ToJavaString(s.GetTypeThroughNamespace(ns)) }) +
				" lookup=" + typedScopeVectorsStr(func() string {
				return jutil.ToJavaString(s.LookupQualifiedName(rhino.QualifiedName_Of(ns)))
			}))
		}
	}

	ok := func(f func()) func() string { return func() string { f(); return "ok" } }
	v.println("validate fBlock " + typedScopeVectorsStr(ok(fBlock.ValidateCompletelyBuilt)))
	v.println("validate fBlock " + typedScopeVectorsStr(ok(fBlock.ValidateCompletelyBuilt)))
	v.println("validate global " + typedScopeVectorsStr(ok(global.ValidateCompletelyBuilt)))

	u := global.GetVar("u")
	u.SetType(num)
	u.ResolveType(scopeDumpReporter{})
	u.MarkEscaped()
	v.println("u " + v.describe(u) + " escaped=" + scopeDumpBool(u.IsMarkedEscaped()) +
		" once=" + scopeDumpBool(u.IsMarkedAssignedExactlyOnce()))
	u.MarkAssignedExactlyOnce()
	symbol, _ := u.GetSymbol().(*TypedVar)
	decl, _ := u.GetDeclaration().(*TypedVar)
	v.println("u once=" + scopeDumpBool(u.IsMarkedAssignedExactlyOnce()) +
		" symbol=" + scopeDumpBool(symbol == u) + " decl=" + scopeDumpBool(decl == u))
	v.println("undeclare " + typedScopeVectorsStr(ok(func() { fBlock.Undeclare(fBlock.GetVar("u")) })))
	v.println("after-undeclare " + typedScopeVectorsStr(func() string { return v.describe(fBlock.GetVar("u")) }))
	v.println("undeclare-foreign " + typedScopeVectorsStr(ok(func() { fBlock.Undeclare(global.GetVar("a")) })))
	v.println("bad-name-node " + typedScopeVectorsStr(func() string {
		return v.describe(NewTypedVar(false, "q", fn.GetSecondChild(), nil, global, 0, input))
	}))
	v.println("this-var " + typedScopeVectorsStr(func() string { return v.describe(fScope.GetVar("this")) }) +
		" type=" + typedScopeVectorsStr(func() string { return jutil.ToJavaString(fScope.GetVar("this").GetType()) }))
	v.println("super-var " + typedScopeVectorsStr(func() string { return v.describe(mScope.GetVar("super")) }))
	v.println("contains " + scopeDumpBool(global.Contains(fBlock)) + " " + scopeDumpBool(fBlock.Contains(global)) +
		" common=" + v.label(fBlock.GetCommonParent(mScope)))
	v.println("hoist " + v.label(fBlock.GetClosestHoistScope()) +
		" container=" + v.label(fBlock.GetClosestContainerScope()) +
		" sameContainer=" + scopeDumpBool(fBlock.HasSameContainerScope(fScope)))
	return v.out.String()
}

func TestTypedScopeVectors(t *testing.T) {
	want, err := os.ReadFile("testdata/typed_scope.txt")
	if err != nil {
		t.Fatal(err)
	}
	got := (&typedScopeVectors{labels: map[any]string{}}).run(t)
	if diff := scopeDumpFirstDiff(string(want), got); diff != "" {
		t.Errorf("differs from Java at %s", diff)
	}
}

func TestTypedScopeReservedNamesO1AndSlabAllocs(t *testing.T) {
	compiler := NewCompiler()
	script := compiler.ParseTestCode("function f() {}")
	root := rhino.IR_Root(rhino.IR_Root(), rhino.IR_Root(script))
	fn := script.GetFirstChild()
	body := fn.GetLastChild()

	// Exercise large (> scopeSmallVars) reserved-names set and deterministic insertion order.
	var names []string
	for i := range 24 {
		names = append(names, "resVar_"+strconv.Itoa(i))
	}
	global := TypedScope_CreateGlobalScope(root)
	fnScope := NewTypedScope(global, fn)
	blockScope := newTypedScopeOwnedReservedNames(fnScope, body, names, nil)
	for _, name := range names {
		if got := blockScope.GetTopmostScopeOfEventualDeclaration(name); got != blockScope {
			t.Fatalf("GetTopmostScopeOfEventualDeclaration(%q) = %v, want blockScope", name, got)
		}
	}
	// Declare all except index 5 and 19 in reverse order; ValidateCompletelyBuilt must report [resVar_5, resVar_19].
	for i := len(names) - 1; i >= 0; i-- {
		if i == 5 || i == 19 {
			continue
		}
		blockScope.Declare(names[i], rhino.NilNode, nil, nil, true)
	}
	var gotErr string
	jutil.Try(func() { blockScope.ValidateCompletelyBuilt() }, func(e jutil.Throwable) { gotErr = e.Error() })
	if !strings.Contains(gotErr, "[resVar_5, resVar_19]") {
		t.Fatalf("ValidateCompletelyBuilt error = %q, want [resVar_5, resVar_19]", gotErr)
	}

	// Verify slab batching on local TypedScopes: declaring 64 vars in a slab-backed local scope
	// allocates only the slab chunk(s) and index table, far fewer than 1 heap alloc per TypedVar.
	creator := NewTypedScopeCreator(compiler)
	nameNodes := make([]rhino.Node, 64)
	varNames := make([]string, 64)
	for i := range 64 {
		varNames[i] = "slabVar_" + strconv.Itoa(i)
		nameNodes[i] = rhino.IR_Name(varNames[i])
	}
	var slabFnScope *TypedScope
	allocs := testing.AllocsPerRun(20, func() {
		creator.ClearCommonState()
		s := newTypedScopeWithCreator(creator, global, fn, append([]string(nil), varNames...), nil)
		for i := range 64 {
			s.declareWithCreator(creator, varNames[i], nameNodes[i], nil, nil, true)
		}
		s.ValidateCompletelyBuilt()
		slabFnScope = s
	})
	// Without slabs/varLists, 1 scope + 64 TypedVars + ClearCommonState + reservedSet + index > 80 allocs.
	if allocs > 24 {
		t.Fatalf("slab-backed local TypedScope 64-var build allocs = %v, want <= 24", allocs)
	}
	if slabFnScope.GetVarCount() != 64 {
		t.Fatalf("slabFnScope var count = %d, want 64", slabFnScope.GetVarCount())
	}
}
