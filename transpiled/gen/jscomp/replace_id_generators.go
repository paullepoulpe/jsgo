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
// Port of com/google/javascript/jscomp/ReplaceIdGenerators.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	pkg_sourcemap "github.com/paullepoulpe/jsgo/transpiled/gen/sourcemap"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.ReplaceIdGenerators
type ReplaceIdGenerators struct {
	compiler                      AbstractCompiler                                            `gson:"compiler"`
	nameGenerators                j2gort.OrderedMap[string, ReplaceIdGenerators_NameSupplier] `gson:"nameGenerators"`
	consistNameMap                j2gort.OrderedMap[string, j2gort.Ext_java_util_Map]         `gson:"consistNameMap"`
	idGeneratorMaps               j2gort.Ext_java_util_Map                                    `gson:"idGeneratorMaps"`
	previousMap                   j2gort.Ext_java_util_Map                                    `gson:"previousMap"`
	xidHashFunction               Xid_HashFunction                                            `gson:"xidHashFunction"`
	idHash_                       atomic.Int32
	templateLiteralsAreTranspiled bool `gson:"templateLiteralsAreTranspiled"`
	generatePseudoNames           bool `gson:"generatePseudoNames"`
}

func (this_ *ReplaceIdGenerators) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceIdGenerators) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceIdGenerators) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceIdGenerators", "ReplaceIdGenerators"
}
func (this_ *ReplaceIdGenerators) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ReplaceIdGenerators) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceIdGenerators", "com.google.javascript.jscomp.ReplaceIdGenerators", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceIdGenerators](v_)
	})
}
func (*ReplaceIdGenerators) Implements_jscomp_CompilerPass() {}

var static_ReplaceIdGenerators_NON_GLOBAL_ID_GENERATOR_CALL *DiagnosticType

func StaticGet_ReplaceIdGenerators_NON_GLOBAL_ID_GENERATOR_CALL() *DiagnosticType {
	clinit_ReplaceIdGenerators()
	return static_ReplaceIdGenerators_NON_GLOBAL_ID_GENERATOR_CALL
}
func StaticSet_ReplaceIdGenerators_NON_GLOBAL_ID_GENERATOR_CALL(v *DiagnosticType) {
	clinit_ReplaceIdGenerators()
	static_ReplaceIdGenerators_NON_GLOBAL_ID_GENERATOR_CALL = v
}

var static_ReplaceIdGenerators_CONDITIONAL_ID_GENERATOR_CALL *DiagnosticType

func StaticGet_ReplaceIdGenerators_CONDITIONAL_ID_GENERATOR_CALL() *DiagnosticType {
	clinit_ReplaceIdGenerators()
	return static_ReplaceIdGenerators_CONDITIONAL_ID_GENERATOR_CALL
}
func StaticSet_ReplaceIdGenerators_CONDITIONAL_ID_GENERATOR_CALL(v *DiagnosticType) {
	clinit_ReplaceIdGenerators()
	static_ReplaceIdGenerators_CONDITIONAL_ID_GENERATOR_CALL = v
}

var static_ReplaceIdGenerators_MISSING_NAME_MAP_FOR_GENERATOR *DiagnosticType

func StaticGet_ReplaceIdGenerators_MISSING_NAME_MAP_FOR_GENERATOR() *DiagnosticType {
	clinit_ReplaceIdGenerators()
	return static_ReplaceIdGenerators_MISSING_NAME_MAP_FOR_GENERATOR
}
func StaticSet_ReplaceIdGenerators_MISSING_NAME_MAP_FOR_GENERATOR(v *DiagnosticType) {
	clinit_ReplaceIdGenerators()
	static_ReplaceIdGenerators_MISSING_NAME_MAP_FOR_GENERATOR = v
}

var static_ReplaceIdGenerators_INVALID_GENERATOR_PARAMETER *DiagnosticType

func StaticGet_ReplaceIdGenerators_INVALID_GENERATOR_PARAMETER() *DiagnosticType {
	clinit_ReplaceIdGenerators()
	return static_ReplaceIdGenerators_INVALID_GENERATOR_PARAMETER
}
func StaticSet_ReplaceIdGenerators_INVALID_GENERATOR_PARAMETER(v *DiagnosticType) {
	clinit_ReplaceIdGenerators()
	static_ReplaceIdGenerators_INVALID_GENERATOR_PARAMETER = v
}

var static_ReplaceIdGenerators_INVALID_TEMPLATE_LITERAL_PARAMETER *DiagnosticType

func StaticGet_ReplaceIdGenerators_INVALID_TEMPLATE_LITERAL_PARAMETER() *DiagnosticType {
	clinit_ReplaceIdGenerators()
	return static_ReplaceIdGenerators_INVALID_TEMPLATE_LITERAL_PARAMETER
}
func StaticSet_ReplaceIdGenerators_INVALID_TEMPLATE_LITERAL_PARAMETER(v *DiagnosticType) {
	clinit_ReplaceIdGenerators()
	static_ReplaceIdGenerators_INVALID_TEMPLATE_LITERAL_PARAMETER = v
}

var static_ReplaceIdGenerators_SHORTHAND_FUNCTION_NOT_SUPPORTED_IN_ID_GEN *DiagnosticType

func StaticGet_ReplaceIdGenerators_SHORTHAND_FUNCTION_NOT_SUPPORTED_IN_ID_GEN() *DiagnosticType {
	clinit_ReplaceIdGenerators()
	return static_ReplaceIdGenerators_SHORTHAND_FUNCTION_NOT_SUPPORTED_IN_ID_GEN
}
func StaticSet_ReplaceIdGenerators_SHORTHAND_FUNCTION_NOT_SUPPORTED_IN_ID_GEN(v *DiagnosticType) {
	clinit_ReplaceIdGenerators()
	static_ReplaceIdGenerators_SHORTHAND_FUNCTION_NOT_SUPPORTED_IN_ID_GEN = v
}

var static_ReplaceIdGenerators_COMPUTED_PROP_NOT_SUPPORTED_IN_ID_GEN *DiagnosticType

func StaticGet_ReplaceIdGenerators_COMPUTED_PROP_NOT_SUPPORTED_IN_ID_GEN() *DiagnosticType {
	clinit_ReplaceIdGenerators()
	return static_ReplaceIdGenerators_COMPUTED_PROP_NOT_SUPPORTED_IN_ID_GEN
}
func StaticSet_ReplaceIdGenerators_COMPUTED_PROP_NOT_SUPPORTED_IN_ID_GEN(v *DiagnosticType) {
	clinit_ReplaceIdGenerators()
	static_ReplaceIdGenerators_COMPUTED_PROP_NOT_SUPPORTED_IN_ID_GEN = v
}

var static_ReplaceIdGenerators_CREATE_TEMPLATE_TAG_FIRST_ARG pkg_rhino.QualifiedName

func StaticGet_ReplaceIdGenerators_CREATE_TEMPLATE_TAG_FIRST_ARG() pkg_rhino.QualifiedName {
	clinit_ReplaceIdGenerators()
	return static_ReplaceIdGenerators_CREATE_TEMPLATE_TAG_FIRST_ARG
}
func StaticSet_ReplaceIdGenerators_CREATE_TEMPLATE_TAG_FIRST_ARG(v pkg_rhino.QualifiedName) {
	clinit_ReplaceIdGenerators()
	static_ReplaceIdGenerators_CREATE_TEMPLATE_TAG_FIRST_ARG = v
}

