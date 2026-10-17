let counter = 0;
class F1 { x = counter++; y; static z = counter++; [`k${counter}`] = 1; static [counter++] = 2; }
class F2 extends F1 { a = this.x; constructor() { super(); this.b = 1; } }
class F3 extends F1 { a = 1; }
const F5 = class Named { static self = Named; static other = this; m() { return Named; } };
x.F6 = class { static s = 1; };
use(class F7 { static s = F7; });
class F8 extends (function() { return F1; })() { static q = super.z; }
class F9 { static { let t = 1; this.t = t; } static {} }
