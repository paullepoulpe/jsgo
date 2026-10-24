// Port of com/google/javascript/jscomp/VarCheck.java @ closure-compiler v20261021

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// VarCheck checks that all variables are declared, that file-private
// variables are accessed only in the file that declares them, and that any
// var references that cross chunk boundaries respect declared chunk
// dependencies.
type VarCheck struct {
	// Vars that were referenced in the externs without being declared in externs, even if they were
	// defined in code. These will be declared at the end of this pass.
	undefinedNamesFromExterns *jutil.LinkedHashSet[string]

	compiler AbstractCompiler

	// Whether this is the post-processing validity check.
	validityCheck bool

	// Whether extern checks emit error.
	strictExternCheck bool

	dupHandler *VarCheck_RedeclarationCheckHandler
}

// Java: VarCheck.UNDEFINED_VAR_ERROR
var VarCheck_UNDEFINED_VAR_ERROR = DiagnosticType_Error(
	"JSC_UNDEFINED_VARIABLE",
	"variable {0} is undeclared")

// Java: VarCheck.VIOLATED_CHUNK_DEP_ERROR
var VarCheck_VIOLATED_CHUNK_DEP_ERROR = DiagnosticType_Error(
	"JSC_VIOLATED_CHUNK_DEPENDENCY",
	"chunk {0} cannot reference {2}, defined in chunk {1}, since {1} loads after {0}")

// Java: VarCheck.MISSING_CHUNK_DEP_ERROR
var VarCheck_MISSING_CHUNK_DEP_ERROR = DiagnosticType_Warning(
	"JSC_MISSING_CHUNK_DEPENDENCY",
	"missing chunk dependency; chunk {0} should depend"+
		" on chunk {1} because it references {2}")

// Java: VarCheck.STRICT_CHUNK_DEP_ERROR
var VarCheck_STRICT_CHUNK_DEP_ERROR = DiagnosticType_Disabled(
	"JSC_STRICT_CHUNK_DEPENDENCY",
	// The newline below causes the JS compiler not to complain when the
	// referenced chunk's name changes because, for example, it's a
	// synthetic chunk.
	"cannot reference {2} because of a missing chunk dependency\n"+
		"defined in chunk {1}, referenced from chunk {0}")

// Java: VarCheck.NAME_REFERENCE_IN_EXTERNS_ERROR
var VarCheck_NAME_REFERENCE_IN_EXTERNS_ERROR = DiagnosticType_Warning(
	"JSC_NAME_REFERENCE_IN_EXTERNS",
	"accessing name {0} in externs has no effect."+
		" Perhaps you forgot to add a var keyword?")

// Java: VarCheck.UNDEFINED_EXTERN_VAR_ERROR
var VarCheck_UNDEFINED_EXTERN_VAR_ERROR = DiagnosticType_Warning(
	"JSC_UNDEFINED_EXTERN_VAR_ERROR",
	"name {0} is not defined in the externs.")

// Java: VarCheck.VAR_MULTIPLY_DECLARED_ERROR
var VarCheck_VAR_MULTIPLY_DECLARED_ERROR = DiagnosticType_Error(
	"JSC_VAR_MULTIPLY_DECLARED_ERROR",
	"Variable {0} declared more than once. First occurrence: {1}")

// Java: VarCheck.BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR
var VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR = DiagnosticType_Error(
	"JSC_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR",
	"Block-scoped variable {0} declared more than once. First occurrence: {1}")

// Java: VarCheck.VAR_ARGUMENTS_SHADOWED_ERROR
var VarCheck_VAR_ARGUMENTS_SHADOWED_ERROR = DiagnosticType_Error(
	"JSC_VAR_ARGUMENTS_SHADOWED_ERROR",
	"Shadowing \"arguments\" is not allowed")

// The arguments variable is special, in that it's declared in every local
// scope, but not explicitly declared.
//
// Java: VarCheck.ARGUMENTS
const varCheck_ARGUMENTS = "arguments"

// The IRFactory parse step adds references to this special name for @closureUnaware support, but
// it's not actually defined anywhere - just a compiler internal implementation detail.
//
// Java: VarCheck.JSCOMP_CLOSURE_UNAWARE_CODE_SHADOW_HOST_NAME
const varCheck_JSCOMP_CLOSURE_UNAWARE_CODE_SHADOW_HOST_NAME = "$jscomp_wrap_closure_unaware_code"

