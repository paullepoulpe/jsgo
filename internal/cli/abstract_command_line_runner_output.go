// Port of com/google/javascript/jscomp/AbstractCommandLineRunner.java @ closure-compiler v20261104

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

package cli

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/sourcemap"
)

// parseChunkWrappers parses chunk wrapper specifications of the form
// <name>:<wrapper>, returning a map from chunk name to chunk wrapper (or ""
// for chunks with no wrapper).
//
// Java: AbstractCommandLineRunner.parseChunkWrappers(List<String>, Iterable<JSChunk>)
func parseChunkWrappers(specs []string, chunks []*jscomp.JSChunk) map[string]string {
	if specs == nil {
		panic(jutil.NewIllegalStateException())
	}

	wrappers := make(map[string]string, len(chunks))
	for _, c := range chunks {
		wrappers[c.GetName()] = ""
	}

	for _, spec := range specs {
		name, wrapper, hasColon := strings.Cut(spec, ":")
		if !hasColon {
			panic(jscomp.NewFlagUsageException(
				"Expected chunk wrapper to have <name>:<wrapper> format: " + spec))
		}
		if _, ok := wrappers[name]; !ok {
			panic(jscomp.NewFlagUsageException("Unknown chunk: '" + name + "'"))
		}
		wrapper = strings.ReplaceAll(strings.ReplaceAll(wrapper, "%output%", "%s"), "%n%", "\n")
		if !strings.Contains(wrapper, "%s") {
			panic(jscomp.NewFlagUsageException("No %s placeholder in chunk wrapper: '" + wrapper + "'"))
		}
		wrappers[name] = wrapper
	}
	return wrappers
}

// parseChunkOutputFiles parses chunk output name specifications of the form
// <name>:<output_file>, returning a map from chunk name to declared output
// file name.
//
// Java: AbstractCommandLineRunner.parseChunkOutputFiles(List<String>)
func parseChunkOutputFiles(specs []string) map[string]string {
	if specs == nil {
		panic(jutil.NewIllegalArgumentException())
	}
	if len(specs) == 0 {
		return map[string]string{}
	}

	entries := make([]PipeFlagEntry, 0, len(specs))
	for _, spec := range specs {
		key, val, hasColon := strings.Cut(spec, ":")
		if !hasColon {
			panic(jscomp.NewFlagUsageException(
				"Expected chunk_output_file to have <name>:<output_file> format: " + spec))
		}
		entries = append(entries, PipeFlagEntry{Key: key, Value: val})
	}
	checkNoDuplicatePipeFlagKeys(entries)

	outputFiles := make(map[string]string, len(entries))
	for _, e := range entries {
		outputFiles[e.Key] = e.Value
	}
	return outputFiles
}

// GetChunkOutputFileName returns the output file name for a chunk.
//
// For chunks with predeclared output file names specified using
// --chunk_output_file, the output file name is
// <outputPathPrefix>/<output_file>.
//
// Otherwise, the output file name is <outputPathPrefix>/<chunkName>.js
//
// Java: AbstractCommandLineRunner.getChunkOutputFileName(JSChunk)
func (r *abstractCommandLineRunnerBase) GetChunkOutputFileName(m *jscomp.JSChunk) string {
	if r.parsedChunkOutputFiles == nil {
		r.parsedChunkOutputFiles = parseChunkOutputFiles(r.config.chunkOutputFiles)
	}
	if outputFile, ok := r.parsedChunkOutputFiles[m.GetName()]; ok {
		return r.config.chunkOutputPathPrefix + outputFile
	}
	return r.config.chunkOutputPathPrefix + m.GetName() + ".js"
}

// GetChunkConformanceFileName returns the conformance output file name for a
// chunk, or "" if none was configured via --chunk_conformance_file.
//
// Java: AbstractCommandLineRunner.getChunkConformanceFileName(JSChunk)
func (r *abstractCommandLineRunnerBase) GetChunkConformanceFileName(m *jscomp.JSChunk) string {
	if r.parsedChunkConformanceFiles == nil {
		r.parsedChunkConformanceFiles = parseChunkOutputFiles(r.config.chunkConformanceFiles)
	}
	return r.parsedChunkConformanceFiles[m.GetName()]
}

