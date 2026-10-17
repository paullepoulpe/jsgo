package jutil

// Go-only: language-adapter — java.util.Properties.load(Reader), following
// the line grammar of its Javadoc:
//   - natural lines end at \n, \r or \r\n (or at the end of the input);
//   - a natural line that holds only white space (space, \t, \f) is skipped,
//     and so is one whose first other character is '#' or '!' (a comment;
//     comments never continue on the next line);
//   - a natural line that ends in an odd number of backslashes continues on
//     the next natural line: that last backslash and the line break are
//     dropped, and so is the white space at the start of the next line;
//   - in the resulting logical line the key runs from the first character to
//     the first '=', ':' or white space character that is not escaped by a
//     backslash. White space after the key is skipped, then at most one '='
//     or ':', then white space again; the rest is the value;
//   - in keys and values \t, \n, \r, \f and \uXXXX stand for those
//     characters, and a backslash before any other character is dropped.
//
// The entries are returned in the order their keys first appear; a key
// that appears again replaces the value. (Properties is a hash table, whose
// order no caller depends on.)
//
// Observed behaviour (black-box, JDK 21):
//   - \u must be followed by exactly four hex digits, otherwise load throws
//     IllegalArgumentException("Malformed \uxxxx encoding.");
//   - while nothing has been read for a logical line, a continuation line is
//     read like a new line: it may be blank or a comment ("\\\n#c\nk=v"
//     defines only k). Once something has been read, a continuation line
//     that is empty or holds only white space ends the logical line
//     ("a\\\n   \nc=d" gives a="" and c=d), and '#' or '!' at its start is
//     ordinary text;
//   - a logical line cut short by the end of the input right after the
//     backslash or after a "\n" or "\r" defines a key even when nothing was
//     read (" \\\n" and "\\" define ""=""); after a "\r\n" it does not
//     ("\\\r\n" defines nothing).

// PropertiesLoad is new Properties().load(new StringReader(text)), with the
// entries in the order described above.
func PropertiesLoad(text string) *LinkedHashMap[string, string] {
	in := StringToUTF16(text)
	result := NewLinkedHashMap[string, string]()
	pos := 0
lines:
	for pos < len(in) {
		start, end, next := propertiesNaturalLine(in, pos)
		pos = next
		i := propertiesSkipSpace(in, start, end)
		if i == end || in[i] == '#' || in[i] == '!' {
			continue // blank or comment
		}
		var logical []uint16
		for {
			if !propertiesContinues(in, i, end) {
				logical = append(logical, in[i:end]...)
				break
			}
			logical = append(logical, in[i:end-1]...)
			if pos >= len(in) {
				if len(logical) == 0 && next-end == 2 {
					continue lines // nothing after "\r\n": no entry
				}
				break
			}
			if len(logical) == 0 {
				continue lines // read the next natural line as a new one
			}
			start, end, next = propertiesNaturalLine(in, pos)
			pos = next
			i = propertiesSkipSpace(in, start, end)
			if i == end {
				break // a blank continuation line ends the logical line
			}
		}
		key, value := propertiesSplitLine(logical)
		result.Put(StringFromUTF16(propertiesUnescape(key)), StringFromUTF16(propertiesUnescape(value)))
	}
	return result
}

// propertiesNaturalLine returns the bounds of the natural line starting at
// pos and the start of the following one.
func propertiesNaturalLine(in []uint16, pos int) (start, end, next int) {
	end = pos
	for end < len(in) && in[end] != '\n' && in[end] != '\r' {
		end++
	}
	next = end
	if next < len(in) {
		if in[next] == '\r' && next+1 < len(in) && in[next+1] == '\n' {
			next++
		}
		next++
	}
	return pos, end, next
}

func propertiesIsSpace(c uint16) bool { return c == ' ' || c == '\t' || c == '\f' }

func propertiesSkipSpace(in []uint16, i, end int) int {
	for i < end && propertiesIsSpace(in[i]) {
		i++
	}
	return i
}

// propertiesContinues reports whether in[i:end] ends in an odd number of
// backslashes.
func propertiesContinues(in []uint16, i, end int) bool {
	n := 0
	for j := end - 1; j >= i && in[j] == '\\'; j-- {
		n++
	}
	return n%2 == 1
}

// propertiesSplitLine splits a logical line (escapes not yet decoded) into
// its key and value.
func propertiesSplitLine(line []uint16) (key, value []uint16) {
	k := 0
	for k < len(line) {
		c := line[k]
		if c == '\\' {
			k += 2
			continue
		}
		if c == '=' || c == ':' || propertiesIsSpace(c) {
			break
		}
		k++
	}
	if k > len(line) {
		k = len(line)
	}
	v := propertiesSkipSpace(line, k, len(line))
	if v < len(line) && (line[v] == '=' || line[v] == ':') {
		v = propertiesSkipSpace(line, v+1, len(line))
	}
	return line[:k], line[v:]
}

// propertiesUnescape decodes the backslash escapes of a key or a value.
func propertiesUnescape(s []uint16) []uint16 {
	out := make([]uint16, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c != '\\' {
			out = append(out, c)
			continue
		}
		i++
		if i == len(s) {
			break // a lone backslash at the very end of the input
		}
		switch c = s[i]; c {
		case 't':
			out = append(out, '\t')
		case 'n':
			out = append(out, '\n')
		case 'r':
			out = append(out, '\r')
		case 'f':
			out = append(out, '\f')
		case 'u':
			var v uint16
			for k := 1; k <= 4; k++ {
				d := -1
				if i+k < len(s) {
					d = hexDigitValue(s[i+k])
				}
				if d < 0 {
					panic(NewIllegalArgumentExceptionMsg(`Malformed \uxxxx encoding.`))
				}
				v = v<<4 | uint16(d)
			}
			out = append(out, v)
			i += 4
		default:
			out = append(out, c)
		}
	}
	return out
}

func hexDigitValue(c uint16) int {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c >= 'a' && c <= 'f':
		return int(c-'a') + 10
	case c >= 'A' && c <= 'F':
		return int(c-'A') + 10
	}
	return -1
}
