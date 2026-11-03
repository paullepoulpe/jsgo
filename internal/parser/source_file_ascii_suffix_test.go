package parser

import (
	"fmt"
	"math/rand/v2"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// naiveUTF16SourceFile is an unoptimized reference implementation backed by a
// single full jutil.StringToUTF16(code) slice for differential testing.
type naiveUTF16SourceFile struct {
	code       string
	units      []uint16
	lineStarts []int32
}

func newNaiveUTF16SourceFile(code string) *naiveUTF16SourceFile {
	units := jutil.StringToUTF16(code)
	starts := []int32{0}
	for i := 0; i < len(units); i++ {
		switch units[i] {
		case '\n', 0x2028, 0x2029:
			starts = append(starts, int32(i+1))
		case '\r':
			if i+1 < len(units) && units[i+1] == '\n' {
				i++
			}
			starts = append(starts, int32(i+1))
		}
	}
	return &naiveUTF16SourceFile{
		code:       code,
		units:      units,
		lineStarts: starts,
	}
}

func (r *naiveUTF16SourceFile) length() int {
	return len(r.units)
}

func (r *naiveUTF16SourceFile) charAt(index int) uint16 {
	return r.units[index]
}

func (r *naiveUTF16SourceFile) substring(start, end int) string {
	return jutil.StringFromUTF16(r.units[start:end])
}

func (r *naiveUTF16SourceFile) lineColumn(offset int) (int, int) {
	if offset < 0 {
		return 0, offset
	}
	line := FindLine(r.lineStarts, int32(offset))
	return line, offset - int(r.lineStarts[line])
}

func verifySourceFileAgainstNaive(t *testing.T, label string, sf *SourceFile, code string) {
	t.Helper()
	ref := newNaiveUTF16SourceFile(code)

	if got, want := sf.Length(), ref.length(); got != want {
		t.Fatalf("%s: Length() = %d, want %d", label, got, want)
	}
	n := sf.Length()
	for i := range n {
		if got, want := sf.CharAt(i), ref.charAt(i); got != want {
			t.Fatalf("%s: CharAt(%d) = %#x, want %#x", label, i, got, want)
		}
	}

	if gotStarts := sf.LineStarts(); !slices.Equal(gotStarts, ref.lineStarts) {
		t.Fatalf("%s: LineStarts() = %v, want %v", label, gotStarts, ref.lineStarts)
	}

	// Forward, backward, and out-of-bounds LineColumn / LinePosition / Line checks.
	for off := -2; off <= n+2; off++ {
		gotL, gotC := sf.LineColumn(off)
		wantL, wantC := ref.lineColumn(off)
		if gotL != wantL || gotC != wantC {
			t.Fatalf("%s: forward LineColumn(%d) = (%d, %d), want (%d, %d)", label, off, gotL, gotC, wantL, wantC)
		}
		if gotLine := sf.Line(NewSourcePosition(off)); gotLine != wantL {
			t.Fatalf("%s: Line(%d) = %d, want %d", label, off, gotLine, wantL)
		}
	}
	for off := n + 2; off >= -2; off-- {
		gotL, gotC := sf.LineColumn(off)
		wantL, wantC := ref.lineColumn(off)
		if gotL != wantL || gotC != wantC {
			t.Fatalf("%s: backward LineColumn(%d) = (%d, %d), want (%d, %d)", label, off, gotL, gotC, wantL, wantC)
		}
	}

	// Exhaustive Substring / ASCIISlice / UTF16Slice checks when n <= 80, sampled otherwise.
	checkSpan := func(start, end int) {
		wantSub := ref.substring(start, end)
		gotSub := sf.Substring(start, end)
		if gotSub != wantSub {
			t.Fatalf("%s: Substring(%d, %d) = %q, want %q", label, start, end, gotSub, wantSub)
		}
		if asciiSub, ok := sf.ASCIISlice(start, end); ok {
			if asciiSub != wantSub {
				t.Fatalf("%s: ASCIISlice(%d, %d) = %q, want %q", label, start, end, asciiSub, wantSub)
			}
			for k := range len(asciiSub) {
				if asciiSub[k] >= 0x80 {
					t.Fatalf("%s: ASCIISlice(%d, %d) returned non-ASCII byte %#x", label, start, end, asciiSub[k])
				}
			}
		}
		if u16, ok := sf.UTF16Slice(start, end); ok {
			if !slices.Equal(u16, ref.units[start:end]) {
				t.Fatalf("%s: UTF16Slice(%d, %d) = %v, want %v", label, start, end, u16, ref.units[start:end])
			}
		}
	}
	if n <= 80 {
		for start := 0; start <= n; start++ {
			for end := start; end <= n; end++ {
				checkSpan(start, end)
			}
		}
	} else {
		anchors := []int{
			0, 1,
			max(0, sf.asciiPrefix-2), max(0, sf.asciiPrefix-1), sf.asciiPrefix, min(n, sf.asciiPrefix+1),
			max(0, sf.asciiSuffixUTF16Start-2), max(0, sf.asciiSuffixUTF16Start-1), sf.asciiSuffixUTF16Start, min(n, sf.asciiSuffixUTF16Start+1),
			max(0, n-2), max(0, n-1), n,
		}
		for _, a := range anchors {
			for _, b := range anchors {
				if a <= b && a >= 0 && b <= n {
					checkSpan(a, b)
				}
			}
		}
	}
}

func TestSourceFileASCIISuffixAvoidsTailUTF16Allocation(t *testing.T) {
	header := "// Copyright © 2026 Example Corp.\n"
	body := strings.Repeat("const value_123 = 'hello world';\n", 6000) // ~198 KB pure ASCII
	src := header + body

	runtime.GC()
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	sf := NewSourceFile("header_copyright.js", src)
	runtime.ReadMemStats(&after)

	if len(sf.utf16Mid) != 1 {
		t.Fatalf("len(sf.utf16Mid) = %d, want 1 (only '©')", len(sf.utf16Mid))
	}
	// Ensure ASCIISlice succeeds zero-copy in the trailing suffix!
	suffixStart := sf.asciiSuffixUTF16Start + 10
	slice, ok := sf.ASCIISlice(suffixStart, suffixStart+15)
	if !ok || slice != "value_123 = 'he"[:15] {
		// Verify exact content via Substring parity.
		if !ok || slice != sf.Substring(suffixStart, suffixStart+15) {
			t.Fatalf("ASCIISlice in trailing suffix = (%q, %v), want (%q, true)", slice, ok, sf.Substring(suffixStart, suffixStart+15))
		}
	}
	allocBytes := after.TotalAlloc - before.TotalAlloc
	if allocBytes > 16<<10 {
		t.Fatalf("NewSourceFile allocated %d bytes for 198 KB trailing-ASCII file; expected < 16 KiB", allocBytes)
	}
	verifySourceFileAgainstNaive(t, "copyright-header-large-ascii-tail", sf, src)
}

func TestSourceFileDifferentialEdgeCasesAndArtificialCRLFSplits(t *testing.T) {
	loneHighWTF8 := jutil.StringFromUTF16([]uint16{0xD800})
	loneLowWTF8 := jutil.StringFromUTF16([]uint16{0xDC00})

	cases := []struct {
		name string
		code string
	}{
		{"empty", ""},
		{"pure_ascii_lf", "let a = 1;\nlet b = 2;\n"},
		{"pure_ascii_crlf_cr", "a\r\nb\rc\n\r\nd\r"},
		{"copyright_then_ascii", "// © 2026\r\nconst x = 1;\r\nconst y = 2;\n"},
		{"emoji_surrogate_pair_mid", "prefix\nconst e = '😀🚀';\nsuffix_ident = 42;\r\n"},
		{"non_ascii_at_very_end", "const a = 1;\nconst b = 'λ'"},
		{"non_ascii_at_start_and_end", "/* α */ const mid = 123; /* ω */"},
		{"u2028_u2029_terminators", "a\u2028b\u2029c\r\nd\re\nf"},
		{"lone_high_surrogate_before_ascii", "pre" + loneHighWTF8 + "post_ascii\r\nline2"},
		{"lone_low_surrogate_before_ascii", "pre" + loneLowWTF8 + "\npost_ascii"},
		{"unpaired_high_then_low_separated_by_ascii", loneHighWTF8 + "x" + loneLowWTF8 + "tail_ascii"},
		{"adjacent_surrogates_and_crlf", "\r\n" + loneHighWTF8 + "\r\n" + "ascii_end\r\n"},
		{"invalid_raw_utf8_bytes", "abc\xff\xfedef_trailing_ascii\r\n123"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sf := NewSourceFile(tc.name+".js", tc.code)
			verifySourceFileAgainstNaive(t, tc.name, sf, tc.code)
		})
	}

	// White-box test: artificially force region splits directly between '\r' and '\n'
	// at BOTH asciiPrefix and asciiSuffixUTF16Start to prove CRLF boundary handling.
	rawCRLF := "ab\r\ncd\r\nef\r\ngh"
	fullUnits := jutil.StringToUTF16(rawCRLF)
	for p := 0; p <= len(rawCRLF); p++ {
		for s := p; s <= len(rawCRLF); s++ {
			artificial := &SourceFile{
				Name:                  "artificial_crlf.js",
				Contents:              rawCRLF,
				asciiPrefix:           p,
				asciiSuffixUTF16Start: s,
				asciiSuffixByteStart:  s,
				utf16Len:              len(rawCRLF),
				utf16Mid:              slices.Clone(fullUnits[p:s]),
			}
			verifySourceFileAgainstNaive(t, fmt.Sprintf("artificial_split_%d_%d", p, s), artificial, rawCRLF)
		}
	}
}

