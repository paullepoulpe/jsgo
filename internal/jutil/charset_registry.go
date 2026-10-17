package jutil

import "golang.org/x/text/encoding/charmap"

// Go-only: language-adapter — the charsets jutil supports, with the
// canonical names and aliases the JDK reports for them (Charset.name() and
// Charset.aliases(), observed black-box on JDK 21; see
// tools/oracle/java/JutilCharsetVectors.java). Charset.forName finds a
// charset by any of these names, ignoring case. The JDK supports more
// charsets (EBCDIC, the IBM and Mac code pages not listed here, ...); for
// those, Charset.isSupported returns false and Charset.forName throws
// UnsupportedCharsetException, as a JDK without them would.

var (
	StandardCharsets_US_ASCII   = newCharset("US-ASCII", newUSASCII, "646", "ANSI_X3.4-1968", "ANSI_X3.4-1986", "ASCII", "IBM367", "ISO646-US", "ISO_646.irv:1991", "ascii7", "cp367", "csASCII", "iso-ir-6", "iso_646.irv:1983", "us")
	StandardCharsets_ISO_8859_1 = newCharset("ISO-8859-1", newSBCS(charmap.ISO8859_1, false), "819", "8859_1", "IBM-819", "IBM819", "ISO8859-1", "ISO8859_1", "ISO_8859-1", "ISO_8859-1:1987", "ISO_8859_1", "cp819", "csISOLatin1", "iso-ir-100", "l1", "latin1")
	StandardCharsets_UTF_8      = newCharset("UTF-8", func() charsetCoder { return utf8Coder{} }, "UTF8", "unicode-1-1-utf-8")
	StandardCharsets_UTF_16     = newCharset("UTF-16", func() charsetCoder { return utf16Coder{bigEndian: true, marked: true} }, "UTF_16", "UnicodeBig", "unicode", "utf16")
	StandardCharsets_UTF_16BE   = newCharset("UTF-16BE", func() charsetCoder { return utf16Coder{bigEndian: true} }, "ISO-10646-UCS-2", "UTF_16BE", "UnicodeBigUnmarked", "X-UTF-16BE")
	StandardCharsets_UTF_16LE   = newCharset("UTF-16LE", func() charsetCoder { return utf16Coder{} }, "UTF_16LE", "UnicodeLittleUnmarked", "X-UTF-16LE")
)

func newCharset(name string, build func() charsetCoder, aliases ...string) *Charset {
	return &Charset{name: name, aliases: aliases, build: build}
}

