package javanio_test

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javanio"
)

func TestCharsets(t *testing.T) {
	if !javanio.Charset_IsSupported("utf-8") || !javanio.Charset_IsSupported("US-ASCII") || !javanio.Charset_IsSupported("Shift_JIS") {
		t.Fatal("expected UTF-8, US-ASCII, Shift_JIS to be supported")
	}
	if javanio.Charset_IsSupported("no-such-charset-123") {
		t.Fatal("unexpected supported charset")
	}
	utf8 := javanio.Charset_ForName("UTF8")
	if !utf8.Equals(javanio.StandardCharsets_UTF_8) {
		t.Fatalf("expected UTF8 alias to equal UTF-8")
	}
	ascii := javanio.StandardCharsets_US_ASCII
	if !ascii.CanEncodeChar('A') || ascii.CanEncodeChar(0x00E9) {
		t.Fatalf("unexpected US-ASCII CanEncodeChar")
	}
	enc := ascii.Encode("Aé😀")
	if string(enc) != "A??" {
		t.Fatalf("expected 'A??', got %q", string(enc))
	}
	utf16 := javanio.StandardCharsets_UTF_16
	b16 := utf16.Encode("Hi😀")
	if got := utf16.Decode(b16); got != "Hi😀" {
		t.Fatalf("UTF-16 roundtrip mismatch: %q", got)
	}
}

func TestPathAndRelativize(t *testing.T) {
	p := javanio.Path_Of("a//b/", "c", "../d/file.js").Normalize()
	if p.String() != "a/b/d/file.js" {
		t.Fatalf("got %q", p.String())
	}
	sib := javanio.Path_Of("baz/bam/qux.js").ResolveSiblingString("../foo/bar.js").Normalize()
	if sib.String() != "baz/foo/bar.js" {
		t.Fatalf("ResolveSibling normalize mismatch: got %q", sib.String())
	}
	rel := javanio.Path_Of("/a/b/c").Relativize(javanio.Path_Of("/a/d/e.js"))
	if rel.String() != "../../d/e.js" {
		t.Fatalf("Relativize mismatch: got %q", rel.String())
	}
}

type collectVisitor struct {
	javanio.SimpleFileVisitor
	matched []string
	matcher javanio.PathMatcher
}

func (v *collectVisitor) VisitFile(file *javanio.Path, _ *javanio.BasicFileAttributes) (javanio.FileVisitResult, error) {
	if v.matcher.Matches(file) || v.matcher.Matches(file.Normalize()) {
		v.matched = append(v.matched, file.Normalize().String())
	}
	return javanio.FileVisitResult_CONTINUE, nil
}

func TestFilesAndGlobWalk(t *testing.T) {
	tmp := t.TempDir()
	sub := filepath.Join(tmp, "src", "sub")
	pDir := javanio.Path_Of(sub)
	if _, err := javanio.Files_CreateDirectories(pDir); err != nil {
		t.Fatal(err)
	}
	f1 := javanio.Path_Of(tmp, "src", "a.js")
	f2 := javanio.Path_Of(tmp, "src", "sub", "b.js")
	f3 := javanio.Path_Of(tmp, "src", "sub", "ignore.txt")
	if _, err := javanio.Files_WriteString(f1, "const a = 1;\n"); err != nil {
		t.Fatal(err)
	}
	if _, err := javanio.Files_WriteString(f2, "const b = 2;\n"); err != nil {
		t.Fatal(err)
	}
	if _, err := javanio.Files_WriteString(f3, "hello"); err != nil {
		t.Fatal(err)
	}
	got, err := javanio.Files_ReadString(f1)
	if err != nil || got != "const a = 1;\n" {
		t.Fatalf("Files_ReadString got %q, err %v", got, err)
	}
	fs := javanio.FileSystems_GetDefault()
	matcher := fs.GetPathMatcher("glob:" + tmp + "/src/**.js")
	vis := &collectVisitor{matcher: matcher}
	if _, err := javanio.Files_WalkFileTree(javanio.Path_Of(tmp, "src"), true, 100, vis); err != nil {
		t.Fatal(err)
	}
	if len(vis.matched) != 2 {
		t.Fatalf("expected 2 matched .js files, got %v", vis.matched)
	}
	// Verify APPEND option
	if _, err := javanio.Files_WriteString(f3, " world", javanio.StandardOpenOption_WRITE, javanio.StandardOpenOption_APPEND); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(f3.String())
	if err != nil || string(raw) != "hello world" {
		t.Fatalf("expected 'hello world', got %q (err=%v)", string(raw), err)
	}
}

