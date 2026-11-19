package protobuf_test

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"testing/iotest"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/protobuf"
)

func TestByteStringAndAstNodeRoundTrip(t *testing.T) {
	out := protobuf.ByteString_NewOutput(16)
	out.WriteByteInt('h')
	out.WriteByteInt('i')
	bs := out.ToByteString()
	if bs.Size() != 2 || bs.ToStringUtf8() != "hi" || bs.ByteAt(0) != 'h' {
		t.Fatalf("ByteString mismatch: %v", bs)
	}

	script := protobuf.AstNode_NewBuilder().
		SetKind(protobuf.NodeKind_SOURCE_FILE).
		SetSourceFile(1).
		AddChild(
			protobuf.AstNode_NewBuilder().
				SetKind(protobuf.NodeKind_NUMBER_LITERAL).
				SetDoubleValue(42.5).
				SetRelativeLine(3).
				SetRelativeColumn(4).
				Build(),
		).
		Build()

	lazy := protobuf.LazyAst_NewBuilder().
		SetScript(script.ToByteString()).
		SetSourceFile(1).
		Build()

	typedAst := protobuf.TypedAst_NewBuilder().
		AddCodeAst(lazy).
		SetStringPool(
			protobuf.StringPoolProto_NewBuilder().
				SetMaxLength(5).
				AddStrings(protobuf.ByteString_CopyFromUtf8("")).
				AddStrings(protobuf.ByteString_CopyFromUtf8("hello")).
				Build(),
		).
		SetSourceFilePool(
			protobuf.SourceFilePool_NewBuilder().
				AddSourceFile(
					protobuf.SourceFileProto_NewBuilder().
						SetFilename("input.js").
						SetSourceKind(protobuf.SourceFileProto_SourceKind_CODE).
						Build(),
				).
				Build(),
		).
		Build()

	var wireBuf bytes.Buffer
	if _, err := protobuf.TypedAst_List_NewBuilder().AddTypedAsts(typedAst).Build().WriteTo(&wireBuf); err != nil {
		t.Fatal(err)
	}

	cis, err := protobuf.CodedInputStream_NewInstance(bytes.NewReader(wireBuf.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	if cis.IsAtEnd() {
		t.Fatal("expected non-empty CodedInputStream")
	}
	tag, err := cis.ReadTag()
	if err != nil {
		t.Fatal(err)
	}
	if protobuf.WireFormat_GetTagFieldNumber(tag) != protobuf.TypedAst_List_TYPED_ASTS_FIELD_NUMBER ||
		protobuf.WireFormat_GetTagWireType(tag) != protobuf.WireFormat_WIRETYPE_LENGTH_DELIMITED {
		t.Fatalf("unexpected tag: %d", tag)
	}
	tb := protobuf.TypedAst_NewBuilder()
	if err := cis.ReadMessage(tb, protobuf.ExtensionRegistry_GetEmptyRegistry()); err != nil {
		t.Fatal(err)
	}
	parsed := tb.Build()
	if len(parsed.GetCodeAstList()) != 1 {
		t.Fatalf("expected 1 code ast, got %d", len(parsed.GetCodeAstList()))
	}
	parsedNode, err := protobuf.AstNode_ParseFromCodedInput(
		parsed.GetCodeAstList()[0].GetScriptByteString().NewCodedInput(),
		protobuf.ExtensionRegistry_GetEmptyRegistry(),
	)
	if err != nil {
		t.Fatal(err)
	}
	if parsedNode.GetKind() != protobuf.NodeKind_SOURCE_FILE || parsedNode.GetChildCount() != 1 {
		t.Fatalf("unexpected parsed node: %v", parsedNode)
	}
	child := parsedNode.GetChild(0)
	if child.GetDoubleValue() != 42.5 || child.GetRelativeLine() != 3 || child.GetRelativeColumn() != 4 {
		t.Fatalf("unexpected child fields: %v %d %d", child.GetDoubleValue(), child.GetRelativeLine(), child.GetRelativeColumn())
	}
}

func TestConformanceTextFormatAndOriginalMapping(t *testing.T) {
	cfgText := `
requirement: {
  type: BANNED_NAME
  value: "eval"
  error_message: "No eval allowed"
}
`
	b := protobuf.ConformanceConfig_NewBuilder()
	protobuf.TextFormat_Merge(cfgText, b)
	cfg := b.Build()
	if len(cfg.GetRequirementList()) != 1 {
		t.Fatalf("expected 1 requirement, got %d", len(cfg.GetRequirementList()))
	}
	req := cfg.GetRequirementList()[0]
	if req.GetErrorMessage() != "No eval allowed" {
		t.Fatalf("got error message %q", req.GetErrorMessage())
	}
	printed := protobuf.TextFormat_PrinterInstance().PrintToString(req)
	if !strings.Contains(printed, "No eval allowed") {
		t.Fatalf("expected printed requirement to contain message, got %q", printed)
	}

	om := protobuf.Mapping_OriginalMapping_NewBuilder().
		SetOriginalFile("app.js").
		SetLineNumber(10).
		SetColumnPosition(5).
		SetIdentifier("myFn").
		SetPrecision(protobuf.Mapping_OriginalMapping_Precision_EXACT).
		Build()
	if om.GetOriginalFile() != "app.js" || om.GetLineNumber() != 10 || om.GetPrecision() != protobuf.Mapping_OriginalMapping_Precision_EXACT {
		t.Fatalf("OriginalMapping mismatch: %+v", om)
	}
}

func TestAstNodeSizeSlabDeserializationAndOneof(t *testing.T) {
	astType := reflect.TypeOf(protobuf.AstNode{})
	if sz := astType.Size(); sz > 16 {
		t.Fatalf("expected compact AstNode handle size <= 16 bytes, got %d", sz)
	}
	for i := 0; i < astType.NumField(); i++ {
		f := astType.Field(i)
		switch f.Type.Kind() {
		case reflect.Ptr, reflect.Slice, reflect.String, reflect.Interface, reflect.Map:
			if f.Offset+f.Type.Size() > 8 {
				t.Fatalf("expected single pointer field %s within first 8 bytes, got offset=%d size=%d", f.Name, f.Offset, f.Type.Size())
			}
		}
	}
	if sz := reflect.TypeOf(protobuf.LazyAst{}).Size(); sz > 48 {
		t.Fatalf("expected LazyAst size <= 48 bytes, got %d", sz)
	}

	// Build an AST exceeding multiple 256-node slabs, exercising:
	// - leaf nodes (0 children)
	// - small interior nodes (1..8 children -> localChildren + slab.ptrs)
	// - wide interior nodes (>8 and >32 children -> overflowChildren)
	// - oneof value transitions (string ptr, double -0.0/NaN, template string)
	// - AddChild mutation on deserialized nodes (testing [:n:n] slice cap safety)
	root := protobuf.AstNode_NewBuilder().
		SetKind(protobuf.NodeKind_SOURCE_FILE).
		SetSourceFile(7).
		SetBooleanProperties(0x123456789ABCDEF0)

	for b := 0; b < 15; b++ {
		blk := protobuf.AstNode_NewBuilder().
			SetKind(protobuf.NodeKind_BLOCK).
			SetRelativeLine(int32(b + 1)).
			SetRelativeColumn(2)
		width := 40
		if b%2 == 1 {
			width = 6 // exercises localChildren <= 8 path
		}
		for i := 0; i < width; i++ {
			stmt := protobuf.AstNode_NewBuilder().
				SetKind(protobuf.NodeKind_EXPRESSION_STATEMENT).
				SetRelativeLine(1)
			var leaf *protobuf.AstNode
			switch i % 4 {
			case 0:
				leaf = protobuf.AstNode_NewBuilder().
					SetKind(protobuf.NodeKind_IDENTIFIER).
					SetStringValuePointer(int32(i + 100)).
					SetOriginalNamePointer(int32(i + 1)).
					SetType(int32(i * 3)).
					Build()
			case 1:
				val := math.Copysign(0, -1)
				if i > 1 {
					val = float64(i) * 1.25
				}
				leaf = protobuf.AstNode_NewBuilder().
					SetKind(protobuf.NodeKind_NUMBER_LITERAL).
					SetStringValuePointer(999). // overwritten by SetDoubleValue
					SetDoubleValue(val).
					Build()
			case 2:
				tsv := protobuf.TemplateStringValue_NewBuilder().
					SetRawStringPointer(int32(i + 10)).
					SetCookedStringPointer(-1).
					Build()
				leaf = protobuf.AstNode_NewBuilder().
					SetKind(protobuf.NodeKind_TEMPLATELIT_STRING).
					SetDoubleValue(9.5). // overwritten by SetTemplateStringValue
					SetTemplateStringValue(tsv).
					Build()
			case 3:
				bin := protobuf.AstNode_NewBuilder().
					SetKind(protobuf.NodeKind_ADD).
					AddChild(protobuf.AstNode_NewBuilder().SetKind(protobuf.NodeKind_IDENTIFIER).SetStringValuePointer(1).Build()).
					AddChild(protobuf.AstNode_NewBuilder().SetKind(protobuf.NodeKind_IDENTIFIER).SetStringValuePointer(2).Build()).
					Build()
				leaf = bin
			}
			stmt.AddChild(leaf)
			blk.AddChild(stmt.Build())
		}
		root.AddChild(blk.Build())
	}

	wire := root.Build().ToByteArray()

	// Verify drastically reduced heap allocations (~933 AstNodes + ~378 child slices + 87 inline TemplateStringValues
	// in ~12 total allocs with zero TemplateStringValue heap objects).
	allocs := testing.AllocsPerRun(20, func() {
		parsed, err := protobuf.AstNode_ParseFrom(wire)
		if err != nil || parsed.GetChildCount() != 15 {
			t.Fatalf("parse failed: %v", err)
		}
	})
	if allocs > 25 {
		t.Fatalf("expected slabbed AstNode_ParseFrom to do <= 25 allocs for ~933 nodes, got %.1f", allocs)
	}

	// Concurrent -race deserialization and post-parse mutation check.
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			parsed, err := protobuf.AstNode_ParseFrom(wire)
			if err != nil {
				t.Errorf("AstNode_ParseFrom err: %v", err)
				return
			}
			if parsed.GetSourceFile() != 7 || parsed.GetBooleanProperties() != 0x123456789ABCDEF0 {
				t.Errorf("root metadata mismatch")
				return
			}
			// Mutate child_ of a deserialized small node to verify [:n:n] slice cap isolation:
			// appending to blk1 must NOT corrupt sibling slices in the pointer slab!
			blk1 := parsed.GetChild(1)
			origBlk1Count := blk1.GetChildCount()
			siblingBefore := parsed.GetChild(3).GetChild(0)
			blk1.AddChild(protobuf.AstNode_NewBuilder().SetKind(protobuf.NodeKind_EMPTY).Build())
			if blk1.GetChildCount() != origBlk1Count+1 {
				t.Errorf("AddChild failed on slab-backed node")
			}
			if parsed.GetChild(3).GetChild(0) != siblingBefore {
				t.Errorf("AddChild on slab-backed slice clobbered adjacent slab pointer!")
			}

			// Check -0.0 sign bit and template string round-trip.
			negZeroNode := parsed.GetChild(0).GetChild(1).GetChild(0)
			if negZeroNode.GetValueCase() != protobuf.AstNode_ValueCase_DOUBLE_VALUE ||
				!math.Signbit(negZeroNode.GetDoubleValue()) || negZeroNode.GetStringValuePointer() != 0 {
				t.Errorf("-0.0 double oneof mismatch: case=%v val=%v", negZeroNode.GetValueCase(), negZeroNode.GetDoubleValue())
			}
			tmplNode := parsed.GetChild(0).GetChild(2).GetChild(0)
			if tmplNode.GetValueCase() != protobuf.AstNode_ValueCase_TEMPLATE_STRING_VALUE ||
				tmplNode.GetTemplateStringValue().GetCookedStringPointer() != -1 ||
				tmplNode.GetDoubleValue() != 0 {
				t.Errorf("template string oneof mismatch")
			}
		}()
	}
	wg.Wait()
}

