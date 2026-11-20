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
// Port of com/google/javascript/jscomp/SourceMap.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	pkg_sourcemap "github.com/paullepoulpe/jsgo/transpiled/gen/sourcemap"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.SourceMap
type SourceMap struct {
	generator                pkg_sourcemap.SourceMapGenerator  `gson:"generator"`
	prefixMappings           j2gort.Ext_java_util_List         `gson:"prefixMappings"`
	sourceLocationFixupCache j2gort.OrderedMap[string, string] `gson:"sourceLocationFixupCache"`
	mapping                  SourceFileMapping                 `gson:"mapping"`
	idHash_                  atomic.Int32
}

func (this_ *SourceMap) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceMap) IsJavaNil_() bool { return this_ == nil }
func (*SourceMap) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SourceMap", "SourceMap"
}
func (this_ *SourceMap) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SourceMap) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SourceMap", "com.google.javascript.jscomp.SourceMap", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceMap](v_)
	})
}

func (this_ *SourceMap) initFields_SourceMap() {
	_ = this_
	this_.prefixMappings = j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	this_.sourceLocationFixupCache = j2gort.NewOrderedMap[string, string]()
}

// Java constructor: SourceMap(com.google.debugging.sourcemap.SourceMapGenerator generator)
func NewSourceMap(generator pkg_sourcemap.SourceMapGenerator) *SourceMap {
	this_ := &SourceMap{}
	InitSourceMap(this_, generator)
	return this_
}

func InitSourceMap(this_ *SourceMap, generator pkg_sourcemap.SourceMapGenerator) {
	_ = generator
	this_.initFields_SourceMap()
	this_.generator = generator
}

// Java method: SourceMap.addMapping(com.google.javascript.jscomp.SourceMap.Mapping mapping)
func (this_ *SourceMap) AddMapping_1(mapping *SourceMap_Mapping) {
	_ = mapping
	this_.AddMapping_3(mapping.node, mapping.start, mapping.end)
}

// Java method: SourceMap.addMapping(com.google.javascript.rhino.Node node, com.google.debugging.sourcemap.FilePosition outputStartPosition, com.google.debugging.sourcemap.FilePosition outputEndPosition)
func (this_ *SourceMap) AddMapping_3(node pkg_rhino.Node, outputStartPosition *pkg_sourcemap.FilePosition, outputEndPosition *pkg_sourcemap.FilePosition) {
	_ = node
	_ = outputStartPosition
	_ = outputEndPosition
	this_.AddMapping_5(node, outputStartPosition.GetLine(), outputStartPosition.GetColumn(), outputEndPosition.GetLine(), outputEndPosition.GetColumn())
}

// Java method: SourceMap.addMapping(com.google.javascript.rhino.Node node, int outputStartLine, int outputStartColumn, int outputEndLine, int outputEndColumn)
func (this_ *SourceMap) AddMapping_5(node pkg_rhino.Node, outputStartLine int32, outputStartColumn int32, outputEndLine int32, outputEndColumn int32) {
	_ = node
	_ = outputStartLine
	_ = outputStartColumn
	_ = outputEndLine
	_ = outputEndColumn
	var sourceFile pkg_rhino.StaticSourceFile = node.GetStaticSourceFile()
	_ = sourceFile
	if (sourceFile == nil || sourceFile.IsJavaNil_()) || (node.GetLineno() < int32(0)) {
		return
	}
	var sourceFileName string = sourceFile.GetName()
	_ = sourceFileName
	var lineNo int32 = node.GetLineno()
	_ = lineNo
	var charNo int32 = node.GetCharno()
	_ = charNo
	var originalName string = SourceMap_GetOriginalName(node)
	_ = originalName
	if this_.mapping != nil && !this_.mapping.IsJavaNil_() {
		var sourceMapping *j2gort.Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping = this_.mapping.GetSourceMapping(sourceFileName, lineNo, charNo)
		_ = sourceMapping
		if sourceMapping == nil {
		} else if sourceMapping.Equals(j2gort.PtrToAny(j2gort.Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping_GetDefaultInstance())) {
			return
		} else {
			sourceFileName = sourceMapping.GetOriginalFile()
			lineNo = sourceMapping.GetLineNumber()
			charNo = sourceMapping.GetColumnPosition()
			var identifier string = sourceMapping.GetIdentifier()
			_ = identifier
			if sourceMapping.HasIdentifier() && !(j2gort.CheckNotNullString(identifier) == "") {
				originalName = identifier
			}
		}
	}
	sourceFileName = this_.priv_SourceMap_fixupSourceLocation_5182(sourceFileName)
	var lineBaseOffset int32 = int32(1)
	_ = lineBaseOffset
	this_.generator.AddMapping_8(sourceFileName, originalName, (lineNo - lineBaseOffset), charNo, outputStartLine, outputStartColumn, outputEndLine, outputEndColumn)
}