func TestURI(t *testing.T) {
	u, err := javanio.NewURI("https://example.com:8080/dir/sub/app%201.js?x=1#top")
	if err != nil {
		t.Fatal(err)
	}
	if !u.IsAbsolute() {
		t.Fatal("expected absolute URI")
	}
	p, ok := u.GetPath()
	if !ok || p != "/dir/sub/app 1.js" {
		t.Fatalf("unexpected decoded path: %q (%v)", p, ok)
	}
	rel := u.ResolveString("../other.js")
	if rel.String() != "https://example.com:8080/dir/other.js" {
		t.Fatalf("URI Resolve mismatch: %q", rel.String())
	}
	if _, err := javanio.NewURI("http://bad host/x"); err == nil {
		t.Fatal("expected URISyntaxException on space in host")
	}
}

func TestCharsetFastPathsAndStrictMalformedUTF8(t *testing.T) {
	// 1. Aligned and unaligned pure-ASCII lengths across UTF-8, US-ASCII, and ISO-8859-1.
	for _, n := range []int{0, 1, 7, 8, 9, 15, 16, 17, 4097} {
		want := strings.Repeat("aB9_", n/4+1)[:n]
		b := []byte(want)
		for _, cs := range []*javanio.Charset{
			javanio.StandardCharsets_UTF_8,
			javanio.StandardCharsets_US_ASCII,
			javanio.StandardCharsets_ISO_8859_1,
		} {
			got, err := cs.DecodeStrict(b)
			if err != nil || got != want {
				t.Fatalf("%s DecodeStrict(len=%d) got %q, err=%v", cs.Name(), n, got, err)
			}
			if dec := cs.Decode(b); dec != want {
				t.Fatalf("%s Decode(len=%d) got %q", cs.Name(), n, dec)
			}
			if string(cs.Encode(want)) != want {
				t.Fatalf("%s Encode(len=%d) roundtrip mismatch", cs.Name(), n)
			}
		}
	}

	// 2. ISO-8859-1 full 0x00..0xFF byte range round-trip.
	latin := make([]byte, 256)
	for i := range latin {
		latin[i] = byte(i)
	}
	latinStr, err := javanio.StandardCharsets_ISO_8859_1.DecodeStrict(latin)
	if err != nil {
		t.Fatalf("ISO-8859-1 DecodeStrict unexpected err: %v", err)
	}
	if got := javanio.StandardCharsets_ISO_8859_1.Encode(latinStr); string(got) != string(latin) {
		t.Fatalf("ISO-8859-1 full range roundtrip mismatch")
	}

	// 3. US-ASCII rejects >= 0x80 at every word alignment offset (0..15).
	for pos := 0; pos < 16; pos++ {
		buf := []byte("0123456789abcdef")
		buf[pos] = 0x80
		_, err := javanio.StandardCharsets_US_ASCII.DecodeStrict(buf)
		mal, ok := err.(*javanio.MalformedInputException)
		if !ok || mal.InputLength != 1 {
			t.Fatalf("US-ASCII pos=%d expected MalformedInputException(1), got %T %v", pos, err, err)
		}
		if rep := javanio.StandardCharsets_US_ASCII.Decode(buf); !strings.ContainsRune(rep, '\uFFFD') {
			t.Fatalf("US-ASCII Decode(pos=%d) expected replacement char, got %q", pos, rep)
		}
	}

	// 4. Strict UTF-8 malformed cases (surrogates, overlongs, >U+10FFFF, truncated sequences)
	// and exact MalformedInputException InputLength parity.
	malformedCases := []struct {
		name    string
		in      []byte
		wantLen int
	}{
		{"lone_cont", []byte{'a', 0x80, 'b'}, 1},
		{"overlong_2byte", []byte{0xC0, 0xAF}, 1},
		{"overlong_3byte", []byte{0xE0, 0x80, 0xAF}, 1},
		{"surrogate_high", []byte{0xED, 0xA0, 0x80}, 1},
		{"surrogate_low", []byte{0xED, 0xBF, 0xBF}, 1},
		{"above_max_codepoint", []byte{0xF4, 0x90, 0x80, 0x80}, 1},
		{"invalid_ff", []byte("const x = 1;\xff"), 1},
		{"trunc_2byte", []byte{0xC2}, 1},
		{"trunc_3byte_1", []byte{0xE2}, 1},
		{"trunc_3byte_2", []byte{0xE2, 0x82}, 2},
		{"bad_3byte_third", []byte{0xE2, 0x82, 'x'}, 2},
		{"trunc_4byte_1", []byte{0xF0}, 1},
		{"trunc_4byte_2", []byte{0xF0, 0x9F}, 2},
		{"trunc_4byte_3", []byte{0xF0, 0x9F, 0x98}, 3},
		{"bad_4byte_fourth", []byte{0xF0, 0x9F, 0x98, 'x'}, 3},
	}
	tmp := t.TempDir()
	for _, tc := range malformedCases {
		_, err := javanio.StandardCharsets_UTF_8.DecodeStrict(tc.in)
		mal, ok := err.(*javanio.MalformedInputException)
		if !ok {
			t.Fatalf("%s: expected *MalformedInputException, got %T (%v)", tc.name, err, err)
		}
		if mal.InputLength != tc.wantLen {
			t.Fatalf("%s: MalformedInputException.InputLength = %d, want %d", tc.name, mal.InputLength, tc.wantLen)
		}
		p := javanio.Path_Of(tmp, tc.name+".js")
		if err := os.WriteFile(p.String(), tc.in, 0644); err != nil {
			t.Fatal(err)
		}
		if _, fErr := javanio.Files_ReadString(p); fErr == nil {
			t.Fatalf("%s: expected Files_ReadString to fail on malformed UTF-8", tc.name)
		}
	}

	// 5. Allocation check: UTF-8 DecodeStrict on 64 KiB valid JS does exactly 1 allocation (string(b)).
	jsBytes := []byte(strings.Repeat("export function compute(x) { return x + 'λ😀'; }\n", 1024))
	allocs := testing.AllocsPerRun(20, func() {
		s, err := javanio.StandardCharsets_UTF_8.DecodeStrict(jsBytes)
		if err != nil || len(s) != len(jsBytes) {
			t.Fatalf("UTF-8 DecodeStrict failed: %v", err)
		}
	})
	if allocs > 1 {
		t.Fatalf("expected <= 1 allocation for UTF-8 DecodeStrict, got %.1f", allocs)
	}
}