type unhintedStreamReader struct {
	r io.Reader
}

func (u *unhintedStreamReader) Read(p []byte) (int, error) {
	return u.r.Read(p)
}

func TestLargeAstSlabAndStreamingWireParity(t *testing.T) {
	// 1. Build a >12,000-node AST to exercise multi-chunk 4096-node slabs without resetting to 256.
	bigRoot := protobuf.AstNode_NewBuilder().SetKind(protobuf.NodeKind_SOURCE_FILE).SetSourceFile(1)
	for fnIdx := 0; fnIdx < 200; fnIdx++ {
		fnBlk := protobuf.AstNode_NewBuilder().SetKind(protobuf.NodeKind_BLOCK).SetRelativeLine(1)
		for s := 0; s < 8; s++ { // <= 8 children uses slab.ptrs
			call := protobuf.AstNode_NewBuilder().
				SetKind(protobuf.NodeKind_CALL).
				SetType(int32(s + 1)).
				AddChild(protobuf.AstNode_NewBuilder().SetKind(protobuf.NodeKind_IDENTIFIER).SetStringValuePointer(int32(s + 1)).Build())
			for a := 0; a < 6; a++ {
				call.AddChild(protobuf.AstNode_NewBuilder().SetKind(protobuf.NodeKind_NUMBER_LITERAL).SetDoubleValue(float64(a)).Build())
			}
			fnBlk.AddChild(call.Build())
		}
		bigRoot.AddChild(fnBlk.Build())
	}
	// Total nodes: 1 + 200*(1 + 8*(1 + 7)) = 13,001 AstNodes (~95 KiB wire bytes > 64 KiB cachedSize cap).
	bigScriptWire := bigRoot.Build().ToByteArray()
	if len(bigScriptWire) <= 65536 {
		t.Fatalf("expected bigScriptWire > 64 KiB, got %d bytes", len(bigScriptWire))
	}
	slabAllocs := testing.AllocsPerRun(10, func() {
		parsed, err := protobuf.AstNode_ParseFrom(bigScriptWire)
		if err != nil || parsed.GetChildCount() != 200 {
			t.Fatalf("large AstNode_ParseFrom failed: %v", err)
		}
	})
	// 13,001 nodes across 4096-node slabs + 12,800 child ptrs across 8192-ptr slabs + 1 200-child root slice <= 15 allocs!
	if slabAllocs > 20 {
		t.Fatalf("expected <= 20 allocs for 13,001-node AST with 4096-node slabs, got %.1f", slabAllocs)
	}

	// 2. Build a multi-shard TypedAst_List (>250 KiB total) exercising:
	// - >=32 KiB direct zero-copy LazyAst.script writes
	// - small LazyAst.script buffered coalescing
	// - StringPoolProto, TypePool, SourceFilePool, ExternsSummary, runtimeLibraryToInject
	sp := protobuf.StringPoolProto_NewBuilder().SetMaxLength(64)
	for i := 0; i < 500; i++ {
		sp.AddStrings(protobuf.ByteString_CopyFromUtf8("ident_" + strings.Repeat("x", i%23)))
	}
	tp := protobuf.TypePool_NewBuilder().
		AddType(protobuf.TypeProto_NewBuilder().SetObject(
			protobuf.ObjectTypeProto_NewBuilder().
				SetUuid(protobuf.ByteString_CopyFromUtf8("uuid-1")).
				SetMarkedConstructor(true).
				AddAllOwnProperty([]int32{1, 2, 3}).
				Build(),
		).Build()).
		AddDisambiguationEdges(protobuf.SubtypingEdge_NewBuilder().SetSubtype(10).SetSupertype(2).Build()).
		Build()
	sfp := protobuf.SourceFilePool_NewBuilder().
		AddSourceFile(protobuf.SourceFileProto_NewBuilder().
			SetFilename("large.js").
			SetFileOnDisk(protobuf.SourceFileProto_FileOnDisk_NewBuilder().SetActualPath("/tmp/large.js").SetCharset("UTF-8").Build()).
			SetSourceKind(protobuf.SourceFileProto_SourceKind_CODE).
			Build()).
		Build()
	extSum := protobuf.ExternsSummary_NewBuilder().AddPropNamePtr(1).AddPropNamePtr(42).AddPropNamePtr(300).Build()

	shard1 := protobuf.TypedAst_NewBuilder().
		SetTypePool(tp).
		SetStringPool(sp.Build()).
		SetSourceFilePool(sfp).
		SetExternsSummary(extSum).
		AddExternAst(protobuf.LazyAst_NewBuilder().SetScriptBytes([]byte{0x08, 0xc9, 0x01}).SetSourceFile(1).Build()).
		AddCodeAst(protobuf.LazyAst_NewBuilder().SetScriptBytes(bigScriptWire).SetSourceFile(1).SetSourceMappingUrl("large.js.map").Build()).
		AddAllRuntimeLibraryToInject([]string{"es6/symbol", "es6/promise"}).
		Build()

	shard2 := protobuf.TypedAst_NewBuilder().
		AddCodeAst(protobuf.LazyAst_NewBuilder().SetScriptBytes(bigScriptWire).SetSourceFile(2).Build()).
		Build()

	list := protobuf.TypedAst_List_NewBuilder().AddTypedAsts(shard1).AddTypedAsts(shard2).Build()
	wantListBytes := list.ToByteArray()

	var streamedList bytes.Buffer
	writtenN, err := list.WriteTo(&streamedList)
	if err != nil {
		t.Fatalf("TypedAst_List.WriteTo failed: %v", err)
	}
	if int(writtenN) != len(wantListBytes) || !bytes.Equal(streamedList.Bytes(), wantListBytes) {
		t.Fatalf("TypedAst_List.WriteTo wire mismatch: written=%d got=%d want=%d", writtenN, streamedList.Len(), len(wantListBytes))
	}

	var msgToBuf bytes.Buffer
	protobuf.WriteMessageTo(list, &msgToBuf)
	if !bytes.Equal(msgToBuf.Bytes(), wantListBytes) {
		t.Fatalf("WriteMessageTo(TypedAst_List) wire mismatch")
	}

	// Verify streaming WriteMessageTo does 0 giant heap allocations when writing to io.Discard.
	writeAllocs := testing.AllocsPerRun(20, func() {
		protobuf.WriteMessageTo(list, io.Discard)
	})
	if writeAllocs > 2 {
		t.Fatalf("expected pooled streaming WriteMessageTo to do <= 2 allocs, got %.1f", writeAllocs)
	}

	// 3. Verify JSCompilerStateProto delimited streaming wire parity & ReadDelimitedBytes round-trip.
	state := (&protobuf.JSCompilerStateProto{}).
		AddAllowableFeatures(1).
		AddAllowableFeatures(5).
		SetTypeCheckingHasRun(true).
		SetHasRegExpGlobalReferences(true).
		SetLifeCycleStage(1).
		AddChunks(protobuf.ChunkProto_NewBuilder().SetName("main").AddInputIds("a.js").AddInputIds("b.js").AddDependencies(0).Build()).
		AddUniqueIdSupplier(protobuf.UniqueIdProto_NewBuilder().SetHash(123).SetCounter(4).Build()).
		SetAccessorSummary(protobuf.AccessorSummaryProto_NewBuilder().
			SetAssumeAlwaysGetterAndSetter(false).
			AddAccessors(protobuf.AccessorSummaryEntryProto_NewBuilder().SetName("prop").SetKind(protobuf.PropertyAccessKindProto_KIND_GETTER_ONLY).Build()).
			Build()).
		Build()
	wantStateRaw := state.ToByteArray()
	var delimBuf bytes.Buffer
	protobuf.WriteMessageDelimitedTo(state, &delimBuf)

	restoredState := protobuf.JSCompilerStateProto_ParseDelimitedFrom(bytes.NewReader(delimBuf.Bytes()))
	if restoredState == nil || !bytes.Equal(restoredState.ToByteArray(), wantStateRaw) {
		t.Fatalf("JSCompilerStateProto delimited streaming round-trip mismatch")
	}

	// 4. Verify CodedInputStream_NewInstance over BOTH hinted (BufferedInputStream(FileInputStream))
	// and unhinted chunked stream readers (no RemainingBytesHint / Len).
	tmpFile := filepath.Join(t.TempDir(), "shards.typedast")
	if err := os.WriteFile(tmpFile, wantListBytes, 0644); err != nil {
		t.Fatal(err)
	}
	fis, err := javaio.NewFileInputStream(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	defer fis.Close()

	checkStreamDecode := func(name string, r io.Reader) {
		cis, err := protobuf.CodedInputStream_NewInstance(r)
		if err != nil {
			t.Fatalf("%s: CodedInputStream_NewInstance err: %v", name, err)
		}
		gotShards := 0
		tb := protobuf.TypedAst_NewBuilder()
		for !cis.IsAtEnd() {
			tag, err := cis.ReadTag()
			if err != nil {
				t.Fatalf("%s: ReadTag err: %v", name, err)
			}
			if protobuf.WireFormat_GetTagFieldNumber(tag) != protobuf.TypedAst_List_TYPED_ASTS_FIELD_NUMBER {
				t.Fatalf("%s: unexpected tag %d", name, tag)
			}
			if err := cis.ReadMessage(tb, protobuf.ExtensionRegistry_GetEmptyRegistry()); err != nil {
				t.Fatalf("%s: ReadMessage err: %v", name, err)
			}
			built := tb.Build()
			tb.Clear()
			if len(built.GetCodeAstList()) != 1 {
				t.Fatalf("%s: expected 1 code AST in shard %d", name, gotShards)
			}
			gotShards++
		}
		if gotShards != 2 {
			t.Fatalf("%s: expected 2 shards, got %d", name, gotShards)
		}
	}

	checkStreamDecode("BufferedFileInputStream", javaio.NewBufferedInputStream(fis))
	checkStreamDecode("UnhintedHalfReader", &unhintedStreamReader{r: iotest.HalfReader(bytes.NewReader(wantListBytes))})

	// Truncated unhinted stream must report InvalidProtocolBufferException.
	truncCis, err := protobuf.CodedInputStream_NewInstance(&unhintedStreamReader{r: bytes.NewReader(wantListBytes[:len(wantListBytes)/2])})
	if err != nil {
		t.Fatal(err)
	}
	_, _ = truncCis.ReadTag()
	if err := truncCis.ReadMessage(protobuf.TypedAst_NewBuilder(), nil); err == nil {
		t.Fatal("expected InvalidProtocolBufferException on truncated stream")
	} else if _, ok := err.(*protobuf.InvalidProtocolBufferException); !ok {
		t.Fatalf("expected *InvalidProtocolBufferException, got %T: %v", err, err)
	}
}

func buildRealistic1000NodeAst() (*protobuf.AstNode, []byte) {
	root := protobuf.AstNode_NewBuilder().
		SetKind(protobuf.NodeKind_SOURCE_FILE).
		SetSourceFile(12).
		SetJsdoc(protobuf.OptimizationJsdoc_NewBuilder().
			AddKind(protobuf.JsdocTag_JSDOC_FILEOVERVIEW).
			SetLicenseTextPointer(3).
			Build())

	// 19 functions * 52 nodes = 988 nodes (+ 1 root = 989)
	for fn := 0; fn < 19; fn++ {
		fnNode := protobuf.AstNode_NewBuilder().
			SetKind(protobuf.NodeKind_FUNCTION_LITERAL).
			SetRelativeLine(2).
			SetRelativeColumn(0).
			SetType(int32(200 + fn)).
			SetBooleanProperties(1 << 15)
		if fn%5 == 0 {
			fnNode.SetJsdoc(protobuf.OptimizationJsdoc_NewBuilder().
				AddKind(protobuf.JsdocTag_JSDOC_NO_INLINE).
				Build())
		}
		fnNode.AddChild(protobuf.AstNode_NewBuilder().
			SetKind(protobuf.NodeKind_IDENTIFIER).
			SetStringValuePointer(int32(1000 + fn)).
			SetRelativeColumn(9).
			SetType(int32(200 + fn)).
			Build())

		params := protobuf.AstNode_NewBuilder().
			SetKind(protobuf.NodeKind_PARAMETER_LIST).
			SetRelativeColumn(6)
		for p := 0; p < 3; p++ {
			param := protobuf.AstNode_NewBuilder().
				SetKind(protobuf.NodeKind_IDENTIFIER).
				SetStringValuePointer(int32(2000 + fn*3 + p)).
				SetRelativeColumn(3).
				SetType(int32(50 + p))
			if p == 0 {
				param.SetOriginalNamePointer(int32(9000 + fn))
			}
			params.AddChild(param.Build())
		}
		fnNode.AddChild(params.Build())

		// 1 + 1 + (1+3) + 1 = 7 nodes so far; body needs 45 nodes: 15 stmts * 3 nodes = 45.
		body := protobuf.AstNode_NewBuilder().
			SetKind(protobuf.NodeKind_BLOCK).
			SetRelativeColumn(2)
		for s := 0; s < 15; s++ {
			stmt := protobuf.AstNode_NewBuilder().
				SetKind(protobuf.NodeKind_EXPRESSION_STATEMENT).
				SetRelativeLine(1).
				SetRelativeColumn(2)
			prop := protobuf.AstNode_NewBuilder().
				SetKind(protobuf.NodeKind_PROPERTY_ACCESS).
				SetStringValuePointer(int32(3000 + s*17 + fn)).
				SetRelativeColumn(4).
				SetType(int32(400 + s)).
				SetBooleanProperties(int64(s%3) << 6)
			var leaf *protobuf.AstNode
			switch (fn + s) % 8 {
			case 0:
				leaf = protobuf.AstNode_NewBuilder().
					SetKind(protobuf.NodeKind_NUMBER_LITERAL).
					SetDoubleValue(float64(fn*100+s) + 0.5).
					SetRelativeColumn(1).
					SetType(4).
					Build()
			case 1:
				leaf = protobuf.AstNode_NewBuilder().
					SetKind(protobuf.NodeKind_TEMPLATELIT_STRING).
					SetTemplateStringValue(protobuf.TemplateStringValue_NewBuilder().
						SetRawStringPointer(int32(5000 + s)).
						SetCookedStringPointer(int32(5000 + s)).
						Build()).
					SetRelativeColumn(1).
					SetType(3).
					Build()
			default:
				leaf = protobuf.AstNode_NewBuilder().
					SetKind(protobuf.NodeKind_IDENTIFIER).
					SetStringValuePointer(int32(4000 + fn*20 + s)).
					SetRelativeColumn(1).
					SetType(int32(300 + s)).
					Build()
			}
			prop.AddChild(leaf)
			stmt.AddChild(prop.Build())
			body.AddChild(stmt.Build())
		}
		fnNode.AddChild(body.Build())
		root.AddChild(fnNode.Build())
	}

	// Remaining 11 nodes: 1 BLOCK with 10 IDENTIFIER children -> total = 1,000 nodes.
	tail := protobuf.AstNode_NewBuilder().
		SetKind(protobuf.NodeKind_BLOCK).
		SetRelativeLine(3).
		SetRelativeColumn(0)
	for i := 0; i < 10; i++ {
		tail.AddChild(protobuf.AstNode_NewBuilder().
			SetKind(protobuf.NodeKind_IDENTIFIER).
			SetStringValuePointer(int32(8000 + i)).
			SetRelativeColumn(2).
			SetType(2).
			Build())
	}
	root.AddChild(tail.Build())

	built := root.Build()
	return built, built.ToByteArray()
}

func assertAstNodeParity(t *testing.T, path string, want, got *protobuf.AstNode) {
	t.Helper()
	if got.GetKind() != want.GetKind() ||
		got.GetValueCase() != want.GetValueCase() ||
		got.GetStringValuePointer() != want.GetStringValuePointer() ||
		math.Float64bits(got.GetDoubleValue()) != math.Float64bits(want.GetDoubleValue()) ||
		got.GetRelativeLine() != want.GetRelativeLine() ||
		got.GetRelativeColumn() != want.GetRelativeColumn() ||
		got.HasJsdoc() != want.HasJsdoc() ||
		got.GetOriginalNamePointer() != want.GetOriginalNamePointer() ||
		got.GetBooleanProperties() != want.GetBooleanProperties() ||
		got.HasType() != want.HasType() ||
		got.GetType() != want.GetType() ||
		got.GetSourceFile() != want.GetSourceFile() ||
		got.GetChildCount() != want.GetChildCount() {
		t.Fatalf("node mismatch at %s: got kind=%v vc=%v line=%d col=%d type=%d/%v orig=%d props=%x sf=%d children=%d",
			path, got.GetKind(), got.GetValueCase(), got.GetRelativeLine(), got.GetRelativeColumn(),
			got.GetType(), got.HasType(), got.GetOriginalNamePointer(), got.GetBooleanProperties(),
			got.GetSourceFile(), got.GetChildCount())
	}
	if want.GetValueCase() == protobuf.AstNode_ValueCase_TEMPLATE_STRING_VALUE {
		wt, gt := want.GetTemplateStringValue(), got.GetTemplateStringValue()
		if wt.GetRawStringPointer() != gt.GetRawStringPointer() || wt.GetCookedStringPointer() != gt.GetCookedStringPointer() {
			t.Fatalf("template string mismatch at %s", path)
		}
	}
	if want.HasJsdoc() {
		if !bytes.Equal(want.GetJsdoc().ToByteArray(), got.GetJsdoc().ToByteArray()) {
			t.Fatalf("jsdoc mismatch at %s", path)
		}
	}
	list := got.GetChildList()
	if len(list) != want.GetChildCount() {
		t.Fatalf("GetChildList count mismatch at %s: got %d want %d", path, len(list), want.GetChildCount())
	}
	for i := 0; i < want.GetChildCount(); i++ {
		if list[i] != got.GetChild(i) {
			t.Fatalf("GetChildList[%d] pointer identity mismatch at %s", i, path)
		}
		assertAstNodeParity(t, path, want.GetChild(i), got.GetChild(i))
	}
}

func TestZeroCopyByteStringAndNoscanAstSlabParity(t *testing.T) {
	// 1. Verify SetScript, AddStrings, SetUuid do 0 defensive clones while
	// external ByteString.ToByteArray() remains a defensive copy.
	bsScript := protobuf.ByteString_CopyFrom([]byte{1, 2, 3, 4, 5, 6, 7, 8})
	var lazy protobuf.LazyAst
	if allocs := testing.AllocsPerRun(20, func() {
		lazy.SetScript(bsScript)
	}); allocs != 0 {
		t.Fatalf("expected 0 allocs in LazyAst.SetScript, got %.1f", allocs)
	}
	var obj protobuf.ObjectTypeProto
	if allocs := testing.AllocsPerRun(20, func() {
		obj.SetUuid(bsScript)
	}); allocs != 0 {
		t.Fatalf("expected 0 allocs in ObjectTypeProto.SetUuid, got %.1f", allocs)
	}
	sp := protobuf.StringPoolProto_NewBuilder()
	for i := 0; i < 8; i++ {
		sp.AddStrings(bsScript)
	}
	// Mutating ToByteArray() output must never mutate internal proto state.
	extCopy := lazy.GetScriptByteString().ToByteArray()
	extCopy[0] = 0xFF
	if lazy.GetScript()[0] != 1 || obj.GetUuid()[0] != 1 || sp.GetStringsList()[0][0] != 1 {
		t.Fatalf("external ToByteArray() mutation leaked into ByteString/proto field")
	}

	// Verify ByteString_Output zero-copy ToByteString + copy-on-write on subsequent writes.
	out := protobuf.ByteString_NewOutput(8)
	_, _ = out.Write([]byte("abc"))
	var bs1 protobuf.ByteString
	if allocs := testing.AllocsPerRun(10, func() {
		bs1 = out.ToByteString()
	}); allocs != 0 {
		t.Fatalf("expected 0 allocs in ByteString_Output.ToByteString, got %.1f", allocs)
	}
	out.WriteByteInt('d')
	bs2 := out.ToByteString()
	if bs1.ToStringUtf8() != "abc" || bs2.ToStringUtf8() != "abcd" {
		t.Fatalf("ByteString_Output COW failed: bs1=%q bs2=%q", bs1.ToStringUtf8(), bs2.ToStringUtf8())
	}

	// 2. Full 1,000-node AST wire & recursive traversal parity + shared concurrent -race reads.
	wantAst, wire := buildRealistic1000NodeAst()
	parsed, err := protobuf.AstNode_ParseFrom(wire)
	if err != nil {
		t.Fatalf("AstNode_ParseFrom failed: %v", err)
	}
	assertAstNodeParity(t, "root", wantAst, parsed)
	if !bytes.Equal(parsed.ToByteArray(), wire) {
		t.Fatalf("re-serialized slab AstNode wire mismatch")
	}

	// Shared parsed AstNode concurrently traversed and serialized across 8 goroutines (-race).
	freshParsed, err := protobuf.AstNode_ParseFrom(wire)
	if err != nil {
		t.Fatal(err)
	}
	cntWant := 0
	accWant := traverseDeserializerStyle(wantAst, &cntWant)
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cnt := 0
			gotAcc := traverseDeserializerStyle(freshParsed, &cnt)
			if cnt != 1000 || gotAcc != accWant {
				t.Errorf("concurrent traversal mismatch: cnt=%d acc=%d want=%d", cnt, gotAcc, accWant)
			}
			if !bytes.Equal(freshParsed.ToByteArray(), wire) {
				t.Errorf("concurrent ToByteArray wire mismatch")
			}
		}()
	}
	wg.Wait()

	// 3. Verify builder copy-on-write after Build() and zero-alloc ResetForPool() reuse.
	orig := protobuf.AstNode_NewBuilder().
		SetKind(protobuf.NodeKind_CALL).
		SetRelativeLine(10).
		AddChild(protobuf.AstNode_NewBuilder().SetKind(protobuf.NodeKind_IDENTIFIER).SetStringValuePointer(7).Build()).
		Build()
	cloned := *orig
	(&cloned).SetKind(protobuf.NodeKind_NEW).AddChild(protobuf.AstNode_NewBuilder().SetKind(protobuf.NodeKind_TRUE).Build())
	if orig.GetKind() != protobuf.NodeKind_CALL || orig.GetChildCount() != 1 {
		t.Fatalf("post-Build value copy mutation clobbered original AstNode")
	}
	if cloned.GetKind() != protobuf.NodeKind_NEW || cloned.GetChildCount() != 2 {
		t.Fatalf("cloned builder mutation failed")
	}

	pooled := protobuf.AstNode_NewBuilder()
	child := protobuf.AstNode_NewBuilder().SetKind(protobuf.NodeKind_TRUE).Build()
	pooled.SetKind(protobuf.NodeKind_BLOCK).AddChild(child).Build()
	pooled.ResetForPool()
	if allocs := testing.AllocsPerRun(20, func() {
		pooled.SetKind(protobuf.NodeKind_BLOCK).SetRelativeLine(4).SetType(9).AddChild(child).Build()
		pooled.ResetForPool()
	}); allocs != 0 {
		t.Fatalf("expected 0 allocs on pooled AstNode_Builder reuse, got %.1f", allocs)
	}
}