func chunkWrapperBaseName(fileName string) string {
	if fileName != "" {
		if b := filepath.Base(fileName); b != "/" {
			return b
		}
	}
	return ""
}

// Java: AbstractCommandLineRunner.writeChunkOutput(String, Appendable, LicenseTracker, JSChunk)
func (r *abstractCommandLineRunnerBase) WriteChunkOutput(fileName string, out jutil.Appendable, lt jscomp.CodePrinter_LicenseTracker, m *jscomp.JSChunk) {
	if r.parsedChunkWrappers == nil {
		r.parsedChunkWrappers = parseChunkWrappers(
			r.config.chunkWrapper, r.compiler.GetChunkGraph().GetAllChunks())
	}

	if !r.isOutputInJson() {
		maybeCreateDirsForPath(fileName)
	}
	baseName := chunkWrapperBaseName(fileName)
	r.self.WriteOutput(
		out,
		r.compiler,
		lt,
		m,
		strings.ReplaceAll(r.parsedChunkWrappers[m.GetName()], "%basename%", baseName),
		"%s",
		nil,
		fileName)
}

// WriteOutput writes code to an output stream, optionally wrapping it in an
// arbitrary wrapper that contains a placeholder where the code should be
// inserted. If chunk is nil, it writes the entire AST.
//
// Java: AbstractCommandLineRunner.writeOutput(Appendable, Compiler, LicenseTracker, JSChunk, String, String, Function<String, String>, String)
func (r *abstractCommandLineRunnerBase) WriteOutput(
	out jutil.Appendable,
	compiler *jscomp.Compiler,
	licenseTracker jscomp.CodePrinter_LicenseTracker,
	chunk *jscomp.JSChunk,
	wrapper string,
	codePlaceholder string,
	escaper func(string) string,
	filename string) {
	r.writeOutput(out, compiler, licenseTracker, chunk, wrapper, codePlaceholder, escaper, filename)
}

func (r *abstractCommandLineRunnerBase) writeOutput(
	out jutil.Appendable,
	compiler *jscomp.Compiler,
	licenseTracker jscomp.CodePrinter_LicenseTracker,
	chunk *jscomp.JSChunk,
	wrapper string,
	codePlaceholder string,
	escaper func(string) string,
	filename string) {
	outputJs := compiler.GetOptions().GetOutputJs()
	if outputJs == jscomp.CompilerOptions_OutputJs_SENTINEL {
		out.AppendCharSequence("// No JS output because the compiler was run in checks-only mode.\n")
		return
	}
	if outputJs != jscomp.CompilerOptions_OutputJs_NORMAL {
		panic(jutil.NewIllegalStateException())
	}

	var code string
	if chunk == nil {
		code = compiler.ToSource()
	} else {
		code = compiler.ToSourceOfLicenseTrackerAndJSChunk(licenseTracker, chunk)
	}
	r.writeOutputOfCode(out, compiler, code, wrapper, codePlaceholder, escaper, filename)
}

// writeOutputOfCode writes code to an output stream, optionally wrapping it
// in an arbitrary wrapper that contains a placeholder where the code should
// be inserted.
//
// Java: AbstractCommandLineRunner.writeOutput(Appendable, Compiler, String, String, String, Function<String, String>, String)
func (r *abstractCommandLineRunnerBase) writeOutputOfCode(
	out jutil.Appendable,
	compiler *jscomp.Compiler,
	code string,
	wrapper string,
	codePlaceholder string,
	escaper func(string) string,
	_ string) {
	var sm *jscomp.SourceMap
	if compiler != nil {
		sm = compiler.GetSourceMap()
	}
	r.writeOutputOfCodeWithSourceMap(out, sm, code, wrapper, codePlaceholder, escaper)
}

