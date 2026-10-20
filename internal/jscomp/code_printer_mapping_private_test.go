package jscomp

import "testing"

// Candidate-only ownership assertions; excluded from stock positive controls.
func TestCodePrinterPrivateMappingSinkAndEscape(t *testing.T) {
	root, opts := mappingLifetimeFixture()
	builder := NewCodePrinter_Builder(root).SetCompilerOptions(opts)
	result := builder.buildWithSourceMappingsForCompiler()
	count := 0
	for _, m := range result.mappings {
		if m.node == nil {
			t.Fatal("private sink mapping lost its node")
		}
		count++
	}
	if count <= 544 {
		t.Fatal("private sink fixture did not cross pool chunks")
	}
	remaining := 0
	for range result.mappings {
		remaining++
	}
	if remaining != 0 {
		t.Fatal("private single-consume sink did not release its storage")
	}
	// Even an explicitly requested private lane must stay stable if a custom
	// factory receives its consumer or GetSourceMappings exports pointers.
	factory := &mappingLifetimeFactory{options: opts}
	result = NewCodePrinter_Builder(root).SetCompilerOptions(opts).SetCodeGeneratorFactory(factory).buildWithSourceMappingsForCompiler()
	base := factory.printer.GetCodePrinter_MappedCodePrinterBase()
	if base.releaseMappingsAfterIteration {
		t.Fatal("custom factory received recyclable mappings")
	}
	base.releaseMappingsAfterIteration = true
	retained := factory.printer.GetSourceMappings(result.source)
	if base.releaseMappingsAfterIteration || len(retained) <= 544 {
		t.Fatal("public mapping export did not revoke release")
	}
	base.iterMappings(func(_ int, m *SourceMap_Mapping) bool { return m.node != nil })
	if retained[544].node == nil {
		t.Fatal("escaped mapping was recycled")
	}
}
