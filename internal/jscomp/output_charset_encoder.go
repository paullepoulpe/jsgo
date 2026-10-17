// Port of com/google/javascript/jscomp/OutputCharsetEncoder.java @ closure-compiler v20260415

/*
 * Copyright 2015 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// OutputCharsetEncoder is the output charset encoder for CodeGenerator that
// delegates to a CharsetEncoder.
type OutputCharsetEncoder struct {
	// Go-only: language-adapter — a CharsetEncoder is only used for
	// canEncode(char), which depends on the charset alone, so the charset
	// stands in for it (nil is Java's null encoder).
	encoder *jutil.Charset
}

// Java: OutputCharsetEncoder(Charset)
func NewOutputCharsetEncoder(outputCharset *jutil.Charset) *OutputCharsetEncoder {
	e := &OutputCharsetEncoder{}
	if outputCharset == nil || outputCharset == jutil.StandardCharsets_US_ASCII {
		// If we want our default (pretending to be UTF-8, but escaping
		// anything outside of straight ASCII), then don't use the encoder,
		// but just special-case the code. This keeps the normal path
		// through the code identical to how it's been for years.
		e.encoder = nil
	} else {
		e.encoder = outputCharset
	}
	return e
}

// Java: OutputCharsetEncoder.canEncode(char)
func (e *OutputCharsetEncoder) CanEncode(c uint16) bool {
	return e.encoder != nil && e.encoder.CanEncodeChar(c)
}
