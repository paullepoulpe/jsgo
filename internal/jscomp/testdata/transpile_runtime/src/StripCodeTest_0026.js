class A {
  constructor() {
    this.logger = goog$debug$Logger.getLogger('A');
    this.otherProperty = 3;
  }
}
let a = new A;
a.logger.warning('foobar');
