// Port of com/google/javascript/jscomp/SuppressDocWarningsGuard.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// SuppressDocWarningsGuard filters warnings based on in-code @suppress
// annotations.
//
// Works by looking at the AST node associated with the warning, and looking
// at parents of the node until it finds a node declaring a symbol (class,
// function, variable, property, assignment, object literal key) or a script.
// For this reason, it doesn't work for warnings without an associated AST
// node, eg, the ones in parsing/IRFactory. They can be turned off with
// jscomp_off.
type SuppressDocWarningsGuard struct {
	WarningsGuardBase

	// Warnings guards for each suppressible warnings group, indexed by name.
	suppressors *jutil.LinkedHashMap[string, *DiagnosticGroup]

	compiler AbstractCompiler
}

const suppressDocWarningsGuard_serialVersionUID int64 = 1

// NewSuppressDocWarningsGuard: suppressors are the suppressible groups,
// indexed by name.
//
// Java: SuppressDocWarningsGuard(AbstractCompiler, Map<String, DiagnosticGroup>)
func NewSuppressDocWarningsGuard(
	compiler AbstractCompiler, suppressors *jutil.LinkedHashMap[string, *DiagnosticGroup]) *SuppressDocWarningsGuard {
	g := &SuppressDocWarningsGuard{}
	g.compiler = compiler
	g.suppressors = SuppressDocWarningsGuard_CreateSuppressors(suppressors)
	return g
}

// Java: SuppressDocWarningsGuard.createSuppressors(Map<String, DiagnosticGroup>)
func SuppressDocWarningsGuard_CreateSuppressors(
	suppressors *jutil.LinkedHashMap[string, *DiagnosticGroup]) *jutil.LinkedHashMap[string, *DiagnosticGroup] {
	builder := jutil.NewLinkedHashMapFrom(suppressors)

	// Hack: Allow "@suppress {missingProperties}" to mean
	// "@suppress {strictMissingProperties}".
	// TODO(johnlenz): Delete this when it is enabled with missingProperties
	builder.Put(
		"missingProperties",
		NewDiagnosticGroupFromGroups(
			DiagnosticGroups_MISSING_PROPERTIES, DiagnosticGroups_STRICT_MISSING_PROPERTIES))

	// Hack: Allow "@suppress {checkTypes}" to include
	// "strictCheckTypes".
	// TODO(johnlenz): Delete this when it is enabled with missingProperties
	builder.Put(
		"checkTypes",
		NewDiagnosticGroupFromGroups(DiagnosticGroups_CHECK_TYPES, DiagnosticGroups_STRICT_CHECK_TYPES))

	return jutil.NewLinkedHashMapFrom(builder)
}

// Java: SuppressDocWarningsGuard.level(JSError)
func (g *SuppressDocWarningsGuard) Level(error *JSError) *CheckLevel {
	node := error.Node()
	if node.IsNil() {
		node = g.GetScriptNodeBySourceName(error)
	}
	if node.IsNil() {
		return nil
	}

	level := g.GetCheckLevelFromAncestors(error, node)
	if level != nil {
		return level
	}

	// Some errors are on nodes that do not have the script as a parent.
	// Look up the script node by filename.
	scriptNode := g.GetScriptNodeBySourceName(error)
	if !scriptNode.IsNil() {
		info := scriptNode.GetJSDocInfo()
		if info != nil {
			return g.GetCheckLevelFromInfo(error, info)
		}
	}

	return nil
}

// GetCheckLevelFromAncestors searches for @suppress tags on nodes
// introducing symbols:
//
// class & function declarations, variables, assignments, object literal
// keys, and the top level script node.
//
// Java: SuppressDocWarningsGuard.getCheckLevelFromAncestors(JSError, Node)
func (g *SuppressDocWarningsGuard) GetCheckLevelFromAncestors(error *JSError, node rhino.Node) *CheckLevel {
	for current := node; !current.IsNil(); current = current.GetParent() {
		var info *rhino.JSDocInfo
		if current.IsFunction() || current.IsClass() {
			info = NodeUtil_GetBestJSDocInfo(current)
		} else if current.IsScript() {
			info = current.GetJSDocInfo()
		} else if NodeUtil_IsNameDeclaration(current) ||
			NodeUtil_MayBeObjectLitKey(current) ||
			current.IsComputedProp() ||
			current.IsMemberFieldDef() ||
			current.IsComputedFieldDef() ||
			(current.HasParent() && current.GetParent().IsExprResult()) {
			info = NodeUtil_GetBestJSDocInfo(current)
		}

		if info != nil {
			level := g.GetCheckLevelFromInfo(error, info)
			if level != nil {
				return level
			}
		}
	}

	return nil
}

// GetCheckLevelFromInfo: if the given JSDocInfo has an @suppress for the
// given JSError, returns the new level.
//
// Java: SuppressDocWarningsGuard.getCheckLevelFromInfo(JSError, JSDocInfo)
func (g *SuppressDocWarningsGuard) GetCheckLevelFromInfo(error *JSError, info *rhino.JSDocInfo) *CheckLevel {
	var matched *CheckLevel
	info.ForEachSuppression(func(suppressor string) bool {
		group := g.suppressors.Get(suppressor)
		if group == nil {
			return true // Some @suppress tags are for other tools, and may not have a group.
		}
		if group.Matches(error) {
			matched = CheckLevel_OFF.Ptr()
			return false
		}
		return true
	})
	return matched
}

// Java: SuppressDocWarningsGuard.getScriptNodeBySourceName(JSError)
func (g *SuppressDocWarningsGuard) GetScriptNodeBySourceName(error *JSError) rhino.Node {
	if error.SourceName() == nil {
		return rhino.Node{}
	}

	scriptNode := g.compiler.GetScriptNode(*error.SourceName())
	if scriptNode.IsNil() {
		return rhino.Node{}
	}

	jutil.CheckState(scriptNode.IsScript())
	return scriptNode
}

// Java: SuppressDocWarningsGuard.getPriority()
func (g *SuppressDocWarningsGuard) GetPriority() int {
	// Happens after path-based filtering, but before other times
	// of filtering.
	return WarningsGuard_Priority_SUPPRESS_DOC.value()
}