type recordedTok struct {
	typ   TokenType
	start int32
	end   int32
	text  string
}

type recordedComment struct {
	typ   Comment_Type
	rng   LineRange
	value string
}

func scanAllTokensAndComments(code string, sf *SourceFile) ([]recordedTok, []recordedComment, []string) {
	rep := &recordingErrorReporter{}
	rec := NewParser_CommentRecorder(&Parser{})
	sc := NewScanner(rep, rec, sf, 0)
	var toks []recordedTok
	for {
		tok := sc.NextToken()
		toks = append(toks, recordedTok{
			typ:   tok.Type(),
			start: tok.Location.Start.Offset,
			end:   tok.Location.End.Offset,
			text:  sf.Substring(int(tok.Location.Start.Offset), int(tok.Location.End.Offset)),
		})
		if tok.Type() == TokenType_END_OF_FILE {
			break
		}
	}
	var comments []recordedComment
	for _, c := range rec.GetComments() {
		comments = append(comments, recordedComment{
			typ:   c.Type,
			rng:   c.Location,
			value: c.Value,
		})
	}
	return toks, comments, rep.lines
}

func makeForceUTF16SourceFile(name, code string) *SourceFile {
	units := jutil.StringToUTF16(code)
	return &SourceFile{
		Name:                  name,
		Contents:              code,
		asciiPrefix:           0,
		asciiSuffixUTF16Start: len(units),
		asciiSuffixByteStart:  len(code),
		utf16Len:              len(units),
		utf16Mid:              units,
	}
}

