/**
 * @fileoverview Benchmark sample app (tools/bench/README.md): persistent note
 * storage. A goog.provide file, so that the app itself mixes both Closure
 * dependency styles, like the library it is built on.
 */

goog.provide('benchapp.Note');
goog.provide('benchapp.NoteStore');

goog.require('goog.array');
goog.require('goog.events.EventTarget');
goog.require('goog.storage.Storage');
goog.require('goog.storage.mechanism.HTML5LocalStorage');
goog.require('goog.string');

/**
 * @typedef {{id: string, title: string, html: string, due: ?number}}
 */
benchapp.Note;

/**
 * Keeps the notes in local storage and tells listeners when they change.
 * @param {string} prefix Storage key prefix.
 * @constructor
 * @extends {goog.events.EventTarget}
 * @final
 */
benchapp.NoteStore = function(prefix) {
  benchapp.NoteStore.base(this, 'constructor');

  /** @private @const {string} */
  this.prefix_ = prefix;

  /** @private @const {!goog.storage.mechanism.HTML5LocalStorage} */
  this.mechanism_ = new goog.storage.mechanism.HTML5LocalStorage();

  /** @private @const {?goog.storage.Storage} */
  this.storage_ = this.mechanism_.isAvailable() ?
      new goog.storage.Storage(this.mechanism_) :
      null;

  /** @private {!Array<!benchapp.Note>} */
  this.notes_ = this.load_();
};
goog.inherits(benchapp.NoteStore, goog.events.EventTarget);

/** @const {string} The event type dispatched after every change. */
benchapp.NoteStore.CHANGE = 'change';

/**
 * @return {!Array<!benchapp.Note>}
 * @private
 */
benchapp.NoteStore.prototype.load_ = function() {
  if (!this.storage_) {
    return [];
  }
  var stored = this.storage_.get(this.prefix_ + 'notes');
  return Array.isArray(stored) ? /** @type {!Array<!benchapp.Note>} */ (stored) : [];
};

/** @private */
benchapp.NoteStore.prototype.save_ = function() {
  if (this.storage_) {
    this.storage_.set(this.prefix_ + 'notes', this.notes_);
  }
  this.dispatchEvent(benchapp.NoteStore.CHANGE);
};

/** @return {!Array<!benchapp.Note>} The notes, oldest first. */
benchapp.NoteStore.prototype.getAll = function() {
  return goog.array.clone(this.notes_);
};

/**
 * @param {string} id
 * @return {?benchapp.Note}
 */
benchapp.NoteStore.prototype.get = function(id) {
  return goog.array.find(this.notes_, function(n) {
    return n.id == id;
  });
};

/**
 * Adds a note with a fresh id.
 * @param {string} title
 * @return {!benchapp.Note} The new note.
 */
benchapp.NoteStore.prototype.add = function(title) {
  var note = {
    id: goog.string.getRandomString(),
    title: goog.string.trim(title) || 'Untitled',
    html: '',
    due: null
  };
  this.notes_.push(note);
  this.save_();
  return note;
};

/**
 * @param {string} id
 * @param {function(!benchapp.Note)} change Edits the note in place.
 */
benchapp.NoteStore.prototype.update = function(id, change) {
  var note = this.get(id);
  if (note) {
    change(note);
    this.save_();
  }
};

/**
 * @param {string} id
 * @return {boolean} Whether a note was removed.
 */
benchapp.NoteStore.prototype.remove = function(id) {
  var removed = goog.array.removeIf(this.notes_, function(n) {
    return n.id == id;
  });
  if (removed) {
    this.save_();
  }
  return removed;
};

/**
 * @param {string} query
 * @return {!Array<!benchapp.Note>} The notes whose title contains query.
 */
benchapp.NoteStore.prototype.search = function(query) {
  var q = query.toLowerCase();
  return goog.array.filter(this.notes_, function(n) {
    return goog.string.contains(n.title.toLowerCase(), q);
  });
};
