export var ea = 1, eb = 2;
export let el = 3;
export const EC = 4, ED = 5;
export function ef() {}
export class EK {}
export var {ex, ey: [ez]} = obj;
export default function() {}
var ACONST = 1;
var notConst = 2;
var $jscomp$CONST = 3;
/** @const */ var annotated = 4;
const realConst = 5;
let shadowed = ACONST;
function useConsts() { return ACONST + notConst + annotated + realConst + shadowed; }
export {notConst as nc};
