// Port of com/google/javascript/jscomp/J2clPass.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/J2clAssertRemovalPass.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/J2clStringValueOfRewriterPass.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/J2clUndefinedChecksRewriterPass.java @ closure-compiler v20261021

/*
 * Copyright 2015 The Closure Compiler Authors.
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
 * Copyright 2017 The Closure Compiler Authors.
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
 * Copyright 2018 The Closure Compiler Authors.
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
 * Copyright 2025 The Closure Compiler Authors.
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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// J2clSourceUtils holds static helpers for recognizing J2CL sources.
type J2clSourceUtils struct{}

// J2clSourceUtils_IsJ2clSource: sourceFileName is a nullable Java String.
//
// Java: J2clSourceUtils.isJ2clSource(String)
func J2clSourceUtils_IsJ2clSource(sourceFileName *string) bool {
	return sourceFileName != nil && strings.HasSuffix(*sourceFileName, ".java.js")
}

// Java: J2clSourceUtils.isJ2clSource(Node)
func J2clSourceUtils_IsJ2clSourceOfNode(scriptNode rhino.Node) bool {
	jutil.CheckArgument(scriptNode.IsScript())
	return J2clSourceUtils_IsJ2clSource(scriptNode.GetSourceFileName())
}

// Java: J2clSourceUtils.isInsideJ2clSource(Node)
func J2clSourceUtils_IsInsideJ2clSource(node rhino.Node) bool {
	var scriptNode rhino.Node
	if node.IsScript() {
		scriptNode = node
	} else {
		scriptNode = NodeUtil_GetEnclosingScript(node)
	}
	if scriptNode.IsNil() {
		return false
	}
	return J2clSourceUtils_IsJ2clSourceOfNode(scriptNode)
}

// J2clSourceFileChecker checks whether there are J2CL generated source files
// with pattern "*.java.js".
//
// Java: final class J2clSourceFileChecker implements CompilerPass
type J2clSourceFileChecker struct {
	compiler AbstractCompiler
}

// Java: J2clSourceFileChecker(AbstractCompiler)
func NewJ2clSourceFileChecker(compiler AbstractCompiler) *J2clSourceFileChecker {
	return &J2clSourceFileChecker{compiler: compiler}
}

// Java: J2clSourceFileChecker.hasJ2cl(Node)
func j2clSourceFileChecker_hasJ2cl(root rhino.Node) bool {
	for script := root.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
		jutil.CheckState(script.IsScript())
		if J2clSourceUtils_IsJ2clSourceOfNode(script) {
			return true
		}
	}
	return false
}

// Java: J2clSourceFileChecker.process(Node, Node)
func (c *J2clSourceFileChecker) Process(externs rhino.Node, root rhino.Node) {
	if j2clSourceFileChecker_hasJ2cl(root) {
		J2clSourceFileChecker_MarkToRunJ2clPasses(c.compiler)
	}
}

// Java: J2clSourceFileChecker.markToRunJ2clPasses(AbstractCompiler)
func J2clSourceFileChecker_MarkToRunJ2clPasses(compiler AbstractCompiler) {
	compiler.SetRunJ2clPasses(true)
}

// J2clSourceFileChecker_ShouldRunJ2clPasses indicates whether it should run
// future J2CL passes with information from the compiler. For example, if the
// compiler's HAS_J2CL annotation is false, it should.
//
// Java: J2clSourceFileChecker.shouldRunJ2clPasses(AbstractCompiler)
func J2clSourceFileChecker_ShouldRunJ2clPasses(compiler AbstractCompiler) bool {
	return compiler.RunJ2clPasses()
}

const j2clPass_ALL_CLASS_FILE_NAMES = "*"

// J2clPass is a normalization pass to inline some J2CL calls to enable other
// optimizations.
//
// Java: public class J2clPass implements CompilerPass
type J2clPass struct {
	compiler           AbstractCompiler
	safeNameIdSupplier func() string
	inliner            *J2clPass_ClassStaticFunctionsInliner
}

// Java: J2clPass(AbstractCompiler)
func NewJ2clPass(compiler AbstractCompiler) *J2clPass {
	return &J2clPass{
		compiler:           compiler,
		safeNameIdSupplier: compiler.GetUniqueNameIdSupplier(),
	}
}

var (
	j2clPass_defaultInlinings = []*J2clPass_InlineSpec{
		/*
		 * Inline functions in Arrays that take references to static $isInstance() functions. This
		 * ensures that the references will be fully qualified to work with collapse properties.
		 */
		NewJ2clPass_InlineSpec(
			"Arrays.impl.java.js",
			map[string]bool{
				"$create":                true,
				"$createWithInitializer": true,
				"$init":                  true,
				"$instanceIsOfType":      true,
				"$castTo":                true,
				"$stampType":             true,
			},
			FunctionInjector_InliningMode_DIRECT),
		NewJ2clPass_InlineSpec("Casts.impl.java.js", map[string]bool{"$to": true}, FunctionInjector_InliningMode_DIRECT),
		/*
		 * Inlines all Interface.$markImplementor(FooClass) metaclass calls so that FooClass and others
		 * like it are not unnecessarily retained and so that static analysis of interface instanceof
		 * calls becomes possible.
		 *
		 * Note that this pass should NOT be restricted to j2cl .java.js files because JavaScript code
		 * implementing Java interfaces (not recommended but widely used in xplat) needs calls to
		 * $markImplementor.
		 */
		NewJ2clPass_InlineSpec(
			j2clPass_ALL_CLASS_FILE_NAMES, map[string]bool{"$markImplementor": true}, FunctionInjector_InliningMode_BLOCK),
		/*
		 * Inlines class metadata calls so they become optimizable and avoids escaping of constructor.
		 */
		NewJ2clPass_InlineSpec(
			"Util.impl.java.js",
			map[string]bool{
				"$setClassMetadata":             true,
				"$setClassMetadataForInterface": true,
				"$setClassMetadataForEnum":      true,
				"$setClassMetadataForPrimitive": true,
			},
			FunctionInjector_InliningMode_BLOCK),
	}
	j2clPass_defaultFnNamesToInline = func() map[string]bool {
		s := make(map[string]bool, 12)
		for _, b := range j2clPass_defaultInlinings {
			// order-independent: builds a set
			for name := range b.fnNamesToInline {
				s[name] = true
			}
		}
		return s
	}()
)

