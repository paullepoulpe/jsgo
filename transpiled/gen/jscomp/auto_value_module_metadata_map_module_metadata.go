// Port of com/google/javascript/jscomp/modules/AutoValue_ModuleMetadataMap_ModuleMetadata.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.modules.AutoValue_ModuleMetadataMap_ModuleMetadata
type AutoValue_ModuleMetadataMap_ModuleMetadata struct {
	ModuleMetadataMap_ModuleMetadataBase
	googNamespaces                    j2gort.Ext_com_google_common_collect_ImmutableMultiset `gson:"googNamespaces"`
	stronglyRequiredGoogNamespaces    j2gort.Ext_com_google_common_collect_ImmutableMultiset `gson:"stronglyRequiredGoogNamespaces"`
	dynamicallyRequiredGoogNamespaces j2gort.Ext_com_google_common_collect_ImmutableMultiset `gson:"dynamicallyRequiredGoogNamespaces"`
	maybeRequiredGoogNamespaces       j2gort.Ext_com_google_common_collect_ImmutableMultiset `gson:"maybeRequiredGoogNamespaces"`
	weaklyRequiredGoogNamespaces      j2gort.Ext_com_google_common_collect_ImmutableMultiset `gson:"weaklyRequiredGoogNamespaces"`
	es6ImportSpecifiers               j2gort.Ext_com_google_common_collect_ImmutableMultiset `gson:"es6ImportSpecifiers"`
	nestedModules                     j2gort.Ext_com_google_common_collect_ImmutableList     `gson:"nestedModules"`
	readToggles                       j2gort.Ext_com_google_common_collect_ImmutableMultiset `gson:"readToggles"`
	path                              *ModuleLoader_ModulePath                               `gson:"path"`
	rootNode                          pkg_rhino.Node                                         `gson:"rootNode"`
	moduleType                        ModuleMetadataMap_ModuleType                           `gson:"moduleType"`
	usesClosure                       bool                                                   `gson:"usesClosure"`
	isTestOnly                        bool                                                   `gson:"isTestOnly"`
}

func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata) BindSelf_(s any) {
	_ = s
	this_.ModuleMetadataMap_ModuleMetadataBase.BindSelf_(s)
}
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata) IsJavaNil_() bool { return this_ == nil }
func (*AutoValue_ModuleMetadataMap_ModuleMetadata) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.AutoValue_ModuleMetadataMap_ModuleMetadata", "AutoValue_ModuleMetadataMap_ModuleMetadata"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.AutoValue_ModuleMetadataMap_ModuleMetadata", "com.google.javascript.jscomp.modules.AutoValue_ModuleMetadataMap_ModuleMetadata", func(v_ any) bool {
		return j2gort.InstanceOf[*AutoValue_ModuleMetadataMap_ModuleMetadata](v_)
	})
}
func (*AutoValue_ModuleMetadataMap_ModuleMetadata) Implements_jscomp_ModuleMetadataMap_ModuleMetadata() {
}

func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata) initFields_AutoValue_ModuleMetadataMap_ModuleMetadata() {
	_ = this_
}

