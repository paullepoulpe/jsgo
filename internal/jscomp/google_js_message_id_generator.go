// Port of com/google/javascript/jscomp/GoogleJsMessageIdGenerator.java @ closure-compiler v20260415

/*
 * Copyright 2008 The Closure Compiler Authors.
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
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// GoogleJsMessageIdGenerator is an IdGenerator designed to play nicely
// with Google's Translation systems. Each message is scoped to a project
// id, so that it does not conflict with other messages at Google.
//
// Just as reminder what key type used in different formats:
//  1. XMB - id. We export using this format.
//  2. XTB - id. Internal, result of translation.
//  3. XLB - name. External, use it if we need to share translation with
//     third part.
//  4. PROPERTIES - name.
//
// See http://cldr.unicode.org/development/development-process/design-proposals/xmb
type GoogleJsMessageIdGenerator struct {
	projectId *string
}

// NewGoogleJsMessageIdGenerator creates an instance.
//
// projectId is a TC project name (e.g. "MyProject").
//
// Java: GoogleJsMessageIdGenerator(String)
func NewGoogleJsMessageIdGenerator(projectId *string) *GoogleJsMessageIdGenerator {
	return &GoogleJsMessageIdGenerator{projectId: projectId}
}

// Java: GoogleJsMessageIdGenerator.generateId(String, List<Part>)
func (g *GoogleJsMessageIdGenerator) GenerateId(meaning *string, messageParts []JsMessage_Part) string {
	jutil.CheckState(meaning != nil)

	sb := jutil.NewStringBuilder()
	for _, part := range messageParts {
		if part.IsPlaceholder() {
			sb.Append(part.GetCanonicalPlaceholderName())
		} else {
			sb.Append(part.GetString())
		}
	}
	tcValue := sb.String()

	projectScopedMeaning := ""
	if !(g.projectId == nil || *g.projectId == "") {
		projectScopedMeaning = *g.projectId + ": "
	}
	projectScopedMeaning += *meaning

	return strconv.FormatInt(
		googleJsMessageIdGenerator_MessageId_generateId(tcValue, &projectScopedMeaning), 10)
}

// googleJsMessageIdGenerator_FP is the 64-bit fingerprint support.
//
// Forked from the guava-internal library.
type googleJsMessageIdGenerator_FP struct{}

// googleJsMessageIdGenerator_FP_fingerprintRange generates fingerprint of
// "byte[start,limit-1]".
//
// Java: GoogleJsMessageIdGenerator.FP.fingerprint(byte[], int, int)
func googleJsMessageIdGenerator_FP_fingerprintRange(str []byte, start int, limit int) int64 {
	hi := googleJsMessageIdGenerator_FP_hash32(str, start, limit, 0)
	lo := googleJsMessageIdGenerator_FP_hash32(str, start, limit, 102072)
	if (hi == 0) && (lo == 0 || lo == 1) {
		// Turn 0/1 into another fingerprint
		hi ^= 0x130f9bef
		lo ^= -0x6b5f56d8 // 0x94a0a928
	}
	return (int64(hi) << 32) | (int64(lo) & 0xffffffff)
}

// googleJsMessageIdGenerator_FP_fingerprint generates fingerprint of
// "str". Equivalent to UTF-encoding "str" into bytes and then
// fingerprinting those bytes.
//
// Java: GoogleJsMessageIdGenerator.FP.fingerprint(String)
func googleJsMessageIdGenerator_FP_fingerprint(str string) int64 {
	tmp := jutil.StringGetBytesUTF8(str)
	return googleJsMessageIdGenerator_FP_fingerprintRange(tmp, 0, len(tmp))
}

// Java: GoogleJsMessageIdGenerator.FP.hash32(byte[], int, int, int)
func googleJsMessageIdGenerator_FP_hash32(str []byte, start int, limit int, c int32) int32 {
	var a int32 = -0x61c88647 // 0x9e3779b9
	var b int32 = -0x61c88647 // 0x9e3779b9
	at := func(i int) int32 { return int32(str[i]) & 0xff }
	var i int
	for i = start; i+12 <= limit; i += 12 {
		a += (at(i+0) << 0) |
			(at(i+1) << 8) |
			(at(i+2) << 16) |
			(at(i+3) << 24)
		b += (at(i+4) << 0) |
			(at(i+5) << 8) |
			(at(i+6) << 16) |
			(at(i+7) << 24)
		c += (at(i+8) << 0) |
			(at(i+9) << 8) | (at(i+10) << 16) |
			(at(i+11) << 24)

		// Mix
		a -= b
		a -= c
		a ^= jutil.Ushr32(c, 13)
		b -= c
		b -= a
		b ^= a << 8
		c -= a
		c -= b
		c ^= jutil.Ushr32(b, 13)
		a -= b
		a -= c
		a ^= jutil.Ushr32(c, 12)
		b -= c
		b -= a
		b ^= a << 16
		c -= a
		c -= b
		c ^= jutil.Ushr32(b, 5)
		a -= b
		a -= c
		a ^= jutil.Ushr32(c, 3)
		b -= c
		b -= a
		b ^= a << 10
		c -= a
		c -= b
		c ^= jutil.Ushr32(b, 15)
	}

	c += int32(limit - start)
	tmp := limit - i
	if tmp == 11 {
		c += at(i+10) << 24
	}
	if tmp >= 10 {
		c += at(i+9) << 16
	}
	if tmp >= 9 {
		c += at(i+8) << 8
		// the first byte of c is reserved for the length
	}
	if tmp >= 8 {
		b += at(i+7) << 24
	}
	if tmp >= 7 {
		b += at(i+6) << 16
	}
	if tmp >= 6 {
		b += at(i+5) << 8
	}
	if tmp >= 5 {
		b += at(i + 4)
	}
	if tmp >= 4 {
		a += at(i+3) << 24
	}
	if tmp >= 3 {
		a += at(i+2) << 16
	}
	if tmp >= 2 {
		a += at(i+1) << 8
	}
	if tmp >= 1 {
		a += at(i + 0)
		// case 0 : nothing left to add
	}

	// Mix
	a -= b
	a -= c
	a ^= jutil.Ushr32(c, 13)
	b -= c
	b -= a
	b ^= a << 8
	c -= a
	c -= b
	c ^= jutil.Ushr32(b, 13)
	a -= b
	a -= c
	a ^= jutil.Ushr32(c, 12)
	b -= c
	b -= a
	b ^= a << 16
	c -= a
	c -= b
	c ^= jutil.Ushr32(b, 5)
	a -= b
	a -= c
	a ^= jutil.Ushr32(c, 3)
	b -= c
	b -= a
	b ^= a << 10
	c -= a
	c -= b
	c ^= jutil.Ushr32(b, 15)
	return c
}

// googleJsMessageIdGenerator_MessageId_generateId generates fingerprint
// for an English message using the FP package. This supersedes the message
// id generation using C fingerprint functions and JNI. This is slower than
// the C implementation (we're talking about microseconds here) but it
// avoids using JNI and shared libraries.
//
// Forked from the i18n library.
//
// Java: GoogleJsMessageIdGenerator.MessageId.generateId(String, String)
func googleJsMessageIdGenerator_MessageId_generateId(message string, meaning *string) int64 {
	fp := googleJsMessageIdGenerator_FP_fingerprint(message)
	if nil != meaning && len(*meaning) > 0 {
		// combine the fingerprints of message and meaning
		fp2 := googleJsMessageIdGenerator_FP_fingerprint(*meaning)
		var carry int64 = 0
		if fp < 0 {
			carry = 1
		}
		fp = fp2 + (fp << 1) + carry
	}
	// To avoid negative ids we strip the high-order bit
	return fp & 0x7fffffffffffffff
}