var clinitDone_ReplaceIdGenerators uint32
var clinitState_ReplaceIdGenerators j2gort.ClinitState

func clinit_ReplaceIdGenerators() {
	if atomic.LoadUint32(&clinitDone_ReplaceIdGenerators) != 0 {
		return
	}
	clinitSlow_ReplaceIdGenerators()
}

func clinitSlow_ReplaceIdGenerators() {
	if !j2gort.EnterClinit(&clinitState_ReplaceIdGenerators, &clinitDone_ReplaceIdGenerators) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ReplaceIdGenerators, &clinitDone_ReplaceIdGenerators)
	static_ReplaceIdGenerators_NON_GLOBAL_ID_GENERATOR_CALL = DiagnosticType_Error("JSC_NON_GLOBAL_ID_GENERATOR_CALL", "Id generator call must be in the global scope")
	static_ReplaceIdGenerators_CONDITIONAL_ID_GENERATOR_CALL = DiagnosticType_Error("JSC_CONDITIONAL_ID_GENERATOR_CALL", "Id generator call must be unconditional")
	static_ReplaceIdGenerators_MISSING_NAME_MAP_FOR_GENERATOR = DiagnosticType_Warning("JSC_MISSING_NAME_MAP_FOR_GENERATOR", "The mapped id generator, does not have a renaming map supplied.")
	static_ReplaceIdGenerators_INVALID_GENERATOR_PARAMETER = DiagnosticType_Warning("JSC_INVALID_GENERATOR_PARAMETER", "An id generator must be called with a literal.")
	static_ReplaceIdGenerators_INVALID_TEMPLATE_LITERAL_PARAMETER = DiagnosticType_Warning("JSC_INVALID_GENERATOR_PARAMETER", "An id generator must be called with a template literals with no invalid escape sequences.")
	static_ReplaceIdGenerators_SHORTHAND_FUNCTION_NOT_SUPPORTED_IN_ID_GEN = DiagnosticType_Error("JSC_SHORTHAND_FUNCTION_NOT_SUPPORTED_IN_ID_GEN", "Object literal shorthand functions is not allowed in the arguments of an id generator")
	static_ReplaceIdGenerators_COMPUTED_PROP_NOT_SUPPORTED_IN_ID_GEN = DiagnosticType_Error("JSC_COMPUTED_PROP_NOT_SUPPORTED_IN_ID_GEN", "Object literal computed property name is not allowed in the arguments of an id generator")
	static_ReplaceIdGenerators_CREATE_TEMPLATE_TAG_FIRST_ARG = pkg_rhino.QualifiedName_Of("$jscomp.createTemplateTagFirstArg")
}

func (this_ *ReplaceIdGenerators) initFields_ReplaceIdGenerators() {
	_ = this_
}

// Java constructor: ReplaceIdGenerators(com.google.javascript.jscomp.AbstractCompiler compiler, boolean templateLiteralsAreTranspiled, java.util.Map<java.lang.String,com.google.javascript.jscomp.RenamingMap> idGens, boolean generatePseudoNames, java.lang.String previousMapSerialized, com.google.javascript.jscomp.Xid.HashFunction xidHashFunction)
func NewReplaceIdGenerators(compiler AbstractCompiler, templateLiteralsAreTranspiled bool, idGens j2gort.Ext_java_util_Map, generatePseudoNames bool, previousMapSerialized string, xidHashFunction Xid_HashFunction) *ReplaceIdGenerators {
	clinit_ReplaceIdGenerators()
	this_ := &ReplaceIdGenerators{}
	InitReplaceIdGenerators(this_, compiler, templateLiteralsAreTranspiled, idGens, generatePseudoNames, previousMapSerialized, xidHashFunction)
	return this_
}

func InitReplaceIdGenerators(this_ *ReplaceIdGenerators, compiler AbstractCompiler, templateLiteralsAreTranspiled bool, idGens j2gort.Ext_java_util_Map, generatePseudoNames bool, previousMapSerialized string, xidHashFunction Xid_HashFunction) {
	_ = compiler
	_ = templateLiteralsAreTranspiled
	_ = idGens
	_ = generatePseudoNames
	_ = previousMapSerialized
	_ = xidHashFunction
	this_.initFields_ReplaceIdGenerators()
	this_.compiler = compiler
	this_.templateLiteralsAreTranspiled = templateLiteralsAreTranspiled
	this_.generatePseudoNames = generatePseudoNames
	this_.xidHashFunction = xidHashFunction
	this_.nameGenerators = j2gort.NewOrderedMap[string, ReplaceIdGenerators_NameSupplier]()
	this_.idGeneratorMaps = j2gort.NewExt_java_util_LinkedHashMap()
	this_.consistNameMap = j2gort.NewOrderedMap[string, j2gort.Ext_java_util_Map]()
	var previousMap j2gort.Ext_java_util_Map = IdMappingUtil_ParseSerializedIdMappings_1_string(previousMapSerialized)
	_ = previousMap
	this_.previousMap = previousMap
	if idGens != nil {
		{
			cur_1 := j2gort.NewMapEntryCursor(idGens)
			for cur_1.NextEntry() {
				var name string = cur_1.KeyStr()
				_ = name
				var map_ RenamingMap = j2gort.AsIface[RenamingMap](cur_1.GetValue())
				_ = map_
				var renamingToken *RenamingToken
				_ = renamingToken
				var ok_2 bool
				renamingToken, ok_2 = j2gort.TryCast[*RenamingToken](map_)
				if ok_2 {
					switch j2gort.CheckNotNullPtr(renamingToken).Ordinal() {
					case int32(2) /* RenamingToken.DISABLE */ :
						this_.nameGenerators.Put(name, nil)
						continue
					case int32(0) /* RenamingToken.INCONSISTENT */ :
						this_.nameGenerators.Put(name, this_.priv_ReplaceIdGenerators_createNameSupplier_4824(StaticGet_ReplaceIdGenerators_RenameStrategy_INCONSISTENT(), j2gort.AsIface[j2gort.Ext_com_google_common_collect_BiMap](j2gort.MapGetStr(previousMap, name))))
					case int32(1) /* RenamingToken.STABLE */ :
						this_.nameGenerators.Put(name, this_.priv_ReplaceIdGenerators_createNameSupplier_4824(StaticGet_ReplaceIdGenerators_RenameStrategy_STABLE(), j2gort.AsIface[j2gort.Ext_com_google_common_collect_BiMap](j2gort.MapGetStr(previousMap, name))))
					}
				} else {
					this_.nameGenerators.Put(name, ReplaceIdGenerators_CreateNameSupplier(StaticGet_ReplaceIdGenerators_RenameStrategy_MAPPED(), map_))
				}
				j2gort.MapPutStr(this_.idGeneratorMaps, name, j2gort.NewExt_java_util_LinkedHashMap())
			}
		}
	}
}

