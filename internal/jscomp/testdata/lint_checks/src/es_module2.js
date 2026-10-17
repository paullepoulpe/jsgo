goog.declareModuleId('my.es');
import * as ns from './a.js';
import def, {a as b} from './a.js';
export {b, def as c};
export * from './c.js';
export class ExportedClass {}
export function exportedFn() {}
export default class Named {}
export const [d1, d2] = arr, {e1} = obj;
let MUT = 1, mut2 = 2;
export {MUT, mut2};
function change() { MUT = 2; mut2++; d1 = 3; ExportedClass = null; exportedFn = null; Named = 1; e1 = 0; ns.x = 1; }
var LOCAL_VAR = 4;
let LET_CONST = 5;
/** @const */ let JSDOC_CONST = 6;
function shadow() { let LET_CONST = 7; LET_CONST = 8; }
const goog = {require() {}};
goog.require('late.require');
const req = goog.require('late2');