// Java: VarCheck.googForwardDeclare
var varCheck_googForwardDeclare = rhino.IR_Getprop(rhino.IR_Name("goog"), "forwardDeclare")

// Java: VarCheck(AbstractCompiler)
func NewVarCheck(compiler AbstractCompiler) *VarCheck {
	return NewVarCheckWithValidityCheck(compiler, false)
}

// Java: VarCheck(AbstractCompiler, boolean)
func NewVarCheckWithValidityCheck(compiler AbstractCompiler, validityCheck bool) *VarCheck {
	v := &VarCheck{}
	v.undefinedNamesFromExterns = jutil.NewLinkedHashSet[string]()
	v.compiler = compiler
	empty := ""
	// getErrorLevel returns a nullable CheckLevel; CheckLevel.Ptr() pointers are canonical, so
	// pointer identity is Java's enum identity.
	v.strictExternCheck = compiler.GetErrorLevel(
		JSError_MakeWithSourceLocation(&empty, 0, 0, VarCheck_UNDEFINED_EXTERN_VAR_ERROR)) == CheckLevel_ERROR.Ptr()
	v.validityCheck = validityCheck
	return v
}

// CreateScopeCreator creates the scope creator used by this pass. If not in
// validity check mode, use a RedeclarationCheckHandler to check var
// redeclarations.
//
// Java: VarCheck.createScopeCreator()
func (v *VarCheck) CreateScopeCreator() *SyntacticScopeCreator {
	if v.validityCheck {
		return NewSyntacticScopeCreator(v.compiler)
	} else {
		v.dupHandler = NewVarCheck_RedeclarationCheckHandler(v)
		return NewSyntacticScopeCreatorWithRedeclarationHandler(v.compiler, v.dupHandler)
	}
}

// Java: VarCheck.process(Node, Node)
func (v *VarCheck) Process(externs *rhino.Node, root *rhino.Node) {
	sc := v.CreateScopeCreator()
	var scopeCreator ScopeCreator = sc

	// Don't run externs-checking in sanity check mode. Normalization will
	// remove duplicate VAR declarations, which will make
	// externs look like they have assigns.
	if !v.validityCheck {
		NodeTraversal_Builder_().
			SetCompiler(v.compiler).
			SetCallback(NewVarCheck_NameRefInExternsCheck(v)).
			SetScopeCreator(scopeCreator).
			Traverse(externs)
	}

	NodeTraversal_Builder_().
		SetCompiler(v.compiler).
		SetCallback(v).
		SetScopeCreator(scopeCreator).
		TraverseRoots(externs, root)

	for varName := range v.undefinedNamesFromExterns.All() {
		v.CreateSynthesizedExternVar(varName, false /* isFromUndefinedCodeRef */)
	}

	if v.dupHandler != nil {
		v.dupHandler.RemoveDuplicates()
	}
}

func (*VarCheck) noShouldTraverse() {}
func (*VarCheck) onlyNameLeaves()   {}

// Java: VarCheck.shouldTraverse(NodeTraversal, Node, Node)
func (v *VarCheck) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	return true
}

// Java: VarCheck.visit(NodeTraversal, Node, Node)
func (v *VarCheck) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsName() {
		v.CheckName(t, n, parent)
	}
}

