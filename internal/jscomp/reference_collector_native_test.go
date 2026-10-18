package jscomp

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func nativeReferenceTestContext(scope *Scope, input *CompilerInput) *NodeTraversal {
	return &NodeTraversal{scopes: []any{scope}, inputId: input.GetInputId(), compilerInput: input}
}

func TestNativeReferenceLookupPreservesScopedEqualityAndIdentity(t *testing.T) {
	t.Setenv("JSGO_EXPERIMENT_NATIVE_REFERENCES", "1")
	root := rhino.IR_Root()
	scope, equivalentScope := Scope_CreateGlobalScope(root), Scope_CreateGlobalScope(root)
	otherScope := Scope_CreateGlobalScope(rhino.IR_Root())
	first := NewVar("same", nil, scope, 0, nil, nil)
	equal := NewVar("same", nil, equivalentScope, 0, nil, nil)
	other := NewVar("same", nil, otherScope, 0, nil, nil)
	empty := NewVar("", nil, scope, 1, nil, nil)
	vars := []*Var{first, equal, other, first, empty, nil}
	inputA := NewCompilerInput(SourceFile_FromCode("a.js", ""))
	inputB := NewCompilerInput(SourceFile_FromCode("b.js", ""))
	contexts := []*NodeTraversal{nativeReferenceTestContext(scope, inputA), nativeReferenceTestContext(equivalentScope, inputB)}
	block := NewBasicBlock(nil, root)
	native := NewReferenceCollector(nil, ReferenceCollector_DO_NOTHING_BEHAVIOR, nil)
	baseline := NewReferenceCollector(nil, ReferenceCollector_DO_NOTHING_BEHAVIOR, nil)
	baseline.nativeReferences = nil
	native.blockStack, baseline.blockStack = []*BasicBlock{block}, []*BasicBlock{block}
	for index, variable := range vars {
		node := rhino.IR_Name(fmt.Sprintf("reference%d", index))
		context := contexts[index%2]
		native.AddReference(variable, node, context)
		baseline.AddReference(variable, node, context)
	}
	if native.nativeReferences == nil {
		t.Fatal("test failed to enable native lookup")
	}
	var gotSymbols, wantSymbols []*Var
	for symbol := range native.GetAllSymbols() {
		gotSymbols = append(gotSymbols, symbol)
	}
	for symbol := range baseline.GetAllSymbols() {
		wantSymbols = append(wantSymbols, symbol)
	}
	if !reflect.DeepEqual(gotSymbols, wantSymbols) || !reflect.DeepEqual(gotSymbols, []*Var{first, other, empty, nil}) {
		t.Fatalf("first symbol identities/order changed: got %v, want %v", gotSymbols, wantSymbols)
	}
	for _, variable := range vars {
		got, want := native.GetReferences(variable).references.ToSlice(), baseline.GetReferences(variable).references.ToSlice()
		if len(got) != len(want) {
			t.Fatal("reference count changed")
		}
		for index := range want {
			if got[index].nameNode != want[index].nameNode || got[index].scope != want[index].scope ||
				got[index].inputId != want[index].inputId || got[index].basicBlock != want[index].basicBlock {
				t.Fatal("reference node/scope/input/block identity or order changed")
			}
		}
	}
	if native.GetReferences(first) != native.GetReferences(equal) || native.GetReferences(first) == native.GetReferences(other) {
		t.Fatal("scoped-name equality changed")
	}
}

func TestNativeReferenceLookupRawMapEscapeRetiresIndex(t *testing.T) {
	t.Setenv("JSGO_EXPERIMENT_NATIVE_REFERENCES", "1")
	root := rhino.IR_Root()
	scope := Scope_CreateGlobalScope(root)
	variable := NewVar("x", nil, scope, 0, nil, nil)
	input := NewCompilerInput(SourceFile_FromCode("test.js", ""))
	context := nativeReferenceTestContext(scope, input)
	collector := NewReferenceCollector(nil, ReferenceCollector_DO_NOTHING_BEHAVIOR, nil)
	collector.blockStack = []*BasicBlock{NewBasicBlock(nil, root)}
	collector.AddReference(variable, rhino.IR_Name("first"), context)
	retained := collector.referenceMapWrapper()
	wrapper := collector.referenceMapWrapper()
	raw := wrapper.GetRawReferenceMap()
	if collector.nativeReferences != nil {
		t.Fatal("raw-map escape left a possibly stale index active")
	}
	replacement := NewReferenceCollection()
	raw.Put(variable, replacement)
	collector.AddReference(variable, rhino.IR_Name("second"), context)
	if collector.GetReferences(variable) != replacement || retained.GetReferences(variable) != replacement || replacement.references.Size() != 1 {
		t.Fatal("replacement was hidden by a stale index or retained wrapper")
	}
	raw.Remove(variable)
	collector.AddReference(variable, rhino.IR_Name("third"), context)
	if collector.GetReferences(variable) == replacement || collector.GetReferences(variable).references.Size() != 1 {
		t.Fatal("removal was hidden by a stale index")
	}
	raw.Clear()
	raw.Put(variable, nil)
	collector.AddReference(variable, rhino.IR_Name("fourth"), context)
	if collector.GetReferences(variable) == nil || collector.GetReferences(variable).references.Size() != 1 || raw.Size() != 1 {
		t.Fatal("clear/null replacement semantics changed")
	}
	plain := NewReferenceCollector_ReferenceMapWrapper(raw)
	if plain.GetRawReferenceMap() != raw || plain.GetReferences(variable) != collector.GetReferences(variable) || plain.String() != wrapper.String() {
		t.Fatal("public wrapper changed its authoritative map")
	}
}

