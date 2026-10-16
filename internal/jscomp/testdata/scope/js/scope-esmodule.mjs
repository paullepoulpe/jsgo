import def, {a as b, c} from './dep.js';
import * as ns from './ns.js';
export {b as renamed, c};
export default function defFn() { return b + c + ns.x; }
export const k = 1;
export let l = 2;
export var m = 3;
export class ExportedClass {}
export function exportedFn(x = l) { return x; }
let local = def;
local++;
export {local as default2};
export * from './star.js';
export * from "./star.js";
const top = this;
await Promise.resolve(top);
