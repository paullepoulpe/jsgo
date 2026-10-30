package jscomp

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// testdata/proto_vectors.tsv is written by tools/oracle/java/ProtoVectors.java:
// random messages serialized by the Java protobuf runtime (with unknown
// fields), the re-serialization of the parsed message, and its toString().

type protoTestMessage interface {
	String() string
	ToByteArray() []byte
}

func TestProtoRoundTripMatchesJava(t *testing.T) {
	f, err := os.Open("testdata/proto_vectors.tsv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1<<20), 1<<22)
	n := 0
	for sc.Scan() {
		cols := strings.Split(sc.Text(), "\t")
		raw, _ := hex.DecodeString(cols[1])
		reser := cols[2]
		text := strings.NewReplacer(`\\`, `\`, `\n`, "\n", `\t`, "\t").Replace(cols[3])
		if cols[0] == "AstNode" {
			m := parseTestAstNode(raw)
			if got := m.String(); got != text {
				t.Errorf("AstNode %s: toString\n got %q\nwant %q", cols[1], got, text)
			}
			reserBytes, _ := hex.DecodeString(reser)
			if got := parseTestAstNode(reserBytes).String(); got != text {
				t.Errorf("AstNode %s: reser toString\n got %q\nwant %q", cols[1], got, text)
			}
			n++
			continue
		}
		parse := func() protoTestMessage {
			switch cols[0] {
			case "TypePool":
				return TypePool_ParseFrom(raw)
			case "SourceFileProto":
				return SourceFileProto_ParseFrom(raw)
			case "Requirement":
				return Requirement_ParseFrom(raw)
			case "OptimizationJsdoc":
				return OptimizationJsdoc_ParseFrom(raw)
			case "ConformanceConfig":
				return ConformanceConfig_ParseFrom(raw)
			}
			t.Fatalf("unknown type %s", cols[0])
			return nil
		}
		m := parse()
		if got := m.String(); got != text {
			t.Errorf("%s %s: toString\n got %q\nwant %q", cols[0], cols[1], got, text)
		}
		if got := hex.EncodeToString(m.ToByteArray()); got != reser {
			t.Errorf("%s %s: bytes\n got %s\nwant %s", cols[0], cols[1], got, reser)
		}
		m2 := parse()
		if !bytes.Equal(m.ToByteArray(), m2.ToByteArray()) || m.String() != m2.String() {
			t.Errorf("%s %s: round-trip mismatch", cols[0], cols[1])
		}
		n++
	}
	if n < 300 {
		t.Fatalf("only %d vectors", n)
	}
}

func TestProtoBuilders(t *testing.T) {
	first := &TypePool{
		type_: []*TypeProto{
			{
				kindCase_: TypeProto_OBJECT_FIELD_NUMBER,
				object_:   &ObjectTypeProto{uuid_: []byte{1, 2}, isInvalidating_: true},
			},
			{
				kindCase_: TypeProto_UNION_FIELD_NUMBER,
				union_:    &UnionTypeProto{unionMember_: []int32{3, 4}},
			},
		},
		disambiguationEdges_: []*SubtypingEdge{
			{subtype_: 8, supertype_: 9},
		},
	}
	want := "type {\n  object {\n    is_invalidating: true\n    uuid: \"\\001\\002\"\n  }\n}\ntype {\n  union {\n    union_member: 3\n    union_member: 4\n  }\n}\ndisambiguation_edges {\n  subtype: 8\n  supertype: 9\n}\n"
	if first.String() != want {
		t.Errorf("got %q", first.String())
	}
	again := TypePool_ParseFrom(first.ToByteArray())
	if again.String() != want || !bytes.Equal(again.ToByteArray(), first.ToByteArray()) {
		t.Errorf("round-trip mismatch: %v", again)
	}
	if len(OptimizationJsdoc_GetDefaultInstance().ToByteArray()) != 0 {
		t.Errorf("default instance non-empty")
	}
	cb := ConformanceConfig_NewBuilder()
	cb.AddRequirementBuilder().SetErrorMessage("first")
	cb.AddRequirementBuilder().SetErrorMessage("second")
	cfg := cb.Build()
	if len(cfg.GetRequirementList()) != 2 {
		t.Errorf("builder count: %v", cfg)
	}
	sf := SourceFileProto_NewBuilder().SetFilename("a.js").Build()
	if !sf.Equals(SourceFileProto_ParseFrom(sf.ToByteArray())) {
		t.Errorf("SourceFileProto equality after round trip")
	}
	r := Requirement_NewBuilder().Build()
	if r.GetTypeMatchingStrategy() != Requirement_TypeMatchingStrategy_LOOSE || !r.GetReportLooseTypeViolations() ||
		r.GetSeverity() != Requirement_Severity_WARNING || r.HasSeverity() {
		t.Errorf("proto2 defaults")
	}
}

// TestProtoRecursionLimit checks CodedInputStream's default recursion limit
// of 100, as observed with the jar's protobuf runtime: 100 nested messages
// (or unknown groups) parse, 101 throw InvalidProtocolBufferException.
func TestProtoRecursionLimit(t *testing.T) {
	const msg = "Protocol message had too many levels of nesting.  May be malicious.  Use setRecursionLimit() to increase the recursion depth limit."
	parse := func(data []byte) (err string) {
		defer func() {
			if r := recover(); r != nil {
				err = r.(*InvalidProtocolBufferException).GetMessage()
			}
		}()
		parseTestAstNode(data)
		return ""
	}
	for _, depth := range []int{99, 100, 101} {
		msgBytes := []byte{8, 1}
		for i := 0; i < depth; i++ {
			wrapped := protoAppendVarint([]byte{8, 0xa0, 0x03, 18}, uint64(len(msgBytes)))
			msgBytes = append(wrapped, msgBytes...)
		}
		want := ""
		if depth > 100 {
			want = msg
		}
		if got := parse(msgBytes); got != want {
			t.Errorf("message depth %d: got %q", depth, got)
		}
		var b []byte
		for i := 0; i < depth; i++ {
			b = append(b, (1000<<3|3)&0x7f|0x80, (1000<<3|3)>>7)
		}
		b = append(b, 8, 1)
		for i := 0; i < depth; i++ {
			b = append(b, (1000<<3|4)&0x7f|0x80, (1000<<3|4)>>7)
		}
		if got := parse(b); got != want {
			t.Errorf("group depth %d: got %q", depth, got)
		}
	}
}

type testAstNode struct {
	kind_                int32
	child_               []*testAstNode
	stringValuePointer_  int32
	doubleValue_         float64
	templateStringValue_ *testTemplateStringValue
	relativeLine_        int32
	relativeColumn_      int32
	jsdoc_               *OptimizationJsdoc
	originalNamePointer_ int32
	booleanProperties_   int64
	type_                int32
	sourceFile_          int32
	bitField0_           uint64
	valueCase_           uint8
	unknownFields        []byte
}

func parseTestAstNode(data []byte) *testAstNode {
	m := &testAstNode{}
	r := protoReader{b: data}
	m.mergeFromReader(&r)
	return m
}

func nodeKind_text(value int32) string {
	if e, ok := NodeKind_ForNumber(value); ok {
		return e.String()
	}
	return protoTextInt32(value)
}

func (m *testAstNode) mergeFromReader(r *protoReader) {
	for {
		start := r.pos
		tag := r.readTag()
		switch tag {
		case 0:
			return
		case 8:
			m.kind_ = r.readInt32()
		case 18:
			data := r.readMessageBytes()
			x := &testAstNode{}
			sub := r.sub(data)
			x.mergeFromReader(&sub)
			m.child_ = append(m.child_, x)
		case 24:
			m.stringValuePointer_ = r.readUint32()
			m.valueCase_ = 3
		case 33:
			m.doubleValue_ = r.readDouble()
			m.valueCase_ = 4
		case 40:
			m.relativeLine_ = r.readSint32()
		case 48:
			m.relativeColumn_ = r.readSint32()
		case 58:
			data := r.readMessageBytes()
			var x *OptimizationJsdoc
			if m.jsdoc_ != nil {
				x = m.jsdoc_
				x.cachedSize = 0
			} else {
				x = &OptimizationJsdoc{}
			}
			sub := r.sub(data)
			x.mergeFromReader(&sub)
			m.jsdoc_ = x
		case 66:
			data := r.readMessageBytes()
			var x *testTemplateStringValue
			if m.valueCase_ == 8 && m.templateStringValue_ != nil {
				x = m.templateStringValue_
			} else {
				x = &testTemplateStringValue{}
			}
			sub := r.sub(data)
			x.mergeFromReader(&sub)
			m.templateStringValue_ = x
			m.valueCase_ = 8
		case 72:
			m.originalNamePointer_ = r.readUint32()
		case 112:
			m.booleanProperties_ = r.readInt64()
		case 120:
			m.type_ = r.readInt32()
			m.bitField0_ |= 1 << 0
		case 128:
			m.sourceFile_ = r.readUint32()
		default:
			r.skipField(tag)
			m.unknownFields = append(m.unknownFields, r.b[start:r.pos]...)
		}
	}
}

func (m *testAstNode) printText(p *protoTextPrinter) {
	if m.kind_ != 0 {
		p.scalar("kind", nodeKind_text(m.kind_))
	}
	for _, x := range m.child_ {
		p.message("child", x.printText)
	}
	if m.valueCase_ == 3 {
		p.scalar("string_value_pointer", protoTextUint32(m.stringValuePointer_))
	}
	if m.valueCase_ == 4 {
		p.scalar("double_value", protoTextDouble(m.doubleValue_))
	}
	if m.relativeLine_ != 0 {
		p.scalar("relative_line", protoTextInt32(m.relativeLine_))
	}
	if m.relativeColumn_ != 0 {
		p.scalar("relative_column", protoTextInt32(m.relativeColumn_))
	}
	if m.jsdoc_ != nil {
		p.message("jsdoc", m.jsdoc_.printText)
	}
	if m.valueCase_ == 8 {
		p.message("template_string_value", m.templateStringValue_.printText)
	}
	if m.originalNamePointer_ != 0 {
		p.scalar("original_name_pointer", protoTextUint32(m.originalNamePointer_))
	}
	if m.booleanProperties_ != 0 {
		p.scalar("boolean_properties", protoTextInt64(m.booleanProperties_))
	}
	if m.bitField0_&(1<<0) != 0 {
		p.scalar("type", protoTextInt32(m.type_))
	}
	if m.sourceFile_ != 0 {
		p.scalar("source_file", protoTextUint32(m.sourceFile_))
	}
	p.unknownFields(m.unknownFields)
}

func (m *testAstNode) String() string {
	p := &protoTextPrinter{}
	m.printText(p)
	return p.sb.String()
}

type testTemplateStringValue struct {
	rawStringPointer_    int32
	cookedStringPointer_ int32
	unknownFields        []byte
}

func (m *testTemplateStringValue) mergeFromReader(r *protoReader) {
	for {
		start := r.pos
		tag := r.readTag()
		switch tag {
		case 0:
			return
		case 8:
			m.rawStringPointer_ = r.readUint32()
		case 16:
			m.cookedStringPointer_ = r.readInt32()
		default:
			r.skipField(tag)
			m.unknownFields = append(m.unknownFields, r.b[start:r.pos]...)
		}
	}
}

func (m *testTemplateStringValue) printText(p *protoTextPrinter) {
	if m.rawStringPointer_ != 0 {
		p.scalar("raw_string_pointer", protoTextUint32(m.rawStringPointer_))
	}
	if m.cookedStringPointer_ != 0 {
		p.scalar("cooked_string_pointer", protoTextInt32(m.cookedStringPointer_))
	}
	p.unknownFields(m.unknownFields)
}

func (m *testTemplateStringValue) String() string {
	p := &protoTextPrinter{}
	m.printText(p)
	return p.sb.String()
}

func TestStreamLargeShardMatchesDecodeShardBytes(t *testing.T) {
	shard := &TypedAst{
		typePool_: &TypePool{
			type_: []*TypeProto{
				{
					kindCase_: TypeProto_OBJECT_FIELD_NUMBER,
					object_:   &ObjectTypeProto{uuid_: []byte{9, 8, 7}, isInvalidating_: true, prototype_: []int32{1}},
				},
				{
					kindCase_: TypeProto_UNION_FIELD_NUMBER,
					union_:    &UnionTypeProto{unionMember_: []int32{0, 2}},
				},
			},
			disambiguationEdges_: []*SubtypingEdge{
				{subtype_: 1, supertype_: 2},
			},
			debugInfo_: &TypePool_DebugInfo{
				mismatch_: []*TypePool_DebugInfo_Mismatch{
					{sourceRef_: "a.js:1", involvedColor_: []int32{0, 1}},
				},
			},
		},
		stringPool_: &StringPoolProto{
			strings_:   [][]byte{[]byte("alpha"), []byte("caf\xc3\xa9")},
			maxLength_: 16,
		},
		externAst_: []*LazyAst{
			{script_: []byte{8, 1}, sourceFile_: 1},
		},
		codeAst_: []*LazyAst{
			{script_: []byte{8, 1, 40, 2}, sourceFile_: 2, sourceMappingUrl_: "a.js.map"},
		},
		sourceFilePool_: &SourceFilePool{
			sourceFile_: []*SourceFileProto{
				{filename_: "externs.js", sourceKind_: SourceFileProto_SourceKind_EXTERN.GetNumber()},
				{filename_: "a.js", loaderCase_: 2, preloadedContents_: "var x = 1;"},
			},
		},
		externsSummary_: &ExternsSummary{
			propNamePtr_: []int32{1, 2},
		},
		runtimeLibraryToInject_: []string{"es6/symbol"},
	}

	raw := shard.ToByteArray()
	wantAst, wantPool := typedAstDeserializer_decodeShardBytes(raw, nil)
	gotAst, gotPool := typedAstDeserializer_streamLargeShard(bytes.NewReader(raw), len(raw), nil)

	wantAstNoStr := *wantAst
	wantAstNoStr.cachedSize = 0
	wantAstNoStr.stringPool_ = nil
	gotAstNoStr := *gotAst
	gotAstNoStr.cachedSize = 0
	gotAstNoStr.stringPool_ = nil
	if gotAstNoStr.String() != wantAstNoStr.String() {
		t.Fatalf("TypedAst mismatch:\ngot:\n%s\nwant:\n%s", gotAstNoStr.String(), wantAstNoStr.String())
	}
	for i := 0; i <= 2; i++ {
		if gotPool.Get(i) != wantPool.Get(i) {
			t.Fatalf("StringPool[%d] = %q, want %q", i, gotPool.Get(i), wantPool.Get(i))
		}
	}

	streamBytes := (&TypedAst_List{typedAsts_: []*TypedAst{shard, shard}}).ToByteArray()
	var bufShards, readerShards int
	typedAstDeserializer_scanTypedAstStream(
		bytes.NewBuffer(bytes.Clone(streamBytes)),
		func(batch [][]byte) { bufShards += len(batch) },
		func(io.Reader, int) { t.Fatal("unexpected large shard on small test") },
	)
	typedAstDeserializer_scanTypedAstStream(
		bytes.NewReader(streamBytes),
		func(batch [][]byte) { readerShards += len(batch) },
		func(io.Reader, int) { t.Fatal("unexpected large shard on small test") },
	)
	if bufShards != 2 || readerShards != 2 {
		t.Fatalf("shard counts: buf=%d reader=%d, want 2", bufShards, readerShards)
	}
}

func TypePool_ParseFrom(data []byte) *TypePool {
	m := &TypePool{}
	r := protoReader{b: data}
	m.mergeFromReader(&r)
	return m
}

func SourceFileProto_ParseFrom(data []byte) *SourceFileProto {
	m := &SourceFileProto{}
	r := protoReader{b: data}
	m.mergeFromReader(&r)
	return m
}

func ConformanceConfig_ParseFrom(data []byte) *ConformanceConfig {
	m := &ConformanceConfig{}
	r := protoReader{b: data}
	m.mergeFromReader(&r)
	return m
}

func Requirement_ParseFrom(data []byte) *Requirement {
	m := &Requirement{}
	r := protoReader{b: data}
	m.mergeFromReader(&r)
	return m
}

func OptimizationJsdoc_ParseFrom(data []byte) *OptimizationJsdoc {
	m := &OptimizationJsdoc{}
	r := protoReader{b: data}
	m.mergeFromReader(&r)
	return m
}

func NodeKind_ForNumber(value int32) (NodeKind, bool) {
	if uint32(value) < uint32(len(nodeKind_byNumber)) {
		if e := nodeKind_byNumber[value]; e != NodeKind_UNRECOGNIZED {
			return e, true
		}
	}
	return 0, false
}

func (m *Requirement) ClearWhitelistEntry() *Requirement {
	m.whitelistEntry_ = nil
	return m
}

func (m *Requirement) ClearAllowlistEntry() *Requirement {
	m.allowlistEntry_ = nil
	return m
}

func SourceFileProto_FileOnDisk_GetDefaultInstance() *SourceFileProto_FileOnDisk {
	return sourceFileProto_FileOnDisk_DEFAULT_INSTANCE
}

func SourceFileProto_NewBuilder() *SourceFileProto {
	return &SourceFileProto{}
}

func (m *SourceFileProto) Build() *SourceFileProto {
	m.cachedSize = 0
	return m
}

func (m *SourceFileProto) SetFilename(v string) *SourceFileProto {
	m.filename_ = v
	return m
}

func (m *SourceFileProto) SetPreloadedContents(v string) *SourceFileProto {
	m.preloadedContents_ = v
	m.fileOnDisk_ = nil
	m.zipEntry_ = nil
	m.stubFile_ = false
	m.loaderCase_ = 2
	return m
}

func (m *SourceFileProto) SetFileOnDisk(v *SourceFileProto_FileOnDisk) *SourceFileProto {
	if v == nil {
		panic(jutil.NewNullPointerException())
	}
	m.fileOnDisk_ = v
	m.preloadedContents_ = ""
	m.zipEntry_ = nil
	m.stubFile_ = false
	m.loaderCase_ = 3
	return m
}

func (m *SourceFileProto) SetSourceKind(v SourceFileProto_SourceKind) *SourceFileProto {
	m.sourceKind_ = v.GetNumber()
	return m
}

func (m *SourceFileProto) SetIsClosureUnawareCode(v bool) *SourceFileProto {
	m.isClosureUnawareCode_ = v
	return m
}

func (m *SourceFileProto) SetNumBytesPlusOne(v int32) *SourceFileProto {
	m.numBytesPlusOne_ = v
	return m
}

func (m *SourceFileProto) SetNumLinesPlusOne(v int32) *SourceFileProto {
	m.numLinesPlusOne_ = v
	return m
}

func ConformanceViolationsProto_NewBuilder() *ConformanceViolationsProto {
	return &ConformanceViolationsProto{}
}

func (m *ConformanceViolationsProto) GetViolations(index int) *ViolationProto {
	return m.violations_[index]
}

func (m *ConformanceViolationsProto) GetViolationsCount() int { return len(m.violations_) }

func (m *ConformanceViolationsProto) Build() *ConformanceViolationsProto {
	m.cachedSize = 0
	return m
}

func (m *ConformanceViolationsProto) AddViolations(v *ViolationProto) *ConformanceViolationsProto {
	if v == nil {
		panic(jutil.NewNullPointerException())
	}
	m.violations_ = append(m.violations_, v)
	return m
}

func ViolationProto_NewBuilder() *ViolationProto { return &ViolationProto{} }

func (m *ViolationProto) GetSourceName() string { return m.sourceName_ }

func (m *ViolationProto) GetLineNumber() int32 { return m.lineNumber_ }

func (m *ViolationProto) Build() *ViolationProto {
	m.cachedSize = 0
	return m
}

func (m *ViolationProto) SetSourceName(v string) *ViolationProto {
	m.sourceName_ = v
	m.bitField0_ |= 1 << 0
	return m
}

func (m *ViolationProto) SetLineNumber(v int32) *ViolationProto {
	m.lineNumber_ = v
	return m
}

func (m *ViolationProto) SetCharNumber(v int32) *ViolationProto {
	m.charNumber_ = v
	return m
}

func (m *SourceFilePool) appendTo(b []byte) []byte {
	for _, x := range m.sourceFile_ {
		b = protoAppendTag(b, 1, 2)
		b = protoAppendVarint(b, uint64(x.serializedSize()))
		b = x.appendTo(b)
	}
	return b
}

func (m *SourceFilePool) printText(p *protoTextPrinter) {
	for _, x := range m.sourceFile_ {
		p.message("source_file", x.printText)
	}
}

func (m *ExternsSummary) printText(p *protoTextPrinter) {
	for _, x := range m.propNamePtr_ {
		p.scalar("prop_name_ptr", protoTextInt32(x))
	}
}

func (m *LazyAst) printText(p *protoTextPrinter) {
	if len(m.script_) != 0 {
		p.scalar("script", protoTextBytes(m.script_))
	}
	if m.sourceFile_ != 0 {
		p.scalar("source_file", protoTextUint32(m.sourceFile_))
	}
	if m.sourceMappingUrl_ != "" {
		p.scalar("source_mapping_url", protoTextString(m.sourceMappingUrl_))
	}
}

func (m *LazyAst) String() string {
	p := &protoTextPrinter{}
	m.printText(p)
	return p.sb.String()
}

func StringPoolProto_GetDefaultInstance() *StringPoolProto { return stringPoolProto_DEFAULT_INSTANCE }

func StringPoolProto_ParseFrom(data []byte) *StringPoolProto {
	m := &StringPoolProto{}
	r := protoReader{b: data}
	m.mergeFromReader(&r)
	return m
}

func (m *StringPoolProto) appendTo(b []byte) []byte {
	for _, x := range m.strings_ {
		b = protoAppendTag(b, 1, 2)
		b = protoAppendBytes(b, x)
	}
	if m.maxLength_ != 0 {
		b = protoAppendTag(b, 2, 0)
		b = protoAppendInt32(b, m.maxLength_)
	}
	return b
}

func (m *StringPoolProto) ToByteArray() []byte {
	return m.appendTo(make([]byte, 0, m.serializedSize()))
}

func (m *StringPoolProto) printText(p *protoTextPrinter) {
	for _, x := range m.strings_ {
		p.scalar("strings", protoTextBytes(x))
	}
	if m.maxLength_ != 0 {
		p.scalar("max_length", protoTextInt32(m.maxLength_))
	}
}

func (m *StringPoolProto) String() string {
	p := &protoTextPrinter{}
	m.printText(p)
	return p.sb.String()
}

func (m *TypedAst) appendTo(b []byte) []byte {
	if m.typePool_ != nil {
		b = protoAppendTag(b, 1, 2)
		b = protoAppendVarint(b, uint64(m.typePool_.serializedSize()))
		b = m.typePool_.appendTo(b)
	}
	if m.stringPool_ != nil {
		b = protoAppendTag(b, 2, 2)
		b = protoAppendVarint(b, uint64(m.stringPool_.serializedSize()))
		b = m.stringPool_.appendTo(b)
	}
	for _, x := range m.externAst_ {
		b = protoAppendTag(b, 3, 2)
		b = protoAppendVarint(b, uint64(x.serializedSize()))
		b = x.appendTo(b)
	}
	for _, x := range m.codeAst_ {
		b = protoAppendTag(b, 4, 2)
		b = protoAppendVarint(b, uint64(x.serializedSize()))
		b = x.appendTo(b)
	}
	if m.sourceFilePool_ != nil {
		b = protoAppendTag(b, 5, 2)
		b = protoAppendVarint(b, uint64(m.sourceFilePool_.serializedSize()))
		b = m.sourceFilePool_.appendTo(b)
	}
	if m.externsSummary_ != nil {
		b = protoAppendTag(b, 6, 2)
		b = protoAppendVarint(b, uint64(m.externsSummary_.serializedSize()))
		b = m.externsSummary_.appendTo(b)
	}
	for _, x := range m.runtimeLibraryToInject_ {
		b = protoAppendTag(b, 7, 2)
		b = protoAppendString(b, x)
	}
	return b
}

func (m *TypedAst) ToByteArray() []byte { return m.appendTo(make([]byte, 0, m.serializedSize())) }

func (m *TypedAst) printText(p *protoTextPrinter) {
	if m.typePool_ != nil {
		p.message("type_pool", m.typePool_.printText)
	}
	if m.stringPool_ != nil {
		p.message("string_pool", m.stringPool_.printText)
	}
	for _, x := range m.externAst_ {
		p.message("extern_ast", x.printText)
	}
	for _, x := range m.codeAst_ {
		p.message("code_ast", x.printText)
	}
	if m.sourceFilePool_ != nil {
		p.message("source_file_pool", m.sourceFilePool_.printText)
	}
	if m.externsSummary_ != nil {
		p.message("externs_summary", m.externsSummary_.printText)
	}
	for _, x := range m.runtimeLibraryToInject_ {
		p.scalar("runtime_library_to_inject", protoTextString(x))
	}
}

func (m *TypedAst) String() string {
	p := &protoTextPrinter{}
	m.printText(p)
	return p.sb.String()
}

type TypedAst_List struct {
	typedAsts_ []*TypedAst
	cachedSize int
}

func TypedAst_List_ParseFrom(data []byte) *TypedAst_List {
	m := &TypedAst_List{}
	r := protoReader{b: data}
	m.mergeFromReader(&r)
	return m
}

func (m *TypedAst_List) GetTypedAsts(index int) *TypedAst { return m.typedAsts_[index] }

func (m *TypedAst_List) appendTo(b []byte) []byte {
	for _, x := range m.typedAsts_ {
		b = protoAppendTag(b, 1, 2)
		b = protoAppendVarint(b, uint64(x.serializedSize()))
		b = x.appendTo(b)
	}
	return b
}

func (m *TypedAst_List) serializedSize() int {
	if m.cachedSize > 0 {
		return m.cachedSize - 1
	}
	size := 0
	for _, x := range m.typedAsts_ {
		s := x.serializedSize()
		size += 1 + protoVarintLen(uint64(s)) + s
	}
	m.cachedSize = size + 1
	return size
}

func (m *TypedAst_List) ToByteArray() []byte { return m.appendTo(make([]byte, 0, m.serializedSize())) }

func (m *TypedAst_List) mergeFromReader(r *protoReader) {
	for {
		tag := r.readTag()
		switch tag {
		case 0:
			return
		case 10:
			data := r.readMessageBytes()
			x := &TypedAst{}
			sub := r.sub(data)
			x.mergeFromReader(&sub)
			m.typedAsts_ = append(m.typedAsts_, x)
		default:
			r.skipField(tag)
		}
	}
}
