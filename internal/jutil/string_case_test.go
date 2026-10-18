package jutil

import (
	unicode "github.com/paullepoulpe/jsgo/internal/unicode15"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// Vectors from tools/oracle/java/JutilStringCaseVectors.java (JDK 21): the
// JDK's Character.toLowerCase(int) and String.toLowerCase(Locale.ROOT).

type stringCaseVectors struct {
	lower  map[int]int              // cp rows
	single map[rune][]uint16        // str1 rows
	sigma  map[string][]bool        // sigma rows, per context
	random map[string][][2][]uint16 // mixed and text rows
}

func loadStringCaseVectors(t *testing.T) *stringCaseVectors {
	t.Helper()
	v := &stringCaseVectors{
		lower:  map[int]int{},
		single: map[rune][]uint16{},
		sigma:  map[string][]bool{},
		random: map[string][][2][]uint16{},
	}
	hex := func(s string) int {
		n, err := strconv.ParseInt(s, 16, 64)
		if err != nil {
			t.Fatal(err)
		}
		return int(int32(n))
	}
	for _, r := range readTSV(t, "string_case.tsv") {
		switch r[0] {
		case "cp":
			v.lower[hex(r[1])] = hex(r[2])
		case "str1":
			v.single[rune(hex(r[1]))] = decodeUnits(t, r[2])
		case "sigma":
			var bits []bool
			cur := false
			for _, f := range strings.Split(r[2], ",") {
				n, err := strconv.Atoi(f)
				if err != nil {
					t.Fatal(err)
				}
				for ; n > 0; n-- {
					bits = append(bits, cur)
				}
				cur = !cur
			}
			if len(bits) != unicode.MaxRune+1 {
				t.Fatalf("sigma %s: %d bits", r[1], len(bits))
			}
			v.sigma[r[1]] = bits
		case "mixed", "text":
			v.random[r[0]] = append(v.random[r[0]], [2][]uint16{decodeUnits(t, r[1]), decodeUnits(t, r[2])})
		default:
			t.Fatalf("bad row %q", r)
		}
	}
	if len(v.lower) < 1400 || len(v.sigma) != 4 || len(v.random["mixed"]) < 1000 || len(v.random["text"]) < 1000 {
		t.Fatalf("incomplete vectors: %d cp, %d sigma, %d mixed, %d text",
			len(v.lower), len(v.sigma), len(v.random["mixed"]), len(v.random["text"]))
	}
	return v
}

func codePointString(cp rune) string {
	return StringFromUTF16(appendUTF16CodePoint(nil, cp))
}

func TestCharacterToLowerCaseCodePointMatchesJava(t *testing.T) {
	v := loadStringCaseVectors(t)
	for cp := 0; cp <= unicode.MaxRune; cp++ {
		want, ok := v.lower[cp]
		if !ok {
			want = cp
		}
		if got := CharacterToLowerCaseCodePoint(cp); got != want {
			t.Errorf("Character.toLowerCase(%X) = %X, want %X", cp, got, want)
		}
	}
	for cp, want := range v.lower {
		if cp < 0 || cp > unicode.MaxRune {
			if got := CharacterToLowerCaseCodePoint(cp); got != want {
				t.Errorf("Character.toLowerCase(%X) = %X, want %X", cp, got, want)
			}
		}
	}
}

// Every one-code-point string (lone surrogates included) lowercases as in
// Java: the simple mapping, except U+0130.
func TestStringToLowerCaseRootSingleCodePointsMatchJava(t *testing.T) {
	v := loadStringCaseVectors(t)
	if !slices.Equal(v.single[0x130], []uint16{'i', 0x307}) || len(v.single) != 1 {
		t.Fatalf("unexpected str1 rows: %v", v.single)
	}
	for cp := rune(0); cp <= unicode.MaxRune; cp++ {
		want, ok := v.single[cp]
		if !ok {
			want = appendUTF16CodePoint(nil, rune(CharacterToLowerCaseCodePoint(int(cp))))
		}
		if got := StringToUTF16(StringToLowerCaseRoot(codePointString(cp))); !slices.Equal(got, want) {
			t.Errorf("toLowerCase(U+%04X) = %04x, want %04x", cp, got, want)
		}
	}
}

// sigmaAgreesWithJava is the set of characters for which the Unicode
// Final_Sigma rule that StringToLowerCaseRoot implements is known to give
// Java's result when the character stands right before or after the sigma
// (the four contexts of the sigma rows): BMP letters with case (Lu, Ll,
// Lt), BMP nonspacing and enclosing marks and format characters that are
// not themselves cased, the apostrophe and the full stop, and the
// characters that are neither cased nor case-ignorable except, in the BMP,
// other letters (Lo), spacing marks (Mc), numbers (Nd, Nl, No), connector
// punctuation (Pc), dashes (Pd) and the quotation mark.
func sigmaAgreesWithJava(cp rune) bool {
	if cp == '\'' || cp == '.' {
		return true
	}
	cased, ignorable := isCased(cp), isCaseIgnorable(cp)
	if cp > 0xFFFF {
		return !cased && !ignorable
	}
	switch {
	case cased:
		return unicode.In(cp, unicode.Lu, unicode.Ll, unicode.Lt)
	case ignorable:
		return unicode.In(cp, unicode.Mn, unicode.Me, unicode.Cf)
	}
	return cp != '"' && !unicode.In(cp, unicode.Lo, unicode.Mc, unicode.Nd, unicode.Nl, unicode.No, unicode.Pc, unicode.Pd)
}

// sigmaAgreesWithJavaInSequence is the subset of sigmaAgreesWithJava on
// which the two rules also agree for any sequence of such characters around
// a sigma. The apostrophe, the full stop and the format characters agree
// alone but not in runs: for "Ἇ\u0342'.Σ," and "AΣ\u00ad\ua67cἍ" Unicode
// skips the whole run and Java does not.
func sigmaAgreesWithJavaInSequence(cp rune) bool {
	return cp != '\'' && cp != '.' && !unicode.Is(unicode.Cf, cp) && sigmaAgreesWithJava(cp)
}

// The sigma contexts, over every code point c. StringToLowerCaseRoot decides
// final sigma by the Unicode definition and Java by its own word-boundary
// logic (see string_to_lower_case.go); the number of code points on which
// they differ is pinned here, and there must be none inside
// sigmaAgreesWithJava.
func TestStringToLowerCaseRootFinalSigmaContexts(t *testing.T) {
	v := loadStringCaseVectors(t)
	contexts := []struct {
		name      string
		prefix    string
		suffix    string
		divergent int // code points where Go and Java differ
	}{
		{"pre", "A", "Σ", 28720},
		{"pre0", "", "Σ", 328},
		{"post", "AΣ", "a", 27228},
		{"post0", "AΣ", "", 327},
	}
	for _, c := range contexts {
		bits := v.sigma[c.name]
		n := 0
		var examples []string
		for cp := rune(0); cp <= unicode.MaxRune; cp++ {
			in := c.prefix + codePointString(cp) + c.suffix
			out := StringToUTF16(StringToLowerCaseRoot(in))
			i := len(out) - 1 // pre, pre0: the sigma is last
			if strings.HasPrefix(c.prefix, "AΣ") {
				i = 1
			}
			if out[i] != 0x3C2 && out[i] != 0x3C3 {
				t.Fatalf("%s U+%04X: no sigma at %d in %04x", c.name, cp, i, out)
			}
			if (out[i] == 0x3C2) == bits[cp] {
				continue
			}
			n++
			if sigmaAgreesWithJava(cp) {
				t.Errorf("%s: %q lowercases to %q, Java gives the other sigma", c.name, in, StringFromUTF16(out))
			}
			if len(examples) < 6 {
				examples = append(examples, strconv.QuoteToASCII(in))
			}
		}
		if n != c.divergent {
			t.Errorf("%s: %d code points give a different sigma than Java, want %d", c.name, n, c.divergent)
		}
		t.Logf("%s: %d code points differ from Java, e.g. %s", c.name, n, strings.Join(examples, " "))
	}
}

// Documented divergences (Java's results observed black-box on JDK 21); the
// Go results are the Unicode Final_Sigma rule's.
func TestStringToLowerCaseRootFinalSigmaKnownDivergences(t *testing.T) {
	cases := []struct{ in, java, goWant string }{
		// Java joins words across '-', '_', digits and '"': the first sigma is
		// not word-final for Java.
		{"ΟΔΟΣ-ΑΘΗΝΩΝ", "οδοσ-αθηνων", "οδος-αθηνων"},
		{"ΤΕΛΟΣ_ΑΡΧΗ", "τελοσ_αρχη", "τελος_αρχη"},
		{"ΑΣ1Β", "ασ1β", "ας1β"},
		{"Α1Σ", "α1ς", "α1σ"},
		// Java stops at ':' and modifier symbols that Unicode skips.
		{"ΟΔΟΣ:ΑΘΗΝΑ", "οδος:αθηνα", "οδοσ:αθηνα"},
		{"Α^Σ", "α^σ", "α^ς"},
		// Java does not skip supplementary marks, nor count supplementary
		// letters before the sigma.
		{"Α\U00010400Σ", "α\U00010428σ", "α\U00010428ς"},
		// Cased characters that are not letters (circled letters, ª).
		{"ⒶΣ", "ⓐσ", "ⓐς"},
	}
	for _, c := range cases {
		if got := StringToLowerCaseRoot(c.in); got != c.goWant {
			t.Errorf("StringToLowerCaseRoot(%q) = %q, want %q (Java: %q)", c.in, got, c.goWant, c.java)
		}
	}
	// And contexts where both agree.
	for in, want := range map[string]string{
		"ΟΔΟΣ ΑΘΗΝΩΝ": "οδος αθηνων",
		"ΟΔΟΣ, ΟΔΟΣ.": "οδος, οδος.",
		"Σ":           "σ",
		"ΣΑ":          "σα",
		"ΑΣ\u0301":    "ας\u0301",
		"ΑΣ\u0301Β":   "ασ\u0301β",
		"Α'Σ":         "α'ς",
		"ΑΣ'Β":        "ασ'β",
		"ΑΣ.":         "ας.",
		"ΑΣ.Β":        "ασ.β",
		"ΑΣΣ":         "ασς",
		"ΑΣ\u200bΒ":   "ασ\u200bβ",
		"\U00010400Σ": "\U00010428ς",
	} {
		if got := StringToLowerCaseRoot(in); got != want {
			t.Errorf("StringToLowerCaseRoot(%q) = %q, want %q", in, got, want)
		}
	}
}

// The random strings: everything but final sigma matches Java exactly; a
// sigma may differ only in a string holding a character outside
// sigmaAgreesWithJavaInSequence.
func TestStringToLowerCaseRootRandomMatchesJava(t *testing.T) {
	v := loadStringCaseVectors(t)
	isSigma := func(c uint16) bool { return c == 0x3C2 || c == 0x3C3 }
	for _, kind := range []string{"mixed", "text"} {
		sigmas, differ := 0, 0
		for _, row := range v.random[kind] {
			in, want := row[0], row[1]
			got := StringToUTF16(StringToLowerCaseRoot(StringFromUTF16(in)))
			// Where each capital sigma of the input lands in the output
			// (U+0130 grows by one unit), and whether the string holds a
			// character on which the final-sigma rules may disagree.
			sigmaAt := map[int]bool{}
			explained := false
			for i, o := 0, 0; i < len(in); {
				cp, n := utf16CodePointAt(in, i)
				explained = explained || !sigmaAgreesWithJavaInSequence(cp)
				switch cp {
				case 0x3A3:
					sigmaAt[o] = true
					sigmas++
					o++
				case 0x130:
					o += 2
				default:
					o += len(appendUTF16CodePoint(nil, rune(CharacterToLowerCaseCodePoint(int(cp)))))
				}
				i += n
			}
			if len(got) != len(want) {
				t.Errorf("toLowerCase(%04x) = %04x, want %04x", in, got, want)
				continue
			}
			for i := range got {
				switch {
				case got[i] == want[i]:
				case sigmaAt[i] && isSigma(got[i]) && isSigma(want[i]):
					differ++
					if !explained {
						t.Errorf("toLowerCase(%q) = %q, Java: %q", StringFromUTF16(in), StringFromUTF16(got), StringFromUTF16(want))
					}
				default:
					t.Errorf("toLowerCase(%04x) = %04x, want %04x", in, got, want)
				}
			}
		}
		t.Logf("%s: %d of %d capital sigmas lowercase differently from Java", kind, differ, sigmas)
	}
}

func TestStringToLowerCaseRootFastPathsAndSurrogates(t *testing.T) {
	for _, s := range []string{"", "abc", "already lower 123", "é", "ς"} {
		if got := StringToLowerCaseRoot(s); got != s {
			t.Errorf("StringToLowerCaseRoot(%q) = %q", s, got)
		}
	}
	for in, want := range map[string]string{
		"ABC":                  "abc",
		"Hello, World":         "hello, world",
		"ÀÉÎ":                  "àéî",
		"İSTANBUL":             "i\u0307stanbul",
		"\U00010400\U0001E900": "\U00010428\U0001E922",
	} {
		if got := StringToLowerCaseRoot(in); got != want {
			t.Errorf("StringToLowerCaseRoot(%q) = %q, want %q", in, got, want)
		}
	}
	// Lone surrogates are kept; a high and a low surrogate stored as two
	// 3-byte WTF-8 sequences (as Go concatenation can leave them) are one
	// code point, as the two chars are in Java.
	lone := StringFromUTF16([]uint16{'A', 0xD801, 'B', 0xDC00})
	if got := StringToUTF16(StringToLowerCaseRoot(lone)); !slices.Equal(got, []uint16{'a', 0xD801, 'b', 0xDC00}) {
		t.Errorf("lone surrogates: %04x", got)
	}
	split := StringFromUTF16([]uint16{0xD801}) + StringFromUTF16([]uint16{0xDC00}) + "Σ"
	if got := StringToUTF16(StringToLowerCaseRoot(split)); !slices.Equal(got, []uint16{0xD801, 0xDC28, 0x3C2}) {
		t.Errorf("split pair: %04x", got)
	}
}