// Java method: ReplaceIdGenerators.createNameSupplier(com.google.javascript.jscomp.ReplaceIdGenerators.RenameStrategy renameStrategy, com.google.common.collect.BiMap<java.lang.String,java.lang.String> previousMappings)
func (this_ *ReplaceIdGenerators) priv_ReplaceIdGenerators_createNameSupplier_4824(renameStrategy ReplaceIdGenerators_RenameStrategy, previousMappings j2gort.Ext_com_google_common_collect_BiMap) ReplaceIdGenerators_NameSupplier {
	_ = renameStrategy
	_ = previousMappings
	if previousMappings != nil {
		previousMappings = previousMappings
	} else {
		previousMappings = j2gort.AsIface[j2gort.Ext_com_google_common_collect_BiMap](j2gort.Ext_com_google_common_collect_ImmutableBiMap_Of_0())
	}
	if renameStrategy == StaticGet_ReplaceIdGenerators_RenameStrategy_STABLE() {
		return NewReplaceIdGenerators_StableNameSupplier()
	} else if renameStrategy == StaticGet_ReplaceIdGenerators_RenameStrategy_XID() {
		return NewReplaceIdGenerators_XidNameSupplier(this_.xidHashFunction)
	} else if this_.generatePseudoNames {
		return NewReplaceIdGenerators_PseudoNameSupplier(renameStrategy)
	} else {
		return NewReplaceIdGenerators_ObfuscatedNameSupplier(renameStrategy, previousMappings)
	}
}

// Java static method: ReplaceIdGenerators.createNameSupplier(com.google.javascript.jscomp.ReplaceIdGenerators.RenameStrategy renameStrategy, com.google.javascript.jscomp.RenamingMap mappings)
func ReplaceIdGenerators_CreateNameSupplier(renameStrategy ReplaceIdGenerators_RenameStrategy, mappings RenamingMap) ReplaceIdGenerators_NameSupplier {
	clinit_ReplaceIdGenerators()
	_ = renameStrategy
	_ = mappings
	if !(renameStrategy == StaticGet_ReplaceIdGenerators_RenameStrategy_MAPPED()) {
		j2gort.CheckState(false)
	}
	return NewReplaceIdGenerators_MappedNameSupplier(mappings)
}

// Java method: ReplaceIdGenerators.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ReplaceIdGenerators) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, NewReplaceIdGenerators_GatherGenerators(this_))
	if !(this_.nameGenerators.IsEmpty()) {
		NodeTraversal_Traverse(this_.compiler, root, NewReplaceIdGenerators_ReplaceGenerators(this_))
	}
}

// Java method: ReplaceIdGenerators.getSerializedIdMappings()
func (this_ *ReplaceIdGenerators) GetSerializedIdMappings() string {
	return IdMappingUtil_GenerateSerializedIdMappings(this_.idGeneratorMaps)
}

// Java static method: ReplaceIdGenerators.getIdForGeneratorNode(boolean consistent, com.google.javascript.rhino.Node n, java.lang.String name)
func ReplaceIdGenerators_GetIdForGeneratorNode(consistent bool, n pkg_rhino.Node, name string) string {
	clinit_ReplaceIdGenerators()
	_ = consistent
	_ = n
	_ = name
	if consistent {
		return name
	} else {
		return (j2gort.StrValString(n.GetSourceFileName()) + ":" + j2gort.StrValInt32(n.GetLineno()) + ":" + j2gort.StrValInt32(n.GetCharno()))
	}
}

// Java enum: com.google.javascript.jscomp.ReplaceIdGenerators.RenameStrategy
type ReplaceIdGenerators_RenameStrategy uint8

const (
	enum_ReplaceIdGenerators_RenameStrategy_CONSISTENT   ReplaceIdGenerators_RenameStrategy = 1
	enum_ReplaceIdGenerators_RenameStrategy_INCONSISTENT ReplaceIdGenerators_RenameStrategy = 2
	enum_ReplaceIdGenerators_RenameStrategy_MAPPED       ReplaceIdGenerators_RenameStrategy = 3
	enum_ReplaceIdGenerators_RenameStrategy_STABLE       ReplaceIdGenerators_RenameStrategy = 4
	enum_ReplaceIdGenerators_RenameStrategy_XID          ReplaceIdGenerators_RenameStrategy = 5
)

var enumBaseTable_ReplaceIdGenerators_RenameStrategy = [...]j2gort.EnumBase{
	{Name_: "CONSISTENT", Ordinal_: 0},
	{Name_: "INCONSISTENT", Ordinal_: 1},
	{Name_: "MAPPED", Ordinal_: 2},
	{Name_: "STABLE", Ordinal_: 3},
	{Name_: "XID", Ordinal_: 4},
}
var enumValues_ReplaceIdGenerators_RenameStrategy = []ReplaceIdGenerators_RenameStrategy{1, 2, 3, 4, 5}

func StaticGet_ReplaceIdGenerators_RenameStrategy_CONSISTENT() ReplaceIdGenerators_RenameStrategy {
	return enum_ReplaceIdGenerators_RenameStrategy_CONSISTENT
}
func StaticGet_ReplaceIdGenerators_RenameStrategy_INCONSISTENT() ReplaceIdGenerators_RenameStrategy {
	return enum_ReplaceIdGenerators_RenameStrategy_INCONSISTENT
}
func StaticGet_ReplaceIdGenerators_RenameStrategy_MAPPED() ReplaceIdGenerators_RenameStrategy {
	return enum_ReplaceIdGenerators_RenameStrategy_MAPPED
}
func StaticGet_ReplaceIdGenerators_RenameStrategy_STABLE() ReplaceIdGenerators_RenameStrategy {
	return enum_ReplaceIdGenerators_RenameStrategy_STABLE
}
func StaticGet_ReplaceIdGenerators_RenameStrategy_XID() ReplaceIdGenerators_RenameStrategy {
	return enum_ReplaceIdGenerators_RenameStrategy_XID
}

func ReplaceIdGenerators_RenameStrategy_Values() []ReplaceIdGenerators_RenameStrategy {
	return j2gort.CloneSlice(enumValues_ReplaceIdGenerators_RenameStrategy)
}

func ReplaceIdGenerators_RenameStrategy_ValueOf(name string) ReplaceIdGenerators_RenameStrategy {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_ReplaceIdGenerators_RenameStrategy {
		if enumBaseTable_ReplaceIdGenerators_RenameStrategy[i_].Name_ == name {
			return ReplaceIdGenerators_RenameStrategy(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.ReplaceIdGenerators.RenameStrategy." + name))
}

func (this_ ReplaceIdGenerators_RenameStrategy) IsJavaNil_() bool                              { return this_ == 0 }
func (ReplaceIdGenerators_RenameStrategy) HasIdentityEquals_()                                 {}
func (ReplaceIdGenerators_RenameStrategy) Implements_j2gort_Ext_java_lang_Enum()               {}
func (ReplaceIdGenerators_RenameStrategy) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ ReplaceIdGenerators_RenameStrategy) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (ReplaceIdGenerators_RenameStrategy) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceIdGenerators$RenameStrategy", "RenameStrategy"
}
func (this_ ReplaceIdGenerators_RenameStrategy) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_ReplaceIdGenerators_RenameStrategy[this_-1]
}
func (this_ ReplaceIdGenerators_RenameStrategy) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ ReplaceIdGenerators_RenameStrategy) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ ReplaceIdGenerators_RenameStrategy) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ReplaceIdGenerators_RenameStrategy[this_-1].Name_
}
func (this_ ReplaceIdGenerators_RenameStrategy) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_ReplaceIdGenerators_RenameStrategy[this_-1].Name_
}
func (this_ ReplaceIdGenerators_RenameStrategy) Equals(other_ any) bool {
	o_, ok_ := other_.(ReplaceIdGenerators_RenameStrategy)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ ReplaceIdGenerators_RenameStrategy) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_ReplaceIdGenerators_RenameStrategy[this_-1].HashCode()
}
func (this_ ReplaceIdGenerators_RenameStrategy) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ReplaceIdGenerators_RenameStrategy[this_-1].CompareTo(other_)
}
func (this_ ReplaceIdGenerators_RenameStrategy) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_ReplaceIdGenerators_RenameStrategy[this_-1].Name_)
}
func init() {
	vals_ := []any{ReplaceIdGenerators_RenameStrategy(1), ReplaceIdGenerators_RenameStrategy(2), ReplaceIdGenerators_RenameStrategy(3), ReplaceIdGenerators_RenameStrategy(4), ReplaceIdGenerators_RenameStrategy(5)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ReplaceIdGenerators.RenameStrategy", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ReplaceIdGenerators$RenameStrategy", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceIdGenerators.RenameStrategy", "com.google.javascript.jscomp.ReplaceIdGenerators$RenameStrategy", func(v_ any) bool {
		return j2gort.InstanceOf[ReplaceIdGenerators_RenameStrategy](v_)
	})
}

