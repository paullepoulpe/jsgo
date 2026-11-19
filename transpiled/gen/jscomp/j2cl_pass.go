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
// Port of com/google/javascript/jscomp/J2clPass.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_j2cl_pass_1 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_j2cl_pass_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_2 any) any {
		b := j2gort.CastPtr[J2clPass_InlineSpec](lp_2)
		_ = b
		return b.fnNamesToInline.Stream()
	})
}

// Java class: com.google.javascript.jscomp.J2clPass
type J2clPass struct {
	compiler           AbstractCompiler                       `gson:"compiler"`
	safeNameIdSupplier j2gort.Ext_java_util_function_Supplier `gson:"safeNameIdSupplier"`
	idHash_            atomic.Int32
}

func (this_ *J2clPass) BindSelf_(s any) {
	_ = s
}
func (this_ *J2clPass) IsJavaNil_() bool { return this_ == nil }
func (*J2clPass) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.J2clPass", "J2clPass"
}
func (this_ *J2clPass) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*J2clPass) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clPass", "com.google.javascript.jscomp.J2clPass", func(v_ any) bool {
		return j2gort.InstanceOf[*J2clPass](v_)
	})
}
func (*J2clPass) Implements_jscomp_CompilerPass() {}

func StaticGet_J2clPass_ALL_CLASS_FILE_NAMES() string { return "*" }
func (this_ *J2clPass) initFields_J2clPass() {
	_ = this_
}

// Java constructor: J2clPass(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewJ2clPass(compiler AbstractCompiler) *J2clPass {
	this_ := &J2clPass{}
	InitJ2clPass(this_, compiler)
	return this_
}

func InitJ2clPass(this_ *J2clPass, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_J2clPass()
	this_.compiler = compiler
	this_.safeNameIdSupplier = j2gort.AsIface[j2gort.Ext_java_util_function_Supplier](compiler.GetUniqueNameIdSupplier())
}

// Java method: J2clPass.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *J2clPass) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if !(J2clSourceFileChecker_ShouldRunJ2clPasses(this_.compiler)) {
		return
	}
	var inlinings j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = inlinings
	inlinings.Add_1(NewJ2clPass_InlineSpec("Arrays.impl.java.js", j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_7("$create", "$createWithInitializer", "$init", "$instanceIsOfType", "$castTo", "$stampType", []any{})), StaticGet_FunctionInjector_InliningMode_DIRECT()))
	inlinings.Add_1(NewJ2clPass_InlineSpec("Casts.impl.java.js", j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.ImmSetOfStr1("$to")), StaticGet_FunctionInjector_InliningMode_DIRECT()))
	inlinings.Add_1(NewJ2clPass_InlineSpec("*", j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.ImmSetOfStr1("$markImplementor")), StaticGet_FunctionInjector_InliningMode_BLOCK()))
	inlinings.Add_1(NewJ2clPass_InlineSpec("Util.impl.java.js", j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_4("$setClassMetadata", "$setClassMetadataForInterface", "$setClassMetadataForEnum", "$setClassMetadataForPrimitive")), StaticGet_FunctionInjector_InliningMode_BLOCK()))
	{
		var val_1 J2clPass_ClassStaticFunctionsInliner
		val_1.outer_ = this_
		InitJ2clPass_ClassStaticFunctionsInliner(&val_1, root, inlinings)
		(&val_1).priv_ClassStaticFunctionsInliner_run_3487()
	}
}

// Java record: com.google.javascript.jscomp.J2clPass.InlineSpec
type J2clPass_InlineSpec struct {
	classFileName   string                        `gson:"classFileName"`
	fnNamesToInline j2gort.Ext_java_util_Set      `gson:"fnNamesToInline"`
	inliningMode    FunctionInjector_InliningMode `gson:"inliningMode"`
}

func (this_ *J2clPass_InlineSpec) BindSelf_(s any) {
	_ = s
}
func (this_ *J2clPass_InlineSpec) IsJavaNil_() bool { return this_ == nil }
func (*J2clPass_InlineSpec) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.J2clPass$InlineSpec", "InlineSpec"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clPass.InlineSpec", "com.google.javascript.jscomp.J2clPass$InlineSpec", func(v_ any) bool {
		return j2gort.InstanceOf[*J2clPass_InlineSpec](v_)
	})
}
func (*J2clPass_InlineSpec) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *J2clPass_InlineSpec) initFields_J2clPass_InlineSpec() {
	_ = this_
}

