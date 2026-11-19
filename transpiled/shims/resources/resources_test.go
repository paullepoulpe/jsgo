package resources_test

import (
	"io"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javazip"
	"github.com/paullepoulpe/jsgo/transpiled/shims/protobuf"
	"github.com/paullepoulpe/jsgo/transpiled/shims/resources"
)

func TestExternsZipEmbeddedAndReadable(t *testing.T) {
	in := resources.GetResourceAsStream("com.google.javascript.jscomp", "/externs.zip")
	if in == nil {
		t.Fatal("expected /externs.zip resource stream")
	}
	defer in.Close()

	zis := javazip.NewZipInputStream(in)
	defer zis.Close()

	jsEntries := 0
	totalBytes := 0
	for {
		entry, err := zis.GetNextEntry()
		if err != nil {
			t.Fatalf("GetNextEntry failed: %v", err)
		}
		if entry == nil {
			break
		}
		if strings.HasSuffix(entry.GetName(), ".js") {
			jsEntries++
			data, err := io.ReadAll(zis)
			if err != nil {
				t.Fatalf("reading zip entry %s failed: %v", entry.GetName(), err)
			}
			totalBytes += len(data)
		}
		if err := zis.CloseEntry(); err != nil {
			t.Fatalf("CloseEntry failed: %v", err)
		}
	}
	if jsEntries < 20 || totalBytes < 100000 {
		t.Fatalf("unexpectedly small externs.zip: entries=%d bytes=%d", jsEntries, totalBytes)
	}
}

func TestRuntimeLibsTypedAstEndToEndDeserialization(t *testing.T) {
	in := resources.GetResourceAsStream("com.google.javascript.jscomp", "/runtime_libs.typedast")
	if in == nil {
		t.Fatal("expected /runtime_libs.typedast resource stream")
	}
	defer in.Close()

	cis, err := protobuf.CodedInputStream_NewInstance(in)
	if err != nil {
		t.Fatalf("CodedInputStream_NewInstance failed: %v", err)
	}

	// Mirror the exact TypedAstDeserializer.java loop, including Build() + Clear() reuse.
	typedAstBuilder := protobuf.TypedAst_NewBuilder()
	shards := 0
	totalScripts := 0
	for !cis.IsAtEnd() {
		tag, err := cis.ReadTag()
		if err != nil {
			t.Fatalf("ReadTag failed: %v", err)
		}
		if protobuf.WireFormat_GetTagFieldNumber(tag) != protobuf.TypedAst_List_TYPED_ASTS_FIELD_NUMBER ||
			protobuf.WireFormat_GetTagWireType(tag) != protobuf.WireFormat_WIRETYPE_LENGTH_DELIMITED {
			t.Fatalf("unexpected tag in runtime_libs.typedast: %d", tag)
		}
		if err := cis.ReadMessage(typedAstBuilder, protobuf.ExtensionRegistry_GetEmptyRegistry()); err != nil {
			t.Fatalf("ReadMessage failed: %v", err)
		}
		typedAst := typedAstBuilder.Build()
		typedAstBuilder.Clear()
		cis.ResetSizeCounter()

		shards++
		codeAsts := typedAst.GetCodeAstList()
		if len(codeAsts) == 0 && len(typedAst.GetExternAstList()) == 0 {
			t.Fatalf("shard %d was empty (check Build/Clear isolation!)", shards)
		}
		if len(codeAsts) > 0 && totalScripts == 0 {
			root, err := protobuf.AstNode_ParseFromCodedInput(
				codeAsts[0].GetScriptByteString().NewCodedInput(),
				protobuf.ExtensionRegistry_GetEmptyRegistry(),
			)
			if err != nil {
				t.Fatalf("AstNode_ParseFromCodedInput failed on runtime lib script: %v", err)
			}
			if root.GetKind() != protobuf.NodeKind_SOURCE_FILE {
				t.Fatalf("expected NodeKind_SOURCE_FILE, got %v", root.GetKind())
			}
		}
		totalScripts += len(codeAsts)
	}

	if shards == 0 || totalScripts < 50 {
		t.Fatalf("unexpected runtime_libs.typedast counts: shards=%d scripts=%d", shards, totalScripts)
	}
}

func TestTextResourcesAndParserConfigResourceBundle(t *testing.T) {
	polyfills := resources.LoadTextResource("com.google.javascript.jscomp", "js/polyfills.txt")
	if !strings.Contains(polyfills, "Array.prototype") {
		t.Fatalf("unexpected polyfills.txt content: %d bytes", len(polyfills))
	}

	transpilationLibs := resources.LoadTextResource("com.google.javascript.jscomp.js", "transpilation_libs.txt")
	if len(transpilationLibs) == 0 {
		t.Fatal("expected non-empty transpilation_libs.txt")
	}

	if !resources.ResourceExists("com.google.javascript.jscomp", "js/base.js") {
		t.Fatal("expected js/base.js to exist")
	}
	if resources.ResourceExists("com.google.javascript.jscomp", "js/does_not_exist.js") {
		t.Fatal("expected js/does_not_exist.js not to exist")
	}
	if resources.GetResourceAsStream("com.google.javascript.jscomp", "js/does_not_exist.js") != nil {
		t.Fatal("expected nil InputStream for missing resource")
	}

	bundle := resources.ResourceBundle_GetBundle("com.google.javascript.jscomp.parsing.ParserConfig")
	for _, key := range []string{"jsdoc.annotations", "jsdoc.suppressions", "jsdoc.primitives", "compiler.reserved.vars"} {
		val := bundle.GetString(key)
		if strings.TrimSpace(val) == "" {
			t.Fatalf("expected non-empty ParserConfig key %q", key)
		}
	}
}