// Java interface: com.google.javascript.jscomp.ReplaceIdGenerators.NameSupplier
type ReplaceIdGenerators_NameSupplier interface {
	Implements_jscomp_ReplaceIdGenerators_NameSupplier()
	IsJavaNil_() bool
	GetName(string, string) string
	GetRenameStrategy() ReplaceIdGenerators_RenameStrategy
}

func ToIface_ReplaceIdGenerators_NameSupplier[T interface {
	ReplaceIdGenerators_NameSupplier
	comparable
}](v T) ReplaceIdGenerators_NameSupplier {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceIdGenerators.NameSupplier", "com.google.javascript.jscomp.ReplaceIdGenerators$NameSupplier", func(v_ any) bool {
		return j2gort.InstanceOf[ReplaceIdGenerators_NameSupplier](v_)
	})
}

// Java class: com.google.javascript.jscomp.ReplaceIdGenerators.ObfuscatedNameSupplier
type ReplaceIdGenerators_ObfuscatedNameSupplier struct {
	generator        NameGenerator            `gson:"generator"`
	previousMappings j2gort.Ext_java_util_Map `gson:"previousMappings"`
	idHash_          atomic.Int32
	renameStrategy   ReplaceIdGenerators_RenameStrategy `gson:"renameStrategy"`
}

func (this_ *ReplaceIdGenerators_ObfuscatedNameSupplier) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceIdGenerators_ObfuscatedNameSupplier) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceIdGenerators_ObfuscatedNameSupplier) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceIdGenerators$ObfuscatedNameSupplier", "ObfuscatedNameSupplier"
}
func (this_ *ReplaceIdGenerators_ObfuscatedNameSupplier) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ReplaceIdGenerators_ObfuscatedNameSupplier) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceIdGenerators.ObfuscatedNameSupplier", "com.google.javascript.jscomp.ReplaceIdGenerators$ObfuscatedNameSupplier", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceIdGenerators_ObfuscatedNameSupplier](v_)
	})
}
func (*ReplaceIdGenerators_ObfuscatedNameSupplier) Implements_jscomp_ReplaceIdGenerators_NameSupplier() {
}

func (this_ *ReplaceIdGenerators_ObfuscatedNameSupplier) initFields_ReplaceIdGenerators_ObfuscatedNameSupplier() {
	_ = this_
}

// Java constructor: ObfuscatedNameSupplier(com.google.javascript.jscomp.ReplaceIdGenerators.RenameStrategy renameStrategy, com.google.common.collect.BiMap<java.lang.String,java.lang.String> previousMappings)
func NewReplaceIdGenerators_ObfuscatedNameSupplier(renameStrategy ReplaceIdGenerators_RenameStrategy, previousMappings j2gort.Ext_com_google_common_collect_BiMap) *ReplaceIdGenerators_ObfuscatedNameSupplier {
	this_ := &ReplaceIdGenerators_ObfuscatedNameSupplier{}
	InitReplaceIdGenerators_ObfuscatedNameSupplier(this_, renameStrategy, previousMappings)
	return this_
}

func InitReplaceIdGenerators_ObfuscatedNameSupplier(this_ *ReplaceIdGenerators_ObfuscatedNameSupplier, renameStrategy ReplaceIdGenerators_RenameStrategy, previousMappings j2gort.Ext_com_google_common_collect_BiMap) {
	_ = renameStrategy
	_ = previousMappings
	this_.initFields_ReplaceIdGenerators_ObfuscatedNameSupplier()
	this_.previousMappings = j2gort.AsIface[j2gort.Ext_com_google_common_collect_BiMap](previousMappings.Inverse())
	this_.generator = NewDefaultNameGenerator_3(j2gort.AsIface[j2gort.Ext_java_util_Set](previousMappings.KeySet_0()), "", j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()))
	this_.renameStrategy = renameStrategy
}

// Java method: ObfuscatedNameSupplier.getName(java.lang.String id, java.lang.String name)
func (this_ *ReplaceIdGenerators_ObfuscatedNameSupplier) GetName(id string, name string) string {
	_ = id
	_ = name
	var newName string = j2gort.CastString(j2gort.MapGetStr(this_.previousMappings, id))
	_ = newName
	if newName == j2gort.NullString {
		newName = this_.generator.GenerateNextName()
	}
	return newName
}

// Java method: ObfuscatedNameSupplier.getRenameStrategy()
func (this_ *ReplaceIdGenerators_ObfuscatedNameSupplier) GetRenameStrategy() ReplaceIdGenerators_RenameStrategy {
	return this_.renameStrategy
}

// Java class: com.google.javascript.jscomp.ReplaceIdGenerators.PseudoNameSupplier
type ReplaceIdGenerators_PseudoNameSupplier struct {
	counter        int32 `gson:"counter"`
	idHash_        atomic.Int32
	renameStrategy ReplaceIdGenerators_RenameStrategy `gson:"renameStrategy"`
}

func (this_ *ReplaceIdGenerators_PseudoNameSupplier) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceIdGenerators_PseudoNameSupplier) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceIdGenerators_PseudoNameSupplier) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceIdGenerators$PseudoNameSupplier", "PseudoNameSupplier"
}
func (this_ *ReplaceIdGenerators_PseudoNameSupplier) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ReplaceIdGenerators_PseudoNameSupplier) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceIdGenerators.PseudoNameSupplier", "com.google.javascript.jscomp.ReplaceIdGenerators$PseudoNameSupplier", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceIdGenerators_PseudoNameSupplier](v_)
	})
}
func (*ReplaceIdGenerators_PseudoNameSupplier) Implements_jscomp_ReplaceIdGenerators_NameSupplier() {}

func (this_ *ReplaceIdGenerators_PseudoNameSupplier) initFields_ReplaceIdGenerators_PseudoNameSupplier() {
	_ = this_

}

