class P { static sm() { return 1; } im() { return 2; } get acc() { return 3; } }
class Q extends P {
  static sm() { return super.sm() + super['sm'](); }
  im() { return super.im() + super.acc + super['im'](); }
  static f() { return super.sm; }

}

class R { constructor() { this.x = 1; } }
class S extends R { constructor() { var f = () => this.x; super(); f(); } }
class T extends R { constructor(a) { if (a) { super(); } else { super(); } } }
