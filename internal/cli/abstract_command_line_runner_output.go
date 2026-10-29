// Port of com/google/javascript/jscomp/AbstractCommandLineRunner.java @ closure-compiler v20261021

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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// AbstractCommandLineRunner_ParseChunkWrappers parses chunk wrapper
// specifications.
//
// specs is a list of chunk wrapper specifications, not null. The spec
// format is: name:wrapper. Wrappers. chunks are the JS chunks whose
// wrappers are specified. It returns a map from chunk name to chunk wrapper.
// Chunks with no wrapper will have the empty string as their value in this
// map.
//
// Java: AbstractCommandLineRunner.parseChunkWrappers(List<String>, Iterable<JSChunk>)
func AbstractCommandLineRunner_ParseChunkWrappers(specs []string, chunks []*jscomp.JSChunk) map[string]string {
	jutil.CheckState(specs != nil)

	wrappers := make(map[string]string, len(chunks))

	// Prepopulate the map with chunk names.
	for _, c := range chunks {
		wrappers[c.GetName()] = ""
	}

	for _, spec := range specs {
		// Format is "<name>:<wrapper>".
		pos := strings.IndexByte(spec, ':')
		if pos == -1 {
			panic(jscomp.NewFlagUsageException(
				"Expected chunk wrapper to have " + "<name>:<wrapper> format: " + spec))
		}

		// Parse chunk name.
		name := spec[:pos]
		if _, ok := wrappers[name]; !ok {
			panic(jscomp.NewFlagUsageException("Unknown chunk: '" + name + "'"))
		}
		wrapper := spec[pos+1:]
		// Support for %n% and %output%
		wrapper = strings.ReplaceAll(strings.ReplaceAll(wrapper, "%output%", "%s"), "%n%", "\n")
		if !strings.Contains(wrapper, "%s") {
			panic(jscomp.NewFlagUsageException("No %s placeholder in chunk wrapper: '" + wrapper + "'"))
		}

		wrappers[name] = wrapper
	}
	return wrappers
}

