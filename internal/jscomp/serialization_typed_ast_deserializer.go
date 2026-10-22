// Port of com/google/javascript/jscomp/serialization/TypedAstDeserializer.java (with the
// AutoValue-generated AutoValue_TypedAstDeserializer_DeserializedAst.java) @ closure-compiler v20261015

/*
 * Copyright 2021 The Closure Compiler Authors.
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
	"bytes"
	"io"
	"math"
	"runtime"
	"strconv"
	"sync"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type runtimeLibsTypedAstStream struct {
	data []byte
}

func (r *runtimeLibsTypedAstStream) Read(p []byte) (int, error) {
	return 0, io.EOF
}

type cachedTypedAstShard struct {
	proto       *TypedAst
	stringShard *StringPool
}

var (
	runtimeLibsTypedAstOnce   sync.Once
	runtimeLibsTypedAstShards []cachedTypedAstShard
)

// TypedAstDeserializer deserializes a list of TypedAst protos into the
// JSCompiler AST structure.
//
// Go-only: language-adapter — Guava/JDK Optional values are nil-able
// pointers (nil is absent); java.util.function.Supplier<Node> is
// func() *rhino.Node; the ConcurrentMap filesystem is a
// *jutil.LinkedHashMap (Java's ConcurrentHashMap iteration order is
// unspecified, and Closure only looks suppliers up by file).
type TypedAstDeserializer struct {
	mode             typedAstDeserializer_Mode
	syntheticExterns *SourceFile
	colorPoolBuilder *ColorPool_Builder // Optional: nil when absent

	requiredInputFiles            *jutil.LinkedHashSet[*SourceFile] // Optional: nil when absent
	filePoolBuilder               *jutil.LinkedHashMap[string, *SourceFile]
	canonicalFilePool             *jutil.LinkedHashMap[string, *SourceFile]
	typedAstFilesystem            *jutil.LinkedHashMap[*SourceFile, func() *rhino.Node]
	externProperties              *jutil.LinkedHashSet[string]
	runtimeLibraries              *jutil.LinkedHashSet[string]
	syntheticExternsDeserializers *jutil.ArrayList[*ScriptNodeDeserializer]
}

// Java: TypedAstDeserializer(SourceFile, Optional<ColorPool.Builder>, Optional<ImmutableSet<SourceFile>>, Mode, boolean)
func newTypedAstDeserializer(
	syntheticExterns *SourceFile,
	existingColorPool *ColorPool_Builder,
	requiredInputFiles *jutil.LinkedHashSet[*SourceFile],
	mode typedAstDeserializer_Mode,
	includeTypeInformation bool) *TypedAstDeserializer {
	d := &TypedAstDeserializer{
		filePoolBuilder:               jutil.NewLinkedHashMap[string, *SourceFile](),
		canonicalFilePool:             jutil.NewLinkedHashMap[string, *SourceFile](),
		typedAstFilesystem:            jutil.NewLinkedHashMap[*SourceFile, func() *rhino.Node](),
		externProperties:              jutil.NewLinkedHashSet[string](),
		runtimeLibraries:              jutil.NewLinkedHashSet[string](),
		syntheticExternsDeserializers: jutil.NewArrayList[*ScriptNodeDeserializer](),
	}
	d.syntheticExterns = syntheticExterns
	d.mode = mode
	if includeTypeInformation {
		// Optional.of(existingColorPool.or(ColorPool.builder()))
		if existingColorPool != nil {
			d.colorPoolBuilder = existingColorPool
		} else {
			d.colorPoolBuilder = ColorPool_Builder_()
		}
	} else {
		d.colorPoolBuilder = nil
	}
	d.requiredInputFiles = requiredInputFiles
	return d
}

// Java: private enum TypedAstDeserializer.Mode
type typedAstDeserializer_Mode int32

const (
	typedAstDeserializer_Mode_RUNTIME_LIBRARY_ONLY typedAstDeserializer_Mode = iota
	typedAstDeserializer_Mode_FULL_AST
)

// TypedAstDeserializer_DeserializeFullAst transforms a given TypedAst
// delimited stream into a compiler AST
//
// requiredInputFiles are all SourceFiles that should go into the
// typedAstFilesystem. If a TypedAst contains a file not in this set, that
// file will not be added to the typedAstFilesystem to avoid wasting memory.
// includeTypeInformation is whether to deserialize the "Typed" half of a
// "TypedAst". If false ignores the TypePool, any TypePointers on AstNodes,
// and does not create a ColorRegistry. resolveSourceMapAnnotations is
// whether to create and register a SourceMapInput for a given
// `//# sourceMappingURL` (stored in LazyAst) when deserializing.
// parseInlineSourceMaps is whether a given `//# sourceMappingURL` should be
// registered as a Base64 encoded source map.
//
// Java: TypedAstDeserializer.deserializeFullAst(AbstractCompiler, SourceFile, ImmutableSet<SourceFile>, InputStream, boolean, boolean, boolean)
func TypedAstDeserializer_DeserializeFullAst(
	compiler AbstractCompiler,
	syntheticExterns *SourceFile,
	requiredInputFiles *jutil.LinkedHashSet[*SourceFile],
	typedAstsStream io.Reader,
	includeTypeInformation bool,
	resolveSourceMapAnnotations bool,
	parseInlineSourceMaps bool) *TypedAstDeserializer_DeserializedAst {
	// requiredInputFiles.stream().collect(toImmutableMap(SourceFile::getName, Function.identity()))
	sourceFilesByNameBuilder := jutil.NewImmutableMapBuilder[string, *SourceFile]()
	for f := range requiredInputFiles.All() {
		sourceFilesByNameBuilder.Put(f.GetName(), f)
	}
	sourceFilesByName := sourceFilesByNameBuilder.BuildOrThrow()
	return typedAstDeserializer_deserialize(
		compiler,
		syntheticExterns,
		requiredInputFiles,
		sourceFilesByName,
		nil,
		typedAstsStream,
		typedAstDeserializer_Mode_FULL_AST,
		includeTypeInformation,
		resolveSourceMapAnnotations,
		parseInlineSourceMaps)
}

// TypedAstDeserializer_DeserializeRuntimeLibraries transforms the special
// runtime library TypedAst
//
// colorPool is a ColorPool.Builder holding the colors on the full AST (nil
// when absent). We want to merge these colors with the runtime library
// colors to allow injecting runtime libraries without re-typechecking them.
//
// Java: TypedAstDeserializer.deserializeRuntimeLibraries(AbstractCompiler, SourceFile, Optional<ColorPool.Builder>, InputStream, boolean, boolean)
func TypedAstDeserializer_DeserializeRuntimeLibraries(
	compiler AbstractCompiler,
	syntheticExterns *SourceFile,
	colorPool *ColorPool_Builder,
	typedAstsStream io.Reader,
	resolveSourceMapAnnotations bool,
	parseInlineSourceMaps bool) *TypedAstDeserializer_DeserializedAst {
	return typedAstDeserializer_deserialize(
		compiler,
		syntheticExterns,
		// we don't a priori know the SourceFiles corresponding to the runtime libraries like we
		// do for normal CompilerInputs
		nil,
		jutil.NewLinkedHashMap[string, *SourceFile](),
		colorPool,
		typedAstsStream,
		typedAstDeserializer_Mode_RUNTIME_LIBRARY_ONLY,
		colorPool != nil,
		resolveSourceMapAnnotations,
		parseInlineSourceMaps)
}

// Java: TypedAstDeserializer.deserialize(AbstractCompiler, SourceFile, Optional<ImmutableSet<SourceFile>>, ImmutableMap<String, SourceFile>, Optional<ColorPool.Builder>, InputStream, Mode, boolean, boolean, boolean)
func typedAstDeserializer_deserialize(
	compiler AbstractCompiler,
	syntheticExterns *SourceFile,
	requiredInputFiles *jutil.LinkedHashSet[*SourceFile],
	scriptSourceFiles *jutil.LinkedHashMap[string, *SourceFile],
	colorPool *ColorPool_Builder,
	typedAstStream io.Reader,
	mode typedAstDeserializer_Mode,
	includeTypeInformation bool,
	resolveSourceMapAnnotations bool,
	parseInlineSourceMaps bool) *TypedAstDeserializer_DeserializedAst {
	jutil.CheckArgumentMsg(
		(colorPool != nil) == (mode == typedAstDeserializer_Mode_RUNTIME_LIBRARY_ONLY && includeTypeInformation),
		"ColorPool.Builder required iff deserializing runtime libraries & including types")

	deserializer :=
		newTypedAstDeserializer(
			syntheticExterns, colorPool, requiredInputFiles, mode, includeTypeInformation)
	deserializer.addToFilePool(syntheticExterns.GetName(), syntheticExterns)
	for name, file := range scriptSourceFiles.All() {
		deserializer.addToFilePool(name, file)
	}

	if mode != typedAstDeserializer_Mode_RUNTIME_LIBRARY_ONLY {
		// skip this step if deserializing the runtime libraries to avoid an infinite loop
		// runtime library initialization needs an unbuilt ColorPool.Builder
		compiler.InitRuntimeLibraryTypedAsts(deserializer.colorPoolBuilder)
	}

	typedAstDeserializer_deserializeTypedAsts(
		typedAstStream, deserializer, compiler, resolveSourceMapAnnotations, parseInlineSourceMaps)

	deserializer.typedAstFilesystem.Put(
		syntheticExterns,
		newTypedAstDeserializer_SyntheticExternsDeserializer(
			syntheticExterns, deserializer.syntheticExternsDeserializers).
			Deserialize)

	return deserializer.toDeserializedAst()
}

// typedAstDeserializer_SyntheticExternsDeserializer is a helper class that
// is a lazy deserializer for the synthetic externs script
//
// The "checks" phase of JSCompiler may synthesize new extern nodes. This
// means that if we are merging multiple TypedAST shards from different
// independent "checks" phases, each shard may have its own synthetic
// externs script. This class concatenates each synthetic externs script
// into one.
//
// Note: while it's possible that other files appear in multiple TypedAST
// shards, we assume that those files are identical, and arbitrarily choose
// one copy of each file. The synthetic externs are the only case where we
// concatenate multiple versions of the same file rather than just choosing
// one.
//
// Java: private static final class TypedAstDeserializer.SyntheticExternsDeserializer
type typedAstDeserializer_SyntheticExternsDeserializer struct {
	syntheticExterns              *SourceFile
	syntheticExternsDeserializers *jutil.ArrayList[*ScriptNodeDeserializer]
}

// Java: TypedAstDeserializer.SyntheticExternsDeserializer(SourceFile, ArrayList<ScriptNodeDeserializer>)
func newTypedAstDeserializer_SyntheticExternsDeserializer(
	syntheticExterns *SourceFile,
	syntheticExternsDeserializers *jutil.ArrayList[*ScriptNodeDeserializer]) *typedAstDeserializer_SyntheticExternsDeserializer {
	return &typedAstDeserializer_SyntheticExternsDeserializer{
		syntheticExterns:              syntheticExterns,
		syntheticExternsDeserializers: syntheticExternsDeserializers,
	}
}

// Java: TypedAstDeserializer.SyntheticExternsDeserializer.deserialize()
func (s *typedAstDeserializer_SyntheticExternsDeserializer) Deserialize() *rhino.Node {
	script := rhino.IR_Script()
	script.SetStaticSourceFile(s.syntheticExterns)
	for d := range s.syntheticExternsDeserializers.All() {
		script.AddChildrenToBack(d.DeserializeNew().RemoveChildren())
	}
	return script
}

// Java: TypedAstDeserializer.toDeserializedAst()
func (d *TypedAstDeserializer) toDeserializedAst() *TypedAstDeserializer_DeserializedAst {
	var registry *colors.ColorRegistry // Optional: nil when absent
	if d.mode == typedAstDeserializer_Mode_RUNTIME_LIBRARY_ONLY || d.colorPoolBuilder == nil {
		registry = nil
	} else {
		registry = d.colorPoolBuilder.Build().GetRegistry()
	}
	return typedAstDeserializer_DeserializedAst_create(
		d.typedAstFilesystem, registry,
		jutil.NewLinkedHashSet(d.externProperties.ToSlice()...),
		jutil.NewLinkedHashSet(d.runtimeLibraries.ToSlice()...))
}

// typedAstDeserializer_lazyAsts is Iterables.concat(typedAstProto.getExternAstList(),
// typedAstProto.getCodeAstList()).
func typedAstDeserializer_lazyAsts(typedAstProto *TypedAst) []*LazyAst {
	return append(append([]*LazyAst(nil), typedAstProto.GetExternAstList()...), typedAstProto.GetCodeAstList()...)
}

// Java: TypedAstDeserializer.deserializeTypedAst(TypedAst, AbstractCompiler, boolean, boolean)
func (d *TypedAstDeserializer) deserializeTypedAst(
	typedAstProto *TypedAst,
	compiler AbstractCompiler,
	resolveSourceMapAnnotations bool,
	parseInlineSourceMaps bool) {
	d.deserializeTypedAstWithStringShard(
		typedAstProto, nil, compiler, resolveSourceMapAnnotations, parseInlineSourceMaps)
}

func (d *TypedAstDeserializer) deserializeTypedAstWithStringShard(
	typedAstProto *TypedAst,
	stringShard *StringPool,
	compiler AbstractCompiler,
	resolveSourceMapAnnotations bool,
	parseInlineSourceMaps bool) {
	fileShard := d.toFileShard(typedAstProto)

	if d.mode == typedAstDeserializer_Mode_FULL_AST {
		// Exclude any TypedAST shards passed to the compiler that don't contain any actual
		// sources we care about. For example: if we're passed "--js=a.js --js=b.js", and see a
		// TypedAST only containing ["c.js", "synthetic:externs"], it's a drain on performance to
		// include it.
		containsRequiredInput := false
		for _, lazyAst := range typedAstDeserializer_lazyAsts(typedAstProto) {
			file := fileShard[lazyAst.GetSourceFile()-1]
			if jutil.CheckNotNull(d.requiredInputFiles).Contains(file) {
				containsRequiredInput = true
				break
			}
		}
		if !containsRequiredInput {
			return
		}
	}

	runtimeLibraryList := typedAstProto.GetRuntimeLibraryToInjectList()
	for i := 0; i < len(runtimeLibraryList); i++ {
		d.runtimeLibraries.Add(runtimeLibraryList[i])
	}

	// TODO(b/248351234): can we avoid some of this work if the shard only contains weak srcs?
	// one risk: could checks passes synthesize new externProperties even for weak srcs?
	if stringShard == nil {
		stringShard = StringPool_FromProto(typedAstProto.GetStringPool())
	}
	var colorShard *ColorPool_ShardView // Optional: nil when absent
	if d.colorPoolBuilder != nil {
		colorShard = d.colorPoolBuilder.AddShard(typedAstProto.GetTypePool(), stringShard)
	}
	for _, x := range typedAstProto.GetExternsSummary().GetPropNamePtrList() {
		d.externProperties.Add(stringShard.Get(int(x)))
	}

	for _, lazyAst := range typedAstDeserializer_lazyAsts(typedAstProto) {
		d.initLazyAstDeserializer(
			lazyAst,
			fileShard,
			colorShard,
			stringShard,
			compiler,
			resolveSourceMapAnnotations,
			parseInlineSourceMaps)
	}
}

// Java: TypedAstDeserializer.addToFilePool(String, SourceFile)
func (d *TypedAstDeserializer) addToFilePool(filename string, file *SourceFile) {
	d.filePoolBuilder.Put(filename, file)
	canonicalPath := SourceFile_CanonicalizePathForMatching(filename)
	if canonicalPath != filename {
		d.canonicalFilePool.PutIfAbsent(canonicalPath, file)
	}
}

// Java: TypedAstDeserializer.findExistingFile(String)
func (d *TypedAstDeserializer) findExistingFile(filename string) *SourceFile {
	exact := d.filePoolBuilder.Get(filename)
	if exact != nil {
		return exact
	}
	canonicalPath := SourceFile_CanonicalizePathForMatching(filename)
	canonical := d.canonicalFilePool.Get(canonicalPath)
	if canonical != nil {
		return canonical
	}
	return d.filePoolBuilder.Get(canonicalPath)
}

// Java: TypedAstDeserializer.toFileShard(TypedAst)
func (d *TypedAstDeserializer) toFileShard(typedAstProto *TypedAst) []*SourceFile {
	var fileShardBuilder []*SourceFile
	protos := typedAstProto.GetSourceFilePool().GetSourceFileList()
	for i := 0; i < len(protos); i++ {
		p := protos[i]

		existingFile := d.findExistingFile(p.GetFilename())
		if existingFile != nil {
			// Merge the existing SourceFile with the information serialized in the TypedAST.
			existingFile.RestoreCachedStateFrom(p)
			fileShardBuilder = append(fileShardBuilder, existingFile)
		} else {
			protoFile := SourceFile_FromProto(p)
			d.addToFilePool(p.GetFilename(), protoFile)
			fileShardBuilder = append(fileShardBuilder, protoFile)
		}
	}
	return fileShardBuilder
}

// Java: TypedAstDeserializer.initLazyAstDeserializer(LazyAst, ImmutableList<SourceFile>, Optional<ColorPool.ShardView>, StringPool, AbstractCompiler, boolean, boolean)
func (d *TypedAstDeserializer) initLazyAstDeserializer(
	lazyAst *LazyAst,
	fileShard []*SourceFile,
	colorShard *ColorPool_ShardView,
	stringShard *StringPool,
	compiler AbstractCompiler,
	resolveSourceMapAnnotations bool,
	parseInlineSourceMaps bool) {
	file := fileShard[lazyAst.GetSourceFile()-1]

	if d.requiredInputFiles != nil &&
		!d.requiredInputFiles.Contains(file) &&
		// Synthetic externs bypass the normal dependency system and are always included in a
		// compilation.
		d.syntheticExterns != file {
		return
	}

	if file.IsWeak() {
		d.typedAstFilesystem.ComputeIfAbsent(file, typedAstDeserializer_createStubWeakScriptNode)
		return
	}

	deserializer :=
		NewScriptNodeDeserializer(lazyAst, stringShard, colorShard, fileShard)

	if d.syntheticExterns == file {
		d.syntheticExternsDeserializers.Add(deserializer)
	} else {
		d.typedAstFilesystem.ComputeIfAbsent(file, func(f *SourceFile) func() *rhino.Node { return deserializer.DeserializeNew })
	}

	typedAstDeserializer_addInputSourceMap(deserializer, compiler, resolveSourceMapAnnotations, parseInlineSourceMaps)
}

// typedAstDeserializer_addInputSourceMap: while deserializing a TypedAST
// script, if it has a corresponding sourceMappingURL, create a
// SourceMapInput and register it.
//
// resolveSourceMapAnnotations is whether to create and register a
// SourceMapInput for a given `//# sourceMappingURL` (stored in LazyAst)
// when deserializing. parseInlineSourceMaps is whether a given
// `//# sourceMappingURL` should be registered as a Base64 encoded source
// map.
//
// Java: TypedAstDeserializer.addInputSourceMap(ScriptNodeDeserializer, AbstractCompiler, boolean, boolean)
func typedAstDeserializer_addInputSourceMap(
	deserializer *ScriptNodeDeserializer,
	compiler AbstractCompiler,
	resolveSourceMapAnnotations bool,
	parseInlineSourceMaps bool) {
	sourceFile := deserializer.GetSourceFile()
	sourceMappingURL := deserializer.GetSourceMappingURL() // This is the encoded source map.

	if jutil.StringLength(sourceMappingURL) > 0 && resolveSourceMapAnnotations {
		// base64EncodedSourceMap adds "data:application/json;base64," prefix to the sourceMappingURL
		base64EncodedSourceMap :=
			SourceMapResolver_AddBase64PrefixToEncodedSourceMap(sourceMappingURL)
		sourceMapSourceFile :=
			SourceMapResolver_ExtractSourceMap(
				sourceFile, base64EncodedSourceMap, parseInlineSourceMaps)
		if sourceMapSourceFile != nil {
			compiler.AddInputSourceMap(sourceFile.GetName(), NewSourceMapInput(sourceMapSourceFile))
		}
	}
}

// typedAstDeserializer_createStubWeakScriptNode adds a stub deserializer for
// weak files to the typed ast files system.
//
// Weak files don't actually need to be deserialized. They're only needed for
// typechecking and by definition a TypedAST has already been typechecked.
//
// Java: TypedAstDeserializer.createStubWeakScriptNode(SourceFile)
func typedAstDeserializer_createStubWeakScriptNode(file *SourceFile) func() *rhino.Node {
	return func() *rhino.Node {
		stubScript := rhino.IR_Script().SetStaticSourceFile(file)
		stubScript.PutProp(rhino.Node_FEATURE_SET, parser.FeatureSet_BARE_MINIMUM)
		return stubScript
	}
}

func typedAstDeserializer_decodeShardBytes(messageBytes []byte, buildStringPool func(*TypedAst) bool) (*TypedAst, *StringPool) {
	rd := &protoReader{b: messageBytes, depth: 1, limit: math.MaxInt32}
	m := &TypedAst{memoizedSize: -1}
	m.mergeFromReader(rd)
	typedAst := m.freeze()
	var stringShard *StringPool
	if buildStringPool == nil || buildStringPool(typedAst) {
		stringShard = StringPool_FromProto(typedAst.GetStringPool())
	}
	return typedAst, stringShard
}

// Java: TypedAstDeserializer.deserializeTypedAsts(InputStream, TypedAstDeserializer, AbstractCompiler, boolean, boolean)
func typedAstDeserializer_deserializeTypedAsts(
	typedAstsStream io.Reader,
	deserializer *TypedAstDeserializer,
	compiler AbstractCompiler,
	resolveSourceMapAnnotations bool,
	parseInlineSourceMaps bool) {
	if rlStream, ok := typedAstsStream.(*runtimeLibsTypedAstStream); ok {
		runtimeLibsTypedAstOnce.Do(func() {
			codedInput := &protoReader{b: rlStream.data, limit: math.MaxInt32}
			for !codedInput.atEnd() {
				tag := codedInput.readTag()
				if jutil.Ushr32(int32(tag), 3) != TypedAst_List_TYPED_ASTS_FIELD_NUMBER ||
					tag&7 != typedAstDeserializer_WIRETYPE_LENGTH_DELIMITED {
					panic(NewInvalidProtocolBufferException(
						"Unexpected field number " +
							strconv.Itoa(int(jutil.Ushr32(int32(tag), 3))) +
							" or wire type " +
							strconv.Itoa(tag&7)))
				}
				messageBytes := codedInput.readMessageBytes()
				typedAst, stringShard := typedAstDeserializer_decodeShardBytes(messageBytes, nil)
				runtimeLibsTypedAstShards = append(runtimeLibsTypedAstShards, cachedTypedAstShard{
					proto:       typedAst,
					stringShard: stringShard,
				})
			}
		})
		for i := range runtimeLibsTypedAstShards {
			sh := &runtimeLibsTypedAstShards[i]
			deserializer.deserializeTypedAstWithStringShard(
				sh.proto, sh.stringShard, compiler, resolveSourceMapAnnotations, parseInlineSourceMaps)
		}
		return
	}

	jutil.Try(func() {
		numWorkers := 1
		if compiler != nil && compiler.GetOptions() != nil && compiler.GetOptions().GetNumParallelThreads() > 1 {
			numWorkers = compiler.GetOptions().GetNumParallelThreads()
		} else if n := runtime.GOMAXPROCS(0); n > 1 {
			numWorkers = n
		}

		var requiredNames map[string]struct{}
		buildStringPool := func(typedAstProto *TypedAst) bool {
			if deserializer.mode != typedAstDeserializer_Mode_FULL_AST || deserializer.requiredInputFiles == nil {
				return true
			}
			if requiredNames == nil {
				requiredNames = make(map[string]struct{}, deserializer.requiredInputFiles.Size()*2)
				for f := range deserializer.requiredInputFiles.All() {
					name := f.GetName()
					requiredNames[name] = struct{}{}
					if canon := SourceFile_CanonicalizePathForMatching(name); canon != name {
						requiredNames[canon] = struct{}{}
					}
				}
			}
			protos := typedAstProto.GetSourceFilePool().GetSourceFileList()
			for _, lazyAst := range typedAstDeserializer_lazyAsts(typedAstProto) {
				idx := int(lazyAst.GetSourceFile()) - 1
				if idx >= 0 && idx < len(protos) {
					fn := protos[idx].GetFilename()
					if _, ok := requiredNames[fn]; ok {
						return true
					}
					if _, ok := requiredNames[SourceFile_CanonicalizePathForMatching(fn)]; ok {
						return true
					}
				}
			}
			return false
		}

		processShardBatch := func(batch [][]byte) {
			if len(batch) == 1 {
				typedAst, stringShard := typedAstDeserializer_decodeShardBytes(batch[0], buildStringPool)
				deserializer.deserializeTypedAstWithStringShard(
					typedAst, stringShard, compiler, resolveSourceMapAnnotations, parseInlineSourceMaps)
				return
			}
			// Trigger lazy initialization of requiredNames before launching workers.
			if deserializer.mode == typedAstDeserializer_Mode_FULL_AST && deserializer.requiredInputFiles != nil && requiredNames == nil {
				requiredNames = make(map[string]struct{}, deserializer.requiredInputFiles.Size()*2)
				for f := range deserializer.requiredInputFiles.All() {
					name := f.GetName()
					requiredNames[name] = struct{}{}
					if canon := SourceFile_CanonicalizePathForMatching(name); canon != name {
						requiredNames[canon] = struct{}{}
					}
				}
			}
			outcomes := make([]cachedTypedAstShard, len(batch))
			failures := parallelPrebuild(len(batch), min(numWorkers, len(batch)), func(i int) {
				typedAst, stringShard := typedAstDeserializer_decodeShardBytes(batch[i], buildStringPool)
				outcomes[i] = cachedTypedAstShard{proto: typedAst, stringShard: stringShard}
			})
			for i := range batch {
				if failures[i] != nil {
					panic(failures[i])
				}
				deserializer.deserializeTypedAstWithStringShard(
					outcomes[i].proto,
					outcomes[i].stringShard,
					compiler,
					resolveSourceMapAnnotations,
					parseInlineSourceMaps)
			}
		}

		if br, ok := typedAstsStream.(*bytes.Reader); ok {
			data := make([]byte, br.Len())
			if _, err := io.ReadFull(br, data); err != nil {
				panic(jutil.NewIOException(err.Error()))
			}
			codedInput := &protoReader{b: data, limit: math.MaxInt32}
			readNextFromCoded := func() ([]byte, bool) {
				if codedInput.atEnd() {
					return nil, false
				}
				tag := codedInput.readTag()
				if jutil.Ushr32(int32(tag), 3) != TypedAst_List_TYPED_ASTS_FIELD_NUMBER ||
					tag&7 != typedAstDeserializer_WIRETYPE_LENGTH_DELIMITED {
					panic(NewInvalidProtocolBufferException(
						"Unexpected field number " +
							strconv.Itoa(int(jutil.Ushr32(int32(tag), 3))) +
							" or wire type " +
							strconv.Itoa(tag&7)))
				}
				return codedInput.readMessageBytes(), true
			}
			firstMsg, ok := readNextFromCoded()
			if !ok {
				return
			}
			if codedInput.atEnd() {
				processShardBatch([][]byte{firstMsg})
				return
			}
			const batchSize = 32
			batch := make([][]byte, 0, batchSize)
			batch = append(batch, firstMsg)
			for {
				var readErr jutil.Throwable
				for len(batch) < batchSize {
					var msg []byte
					var hasMore bool
					jutil.Try(func() {
						msg, hasMore = readNextFromCoded()
					}, func(t jutil.Throwable) {
						readErr = t
					})
					if readErr != nil || !hasMore {
						break
					}
					batch = append(batch, msg)
				}
				if len(batch) > 0 {
					processShardBatch(batch)
					batch = batch[:0]
				}
				if readErr != nil {
					panic(readErr)
				}
				if codedInput.atEnd() {
					break
				}
			}
			return
		}

		var oneByte [1]byte
		readByte := func(allowEOF bool) (byte, bool) {
			for {
				n, err := typedAstsStream.Read(oneByte[:])
				if n == 1 {
					return oneByte[0], true
				}
				if err == io.EOF {
					if allowEOF {
						return 0, false
					}
					panic(protoTruncatedMessage())
				}
				if err != nil {
					panic(jutil.NewIOException(err.Error()))
				}
			}
		}
		readVarint64 := func(firstByte byte) uint64 {
			x := uint64(firstByte & 0x7f)
			if firstByte < 0x80 {
				return x
			}
			for shift := uint(7); shift < 64; shift += 7 {
				b, _ := readByte(false)
				x |= uint64(b&0x7f) << shift
				if b < 0x80 {
					return x
				}
			}
			panic(NewInvalidProtocolBufferException("CodedInputStream encountered a malformed varint."))
		}
		readNextFromStream := func() ([]byte, bool) {
			first, ok := readByte(true)
			if !ok {
				return nil, false
			}
			uTag := readVarint64(first)
			if uTag>>3 == 0 || uTag > math.MaxUint32 {
				panic(NewInvalidProtocolBufferException("Protocol message contained an invalid tag (zero)."))
			}
			tag := int(int32(uTag))
			if jutil.Ushr32(int32(tag), 3) != TypedAst_List_TYPED_ASTS_FIELD_NUMBER ||
				tag&7 != typedAstDeserializer_WIRETYPE_LENGTH_DELIMITED {
				panic(NewInvalidProtocolBufferException(
					"Unexpected field number " +
						strconv.Itoa(int(jutil.Ushr32(int32(tag), 3))) +
						" or wire type " +
						strconv.Itoa(tag&7)))
			}
			lenFirst, _ := readByte(false)
			msgLen := int(int32(readVarint64(lenFirst)))
			if msgLen < 0 {
				panic(protoNegativeSize())
			}
			// Do not pool messageBytes buffers because LazyAst.script uses zero-copy readBytesView().
			var messageBytes []byte
			const probeSize = 64 * 1024
			if msgLen <= probeSize {
				messageBytes = make([]byte, msgLen)
				if _, err := io.ReadFull(typedAstsStream, messageBytes); err != nil {
					if err == io.EOF || err == io.ErrUnexpectedEOF {
						panic(protoTruncatedMessage())
					}
					panic(jutil.NewIOException(err.Error()))
				}
			} else {
				var probe [probeSize]byte
				if _, err := io.ReadFull(typedAstsStream, probe[:]); err != nil {
					if err == io.EOF || err == io.ErrUnexpectedEOF {
						panic(protoTruncatedMessage())
					}
					panic(jutil.NewIOException(err.Error()))
				}
				messageBytes = make([]byte, msgLen)
				copy(messageBytes, probe[:])
				if _, err := io.ReadFull(typedAstsStream, messageBytes[probeSize:]); err != nil {
					if err == io.EOF || err == io.ErrUnexpectedEOF {
						panic(protoTruncatedMessage())
					}
					panic(jutil.NewIOException(err.Error()))
				}
			}
			return messageBytes, true
		}

		firstMsg, ok := readNextFromStream()
		if !ok {
			return
		}
		const batchSize = 32
		batch := make([][]byte, 0, batchSize)
		batch = append(batch, firstMsg)
		for {
			var readErr jutil.Throwable
			var reachedEOF bool
			for len(batch) < batchSize {
				var msg []byte
				var hasMore bool
				jutil.Try(func() {
					msg, hasMore = readNextFromStream()
				}, func(t jutil.Throwable) {
					readErr = t
				})
				if readErr != nil {
					break
				}
				if !hasMore {
					reachedEOF = true
					break
				}
				batch = append(batch, msg)
			}
			if len(batch) > 0 {
				processShardBatch(batch)
				batch = batch[:0]
			}
			if readErr != nil {
				panic(readErr)
			}
			if reachedEOF {
				break
			}
		}
	}, func(ex jutil.IIOException) {
		message := "Cannot read from TypedAST input stream"
		e := &jutil.IllegalArgumentException{}
		jutil.InitThrowable(e, "java.lang.IllegalArgumentException", &message, ex)
		panic(e)
	})
}

// typedAstDeserializer_WIRETYPE_LENGTH_DELIMITED is
// WireFormat.WIRETYPE_LENGTH_DELIMITED.
const typedAstDeserializer_WIRETYPE_LENGTH_DELIMITED = 2

// TypedAstDeserializer_DeserializedAst is the result of deserializing a
// TypedAst.List
//
// Go-only: language-adapter — an @AutoValue: the abstract class and the
// generated AutoValue_TypedAstDeserializer_DeserializedAst (fields, equals,
// hashCode, toString) are one struct (PORTING.md §4).
//
// Java: public abstract static class TypedAstDeserializer.DeserializedAst
type TypedAstDeserializer_DeserializedAst struct {
	filesystem       *jutil.LinkedHashMap[*SourceFile, func() *rhino.Node]
	colorRegistry    *colors.ColorRegistry // Optional: nil when absent
	externProperties *jutil.LinkedHashSet[string]
	runtimeLibraries *jutil.LinkedHashSet[string]
}

// GetFilesystem maps from SourceFile to a lazy deserializer of the SCRIPT
// node for that file
//
// The supplier creates a new Node whenever called (but the results should
// be .equals)
//
// Java: TypedAstDeserializer.DeserializedAst.getFilesystem()
func (a *TypedAstDeserializer_DeserializedAst) GetFilesystem() *jutil.LinkedHashMap[*SourceFile, func() *rhino.Node] {
	return a.filesystem
}

// GetColorRegistry is the built ColorRegistry (nil when absent).
//
// Note that this is absent if either a)
// TypedAstDeserializer#deserializeRuntimeLibraries(AbstractCompiler,
// SourceFile, ColorPool.Builder, InputStream) was called, as this sort of
// deserialization does not build a complete AST + colors, just a shard of
// it, or b) type information was not requested.
//
// Java: TypedAstDeserializer.DeserializedAst.getColorRegistry()
func (a *TypedAstDeserializer_DeserializedAst) GetColorRegistry() *colors.ColorRegistry {
	return a.colorRegistry
}

// GetExternProperties returns a list of all known extern properties,
// including properties that were present in type annotations in source code
// but not serialized on the AST
//
// Java: TypedAstDeserializer.DeserializedAst.getExternProperties()
func (a *TypedAstDeserializer_DeserializedAst) GetExternProperties() *jutil.LinkedHashSet[string] {
	return a.externProperties
}

// Java: TypedAstDeserializer.DeserializedAst.getRuntimeLibraries()
func (a *TypedAstDeserializer_DeserializedAst) GetRuntimeLibraries() *jutil.LinkedHashSet[string] {
	return a.runtimeLibraries
}

// Java: TypedAstDeserializer.DeserializedAst.create(ConcurrentMap<SourceFile, Supplier<Node>>, Optional<ColorRegistry>, ImmutableSet<String>, ImmutableSet<String>)
func typedAstDeserializer_DeserializedAst_create(
	filesystem *jutil.LinkedHashMap[*SourceFile, func() *rhino.Node],
	colorRegistry *colors.ColorRegistry,
	externProperties *jutil.LinkedHashSet[string],
	runtimeLibraries *jutil.LinkedHashSet[string]) *TypedAstDeserializer_DeserializedAst {
	return newAutoValue_TypedAstDeserializer_DeserializedAst(
		filesystem, colorRegistry, externProperties, runtimeLibraries)
}

// Java: AutoValue_TypedAstDeserializer_DeserializedAst(ConcurrentMap<SourceFile, Supplier<Node>>, Optional<ColorRegistry>, ImmutableSet<String>, ImmutableSet<String>)
func newAutoValue_TypedAstDeserializer_DeserializedAst(
	filesystem *jutil.LinkedHashMap[*SourceFile, func() *rhino.Node],
	colorRegistry *colors.ColorRegistry,
	externProperties *jutil.LinkedHashSet[string],
	runtimeLibraries *jutil.LinkedHashSet[string]) *TypedAstDeserializer_DeserializedAst {
	if filesystem == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null filesystem"))
	}
	// colorRegistry: the Optional itself is never null; absence is nil here.
	if runtimeLibraries == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null runtimeLibraries"))
	}
	return &TypedAstDeserializer_DeserializedAst{
		filesystem:       filesystem,
		colorRegistry:    colorRegistry,
		externProperties: externProperties,
		runtimeLibraries: runtimeLibraries,
	}
}

// String is the AutoValue toString().
//
// Java: AutoValue_TypedAstDeserializer_DeserializedAst.toString()
func (a *TypedAstDeserializer_DeserializedAst) String() string {
	colorRegistry := "Optional.absent()"
	if a.colorRegistry != nil {
		colorRegistry = "Optional.of(" + jutil.ToJavaString(a.colorRegistry) + ")"
	}
	externProperties := "null"
	if a.externProperties != nil {
		externProperties = a.externProperties.String()
	}
	return "DeserializedAst{" +
		"filesystem=" + a.filesystem.String() + ", " +
		"colorRegistry=" + colorRegistry + ", " +
		"externProperties=" + externProperties + ", " +
		"runtimeLibraries=" + a.runtimeLibraries.String() +
		"}"
}

// Equals is the AutoValue equals(Object).
//
// Go-only: language-adapter — the filesystem values are Supplier lambdas,
// which Java compares by identity and Go funcs cannot be compared at all:
// two filesystems are equal here only when they are the same map.
//
// Java: AutoValue_TypedAstDeserializer_DeserializedAst.equals(Object)
func (a *TypedAstDeserializer_DeserializedAst) Equals(o any) bool {
	if o == a {
		return true
	}
	if that, ok := o.(*TypedAstDeserializer_DeserializedAst); ok {
		externPropertiesEqual := false
		if a.externProperties == nil {
			externPropertiesEqual = that.GetExternProperties() == nil
		} else {
			externPropertiesEqual = a.externProperties.Equals(that.GetExternProperties())
		}
		return a.filesystem == that.GetFilesystem() &&
			a.colorRegistry == that.GetColorRegistry() &&
			externPropertiesEqual &&
			a.runtimeLibraries.Equals(that.GetRuntimeLibraries())
	}
	return false
}

// HashCode is the AutoValue hashCode().
//
// Go-only: language-adapter — the filesystem and color registry contribute
// identity hashes in Java (lambdas and ColorRegistry do not override
// hashCode); here they contribute their sizes / presence.
//
// Java: AutoValue_TypedAstDeserializer_DeserializedAst.hashCode()
func (a *TypedAstDeserializer_DeserializedAst) HashCode() int32 {
	h := int32(1)
	h *= 1000003
	h ^= int32(a.filesystem.Size())
	h *= 1000003
	if a.colorRegistry != nil {
		h ^= 1
	}
	h *= 1000003
	if a.externProperties != nil {
		h ^= a.externProperties.HashCode()
	}
	h *= 1000003
	h ^= a.runtimeLibraries.HashCode()
	return h
}