func TestScannerParityAcrossPrefixMidAndTrailingASCIISuffix(t *testing.T) {
	programs := []string{
		// Copyright © in header followed by comments, strings, templates, numbers, keywords, private idents, block comments.
		"// Copyright © 2026\r\n" +
			"/* block comment with * stars * inside */\n" +
			"/** @jsdoc comment */\n" +
			"/*! important comment */\n" +
			"class Foo {\n" +
			"  #secretField = 123_456.789e2;\n" +
			"  method(alpha, beta_long_identifier_name_exceeding_sixty_four_characters_0123456789abcdef) {\n" +
			"    const s1 = 'single \\n escaped quote \\' end';\n" +
			"    const s2 = \"double quote string\";\n" +
			"    const tpl = `head ${alpha + 0xdeadbeef} mid ${`nested`} tail`;\n" +
			"    return this.#secretField + s1.length + s2.length + tpl.length; // trailing line comment\n" +
			"  }\n" +
			"}/*eof*/",
		// Tokens straddling the middle non-ASCII span into the trailing ASCII suffix!
		"const caf\u00e9_suffixIdent = 'str_before_\u00a9_after_ascii';\n" +
			"const \u03b1_startNonASCII = `tpl_\U0001f680_${caf\u00e9_suffixIdent}_ascii_tail`;\n" +
			"/* comment \u2603 snow * still comment */ let afterSnow = 987654321;\r\n" +
			"// single line with \u00e9 and trailing ascii",
	}

	for idx, src := range programs {
		optSF := NewSourceFile(fmt.Sprintf("prog_%d.js", idx), src)
		refSF := makeForceUTF16SourceFile(fmt.Sprintf("prog_%d.js", idx), src)

		gotToks, gotComments, gotErrs := scanAllTokensAndComments(src, optSF)
		wantToks, wantComments, wantErrs := scanAllTokensAndComments(src, refSF)

		if !slices.Equal(gotToks, wantToks) {
			t.Fatalf("program %d token mismatch:\n got: %+v\nwant: %+v", idx, gotToks, wantToks)
		}
		if !slices.Equal(gotComments, wantComments) {
			t.Fatalf("program %d comment mismatch:\n got: %+v\nwant: %+v", idx, gotComments, wantComments)
		}
		if !slices.Equal(gotErrs, wantErrs) {
			t.Fatalf("program %d error mismatch:\n got: %v\nwant: %v", idx, gotErrs, wantErrs)
		}
	}
}

