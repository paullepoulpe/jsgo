// Port of com/google/javascript/refactoring/testing/SuggestedFixes.java @ closure-compiler v20261021
// (src/com/google/javascript/refactoring/testing/SuggestedFixes.java, a JUnit/Truth test utility)

/*
 * Copyright 2014 The Closure Compiler Authors.
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

package refactoring

import (
	"slices"

	"github.com/paullepoulpe/jsgo/internal/testjunit"
)

// SuggestedFixes holds utilities for testing SuggestedFix-related code.
//
// Java: public final class refactoring.testing.SuggestedFixes
type SuggestedFixes struct {
	_ byte // Go-only: language-adapter — distinct identity
}

// Java: SuggestedFixes()
func newSuggestedFixes() *SuggestedFixes {
	return &SuggestedFixes{}
}

// Java: SuggestedFixes.assertReplacement(SuggestedFix, CodeReplacement)
func SuggestedFixes_AssertReplacement(fix *SuggestedFix, expectedReplacement *CodeReplacement) {
	SuggestedFixes_AssertReplacements(fix, []*CodeReplacement{expectedReplacement})
}

// Java: SuggestedFixes.assertReplacements(SuggestedFix, Set<CodeReplacement>)
func SuggestedFixes_AssertReplacements(
	fix *SuggestedFix, expectedReplacements []*CodeReplacement) {
	replacementMap := fix.GetReplacements()
	// MultimapSubject.hasSize(1): Multimap.size(), the number of entries.
	testjunit.AssertThat(replacementMap.Size()).IsEqualTo(1)
	replacements := replacementMap.GetSlice("test")
	testjunit.AssertThat(replacements).HasSize(len(expectedReplacements))
	testjunit.AssertThat(replacements).ContainsExactlyElementsIn(expectedReplacements)
}

// Java: SuggestedFixes.assertChanges(SuggestedFix, String, String, String)
func SuggestedFixes_AssertChanges(
	fix *SuggestedFix, externs string, originalCode string, expectedCode string) {
	SuggestedFixes_AssertChangesOfList([]*SuggestedFix{fix}, externs, originalCode, expectedCode)
}

// Java: SuggestedFixes.assertChanges(List<SuggestedFix>, String, String, String)
func SuggestedFixes_AssertChangesOfList(
	fixes []*SuggestedFix, externs string, originalCode string, expectedCode string) {
	testjunit.AssertThat(fixes).IsNotEmpty()
	codeMap := map[string]string{"test": originalCode}
	res := ApplySuggestedFixes_ApplySuggestedFixesToCode(slices.Values(fixes), codeMap)
	var newCode any
	if v, ok := res["test"]; ok {
		newCode = v
	}
	testjunit.AssertThat(newCode).IsEqualTo(expectedCode)
}
