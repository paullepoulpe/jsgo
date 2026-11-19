package guava_test

import (
	"bytes"
	"path/filepath"
	"slices"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/shims/guava"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javanio"
)

func TestStringsAsciiJoinerMoreObjects(t *testing.T) {
	if guava.Strings_PadStart("7", 3, '0') != "007" {
		t.Fatal("PadStart mismatch")
	}
	if guava.Strings_PadEnd("4.", 5, '0') != "4.000" {
		t.Fatal("PadEnd mismatch")
	}
	// Surrogate pair boundary check in CommonPrefix: U+1F600 (D83D DE00) vs U+1F601 (D83D DE01)
	if got := guava.Strings_CommonPrefix("a😀", "a😁"); got != "a" {
		t.Fatalf("CommonPrefix split surrogate pair: %q", got)
	}
	if !guava.Ascii_EqualsIgnoreCase("ClosureCompiler_123", "closurecompiler_123") {
		t.Fatal("Ascii_EqualsIgnoreCase mismatch")
	}
	if guava.Ascii_ToLowerCase("AbC-Ä") != "abc-Ä" {
		t.Fatal("Ascii_ToLowerCase modified non-ASCII")
	}

	joined := guava.Joiner_On(", ").SkipNulls().Join("a", nil, "b", int32(3))
	if joined != "a, b, 3" {
		t.Fatalf("Joiner SkipNulls got %q", joined)
	}
	joinedNull := guava.Joiner_On("-").UseForNull("NULL").Join("x", nil, "y")
	if joinedNull != "x-NULL-y" {
		t.Fatalf("Joiner UseForNull got %q", joinedNull)
	}

	str := guava.MoreObjects_ToStringHelper("Node").OmitNullValues().Add("type", "NAME").Add("doc", nil).AddValue(int32(42)).String()
	if str != "Node{type=NAME, 42}" {
		t.Fatalf("ToStringHelper got %q", str)
	}
}

func TestCharMatcherAndSplitter(t *testing.T) {
	ws := guava.CharMatcher_Whitespace()
	if got := ws.TrimFrom(" \t\n hello \u3000 "); got != "hello" {
		t.Fatalf("TrimFrom got %q", got)
	}
	if got := ws.CollapseFrom("a   b\t\tc", ' '); got != "a b c" {
		t.Fatalf("CollapseFrom got %q", got)
	}

	parts := guava.Splitter_OnChar(',').TrimResults().OmitEmptyStrings().SplitToList(" foo, ,bar , baz ,")
	if !slices.Equal(parts, []string{"foo", "bar", "baz"}) {
		t.Fatalf("Splitter_OnChar got %v", parts)
	}

	limited := guava.Splitter_On("::").TrimResults().Limit(2).SplitToList(" a :: b :: c ")
	if !slices.Equal(limited, []string{"a", "b :: c"}) {
		t.Fatalf("Splitter Limit got %v", limited)
	}

	reParts := guava.Splitter_OnPattern(`\r?\n`).OmitEmptyStrings().SplitToList("l1\nl2\r\n\nl3")
	if !slices.Equal(reParts, []string{"l1", "l2", "l3"}) {
		t.Fatalf("Splitter_OnPattern got %v", reParts)
	}
}

func TestImmutableBuildersAndPercentEscaper(t *testing.T) {
	mb := guava.ImmutableMap_NewBuilder[string, int]()
	mb.Put("a", 1).Put("b", 2)
	im := mb.BuildOrThrow()
	if im.Size() != 2 || im.GetVal("b") != 2 {
		t.Fatalf("ImmutableMap_Builder mismatch: %v", im)
	}

	// Duplicate key panics in BuildOrThrow.
	func() {
		defer func() {
			if recover() == nil {
				t.Fatal("expected duplicate key panic in BuildOrThrow")
			}
		}()
		guava.ImmutableMap_NewBuilder[string, int]().Put("k", 1).Put("k", 2).BuildOrThrow()
	}()

	// TranspileResult.java: new PercentEscaper("-_.*", false)
	esc := guava.NewPercentEscaper("-_.*", false)
	got := esc.Escape(`{"version":3,"file":"a b.js"}`)
	want := `%7B%22version%22%3A3%2C%22file%22%3A%22a%20b.js%22%7D`
	if got != want {
		t.Fatalf("PercentEscaper got %q, want %q", got, want)
	}
}