// Java constructor: PseudoNameSupplier(com.google.javascript.jscomp.ReplaceIdGenerators.RenameStrategy renameStrategy)
func NewReplaceIdGenerators_PseudoNameSupplier(renameStrategy ReplaceIdGenerators_RenameStrategy) *ReplaceIdGenerators_PseudoNameSupplier {
	this_ := &ReplaceIdGenerators_PseudoNameSupplier{}
	InitReplaceIdGenerators_PseudoNameSupplier(this_, renameStrategy)
	return this_
}

func InitReplaceIdGenerators_PseudoNameSupplier(this_ *ReplaceIdGenerators_PseudoNameSupplier, renameStrategy ReplaceIdGenerators_RenameStrategy) {
	_ = renameStrategy
	this_.initFields_ReplaceIdGenerators_PseudoNameSupplier()
	this_.renameStrategy = renameStrategy
}

// Java method: PseudoNameSupplier.getName(java.lang.String id, java.lang.String name)
func (this_ *ReplaceIdGenerators_PseudoNameSupplier) GetName(id string, name string) string {
	_ = id
	_ = name
	if this_.renameStrategy == StaticGet_ReplaceIdGenerators_RenameStrategy_INCONSISTENT() {
		return (j2gort.StrValString(name) + "$" + j2gort.StrValInt32(func() int32 { old_3 := this_.counter; nw_4 := old_3 + 1; this_.counter = nw_4; return old_3 }()))
	}
	return (j2gort.StrValString(name) + "$0")
}

// Java method: PseudoNameSupplier.getRenameStrategy()
func (this_ *ReplaceIdGenerators_PseudoNameSupplier) GetRenameStrategy() ReplaceIdGenerators_RenameStrategy {
	return this_.renameStrategy
}

// Java class: com.google.javascript.jscomp.ReplaceIdGenerators.StableNameSupplier
type ReplaceIdGenerators_StableNameSupplier struct {
	idHash_ atomic.Int32
}

func (this_ *ReplaceIdGenerators_StableNameSupplier) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceIdGenerators_StableNameSupplier) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceIdGenerators_StableNameSupplier) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceIdGenerators$StableNameSupplier", "StableNameSupplier"
}
func (this_ *ReplaceIdGenerators_StableNameSupplier) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ReplaceIdGenerators_StableNameSupplier) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceIdGenerators.StableNameSupplier", "com.google.javascript.jscomp.ReplaceIdGenerators$StableNameSupplier", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceIdGenerators_StableNameSupplier](v_)
	})
}
func (*ReplaceIdGenerators_StableNameSupplier) Implements_jscomp_ReplaceIdGenerators_NameSupplier() {}

func (this_ *ReplaceIdGenerators_StableNameSupplier) initFields_ReplaceIdGenerators_StableNameSupplier() {
	_ = this_
}

// Java constructor: StableNameSupplier()
func NewReplaceIdGenerators_StableNameSupplier() *ReplaceIdGenerators_StableNameSupplier {
	this_ := &ReplaceIdGenerators_StableNameSupplier{}
	InitReplaceIdGenerators_StableNameSupplier(this_)
	return this_
}

func InitReplaceIdGenerators_StableNameSupplier(this_ *ReplaceIdGenerators_StableNameSupplier) {
	this_.initFields_ReplaceIdGenerators_StableNameSupplier()
}

// Java method: StableNameSupplier.getName(java.lang.String id, java.lang.String name)
func (this_ *ReplaceIdGenerators_StableNameSupplier) GetName(id string, name string) string {
	_ = id
	_ = name
	return pkg_sourcemap.Base64_Base64EncodeInt(javalang.StringHashCode(j2gort.CheckNotNullString(name)))
}

// Java method: StableNameSupplier.getRenameStrategy()
func (this_ *ReplaceIdGenerators_StableNameSupplier) GetRenameStrategy() ReplaceIdGenerators_RenameStrategy {
	return StaticGet_ReplaceIdGenerators_RenameStrategy_STABLE()
}

// Java class: com.google.javascript.jscomp.ReplaceIdGenerators.XidNameSupplier
type ReplaceIdGenerators_XidNameSupplier struct {
	xid     *Xid `gson:"xid"`
	idHash_ atomic.Int32
}

func (this_ *ReplaceIdGenerators_XidNameSupplier) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceIdGenerators_XidNameSupplier) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceIdGenerators_XidNameSupplier) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceIdGenerators$XidNameSupplier", "XidNameSupplier"
}
func (this_ *ReplaceIdGenerators_XidNameSupplier) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ReplaceIdGenerators_XidNameSupplier) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceIdGenerators.XidNameSupplier", "com.google.javascript.jscomp.ReplaceIdGenerators$XidNameSupplier", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceIdGenerators_XidNameSupplier](v_)
	})
}
func (*ReplaceIdGenerators_XidNameSupplier) Implements_jscomp_ReplaceIdGenerators_NameSupplier() {}

func (this_ *ReplaceIdGenerators_XidNameSupplier) initFields_ReplaceIdGenerators_XidNameSupplier() {
	_ = this_
}

// Java constructor: XidNameSupplier(com.google.javascript.jscomp.Xid.HashFunction hashFunction)
func NewReplaceIdGenerators_XidNameSupplier(hashFunction Xid_HashFunction) *ReplaceIdGenerators_XidNameSupplier {
	this_ := &ReplaceIdGenerators_XidNameSupplier{}
	InitReplaceIdGenerators_XidNameSupplier(this_, hashFunction)
	return this_
}

func InitReplaceIdGenerators_XidNameSupplier(this_ *ReplaceIdGenerators_XidNameSupplier, hashFunction Xid_HashFunction) {
	_ = hashFunction
	this_.initFields_ReplaceIdGenerators_XidNameSupplier()
	if hashFunction == nil || hashFunction.IsJavaNil_() {
		this_.xid = NewXid()
	} else {
		this_.xid = NewXid_1(hashFunction)
	}
}

// Java method: XidNameSupplier.getName(java.lang.String id, java.lang.String name)
func (this_ *ReplaceIdGenerators_XidNameSupplier) GetName(id string, name string) string {
	_ = id
	_ = name
	return this_.xid.Get(name)
}

// Java method: XidNameSupplier.getRenameStrategy()
func (this_ *ReplaceIdGenerators_XidNameSupplier) GetRenameStrategy() ReplaceIdGenerators_RenameStrategy {
	return StaticGet_ReplaceIdGenerators_RenameStrategy_XID()
}

// Java class: com.google.javascript.jscomp.ReplaceIdGenerators.MappedNameSupplier
type ReplaceIdGenerators_MappedNameSupplier struct {
	map_    RenamingMap `gson:"map"`
	idHash_ atomic.Int32
}

func (this_ *ReplaceIdGenerators_MappedNameSupplier) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceIdGenerators_MappedNameSupplier) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceIdGenerators_MappedNameSupplier) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceIdGenerators$MappedNameSupplier", "MappedNameSupplier"
}
func (this_ *ReplaceIdGenerators_MappedNameSupplier) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ReplaceIdGenerators_MappedNameSupplier) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceIdGenerators.MappedNameSupplier", "com.google.javascript.jscomp.ReplaceIdGenerators$MappedNameSupplier", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceIdGenerators_MappedNameSupplier](v_)
	})
}
func (*ReplaceIdGenerators_MappedNameSupplier) Implements_jscomp_ReplaceIdGenerators_NameSupplier() {}

func (this_ *ReplaceIdGenerators_MappedNameSupplier) initFields_ReplaceIdGenerators_MappedNameSupplier() {
	_ = this_
}

