// Port of com/google/javascript/jscomp/CheckRegExp.java @ closure-compiler v20261021

/*
 * Copyright 2010 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/regex"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

var checkRegExpID = uint32(strpool.Intern("RegExp"))

var checkRegExpVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_NAME)
	s.Add(rhino.Token_REGEXP)
	return s
}()

var _ nodeTraversal_visitTokens = (*CheckRegExp)(nil)

func (*CheckRegExp) visitTokens() *ast.TokenSet { return &checkRegExpVisitTokens }

// CheckRegExp looks for references to the global RegExp object that would
// cause regular expressions to be unoptimizable, and checks that regular
// expressions are syntactically valid.
type CheckRegExp struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler                           AbstractCompiler
	reportErrors                       bool
	globalRegExpPropertiesUsed         bool
	assumeAllGlobalRegexpUsagesVisible bool
}

// Java: CheckRegExp.REGEXP_REFERENCE
var CheckRegExp_REGEXP_REFERENCE = DiagnosticType_Warning("JSC_REGEXP_REFERENCE",
	"References to the global RegExp object prevents "+
		"optimization of regular expressions.")

// Java: CheckRegExp.MALFORMED_REGEXP
var CheckRegExp_MALFORMED_REGEXP = DiagnosticType_Warning(
	"JSC_MALFORMED_REGEXP",
	"Malformed Regular Expression: {0}")

// Java: CheckRegExp.REGEXP_PROPERTY_SKIPLIST
func checkRegExp_isSkippedProperty(prop string) bool {
	switch prop {
	case "$1", "$2", "$3", "$4", "$5", "$6", "$7", "$8", "$9", "$_", "$input",
		// The following would also be skiplisted, but they aren't valid
		// identifiers, so can't be accessed with the '.' operator anyway.
		// "$*", "$&", "$+", "$`", "$'",
		"input", "lastMatch", "lastParen", "leftContext", "rightContext",
		"global", "ignoreCase", "lastIndex", "multiline", "source":
		return true
	default:
		return false
	}
}

// Java: CheckRegExp.isGlobalRegExpPropertiesUsed()
func (c *CheckRegExp) IsGlobalRegExpPropertiesUsed() bool {
	if c.assumeAllGlobalRegexpUsagesVisible {
		return c.globalRegExpPropertiesUsed
	}
	return true
}

// Java: CheckRegExp(AbstractCompiler, boolean, boolean)
func NewCheckRegExp(
	compiler AbstractCompiler, assumeAllGlobalRegexpUsagesVisible bool, reportErrors bool) *CheckRegExp {
	c := &CheckRegExp{}
	c.globalRegExpPropertiesUsed = false
	c.compiler = compiler
	c.assumeAllGlobalRegexpUsagesVisible = assumeAllGlobalRegexpUsagesVisible
	c.reportErrors = reportErrors
	return c
}

// Java: CheckRegExp.process(Node, Node)
func (c *CheckRegExp) Process(externs rhino.Node, root rhino.Node) {
	if c.processParallel(root) {
		return
	}
	NodeTraversal_Traverse(c.compiler, root, c)
}

func (c *CheckRegExp) processParallel(root rhino.Node) bool {
	var usedByWorker []bool
	ok := traverseScriptsParallel(c.compiler, root, parallelScriptOptions[*CheckRegExp]{
		useSharedGlobalScope: true,
		prepare: func(numWorkers, _ int) bool {
			usedByWorker = make([]bool, numWorkers)
			return true
		},
		initWorker: func(_ int, rc *Compiler) *CheckRegExp {
			return NewCheckRegExp(rc, c.assumeAllGlobalRegexpUsagesVisible, c.reportErrors)
		},
		finishWorker: func(workerID int, w *CheckRegExp) {
			if w.globalRegExpPropertiesUsed {
				usedByWorker[workerID] = true
			}
		},
	})
	if !ok {
		return false
	}
	for _, used := range usedByWorker {
		if used {
			c.globalRegExpPropertiesUsed = true
			break
		}
	}
	return true
}

// Java: CheckRegExp.visit(NodeTraversal, Node, Node)
func (c *CheckRegExp) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if NodeUtil_IsReferenceName(n) {
		if n.GetStringID() == checkRegExpID && t.GetScope().GetVarByID(checkRegExpID) == nil {
			parentType := parent.GetToken()
			first := n == parent.GetFirstChild()
			if !((parentType == rhino.Token_NEW && first) ||
				(parentType == rhino.Token_CALL && first) ||
				(parentType == rhino.Token_INSTANCEOF && !first) ||
				parentType == rhino.Token_EQ ||
				parentType == rhino.Token_NE ||
				parentType == rhino.Token_SHEQ ||
				parentType == rhino.Token_SHNE ||
				parentType == rhino.Token_CASE ||
				(parentType == rhino.Token_GETPROP &&
					first &&
					!checkRegExp_isSkippedProperty(parent.GetString()))) {
				if c.reportErrors {
					t.Report(n, CheckRegExp_REGEXP_REFERENCE)
				}
				c.globalRegExpPropertiesUsed = true
			}
		}

		// Check the syntax of regular expression patterns.
	} else if c.reportErrors && n.IsRegExp() {
		pattern := n.GetFirstChild().GetString()
		flags := ""
		if n.HasTwoChildren() {
			flags = n.GetLastChild().GetString()
		}
		if _, err := regex.ParseRegExp(pattern, flags); err != nil {
			t.Report(n, CheckRegExp_MALFORMED_REGEXP, err.Message)
		}
	}
}
