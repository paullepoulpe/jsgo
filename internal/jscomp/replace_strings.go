// Port of com/google/javascript/jscomp/ReplaceStrings.java @ closure-compiler v20261015

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
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ReplaceStrings replaces JavaScript strings in the list of supplied
// functions with shortened forms. Useful for replacing debug message such as:
// throw new Error("Something bad happened"); with generated codes like: throw
// new Error("a"); This makes the compiled JavaScript smaller and prevents us
// from leaking details about the source code.
//
// Based in concept on the work by Jared Jacobs.
//
// Java: class ReplaceStrings extends AbstractPostOrderCallback implements CompilerPass
type ReplaceStrings struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	placeholderToken string

	compiler AbstractCompiler

	functions     *jutil.LinkedHashMap[string, *ReplaceStrings_Config]
	nameGenerator *DefaultNameGenerator
	results       *jutil.LinkedHashMap[string, *ReplaceStrings_Result]
}

// Java: ReplaceStrings.BAD_REPLACEMENT_CONFIGURATION
var ReplaceStrings_BAD_REPLACEMENT_CONFIGURATION = DiagnosticType_Error(
	"JSC_BAD_REPLACEMENT_CONFIGURATION", "Bad replacement configuration \"{0}\": {1}")

// Java: ReplaceStrings.STRING_REPLACEMENT_TAGGED_TEMPLATE
var ReplaceStrings_STRING_REPLACEMENT_TAGGED_TEMPLATE = DiagnosticType_Warning(
	"JSC_STRING_REPLACEMENT_TAGGED_TEMPLATE",
	"Cannot string-replace arguments of a template literal tag function.")

const replaceStrings_DEFAULT_PLACEHOLDER_TOKEN = "`"

// Java: ReplaceStrings.EXCLUSION_PREFIX
const ReplaceStrings_EXCLUSION_PREFIX = ":!"

const replaceStrings_REPLACE_ONE_MARKER = "?"
const replaceStrings_REPLACE_ALL_MARKER = "*"

// ReplaceStrings_Config describes a function to look for a which parameters
// to replace.
//
// Java: private static class ReplaceStrings.Config
type ReplaceStrings_Config struct {
	// TODO(johnlenz): Support name "groups" so that unrelated strings can
	// reuse strings.  For example, event-id can reuse the names used for logger
	// classes.
	name                     string
	parameters               []int
	excludedFilenameSuffixes *jutil.LinkedHashSet[string]
}

// Java: ReplaceStrings.Config.REPLACE_ALL_VALUE
const ReplaceStrings_Config_REPLACE_ALL_VALUE = 0

// Java: ReplaceStrings.Config(String, List<Integer>, ImmutableSet<String>)
func NewReplaceStrings_Config(
	name string,
	replacementParameters []int,
	excludedFilenameSuffixes *jutil.LinkedHashSet[string]) *ReplaceStrings_Config {
	c := &ReplaceStrings_Config{}
	c.name = name
	c.parameters = replacementParameters
	c.excludedFilenameSuffixes = excludedFilenameSuffixes
	return c
}

// Java: ReplaceStrings.Config.isReplaceAll()
func (c *ReplaceStrings_Config) IsReplaceAll() bool {
	return len(c.parameters) == 1 && slices.Contains(c.parameters, ReplaceStrings_Config_REPLACE_ALL_VALUE)
}

// ReplaceStrings_Result describes a replacement that occurred.
//
// Java: static class ReplaceStrings.Result
type ReplaceStrings_Result struct {
	// The original message with non-static content replaced with
	// placeholderToken.
	original       string
	replacement    string
	didReplacement bool
}

// Java: ReplaceStrings.Result(String, String)
func NewReplaceStrings_Result(original string, replacement string) *ReplaceStrings_Result {
	r := &ReplaceStrings_Result{}
	r.didReplacement = false
	r.original = original
	r.replacement = replacement
	return r
}

// NewReplaceStrings: placeholderToken is the separator to use between string
// parts, used to replace non-static string content; functionsToInspect is a
// list of function configurations in the form of
// function($,,,):exclued_filename_suffix1,excluded_filename_suffix2,... or
// class.prototype.method($,,,):exclued_filename_suffix1,excluded_filename_suffix2,...
//
// Java: ReplaceStrings(AbstractCompiler, String, List<String>)
func NewReplaceStrings(
	compiler AbstractCompiler, placeholderToken string, functionsToInspect []string) *ReplaceStrings {
	r := &ReplaceStrings{
		functions: jutil.NewLinkedHashMap[string, *ReplaceStrings_Config](),
		results:   jutil.NewLinkedHashMap[string, *ReplaceStrings_Result](),
	}
	r.compiler = compiler
	if placeholderToken == "" {
		r.placeholderToken = replaceStrings_DEFAULT_PLACEHOLDER_TOKEN
	} else {
		r.placeholderToken = placeholderToken
	}

	r.nameGenerator = ReplaceStrings_CreateNameGenerator()

	// Initialize the map of functions to inspect for renaming candidates.
	r.ParseConfiguration(functionsToInspect)
	return r
}

