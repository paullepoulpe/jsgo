// Port of com/google/javascript/refactoring/Scanner.java @ closure-compiler v20261021

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

package jscomp

import "github.com/paullepoulpe/jsgo/internal/rhino"

// Scanner is the class that developers should implement to perform a JsFlume
// refactoring.
//
// Java: public abstract class Scanner implements Serializable
type Scanner interface {
	// Matches returns true if the given node and node traversal should match
	// for this particular scanner. Typically this function uses the Matcher
	// class or predefined matchers from Matchers to match against the Node and
	// NodeMetadata.
	//
	// If this function returns true, a Match for this node will be passed to
	// ProcessMatch and all matches will be passed to ProcessAllMatches at the
	// end of the traversal.
	//
	// Java: Scanner.matches(Node, NodeMetadata)
	Matches(node *rhino.Node, t *NodeMetadata) bool

	// Java: Scanner.processMatch(Match)
	ProcessMatch(match *Match) []*SuggestedFix

	// Java: Scanner.processAllMatches(Collection<Match>)
	ProcessAllMatches(matches []*Match) []*SuggestedFix

	// GetScannerBase gives access to the base class.
	GetScannerBase() *ScannerBase
}

// ScannerBase holds Scanner's method bodies; subclasses embed it and define
// Matches.
type ScannerBase struct {
	_ byte // Go-only: language-adapter — distinct identity
}

// GetScannerBase gives access to the base class.
func (s *ScannerBase) GetScannerBase() *ScannerBase { return s }

// ProcessMatch processes one Match at a time. There is no order guaranteed for
// when this function will be called with the Match.
//
// match is the Match from the node and traversal for any match that Matches
// returned true for. It returns the list of SuggestedFix classes that will be
// applied to the source files at the end of the run to create the
// refactoring CL.
//
// Java: Scanner.processMatch(Match)
func (s *ScannerBase) ProcessMatch(match *Match) []*SuggestedFix {
	return []*SuggestedFix{}
}

// ProcessAllMatches processes every given match at one time. This function
// can be used when the refactoring needs the information from the entire run
// to perform the refactoring, such as moving functions around.
//
// matches are all the Match matches that were collected when the Matches
// function returned true. It returns the list of SuggestedFix classes that
// will be applied to the source files at the end of the run to create the
// refactoring CL.
//
// Java: Scanner.processAllMatches(Collection<Match>)
func (s *ScannerBase) ProcessAllMatches(matches []*Match) []*SuggestedFix {
	return []*SuggestedFix{}
}
