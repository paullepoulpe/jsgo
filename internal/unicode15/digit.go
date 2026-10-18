// Derived from src/unicode/digit.go (Go 1.24.7)
// Licence: LICENSES/BSD-3-Clause-Go.txt.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unicode15

// IsDigit reports whether the rune is a decimal digit.
func IsDigit(r rune) bool {
	if r <= MaxLatin1 {
		return '0' <= r && r <= '9'
	}
	return isExcludingLatin(Digit, r)
}