// Java constructor: InlineSpec(java.lang.String classFileName, java.util.Set<java.lang.String> fnNamesToInline, com.google.javascript.jscomp.FunctionInjector.InliningMode inliningMode)
func NewJ2clPass_InlineSpec(classFileName string, fnNamesToInline j2gort.Ext_java_util_Set, inliningMode FunctionInjector_InliningMode) *J2clPass_InlineSpec {
	this_ := &J2clPass_InlineSpec{}
	InitJ2clPass_InlineSpec(this_, classFileName, fnNamesToInline, inliningMode)
	return this_
}

func InitJ2clPass_InlineSpec(this_ *J2clPass_InlineSpec, classFileName string, fnNamesToInline j2gort.Ext_java_util_Set, inliningMode FunctionInjector_InliningMode) {
	_ = classFileName
	_ = fnNamesToInline
	_ = inliningMode
	this_.initFields_J2clPass_InlineSpec()
	this_.classFileName = classFileName
	this_.fnNamesToInline = fnNamesToInline
	this_.inliningMode = inliningMode
}

func (this_ *J2clPass_InlineSpec) ClassFileName() string { return this_.classFileName }
func (this_ *J2clPass_InlineSpec) FnNamesToInline() j2gort.Ext_java_util_Set {
	return this_.fnNamesToInline
}
func (this_ *J2clPass_InlineSpec) InliningMode() FunctionInjector_InliningMode {
	return this_.inliningMode
}
func (this_ *J2clPass_InlineSpec) Equals(other_ any) bool {
	o_, ok_ := other_.(*J2clPass_InlineSpec)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.classFileName != o_.classFileName {
		return false
	}
	if !javalang.ObjectsEquals(this_.fnNamesToInline, o_.fnNamesToInline) {
		return false
	}
	if this_.inliningMode != o_.inliningMode {
		return false
	}
	return true
}
func (this_ *J2clPass_InlineSpec) HashCode() int32 {
	var h_ int32 = 0
	{
		hs_2 := this_.classFileName
		h_ = 31 * h_
		if hs_2 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_2)
		}
	}
	h_ = 31*h_ + j2gort.ObjectHash(this_.fnNamesToInline)
	h_ = 31*h_ + (this_.inliningMode).HashCode()
	return h_
}
func (this_ *J2clPass_InlineSpec) String() string {
	_ = fmt.Sprint
	return "InlineSpec[" + "classFileName=" + j2gort.StrVal(j2gort.BoxString(this_.classFileName)) + ", " + "fnNamesToInline=" + j2gort.StrVal(this_.fnNamesToInline) + ", " + "inliningMode=" + j2gort.StrVal(j2gort.BoxNullableVal(this_.inliningMode)) + "]"
}

// Java class: com.google.javascript.jscomp.J2clPass.ClassStaticFunctionsInliner
type J2clPass_ClassStaticFunctionsInliner struct {
	outer_                     *J2clPass
	fnsToInlineByQualifiedName j2gort.OrderedMap[string, *J2clPass_ClassStaticFunctionsInliner_InlinableFunction] `gson:"fnsToInlineByQualifiedName"`
	injector                   *FunctionInjector                                                                  `gson:"injector"`
	inlineSpecs                j2gort.Ext_java_util_List                                                          `gson:"inlineSpecs"`
	root                       pkg_rhino.Node                                                                     `gson:"root"`
	idHash_                    atomic.Int32
}

func (this_ *J2clPass_ClassStaticFunctionsInliner) BindSelf_(s any) {
	_ = s
}
func (this_ *J2clPass_ClassStaticFunctionsInliner) IsJavaNil_() bool { return this_ == nil }
func (*J2clPass_ClassStaticFunctionsInliner) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.J2clPass$ClassStaticFunctionsInliner", "ClassStaticFunctionsInliner"
}
func (this_ *J2clPass_ClassStaticFunctionsInliner) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*J2clPass_ClassStaticFunctionsInliner) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clPass.ClassStaticFunctionsInliner", "com.google.javascript.jscomp.J2clPass$ClassStaticFunctionsInliner", func(v_ any) bool {
		return j2gort.InstanceOf[*J2clPass_ClassStaticFunctionsInliner](v_)
	})
}
func (this_ *J2clPass_ClassStaticFunctionsInliner) SetOuter_(o *J2clPass) { this_.outer_ = o }
func (this_ *J2clPass_ClassStaticFunctionsInliner) GetOuter_() *J2clPass  { return this_.outer_ }

