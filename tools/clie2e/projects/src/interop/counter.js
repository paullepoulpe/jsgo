/**
 * @fileoverview clie2e projects suite: a CommonJS module that requires
 * lodash's CommonJS modules.
 */
'use strict';

var isObject = require('lodash/isObject');
var assign = require('lodash/assign');

var defaults = {start: 0, step: 1};

function createCounter(options) {
  var settings = assign({}, defaults, isObject(options) ? options : {});
  var value = settings.start;
  return {
    increment: function() {
      value += settings.step;
      return value;
    },
    value: function() {
      return value;
    }
  };
}

exports.createCounter = createCounter;
exports.defaults = defaults;
