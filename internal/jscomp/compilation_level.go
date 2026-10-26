// Port of com/google/javascript/jscomp/CompilationLevel.java @ closure-compiler v20261021

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

// CompilationLevel represents the level of optimization that should be
// applied when compiling JavaScript code.
type CompilationLevel int32

const (
	// CompilationLevel_BUNDLE simply orders and concatenates files to the
	// output.
	CompilationLevel_BUNDLE CompilationLevel = iota

	// CompilationLevel_WHITESPACE_ONLY removes comments and extra whitespace
	// in the input JS.
	CompilationLevel_WHITESPACE_ONLY

	// CompilationLevel_SIMPLE_OPTIMIZATIONS performs transformations to the
	// input JS that do not require any changes to JS that depend on the input
	// JS. For example, function arguments are renamed (which should not
	// matter to code that depends on the input JS), but functions themselves
	// are not renamed (which would otherwise require external code to change
	// to use the renamed function names).
	CompilationLevel_SIMPLE_OPTIMIZATIONS

	// CompilationLevel_ADVANCED_OPTIMIZATIONS aggressively reduces code size
	// by renaming function names and variables, removing code which is never
	// called, etc.
	CompilationLevel_ADVANCED_OPTIMIZATIONS

	// CompilationLevel_TRANSPILE_ONLY only transpiles the code down to the
	// requested `--language_out level`, including the addition of polyfills.
	// This mode will fully parse the JS files it is running, so that it can
	// find syntax that needs to be transpiled down and discover required
	// polyfills. This mode will not perform any optimizations.
	CompilationLevel_TRANSPILE_ONLY
)

var compilationLevel_names = [...]string{
	"BUNDLE", "WHITESPACE_ONLY", "SIMPLE_OPTIMIZATIONS", "ADVANCED_OPTIMIZATIONS", "TRANSPILE_ONLY",
}

// String is Enum.name().
func (l CompilationLevel) String() string { return compilationLevel_names[l] }

// Ordinal is Enum.ordinal().
func (l CompilationLevel) Ordinal() int { return int(l) }

var compilationLevel_ptrs = [...]CompilationLevel{
	CompilationLevel_BUNDLE, CompilationLevel_WHITESPACE_ONLY, CompilationLevel_SIMPLE_OPTIMIZATIONS,
	CompilationLevel_ADVANCED_OPTIMIZATIONS, CompilationLevel_TRANSPILE_ONLY,
}

// Ptr returns a shared pointer to the constant, for nullable
// CompilationLevel values.
//
// Go-only: language-adapter — a Java enum reference may be null.
func (l CompilationLevel) Ptr() *CompilationLevel { return &compilationLevel_ptrs[l] }

// CompilationLevel_FromString: value may be nil, and so may the result.
//
// Java: CompilationLevel.fromString(String)
func CompilationLevel_FromString(value *string) *CompilationLevel {
	if value == nil {
		return nil
	}
	switch *value {
	case "BUNDLE":
		return CompilationLevel_BUNDLE.Ptr()
	case "WHITESPACE_ONLY", "WHITESPACE":
		return CompilationLevel_WHITESPACE_ONLY.Ptr()
	case "TRANSPILE_ONLY":
		return CompilationLevel_TRANSPILE_ONLY.Ptr()
	case "SIMPLE_OPTIMIZATIONS", "SIMPLE":
		return CompilationLevel_SIMPLE_OPTIMIZATIONS.Ptr()
	case "ADVANCED_OPTIMIZATIONS", "ADVANCED":
		return CompilationLevel_ADVANCED_OPTIMIZATIONS.Ptr()
	default:
		return nil
	}
}

// Java: CompilationLevel.setOptionsForCompilationLevel(CompilerOptions)
func (l CompilationLevel) SetOptionsForCompilationLevel(options *CompilerOptions) {
	switch l {
	case CompilationLevel_BUNDLE:
	case CompilationLevel_WHITESPACE_ONLY:
		CompilationLevel_ApplyBasicCompilationOptions(options)
	case CompilationLevel_SIMPLE_OPTIMIZATIONS:
		CompilationLevel_ApplySafeCompilationOptions(options)
	case CompilationLevel_TRANSPILE_ONLY:
		CompilationLevel_ApplyTranspileOnlyOptions(options)
	case CompilationLevel_ADVANCED_OPTIMIZATIONS:
		CompilationLevel_ApplyFullCompilationOptions(options)
	}
}