func (this_ *J2clPass_ClassStaticFunctionsInliner) initFields_J2clPass_ClassStaticFunctionsInliner() {
	_ = this_
	this_.fnsToInlineByQualifiedName = j2gort.NewOrderedMap[string, *J2clPass_ClassStaticFunctionsInliner_InlinableFunction]()
}

// Java constructor: ClassStaticFunctionsInliner(com.google.javascript.rhino.Node root, java.util.List<com.google.javascript.jscomp.J2clPass.InlineSpec> inlineSpecs)
func NewJ2clPass_ClassStaticFunctionsInliner(outer_ *J2clPass, root pkg_rhino.Node, inlineSpecs j2gort.Ext_java_util_List) *J2clPass_ClassStaticFunctionsInliner {
	this_ := &J2clPass_ClassStaticFunctionsInliner{}
	this_.outer_ = outer_
	InitJ2clPass_ClassStaticFunctionsInliner(this_, root, inlineSpecs)
	return this_
}

func InitJ2clPass_ClassStaticFunctionsInliner(this_ *J2clPass_ClassStaticFunctionsInliner, root pkg_rhino.Node, inlineSpecs j2gort.Ext_java_util_List) {
	_ = root
	_ = inlineSpecs
	this_.initFields_J2clPass_ClassStaticFunctionsInliner()
	this_.root = root
	this_.inlineSpecs = inlineSpecs
	this_.injector = NewFunctionInjector_Builder(this_.outer_.compiler).SafeNameIdSupplier(this_.outer_.safeNameIdSupplier).AssumeStrictThis(true).Build()
	var fnNamesToInline j2gort.Ext_com_google_common_collect_ImmutableSet = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](inlineSpecs.Stream().FlatMap(lambdaSingleton_j2cl_pass_1).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableSet_ToImmutableSet()))
	_ = fnNamesToInline
	this_.injector.SetKnownConstantFunctions(fnNamesToInline)
}

// Java method: ClassStaticFunctionsInliner.run()
func (this_ *J2clPass_ClassStaticFunctionsInliner) priv_ClassStaticFunctionsInliner_run_3487() {
	NodeTraversal_Traverse(this_.outer_.compiler, this_.root, NewJ2clPass_ClassStaticFunctionsInliner_FunctionDefsCollector(this_))
	NodeTraversal_Traverse(this_.outer_.compiler, this_.root, NewJ2clPass_ClassStaticFunctionsInliner_StaticCallInliner(this_))
}

// Java method: ClassStaticFunctionsInliner.findMatchingSpecsForFile(java.lang.String sourceFileName)
func (this_ *J2clPass_ClassStaticFunctionsInliner) priv_ClassStaticFunctionsInliner_findMatchingSpecsForFile_3488(sourceFileName string) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = sourceFileName
	var matchingSpecs []*J2clPass_InlineSpec = []*J2clPass_InlineSpec{}
	_ = matchingSpecs
	{
		cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.inlineSpecs))
		for cur_3.HasNext() {
			rawElem_4 := cur_3.Next()
			inlineSpec := j2gort.CastPtr[J2clPass_InlineSpec](rawElem_4)
			_ = inlineSpec
			{
				var iifeCond_1 bool
				{

					sr_5 := inlineSpec.classFileName
					se_6 := "*"
					iifeCond_1 = j2gort.CheckNotNullString(sr_5) == se_6
				}
				if iifeCond_1 || javalang.StringEndsWith(j2gort.CheckNotNullString(sourceFileName), j2gort.CheckNotNullString(inlineSpec.classFileName)) {
					matchingSpecs = append(matchingSpecs, j2gort.CheckNotNullElem(inlineSpec))
				}
			}
		}
	}
	return j2gort.ImmListFromTypedSlice(matchingSpecs)
}

// Java record: com.google.javascript.jscomp.J2clPass.ClassStaticFunctionsInliner.InlinableFunction
type J2clPass_ClassStaticFunctionsInliner_InlinableFunction struct {
	impl         pkg_rhino.Node                `gson:"impl"`
	inliningMode FunctionInjector_InliningMode `gson:"inliningMode"`
}