// abstractCommandLineRunner_parseChunkOutputFiles parses chunk output name
// specifications.
//
// specs is a list of chunk output name specifications, not null. The spec
// format is: name:output_file. It returns a map from chunk name to chunk
// output file name, if declared. Chunks with no predeclared output file name
// will have no entry in this map.
//
// Java: AbstractCommandLineRunner.parseChunkOutputFiles(List<String>)
func abstractCommandLineRunner_parseChunkOutputFiles(specs []string) map[string]string {
	jutil.CheckArgument(specs != nil)

	if len(specs) == 0 {
		return map[string]string{}
	}

	entries := make([]PipeFlagEntry, 0, len(specs))
	for _, spec := range specs {
		// Format is "<name>:<output_file>".
		pos := strings.IndexByte(spec, ':')
		if pos == -1 {
			panic(jscomp.NewFlagUsageException(
				"Expected chunk_output_file to have " + "<name>:<output_file> format: " + spec))
		}
		entries = append(entries, PipeFlagEntry{Key: spec[:pos], Value: spec[pos+1:]})
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
func (r *AbstractCommandLineRunnerBase) GetChunkOutputFileName(m *jscomp.JSChunk) string {
	if r.parsedChunkOutputFiles == nil {
		r.parsedChunkOutputFiles = abstractCommandLineRunner_parseChunkOutputFiles(r.config.chunkOutputFiles)
	}
	if outputFile, ok := r.parsedChunkOutputFiles[m.GetName()]; ok {
		return r.config.chunkOutputPathPrefix + outputFile
	}
	return r.config.chunkOutputPathPrefix + m.GetName() + ".js"
}

// GetChunkConformanceFileName returns the conformance file name for a
// chunk.
//
// Java: AbstractCommandLineRunner.getChunkConformanceFileName(JSChunk)
func (r *AbstractCommandLineRunnerBase) GetChunkConformanceFileName(m *jscomp.JSChunk) *string {
	if r.parsedChunkConformanceFiles == nil {
		r.parsedChunkConformanceFiles = abstractCommandLineRunner_parseChunkOutputFiles(r.config.chunkConformanceFiles)
	}
	if name, ok := r.parsedChunkConformanceFiles[m.GetName()]; ok {
		return &name
	}
	return nil
}

// Java: AbstractCommandLineRunner.writeChunkOutput(String, Appendable, LicenseTracker, JSChunk)
func (r *AbstractCommandLineRunnerBase) WriteChunkOutput(fileName string, out jutil.Appendable, lt jscomp.CodePrinter_LicenseTracker, m *jscomp.JSChunk) {
	if r.parsedChunkWrappers == nil {
		r.parsedChunkWrappers = AbstractCommandLineRunner_ParseChunkWrappers(
			r.config.chunkWrapper, append([]*jscomp.JSChunk{}, r.compiler.GetChunkGraph().GetAllChunks()...))
	}

	if !r.IsOutputInJson() {
		abstractCommandLineRunner_maybeCreateDirsForPath(fileName)
	}
	baseName := jutil.NewFile(fileName).GetName()
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
// inserted.
//
// chunk is which chunk to write. If this is null, write the entire AST.
//
// Java: AbstractCommandLineRunner.writeOutput(Appendable, Compiler, LicenseTracker, JSChunk, String, String, Function<String, String>, String)
func (r *AbstractCommandLineRunnerBase) WriteOutput(
	out jutil.Appendable,
	compiler *jscomp.Compiler,
	licenseTracker jscomp.CodePrinter_LicenseTracker,
	chunk *jscomp.JSChunk,
	wrapper string,
	codePlaceholder string,
	escaper func(string) string,
	filename string) {
	if compiler.GetOptions().GetOutputJs() == jscomp.CompilerOptions_OutputJs_SENTINEL {
		out.AppendCharSequence("// No JS output because the compiler was run in checks-only mode.\n")
		return
	}
	jutil.CheckState(compiler.GetOptions().GetOutputJs() == jscomp.CompilerOptions_OutputJs_NORMAL)

	var code string
	if chunk == nil {
		code = compiler.ToSource()
	} else {
		code = compiler.ToSourceOfLicenseTrackerAndJSChunk(licenseTracker, chunk)
	}
	r.self.WriteOutputOfCode(out, compiler, code, wrapper, codePlaceholder, escaper, filename)
}

// WriteOutputOfCode writes code to an output stream, optionally wrapping it
// in an arbitrary wrapper that contains a placeholder where the code should
// be inserted.
//
// Java: AbstractCommandLineRunner.writeOutput(Appendable, Compiler, String, String, String, Function<String, String>, String)
func (r *AbstractCommandLineRunnerBase) WriteOutputOfCode(
	out jutil.Appendable,
	compiler *jscomp.Compiler,
	code string,
	wrapper string,
	codePlaceholder string,
	escaper func(string) string,
	filename string) {
	var sm *jscomp.SourceMap
	if compiler != nil {
		sm = compiler.GetSourceMap()
	}
	r.writeOutputOfCodeWithSourceMap(out, sm, code, wrapper, codePlaceholder, escaper)
}

func (r *AbstractCommandLineRunnerBase) writeOutputOfCodeWithSourceMap(
	out jutil.Appendable,
	sm *jscomp.SourceMap,
	code string,
	wrapper string,
	codePlaceholder string,
	escaper func(string) string) {
	pos := strings.Index(wrapper, codePlaceholder)
	if pos != -1 {
		prefix := ""

		if pos > 0 {
			prefix = wrapper[:pos]
			out.AppendCharSequence(prefix)
		}

		if escaper == nil {
			out.AppendCharSequence(code)
		} else {
			out.AppendCharSequence(escaper(code))
		}

		suffixStart := pos + len(codePlaceholder)
		if suffixStart != len(wrapper) {
			// Something after placeholder?
			out.AppendCharSequence(wrapper[suffixStart:])
		}
		if r.self.GetCommandLineConfig().includeTrailingNewline {
			out.AppendCharacter('\n')
		}

		// If we have a source map, adjust its offsets to match
		// the code WITHIN the wrapper.
		if sm != nil {
			sm.SetWrapperPrefix(prefix)
		}
	} else {
		out.AppendCharSequence(code)
		if r.self.GetCommandLineConfig().includeTrailingNewline {
			out.AppendCharacter('\n')
		}
	}
}

// abstractCommandLineRunner_maybeCreateDirsForPath creates any directories
// necessary to write a file that will have a given path prefix.
//
// Java: AbstractCommandLineRunner.maybeCreateDirsForPath(String)
func abstractCommandLineRunner_maybeCreateDirsForPath(pathPrefix string) {
	if pathPrefix != "" {
		parent := jutil.NewFile(pathPrefix).GetParentFile()
		if parent != nil {
			parent.Mkdirs()
		}
	}
}

// Java: AbstractCommandLineRunner.createDefaultOutput()
func (r *AbstractCommandLineRunnerBase) CreateDefaultOutput() *cliBufferedWriter {
	if r.config.jsOutputFile != "" {
		return r.FileNameToLegacyOutputWriter(r.config.jsOutputFile)
	}
	cs := r.legacyOutputCharset
	if cs == nil {
		cs = jutil.StandardCharsets_UTF_8
	}
	return newCliBufferedWriter(r.defaultJsOutput, nil, r.defaultJsOutput.Flush, cs)
}

// Java: AbstractCommandLineRunner.closeAppendable(Appendable)
func abstractCommandLineRunner_closeAppendable(output jutil.Appendable) {
	if flushable, ok := output.(interface{ Flush() }); ok {
		flushable.Flush()
	}
	if closeable, ok := output.(interface{ Close() }); ok {
		closeable.Close()
	}
}

// Java: AbstractCommandLineRunner.getJavascriptEscaper()
func (r *AbstractCommandLineRunnerBase) GetJavascriptEscaper() func(string) string {
	return jscomp.SourceCodeEscapers_JavascriptEscaper().Escape
}

// Java: AbstractCommandLineRunner.outputSingleBinary(B)
func (r *AbstractCommandLineRunnerBase) OutputSingleBinary(options *jscomp.CompilerOptions) {
	config := r.config
	var escaper func(string) string
	marker := AbstractCommandLineRunner_OUTPUT_MARKER
	if strings.Contains(config.outputWrapper, abstractCommandLineRunner_OUTPUT_MARKER_JS_STRING) {
		marker = abstractCommandLineRunner_OUTPUT_MARKER_JS_STRING
		escaper = r.self.GetJavascriptEscaper()
	}

	if r.IsOutputInJson() {
		r.filesToStreamOut = append(r.filesToStreamOut, r.self.CreateJsonFile(options, marker, escaper))
	} else {
		if config.jsOutputFile != "" {
			abstractCommandLineRunner_maybeCreateDirsForPath(config.jsOutputFile)
		}

		jsOutput := r.CreateDefaultOutput()
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
		abstractCommandLineRunner_closeAppendable(jsOutput)
	}
}

// CreateJsonFile saves the compiler output to a JsonFileSpec to be later
// written to stdout.
//
// Java: AbstractCommandLineRunner.createJsonFile(B, String, Function<String, String>)
func (r *AbstractCommandLineRunnerBase) CreateJsonFile(options *jscomp.CompilerOptions, outputMarker string, escaper func(string) string) *AbstractCommandLineRunner_JsonFileSpec {
	config := r.config
	jsOutput := jutil.NewStringWriter()
	r.self.WriteOutput(
		jsOutput,
		r.compiler,
		// So long as the JSChunk arg is null the compiler will write all sources to jsOutput
		// Use single-binary license tracking to dedupe licenses among all the inputs
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
		sourcemap := &strings.Builder{} // Go-only: language-adapter — StringBuilder as SourceMap.appendTo's Appendable
		r.compiler.GetSourceMap().AppendTo(sourcemap, jsonOutput.GetPath())
		jsonOutput.SetSourceMap(sourcemap.String())
	}

	return jsonOutput
}

// Java: AbstractCommandLineRunner.outputJsonStream()
func (r *AbstractCommandLineRunnerBase) OutputJsonStream() {
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
		jsonWriter.Name("src").ValueString(spec.Src)
		jsonWriter.Name("path").ValueString(spec.Path)
		jsonWriter.Name("source_map").ValueString(spec.SourceMap)
		jsonWriter.Name("webpack_id").ValueString(spec.WebpackId)
		jsonWriter.EndObject()
	}
	jsonWriter.EndArray()
}

// Java: AbstractCommandLineRunner.outputChunkBinaryAndSourceMaps(JSChunkGraph, B)
func (r *AbstractCommandLineRunnerBase) OutputChunkBinaryAndSourceMaps(
	chunkGraph *jscomp.JSChunkGraph, options *jscomp.CompilerOptions) *jscomp.DiagnosticType {
	config := r.config
	chunks := chunkGraph.GetAllChunks()
	r.parsedChunkWrappers = AbstractCommandLineRunner_ParseChunkWrappers(config.chunkWrapper, chunks)
	if !r.IsOutputInJson() {
		// make sure the method generates all dirs up to the latest /
		abstractCommandLineRunner_maybeCreateDirsForPath(config.chunkOutputPathPrefix + "dummy")
	}

	// If the source map path is in fact a pattern for each
	// chunk, create a stream per-chunk. Otherwise, create
	// a single source map.
	var mapFileOut *cliBufferedWriter

	// When the json_streams flag is specified, sourcemaps are always generated
	// per chunk
	if !(r.self.ShouldGenerateMapPerChunk(options) ||
		!options.ShouldGatherSourceMapInfo() ||
		config.jsonStreamMode == jscomp.CompilerOptions_JsonStreamMode_OUT ||
		config.jsonStreamMode == jscomp.CompilerOptions_JsonStreamMode_BOTH) {
		// warn that this is not supported
		return AbstractCommandLineRunner_INVALID_CHUNK_SOURCEMAP_PATTERN
	}

	mlicenseTracker := jscomp.NewCompiler_ChunkGraphAwareLicenseTracker(r.compiler)
	if !r.IsOutputInJson() &&
		!r.self.IsInTestMode() &&
		(!options.ShouldGatherSourceMapInfo() || r.self.ShouldGenerateMapPerChunk(options)) &&
		options.GetOutputJs() == jscomp.CompilerOptions_OutputJs_NORMAL &&
		options.GetNumParallelThreads() > 1 &&
		len(chunks) > 1 {
		r.outputChunksParallel(chunks, options, mlicenseTracker)
		return nil
	}
	for _, m := range chunks {
		if m.GetName() == jscomp.JSChunk_WEAK_CHUNK_NAME {
			// Skip the weak chunk, which is always empty.
			continue
		}
		if r.IsOutputInJson() {
			r.filesToStreamOut = append(r.filesToStreamOut, r.CreateJsonFileFromChunk(m))
		} else {
			if r.self.ShouldGenerateMapPerChunk(options) {
				mapFileOut = r.FileNameToOutputWriter2(r.self.ExpandSourceMapPath(options, m))
			}

			chunkFilename := r.self.GetChunkOutputFileName(m)
			abstractCommandLineRunner_maybeCreateDirsForPath(chunkFilename)
			func() {
				writer := r.FileNameToLegacyOutputWriter(chunkFilename)
				defer abstractCommandLineRunner_closeWriter(writer)
				if options.ShouldGatherSourceMapInfo() {
					r.compiler.ResetAndIntitializeSourceMap()
				}
				mlicenseTracker.SetCurrentChunkContext(m)
				r.self.WriteChunkOutput(chunkFilename, writer, mlicenseTracker, m)
				if options.ShouldGatherSourceMapInfo() {
					r.compiler.GetSourceMap().AppendTo(mapFileOut, &chunkFilename)
				}
			}()

			if r.self.ShouldGenerateMapPerChunk(options) && mapFileOut != nil {
				mapFileOut.Close()
				mapFileOut = nil
			}
		}
	}

	if mapFileOut != nil {
		mapFileOut.Close()
	}
	return nil
}

func (r *AbstractCommandLineRunnerBase) outputChunksParallel(
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
		r.parsedChunkOutputFiles = abstractCommandLineRunner_parseChunkOutputFiles(r.config.chunkOutputFiles)
	}
	createdDirs := make(map[string]struct{}, 16)
	r.compiler.OutputChunksParallel(chunks, options, mlicenseTracker, jscomp.ParallelChunkOutputCallbacks{
		ChunkFileName: func(m *jscomp.JSChunk) string {
			return r.self.GetChunkOutputFileName(m)
		},
		SourceMapPath: func(m *jscomp.JSChunk) string {
			if !options.ShouldGatherSourceMapInfo() {
				return ""
			}
			return r.self.ExpandSourceMapPath(options, m)
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
				abstractCommandLineRunner_maybeCreateDirsForPath(filename)
			}
		},
		WriteChunkWithMap: func(filename string, mapFilename string, sm *jscomp.SourceMap, code string, wrapper string) {
			var mapFileOut *cliBufferedWriter
			if mapFilename != "" {
				mapFileOut = r.FileNameToOutputWriter2(mapFilename)
				if mapFileOut != nil {
					defer mapFileOut.Close()
				}
			}
			writer := r.FileNameToLegacyOutputWriter(filename)
			defer abstractCommandLineRunner_closeWriter(writer)
			r.writeOutputOfCodeWithSourceMap(writer, sm, code, wrapper, "%s", nil)
			if mapFileOut != nil {
				sm.AppendTo(mapFileOut, &filename)
			}
		},
		WriteChunkNoMap: func(filename string, code string, wrapper string) {
			writer := r.FileNameToLegacyOutputWriter(filename)
			defer abstractCommandLineRunner_closeWriter(writer)
			r.self.WriteOutputOfCode(writer, r.compiler, code, wrapper, "%s", nil, filename)
		},
	})
}

// CreateJsonFileFromChunk converts the given output chunk to a JSONFileSpec
// with associated sourcemap.
//
// Java: AbstractCommandLineRunner.createJsonFileFromChunk(JSChunk)
func (r *AbstractCommandLineRunnerBase) CreateJsonFileFromChunk(chunk *jscomp.JSChunk) *AbstractCommandLineRunner_JsonFileSpec {
	r.compiler.ResetAndIntitializeSourceMap()

	filename := r.self.GetChunkOutputFileName(chunk)
	output := jutil.NewStringWriter()
	r.self.WriteChunkOutput(filename, output, jscomp.NewCompiler_ScriptNodeLicensesOnlyTracker(r.compiler), chunk)

	jsonFile := NewAbstractCommandLineRunner_JsonFileSpec(output.String(), filename)

	chunkSourceMap := &strings.Builder{} // Go-only: language-adapter — StringBuilder as SourceMap.appendTo's Appendable

	chunkOutputFileName := r.self.GetChunkOutputFileName(chunk)
	r.compiler.GetSourceMap().AppendTo(chunkSourceMap, &chunkOutputFileName)

	jsonFile.SetSourceMap(chunkSourceMap.String())

	return jsonFile
}
