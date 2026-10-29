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
	"io"
	"slices"
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
		abstractCommandLineRunner_maybeCreateDirsForPath(&fileName)
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
	pos := jutil.StringIndexOf(wrapper, codePlaceholder)
	if pos != -1 {
		prefix := ""

		if pos > 0 {
			prefix = jutil.StringSubstring(wrapper, 0, pos)
			out.AppendCharSequence(prefix)
		}

		if escaper == nil {
			out.AppendCharSequence(code)
		} else {
			out.AppendCharSequence(escaper(code))
		}

		suffixStart := pos + jutil.StringLength(codePlaceholder)
		if suffixStart != jutil.StringLength(wrapper) {
			// Something after placeholder?
			out.AppendCharSequence(jutil.StringSubstringFrom(wrapper, suffixStart))
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
func abstractCommandLineRunner_maybeCreateDirsForPath(pathPrefix *string) {
	if pathPrefix != nil && *pathPrefix != "" {
		parent := jutil.NewFile(*pathPrefix).GetParentFile()
		if parent != nil {
			parent.Mkdirs()
		}
	}
}

// Java: AbstractCommandLineRunner.createDefaultOutput()
func (r *AbstractCommandLineRunnerBase) CreateDefaultOutput() AbstractCommandLineRunner_Writer {
	writeOutputToFile := r.config.jsOutputFile != ""
	if writeOutputToFile {
		return r.FileNameToLegacyOutputWriter(&r.config.jsOutputFile)
	}
	return r.StreamToLegacyOutputWriter(abstractCommandLineRunner_printStreamAsOutputStream{r.defaultJsOutput})
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
	return jutil.Escaper_AsFunction(jscomp.SourceCodeEscapers_JavascriptEscaper())
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
			abstractCommandLineRunner_maybeCreateDirsForPath(&config.jsOutputFile)
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
	jsOutput := jutil.NewStringBuilder()
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
	jsonWriter := jutil.NewJsonWriter(jutil.NewBufferedWriter(jutil.NewOutputStreamWriterWithCharset(
		abstractCommandLineRunner_printStreamAsOutputStream{r.defaultJsOutput}, jutil.StandardCharsets_UTF_8)))
	defer jsonWriter.Close()
	jsonWriter.SetLenient(true)
	jsonWriter.SetHtmlSafe(false)
	jsonWriter.SetSerializeNulls(false)
	jutil.Try(func() {
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
	}, func(e jutil.IIOException) {
		panic(jutil.NewJsonIOExceptionWithCause(e))
	})
}

// Java: AbstractCommandLineRunner.outputChunkBinaryAndSourceMaps(JSChunkGraph, B)
func (r *AbstractCommandLineRunnerBase) OutputChunkBinaryAndSourceMaps(
	chunkGraph *jscomp.JSChunkGraph, options *jscomp.CompilerOptions) *jscomp.DiagnosticType {
	config := r.config
	chunks := chunkGraph.GetAllChunks()
	r.parsedChunkWrappers = AbstractCommandLineRunner_ParseChunkWrappers(config.chunkWrapper, chunks)
	if !r.IsOutputInJson() {
		// make sure the method generates all dirs up to the latest /
		dummy := config.chunkOutputPathPrefix + "dummy"
		abstractCommandLineRunner_maybeCreateDirsForPath(&dummy)
	}

	// If the source map path is in fact a pattern for each
	// chunk, create a stream per-chunk. Otherwise, create
	// a single source map.
	var mapFileOut AbstractCommandLineRunner_Writer

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
			abstractCommandLineRunner_maybeCreateDirsForPath(&chunkFilename)
			func() {
				writer := r.FileNameToLegacyOutputWriter(&chunkFilename)
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
		SourceMapPath: func(m *jscomp.JSChunk) *string {
			return r.self.ExpandSourceMapPath(options, m)
		},
		WrapperTemplate: func(chunkName string) string {
			return r.parsedChunkWrappers[chunkName]
		},
		EnsureDir: func(filename *string) {
			if filename == nil || *filename == "" {
				return
			}
			dir := ""
			if idx := strings.LastIndexAny(*filename, "/\\"); idx >= 0 {
				dir = (*filename)[:idx]
			}
			if _, ok := createdDirs[dir]; !ok {
				createdDirs[dir] = struct{}{}
				abstractCommandLineRunner_maybeCreateDirsForPath(filename)
			}
		},
		WriteChunkWithMap: func(filename string, mapFilename string, sm *jscomp.SourceMap, code string, wrapper string) {
			var mapFileOut AbstractCommandLineRunner_Writer
			if mapFilename != "" {
				mapFileOut = r.FileNameToOutputWriter2(&mapFilename)
				if mapFileOut != nil {
					defer mapFileOut.Close()
				}
			}
			writer := r.FileNameToLegacyOutputWriter(&filename)
			defer abstractCommandLineRunner_closeWriter(writer)
			r.writeOutputOfCodeWithSourceMap(writer, sm, code, wrapper, "%s", nil)
			if mapFileOut != nil {
				sm.AppendTo(mapFileOut, &filename)
			}
		},
		WriteChunkNoMap: func(filename string, code string, wrapper string) {
			writer := r.FileNameToLegacyOutputWriter(&filename)
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
	output := jutil.NewStringBuilder()
	r.self.WriteChunkOutput(filename, output, jscomp.NewCompiler_ScriptNodeLicensesOnlyTracker(r.compiler), chunk)

	jsonFile := NewAbstractCommandLineRunner_JsonFileSpec(output.String(), filename)

	chunkSourceMap := &strings.Builder{} // Go-only: language-adapter — StringBuilder as SourceMap.appendTo's Appendable

	chunkOutputFileName := r.self.GetChunkOutputFileName(chunk)
	r.compiler.GetSourceMap().AppendTo(chunkSourceMap, &chunkOutputFileName)

	jsonFile.SetSourceMap(chunkSourceMap.String())

	return jsonFile
}

// ShouldGenerateMapPerChunk returns true if and only if a source map file
// should be generated for each chunk, as opposed to one unified map. This is
// specified by having the source map pattern include the %outname%
// variable.
//
// Java: AbstractCommandLineRunner.shouldGenerateMapPerChunk(B)
func (r *AbstractCommandLineRunnerBase) ShouldGenerateMapPerChunk(options *jscomp.CompilerOptions) bool {
	return options.ShouldGatherSourceMapInfo() &&
		strings.Contains(*options.GetSourceMapOutputPath(), "%outname%")
}

// OpenExternExportsStream returns a stream for outputting the generated
// externs file.
//
// options are the options to the Compiler; path is the path of the
// generated JS source file. It returns the stream or null if no extern-ed
// exports are being generated.
//
// Java: AbstractCommandLineRunner.openExternExportsStream(B, String)
func (r *AbstractCommandLineRunnerBase) OpenExternExportsStream(options *jscomp.CompilerOptions, path string) AbstractCommandLineRunner_Writer {
	externExportsPath := options.GetExternExportsPath()
	if externExportsPath == nil {
		return nil
	}

	exPath := *externExportsPath

	if !strings.Contains(exPath, jutil.File_separator) {
		outputFile := jutil.NewFile(path)
		parent := "null" // Java: a null parent concatenates as "null".
		if p := outputFile.GetParent(); p != nil {
			parent = *p
		}
		exPath = parent + string(rune(jutil.File_separatorChar)) + exPath
	}

	return r.FileNameToOutputWriter2(&exPath)
}

// ExpandCommandLinePath expands a file path specified on the command-line.
//
// Most file paths on the command-line allow an %outname% placeholder. The
// placeholder will expand to a different value depending on the current
// output mode. There are three scenarios:
//
// 1) Single JS output, single extra output: sub in jsOutputPath. 2) Multiple
// JS output, single extra output: sub in the base chunk name. 3) Multiple JS
// output, multiple extra output: sub in the module output file.
//
// Passing a JSChunk to this function automatically triggers case #3.
// Otherwise, we'll use strategy #1 or #2 based on the current output mode.
//
// Java: AbstractCommandLineRunner.expandCommandLinePath(String, JSChunk)
func (r *AbstractCommandLineRunnerBase) ExpandCommandLinePath(path string, forChunk *jscomp.JSChunk) string {
	var sub string
	if forChunk != nil {
		sub = r.self.GetChunkOutputFileName(forChunk)
	} else if len(r.config.chunk) != 0 {
		sub = r.config.chunkOutputPathPrefix
	} else {
		sub = r.config.jsOutputFile
	}
	return strings.ReplaceAll(path, "%outname%", sub)
}

// ExpandSourceMapPath is the expansion function for source map.
//
// Java: AbstractCommandLineRunner.expandSourceMapPath(B, JSChunk)
func (r *AbstractCommandLineRunnerBase) ExpandSourceMapPath(options *jscomp.CompilerOptions, forChunk *jscomp.JSChunk) *string {
	if !options.ShouldGatherSourceMapInfo() {
		return nil
	}
	expanded := r.ExpandCommandLinePath(*options.GetSourceMapOutputPath(), forChunk)
	return &expanded
}

// FileNameToLegacyOutputWriter converts a file name into a Writer taking in
// account the output charset. Returns null if the file name is null.
//
// Java: AbstractCommandLineRunner.fileNameToLegacyOutputWriter(String)
func (r *AbstractCommandLineRunnerBase) FileNameToLegacyOutputWriter(fileName *string) AbstractCommandLineRunner_Writer {
	if fileName == nil {
		return nil
	}
	if r.self.IsInTestMode() {
		return jutil.NewStringWriter()
	}

	return r.StreamToLegacyOutputWriter(r.self.FilenameToOutputStream(fileName))
}

// FileNameToOutputWriter2 converts a file name into a Writer taking in
// account the output charset. Returns null if the file name is null.
//
// Java: AbstractCommandLineRunner.fileNameToOutputWriter2(String)
func (r *AbstractCommandLineRunnerBase) FileNameToOutputWriter2(fileName *string) AbstractCommandLineRunner_Writer {
	if fileName == nil {
		return nil
	}
	if r.self.IsInTestMode() {
		return jutil.NewStringWriter()
	}

	return r.StreamToOutputWriter2(r.self.FilenameToOutputStream(fileName))
}

// FilenameToOutputStream converts a file name into a Outputstream. Returns
// null if the file name is null.
//
// Java: AbstractCommandLineRunner.filenameToOutputStream(String)
func (r *AbstractCommandLineRunnerBase) FilenameToOutputStream(fileName *string) io.WriteCloser {
	if fileName == nil {
		return nil
	}
	return jutil.NewFileOutputStream(*fileName)
}

// StreamToLegacyOutputWriter creates a writer with the legacy output
// charset.
//
// Java: AbstractCommandLineRunner.streamToLegacyOutputWriter(OutputStream)
func (r *AbstractCommandLineRunnerBase) StreamToLegacyOutputWriter(stream io.WriteCloser) AbstractCommandLineRunner_Writer {
	if r.legacyOutputCharset == nil {
		return r.self.CreateWriter(stream, jutil.StandardCharsets_UTF_8)
	}
	return r.self.CreateWriter(stream, r.legacyOutputCharset)
}

// StreamToOutputWriter2 creates a writer with the newer output charset.
//
// Java: AbstractCommandLineRunner.streamToOutputWriter2(OutputStream)
func (r *AbstractCommandLineRunnerBase) StreamToOutputWriter2(stream io.WriteCloser) AbstractCommandLineRunner_Writer {
	if r.outputCharset2 == nil {
		return r.self.CreateWriter(stream, jutil.StandardCharsets_UTF_8)
	}
	return r.self.CreateWriter(stream, r.outputCharset2)
}

// CreateWriter creates a buffered Writer that writes to the given stream
// using the given encoding.
//
// Java: AbstractCommandLineRunner.createWriter(OutputStream, Charset)
func (r *AbstractCommandLineRunnerBase) CreateWriter(stream io.WriteCloser, charset *jutil.Charset) AbstractCommandLineRunner_Writer {
	return jutil.NewBufferedWriter(jutil.NewOutputStreamWriterWithCharset(stream, charset))
}

// OutputSourceMap outputs the source map found in the compiler to the
// proper path if one exists.
//
// Java: AbstractCommandLineRunner.outputSourceMap(B, String)
func (r *AbstractCommandLineRunnerBase) OutputSourceMap(options *jscomp.CompilerOptions, associatedName string) {
	if !options.ShouldGatherSourceMapInfo() ||
		*options.GetSourceMapOutputPath() == "/dev/null" {
		return
	}

	outName := r.self.ExpandSourceMapPath(options, nil)
	abstractCommandLineRunner_maybeCreateDirsForPath(outName)
	out := r.FileNameToOutputWriter2(outName)
	defer abstractCommandLineRunner_closeWriter(out)
	r.compiler.GetSourceMap().AppendTo(out, &associatedName)
}

// GetMapPath returns the path at which to output map file(s) based on the
// path at which the JS binary will be placed. It returns the path in which
// to place the generated map file(s).
//
// Java: AbstractCommandLineRunner.getMapPath(String)
func (r *AbstractCommandLineRunnerBase) GetMapPath(outputFile string) string {
	basePath := ""

	if outputFile == "" {
		// If we have a js_module_binary rule, output the maps
		// at chunkname_props_map.out, etc.
		if r.config.chunkOutputPathPrefix != "" {
			basePath = r.config.chunkOutputPathPrefix
		} else {
			basePath = "jscompiler"
		}
	} else {
		// Get the path of the output file.
		file := jutil.NewFile(outputFile)

		outputFileName := file.GetName()

		// Strip the .js from the name.
		if strings.HasSuffix(outputFileName, ".js") {
			outputFileName = outputFileName[:len(outputFileName)-3]
		}

		fileParent := file.GetParent()
		if fileParent == nil {
			basePath = outputFileName
		} else {
			basePath = *file.GetParent() + string(rune(jutil.File_separatorChar)) + outputFileName
		}
	}

	return basePath
}

// OutputNameMaps outputs the variable and property name maps for the
// specified compiler if the proper FLAGS are set.
//
// Java: AbstractCommandLineRunner.outputNameMaps()
func (r *AbstractCommandLineRunnerBase) OutputNameMaps() {
	config := r.config

	var propertyMapOutputPath *string
	var variableMapOutputPath *string

	// Check the create_name_map_files FLAG.
	if config.createNameMapFiles {
		basePath := r.GetMapPath(config.jsOutputFile)

		propertyMapOutputPath = jutil.StringPtr(basePath + "_props_map.out")
		variableMapOutputPath = jutil.StringPtr(basePath + "_vars_map.out")
	}

	// Check the individual FLAGS.
	if config.variableMapOutputFile != "" {
		if variableMapOutputPath != nil {
			panic(jscomp.NewFlagUsageException(
				"The flags variable_map_output_file and " +
					"create_name_map_files cannot both be used simultaneously."))
		}

		variableMapOutputPath = &config.variableMapOutputFile
	}

	if config.propertyMapOutputFile != "" {
		if propertyMapOutputPath != nil {
			panic(jscomp.NewFlagUsageException(
				"The flags property_map_output_file and " +
					"create_name_map_files cannot both be used simultaneously."))
		}

		propertyMapOutputPath = &config.propertyMapOutputFile
	}

	// Output the maps.
	if variableMapOutputPath != nil && r.compiler.GetVariableMap() != nil {
		r.compiler.GetVariableMap().Save(*variableMapOutputPath)
	}

	if propertyMapOutputPath != nil && r.compiler.GetPropertyMap() != nil {
		r.compiler.GetPropertyMap().Save(*propertyMapOutputPath)
	}
}

// OutputStringMap outputs the string map generated by the ReplaceStrings
// pass if an output path exists.
//
// Java: AbstractCommandLineRunner.outputStringMap()
func (r *AbstractCommandLineRunnerBase) OutputStringMap() {
	config := r.config
	if config.stringMapOutputPath != "" {
		if r.compiler.GetStringMap() == nil {
			// Ensure an empty file is created if there is no string map.
			// This avoids confusing some build tools that expect to see the file, even if it is empty.
			if !jutil.NewFile(config.stringMapOutputPath).CreateNewFile() {
				panic(jutil.NewIOException("Could not create file: " + config.stringMapOutputPath))
			}
		} else {
			r.compiler.GetStringMap().Save(config.stringMapOutputPath)
		}
	}
}

// ShouldGenerateOutputPerChunk returns true if and only if a manifest or
// bundle should be generated for each chunk, as opposed to one unified
// manifest.
//
// Java: AbstractCommandLineRunner.shouldGenerateOutputPerChunk(String)
func (r *AbstractCommandLineRunnerBase) ShouldGenerateOutputPerChunk(output string) bool {
	return len(r.config.chunk) != 0 && strings.Contains(output, "%outname%")
}

// Java: AbstractCommandLineRunner.outputManifest()
func (r *AbstractCommandLineRunnerBase) OutputManifest() {
	r.OutputManifestOrBundle(r.config.outputManifests, true)
}

// Java: AbstractCommandLineRunner.outputBundle()
func (r *AbstractCommandLineRunnerBase) OutputBundle() {
	r.OutputManifestOrBundle(r.config.outputBundles, false)
}

// Java: AbstractCommandLineRunner.outputInstrumentationMapping()
func (r *AbstractCommandLineRunnerBase) OutputInstrumentationMapping() {
	if r.config.instrumentationMappingFile != "" {
		path := r.ExpandCommandLinePath(r.config.instrumentationMappingFile /* forChunk= */, nil)
		r.compiler.GetInstrumentationMapping().Save(path)
	}
}

// OutputManifestOrBundle writes the manifest or bundle of all compiler
// input files that were included as controlled by --dependency_mode, if
// requested.
//
// Java: AbstractCommandLineRunner.outputManifestOrBundle(List<String>, boolean)
func (r *AbstractCommandLineRunnerBase) OutputManifestOrBundle(outputFiles []string, isManifest bool) {
	if len(outputFiles) == 0 {
		return
	}

	for _, output := range outputFiles {
		if output == "" {
			continue
		}

		if r.ShouldGenerateOutputPerChunk(output) {
			// Generate per-chunk manifests or bundles.
			chunks := r.compiler.GetChunkGraph().GetAllChunks()
			for _, chunk := range chunks {
				func() {
					path := r.ExpandCommandLinePath(output, chunk)
					out := r.FileNameToOutputWriter2(&path)
					defer abstractCommandLineRunner_closeWriter(out)
					if isManifest {
						r.self.PrintManifestTo(chunk, out)
					} else {
						r.self.PrintBundleTo(chunk, out)
					}
				}()
			}
		} else {
			// Generate a single file manifest or bundle.
			func() {
				path := r.ExpandCommandLinePath(output, nil)
				out := r.FileNameToOutputWriter2(&path)
				defer abstractCommandLineRunner_closeWriter(out)
				if len(r.config.chunk) == 0 {
					// For a single-chunk compilation, generate a single headerless manifest or bundle
					// containing only the strong files.
					chunk := r.compiler.GetChunkGraph().GetChunkByName(jscomp.JSChunk_STRONG_CHUNK_NAME)
					if isManifest {
						r.self.PrintManifestTo(chunk, out)
					} else {
						r.self.PrintBundleTo(chunk, out)
					}
				} else {
					// For a multi-chunk compilation, generate a single manifest file with chunk headers.
					r.self.PrintChunkGraphManifestOrBundleTo(r.compiler.GetChunkGraph(), out, isManifest)
				}
			}()
		}
	}
}

// OutputChunkGraphJson creates a file containing the current chunk graph in
// JSON serialization.
//
// Java: AbstractCommandLineRunner.outputChunkGraphJson()
func (r *AbstractCommandLineRunnerBase) OutputChunkGraphJson() {
	if r.config.outputChunkDependencies != nil && jutil.StringLength(*r.config.outputChunkDependencies) != 0 {
		out := r.FileNameToOutputWriter2(r.config.outputChunkDependencies)
		defer abstractCommandLineRunner_closeWriter(out)
		r.self.PrintChunkGraphJsonTo(out)
	}
}

// PrintChunkGraphJsonTo prints the current chunk graph as JSON.
//
// Java: AbstractCommandLineRunner.printChunkGraphJsonTo(Appendable)
func (r *AbstractCommandLineRunnerBase) PrintChunkGraphJsonTo(out jutil.Appendable) {
	out.AppendCharSequence(r.compiler.GetChunkGraph().ToJson().String())
}

// PrintChunkGraphManifestOrBundleTo prints a set of chunks to the manifest
// or bundle file.
//
// Java: AbstractCommandLineRunner.printChunkGraphManifestOrBundleTo(JSChunkGraph, Appendable, boolean)
func (r *AbstractCommandLineRunnerBase) PrintChunkGraphManifestOrBundleTo(graph *jscomp.JSChunkGraph, out jutil.Appendable, isManifest bool) {
	requiresNewline := false
	for _, chunk := range graph.GetAllChunks() {
		if !isManifest && chunk.IsWeak() {
			// Skip the weak chunk on a multi-chunk bundle, but not a multi-chunk manifest.
			continue
		}

		if requiresNewline {
			out.AppendCharSequence("\n")
		}

		if isManifest {
			// See CommandLineRunnerTest to see what the format of this
			// manifest looks like.
			dependencies := strings.Join(chunk.GetSortedDependencyNames(), ",")
			deps := ""
			if dependencies != "" {
				deps = ":" + dependencies
			}
			out.AppendCharSequence("{" + chunk.GetName() + deps + "}\n")
			r.self.PrintManifestTo(chunk, out)
		} else {
			r.self.PrintBundleTo(chunk, out)
		}
		requiresNewline = true
	}
}

// PrintManifestTo prints a list of input names (using root-relative paths),
// delimited by newlines, to the manifest file.
//
// Java: AbstractCommandLineRunner.printManifestTo(JSChunk, Appendable)
func (r *AbstractCommandLineRunnerBase) PrintManifestTo(chunk *jscomp.JSChunk, out jutil.Appendable) {
	for _, input := range chunk.GetInputs() {
		name := input.GetName()

		// Ignore fill files created by the compiler to facilitate cross-chunk code motion.
		if jscomp.AbstractCompiler_IsFillFileName(name) {
			continue
		}

		displayName := name
		if rootRelativePath, ok := r.rootRelativePathsMap[name]; ok {
			displayName = rootRelativePath
		}
		out.AppendCharSequence(displayName)
		out.AppendCharSequence("\n")
	}
}

// PrintBundleTo prints all the input contents, starting with a comment that
// specifies the input file name (using root-relative paths) before each
// file.
//
// Java: AbstractCommandLineRunner.printBundleTo(JSChunk, Appendable)
func (r *AbstractCommandLineRunnerBase) PrintBundleTo(chunk *jscomp.JSChunk, out jutil.Appendable) {
	inputs := chunk.GetInputs()
	options := r.compiler.GetOptions()
	// Prebuild ASTs before they're needed in getLoadFlags, for performance and because
	// StackOverflowErrors can be hit if not prebuilt.
	if options.GetNumParallelThreads() > 1 {
		jscomp.NewPrebuildAst(r.compiler, r.compiler.GetOptions().GetNumParallelThreads()).Prebuild(slices.Values(inputs))
	}
	if options.GetRuntimeLibraryMode() == jscomp.RuntimeJsLibManager_RuntimeLibraryMode_INJECT {
		// ES6 modules will need a runtime in a bundle. Skip appending this runtime if there are no
		// ES6 modules to cut down on size.
		for _, input := range inputs {
			if input.IsEs6Module() {
				r.self.AppendRuntimeTo(out)
				break
			}
		}
	}

	for _, input := range inputs {
		name := input.GetName()
		code := input.GetSourceFile().GetCode()

		// Ignore empty fill files created by the compiler to facilitate cross-chunk code motion.
		// Note that non-empty fill files (ones whose code has actually been moved into) are still
		// emitted. In particular, this ensures that if there are no (real) inputs the bundle will be
		// empty.
		if jscomp.AbstractCompiler_IsFillFileName(name) && code == "" {
			continue
		}

		displayName := input.GetName()
		if rootRelativePath, ok := r.rootRelativePathsMap[name]; ok {
			displayName = rootRelativePath
		}
		out.AppendCharSequence("//")
		out.AppendCharSequence(displayName)
		out.AppendCharSequence("\n")

		r.self.PrepForBundleAndAppendTo(out, input, code)

		out.AppendCharSequence("\n")
	}
}

// AbstractCommandLineRunner_Writer is a java.io.Writer as this class uses
// it: Writer.write/flush/close, the Appendable methods it is passed through,
// and io.StringWriter for SourceMap.appendTo.
//
// Go-only: language-adapter — java.io.Writer implements Appendable.
type AbstractCommandLineRunner_Writer interface {
	jutil.Writer
	jutil.Appendable
	io.StringWriter
}

// abstractCommandLineRunner_closeWriter is the close() of a
// try-with-resources statement, which skips a null resource.
//
// Go-only: language-adapter — try-with-resources.
func abstractCommandLineRunner_closeWriter(w AbstractCommandLineRunner_Writer) {
	if w != nil {
		w.Close()
	}
}

// abstractCommandLineRunner_printStreamAsOutputStream passes a PrintStream
// where Java passes it as an OutputStream: bytes go through
// PrintStream.write(byte[]) unchanged, and close() closes the PrintStream.
//
// Go-only: language-adapter — PrintStream is an OutputStream.
type abstractCommandLineRunner_printStreamAsOutputStream struct {
	ps *jutil.PrintStream
}

func (s abstractCommandLineRunner_printStreamAsOutputStream) Write(b []byte) (int, error) {
	return s.ps.Write(b)
}

func (s abstractCommandLineRunner_printStreamAsOutputStream) Close() error {
	s.ps.Close()
	return nil
}
