// Port of com/google/javascript/jscomp/J2clPass.java @ closure-compiler v20261104
// Port of com/google/javascript/jscomp/J2clAssertRemovalPass.java @ closure-compiler v20261104
// Port of com/google/javascript/jscomp/J2clStringValueOfRewriterPass.java @ closure-compiler v20261104
// Port of com/google/javascript/jscomp/J2clUndefinedChecksRewriterPass.java @ closure-compiler v20261104

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

func nodeHasJ2clSource(n rhino.Node) bool {
	sf := n.GetStaticSourceFile()
	return sf != nil && J2clSourceUtils_IsJ2clSource(sf.GetName())
}

// J2clSourceUtils_IsJ2clSource reports whether sourceFileName ends with ".java.js".
//
// Java: J2clSourceUtils.isJ2clSource(String)
func J2clSourceUtils_IsJ2clSource(sourceFileName string) bool {
	return strings.HasSuffix(sourceFileName, ".java.js")
}

// Java: J2clSourceUtils.isJ2clSource(Node)
func J2clSourceUtils_IsJ2clSourceOfNode(scriptNode rhino.Node) bool {
	jutil.CheckArgument(scriptNode.IsScript())
	return nodeHasJ2clSource(scriptNode)
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

// J2clSourceUtils_IsJ2clHeaderSource reports whether sourceFileName is a J2CL
// header file (*.java.js, excluding *.impl.java.js).
//
// Java: J2clSourceUtils.isJ2clHeaderSource(String)
func J2clSourceUtils_IsJ2clHeaderSource(sourceFileName string) bool {
	return strings.HasSuffix(sourceFileName, ".java.js") &&
		!strings.HasSuffix(sourceFileName, ".impl.java.js")
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
func hasJ2clSource(root rhino.Node) bool {
	for item := range ast.MustRoot(root).Items() {
		script := item.Node()
		jutil.CheckState(script.IsScript())
		if J2clSourceUtils_IsJ2clSourceOfNode(script) {
			return true
		}
	}
	return false
}

// Java: J2clSourceFileChecker.process(Node, Node)
func (c *J2clSourceFileChecker) Process(externs rhino.Node, root rhino.Node) {
	if hasJ2clSource(root) {
		c.compiler.SetRunJ2clPasses(true)
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

const j2clAllClassFileNames = "*"

// J2clPass is a normalization pass to inline some J2CL calls to enable other
// optimizations.
//
// Java: public class J2clPass implements CompilerPass
type J2clPass struct {
	compiler           AbstractCompiler
	safeNameIdSupplier func() string
	inliner            *j2clClassStaticFunctionsInliner
}

// Java: J2clPass(AbstractCompiler)
func NewJ2clPass(compiler AbstractCompiler) *J2clPass {
	return &J2clPass{
		compiler:           compiler,
		safeNameIdSupplier: compiler.GetUniqueNameIdSupplier(),
	}
}

var (
	j2clDefaultInlinings = []*j2clInlineSpec{
		/*
		 * Inline functions in Arrays that take references to static $isInstance() functions. This
		 * ensures that the references will be fully qualified to work with collapse properties.
		 */
		{
			classFileName: "Arrays.impl.java.js",
			fnNamesToInline: map[string]bool{
				"$create":                true,
				"$createWithInitializer": true,
				"$init":                  true,
				"$instanceIsOfType":      true,
				"$castTo":                true,
				"$stampType":             true,
			},
			inliningMode: FunctionInjector_InliningMode_DIRECT,
		},
		{
			classFileName:   "Casts.impl.java.js",
			fnNamesToInline: map[string]bool{"$to": true},
			inliningMode:    FunctionInjector_InliningMode_DIRECT,
		},
		/*
		 * Inlines all Interface.$markImplementor(FooClass) metaclass calls so that FooClass and others
		 * like it are not unnecessarily retained and so that static analysis of interface instanceof
		 * calls becomes possible.
		 *
		 * Note that this pass should NOT be restricted to j2cl .java.js files because JavaScript code
		 * implementing Java interfaces (not recommended but widely used in xplat) needs calls to
		 * $markImplementor.
		 */
		{
			classFileName:   j2clAllClassFileNames,
			fnNamesToInline: map[string]bool{"$markImplementor": true},
			inliningMode:    FunctionInjector_InliningMode_BLOCK,
		},
		/*
		 * Inlines class metadata calls so they become optimizable and avoids escaping of constructor.
		 */
		{
			classFileName: "Util.impl.java.js",
			fnNamesToInline: map[string]bool{
				"$setClassMetadata":             true,
				"$setClassMetadataForInterface": true,
				"$setClassMetadataForEnum":      true,
				"$setClassMetadataForPrimitive": true,
			},
			inliningMode: FunctionInjector_InliningMode_BLOCK,
		},
	}
	j2clDefaultFnNamesToInline = func() map[string]bool {
		s := make(map[string]bool, 12)
		for _, b := range j2clDefaultInlinings {
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
	if !p.compiler.RunJ2clPasses() {
		return
	}

	p.getOrCreateInliner(root, j2clDefaultInlinings).run()
}

// j2clInlineSpec is the private record InlineSpec.
//
// Java: private record J2clPass.InlineSpec(String, Set<String>, InliningMode)
type j2clInlineSpec struct {
	classFileName   string
	fnNamesToInline map[string]bool
	inliningMode    FunctionInjector_InliningMode
}

// j2clClassStaticFunctionsInliner collects references to certain function
// definitions in a certain class and then inlines fully qualified static method
// calls to those functions anywhere in the program.
//
// Assumes that the set of provided short function names will not collide with
// any of the collected fully qualified function names once the module prefix
// has been added.
//
// Java: private class J2clPass.ClassStaticFunctionsInliner
type j2clClassStaticFunctionsInliner struct {
	outer                      *J2clPass
	fnsToInlineByQualifiedName map[string]*j2clInlinableFunction
	injector                   *FunctionInjector
	root                       rhino.Node
	inlineSpecs                []*j2clInlineSpec
	collector                  j2clFunctionDefsCollector
	callInliner                j2clStaticCallInliner
}

// j2clInlinableFunction is the private record InlinableFunction.
//
// Java: private record J2clPass.ClassStaticFunctionsInliner.InlinableFunction(Node, InliningMode)
type j2clInlinableFunction struct {
	impl           rhino.Node
	inliningMode   FunctionInjector_InliningMode
	summary        calleeSummary
	summaryScanned bool
}

func (f *j2clInlinableFunction) calleeSummary(fnName string) *calleeSummary {
	if !f.summaryScanned {
		scanCalleeSummary(fnName, f.impl, &f.summary)
		f.summaryScanned = true
	}
	return &f.summary
}

// String is the record's toString().
//
// Java: J2clPass.ClassStaticFunctionsInliner.InlinableFunction.toString()
func (f *j2clInlinableFunction) String() string {
	implStr := "null"
	if !f.impl.IsNil() {
		implStr = f.impl.String()
	}
	return "InlinableFunction[impl=" + implStr + ", inliningMode=" + f.inliningMode.String() + "]"
}

func (p *J2clPass) getOrCreateInliner(root rhino.Node, inlineSpecs []*j2clInlineSpec) *j2clClassStaticFunctionsInliner {
	c := p.inliner
	if c == nil {
		c = &j2clClassStaticFunctionsInliner{
			outer:                      p,
			fnsToInlineByQualifiedName: make(map[string]*j2clInlinableFunction),
		}
		c.injector = NewFunctionInjector(p.compiler, FunctionInjectorOptions{
			SafeNameIdSupplier: p.safeNameIdSupplier,
			AssumeStrictThis:   true,
			AllowDecomposition: true,
		})
		c.injector.SetKnownConstantFunctions(j2clDefaultFnNamesToInline)
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
func (c *j2clClassStaticFunctionsInliner) run() {
	NodeTraversal_Traverse(c.outer.compiler, c.root, &c.collector)
	if len(c.fnsToInlineByQualifiedName) > 0 {
		NodeTraversal_Traverse(c.outer.compiler, c.root, &c.callInliner)
	}
}

// Java: J2clPass.ClassStaticFunctionsInliner.findMatchingSpecsForFile(String)
func (c *j2clClassStaticFunctionsInliner) findMatchingSpecsForFile(
	sourceFileName string, hasSource bool, dst []*j2clInlineSpec,
) []*j2clInlineSpec {
	matchingSpecs := dst[:0]
	for _, inlineSpec := range c.inlineSpecs {
		if inlineSpec.classFileName == j2clAllClassFileNames {
			matchingSpecs = append(matchingSpecs, inlineSpec)
			continue
		}
		jutil.CheckStateMsg(hasSource, "expected script node to have a source file name")
		if strings.HasSuffix(sourceFileName, inlineSpec.classFileName) {
			matchingSpecs = append(matchingSpecs, inlineSpec)
		}
	}
	return matchingSpecs
}

var (
	j2clFunctionDefsCollectorVisitTokens = makeVisitTokenSet(rhino.Token_ASSIGN, rhino.Token_MEMBER_FUNCTION_DEF)
	j2clCallOnlyVisitTokens              = makeVisitTokenSet(rhino.Token_CALL)
)

// Java: private class J2clPass.ClassStaticFunctionsInliner.FunctionDefsCollector implements NodeTraversal.Callback
type j2clFunctionDefsCollector struct {
	outer *j2clClassStaticFunctionsInliner
	// This list is initialized every time the node traversal enters a new script. It contains
	// all InlineSpecs whose classFileName matches the script's filename.
	currFileInlineSpecs []*j2clInlineSpec
}

var (
	_ nodeTraversal_visitTokens                         = (*j2clFunctionDefsCollector)(nil)
	_ nodeTraversal_noVisitLeaves                       = (*j2clFunctionDefsCollector)(nil)
	_ nodeTraversal_onlyScriptAndFunctionShouldTraverse = (*j2clFunctionDefsCollector)(nil)
)

func (*j2clFunctionDefsCollector) noVisitLeaves() {}

func (*j2clFunctionDefsCollector) visitTokens() *ast.TokenSet {
	return &j2clFunctionDefsCollectorVisitTokens
}

func (*j2clFunctionDefsCollector) shouldTraverseOnlyScriptAndFunction() bool {
	return true
}

// Java: J2clPass.ClassStaticFunctionsInliner.FunctionDefsCollector.shouldTraverse(NodeTraversal, Node, Node)
func (f *j2clFunctionDefsCollector) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsScript() {
		srcName, hasSrc := n.SourceFileName()
		f.currFileInlineSpecs = f.outer.findMatchingSpecsForFile(srcName, hasSrc, f.currFileInlineSpecs)
		return len(f.currFileInlineSpecs) != 0
	}
	return true
}

// Java: J2clPass.ClassStaticFunctionsInliner.FunctionDefsCollector.visit(NodeTraversal, Node, Node)
func (f *j2clFunctionDefsCollector) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	// Only look at named function declarations that are fully qualified.
	var fnName string
	var fnImpl rhino.Node
	var assignTarget rhino.Node
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
		assignTarget = qualifiedProp.Node()
	} else if memberFn, ok := ast.AsMemberFunctionDef(n); ok {
		fnName = memberFn.Key()
		fnImpl = memberFn.Function().Node()
	} else {
		return
	}
	for _, inlineSpec := range f.currFileInlineSpecs {
		if !inlineSpec.fnNamesToInline[fnName] {
			continue
		}
		// Then store a reference to it.
		inlinableFunction := &j2clInlinableFunction{
			impl:         fnImpl,
			inliningMode: inlineSpec.inliningMode,
		}
		var qualifiedFnName string
		if !assignTarget.IsNil() {
			qname, ok := assignTarget.GetQualifiedNameString()
			jutil.CheckStateMsg(ok, "expected qualified name on assignment target")
			qualifiedFnName = qname
		} else {
			qname, ok := NodeUtil_GetBestLValueName(n)
			if !ok {
				panic(jutil.NewNullPointerException())
			}
			qualifiedFnName = qname
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
type j2clStaticCallInliner struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *j2clClassStaticFunctionsInliner
}

var (
	_ nodeTraversal_visitTokens   = (*j2clStaticCallInliner)(nil)
	_ nodeTraversal_noVisitLeaves = (*j2clStaticCallInliner)(nil)
)

func (*j2clStaticCallInliner) noVisitLeaves() {}

func (*j2clStaticCallInliner) visitTokens() *ast.TokenSet {
	return &j2clCallOnlyVisitTokens
}

// Java: J2clPass.ClassStaticFunctionsInliner.StaticCallInliner.visit(NodeTraversal, Node, Node)
func (s *j2clStaticCallInliner) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
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
	fnImpl := fn.impl
	inliningMode := fn.inliningMode

	// Ensure that the function only has a single return statement when direct inlining.
	if inliningMode == FunctionInjector_InliningMode_DIRECT {
		firstStmt, hasStmt := ast.MustFunction(fnImpl).Body().MustBlock().FirstStatement()
		if !hasStmt || !firstStmt.Node().IsReturn() {
			panic(jutil.NewIllegalStateExceptionMsg(
				"Attempted to direct inline function " +
					qualifiedFnName +
					", but function is not a simple return."))
		}
	}

	// Otherwise inline the call.
	// Note: This pass has to run before normalization, so we must use the unsafeInline method.
	// It is safe because these are strictly controlled trivial bootstrap methods that are
	// written with inlining in mind (e.g. doesn't read/write local/global variables).
	// TODO(goktug): Add a check that will ensure safety of this.
	var summary *calleeSummary
	if inliningMode == FunctionInjector_InliningMode_BLOCK {
		summary = fn.calleeSummary(fnName)
	}
	inlinedCall := s.outer.injector.internalInlineWithSummary(
		NewFunctionInjector_Reference(n, t.GetScope(), t.GetChunk(), inliningMode), fnName, fnImpl, summary)
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

var (
	_ nodeTraversal_visitTokens   = (*J2clAssertRemovalPass)(nil)
	_ nodeTraversal_noVisitLeaves = (*J2clAssertRemovalPass)(nil)
)

func (*J2clAssertRemovalPass) noVisitLeaves() {}

func (*J2clAssertRemovalPass) visitTokens() *ast.TokenSet {
	return &j2clCallOnlyVisitTokens
}

// Java: J2clAssertRemovalPass(AbstractCompiler)
func NewJ2clAssertRemovalPass(compiler AbstractCompiler) *J2clAssertRemovalPass {
	return &J2clAssertRemovalPass{compiler: compiler}
}

// Java: J2clAssertRemovalPass.process(Node, Node)
func (p *J2clAssertRemovalPass) Process(externs rhino.Node, root rhino.Node) {
	if !p.compiler.RunJ2clPasses() {
		return
	}

	NodeTraversal_Traverse(p.compiler, root, p)
}

// Java: J2clAssertRemovalPass.visit(NodeTraversal, Node, Node)
func (p *J2clAssertRemovalPass) Visit(t *NodeTraversal, node rhino.Node, parent rhino.Node) {
	if isJ2clAssertionCall(node) {
		NodeUtil_DeleteFunctionCall(node, p.compiler)
	}
}

// Java: J2clAssertRemovalPass.isAssertionCall(Node)
func isJ2clAssertionCall(node rhino.Node) bool {
	call, ok := ast.AsCall(node)
	return ok && isJ2clAssertionMethodName(call.Callee().Node())
}

// Java: J2clAssertRemovalPass.isAssertionMethodName(Node)
func isJ2clAssertionMethodName(fnName rhino.Node) bool {
	originalQname, ok := fnName.OriginalQualifiedName()
	return ok && strings.Contains(originalQname, "Asserts") && strings.Contains(originalQname, ".$assert")
}

const j2clStringValueOfMethodName = "module$exports$java$lang$String$impl.m_valueOf__java_lang_Object__java_lang_String"

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
	p.shouldRunJ2clPasses = compiler.RunJ2clPasses()
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
	if !ok || !isJ2clStringValueOfMethodName(call.Callee().Node()) {
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
func isJ2clStringValueOfMethodName(fnName rhino.Node) bool {
	originalQname, ok := fnName.OriginalQualifiedName()
	return ok && originalQname == j2clStringValueOfMethodName
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
	p.shouldRunJ2clPasses = compiler.RunJ2clPasses()
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
	qualifiedName, ok := callee.OriginalQualifiedName()
	if !ok {
		return subtree
	}
	if !strings.Contains(qualifiedName, "nativebootstrap$Util") {
		return subtree
	}

	var replacement rhino.Node
	if strings.HasSuffix(qualifiedName, "$isUndefined") {
		replacement = maybeOptimizeJ2clIsUndefinedCall(arg.Node())
	} else if strings.HasSuffix(qualifiedName, "$coerceToNull") {
		replacement = maybeOptimizeJ2clCoerceToNullCall(arg.Node())
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
func maybeOptimizeJ2clIsUndefinedCall(valueExpression rhino.Node) rhino.Node {
	rhino.CheckNotNull(valueExpression)
	if NodeUtil_IsUndefined(valueExpression) {
		return rhino.IR_TrueNode()
	} else if NodeUtil_IsDefinedValue(valueExpression) {
		return rhino.IR_Comma(valueExpression.Detach(), rhino.IR_FalseNode())
	}
	return rhino.NilNode
}

// Java: J2clUndefinedChecksRewriterPass.maybeOptimizeCoerceToNullCall(Node)
func maybeOptimizeJ2clCoerceToNullCall(valueExpression rhino.Node) rhino.Node {
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
