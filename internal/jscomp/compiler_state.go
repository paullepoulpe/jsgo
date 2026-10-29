// Port of com/google/javascript/jscomp/Compiler.java @ closure-compiler v20261021

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

package jscomp

import (
	"bufio"
	"io"
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Java: Compiler.toProto(LifeCycleStage)
func compiler_toProtoOfLifeCycleStage(stage AbstractCompiler_LifeCycleStage) LifeCycleStageProto {
	switch stage {
	case AbstractCompiler_LifeCycleStage_RAW:
		return LifeCycleStageProto_LIFE_CYCLE_STAGE_RAW
	case AbstractCompiler_LifeCycleStage_COLORS_AND_SIMPLIFIED_JSDOC:
		return LifeCycleStageProto_LIFE_CYCLE_STAGE_COLORS_AND_SIMPLIFIED_JSDOC
	case AbstractCompiler_LifeCycleStage_NORMALIZED:
		return LifeCycleStageProto_LIFE_CYCLE_STAGE_NORMALIZED
	case AbstractCompiler_LifeCycleStage_NORMALIZED_OBFUSCATED:
		return LifeCycleStageProto_LIFE_CYCLE_STAGE_NORMALIZED_OBFUSCATED
	}
	// Go-only: language-adapter — Java's exhaustive enum switch expression
	// throws MatchException for an unknown constant.
	panic(jutil.NewIllegalStateExceptionMsg("unexpected LifeCycleStage " + stage.String()))
}

// Java: Compiler.toVariableMapEntries(VariableMap)
func compiler_toVariableMapEntries(m *VariableMap) []*VariableMapEntryProto {
	entries := []*VariableMapEntryProto{}
	if m != nil {
		for key, value := range m.GetOriginalNameToNewNameMap().All() {
			entries = append(entries, &VariableMapEntryProto{
				originalName_: key,
				newName_:      value,
			})
		}
	}
	return entries
}

// Java: Compiler.fromProto(List<FeatureProto>)
func compiler_fromProto(protos []FeatureProto) *parser.FeatureSet {
	fs := parser.FeatureSet_BARE_MINIMUM
	for _, p := range protos {
		if p != FeatureProto_FEATURE_UNKNOWN {
			if name, ok := strings.CutPrefix(p.String(), "FEATURE_"); ok {
				if feat, ok := parser.FeatureSet_Feature_Lookup(name); ok {
					fs = fs.With(feat)
				}
			}
		}
	}
	return fs
}

// Java: Compiler.fromProto(LifeCycleStageProto)
func compiler_fromProtoOfLifeCycleStage(stage LifeCycleStageProto) AbstractCompiler_LifeCycleStage {
	switch stage {
	case LifeCycleStageProto_LIFE_CYCLE_STAGE_RAW:
		return AbstractCompiler_LifeCycleStage_RAW
	case LifeCycleStageProto_LIFE_CYCLE_STAGE_COLORS_AND_SIMPLIFIED_JSDOC:
		return AbstractCompiler_LifeCycleStage_COLORS_AND_SIMPLIFIED_JSDOC
	case LifeCycleStageProto_LIFE_CYCLE_STAGE_NORMALIZED:
		return AbstractCompiler_LifeCycleStage_NORMALIZED
	case LifeCycleStageProto_LIFE_CYCLE_STAGE_NORMALIZED_OBFUSCATED:
		return AbstractCompiler_LifeCycleStage_NORMALIZED_OBFUSCATED
	default:
		return AbstractCompiler_LifeCycleStage_RAW
	}
}

// Java: Compiler.fromVariableMapEntries(List<VariableMapEntryProto>)
func compiler_fromVariableMapEntries(entries []*VariableMapEntryProto) *VariableMap {
	m := jutil.NewLinkedHashMap[string, string]()
	for _, entry := range entries {
		m.Put(entry.GetOriginalName(), entry.GetNewName())
	}
	return NewVariableMap(m)
}

// GetCompilerStateProtoBuilder returns a builder for the serializable state
// of the compiler.
//
// Java: Compiler.getCompilerStateProtoBuilder()
func (c *Compiler) GetCompilerStateProtoBuilder() *JSCompilerStateProto_Builder {
	builder := JSCompilerStateProto_NewBuilder()

	for feature := range c.allowableFeatures.All() {
		builder.AddAllowableFeatures(compiler_toProto(feature))
	}
	builder.
		SetTypeCheckingHasRun(c.typeCheckingHasRun).
		SetHasRegExpGlobalReferences(c.hasRegExpGlobalReferences).
		SetLifeCycleStage(compiler_toProtoOfLifeCycleStage(c.GetLifeCycleStage())).
		SetMergedPrecompiledLibraries(c.options.GetMergedPrecompiledLibraries())

	chunkList := c.chunkGraph.GetIntegralChunkArrayForSerialization()
	chunkNameMap := make(map[string]int, len(chunkList))
	for i, chunk := range chunkList {
		chunkNameMap[chunk.GetName()] = i
	}

	for _, chunk := range chunkList {
		deps := make([]int32, 0, len(chunk.GetDependencies()))
		for _, dep := range chunk.GetDependencies() {
			depIdx, ok := chunkNameMap[dep.GetName()]
			jutil.CheckStateMsg2(
				ok, "Missing dependency %s for chunk %s", dep.GetName(), chunk.GetName())
			deps = append(deps, int32(depIdx))
		}
		inputs := chunk.GetInputs()
		inputIds := make([]string, len(inputs))
		for i, input := range inputs {
			inputIds[i] = input.GetInputId().GetIdName()
		}
		builder.AddChunks(&ChunkProto{
			name_:         chunk.GetName(),
			dependencies_: deps,
			inputIds_:     inputIds,
		})
	}

	builder.
		SetUniqueNameId(int32(c.uniqueNameId)).
		AddAllUniqueIdSupplier(c.uniqueIdSupplier.ToProto()).
		AddAllExportedNames(c.exportedNames.ToSlice())

	builder.SetHasCssNames(c.cssNames != nil)
	if c.cssNames != nil {
		builder.AddAllCssNames(c.cssNames.ToSlice())
	}

	if c.idGeneratorMap != nil {
		builder.SetIdGeneratorMap(*c.idGeneratorMap)
	}

	builder.
		SetTranspiledFiles(c.transpiledFiles).
		SetIdGeneratorCurrentId(int32(c.crossChunkIdGenerator.GetCurrentId())).
		SetRunJ2ClPasses(c.runJ2clPasses)

	for _, extern := range c.externs {
		builder.AddExterns(extern.GetInputId().GetIdName())
	}

	builder.AddAllInjectedLibraries(c.runtimeJsLibManager.GetInjectedLibraries())
	lastInjectedLibrary := c.runtimeJsLibManager.GetLastInjectedLibrary()
	if !lastInjectedLibrary.IsNil() && !c.jsRoot.IsNil() && c.jsRoot.HasChildren() {
		builder.SetLastInjectedLibraryIndexInFirstScript(
			int32(c.jsRoot.GetFirstChild().GetIndexOfChild(lastInjectedLibrary)))
	} else {
		builder.SetLastInjectedLibraryIndexInFirstScript(-1)
	}

	if c.accessorSummary != nil {
		builder.SetAccessorSummary(c.accessorSummary.ToProto())
	}

	builder.SetHasStringMap(c.stringMap != nil)
	if c.stringMap != nil {
		builder.AddAllStringMap(compiler_toVariableMapEntries(c.stringMap))
	}

	builder.SetHasInstrumentationMapping(c.instrumentationMapping != nil)
	if c.instrumentationMapping != nil {
		builder.AddAllInstrumentationMapping(compiler_toVariableMapEntries(c.instrumentationMapping))
	}

	if len(c.conformanceViolations) > 0 {
		builder.AddAllConformanceViolations(c.conformanceViolations)
	}
	if len(c.unattributedConformanceViolations) > 0 {
		builder.AddAllUnattributedConformanceViolations(c.unattributedConformanceViolations)
	}

	return builder
}

// RestoreFromState restores the portions of the compiler state that don't
// require access to the serialized AST.
//
// Java: Compiler.restoreFromState(JSCompilerStateProto)
func (c *Compiler) RestoreFromState(compilerState *JSCompilerStateProto) {
	c.allowableFeatures = compiler_fromProto(compilerState.GetAllowableFeaturesList())
	c.typeCheckingHasRun = compilerState.GetTypeCheckingHasRun()
	clear(c.scriptNodeByFilename)
	c.hasRegExpGlobalReferences = compilerState.GetHasRegExpGlobalReferences()

	for _, library := range compilerState.GetInjectedLibrariesList() {
		c.GetRuntimeJsLibManager().RecordLibraryInjected(library)
	}

	stage := AbstractCompiler_LifeCycleStage_COLORS_AND_SIMPLIFIED_JSDOC
	if compilerState.GetLifeCycleStage() != LifeCycleStageProto_LIFE_CYCLE_STAGE_RAW {
		stage = compiler_fromProtoOfLifeCycleStage(compilerState.GetLifeCycleStage())
	}
	c.SetLifeCycleStage(stage)
	c.GetOptions().SetMergedPrecompiledLibraries(compilerState.GetMergedPrecompiledLibraries())

	c.uniqueNameId = int(compilerState.GetUniqueNameId())
	c.uniqueIdSupplier = UniqueIdSupplier_FromProto(compilerState.GetUniqueIdSupplierList())

	c.exportedNames.Clear()
	c.exportedNames.AddAll(compilerState.GetExportedNamesList()...)

	if compilerState.GetHasCssNames() {
		c.cssNames = jutil.NewLinkedHashSet(compilerState.GetCssNamesList()...)
	} else {
		c.cssNames = nil
	}

	c.variableMap = nil
	c.propertyMap = nil

	if compilerState.HasIdGeneratorMap() {
		idGeneratorMap := compilerState.GetIdGeneratorMap()
		c.idGeneratorMap = &idGeneratorMap
	} else {
		c.idGeneratorMap = nil
	}
	c.transpiledFiles = compilerState.GetTranspiledFiles()

	c.crossChunkIdGenerator = NewIdGenerator()
	c.crossChunkIdGenerator.SetCurrentId(int(compilerState.GetIdGeneratorCurrentId()))

	c.runJ2clPasses = compilerState.GetRunJ2ClPasses()

	if compilerState.HasAccessorSummary() {
		c.accessorSummary = AccessorSummary_FromProto(compilerState.GetAccessorSummary())
	} else {
		c.accessorSummary = nil
	}

	if compilerState.GetHasStringMap() {
		c.stringMap = compiler_fromVariableMapEntries(compilerState.GetStringMapList())
	} else {
		c.stringMap = nil
	}

	if compilerState.GetHasInstrumentationMapping() {
		c.instrumentationMapping =
			compiler_fromVariableMapEntries(compilerState.GetInstrumentationMappingList())
	} else {
		c.instrumentationMapping = nil
	}

	c.conformanceViolations = slices.Clone(compilerState.GetConformanceViolationsList())
	c.unattributedConformanceViolations = slices.Clone(compilerState.GetUnattributedConformanceViolationsList())

	// We don't save changeStamp, because its value turned out to be non-deterministic
	// (Reasons for this are unknown.), and there's no benefit to saving it anyway.
	// We don't save the change stamps that are stored on the AST nodes,
	// so all the AST Nodes we read in effectively have a change stamp of 0,
	// and we can just start the compiler's counter over at 1.
	c.changeTracker.ResetChangeStamp()
}

// Java: Compiler.saveState(OutputStream)
func (c *Compiler) SaveState(outputStream io.Writer) {
	// Do not close the outputstream, caller is responsible for closing it.
	c.runInCompilerThread(func() {
		tracer := c.NewTracer("serializeCompilerState")
		nonClosingOut := NewCompiler_NonClosingOutputStream(outputStream)
		var compressionStream io.Writer
		wrapper := c.options.GetStateCompressionWrapper()
		if wrapper != nil {
			compressionStream = wrapper(nonClosingOut)
		} else {
			compressionStream = newSaveStateGzipOutputStream(nonClosingOut, c.options.GetNumParallelThreads())
		}

		bufferedCompressionStream := compiler_newBufferedOutputStream(compressionStream, 64*1024)

		builder := c.GetCompilerStateProtoBuilder()
		compiler_writeDelimitedTo(builder.Build().ToByteArray(), bufferedCompressionStream)
		c.StopTracer(tracer, "serializeCompilerState")
		tracer = c.NewTracer("serializeTypedAst")
		SerializeTypedAstPass_CreateFromOutputStream(
			c,
			bufferedCompressionStream,
			&SerializationOptions{
				includeDebugInfo: c.GetOptions().ShouldSerializeExtraDebugInfo(),
				runtimeLibraries: emptyRuntimeLibraries,
			}).
			Process(c.externsRoot, c.jsRoot)
		c.StopTracer(tracer, "serializeTypedAst")
		// Close the buffered stream to flush uncompressed bytes and close the compressor
		bufferedCompressionStream.Close()
	})
}

// Java: Compiler.restoreState(InputStream)
func (c *Compiler) RestoreState(inputStream io.Reader) {
	c.initWarningsGuard(c.options.GetWarningsGuard())
	c.MaybeSetTracker()

	c.runInCompilerThread(func() {
		tracer := c.NewTracer(PassNames_DESERIALIZE_COMPILER_STATE)
		Compiler_logger.Fine("Deserializing the CompilerState")
		nonClosingIn := NewCompiler_NonClosingInputStream(inputStream)
		var decompressStream io.Reader
		wrapper := c.options.GetStateDecompressionWrapper()
		if wrapper != nil {
			decompressStream = wrapper(nonClosingIn)
		} else {
			decompressStream = jutil.NewGZIPInputStreamWithSize(nonClosingIn, 64*1024)
		}
		bufferedDecompressStream := compiler_newBufferedInputStream(decompressStream, 64*1024)
		defer func() {
			bufferedDecompressStream.Close()
			Compiler_logger.Fine("Finished deserializing CompilerState")
			c.StopTracer(tracer, PassNames_DESERIALIZE_COMPILER_STATE)
		}()
		c.deserializeCompilerState(bufferedDecompressStream)
	})

	if c.tracker != nil {
		c.tracker.UpdateAfterDeserialize(c.jsRoot)
	}
}

// deserializeCompilerState must be called from within a "compiler thread"
// with a larger stack
//
// Java: Compiler.deserializeCompilerState(InputStream)
func (c *Compiler) deserializeCompilerState(inputStream io.Reader) {
	// Do not close the input stream, caller is responsible for closing it.
	stateProto := compiler_parseDelimitedFrom(inputStream)

	jutil.CheckNotNullMsg(
		c.chunkGraph, "Did you forget to call .init or .initChunks before restoreState?")
	allInputs := c.chunkGraph.GetAllInputs()
	externFiles := make(map[string]*SourceFile, len(c.externs))
	codeFiles := make(map[string]*SourceFile, len(allInputs))
	allInputFiles := make([]*SourceFile, 0, len(allInputs)+len(c.externs))
	compiler_putOrThrow := func(m map[string]*SourceFile, k string, v *SourceFile) {
		// ImmutableMap.Builder.buildOrThrow: duplicate keys are an error.
		if _, dup := m[k]; dup {
			panic(jutil.NewIllegalArgumentExceptionMsg("Multiple entries with same key: " + k))
		}
		m[k] = v
	}
	for _, input := range allInputs {
		sf := input.GetSourceFile()
		allInputFiles = append(allInputFiles, sf)
		compiler_putOrThrow(codeFiles, input.GetInputId().GetIdName(), sf)
	}
	for _, extern := range c.externs {
		sf := extern.GetSourceFile()
		allInputFiles = append(allInputFiles, sf)
		compiler_putOrThrow(externFiles, extern.GetInputId().GetIdName(), sf)
	}

	deserializedAst := TypedAstDeserializer_DeserializeFullAst(
		c,
		c.SYNTHETIC_EXTERNS_FILE,
		allInputFiles,
		inputStream,
		stateProto.GetTypeCheckingHasRun(),
		c.GetOptions().GetResolveSourceMapAnnotations(),
		c.GetOptions().GetParseInlineSourceMaps())

	c.RestoreFromState(stateProto)

	// Restore the chunk graph
	chunkProtos := stateProto.GetChunksList()
	chunks := make([]*JSChunk, 0, len(chunkProtos))
	for _, chunkProto := range chunkProtos {
		chunks = append(chunks, NewJSChunk(chunkProto.GetName()))
	}
	for i := 0; i < len(chunks); i++ {
		chunk := chunks[i]
		chunkProto := chunkProtos[i]
		for _, depIdx := range chunkProto.GetDependenciesList() {
			chunk.AddDependency(chunks[depIdx])
		}
	}
	c.chunkGraph = NewJSChunkGraph(chunks)

	// Restore TypedAST and related fields
	c.externProperties = deserializedAst.GetExternProperties()
	c.externAndJsRoot = rhino.IR_Root(rhino.IR_Root(), rhino.IR_Root())
	c.externsRoot = c.externAndJsRoot.GetFirstChild()
	c.jsRoot = c.externAndJsRoot.GetLastChild()
	clear(c.inputsById)
	clear(c.externs)
	c.externs = c.externs[:0]

	c.colorRegistry = deserializedAst.GetColorRegistry()
	runtimeLibraries := deserializedAst.GetRuntimeLibraries()

	// Tells CompilerInput::getRoot to deserialize an AST rather than re-parsing the file
	c.typedAstFilesystem = deserializedAst.filesystem
	deserializedAst = nil

	// overwrite any existing CompilerInput instances. Reuse the SourceFiles created by
	// TypedAstDeserializer; otherwise lookups in this.typedAstFilesystem will fail.
	for _, externName := range stateProto.GetExternsList() {
		if externName == c.SYNTHETIC_EXTERNS_FILE.GetName() {
			c.GetSynthesizedExternsInput()
			continue
		}
		externFile := externFiles[externName]
		if externFile == nil {
			// Extern files may passed in with regular source files, and later 'hoisted' to externs
			// because they are annotated `@externs`.
			externFile = jutil.CheckNotNullMsg1(codeFiles[externName], "Missing %s", externName)
		}
		input := NewCompilerInputWithIsExtern(externFile /* isExtern= */, true)
		script := input.GetAstRoot(c) // accesses this.typedAstFilesystem

		c.externsRoot.AddChildToBack(script)
		c.putCompilerInput(input)
		c.scriptNodeByFilename[externFile.GetName()] = script
		c.externs = append(c.externs, input)
	}
	c.typedAstFilesystem.Remove(c.SYNTHETIC_EXTERNS_FILE)
	externFiles = nil

	// Restore the inputs for each chunk.
	type restoredChunkInput struct {
		chunkIdx int
		src      *SourceFile
		input    *CompilerInput
	}
	var restoredInputs []restoredChunkInput
	var prebuildInputs []*CompilerInput
	parallelRestore := c.options != nil && c.options.GetNumParallelThreads() > 1
	for i := 0; i < len(chunks); i++ {
		chunkProto := chunkProtos[i]
		for _, inputId := range chunkProto.GetInputIdsList() {
			src := codeFiles[inputId]
			if src == nil {
				// The auto-generated empty fill files used to facilitate CCCM for
				// empty chunks may not have gotten serialized, but all the others
				// should have.
				jutil.CheckStateMsg1(AbstractCompiler_IsFillFileName(inputId), "Missing %s", inputId)
				continue
			}
			input := NewCompilerInput(src)
			if parallelRestore {
				restoredInputs = append(restoredInputs, restoredChunkInput{chunkIdx: i, src: src, input: input})
				prebuildInputs = append(prebuildInputs, input)
			} else {
				script := input.GetAstRoot(c) // accesses this.typedAstFilesystem

				c.jsRoot.AddChildToBack(script)
				c.scriptNodeByFilename[src.GetName()] = script

				c.putCompilerInput(input) // overwrite the old input
				chunks[i].Add(input)
			}
		}
	}
	codeFiles = nil
	if parallelRestore && len(restoredInputs) > 0 {
		NewPrebuildAst(c, c.options.GetNumParallelThreads()).Prebuild(slices.Values(prebuildInputs))
		prebuildInputs = nil
		for _, item := range restoredInputs {
			script := item.input.GetAstRoot(c)

			c.jsRoot.AddChildToBack(script)
			c.scriptNodeByFilename[item.src.GetName()] = script

			c.putCompilerInput(item.input) // overwrite the old input
			chunks[item.chunkIdx].Add(item.input)
		}
		restoredInputs = nil
	}

	for library := range runtimeLibraries.All() {
		c.runtimeJsLibManager.EnsureLibraryInjected(library, false)
	}

	c.typedAstFilesystem = nil // allow garbage collection

	if stateProto.GetLastInjectedLibraryIndexInFirstScript() != -1 {
		lastInjectedLibrary := c.jsRoot.
			GetFirstChild().
			GetChildAtIndex(int(stateProto.GetLastInjectedLibraryIndexInFirstScript()))
		c.runtimeJsLibManager.SetLastInjectedLibrary(lastInjectedLibrary)
	}
}

// compiler_writeDelimitedTo is message.writeDelimitedTo(out): the message's
// size as a varint, then the message.
//
// Go-only: language-adapter — MessageLite.writeDelimitedTo(OutputStream) of
// the protoc-generated message.
func compiler_writeDelimitedTo(message []byte, out io.Writer) {
	buf := protoAppendVarint(make([]byte, 0, 10+len(message)), uint64(len(message)))
	buf = append(buf, message...)
	if _, err := out.Write(buf); err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
}

// compiler_parseDelimitedFrom is JSCompilerStateProto.parseDelimitedFrom(in):
// it reads a varint size, then parses at most that many bytes of in (fewer if
// in ends first, as protobuf's LimitedInputStream does), and leaves in right
// after them. It returns nil (Java null) if in is at its end.
//
// Go-only: language-adapter — Parser.parseDelimitedFrom(InputStream) of the
// protoc-generated message.
func compiler_parseDelimitedFrom(in io.Reader) *JSCompilerStateProto {
	readByte := func() (byte, bool) {
		var b [1]byte
		for {
			n, err := in.Read(b[:])
			if n == 1 {
				return b[0], true
			}
			if err == io.EOF {
				return 0, false
			}
			if err != nil {
				panic(jutil.NewIOException(err.Error()))
			}
		}
	}
	firstByte, ok := readByte()
	if !ok {
		return nil
	}
	// CodedInputStream.readRawVarint32(int, InputStream)
	size := uint32(firstByte & 0x7f)
	if firstByte&0x80 != 0 {
		offset := 7
		for ; offset < 32; offset += 7 {
			b, ok := readByte()
			if !ok {
				panic(protoTruncatedMessage())
			}
			size |= uint32(b&0x7f) << offset
			if b&0x80 == 0 {
				break
			}
		}
		if offset >= 32 {
			// Keep reading up to 64 bits.
			for ; offset < 64; offset += 7 {
				b, ok := readByte()
				if !ok {
					panic(protoTruncatedMessage())
				}
				if b&0x80 == 0 {
					break
				}
			}
			if offset >= 64 {
				panic(protoMalformedVarint())
			}
		}
	}
	if int32(size) < 0 {
		panic(protoNegativeSize())
	}
	data, err := io.ReadAll(io.LimitReader(in, int64(size)))
	if err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
	return JSCompilerStateProto_ParseFrom(data)
}

// compiler_bufferedOutputStream is new BufferedOutputStream(out, size).
//
// Go-only: language-adapter — java.io.BufferedOutputStream over bufio:
// close() flushes the buffer, then closes out (FilterOutputStream.close).
type compiler_bufferedOutputStream struct {
	w   *bufio.Writer
	out io.Writer
}

func compiler_newBufferedOutputStream(out io.Writer, size int) *compiler_bufferedOutputStream {
	return &compiler_bufferedOutputStream{w: bufio.NewWriterSize(out, size), out: out}
}

func (s *compiler_bufferedOutputStream) Write(b []byte) (int, error) {
	return s.w.Write(b)
}

func (s *compiler_bufferedOutputStream) Close() {
	if err := s.w.Flush(); err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
	if closer, ok := s.out.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			panic(jutil.NewIOException(err.Error()))
		}
	}
}

// compiler_bufferedInputStream is new BufferedInputStream(in, size).
//
// Go-only: language-adapter — java.io.BufferedInputStream over bufio:
// close() closes in.
type compiler_bufferedInputStream struct {
	r  *bufio.Reader
	in io.Reader
}

func compiler_newBufferedInputStream(in io.Reader, size int) *compiler_bufferedInputStream {
	return &compiler_bufferedInputStream{r: bufio.NewReaderSize(in, size), in: in}
}

func (s *compiler_bufferedInputStream) Read(b []byte) (int, error) {
	return s.r.Read(b)
}

func (s *compiler_bufferedInputStream) Close() {
	if closer, ok := s.in.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			panic(jutil.NewIOException(err.Error()))
		}
	}
}

// GetModuleTypeByName returns the module type for the provided namespace.
//
// Java: Compiler.getModuleTypeByName(String)
func (c *Compiler) GetModuleTypeByName(moduleName string) *CompilerInput_ModuleType {
	// Map.get: null when absent.
	if moduleType, ok := c.moduleTypesByName[moduleName]; ok {
		return &moduleType
	}
	return nil
}
