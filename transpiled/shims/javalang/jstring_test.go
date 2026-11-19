package javalang

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"testing"
	"unicode/utf16"
)

func TestJStringASCIIAndSurrogates(t *testing.T) {
	ascii := "hello world"
	if !IsASCII(ascii) {
		t.Errorf("IsASCII(%q) should be true", ascii)
	}
	if StringLength(ascii) != len(ascii) {
		t.Errorf("StringLength(%q) should be %d, got %d", ascii, len(ascii), StringLength(ascii))
	}
	if StringCharAt(ascii, 1) != 'e' {
		t.Errorf("StringCharAt(%q, 1) should be 'e', got %c", ascii, StringCharAt(ascii, 1))
	}
	if StringSubstring(ascii, 0, 5) != "hello" {
		t.Errorf("StringSubstring(%q, 0, 5) should be \"hello\", got %q", ascii, StringSubstring(ascii, 0, 5))
	}

	// Test surrogate pair: emoji U+1F600 (Grinning Face): \uD83D \uDE00 (length 2 in UTF-16, length 4 in UTF-8)
	emoji := "a😀b"
	if IsASCII(emoji) {
		t.Errorf("IsASCII(%q) should be false", emoji)
	}
	if StringLength(emoji) != 4 { // 'a' (1) + emoji (2) + 'b' (1)
		t.Errorf("StringLength(%q) should be 4, got %d", emoji, StringLength(emoji))
	}
	if StringCharAt(emoji, 0) != 'a' {
		t.Errorf("char at 0 should be 'a'")
	}
	if StringCharAt(emoji, 1) != 0xD83D {
		t.Errorf("char at 1 should be high surrogate 0xD83D, got 0x%X", StringCharAt(emoji, 1))
	}
	if StringCharAt(emoji, 2) != 0xDE00 {
		t.Errorf("char at 2 should be low surrogate 0xDE00, got 0x%X", StringCharAt(emoji, 2))
	}
	if StringCharAt(emoji, 3) != 'b' {
		t.Errorf("char at 3 should be 'b'")
	}
	if StringSubstring(emoji, 1, 3) != "😀" {
		t.Errorf("StringSubstring(%q, 1, 3) should be emoji, got %q", emoji, StringSubstring(emoji, 1, 3))
	}

	// Test lone surrogate roundtrip
	lone := []uint16{'x', 0xD800, 'y'}
	sLone := StringFromUTF16(lone)
	u16Out := StringToUTF16(sLone)
	if len(u16Out) != 3 || u16Out[0] != 'x' || u16Out[1] != 0xD800 || u16Out[2] != 'y' {
		t.Errorf("lone surrogate roundtrip failed: got %v, want %v", u16Out, lone)
	}
}

func TestJStringSearchingAndComparing(t *testing.T) {
	s := "banana"
	if StringIndexOf(s, "an") != 1 {
		t.Errorf("StringIndexOf: got %d, want 1", StringIndexOf(s, "an"))
	}
	if StringLastIndexOf(s, "an") != 3 {
		t.Errorf("StringLastIndexOf: got %d, want 3", StringLastIndexOf(s, "an"))
	}
	if StringIndexOfChar(s, 'n') != 2 {
		t.Errorf("StringIndexOfChar: got %d, want 2", StringIndexOfChar(s, 'n'))
	}

	if StringCompareTo("abc", "abd") >= 0 {
		t.Errorf("StringCompareTo(\"abc\", \"abd\") should be negative")
	}
	if StringCompareTo("abc", "abc") != 0 {
		t.Errorf("StringCompareTo equal strings should be 0")
	}

	// Verify Java polynomial hash code on known vector
	if StringHashCode("test") != 3556498 {
		t.Errorf("StringHashCode(\"test\") got %d, want 3556498", StringHashCode("test"))
	}
}

func naiveUTF16CompareTo(a, b string) int {
	ua := StringToUTF16(a)
	ub := StringToUTF16(b)
	n := min(len(ua), len(ub))
	for i := 0; i < n; i++ {
		if ua[i] != ub[i] {
			return int(ua[i]) - int(ub[i])
		}
	}
	return len(ua) - len(ub)
}

func naiveUTF16HashCode(s string) int32 {
	var h int32
	for _, u := range StringToUTF16(s) {
		h = 31*h + int32(u)
	}
	return h
}

