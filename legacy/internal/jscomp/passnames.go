// Port of com/google/javascript/jscomp/PassNames.java
// @ closure-compiler v20260415

package jscomp

// PassNames contains symbolic constants for compiler pass names used in
// multiple places. Mirrors Java's PassNames final class.
//
// Java's PassNames is a non-instantiable final class with a private no-arg
// constructor and 79 public static final String fields. Go idiomatically
// represents these as package-level constants in a const block.

const (
	PassNameAfterEarlyOptimizationLoop            = "afterEarlyOptimizationLoop"
	PassNameAfterMainOptimizations                = "afterMainOptimizations"
	PassNameAfterStandardChecks                   = "afterStandardChecks"
	PassNameAmbiguateProperties                   = "ambiguateProperties"
	PassNameAnalyzerChecks                        = "analyzerChecks"
	PassNameAngularPass                           = "angularPass"
	PassNameBeforeEarlyOptimizationLoop           = "beforeEarlyOptimizationLoop"
	PassNameBeforeEarlyOptimizationsTranspilation = "beforeEarlyOptimizationsTranspilation"
	PassNameBeforeStandardOptimizations           = "beforeStandardOptimizations"
	PassNameObfuscationPassMarker                 = "obfuscationPassMarker"
	PassNameBeforeMainOptimizations               = "beforeMainOptimizations"
	PassNameBeforeTypeChecking                    = "beforeTypeChecking"
	PassNameBeforeSerialization                   = "beforeSerialization"
	PassNameBeforeVariableRenaming                = "beforeVariableRenaming"
	PassNameBeforeExtractPrototypeMemberDeclarations = "beforeExtractPrototypeMemberDeclarations"
	PassNameCheckConformance                      = "checkConformance"
	PassNameCheckRegExp                           = "checkRegExp"
	PassNameCheckTypes                            = "checkTypes"
	PassNameCheckVariableReferences               = "checkVariableReferences"
	PassNameCheckVars                             = "checkVars"
	PassNameCoalesceVariableNames                 = "coalesceVariableNames"
	PassNameCollapseAnonymousFunctions            = "collapseAnonymousFunctions"
	PassNameCollapseObjectLiterals                = "collapseObjectLiterals"
	PassNameCollapseProperties                    = "collapseProperties"
	PassNameCollapseVariableDeclarations          = "collapseVariableDeclarations"
	PassNameConvertToDottedProperties             = "convertToDottedProperties"
	PassNameCreateModuleMap                       = "createModuleMap"
	PassNameCrossChunkCodeMotion                  = "crossChunkCodeMotion"
	PassNameCrossChunkMethodMotion                = "crossChunkMethodMotion"
	PassNameDeadAssignmentElimination             = "deadAssignmentsElimination"
	PassNameDeclaredGlobalExternsOnWindow         = "declaredGlobalExternsOnWindow"
	PassNameDeserializeCompilerState              = "deserializeCompilerState"
	PassNameDevirtualizeMethods                   = "devirtualizeMethods"
	PassNameDisambiguateProperties                = "disambiguateProperties"
	PassNameEs6NormalizeClasses                   = "es6NormalizeClasses"
	PassNameExploitAssign                         = "exploitAssign"
	PassNameExportTestFunctions                   = "exportTestFunctions"
	PassNameExternExports                         = "externExports"
	PassNameExtractPrototypeMemberDeclarations    = "extractPrototypeMemberDeclarations"
	PassNameFlowSensitiveInlineVariables          = "flowSensitiveInlineVariables"
	PassNameGatherGettersAndSetters               = "gatherGettersAndSetters"
	PassNameGatherModuleMetadata                  = "gatherModuleMetadata"
	PassNameGatherRawExports                      = "gatherRawExports"
	PassNameGenerateExports                       = "generateExports"
	PassNameInferConsts                           = "inferConsts"
	PassNameInferTypes                            = "inferTypes"
	PassNameInlineFunctions                       = "inlineFunctions"
	PassNameInlineProperties                      = "inlineProperties"
	PassNameInlineTypeAliases                     = "inlineTypeAliases"
	PassNameInlineVariables                       = "inlineVariables"
	PassNameLatePeepholeOptimizations             = "latePeepholeOptimizations"
	PassNameLintChecks                            = "lintChecks"
	PassNameMarkUnnormalized                      = "markUnnormalized"
	PassNameNormalize                             = "normalize"
	PassNameOptimizationsHalfwayPoint             = "optimizationsHalfwayPoint"
	PassNameOptimizeCalls                         = "optimizeCalls"
	PassNameParseInputs                           = "parseInputs"
	PassNamePeepholeOptimizations                 = "peepholeOptimizations"
	PassNameRewriteCommonJsModules                = "rewriteCommonJsModules"
	PassNameRewriteScriptsToEs6Modules            = "rewriteScriptsToEs6Modules"
	PassNameRemoveUnusedCode                      = "removeUnusedCode"
	PassNameReplaceIdGenerators                   = "replaceIdGenerators"
	PassNameReplaceMessages                       = "replaceMessages"
	PassNameResolveTypes                          = "resolveTypes"
	PassNameRewriteFunctionExpressions            = "rewriteFunctionExpressions"
	PassNameRenameProperties                      = "renameProperties"
	PassNameBeforeRenameProperties                = "beforeRenameProperties"
	PassNameStripSideEffectProtection             = "stripSideEffectProtection"
	PassNameWizPass                               = "wizPass"
)