// Java constructor: MappedNameSupplier(com.google.javascript.jscomp.RenamingMap map)
func NewReplaceIdGenerators_MappedNameSupplier(map_ RenamingMap) *ReplaceIdGenerators_MappedNameSupplier {
	this_ := &ReplaceIdGenerators_MappedNameSupplier{}
	InitReplaceIdGenerators_MappedNameSupplier(this_, map_)
	return this_
}

func InitReplaceIdGenerators_MappedNameSupplier(this_ *ReplaceIdGenerators_MappedNameSupplier, map_ RenamingMap) {
	_ = map_
	this_.initFields_ReplaceIdGenerators_MappedNameSupplier()
	this_.map_ = map_
}

// Java method: MappedNameSupplier.getName(java.lang.String id, java.lang.String name)
func (this_ *ReplaceIdGenerators_MappedNameSupplier) GetName(id string, name string) string {
	_ = id
	_ = name
	return this_.map_.Get(name)
}

// Java method: MappedNameSupplier.getRenameStrategy()
func (this_ *ReplaceIdGenerators_MappedNameSupplier) GetRenameStrategy() ReplaceIdGenerators_RenameStrategy {
	return StaticGet_ReplaceIdGenerators_RenameStrategy_MAPPED()
}

// Java class: com.google.javascript.jscomp.ReplaceIdGenerators.GatherGenerators
type ReplaceIdGenerators_GatherGenerators struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *ReplaceIdGenerators
}

func (this_ *ReplaceIdGenerators_GatherGenerators) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceIdGenerators_GatherGenerators) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceIdGenerators_GatherGenerators) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceIdGenerators$GatherGenerators", "GatherGenerators"
}
func (this_ *ReplaceIdGenerators_GatherGenerators) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*ReplaceIdGenerators_GatherGenerators) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceIdGenerators.GatherGenerators", "com.google.javascript.jscomp.ReplaceIdGenerators$GatherGenerators", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceIdGenerators_GatherGenerators](v_)
	})
}
func (this_ *ReplaceIdGenerators_GatherGenerators) SetOuter_(o *ReplaceIdGenerators) {
	this_.outer_ = o
}
func (this_ *ReplaceIdGenerators_GatherGenerators) GetOuter_() *ReplaceIdGenerators {
	return this_.outer_
}
func (*ReplaceIdGenerators_GatherGenerators) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*ReplaceIdGenerators_GatherGenerators) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *ReplaceIdGenerators_GatherGenerators) initFields_ReplaceIdGenerators_GatherGenerators() {
	_ = this_
}

// Java constructor: GatherGenerators()
func NewReplaceIdGenerators_GatherGenerators(outer_ *ReplaceIdGenerators) *ReplaceIdGenerators_GatherGenerators {
	this_ := &ReplaceIdGenerators_GatherGenerators{}
	this_.outer_ = outer_
	InitReplaceIdGenerators_GatherGenerators(this_)
	return this_
}

func InitReplaceIdGenerators_GatherGenerators(this_ *ReplaceIdGenerators_GatherGenerators) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_ReplaceIdGenerators_GatherGenerators()
}

// Java method: GatherGenerators.visit(com.google.javascript.jscomp.NodeTraversal unused, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ReplaceIdGenerators_GatherGenerators) Visit_3(unused *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = unused
	_ = n
	_ = parent
	var doc *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = doc
	if (doc == nil) || !(doc.IsAnyIdGenerator()) {
		return
	}
	var name string = j2gort.NullString
	_ = name
	if n.IsAssign() {
		name = n.GetFirstChild().GetQualifiedName()
	} else if NodeUtil_IsNameDeclaration(n) {
		name = n.GetFirstChild().GetString()
	} else if n.IsFunction() {
		name = n.GetFirstChild().GetString()
		if j2gort.CheckNotNullString(name) == "" {
			return
		}
	}
	if this_.outer_.nameGenerators.ContainsKey(name) {
		return
	}
	if doc.IsConsistentIdGenerator() {
		this_.outer_.consistNameMap.Put(name, j2gort.NewExt_java_util_LinkedHashMap())
		this_.outer_.nameGenerators.Put(name, this_.outer_.priv_ReplaceIdGenerators_createNameSupplier_4824(StaticGet_ReplaceIdGenerators_RenameStrategy_CONSISTENT(), j2gort.AsIface[j2gort.Ext_com_google_common_collect_BiMap](j2gort.MapGetStr(this_.outer_.previousMap, name))))
	} else if doc.IsStableIdGenerator() {
		this_.outer_.nameGenerators.Put(name, this_.outer_.priv_ReplaceIdGenerators_createNameSupplier_4824(StaticGet_ReplaceIdGenerators_RenameStrategy_STABLE(), j2gort.AsIface[j2gort.Ext_com_google_common_collect_BiMap](j2gort.MapGetStr(this_.outer_.previousMap, name))))
	} else if doc.IsXidGenerator() {
		this_.outer_.nameGenerators.Put(name, this_.outer_.priv_ReplaceIdGenerators_createNameSupplier_4824(StaticGet_ReplaceIdGenerators_RenameStrategy_XID(), j2gort.AsIface[j2gort.Ext_com_google_common_collect_BiMap](j2gort.MapGetStr(this_.outer_.previousMap, name))))
	} else if doc.IsIdGenerator() {
		this_.outer_.nameGenerators.Put(name, this_.outer_.priv_ReplaceIdGenerators_createNameSupplier_4824(StaticGet_ReplaceIdGenerators_RenameStrategy_INCONSISTENT(), j2gort.AsIface[j2gort.Ext_com_google_common_collect_BiMap](j2gort.MapGetStr(this_.outer_.previousMap, name))))
	} else if doc.IsMappedIdGenerator() {
		var supplier ReplaceIdGenerators_NameSupplier = this_.outer_.nameGenerators.Get(name)
		_ = supplier
		if (supplier == nil || supplier.IsJavaNil_()) || (supplier.GetRenameStrategy() != StaticGet_ReplaceIdGenerators_RenameStrategy_MAPPED()) {
			this_.outer_.compiler.Report_1(JSError_Make_3(n, StaticGet_ReplaceIdGenerators_MISSING_NAME_MAP_FOR_GENERATOR(), []string{}))
			return
		}
	} else {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "unexpected")))
	}
	j2gort.MapPutStr(this_.outer_.idGeneratorMaps, name, j2gort.NewExt_java_util_LinkedHashMap())
}

// Java class: com.google.javascript.jscomp.ReplaceIdGenerators.ReplaceGenerators
type ReplaceIdGenerators_ReplaceGenerators struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *ReplaceIdGenerators
}

func (this_ *ReplaceIdGenerators_ReplaceGenerators) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceIdGenerators_ReplaceGenerators) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceIdGenerators_ReplaceGenerators) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceIdGenerators$ReplaceGenerators", "ReplaceGenerators"
}
func (this_ *ReplaceIdGenerators_ReplaceGenerators) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*ReplaceIdGenerators_ReplaceGenerators) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceIdGenerators.ReplaceGenerators", "com.google.javascript.jscomp.ReplaceIdGenerators$ReplaceGenerators", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceIdGenerators_ReplaceGenerators](v_)
	})
}
func (this_ *ReplaceIdGenerators_ReplaceGenerators) SetOuter_(o *ReplaceIdGenerators) {
	this_.outer_ = o
}
func (this_ *ReplaceIdGenerators_ReplaceGenerators) GetOuter_() *ReplaceIdGenerators {
	return this_.outer_
}
func (*ReplaceIdGenerators_ReplaceGenerators) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*ReplaceIdGenerators_ReplaceGenerators) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *ReplaceIdGenerators_ReplaceGenerators) initFields_ReplaceIdGenerators_ReplaceGenerators() {
	_ = this_
}

