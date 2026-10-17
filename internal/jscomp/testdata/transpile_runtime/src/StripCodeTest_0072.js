var a = {};
a.b = {
  x: 0,
  ['logger']: null // Don't strip computed properties
};
a.b['logger']()
