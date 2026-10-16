const a = { b: [3] };
label: for (const val of a?.b) {
  if (val != 3) {
    continue label;
  }
}