func TestStringHashCodeCompareToAndBlockMeta(t *testing.T) {
	loneHi := StringFromChar(0xD800)
	loneLo := StringFromChar(0xDC00)
	samples := []string{
		"",
		"a",
		"ab",
		"abc",
		"abcd",
		"abcde",
		"identifier_with_long_common_prefix_alpha",
		"identifier_with_long_common_prefix_beta",
		"hello\uFF21world", // U+FF21 > 0xD83D in UTF-16, < U+1F600 in UTF-8!
		"hello\U0001F600world",
		"prefix" + loneHi + "tail",
		"prefix" + loneLo + "tail",
		"prefix\U00010000tail",
		"\xF0\x90\x80\x80\x80A",
		"\xF0\x90\x80\x80\x80\xC2\x80",
		strings.Repeat("a", 63) + "😀" + strings.Repeat("b", 70) + "世界" + strings.Repeat("c", 80),
	}

	for _, s := range samples {
		if got, want := StringHashCode(s), naiveUTF16HashCode(s); got != want {
			t.Errorf("StringHashCode(%q) = %d, want %d", s, got, want)
		}
	}
	for _, a := range samples {
		for _, b := range samples {
			if got, want := StringCompareTo(a, b), naiveUTF16CompareTo(a, b); got != want {
				t.Errorf("StringCompareTo(%q, %q) = %d, want %d", a, b, got, want)
			}
		}
	}

	// Specifically exercise block boundary where UTF-16 index 64 lands on a low surrogate!
	// 63 ASCII chars (u=0..62), '😀' at u=63..64, 70 'b's (u=65..134), '世界' (u=135..136), 80 'c's (u=137..216).
	big := strings.Repeat("a", 63) + "😀" + strings.Repeat("b", 70) + "世界" + strings.Repeat("c", 80)
	uBig := StringToUTF16(big)
	for begin := 0; begin <= len(uBig); begin += 7 {
		for end := begin; end <= len(uBig); end += 13 {
			got := StringSubstring(big, begin, end)
			want := StringFromUTF16(uBig[begin:end])
			if got != want {
				t.Fatalf("StringSubstring(big, %d, %d) mismatch: got %q, want %q", begin, end, got, want)
			}
		}
		if got, want := StringSubstringFrom(big, begin), StringFromUTF16(uBig[begin:]); got != want {
			t.Fatalf("StringSubstringFrom(big, %d) mismatch", begin)
		}
	}
	// Check exact surrogate-splitting boundaries around u=63, 64, 65
	for _, br := range [][2]int{{62, 65}, {63, 65}, {63, 64}, {64, 65}, {65, 137}} {
		if got, want := StringSubstring(big, br[0], br[1]), StringFromUTF16(uBig[br[0]:br[1]]); got != want {
			t.Fatalf("StringSubstring(big, %d, %d) = %q, want %q", br[0], br[1], got, want)
		}
	}
	if idx := StringIndexOf(big, "世界"); idx != 135 {
		t.Errorf("StringIndexOf(big, '世界') = %d, want 135", idx)
	}
	if idx := StringIndexOfFrom(big, "c", 65); idx != 137 {
		t.Errorf("StringIndexOfFrom(big, 'c', 65) = %d, want 137", idx)
	}
	if idx := StringLastIndexOf(big, "世界"); idx != 135 {
		t.Errorf("StringLastIndexOf(big, '世界') = %d, want 135", idx)
	}

	// Prime cache, then verify 0 allocs on hot paths!
	_ = StringLength(big)
	s1 := "com.google.javascript.jscomp.TypeCheck$1"
	s2 := "com.google.javascript.jscomp.TypeCheck$2"
	allocs := testing.AllocsPerRun(100, func() {
		_ = StringHashCode(big)
		_ = StringCompareTo(s1, s2)
		_ = StringCompareTo("hello\uFF21world", "hello\U0001F600world")
		_ = StringSubstring(big, 65, 137)
		_ = StringIndexOf(big, "世界")
		_ = StringIndexOfFrom(big, "c", 65)
	})
	if allocs != 0 {
		t.Errorf("expected 0 allocs/op on cached/zero-alloc jstring hot paths, got %v", allocs)
	}
}

func TestLargeCopyrightJSFileZeroUTF16SliceAndNoASCIIEviction(t *testing.T) {
	// Simulate a 200 KB googleapis JS file with a single '©' in the header.
	jsFile := "// Copyright © 2026 Google LLC. All rights reserved.\n" + strings.Repeat("goog.provide('gapi.client.compute.Instance');\n", 4400)
	m := tryGetStrMeta(jsFile)
	if m == nil {
		t.Fatal("expected cached strMeta for 200 KB non-ASCII JS file")
	}
	if len(m.checkpoints) != 0 {
		t.Fatalf("expected 0 checkpoints for single header symbol, got %d", len(m.checkpoints))
	}
	wantLen := len(jsFile) - 1 // '©' is 2 UTF-8 bytes -> 1 UTF-16 unit
	if StringLength(jsFile) != wantLen {
		t.Fatalf("StringLength = %d, want %d", StringLength(jsFile), wantLen)
	}
	if ch := StringCharAt(jsFile, 13); ch != 0x00A9 {
		t.Fatalf("StringCharAt(13) = 0x%X, want 0x00A9", ch)
	}
	if ch := StringCharAt(jsFile, wantLen-2); ch != ';' {
		t.Fatalf("StringCharAt(end-2) = %q, want ';'", rune(ch))
	}

	// Pre-build 5,000 distinct 80..200-byte pure-ASCII qualified type names.
	typeNames := make([]string, 5000)
	for i := range typeNames {
		typeNames[i] = fmt.Sprintf("com.google.api.services.compute.v1.model.AutoscalingPolicyMetricConfiguration_%04d", i)
	}

	// Accessing thousands of pure-ASCII qualified names must allocate 0 objects
	// and MUST NOT evict jsFile from strMetaCache!
	asciiAllocs := testing.AllocsPerRun(1, func() {
		for _, tn := range typeNames {
			if !IsASCII(tn) || StringLength(tn) != len(tn) {
				t.Fatal("ASCII typeName check failed")
			}
			if StringCharAt(tn, 0) != 'c' || StringCharAt(tn, len(tn)-1) != uint16(tn[len(tn)-1]) {
				t.Fatal("StringCharAt on ASCII typeName failed")
			}
			_ = StringSubstring(tn, 4, 30)
		}
	})
	if asciiAllocs != 0 {
		t.Fatalf("pure-ASCII qualified names allocated %v objects, want 0", asciiAllocs)
	}

	postAllocs := testing.AllocsPerRun(20, func() {
		_ = StringLength(jsFile)
		_ = StringCharAt(jsFile, 0)
		_ = StringCharAt(jsFile, 13)
		_ = StringCharAt(jsFile, 100000)
		_ = StringSubstring(jsFile, 10, wantLen-10)
		_ = StringLastIndexOfCharFrom(jsFile, '\n', 150000)
		_ = StringIndexOfChar(jsFile, 0xFEFF)
		_ = StringRegionMatches(jsFile, 53, "goog.provide", 0, 12)
	})
	if postAllocs != 0 {
		t.Fatalf("expected 0 allocs/op on cached 200 KB non-ASCII JS file after 5,000 ASCII type lookups, got %v", postAllocs)
	}

	// Token substrings (<= 2048 bytes) from 200 KB JS file must allocate 0 objects on
	// warm deduplicated hits, while > 2048-byte small slices allocate 1 clone (no 200 KB []uint16 slice).
	subAllocs := testing.AllocsPerRun(20, func() {
		tok45 := StringSubstring(jsFile, 1000, 1045)
		tok120 := StringSubstring(jsFile, 1000, 1120)
		tok600 := StringSubstring(jsFile, 1000, 1600)
		if len(tok45) != 45 || len(tok120) != 120 || len(tok600) != 600 {
			t.Fatalf("unexpected token lengths %d, %d, %d", len(tok45), len(tok120), len(tok600))
		}
	})
	if subAllocs != 0 {
		t.Fatalf("expected 0 allocs/op on deduplicated <=2048B token/comment slices from 200 KB file, got %v", subAllocs)
	}
	longSubAllocs := testing.AllocsPerRun(20, func() {
		tok := StringSubstring(jsFile, 1000, 3500)
		if len(tok) != 2500 {
			t.Fatalf("unexpected token length %d", len(tok))
		}
	})
	if longSubAllocs != 1 {
		t.Fatalf("expected 1 small clone alloc/op on 2500B (>2048B) slice from 200 KB file, got %v", longSubAllocs)
	}
}

