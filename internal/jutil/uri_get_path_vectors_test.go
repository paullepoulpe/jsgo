package jutil

import (
	"strconv"
	"strings"
	"testing"
)

// Vectors generated with tools/oracle/java/JutilUriGetPathVectors.java.

func uriGetPathHex(t *testing.T, h string) *string {
	t.Helper()
	switch h {
	case "null":
		return nil
	case "-":
		s := ""
		return &s
	}
	var u []uint16
	for i := 0; i < len(h); i += 4 {
		v, err := strconv.ParseUint(h[i:i+4], 16, 16)
		if err != nil {
			t.Fatal(err)
		}
		u = append(u, uint16(v))
	}
	s := StringFromUTF16(u)
	return &s
}

func TestURIGetPathMatchesJDK(t *testing.T) {
	rows := readTSV(t, "uri_get_path.tsv")
	if len(rows) < 600 {
		t.Fatalf("only %d vectors", len(rows))
	}
	for _, r := range rows {
		in := *uriGetPathHex(t, r[0])
		var got []string
		if uri, err := ParseURI(in); err != nil {
			got = []string{"error"}
		} else {
			got = []string{uriVectorHex(uri.GetPath()), uriVectorHex(uri.path)}
		}
		if strings.Join(got, "\t") != strings.Join(r[1:], "\t") {
			t.Errorf("%q: got %v, want %v", in, got, r[1:])
		}
	}
}

func uriVectorHex(s *string) string {
	if s == nil {
		return "null"
	}
	if *s == "" {
		return "-"
	}
	var b strings.Builder
	for _, c := range StringToUTF16(*s) {
		b.WriteString(strings.Repeat("0", 4-len(strconv.FormatUint(uint64(c), 16))))
		b.WriteString(strconv.FormatUint(uint64(c), 16))
	}
	return b.String()
}
