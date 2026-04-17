/**@param {string} a*/function foo(a){
  /**@param {Array|string} a*/function bar(a){
    if (a instanceof Array)return;
  }
}
