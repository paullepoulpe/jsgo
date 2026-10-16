var f1=function(){
  var x
};
(function() {
  function f2() {
    alert(x)
  }
  {
    var x=0
  }
  f2()
})()