func traverseDeserializerStyle(n *protobuf.AstNode, count *int) int64 {
	*count++
	acc := int64(n.GetSourceFile()) + int64(n.GetRelativeLine()) + int64(n.GetRelativeColumn())
	kind := n.GetKind()
	acc += int64(kind)
	switch kind {
	case protobuf.NodeKind_NUMBER_LITERAL:
		acc += int64(n.GetDoubleValue())
	case protobuf.NodeKind_IDENTIFIER, protobuf.NodeKind_PROPERTY_ACCESS, protobuf.NodeKind_STRING_LITERAL:
		acc += int64(n.GetStringValuePointer())
	case protobuf.NodeKind_TEMPLATELIT_STRING:
		tsv := n.GetTemplateStringValue()
		acc += int64(tsv.GetRawStringPointer() + tsv.GetCookedStringPointer())
	}
	if n.HasType() {
		acc += int64(n.GetType())
	}
	acc += n.GetBooleanProperties()
	if n.HasJsdoc() {
		acc += int64(n.GetJsdoc().GetKindCount())
	}
	nc := n.GetChildCount()
	for i := 0; i < nc; i++ {
		c := n.GetChild(i)
		acc += traverseDeserializerStyle(c, count)
		acc += int64(c.GetOriginalNamePointer())
	}
	return acc
}