// supportedCharsets are the charsets Charset.forName finds.
var supportedCharsets = []*Charset{
	StandardCharsets_US_ASCII, StandardCharsets_ISO_8859_1, StandardCharsets_UTF_8,
	StandardCharsets_UTF_16, StandardCharsets_UTF_16BE, StandardCharsets_UTF_16LE,

	newCharset("UTF-32", func() charsetCoder { return utf32Coder{bigEndian: true, detect: true} }, "UTF32", "UTF_32"),
	newCharset("UTF-32BE", func() charsetCoder { return utf32Coder{bigEndian: true} }, "UTF_32BE", "X-UTF-32BE"),
	newCharset("UTF-32LE", func() charsetCoder { return utf32Coder{} }, "UTF_32LE", "X-UTF-32LE"),
	newCharset("X-UTF-32BE-BOM", func() charsetCoder { return utf32Coder{bigEndian: true, writeMark: true} }, "UTF-32BE-BOM", "UTF_32BE_BOM"),
	newCharset("X-UTF-32LE-BOM", func() charsetCoder { return utf32Coder{writeMark: true} }, "UTF-32LE-BOM", "UTF_32LE_BOM"),

	newCharset("windows-1250", newSBCS(charmap.Windows1250, false), "cp1250", "cp5346"),
	newCharset("windows-1251", newSBCS(charmap.Windows1251, false), "ansi-1251", "cp1251", "cp5347"),
	newCharset("windows-1252", newSBCS(charmap.Windows1252, false), "cp1252", "cp5348", "ibm-1252", "ibm1252"),
	newCharset("windows-1253", newSBCS(charmap.Windows1253, false), "cp1253", "cp5349"),
	newCharset("windows-1254", newSBCS(charmap.Windows1254, false), "cp1254", "cp5350"),
	newCharset("windows-1255", sbcsWindows1255, "cp1255"),
	newCharset("windows-1256", newSBCS(charmap.Windows1256, false), "cp1256"),
	newCharset("windows-1257", newSBCS(charmap.Windows1257, false), "cp1257", "cp5353"),
	newCharset("windows-1258", newSBCS(charmap.Windows1258, false), "cp1258"),
	newCharset("x-windows-874", newSBCS(charmap.Windows874, false), "ms-874", "ms874", "windows-874"),

	newCharset("ISO-8859-2", newSBCS(charmap.ISO8859_2, true), "8859_2", "912", "ISO8859-2", "ISO_8859-2", "ISO_8859-2:1987", "cp912", "csISOLatin2", "ibm-912", "ibm912", "iso-ir-101", "iso8859_2", "l2", "latin2"),
	newCharset("ISO-8859-3", newSBCS(charmap.ISO8859_3, true), "8859_3", "913", "ISO8859-3", "ISO_8859-3", "ISO_8859-3:1988", "cp913", "csISOLatin3", "ibm-913", "ibm913", "iso-ir-109", "iso8859_3", "l3", "latin3"),
	newCharset("ISO-8859-4", newSBCS(charmap.ISO8859_4, true), "8859_4", "914", "ISO_8859-4", "ISO_8859-4:1988", "cp914", "csISOLatin4", "ibm-914", "ibm914", "iso-ir-110", "iso8859-4", "iso8859_4", "l4", "latin4"),
	newCharset("ISO-8859-5", newSBCS(charmap.ISO8859_5, true), "8859_5", "915", "ISO8859-5", "ISO_8859-5", "ISO_8859-5:1988", "cp915", "csISOLatinCyrillic", "cyrillic", "ibm-915", "ibm915", "iso-ir-144", "iso8859_5"),
	newCharset("ISO-8859-6", newSBCS(charmap.ISO8859_6, true), "1089", "8859_6", "ASMO-708", "ECMA-114", "ISO8859-6", "ISO_8859-6", "ISO_8859-6:1987", "arabic", "cp1089", "csISOLatinArabic", "ibm-1089", "ibm1089", "iso-ir-127", "iso8859_6"),
	newCharset("ISO-8859-7", newSBCS(charmap.ISO8859_7, true), "813", "8859_7", "ECMA-118", "ELOT_928", "ISO_8859-7", "ISO_8859-7:1987", "cp813", "csISOLatinGreek", "greek", "greek8", "ibm-813", "ibm813", "iso-ir-126", "iso8859-7", "iso8859_7", "sun_eu_greek"),
	newCharset("ISO-8859-8", newSBCS(charmap.ISO8859_8, true), "8859_8", "916", "ISO8859-8", "ISO_8859-8", "ISO_8859-8:1988", "cp916", "csISOLatinHebrew", "hebrew", "ibm-916", "ibm916", "iso-ir-138", "iso8859_8"),
	newCharset("ISO-8859-9", newSBCS(charmap.ISO8859_9, true), "8859_9", "920", "ISO8859-9", "ISO_8859-9", "ISO_8859-9:1989", "cp920", "csISOLatin5", "ibm-920", "ibm920", "iso-ir-148", "iso8859_9", "l5", "latin5"),
	// x-iso-8859-11 is windows-874 with the C1 controls at 0x80 to 0x9F.
	newCharset("x-iso-8859-11", newSBCS(charmap.Windows874, true), "iso-8859-11", "iso8859_11"),
	newCharset("ISO-8859-13", newSBCS(charmap.ISO8859_13, true), "8859_13", "ISO8859-13", "iso8859_13", "iso_8859-13"),
	newCharset("ISO-8859-15", newSBCS(charmap.ISO8859_15, true), "8859_15", "923", "IBM-923", "IBM923", "ISO-8859-15", "ISO8859-15", "ISO8859_15", "ISO8859_15_FDIS", "ISO_8859-15", "L9", "LATIN0", "LATIN9", "Latin-9", "cp923", "csISO885915", "csISOlatin0", "csISOlatin9"),
	newCharset("ISO-8859-16", newSBCS(charmap.ISO8859_16, true), "ISO8859_16", "ISO_8859-16", "ISO_8859-16:2001", "csISO885916", "iso-ir-226", "l10", "latin10"),
	newCharset("TIS-620", sbcsTIS620, "tis620", "tis620.2533"),

	newCharset("KOI8-R", newSBCS(charmap.KOI8R, false), "cskoi8r", "koi8", "koi8_r"),
	newCharset("KOI8-U", sbcsKOI8U, "koi8_u"),

	newCharset("IBM437", newSBCS(charmap.CodePage437, false), "437", "cp437", "cspc8codepage437", "ibm-437", "ibm437", "windows-437"),
	newCharset("IBM850", newSBCS(charmap.CodePage850, false), "850", "cp850", "cspc850multilingual", "ibm-850", "ibm850"),
	newCharset("IBM852", newSBCS(charmap.CodePage852, false), "852", "cp852", "csPCp852", "ibm-852", "ibm852"),
	newCharset("IBM855", newSBCS(charmap.CodePage855, false), "855", "cp855", "cspcp855", "ibm-855", "ibm855"),
	newCharset("IBM00858", newSBCS(charmap.CodePage858, false), "858", "PC-Multilingual-850+euro", "ccsid00858", "cp00858", "cp858", "ibm-858", "ibm858"),
	newCharset("IBM860", newSBCS(charmap.CodePage860, false), "860", "cp860", "csIBM860", "ibm-860", "ibm860"),
	newCharset("IBM862", newSBCS(charmap.CodePage862, false), "862", "cp862", "csIBM862", "cspc862latinhebrew", "ibm-862", "ibm862"),
	newCharset("IBM863", newSBCS(charmap.CodePage863, false), "863", "cp863", "csIBM863", "ibm-863", "ibm863"),
	newCharset("IBM865", newSBCS(charmap.CodePage865, false), "865", "cp865", "csIBM865", "ibm-865", "ibm865"),
	newCharset("IBM866", newSBCS(charmap.CodePage866, false), "866", "cp866", "csIBM866", "ibm-866", "ibm866"),

	newCharset("x-MacRoman", newSBCS(charmap.Macintosh, false), "MacRoman"),
	newCharset("x-MacCyrillic", sbcsMacCyrillic, "MacCyrillic"),

	newCharset("Shift_JIS", newShiftJIS, "csShiftJIS", "ms_kanji", "shift-jis", "shift_jis", "sjis", "x-sjis"),
	newCharset("windows-31j", newWindows31J, "MS932", "csWindows31J", "windows-932"),
	newCharset("EUC-JP", newEUCJP, "Extended_UNIX_Code_Packed_Format_for_Japanese", "csEUCPkdFmtjapanese", "euc_jp", "eucjis", "eucjp", "x-euc-jp", "x-eucjp"),
	newCharset("ISO-2022-JP", newISO2022JP, "csISO2022JP", "csjisencoding", "iso2022jp", "jis", "jis_encoding"),
	newCharset("GBK", newGBK, "CP936", "windows-936"),
	newCharset("GB2312", newGB2312, "EUC_CN", "euc-cn", "euccn", "gb2312", "gb2312-1980", "gb2312-80", "x-EUC-CN"),
	newCharset("GB18030", newGB18030, "gb18030-2022"),
	newCharset("Big5", newBig5, "csBig5"),
	newCharset("EUC-KR", newEUCKR, "5601", "csEUCKR", "euc_kr", "euckr", "ks_c_5601-1987", "ksc5601", "ksc5601-1987", "ksc5601_1987", "ksc_5601"),
	newCharset("x-windows-949", newWindows949, "ms949", "ms_949", "windows-949", "windows949"),
}
