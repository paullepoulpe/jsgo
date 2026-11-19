/*
 * Copyright 2007 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/BasicErrorManager.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java polymorphic class interface: com.google.javascript.jscomp.BasicErrorManager
type BasicErrorManager interface {
	Implements_jscomp_BasicErrorManager()
	IsJavaNil_() bool
	GetBasicErrorManagerBase() *BasicErrorManagerBase
	SortingErrorManager
	Println(*CheckLevel, *JSError)
	PrintSummary()
}
func ToIface_BasicErrorManager[T interface {
	BasicErrorManager
	comparable
}](v T) BasicErrorManager {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_BasicErrorManager[T BasicErrorManager](v T) BasicErrorManager {
	var vi BasicErrorManager = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.BasicErrorManager
type BasicErrorManagerBase struct {
	SortingErrorManagerBase
}

func (this_ *BasicErrorManagerBase) BindSelf_(s any) {
	_ = s
	this_.SortingErrorManagerBase.BindSelf_(s)
}
func (this_ *BasicErrorManagerBase) IsJavaNil_() bool { return this_ == nil }
func (*BasicErrorManagerBase) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.BasicErrorManager", "BasicErrorManager" }
func (this_ *BasicErrorManagerBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.SortingErrorManagerBase.IdentityHash_()
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.BasicErrorManager", "com.google.javascript.jscomp.BasicErrorManager", func(v_ any) bool {
		return j2gort.InstanceOf[BasicErrorManager](v_)
	})
}
func (b *BasicErrorManagerBase) GetBasicErrorManagerBase() *BasicErrorManagerBase { return b }
func (this_ *BasicErrorManagerBase) self_BasicErrorManager_() BasicErrorManager {
	if s := this_.SortingErrorManagerBase.self_; s != nil {
		return s.(BasicErrorManager)
	}
	return nil
}
func (*BasicErrorManagerBase) Implements_jscomp_BasicErrorManager() {}
func (*BasicErrorManagerBase) Implements_jscomp_SortingErrorManager() {}
func (*BasicErrorManagerBase) Implements_jscomp_ErrorManager() {}
func (*BasicErrorManagerBase) Implements_jscomp_ErrorHandler() {}

func (this_ *BasicErrorManagerBase) initFields_BasicErrorManager() {
	_ = this_
}

func InitBasicErrorManager(this_ *BasicErrorManagerBase) {
	InitSortingErrorManager(&this_.SortingErrorManagerBase, j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()))
	this_.initFields_BasicErrorManager()
}

// Java method: BasicErrorManager.generateReport()
func (this_ *BasicErrorManagerBase) GenerateReport() {
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.SortingErrorManagerBase.GetSortedDiagnostics()))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			message := j2gort.CastPtr[SortingErrorManager_ErrorWithLevel](rawElem_2)
			_ = message
			this_.self_BasicErrorManager_().Println(message.level, message.error_)
		}
	}
	this_.self_BasicErrorManager_().PrintSummary()
}

// Java method: BasicErrorManager.println(com.google.javascript.jscomp.CheckLevel level, com.google.javascript.jscomp.JSError error)
func (this_ *BasicErrorManagerBase) Println(level *CheckLevel, error_ *JSError) {
	_ = level
	_ = error_
	panic("abstract or native method invoked")
}

// Java method: BasicErrorManager.printSummary()
func (this_ *BasicErrorManagerBase) PrintSummary() {
	panic("abstract or native method invoked")
}