// CheckName validates that a NAME node does not refer to an undefined name.
//
// Java: VarCheck.checkName(NodeTraversal, Node, Node)
func (v *VarCheck) CheckName(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	varName := n.GetString()

	// Only a function can have an empty name.
	if varName == "" {
		// Name is optional for function expressions
		// x = function() {...}
		// Arrow functions are also expressions and cannot have a name
		// x = () => {...}
		// Member functions have an empty NAME node string, because the actual name is stored on the
		// MEMBER_FUNCTION_DEF object that contains the FUNCTION.
		// class C { foo() {...} }
		// x = { foo() {...} }
		jutil.CheckState(NodeUtil_IsFunctionExpression(parent) || NodeUtil_IsMethodDeclaration(parent))
		return
	}

	var_ := t.LookupVar(varName)

	// Check that the var has been declared.
	if var_ == nil {
		if (NodeUtil_IsFunctionExpression(parent) || NodeUtil_IsClassExpression(parent)) &&
			n.IsFirstChildOf(parent) {
			// e.g. [ function foo() {} ], it's okay if "foo" isn't defined in the
			// current scope.
			return
		}

		if NodeUtil_IsNonlocalModuleExportName(n) {
			// e.g. "export {a as b}" or "import {b as a} from './foo.js'
			// where b is defined in a module's export entries but not in any module scope.
			return
		}

		v.HandleUndeclaredVariableRef(t, n)
		t.GetScope().GetGlobalScope().Declare(varName, n, v.compiler.GetSynthesizedExternsInput())

		return
	}

	if var_.IsImplicitGoogNamespace() &&
		var_.GetImplicitGoogNamespaceStrength() == rhino.StaticSourceFile_SourceKind_WEAK &&
		VarCheck_StrengthOf(n) == rhino.StaticSourceFile_SourceKind_STRONG {
		// This use will be retained but its definition will be deleted.
		v.HandleUndeclaredVariableRef(t, n)
		return
	}

	currInput := t.GetInput()
	varInput := var_.GetInput()
	if currInput == varInput || currInput == nil || varInput == nil {
		// The variable was defined in the same file. This is fine.
		return
	}

	// Check chunk dependencies.
	currChunk := currInput.GetChunk()
	varChunk := varInput.GetChunk()
	chunkGraph := v.compiler.GetChunkGraph()
	if !v.validityCheck && varChunk != currChunk && varChunk != nil && currChunk != nil {
		if varChunk.IsWeak() {
			v.HandleUndeclaredVariableRef(t, n)
		}

		if chunkGraph.DependsOn(currChunk, varChunk) {
			// The chunk dependency was properly declared.
		} else {
			if t.InGlobalScope() {
				if chunkGraph.DependsOn(varChunk, currChunk) {
					// The variable reference violates a declared chunk dependency.
					t.Report(n, VarCheck_VIOLATED_CHUNK_DEP_ERROR, currChunk.GetName(), varChunk.GetName(), varName)
				} else {
					// The variable reference is between two chunks that have no dependency relationship.
					// This should probably be considered an error, but just issue a warning for now.
					t.Report(n, VarCheck_MISSING_CHUNK_DEP_ERROR, currChunk.GetName(), varChunk.GetName(), varName)
				}
			} else {
				t.Report(n, VarCheck_STRICT_CHUNK_DEP_ERROR, currChunk.GetName(), varChunk.GetName(), varName)
			}
		}
	}
}

// Java: VarCheck.strengthOf(Node)
func VarCheck_StrengthOf(n *rhino.Node) rhino.StaticSourceFile_SourceKind {
	source := n.GetStaticSourceFile()
	if source == nil {
		return rhino.StaticSourceFile_SourceKind_EXTERN
	}

	return source.GetKind()
}

// Java: VarCheck.handleUndeclaredVariableRef(NodeTraversal, Node)
func (v *VarCheck) HandleUndeclaredVariableRef(t *NodeTraversal, n *rhino.Node) {
	jutil.CheckState(n.IsName())

	varName := n.GetString()

	if n.GetParent().IsTypeOf() {
		// `typeof` is used for existence checks.
	} else if varName == varCheck_JSCOMP_CLOSURE_UNAWARE_CODE_SHADOW_HOST_NAME {
		// Compiler internal, we just create an extern for it.
	} else if v.strictExternCheck && t.GetInput().IsExtern() {
		// The extern checks are stricter, don't report a second error.
	} else {
		t.Report(n, VarCheck_UNDEFINED_VAR_ERROR, varName)
	}

	if v.validityCheck {
		// When the code is initially traversed, any undeclared variables are treated as
		// externs. During this sanity check, we ensure that all variables have either been
		// declared or marked as an extern. A failure at this point means that we have created
		// some variable/generated some code with an undefined reference.
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected variable " + varName))
	} else if !v.undefinedNamesFromExterns.Contains(varName) {
		// Skip this case if the name is already going to be added as an "undefined name from extern"
		// That declaration must take priority, to avoid the RemoveUnnecessarySyntheticExterns pass
		// accidentally treating the synthetic extern as unnecessary.
		v.CreateSynthesizedExternVar(varName, true /* isFromUndefinedCodeRef */)
	}
}

