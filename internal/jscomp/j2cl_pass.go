// Port of com/google/javascript/jscomp/J2clPass.java @ closure-compiler v20261021

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

package jscomp

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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
	return "InlinableFunction[impl=" + jutil.ToJavaString(f.impl) + ", inliningMode=" + f.inliningMode.String() + "]"
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
func (c *J2clPass_ClassStaticFunctionsInliner) findMatchingSpecsForFile(sourceFileName *string, dst []*J2clPass_InlineSpec) []*J2clPass_InlineSpec {
	matchingSpecs := dst[:0]
	for _, inlineSpec := range c.inlineSpecs {
		if inlineSpec.classFileName == j2clPass_ALL_CLASS_FILE_NAMES ||
			strings.HasSuffix(*sourceFileName, inlineSpec.classFileName) {
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
		f.currFileInlineSpecs = f.outer.findMatchingSpecsForFile(n.GetSourceFileName(), f.currFileInlineSpecs)
		return len(f.currFileInlineSpecs) != 0
	}
	return true
}

// Java: J2clPass.ClassStaticFunctionsInliner.FunctionDefsCollector.visit(NodeTraversal, Node, Node)
func (f *J2clPass_ClassStaticFunctionsInliner_FunctionDefsCollector) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	// Only look at named function declarations that are fully qualified.
	var fnName string
	switch n.GetToken() {
	case rhino.Token_ASSIGN:
		// TODO(b/69730966): Delete this branch when ES5 syntax support is no longer needed.
		if !n.GetLastChild().IsFunction() {
			return
		}

		qualifiedNameNode := n.GetFirstChild()
		if !qualifiedNameNode.IsGetProp() || !qualifiedNameNode.IsQualifiedName() {
			return
		}
		fnName = qualifiedNameNode.GetString()
	case rhino.Token_MEMBER_FUNCTION_DEF:
		fnName = n.GetString()
	default:
		return
	}
	for _, inlineSpec := range f.currFileInlineSpecs {
		if !inlineSpec.fnNamesToInline[fnName] {
			continue
		}
		// Then store a reference to it.
		inlinableFunction := NewJ2clPass_ClassStaticFunctionsInliner_InlinableFunction(n.GetLastChild(), inlineSpec.inliningMode)
		var qualifiedFnName *string
		switch n.GetToken() {
		case rhino.Token_ASSIGN:
			qualifiedFnName = n.GetFirstChild().GetQualifiedName()
		case rhino.Token_MEMBER_FUNCTION_DEF:
			qualifiedFnName = NodeUtil_GetBestLValueName(n)
		default:
			panic(jutil.NewAssertionErrorNoMessage())
		}
		previousValue := f.outer.fnsToInlineByQualifiedName[*qualifiedFnName]
		f.outer.fnsToInlineByQualifiedName[*qualifiedFnName] = inlinableFunction
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
	if !n.IsCall() {
		return
	}

	// ... that are fully qualified
	qualifiedNameNode := n.GetFirstChild()
	if !qualifiedNameNode.IsGetProp() || !qualifiedNameNode.IsQualifiedName() {
		return
	}

	// ... and that reference a function definition we want to inline
	qualifiedFnName := *qualifiedNameNode.GetQualifiedName()
	fn := s.outer.fnsToInlineByQualifiedName[qualifiedFnName]
	if fn == nil {
		return
	}
	fnName := qualifiedNameNode.GetString()
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