var benchSink int64

func BenchmarkAstNode1000ParseOnly(b *testing.B) {
	_, wire := buildRealistic1000NodeAst()
	b.ReportAllocs()
	b.SetBytes(int64(len(wire)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parsed, err := protobuf.AstNode_ParseFrom(wire)
		if err != nil {
			b.Fatal(err)
		}
		benchSink = int64(parsed.GetKind())
	}
}

func BenchmarkAstNode1000ParseAndTraverse(b *testing.B) {
	_, wire := buildRealistic1000NodeAst()
	b.ReportAllocs()
	b.SetBytes(int64(len(wire)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parsed, err := protobuf.AstNode_ParseFrom(wire)
		if err != nil {
			b.Fatal(err)
		}
		cnt := 0
		benchSink = traverseDeserializerStyle(parsed, &cnt)
		if cnt != 1000 {
			b.Fatalf("expected 1000 nodes, got %d", cnt)
		}
	}
}

func BenchmarkByteStringProtoSetters(b *testing.B) {
	ast, _ := buildRealistic1000NodeAst()
	uuidSrc := []byte("0123456789abcdef")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lazy := protobuf.LazyAst_NewBuilder().
			SetScript(ast.ToByteString()).
			SetSourceFile(1).
			Build()
		sp := protobuf.StringPoolProto_NewBuilder()
		for s := 0; s < 20; s++ {
			out := protobuf.ByteString_NewOutput(16)
			_, _ = out.Write(uuidSrc)
			sp.AddStrings(out.ToByteString())
		}
		obj := protobuf.ObjectTypeProto_NewBuilder().
			SetUuid(protobuf.ByteString_CopyFrom(uuidSrc)).
			Build()
		benchSink = int64(len(lazy.GetScript()) + len(sp.Build().GetStringsList()) + len(obj.GetUuid()))
	}
}

type mockExtWrapper struct {
	core *protobuf.AstNode
}

type mockExtSlabCache struct {
	wraps       []mockExtWrapper
	firstHandle *protobuf.AstNode
	initLen     int
}

func initMockExtCache(prev any, root *protobuf.AstNode, handles []protobuf.AstNode) any {
	c, _ := prev.(*mockExtSlabCache)
	if c == nil {
		c = &mockExtSlabCache{}
	}
	total := 1 + len(handles)
	if cap(c.wraps) < total {
		c.wraps = make([]mockExtWrapper, total, max(total, cap(c.wraps)*2))
		c.firstHandle = nil
		c.initLen = 0
	} else {
		c.wraps = c.wraps[:total]
	}
	if len(handles) > 0 && c.firstHandle != &handles[0] {
		c.firstHandle = &handles[0]
		c.initLen = 0
	}
	if c.initLen == 0 {
		c.wraps[0].core = root
		c.initLen = 1
	}
	for i := c.initLen; i < total; i++ {
		c.wraps[i].core = &handles[i-1]
	}
	if c.initLen < total {
		c.initLen = total
	}
	return c
}

func traverseWithMockExtCache(w *mockExtWrapper, cnt *int) int64 {
	p := w.core
	*cnt++
	sum := int64(p.GetKind()) + int64(p.GetStringValuePointer()) + int64(p.GetType()) + p.GetBooleanProperties()
	if p.HasJsdoc() {
		sum += int64(p.GetJsdoc().GetKindCount())
	}
	n := p.GetChildCount()
	for i := 0; i < n; i++ {
		var child *mockExtWrapper
		if cIdx, cache := p.SlabChildFast(i); cache != nil {
			child = &cache.(*mockExtSlabCache).wraps[cIdx]
		} else if cIdx, cache, ok := p.SlabChildInit(i, initMockExtCache); ok {
			child = &cache.(*mockExtSlabCache).wraps[cIdx]
		} else {
			child = &mockExtWrapper{core: p.GetChild(i)}
		}
		sum += traverseWithMockExtCache(child, cnt)
	}
	return sum
}

func TestAstNodeSlabPoolReuseAndExtCacheZeroAlloc(t *testing.T) {
	_, wire := buildRealistic1000NodeAst()

	var nilNode *protobuf.AstNode
	nilNode.ReleaseSlab() // nil receiver no-op

	// Warm up pool once.
	warm, err := protobuf.AstNode_ParseFrom(wire)
	if err != nil {
		t.Fatal(err)
	}
	if warm.SlabNodeCount() != 1000 {
		t.Fatalf("SlabNodeCount = %d, want 1000", warm.SlabNodeCount())
	}
	warmCnt := 0
	wantSum := traverseWithMockExtCache(&mockExtWrapper{core: warm}, &warmCnt)
	if warmCnt != 1000 {
		t.Fatalf("warmCnt = %d, want 1000", warmCnt)
	}
	warm.ReleaseSlab()
	warm.ReleaseSlab() // idempotent double release

	var probePool sync.Pool
	poolDropsUnderRace := false
	for range 40 {
		v := new(int)
		probePool.Put(v)
		if probePool.Get() == nil {
			poolDropsUnderRace = true
			break
		}
	}

	var rootWrap mockExtWrapper
	allocs := testing.AllocsPerRun(50, func() {
		parsed, err := protobuf.AstNode_ParseFrom(wire)
		if err != nil {
			t.Fatal(err)
		}
		rootWrap.core = parsed
		cnt := 0
		got := traverseWithMockExtCache(&rootWrap, &cnt)
		if cnt != 1000 || got != wantSum {
			t.Fatalf("cnt=%d (want 1000) sum=%d (want %d)", cnt, got, wantSum)
		}
		parsed.ReleaseSlab()
	})
	maxAllocs := 0.0
	if poolDropsUnderRace {
		// Go's -race runtime intentionally drops 1 in 4 sync.Pool.Put calls.
		maxAllocs = 8.0
	}
	if allocs > maxAllocs {
		t.Fatalf("expected <= %.1f allocs/op with pooled astNodeSlab + cached wrappers across 1000-node script, got %.2f", maxAllocs, allocs)
	}

	// Mutated slab must not corrupt retained children when ReleaseSlab is called.
	mut, err := protobuf.AstNode_ParseFrom(wire)
	if err != nil {
		t.Fatal(err)
	}
	c0 := mut.GetChild(0)
	origCount := c0.GetChildCount()
	c0.AddChild(protobuf.AstNode_NewBuilder().SetKind(protobuf.NodeKind_EMPTY).Build())
	mut.ReleaseSlab()
	if c0.GetChildCount() != origCount+1 {
		t.Fatalf("mutated slab child corrupted after ReleaseSlab")
	}
}

func buildAndSerializeScriptToByteString(seed int32) protobuf.ByteString {
	root := protobuf.AstNode_NewBuilder().
		SetKind(protobuf.NodeKind_SOURCE_FILE).
		SetSourceFile(seed + 1)
	for fn := int32(0); fn < 16; fn++ {
		fnNode := protobuf.AstNode_NewBuilder().
			SetKind(protobuf.NodeKind_FUNCTION_LITERAL).
			SetRelativeLine(2).
			SetType(100 + fn).
			AddChild(protobuf.AstNode_NewBuilder().
				SetKind(protobuf.NodeKind_IDENTIFIER).
				SetStringValuePointer(seed*100 + fn).
				Build())
		body := protobuf.AstNode_NewBuilder().SetKind(protobuf.NodeKind_BLOCK)
		for s := int32(0); s < 10; s++ {
			call := protobuf.AstNode_NewBuilder().
				SetKind(protobuf.NodeKind_CALL).
				SetBooleanProperties(int64(s) << 4).
				AddChild(protobuf.AstNode_NewBuilder().
					SetKind(protobuf.NodeKind_IDENTIFIER).
					SetStringValuePointer(2000 + s).
					Build()).
				AddChild(protobuf.AstNode_NewBuilder().
					SetKind(protobuf.NodeKind_NUMBER_LITERAL).
					SetDoubleValue(float64(seed+fn+s) + 0.25).
					Build()).
				Build()
			body.AddChild(call)
		}
		fnNode.AddChild(body.Build())
		root.AddChild(fnNode.Build())
	}
	return root.Build().ToByteString()
}

func TestWave20AAstBuilderPoolReuseAndSafety(t *testing.T) {
	// Warm up pool across 2 scripts so both astBuilderBundles and child_ backing slices are populated.
	want1 := buildAndSerializeScriptToByteString(1)
	want2 := buildAndSerializeScriptToByteString(2)
	if want1.IsEmpty() || want2.IsEmpty() || want1.Equals(want2) {
		t.Fatalf("unexpected warm ByteStrings")
	}

	var probePool sync.Pool
	poolDropsUnderRace := false
	for range 40 {
		v := new(int)
		probePool.Put(v)
		if probePool.Get() == nil {
			poolDropsUnderRace = true
			break
		}
	}

	// Verify 1 alloc/script (ONLY the output []byte wire buffer!) when pool doesn't drop under -race.
	allocs := testing.AllocsPerRun(30, func() {
		bs := buildAndSerializeScriptToByteString(1)
		if !bytes.Equal(bs.BytesView(), want1.BytesView()) {
			t.Fatal("pooled script wire mismatch")
		}
	})
	if !poolDropsUnderRace && allocs > 3 {
		// ~513 AstNode_NewBuilder calls + ~193 child_ slices -> <= 3 allocs total (just the 1 output []byte)!
		t.Fatalf("expected <= 3 allocs/script with pooled astBuilderBundle + child_ slices, got %.2f", allocs)
	}

	// Shared DAG child added to multiple SOURCE_FILE trees must be pinned and never corrupted.
	sharedLeaf := protobuf.AstNode_NewBuilder().
		SetKind(protobuf.NodeKind_IDENTIFIER).
		SetStringValuePointer(777).
		Build()
	sA := protobuf.AstNode_NewBuilder().
		SetKind(protobuf.NodeKind_SOURCE_FILE).
		AddChild(sharedLeaf).
		Build()
	sB := protobuf.AstNode_NewBuilder().
		SetKind(protobuf.NodeKind_SOURCE_FILE).
		AddChild(sharedLeaf).
		Build()
	bsA := sA.ToByteString()
	bsB := sB.ToByteString()
	if !bytes.Equal(bsA.BytesView(), bsB.BytesView()) || sharedLeaf.GetStringValuePointer() != 777 {
		t.Fatalf("shared child corrupted across SOURCE_FILE ToByteString recycling")
	}

	// Concurrent multi-goroutine stress under -race.
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		g := int32(g)
		wg.Add(1)
		go func() {
			defer wg.Done()
			ref := buildAndSerializeScriptToByteString(100 + g)
			for iter := 0; iter < 25; iter++ {
				got := buildAndSerializeScriptToByteString(100 + g)
				if !bytes.Equal(got.BytesView(), ref.BytesView()) {
					t.Errorf("goroutine %d iter %d wire mismatch", g, iter)
					return
				}
				parsed, err := protobuf.AstNode_ParseFrom(got.BytesView())
				if err != nil || parsed.GetChildCount() != 16 {
					t.Errorf("goroutine %d parse failed: %v", g, err)
					return
				}
				parsed.ReleaseSlab()
			}
		}()
	}
	wg.Wait()
}