// Java: VarCheck.enterScope(NodeTraversal)
func (v *VarCheck) EnterScope(t *NodeTraversal) {}

// Java: VarCheck.exitScope(NodeTraversal)
func (v *VarCheck) ExitScope(t *NodeTraversal) {
	if !v.validityCheck && t.InGlobalScope() {
		scope := t.GetScope()
		// Add symbols that are known to be needed to the standard injected code (polyfills, etc).
		for requiredSymbol := range VarCheck_REQUIRED_SYMBOLS.All() {
			var_ := scope.GetVar(requiredSymbol)
			if var_ == nil {
				v.undefinedNamesFromExterns.Add(requiredSymbol)
			}
		}
	}
}

// VarCheck_REQUIRED_SYMBOLS is the list of symbols that must always be
// externed even if they are not referenced anywhere (yet). These are used by
// runtime libraries that might not be present when the first VarCheck runs.
//
// Java: VarCheck.REQUIRED_SYMBOLS
var VarCheck_REQUIRED_SYMBOLS = jutil.NewLinkedHashSet(
	// go/keep-sorted start
	"AggregateError",
	"Array",
	"Error",
	"Float32Array",
	"Function",
	"Infinity",
	"Iterator",
	"JSCOMPILER_PRESERVE", // added by CheckSideEffects
	"JSCompiler_renameProperty",
	"Map",
	"Math",
	"NaN",
	"Number",
	"Object",
	"Promise",
	"RangeError",
	"ReferenceError",
	"Reflect",
	"RegExp",
	"Set",
	"String",
	"SuppressedError",
	"Symbol",
	"TypeError",
	"WeakMap",
	"global",
	"globalThis",
	"isNaN",
	"parseFloat",
	"parseInt",
	"self",
	"undefined",
	"window",
	// go/keep-sorted end
)

// CreateSynthesizedExternVar creates a new variable in a synthetic script.
// This will prevent subsequent compiler passes from crashing.
//
// Java: VarCheck.createSynthesizedExternVar(String, boolean)
func (v *VarCheck) CreateSynthesizedExternVar(varName string, isFromUndefinedCodeRef bool) {
	nameNode := rhino.IR_Name(varName)

	// Mark the variable as constant if it matches the coding convention
	// for constant vars.
	// NOTE(nicksantos): honestly, I'm not sure how much this matters.
	// AFAIK, all people who use the CONST coding convention also
	// compile with undeclaredVars as errors. We have some test
	// cases for this configuration though, and it makes them happier.
	if v.compiler.GetCodingConvention().IsConstant(varName) {
		nameNode.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	}

	syntheticExternVar := rhino.IR_Var(nameNode)
	syntheticExternVar.SetIsSynthesizedUnfulfilledNameDeclaration(isFromUndefinedCodeRef)
	VarCheck_GetSynthesizedExternsRoot(v.compiler).AddChildToBack(syntheticExternVar)
	v.compiler.ReportChangeToEnclosingScope(syntheticExternVar)
}

// VarCheck_NameRefInExternsCheck is a check for name references in the
// externs inputs. These used to prevent a variable from getting renamed, but
// no longer have any effect.
//
// Java: VarCheck.NameRefInExternsCheck
type VarCheck_NameRefInExternsCheck struct {
	// Go-only: language-adapter — the enclosing VarCheck instance of the
	// inner class.
	outer *VarCheck
}

// Java: VarCheck.NameRefInExternsCheck()
func NewVarCheck_NameRefInExternsCheck(outer *VarCheck) *VarCheck_NameRefInExternsCheck {
	return &VarCheck_NameRefInExternsCheck{outer: outer}
}

// Java: VarCheck.NameRefInExternsCheck.shouldTraverse(NodeTraversal, Node, Node)
func (c *VarCheck_NameRefInExternsCheck) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	// Type summaries are generated from code rather than hand-written,
	// so warning about name references there would usually not be helpful.
	return !n.IsScript() || !NodeUtil_IsFromTypeSummary(n)
}

