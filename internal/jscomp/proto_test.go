package jscomp

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"os"
	"strings"
	"testing"
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