func TestFilesReadStringZeroCopyFastPath(t *testing.T) {
	tmp := t.TempDir()

	// 1. Empty file across UTF-8, US-ASCII, and ISO-8859-1.
	emptyPath := javanio.Path_Of(tmp, "empty.js")
	if err := os.WriteFile(emptyPath.String(), nil, 0644); err != nil {
		t.Fatal(err)
	}
	for _, cs := range []*javanio.Charset{
		javanio.StandardCharsets_UTF_8,
		javanio.StandardCharsets_US_ASCII,
		javanio.StandardCharsets_ISO_8859_1,
	} {
		got, err := javanio.Files_ReadStringCharset(emptyPath, cs)
		if err != nil || got != "" {
			t.Fatalf("empty file (%s): got %q, err=%v", cs.Name(), got, err)
		}
	}

	// 2. >64KB multi-chunk UTF-8 file with 2-byte, 3-byte, and 4-byte UTF-8 runes
	// deliberately straddling the 32768-byte and 65536-byte scratch chunk boundaries.
	prefix1 := strings.Repeat("x", 32767) // 1 byte before 32768
	mid := "é" + strings.Repeat("y", 32765) + "😀" + strings.Repeat("z", 15000) + "λ€🚀\n"
	wantUTF8 := prefix1 + mid
	if len(wantUTF8) <= 65536 {
		t.Fatalf("expected >64KB test payload, got %d", len(wantUTF8))
	}
	utf8Path := javanio.Path_Of(tmp, "large_utf8.js")
	if err := os.WriteFile(utf8Path.String(), []byte(wantUTF8), 0644); err != nil {
		t.Fatal(err)
	}
	gotUTF8, err := javanio.Files_ReadString(utf8Path)
	if err != nil || gotUTF8 != wantUTF8 {
		t.Fatalf("large multi-chunk UTF-8 mismatch: len(got)=%d len(want)=%d err=%v", len(gotUTF8), len(wantUTF8), err)
	}

	// 3. Malformed UTF-8 in >64KB file (truncated 3-byte sequence straddling 32768 boundary)
	// must still return *MalformedInputException with exact InputLength!
	badUTF8 := make([]byte, 70000)
	for i := range badUTF8 {
		badUTF8[i] = 'a'
	}
	badUTF8[32767] = 0xE2
	badUTF8[32768] = 0x82
	badUTF8[32769] = 'z' // invalid 3rd byte -> InputLength = 2
	badUTF8Path := javanio.Path_Of(tmp, "bad_utf8.js")
	if err := os.WriteFile(badUTF8Path.String(), badUTF8, 0644); err != nil {
		t.Fatal(err)
	}
	_, err = javanio.Files_ReadString(badUTF8Path)
	mal, ok := err.(*javanio.MalformedInputException)
	if !ok || mal.InputLength != 2 {
		t.Fatalf("expected *MalformedInputException(InputLength=2) on malformed >64KB UTF-8, got %T (%v)", err, err)
	}

	// 4. US-ASCII >64KB valid + non-ASCII rejection in second chunk.
	wantASCII := strings.Repeat("const value = 1234567890;\n", 3000) // ~78 KB
	asciiPath := javanio.Path_Of(tmp, "large_ascii.js")
	if err := os.WriteFile(asciiPath.String(), []byte(wantASCII), 0644); err != nil {
		t.Fatal(err)
	}
	gotASCII, err := javanio.Files_ReadStringCharset(asciiPath, javanio.StandardCharsets_US_ASCII)
	if err != nil || gotASCII != wantASCII {
		t.Fatalf("US-ASCII >64KB read failed: err=%v len=%d", err, len(gotASCII))
	}
	_, err = javanio.Files_ReadStringCharset(utf8Path, javanio.StandardCharsets_US_ASCII)
	if malASCII, ok := err.(*javanio.MalformedInputException); !ok || malASCII.InputLength != 1 {
		t.Fatalf("expected *MalformedInputException(1) reading UTF-8 as US-ASCII, got %T (%v)", err, err)
	}

	// 5. ISO-8859-1 >64KB pure-ASCII fast path + >64KB non-ASCII Latin-1 decoding.
	gotISOASCII, err := javanio.Files_ReadStringCharset(asciiPath, javanio.StandardCharsets_ISO_8859_1)
	if err != nil || gotISOASCII != wantASCII {
		t.Fatalf("ISO-8859-1 ASCII >64KB read failed: err=%v", err)
	}
	latinBytes := make([]byte, 70000)
	for i := range latinBytes {
		latinBytes[i] = byte(i & 0xFF)
	}
	latinPath := javanio.Path_Of(tmp, "large_latin1.bin")
	if err := os.WriteFile(latinPath.String(), latinBytes, 0644); err != nil {
		t.Fatal(err)
	}
	wantLatin, err := javanio.StandardCharsets_ISO_8859_1.DecodeStrict(latinBytes)
	if err != nil {
		t.Fatal(err)
	}
	gotLatin, err := javanio.Files_ReadStringCharset(latinPath, javanio.StandardCharsets_ISO_8859_1)
	if err != nil || gotLatin != wantLatin {
		t.Fatalf("ISO-8859-1 non-ASCII >64KB mismatch: err=%v len(got)=%d len(want)=%d", err, len(gotLatin), len(wantLatin))
	}
}

