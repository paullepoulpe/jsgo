var Di = I(() => {
  function zv() {
    var temp_const$jscomp$0 = JSCOMPILER_PRESERVE(e);
// TODO(b/339701959): We decided to back off here because of b/338660589. But we should
// optimize this to:
//   var temp_const$jscomp$1 = CN();
//   var temp_const$jscomp$2 = getObj();
//   temp_const$jscomp$0, [temp_const$jscomp$2.propName] = temp_const$jscomp$1;
    temp_const$jscomp$0, [getObj().propName] = CN();
  }
  function CN() {
    return [1];
  }
});
