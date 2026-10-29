// Port of com/google/javascript/refactoring/ApplySuggestedFixes.java @ closure-compiler v20261021
// Port of com/google/javascript/refactoring/CodeReplacement.java (with the AutoValue-generated
// AutoValue_CodeReplacement.java) @ closure-compiler v20261021

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
	"iter"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// applySuggestedFixes_DOUBLE_LINE_JOINER_Join is
// DOUBLE_LINE_JOINER.join(parts) with DOUBLE_LINE_JOINER =
// Joiner.on("\n\n"): the parts' toString() separated by "\n\n".
func applySuggestedFixes_DOUBLE_LINE_JOINER_Join(parts []*CodeReplacement) string {
	var sb strings.Builder
	for i, part := range parts {
		if i > 0 {
			sb.WriteString("\n\n")
		}
		sb.WriteString(part.String())
	}
	return sb.String()
}

// ApplySuggestedFixes_ApplySuggestedFixesToFiles applies the provided set of
// suggested fixes to the files listed in the suggested fixes. The fixes can
// be provided in any order, but they may not have any overlapping
// modifications for the same file. It throws IOException.
//
// Java: ApplySuggestedFixes.applySuggestedFixesToFiles(Iterable<SuggestedFix>)
func ApplySuggestedFixes_ApplySuggestedFixesToFiles(fixes iter.Seq[*SuggestedFix]) {
	fixList := slices.Collect(fixes)
	var filenames []string
	for _, fix := range fixList {
		for filename := range fix.GetReplacements().KeySet() {
			if !slices.Contains(filenames, filename) {
				filenames = append(filenames, filename)
			}
		}
	}

	filenameToCodeMap := make(map[string]string, len(filenames))
	for _, filename := range filenames {
		filenameToCodeMap[filename] = jutil.Files_AsCharSource(jutil.NewFile(filename), jutil.StandardCharsets_UTF_8).Read()
	}

	newCode := ApplySuggestedFixes_ApplySuggestedFixesToCode(slices.Values(fixList), filenameToCodeMap)
	for _, filename := range filenames {
		if value, ok := newCode[filename]; ok {
			if err := os.WriteFile(filename, jutil.StringGetBytesUTF8(value), 0o666); err != nil {
				panic(jutil.FileNotFoundExceptionFromError(err, filename))
			}
		}
	}
}

// ApplySuggestedFixes_ApplyAllSuggestedFixChoicesToCode applies all possible
// options from each SuggestedFixAlternative to the provided code and returns
// the new code. This only makes sense if all the SuggestedFixAlternatives
// come from the same checker, i.e. they offer the same number of choices and
// the same index corresponds to similar fixes. The fileNameToCodeMap must
// contain all the files that the provided fixes apply to. The fixes can be
// provided in any order, but they may not have any overlapping
// modifications for the same file. This function will return new code only
// for the files that have been modified.
//
// Java: ApplySuggestedFixes.applyAllSuggestedFixChoicesToCode(Iterable<SuggestedFix>, Map<String, String>)
func ApplySuggestedFixes_ApplyAllSuggestedFixChoicesToCode(
	fixChoices iter.Seq[*SuggestedFix], fileNameToCodeMap map[string]string) []map[string]string {
	choiceList := slices.Collect(fixChoices)
	if len(choiceList) == 0 {
		return []map[string]string{{}}
	}
	alternativeCount := len(choiceList[0].GetAlternatives())
	allMatch := true
	for _, f := range choiceList[1:] {
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
		result = append(result, ApplySuggestedFixes_ApplySuggestedFixChoicesToCode(slices.Values(choiceList), i, fileNameToCodeMap))
	}
	return result
}

// Java: ApplySuggestedFixes.applySuggestedFixChoicesToCode(Iterable<SuggestedFix>, int, Map<String, String>)
func ApplySuggestedFixes_ApplySuggestedFixChoicesToCode(
	fixChoices iter.Seq[*SuggestedFix],
	choiceIndex int,
	fileNameToCodeMap map[string]string) map[string]string {
	var chosenFixes []*SuggestedFix
	for choices := range fixChoices {
		chosenFixes = append(chosenFixes, choices.GetAlternatives()[choiceIndex])
	}
	return ApplySuggestedFixes_ApplySuggestedFixesToCode(slices.Values(chosenFixes), fileNameToCodeMap)
}