func TestSubstringUnpinsLargeSourceBuffers(t *testing.T) {
	runtime.GC()
	runtime.GC()
	var before runtime.MemStats
	runtime.ReadMemStats(&before)

	const numFiles = 400
	retained := make([]string, 0, numFiles*4)
	func() {
		for i := 0; i < numFiles; i++ {
			// ~100 KB non-ASCII and pure-ASCII JS source buffers (~40 MB total).
			var src string
			if i%2 == 0 {
				src = fmt.Sprintf("// Copyright © 2026 File %04d\n/** @suppress {checkTypes} */\n", i) +
					strings.Repeat("goog.provide('gapi.client.service.ResourceIdentifier');\n", 1800) +
					fmt.Sprintf("export const tail_%04d = 1;", i)
			} else {
				src = fmt.Sprintf("// Pure ASCII File %04d\n/** @suppress {extraRequire} */\n", i) +
					strings.Repeat("const localVariableIdentifier = window.document.body;\n", 1850) +
					fmt.Sprintf("export const end_%04d = 2;", i)
			}
			uLen := StringLength(src)
			// 1. Small identifier token via StringSubstring
			retained = append(retained, StringSubstring(src, 60, 83))
			// 2. Small trailing token via StringSubstringFrom
			retained = append(retained, StringSubstringFrom(src, uLen-20))
			// 3. Multi-hop comment -> @suppress tag slice
			comment := StringSubstring(src, 0, 600)
			retained = append(retained, StringSubstring(comment, 28, 52))
			// 4. Trimmed small slice
			retained = append(retained, StringTrim("   "+StringSubstring(src, 60, 80)+"   \n"))
		}
		// Flush strMetaCache / hotStrMeta0..1 with tiny 600-byte strings so cache slots don't hold the last 256 files.
		for i := 0; i < 512; i++ {
			evict := fmt.Sprintf("// evict © %04d ", i) + strings.Repeat("x", 550)
			_ = StringLength(evict)
		}
	}()

	runtime.GC()
	runtime.GC()
	var after runtime.MemStats
	runtime.ReadMemStats(&after)

	if len(retained) != numFiles*4 || retained[0] == "" {
		t.Fatalf("unexpected retained tokens: %d", len(retained))
	}
	var liveDelta int64
	if after.HeapAlloc > before.HeapAlloc {
		liveDelta = int64(after.HeapAlloc - before.HeapAlloc)
	}
	// 400 * 100 KB = 40 MB of source files; retained tokens + small cache strings are < 1 MB.
	if liveDelta > 8<<20 {
		t.Fatalf("large source buffers still pinned in heap: HeapAlloc grew by %d bytes (%.2f MB)", liveDelta, float64(liveDelta)/(1<<20))
	}
	runtime.KeepAlive(retained)
}

