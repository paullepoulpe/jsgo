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
	"bufio"
	"errors"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

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

	return r.FileNameToOutputWriter2(exPath)
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
func (r *AbstractCommandLineRunnerBase) ExpandSourceMapPath(options *jscomp.CompilerOptions, forChunk *jscomp.JSChunk) string {
	if !options.ShouldGatherSourceMapInfo() {
		return ""
	}
	return r.ExpandCommandLinePath(*options.GetSourceMapOutputPath(), forChunk)
}

// FileNameToLegacyOutputWriter converts a file name into a Writer taking in
// account the output charset.
//
// Java: AbstractCommandLineRunner.fileNameToLegacyOutputWriter(String)
func (r *AbstractCommandLineRunnerBase) FileNameToLegacyOutputWriter(fileName string) AbstractCommandLineRunner_Writer {
	if r.self.IsInTestMode() {
		return jutil.NewStringWriter()
	}

	return r.StreamToLegacyOutputWriter(r.self.FilenameToOutputStream(fileName))
}

// FileNameToOutputWriter2 converts a file name into a Writer taking in
// account the output charset.
//
// Java: AbstractCommandLineRunner.fileNameToOutputWriter2(String)
func (r *AbstractCommandLineRunnerBase) FileNameToOutputWriter2(fileName string) AbstractCommandLineRunner_Writer {
	if r.self.IsInTestMode() {
		return jutil.NewStringWriter()
	}

	return r.StreamToOutputWriter2(r.self.FilenameToOutputStream(fileName))
}

// FilenameToOutputStream converts a file name into an OutputStream.
//
// Java: AbstractCommandLineRunner.filenameToOutputStream(String)
func (r *AbstractCommandLineRunnerBase) FilenameToOutputStream(fileName string) io.WriteCloser {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	if err != nil {
		panic(jutil.FileNotFoundExceptionFromError(err, fileName))
	}
	return f
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
	return newCliBufferedWriter(stream, stream, nil, charset)
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

	var propertyMapOutputPath string
	var variableMapOutputPath string

	// Check the create_name_map_files FLAG.
	if config.createNameMapFiles {
		basePath := r.GetMapPath(config.jsOutputFile)

		propertyMapOutputPath = basePath + "_props_map.out"
		variableMapOutputPath = basePath + "_vars_map.out"
	}

	// Check the individual FLAGS.
	if config.variableMapOutputFile != "" {
		if variableMapOutputPath != "" {
			panic(jscomp.NewFlagUsageException(
				"The flags variable_map_output_file and " +
					"create_name_map_files cannot both be used simultaneously."))
		}

		variableMapOutputPath = config.variableMapOutputFile
	}

	if config.propertyMapOutputFile != "" {
		if propertyMapOutputPath != "" {
			panic(jscomp.NewFlagUsageException(
				"The flags property_map_output_file and " +
					"create_name_map_files cannot both be used simultaneously."))
		}

		propertyMapOutputPath = config.propertyMapOutputFile
	}

	// Output the maps.
	if variableMapOutputPath != "" && r.compiler.GetVariableMap() != nil {
		r.compiler.GetVariableMap().Save(variableMapOutputPath)
	}

	if propertyMapOutputPath != "" && r.compiler.GetPropertyMap() != nil {
		r.compiler.GetPropertyMap().Save(propertyMapOutputPath)
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
					out := r.FileNameToOutputWriter2(path)
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
				out := r.FileNameToOutputWriter2(path)
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
	if r.config.outputChunkDependencies != nil && *r.config.outputChunkDependencies != "" {
		out := r.FileNameToOutputWriter2(*r.config.outputChunkDependencies)
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

type cliBufferedWriter struct {
	bw     *bufio.Writer
	closer io.Closer
	flush  func()
	enc    jutil.StreamEncoder
	closed bool
}

func newCliBufferedWriter(out io.Writer, closer io.Closer, flush func(), charset *jutil.Charset) *cliBufferedWriter {
	if out == nil || charset == nil {
		panic(jutil.NewNullPointerException())
	}
	return &cliBufferedWriter{
		bw:     bufio.NewWriterSize(out, 8192),
		closer: closer,
		flush:  flush,
		enc:    charset.NewStreamEncoder(),
	}
}

func (b *cliBufferedWriter) Write(s string) {
	if ex := b.TryWrite(s); ex != nil {
		panic(ex)
	}
}

func (b *cliBufferedWriter) TryWrite(s string) jutil.IIOException {
	if _, err := b.WriteString(s); err != nil {
		return jutil.NewIOException(err.Error())
	}
	return nil
}

func (b *cliBufferedWriter) AppendCharSequence(csq string) {
	b.Write(csq)
}

func (b *cliBufferedWriter) AppendCharacter(c uint16) {
	if b.closed {
		panic(jutil.NewIOException("Stream closed"))
	}
	if err := b.enc.WriteCharTo(b.bw, c); err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
}

func (b *cliBufferedWriter) WriteString(s string) (int, error) {
	if b.closed {
		return 0, errors.New("Stream closed")
	}
	if err := b.enc.WriteTo(b.bw, s); err != nil {
		return 0, err
	}
	return len(s), nil
}

func (b *cliBufferedWriter) Flush() {
	if ex := b.TryFlush(); ex != nil {
		panic(ex)
	}
}

func (b *cliBufferedWriter) TryFlush() jutil.IIOException {
	if b.closed {
		return jutil.NewIOException("Stream closed")
	}
	if err := b.bw.Flush(); err != nil {
		return jutil.NewIOException(err.Error())
	}
	if b.flush != nil {
		b.flush()
	}
	return nil
}

func (b *cliBufferedWriter) Close() {
	if b.closed {
		return
	}
	b.closed = true
	finishErr := b.enc.FinishTo(b.bw)
	flushErr := b.bw.Flush()
	if b.flush != nil {
		b.flush()
	}
	var closeErr error
	if b.closer != nil {
		closeErr = b.closer.Close()
	}
	if finishErr != nil {
		panic(jutil.NewIOException(finishErr.Error()))
	}
	if flushErr != nil {
		panic(jutil.NewIOException(flushErr.Error()))
	}
	if closeErr != nil {
		panic(jutil.NewIOException(closeErr.Error()))
	}
}
