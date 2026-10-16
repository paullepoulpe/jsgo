// Port of test/com/google/javascript/jscomp/serialization/Wtf8Test.java @ closure-compiler v20260415

package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func TestWtf8_Ascii(t *testing.T) {
	jsaWtf8Roundtrip(t, "hello world")
}

func TestWtf8_TwoByteLatin1SupplementCharacter(t *testing.T) {
	// Two-bytes in UTF-8: https://www.compart.com/en/unicode/U+00BC
	jsaWtf8Roundtrip(t, "¼")
}

func TestWtf8_MultibyteBMPCharacters(t *testing.T) {
	// Korean characters are 3 bytes in UTF-8. e.g. https://www.compart.com/en/unicode/U+C138
	jsaWtf8Roundtrip(t, "세계님 안녕")
}

func TestWtf8_EscapedPairedSurrogate(t *testing.T) {
	// This is the four-byte character "💅" in UTF-8: https://www.compart.com/en/unicode/U+1F485
	jsaWtf8Roundtrip(t, jutil.StringFromUTF16([]uint16{0xd83d, 0xdc85}))
}

func TestWtf8_UnpairedSurrogates(t *testing.T) {
	jsaWtf8Roundtrip(t, "Hello "+jutil.StringFromUTF16([]uint16{0xd800}))
	jsaWtf8Roundtrip(t, jutil.StringFromUTF16([]uint16{0xd800})+", hello")
}

func TestWtf8_InvalidSurrogatePairings(t *testing.T) {
	// This is an invalid pair: low-high
	jsaWtf8Roundtrip(t, jutil.StringFromUTF16([]uint16{0xdc37, 0xd801}))
	// This is an invalid pair: low-low
	jsaWtf8Roundtrip(t, jutil.StringFromUTF16([]uint16{0xd802, 0xd801}))
	// This is an invalid pair: high-high
	jsaWtf8Roundtrip(t, jutil.StringFromUTF16([]uint16{0xdc37, 0xdc39}))
}

func TestWtf8_UnescapedPairedSurrogates(t *testing.T) {
	// Emojis are always 4 bytes in UTF-8 and thus represented as paired surrogates in UTF-16.
	jsaWtf8Roundtrip(t, "💅🏽")
}

var jsaWtf8Decoder = Wtf8_Decoder_(256)

func jsaWtf8Roundtrip(t *testing.T, s string) {
	t.Helper()
	serialized := Wtf8_EncodeToWtf8(s)
	roundtripped := jsaWtf8Decoder.Decode(serialized)
	jsaEq(t, roundtripped, s)
}
