// Port of com/google/javascript/refactoring/ApplySuggestedFixes.java @ closure-compiler v20261028
// Port of com/google/javascript/refactoring/CodeReplacement.java (with the AutoValue-generated
// AutoValue_CodeReplacement.java) @ closure-compiler v20261028

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
	"cmp"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func joinCodeReplacements(parts []*CodeReplacement) string {
	var sb strings.Builder
	for i, part := range parts {
		if i > 0 {
			sb.WriteString("\n\n")
		}
		sb.WriteString(part.String())
	}
	return sb.String()
}

// ApplyToFiles applies the provided suggested fixes to the files listed in the
// fixes. The fixes can be provided in any order, but they may not have any
// overlapping modifications for the same file.
func ApplyToFiles(fixes []*SuggestedFix) error {
	var filenames []string
	for _, fix := range fixes {
		for filename := range fix.GetReplacements().KeySet() {
			if !slices.Contains(filenames, filename) {
				filenames = append(filenames, filename)
			}
		}
	}

	filenameToCodeMap := make(map[string]string, len(filenames))
	for _, filename := range filenames {
		b, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		filenameToCodeMap[filename] = jutil.StandardCharsets_UTF_8.Decode(b)
	}

	newCode := ApplyToCode(fixes, filenameToCodeMap)
	for _, filename := range filenames {
		if value, ok := newCode[filename]; ok {
			if err := os.WriteFile(filename, jutil.StringGetBytesUTF8(value), 0o666); err != nil {
				return err
			}
		}
	}
	return nil
}

// ApplyAllChoicesToCode applies all possible options from each
// SuggestedFixAlternative to the provided code and returns the new code. This
// only makes sense if all the SuggestedFixAlternatives come from the same
// checker, i.e. they offer the same number of choices and the same index
// corresponds to similar fixes. The fileNameToCodeMap must contain all the
// files that the provided fixes apply to. The fixes can be provided in any
// order, but they may not have any overlapping modifications for the same file.
// This function will return new code only for the files that have been modified.
func ApplyAllChoicesToCode(
	fixChoices []*SuggestedFix, fileNameToCodeMap map[string]string) []map[string]string {
	if len(fixChoices) == 0 {
		return []map[string]string{{}}
	}
	alternativeCount := len(fixChoices[0].GetAlternatives())
	allMatch := true
	for _, f := range fixChoices[1:] {
		if len(f.GetAlternatives()) != alternativeCount {
			allMatch = false
			break
		}
	}
	jutil.CheckArgumentMsg(
		allMatch,
		"All SuggestedFixAlternatives must offer an equal number of choices for this "+
			"utility to make sense")
	result := make([]map[string]string, 0, alternativeCount)
	for i := 0; i < alternativeCount; i++ {
		result = append(result, ApplyChoicesToCode(fixChoices, i, fileNameToCodeMap))
	}
	return result
}

// ApplyChoicesToCode applies the choiceIndex-th alternative from each fix in
// fixChoices to fileNameToCodeMap and returns the modified files.
func ApplyChoicesToCode(
	fixChoices []*SuggestedFix,
	choiceIndex int,
	fileNameToCodeMap map[string]string) map[string]string {
	chosenFixes := make([]*SuggestedFix, 0, len(fixChoices))
	for _, choices := range fixChoices {
		chosenFixes = append(chosenFixes, choices.GetAlternatives()[choiceIndex])
	}
	return ApplyToCode(chosenFixes, fileNameToCodeMap)
}

