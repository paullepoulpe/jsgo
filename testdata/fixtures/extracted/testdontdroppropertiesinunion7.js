/** @param {{a: number}|{a:number, b:string}} x */
function f(x) {}
/** @param {{c: number}|{c:number, d:string}} x */
function g(x) { return x.b; }
