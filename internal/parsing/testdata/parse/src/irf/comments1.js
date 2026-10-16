// leading line comment
/* leading block */
var a = 1; // trailing on a
var b = 2; /* trailing block on b */ // and line
function f(/* p1 */ x, // after x
    y /* after y */,
    /** @type {number} */ z) { // after brace
  // inside body
  return x; // after return
  // last in body
}
f(1 /* a1 */, 2, // a2
  3); // after call
{
  // only comment in block
}
{
  stmt(); /* t1 */
  /* end1 */ // end2
}
if (a) /* c */ { b(); } else /* d */ { c(); }
let {x: /* k */ y, /* s */ z = 1, ...rest} = o;
let [/* e1 */ p, , /* e2 */ q = 2, ...r] = arr;
class K { /* m */ m() {} /* f */ field = 1; /* g */ get g() { return 1; } static /* s */ s() {} }
obj./* prop */foo; obj?./* oc */bar;
label: /* lbl */ for (;;) { break /* b */ label; }
import(/* dyn */ 'x');
x = /* rhs */ y + /* rhs2 */ z;
/** jsdoc not attached */
// end of file comment