// Java: VarCheck.NameRefInExternsCheck.visit(NodeTraversal, Node, Node)
func (c *VarCheck_NameRefInExternsCheck) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsName() {
		switch parent.GetToken() {
		case rhino.Token_VAR,
			rhino.Token_LET,
			rhino.Token_CONST,
			rhino.Token_FUNCTION,
			rhino.Token_CLASS,
			rhino.Token_PARAM_LIST,
			rhino.Token_DEFAULT_VALUE,
			rhino.Token_ITER_REST,
			rhino.Token_OBJECT_REST,
			rhino.Token_ARRAY_PATTERN:
			// These are okay.
			return
		case rhino.Token_STRING_KEY:
			if parent.GetParent().IsObjectPattern() {
				return
			}
		case rhino.Token_GETPROP:
			if n == parent.GetFirstChild() {
				scope := t.GetScope()
				var_ := scope.GetVar(n.GetString())
				if var_ != nil {
					return
				}
				if parent.MatchesQualifiedNameNode(varCheck_googForwardDeclare) {
					// Allow using `goog.forwardDeclare` in the externs without an externs definition
					// of goog.
					return
				}
				t.Report(n, VarCheck_UNDEFINED_EXTERN_VAR_ERROR, n.GetString())
				c.outer.undefinedNamesFromExterns.Add(n.GetString())
			}
			return
		case rhino.Token_ASSIGN:
			// Don't warn for the "window.foo = foo;" nodes added by
			// DeclaredGlobalExternsOnWindow, nor for alias declarations
			// of the form "/** @const */ ns.Foo = Bar;"
			if n == parent.GetLastChild() &&
				n.IsQualifiedName() &&
				parent.GetFirstChild().IsQualifiedName() {
				return
			}
		case rhino.Token_NAME:
			// Don't warn for simple var assignments "/** @const */ var foo = bar;"
			// They are used to infer the types of namespace aliases.
			if NodeUtil_IsNameDeclaration(parent.GetParent()) {
				return
			}
		case rhino.Token_OR:
			// Don't warn for namespace declarations: "/** @const */ var ns = ns || {};"
			if NodeUtil_IsNamespaceDecl(parent.GetParent()) {
				return
			}
		default:
		}
		t.Report(n, VarCheck_NAME_REFERENCE_IN_EXTERNS_ERROR, n.GetString())
		scope := t.GetScope()
		var_ := scope.GetVar(n.GetString())
		if var_ == nil {
			c.outer.undefinedNamesFromExterns.Add(n.GetString())
		}
	}
}

// VarCheck_HasDuplicateDeclarationSuppression returns true if duplication
// warnings are suppressed on either n or origVar.
//
// Java: VarCheck.hasDuplicateDeclarationSuppression(AbstractCompiler, Node, Node)
func VarCheck_HasDuplicateDeclarationSuppression(
	compiler AbstractCompiler, n *rhino.Node, origVar *rhino.Node) bool {
	// For VarCheck and VariableReferenceCheck, variables in externs do not generate duplicate
	// warnings.
	if VarCheck_IsExternNamespace(n) {
		return true
	}
	return TypeValidator_HasDuplicateDeclarationSuppression(compiler, origVar)
}

// VarCheck_IsExternNamespace returns true if n is the name of a variable that
// declares a namespace in an externs file.
//
// Java: VarCheck.isExternNamespace(Node)
func VarCheck_IsExternNamespace(n *rhino.Node) bool {
	return n.GetParent().IsVar() && n.IsFromExterns() && NodeUtil_IsNamespaceDecl(n)
}

// VarCheck_RedeclarationCheckHandler is the handler for duplicate
// declarations.
//
// Java: VarCheck.RedeclarationCheckHandler
type VarCheck_RedeclarationCheckHandler struct {
	// Go-only: language-adapter — the enclosing VarCheck instance of the
	// inner class.
	outer *VarCheck

	dupDeclNodes []*rhino.Node
}

// Java: VarCheck.RedeclarationCheckHandler()
func NewVarCheck_RedeclarationCheckHandler(outer *VarCheck) *VarCheck_RedeclarationCheckHandler {
	return &VarCheck_RedeclarationCheckHandler{outer: outer}
}