// Java constructor: AutoValue_ModuleMetadataMap_ModuleMetadata(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleType moduleType, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node rootNode, boolean usesClosure, boolean isTestOnly, com.google.common.collect.ImmutableMultiset<java.lang.String> googNamespaces, com.google.common.collect.ImmutableMultiset<java.lang.String> stronglyRequiredGoogNamespaces, com.google.common.collect.ImmutableMultiset<java.lang.String> dynamicallyRequiredGoogNamespaces, com.google.common.collect.ImmutableMultiset<java.lang.String> maybeRequiredGoogNamespaces, com.google.common.collect.ImmutableMultiset<java.lang.String> weaklyRequiredGoogNamespaces, com.google.common.collect.ImmutableMultiset<java.lang.String> es6ImportSpecifiers, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata> nestedModules, com.google.common.collect.ImmutableMultiset<java.lang.String> readToggles, com.google.javascript.jscomp.deps.ModuleLoader.@org.jspecify.annotations.Nullable ModulePath path)
func NewAutoValue_ModuleMetadataMap_ModuleMetadata(moduleType ModuleMetadataMap_ModuleType, rootNode pkg_rhino.Node, usesClosure bool, isTestOnly bool, googNamespaces j2gort.Ext_com_google_common_collect_ImmutableMultiset, stronglyRequiredGoogNamespaces j2gort.Ext_com_google_common_collect_ImmutableMultiset, dynamicallyRequiredGoogNamespaces j2gort.Ext_com_google_common_collect_ImmutableMultiset, maybeRequiredGoogNamespaces j2gort.Ext_com_google_common_collect_ImmutableMultiset, weaklyRequiredGoogNamespaces j2gort.Ext_com_google_common_collect_ImmutableMultiset, es6ImportSpecifiers j2gort.Ext_com_google_common_collect_ImmutableMultiset, nestedModules j2gort.Ext_com_google_common_collect_ImmutableList, readToggles j2gort.Ext_com_google_common_collect_ImmutableMultiset, path *ModuleLoader_ModulePath) *AutoValue_ModuleMetadataMap_ModuleMetadata {
	this_ := &AutoValue_ModuleMetadataMap_ModuleMetadata{}
	this_.BindSelf_(this_)
	InitAutoValue_ModuleMetadataMap_ModuleMetadata(this_, moduleType, rootNode, usesClosure, isTestOnly, googNamespaces, stronglyRequiredGoogNamespaces, dynamicallyRequiredGoogNamespaces, maybeRequiredGoogNamespaces, weaklyRequiredGoogNamespaces, es6ImportSpecifiers, nestedModules, readToggles, path)
	return this_
}

func InitAutoValue_ModuleMetadataMap_ModuleMetadata(this_ *AutoValue_ModuleMetadataMap_ModuleMetadata, moduleType ModuleMetadataMap_ModuleType, rootNode pkg_rhino.Node, usesClosure bool, isTestOnly bool, googNamespaces j2gort.Ext_com_google_common_collect_ImmutableMultiset, stronglyRequiredGoogNamespaces j2gort.Ext_com_google_common_collect_ImmutableMultiset, dynamicallyRequiredGoogNamespaces j2gort.Ext_com_google_common_collect_ImmutableMultiset, maybeRequiredGoogNamespaces j2gort.Ext_com_google_common_collect_ImmutableMultiset, weaklyRequiredGoogNamespaces j2gort.Ext_com_google_common_collect_ImmutableMultiset, es6ImportSpecifiers j2gort.Ext_com_google_common_collect_ImmutableMultiset, nestedModules j2gort.Ext_com_google_common_collect_ImmutableList, readToggles j2gort.Ext_com_google_common_collect_ImmutableMultiset, path *ModuleLoader_ModulePath) {
	_ = moduleType
	_ = rootNode
	_ = usesClosure
	_ = isTestOnly
	_ = googNamespaces
	_ = stronglyRequiredGoogNamespaces
	_ = dynamicallyRequiredGoogNamespaces
	_ = maybeRequiredGoogNamespaces
	_ = weaklyRequiredGoogNamespaces
	_ = es6ImportSpecifiers
	_ = nestedModules
	_ = readToggles
	_ = path
	InitModuleMetadataMap_ModuleMetadata(&this_.ModuleMetadataMap_ModuleMetadataBase)
	this_.initFields_AutoValue_ModuleMetadataMap_ModuleMetadata()
	this_.moduleType = moduleType
	this_.rootNode = rootNode
	this_.usesClosure = usesClosure
	this_.isTestOnly = isTestOnly
	this_.googNamespaces = googNamespaces
	this_.stronglyRequiredGoogNamespaces = stronglyRequiredGoogNamespaces
	this_.dynamicallyRequiredGoogNamespaces = dynamicallyRequiredGoogNamespaces
	this_.maybeRequiredGoogNamespaces = maybeRequiredGoogNamespaces
	this_.weaklyRequiredGoogNamespaces = weaklyRequiredGoogNamespaces
	this_.es6ImportSpecifiers = es6ImportSpecifiers
	this_.nestedModules = nestedModules
	this_.readToggles = readToggles
	this_.path = path
}