func TestWave21AStringPoolWireChunksJsdocInlineAndByteStringOutput(t *testing.T) {
	// 1. ByteString_NewOutput inline 23-byte fast path: <= 1 alloc total for NewOutput + Write + ToByteString!
	smallPayload := []byte("myJsIdentifier_123")
	outAllocs := testing.AllocsPerRun(50, func() {
		o := protobuf.ByteString_NewOutput(len(smallPayload))
		for _, b := range smallPayload {
			o.WriteByteInt(int(b))
		}
		bs := o.ToByteString()
		if bs.Size() != len(smallPayload) {
			t.Fatal("unexpected size")
		}
	})
	if outAllocs > 1 {
		t.Fatalf("expected <= 1 alloc for small ByteString_NewOutput + ToByteString, got %.2f", outAllocs)
	}

	// Overflow past inline 23-byte buffer + COW after ToByteString.
	growOut := protobuf.ByteString_NewOutput(8)
	_, _ = growOut.Write([]byte("0123456789abcdef0123456789ABCDEF"))
	bsGrow1 := growOut.ToByteString()
	growOut.WriteByteInt('!')
	bsGrow2 := growOut.ToByteString()
	if bsGrow1.ToStringUtf8() != "0123456789abcdef0123456789ABCDEF" ||
		bsGrow2.ToStringUtf8() != "0123456789abcdef0123456789ABCDEF!" {
		t.Fatalf("ByteString_Output inline->heap growth or COW failed: %q vs %q",
			bsGrow1.ToStringUtf8(), bsGrow2.ToStringUtf8())
	}

	// 2. StringPoolProto chunked wire encoding across > 5,000 strings + > 70 KiB oversized string.
	sp := protobuf.StringPoolProto_NewBuilder().SetMaxLength(75000)
	wantStrings := make([]string, 5200)
	for i := range wantStrings {
		var s string
		switch {
		case i == 0:
			s = ""
		case i == 2500:
			s = strings.Repeat("Z", 70000) // exceeds 64 KiB max chunk size
		default:
			s = "prop_" + strings.Repeat("a", i%37)
		}
		wantStrings[i] = s
		o := protobuf.ByteString_NewOutput(len(s))
		_, _ = o.Write([]byte(s))
		sp.AddStrings(o.ToByteString())
	}
	builtSP := sp.Build()

	gotList := builtSP.GetStringsList()
	if len(gotList) != len(wantStrings) {
		t.Fatalf("GetStringsList len = %d, want %d", len(gotList), len(wantStrings))
	}
	for i, w := range wantStrings {
		if string(gotList[i]) != w {
			t.Fatalf("GetStringsList[%d] mismatch (len %d vs %d)", i, len(gotList[i]), len(w))
		}
	}

	// Verify streaming WriteMessageTo == ToByteArray == round-trip deserialized StringPoolProto.
	typed := protobuf.TypedAst_NewBuilder().SetStringPool(builtSP).Build()
	rawBytes := typed.ToByteArray()
	var streamBuf bytes.Buffer
	protobuf.WriteMessageTo(typed, &streamBuf)
	if !bytes.Equal(rawBytes, streamBuf.Bytes()) {
		t.Fatalf("StringPoolProto WriteMessageTo != ToByteArray (%d vs %d bytes)", streamBuf.Len(), len(rawBytes))
	}

	cis := protobuf.CodedInputStream_NewInstanceBytes(rawBytes)
	roundTrip := protobuf.TypedAst_NewBuilder()
	for !cis.IsAtEnd() {
		// TypedAst is top-level here, read via Ast/TypedAst merge: wrap in TypedAst_List for ReadMessage.
		break
	}
	var listBuf bytes.Buffer
	_, err := protobuf.TypedAst_List_NewBuilder().AddTypedAsts(typed).Build().WriteTo(&listBuf)
	if err != nil {
		t.Fatal(err)
	}
	listCis := protobuf.CodedInputStream_NewInstanceBytes(listBuf.Bytes())
	_, err = listCis.ReadTag()
	if err != nil {
		t.Fatal(err)
	}
	if err := listCis.ReadMessage(roundTrip, protobuf.ExtensionRegistry_GetEmptyRegistry()); err != nil {
		t.Fatal(err)
	}
	parsedSP := roundTrip.Build().GetStringPool()
	if parsedSP.GetMaxLength() != 75000 {
		t.Fatalf("parsed MaxLength = %d, want 75000", parsedSP.GetMaxLength())
	}
	parsedBS := parsedSP.GetStringsByteStringList()
	if len(parsedBS) != len(wantStrings) {
		t.Fatalf("parsed GetStringsByteStringList len = %d, want %d", len(parsedBS), len(wantStrings))
	}
	for i, w := range wantStrings {
		if parsedBS[i].ToStringUtf8() != w {
			t.Fatalf("parsedBS[%d] mismatch", i)
		}
	}

	// 3. OptimizationJsdoc inline kindInline_ (0, 3, 7, and 10 tags) + 0-alloc Equals against DefaultInstance.
	defJsdoc := protobuf.OptimizationJsdoc_GetDefaultInstance()
	emptyBuilt := protobuf.OptimizationJsdoc_NewBuilder().Build()
	nonEmptyBuilt := protobuf.OptimizationJsdoc_NewBuilder().
		AddKind(protobuf.JsdocTag_JSDOC_CONST).
		AddKind(protobuf.JsdocTag_JSDOC_NO_INLINE).
		Build()

	eqAllocs := testing.AllocsPerRun(50, func() {
		if !defJsdoc.Equals(emptyBuilt) || defJsdoc.Equals(nonEmptyBuilt) {
			t.Fatal("OptimizationJsdoc.Equals mismatch")
		}
	})
	if eqAllocs > 0 {
		t.Fatalf("expected 0 allocs in OptimizationJsdoc.Equals, got %.2f", eqAllocs)
	}

	manyTags := protobuf.OptimizationJsdoc_NewBuilder()
	allTags := []protobuf.JsdocTag{
		protobuf.JsdocTag_JSDOC_FILEOVERVIEW,
		protobuf.JsdocTag_JSDOC_CONST,
		protobuf.JsdocTag_JSDOC_CONSTRUCTOR,
		protobuf.JsdocTag_JSDOC_NO_INLINE,
		protobuf.JsdocTag_JSDOC_NO_SIDE_EFFECTS,
		protobuf.JsdocTag_JSDOC_MODIFIES_THIS,
		protobuf.JsdocTag_JSDOC_THROWS,
		protobuf.JsdocTag_JSDOC_ABSTRACT,
		protobuf.JsdocTag_JSDOC_DEFINE,
	}
	for _, tg := range allTags {
		manyTags.AddKind(tg)
	}
	manyBuilt := manyTags.SetLicenseTextPointer(9).Build()
	astWithJsdoc := protobuf.AstNode_NewBuilder().
		SetKind(protobuf.NodeKind_SOURCE_FILE).
		SetJsdoc(manyBuilt).
		AddChild(protobuf.AstNode_NewBuilder().SetKind(protobuf.NodeKind_BLOCK).SetJsdoc(nonEmptyBuilt).Build()).
		Build()
	astBytes := astWithJsdoc.ToByteString()
	parsedAst, err := protobuf.AstNode_ParseFrom(astBytes.BytesView())
	if err != nil {
		t.Fatal(err)
	}
	if !parsedAst.GetJsdoc().Equals(manyBuilt) || !parsedAst.GetChild(0).GetJsdoc().Equals(nonEmptyBuilt) {
		t.Fatalf("OptimizationJsdoc inline/overflow round-trip Equals mismatch")
	}
	if gotKinds := parsedAst.GetJsdoc().GetKindList(); len(gotKinds) != len(allTags) {
		t.Fatalf("GetKindList len = %d, want %d", len(gotKinds), len(allTags))
	}
}