// Go-only: language-adapter
// passNamesSet returns a set of all known pass-name constants. Java's PassNames
// class simply declares the constants; this Go helper allows pass infrastructure
// to enumerate or validate registered pass names.
func passNamesSet() map[string]bool {
	return map[string]bool{
		PassNameAfterEarlyOptimizationLoop:            true,
		PassNameAfterMainOptimizations:                true,
		PassNameAfterStandardChecks:                   true,
		PassNameAmbiguateProperties:                   true,
		PassNameAnalyzerChecks:                        true,
		PassNameAngularPass:                           true,
		PassNameBeforeEarlyOptimizationLoop:           true,
		PassNameBeforeEarlyOptimizationsTranspilation: true,
		PassNameBeforeStandardOptimizations:           true,
		PassNameObfuscationPassMarker:                 true,
		PassNameBeforeMainOptimizations:               true,
		PassNameBeforeTypeChecking:                    true,
		PassNameBeforeSerialization:                   true,
		PassNameBeforeVariableRenaming:                true,
		PassNameBeforeExtractPrototypeMemberDeclarations: true,
		PassNameCheckConformance:                      true,
		PassNameCheckRegExp:                           true,
		PassNameCheckTypes:                            true,
		PassNameCheckVariableReferences:               true,
		PassNameCheckVars:                             true,
		PassNameCoalesceVariableNames:                 true,
		PassNameCollapseAnonymousFunctions:            true,
		PassNameCollapseObjectLiterals:                true,
		PassNameCollapseProperties:                    true,
		PassNameCollapseVariableDeclarations:          true,
		PassNameConvertToDottedProperties:             true,
		PassNameCreateModuleMap:                       true,
		PassNameCrossChunkCodeMotion:                  true,
		PassNameCrossChunkMethodMotion:                true,
		PassNameDeadAssignmentElimination:             true,
		PassNameDeclaredGlobalExternsOnWindow:         true,
		PassNameDeserializeCompilerState:              true,
		PassNameDevirtualizeMethods:                   true,
		PassNameDisambiguateProperties:                true,
		PassNameEs6NormalizeClasses:                   true,
		PassNameExploitAssign:                         true,
		PassNameExportTestFunctions:                   true,
		PassNameExternExports:                         true,
		PassNameExtractPrototypeMemberDeclarations:    true,
		PassNameFlowSensitiveInlineVariables:          true,
		PassNameGatherGettersAndSetters:               true,
		PassNameGatherModuleMetadata:                  true,
		PassNameGatherRawExports:                      true,
		PassNameGenerateExports:                       true,
		PassNameInferConsts:                           true,
		PassNameInferTypes:                            true,
		PassNameInlineFunctions:                       true,
		PassNameInlineProperties:                      true,
		PassNameInlineTypeAliases:                     true,
		PassNameInlineVariables:                       true,
		PassNameLatePeepholeOptimizations:             true,
		PassNameLintChecks:                            true,
		PassNameMarkUnnormalized:                      true,
		PassNameNormalize:                             true,
		PassNameOptimizationsHalfwayPoint:             true,
		PassNameOptimizeCalls:                         true,
		PassNameParseInputs:                           true,
		PassNamePeepholeOptimizations:                 true,
		PassNameRewriteCommonJsModules:                true,
		PassNameRewriteScriptsToEs6Modules:            true,
		PassNameRemoveUnusedCode:                      true,
		PassNameReplaceIdGenerators:                   true,
		PassNameReplaceMessages:                       true,
		PassNameResolveTypes:                          true,
		PassNameRewriteFunctionExpressions:            true,
		PassNameRenameProperties:                      true,
		PassNameBeforeRenameProperties:                true,
		PassNameStripSideEffectProtection:             true,
		PassNameWizPass:                               true,
	}
}

// Go-only: language-adapter
// IsKnownPassName reports whether name is a registered pass-name constant.
// Mirrors Java's PassNames class as a whole: if a name is declared there, it
// is a "known" pass name. Go exposes this as a function for use by pass
// infrastructure that needs to validate pass names at runtime.
func IsKnownPassName(name string) bool {
	return passNamesSet()[name]
}
