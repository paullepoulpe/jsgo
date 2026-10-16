var t1 = `plain`;
var t2 = `with ${x} and ${y + z} subs`;
var t3 = `nested ${`inner ${deep}`} done`;
var t4 = tag`raw \n ${a} \u{1F600} \xff`;
var t5 = `multi
line
  template`;
var t6 = `caf\u00e9 café 😀 \u2028   ${"str"}`;
var t7 = String.raw`\unicode and \x`;
var t8 = `$ {not} $${yes} \${escaped}`;
var t9 = a.b`member tag`;
var t10 = (() => 1)`arrow tag`;
var t11 = `${{a: 1}}`;
var t12 = `${function () {}}`;
