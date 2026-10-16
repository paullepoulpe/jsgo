/**
 * @fileoverview Benchmark sample app (tools/bench/README.md): the entry point
 * (--entry_point=goog:benchapp.main).
 */
goog.module('benchapp.main');

const NoteStore = goog.require('benchapp.NoteStore');
const NotesView = goog.require('benchapp.NotesView');
const dom = goog.require('goog.dom');
const events = goog.require('goog.events');

/**
 * Starts the app in the element with id "app".
 * @param {string} syncUrl
 * @return {!NotesView}
 */
function main(syncUrl) {
  const store = new NoteStore('benchapp.');
  if (!store.getAll().length) {
    store.add('Welcome');
  }
  const view = new NotesView(store, syncUrl);
  view.render(dom.getRequiredElement('app'));
  return view;
}

events.listenOnce(window, events.EventType.LOAD, () => {
  main('/notes/sync');
});

exports = {main};