// Java: ReplaceStrings.USED_RESULTS
var ReplaceStrings_USED_RESULTS = func(result *ReplaceStrings_Result) bool { return result.didReplacement }

// GetResult gets the list of all replacements performed.
//
// Java: ReplaceStrings.getResult()
func (r *ReplaceStrings) GetResult() []*ReplaceStrings_Result {
	var out []*ReplaceStrings_Result
	for result := range r.results.Values() {
		if ReplaceStrings_USED_RESULTS(result) {
			out = append(out, result)
		}
	}
	return out
}

// GetStringMap gets the list of replaces as a VariableMap.
//
// Java: ReplaceStrings.getStringMap()
func (r *ReplaceStrings) GetStringMap() *VariableMap {
	map_ := jutil.NewImmutableMapBuilder[string, string]()
	for result := range r.results.Values() {
		if ReplaceStrings_USED_RESULTS(result) {
			map_.Put(result.replacement, result.original)
		}
	}

	return NewVariableMap(map_.BuildOrThrow())
}

// Java: ReplaceStrings.process(Node, Node)
func (r *ReplaceStrings) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(r.compiler, root, r)
}

// Java: ReplaceStrings.visit(NodeTraversal, Node, Node)
func (r *ReplaceStrings) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	// TODO(johnlenz): Determine if it is necessary to support ".call" or ".apply".
	switch n.GetToken() {
	case rhino.Token_NEW, // e.g. new Error('msg');
		rhino.Token_CALL,               // e.g. Error('msg');
		rhino.Token_TAGGED_TEMPLATELIT: // e.g. Error`msg` - not supported!
		calledFn := n.GetFirstChild()

		// Look for calls to static functions.
		name := calledFn.GetOriginalQualifiedName()
		if name != nil {
			config := r.FindMatching(*name, n.GetSourceFileName())
			if config != nil {
				r.DoSubstitutions(t, config, n)
				return
			}
		}
	default:
	}
}

// FindMatching: name is the function name to find; callsiteSourceFileName is
// the filename containing the callsite. It returns the Config object for the
// name or null if no match was found.
//
// Java: ReplaceStrings.findMatching(String, String)
func (r *ReplaceStrings) FindMatching(name string, callsiteSourceFileName *string) *ReplaceStrings_Config {
	config := r.functions.Get(name)
	if config == nil {
		name = strings.ReplaceAll(name, "$", ".")
		config = r.functions.Get(name)
	}
	if config != nil {
		for excludedSuffix := range config.excludedFilenameSuffixes.All() {
			// callsiteSourceFileName.endsWith(excludedSuffix): a null name throws.
			if callsiteSourceFileName == nil {
				panic(jutil.NewNullPointerException())
			}
			if strings.HasSuffix(*callsiteSourceFileName, excludedSuffix) {
				return nil
			}
		}
	}
	return config
}

// DoSubstitutions replaces the parameters specified in the config, if
// possible.
//
// Java: ReplaceStrings.doSubstitutions(NodeTraversal, ReplaceStrings.Config, Node)
func (r *ReplaceStrings) DoSubstitutions(t *NodeTraversal, config *ReplaceStrings_Config, n *rhino.Node) {
	if n.IsTaggedTemplateLit() {
		// This is currently not supported, since tagged template literals have a different calling
		// convention than ordinary functions, so it's unclear which arguments are expected to be
		// replaced. Specifically, there are no direct string arguments, and for arbitrary tag
		// functions it's not clear that it's safe to inline any constant placeholders.
		r.compiler.Report(JSError_MakeWithNode(n, ReplaceStrings_STRING_REPLACEMENT_TAGGED_TEMPLATE))
		return
	}
	jutil.CheckState(n.IsNew() || n.IsCall())

	if !config.IsReplaceAll() {
		// Note: the first child is the function, but the parameter id is 1 based.
		for _, parameter := range config.parameters {
			arg := n.GetChildAtIndex(parameter)
			if arg != nil {
				r.ReplaceExpression(t, arg)
			}
		}
	} else {
		// Replace all parameters.
		firstParam := n.GetSecondChild()
		for arg := firstParam; arg != nil; arg = arg.GetNext() {
			arg = r.ReplaceExpression(t, arg)
		}
	}
}