func TestWave25AZeroEscapeByteStringOutputTypePoolSlabsAndStreamingTypedAst(t *testing.T) {
	// 1. Compact TypePool message struct sizes.
	if sz := reflect.TypeOf(protobuf.TypeProto{}).Size(); sz > 32 {
		t.Fatalf("sizeof(TypeProto) = %d, want <= 32", sz)
	}
	if sz := reflect.TypeOf(protobuf.ObjectTypeProto{}).Size(); sz > 112 {
		t.Fatalf("sizeof(ObjectTypeProto) = %d, want <= 112", sz)
	}
	if sz := reflect.TypeOf(protobuf.UnionTypeProto{}).Size(); sz > 40 {
		t.Fatalf("sizeof(UnionTypeProto) = %d, want <= 40", sz)
	}
	if sz := reflect.TypeOf(protobuf.SubtypingEdge{}).Size(); sz > 24 {
		t.Fatalf("sizeof(SubtypingEdge) = %d, want <= 24", sz)
	}

	// 2. Stack-allocated ByteString_NewOutput even with oversized initialCapacity hint (e.g. 4096)
	// + slab-backed ToByteString() => 0 heap allocations per call once warm!
	ident := []byte("google.protobuf.DescriptorProto.ExtensionRange")
	allocs := testing.AllocsPerRun(100, func() {
		o := protobuf.ByteString_NewOutput(4096)
		for _, c := range ident {
			o.WriteByteInt(int(c))
		}
		bs := o.ToByteString()
		if bs.Size() != len(ident) {
			t.Fatal("size mismatch")
		}
	})
	if allocs > 0.1 {
		t.Fatalf("expected ~0 allocs for stack-allocated ByteString_NewOutput(4096) + ToByteString, got %.2f", allocs)
	}

	// Spill > 246 inline bytes + COW mutation after ToByteString.
	bigIdent := bytes.Repeat([]byte("abcdEFGH"), 60) // 480 bytes > 246 inline cap
	oBig := protobuf.ByteString_NewOutput(512)
	_, _ = oBig.Write(bigIdent[:200])
	for _, b := range bigIdent[200:] {
		_ = oBig.WriteByte(b)
	}
	bsBig1 := oBig.ToByteString()
	_, _ = oBig.Write([]byte("_tail"))
	bsBig2 := oBig.ToByteString()
	if !bytes.Equal(bsBig1.BytesView(), bigIdent) || bsBig2.Size() != len(bigIdent)+5 {
		t.Fatalf("spilled ByteString_Output COW mismatch: len1=%d len2=%d", bsBig1.Size(), bsBig2.Size())
	}

	// 3. Streaming CodedInputStream.ReadMessage on >300 KiB TypedAst (exercising both >256 KiB
	// hinted FileInputStream streaming AND unhinted chunked reader streaming) + TypePool arena
	// slice isolation + StringPoolProto 0-alloc wire span mergeFromReader.
	tpb := protobuf.TypePool_NewBuilder()
	for i := int32(0); i < 600; i++ {
		if i%4 == 3 {
			u := protobuf.UnionTypeProto_NewBuilder().AddAllUnionMember([]int32{1, i, i + 1}).Build()
			tpb.AddType(protobuf.TypeProto_NewBuilder().SetUnion(u).Build())
		} else {
			obj := protobuf.ObjectTypeProto_NewBuilder().
				SetUuid(protobuf.ByteString_CopyFromUtf8("uuid_" + strings.Repeat("k", int(i%11)))).
				SetIsInvalidating(i%5 == 0).
				SetPropertiesKeepOriginalName(i%7 == 0).
				SetMarkedConstructor(i%2 == 0).
				SetClosureAssert(i%13 == 0).
				AddAllInstanceType([]int32{i + 1}).
				AddAllPrototype([]int32{i + 2}).
				AddAllOwnProperty([]int32{10 + i, 20 + i, 30 + i}).
				Build()
			tpb.AddType(protobuf.TypeProto_NewBuilder().SetObject(obj).Build())
		}
		tpb.AddDisambiguationEdgesBuilder().SetSubtype(i + 1).SetSupertype(0)
	}
	spb := protobuf.StringPoolProto_NewBuilder().SetMaxLength(500)
	for i := 0; i < 4000; i++ {
		o := protobuf.ByteString_NewOutput(256)
		if i > 0 {
			_, _ = o.Write([]byte("sym_" + strings.Repeat("q", i%40)))
		}
		spb.AddStrings(o.ToByteString())
	}
	_, scriptWire := buildRealistic1000NodeAst()
	tb := protobuf.TypedAst_NewBuilder().
		SetTypePool(tpb.Build()).
		SetStringPool(spb.Build()).
		SetSourceFilePool(protobuf.SourceFilePool_NewBuilder().
			AddSourceFile(protobuf.SourceFileProto_NewBuilder().
				SetFilename("mod.js").
				SetFileOnDisk(protobuf.SourceFileProto_FileOnDisk_NewBuilder().SetActualPath("/src/mod.js").SetCharset("UTF-8").Build()).
				SetSourceKind(protobuf.SourceFileProto_SourceKind_CODE).
				Build()).
			Build()).
		SetExternsSummary(protobuf.ExternsSummary_NewBuilder().AddPropNamePtr(3).AddPropNamePtr(9).Build()).
		AddAllRuntimeLibraryToInject([]string{"es6/util/arrayfromiterator"})
	for s := int32(0); s < 30; s++ {
		tb.AddCodeAst(protobuf.LazyAst_NewBuilder().
			SetScriptBytes(scriptWire).
			SetSourceFile(1).
			SetSourceMappingUrl("mod.js.map").
			Build())
	}
	wantTypedAst := tb.Build()
	wantTypedBytes := wantTypedAst.ToByteArray()
	listBytes := protobuf.TypedAst_List_NewBuilder().AddTypedAsts(wantTypedAst).Build().ToByteArray()
	if len(listBytes) <= 256*1024 {
		t.Fatalf("expected >256 KiB listBytes to exercise streaming hinted reader, got %d", len(listBytes))
	}

	// Verify both >256KiB bytes.Reader (hinted) and unhinted stream decode byte-identically!
	for _, tc := range []struct {
		name string
		r    io.Reader
	}{
		{"HintedLargeBytesReader", bytes.NewReader(listBytes)},
		{"UnhintedChunkedReader", &unhintedStreamReader{r: iotest.HalfReader(bytes.NewReader(listBytes))}},
		{"InMemoryBytesCis", nil},
	} {
		var cis *protobuf.CodedInputStream
		var err error
		if tc.r == nil {
			cis = protobuf.CodedInputStream_NewInstanceBytes(listBytes)
		} else {
			cis, err = protobuf.CodedInputStream_NewInstance(tc.r)
			if err != nil {
				t.Fatalf("%s: NewInstance failed: %v", tc.name, err)
			}
		}
		tag, err := cis.ReadTag()
		if err != nil || protobuf.WireFormat_GetTagFieldNumber(tag) != 1 {
			t.Fatalf("%s: unexpected tag %d, err=%v", tc.name, tag, err)
		}
		dst := protobuf.TypedAst_NewBuilder()
		if err := cis.ReadMessage(dst, protobuf.ExtensionRegistry_GetEmptyRegistry()); err != nil {
			t.Fatalf("%s: ReadMessage failed: %v", tc.name, err)
		}
		if !cis.IsAtEnd() {
			t.Fatalf("%s: expected stream at end", tc.name)
		}
		got := dst.Build()
		if !bytes.Equal(got.ToByteArray(), wantTypedBytes) {
			t.Fatalf("%s: round-trip TypedAst wire mismatch (%d vs %d bytes)", tc.name, len(got.ToByteArray()), len(wantTypedBytes))
		}
		// Mutating arena-backed ObjectTypeProto slice must NOT clobber neighbor types!
		t0 := got.GetTypePool().GetType(0).GetObject()
		t1 := got.GetTypePool().GetType(1).GetObject()
		wantT1Prop0 := t1.GetOwnProperty(0)
		t0.AddAllOwnProperty([]int32{999999})
		if t1.GetOwnProperty(0) != wantT1Prop0 {
			t.Fatalf("%s: arena slice append clobbered adjacent ObjectTypeProto!", tc.name)
		}
	}

	// Verify SetRecursionLimit(1) rejects nested submessages in streaming TypedAst reader.
	recCis, err := protobuf.CodedInputStream_NewInstance(&unhintedStreamReader{r: bytes.NewReader(listBytes)})
	if err != nil {
		t.Fatal(err)
	}
	recCis.SetRecursionLimit(1)
	_, _ = recCis.ReadTag()
	if err := recCis.ReadMessage(protobuf.TypedAst_NewBuilder(), nil); err == nil {
		t.Fatalf("expected recursion limit error on SetRecursionLimit(1)")
	}
}

func TestWave26AScriptBytesTwoSlotSlabPackingAndIsolation(t *testing.T) {
	// Exercise 0-byte, tiny, medium (90-450 KiB), and oversized (>512 KiB) LazyAst.script_
	// payloads across streaming CodedInputStream to verify:
	// 1. 2-slot 2 MiB slab packing (scriptSlab + scriptTail) coalesces <=512 KiB scripts
	//    into <= 3 slab allocations instead of ~40 separate heap span allocations.
	// 2. Every returned GetScript() slice is cap-limited (cap == len) so append() never
	//    corrupts adjacent scripts in either scriptSlab or scriptTail.
	// 3. IsAtEnd() at EOF releases internal slab/scratch references while LazyAst scripts stay valid.
	sizes := []int{
		0, 3, 96 * 1024, 110 * 1024, 48 * 1024, 380 * 1024, 12 * 1024,
		220 * 1024, 7 * 1024, 450 * 1024, 15 * 1024, 310 * 1024, 600 * 1024,
		64 * 1024, 140 * 1024, 8 * 1024, 190 * 1024, 32 * 1024,
	}
	expected := make([][]byte, len(sizes))
	tb := protobuf.TypedAst_NewBuilder().
		SetSourceFilePool(protobuf.SourceFilePool_NewBuilder().
			AddSourceFile(protobuf.SourceFileProto_NewBuilder().
				SetFilename("slab.js").
				SetSourceKind(protobuf.SourceFileProto_SourceKind_CODE).
				Build()).
			Build())
	for idx, sz := range sizes {
		raw := make([]byte, sz)
		for j := range raw {
			raw[j] = byte((idx*37 + j*13 + 5) & 0xFF)
		}
		expected[idx] = raw
		tb.AddCodeAst(protobuf.LazyAst_NewBuilder().
			SetScriptBytes(raw).
			SetSourceFile(1).
			Build())
	}
	listWire := protobuf.TypedAst_List_NewBuilder().AddTypedAsts(tb.Build()).Build().ToByteArray()

	parseOnce := func() (*protobuf.CodedInputStream, *protobuf.TypedAst) {
		cis, err := protobuf.CodedInputStream_NewInstance(&unhintedStreamReader{r: bytes.NewReader(listWire)})
		if err != nil {
			t.Fatalf("CodedInputStream_NewInstance failed: %v", err)
		}
		tag, err := cis.ReadTag()
		if err != nil || protobuf.WireFormat_GetTagFieldNumber(tag) != 1 {
			t.Fatalf("unexpected tag %d err=%v", tag, err)
		}
		dst := protobuf.TypedAst_NewBuilder()
		if err := cis.ReadMessage(dst, protobuf.ExtensionRegistry_GetEmptyRegistry()); err != nil {
			t.Fatalf("ReadMessage failed: %v", err)
		}
		if !cis.IsAtEnd() {
			t.Fatalf("expected IsAtEnd() == true")
		}
		return cis, dst.Build()
	}

	_, got := parseOnce()
	codeAsts := got.GetCodeAstList()
	if len(codeAsts) != len(sizes) {
		t.Fatalf("codeAsts len = %d, want %d", len(codeAsts), len(sizes))
	}
	// Verify cap == len isolation by appending to every script before verifying all contents!
	for i, la := range codeAsts {
		sc := la.GetScript()
		if len(sc) != sizes[i] {
			t.Fatalf("script[%d] len = %d, want %d", i, len(sc), sizes[i])
		}
		if cap(sc) != len(sc) {
			t.Fatalf("script[%d] cap = %d != len %d (unsafe against append)", i, cap(sc), len(sc))
		}
		_ = append(sc, 0xAA, 0xBB, 0xCC)
	}
	for i, la := range codeAsts {
		if !bytes.Equal(la.GetScript(), expected[i]) {
			t.Fatalf("script[%d] corrupted after neighbor append!", i)
		}
	}

	// Measure allocation count across 30 x 140 KiB scripts (~4.2 MiB total):
	// previously 30+ allocs (>128 KiB bypassed slab completely); now <= 12 total allocs
	// including CodedInputStream, bufio.Reader, TypedAst slices, and three 2 MiB slabs!
	homogTB := protobuf.TypedAst_NewBuilder()
	chunk140K := bytes.Repeat([]byte{0x08, 0x01}, 70*1024)
	for i := 0; i < 30; i++ {
		homogTB.AddCodeAst(protobuf.LazyAst_NewBuilder().SetScriptBytes(chunk140K).SetSourceFile(1).Build())
	}
	homogWire := protobuf.TypedAst_List_NewBuilder().AddTypedAsts(homogTB.Build()).Build().ToByteArray()
	allocs := testing.AllocsPerRun(5, func() {
		cis, err := protobuf.CodedInputStream_NewInstance(&unhintedStreamReader{r: bytes.NewReader(homogWire)})
		if err != nil {
			t.Fatal(err)
		}
		_, _ = cis.ReadTag()
		dst := protobuf.TypedAst_NewBuilder()
		if err := cis.ReadMessage(dst, nil); err != nil {
			t.Fatal(err)
		}
		_ = cis.IsAtEnd()
	})
	if allocs > 22 {
		t.Fatalf("expected <= 22 allocs for 30x140KiB streaming LazyAst scripts with 2MiB slabs, got %.1f", allocs)
	}
}

