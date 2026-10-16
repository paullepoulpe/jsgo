class N1 { constructor() { this.t = new.target; var f = () => new.target; } }
class N2 extends N1 { constructor() { super(); if (new.target === N2) {} } }
