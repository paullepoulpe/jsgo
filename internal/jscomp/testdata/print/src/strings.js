// String literal escaping: quotes, control characters, HTML-sensitive
// sequences, line separators, non-ASCII and astral characters.
var a = "double \" and 'single'";
var b = 'only "doubles" here';
var c = 'it\'s';
var d = "tab\there\nnewline\rcr\fformfeed\bbackspace\vvertical\0nul";
var e = "\x00\x01\x1f\x7f\x80\xff";
var f = "</script> </SCRIPT> <!-- --> ]]> a=b&c<d>e";
var g = "\u2028 and \u2029 separators";
var g2 = "  raw separators  ";
var h = "caf\u00e9 \u4e2d\u6587 \ud83d\ude00 emoji";
var h2 = "café 中文 😀 raw";
var i = "lone \ud800 high and \udc00 low";
var j = "back\\slash and $dollar and `backtick`";
var k = '\'\'"';
var l = "\"\"'";
var m = "";
var n = "a\
b";
var o = {"quoted key": 1, "if": 2, "123": 3, "0x10": 4, "01": 5, "1e3": 6, "\u00e9": 7, 'valid': 8, "é": 9};
var p = o["prop"] + o['if'] + o["123"];