// Java method: SourceMap.addSourceFile(java.lang.String name, java.lang.String code)
func (this_ *SourceMap) AddSourceFile(name string, code string) {
	_ = name
	_ = code
	this_.generator.AddSourcesContent(this_.priv_SourceMap_fixupSourceLocation_5182(name), code)
}

// Java static method: SourceMap.getOriginalName(com.google.javascript.rhino.Node node)
func SourceMap_GetOriginalName(node pkg_rhino.Node) string {
	_ = node
	if node.GetOriginalName() != j2gort.NullString {
		return node.GetOriginalName()
	}
	if node.IsName() {
		var parent pkg_rhino.Node = node.GetParent()
		_ = parent
		if ((parent != 0) && parent.IsFunction()) && (parent.GetFirstChild() == node) {
			return parent.GetOriginalName()
		}
	}
	if node.IsMemberFunctionDef() {
		return node.GetFirstChild().GetOriginalName()
	}
	return j2gort.NullString
}

// Java method: SourceMap.fixupSourceLocation(java.lang.String sourceFile)
func (this_ *SourceMap) priv_SourceMap_fixupSourceLocation_5182(sourceFile string) string {
	_ = sourceFile
	if this_.prefixMappings.IsEmpty() {
		return sourceFile
	}
	var fixed string = j2gort.OrderedMapGetStr(this_.sourceLocationFixupCache, sourceFile)
	_ = fixed
	if fixed != j2gort.NullString {
		return fixed
	}
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.prefixMappings))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			mapping := j2gort.AsIface[SourceMap_LocationMapping](rawElem_2)
			_ = mapping
			fixed = mapping.Map(sourceFile)
			if fixed != j2gort.NullString {
				break
			}
		}
	}
	if fixed == j2gort.NullString {
		fixed = sourceFile
	}
	this_.sourceLocationFixupCache.Put(sourceFile, fixed)
	return fixed
}

// Java method: SourceMap.appendTo(java.lang.Appendable out, java.lang.String name)
func (this_ *SourceMap) AppendTo(out j2gort.Ext_java_lang_Appendable, name string) {
	_ = out
	_ = name
	this_.generator.AppendTo(out, this_.priv_SourceMap_fixupSourceLocation_5182(name))
}

// Java method: SourceMap.reset()
func (this_ *SourceMap) Reset() {
	this_.generator.Reset()
	this_.sourceLocationFixupCache.Clear()
}

// Java method: SourceMap.setStartingPosition(int offsetLine, int offsetIndex)
func (this_ *SourceMap) SetStartingPosition(offsetLine int32, offsetIndex int32) {
	_ = offsetLine
	_ = offsetIndex
	this_.generator.SetStartingPosition(offsetLine, offsetIndex)
}

// Java method: SourceMap.setWrapperPrefix(java.lang.String prefix)
func (this_ *SourceMap) SetWrapperPrefix(prefix string) {
	_ = prefix
	this_.generator.SetWrapperPrefix(prefix)
}

