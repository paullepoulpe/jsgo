function f(/** null */ x) {}
(function(/** !Array<string>|!Iterable<number> */ it) {
  for (let elem of it) {
    f(elem);
  }
})(['']);