func (r *abstractCommandLineRunnerBase) writeOutputOfCodeWithSourceMap(
	out jutil.Appendable,
	sm *jscomp.SourceMap,
	code string,
	wrapper string,
	codePlaceholder string,
	escaper func(string) string) {
	if prefix, suffix, found := strings.Cut(wrapper, codePlaceholder); found {
		if prefix != "" {
			out.AppendCharSequence(prefix)
		}
		if escaper == nil {
			out.AppendCharSequence(code)
		} else {
			out.AppendCharSequence(escaper(code))
		}
		if suffix != "" {
			out.AppendCharSequence(suffix)
		}
		if r.config.includeTrailingNewline {
			out.AppendCharacter('\n')
		}
		// If we have a source map, adjust its offsets to match
		// the code WITHIN the wrapper.
		if sm != nil {
			sm.SetWrapperPrefix(prefix)
		}
		return
	}

	out.AppendCharSequence(code)
	if r.config.includeTrailingNewline {
		out.AppendCharacter('\n')
	}
}

// maybeCreateDirsForPath creates any directories necessary to write a file
// that will have a given path prefix.
//
// Java: AbstractCommandLineRunner.maybeCreateDirsForPath(String)
func maybeCreateDirsForPath(pathPrefix string) {
	if trimmed := strings.TrimRight(pathPrefix, "/"); trimmed != "" {
		if dir := filepath.Dir(trimmed); dir != "." {
			_ = os.MkdirAll(dir, 0o777)
		}
	}
}

// Java: AbstractCommandLineRunner.createDefaultOutput()
func (r *abstractCommandLineRunnerBase) createDefaultOutput() *cliBufferedWriter {
	if r.config.jsOutputFile != "" {
		return r.fileNameToLegacyOutputWriter(r.config.jsOutputFile)
	}
	cs := r.legacyOutputCharset
	if cs == nil {
		cs = jutil.StandardCharsets_UTF_8
	}
	return newCliBufferedWriter(r.defaultJsOutput, nil, r.defaultJsOutput.Flush, cs)
}

// Java: AbstractCommandLineRunner.closeAppendable(Appendable)
func closeAppendable(output jutil.Appendable) {
	if flushable, ok := output.(interface{ Flush() }); ok {
		flushable.Flush()
	}
	if closeable, ok := output.(interface{ Close() }); ok {
		closeable.Close()
	}
}

// Java: AbstractCommandLineRunner.getJavascriptEscaper()
func (r *abstractCommandLineRunnerBase) getJavascriptEscaper() func(string) string {
	return jscomp.SourceCodeEscapers_JavascriptEscaper().Escape
}

// Java: AbstractCommandLineRunner.outputSingleBinary(B)
func (r *abstractCommandLineRunnerBase) outputSingleBinary(options *jscomp.CompilerOptions) {
	config := r.config
	var escaper func(string) string
	marker := defaultOutputMarker
	if strings.Contains(config.outputWrapper, outputMarkerJsString) {
		marker = outputMarkerJsString
		escaper = r.getJavascriptEscaper()
	}

	if r.isOutputInJson() {
		r.filesToStreamOut = append(r.filesToStreamOut, r.createJsonFile(options, marker, escaper))
		return
	}

	if config.jsOutputFile != "" {
		maybeCreateDirsForPath(config.jsOutputFile)
	}

	jsOutput := r.createDefaultOutput()
	r.self.WriteOutput(
		jsOutput,
		r.compiler,
		// So long as the JSChunk arg is null the compiler will write all sources to jsOutput
		// Use single-binary license tracking to dedupe licenses among all the inputs
		jscomp.NewCompiler_SingleBinaryLicenseTracker(r.compiler),
		nil,
		config.outputWrapper,
		marker,
		escaper,
		config.jsOutputFile)
	closeAppendable(jsOutput)
}

// createJsonFile saves the compiler output to a JsonFileSpec to be later
// written to stdout.
//
// Java: AbstractCommandLineRunner.createJsonFile(B, String, Function<String, String>)
func (r *abstractCommandLineRunnerBase) createJsonFile(options *jscomp.CompilerOptions, outputMarker string, escaper func(string) string) *AbstractCommandLineRunner_JsonFileSpec {
	config := r.config
	jsOutput := jutil.NewStringWriter()
	r.self.WriteOutput(
		jsOutput,
		r.compiler,
		jscomp.NewCompiler_SingleBinaryLicenseTracker(r.compiler),
		nil,
		config.outputWrapper,
		outputMarker,
		escaper,
		config.jsOutputFile)

	path := config.jsOutputFile
	if path == "" {
		path = "compiled.js"
	}
	jsonOutput := NewAbstractCommandLineRunner_JsonFileSpec(jsOutput.String(), path)

	if options.ShouldGatherSourceMapInfo() {
		smBuf := &strings.Builder{}
		if p, ok := jsonOutput.GetPath(); ok {
			r.compiler.GetSourceMap().AppendTo(smBuf, p)
		} else {
			r.compiler.GetSourceMap().AppendToOptional(smBuf, sourcemap.OptionalString{})
		}
		jsonOutput.SetSourceMap(smBuf.String())
	}

	return jsonOutput
}

