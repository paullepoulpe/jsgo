// Port of com/google/javascript/jscomp/ValidityCheck.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ValidityCheck is a compiler pass that verifies the structure of the AST
// conforms to a number of invariants. Because this can add a lot of
// overhead, we only run this in development mode.
type ValidityCheck struct {
	compiler      AbstractCompiler
	changeTracker *ChangeTracker
	astValidator  *AstValidator
}

var ValidityCheck_EXTERN_PROPERTIES_CHANGED = DiagnosticType_Error(
	"JSC_EXTERN_PROPERTIES_CHANGED",
	"Internal compiler error. Extern properties modified from:\n{0}\nto:\n{1}")

// Java: ValidityCheck(AbstractCompiler)
func NewValidityCheck(compiler AbstractCompiler) *ValidityCheck {
	v := &ValidityCheck{}
	v.compiler = compiler
	v.changeTracker = compiler.GetChangeTracker()
	v.astValidator = NewAstValidator(compiler)
	return v
}

// Java: ValidityCheck.process(Node, Node)
func (v *ValidityCheck) Process(externs *rhino.Node, root *rhino.Node) {
	v.CheckAst(externs, root)
	v.CheckNormalization(externs, root)
	v.CheckVars(externs, root)
	v.CheckExternProperties(externs)
}

// CheckAst checks that the AST is structurally accurate.
//
// Java: ValidityCheck.checkAst(Node, Node)
func (v *ValidityCheck) CheckAst(externs *rhino.Node, root *rhino.Node) {
	v.astValidator.ValidateCodeRoot(externs)
	v.astValidator.ValidateCodeRoot(root)
}

// Java: ValidityCheck.checkVars(Node, Node)
func (v *ValidityCheck) CheckVars(externs *rhino.Node, root *rhino.Node) {
	if v.compiler.GetLifeCycleStage().IsNormalized() {
		// TODO(rishipal): Why is VarCheck only run when AST is normalized?
		(NewVarCheckWithValidityCheck(v.compiler, true)).Process(externs, root)
	}
}

// CheckNormalization verifies that the normalization pass does nothing on an
// already-normalized tree.
//
// Java: ValidityCheck.checkNormalization(Node, Node)
func (v *ValidityCheck) CheckNormalization(externs *rhino.Node, root *rhino.Node) {
	// Verify nothing has inappropriately denormalize the AST.
	var handler CodeChangeHandler = NewForbiddenChange()
	v.changeTracker.AddChangeHandler(handler)

	// TODO(johnlenz): Change these normalization checks Preconditions and
	// Exceptions into Errors so that it is easier to find the root cause
	// when there are cascading issues.
	if v.compiler.GetLifeCycleStage().IsNormalized() {
		Normalize_Builder_(v.compiler).AssertOnChange(true).Build().Process(externs, root)

		if v.compiler.GetLifeCycleStage().IsNormalizedUnobfuscated() {
			checkUserDeclarations := true
			var pass CompilerPass = NewValidityCheck_VerifyConstants(v.compiler, checkUserDeclarations)
			pass.Process(externs, root)
		}
	}

	v.changeTracker.RemoveChangeHandler(handler)
}

// Java: ValidityCheck.checkExternProperties(Node)
func (v *ValidityCheck) CheckExternProperties(externs *rhino.Node) {
	externProperties := v.compiler.GetExternProperties()
	if externProperties == nil {
		// GatherExternProperties hasn't run yet. Don't report a violation.
		return
	}
	NewGatherExternProperties(v.compiler, GatherExternProperties_Mode_CHECK).Process(externs, nil)
	if !v.compiler.GetExternProperties().Equals(externProperties) {
		v.compiler.Report(
			JSError_Make(
				ValidityCheck_EXTERN_PROPERTIES_CHANGED,
				externProperties.String(),
				v.compiler.GetExternProperties().String()))
		// Throw an exception, so that the infrastructure will tell us which pass violated the check.
		panic(jutil.NewIllegalStateExceptionMsg(
			"Validity Check failed: Extern properties changed from:\n" +
				externProperties.String() +
				"\nto:\n" +
				v.compiler.GetExternProperties().String()))
	}
}

// ValidityCheck_VerifyConstants walks the AST tree and verifies that
// constant names are used consistently.
type ValidityCheck_VerifyConstants struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler              AbstractCompiler
	checkUserDeclarations bool

	constantMap *jutil.LinkedHashMap[string, bool]
}

// Java: ValidityCheck.VerifyConstants(AbstractCompiler, boolean)
func NewValidityCheck_VerifyConstants(compiler AbstractCompiler, checkUserDeclarations bool) *ValidityCheck_VerifyConstants {
	c := &ValidityCheck_VerifyConstants{}
	c.constantMap = jutil.NewLinkedHashMap[string, bool]()
	c.compiler = compiler
	c.checkUserDeclarations = checkUserDeclarations
	return c
}

// Java: ValidityCheck.VerifyConstants.process(Node, Node)
func (c *ValidityCheck_VerifyConstants) Process(externs *rhino.Node, root *rhino.Node) {
	externsAndJs := root.GetParent()
	jutil.CheckState(externsAndJs != nil)
	jutil.CheckState(externsAndJs.HasChild(externs))
	NodeTraversal_TraverseRoots(c.compiler, c, externs, root)
}

// Java: ValidityCheck.VerifyConstants.visit(NodeTraversal, Node, Node)
func (c *ValidityCheck_VerifyConstants) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsName() {
		name := n.GetString()
		if n.GetString() == "" {
			return
		}

		isConst := n.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME)
		if c.checkUserDeclarations {
			expectedConst := false
			convention := c.compiler.GetCodingConvention()
			if NodeUtil_IsConstantName(n) || NodeUtil_IsConstantByConvention(convention, n) {
				expectedConst = true
			} else {
				expectedConst = false

				var info *rhino.JSDocInfo
				var_ := t.GetScope().GetVar(n.GetString())
				if var_ != nil {
					info = var_.GetJSDocInfo()
				}

				if info != nil && info.IsConstant() {
					expectedConst = true
				} else {
					expectedConst = false
				}
			}

			if expectedConst {
				jutil.CheckStateMsg1(expectedConst == isConst, "The name %s is not annotated as constant.", name)
			} else {
				jutil.CheckStateMsg1(
					expectedConst == isConst, "The name %s should not be annotated as constant.", name)
			}
		}

		value, ok := c.constantMap.Lookup(name)
		if !ok {
			c.constantMap.Put(name, isConst)
		} else if value != isConst {
			panic(jutil.NewIllegalStateExceptionMsg(
				"The name " +
					name +
					" is not consistently annotated as constant. Current constness: " +
					jutil.ToJavaString(isConst) +
					", previous constness: " +
					jutil.ToJavaString(value) +
					"\n\nAt node: " +
					n.String()))
		}
	}
}