// Java constructor: ReplaceGenerators()
func NewReplaceIdGenerators_ReplaceGenerators(outer_ *ReplaceIdGenerators) *ReplaceIdGenerators_ReplaceGenerators {
	this_ := &ReplaceIdGenerators_ReplaceGenerators{}
	this_.outer_ = outer_
	InitReplaceIdGenerators_ReplaceGenerators(this_)
	return this_
}

func InitReplaceIdGenerators_ReplaceGenerators(this_ *ReplaceIdGenerators_ReplaceGenerators) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_ReplaceIdGenerators_ReplaceGenerators()
}

// Java method: ReplaceGenerators.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ReplaceIdGenerators_ReplaceGenerators) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !(n.IsCall()) && !(n.IsTaggedTemplateLit()) {
		return
	}
	var qname pkg_rhino.Node = NodeUtil_GetCallTargetResolvingIndirectCalls(n)
	_ = qname
	var callName string = qname.GetQualifiedName()
	_ = callName
	var nameGenerator ReplaceIdGenerators_NameSupplier = this_.outer_.nameGenerators.Get(callName)
	_ = nameGenerator
	if nameGenerator == nil || nameGenerator.IsJavaNil_() {
		return
	}
	if !(t.InGlobalHoistScope()) && (nameGenerator.GetRenameStrategy() == StaticGet_ReplaceIdGenerators_RenameStrategy_INCONSISTENT()) {
		this_.outer_.compiler.Report_1(JSError_Make_3(n, StaticGet_ReplaceIdGenerators_NON_GLOBAL_ID_GENERATOR_CALL(), []string{}))
		return
	}
	if nameGenerator.GetRenameStrategy() == StaticGet_ReplaceIdGenerators_RenameStrategy_INCONSISTENT() {
		{
			curAnc_ := n
			if curAnc_ == 0 {
				j2gort.ThrowNPE()
			}
			curAnc_ = curAnc_.GetParent()
			for curAnc_ != 0 {
				ancestor := curAnc_
				curAnc_ = curAnc_.GetParent()

				_ = ancestor
				if NodeUtil_IsControlStructure(ancestor) {
					this_.outer_.compiler.Report_1(JSError_Make_3(n, StaticGet_ReplaceIdGenerators_CONDITIONAL_ID_GENERATOR_CALL(), []string{}))
					return
				}
			}
		}
	}
	if n.IsCall() {
		this_.priv_ReplaceGenerators_maybeReplaceCall_4829(t, n, callName, nameGenerator)
	} else {
		this_.priv_ReplaceGenerators_maybeReplaceTaggedTemplateLit_4828(t, n, callName, nameGenerator)
	}
}

// Java method: ReplaceGenerators.maybeReplaceTaggedTemplateLit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, java.lang.String callName, com.google.javascript.jscomp.ReplaceIdGenerators.NameSupplier nameGenerator)
func (this_ *ReplaceIdGenerators_ReplaceGenerators) priv_ReplaceGenerators_maybeReplaceTaggedTemplateLit_4828(t *NodeTraversal, n pkg_rhino.Node, callName string, nameGenerator ReplaceIdGenerators_NameSupplier) {
	_ = t
	_ = n
	_ = callName
	_ = nameGenerator
	var arg pkg_rhino.Node = n.GetLastChild()
	_ = arg
	if (arg == 0) || !(arg.IsTemplateLit()) {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException")))
	} else if arg.HasOneChild() {
		var cooked string = arg.GetFirstChild().GetCookedString()
		_ = cooked
		if cooked == j2gort.NullString {
			this_.outer_.compiler.Report_1(JSError_Make_3(n, StaticGet_ReplaceIdGenerators_INVALID_TEMPLATE_LITERAL_PARAMETER(), []string{}))
			return
		}
		var rename string = this_.priv_ReplaceGenerators_getObfuscatedName_4831(arg, callName, nameGenerator, cooked)
		_ = rename
		n.ReplaceWith(pkg_rhino.IR_String(rename))
		t.ReportCodeChange_0()
	} else {
		var newTemplateLit pkg_rhino.Node = pkg_rhino.IR_TemplateLiteral()
		_ = newTemplateLit
		{
			var child pkg_rhino.Node = arg.GetFirstChild()
			_ = child
			for ; child != 0; child = child.GetNext() {
				if child.IsTemplateLitString() {
					var cooked_2 string = child.GetCookedString()
					_ = cooked_2
					if cooked_2 == j2gort.NullString {
						this_.outer_.compiler.Report_1(JSError_Make_3(n, StaticGet_ReplaceIdGenerators_INVALID_TEMPLATE_LITERAL_PARAMETER(), []string{}))
						return
					}
					var rename_2 string = this_.priv_ReplaceGenerators_getObfuscatedName_4831(child, callName, nameGenerator, cooked_2)
					_ = rename_2
					newTemplateLit.AddChildToBack(pkg_rhino.IR_TemplateLiteralString(rename_2, rename_2).SrcrefIfMissing(child))
				} else {
					newTemplateLit.AddChildToBack(pkg_rhino.IR_TemplateLiteralSubstitution(child.RemoveFirstChild()).SrcrefIfMissing(child))
				}
			}
		}
		arg.ReplaceWith(newTemplateLit)
		t.ReportCodeChange_1(newTemplateLit)
	}
}