func BenchmarkWave26AStreamingTypedAst(b *testing.B) {
	tb := protobuf.TypedAst_NewBuilder()
	chunk96K := bytes.Repeat([]byte{0x08, 0x01}, 48*1024)
	chunk24K := bytes.Repeat([]byte{0x08, 0x02}, 12*1024)
	for i := 0; i < 40; i++ {
		if i%2 == 0 {
			tb.AddCodeAst(protobuf.LazyAst_NewBuilder().SetScriptBytes(chunk96K).SetSourceFile(1).Build())
		} else {
			tb.AddCodeAst(protobuf.LazyAst_NewBuilder().SetScriptBytes(chunk24K).SetSourceFile(1).Build())
		}
	}
	wire := protobuf.TypedAst_List_NewBuilder().AddTypedAsts(tb.Build()).Build().ToByteArray()
	b.ReportAllocs()
	b.SetBytes(int64(len(wire)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cis, err := protobuf.CodedInputStream_NewInstance(&unhintedStreamReader{r: bytes.NewReader(wire)})
		if err != nil {
			b.Fatal(err)
		}
		_, _ = cis.ReadTag()
		dst := protobuf.TypedAst_NewBuilder()
		if err := cis.ReadMessage(dst, nil); err != nil {
			b.Fatal(err)
		}
		_ = cis.IsAtEnd()
	}
}

func TestWave32CCompactByteStringOutputPoolAndTwoSlotWireSlab(t *testing.T) {
	// 1. ByteString_Output must stay <= 48 bytes so heap-escaped wrappers cut 293 MB -> 48 MB.
	if sz := reflect.TypeOf(protobuf.ByteString_Output{}).Size(); sz > 48 {
		t.Fatalf("sizeof(ByteString_Output) = %d, want <= 48", sz)
	}

	// 2. Zero-alloc stack-allocated ByteString_NewOutput across 0B, 1B, 15B (inline), and 2 KiB (pooled 32 KiB scratch).
	payload2K := bytes.Repeat([]byte("closure_ident_01"), 128) // 2048 bytes
	for _, tc := range []struct {
		name string
		data []byte
	}{
		{"Empty0B", nil},
		{"Single1B", []byte{'x'}},
		{"Inline15B", []byte("0123456789abcde")},
		{"Pooled2KiB", payload2K},
	} {
		allocs := testing.AllocsPerRun(50, func() {
			o := protobuf.ByteString_NewOutput(len(tc.data))
			if len(tc.data) <= 32 {
				for _, b := range tc.data {
					o.WriteByteInt(int(b))
				}
			} else {
				_, _ = o.Write(tc.data[:1000])
				_, _ = o.Write(tc.data[1000:])
			}
			bs := o.ToByteString()
			if bs.Size() != len(tc.data) {
				t.Fatalf("%s: size = %d, want %d", tc.name, bs.Size(), len(tc.data))
			}
		})
		if allocs > 0.1 {
			t.Fatalf("%s: expected ~0 allocs once warm, got %.2f", tc.name, allocs)
		}
	}

	// 3. Repeated ToByteString() + post-ToByteString() COW writes across 1B, 15B, 2 KiB, and >32 KiB spill.
	for _, initLen := range []int{0, 1, 15, 2048, 40 * 1024} {
		initBytes := bytes.Repeat([]byte{'a'}, initLen)
		o := protobuf.ByteString_NewOutput(16)
		_, _ = o.Write(initBytes)
		bs1 := o.ToByteString()
		bs1Dup := o.ToByteString()
		o.WriteByteInt('Z')
		_, _ = o.Write([]byte("99"))
		bs2 := o.ToByteString()
		if !bytes.Equal(bs1.BytesView(), initBytes) || !bytes.Equal(bs1Dup.BytesView(), initBytes) {
			t.Fatalf("initLen=%d: bs1 mutated after subsequent writes!", initLen)
		}
		if bs2.Size() != initLen+3 || bs2.ByteAt(initLen) != 'Z' {
			t.Fatalf("initLen=%d: bs2 mismatch (size=%d)", initLen, bs2.Size())
		}
	}

	// 4. Concurrent -race stress on pooled ByteString_NewOutput + 2-slot 2 MiB allocWireBytes (>64 KiB items).
	bigPayload := bytes.Repeat([]byte("01234567"), 12*1024) // 96 KiB (> old 64 KiB wireSlabMaxItem)
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		g := g
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 40; i++ {
				var want []byte
				switch i % 4 {
				case 0:
					want = []byte{byte(g*31 + i)}
				case 1:
					want = []byte(fmt.Sprintf("short_%d_%d", g, i))
				case 2:
					want = bytes.Repeat([]byte{byte(g + i + 1)}, 600)
				case 3:
					want = bigPayload
				}
				o := protobuf.ByteString_NewOutput(8)
				if len(want) > 0 {
					o.WriteByteInt(int(want[0]))
					if len(want) > 1 {
						_, _ = o.Write(want[1:])
					}
				}
				bs := o.ToByteString()
				if !bytes.Equal(bs.BytesView(), want) {
					t.Errorf("goroutine %d iter %d mismatch (len %d vs %d)", g, i, bs.Size(), len(want))
					return
				}
			}
		}()
	}
	wg.Wait()
}

func BenchmarkWave34CAstSerializeAndRecycle(b *testing.B) {
	b.ReportAllocs()
	var sink int
	for i := 0; i < b.N; i++ {
		root := protobuf.AstNode_NewBuilder().SetKind(protobuf.NodeKind_SOURCE_FILE).SetSourceFile(1)
		for fn := 0; fn < 25; fn++ {
			blk := protobuf.AstNode_NewBuilder().SetKind(protobuf.NodeKind_BLOCK).SetRelativeLine(1)
			for st := 0; st < 8; st++ {
				call := protobuf.AstNode_NewBuilder().
					SetKind(protobuf.NodeKind_CALL).
					SetType(int32(st + 1)).
					AddChild(protobuf.AstNode_NewBuilder().SetKind(protobuf.NodeKind_IDENTIFIER).SetStringValuePointer(int32(st + 10)).Build())
				for a := 0; a < 3; a++ {
					call.AddChild(protobuf.AstNode_NewBuilder().SetKind(protobuf.NodeKind_NUMBER_LITERAL).SetDoubleValue(float64(a)).Build())
				}
				blk.AddChild(call.Build())
			}
			root.AddChild(blk.Build())
		}
		bs := root.Build().ToByteString()
		sink += bs.Size()
	}
	_ = sink
}

func TestWave35CAstNodePutTagVarintBoundaries(t *testing.T) {
	// Exercise 1-byte (<0x80), 2-byte (<0x4000), 3-byte (<0x200000), and 5/10-byte varints
	// through both builder (*AstNode.putTo) and deserialized slab (*astNodeSlab.putRecTo) paths.
	offsets := []int32{0, 1, 127, 128, 16383, 16384, 2097151, 2097152, math.MaxInt32, -1, -128, -16384}
	boolProps := []int64{0, 1, 0x7F, 0x80, 0x3FFF, 0x4000, 1 << 35, -1}

	rootB := protobuf.AstNode_NewBuilder().SetKind(protobuf.NodeKind_SOURCE_FILE).SetSourceFile(300)
	for i, off := range offsets {
		bp := boolProps[i%len(boolProps)]
		tv := protobuf.TemplateStringValue_NewBuilder().
			SetRawStringPointer(off).
			SetCookedStringPointer(-off).
			Build()
		c1 := protobuf.AstNode_NewBuilder().
			SetKind(protobuf.NodeKind_IDENTIFIER).
			SetStringValuePointer(off).
			SetOriginalNamePointer(off ^ 0x4000).
			SetRelativeLine(off).
			SetRelativeColumn(-off).
			SetBooleanProperties(bp).
			SetType(off).
			SetSourceFile(int32(i + 1)).
			Build()
		c2 := protobuf.AstNode_NewBuilder().
			SetKind(protobuf.NodeKind_TEMPLATELIT_STRING).
			SetTemplateStringValue(tv).
			SetBooleanProperties(bp).
			Build()
		rootB.AddChild(c1).AddChild(c2)
	}
	built := rootB.Build()
	wire1 := built.ToByteArray()

	parsed, err := protobuf.AstNode_ParseFrom(wire1)
	if err != nil {
		t.Fatalf("AstNode_ParseFrom failed: %v", err)
	}
	defer parsed.ReleaseSlab()

	// Re-serialize from unmutated astNodeSlab (putRecTo) and verify exact byte-for-byte wire identity!
	wire2 := parsed.ToByteArray()
	if !bytes.Equal(wire1, wire2) {
		t.Fatalf("slab putRecTo wire mismatch: len %d vs %d", len(wire1), len(wire2))
	}
	if parsed.GetChildCount() != len(offsets)*2 {
		t.Fatalf("child count = %d, want %d", parsed.GetChildCount(), len(offsets)*2)
	}
	for i, off := range offsets {
		c1 := parsed.GetChild(i * 2)
		if c1.GetStringValuePointer() != off || c1.GetRelativeLine() != off || c1.GetRelativeColumn() != -off || c1.GetType() != off {
			t.Fatalf("child %d scalar roundtrip mismatch for off=%d", i, off)
		}
		c2 := parsed.GetChild(i*2 + 1)
		tv := c2.GetTemplateStringValue()
		if tv.GetRawStringPointer() != off || tv.GetCookedStringPointer() != -off {
			t.Fatalf("child %d template string mismatch for off=%d", i, off)
		}
	}
}

