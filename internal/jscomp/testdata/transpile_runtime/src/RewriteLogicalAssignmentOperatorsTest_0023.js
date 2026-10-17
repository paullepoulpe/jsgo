foo().x = 4;
do {
  foo().x = foo().x - 1;
} while (foo().x ||= null);