func TestDenseCheckpointsDirtyWTF8AndConcurrentRace(t *testing.T) {
	loneHi := StringFromChar(0xD83D)
	loneLo := StringFromChar(0xDE00)
	cases := []string{
		// Dense multi-block emojis + BMP CJK + short (<64) and long (>=64) ASCII runs
		strings.Repeat("🚀漢字αβ", 45) + strings.Repeat("x", 150) + strings.Repeat("😀Ω", 50) + strings.Repeat("z", 90),
		// Uncached (<48) and cached (>=48) WTF-8 lone surrogates & invalid UTF-8 bytes
		"short" + loneHi + "mid\U0001F680" + loneLo,
		strings.Repeat("p", 30) + loneHi + strings.Repeat("q", 70) + "\xff\xfe" + strings.Repeat("😀", 40) + loneLo + strings.Repeat("r", 80),
	}

	for ci, s := range cases {
		u := StringToUTF16(s)
		if got := StringLength(s); got != len(u) {
			t.Fatalf("case %d: StringLength = %d, want %d", ci, got, len(u))
		}
		for idx := 0; idx < len(u); idx++ {
			if got := StringCharAt(s, idx); got != u[idx] {
				t.Fatalf("case %d: StringCharAt(%d) = 0x%X, want 0x%X", ci, idx, got, u[idx])
			}
			wantCP := int32(u[idx])
			if u[idx] >= 0xD800 && u[idx] <= 0xDBFF && idx+1 < len(u) && u[idx+1] >= 0xDC00 && u[idx+1] <= 0xDFFF {
				wantCP = 0x10000 + ((int32(u[idx]) - 0xD800) << 10) + (int32(u[idx+1]) - 0xDC00)
			}
			if got := StringCodePointAt(s, idx); got != wantCP {
				t.Fatalf("case %d: StringCodePointAt(%d) = 0x%X, want 0x%X", ci, idx, got, wantCP)
			}
		}
		for b := 0; b <= len(u); b += 11 {
			for e := b; e <= len(u); e += 19 {
				if got, want := StringSubstring(s, b, e), StringFromUTF16(u[b:e]); got != want {
					t.Fatalf("case %d: StringSubstring(%d, %d) mismatch", ci, b, e)
				}
			}
			if got, want := StringSubstringFrom(s, b), StringFromUTF16(u[b:]); got != want {
				t.Fatalf("case %d: StringSubstringFrom(%d) mismatch", ci, b)
			}
		}
		if ci == 0 {
			if got, want := StringLastIndexOfFrom(s, "😀Ω", 300), utf16LastIndexFromRef(u, StringToUTF16("😀Ω"), 300); got != want {
				t.Fatalf("StringLastIndexOfFrom = %d, want %d", got, want)
			}
			for _, ch := range []uint16{'x', 'z', 0x6F22, 0xD83D, 0xDE00, 0xFEFF} {
				if got, want := StringIndexOfChar(s, ch), slicesIndex(u, ch); got != want {
					t.Fatalf("StringIndexOfChar(0x%X) = %d, want %d", ch, got, want)
				}
				if got, want := StringLastIndexOfCharFrom(s, ch, 260), slicesLastIndexFrom(u, ch, 260); got != want {
					t.Fatalf("StringLastIndexOfCharFrom(0x%X, 260) = %d, want %d", ch, got, want)
				}
			}
		}
	}

	// Concurrent race stress across 16 goroutines
	pool := make([]string, 32)
	for i := range pool {
		pool[i] = fmt.Sprintf("// File %d © 2026\n", i) + strings.Repeat("var x = 'hello 世界 🚀';\n", 40+i)
	}
	var wg sync.WaitGroup
	for g := 0; g < 16; g++ {
		g := g
		wg.Add(1)
		go func() {
			defer wg.Done()
			for iter := 0; iter < 200; iter++ {
				s := pool[(g*7+iter)%len(pool)]
				n := StringLength(s)
				_ = StringCharAt(s, n/2)
				_ = StringSubstring(s, 5, n-5)
				_ = StringIndexOf(s, "世界")
				_ = StringLastIndexOfCharFrom(s, '\n', n-1)
			}
		}()
	}
	wg.Wait()
}

func slicesIndex(u []uint16, ch uint16) int {
	for i, c := range u {
		if c == ch {
			return i
		}
	}
	return -1
}

func slicesLastIndexFrom(u []uint16, ch uint16, from int) int {
	for i := min(from, len(u)-1); i >= 0; i-- {
		if u[i] == ch {
			return i
		}
	}
	return -1
}

func utf16LastIndexFromRef(u, sub []uint16, from int) int {
	for k := min(len(u)-len(sub), from); k >= 0; k-- {
		match := true
		for j := range sub {
			if u[k+j] != sub[j] {
				match = false
				break
			}
		}
		if match {
			return k
		}
	}
	return -1
}

func TestWordScanASCIIPrefixHashStartsEndsAndNegativeArraySizeEx(t *testing.T) {
	// Exercise every 32-byte / 8-byte / 1-byte lane in asciiPrefixLenRaw and StringHashCode.
	loneHi := StringFromChar(0xD83D)
	loneLo := StringFromChar(0xDE00)
	for n := 0; n <= 75; n++ {
		pure := strings.Repeat("a", n)
		if got := ASCIIPrefixLen(pure); got != n {
			t.Fatalf("ASCIIPrefixLen(len=%d) = %d, want %d", n, got, n)
		}
		if got, want := StringHashCode(pure), naiveUTF16HashCode(pure); got != want {
			t.Fatalf("StringHashCode(pure len=%d) = %d, want %d", n, got, want)
		}
		for _, tail := range []string{"©tail", "🚀end", loneHi + "x", loneLo + "y"} {
			mixed := pure + tail
			if got := asciiPrefixLenRaw(mixed); got != n {
				t.Fatalf("asciiPrefixLenRaw(n=%d, tail=%q) = %d, want %d", n, tail, got, n)
			}
			if got, want := StringHashCode(mixed), naiveUTF16HashCode(mixed); got != want {
				t.Fatalf("StringHashCode(n=%d, tail=%q) = %d, want %d", n, tail, got, want)
			}
			if !StringStartsWith(mixed, pure) || !StringEndsWith(mixed, tail) {
				t.Fatalf("StartsWith/EndsWith failed for n=%d tail=%q", n, tail)
			}
			if StringStartsWith(mixed, pure+"z") || StringEndsWith(mixed, "z"+tail) {
				t.Fatalf("StartsWith/EndsWith false positive for n=%d tail=%q", n, tail)
			}
		}
	}

	// Verify NegativeArraySizeException class info.
	ex := NewNegativeArraySizeException("nbits < 0: -5")
	if ex.Error() != "java.lang.NegativeArraySizeException: nbits < 0: -5" || ex.GetMessage() != "nbits < 0: -5" {
		t.Fatalf("unexpected NegativeArraySizeException: err=%q msg=%q", ex.Error(), ex.GetMessage())
	}
}

