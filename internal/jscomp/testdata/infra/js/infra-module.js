import def, {a as b, c} from './dep.js';
import * as ns from './ns.js';
export {b, c as d};
export * from './all.js';
export default function () { return ns.x + def; }
export const meta = import.meta.url;
export class K {}
let top = await Promise.resolve(1);