// ApplySuggestedFixes_ApplySuggestedFixesToCode applies the provided set of
// suggested fixes to the provided code and returns the new code, ignoring
// alternative fixes. The filenameToCodeMap must contain all the files that
// the provided fixes apply to. The fixes can be provided in any order, but
// they may not have any overlapping modifications for the same file. This
// function will return new code only for the files that have been modified.
//
// Java: ApplySuggestedFixes.applySuggestedFixesToCode(Iterable<SuggestedFix>, Map<String, String>)
func ApplySuggestedFixes_ApplySuggestedFixesToCode(
	fixes iter.Seq[*SuggestedFix], filenameToCodeMap map[string]string) map[string]string {
	map_ := NewApplySuggestedFixes_ReplacementMap()
	for fix := range fixes {
		map_.PutIfNoOverlap(fix)
	}
	newCodeMap := make(map[string]string, len(map_.keys))
	for filename, replacements := range map_.EntrySet() {
		code, ok := filenameToCodeMap[filename]
		if !ok {
			panic(jutil.NewIllegalArgumentExceptionMsg("filenameToCodeMap missing code for file: " + filename))
		}
		if _, dup := newCodeMap[filename]; dup {
			panic(jutil.NewIllegalArgumentExceptionMsg("Multiple entries with same key: " + filename))
		}
		newCodeMap[filename] = ApplySuggestedFixes_ApplyCodeReplacements(replacements, code)
	}
	return newCodeMap
}

// ApplySuggestedFixes_ApplyCodeReplacements applies the provided set of
// code replacements to the code and returns the transformed code. The code
// replacements may not have any overlap.
//
// Java: ApplySuggestedFixes.applyCodeReplacements(Collection<CodeReplacement>, String)
func ApplySuggestedFixes_ApplyCodeReplacements(
	replacements []*CodeReplacement, code string) string {
	sortedReplacements :=
		applySuggestedFixes_ImmutableSortedSetCopyOf(replacements)
	if ApplySuggestedFixes_ContainsOverlaps(sortedReplacements) {
		jutil.CheckArgumentMsg1(
			false,
			"Found overlap between code replacements:\n%s",
			applySuggestedFixes_DOUBLE_LINE_JOINER_Join(sortedReplacements))
	}

	if jutil.IsASCII(code) {
		var sb strings.Builder
		sb.Grow(len(code))
		lastIndex := 0
		for _, replacement := range sortedReplacements {
			start := replacement.GetStartPosition()
			end := replacement.GetEndPosition()
			if lastIndex < 0 || lastIndex > start || start > len(code) {
				panic(jutil.NewIndexOutOfBoundsExceptionForRange(lastIndex, start, len(code)))
			}
			sb.WriteString(code[lastIndex:start])
			sb.WriteString(replacement.GetNewContent())
			lastIndex = end
		}
		if lastIndex <= len(code) {
			if lastIndex < 0 {
				panic(jutil.NewIndexOutOfBoundsExceptionForRange(lastIndex, len(code), len(code)))
			}
			sb.WriteString(code[lastIndex:])
		}
		return sb.String()
	}

	codeUnits := jutil.StringToUTF16(code)
	out := make([]uint16, 0, len(codeUnits))
	lastIndex := 0
	for _, replacement := range sortedReplacements {
		start := replacement.GetStartPosition()
		if lastIndex < 0 || lastIndex > start || start > len(codeUnits) {
			panic(jutil.NewIndexOutOfBoundsExceptionForRange(lastIndex, start, len(codeUnits)))
		}
		out = append(out, codeUnits[lastIndex:start]...)
		out = append(out, jutil.StringToUTF16(replacement.GetNewContent())...)
		lastIndex = replacement.GetEndPosition()
	}
	if lastIndex <= len(codeUnits) {
		if lastIndex < 0 {
			panic(jutil.NewIndexOutOfBoundsExceptionForRange(lastIndex, len(codeUnits), len(codeUnits)))
		}
		out = append(out, codeUnits[lastIndex:]...)
	}
	return jutil.StringFromUTF16(out)
}