func TestWave29AUnpinSubstringShortDedup(t *testing.T) {
	// 1. 1-byte ASCII fast-path: 0 hashes, 0 atomics, 0 allocs across all 128 ASCII code points.
	bigASCII := strings.Repeat("padding_", 200) // 1600 bytes >= 1024
	var asciiBuf strings.Builder
	asciiBuf.WriteString(bigASCII)
	for c := 0; c < 128; c++ {
		asciiBuf.WriteByte(byte(c))
	}
	srcASCII := asciiBuf.String()
	baseOff := len(bigASCII)

	ascii1Allocs := testing.AllocsPerRun(20, func() {
		for c := 0; c < 128; c++ {
			got := StringSubstring(srcASCII, baseOff+c, baseOff+c+1)
			if len(got) != 1 || got[0] != byte(c) {
				t.Fatalf("1-byte ASCII mismatch at %d: %q", c, got)
			}
			gotChar := StringFromChar(uint16(c))
			if gotChar != got {
				t.Fatalf("StringFromChar(%d) = %q, want %q", c, gotChar, got)
			}
		}
	})
	if ascii1Allocs != 0 {
		t.Fatalf("expected 0 allocs/op on 1-byte ASCII fast path, got %v", ascii1Allocs)
	}

	// Also verify 1-byte non-ASCII raw byte (>= 0x80) in unpinSubstring does not panic or corrupt.
	rawInvalid := bigASCII + "\x80" + bigASCII
	if got := unpinSubstring(len(rawInvalid), rawInvalid[len(bigASCII):len(bigASCII)+1]); got != "\x80" {
		t.Fatalf("1-byte >=0x80 unpinSubstring = %x, want 80", got)
	}

	// 2. Test all short lengths 0..120 (covering 0, 1, 2..3, 4..7, 8..16, 17..24, 25..96, 97..120).
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_.$"
	var patBuilder strings.Builder
	for i := 0; i < 200; i++ {
		patBuilder.WriteByte(alphabet[(i*13+7)%len(alphabet)])
	}
	pattern := patBuilder.String()
	src1 := strings.Repeat("/*", 700) + pattern + strings.Repeat("*/", 700)
	src2 := strings.Repeat("//", 800) + pattern + strings.Repeat("\n ", 800)
	off1 := 1400
	off2 := 1600

	for length := 0; length <= 120; length++ {
		want := pattern[:length]
		g1 := StringSubstring(src1, off1, off1+length)
		g2 := StringSubstring(src2, off2, off2+length)
		if g1 != want || g2 != want {
			t.Fatalf("len=%d mismatch: g1=%q g2=%q want=%q", length, g1, g2, want)
		}
		if length <= unpinDedupMaxLen {
			allocs := testing.AllocsPerRun(50, func() {
				if got := StringSubstring(src2, off2, off2+length); got != want {
					t.Fatal("dedup hit mismatch")
				}
			})
			if allocs != 0 {
				t.Fatalf("len=%d (<= %d): expected 0 allocs on warm dedup hit, got %v", length, unpinDedupMaxLen, allocs)
			}
		}
	}

	// 3. Overflow / collision test: > 45,000 unique short tokens (> 32,768 table capacity)
	// must never return a mismatched string and must remain bounded.
	prefix := strings.Repeat("x", 1200)
	for i := 0; i < 45000; i++ {
		tok := fmt.Sprintf("gapi.client.service.v1.ResourceModel_%06d", i)
		buf := prefix + tok
		got := unpinSubstring(len(buf), buf[len(prefix):])
		if got != tok {
			t.Fatalf("overflow/collision corruption at %d: got %q, want %q", i, got, tok)
		}
	}

	// 4. Multi-goroutine -race stress test across 32 goroutines.
	hotTokens := []string{
		"this", "prototype", "string", "number", "boolean", "undefined",
		"Object", "Array", "Function", "goog", "opt_data", "disposeInternal",
		"com.google.javascript.jscomp.Compiler",
		"module$exports$google3$javascript$wiz$core$Component_render_v29a",
		"変数名_日本語_🚀",
	}
	sources := make([]string, len(hotTokens))
	offsets := make([]int, len(hotTokens))
	uLengths := make([]int, len(hotTokens))
	for i, tok := range hotTokens {
		pad := fmt.Sprintf("/* file_%02d */%s", i, strings.Repeat(" ", 1200))
		sources[i] = pad + tok + strings.Repeat(" ", 400)
		offsets[i] = len(pad)
		uLengths[i] = StringLength(tok)
	}

	var wg sync.WaitGroup
	const workers = 32
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		w := w
		go func() {
			defer wg.Done()
			localPad := strings.Repeat("p", 1100)
			for iter := 0; iter < 600; iter++ {
				idx := (w*7 + iter) % len(hotTokens)
				got := StringSubstring(sources[idx], offsets[idx], offsets[idx]+uLengths[idx])
				if got != hotTokens[idx] {
					t.Errorf("worker %d hot token mismatch: got %q, want %q", w, got, hotTokens[idx])
					return
				}
				// Interleave cold/evicting misses on shared shards
				if iter%8 == 0 {
					cold := fmt.Sprintf("cold_w%02d_i%04d_ident", w, iter)
					full := localPad + cold
					if gotCold := StringSubstring(full, len(localPad), len(full)); gotCold != cold {
						t.Errorf("worker %d cold token mismatch: got %q, want %q", w, gotCold, cold)
						return
					}
				}
				// Interleave 1-byte ASCII fast path
				ch := byte('a' + (iter % 26))
				if got1 := StringFromChar(uint16(ch)); len(got1) != 1 || got1[0] != ch {
					t.Errorf("StringFromChar(%c) = %q", ch, got1)
					return
				}
			}
		}()
	}
	wg.Wait()
}

