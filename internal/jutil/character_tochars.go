package jutil

import (
	"fmt"
	"strings"
)

// Go-only: language-adapter — java.lang.Character code point conversions.

// CharacterToChars is Character.toChars(int codePoint): the UTF-16
// representation of a code point (a lone surrogate code point is returned
// as itself). It throws IllegalArgumentException for an invalid code point.
func CharacterToChars(codePoint int) []uint16 {
	switch {
	case codePoint >= 0 && codePoint < 0x10000:
		return []uint16{uint16(codePoint)}
	case codePoint >= 0x10000 && codePoint <= 0x10FFFF:
		cp := codePoint - 0x10000
		return []uint16{uint16(0xD800 + (cp >> 10)), uint16(0xDC00 + (cp & 0x3FF))}
	}
	panic(NewIllegalArgumentExceptionMsg(
		"Not a valid Unicode code point: 0x" + strings.ToUpper(fmt.Sprintf("%x", uint32(int32(codePoint))))))
}