// Java method: AutoValue_ModuleMetadataMap_ModuleMetadata.moduleType()
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata) ModuleType() ModuleMetadataMap_ModuleType {
	return this_.moduleType
}

// Java method: AutoValue_ModuleMetadataMap_ModuleMetadata.rootNode()
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata) RootNode() pkg_rhino.Node {
	return this_.rootNode
}

// Java method: AutoValue_ModuleMetadataMap_ModuleMetadata.usesClosure()
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata) UsesClosure() bool {
	return this_.usesClosure
}

// Java method: AutoValue_ModuleMetadataMap_ModuleMetadata.isTestOnly()
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata) IsTestOnly() bool {
	return this_.isTestOnly
}

// Java method: AutoValue_ModuleMetadataMap_ModuleMetadata.googNamespaces()
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata) GoogNamespaces() j2gort.Ext_com_google_common_collect_ImmutableMultiset {
	return this_.googNamespaces
}

// Java method: AutoValue_ModuleMetadataMap_ModuleMetadata.stronglyRequiredGoogNamespaces()
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata) StronglyRequiredGoogNamespaces() j2gort.Ext_com_google_common_collect_ImmutableMultiset {
	return this_.stronglyRequiredGoogNamespaces
}

// Java method: AutoValue_ModuleMetadataMap_ModuleMetadata.dynamicallyRequiredGoogNamespaces()
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata) DynamicallyRequiredGoogNamespaces() j2gort.Ext_com_google_common_collect_ImmutableMultiset {
	return this_.dynamicallyRequiredGoogNamespaces
}

// Java method: AutoValue_ModuleMetadataMap_ModuleMetadata.maybeRequiredGoogNamespaces()
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata) MaybeRequiredGoogNamespaces() j2gort.Ext_com_google_common_collect_ImmutableMultiset {
	return this_.maybeRequiredGoogNamespaces
}

// Java method: AutoValue_ModuleMetadataMap_ModuleMetadata.weaklyRequiredGoogNamespaces()
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata) WeaklyRequiredGoogNamespaces() j2gort.Ext_com_google_common_collect_ImmutableMultiset {
	return this_.weaklyRequiredGoogNamespaces
}

// Java method: AutoValue_ModuleMetadataMap_ModuleMetadata.es6ImportSpecifiers()
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata) Es6ImportSpecifiers() j2gort.Ext_com_google_common_collect_ImmutableMultiset {
	return this_.es6ImportSpecifiers
}

// Java method: AutoValue_ModuleMetadataMap_ModuleMetadata.nestedModules()
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata) NestedModules() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.nestedModules
}

// Java method: AutoValue_ModuleMetadataMap_ModuleMetadata.readToggles()
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata) ReadToggles() j2gort.Ext_com_google_common_collect_ImmutableMultiset {
	return this_.readToggles
}

// Java method: AutoValue_ModuleMetadataMap_ModuleMetadata.path()
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata) Path() *ModuleLoader_ModulePath {
	return this_.path
}

// Java method: AutoValue_ModuleMetadataMap_ModuleMetadata.toString()
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata) String() string {
	return ("ModuleMetadata{moduleType=" + j2gort.StrVal(this_.moduleType) + ", rootNode=" + j2gort.StrVal(this_.rootNode) + ", usesClosure=" + j2gort.StrValBool(this_.usesClosure) + ", isTestOnly=" + j2gort.StrValBool(this_.isTestOnly) + ", googNamespaces=" + j2gort.StrVal(this_.googNamespaces) + ", stronglyRequiredGoogNamespaces=" + j2gort.StrVal(this_.stronglyRequiredGoogNamespaces) + ", dynamicallyRequiredGoogNamespaces=" + j2gort.StrVal(this_.dynamicallyRequiredGoogNamespaces) + ", maybeRequiredGoogNamespaces=" + j2gort.StrVal(this_.maybeRequiredGoogNamespaces) + ", weaklyRequiredGoogNamespaces=" + j2gort.StrVal(this_.weaklyRequiredGoogNamespaces) + ", es6ImportSpecifiers=" + j2gort.StrVal(this_.es6ImportSpecifiers) + ", nestedModules=" + j2gort.StrVal(this_.nestedModules) + ", readToggles=" + j2gort.StrVal(this_.readToggles) + ", path=" + j2gort.StrVal(this_.path) + "}")
}

