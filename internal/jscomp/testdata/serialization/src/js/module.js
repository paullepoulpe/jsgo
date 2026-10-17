export const x = 1 + 2;
export default function f() { return x + import.meta.url; }
export {x as y};
export class C { m() { return import('./lazy.js'); } }