// ReplaceExpression replaces a string expression with a short encoded string
// expression. It returns the replacement node (or the original expression if
// no replacement is made).
//
// Java: ReplaceStrings.replaceExpression(NodeTraversal, Node)
func (r *ReplaceStrings) ReplaceExpression(t *NodeTraversal, expr *rhino.Node) *rhino.Node {
	var replacement *rhino.Node
	var key string
	var replacementString string
	switch expr.GetToken() {
	case rhino.Token_STRINGLIT:
		key = expr.GetString()
		replacementString = r.GetReplacement(key)
		replacement = rhino.IR_String(replacementString)
	case rhino.Token_TEMPLATELIT, rhino.Token_ADD, rhino.Token_NAME:
		var keyBuilder strings.Builder
		keyNode := rhino.IR_String("")
		replacement = r.BuildReplacement(t, expr, keyNode, &keyBuilder)
		key = keyBuilder.String()
		if key == r.placeholderToken {
			// There is no static text in expr - only a placeholder - so just return expr directly.
			// In this case, replacement is just the string join ('`' + expr), which is not useful.
			return expr
		}
		replacementString = r.GetReplacement(key)
		keyNode.SetString(replacementString)
	default:
		// This may be a function call or a variable reference. We don't
		// replace these.
		return expr
	}

	// checkNotNull(key); checkNotNull(replacementString): both are set above.
	r.RecordReplacement(key)

	replacement.SrcrefTreeIfMissing(expr)
	expr.ReplaceWith(replacement)
	t.ReportCodeChange()
	return replacement
}

// GetReplacement gets a replacement string for the provide key text.
//
// Java: ReplaceStrings.getReplacement(String)
func (r *ReplaceStrings) GetReplacement(key string) string {
	result := r.results.Get(key)
	if result != nil {
		return result.replacement
	}

	replacement := r.nameGenerator.GenerateNextName()
	result = NewReplaceStrings_Result(key, replacement)
	r.results.Put(key, result)
	return replacement
}

// RecordReplacement records the location the replacement was made.
//
// Java: ReplaceStrings.recordReplacement(String)
func (r *ReplaceStrings) RecordReplacement(key string) {
	result := r.results.Get(key)
	jutil.CheckState(result != nil)

	result.didReplacement = true
}

// BuildReplacement builds a replacement abstract syntax tree for the string
// expression expr. Appends any string literal values that are encountered to
// keyBuilder, to build the expression's replacement key.
//
// expr is a JS expression that evaluates to a string value; prefix is the JS
// expression to which expr's replacement is logically being concatenated. It
// is a partial solution to the problem at hand and will either be this
// method's return value or a descendant of it. keyBuilder is a builder of the
// string expression's replacement key. It returns the abstract syntax tree
// that should replace expr.
//
// Java: ReplaceStrings.buildReplacement(NodeTraversal, Node, Node, StringBuilder)
func (r *ReplaceStrings) BuildReplacement(t *NodeTraversal, expr *rhino.Node, prefix *rhino.Node, keyBuilder *strings.Builder) *rhino.Node {
	switch expr.GetToken() {
	case rhino.Token_ADD:
		left := expr.GetFirstChild()
		right := left.GetNext()
		prefix = r.BuildReplacement(t, left, prefix, keyBuilder)
		return r.BuildReplacement(t, right, prefix, keyBuilder)
	case rhino.Token_TEMPLATELIT:
		for child := expr.GetFirstChild(); child != nil; child = child.GetNext() {
			switch child.GetToken() {
			case rhino.Token_TEMPLATELIT_STRING:
				keyBuilder.WriteString(jutil.ToJavaString(child.GetCookedString()))
			case rhino.Token_TEMPLATELIT_SUB:
				prefix = r.BuildReplacement(t, child.GetFirstChild(), prefix, keyBuilder)
			default:
				panic(jutil.NewIllegalStateExceptionMsg("Unexpected TEMPLATELIT child: " + child.String()))
			}
		}
		return prefix
	case rhino.Token_STRINGLIT:
		keyBuilder.WriteString(expr.GetString())
		return prefix
	case rhino.Token_NAME:
		// If the referenced variable is a constant, use its value.
		var_ := t.GetScope().GetVar(expr.GetString())
		if var_ != nil && (var_.IsDeclaredOrInferredConst() || var_.IsConst()) {
			initialValue := var_.GetInitialValue()
			if initialValue != nil {
				newKeyNode := rhino.IR_String("")
				var newKeyBuilder strings.Builder
				replacement := r.BuildReplacement(t, initialValue, newKeyNode, &newKeyBuilder)
				if replacement == newKeyNode {
					keyBuilder.WriteString(newKeyBuilder.String())
					return prefix
				}
			}
			// Not a simple string constant.
		}
		fallthrough
	default:
		keyBuilder.WriteString(r.placeholderToken)
		prefix = rhino.IR_Add(prefix, rhino.IR_String(r.placeholderToken))
		return rhino.IR_Add(prefix, expr.CloneTree())
	}
}

