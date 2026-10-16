let temp_const$jscomp$0;
let temp_const$jscomp$1;
if ((temp_const$jscomp$0 = nonOptional.part) == null) {
  void 0;
} else {
// `temp_const$jscomp$0.optional` is the expression we were trying to expose, and it
// is now movable.
  (temp_const$jscomp$1 = temp_const$jscomp$0.optional.chain) == null
      ? void 0 : temp_const$jscomp$1.continues();
}