// ApplyToCode applies the provided suggested fixes to the provided code and
// returns the new code, ignoring alternative fixes. The filenameToCodeMap must
// contain all the files that the provided fixes apply to. The fixes can be
// provided in any order, but they may not have any overlapping modifications
// for the same file. This function will return new code only for the files that
// have been modified.
func ApplyToCode(
	fixes []*SuggestedFix, filenameToCodeMap map[string]string) map[string]string {
	var rm replacementMap
	for _, fix := range fixes {
		rm.putIfNoOverlap(fix)
	}
	newCodeMap := make(map[string]string, len(rm.keys))
	for _, filename := range rm.keys {
		code, ok := filenameToCodeMap[filename]
		if !ok {
			panic(jutil.NewIllegalArgumentExceptionMsg("filenameToCodeMap missing code for file: " + filename))
		}
		newCodeMap[filename] = ApplyCodeReplacements(rm.byFile[filename], code)
	}
	return newCodeMap
}

// ApplyCodeReplacements applies the provided set of code replacements to the
// code and returns the transformed code. The code replacements may not have
// any overlap.
func ApplyCodeReplacements(
	replacements []*CodeReplacement, code string) string {
	sortedReplacements := sortedUniqueReplacements(replacements)
	if ContainsOverlaps(sortedReplacements) {
		jutil.CheckArgumentMsg1(
			false,
			"Found overlap between code replacements:\n%s",
			joinCodeReplacements(sortedReplacements))
	}

	isASCII := jutil.IsASCII(code)
	codeLen := len(code)
	var codeUnits []uint16
	if !isASCII {
		codeUnits = jutil.StringToUTF16(code)
		codeLen = len(codeUnits)
	}
	validateReplacementRanges(sortedReplacements, codeLen)

	if isASCII {
		var sb strings.Builder
		sb.Grow(codeLen)
		lastIndex := 0
		for _, replacement := range sortedReplacements {
			sb.WriteString(code[lastIndex:replacement.GetStartPosition()])
			sb.WriteString(replacement.GetNewContent())
			lastIndex = replacement.GetEndPosition()
		}
		if lastIndex < codeLen {
			sb.WriteString(code[lastIndex:])
		}
		return sb.String()
	}

	out := make([]uint16, 0, codeLen)
	lastIndex := 0
	for _, replacement := range sortedReplacements {
		out = append(out, codeUnits[lastIndex:replacement.GetStartPosition()]...)
		out = append(out, jutil.StringToUTF16(replacement.GetNewContent())...)
		lastIndex = replacement.GetEndPosition()
	}
	if lastIndex < codeLen {
		out = append(out, codeUnits[lastIndex:]...)
	}
	return jutil.StringFromUTF16(out)
}

func validateReplacementRanges(sortedReplacements []*CodeReplacement, codeLen int) {
	lastIndex := 0
	for _, replacement := range sortedReplacements {
		start := replacement.GetStartPosition()
		if lastIndex < 0 || lastIndex > start || start > codeLen {
			panic(jutil.NewIndexOutOfBoundsExceptionForRange(lastIndex, start, codeLen))
		}
		lastIndex = replacement.GetEndPosition()
	}
	if lastIndex < 0 || lastIndex > codeLen {
		panic(jutil.NewIndexOutOfBoundsExceptionForRange(lastIndex, codeLen, codeLen))
	}
}

func sortedUniqueReplacements(replacements []*CodeReplacement) []*CodeReplacement {
	sorted := slices.Clone(replacements)
	slices.SortStableFunc(sorted, func(a, b *CodeReplacement) int { return a.CompareTo(b) })
	return slices.CompactFunc(sorted, func(a, b *CodeReplacement) bool { return a.CompareTo(b) == 0 })
}

// ContainsOverlaps reports whether sorted CodeReplacements have any overlap.
// Equivalent CodeReplacements should be deduplicated before calling this.
func ContainsOverlaps(replacements []*CodeReplacement) bool {
	start := -1
	for _, replacement := range replacements {
		if replacement.GetStartPosition() < start {
			return true
		}
		start = replacement.GetEndPosition()
	}
	return false
}

type replacementMap struct {
	byFile map[string][]*CodeReplacement
	keys   []string
}