func BenchmarkWave29AUnpinSubstring(b *testing.B) {
	pad := strings.Repeat(" ", 2048)
	src1 := pad + ";" + pad
	src8 := pad + "opt_data" + pad
	src48 := pad + "gapi.client.compute.v1.AutoscalingPolicy.Builder" + pad
	off := len(pad)

	// Warm deduplication table
	_ = StringSubstring(src1, off, off+1)
	_ = StringSubstring(src8, off, off+8)
	_ = StringSubstring(src48, off, off+48)

	b.Run("ASCII1Hit", func(b *testing.B) {
		b.ReportAllocs()
		var sink int
		for i := 0; i < b.N; i++ {
			sink += len(StringSubstring(src1, off, off+1))
		}
		_ = sink
	})
	b.Run("Ident8BHit", func(b *testing.B) {
		b.ReportAllocs()
		var sink int
		for i := 0; i < b.N; i++ {
			sink += len(StringSubstring(src8, off, off+8))
		}
		_ = sink
	})
	b.Run("Qual48BHit", func(b *testing.B) {
		b.ReportAllocs()
		var sink int
		for i := 0; i < b.N; i++ {
			sink += len(StringSubstring(src48, off, off+48))
		}
		_ = sink
	})
}

func TestWave32CShortASCIICharAtAndCodePointAt(t *testing.T) {
	jsdoc38 := "/** @param {string} userIdentifier */"
	jsdoc260 := "/**\n * @param {!gapi.client.compute.Instance} instance\n * @param {string=} opt_projectId\n * @return {!Promise<!gapi.client.Request<!Object>>}\n * @suppress {checkTypes,extraRequire}\n */\n" + strings.Repeat(" * detail line\n", 5)
	mixed42 := strings.Repeat("a", 36) + "β🚀z"
	mixed220 := strings.Repeat("b", 180) + "©漢字😀" + strings.Repeat("c", 30)

	for _, s := range []string{jsdoc38, jsdoc260, mixed42, mixed220} {
		u := StringToUTF16(s)
		if got := StringLength(s); got != len(u) {
			t.Fatalf("StringLength(%dB) = %d, want %d", len(s), got, len(u))
		}
		for i := 0; i < len(u); i++ {
			if got := StringCharAt(s, i); got != u[i] {
				t.Fatalf("StringCharAt(%dB, %d) = 0x%X, want 0x%X", len(s), i, got, u[i])
			}
		}
		for i := 0; i < len(u); i++ {
			cp := StringCodePointAt(s, i)
			if u[i] >= 0xD800 && u[i] <= 0xDBFF && i+1 < len(u) && u[i+1] >= 0xDC00 && u[i+1] <= 0xDFFF {
				want := int32(utf16.DecodeRune(rune(u[i]), rune(u[i+1])))
				if cp != want {
					t.Fatalf("StringCodePointAt(%dB, %d) = 0x%X, want 0x%X", len(s), i, cp, want)
				}
				i++
			} else if cp != int32(u[i]) {
				t.Fatalf("StringCodePointAt(%dB, %d) = 0x%X, want 0x%X", len(s), i, cp, u[i])
			}
		}
	}

	// Cold & warm full scan of 40..400B pure-ASCII JSDoc comments must allocate 0 objects.
	comments := make([]string, 128)
	for i := range comments {
		comments[i] = fmt.Sprintf("/** @typedef {{field_%03d: string, count_%03d: number, active: boolean}} */", i, i) + strings.Repeat(" ", i%150)
	}
	allocs := testing.AllocsPerRun(1, func() {
		var sum int32
		for _, c := range comments {
			n := StringLength(c)
			for i := 0; i < n; i++ {
				sum += int32(StringCharAt(c, i)) + StringCodePointAt(c, i)
			}
		}
		if sum == 0 {
			t.Fatal("unexpected zero sum")
		}
	})
	if allocs != 0 {
		t.Fatalf("expected 0 allocs scanning 128 ASCII JSDoc strings, got %v", allocs)
	}

	// Concurrent -race stress across workers scanning shared & distinct short ASCII strings.
	var wg sync.WaitGroup
	for w := 0; w < 8; w++ {
		w := w
		wg.Add(1)
		go func() {
			defer wg.Done()
			for iter := 0; iter < 100; iter++ {
				s := comments[(w*17+iter)%len(comments)]
				for i := 0; i < len(s); i++ {
					if StringCharAt(s, i) != uint16(s[i]) || StringCodePointAt(s, i) != int32(s[i]) {
						t.Errorf("worker %d mismatch at %d", w, i)
						return
					}
				}
			}
		}()
	}
	wg.Wait()

	// Verify >=16 KB JS source file scanning >8,192 characters allocates 0 objects
	// (no periodic rebindStrMeta heap allocations or atomic tick RMWs).
	src32K := strings.Repeat("goog.provide('gapi.client.compute.v1.Instance');\n", 700) // ~34.3 KB
	_ = StringLength(src32K)
	scanAllocs := testing.AllocsPerRun(5, func() {
		n := StringLength(src32K)
		var acc uint32
		for i := 0; i < n; i++ {
			acc += uint32(StringCharAt(src32K, i))
		}
		if acc == 0 {
			t.Fatal("unexpected zero acc")
		}
	})
	if scanAllocs != 0 {
		t.Fatalf("expected 0 allocs scanning 34 KB ASCII JS file, got %v", scanAllocs)
	}
}