func TestWave27AFilesReadStringNamedPipeFallback(t *testing.T) {
	tmp := t.TempDir()
	fifoPath := filepath.Join(tmp, "stream.fifo")
	if err := syscall.Mkfifo(fifoPath, 0o600); err != nil {
		t.Skipf("Mkfifo unsupported in temp dir: %v", err)
	}
	want := "export const answer = 'λ-fifo-42';\n"
	errCh := make(chan error, 1)
	go func() {
		f, err := os.OpenFile(fifoPath, os.O_WRONLY, 0)
		if err != nil {
			errCh <- err
			return
		}
		_, err = f.WriteString(want)
		_ = f.Close()
		errCh <- err
	}()
	got, err := javanio.Files_ReadString(javanio.Path_Of(fifoPath))
	if werr := <-errCh; werr != nil {
		t.Fatalf("fifo writer failed: %v", werr)
	}
	if err != nil || got != want {
		t.Fatalf("FIFO Files_ReadString got %q, err=%v (want %q)", got, err, want)
	}
}

func BenchmarkWave27AFilesReadString(b *testing.B) {
	tmp := b.TempDir()
	smallJS := strings.Repeat("export function add(a, b) { return a + b; }\n", 90) // ~4 KiB (<= 32 KiB scratch fast path)
	largeJS := strings.Repeat("export function mul(a, b) { return a * b; }\n", 2200) // ~96 KiB (> 32 KiB exact-size ReadFull)
	smallPath := javanio.Path_Of(tmp, "small.js")
	largePath := javanio.Path_Of(tmp, "large.js")
	if err := os.WriteFile(smallPath.String(), []byte(smallJS), 0o644); err != nil {
		b.Fatal(err)
	}
	if err := os.WriteFile(largePath.String(), []byte(largeJS), 0o644); err != nil {
		b.Fatal(err)
	}

	b.Run("Small4KB_UTF8", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(smallJS)))
		for i := 0; i < b.N; i++ {
			s, err := javanio.Files_ReadString(smallPath)
			if err != nil || len(s) != len(smallJS) {
				b.Fatal(err)
			}
		}
	})
	b.Run("Large96KB_UTF8", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(largeJS)))
		for i := 0; i < b.N; i++ {
			s, err := javanio.Files_ReadString(largePath)
			if err != nil || len(s) != len(largeJS) {
				b.Fatal(err)
			}
		}
	})
}

