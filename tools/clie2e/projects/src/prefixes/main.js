/**
 * @fileoverview clie2e projects suite: imports through prefixes that
 * --browser_resolver_prefix_replacements maps to source directories.
 */
import {square, cube} from 'util/math.js';
import {shout} from '~text/text.js';
import * as math from './util/math.js';

export const answer = shout(String(square(7) + cube(2)));
export const viaNamespace = math.square(3);
