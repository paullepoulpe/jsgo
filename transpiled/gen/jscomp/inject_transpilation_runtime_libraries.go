/*
 * Copyright 2014 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/InjectTranspilationRuntimeLibraries.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.InjectTranspilationRuntimeLibraries
type InjectTranspilationRuntimeLibraries struct {
	compiler AbstractCompiler `gson:"compiler"`
	runtimeLibs *RuntimeJsLibManager `gson:"runtimeLibs"`
	idHash_ atomic.Int32
	injectedClassExtendsLibraries bool `gson:"injectedClassExtendsLibraries"`
	shouldInstrumentAsyncContext bool `gson:"shouldInstrumentAsyncContext"`
}

func (this_ *InjectTranspilationRuntimeLibraries) BindSelf_(s any) {
	_ = s
}
func (this_ *InjectTranspilationRuntimeLibraries) IsJavaNil_() bool { return this_ == nil }
func (*InjectTranspilationRuntimeLibraries) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.InjectTranspilationRuntimeLibraries", "InjectTranspilationRuntimeLibraries" }
func (this_ *InjectTranspilationRuntimeLibraries) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InjectTranspilationRuntimeLibraries) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InjectTranspilationRuntimeLibraries", "com.google.javascript.jscomp.InjectTranspilationRuntimeLibraries", func(v_ any) bool {
		return j2gort.InstanceOf[*InjectTranspilationRuntimeLibraries](v_)
	})
}
func (*InjectTranspilationRuntimeLibraries) Implements_jscomp_CompilerPass() {}

func (this_ *InjectTranspilationRuntimeLibraries) initFields_InjectTranspilationRuntimeLibraries() {
	_ = this_
}

// Java constructor: InjectTranspilationRuntimeLibraries(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewInjectTranspilationRuntimeLibraries_1(compiler AbstractCompiler) *InjectTranspilationRuntimeLibraries {
	this_ := &InjectTranspilationRuntimeLibraries{}
	InitInjectTranspilationRuntimeLibraries_1(this_, compiler)
	return this_
}

func InitInjectTranspilationRuntimeLibraries_1(this_ *InjectTranspilationRuntimeLibraries, compiler AbstractCompiler) {
	_ = compiler
	InitInjectTranspilationRuntimeLibraries_2(this_, compiler, compiler.GetOptions().GetInstrumentAsyncContext())
}

// Java constructor: InjectTranspilationRuntimeLibraries(com.google.javascript.jscomp.AbstractCompiler compiler, boolean shouldInstrumentAsyncContext)
func NewInjectTranspilationRuntimeLibraries_2(compiler AbstractCompiler, shouldInstrumentAsyncContext bool) *InjectTranspilationRuntimeLibraries {
	this_ := &InjectTranspilationRuntimeLibraries{}
	InitInjectTranspilationRuntimeLibraries_2(this_, compiler, shouldInstrumentAsyncContext)
	return this_
}

func InitInjectTranspilationRuntimeLibraries_2(this_ *InjectTranspilationRuntimeLibraries, compiler AbstractCompiler, shouldInstrumentAsyncContext bool) {
	_ = compiler
	_ = shouldInstrumentAsyncContext
	this_.initFields_InjectTranspilationRuntimeLibraries()
	this_.compiler = compiler
	this_.runtimeLibs = compiler.GetRuntimeJsLibManager()
	this_.injectedClassExtendsLibraries = false
	this_.shouldInstrumentAsyncContext = shouldInstrumentAsyncContext
}

// Java method: InjectTranspilationRuntimeLibraries.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *InjectTranspilationRuntimeLibraries) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	var used *pkg_parser.FeatureSet = pkg_parser.StaticGet_FeatureSet_ES3()
	_ = used
	{
		var script pkg_rhino.Node = root.GetFirstChild()
		_ = script
		for ; (script != 0); script = script.GetNext() {
			used = used.With_1_FeatureSet(InjectTranspilationRuntimeLibraries_GetScriptFeatures(script))
		}
	}
	var outputFeatures *pkg_parser.FeatureSet = this_.compiler.GetOptions().GetOutputFeatureSet()
	_ = outputFeatures
	if (!(outputFeatures.Contains_1_FeatureSet(used))) {
		NodeUtil_VisitPostOrder_3(root, func() NodeUtil_Visitor { bound_1 := this_; return NewLambda_NodeUtil_Visitor(func(mrp0_ pkg_rhino.Node) { bound_1.priv_InjectTranspilationRuntimeLibraries_checkForClassExtends_3207(mrp0_) }) }(), j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_2 uint32) bool {
	unused := pkg_rhino.Node(rawNode_2)
	_ = unused
	return !(this_.injectedClassExtendsLibraries)
}))
	}
	var mustBeCompiledAway *pkg_parser.FeatureSet = used.Without_1(outputFeatures)
	_ = mustBeCompiledAway
	if (mustBeCompiledAway.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_TEMPLATE_LITERALS())) {
		this_.runtimeLibs.InjectLibForField("$jscomp.createTemplateTagFirstArg")
		this_.runtimeLibs.InjectLibForField("$jscomp.createTemplateTagFirstArgWithRaw")
	}
	if (((mustBeCompiledAway.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_FOR_OF()) || mustBeCompiledAway.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_ARRAY_DESTRUCTURING())) || mustBeCompiledAway.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_OBJECT_PATTERN_REST()))) {
		this_.runtimeLibs.InjectLibForField("$jscomp.makeIterator")
	}
	if (mustBeCompiledAway.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_FOR_OF())) {
		this_.runtimeLibs.InjectLibForField("$jscomp.iteratorClose")
	}
	if (mustBeCompiledAway.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_ARRAY_PATTERN_REST())) {
		this_.runtimeLibs.InjectLibForField("$jscomp.arrayFromIterator")
	}
	if (mustBeCompiledAway.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_SPREAD_EXPRESSIONS())) {
		this_.runtimeLibs.InjectLibForField("$jscomp.arrayFromIterable")
	}
	if ((((mustBeCompiledAway.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD()) || mustBeCompiledAway.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_OBJECT_PATTERN_REST()))) && !(outputFeatures.Contains_1_FeatureSet(pkg_parser.StaticGet_FeatureSet_ES2015())))) {
		this_.runtimeLibs.InjectLibForField("$jscomp.assign")
	}
	if (mustBeCompiledAway.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_CLASS_GETTER_SETTER())) {
		this_.runtimeLibs.InjectLibForField("$jscomp.global")
	}
	if (mustBeCompiledAway.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_GENERATORS())) {
		this_.runtimeLibs.InjectLibForField("$jscomp.generator.createGenerator")
		this_.runtimeLibs.InjectLibForField("$jscomp.generator.Context")
	}
	if (mustBeCompiledAway.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_ASYNC_FUNCTIONS())) {
		this_.runtimeLibs.InjectLibForField("$jscomp.asyncExecutePromiseGeneratorFunction")
		if (!(outputFeatures.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_GENERATORS()))) {
			this_.runtimeLibs.InjectLibForField("$jscomp.asyncExecutePromiseGeneratorProgram")
			this_.runtimeLibs.InjectLibForField("$jscomp.generator.Context")
		}
	}
	if (mustBeCompiledAway.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_ASYNC_GENERATORS())) {
		this_.runtimeLibs.InjectLibForField("$jscomp.asyncExecutePromiseGeneratorFunction")
		this_.runtimeLibs.InjectLibForField("$jscomp.AsyncGeneratorWrapper")
		this_.runtimeLibs.InjectLibForField("$jscomp.AsyncGeneratorWrapper$ActionRecord")
		this_.runtimeLibs.InjectLibForField("$jscomp.AsyncGeneratorWrapper$ActionEnum.AWAIT_VALUE")
		this_.runtimeLibs.InjectLibForField("$jscomp.AsyncGeneratorWrapper$ActionEnum.YIELD_VALUE")
		this_.runtimeLibs.InjectLibForField("$jscomp.AsyncGeneratorWrapper$ActionEnum.YIELD_STAR")
		if (!(outputFeatures.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_GENERATORS()))) {
			this_.runtimeLibs.InjectLibForField("$jscomp.asyncExecutePromiseGeneratorProgram")
			this_.runtimeLibs.InjectLibForField("$jscomp.generator.Context")
		}
	}
	if (mustBeCompiledAway.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_FOR_AWAIT_OF())) {
		this_.runtimeLibs.InjectLibForField("$jscomp.makeAsyncIterator")
	}
	if (mustBeCompiledAway.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_REST_PARAMETERS())) {
		this_.runtimeLibs.InjectLibForField("$jscomp.getRestArguments")
	}
	if ((used.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_PRIVATE_ELEMENTS()) && mustBeCompiledAway.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_PRIVATE_ELEMENTS()))) {
		this_.runtimeLibs.InjectLibForField("$jscomp.PrivateMap")
	}
	if ((this_.shouldInstrumentAsyncContext && (((outputFeatures.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_ASYNC_FUNCTIONS()) || used.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_GENERATORS())) || used.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_ASYNC_GENERATORS()))))) {
		this_.runtimeLibs.InjectLibForField("$jscomp.asyncContextStart")
	}
}

// Java static method: InjectTranspilationRuntimeLibraries.getScriptFeatures(com.google.javascript.rhino.Node script)
func InjectTranspilationRuntimeLibraries_GetScriptFeatures(script pkg_rhino.Node) *pkg_parser.FeatureSet {
	_ = script
	var features *pkg_parser.FeatureSet = NodeUtil_GetFeatureSetOfScript(script)
	_ = features
	if (features != nil) {
		return features
	} else {
		return pkg_parser.StaticGet_FeatureSet_ES3()
	}
}

// Java method: InjectTranspilationRuntimeLibraries.checkForClassExtends(com.google.javascript.rhino.Node n)
func (this_ *InjectTranspilationRuntimeLibraries) priv_InjectTranspilationRuntimeLibraries_checkForClassExtends_3207(n pkg_rhino.Node) {
	_ = n
	if (!(n.IsClass())) {
		return
	}
	var superclass pkg_rhino.Node = n.GetSecondChild()
	_ = superclass
	if ((!(this_.injectedClassExtendsLibraries) && !(superclass.IsEmpty()))) {
		this_.runtimeLibs.InjectLibForField("$jscomp.construct")
		this_.runtimeLibs.InjectLibForField("$jscomp.inherits")
		this_.runtimeLibs.InjectLibForField("$jscomp.arrayFromIterable")
		this_.injectedClassExtendsLibraries = true
	}
}