// Java method: SourceMap.validate(boolean validate)
func (this_ *SourceMap) Validate(validate bool) {
	_ = validate
	this_.generator.Validate(validate)
}

// Java method: SourceMap.setPrefixMappings(java.util.List<? extends com.google.javascript.jscomp.SourceMap.LocationMapping> sourceMapLocationMappings)
func (this_ *SourceMap) SetPrefixMappings(sourceMapLocationMappings j2gort.Ext_java_util_List) {
	_ = sourceMapLocationMappings
	this_.prefixMappings = sourceMapLocationMappings
}

// Java method: SourceMap.setSourceFileMapping(com.google.javascript.jscomp.SourceFileMapping mapping)
func (this_ *SourceMap) SetSourceFileMapping(mapping SourceFileMapping) {
	_ = mapping
	this_.mapping = mapping
}

// Java enum: com.google.javascript.jscomp.SourceMap.Format
type SourceMap_Format struct {
	j2gort.EnumBase
	fn_GetInstance func(*SourceMap_Format) *SourceMap
}

func (this_ *SourceMap_Format) GetEnumBase() *j2gort.EnumBase {
	if this_ == nil {
		return nil
	}
	return &this_.EnumBase
}
func (this_ *SourceMap_Format) EnumOrdinal_() int32 {
	if this_ == nil {
		return -1
	}
	return this_.EnumBase.Ordinal_
}
func (*SourceMap_Format) Implements_j2gort_Ext_java_lang_Enum()               {}
func (*SourceMap_Format) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ *SourceMap_Format) initFields_SourceMap_Format() {
	_ = this_
}

func initEnumSourceMap_Format(this_ *SourceMap_Format, enumName_ string, enumOrd_ int32) *SourceMap_Format {
	this_.EnumBase.Name_ = enumName_
	this_.EnumBase.Ordinal_ = enumOrd_
	this_.initFields_SourceMap_Format()
	return this_
}

var enum_SourceMap_Format_DEFAULT *SourceMap_Format

func StaticGet_SourceMap_Format_DEFAULT() *SourceMap_Format {
	clinit_SourceMap_Format()
	return enum_SourceMap_Format_DEFAULT
}

var enum_SourceMap_Format_V3 *SourceMap_Format

func StaticGet_SourceMap_Format_V3() *SourceMap_Format {
	clinit_SourceMap_Format()
	return enum_SourceMap_Format_V3
}

var clinitDone_SourceMap_Format uint32
var clinitState_SourceMap_Format j2gort.ClinitState

func clinit_SourceMap_Format() {
	if atomic.LoadUint32(&clinitDone_SourceMap_Format) != 0 {
		return
	}
	clinitSlow_SourceMap_Format()
}

func clinitSlow_SourceMap_Format() {
	if !j2gort.EnterClinit(&clinitState_SourceMap_Format, &clinitDone_SourceMap_Format) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_SourceMap_Format, &clinitDone_SourceMap_Format)
	enum_SourceMap_Format_DEFAULT = &SourceMap_Format{}
	enum_SourceMap_Format_DEFAULT.fn_GetInstance = func(this_ *SourceMap_Format) *SourceMap {
		return NewSourceMap(pkg_sourcemap.SourceMapGeneratorFactory_GetInstance(pkg_sourcemap.StaticGet_SourceMapFormat_DEFAULT()))
	}
	initEnumSourceMap_Format(enum_SourceMap_Format_DEFAULT, "DEFAULT", int32(0))
	enum_SourceMap_Format_V3 = &SourceMap_Format{}
	enum_SourceMap_Format_V3.fn_GetInstance = func(this_ *SourceMap_Format) *SourceMap {
		return NewSourceMap(pkg_sourcemap.SourceMapGeneratorFactory_GetInstance(pkg_sourcemap.StaticGet_SourceMapFormat_V3()))
	}
	initEnumSourceMap_Format(enum_SourceMap_Format_V3, "V3", int32(1))
	enumValues_SourceMap_Format = []*SourceMap_Format{enum_SourceMap_Format_DEFAULT, enum_SourceMap_Format_V3}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.SourceMap.Format", []any{enum_SourceMap_Format_DEFAULT, enum_SourceMap_Format_V3})
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.SourceMap$Format", []any{enum_SourceMap_Format_DEFAULT, enum_SourceMap_Format_V3})
}