// Java: J2clPass.process(Node, Node)
func (p *J2clPass) Process(externs rhino.Node, root rhino.Node) {
	if !J2clSourceFileChecker_ShouldRunJ2clPasses(p.compiler) {
		return
	}

	p.getOrCreateInliner(root, j2clPass_defaultInlinings).run()
}

// J2clPass_InlineSpec is the private record InlineSpec.
//
// Java: private record J2clPass.InlineSpec(String, Set<String>, InliningMode)
type J2clPass_InlineSpec struct {
	classFileName   string
	fnNamesToInline map[string]bool
	inliningMode    FunctionInjector_InliningMode
}

// Java: J2clPass.InlineSpec(String, Set<String>, InliningMode)
func NewJ2clPass_InlineSpec(classFileName string, fnNamesToInline map[string]bool, inliningMode FunctionInjector_InliningMode) *J2clPass_InlineSpec {
	return &J2clPass_InlineSpec{classFileName: classFileName, fnNamesToInline: fnNamesToInline, inliningMode: inliningMode}
}

// J2clPass_ClassStaticFunctionsInliner collects references to certain
// function definitions in a certain class and then inlines fully qualified
// static method calls to those functions anywhere in the program.
//
// Assumes that the set of provided short function names will not collide
// with any of the collected fully qualified function names once the module
// prefix has been added.
//
// Java: private class J2clPass.ClassStaticFunctionsInliner
type J2clPass_ClassStaticFunctionsInliner struct {
	outer                      *J2clPass
	fnsToInlineByQualifiedName map[string]*J2clPass_ClassStaticFunctionsInliner_InlinableFunction
	injector                   *FunctionInjector
	root                       rhino.Node
	inlineSpecs                []*J2clPass_InlineSpec
	collector                  J2clPass_ClassStaticFunctionsInliner_FunctionDefsCollector
	callInliner                J2clPass_ClassStaticFunctionsInliner_StaticCallInliner
}