func writeJsonOptStringField(w *jutil.JsonWriter, name string, opt base.OptString) {
	if v, ok := opt.Get(); ok {
		w.Name(name).ValueString(v)
	} else {
		w.Name(name).NullValue()
	}
}

// Java: AbstractCommandLineRunner.outputJsonStream()
func (r *abstractCommandLineRunnerBase) outputJsonStream() {
	jsonWriter := jutil.NewJsonWriter(newCliBufferedWriter(
		r.defaultJsOutput, nil, r.defaultJsOutput.Flush, jutil.StandardCharsets_UTF_8))
	defer jsonWriter.Close()
	jsonWriter.SetLenient(true)
	jsonWriter.SetHtmlSafe(false)
	jsonWriter.SetSerializeNulls(false)
	if r.filesToStreamOut == nil {
		jsonWriter.NullValue()
		return
	}
	jsonWriter.BeginArray()
	for _, spec := range r.filesToStreamOut {
		if spec == nil {
			jsonWriter.NullValue()
			continue
		}
		jsonWriter.BeginObject()
		writeJsonOptStringField(jsonWriter, "src", spec.Src)
		writeJsonOptStringField(jsonWriter, "path", spec.Path)
		writeJsonOptStringField(jsonWriter, "source_map", spec.SourceMap)
		writeJsonOptStringField(jsonWriter, "webpack_id", spec.WebpackId)
		jsonWriter.EndObject()
	}
	jsonWriter.EndArray()
}

// Java: AbstractCommandLineRunner.outputChunkBinaryAndSourceMaps(JSChunkGraph, B)
func (r *abstractCommandLineRunnerBase) outputChunkBinaryAndSourceMaps(
	chunkGraph *jscomp.JSChunkGraph, options *jscomp.CompilerOptions) *jscomp.DiagnosticType {
	config := r.config
	chunks := chunkGraph.GetAllChunks()
	r.parsedChunkWrappers = parseChunkWrappers(config.chunkWrapper, chunks)
	if !r.isOutputInJson() {
		// make sure the method generates all dirs up to the latest /
		maybeCreateDirsForPath(config.chunkOutputPathPrefix + "dummy")
	}

	// When the json_streams flag is specified, sourcemaps are always generated
	// per chunk.
	if !(r.shouldGenerateMapPerChunk(options) ||
		!options.ShouldGatherSourceMapInfo() ||
		config.jsonStreamMode == jscomp.CompilerOptions_JsonStreamMode_OUT ||
		config.jsonStreamMode == jscomp.CompilerOptions_JsonStreamMode_BOTH) {
		return AbstractCommandLineRunner_INVALID_CHUNK_SOURCEMAP_PATTERN
	}

	_, isStandardRunner := r.self.(*CommandLineRunner)
	if r.isOutputInJson() &&
		isStandardRunner &&
		!r.testMode &&
		r.compiler.GetSourceMap() != nil &&
		options.GetOutputJs() == jscomp.CompilerOptions_OutputJs_NORMAL &&
		options.GetNumParallelThreads() > 1 &&
		len(chunks) > 1 {
		r.outputJsonChunksParallel(chunks, options)
		return nil
	}

	mlicenseTracker := jscomp.NewCompiler_ChunkGraphAwareLicenseTracker(r.compiler)
	if !r.isOutputInJson() &&
		isStandardRunner &&
		!r.testMode &&
		(!options.ShouldGatherSourceMapInfo() || r.shouldGenerateMapPerChunk(options)) &&
		options.GetOutputJs() == jscomp.CompilerOptions_OutputJs_NORMAL &&
		options.GetNumParallelThreads() > 1 &&
		len(chunks) > 1 {
		r.outputChunksParallel(chunks, options, mlicenseTracker)
		return nil
	}

	var mapFileOut *cliBufferedWriter
	for _, m := range chunks {
		if m.GetName() == jscomp.JSChunk_WEAK_CHUNK_NAME {
			// Skip the weak chunk, which is always empty.
			continue
		}
		if r.isOutputInJson() {
			r.filesToStreamOut = append(r.filesToStreamOut, r.createJsonFileFromChunk(m))
			continue
		}

		if r.shouldGenerateMapPerChunk(options) {
			mapFileOut = r.fileNameToOutputWriter2(r.ExpandSourceMapPath(options, m))
		}

		chunkFilename := r.GetChunkOutputFileName(m)
		maybeCreateDirsForPath(chunkFilename)
		func() {
			writer := r.fileNameToLegacyOutputWriter(chunkFilename)
			defer closeBufferedWriter(writer)
			if options.ShouldGatherSourceMapInfo() {
				r.compiler.ResetAndInitializeSourceMap()
			}
			mlicenseTracker.SetCurrentChunkContext(m)
			r.self.WriteChunkOutput(chunkFilename, writer, mlicenseTracker, m)
			if options.ShouldGatherSourceMapInfo() {
				r.compiler.GetSourceMap().AppendTo(mapFileOut, chunkFilename)
			}
		}()

		if r.shouldGenerateMapPerChunk(options) && mapFileOut != nil {
			mapFileOut.Close()
			mapFileOut = nil
		}
	}

	if mapFileOut != nil {
		mapFileOut.Close()
	}
	return nil
}

