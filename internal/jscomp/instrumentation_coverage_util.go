// Port of com/google/javascript/jscomp/instrumentation/CoverageUtil.java @ closure-compiler v20261021

/*
 * Copyright 2009 The Closure Compiler Authors.
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

// CoverageUtil contains all general purpose static methods that are used
// by coverage instrumentation related classes - not limited to them,
// though.
type CoverageUtil struct{}

// CoverageUtil_CreateIdentifierFromText returns a string with all
// non-alphanumeric characters in the given string replaced with
// underscores. This is to create a valid identifier based on the given
// text.
//
// Java: CoverageUtil.createIdentifierFromText(String)
func CoverageUtil_CreateIdentifierFromText(inputText string) string {
	// Java: inputText.replaceAll("[^\\p{Alnum}]", "_"). \p{Alnum} is ASCII
	// [a-zA-Z0-9]; the regex matches each code point, so a surrogate pair
	// is replaced by one '_' and a lone surrogate by one '_'.
	units := jutil.StringToUTF16(inputText)
	out := make([]uint16, 0, len(units))
	for i := 0; i < len(units); i++ {
		c := units[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			out = append(out, c)
			continue
		}
		if c >= 0xD800 && c <= 0xDBFF && i+1 < len(units) && units[i+1] >= 0xDC00 && units[i+1] <= 0xDFFF {
			i++
		}
		out = append(out, '_')
	}
	return jutil.StringFromUTF16(out)
}
