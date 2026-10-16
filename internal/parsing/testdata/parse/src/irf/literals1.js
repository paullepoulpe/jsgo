var s1 = 'a\'b"c\\d\n\t\r\b\f\v\0\x41\u0042\u{43}\u{1F600}';
var s2 = "line\
continuation";
var s3 = '\1\12\123\1234\8\9';
var s4 = `tpl\n${x}\u{41}\x42
raw\r\nline`;
var s5 = tag`\unicode bad escape ${1}`;
var n1 = 0x1F + 0X1f + 0b101 + 0B1 + 0o17 + 0O7 + 017 + 019 + 08 + 09 + 1_000_000 + .5 + 5. + 1e10 + 1E-5 + 0.0 + 0e1;
var n2 = 123456789012345678901234567890 + 0xFFFFFFFFFFFFFFFFFFFF + 1.7976931348623159e308;
var b1 = 0n + 1n + 0x1Fn + 0b11n + 0o17n + 1_0n + 12345678901234567890123n;
var r1 = /a\/b\.c\d\w\s\S\-\]\[\^\$\q\!/gimsuyd;
var r2 = /[/]/;
var r3 = /x/gz;
var o = {0x10: 1, 1.5: 2, 1e3: 3, .5: 4, 10n: 5, 'str\x41': 6, "\u0041": 7};
var u = -x + +y + !z + ~w + typeof v + void 0 + delete o.p + delete o[1] + delete x;
delete 5; delete (a, b); +1n; -1n;
x++; ++x; x--; --x; (1)++; ++(a.b); a?.b++; (/** @type {number} */ (x))++;
