package jutil

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"testing"
)

// Vectors from tools/oracle/java/JavaCaseTable.java.
func TestCharacterCaseMatchesJava(t *testing.T) {
	f, err := os.Open("testdata/character_case.tsv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	want := map[uint16][2]uint16{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		p := strings.Fields(sc.Text())
		c, _ := strconv.ParseUint(p[0], 16, 16)
		u, _ := strconv.ParseUint(p[1], 16, 16)
		l, _ := strconv.ParseUint(p[2], 16, 16)
		want[uint16(c)] = [2]uint16{uint16(u), uint16(l)}
	}
	if len(want) < 2000 {
		t.Fatalf("only %d entries", len(want))
	}
	for c := 0; c < 0x10000; c++ {
		ch := uint16(c)
		u, l := ch, ch
		if w, ok := want[ch]; ok {
			u, l = w[0], w[1]
		}
		if got := CharacterToUpperCase(ch); got != u {
			t.Errorf("toUpperCase(%04X) = %04X, want %04X", c, got, u)
		}
		if got := CharacterToLowerCase(ch); got != l {
			t.Errorf("toLowerCase(%04X) = %04X, want %04X", c, got, l)
		}
	}
	// regionMatches(true, ...): dotless i and long s fold to ASCII.
	script := StringToUTF16("/script")
	for _, s := range []string{"/SCRIPT", "/scrıpt", "/ſcript", "/scrİpt", "/scripx"} {
		got := UTF16RegionMatches(true, StringToUTF16(s), 0, script, 0, len(script))
		if got != (s != "/scripx") {
			t.Errorf("regionMatches(%q) = %v", s, got)
		}
	}
}