// ParseConfiguration builds the data structures need by this pass from the
// provided list of functions.
//
// Java: ReplaceStrings.parseConfiguration(List<String>)
func (r *ReplaceStrings) ParseConfiguration(functionsToInspect []string) {
	for _, function := range functionsToInspect {
		config := r.ParseConfigurationOfFunction(function)
		if config != nil {
			r.functions.Put(config.name, config)
		}
	}
}

// ParseConfigurationOfFunction converts the provide string into a Config. The
// string can be a variable or static function: foo(,,?) foo.bar(?) (but not a
// prototype method) and is allowed to either replace all parameters using "*"
// or one parameter "?". "," is used as a placeholder for ignored parameters.
//
// It returns null if this is an invalid function, otherwise the Config.
//
// Java: ReplaceStrings.parseConfiguration(String)
func (r *ReplaceStrings) ParseConfigurationOfFunction(function string) *ReplaceStrings_Config {
	// Looks like this function_name(,$,)
	first := jutil.StringIndexOfChar(function, '(')
	last := jutil.StringIndexOfChar(function, ')')
	colon := jutil.StringIndexOf(function, ReplaceStrings_EXCLUSION_PREFIX)

	// TODO(johnlenz): Make parsing precondition checks JSErrors reports.
	jutil.CheckState(first != -1 && last != -1)

	if strings.Contains(function, ".prototype.") {
		r.compiler.Report(
			JSError_Make(
				ReplaceStrings_BAD_REPLACEMENT_CONFIGURATION,
				function,
				"Cannot replace strings passed to prototype methods."))
		return nil
	}

	name := jutil.StringSubstring(function, 0, first)
	params := jutil.StringSubstring(function, first+1, last)

	paramCount := 0
	var replacementParameters []int
	parts := strings.Split(params, ",") // Splitter.on(',').splitToList(params)
	for _, param := range parts {
		paramCount++
		if param == replaceStrings_REPLACE_ALL_MARKER {
			jutil.CheckState(paramCount == 1 && len(parts) == 1)
			replacementParameters = append(replacementParameters, ReplaceStrings_Config_REPLACE_ALL_VALUE)
		} else if param == replaceStrings_REPLACE_ONE_MARKER {
			// TODO(johnlenz): Support multiple.
			jutil.CheckState(!slices.Contains(replacementParameters, ReplaceStrings_Config_REPLACE_ALL_VALUE))
			replacementParameters = append(replacementParameters, paramCount)
		} else {
			// TODO(johnlenz): report an error.
			jutil.CheckStateMsg1(param == "", "Unknown marker (%s)", param)
		}
	}

	jutil.CheckState(len(replacementParameters) != 0)

	excludedFilenameSuffixes := jutil.NewLinkedHashSet[string]()
	if colon != -1 {
		for _, suffix := range jutil.StringSplitLiteral(
			jutil.StringSubstringFrom(function, colon+jutil.StringLength(ReplaceStrings_EXCLUSION_PREFIX)), ",") {
			excludedFilenameSuffixes.Add(suffix)
		}
	}
	return NewReplaceStrings_Config(
		name,
		replacementParameters,
		excludedFilenameSuffixes)
}

// ReplaceStrings_CreateNameGenerator uses a name generate to create names so
// the names overlap with the names used for variable and properties.
//
// Java: ReplaceStrings.createNameGenerator()
func ReplaceStrings_CreateNameGenerator() *DefaultNameGenerator {
	namePrefix := ""
	reservedChars := jutil.NewLinkedHashSet[uint16]()
	return NewDefaultNameGeneratorWithReserved(jutil.NewLinkedHashSet[string](), namePrefix, reservedChars)
}
