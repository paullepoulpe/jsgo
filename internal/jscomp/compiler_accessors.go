// Port of com/google/javascript/jscomp/Compiler.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/Result.java @ closure-compiler v20261021

/*
 * Copyright 2004 The Closure Compiler Authors.
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

import (
	"runtime"
	"slices"
	"weak"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Java: Compiler.getScriptNode(String)
func (c *Compiler) GetScriptNode(filename string) rhino.Node {
	return c.scriptNodeByFilename.Get(jutil.CheckNotNull(filename))
}

// GetResult returns the result of the compilation.
//
// Java: Compiler.getResult()
func (c *Compiler) GetResult() *Result {
	return NewResultWithAllMaps(
		c.GetErrors(),
		c.GetWarnings(),
		c.variableMap,
		c.propertyMap,
		nil,
		c.stringMap,
		c.instrumentationMapping,
		c.sourceMap,
		c.externExports,
		c.cssNames,
		c.idGeneratorMap,
		c.transpiledFiles)
}

// SetDefaultCodingConvention overrides the compiler's fallback CodingConvention
// used when CompilerOptions does not specify one.
func (c *Compiler) SetDefaultCodingConvention(cc CodingConvention) {
	c.defaultCodingConvention = jutil.CheckNotNull(cc)
}

// GetDefaultCodingConvention returns the compiler's fallback CodingConvention.
func (c *Compiler) GetDefaultCodingConvention() CodingConvention {
	return c.defaultCodingConvention
}

// GetIdGeneratorMap returns the serialized ID generator mapping, or "" when
// none has been produced, without allocating a full Result.
func (c *Compiler) GetIdGeneratorMap() string {
	if c.idGeneratorMap == nil {
		return ""
	}
	return *c.idGeneratorMap
}

// GetConformanceViolations returns the attributed conformance violation protos
// persisted on this Compiler across compilation stages.
func (c *Compiler) GetConformanceViolations() []*ConformanceViolationsProto {
	return c.conformanceViolations
}

// SetConformanceViolations replaces the attributed conformance violation protos
// stored on this Compiler.
func (c *Compiler) SetConformanceViolations(violations []*ConformanceViolationsProto) {
	c.conformanceViolations = slices.Clone(violations)
}

// AddConformanceViolations appends an attributed conformance violations proto
// to this Compiler.
func (c *Compiler) AddConformanceViolations(violations *ConformanceViolationsProto) {
	c.conformanceViolations = append(c.conformanceViolations, jutil.CheckNotNull(violations))
}

// GetUnattributedConformanceViolations returns the unattributed conformance
// violation protos persisted on this Compiler across compilation stages.
func (c *Compiler) GetUnattributedConformanceViolations() []*ViolationProto {
	return c.unattributedConformanceViolations
}

// SetUnattributedConformanceViolations replaces the unattributed conformance
// violation protos stored on this Compiler.
func (c *Compiler) SetUnattributedConformanceViolations(violations []*ViolationProto) {
	c.unattributedConformanceViolations = slices.Clone(violations)
}

// AddUnattributedConformanceViolation appends an unattributed conformance
// violation proto to this Compiler.
func (c *Compiler) AddUnattributedConformanceViolation(violation *ViolationProto) {
	c.unattributedConformanceViolations = append(c.unattributedConformanceViolations, jutil.CheckNotNull(violation))
}

// GetErrors returns the list of errors (never null).
//
// Java: Compiler.getErrors()
func (c *Compiler) GetErrors() []*JSError {
	if c.errorManager == nil {
		return []*JSError{}
	}
	return c.errorManager.GetErrors()
}

// GetWarnings returns the list of warnings (never null).
//
// Java: Compiler.getWarnings()
func (c *Compiler) GetWarnings() []*JSError {
	if c.errorManager == nil {
		return []*JSError{}
	}
	return c.errorManager.GetWarnings()
}

// Java: Compiler.getRoot()
func (c *Compiler) GetRoot() rhino.Node {
	return c.externAndJsRoot
}

// Java: Compiler.getAllowableFeatures()
func (c *Compiler) GetAllowableFeatures() *parser.FeatureSet {
	return c.allowableFeatures
}

// Java: Compiler.setAllowableFeatures(FeatureSet)
func (c *Compiler) SetAllowableFeatures(allowableFeatures *parser.FeatureSet) {
	c.allowableFeatures = allowableFeatures
}

// Java: Compiler.markFeatureNotAllowed(Feature)
func (c *Compiler) MarkFeatureNotAllowed(feature parser.FeatureSet_Feature) {
	c.allowableFeatures = c.allowableFeatures.WithoutFeatures(feature)
}

// Java: Compiler.markFeatureSetNotAllowed(FeatureSet)
func (c *Compiler) MarkFeatureSetNotAllowed(featureSet *parser.FeatureSet) {
	c.allowableFeatures = c.allowableFeatures.Without(featureSet)
}

// Java: Compiler.getUniqueIdSupplier()
func (c *Compiler) GetUniqueIdSupplier() *UniqueIdSupplier {
	return c.uniqueIdSupplier
}

// nextUniqueNameId creates a new id for making unique names.
//
// Deprecated.
//
// Java: Compiler.nextUniqueNameId()
func (c *Compiler) nextUniqueNameId() int {
	id := c.uniqueNameId
	c.uniqueNameId++
	return id
}

// GetUniqueNameIdSupplier is the legacy supplier for getting unique names.
//
// Deprecated: This is deprecated because this supplier fails to generate
// unique Ids across input files when used for tagged template literal
// transpilation. Kindly use the new UniqueIdSupplier instead.
//
// Java: Compiler.getUniqueNameIdSupplier()
func (c *Compiler) GetUniqueNameIdSupplier() func() string {
	if c.uniqueNameIdSupplier == nil {
		c.uniqueNameIdSupplier = func() string { return nodeUtilFormatInt(c.nextUniqueNameId()) }
	}
	return c.uniqueNameIdSupplier
}

// Java: Compiler.areNodesEqualForInlining(Node, Node)
func (c *Compiler) AreNodesEqualForInlining(n1 rhino.Node, n2 rhino.Node) bool {
	if c.options.ShouldAmbiguateProperties() || c.options.ShouldDisambiguateProperties() {
		// The type based optimizations require that type information is preserved
		// during other optimizations.
		return n1.IsEquivalentToTyped(n2)
	}
	return n1.IsEquivalentTo(n2)
}

// ------------------------------------------------------------------------
// Inputs
// ------------------------------------------------------------------------

// Java: Compiler.getInput(InputId)
func (c *Compiler) GetInput(id *rhino.InputId) *CompilerInput {
	if id == nil {
		return nil
	}
	return c.inputsById[id.GetIdName()]
}

// PutCompilerInput registers input by its InputId on c.
//
// Java: Compiler.putCompilerInput(CompilerInput)
func (c *Compiler) PutCompilerInput(input *CompilerInput) *CompilerInput {
	return c.putCompilerInput(input)
}

func (c *Compiler) putCompilerInput(input *CompilerInput) *CompilerInput {
	input.SetCompiler(c)
	idName := input.GetInputId().GetIdName()
	prev := c.inputsById[idName]
	c.inputsById[idName] = input
	return prev
}

// GetChunkGraph gets the graph of JS source chunks.
//
// Returns null if #init or #initChunks hasn't been called yet. Otherwise,
// the result is always a chunk graph, even in the degenerate case where
// there's only one chunk.
//
// Java: Compiler.getChunkGraph()
func (c *Compiler) GetChunkGraph() *JSChunkGraph {
	return c.chunkGraph
}

// GetChunks gets the JS source chunks in dependency order.
//
// Returns null if #init or #initChunks hasn't been called yet. Otherwise,
// the result is always non-empty, even in the degenerate case where
// there's only one chunk.
//
// Java: Compiler.getChunks()
func (c *Compiler) GetChunks() []*JSChunk {
	if c.chunkGraph != nil {
		return c.chunkGraph.GetAllChunks()
	}
	return nil
}

// Java: Compiler.clearJSTypeRegistry()
func (c *Compiler) ClearJSTypeRegistry() {
	if c.typeRegistry != nil {
		c.clearedTypeRegistry = weak.Make(c.typeRegistry)
		// Drop the type column: every type a Node referred to is released
		// here, whether or not the node is still reachable.
		c.typeRegistry.ReleaseNodeTypeHandles()
	}
	c.typeRegistry = nil
	c.typeValidator = nil
	c.abstractInterpreter = nil
	// Go-only: the cached analyzer holds the registry it was built with.
	c.cachedAstAnalyzer = nil
	if AfterTypesToColorsHook != nil {
		AfterTypesToColorsHook()
	}
}

// AfterTypesToColorsHook, when set, runs at the end of ClearJSTypeRegistry.
//
// Go-only: measurement aid (cmd/jsgo records the live heap there).
var AfterTypesToColorsHook func()

// typeRegistryReleased reports whether the registry that ClearJSTypeRegistry
// dropped has been garbage collected.
//
// Go-only: debugging aid for tests.
func (c *Compiler) typeRegistryReleased() bool {
	runtime.GC()
	runtime.GC()
	return c.clearedTypeRegistry.Value() == nil
}

// TypeRegistryReleased is typeRegistryReleased for cmd/jsgo's debugging
// hook (Go-only).
func (c *Compiler) TypeRegistryReleased() bool { return c.typeRegistryReleased() }

// Java: Compiler.isTypeRegistryCleared()
func (c *Compiler) IsTypeRegistryCleared() bool {
	return c.typeCheckingHasRun && c.typeRegistry == nil
}

// Java: Compiler.getTypeRegistry()
func (c *Compiler) GetTypeRegistry() *rhino.JSTypeRegistry {
	if c.typeRegistry == nil {
		jutil.CheckStateMsg(
			!c.HasTypeCheckingRun(),
			// This could throw when calling "getTypeRegistry()" during the optimizations phase when
			// JSTypes have been converted to optimization colors
			"Attempted to re-initialize JSTypeRegistry after it had been cleared")
		c.typeRegistry = rhino.NewJSTypeRegistryWithForwardDeclaredTypes(c.oldErrorReporter, &c.forwardDeclaredTypes)
	}
	return c.typeRegistry
}

// Java: Compiler.getColorRegistry()
func (c *Compiler) GetColorRegistry() *colors.ColorRegistry {
	return jutil.CheckNotNullMsg(c.colorRegistry, "Color registry has not been initialized yet")
}

// Java: Compiler.setColorRegistry(ColorRegistry)
func (c *Compiler) SetColorRegistry(colorRegistry *colors.ColorRegistry) {
	jutil.CheckState(c.runtimeLibraryTypedAsts != nil)
	c.colorRegistry = colorRegistry
}

// Java: Compiler.forwardDeclareType(String)
func (c *Compiler) ForwardDeclareType(typeName string) {
	c.forwardDeclaredTypes.Add(typeName)
}

// Java: Compiler.setTypeCheckingHasRun(boolean)
func (c *Compiler) SetTypeCheckingHasRun(hasRun bool) {
	c.typeCheckingHasRun = hasRun
}

// Java: Compiler.hasTypeCheckingRun()
func (c *Compiler) HasTypeCheckingRun() bool {
	return c.typeCheckingHasRun
}

// Java: Compiler.hasOptimizationColors()
func (c *Compiler) HasOptimizationColors() bool {
	return c.colorRegistry != nil
}

// Java: Compiler.getTypedScopeCreator()
func (c *Compiler) GetTypedScopeCreator() ScopeCreator {
	if c.typedScopeCreator == nil {
		jutil.CheckStateMsg(
			!c.HasTypeCheckingRun(),
			// This could throw when calling "getTypedScopeCreator()" during the optimizations phase
			// when JSTypes have been converted to optimization colors
			"Attempted to re-initialize TypedScopeCreator after it had been cleared")
		c.typedScopeCreator = NewTypedScopeCreator(c)
	}
	return c.typedScopeCreator
}

// Java: Compiler.clearTypedScopeCreator()
func (c *Compiler) ClearTypedScopeCreator() {
	if c.typedScopeCreator != nil {
		c.typedScopeCreator.memoized = nil
		c.typedScopeCreator = nil
	}
}

// Java: Compiler.getTopScope()
func (c *Compiler) GetTopScope() *TypedScope {
	return c.topScope
}

// Java: Compiler.setTopScope(TypedScope)
func (c *Compiler) SetTopScope(x *TypedScope) {
	jutil.CheckStateMsg1(x == nil || x.GetParent() == nil, "%s", x)
	c.topScope = x
}

// Java: Compiler.getReverseAbstractInterpreter()
func (c *Compiler) GetReverseAbstractInterpreter() ReverseAbstractInterpreter {
	if c.abstractInterpreter == nil {
		var interpreter ChainableReverseAbstractInterpreter = NewSemanticReverseAbstractInterpreter(c.GetTypeRegistry())
		if c.options.GetClosurePass() {
			interpreter = NewClosureReverseAbstractInterpreter(c.GetTypeRegistry()).Append(interpreter).GetFirst()
		}
		c.abstractInterpreter = interpreter
	}
	return c.abstractInterpreter
}

// Java: Compiler.getTypeValidator()
func (c *Compiler) GetTypeValidator() *TypeValidator {
	if c.typeValidator == nil {
		jutil.CheckStateMsg(
			!c.HasTypeCheckingRun(),
			// This could throw when calling "getTypeValidator()" during the optimizations phase when
			// JSTypes have been converted to optimization colors
			"Attempted to re-initialize TypeValidator after it had been cleared")
		c.typeValidator = NewTypeValidator(c)
	}
	return c.typeValidator
}

// Java: Compiler.getTypeMismatches()
func (c *Compiler) GetTypeMismatches() []TypeMismatch {
	if c.typeCheckingHasRun {
		return c.GetTypeValidator().GetMismatches()
	}
	panic(jutil.NewRuntimeException("Can't ask for type mismatches before type checking."))
}

// Java: Compiler.getTranspilationNamespace()
func (c *Compiler) GetTranspilationNamespace() rhino.StaticScope {
	if c.transpilationNamespace == nil {
		gn := NewGlobalNamespaceWithRoots(c, c.GetExternsRoot(), c.GetJsRoot())
		// Exclude user-defined code names, to save memory and help prevent misuse.
		gn.SetShouldTraverseScriptPredicate(
			func(script rhino.Node) bool { return script.IsFromExterns() || script.IsFirstChildOf(c.GetJsRoot()) })
		gn.EnsureGenerated() // ensure namespace is generated
		// wrap GlobalNamespace to prevent looking up extra information in GlobalNamespace
		c.transpilationNamespace = gn
	}

	return c.transpilationNamespace
}

// Result holds compilation results.
type Result struct {
	// Java's public final fields; the ones other Go packages read are
	// exported (PORTING.md §3).
	Success  bool
	Errors   []*JSError
	Warnings []*JSError

	variableMap             *VariableMap
	propertyMap             *VariableMap
	namedAnonFunctionMap    *VariableMap
	stringMap               *VariableMap                 // nil is null
	instrumentationMappings *VariableMap                 // nil is null
	SourceMap               *SourceMap                   // nil is null
	cssNames                *jutil.LinkedHashSet[string] // nil is null
	ExternExport            *string                      // nil is null
	idGeneratorMap          *string                      // nil is null
	TranspiledFiles         bool
}

// Java: Result(ImmutableList<JSError>, ImmutableList<JSError>, VariableMap, VariableMap, VariableMap, VariableMap, VariableMap, SourceMap, String, Set<String>, String, boolean)
func NewResultWithAllMaps(
	errors []*JSError,
	warnings []*JSError,
	variableMap *VariableMap,
	propertyMap *VariableMap,
	namedAnonFunctionMap *VariableMap,
	stringMap *VariableMap,
	instrumentationMappings *VariableMap,
	sourceMap *SourceMap,
	externExport *string,
	cssNames *jutil.LinkedHashSet[string],
	idGeneratorMap *string,
	transpiledFiles bool) *Result {
	return &Result{
		Success:                 len(errors) == 0,
		Errors:                  errors,
		Warnings:                warnings,
		variableMap:             variableMap,
		propertyMap:             propertyMap,
		namedAnonFunctionMap:    namedAnonFunctionMap,
		stringMap:               stringMap,
		instrumentationMappings: instrumentationMappings,
		SourceMap:               sourceMap,
		ExternExport:            externExport,
		cssNames:                cssNames,
		idGeneratorMap:          idGeneratorMap,
		TranspiledFiles:         transpiledFiles,
	}
}

// NewResult is @VisibleForTesting and deprecated in Java.
//
// Java: Result(ImmutableList<JSError>, ImmutableList<JSError>, VariableMap, VariableMap, VariableMap, SourceMap, String)
func NewResult(
	errors []*JSError,
	warnings []*JSError,
	variableMap *VariableMap,
	propertyMap *VariableMap,
	namedAnonFunctionMap *VariableMap,
	sourceMap *SourceMap,
	externExport *string) *Result {
	return NewResultWithAllMaps(
		errors,
		warnings,
		variableMap,
		propertyMap,
		namedAnonFunctionMap,
		/* stringMap= */ nil,
		/* instrumentationMappings= */ nil,
		sourceMap,
		externExport,
		/* cssNames= */ nil,
		/* idGeneratorMap= */ nil,
		/* transpiledFiles= */ false)
}

// Result_PruneResultForPartialCompilation returns an almost empty result
// that is more appropriate for a partial compilation.
//
// For a partial compilation we only care about errors and warnings. It is
// unnecessary to examine all of the other results.
//
// result is the full `Result` object provided by the compiler.
//
// Java: Result.pruneResultForPartialCompilation(Result)
func Result_PruneResultForPartialCompilation(result *Result) *Result {
	emptyVariableMap := NewVariableMap(jutil.NewLinkedHashMap[string, string]())
	return NewResultWithAllMaps(
		/* errors= */ result.Errors,
		/* warnings= */ result.Warnings,
		/* variableMap= */ emptyVariableMap,
		/* propertyMap= */ emptyVariableMap,
		/* namedAnonFunctionMap= */ emptyVariableMap,
		/* stringMap= */ nil,
		/* instrumentationMappings= */ emptyVariableMap,
		/* sourceMap= */ nil,
		/* externExport= */ jutil.StringPtr(""),
		/* cssNames= */ nil,
		/* idGeneratorMap= */ nil,
		/* transpiledFiles= */ false)
}