// J2clPass_ClassStaticFunctionsInliner_InlinableFunction is the private
// record InlinableFunction.
//
// Java: private record J2clPass.ClassStaticFunctionsInliner.InlinableFunction(Node, InliningMode)
type J2clPass_ClassStaticFunctionsInliner_InlinableFunction struct {
	impl         rhino.Node
	inliningMode FunctionInjector_InliningMode
}

// Java: J2clPass.ClassStaticFunctionsInliner.InlinableFunction(Node, InliningMode)
func NewJ2clPass_ClassStaticFunctionsInliner_InlinableFunction(impl rhino.Node, inliningMode FunctionInjector_InliningMode) *J2clPass_ClassStaticFunctionsInliner_InlinableFunction {
	return &J2clPass_ClassStaticFunctionsInliner_InlinableFunction{impl: impl, inliningMode: inliningMode}
}

// Java: J2clPass.ClassStaticFunctionsInliner.InlinableFunction.impl()
func (f *J2clPass_ClassStaticFunctionsInliner_InlinableFunction) Impl() rhino.Node { return f.impl }

// Java: J2clPass.ClassStaticFunctionsInliner.InlinableFunction.inliningMode()
func (f *J2clPass_ClassStaticFunctionsInliner_InlinableFunction) InliningMode() FunctionInjector_InliningMode {
	return f.inliningMode
}

// String is the record's toString().
//
// Java: J2clPass.ClassStaticFunctionsInliner.InlinableFunction.toString()
func (f *J2clPass_ClassStaticFunctionsInliner_InlinableFunction) String() string {
	implStr := "null"
	if !f.impl.IsNil() {
		implStr = f.impl.String()
	}
	return "InlinableFunction[impl=" + implStr + ", inliningMode=" + f.inliningMode.String() + "]"
}

func (p *J2clPass) getOrCreateInliner(root rhino.Node, inlineSpecs []*J2clPass_InlineSpec) *J2clPass_ClassStaticFunctionsInliner {
	c := p.inliner
	if c == nil {
		c = &J2clPass_ClassStaticFunctionsInliner{
			outer:                      p,
			fnsToInlineByQualifiedName: make(map[string]*J2clPass_ClassStaticFunctionsInliner_InlinableFunction),
		}
		c.injector = NewFunctionInjector_Builder(p.compiler).
			SafeNameIdSupplier(p.safeNameIdSupplier).
			AssumeStrictThis(true).
			Build()
		c.injector.SetKnownConstantFunctions(j2clPass_defaultFnNamesToInline)
		c.collector.outer = c
		c.callInliner.outer = c
		p.inliner = c
	} else {
		clear(c.fnsToInlineByQualifiedName)
	}
	c.root = root
	c.inlineSpecs = inlineSpecs
	return c
}

// Java: J2clPass.ClassStaticFunctionsInliner.run()
func (c *J2clPass_ClassStaticFunctionsInliner) run() {
	NodeTraversal_Traverse(c.outer.compiler, c.root, &c.collector)
	if len(c.fnsToInlineByQualifiedName) > 0 {
		NodeTraversal_Traverse(c.outer.compiler, c.root, &c.callInliner)
	}
}

// Java: J2clPass.ClassStaticFunctionsInliner.findMatchingSpecsForFile(String)
func (c *J2clPass_ClassStaticFunctionsInliner) findMatchingSpecsForFile(
	sourceFileName string, hasSource bool, dst []*J2clPass_InlineSpec,
) []*J2clPass_InlineSpec {
	matchingSpecs := dst[:0]
	for _, inlineSpec := range c.inlineSpecs {
		if inlineSpec.classFileName == j2clPass_ALL_CLASS_FILE_NAMES {
			matchingSpecs = append(matchingSpecs, inlineSpec)
			continue
		}
		if !hasSource {
			panic(jutil.NewNullPointerException())
		}
		if strings.HasSuffix(sourceFileName, inlineSpec.classFileName) {
			matchingSpecs = append(matchingSpecs, inlineSpec)
		}
	}
	return matchingSpecs
}

