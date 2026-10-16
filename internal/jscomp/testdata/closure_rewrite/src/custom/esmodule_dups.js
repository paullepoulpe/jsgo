export var x = 1;
export {x};
export default 1;
export default function f() {}
export * from './esmodule_dups.js';
export {missing} from './esmodule_dups.js';
import {nothere} from './esmodule_dups.js';
export {nothere};
