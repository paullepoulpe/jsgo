async function f() {
  for await (const x in obj) {}
  for await (var y, z of arr) {}
  for await (let w = 1 of arr) {}
}
