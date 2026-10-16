const a = {b:[3]};
let tmp0;
label: for (const val of (tmp0 = a) == null ? void 0 : tmp0.b) {
  if (val != 3) {
  continue label;
  }
}