// Java: private class J2clPass.ClassStaticFunctionsInliner.FunctionDefsCollector implements NodeTraversal.Callback
type J2clPass_ClassStaticFunctionsInliner_FunctionDefsCollector struct {
	outer *J2clPass_ClassStaticFunctionsInliner
	// This list is initialized every time the node traversal enters a new script. It contains
	// all InlineSpecs whose classFileName matches the script's filename.
	currFileInlineSpecs []*J2clPass_InlineSpec
}

// Java: J2clPass.ClassStaticFunctionsInliner.FunctionDefsCollector.shouldTraverse(NodeTraversal, Node, Node)
func (f *J2clPass_ClassStaticFunctionsInliner_FunctionDefsCollector) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsScript() {
		var srcName string
		hasSrc := false
		if s := n.GetSourceFileName(); s != nil {
			srcName = *s
			hasSrc = true
		}
		f.currFileInlineSpecs = f.outer.findMatchingSpecsForFile(srcName, hasSrc, f.currFileInlineSpecs)
		return len(f.currFileInlineSpecs) != 0
	}
	return true
}

// Java: J2clPass.ClassStaticFunctionsInliner.FunctionDefsCollector.visit(NodeTraversal, Node, Node)
func (f *J2clPass_ClassStaticFunctionsInliner_FunctionDefsCollector) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	// Only look at named function declarations that are fully qualified.
	var fnName string
	var fnImpl rhino.Node
	if assign, ok := ast.AsAssign(n); ok {
		// TODO(b/69730966): Delete this branch when ES5 syntax support is no longer needed.
		val := assign.Value()
		if _, isFn := val.AsFunction(); !isFn {
			return
		}
		qualifiedProp, isGetProp := assign.Target().AsGetProp()
		if !isGetProp || !qualifiedProp.Node().IsQualifiedName() {
			return
		}
		fnName = qualifiedProp.Property()
		fnImpl = val.Node()
	} else if n.IsMemberFunctionDef() {
		fnName = n.GetString()
		fnImpl = n.GetLastChild()
	} else {
		return
	}
	for _, inlineSpec := range f.currFileInlineSpecs {
		if !inlineSpec.fnNamesToInline[fnName] {
			continue
		}
		// Then store a reference to it.
		inlinableFunction := NewJ2clPass_ClassStaticFunctionsInliner_InlinableFunction(fnImpl, inlineSpec.inliningMode)
		var qualifiedFnName string
		switch n.GetToken() {
		case rhino.Token_ASSIGN:
			qname, ok := n.GetFirstChild().GetQualifiedNameString()
			if !ok {
				panic(jutil.NewNullPointerException())
			}
			qualifiedFnName = qname
		case rhino.Token_MEMBER_FUNCTION_DEF:
			bestName := NodeUtil_GetBestLValueName(n)
			if bestName == nil {
				panic(jutil.NewNullPointerException())
			}
			qualifiedFnName = *bestName
		default:
			panic(jutil.NewAssertionErrorNoMessage())
		}
		previousValue := f.outer.fnsToInlineByQualifiedName[qualifiedFnName]
		f.outer.fnsToInlineByQualifiedName[qualifiedFnName] = inlinableFunction
		jutil.CheckStateMsg2(
			previousValue == nil,
			"expected each function to match at most one InliningSpec, but found:\n%s\n%s",
			previousValue,
			inlinableFunction)
	}
}

// Java: private class J2clPass.ClassStaticFunctionsInliner.StaticCallInliner extends AbstractPostOrderCallback
type J2clPass_ClassStaticFunctionsInliner_StaticCallInliner struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *J2clPass_ClassStaticFunctionsInliner
}

