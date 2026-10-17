/** @desc A simple message. */
var MSG_A = goog.getMsg('Hello world');

/** @desc With placeholders. */
var MSG_B = goog.getMsg('Hello {$name}, you have {$count} new messages', {name: 'Ann', count: 3});

/** @desc Concatenated text. */
var MSG_C = goog.getMsg('first part ' + 'second part');

/** @desc Unused message. */
var MSG_D = goog.getMsg('Goodbye');

/** @desc Html. */
var MSG_E = goog.getMsg('a < b &amp; c', {}, {html: true});

/** @desc Unescape. */
var MSG_F = goog.getMsg('x &lt; y &amp; z &quot;q&quot;', {}, {unescapeHtmlEntities: true});

/** @desc Fallback base. */
var MSG_G = goog.getMsg('fallback {$x}', {x: 1});
/** @desc Fallback alt. */
var MSG_H = goog.getMsg('other {$x}', {x: 2});
var chosen = goog.getMsgWithFallback(MSG_G, MSG_H);

/**
 * @desc Alternate.
 * @alternateMessageId 1984
 */
var MSG_I = goog.getMsg('Hi {$who}!', {who: chosen});

/** @desc External. */
var MSG_EXTERNAL_12345 = goog.getMsg('External {$v}', {v: 'value'});

/** @desc ICU select. */
var MSG_J = goog.getMsg(
    '{GENDER,select,female{She {NAME}}male{He {NAME}}other{They {NAME}}}');

/** @desc On an object. */
var ns = {};
ns.MSG_K = goog.getMsg('Property message {$p}', {p: MSG_A});

/** @desc Original code placeholders. */
var MSG_L = goog.getMsg(
    'Click {$startLink}here{$endLink}', {startLink: '<a href="#">', endLink: '</a>'},
    {original_code: {startLink: '<a>', endLink: '</a>'}, example: {startLink: '<a href=x>'}});

function getMessages(n) {
  /** @desc In a function. */
  const MSG_M = goog.getMsg('Item {$num}', {num: n});
  return MSG_M + MSG_B + MSG_C + MSG_E + MSG_F + MSG_I + MSG_EXTERNAL_12345 + MSG_J + ns.MSG_K +
      MSG_L;
}
getMessages(5);