func (r *abstractCommandLineRunnerBase) outputChunksParallel(
	chunks []*jscomp.JSChunk,
	options *jscomp.CompilerOptions,
	mlicenseTracker *jscomp.Compiler_ChunkGraphAwareLicenseTracker) {
	hasActive := false
	for _, m := range chunks {
		if m.GetName() != jscomp.JSChunk_WEAK_CHUNK_NAME {
			hasActive = true
			break
		}
	}
	if !hasActive {
		return
	}
	if r.parsedChunkOutputFiles == nil {
		r.parsedChunkOutputFiles = parseChunkOutputFiles(r.config.chunkOutputFiles)
	}
	createdDirs := make(map[string]struct{}, 16)
	r.compiler.OutputChunksParallel(chunks, options, mlicenseTracker, jscomp.ParallelChunkOutputCallbacks{
		ChunkFileName: func(m *jscomp.JSChunk) string {
			return r.GetChunkOutputFileName(m)
		},
		SourceMapPath: func(m *jscomp.JSChunk) string {
			if !options.ShouldGatherSourceMapInfo() {
				return ""
			}
			return r.ExpandSourceMapPath(options, m)
		},
		WrapperTemplate: func(chunkName string) string {
			return r.parsedChunkWrappers[chunkName]
		},
		EnsureDir: func(filename string) {
			if filename == "" {
				return
			}
			dir := ""
			if idx := strings.LastIndexAny(filename, "/\\"); idx >= 0 {
				dir = filename[:idx]
			}
			if _, ok := createdDirs[dir]; !ok {
				createdDirs[dir] = struct{}{}
				maybeCreateDirsForPath(filename)
			}
		},
		WriteChunkWithMap: func(filename string, mapFilename string, sm *jscomp.SourceMap, code string, wrapper string) {
			var mapFileOut *cliBufferedWriter
			if mapFilename != "" {
				mapFileOut = r.fileNameToOutputWriter2(mapFilename)
				defer closeBufferedWriter(mapFileOut)
			}
			writer := r.fileNameToLegacyOutputWriter(filename)
			defer closeBufferedWriter(writer)
			r.writeOutputOfCodeWithSourceMap(writer, sm, code, wrapper, "%s", nil)
			if mapFileOut != nil {
				sm.AppendTo(mapFileOut, filename)
			}
		},
		WriteChunkNoMap: func(filename string, code string, wrapper string) {
			writer := r.fileNameToLegacyOutputWriter(filename)
			defer closeBufferedWriter(writer)
			r.writeOutputOfCode(writer, r.compiler, code, wrapper, "%s", nil, filename)
		},
	})
}

type jsonChunkMeta struct {
	filename string
	prefix   string
	suffix   string
}