func TestSourceFileAndScannerRandomDifferentialProperty(t *testing.T) {
	fragments := []string{
		"a", "bc", "ident_123", "function", "return", " ", "\t",
		"\n", "\r", "\r\n", "\u2028", "\u2029",
		"'str'", "\"esc\\n\"", "`tpl${1}end`",
		"//cmt\n", "/*blk*end*/", "42", "3.14",
		"©", "é", "λ", "☃", "😀", "🚀",
		jutil.StringFromUTF16([]uint16{0xD800}),
		jutil.StringFromUTF16([]uint16{0xDC00}),
	}

	rng := rand.New(rand.NewPCG(93, 2026))
	for iter := range 250 {
		var sb strings.Builder
		nParts := rng.IntN(24)
		for range nParts {
			sb.WriteString(fragments[rng.IntN(len(fragments))])
		}
		// Frequently append a pure-ASCII tail of varying lengths (0..90 bytes, crossing 8-byte SWAR chunks).
		if rng.IntN(2) == 0 {
			sb.WriteString(strings.Repeat("x = 1;\n", rng.IntN(14)))
		}
		code := sb.String()

		sf := NewSourceFile("prop.js", code)
		verifySourceFileAgainstNaive(t, fmt.Sprintf("iter_%d", iter), sf, code)

		refSF := makeForceUTF16SourceFile("prop.js", code)
		gotToks, gotComments, gotErrs := scanAllTokensAndComments(code, sf)
		wantToks, wantComments, wantErrs := scanAllTokensAndComments(code, refSF)
		if !slices.Equal(gotToks, wantToks) || !slices.Equal(gotComments, wantComments) || !slices.Equal(gotErrs, wantErrs) {
			t.Fatalf("iter %d (%q): scanner output diverged from full-UTF16 reference", iter, code)
		}
	}
}
