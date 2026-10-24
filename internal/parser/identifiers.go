// Port of com/google/javascript/jscomp/parsing/parser/Identifiers.java @ closure-compiler v20261021

/*
 * Copyright 2011 The Closure Compiler Authors.
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

package parser

import "github.com/paullepoulpe/jsgo/internal/jutil"

// Identifiers: JS identifier parsing utilities.
//
// Java: final class Identifiers (static methods only; the private
// constructor has no Go counterpart).

const (
	asciiFlagIdentStart uint8 = 1 << iota
	asciiFlagIdentPart
)

var asciiIdentFlags = func() [128]uint8 {
	var t [128]uint8
	for c := uint8('A'); c <= 'Z'; c++ {
		t[c] = asciiFlagIdentStart | asciiFlagIdentPart
	}
	for c := uint8('a'); c <= 'z'; c++ {
		t[c] = asciiFlagIdentStart | asciiFlagIdentPart
	}
	for c := uint8('0'); c <= '9'; c++ {
		t[c] = asciiFlagIdentPart
	}
	t['_'] = asciiFlagIdentStart | asciiFlagIdentPart
	t['$'] = asciiFlagIdentStart | asciiFlagIdentPart
	return t
}()

// Identifiers_IsIdentifierStart reports whether ch can start an identifier.
//
// Java: Identifiers.isIdentifierStart(char)
func Identifiers_IsIdentifierStart(ch uint16) bool {
	// Most code is written in pure ASCII, so create a fast path here.
	if ch <= 127 {
		return asciiIdentFlags[ch]&asciiFlagIdentStart != 0
	}

	// Handle non-ASCII characters.
	// TODO(tjgq): This should include all characters with the ID_Start property.
	return jutil.CharacterIsLetter(ch)
}

// Identifiers_IsIdentifierPart reports whether ch can be part of an
// identifier.
//
// Java: Identifiers.isIdentifierPart(char)
func Identifiers_IsIdentifierPart(ch uint16) bool {
	// Most code is written in pure ASCII, so create a fast path here.
	if ch <= 127 {
		return asciiIdentFlags[ch]&asciiFlagIdentPart != 0
	}

	// Handle non-ASCII characters.
	// TODO(tjgq): This should include all characters with the ID_Continue property, plus
	// Zero Width Non-Joiner and Zero Width Joiner.
	return Identifiers_IsIdentifierStart(ch) || jutil.CharacterIsDigit(ch)
}
