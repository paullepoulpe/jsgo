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
// Port of com/google/javascript/jscomp/JsMessageExtractor.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.JsMessageExtractor
type JsMessageExtractor struct {
	idGenerator             JsMessage_IdGenerator `gson:"idGenerator"`
	options                 *CompilerOptions      `gson:"options"`
	idHash_                 atomic.Int32
	extractExternalMessages bool `gson:"extractExternalMessages"`
}

func (this_ *JsMessageExtractor) BindSelf_(s any) {
	_ = s
}
func (this_ *JsMessageExtractor) IsJavaNil_() bool { return this_ == nil }
func (*JsMessageExtractor) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JsMessageExtractor", "JsMessageExtractor"
}
func (this_ *JsMessageExtractor) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JsMessageExtractor) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JsMessageExtractor", "com.google.javascript.jscomp.JsMessageExtractor", func(v_ any) bool {
		return j2gort.InstanceOf[*JsMessageExtractor](v_)
	})
}

func (this_ *JsMessageExtractor) initFields_JsMessageExtractor() {
	_ = this_
}

// Java constructor: JsMessageExtractor(com.google.javascript.jscomp.JsMessage.IdGenerator idGenerator)
func NewJsMessageExtractor_1(idGenerator JsMessage_IdGenerator) *JsMessageExtractor {
	this_ := &JsMessageExtractor{}
	InitJsMessageExtractor_1(this_, idGenerator)
	return this_
}

func InitJsMessageExtractor_1(this_ *JsMessageExtractor, idGenerator JsMessage_IdGenerator) {
	_ = idGenerator
	InitJsMessageExtractor_3(this_, idGenerator, NewCompilerOptions(), false)
}

// Java constructor: JsMessageExtractor(com.google.javascript.jscomp.JsMessage.IdGenerator idGenerator, com.google.javascript.jscomp.CompilerOptions options, boolean extractExternalMessages)
func NewJsMessageExtractor_3(idGenerator JsMessage_IdGenerator, options *CompilerOptions, extractExternalMessages bool) *JsMessageExtractor {
	this_ := &JsMessageExtractor{}
	InitJsMessageExtractor_3(this_, idGenerator, options, extractExternalMessages)
	return this_
}

func InitJsMessageExtractor_3(this_ *JsMessageExtractor, idGenerator JsMessage_IdGenerator, options *CompilerOptions, extractExternalMessages bool) {
	_ = idGenerator
	_ = options
	_ = extractExternalMessages
	this_.initFields_JsMessageExtractor()
	this_.idGenerator = idGenerator
	this_.options = options
	this_.extractExternalMessages = extractExternalMessages
}

// Java method: JsMessageExtractor.extractMessages(com.google.javascript.jscomp.SourceFile[] inputs)
func (this_ *JsMessageExtractor) ExtractMessages_1_SourceFile(inputs []*SourceFile) j2gort.Ext_java_util_Collection {
	_ = inputs
	var iifeArg_1 []any
iifeSwitch_1:
	switch {
	default:

		arr_1 := inputs
		if arr_1 == nil {
			iifeArg_1 = nil
			break iifeSwitch_1
		}
		out_2 := make([]any, len(arr_1))
		for i_, el_3 := range arr_1 {
			out_2[i_] = j2gort.PtrToAny(el_3)
		}
		iifeArg_1 = out_2
		break iifeSwitch_1
	}
	return this_.ExtractMessages_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_slice_any(iifeArg_1)))
}

// Java method: JsMessageExtractor.extractMessages(java.lang.Iterable<com.google.javascript.jscomp.SourceFile> inputs)
func (this_ *JsMessageExtractor) ExtractMessages_1_Ext_java_lang_Iterable(inputs j2gort.Ext_java_lang_Iterable) j2gort.Ext_java_util_Collection {
	_ = inputs
	var compiler *Compiler = NewCompiler()
	_ = compiler
	compiler.Init(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()), j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_lang_Iterable(inputs)), this_.options)
	compiler.RunInCompilerThread(j2gort.NewLambda_Ext_java_util_concurrent_Callable(func() any {
		compiler.ParseInputs()
		return nil
	}))
	var extractCompilerPass *JsMessageExtractor_ExtractMessagesVisitor = NewJsMessageExtractor_ExtractMessagesVisitor(this_, compiler)
	_ = extractCompilerPass
	if compiler.GetErrors().IsEmpty() {
		extractCompilerPass.Process_2(pkg_rhino.Node(0), compiler.GetRoot())
	}
	var errors []*JSError = j2gort.ImmListToSliceNullable[*JSError](compiler.GetErrors())
	_ = errors
	if len(errors) != 0 {
		var msg j2gort.LocalStringBuilder
		msg.InitStr("JSCompiler errors\n")

		var formatter MessageFormatter = NewLightweightMessageFormatter_1(compiler)
		_ = formatter
		{

			for _, e := range errors {

				_ = e
				msg.AppendStr(formatter.FormatError(e))
			}
		}
		panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", j2gort.BoxString(msg.String()))))
	}
	return extractCompilerPass.GetMessages()
}

