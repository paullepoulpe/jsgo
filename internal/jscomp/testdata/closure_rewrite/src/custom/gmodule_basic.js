goog.module('app.widgets.Button');
'use strict';
const Component = goog.require('app.ui.Component');
const {assert, fail} = goog.require('goog.asserts');
const events = goog.requireType('goog.events');
let Lazy = goog.forwardDeclare('app.Lazy');
var dom = goog.require('goog.dom'), str = goog.require('goog.string');

/** @typedef {{label: string}} */
let Options;

/** @const {number} */
const DEFAULT_WIDTH = 100;

class Button extends Component {
  /** @param {!Options} opts */
  constructor(opts) {
    super();
    /** @type {!events.EventTarget|undefined} */
    this.target = undefined;
    assert(opts.label);
    this.width = DEFAULT_WIDTH;
  }
  /** @return {!Lazy} */
  lazy() { Lazy = goog.module.get('app.Lazy'); return new Lazy(); }
}

function helper(/** !Button */ b) { return dom.getElement(str.trim(b.label)); }

exports = Button;
exports.Options = Options;
exports.helper = helper;