// Java: CompilationLevel.setDebugOptionsForCompilationLevel(CompilerOptions)
func (l CompilationLevel) SetDebugOptionsForCompilationLevel(options *CompilerOptions) {
	options.SetGeneratePseudoNames(true)
	options.SetRemoveClosureAsserts(false)
	options.SetRemoveJ2clAsserts(true)
}

// CompilationLevel_ApplyBasicCompilationOptions gets options that only strip
// whitespace and comments.
//
// Java: CompilationLevel.applyBasicCompilationOptions(CompilerOptions)
func CompilationLevel_ApplyBasicCompilationOptions(options *CompilerOptions) {
	options.SkipAllCompilerPasses()
}

// CompilationLevel_ApplyTranspileOnlyOptions gets options that only
// transpile the code, including polyfills.
//
// Java: CompilationLevel.applyTranspileOnlyOptions(CompilerOptions)
func CompilationLevel_ApplyTranspileOnlyOptions(options *CompilerOptions) {
	// ReplaceIdGenerators is on by default, but should not run in TRANSPILE_ONLY mode.
	options.SetReplaceIdGenerators(false)

	options.SetClosurePass(true)
	options.SetRenamingPolicy(VariableRenamingPolicy_OFF, PropertyRenamingPolicy_OFF)
	options.SetInlineVariablesOfReach(CompilerOptions_Reach_NONE)
	options.SetInlineFunctions(CompilerOptions_Reach_NONE)
	options.SetFoldConstants(false)
	options.SetCoalesceVariableNames(false)
	options.SetDeadAssignmentElimination(false)
	options.SetCollapseVariableDeclarations(false)
	options.SetConvertToDottedProperties(false)
	options.SetLabelRenaming(false)
	options.SetRemoveUnusedVariables(CompilerOptions_Reach_NONE)
	options.SetCollapseObjectLiterals(false)
	// Turn off protecting side-effect free nodes by making them parameters to a extern function
	// call.
	options.SetProtectHiddenSideEffects(false)
}

// CompilationLevel_ApplySafeCompilationOptions adds options that are safe.
// Safe means options that won't break the JavaScript code even if no symbols
// are exported and no coding convention is used.
//
// Java: CompilationLevel.applySafeCompilationOptions(CompilerOptions)
func CompilationLevel_ApplySafeCompilationOptions(options *CompilerOptions) {
	// TODO(tjgq): Remove this.
	options.SetDependencyOptions(DependencyOptions_SortOnly())

	// ReplaceIdGenerators is on by default, but should not run in simple mode.
	options.SetReplaceIdGenerators(false)

	// Does not call applyBasicCompilationOptions(options) because the call to
	// skipAllCompilerPasses() cannot be easily undone.
	options.SetClosurePass(true)
	options.SetRenamingPolicy(VariableRenamingPolicy_LOCAL, PropertyRenamingPolicy_OFF)
	options.SetInlineVariablesOfReach(CompilerOptions_Reach_LOCAL_ONLY)
	options.SetInlineFunctions(CompilerOptions_Reach_LOCAL_ONLY)
	options.SetWarningLevel(DiagnosticGroups_GLOBAL_THIS, CheckLevel_OFF)
	options.SetFoldConstants(true)
	options.SetCoalesceVariableNames(true)
	options.SetDeadAssignmentElimination(true)
	options.SetDeadPropertyAssignmentElimination(false)
	options.SetCollapseVariableDeclarations(true)
	options.SetConvertToDottedProperties(true)
	options.SetLabelRenaming(true)
	options.SetRemoveUnusedVariables(CompilerOptions_Reach_LOCAL_ONLY)
	options.SetCollapseObjectLiterals(true)
	options.SetProtectHiddenSideEffects(true)
}

