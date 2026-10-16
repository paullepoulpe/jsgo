package jutil

// CharStreams_ReadLines is Guava CharStreams.readLines(Readable): the lines
// of the input, split at "\n", "\r" and "\r\n" (Guava's LineReader). The
// line terminators are not included, and an input that ends with a
// terminator has no empty last line.
func CharStreams_ReadLines(r JavaReader) []string {
	units := StringToUTF16(CharStreams_ToString(r))
	var result []string
	start := 0
	for i := 0; i < len(units); i++ {
		switch units[i] {
		case '\n':
			result = append(result, StringFromUTF16(units[start:i]))
			start = i + 1
		case '\r':
			result = append(result, StringFromUTF16(units[start:i]))
			if i+1 < len(units) && units[i+1] == '\n' {
				i++
			}
			start = i + 1
		}
	}
	if start < len(units) {
		result = append(result, StringFromUTF16(units[start:]))
	}
	return result
}