// Java: J2clPass.ClassStaticFunctionsInliner.StaticCallInliner.visit(NodeTraversal, Node, Node)
func (s *J2clPass_ClassStaticFunctionsInliner_StaticCallInliner) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	// Only look at method calls
	call, ok := ast.AsCall(n)
	if !ok {
		return
	}

	// ... that are fully qualified
	calleeProp, ok := call.Callee().AsGetProp()
	if !ok {
		return
	}
	qualifiedFnName, ok := calleeProp.Node().GetQualifiedNameString()
	if !ok {
		return
	}

	// ... and that reference a function definition we want to inline
	fn := s.outer.fnsToInlineByQualifiedName[qualifiedFnName]
	if fn == nil {
		return
	}
	fnName := calleeProp.Property()
	fnImpl := fn.Impl()
	inliningMode := fn.InliningMode()

	// Ensure that the function only has a single return statement when direct inlining.
	if inliningMode == FunctionInjector_InliningMode_DIRECT &&
		!NodeUtil_GetFunctionBody(fnImpl).GetFirstChild().IsReturn() {
		panic(jutil.NewIllegalStateExceptionMsg(
			"Attempted to direct inline function " +
				qualifiedFnName +
				", but function is not a simple return."))
	}

	// Otherwise inline the call.
	// Note: This pass has to run before normalization, so we must use the unsafeInline method.
	// It is safe because these are strictly controlled trivial bootstrap methods that are
	// written with inlining in mind (e.g. doesn't read/write local/global variables).
	// TODO(goktug): Add a check that will ensure safety of this.
	inlinedCall := s.outer.injector.UnsafeInline(
		NewFunctionInjector_Reference(n, t.GetScope(), t.GetChunk(), inliningMode), fnName, fnImpl)
	// Avoid overridding original source information with the helper classes source information.
	// For example; we want a cast to point related Java statement instead of the Casts utility.
	inlinedCall.SrcrefTree(n)
	t.GetCompiler().ReportChangeToEnclosingScope(inlinedCall)
}

// J2clAssertRemovalPass is an optimization pass to remove J2CL
// Asserts.$assert.
//
// Java: public class J2clAssertRemovalPass extends AbstractPostOrderCallback implements CompilerPass
type J2clAssertRemovalPass struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler
}

// Java: J2clAssertRemovalPass(AbstractCompiler)
func NewJ2clAssertRemovalPass(compiler AbstractCompiler) *J2clAssertRemovalPass {
	return &J2clAssertRemovalPass{compiler: compiler}
}

// Java: J2clAssertRemovalPass.process(Node, Node)
func (p *J2clAssertRemovalPass) Process(externs rhino.Node, root rhino.Node) {
	if !J2clSourceFileChecker_ShouldRunJ2clPasses(p.compiler) {
		return
	}

	NodeTraversal_Traverse(p.compiler, root, p)
}

// Java: J2clAssertRemovalPass.visit(NodeTraversal, Node, Node)
func (p *J2clAssertRemovalPass) Visit(t *NodeTraversal, node rhino.Node, parent rhino.Node) {
	if j2clAssertRemovalPass_isAssertionCall(node) {
		NodeUtil_DeleteFunctionCall(node, p.compiler)
	}
}

// Java: J2clAssertRemovalPass.isAssertionCall(Node)
func j2clAssertRemovalPass_isAssertionCall(node rhino.Node) bool {
	call, ok := ast.AsCall(node)
	return ok && j2clAssertRemovalPass_isAssertionMethodName(call.Callee().Node())
}

// Java: J2clAssertRemovalPass.isAssertionMethodName(Node)
func j2clAssertRemovalPass_isAssertionMethodName(fnName rhino.Node) bool {
	if !fnName.IsQualifiedName() {
		return false
	}
	originalQname := *fnName.GetOriginalQualifiedName()
	return strings.Contains(originalQname, "Asserts") && strings.Contains(originalQname, ".$assert")
}