func init() {
	j2gort.RegisterEnumConstantsLoader("com.google.javascript.jscomp.SourceMap.Format", "com.google.javascript.jscomp.SourceMap$Format", func() {
		clinit_SourceMap_Format()
	})
}

var enumValues_SourceMap_Format []*SourceMap_Format

func SourceMap_Format_Values() []*SourceMap_Format {
	clinit_SourceMap_Format()
	return j2gort.CloneSlice(enumValues_SourceMap_Format)
}

func SourceMap_Format_ValueOf(name string) *SourceMap_Format {
	clinit_SourceMap_Format()
	for _, c := range enumValues_SourceMap_Format {
		if c.Name() == name {
			return c
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.SourceMap.Format." + name))
}

// Java method: Format.getInstance()
func (this_ *SourceMap_Format) GetInstance() *SourceMap {
	if this_.fn_GetInstance != nil {
		return this_.fn_GetInstance(this_)
	}
	return this_.baseEnum_GetInstance()
}

func (this_ *SourceMap_Format) baseEnum_GetInstance() *SourceMap {
	panic("abstract or native method invoked")
}

func (this_ *SourceMap_Format) IsJavaNil_() bool { return this_ == nil }
func (*SourceMap_Format) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SourceMap$Format", "Format"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SourceMap.Format", "com.google.javascript.jscomp.SourceMap$Format", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceMap_Format](v_)
	})
}

// Java enum: com.google.javascript.jscomp.SourceMap.DetailLevel
type SourceMap_DetailLevel struct {
	j2gort.EnumBase
	fn_Apply     func(*SourceMap_DetailLevel, any) bool
	fn_TestNode_ func(*SourceMap_DetailLevel, uint32) bool
}

func (this_ *SourceMap_DetailLevel) GetEnumBase() *j2gort.EnumBase {
	if this_ == nil {
		return nil
	}
	return &this_.EnumBase
}
func (this_ *SourceMap_DetailLevel) EnumOrdinal_() int32 {
	if this_ == nil {
		return -1
	}
	return this_.EnumBase.Ordinal_
}
func (*SourceMap_DetailLevel) Implements_j2gort_Ext_java_lang_Enum()                   {}
func (*SourceMap_DetailLevel) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
func (*SourceMap_DetailLevel) Implements_j2gort_Ext_java_lang_constant_Constable()     {}
func (*SourceMap_DetailLevel) Implements_j2gort_Ext_java_util_function_Predicate()     {}
func (this_ *SourceMap_DetailLevel) initFields_SourceMap_DetailLevel() {
	_ = this_
}

func initEnumSourceMap_DetailLevel(this_ *SourceMap_DetailLevel, enumName_ string, enumOrd_ int32) *SourceMap_DetailLevel {
	this_.EnumBase.Name_ = enumName_
	this_.EnumBase.Ordinal_ = enumOrd_
	this_.initFields_SourceMap_DetailLevel()
	return this_
}

var enum_SourceMap_DetailLevel_ALL *SourceMap_DetailLevel

func StaticGet_SourceMap_DetailLevel_ALL() *SourceMap_DetailLevel {
	clinit_SourceMap_DetailLevel()
	return enum_SourceMap_DetailLevel_ALL
}

var enum_SourceMap_DetailLevel_SYMBOLS *SourceMap_DetailLevel

func StaticGet_SourceMap_DetailLevel_SYMBOLS() *SourceMap_DetailLevel {
	clinit_SourceMap_DetailLevel()
	return enum_SourceMap_DetailLevel_SYMBOLS
}

