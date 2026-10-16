import Default, {a, b as bee} from './other.js';
import * as ns from './ns.js';
import {Foo} from 'goog:app.Foo';
import Bar from 'goog:app.Bar';
import * as starGoog from 'goog:app.Star';
import './side_effect.js';
export {a, bee as b2};
export * from './reexport.js';
export {x as y, z} from './from.js';
export default class Klass {}
export function fn() { a2 = 1; }
export let a2 = 0, [c, {d}] = [];
export const e = ns.e;
let mutated = 0;
export {mutated};
function mutate() { mutated++; }
export {Foo, Bar};