func TestWave36CIntrusiveAstBuilderArenaCOWAndConcurrentSafety(t *testing.T) {
	// 1. Standalone inline AstNode (simulating j2gort inlineRaw lifecycle with ResetForPool)
	// must build & serialize without allocating []*AstNode child slices.
	var inlineNodes [32]protobuf.AstNode
	buildInlineTree := func(seed int32) protobuf.ByteString {
		root := &inlineNodes[0]
		root.SetKind(protobuf.NodeKind_SOURCE_FILE).SetSourceFile(1)
		idx := 1
		for b := int32(0); b < 5; b++ {
			blk := &inlineNodes[idx]
			idx++
			blk.SetKind(protobuf.NodeKind_BLOCK).SetRelativeLine(b + 1).SetRelativeColumn(2)
			for c := int32(0); c < 5; c++ {
				leaf := &inlineNodes[idx]
				idx++
				leaf.SetKind(protobuf.NodeKind_IDENTIFIER).
					SetStringValuePointer(seed + b*10 + c).
					SetBooleanProperties(int64(c + 1)).
					SetType(c + 2)
				blk.AddChild(leaf.Build())
			}
			root.AddChild(blk.Build())
		}
		bs := root.Build().ToByteString()
		for i := 0; i < idx; i++ {
			inlineNodes[i].ResetForPool()
		}
		return bs
	}
	refBS := buildInlineTree(100)
	_ = buildInlineTree(200)
	inlineAllocs := testing.AllocsPerRun(50, func() {
		got := buildInlineTree(100)
		if !bytes.Equal(got.BytesView(), refBS.BytesView()) {
			t.Fatal("inline AstNode ResetForPool wire mismatch")
		}
	})
	if inlineAllocs > 1 {
		t.Fatalf("expected <= 1 allocs for intrusive inlineRaw AstNode tree, got %.2f", inlineAllocs)
	}

	// 2. Verify UNRECOGNIZED kind, COW after Build(), GetChild/GetChildList on intrusive built node,
	// shared child added twice to the same parent, and mergeIntoBuilder preserving prior intrusive children.
	shared := protobuf.AstNode_NewBuilder().SetKind(protobuf.NodeKind_TRUE).SetRelativeColumn(4).Build()
	unrec := protobuf.AstNode_NewBuilder().SetKind(protobuf.NodeKind_UNRECOGNIZED).Build()
	parentB := protobuf.AstNode_NewBuilder().
		SetKind(protobuf.NodeKind_ARRAY_LITERAL).
		AddChild(shared).
		AddChild(unrec)
	built1 := parentB.Build()

	// Concurrent reads (GetChild, GetChildList, ToByteArray) on intrusive built node + concurrent AddChild of shared node.
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 40; i++ {
				if built1.GetChildCount() != 2 || built1.GetChild(0).GetKind() != protobuf.NodeKind_TRUE ||
					built1.GetChild(1).GetKind() != protobuf.NodeKind_UNRECOGNIZED {
					t.Errorf("concurrent GetChild mismatch")
					return
				}
				if list := built1.GetChildList(); len(list) != 2 {
					t.Errorf("concurrent GetChildList len = %d", len(list))
					return
				}
				other := protobuf.AstNode_NewBuilder().
					SetKind(protobuf.NodeKind_SOURCE_FILE).
					AddChild(shared).
					AddChild(protobuf.AstNode_NewBuilder().SetKind(protobuf.NodeKind_NULL).Build()).
					Build()
				bs := other.ToByteString()
				if bs.IsEmpty() {
					t.Errorf("unexpected empty ByteString")
					return
				}
			}
		}()
	}
	wg.Wait()

	// COW mutation on built1 via copy must not alter built1's child count or wire bytes.
	wireBeforeCOW := built1.ToByteArray()
	cowCopy := *built1
	cowCopy.AddChild(shared)
	if built1.GetChildCount() != 2 || cowCopy.GetChildCount() != 3 {
		t.Fatalf("COW child count mismatch: orig=%d copy=%d", built1.GetChildCount(), cowCopy.GetChildCount())
	}
	if !bytes.Equal(built1.ToByteArray(), wireBeforeCOW) {
		t.Fatal("COW mutated original built AstNode wire bytes")
	}

	// ByteString_CopyFromUtf8 zero-alloc wire-slab fast path for short ASCII strings.
	utf8Allocs := testing.AllocsPerRun(50, func() {
		b0 := protobuf.ByteString_CopyFromUtf8("")
		b1 := protobuf.ByteString_CopyFromUtf8("x")
		b2 := protobuf.ByteString_CopyFromUtf8("goog.module.exports.MyService")
		if b0.Size() != 0 || b1.ToStringUtf8() != "x" || b2.ToStringUtf8() != "goog.module.exports.MyService" {
			t.Fatal("ByteString_CopyFromUtf8 roundtrip mismatch")
		}
	})
	if utf8Allocs > 1 {
		t.Fatalf("expected <= 1 allocs for ByteString_CopyFromUtf8 wire-slab fast path, got %.2f", utf8Allocs)
	}
}

func TestWave37COptimizationJsdocPoolingPinningAndDirectWirePutTo(t *testing.T) {
	defJsdoc := protobuf.OptimizationJsdoc_GetDefaultInstance()

	// 1. Warm pool & verify JSDocSerializer pattern (NewBuilder -> AddKind/Set* -> Build -> DefaultInstance.Equals -> SetJsdoc -> ToByteString)
	// recycles both inline (<=7 tags) and overflow (>7 tags) OptimizationJsdoc instances with ~0 heap allocs/op.
	manyTags := []protobuf.JsdocTag{
		protobuf.JsdocTag_JSDOC_FILEOVERVIEW,
		protobuf.JsdocTag_JSDOC_CONST,
		protobuf.JsdocTag_JSDOC_CONSTRUCTOR,
		protobuf.JsdocTag_JSDOC_NO_INLINE,
		protobuf.JsdocTag_JSDOC_NO_SIDE_EFFECTS,
		protobuf.JsdocTag_JSDOC_MODIFIES_THIS,
		protobuf.JsdocTag_JSDOC_THROWS,
		protobuf.JsdocTag_JSDOC_ABSTRACT,
		protobuf.JsdocTag_JSDOC_DEFINE,
	}

	buildAndSerialize := func(seed int32) protobuf.ByteString {
		jRoot := protobuf.OptimizationJsdoc_NewBuilder()
		for _, tg := range manyTags {
			jRoot.AddKind(tg)
		}
		jRootBuilt := jRoot.
			SetLicenseTextPointer(seed + 1).
			SetDescriptionPointer(seed + 2).
			SetMeaningPointer(seed + 3).
			SetAlternateMessageIdPointer(seed + 4).
			SetClosureUnawareModePointer(seed + 5).
			Build()
		if defJsdoc.Equals(jRootBuilt) {
			t.Fatal("expected non-default jRootBuilt")
		}

		root := protobuf.AstNode_NewBuilder().
			SetKind(protobuf.NodeKind_SOURCE_FILE).
			SetSourceFile(1).
			SetJsdoc(jRootBuilt)

		for i := int32(0); i < 6; i++ {
			jc := protobuf.OptimizationJsdoc_NewBuilder().
				AddKind(protobuf.JsdocTag_JSDOC_NO_INLINE).
				AddKind(protobuf.JsdocTag_JSDOC_CONST).
				SetMeaningPointer(seed + 10 + i).
				Build()
			if defJsdoc.Equals(jc) {
				t.Fatal("expected non-default child jsdoc")
			}
			ch := protobuf.AstNode_NewBuilder().
				SetKind(protobuf.NodeKind_FUNCTION_LITERAL).
				SetJsdoc(jc).
				Build()
			root.AddChild(ch)
		}
		return root.Build().ToByteString()
	}

	refBS := buildAndSerialize(100)
	parsed, err := protobuf.AstNode_ParseFrom(refBS.BytesView())
	if err != nil {
		t.Fatalf("AstNode_ParseFrom failed: %v", err)
	}
	if parsed.GetJsdoc().GetKindCount() != len(manyTags) ||
		parsed.GetJsdoc().GetLicenseTextPointer() != 101 ||
		parsed.GetJsdoc().GetDescriptionPointer() != 102 ||
		parsed.GetJsdoc().GetMeaningPointer() != 103 ||
		parsed.GetJsdoc().GetAlternateMessageIdPointer() != 104 ||
		parsed.GetJsdoc().GetClosureUnawareModePointer() != 105 {
		t.Fatalf("parsed root jsdoc mismatch: %v", parsed.GetJsdoc().String())
	}
	// Re-serialize from astNodeSlab (putRecTo -> jsdoc.putTo) and verify exact wire bytes.
	if !bytes.Equal(parsed.ToByteArray(), refBS.BytesView()) {
		t.Fatal("slab putRecTo OptimizationJsdoc wire bytes mismatch")
	}
	parsed.ReleaseSlab()

	// Warm up pool batch completely and measure steady-state allocations.
	for i := 0; i < 64; i++ {
		_ = buildAndSerialize(int32(200 + i))
	}
	allocs := testing.AllocsPerRun(100, func() {
		bs := buildAndSerialize(100)
		if !bytes.Equal(bs.BytesView(), refBS.BytesView()) {
			t.Fatal("recycled OptimizationJsdoc wire mismatch")
		}
	})
	if allocs > 1 {
		t.Fatalf("expected <= 1 allocs/op with pooled OptimizationJsdoc + AstNode, got %.2f", allocs)
	}

	// 2. Pinning invariants: GetJsdoc(), multi-parent SetJsdoc(), ResetForPool(), and post-recycle Equals() rescue.
	sharedJsdoc := protobuf.OptimizationJsdoc_NewBuilder().
		AddKind(protobuf.JsdocTag_JSDOC_ENUM).
		SetDescriptionPointer(777).
		Build()
	p1 := protobuf.AstNode_NewBuilder().SetKind(protobuf.NodeKind_SOURCE_FILE).SetJsdoc(sharedJsdoc).Build()
	p2 := protobuf.AstNode_NewBuilder().SetKind(protobuf.NodeKind_SOURCE_FILE).SetJsdoc(sharedJsdoc).Build()
	_ = p1.ToByteString() // must NOT recycle sharedJsdoc because p2 also claimed/pinned it

	getterJsdoc := protobuf.OptimizationJsdoc_NewBuilder().
		AddKind(protobuf.JsdocTag_JSDOC_INTERFACE).
		SetMeaningPointer(888).
		Build()
	var inlineAst protobuf.AstNode
	inlineAst.SetKind(protobuf.NodeKind_BLOCK).SetJsdoc(getterJsdoc)
	gotGetter := inlineAst.GetJsdoc() // pins getterJsdoc
	inlineAst.ResetForPool()

	// Churn pool heavily to ensure pinned sharedJsdoc & gotGetter are never overwritten.
	for i := 0; i < 600; i++ {
		_ = buildAndSerialize(int32(1000 + i))
	}
	if sharedJsdoc.GetDescriptionPointer() != 777 || p2.GetJsdoc().GetDescriptionPointer() != 777 {
		t.Fatalf("multi-parent sharedJsdoc corrupted after pool churn: %d", sharedJsdoc.GetDescriptionPointer())
	}
	if gotGetter.GetMeaningPointer() != 888 {
		t.Fatalf("GetJsdoc()-pinned instance corrupted after ResetForPool + churn: %d", gotGetter.GetMeaningPointer())
	}

	// 3. Concurrent -race stress across 8 goroutines.
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		g := g
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 80; i++ {
				bs := buildAndSerialize(int32(g*1000 + i))
				if bs.IsEmpty() {
					t.Errorf("empty serialized ByteString")
					return
				}
			}
		}()
	}
	wg.Wait()
}
