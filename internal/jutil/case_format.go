// Port of com/google/common/base/CaseFormat.java (Guava 33.4.6-jre)

/*
 * Copyright (C) 2006 The Guava Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 */

package jutil

import "strings"

// Go-only: language-adapter — Guava's com.google.common.base.CaseFormat
// conversions not covered by CaseFormat_UpperUnderscoreToLowerCamel. Both
// conversions only inspect ASCII boundaries and ASCII letter case, so they
// operate on UTF-8/WTF-8 bytes directly in a single pass.

// CaseFormat_LowerCamelToUpperUnderscore is Guava's
// CaseFormat.LOWER_CAMEL.to(CaseFormat.UPPER_UNDERSCORE, s): words split at
// each ASCII uppercase letter after index 0, are upper-cased
// (Ascii.toUpperCase), and are joined with '_'.
func CaseFormat_LowerCamelToUpperUnderscore(s string) string {
	extra := 0
	needsChange := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			needsChange = true
		} else if i > 0 && c >= 'A' && c <= 'Z' {
			extra++
			needsChange = true
		}
	}
	if !needsChange {
		return s
	}
	var sb strings.Builder
	sb.Grow(len(s) + extra)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			sb.WriteByte(c - ('a' - 'A'))
		} else if i > 0 && c >= 'A' && c <= 'Z' {
			sb.WriteByte('_')
			sb.WriteByte(c)
		} else {
			sb.WriteByte(c)
		}
	}
	return sb.String()
}

// CaseFormat_LowerHyphenToUpperCamel is Guava's
// CaseFormat.LOWER_HYPHEN.to(CaseFormat.UPPER_CAMEL, s): words end at each
// '-', each is normalized by UPPER_CAMEL (first char Ascii.toUpperCase, the
// rest Ascii.toLowerCase), and they are joined without a separator.
func CaseFormat_LowerHyphenToUpperCamel(s string) string {
	hyphens := 0
	needsChange := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '-' {
			hyphens++
			needsChange = true
		} else if i == 0 && c >= 'a' && c <= 'z' {
			needsChange = true
		} else if i > 0 && c >= 'A' && c <= 'Z' {
			needsChange = true
		}
	}
	if !needsChange {
		return s
	}
	var sb strings.Builder
	sb.Grow(len(s) - hyphens)
	atWordStart := true
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '-' {
			atWordStart = true
			continue
		}
		if atWordStart {
			atWordStart = false
			if c >= 'a' && c <= 'z' {
				c -= 'a' - 'A'
			}
		} else if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		sb.WriteByte(c)
	}
	return sb.String()
}