func BenchmarkWave34CJString(b *testing.B) {
	src32K := "// Bench 34C\n" + strings.Repeat("goog.provide('gapi.client.compute.v1.Instance');\n", 700)
	jsdoc220 := "/**\n * @param {!gapi.client.compute.Instance} instance\n * @param {string=} opt_projectId\n * @return {!Promise<!gapi.client.Request<!Object>>}\n * @suppress {checkTypes,extraRequire}\n */\n"
	short24 := "gapi.client.compute.v1.X"

	// Warm caches (including SourceFile BOM check on src32K)
	_ = StringStartsWith(src32K, "\ufeff")
	_ = StringLength(src32K)
	for i := 0; i < len(jsdoc220); i++ {
		_ = StringCharAt(jsdoc220, i)
	}

	b.Run("CharAt_32KB_HotFile", func(b *testing.B) {
		b.ReportAllocs()
		var sink uint16
		for i := 0; i < b.N; i++ {
			sink += StringCharAt(src32K, i&16383)
		}
		_ = sink
	})
	b.Run("CharAt_220B_JSDoc", func(b *testing.B) {
		b.ReportAllocs()
		var sink uint16
		for i := 0; i < b.N; i++ {
			sink += StringCharAt(jsdoc220, 40+(i&127))
		}
		_ = sink
	})
	b.Run("CharAt_24B_ShortASCII", func(b *testing.B) {
		b.ReportAllocs()
		var sink uint16
		for i := 0; i < b.N; i++ {
			sink += StringCharAt(short24, i&15)
		}
		_ = sink
	})
	b.Run("Length_32KB_HotFile", func(b *testing.B) {
		b.ReportAllocs()
		var sink int
		for i := 0; i < b.N; i++ {
			sink += StringLength(src32K)
		}
		_ = sink
	})
	b.Run("Length_220B_JSDoc", func(b *testing.B) {
		b.ReportAllocs()
		var sink int
		for i := 0; i < b.N; i++ {
			sink += StringLength(jsdoc220)
		}
		_ = sink
	})
}

func TestWave36CZeroAllocColdShortASCIIScan(t *testing.T) {
	// Generate 6,000 distinct 70..320B ASCII strings (exceeding the 4,096-entry cache capacity
	// so thousands of misses & evictions occur!) and verify ZERO heap allocations on misses!
	pool := make([]string, 6000)
	for i := range pool {
		pool[i] = fmt.Sprintf("module$exports$google3$javascript$closure$w36c$Service_%05d$Method_", i) + strings.Repeat("x", (i*7)%240)
	}
	idx := 0
	allocs := testing.AllocsPerRun(200, func() {
		// Scan 30 strings per iteration across 200 iterations = 6,000 cold/evicting strings!
		var sum int
		for k := 0; k < 30; k++ {
			s := pool[idx%len(pool)]
			idx++
			n := StringLength(s)
			if !IsASCII(s) || n != len(s) {
				t.Fatal("IsASCII / StringLength mismatch")
			}
			for j := 0; j < n; j += 17 {
				sum += int(StringCharAt(s, j)) + int(StringCodePointAt(s, j))
			}
		}
		if sum == 0 {
			t.Fatal("unexpected zero sum")
		}
	})
	if allocs != 0 {
		t.Fatalf("expected 0 heap allocations on cold/evicting >64B ASCII misses, got %v", allocs)
	}
}