func (this_ *J2clPass_ClassStaticFunctionsInliner_InlinableFunction) BindSelf_(s any) {
	_ = s
}
func (this_ *J2clPass_ClassStaticFunctionsInliner_InlinableFunction) IsJavaNil_() bool {
	return this_ == nil
}
func (*J2clPass_ClassStaticFunctionsInliner_InlinableFunction) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.J2clPass$ClassStaticFunctionsInliner$InlinableFunction", "InlinableFunction"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clPass.ClassStaticFunctionsInliner.InlinableFunction", "com.google.javascript.jscomp.J2clPass$ClassStaticFunctionsInliner$InlinableFunction", func(v_ any) bool {
		return j2gort.InstanceOf[*J2clPass_ClassStaticFunctionsInliner_InlinableFunction](v_)
	})
}
func (*J2clPass_ClassStaticFunctionsInliner_InlinableFunction) Implements_j2gort_Ext_java_lang_Record() {
}

func (this_ *J2clPass_ClassStaticFunctionsInliner_InlinableFunction) initFields_J2clPass_ClassStaticFunctionsInliner_InlinableFunction() {
	_ = this_
}

// Java constructor: InlinableFunction(com.google.javascript.rhino.Node impl, com.google.javascript.jscomp.FunctionInjector.InliningMode inliningMode)
func NewJ2clPass_ClassStaticFunctionsInliner_InlinableFunction(impl pkg_rhino.Node, inliningMode FunctionInjector_InliningMode) *J2clPass_ClassStaticFunctionsInliner_InlinableFunction {
	this_ := &J2clPass_ClassStaticFunctionsInliner_InlinableFunction{}
	InitJ2clPass_ClassStaticFunctionsInliner_InlinableFunction(this_, impl, inliningMode)
	return this_
}

func InitJ2clPass_ClassStaticFunctionsInliner_InlinableFunction(this_ *J2clPass_ClassStaticFunctionsInliner_InlinableFunction, impl pkg_rhino.Node, inliningMode FunctionInjector_InliningMode) {
	_ = impl
	_ = inliningMode
	this_.initFields_J2clPass_ClassStaticFunctionsInliner_InlinableFunction()
	this_.impl = impl
	this_.inliningMode = inliningMode
}

func (this_ *J2clPass_ClassStaticFunctionsInliner_InlinableFunction) Impl() pkg_rhino.Node {
	return this_.impl
}
func (this_ *J2clPass_ClassStaticFunctionsInliner_InlinableFunction) InliningMode() FunctionInjector_InliningMode {
	return this_.inliningMode
}
func (this_ *J2clPass_ClassStaticFunctionsInliner_InlinableFunction) Equals(other_ any) bool {
	o_, ok_ := other_.(*J2clPass_ClassStaticFunctionsInliner_InlinableFunction)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.impl != o_.impl {
		return false
	}
	if this_.inliningMode != o_.inliningMode {
		return false
	}
	return true
}
func (this_ *J2clPass_ClassStaticFunctionsInliner_InlinableFunction) HashCode() int32 {
	var h_ int32 = 0
	h_ = 31*h_ + int32(this_.impl)
	h_ = 31*h_ + (this_.inliningMode).HashCode()
	return h_
}
func (this_ *J2clPass_ClassStaticFunctionsInliner_InlinableFunction) String() string {
	_ = fmt.Sprint
	return "InlinableFunction[" + "impl=" + j2gort.StrVal(j2gort.BoxNode(this_.impl)) + ", " + "inliningMode=" + j2gort.StrVal(j2gort.BoxNullableVal(this_.inliningMode)) + "]"
}

// Java class: com.google.javascript.jscomp.J2clPass.ClassStaticFunctionsInliner.FunctionDefsCollector
type J2clPass_ClassStaticFunctionsInliner_FunctionDefsCollector struct {
	outer_              *J2clPass_ClassStaticFunctionsInliner
	currFileInlineSpecs []*J2clPass_InlineSpec `gson:"currFileInlineSpecs"`
	idHash_             atomic.Int32
}

func (this_ *J2clPass_ClassStaticFunctionsInliner_FunctionDefsCollector) BindSelf_(s any) {
	_ = s
}
func (this_ *J2clPass_ClassStaticFunctionsInliner_FunctionDefsCollector) IsJavaNil_() bool {
	return this_ == nil
}
func (*J2clPass_ClassStaticFunctionsInliner_FunctionDefsCollector) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.J2clPass$ClassStaticFunctionsInliner$FunctionDefsCollector", "FunctionDefsCollector"
}
func (this_ *J2clPass_ClassStaticFunctionsInliner_FunctionDefsCollector) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*J2clPass_ClassStaticFunctionsInliner_FunctionDefsCollector) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clPass.ClassStaticFunctionsInliner.FunctionDefsCollector", "com.google.javascript.jscomp.J2clPass$ClassStaticFunctionsInliner$FunctionDefsCollector", func(v_ any) bool {
		return j2gort.InstanceOf[*J2clPass_ClassStaticFunctionsInliner_FunctionDefsCollector](v_)
	})
}
func (this_ *J2clPass_ClassStaticFunctionsInliner_FunctionDefsCollector) SetOuter_(o *J2clPass_ClassStaticFunctionsInliner) {
	this_.outer_ = o
}
func (this_ *J2clPass_ClassStaticFunctionsInliner_FunctionDefsCollector) GetOuter_() *J2clPass_ClassStaticFunctionsInliner {
	return this_.outer_
}
func (*J2clPass_ClassStaticFunctionsInliner_FunctionDefsCollector) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *J2clPass_ClassStaticFunctionsInliner_FunctionDefsCollector) initFields_J2clPass_ClassStaticFunctionsInliner_FunctionDefsCollector() {
	_ = this_
}