func TestNativeReferenceLookupPreservesNullRootException(t *testing.T) {
	variable := NewVar("invalid", nil, &Scope{}, 0, nil, nil)
	for _, native := range []bool{false, true} {
		collector := NewReferenceCollector(nil, ReferenceCollector_DO_NOTHING_BEHAVIOR, nil)
		if native {
			collector.nativeReferences = make(map[referenceCollectorScopedKey]*ReferenceCollection)
		} else {
			collector.nativeReferences = nil
		}
		failed := false
		jutil.Try(func() { collector.GetReferences(variable) }, func(jutil.INullPointerException) { failed = true })
		if !failed {
			t.Fatalf("native=%v accepted a non-null Var with null scope root", native)
		}
	}
}

func TestNativeReferenceLookupMatchesKnownConsumers(t *testing.T) {
	programs := []string{
		`function f(a) { var x = 1; var y = x; return y + a; } use(f(2));`,
		`function outer(a) { return later(a); function later(b) { var x = b + 1; return x + a; } } use(outer(2));`,
		`function f(a) { var x = a; function inner() { var y = x; return y; } return inner(); } use(f(2));`,
		`function f() { var object = {x: side(), y: 2}; use(object.x, object.y); } f();`,
		`function f(a, b = a) { try { throw b; } catch (e) { let x = e; use(x); } return arguments[0]; } f(2);`,
		`function f(a) { let x = a; { let y = x; use(y); } var z; z = x + 1; return z; } f(2);`,
		`var named = function self(a) { var x = a; return x ? self(x - 1) : x; }; use(named(2));`,
	}
	for _, consumer := range []string{"infer", "variables", "objects", "denormalize", "checks"} {
		for index, program := range programs {
			t.Run(fmt.Sprintf("%s/case%d", consumer, index), func(t *testing.T) {
				render := func(native bool) string {
					value := ""
					if native {
						value = "1"
					}
					t.Setenv("JSGO_EXPERIMENT_NATIVE_REFERENCES", value)
					options := NewCompilerOptions()
					options.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
					compiler := NewCompiler()
					compiler.Init([]*SourceFile{SourceFile_FromCode("externs.js", "var external;")},
						[]*SourceFile{SourceFile_FromCode("input.js", program), SourceFile_FromCode("second.js", "use(external);")}, options)
					compiler.ParseInputs()
					externs, root := compiler.GetExternsRoot(), compiler.GetJsRoot()
					Normalize_CreateNormalizeForOptimizations(compiler).Process(externs, root)
					switch consumer {
					case "infer":
						NewInferConsts(compiler).Process(externs, root)
					case "variables":
						NewInlineVariables(compiler, InlineVariables_Mode_ALL).Process(externs, root)
					case "objects":
						NewInlineObjectLiterals(compiler, compiler.GetUniqueNameIdSupplier()).Process(externs, root)
					case "denormalize":
						NewDenormalize(compiler, options.GetOutputFeatureSet()).Process(externs, root)
					case "checks":
						NewVariableReferenceCheck(compiler).Process(externs, root)
					}
					return root.ToStringTree() + "\n" + compiler.ToSource() + "\n" + fmt.Sprint(compiler.GetErrors(), compiler.GetWarnings())
				}
				want := render(false)
				if got := render(true); got != want {
					t.Fatalf("native lookup changed AST/output/diagnostics:\n%s\nwant:\n%s", got, want)
				}
			})
		}
	}
}