// Java class: com.google.javascript.jscomp.JsMessageExtractor.ExtractMessagesVisitor
type JsMessageExtractor_ExtractMessagesVisitor struct {
	JsMessageVisitorBase
	outer_   *JsMessageExtractor
	messages j2gort.Ext_java_util_List `gson:"messages"`
}

func (this_ *JsMessageExtractor_ExtractMessagesVisitor) BindSelf_(s any) {
	_ = s
	this_.JsMessageVisitorBase.BindSelf_(s)
}
func (this_ *JsMessageExtractor_ExtractMessagesVisitor) IsJavaNil_() bool { return this_ == nil }
func (*JsMessageExtractor_ExtractMessagesVisitor) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JsMessageExtractor$ExtractMessagesVisitor", "ExtractMessagesVisitor"
}
func (this_ *JsMessageExtractor_ExtractMessagesVisitor) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.JsMessageVisitorBase.IdentityHash_()
}
func (*JsMessageExtractor_ExtractMessagesVisitor) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JsMessageExtractor.ExtractMessagesVisitor", "com.google.javascript.jscomp.JsMessageExtractor$ExtractMessagesVisitor", func(v_ any) bool {
		return j2gort.InstanceOf[*JsMessageExtractor_ExtractMessagesVisitor](v_)
	})
}
func (this_ *JsMessageExtractor_ExtractMessagesVisitor) SetOuter_(o *JsMessageExtractor) {
	this_.outer_ = o
}
func (this_ *JsMessageExtractor_ExtractMessagesVisitor) GetOuter_() *JsMessageExtractor {
	return this_.outer_
}
func (*JsMessageExtractor_ExtractMessagesVisitor) Implements_jscomp_JsMessageVisitor() {}
func (*JsMessageExtractor_ExtractMessagesVisitor) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*JsMessageExtractor_ExtractMessagesVisitor) Implements_jscomp_CompilerPass()           {}
func (*JsMessageExtractor_ExtractMessagesVisitor) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *JsMessageExtractor_ExtractMessagesVisitor) initFields_JsMessageExtractor_ExtractMessagesVisitor() {
	_ = this_
	this_.messages = j2gort.NewExt_java_util_ArrayList()
}

// Java constructor: ExtractMessagesVisitor(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewJsMessageExtractor_ExtractMessagesVisitor(outer_ *JsMessageExtractor, compiler AbstractCompiler) *JsMessageExtractor_ExtractMessagesVisitor {
	this_ := &JsMessageExtractor_ExtractMessagesVisitor{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitJsMessageExtractor_ExtractMessagesVisitor(this_, compiler)
	return this_
}

func InitJsMessageExtractor_ExtractMessagesVisitor(this_ *JsMessageExtractor_ExtractMessagesVisitor, compiler AbstractCompiler) {
	_ = compiler
	InitJsMessageVisitor(&this_.JsMessageVisitorBase, compiler, this_.outer_.idGenerator)
	this_.initFields_JsMessageExtractor_ExtractMessagesVisitor()
}

// Java method: ExtractMessagesVisitor.processJsMessageDefinition(com.google.javascript.jscomp.JsMessageDefinition definition)
func (this_ *JsMessageExtractor_ExtractMessagesVisitor) ProcessJsMessageDefinition(definition JsMessageDefinition) {
	_ = definition
	this_.priv_ExtractMessagesVisitor_processJsMessage_3670(definition.GetMessage())
}

// Java method: ExtractMessagesVisitor.processIcuTemplateDefinition(com.google.javascript.jscomp.IcuTemplateDefinition definition)
func (this_ *JsMessageExtractor_ExtractMessagesVisitor) ProcessIcuTemplateDefinition(definition IcuTemplateDefinition) {
	_ = definition
	this_.priv_ExtractMessagesVisitor_processJsMessage_3670(definition.GetMessage())
}

// Java method: ExtractMessagesVisitor.processJsMessage(com.google.javascript.jscomp.JsMessage jsMessage)
func (this_ *JsMessageExtractor_ExtractMessagesVisitor) priv_ExtractMessagesVisitor_processJsMessage_3670(jsMessage *JsMessage) {
	_ = jsMessage
	if this_.outer_.extractExternalMessages || !(jsMessage.IsExternal()) {
		this_.messages.Add_1(j2gort.PtrToAny(jsMessage))
	}
}

// Java method: ExtractMessagesVisitor.getMessages()
func (this_ *JsMessageExtractor_ExtractMessagesVisitor) GetMessages() j2gort.Ext_java_util_Collection {
	return j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.messages)
}
func (*JsMessageExtractor_ExtractMessagesVisitor) AlwaysShouldTraverse_() {}