func TestWave40CShardedZeroBucketAllocUnpinDedup256(t *testing.T) {
	// 1. Verify long Google Cloud protobuf qualified names (97..256 bytes, previously > old 96B cap!)
	// deduplicate with 0 allocations on warm hits, while >256B still clones cleanly.
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_.$"
	var pb strings.Builder
	for i := 0; i < 320; i++ {
		pb.WriteByte(alphabet[(i*17+11)%len(alphabet)])
	}
	longPattern := pb.String()
	srcA := strings.Repeat("/*", 900) + longPattern + strings.Repeat("*/", 900)
	srcB := strings.Repeat("//", 1000) + longPattern + strings.Repeat("\n ", 1000)
	offA, offB := 1800, 2000

	for _, length := range []int{97, 128, 180, 224, 255, 256, 257, 300} {
		want := longPattern[:length]
		gA := StringSubstring(srcA, offA, offA+length)
		gB := StringSubstring(srcB, offB, offB+length)
		if gA != want || gB != want {
			t.Fatalf("len=%d mismatch", length)
		}
		if length <= unpinDedupMaxLen {
			allocs := testing.AllocsPerRun(50, func() {
				if got := StringSubstring(srcB, offB, offB+length); got != want {
					t.Fatal("dedup hit mismatch")
				}
			})
			if allocs != 0 {
				t.Fatalf("len=%d (<= %d): expected 0 allocs on warm dedup hit, got %v", length, unpinDedupMaxLen, allocs)
			}
		}
	}

	// 2. Verify cold misses into unpinShards do ONLY 1 heap alloc (strings.Clone) and ZERO unpinBucket struct allocs!
	prefix := strings.Repeat("P", 1500)
	coldBufs := make([]string, 32)
	for i := range coldBufs {
		coldBufs[i] = prefix + fmt.Sprintf("proto.google.cloud.dialogflow.cx.v3beta1.UniqueColdSymbolW40C_%06d_ServiceClientImpl", i)
	}
	missIdx := 0
	missAllocs := testing.AllocsPerRun(25, func() {
		buf := coldBufs[missIdx]
		missIdx++
		sub := buf[len(prefix):]
		if got := unpinSubstring(len(buf), sub); got != sub {
			t.Fatalf("cold miss mismatch: got %q, want %q", got, sub)
		}
	})
	// Exactly 1 strings.Clone per unique miss (previously 2.0 allocs/miss due to &unpinBucket{}!).
	if missAllocs > 1 {
		t.Fatalf("expected <= 1 alloc per cold unpinDedup miss (no bucket struct), got %v", missAllocs)
	}

	// 3. Concurrent -race stress across 24 workers on 100..256B substrings (hits + misses).
	var wg sync.WaitGroup
	const workers = 24
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		w := w
		go func() {
			defer wg.Done()
			pad := strings.Repeat("z", 1300)
			for i := 0; i < 400; i++ {
				var tok string
				if i&1 == 0 {
					tok = longPattern[:120+(i%137)] // 120..256B hot shared slice
					got := StringSubstring(srcA, offA, offA+len(tok))
					if got != tok {
						t.Errorf("worker %d hot >96B mismatch", w)
						return
					}
				} else {
					tok = fmt.Sprintf("proto.google.cloud.compute.v1.Worker%02d_Iter%04d_RegionalInstanceGroupManager", w%6, i)
					buf := pad + tok
					if got := unpinSubstring(len(buf), buf[len(pad):]); got != tok {
						t.Errorf("worker %d miss >96B mismatch", w)
						return
					}
				}
			}
		}()
	}
	wg.Wait()
}


func TestWave42C4WaySLRUAndLargeCommentUnpinDedup(t *testing.T) {
	// 1. Repeated 600-byte Apache license header & JSDoc block comments (257..2048B) across
	// >= 1 KB source files must admit on 2nd occurrence and allocate 0 objects on 3rd+ occurrences!
	license := "/*\n * Copyright 2026 Google LLC\n * Licensed under the Apache License, Version 2.0 (the \"License\");\n * " +
		strings.Repeat("you may not use this file except in compliance with the License. ", 8) + "\n */"
	if len(license) <= 256 || len(license) > 2048 {
		t.Fatalf("expected 257..2048B license comment, got %d", len(license))
	}
	file1 := license + "\n" + strings.Repeat("export const x = 1;\n", 100)
	file2 := license + "\n" + strings.Repeat("export const y = 2;\n", 100)
	file3 := license + "\n" + strings.Repeat("export const z = 3;\n", 100)

	if g1 := StringSubstring(file1, 0, len(license)); g1 != license {
		t.Fatal("1st occurrence mismatch")
	}
	if g2 := StringSubstring(file2, 0, len(license)); g2 != license {
		t.Fatal("2nd occurrence mismatch")
	}
	allocs := testing.AllocsPerRun(50, func() {
		if g3 := StringSubstring(file3, 0, len(license)); g3 != license {
			t.Fatal("3rd+ occurrence mismatch")
		}
	})
	if allocs != 0 {
		t.Fatalf("expected 0 allocs on repeated 600B license comment after 2-hit admission, got %v", allocs)
	}

	// 2. Verify 4-way SLRU probation protection: flooding 50,000 one-hit unique comments must NOT
	// evict hot protected identifiers that have already been promoted to slot[0..2]!
	pad := strings.Repeat(" ", 1500)
	hotSrc := pad + "goog.module.declareLegacyNamespace"
	// Hit twice to ensure promotion into protected ways slot[0..2].
	_ = unpinSubstring(len(hotSrc), hotSrc[len(pad):])
	_ = unpinSubstring(len(hotSrc), hotSrc[len(pad):])
	for i := 0; i < 50000; i++ {
		oneOff := pad + fmt.Sprintf("/* one_off_unique_comment_%06d */", i)
		_ = unpinSubstring(len(oneOff), oneOff[len(pad):])
	}
	if hitAllocs := testing.AllocsPerRun(20, func() {
		if got := unpinSubstring(len(hotSrc), hotSrc[len(pad):]); got != "goog.module.declareLegacyNamespace" {
			t.Fatal("hot protected identifier mismatch")
		}
	}); hitAllocs != 0 {
		t.Fatalf("hot protected identifier was evicted by one-hit wonders: got %v allocs", hitAllocs)
	}
}