// Java constructor: FunctionDefsCollector()
func NewJ2clPass_ClassStaticFunctionsInliner_FunctionDefsCollector(outer_ *J2clPass_ClassStaticFunctionsInliner) *J2clPass_ClassStaticFunctionsInliner_FunctionDefsCollector {
	this_ := &J2clPass_ClassStaticFunctionsInliner_FunctionDefsCollector{}
	this_.outer_ = outer_
	InitJ2clPass_ClassStaticFunctionsInliner_FunctionDefsCollector(this_)
	return this_
}

func InitJ2clPass_ClassStaticFunctionsInliner_FunctionDefsCollector(this_ *J2clPass_ClassStaticFunctionsInliner_FunctionDefsCollector) {
	this_.initFields_J2clPass_ClassStaticFunctionsInliner_FunctionDefsCollector()
}

// Java method: FunctionDefsCollector.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *J2clPass_ClassStaticFunctionsInliner_FunctionDefsCollector) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if n.IsScript() {
		this_.currFileInlineSpecs = j2gort.ImmListToSliceNullable[*J2clPass_InlineSpec](this_.outer_.priv_ClassStaticFunctionsInliner_findMatchingSpecsForFile_3488(n.GetSourceFileName()))
		return len(this_.currFileInlineSpecs) != 0
	}
	return true
}

// Java method: FunctionDefsCollector.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *J2clPass_ClassStaticFunctionsInliner_FunctionDefsCollector) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	var fnName string = j2gort.NullString
	_ = fnName
	switch n.TokenSwitchOrd_() {
	case int32(51) /* Token.ASSIGN */ :
		if !(n.GetLastChild().IsFunction()) {
			return
		}
		var qualifiedNameNode pkg_rhino.Node = n.GetFirstChild()
		_ = qualifiedNameNode
		if !(qualifiedNameNode.IsGetProp()) || !(qualifiedNameNode.IsQualifiedName()) {
			return
		}
		fnName = qualifiedNameNode.GetString()
	case int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
		fnName = n.GetString()
	default:
		return
	}
	{

		for _, inlineSpec := range this_.currFileInlineSpecs {

			_ = inlineSpec
			if !(j2gort.SetContainsStr(inlineSpec.fnNamesToInline, fnName)) {
				continue
			}
			var inlinableFunction *J2clPass_ClassStaticFunctionsInliner_InlinableFunction = NewJ2clPass_ClassStaticFunctionsInliner_InlinableFunction(n.GetLastChild(), inlineSpec.inliningMode)
			_ = inlinableFunction
			var qualifiedFnName string = j2gort.NullString
			_ = qualifiedFnName
			switch n.TokenSwitchOrd_() {
			case int32(51) /* Token.ASSIGN */ :
				qualifiedFnName = n.GetFirstChild().GetQualifiedName()
			case int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
				qualifiedFnName = NodeUtil_GetBestLValueName(n)
			default:
				panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
			}
			var previousValue *J2clPass_ClassStaticFunctionsInliner_InlinableFunction = this_.outer_.fnsToInlineByQualifiedName.Put(qualifiedFnName, inlinableFunction)
			_ = previousValue
			if !(previousValue == nil) {
				j2gort.CheckStateFmt(false, "expected each function to match at most one InliningSpec, but found:\n%s\n%s", j2gort.PtrToAny(previousValue), inlinableFunction)
			}
		}
	}
}

// Java class: com.google.javascript.jscomp.J2clPass.ClassStaticFunctionsInliner.StaticCallInliner
type J2clPass_ClassStaticFunctionsInliner_StaticCallInliner struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *J2clPass_ClassStaticFunctionsInliner
}