// applySuggestedFixes_ImmutableSortedSetCopyOf is
// ImmutableSortedSet.copyOf(replacements) and
// ImmutableSortedSet.naturalOrder().addAll(...).build(): the elements sorted
// by CodeReplacement.compareTo, keeping the first of the elements that
// compare equal.
//
// Go-only: language-adapter — Guava ImmutableSortedSet.
func applySuggestedFixes_ImmutableSortedSetCopyOf(replacements []*CodeReplacement) []*CodeReplacement {
	sorted := slices.Clone(replacements)
	slices.SortStableFunc(sorted, func(a, b *CodeReplacement) int { return a.CompareTo(b) })
	return slices.CompactFunc(sorted, func(a, b *CodeReplacement) bool { return a.CompareTo(b) == 0 })
}

// ApplySuggestedFixes_ContainsOverlaps checks whether the
// `CodeReplacement`s have any overlap.
//
// We intentinoally accept a set here beceuse we want to supprt/ignore
// equivalent `CodeReplacement`s. It's safe if those overlap as long as we
// only apply one of them.
//
// Java: ApplySuggestedFixes.containsOverlaps(ImmutableSortedSet<CodeReplacement>)
func ApplySuggestedFixes_ContainsOverlaps(replacements []*CodeReplacement) bool {
	start := -1
	for _, replacement := range replacements {
		if replacement.GetStartPosition() < start {
			return true
		}
		start = replacement.GetEndPosition()
	}
	return false
}

// Java: private static class ApplySuggestedFixes.ReplacementMap
type ApplySuggestedFixes_ReplacementMap struct {
	map_ map[string][]*CodeReplacement
	keys []string
}

// Java: ApplySuggestedFixes.ReplacementMap()
func NewApplySuggestedFixes_ReplacementMap() *ApplySuggestedFixes_ReplacementMap {
	return &ApplySuggestedFixes_ReplacementMap{map_: make(map[string][]*CodeReplacement)}
}

// Java: ApplySuggestedFixes.ReplacementMap.putIfNoOverlap(SuggestedFix)
func (r *ApplySuggestedFixes_ReplacementMap) PutIfNoOverlap(fix *SuggestedFix) {
	if r.CanPut(fix) {
		for k, v := range fix.GetReplacements().Entries() {
			cur, ok := r.map_[k]
			if !ok {
				r.keys = append(r.keys, k)
			}
			if !slices.ContainsFunc(cur, func(existing *CodeReplacement) bool { return existing.Equals(v) }) {
				r.map_[k] = append(cur, v)
			}
		}
	}
}

// Java: ApplySuggestedFixes.ReplacementMap.canPut(SuggestedFix)
func (r *ApplySuggestedFixes_ReplacementMap) CanPut(fix *SuggestedFix) bool {
	reps := fix.GetReplacements()
	for filename := range reps.KeySet() {
		existing := r.map_[filename]
		incoming := reps.GetSlice(filename)
		combined := make([]*CodeReplacement, 0, len(existing)+len(incoming))
		combined = append(combined, existing...)
		combined = append(combined, incoming...)
		replacements := applySuggestedFixes_ImmutableSortedSetCopyOf(combined)
		if ApplySuggestedFixes_ContainsOverlaps(replacements) {
			return false
		}
	}
	return true
}

// Java: ApplySuggestedFixes.ReplacementMap.entrySet()
func (r *ApplySuggestedFixes_ReplacementMap) EntrySet() iter.Seq2[string, []*CodeReplacement] {
	return func(yield func(string, []*CodeReplacement) bool) {
		for _, k := range r.keys {
			if !yield(k, r.map_[k]) {
				return
			}
		}
	}
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