var clinitDone_SourceMap_DetailLevel uint32
var clinitState_SourceMap_DetailLevel j2gort.ClinitState

func clinit_SourceMap_DetailLevel() {
	if atomic.LoadUint32(&clinitDone_SourceMap_DetailLevel) != 0 {
		return
	}
	clinitSlow_SourceMap_DetailLevel()
}

func clinitSlow_SourceMap_DetailLevel() {
	if !j2gort.EnterClinit(&clinitState_SourceMap_DetailLevel, &clinitDone_SourceMap_DetailLevel) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_SourceMap_DetailLevel, &clinitDone_SourceMap_DetailLevel)
	enum_SourceMap_DetailLevel_ALL = &SourceMap_DetailLevel{}
	enum_SourceMap_DetailLevel_ALL.fn_TestNode_ = func(this_ *SourceMap_DetailLevel, rawNode_3 uint32) bool {
		node := pkg_rhino.Node(rawNode_3)
		_ = node
		return true
	}
	enum_SourceMap_DetailLevel_ALL.fn_Apply = func(this_ *SourceMap_DetailLevel, arg_ any) bool {
		return this_.fn_TestNode_(this_, j2gort.CastNode[uint32](arg_))
	}
	initEnumSourceMap_DetailLevel(enum_SourceMap_DetailLevel_ALL, "ALL", int32(0))
	enum_SourceMap_DetailLevel_SYMBOLS = &SourceMap_DetailLevel{}
	enum_SourceMap_DetailLevel_SYMBOLS.fn_TestNode_ = func(this_ *SourceMap_DetailLevel, rawNode_4 uint32) bool {
		node := pkg_rhino.Node(rawNode_4)
		_ = node
		return ((((((node.IsCall() || node.IsNew()) || node.IsFunction()) || node.IsName()) || NodeUtil_IsNormalOrOptChainGet(node)) || NodeUtil_MayBeObjectLitKey(node)) || node.IsTaggedTemplateLit())
	}
	enum_SourceMap_DetailLevel_SYMBOLS.fn_Apply = func(this_ *SourceMap_DetailLevel, arg_ any) bool {
		return this_.fn_TestNode_(this_, j2gort.CastNode[uint32](arg_))
	}
	initEnumSourceMap_DetailLevel(enum_SourceMap_DetailLevel_SYMBOLS, "SYMBOLS", int32(1))
	enumValues_SourceMap_DetailLevel = []*SourceMap_DetailLevel{enum_SourceMap_DetailLevel_ALL, enum_SourceMap_DetailLevel_SYMBOLS}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.SourceMap.DetailLevel", []any{enum_SourceMap_DetailLevel_ALL, enum_SourceMap_DetailLevel_SYMBOLS})
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.SourceMap$DetailLevel", []any{enum_SourceMap_DetailLevel_ALL, enum_SourceMap_DetailLevel_SYMBOLS})
}

func init() {
	j2gort.RegisterEnumConstantsLoader("com.google.javascript.jscomp.SourceMap.DetailLevel", "com.google.javascript.jscomp.SourceMap$DetailLevel", func() {
		clinit_SourceMap_DetailLevel()
	})
}

var enumValues_SourceMap_DetailLevel []*SourceMap_DetailLevel

func SourceMap_DetailLevel_Values() []*SourceMap_DetailLevel {
	clinit_SourceMap_DetailLevel()
	return j2gort.CloneSlice(enumValues_SourceMap_DetailLevel)
}