// Java: VarCheck.RedeclarationCheckHandler.onRedeclaration(Scope, String, Node, CompilerInput)
func (h *VarCheck_RedeclarationCheckHandler) OnRedeclaration(
	s *Scope, name string, n *rhino.Node, input *CompilerInput) {
	parent := NodeUtil_GetDeclaringParent(n)
	origVar := s.GetVar(name)
	// origNode will be null for `arguments`, since there's no node that declares it.
	origNode := origVar.GetNode()
	var origParent *rhino.Node
	if origNode != nil {
		origParent = NodeUtil_GetDeclaringParent(origNode)
	}

	switch parent.GetToken() {
	case rhino.Token_CLASS, rhino.Token_CONST, rhino.Token_LET:
		h.outer.ReportBlockScopedMultipleDeclaration(n, name, origNode)
		return
	default:
	}

	if origParent != nil {
		switch origParent.GetToken() {
		case rhino.Token_CLASS, rhino.Token_CONST, rhino.Token_LET:
			h.outer.ReportBlockScopedMultipleDeclaration(n, name, origNode)
			return
		case rhino.Token_FUNCTION:
			// Redeclarations of functions in global scope are fairly common, so allow them
			// (at least for now).
			if !s.IsGlobal() && parent.IsFunction() {
				h.outer.ReportBlockScopedMultipleDeclaration(n, name, origNode)
				return
			}
		default:
		}
	}

	// Don't allow multiple variables to be declared at the top-level scope
	if s.IsGlobal() {
		if origParent.IsCatch() && parent.IsCatch() {
			// Okay, both are 'catch(x)' variables.
			return
		}

		allowDupe := VarCheck_HasDuplicateDeclarationSuppression(h.outer.compiler, n, origVar.GetNameNode())
		if VarCheck_IsExternNamespace(n) {
			h.dupDeclNodes = append(h.dupDeclNodes, parent)
			return
		}
		if !allowDupe {
			VarCheck_ReportVarMultiplyDeclared(h.outer.compiler, n, name, origNode)
		}
	} else if name == varCheck_ARGUMENTS &&
		!(NodeUtil_IsNameDeclaration(n.GetParent()) && n.IsName()) {
		// Disallow shadowing "arguments" as we can't handle with our current
		// scope modeling.
		h.outer.compiler.Report(JSError_MakeWithNode(n, VarCheck_VAR_ARGUMENTS_SHADOWED_ERROR))
	}
}

// Java: VarCheck.RedeclarationCheckHandler.removeDuplicates()
func (h *VarCheck_RedeclarationCheckHandler) RemoveDuplicates() {
	for _, n := range h.dupDeclNodes {
		parent := n.GetParent()
		if parent != nil {
			n.Detach()
			h.outer.compiler.ReportChangeToEnclosingScope(parent)
		}
	}
}

// Java: VarCheck.reportVarMultiplyDeclared(AbstractCompiler, Node, String, Node)
func VarCheck_ReportVarMultiplyDeclared(
	compiler AbstractCompiler, current *rhino.Node, name string, original *rhino.Node) {
	compiler.Report(JSError_MakeWithNode(current, VarCheck_VAR_MULTIPLY_DECLARED_ERROR, name, VarCheck_LocationOf(original)))
}

// Java: VarCheck.reportBlockScopedMultipleDeclaration(Node, String, Node)
func (v *VarCheck) ReportBlockScopedMultipleDeclaration(
	current *rhino.Node, name string, original *rhino.Node) {
	v.compiler.Report(
		JSError_MakeWithNode(
			current, VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR, name, VarCheck_LocationOf(original)))
}

// Java: VarCheck.locationOf(Node)
func VarCheck_LocationOf(n *rhino.Node) string {
	if n == nil {
		return "<unknown>"
	}
	return n.GetLocation()
}

// VarCheck_GetSynthesizedExternsRoot lazily creates a "new" externs root for
// undeclared variables.
//
// Java: VarCheck.getSynthesizedExternsRoot(AbstractCompiler)
func VarCheck_GetSynthesizedExternsRoot(compiler AbstractCompiler) *rhino.Node {
	return compiler.GetSynthesizedExternsInput().GetAstRoot(compiler)
}
