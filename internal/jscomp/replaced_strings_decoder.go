// Port of com/google/javascript/jscomp/ReplacedStringsDecoder.java @ closure-compiler v20261021

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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// ReplacedStringsDecoder is a decoder for strings encoded by the
// ReplaceStrings JS compiler pass. This class is immutable.
type ReplacedStringsDecoder struct {
	originalToNewNameMap *jutil.LinkedHashMap[string, string]
}

// ReplacedStringsDecoder_ARGUMENT_PLACE_HOLDER is the default place holder
// for replaced string arguments. Note: We're not allowing this to be
// changed here because it's present in the replacement strings in the
// js_storage protocol buffer.
const ReplacedStringsDecoder_ARGUMENT_PLACE_HOLDER = "`"

// ReplacedStringsDecoder_NULL_DECODER is a null decoder that does no
// mapping.
var ReplacedStringsDecoder_NULL_DECODER = NewReplacedStringsDecoder(VariableMap_FromMap(
	jutil.NewLinkedHashMap[string, string]()))

// Java: ReplacedStringsDecoder(VariableMap)
func NewReplacedStringsDecoder(variableMap *VariableMap) *ReplacedStringsDecoder {
	// VariableMap is not an immutable type, so we extract the map instead of
	// directly using it.
	return &ReplacedStringsDecoder{originalToNewNameMap: variableMap.GetOriginalNameToNewNameMap()}
}

// Decode decodes an encoded string from the JS Compiler ReplaceStrings
// pass.
//
// An original string with args might look like this:
//
//	Error('Some ' + arg1 + ' error ' + arg2 + ' message.');
//
// Which gets replaced with:
//
//	Error('key' + '`' + arg1 + '`' + arg2);
//
// Where ` is the argument place holder. The replacement mapping would be:
//
//	key → 'Some ` error ` message.'
//
// Where key is some arbitrary replacement string. An encoded string, with
// args, from the client will look like:
//
//	'key`arg1`arg2'
//
// Returns the decoded string, or the encoded string if it fails to decode.
//
// See com.google.javascript.jscomp.ReplaceStrings.
//
// Java: ReplacedStringsDecoder.decode(String)
func (d *ReplacedStringsDecoder) Decode(encodedStr string) string {
	// Java: encodedStr.split(ARGUMENT_PLACE_HOLDER, -1)
	suppliedBits := strings.Split(encodedStr, ReplacedStringsDecoder_ARGUMENT_PLACE_HOLDER)
	originalStr, ok := d.originalToNewNameMap.Lookup(suppliedBits[0])
	if !ok {
		return encodedStr // Failed to decode.
	}
	originalBits := strings.Split(originalStr, ReplacedStringsDecoder_ARGUMENT_PLACE_HOLDER)
	sb := jutil.NewStringBuilderFrom(originalBits[0])
	for i := 1; i < max(len(originalBits), len(suppliedBits)); i++ {
		// Replace missing bits with "-". Shouldn't happen except that we
		// aren't escaping the replacement token at the moment.
		if i < len(suppliedBits) {
			sb.Append(suppliedBits[i])
		} else {
			sb.Append("-")
		}
		if i < len(originalBits) {
			sb.Append(originalBits[i])
		} else {
			sb.Append("-")
		}
	}
	return sb.String()
}