func TestWave30CStreamingSingleAllocAndBoundaryMalformedUTF8(t *testing.T) {
	tmp := t.TempDir()

	// 1. Literal U+FFFD (0xEF 0xBF 0xBD) is valid UTF-8 and must NOT be rejected by DecodeStrict or Files_ReadString!
	withReplacement := "export const sym = '\uFFFD'; // valid U+FFFD literal\n" + strings.Repeat("let x = 1;\n", 4000)
	fffdPath := javanio.Path_Of(tmp, "valid_fffd.js")
	if err := os.WriteFile(fffdPath.String(), []byte(withReplacement), 0o644); err != nil {
		t.Fatal(err)
	}
	if s, err := javanio.StandardCharsets_UTF_8.DecodeStrict([]byte(withReplacement)); err != nil || s != withReplacement {
		t.Fatalf("DecodeStrict rejected valid U+FFFD: err=%v", err)
	}
	if s, err := javanio.Files_ReadString(fffdPath); err != nil || s != withReplacement {
		t.Fatalf("Files_ReadString rejected valid U+FFFD: err=%v", err)
	}

	// 2. Truncated multi-byte UTF-8 sequence ending RIGHT AT EOF on a >32KB boundary (32768 + 2 bytes: [0xF0, 0x9F]).
	truncAtBoundary := make([]byte, 32768+2)
	for i := 0; i < 32768; i++ {
		truncAtBoundary[i] = 'a'
	}
	truncAtBoundary[32768] = 0xF0
	truncAtBoundary[32769] = 0x9F
	truncPath := javanio.Path_Of(tmp, "trunc_boundary.js")
	if err := os.WriteFile(truncPath.String(), truncAtBoundary, 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := javanio.Files_ReadString(truncPath)
	mal, ok := err.(*javanio.MalformedInputException)
	if !ok || mal.InputLength != 2 {
		t.Fatalf("expected *MalformedInputException(InputLength=2) at EOF after chunk carry, got %T (%v)", err, err)
	}

	// Also test truncated 3-byte sequence starting at 32767 ([0xE2, 0x82] then immediate EOF at 32769!).
	truncStraddle := make([]byte, 32769)
	for i := 0; i < 32767; i++ {
		truncStraddle[i] = 'b'
	}
	truncStraddle[32767] = 0xE2
	truncStraddle[32768] = 0x82
	truncStraddlePath := javanio.Path_Of(tmp, "trunc_straddle.js")
	if err := os.WriteFile(truncStraddlePath.String(), truncStraddle, 0o644); err != nil {
		t.Fatal(err)
	}
	_, err = javanio.Files_ReadString(truncStraddlePath)
	mal, ok = err.(*javanio.MalformedInputException)
	if !ok || mal.InputLength != 2 {
		t.Fatalf("expected *MalformedInputException(InputLength=2) on straddling EOF truncation, got %T (%v)", err, err)
	}

	// 3. Verify Large96KB Files_ReadString eliminates the 2x []byte + string double allocation!
	largeJS := strings.Repeat("export function mul(a, b) { return a * b; }\n", 2200) // 96,800 B
	largePath := javanio.Path_Of(tmp, "large96k.js")
	if err := os.WriteFile(largePath.String(), []byte(largeJS), 0o644); err != nil {
		t.Fatal(err)
	}
	// Warm pool once
	if _, err := javanio.Files_ReadString(largePath); err != nil {
		t.Fatal(err)
	}
	res := testing.Benchmark(func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			s, err := javanio.Files_ReadString(largePath)
			if err != nil || len(s) != len(largeJS) {
				b.Fatal(err)
			}
		}
	})
	// Single-buffer strings.Builder allocation is ~98.3 KB + ~400 B os.File/Stat overhead (< 115 KB, vs 197 KB before!).
	if res.AllocedBytesPerOp() > 115_000 {
		t.Fatalf("Large96KB Files_ReadString double-allocated %d B/op (want <= 115,000 B/op for 96,800 B file)", res.AllocedBytesPerOp())
	}
}
