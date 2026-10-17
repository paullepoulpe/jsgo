/**
 * @fileoverview clie2e projects suite: a dynamic import of a sibling
 * module (--allow_dynamic_import, --dynamic_import_alias).
 */
import {shout} from './util/text.js';

/** @return {!Promise<string>} */
export async function load() {
  const math = await import('./util/math.js');
  return shout(String(math.cube(3)));
}
