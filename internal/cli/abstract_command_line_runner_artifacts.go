// Port of com/google/javascript/jscomp/AbstractCommandLineRunner.java @ closure-compiler v20261028

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

func mustSourceMapOutputPath(options *jscomp.CompilerOptions) string {
	p, ok := options.GetSourceMapOutputPath()
	if !ok {
		panic(jutil.NewNullPointerException())
	}
	return p
}

// ShouldGenerateMapPerChunk returns true if and only if a source map file
// should be generated for each chunk, as opposed to one unified map. This is
// specified by having the source map pattern include the %outname%
// variable.
//
// Java: AbstractCommandLineRunner.shouldGenerateMapPerChunk(B)
func (r *abstractCommandLineRunnerBase) ShouldGenerateMapPerChunk(options *jscomp.CompilerOptions) bool {
	return options.ShouldGatherSourceMapInfo() &&
		strings.Contains(mustSourceMapOutputPath(options), "%outname%")
}

func (r *abstractCommandLineRunnerBase) shouldGenerateMapPerChunk(options *jscomp.CompilerOptions) bool {
	return r.self.ShouldGenerateMapPerChunk(options)
}

// openExternExportsStream returns a stream for outputting the generated
// externs file, or nil if no extern-ed exports are being generated.
//
// Java: AbstractCommandLineRunner.openExternExportsStream(B, String)
func (r *abstractCommandLineRunnerBase) openExternExportsStream(options *jscomp.CompilerOptions, path string) *cliBufferedWriter {
	exPath, ok := options.GetExternExportsPath()
	if !ok {
		return nil
	}

	if !strings.Contains(exPath, jutil.File_separator) {
		outputFile := jutil.NewFile(path)
		parent := "null" // Java: a null parent concatenates as "null".
		if p, hasParent := outputFile.Parent(); hasParent {
			parent = p
		}
		exPath = parent + jutil.File_separator + exPath
	}

	return r.fileNameToOutputWriter2(exPath)
}

