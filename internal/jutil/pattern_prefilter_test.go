package jutil

import (
	"strings"
	"testing"
)

func TestPatternPrefilterAndDotStarEquivalence(t *testing.T) {
	patterns := []string{
		`.*`,
		`(?s).*`,
		`.*foo.*`,
		`(?s).*foo.*`,
		`.*foo`,
		`foo.*`,
		`(?s).*foo`,
		`(?s)foo.*`,
		`.*(eval|Function|setTimeout).*`,
		`(?i).*eval.*`,
		`.*foo[0-9]+bar.*`,
		`.*\bfoo\b.*`,
		`.*^foo.*`,
		`.*foo$.*`,
		`(?m).*^foo.*`,
		`.*a\nb.*`,
		`.*foo|bar`,
		`foo|.*`,
	}

	inputs := []string{
		"",
		"foo",
		"FOO",
		"xxfooyy",
		"xxFOOyy",
		"foobar",
		"a foo b",
		"xxfoo123baryy",
		"xxfoo123\nbaryy",
		"bar",
		"xxbar",
		"bar\nfoo",
		"xxa\nbyy",
		"xxfoo\nyy",
		"xxfoo\ryy",
		"xxfoo\u0085yy",
		"xxfoo\u2028yy",
		"xxfoo\u2029yy",
		"before setTimeout( after",
		"no banned token here\r\n",
		"\xfffoo\xfe",
	}

	for _, pat := range patterns {
		p := patternCompile(pat)
		for _, in := range inputs {
			wantMatches := p.full.MatchString(in)
			if got := p.Matches(in); got != wantMatches {
				t.Errorf("Pattern(%q).Matches(%q) = %v, want %v (mode=%d canTerm=%v pre=%+v)",
					pat, in, got, wantMatches, p.matchesCoreMode, p.matchesCoreCanMatchTerm, p.fullPrefilter)
			}
			m := p.Matcher(in)
			if got := m.Matches(); got != wantMatches {
				t.Errorf("Matcher(%q).Matches(%q) = %v, want %v", pat, in, got, wantMatches)
			}
			wantFind := p.searcherFor(in, 0).re.MatchString(in)
			if got := p.Find(in); got != wantFind {
				t.Errorf("Pattern(%q).Find(%q) = %v, want %v (pre=%+v)",
					pat, in, got, wantFind, p.findPrefilter)
			}
		}
	}
}

func BenchmarkPatternMatchesDotStarLargeLiteral(b *testing.B) {
	payload := strings.Repeat("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/", 65536) // 4 MiB
	pSingle := Pattern_Compile(`.*dangerouslySetInnerHTML.*`)
	pAlt := Pattern_Compile(`.*(dangerouslySetInnerHTML|javascript:|vbscript:).*`)
	payloadHit := payload[:len(payload)/2] + "dangerouslySetInnerHTML" + payload[len(payload)/2:]

	b.Run("Miss4MiB_Single", func(b *testing.B) {
		b.SetBytes(int64(len(payload)))
		b.ReportAllocs()
		for b.Loop() {
			if pSingle.Matches(payload) {
				b.Fatal("unexpected match")
			}
		}
	})
	b.Run("Miss4MiB_Alt3", func(b *testing.B) {
		b.SetBytes(int64(len(payload)))
		b.ReportAllocs()
		for b.Loop() {
			if pAlt.Matches(payload) {
				b.Fatal("unexpected match")
			}
		}
	})
	b.Run("Hit4MiB_Single", func(b *testing.B) {
		b.SetBytes(int64(len(payloadHit)))
		b.ReportAllocs()
		for b.Loop() {
			if !pSingle.Matches(payloadHit) {
				b.Fatal("expected match")
			}
		}
	})
}
