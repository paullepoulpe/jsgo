/**
 * @fileoverview clie2e projects suite: an ES module entry point that imports
 * lodash's CommonJS modules by bare specifier (NODE resolution), a local
 * CommonJS module and a local ES module.
 */
import debounce from 'lodash/debounce';
import throttle from 'lodash/throttle';
import counterModule from './counter.js';
import {formatCount, LABEL} from './format.js';

const counter = counterModule.createCounter({step: 2});

const report = debounce(() => {
  console.log(formatCount(LABEL, counter.value()));
}, 50);

export const tick = throttle(() => {
  counter.increment();
  report();
}, 10);

export default function run(times) {
  for (let i = 0; i < times; i++) {
    tick();
  }
  return counter.value();
}
