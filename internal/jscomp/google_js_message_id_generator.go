// Port of com/google/javascript/jscomp/GoogleJsMessageIdGenerator.java @ closure-compiler v20261021

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
	"encoding/binary"
	"strconv"
	"strings"

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
	projectId string
}

// NewGoogleJsMessageIdGenerator creates an instance.
//
// projectId is a TC project name (e.g. "MyProject").
//
// Java: GoogleJsMessageIdGenerator(String)
func NewGoogleJsMessageIdGenerator(projectId *string) *GoogleJsMessageIdGenerator {
	var pid string
	if projectId != nil {
		pid = *projectId
	}
	return &GoogleJsMessageIdGenerator{projectId: pid}
}

// Java: GoogleJsMessageIdGenerator.generateId(String, List<Part>)
func (g *GoogleJsMessageIdGenerator) GenerateId(meaning *string, messageParts []JsMessage_Part) string {
	jutil.CheckState(meaning != nil)
	return g.generateId(*meaning, messageParts)
}

func (g *GoogleJsMessageIdGenerator) generateId(meaning string, messageParts []JsMessage_Part) string {
	var tcValue string
	if len(messageParts) == 1 && !messageParts[0].IsPlaceholder() {
		tcValue = messageParts[0].GetString()
	} else if len(messageParts) > 0 {
		var sb strings.Builder
		for _, part := range messageParts {
			if part.IsPlaceholder() {
				sb.WriteString(part.GetCanonicalPlaceholderName())
			} else {
				sb.WriteString(part.GetString())
			}
		}
		tcValue = sb.String()
	}

	projectScopedMeaning := meaning
	if g.projectId != "" {
		projectScopedMeaning = g.projectId + ": " + projectScopedMeaning
	}

	return strconv.FormatInt(
		googleJsMessageIdGenerator_MessageId_generateId(tcValue, projectScopedMeaning), 10)
}

// googleJsMessageIdGenerator_FP_fingerprintRange generates fingerprint of
// "byte[start,limit-1]".
//
// Java: GoogleJsMessageIdGenerator.FP.fingerprint(byte[], int, int)
func googleJsMessageIdGenerator_FP_fingerprintRange(str []byte, start int, limit int) int64 {
	b := str[start:limit]
	hi := googleJsMessageIdGenerator_FP_hash32(b, 0)
	lo := googleJsMessageIdGenerator_FP_hash32(b, 102072)
	if hi == 0 && (lo == 0 || lo == 1) {
		// Turn 0/1 into another fingerprint
		hi ^= 0x130f9bef
		lo ^= 0x94a0a928
	}
	return int64((uint64(hi) << 32) | uint64(lo))
}

// googleJsMessageIdGenerator_FP_fingerprint generates fingerprint of
// "str". Equivalent to UTF-encoding "str" into bytes and then
// fingerprinting those bytes.
//
// Java: GoogleJsMessageIdGenerator.FP.fingerprint(String)
func googleJsMessageIdGenerator_FP_fingerprint(str string) int64 {
	var buf [128]byte
	tmp := jutil.AppendStringBytesUTF8(buf[:0], str)
	return googleJsMessageIdGenerator_FP_fingerprintRange(tmp, 0, len(tmp))
}

// Java: GoogleJsMessageIdGenerator.FP.hash32(byte[], int, int, int)
func googleJsMessageIdGenerator_FP_hash32(str []byte, c uint32) uint32 {
	var a uint32 = 0x9e3779b9
	var b uint32 = 0x9e3779b9
	i := 0
	for ; i+12 <= len(str); i += 12 {
		a += binary.LittleEndian.Uint32(str[i:])
		b += binary.LittleEndian.Uint32(str[i+4:])
		c += binary.LittleEndian.Uint32(str[i+8:])
		a, b, c = jenkinsMix32(a, b, c)
	}

	c += uint32(len(str))
	tail := str[i:]
	switch len(tail) {
	case 11:
		c += uint32(tail[10]) << 24
		fallthrough
	case 10:
		c += uint32(tail[9]) << 16
		fallthrough
	case 9:
		c += uint32(tail[8]) << 8
		// the first byte of c is reserved for the length
		fallthrough
	case 8:
		b += binary.LittleEndian.Uint32(tail[4:])
		a += binary.LittleEndian.Uint32(tail)
	case 7:
		b += uint32(tail[6]) << 16
		fallthrough
	case 6:
		b += uint32(tail[5]) << 8
		fallthrough
	case 5:
		b += uint32(tail[4])
		fallthrough
	case 4:
		a += binary.LittleEndian.Uint32(tail)
	case 3:
		a += uint32(tail[2]) << 16
		fallthrough
	case 2:
		a += uint32(tail[1]) << 8
		fallthrough
	case 1:
		a += uint32(tail[0])
		// case 0 : nothing left to add
	}

	_, _, c = jenkinsMix32(a, b, c)
	return c
}

func jenkinsMix32(a, b, c uint32) (uint32, uint32, uint32) {
	a -= b
	a -= c
	a ^= c >> 13
	b -= c
	b -= a
	b ^= a << 8
	c -= a
	c -= b
	c ^= b >> 13
	a -= b
	a -= c
	a ^= c >> 12
	b -= c
	b -= a
	b ^= a << 16
	c -= a
	c -= b
	c ^= b >> 5
	a -= b
	a -= c
	a ^= c >> 3
	b -= c
	b -= a
	b ^= a << 10
	c -= a
	c -= b
	c ^= b >> 15
	return a, b, c
}

// googleJsMessageIdGenerator_MessageId_generateId generates fingerprint
// for an English message using the FP package.
//
// Java: GoogleJsMessageIdGenerator.MessageId.generateId(String, String)
func googleJsMessageIdGenerator_MessageId_generateId(message string, meaning string) int64 {
	fp := uint64(googleJsMessageIdGenerator_FP_fingerprint(message))
	if len(meaning) > 0 {
		// combine the fingerprints of message and meaning
		fp2 := uint64(googleJsMessageIdGenerator_FP_fingerprint(meaning))
		fp = fp2 + (fp << 1) + (fp >> 63)
	}
	// To avoid negative ids we strip the high-order bit
	return int64(fp & 0x7fffffffffffffff)
}