// Java class: com.google.javascript.jscomp.modules.AutoValue_ModuleMetadataMap_ModuleMetadata.Builder
type AutoValue_ModuleMetadataMap_ModuleMetadata_Builder struct {
	ModuleMetadataMap_ModuleMetadata_BuilderBase
	googNamespacesBuilder_d                    *j2gort.Ext_com_google_common_collect_ImmutableMultiset_Builder `gson:"googNamespacesBuilder$"`
	googNamespaces                             j2gort.Ext_com_google_common_collect_ImmutableMultiset          `gson:"googNamespaces"`
	stronglyRequiredGoogNamespacesBuilder_d    *j2gort.Ext_com_google_common_collect_ImmutableMultiset_Builder `gson:"stronglyRequiredGoogNamespacesBuilder$"`
	stronglyRequiredGoogNamespaces             j2gort.Ext_com_google_common_collect_ImmutableMultiset          `gson:"stronglyRequiredGoogNamespaces"`
	dynamicallyRequiredGoogNamespacesBuilder_d *j2gort.Ext_com_google_common_collect_ImmutableMultiset_Builder `gson:"dynamicallyRequiredGoogNamespacesBuilder$"`
	dynamicallyRequiredGoogNamespaces          j2gort.Ext_com_google_common_collect_ImmutableMultiset          `gson:"dynamicallyRequiredGoogNamespaces"`
	maybeRequiredGoogNamespacesBuilder_d       *j2gort.Ext_com_google_common_collect_ImmutableMultiset_Builder `gson:"maybeRequiredGoogNamespacesBuilder$"`
	maybeRequiredGoogNamespaces                j2gort.Ext_com_google_common_collect_ImmutableMultiset          `gson:"maybeRequiredGoogNamespaces"`
	weaklyRequiredGoogNamespacesBuilder_d      *j2gort.Ext_com_google_common_collect_ImmutableMultiset_Builder `gson:"weaklyRequiredGoogNamespacesBuilder$"`
	weaklyRequiredGoogNamespaces               j2gort.Ext_com_google_common_collect_ImmutableMultiset          `gson:"weaklyRequiredGoogNamespaces"`
	es6ImportSpecifiersBuilder_d               *j2gort.Ext_com_google_common_collect_ImmutableMultiset_Builder `gson:"es6ImportSpecifiersBuilder$"`
	es6ImportSpecifiers                        j2gort.Ext_com_google_common_collect_ImmutableMultiset          `gson:"es6ImportSpecifiers"`
	nestedModulesBuilder_d                     *j2gort.Ext_com_google_common_collect_ImmutableList_Builder     `gson:"nestedModulesBuilder$"`
	nestedModules                              j2gort.Ext_com_google_common_collect_ImmutableList              `gson:"nestedModules"`
	readTogglesBuilder_d                       *j2gort.Ext_com_google_common_collect_ImmutableMultiset_Builder `gson:"readTogglesBuilder$"`
	readToggles                                j2gort.Ext_com_google_common_collect_ImmutableMultiset          `gson:"readToggles"`
	path                                       *ModuleLoader_ModulePath                                        `gson:"path"`
	rootNode                                   pkg_rhino.Node                                                  `gson:"rootNode"`
	moduleType                                 ModuleMetadataMap_ModuleType                                    `gson:"moduleType"`
	usesClosure                                bool                                                            `gson:"usesClosure"`
	isTestOnly                                 bool                                                            `gson:"isTestOnly"`
	set_d0                                     int8                                                            `gson:"set$0"`
}