func (r *replacementMap) putIfNoOverlap(fix *SuggestedFix) {
	if !r.canPut(fix) {
		return
	}
	if r.byFile == nil {
		r.byFile = make(map[string][]*CodeReplacement)
	}
	for k, v := range fix.GetReplacements().Entries() {
		cur, ok := r.byFile[k]
		if !ok {
			r.keys = append(r.keys, k)
		}
		if !slices.ContainsFunc(cur, func(existing *CodeReplacement) bool { return existing.Equals(v) }) {
			r.byFile[k] = append(cur, v)
		}
	}
}

func (r *replacementMap) canPut(fix *SuggestedFix) bool {
	reps := fix.GetReplacements()
	for filename := range reps.KeySet() {
		existing := r.byFile[filename]
		incoming := reps.GetSlice(filename)
		combined := make([]*CodeReplacement, 0, len(existing)+len(incoming))
		combined = append(combined, existing...)
		combined = append(combined, incoming...)
		if ContainsOverlaps(sortedUniqueReplacements(combined)) {
			return false
		}
	}
	return true
}

// CodeReplacement represents a set of changes to make to the code.
type CodeReplacement struct {
	startPosition int
	length        int
	newContent    string
	sortKey       string
}

// Java: CodeReplacement.create(int, int, String)
func CodeReplacement_Create(startPosition int, length int, newContent string) *CodeReplacement {
	return CodeReplacement_CreateWithSortKey(startPosition, length, newContent, "")
}

// Java: CodeReplacement.create(int, int, String, String)
func CodeReplacement_CreateWithSortKey(startPosition int, length int, newContent string, sortKey string) *CodeReplacement {
	return &CodeReplacement{
		startPosition: startPosition,
		length:        length,
		newContent:    newContent,
		sortKey:       sortKey,
	}
}

// Java: CodeReplacement.getStartPosition()
func (r *CodeReplacement) GetStartPosition() int { return r.startPosition }

// Java: CodeReplacement.getLength()
func (r *CodeReplacement) GetLength() int { return r.length }

// Java: CodeReplacement.getEndPosition()
func (r *CodeReplacement) GetEndPosition() int { return r.startPosition + r.length }

// Java: CodeReplacement.getNewContent()
func (r *CodeReplacement) GetNewContent() string { return r.newContent }

// Java: CodeReplacement.getSortKey()
func (r *CodeReplacement) GetSortKey() string { return r.sortKey }

// Java: CodeReplacement.compareTo(CodeReplacement)
func (r *CodeReplacement) CompareTo(x *CodeReplacement) int {
	if c := cmp.Compare(r.startPosition, x.startPosition); c != 0 {
		return c
	}
	if c := cmp.Compare(r.length, x.length); c != 0 {
		return c
	}
	if c := jutil.StringCompareTo(r.sortKey, x.sortKey); c != 0 {
		return c
	}
	return jutil.StringCompareTo(r.newContent, x.newContent)
}

// Java: AutoValue_CodeReplacement.toString()
func (r *CodeReplacement) String() string {
	return "CodeReplacement{" +
		"startPosition=" + strconv.Itoa(r.startPosition) + ", " +
		"length=" + strconv.Itoa(r.length) + ", " +
		"newContent=" + r.newContent + ", " +
		"sortKey=" + r.sortKey +
		"}"
}

// Java: AutoValue_CodeReplacement.equals(Object)
func (r *CodeReplacement) Equals(o any) bool {
	that, ok := o.(*CodeReplacement)
	if !ok || that == nil {
		return false
	}
	if that == r {
		return true
	}
	return r.startPosition == that.startPosition &&
		r.length == that.length &&
		r.newContent == that.newContent &&
		r.sortKey == that.sortKey
}

// Java: AutoValue_CodeReplacement.hashCode()
func (r *CodeReplacement) HashCode() int32 {
	var h int32 = 1
	h = h*1000003 ^ int32(r.startPosition)
	h = h*1000003 ^ int32(r.length)
	h = h*1000003 ^ jutil.StringHashCode(r.newContent)
	h = h*1000003 ^ jutil.StringHashCode(r.sortKey)
	return h
}
