/** Real generated J2CL JRE collection operations, with an observable export. */
goog.module('enterprise.j2cl');

const HashMap = goog.require('java.util.HashMap');
const ArrayList = goog.require('java.util.ArrayList');

/** @return {string} */
function exerciseCollections() {
  const map = HashMap.$create__();
  map.put('alpha', 7);
  map.put('beta', 11);
  map.put('alpha', 13);
  const list = ArrayList.$create__();
  list.add(map.get('alpha'));
  list.add(map.get('beta'));
  if (map.size() !== 2 || list.size() !== 2 || list.getAtIndex(0) !== 13 ||
      list.getAtIndex(1) !== 11 || map.remove('alpha') !== 13 || map.size() !== 1) {
    throw new Error('J2CL generated collection semantics failed');
  }
  return 'J2CL:13,11;remaining=1';
}

globalThis['jsgoEnterpriseJ2cl'] = exerciseCollections;