func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata_Builder) BindSelf_(s any) {
	_ = s
	this_.ModuleMetadataMap_ModuleMetadata_BuilderBase.BindSelf_(s)
}
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata_Builder) IsJavaNil_() bool {
	return this_ == nil
}
func (*AutoValue_ModuleMetadataMap_ModuleMetadata_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.AutoValue_ModuleMetadataMap_ModuleMetadata$Builder", "Builder"
}
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ModuleMetadataMap_ModuleMetadata_BuilderBase.IdentityHash_()
}
func (*AutoValue_ModuleMetadataMap_ModuleMetadata_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.AutoValue_ModuleMetadataMap_ModuleMetadata.Builder", "com.google.javascript.jscomp.modules.AutoValue_ModuleMetadataMap_ModuleMetadata$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*AutoValue_ModuleMetadataMap_ModuleMetadata_Builder](v_)
	})
}
func (*AutoValue_ModuleMetadataMap_ModuleMetadata_Builder) Implements_jscomp_ModuleMetadataMap_ModuleMetadata_Builder() {
}

func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata_Builder) initFields_AutoValue_ModuleMetadataMap_ModuleMetadata_Builder() {
	_ = this_
}

// Java constructor: Builder()
func NewAutoValue_ModuleMetadataMap_ModuleMetadata_Builder() *AutoValue_ModuleMetadataMap_ModuleMetadata_Builder {
	this_ := &AutoValue_ModuleMetadataMap_ModuleMetadata_Builder{}
	this_.BindSelf_(this_)
	InitAutoValue_ModuleMetadataMap_ModuleMetadata_Builder(this_)
	return this_
}

func InitAutoValue_ModuleMetadataMap_ModuleMetadata_Builder(this_ *AutoValue_ModuleMetadataMap_ModuleMetadata_Builder) {
	InitModuleMetadataMap_ModuleMetadata_Builder(&this_.ModuleMetadataMap_ModuleMetadata_BuilderBase)
	this_.initFields_AutoValue_ModuleMetadataMap_ModuleMetadata_Builder()
}

// Java method: Builder.moduleType(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleType moduleType)
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata_Builder) ModuleType_1(moduleType ModuleMetadataMap_ModuleType) ModuleMetadataMap_ModuleMetadata_Builder {
	_ = moduleType
	if moduleType == 0 {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null moduleType")))
	}
	this_.moduleType = moduleType
	return this_
}

// Java method: Builder.moduleType()
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata_Builder) ModuleType_0() ModuleMetadataMap_ModuleType {
	if this_.moduleType == 0 {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Property \"moduleType\" has not been set")))
	}
	return this_.moduleType
}

// Java method: Builder.rootNode(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node rootNode)
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata_Builder) RootNode(rootNode pkg_rhino.Node) ModuleMetadataMap_ModuleMetadata_Builder {
	_ = rootNode
	this_.rootNode = rootNode
	return this_
}

// Java method: Builder.usesClosure(boolean usesClosure)
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata_Builder) UsesClosure(usesClosure bool) ModuleMetadataMap_ModuleMetadata_Builder {
	_ = usesClosure
	this_.usesClosure = usesClosure
	this_.set_d0 = int8((int32(this_.set_d0) | int32(int8(1))))
	return this_
}

// Java method: Builder.isTestOnly(boolean isTestOnly)
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata_Builder) IsTestOnly(isTestOnly bool) ModuleMetadataMap_ModuleMetadata_Builder {
	_ = isTestOnly
	this_.isTestOnly = isTestOnly
	this_.set_d0 = int8((int32(this_.set_d0) | int32(int8(2))))
	return this_
}

// Java method: Builder.googNamespacesBuilder()
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata_Builder) GoogNamespacesBuilder() *j2gort.Ext_com_google_common_collect_ImmutableMultiset_Builder {
	if this_.googNamespacesBuilder_d == nil {
		this_.googNamespacesBuilder_d = j2gort.Ext_com_google_common_collect_ImmutableMultiset_S_Builder()
	}
	return this_.googNamespacesBuilder_d
}