// expandCommandLinePath expands a file path specified on the command-line.
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
func (r *abstractCommandLineRunnerBase) expandCommandLinePath(path string, forChunk *jscomp.JSChunk) string {
	var sub string
	if forChunk != nil {
		sub = r.GetChunkOutputFileName(forChunk)
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
func (r *abstractCommandLineRunnerBase) ExpandSourceMapPath(options *jscomp.CompilerOptions, forChunk *jscomp.JSChunk) string {
	if !options.ShouldGatherSourceMapInfo() {
		return ""
	}
	return r.expandCommandLinePath(mustSourceMapOutputPath(options), forChunk)
}

// fileNameToLegacyOutputWriter converts a file name into a Writer taking into
// account the legacy output charset.
//
// Java: AbstractCommandLineRunner.fileNameToLegacyOutputWriter(String)
func (r *abstractCommandLineRunnerBase) fileNameToLegacyOutputWriter(fileName string) *cliBufferedWriter {
	if r.testMode {
		return newCliBufferedWriter(io.Discard, nil, nil, jutil.StandardCharsets_UTF_8)
	}
	return r.streamToLegacyOutputWriter(r.filenameToOutputStream(fileName))
}

// fileNameToOutputWriter2 converts a file name into a Writer taking into
// account the newer output charset.
//
// Java: AbstractCommandLineRunner.fileNameToOutputWriter2(String)
func (r *abstractCommandLineRunnerBase) fileNameToOutputWriter2(fileName string) *cliBufferedWriter {
	if r.testMode {
		return newCliBufferedWriter(io.Discard, nil, nil, jutil.StandardCharsets_UTF_8)
	}
	return r.streamToOutputWriter2(r.filenameToOutputStream(fileName))
}

// filenameToOutputStream converts a file name into an OutputStream.
//
// Java: AbstractCommandLineRunner.filenameToOutputStream(String)
func (r *abstractCommandLineRunnerBase) filenameToOutputStream(fileName string) io.WriteCloser {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	if err != nil {
		panic(jutil.FileNotFoundExceptionFromError(err, fileName))
	}
	return f
}

// streamToLegacyOutputWriter creates a writer with the legacy output charset.
//
// Java: AbstractCommandLineRunner.streamToLegacyOutputWriter(OutputStream)
func (r *abstractCommandLineRunnerBase) streamToLegacyOutputWriter(stream io.WriteCloser) *cliBufferedWriter {
	charset := r.legacyOutputCharset
	if charset == nil {
		charset = jutil.StandardCharsets_UTF_8
	}
	return newCliBufferedWriter(stream, stream, nil, charset)
}

// streamToOutputWriter2 creates a writer with the newer output charset.
//
// Java: AbstractCommandLineRunner.streamToOutputWriter2(OutputStream)
func (r *abstractCommandLineRunnerBase) streamToOutputWriter2(stream io.WriteCloser) *cliBufferedWriter {
	charset := r.outputCharset2
	if charset == nil {
		charset = jutil.StandardCharsets_UTF_8
	}
	return newCliBufferedWriter(stream, stream, nil, charset)
}

// outputSourceMap outputs the source map found in the compiler to the
// proper path if one exists.
//
// Java: AbstractCommandLineRunner.outputSourceMap(B, String)
func (r *abstractCommandLineRunnerBase) outputSourceMap(options *jscomp.CompilerOptions, associatedName string) {
	if !options.ShouldGatherSourceMapInfo() ||
		mustSourceMapOutputPath(options) == "/dev/null" {
		return
	}

	outName := r.ExpandSourceMapPath(options, nil)
	maybeCreateDirsForPath(outName)
	out := r.fileNameToOutputWriter2(outName)
	defer closeBufferedWriter(out)
	r.compiler.GetSourceMap().AppendTo(out, associatedName)
}

// getMapPath returns the path prefix in which to place generated map file(s)
// based on the path at which the JS binary will be placed.
//
// Java: AbstractCommandLineRunner.getMapPath(String)
func (r *abstractCommandLineRunnerBase) getMapPath(outputFile string) string {
	if outputFile == "" {
		// If we have a js_module_binary rule, output the maps
		// at chunkname_props_map.out, etc.
		if r.config.chunkOutputPathPrefix != "" {
			return r.config.chunkOutputPathPrefix
		}
		return "jscompiler"
	}

	file := jutil.NewFile(outputFile)
	outputFileName := strings.TrimSuffix(file.GetName(), ".js")
	if fileParent, hasParent := file.Parent(); hasParent {
		return fileParent + jutil.File_separator + outputFileName
	}
	return outputFileName
}

// outputNameMaps outputs the variable and property name maps for the
// specified compiler if the proper FLAGS are set.
//
// Java: AbstractCommandLineRunner.outputNameMaps()
func (r *abstractCommandLineRunnerBase) outputNameMaps() {
	config := r.config

	var propertyMapOutputPath string
	var variableMapOutputPath string

	if config.createNameMapFiles {
		basePath := r.getMapPath(config.jsOutputFile)
		propertyMapOutputPath = basePath + "_props_map.out"
		variableMapOutputPath = basePath + "_vars_map.out"
	}

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

	if variableMapOutputPath != "" && r.compiler.GetVariableMap() != nil {
		r.compiler.GetVariableMap().Save(variableMapOutputPath)
	}

	if propertyMapOutputPath != "" && r.compiler.GetPropertyMap() != nil {
		r.compiler.GetPropertyMap().Save(propertyMapOutputPath)
	}
}

// outputStringMap outputs the string map generated by the ReplaceStrings
// pass if an output path exists.
//
// Java: AbstractCommandLineRunner.outputStringMap()
func (r *abstractCommandLineRunnerBase) outputStringMap() {
	config := r.config
	if config.stringMapOutputPath == "" {
		return
	}
	if r.compiler.GetStringMap() != nil {
		r.compiler.GetStringMap().Save(config.stringMapOutputPath)
		return
	}
	// Ensure an empty file is created if there is no string map.
	// This avoids confusing some build tools that expect to see the file, even if it is empty.
	f, err := os.OpenFile(config.stringMapOutputPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o666)
	if err == nil {
		_ = f.Close()
	} else if errors.Is(err, os.ErrExist) {
		panic(jutil.NewIOException("Could not create file: " + config.stringMapOutputPath))
	} else {
		reason := jutil.FileNotFoundExceptionFromError(err, "").GetMessage()
		panic(jutil.NewIOException(strings.TrimSuffix(strings.TrimPrefix(reason, " ("), ")")))
	}
}

// shouldGenerateOutputPerChunk returns true if and only if a manifest or
// bundle should be generated for each chunk, as opposed to one unified
// manifest.
//
// Java: AbstractCommandLineRunner.shouldGenerateOutputPerChunk(String)
func (r *abstractCommandLineRunnerBase) shouldGenerateOutputPerChunk(output string) bool {
	return len(r.config.chunk) != 0 && strings.Contains(output, "%outname%")
}

// Java: AbstractCommandLineRunner.outputManifest()
func (r *abstractCommandLineRunnerBase) outputManifest() {
	r.outputManifestOrBundle(r.config.outputManifests, true)
}

// Java: AbstractCommandLineRunner.outputBundle()
func (r *abstractCommandLineRunnerBase) outputBundle() {
	r.outputManifestOrBundle(r.config.outputBundles, false)
}

// Java: AbstractCommandLineRunner.outputInstrumentationMapping()
func (r *abstractCommandLineRunnerBase) outputInstrumentationMapping() {
	if r.config.instrumentationMappingFile != "" {
		path := r.expandCommandLinePath(r.config.instrumentationMappingFile, nil)
		r.compiler.GetInstrumentationMapping().Save(path)
	}
}

// outputManifestOrBundle writes the manifest or bundle of all compiler
// input files that were included as controlled by --dependency_mode, if
// requested.
//
// Java: AbstractCommandLineRunner.outputManifestOrBundle(List<String>, boolean)
func (r *abstractCommandLineRunnerBase) outputManifestOrBundle(outputFiles []string, isManifest bool) {
	if len(outputFiles) == 0 {
		return
	}

	for _, output := range outputFiles {
		if output == "" {
			continue
		}

		if r.shouldGenerateOutputPerChunk(output) {
			// Generate per-chunk manifests or bundles.
			for _, chunk := range r.compiler.GetChunkGraph().GetAllChunks() {
				func() {
					path := r.expandCommandLinePath(output, chunk)
					out := r.fileNameToOutputWriter2(path)
					defer closeBufferedWriter(out)
					if isManifest {
						r.printManifestTo(chunk, out)
					} else {
						r.printBundleTo(chunk, out)
					}
				}()
			}
			continue
		}

		// Generate a single file manifest or bundle.
		func() {
			path := r.expandCommandLinePath(output, nil)
			out := r.fileNameToOutputWriter2(path)
			defer closeBufferedWriter(out)
			if len(r.config.chunk) == 0 {
				// For a single-chunk compilation, generate a single headerless manifest or bundle
				// containing only the strong files.
				chunk := r.compiler.GetChunkGraph().GetChunkByName(jscomp.JSChunk_STRONG_CHUNK_NAME)
				if isManifest {
					r.printManifestTo(chunk, out)
				} else {
					r.printBundleTo(chunk, out)
				}
			} else {
				// For a multi-chunk compilation, generate a single manifest file with chunk headers.
				r.PrintChunkGraphManifestOrBundleTo(r.compiler.GetChunkGraph(), out, isManifest)
			}
		}()
	}
}

// outputChunkGraphJson creates a file containing the current chunk graph in
// JSON serialization.
//
// Java: AbstractCommandLineRunner.outputChunkGraphJson()
func (r *abstractCommandLineRunnerBase) outputChunkGraphJson() {
	if r.config.outputChunkDependencies != "" {
		out := r.fileNameToOutputWriter2(r.config.outputChunkDependencies)
		defer closeBufferedWriter(out)
		r.PrintChunkGraphJsonTo(out)
	}
}

// PrintChunkGraphJsonTo prints the current chunk graph as JSON.
//
// Java: AbstractCommandLineRunner.printChunkGraphJsonTo(Appendable)
func (r *abstractCommandLineRunnerBase) PrintChunkGraphJsonTo(out jutil.Appendable) {
	out.AppendCharSequence(r.compiler.GetChunkGraph().ToJson().String())
}

// PrintChunkGraphManifestOrBundleTo prints a set of chunks to the manifest
// or bundle file.
//
// Java: AbstractCommandLineRunner.printChunkGraphManifestOrBundleTo(JSChunkGraph, Appendable, boolean)
func (r *abstractCommandLineRunnerBase) PrintChunkGraphManifestOrBundleTo(graph *jscomp.JSChunkGraph, out jutil.Appendable, isManifest bool) {
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
			r.printManifestTo(chunk, out)
		} else {
			r.printBundleTo(chunk, out)
		}
		requiresNewline = true
	}
}

// printManifestTo prints a list of input names (using root-relative paths),
// delimited by newlines, to the manifest file.
//
// Java: AbstractCommandLineRunner.printManifestTo(JSChunk, Appendable)
func (r *abstractCommandLineRunnerBase) printManifestTo(chunk *jscomp.JSChunk, out jutil.Appendable) {
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

// printBundleTo prints all the input contents, starting with a comment that
// specifies the input file name (using root-relative paths) before each
// file.
//
// Java: AbstractCommandLineRunner.printBundleTo(JSChunk, Appendable)
func (r *abstractCommandLineRunnerBase) printBundleTo(chunk *jscomp.JSChunk, out jutil.Appendable) {
	inputs := chunk.GetInputs()
	options := r.compiler.GetOptions()
	// Prebuild ASTs before they're needed in getLoadFlags, for performance and because
	// StackOverflowErrors can be hit if not prebuilt.
	if options.GetNumParallelThreads() > 1 {
		jscomp.NewPrebuildAst(r.compiler, options.GetNumParallelThreads()).Prebuild(slices.Values(inputs))
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

		displayName := name
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

// closeBufferedWriter closes w if non-nil.
func closeBufferedWriter(w *cliBufferedWriter) {
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