const j2clStringValueOfRewriterPass_METHOD_NAME = "module$exports$java$lang$String$impl.m_valueOf__java_lang_Object__java_lang_String"

// J2clStringValueOfRewriterPass is an optimization pass to rewrite J2CL
// String.valueOf.
//
// Java: final class J2clStringValueOfRewriterPass extends AbstractPeepholeOptimization
type J2clStringValueOfRewriterPass struct {
	AbstractPeepholeOptimizationBase
	shouldRunJ2clPasses bool
}

// Java: J2clStringValueOfRewriterPass()
func NewJ2clStringValueOfRewriterPass() *J2clStringValueOfRewriterPass {
	return &J2clStringValueOfRewriterPass{}
}

func (p *J2clStringValueOfRewriterPass) peepholeConfig() peepholeOptConfig {
	if p.shouldRunJ2clPasses {
		return peepholeOptConfig{
			tokens:         peepholeCallOnlyTokens,
			isStandard:     true,
			isJ2clCallPass: true,
		}
	}
	return peepholeOptConfig{isStandard: true}
}

func (p *J2clStringValueOfRewriterPass) cloneForWorker() AbstractPeepholeOptimization {
	return NewJ2clStringValueOfRewriterPass()
}

// Java: J2clStringValueOfRewriterPass.beginTraversal(AbstractCompiler)
func (p *J2clStringValueOfRewriterPass) BeginTraversal(compiler AbstractCompiler) {
	p.AbstractPeepholeOptimizationBase.BeginTraversal(compiler)
	p.shouldRunJ2clPasses = J2clSourceFileChecker_ShouldRunJ2clPasses(compiler)
}

// Java: J2clStringValueOfRewriterPass.optimizeSubtree(Node)
func (p *J2clStringValueOfRewriterPass) OptimizeSubtree(node rhino.Node) rhino.Node {
	if !p.shouldRunJ2clPasses {
		return node
	}
	call, ok := ast.AsCall(node)
	if !ok {
		return node
	}
	// Do not optimize if the parameter was removed
	param, ok := call.OnlyArg()
	if !ok || !j2clStringValueOfRewriterPass_isStringValueOfMethodName(call.Callee().Node()) {
		return node
	}

	replacement := p.tryRewriteStringValueOfCall(node, param.Node())
	if replacement != node {
		replacement = replacement.SrcrefIfMissing(node)
		node.ReplaceWith(replacement)
		p.ReportChangeToEnclosingScope(replacement)
	}
	return replacement
}

// Cases to rewrite:
//
//	String.valueOf("foo") -> String("foo")
//	String.valueOf(null) -> String("null")
//	String.valueOf(undefined) -> "null"
//	String.valueOf("foo" + bar) -> String("foo" + bar)
//	String.valueOf(bar + "foo") -> String(bar + "foo")
//	String.valueOf("foo" + String.valueOf(bar)) -> String("foo" + String.valueOf(bar))
//
// Java: J2clStringValueOfRewriterPass.tryRewriteStringValueOfCall(Node)
func (p *J2clStringValueOfRewriterPass) tryRewriteStringValueOfCall(n rhino.Node, param rhino.Node) rhino.Node {
	if NodeUtil_IsUndefined(param) {
		return rhino.IR_String("null")
	} else if NodeUtil_IsDefinedValue(param) && !param.IsArrayLit() {
		// Generate String(param), let other peephole optimizations handle the rest when safe
		return NodeUtil_NewCallNode(rhino.IR_Name("String").Srcref(n), param.Detach())
	}
	return n
}

// Java: J2clStringValueOfRewriterPass.isStringValueOfMethodName(Node)
func j2clStringValueOfRewriterPass_isStringValueOfMethodName(fnName rhino.Node) bool {
	return fnName.IsQualifiedName() && *fnName.GetOriginalQualifiedName() == j2clStringValueOfRewriterPass_METHOD_NAME
}

