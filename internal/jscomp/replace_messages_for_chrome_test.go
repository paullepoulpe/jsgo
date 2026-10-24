// Port of com/google/javascript/jscomp/ReplaceMessagesForChromeTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/ReplaceMessagesForChromeTest.java)

/*
 * Copyright 2012 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import "testing"

// ReplaceMessagesForChromeTest tests ReplaceMessagesForChrome.
//
// Java: public final class ReplaceMessagesForChromeTest extends CompilerTestCase
type ReplaceMessagesForChromeTest struct {
	CompilerTestCaseBase
}

// Java: ReplaceMessagesForChromeTest()
func NewReplaceMessagesForChromeTest() *ReplaceMessagesForChromeTest {
	c := &ReplaceMessagesForChromeTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestReplaceMessagesForChromeTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewReplaceMessagesForChromeTest() })
}

// Java: ReplaceMessagesForChromeTest.getProcessor(Compiler)
func (c *ReplaceMessagesForChromeTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewReplaceMessagesForChrome(compiler, NewGoogleJsMessageIdGenerator(nil))
}

// Java: ReplaceMessagesForChromeTest.setUp()
func (c *ReplaceMessagesForChromeTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableTypeCheck()
	c.ReplaceTypesWithColors()
	c.EnableTypeInfoValidation()
	c.IgnoreWarnings(DiagnosticGroups_MISSING_PROPERTIES)
}

// Java: ReplaceMessagesForChromeTest.testReplaceSimpleMessage()
func (c *ReplaceMessagesForChromeTest) TestReplaceSimpleMessage() {
	c.Test(
		`/** @desc A simple message. */

var MSG_A = goog.getMsg('Hello world');
`,
		`/** @desc A simple message. */

var MSG_A=chrome.i18n.getMessage('8660696502365331902');
`)

	c.Test(
		`/** @desc A message attached to an object. */

foo.bar.MSG_B = goog.getMsg('Goodbye world');
`,
		`/** @desc A message attached to an object. */

foo.bar.MSG_B=chrome.i18n.getMessage('2356086230621084760');
`)
}

// Java: ReplaceMessagesForChromeTest.testReplaceSinglePlaceholder()
func (c *ReplaceMessagesForChromeTest) TestReplaceSinglePlaceholder() {
	c.Test(
		`/** @desc A message with one placeholder. */

var MSG_C = goog.getMsg('Hello, {$name}', {name: 'Tyler'});
`,
		`/** @desc A message with one placeholder. */

var MSG_C=chrome.i18n.getMessage('4985325380591528435', ['Tyler']);
`)
}

// Java: ReplaceMessagesForChromeTest.testReplaceTwoPlaceholders()
func (c *ReplaceMessagesForChromeTest) TestReplaceTwoPlaceholders() {
	c.Test(
		`/** @desc A message with two placeholders. */

var MSG_D = goog.getMsg('{$greeting}, {$name}',
{greeting: 'Hi', name: 'Tyler'});
`,
		`/** @desc A message with two placeholders. */

var MSG_D=chrome.i18n.getMessage('3605047247574980322',
['Hi', 'Tyler']);
`)

	c.Test(
		`/** @desc A message with two placeholders, but their order is

 * reversed in the object literal. (Shouldn't make a difference.)

 */

var MSG_E = goog.getMsg('{$greeting}, {$name}!',
{name: 'Tyler', greeting: 'Hi'});
`,
		`/** @desc A message with two placeholders, but their order is

 * reversed in the object literal. (Shouldn't make a difference.)

 */

var MSG_E=chrome.i18n.getMessage('691522386483664339',
['Hi', 'Tyler']);
`)
}

// Java: ReplaceMessagesForChromeTest.testReplacePlaceholderMissingValue()
func (c *ReplaceMessagesForChromeTest) TestReplacePlaceholderMissingValue() {
	c.TestError(
		`/** @desc A message with two placeholders, but one is missing. */
var MSG_F = goog.getMsg('{$greeting}, {$name}!', {name: 'Tyler'});
`,
		JsMessageVisitor_MESSAGE_TREE_MALFORMED)
}

// Java: ReplaceMessagesForChromeTest.testReplaceTwoPlaceholdersNonAlphaOrder()
func (c *ReplaceMessagesForChromeTest) TestReplaceTwoPlaceholdersNonAlphaOrder() {
	c.Test(
		`/** @desc A message with two placeholders not in order .*/

var MSG_G = goog.getMsg('{$name}: {$greeting}',
{greeting: 'Salutations', name: 'Tyler'});
`,
		`/** @desc A message with two placeholders not in order .*/

var MSG_G=chrome.i18n.getMessage('7437383242562773138',
['Salutations', 'Tyler']);
`)
}

// Java: ReplaceMessagesForChromeTest.testReplaceSinglePlaceholderComputedProp()
func (c *ReplaceMessagesForChromeTest) TestReplaceSinglePlaceholderComputedProp() {
	c.TestError(
		`/** @desc A message with one placeholder. */

