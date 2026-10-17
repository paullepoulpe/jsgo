// Port of com/google/javascript/jscomp/SimpleRegion.java @ closure-compiler v20261015

/*
 * Copyright 2007 The Closure Compiler Authors.
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

// SimpleRegion is a simple region.
type SimpleRegion struct {
	beginningLineNumber int
	endingLineNumber    int
	source              string
}

// Java: SimpleRegion(int, int, String)
func NewSimpleRegion(beginningLineNumber int, endingLineNumber int, source string) *SimpleRegion {
	return &SimpleRegion{beginningLineNumber: beginningLineNumber, endingLineNumber: endingLineNumber, source: source}
}

// Java: SimpleRegion.getBeginningLineNumber()
func (r *SimpleRegion) GetBeginningLineNumber() int { return r.beginningLineNumber }

// Java: SimpleRegion.getEndingLineNumber()
func (r *SimpleRegion) GetEndingLineNumber() int { return r.endingLineNumber }

// Java: SimpleRegion.getSourceExcerpt()
func (r *SimpleRegion) GetSourceExcerpt() string { return r.source }
