import {a} from './a.js';
import * as ns from './b.js';
export const c = a + 1;
export function f() { if (c) return ns.x; }
export default class { m() { return f(); } }
for (const k of [1, 2]) { if (k) break; }
