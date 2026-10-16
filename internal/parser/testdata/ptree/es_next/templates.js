var t1 = `plain`;
var t2 = `a${b}c${d}e`;
var t3 = `nested ${`inner ${x}`} outer`;
var t4 = tag`raw ${x} \unicode and \xerror`;
var t5 = `escapes \n \t \\ \` \$ \u0041 \u{1F600} \x41 \0`;
var t6 = `line
continuation \
done`;
var t7 = `\'`;
var t8 = `\a`;
var t9 = a.b`member tag`;
var t10 = f()`call tag`;
var t11 = `${a}${b}`;
var t12 = ``;
var t13 = `$notsub ${'$'} $`;
var t14 = tag`\u{} \xz \01 \1`;
var t15 = `\u{110000} \u{0000041}`;
