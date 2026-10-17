package jutil

import "strings"

// Go-only: language-adapter — java.net.URI.getPath(), written from its
// Javadoc ("the decoded path component": every sequence of escaped octets
// is replaced by the characters it represents in UTF-8, malformed input
// replaced as the UTF-8 decoder in REPLACE mode does) and checked against
// JDK output (tools/oracle/java/JutilUriGetPathVectors.java).

// GetPath is URI.getPath(): the decoded path component, nil if the path is
// undefined.
func (u *URI) GetPath() *string {
	raw := u.GetRawPath()
	if raw == nil {
		return nil
	}
	s := uriDecodeEscapedOctets(*raw)
	return &s
}

// uriDecodeEscapedOctets replaces each maximal run of escaped octets
// ("%XX%YY...") in s by the UTF-8 decoding of those octets. A URI that
// parsed only contains well-formed escapes.
func uriDecodeEscapedOctets(s string) string {
	if !strings.Contains(s, "%") {
		return s
	}
	units := StringToUTF16(s)
	out := make([]uint16, 0, len(units))
	for i := 0; i < len(units); {
		if units[i] != '%' {
			out = append(out, units[i])
			i++
			continue
		}
		var octets []byte
		for i+2 < len(units) && units[i] == '%' {
			octets = append(octets, byte(uriGetPathHexValue(units[i+1])<<4|uriGetPathHexValue(units[i+2])))
			i += 3
		}
		out = append(out, StringToUTF16(StandardCharsets_UTF_8.Decode(octets))...)
	}
	return StringFromUTF16(out)
}

// uriGetPathHexValue is the value of the hexadecimal digit c.
func uriGetPathHexValue(c uint16) int {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c >= 'a' && c <= 'f':
		return int(c-'a') + 10
	default:
		return int(c-'A') + 10
	}
}