func SourceMap_DetailLevel_ValueOf(name string) *SourceMap_DetailLevel {
	clinit_SourceMap_DetailLevel()
	for _, c := range enumValues_SourceMap_DetailLevel {
		if c.Name() == name {
			return c
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.SourceMap.DetailLevel." + name))
}

func (this_ *SourceMap_DetailLevel) TestNode_(rawNode_ uint32) bool {
	if this_ == nil {
		j2gort.ThrowNPE()
	}
	if this_.fn_TestNode_ != nil {
		return this_.fn_TestNode_(this_, rawNode_)
	}
	return this_.baseEnum_TestNode_(rawNode_)
}

func (this_ *SourceMap_DetailLevel) baseEnum_TestNode_(rawNode_ uint32) bool {
	_ = rawNode_
	panic("abstract or native method invoked")
}

func (this_ *SourceMap_DetailLevel) Apply(arg_ any) bool {
	if this_ == nil {
		j2gort.ThrowNPE()
	}
	return this_.TestNode_(j2gort.CastNode[uint32](arg_))
}

func (this_ *SourceMap_DetailLevel) Test(p0_ any) bool {
	return this_.Apply(p0_)
}

func (this_ *SourceMap_DetailLevel) IsJavaNil_() bool { return this_ == nil }
func (*SourceMap_DetailLevel) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SourceMap$DetailLevel", "DetailLevel"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SourceMap.DetailLevel", "com.google.javascript.jscomp.SourceMap$DetailLevel", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceMap_DetailLevel](v_)
	})
}

// Java interface: com.google.javascript.jscomp.SourceMap.LocationMapping
type SourceMap_LocationMapping interface {
	Implements_jscomp_SourceMap_LocationMapping()
	IsJavaNil_() bool
	Map(string) string
}

func ToIface_SourceMap_LocationMapping[T interface {
	SourceMap_LocationMapping
	comparable
}](v T) SourceMap_LocationMapping {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_SourceMap_LocationMapping[T SourceMap_LocationMapping](v T) SourceMap_LocationMapping {
	var vi SourceMap_LocationMapping = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SourceMap.LocationMapping", "com.google.javascript.jscomp.SourceMap$LocationMapping", func(v_ any) bool {
		return j2gort.InstanceOf[SourceMap_LocationMapping](v_)
	})
}

type lambda_SourceMap_LocationMapping struct {
	fn func(string) string
}

func (l *lambda_SourceMap_LocationMapping) IsJavaNil_() bool { return l == nil }
func NewLambda_SourceMap_LocationMapping(fn func(string) string) SourceMap_LocationMapping {
	if fn == nil {
		return nil
	}
	return &lambda_SourceMap_LocationMapping{fn: fn}
}

func (*lambda_SourceMap_LocationMapping) Implements_jscomp_SourceMap_LocationMapping() {}
func (l *lambda_SourceMap_LocationMapping) Map(p0_ string) string {
	return l.fn(p0_)
}

// Java class: com.google.javascript.jscomp.SourceMap.PrefixLocationMapping
type SourceMap_PrefixLocationMapping struct {
	prefix      string `gson:"prefix"`
	replacement string `gson:"replacement"`
}

func (this_ *SourceMap_PrefixLocationMapping) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceMap_PrefixLocationMapping) IsJavaNil_() bool { return this_ == nil }
func (*SourceMap_PrefixLocationMapping) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SourceMap$PrefixLocationMapping", "PrefixLocationMapping"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SourceMap.PrefixLocationMapping", "com.google.javascript.jscomp.SourceMap$PrefixLocationMapping", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceMap_PrefixLocationMapping](v_)
	})
}
func (*SourceMap_PrefixLocationMapping) Implements_jscomp_SourceMap_LocationMapping() {}

func (this_ *SourceMap_PrefixLocationMapping) initFields_SourceMap_PrefixLocationMapping() {
	_ = this_
	this_.prefix = j2gort.NullString
	this_.replacement = j2gort.NullString
}

// Java constructor: PrefixLocationMapping(java.lang.String prefix, java.lang.String replacement)
func NewSourceMap_PrefixLocationMapping(prefix string, replacement string) *SourceMap_PrefixLocationMapping {
	this_ := &SourceMap_PrefixLocationMapping{}
	InitSourceMap_PrefixLocationMapping(this_, prefix, replacement)
	return this_
}