var MSG_H = goog.getMsg('Hello, {$name}', {['name']: 'Tyler'});
`,
		JsMessageVisitor_MESSAGE_TREE_MALFORMED)
}

// Java: ReplaceMessagesForChromeTest.testReplaceSimpleMessageWithLet()
func (c *ReplaceMessagesForChromeTest) TestReplaceSimpleMessageWithLet() {
	c.Test(
		`/** @desc A simple message. */

let MSG_I = goog.getMsg('Hello world');
`,
		`/** @desc A simple message. */

let MSG_I = chrome.i18n.getMessage('987871171253827787');
`)
}

// Java: ReplaceMessagesForChromeTest.testReplaceSimpleMessageWithConst()
func (c *ReplaceMessagesForChromeTest) TestReplaceSimpleMessageWithConst() {
	c.Test(
		`/** @desc A simple message. */

const MSG_J = goog.getMsg('Hello world');
`,
		`/** @desc A simple message. */

const MSG_J =chrome.i18n.getMessage('3477894568604521782');
`)
}

// Java: ReplaceMessagesForChromeTest.testReplaceTemplatedMessage()
func (c *ReplaceMessagesForChromeTest) TestReplaceTemplatedMessage() {
	c.TestError(
		"const greeting = '{$greeting}'\n/** @desc A simple message */\n\nvar MSG_K = goog.getMsg(`${greeting}, Tyler`, {greeting: 'Hello'});\n",
		JsMessageVisitor_MESSAGE_TREE_MALFORMED)
}

// Java: ReplaceMessagesForChromeTest.testReplaceExternalMessage()
func (c *ReplaceMessagesForChromeTest) TestReplaceExternalMessage() {
	c.Test(
		`/** @desc A message that was extracted with SoyMsgExtractor. */

var MSG_EXTERNAL_1357902468 = goog.getMsg('Hello world');
`,
		`/** @desc A message that was extracted with SoyMsgExtractor. */

var MSG_EXTERNAL_1357902468 = chrome.i18n.getMessage('1357902468');
`)
}

// Java: ReplaceMessagesForChromeTest.testReplaceMessageWithDuplicatePlaceholders()
func (c *ReplaceMessagesForChromeTest) TestReplaceMessageWithDuplicatePlaceholders() {
	c.Test(
		`/** @desc A message that contains two instances of the same placeholder. */

var MSG_EXTERNAL_987654321 = goog.getMsg(
'{$startDiv_1}You are signed in as{$endDiv}{$img}{$startDiv_2}{$name}{$endDiv}',
{'startDiv_1': '<div>',
'endDiv': '</div>',
'img': '<img src="http://example.com/photo.png">',
'startDiv_2': '<div class="name">',
'name': name});
`,
		`/** @desc A message that contains two instances of the same placeholder. */

var MSG_EXTERNAL_987654321 = chrome.i18n.getMessage('987654321',
[
'</div>',  // endDiv, only included once even though it appears twice in the message.
'<img src="http://example.com/photo.png">',  // img
name,  // name
'<div>',  // startDiv_1
'<div class="name">' // startDiv_2
]);
`)
}

// Java: ReplaceMessagesForChromeTest.testReplaceMessageWithHtml()
func (c *ReplaceMessagesForChromeTest) TestReplaceMessageWithHtml() {
	c.Test(
		`/** @desc A message with one placeholder. */

var MSG_C = goog.getMsg('Hello, {$name}', {name: 'Tyler'}, {html: true});
`,
		`/** @desc A message with one placeholder. */

var MSG_C=chrome.i18n.getMessage.apply(null, ['4985325380591528435', ['Tyler']].concat((/Chrome\/(\d+)/.exec(navigator.userAgent) || [])[1] >= 79 ? [{escapeLt: true}] : []));
`)

	c.Test(
		`/** @desc A simple message. */

var MSG_A = goog.getMsg('Hello world', {}, {html: true});
`,
		`/** @desc A simple message. */

var MSG_A=chrome.i18n.getMessage.apply(null, ['8660696502365331902', []].concat((/Chrome\/(\d+)/.exec(navigator.userAgent) || [])[1] >= 79 ? [{escapeLt: true}] : []));
`)

	c.Test(
		`/** @desc A simple message. */

var MSG_A = goog.getMsg('Hello world', {}, {html: false});
`,
		`/** @desc A simple message. */

var MSG_A=chrome.i18n.getMessage('8660696502365331902');
`)

	c.Test(
		`/** @desc A simple message. */

var MSG_A = goog.getMsg('Hello world', {}, {});
`,
		`/** @desc A simple message. */

var MSG_A=chrome.i18n.getMessage('8660696502365331902');
`)
}
