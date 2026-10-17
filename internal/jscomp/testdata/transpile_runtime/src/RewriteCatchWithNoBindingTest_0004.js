try {
  stuff();
} catch {
  try {
    onError();
  } catch {
    shruggie();
  }
}