func InitSourceMap_PrefixLocationMapping(this_ *SourceMap_PrefixLocationMapping, prefix string, replacement string) {
	_ = prefix
	_ = replacement
	this_.initFields_SourceMap_PrefixLocationMapping()
	this_.prefix = prefix
	this_.replacement = replacement
}

// Java method: PrefixLocationMapping.map(java.lang.String location)
func (this_ *SourceMap_PrefixLocationMapping) Map(location string) string {
	_ = location
	if javalang.StringStartsWith(j2gort.CheckNotNullString(location), j2gort.CheckNotNullString(this_.prefix)) {
		return (j2gort.StrValString(this_.replacement) + j2gort.StrValString(j2gort.JStrSubFrom(location, j2gort.JStrLen(this_.prefix))))
	}
	return j2gort.NullString
}

// Java method: PrefixLocationMapping.toString()
func (this_ *SourceMap_PrefixLocationMapping) String() string {
	return ("(" + j2gort.StrValString(this_.prefix) + "|" + j2gort.StrValString(this_.replacement) + ")")
}

// Java method: PrefixLocationMapping.equals(java.lang.Object other)
func (this_ *SourceMap_PrefixLocationMapping) Equals(arg_5 any) bool {
	other := arg_5
	_ = other
	var prefixLocationMapping *SourceMap_PrefixLocationMapping
	_ = prefixLocationMapping
	var ok_6 bool
	prefixLocationMapping, ok_6 = j2gort.TryCast[*SourceMap_PrefixLocationMapping](other)
	if ok_6 {
		return (func() bool {
			sr_7 := prefixLocationMapping.prefix
			se_8 := this_.prefix
			return j2gort.CheckNotNullString(sr_7) == se_8
		}() && func() bool {
			sr_9 := prefixLocationMapping.replacement
			se_10 := this_.replacement
			return j2gort.CheckNotNullString(sr_9) == se_10
		}())
	} else {
		return false
	}
}

// Java method: PrefixLocationMapping.hashCode()
func (this_ *SourceMap_PrefixLocationMapping) HashCode() int32 {
	{

		var oh_11 int32 = 1
		{
			hs_12 := this_.prefix
			oh_11 = 31 * oh_11
			if hs_12 != j2gort.NullString {
				oh_11 += javalang.StringHashCode(hs_12)
			}
		}
		{
			hs_13 := this_.replacement
			oh_11 = 31 * oh_11
			if hs_13 != j2gort.NullString {
				oh_11 += javalang.StringHashCode(hs_13)
			}
		}
		return oh_11
	}

}

// Java class: com.google.javascript.jscomp.SourceMap.Mapping
type SourceMap_Mapping struct {
	start   *pkg_sourcemap.FilePosition `gson:"start"`
	end     *pkg_sourcemap.FilePosition `gson:"end"`
	node    pkg_rhino.Node              `gson:"node"`
	idHash_ atomic.Int32
}

func (this_ *SourceMap_Mapping) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceMap_Mapping) IsJavaNil_() bool { return this_ == nil }
func (*SourceMap_Mapping) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SourceMap$Mapping", "Mapping"
}
func (this_ *SourceMap_Mapping) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SourceMap_Mapping) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SourceMap.Mapping", "com.google.javascript.jscomp.SourceMap$Mapping", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceMap_Mapping](v_)
	})
}

func (this_ *SourceMap_Mapping) initFields_SourceMap_Mapping() {
	_ = this_
}

// Java constructor: Mapping()
func NewSourceMap_Mapping() *SourceMap_Mapping {
	this_ := &SourceMap_Mapping{}
	InitSourceMap_Mapping(this_)
	return this_
}

func InitSourceMap_Mapping(this_ *SourceMap_Mapping) {
	this_.initFields_SourceMap_Mapping()
}

// Java method: Mapping.toString()
func (this_ *SourceMap_Mapping) String() string {
	return ("Mapping: start " + j2gort.StrVal(this_.start) + ", end " + j2gort.StrVal(this_.end) + ", node " + j2gort.StrVal(this_.node))
}