// Java: final class J2clUndefinedChecksRewriterPass extends AbstractPeepholeOptimization
type J2clUndefinedChecksRewriterPass struct {
	AbstractPeepholeOptimizationBase
	shouldRunJ2clPasses bool
}

// Java: J2clUndefinedChecksRewriterPass()
func NewJ2clUndefinedChecksRewriterPass() *J2clUndefinedChecksRewriterPass {
	return &J2clUndefinedChecksRewriterPass{}
}

func (p *J2clUndefinedChecksRewriterPass) peepholeConfig() peepholeOptConfig {
	if p.shouldRunJ2clPasses {
		return peepholeOptConfig{
			tokens:         peepholeCallOnlyTokens,
			isStandard:     true,
			isJ2clCallPass: true,
		}
	}
	return peepholeOptConfig{isStandard: true}
}

func (p *J2clUndefinedChecksRewriterPass) cloneForWorker() AbstractPeepholeOptimization {
	return NewJ2clUndefinedChecksRewriterPass()
}

// Java: J2clUndefinedChecksRewriterPass.beginTraversal(AbstractCompiler)
func (p *J2clUndefinedChecksRewriterPass) BeginTraversal(compiler AbstractCompiler) {
	p.AbstractPeepholeOptimizationBase.BeginTraversal(compiler)
	p.shouldRunJ2clPasses = J2clSourceFileChecker_ShouldRunJ2clPasses(compiler)
}

// Java: J2clUndefinedChecksRewriterPass.optimizeSubtree(Node)
func (p *J2clUndefinedChecksRewriterPass) OptimizeSubtree(subtree rhino.Node) rhino.Node {
	if !p.shouldRunJ2clPasses {
		return subtree
	}
	call, ok := ast.AsCall(subtree)
	if !ok {
		return subtree
	}
	arg, ok := call.OnlyArg()
	if !ok {
		return subtree
	}
	callee := call.Callee().Node()
	if !callee.IsQualifiedName() {
		return subtree
	}
	qualifiedName := *callee.GetOriginalQualifiedName()
	if !strings.Contains(qualifiedName, "nativebootstrap$Util") {
		return subtree
	}

	var replacement rhino.Node
	if strings.HasSuffix(qualifiedName, "$isUndefined") {
		replacement = j2clUndefinedChecksRewriterPass_maybeOptimizeIsUndefinedCall(arg.Node())
	} else if strings.HasSuffix(qualifiedName, "$coerceToNull") {
		replacement = j2clUndefinedChecksRewriterPass_maybeOptimizeCoerceToNullCall(arg.Node())
	}

	if !replacement.IsNil() {
		jutil.CheckState(replacement != subtree)
		replacement.SrcrefTreeIfMissing(subtree)
		subtree.ReplaceWith(replacement)
		p.ReportChangeToEnclosingScope(replacement)
		return replacement
	}

	return subtree
}

// Java: J2clUndefinedChecksRewriterPass.maybeOptimizeIsUndefinedCall(Node)
func j2clUndefinedChecksRewriterPass_maybeOptimizeIsUndefinedCall(valueExpression rhino.Node) rhino.Node {
	rhino.CheckNotNull(valueExpression)
	if NodeUtil_IsUndefined(valueExpression) {
		return rhino.IR_TrueNode()
	} else if NodeUtil_IsDefinedValue(valueExpression) {
		return rhino.IR_Comma(valueExpression.Detach(), rhino.IR_FalseNode())
	}
	return rhino.NilNode
}

// Java: J2clUndefinedChecksRewriterPass.maybeOptimizeCoerceToNullCall(Node)
func j2clUndefinedChecksRewriterPass_maybeOptimizeCoerceToNullCall(valueExpression rhino.Node) rhino.Node {
	rhino.CheckNotNull(valueExpression)
	if NodeUtil_IsUndefined(valueExpression) {
		return rhino.IR_NullNode()
	}
	if NodeUtil_IsDefinedValue(valueExpression) {
		// detach the valueExpression as we'll be using it as the replacement node.
		return valueExpression.Detach()
	}
	return rhino.NilNode
}
