// Port of com/google/javascript/jscomp/ReplacedStringsDecoderTest.java @ closure-compiler v20261015

/*
 * Copyright 2016 The Closure Compiler Authors.
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

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func verifyDecoder(t *testing.T, original, replacement string) {
	t.Helper()
	verifyDecoderWithArgs(t, original, replacement, "")
}

func verifyDecoderWithArgs(t *testing.T, original, replacement, args string) {
	t.Helper()
	encoding := "Xy" // an arbitrary encoding for the test
	m := jutil.NewLinkedHashMap[string, string]()
	m.Put(encoding, replacement)
	variableMap := VariableMap_FromMap(m)
	decoder := NewReplacedStringsDecoder(variableMap)
	assertEqual(t, "decode", decoder.Decode(encoding+args), original)
}

func TestReplacedStringsDecoder_SimpleDecoding(t *testing.T) {
	verifyDecoder(t, "A Message.", "A Message.")
}

func TestReplacedStringsDecoder_DecodingWithArgs(t *testing.T) {
	verifyDecoderWithArgs(t, "A foo B bar C.", "A ` B ` C.", "`foo`bar")
}

func TestReplacedStringsDecoder_DecodingWithExcessSlots(t *testing.T) {
	verifyDecoderWithArgs(t, "A foo B bar C - D -.", "A ` B ` C ` D `.", "`foo`bar")
}

func TestReplacedStringsDecoder_DecodingWithExcessArgs(t *testing.T) {
	verifyDecoderWithArgs(t, "A foo B bar C.baz-bam-", "A ` B ` C.", "`foo`bar`baz`bam")
}

func TestReplacedStringsDecoder_TrailingArg(t *testing.T) {
	verifyDecoderWithArgs(t, "foo bar", "foo `", "`bar")
}

func TestReplacedStringsDecoder_EmptyArgs(t *testing.T) {
	verifyDecoderWithArgs(t, "foo bar", "foo `bar``", "```")
}

func TestReplacedStringsDecoder_NullDecoder(t *testing.T) {
	nullDecoder := ReplacedStringsDecoder_NULL_DECODER
	assertEqual(t, "decode", nullDecoder.Decode("foo"), "foo")
}