// CompilationLevel_ApplyFullCompilationOptions adds the options that will
// work only if the user exported all the symbols correctly.
//
// Java: CompilationLevel.applyFullCompilationOptions(CompilerOptions)
func CompilationLevel_ApplyFullCompilationOptions(options *CompilerOptions) {
	// TODO(tjgq): Remove this.
	options.SetDependencyOptions(DependencyOptions_SortOnly())

	// Do not call applySafeCompilationOptions(options) because the call can
	// create possible conflicts between multiple diagnostic groups.

	options.SetCheckSymbols(true)
	options.SetCheckTypes(true)

	// All the safe optimizations.
	options.SetClosurePass(true)
	options.SetFoldConstants(true)
	options.SetCoalesceVariableNames(true)
	options.SetDeadAssignmentElimination(true)
	options.SetExtractPrototypeMemberDeclarations(true)
	options.SetCollapseVariableDeclarations(true)
	options.SetConvertToDottedProperties(true)
	options.SetLabelRenaming(true)
	options.SetCollapseObjectLiterals(true)
	options.SetProtectHiddenSideEffects(true)

	// All the advanced optimizations.
	options.SetRemoveClosureAsserts(true)
	options.SetRemoveAbstractMethods(true)
	options.SetReserveRawExports(true)
	options.SetRenamingPolicy(VariableRenamingPolicy_ALL, PropertyRenamingPolicy_ALL_UNQUOTED)
	options.SetRemoveUnusedPrototypeProperties(true)
	options.SetRemoveUnusedClassProperties(true)
	options.SetCollapseAnonymousFunctions(true)
	options.SetCollapsePropertiesLevel(CompilerOptions_PropertyCollapseLevel_ALL)
	options.SetWarningLevel(DiagnosticGroups_GLOBAL_THIS, CheckLevel_WARNING)
	options.SetRewriteFunctionExpressions(false)
	options.SetSmartNameRemoval(true)
	options.SetInlineConstantVars(true)
	options.SetInlineFunctions(CompilerOptions_Reach_ALL)
	options.SetInlineVariablesOfReach(CompilerOptions_Reach_ALL)
	options.SetComputeFunctionSideEffects(true)
	options.SetAssumeStrictThis(true)

	// Remove unused vars also removes unused functions.
	options.SetRemoveUnusedVariables(CompilerOptions_Reach_ALL)

	// Move code around based on the defined modules.
	options.SetCrossChunkCodeMotion(true)
	options.SetCrossChunkMethodMotion(true)

	// Call optimizations
	options.SetDevirtualizeMethods(true)
	options.SetOptimizeCalls(true)
	options.SetOptimizeESClassConstructors(true)
}

// SetTypeBasedOptimizationOptions enables additional optimizations that use
// type information. Only has an effect for ADVANCED_OPTIMIZATIONS; this is a
// no-op for other modes.
//
// Java: CompilationLevel.setTypeBasedOptimizationOptions(CompilerOptions)
func (l CompilationLevel) SetTypeBasedOptimizationOptions(options *CompilerOptions) {
	switch l {
	case CompilationLevel_ADVANCED_OPTIMIZATIONS:
		options.SetDisambiguateProperties(true)
		options.SetAmbiguateProperties(true)
		options.SetInlineProperties(true)
		options.SetUseTypesForLocalOptimization(true)
	case CompilationLevel_SIMPLE_OPTIMIZATIONS, CompilationLevel_WHITESPACE_ONLY, CompilationLevel_BUNDLE,
		CompilationLevel_TRANSPILE_ONLY:
	}
}

// SetWrappedOutputOptimizations enables additional optimizations that
// operate on global declarations. Advanced mode does this by default, but
// this isn't valid in simple mode in the general case and should only be
// enabled when code is self contained (such as when it is enclosed by a
// function wrapper.
//
// Java: CompilationLevel.setWrappedOutputOptimizations(CompilerOptions)
func (l CompilationLevel) SetWrappedOutputOptimizations(options *CompilerOptions) {
	// Global variables and properties names can't conflict.
	options.SetReserveRawExports(false)
	switch l {
	case CompilationLevel_SIMPLE_OPTIMIZATIONS:
		// Enable global variable optimizations (but not property optimizations)
		options.SetVariableRenaming(VariableRenamingPolicy_ALL)
		options.SetCollapsePropertiesLevel(CompilerOptions_PropertyCollapseLevel_MODULE_EXPORT)
		options.SetCollapseAnonymousFunctions(true)
		options.SetInlineConstantVars(true)
		options.SetInlineFunctions(CompilerOptions_Reach_ALL)
		options.SetInlineVariablesOfReach(CompilerOptions_Reach_ALL)
		options.SetRemoveUnusedVariables(CompilerOptions_Reach_ALL)
	case CompilationLevel_ADVANCED_OPTIMIZATIONS, CompilationLevel_WHITESPACE_ONLY, CompilationLevel_BUNDLE,
		CompilationLevel_TRANSPILE_ONLY:
	}
}