// Java method: Builder.stronglyRequiredGoogNamespacesBuilder()
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata_Builder) StronglyRequiredGoogNamespacesBuilder() *j2gort.Ext_com_google_common_collect_ImmutableMultiset_Builder {
	if this_.stronglyRequiredGoogNamespacesBuilder_d == nil {
		this_.stronglyRequiredGoogNamespacesBuilder_d = j2gort.Ext_com_google_common_collect_ImmutableMultiset_S_Builder()
	}
	return this_.stronglyRequiredGoogNamespacesBuilder_d
}

// Java method: Builder.dynamicallyRequiredGoogNamespacesBuilder()
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata_Builder) DynamicallyRequiredGoogNamespacesBuilder() *j2gort.Ext_com_google_common_collect_ImmutableMultiset_Builder {
	if this_.dynamicallyRequiredGoogNamespacesBuilder_d == nil {
		this_.dynamicallyRequiredGoogNamespacesBuilder_d = j2gort.Ext_com_google_common_collect_ImmutableMultiset_S_Builder()
	}
	return this_.dynamicallyRequiredGoogNamespacesBuilder_d
}

// Java method: Builder.maybeRequiredGoogNamespacesBuilder()
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata_Builder) MaybeRequiredGoogNamespacesBuilder() *j2gort.Ext_com_google_common_collect_ImmutableMultiset_Builder {
	if this_.maybeRequiredGoogNamespacesBuilder_d == nil {
		this_.maybeRequiredGoogNamespacesBuilder_d = j2gort.Ext_com_google_common_collect_ImmutableMultiset_S_Builder()
	}
	return this_.maybeRequiredGoogNamespacesBuilder_d
}

// Java method: Builder.weaklyRequiredGoogNamespacesBuilder()
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata_Builder) WeaklyRequiredGoogNamespacesBuilder() *j2gort.Ext_com_google_common_collect_ImmutableMultiset_Builder {
	if this_.weaklyRequiredGoogNamespacesBuilder_d == nil {
		this_.weaklyRequiredGoogNamespacesBuilder_d = j2gort.Ext_com_google_common_collect_ImmutableMultiset_S_Builder()
	}
	return this_.weaklyRequiredGoogNamespacesBuilder_d
}

// Java method: Builder.es6ImportSpecifiersBuilder()
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata_Builder) Es6ImportSpecifiersBuilder() *j2gort.Ext_com_google_common_collect_ImmutableMultiset_Builder {
	if this_.es6ImportSpecifiersBuilder_d == nil {
		this_.es6ImportSpecifiersBuilder_d = j2gort.Ext_com_google_common_collect_ImmutableMultiset_S_Builder()
	}
	return this_.es6ImportSpecifiersBuilder_d
}

// Java method: Builder.nestedModulesBuilder()
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata_Builder) NestedModulesBuilder() *j2gort.Ext_com_google_common_collect_ImmutableList_Builder {
	if this_.nestedModulesBuilder_d == nil {
		this_.nestedModulesBuilder_d = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
	}
	return this_.nestedModulesBuilder_d
}

// Java method: Builder.readTogglesBuilder()
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata_Builder) ReadTogglesBuilder() *j2gort.Ext_com_google_common_collect_ImmutableMultiset_Builder {
	if this_.readTogglesBuilder_d == nil {
		this_.readTogglesBuilder_d = j2gort.Ext_com_google_common_collect_ImmutableMultiset_S_Builder()
	}
	return this_.readTogglesBuilder_d
}

// Java method: Builder.path(com.google.javascript.jscomp.deps.ModuleLoader.@org.jspecify.annotations.Nullable ModulePath path)
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata_Builder) Path(path *ModuleLoader_ModulePath) ModuleMetadataMap_ModuleMetadata_Builder {
	_ = path
	this_.path = path
	return this_
}

