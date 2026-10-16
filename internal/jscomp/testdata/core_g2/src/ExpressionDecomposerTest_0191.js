function first() {
  alert('first');
      return '';
}
var temp_const$jscomp$0 = alert;
temp_const$jscomp$0(first().method(
    alert('second')).method(alert('third')));