func (this_ *J2clPass_ClassStaticFunctionsInliner_StaticCallInliner) BindSelf_(s any) {
	_ = s
}
func (this_ *J2clPass_ClassStaticFunctionsInliner_StaticCallInliner) IsJavaNil_() bool {
	return this_ == nil
}
func (*J2clPass_ClassStaticFunctionsInliner_StaticCallInliner) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.J2clPass$ClassStaticFunctionsInliner$StaticCallInliner", "StaticCallInliner"
}
func (this_ *J2clPass_ClassStaticFunctionsInliner_StaticCallInliner) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*J2clPass_ClassStaticFunctionsInliner_StaticCallInliner) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clPass.ClassStaticFunctionsInliner.StaticCallInliner", "com.google.javascript.jscomp.J2clPass$ClassStaticFunctionsInliner$StaticCallInliner", func(v_ any) bool {
		return j2gort.InstanceOf[*J2clPass_ClassStaticFunctionsInliner_StaticCallInliner](v_)
	})
}
func (this_ *J2clPass_ClassStaticFunctionsInliner_StaticCallInliner) SetOuter_(o *J2clPass_ClassStaticFunctionsInliner) {
	this_.outer_ = o
}
func (this_ *J2clPass_ClassStaticFunctionsInliner_StaticCallInliner) GetOuter_() *J2clPass_ClassStaticFunctionsInliner {
	return this_.outer_
}
func (*J2clPass_ClassStaticFunctionsInliner_StaticCallInliner) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*J2clPass_ClassStaticFunctionsInliner_StaticCallInliner) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *J2clPass_ClassStaticFunctionsInliner_StaticCallInliner) initFields_J2clPass_ClassStaticFunctionsInliner_StaticCallInliner() {
	_ = this_
}

// Java constructor: StaticCallInliner()
func NewJ2clPass_ClassStaticFunctionsInliner_StaticCallInliner(outer_ *J2clPass_ClassStaticFunctionsInliner) *J2clPass_ClassStaticFunctionsInliner_StaticCallInliner {
	this_ := &J2clPass_ClassStaticFunctionsInliner_StaticCallInliner{}
	this_.outer_ = outer_
	InitJ2clPass_ClassStaticFunctionsInliner_StaticCallInliner(this_)
	return this_
}

func InitJ2clPass_ClassStaticFunctionsInliner_StaticCallInliner(this_ *J2clPass_ClassStaticFunctionsInliner_StaticCallInliner) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_J2clPass_ClassStaticFunctionsInliner_StaticCallInliner()
}

// Java method: StaticCallInliner.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *J2clPass_ClassStaticFunctionsInliner_StaticCallInliner) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !(n.IsCall()) {
		return
	}
	var qualifiedNameNode pkg_rhino.Node = n.GetFirstChild()
	_ = qualifiedNameNode
	if !(qualifiedNameNode.IsGetProp()) || !(qualifiedNameNode.IsQualifiedName()) {
		return
	}
	var qualifiedFnName string = qualifiedNameNode.GetQualifiedName()
	_ = qualifiedFnName
	var fn *J2clPass_ClassStaticFunctionsInliner_InlinableFunction = this_.outer_.fnsToInlineByQualifiedName.Get(qualifiedFnName)
	_ = fn
	if fn == nil {
		return
	}
	var fnName string = qualifiedNameNode.GetString()
	_ = fnName
	var fnImpl pkg_rhino.Node = fn.Impl()
	_ = fnImpl
	var inliningMode FunctionInjector_InliningMode = fn.InliningMode()
	_ = inliningMode
	if (inliningMode == StaticGet_FunctionInjector_InliningMode_DIRECT()) && !(NodeUtil_GetFunctionBody(fnImpl).GetFirstChild().IsReturn()) {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Attempted to direct inline function "+j2gort.StrValString(qualifiedFnName)+", but function is not a simple return.")))
	}
	var inlinedCall pkg_rhino.Node = this_.outer_.injector.UnsafeInline(NewFunctionInjector_Reference(n, t.GetScope(), t.GetChunk(), inliningMode), fnName, fnImpl)
	_ = inlinedCall
	inlinedCall.SrcrefTree(n)
	t.GetCompiler().ReportChangeToEnclosingScope(inlinedCall)
}
func (*J2clPass_ClassStaticFunctionsInliner_FunctionDefsCollector) CanReleaseSyntacticScopes_() {}
func (*J2clPass_ClassStaticFunctionsInliner_StaticCallInliner) AlwaysShouldTraverse_()          {}
