package jscomp

import (
	"bufio"
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
	Equals(any) bool
	HashCode() int32
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
		parse := func() protoTestMessage {
			switch cols[0] {
			case "AstNode":
				return AstNode_ParseFrom(raw)
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
		if !m.Equals(m2) || m.HashCode() != m2.HashCode() {
			t.Errorf("%s %s: not equal to itself", cols[0], cols[1])
		}
		n++
	}
	if n < 300 {
		t.Fatalf("only %d vectors", n)
	}
}

func TestProtoBuilders(t *testing.T) {
	b := TypePool_NewBuilder()
	b.AddTypeBuilder().GetObjectBuilder().SetUuid([]byte{1, 2}).SetIsInvalidating(true)
	u := b.AddTypeBuilder().GetUnionBuilder()
	u.AddUnionMember(3).AddUnionMember(4)
	b.AddDisambiguationEdgesBuilder().SetSubtype(8).SetSupertype(9)
	first := b.Build()
	u.AddUnionMember(5) // must not affect the built message
	second := b.Build()
	if first.GetType(1).GetUnion().GetUnionMemberCount() != 2 || second.GetType(1).GetUnion().GetUnionMemberCount() != 3 {
		t.Errorf("builder isolation: %v / %v", first, second)
	}
	want := "type {\n  object {\n    is_invalidating: true\n    uuid: \"\\001\\002\"\n  }\n}\ntype {\n  union {\n    union_member: 3\n    union_member: 4\n  }\n}\ndisambiguation_edges {\n  subtype: 8\n  supertype: 9\n}\n"
	if first.String() != want {
		t.Errorf("got %q", first.String())
	}
	again := TypePool_ParseFrom(first.ToByteArray())
	if !again.Equals(first) || again.Equals(second) {
		t.Errorf("equality after round trip")
	}
	if !OptimizationJsdoc_GetDefaultInstance().Equals(OptimizationJsdoc_NewBuilder().Build()) {
		t.Errorf("default instance equality")
	}
	if TypeProto_NewBuilder().Build().GetKindCase() != TypeProto_KindCase_KIND_NOT_SET {
		t.Errorf("oneof case")
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
		AstNode_ParseFrom(data)
		return ""
	}
	for _, depth := range []int{99, 100, 101} {
		n := AstNode_NewBuilder().SetKind(NodeKind_NUMBER_LITERAL).Build()
		for i := 0; i < depth; i++ {
			n = AstNode_NewBuilder().SetKind(NodeKind_BLOCK).AddChild(n).Build()
		}
		want := ""
		if depth > 100 {
			want = msg
		}
		if got := parse(n.ToByteArray()); got != want {
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