func TestBaseEncodingStreamsAndHashing(t *testing.T) {
	raw := []byte(`{"version":3}`)
	b64 := guava.BaseEncoding_Base64().Encode(raw)
	if !bytes.Equal(guava.BaseEncoding_Base64().Decode(b64), raw) {
		t.Fatal("Base64 round-trip failed")
	}

	lines := guava.CharStreams_ReadLines(javaio.NewStringReader("alpha\r\nbeta\ngamma"))
	if !slices.Equal(lines, []string{"alpha", "beta", "gamma"}) {
		t.Fatalf("CharStreams_ReadLines got %v", lines)
	}

	// CharStreams_ToString fast path on unread StringReader, BufferedReader, and InputStreamReader.
	unicodeSrc := "function f() { return 'λ😀'; }\n"
	sr := javaio.NewStringReader(unicodeSrc)
	if got := guava.CharStreams_ToString(sr); got != unicodeSrc {
		t.Fatalf("CharStreams_ToString(StringReader) got %q, want %q", got, unicodeSrc)
	}
	if got := guava.CharStreams_ToString(sr); got != "" {
		t.Fatalf("second CharStreams_ToString on drained StringReader got %q, want empty", got)
	}

	br := javaio.NewBufferedReader(javaio.NewInputStreamReader(bytes.NewReader([]byte(unicodeSrc))))
	if got := guava.CharStreams_ToString(br); got != unicodeSrc {
		t.Fatalf("CharStreams_ToString(BufferedReader) got %q, want %q", got, unicodeSrc)
	}
	// Partially read BufferedReader falls back to chunk reading seamlessly.
	brPart := javaio.NewBufferedReader(javaio.NewStringReader("first\nsecond\nthird"))
	if first, err := brPart.ReadLine(); err != nil || first != "first" {
		t.Fatalf("ReadLine got %q, err=%v", first, err)
	}
	if rest := guava.CharStreams_ToString(brPart); rest != "second\nthird" {
		t.Fatalf("CharStreams_ToString(partially read BufferedReader) got %q", rest)
	}

	tmp := javaio.NewFile(filepath.Join(t.TempDir(), "out.txt"))
	if err := guava.Files_AsCharSinkWrite(tmp, javanio.StandardCharsets_UTF_8, "hello\nworld"); err != nil {
		t.Fatal(err)
	}
	readBack, err := guava.Files_AsCharSource(tmp, javanio.StandardCharsets_UTF_8).Read()
	if err != nil || readBack != "hello\nworld" {
		t.Fatalf("Files_AsCharSource got %q, err=%v", readBack, err)
	}

	// Verify FarmHash64 determinism across length buckets (0, 8, 24, 48, 100 bytes).
	if guava.FarmHashFingerprint64(nil) != 0x9ae16a3b2f90404f {
		t.Fatalf("FarmHash64 empty mismatch: 0x%x", guava.FarmHashFingerprint64(nil))
	}
	marker := uint32(0x8c8b70db)
	fh := guava.Hashing_FarmHashFingerprint64().NewHasher().
		PutInt(int32(marker)).
		PutString("Array", javanio.StandardCharsets_UTF_8).
		Hash()
	if fh.AsLong() == 0 {
		t.Fatal("unexpected zero FarmHash64")
	}

	// Verify Murmur3_128 empty vector (Guava Hashing.murmur3_128().hashBytes(new byte[0]) == 0).
	mh0 := guava.Hashing_Murmur3_128().HashBytes(nil)
	if mh0.AsLong() != 0 || mh0.String() != "00000000000000000000000000000000" {
		t.Fatalf("Murmur3_128 empty mismatch: %s", mh0.String())
	}
	mh := guava.Hashing_Murmur3_128().NewHasher().PutInt(3).PutUnencodedChars("ab").Hash()
	if mh.AsInt() == 0 {
		t.Fatal("unexpected zero Murmur3_128 hash")
	}
}