// Java method: Builder.build()
func (this_ *AutoValue_ModuleMetadataMap_ModuleMetadata_Builder) Build() ModuleMetadataMap_ModuleMetadata {
	if this_.googNamespacesBuilder_d != nil {
		this_.googNamespaces = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultiset](this_.googNamespacesBuilder_d.Build())
	} else if this_.googNamespaces == nil {
		this_.googNamespaces = j2gort.Ext_com_google_common_collect_ImmutableMultiset_Of_0()
	}
	if this_.stronglyRequiredGoogNamespacesBuilder_d != nil {
		this_.stronglyRequiredGoogNamespaces = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultiset](this_.stronglyRequiredGoogNamespacesBuilder_d.Build())
	} else if this_.stronglyRequiredGoogNamespaces == nil {
		this_.stronglyRequiredGoogNamespaces = j2gort.Ext_com_google_common_collect_ImmutableMultiset_Of_0()
	}
	if this_.dynamicallyRequiredGoogNamespacesBuilder_d != nil {
		this_.dynamicallyRequiredGoogNamespaces = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultiset](this_.dynamicallyRequiredGoogNamespacesBuilder_d.Build())
	} else if this_.dynamicallyRequiredGoogNamespaces == nil {
		this_.dynamicallyRequiredGoogNamespaces = j2gort.Ext_com_google_common_collect_ImmutableMultiset_Of_0()
	}
	if this_.maybeRequiredGoogNamespacesBuilder_d != nil {
		this_.maybeRequiredGoogNamespaces = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultiset](this_.maybeRequiredGoogNamespacesBuilder_d.Build())
	} else if this_.maybeRequiredGoogNamespaces == nil {
		this_.maybeRequiredGoogNamespaces = j2gort.Ext_com_google_common_collect_ImmutableMultiset_Of_0()
	}
	if this_.weaklyRequiredGoogNamespacesBuilder_d != nil {
		this_.weaklyRequiredGoogNamespaces = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultiset](this_.weaklyRequiredGoogNamespacesBuilder_d.Build())
	} else if this_.weaklyRequiredGoogNamespaces == nil {
		this_.weaklyRequiredGoogNamespaces = j2gort.Ext_com_google_common_collect_ImmutableMultiset_Of_0()
	}
	if this_.es6ImportSpecifiersBuilder_d != nil {
		this_.es6ImportSpecifiers = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultiset](this_.es6ImportSpecifiersBuilder_d.Build())
	} else if this_.es6ImportSpecifiers == nil {
		this_.es6ImportSpecifiers = j2gort.Ext_com_google_common_collect_ImmutableMultiset_Of_0()
	}
	if this_.nestedModulesBuilder_d != nil {
		this_.nestedModules = this_.nestedModulesBuilder_d.BuildImmList_()
	} else if this_.nestedModules == nil {
		this_.nestedModules = j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	}
	if this_.readTogglesBuilder_d != nil {
		this_.readToggles = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultiset](this_.readTogglesBuilder_d.Build())
	} else if this_.readToggles == nil {
		this_.readToggles = j2gort.Ext_com_google_common_collect_ImmutableMultiset_Of_0()
	}
	if (int32(this_.set_d0) != int32(3)) || (this_.moduleType == 0) {
		var missing *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
		_ = missing
		if this_.moduleType == 0 {
			missing.Append_1_string(" moduleType")
		}
		if (int32(this_.set_d0) & int32(1)) == int32(0) {
			missing.Append_1_string(" usesClosure")
		}
		if (int32(this_.set_d0) & int32(2)) == int32(0) {
			missing.Append_1_string(" isTestOnly")
		}
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Missing required properties:"+j2gort.StrVal(missing))))
	}
	return NewAutoValue_ModuleMetadataMap_ModuleMetadata(this_.moduleType, this_.rootNode, this_.usesClosure, this_.isTestOnly, this_.googNamespaces, this_.stronglyRequiredGoogNamespaces, this_.dynamicallyRequiredGoogNamespaces, this_.maybeRequiredGoogNamespaces, this_.weaklyRequiredGoogNamespaces, this_.es6ImportSpecifiers, this_.nestedModules, this_.readToggles, this_.path)
}