// Java method: ReplaceGenerators.maybeReplaceCall(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, java.lang.String callName, com.google.javascript.jscomp.ReplaceIdGenerators.NameSupplier nameGenerator)
func (this_ *ReplaceIdGenerators_ReplaceGenerators) priv_ReplaceGenerators_maybeReplaceCall_4829(t *NodeTraversal, n pkg_rhino.Node, callName string, nameGenerator ReplaceIdGenerators_NameSupplier) {
	_ = t
	_ = n
	_ = callName
	_ = nameGenerator
	var arg pkg_rhino.Node = n.GetSecondChild()
	_ = arg
	if arg == 0 {
		this_.outer_.compiler.Report_1(JSError_Make_3(n, StaticGet_ReplaceIdGenerators_INVALID_GENERATOR_PARAMETER(), []string{}))
	} else if arg.IsStringLit() {
		var rename string = this_.priv_ReplaceGenerators_getObfuscatedName_4831(arg, callName, nameGenerator, arg.GetString())
		_ = rename
		n.ReplaceWith(pkg_rhino.IR_String(rename))
		t.ReportCodeChange_0()
	} else if arg.IsTemplateLit() && arg.HasOneChild() {
		var cooked string = arg.GetFirstChild().GetCookedString()
		_ = cooked
		if cooked == j2gort.NullString {
			this_.outer_.compiler.Report_1(JSError_Make_3(n, StaticGet_ReplaceIdGenerators_INVALID_GENERATOR_PARAMETER(), []string{}))
		} else {
			var rename_2 string = this_.priv_ReplaceGenerators_getObfuscatedName_4831(arg, callName, nameGenerator, cooked)
			_ = rename_2
			n.ReplaceWith(pkg_rhino.IR_String(rename_2))
			t.ReportCodeChange_0()
		}
	} else if arg.IsObjectLit() {
		{
			var key pkg_rhino.Node = arg.GetFirstChild()
			_ = key
			for ; key != 0; key = key.GetNext() {
				if key.IsMemberFunctionDef() {
					this_.outer_.compiler.Report_1(JSError_Make_3(n, StaticGet_ReplaceIdGenerators_SHORTHAND_FUNCTION_NOT_SUPPORTED_IN_ID_GEN(), []string{}))
					return
				}
				if key.IsComputedProp() {
					this_.outer_.compiler.Report_1(JSError_Make_3(n, StaticGet_ReplaceIdGenerators_COMPUTED_PROP_NOT_SUPPORTED_IN_ID_GEN(), []string{}))
					return
				}
				var rename_3 string = this_.priv_ReplaceGenerators_getObfuscatedName_4831(key, callName, nameGenerator, key.GetString())
				_ = rename_3
				key.SetString(rename_3)
				key.PutBooleanProp(pkg_rhino.StaticGet_Node_QUOTED_PROP(), true)
			}
		}
		arg.Detach()
		n.ReplaceWith(arg)
		t.ReportCodeChange_0()
	} else if this_.outer_.templateLiteralsAreTranspiled && arg.IsName() {
		var ttlVar pkg_rhino.Node = this_.priv_ReplaceGenerators_findTtlVar_4830(t, arg)
		_ = ttlVar
		if ttlVar == 0 {
			this_.outer_.compiler.Report_1(JSError_Make_3(n, StaticGet_ReplaceIdGenerators_INVALID_GENERATOR_PARAMETER(), []string{}))
			return
		}
		var ttlFunctionCall pkg_rhino.Node = ttlVar.GetFirstFirstChild()
		_ = ttlFunctionCall
		var paramsArrayLiteral pkg_rhino.Node = ttlFunctionCall.GetSecondChild()
		_ = paramsArrayLiteral
		if (paramsArrayLiteral == 0) || !(paramsArrayLiteral.IsArrayLit()) {
			panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "bad transpiled structure")))
		}
		if paramsArrayLiteral.GetChildCount() == int32(1) {
			var originalName string = paramsArrayLiteral.GetFirstChild().GetString()
			_ = originalName
			var rename_4 string = this_.priv_ReplaceGenerators_getObfuscatedName_4831(arg, callName, nameGenerator, originalName)
			_ = rename_4
			n.ReplaceWith(pkg_rhino.IR_String(rename_4))
			t.ReportCodeChange_0()
			t.ReportCodeChange_1(ttlVar)
			ttlVar.Detach()
		} else {
			{
				var param pkg_rhino.Node = paramsArrayLiteral.GetFirstChild()
				_ = param
				for ; param != 0; param = param.GetNext() {
					var originalName_2 string = param.GetString()
					_ = originalName_2
					var rename_5 string = this_.priv_ReplaceGenerators_getObfuscatedName_4831(arg, callName, nameGenerator, originalName_2)
					_ = rename_5
					param.SetString(rename_5)
				}
			}
			t.ReportCodeChange_1(paramsArrayLiteral)
		}
	} else {
		this_.outer_.compiler.Report_1(JSError_Make_3(n, StaticGet_ReplaceIdGenerators_INVALID_GENERATOR_PARAMETER(), []string{}))
	}
}

// Java method: ReplaceGenerators.findTtlVar(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node arg)
func (this_ *ReplaceIdGenerators_ReplaceGenerators) priv_ReplaceGenerators_findTtlVar_4830(t *NodeTraversal, arg pkg_rhino.Node) pkg_rhino.Node {
	_ = t
	_ = arg
	var name string = arg.GetString()
	_ = name
	var ttlVarName pkg_rhino.Node = j2gort.CastPtr[Var](t.GetScope().GetVar(name)).GetNode()
	_ = ttlVarName
	var parent pkg_rhino.Node = ttlVarName.GetParent()
	_ = parent
	if !(parent.IsVar()) || !(parent.GetParent().IsScript()) {
		return pkg_rhino.Node(0)
	}
	if !(ttlVarName.GetParent().HasOneChild()) {
		j2gort.CheckStateStr(false, "The AST is normalized so there should only be one name per var")
	}
	var callNode pkg_rhino.Node = ttlVarName.GetFirstChild()
	_ = callNode
	if !(callNode.IsCall()) {
		return pkg_rhino.Node(0)
	}
	var callExpr pkg_rhino.Node = callNode.GetFirstChild()
	_ = callExpr
	if callExpr.MatchesName_1_string("$jscomp$createTemplateTagFirstArg") || StaticGet_ReplaceIdGenerators_CREATE_TEMPLATE_TAG_FIRST_ARG().Matches(callExpr) {
		return ttlVarName.GetParent()
	}
	return pkg_rhino.Node(0)
}

// Java method: ReplaceGenerators.getObfuscatedName(com.google.javascript.rhino.Node id, java.lang.String callName, com.google.javascript.jscomp.ReplaceIdGenerators.NameSupplier nameGenerator, java.lang.String name)
func (this_ *ReplaceIdGenerators_ReplaceGenerators) priv_ReplaceGenerators_getObfuscatedName_4831(id pkg_rhino.Node, callName string, nameGenerator ReplaceIdGenerators_NameSupplier, name string) string {
	_ = id
	_ = callName
	_ = nameGenerator
	_ = name
	var rename string = j2gort.NullString
	_ = rename
	var idGeneratorMap j2gort.Ext_java_util_Map = j2gort.AsIface[j2gort.Ext_java_util_Map](j2gort.MapGetStr(this_.outer_.idGeneratorMaps, callName))
	_ = idGeneratorMap
	var instanceId string = ReplaceIdGenerators_GetIdForGeneratorNode((nameGenerator.GetRenameStrategy() != StaticGet_ReplaceIdGenerators_RenameStrategy_INCONSISTENT()), id, name)
	_ = instanceId
	if nameGenerator.GetRenameStrategy() == StaticGet_ReplaceIdGenerators_RenameStrategy_CONSISTENT() {
		var entry j2gort.Ext_java_util_Map = this_.outer_.consistNameMap.Get(callName)
		_ = entry
		rename = j2gort.CastString(j2gort.MapGetStr(entry, instanceId))
		if rename == j2gort.NullString {
			rename = nameGenerator.GetName(instanceId, name)
			j2gort.MapPutStr(entry, instanceId, j2gort.BoxString(rename))
		}
	} else {
		rename = nameGenerator.GetName(instanceId, name)
	}
	j2gort.MapPutStr(idGeneratorMap, rename, j2gort.BoxString(instanceId))
	return rename
}
func (*ReplaceIdGenerators_GatherGenerators) AlwaysShouldTraverse_()       {}
func (*ReplaceIdGenerators_GatherGenerators) CanReleaseSyntacticScopes_()  {}
func (*ReplaceIdGenerators_ReplaceGenerators) AlwaysShouldTraverse_()      {}
func (*ReplaceIdGenerators_ReplaceGenerators) CanReleaseSyntacticScopes_() {}