func formatJsonChunkSrc(code, prefix, suffix string, includeTrailingNewline bool) string {
	if prefix == "" && suffix == "" && !includeTrailingNewline {
		return code
	}
	extra := 0
	if includeTrailingNewline {
		extra = 1
	}
	var sb strings.Builder
	sb.Grow(len(prefix) + len(code) + len(suffix) + extra)
	sb.WriteString(prefix)
	sb.WriteString(code)
	sb.WriteString(suffix)
	if includeTrailingNewline {
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (r *abstractCommandLineRunnerBase) outputJsonChunksParallel(
	chunks []*jscomp.JSChunk,
	options *jscomp.CompilerOptions,
) {
	activeCount := 0
	for _, m := range chunks {
		if m != nil && m.GetName() != jscomp.JSChunk_WEAK_CHUNK_NAME {
			activeCount++
		}
	}
	if activeCount == 0 {
		return
	}
	if r.parsedChunkOutputFiles == nil {
		r.parsedChunkOutputFiles = parseChunkOutputFiles(r.config.chunkOutputFiles)
	}
	metaByChunk := make(map[*jscomp.JSChunk]jsonChunkMeta, activeCount)
	metas := make([]jsonChunkMeta, 0, activeCount)
	for _, m := range chunks {
		if m == nil || m.GetName() == jscomp.JSChunk_WEAK_CHUNK_NAME {
			continue
		}
		filename := r.GetChunkOutputFileName(m)
		wrapper := strings.ReplaceAll(
			r.parsedChunkWrappers[m.GetName()],
			"%basename%",
			chunkWrapperBaseName(filename))
		meta := jsonChunkMeta{filename: filename}
		if prefix, suffix, found := strings.Cut(wrapper, "%s"); found {
			meta.prefix = prefix
			meta.suffix = suffix
		}
		metaByChunk[m] = meta
		metas = append(metas, meta)
	}

	rendered, err := r.compiler.RenderChunksAndSourceMaps(chunks, jscomp.RenderChunkOptions{
		PerChunkScriptLicenses: true,
		GenerateSourceMaps:     true,
		ChunkFileName: func(m *jscomp.JSChunk) string {
			return metaByChunk[m].filename
		},
		WrapperPrefix: func(m *jscomp.JSChunk) string {
			return metaByChunk[m].prefix
		},
		NumWorkers: options.GetNumParallelThreads(),
	})
	if err != nil {
		panic(err)
	}

	includeTrailingNewline := r.config.includeTrailingNewline
	if cap(r.filesToStreamOut)-len(r.filesToStreamOut) < len(rendered) {
		grown := make([]*AbstractCommandLineRunner_JsonFileSpec, len(r.filesToStreamOut), len(r.filesToStreamOut)+len(rendered))
		copy(grown, r.filesToStreamOut)
		r.filesToStreamOut = grown
	}
	for i, rc := range rendered {
		meta := &metas[i]
		src := formatJsonChunkSrc(rc.Code, meta.prefix, meta.suffix, includeTrailingNewline)
		r.filesToStreamOut = append(
			r.filesToStreamOut,
			NewAbstractCommandLineRunner_JsonFileSpecWithSourceMap(src, meta.filename, rc.SourceMap))
	}
}

// createJsonFileFromChunk converts the given output chunk to a JsonFileSpec
// with associated sourcemap.
//
// Java: AbstractCommandLineRunner.createJsonFileFromChunk(JSChunk)
func (r *abstractCommandLineRunnerBase) createJsonFileFromChunk(chunk *jscomp.JSChunk) *AbstractCommandLineRunner_JsonFileSpec {
	r.compiler.ResetAndInitializeSourceMap()

	filename := r.GetChunkOutputFileName(chunk)
	output := jutil.NewStringWriter()
	r.self.WriteChunkOutput(filename, output, jscomp.NewCompiler_ScriptNodeLicensesOnlyTracker(r.compiler), chunk)

	jsonFile := NewAbstractCommandLineRunner_JsonFileSpec(output.String(), filename)

	var chunkSourceMap strings.Builder
	r.compiler.GetSourceMap().AppendTo(&chunkSourceMap, filename)
	jsonFile.SetSourceMap(chunkSourceMap.String())

	return jsonFile
}
