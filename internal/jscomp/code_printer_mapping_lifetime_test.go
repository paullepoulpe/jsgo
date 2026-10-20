package jscomp

import (
	"fmt"
	"sync"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func mappingLifetimeFixture() (*rhino.Node, *CompilerOptions) {
	file := rhino.NewSimpleSourceFile("mapping-lifetime.js", rhino.StaticSourceFile_SourceKind_STRONG)
	root := rhino.NewNode(rhino.Token_SCRIPT).SetStaticSourceFile(file)
	root.SetLinenoCharno(1, 0)
	for i := 0; i < 700; i++ {
		name := rhino.Node_NewStringWithToken(rhino.Token_NAME, fmt.Sprintf("value%d", i)).SetStaticSourceFile(file)
		name.SetLinenoCharno(i+1, 0)
		name.SetLength(8)
		expr := rhino.NewNodeWithChild(rhino.Token_EXPR_RESULT, name).SetStaticSourceFile(file)
		expr.SetLinenoCharno(i+1, 0)
		root.AddChildToBack(expr)
	}
	opts := NewCompilerOptions()
	opts.SetSourceMapOutputPath(jutil.StringPtr("out.map"))
	opts.SetSourceMapDetailLevel(SourceMap_DetailLevel_ALL)
	return root, opts
}

// These controls intentionally use only APIs/fields already present on stock1fe.
// Range works for the old slice and the new sequence, without implementation tags.
func TestCodePrinterMappingsRepeatAndRetain(t *testing.T) {
	root, opts := mappingLifetimeFixture()
	result := NewCodePrinter_Builder(root).SetCompilerOptions(opts).BuildWithSourceMappings()
	var retained []*SourceMap_Mapping
	var values []SourceMap_Mapping
	for _, m := range result.mappings {
		if m == nil || m.node == nil {
			t.Fatal("missing mapping node")
		}
		retained = append(retained, m)
		values = append(values, *m)
	}
	if len(values) <= 544 {
		t.Fatalf("fixture did not cross multiple mapping chunks: %d", len(values))
	}
	for i, m := range retained {
		if *m != values[i] {
			t.Fatalf("retained mapping%d changed after iteration", i)
		}
	}
	copyResult := *result
	count := 0
	for i, m := range copyResult.mappings {
		if i >= len(values) || m != retained[i] || *m != values[i] {
			t.Fatalf("copied result mapping%d changed", i)
		}
		count++
	}
	if count != len(values) {
		t.Fatalf("second iteration count=%d, want%d", count, len(values))
	}
}

func TestCodePrinterMappingsEarlyBreakAndPanicRemainRepeatable(t *testing.T) {
	root, opts := mappingLifetimeFixture()
	result := NewCodePrinter_Builder(root).SetCompilerOptions(opts).BuildWithSourceMappings()
	var first *SourceMap_Mapping
	for _, m := range result.mappings {
		first = m
		break
	}
	if first == nil {
		t.Fatal("empty mappings")
	}
	before := *first
	marker := new(int)
	func() {
		defer func() {
			if got := recover(); got != marker {
				t.Fatalf("callback panic changed: %v", got)
			}
		}()
		for _, m := range result.mappings {
			if m != first || *m != before {
				t.Fatal("iteration did not restart at original mapping")
			}
			panic(marker)
		}
	}()
	for pass := 0; pass < 2; pass++ {
		count := 0
		for _, m := range result.mappings {
			if m == nil || m.node == nil {
				t.Fatal("mapping was recycled after interrupted iteration")
			}
			count++
		}
		if count <= 544 {
			t.Fatalf("pass%d lost mappings: %d", pass, count)
		}
	}
}

type mappingLifetimeFactory struct {
	options *CompilerOptions
	printer CodePrinter_MappedCodePrinter
}

func (f *mappingLifetimeFactory) GetCodeGenerator(_ CodePrinter_Format, cc CodeConsumer) CodeGenerator {
	f.printer = cc.(CodePrinter_MappedCodePrinter)
	return NewCodeGeneratorWithOptions(cc, f.options)
}

func TestCodePrinterCustomFactoryRetainsMappings(t *testing.T) {
	for _, build := range []string{"Build", "BuildWithSourceMappings", "ToSource"} {
		t.Run(build, func(t *testing.T) {
			root, opts := mappingLifetimeFixture()
			factory := &mappingLifetimeFactory{options: opts}
			builder := NewCodePrinter_Builder(root).SetCompilerOptions(opts).SetCodeGeneratorFactory(factory)
			var result *CodePrinter_SourceAndMappings
			var source string
			switch build {
			case "Build":
				source = builder.Build()
			case "BuildWithSourceMappings":
				result = builder.BuildWithSourceMappings()
				source = result.source
			case "ToSource":
				result = CodePrinter_ToSource(root, CodePrinter_Format_COMPACT, opts, nil, false, false, false, factory)
				source = result.source
			}
			mappings := factory.printer.GetSourceMappings(source)
			if len(mappings) <= 544 {
				t.Fatalf("custom callback lost requested source mappings: %d", len(mappings))
			}
			values := make([]SourceMap_Mapping, len(mappings))
			for i, m := range mappings {
				values[i] = *m
			}
			if result != nil {
				for _, m := range result.mappings {
					if m.node == nil {
						t.Fatal("missing source node during consume")
					}
				}
			}
			again := factory.printer.GetSourceMappings(source)
			if len(again) != len(mappings) {
				t.Fatal("retained custom printer lost mapping list")
			}
			for i, m := range mappings {
				if *m != values[i] || again[i] != m {
					t.Fatalf("retained custom mapping%d changed", i)
				}
			}
			if result != nil {
				factory.printer.StartSourceMapping(root)
				factory.printer.EndSourceMapping(root)
				count := 0
				for range result.mappings {
					count++
				}
				if count != len(mappings) {
					t.Fatal("retained printer append changed returned result bounds")
				}
			}
		})
	}
}

func TestCodePrinterParallelPublicMappingsRemainStable(t *testing.T) {
	type observation struct {
		source   string
		mappings []*SourceMap_Mapping
		values   []SourceMap_Mapping
	}
	var results [4]observation
	var group sync.WaitGroup
	for worker := range results {
		group.Add(1)
		go func() {
			defer group.Done()
			root, opts := mappingLifetimeFixture()
			result := NewCodePrinter_Builder(root).SetCompilerOptions(opts).BuildWithSourceMappings()
			results[worker].source = result.source
			for _, m := range result.mappings {
				results[worker].mappings = append(results[worker].mappings, m)
				results[worker].values = append(results[worker].values, *m)
			}
		}()
	}
	group.Wait()
	for worker, result := range results {
		if result.source != results[0].source || len(result.mappings) <= 544 {
			t.Fatalf("worker%d output/count differs", worker)
		}
		for i, m := range result.mappings {
			if *m != result.values[i] {
				t.Fatalf("worker%d retained mapping%d was overwritten", worker, i)
			}
		}
	}
}
