function first() {
  alert('first');
  return '';
}
// alert must be preserved before the first side-effect
alert(first().method(alert('second')).method(alert('third')));
