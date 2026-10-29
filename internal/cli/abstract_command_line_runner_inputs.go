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
	"archive/zip"
	"bytes"
	"math"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/resources"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Java: AbstractCommandLineRunner.setWarningGuardOptions(CompilerOptions, ArrayList<FlagEntry<CheckLevel>>, DiagnosticGroups)
func AbstractCommandLineRunner_SetWarningGuardOptions(
	options *jscomp.CompilerOptions,
	warningGuards []*AbstractCommandLineRunner_FlagEntry[jscomp.CheckLevel],
	diagnosticGroups *jscomp.DiagnosticGroups) {
	if warningGuards != nil {
		groupNames := jscomp.DiagnosticGroups_GetRegisteredGroups() // Java: .keySet()
		for _, entry := range warningGuards {
			if "*" == entry.value {
				for groupName := range groupNames.Keys() {
					if !jscomp.DiagnosticGroups_wildcardExcludedGroups.Contains(groupName) {
						diagnosticGroups.SetWarningLevel(options, groupName, entry.flag)
					}
				}
			} else if groupNames.ContainsKey(entry.value) {
				diagnosticGroups.SetWarningLevel(options, entry.value, entry.flag)
			} else {
				panic(jscomp.NewFlagUsageException("Unknown diagnostic group: '" + entry.value + "'"))
			}
		}
	}
}

// AbstractCommandLineRunner_GetBuiltinExterns returns a mutable list.
//
// Java: AbstractCommandLineRunner.getBuiltinExterns(CompilerOptions.Environment)
func AbstractCommandLineRunner_GetBuiltinExterns(env jscomp.CompilerOptions_Environment) []*jscomp.SourceFile {
	input := abstractCommandLineRunner_getExternsInput()
	zipReader, err := zip.NewReader(bytes.NewReader(input), int64(len(input)))
	if err != nil {
		panic(jutil.NewZipException(err.Error()))
	}
	envPrefix := asciiLower(env.String()) + "/"
	mapFromExternsZip := jutil.NewLinkedHashMap[string, *jscomp.SourceFile]()
	for _, entry := range zipReader.File {
		filename := entry.Name

		// Always load externs in the root folder.
		// If the non-core-JS externs are organized in subfolders, only load
		// the ones in a subfolder matching the specified environment. Strip the subfolder.
		if strings.Contains(filename, "/") {
			if !strings.HasPrefix(filename, envPrefix) {
				continue
			}
			filename = filename[len(envPrefix):] // remove envPrefix, including '/'
		}

		mapFromExternsZip.Put(
			filename,
			// Give the files an odd prefix, so that they do not conflict
			// with the user's files.
			//
			// Go-only: performance — the content is read from the zip
			// (embedded in the binary) when needed, and again if needed
			// after parsing dropped it (Java: withContentFromStream, which
			// keeps it for the whole compilation).
			jscomp.SourceFile_Builder_().
				WithPath("externs.zip//"+filename).
				WithReloadableContent(func() string {
					entryStream, err := entry.Open()
					if err != nil {
						panic(jutil.NewZipException(err.Error()))
					}
					defer entryStream.Close()
					return jutil.CharStreams_ToString(jutil.NewInputStreamReader(entryStream, jutil.StandardCharsets_UTF_8))
				}).
				Build())
	}
	return DefaultExterns_PrepareExterns(env, mapFromExternsZip)
}

// Java: AbstractCommandLineRunner.getExternsInput()
func abstractCommandLineRunner_getExternsInput() []byte {
	input := resources.ResourceLoader_GetResourceAsStream("com.google.javascript.jscomp.AbstractCommandLineRunner", "/externs.zip")
	if input == nil {
		// In some environments, the externs.zip is relative to this class.
		input = resources.ResourceLoader_GetResourceAsStream("com.google.javascript.jscomp.AbstractCommandLineRunner", "externs.zip")
	}
	jutil.CheckNotNull(input)
	return input
}

// Java: AbstractCommandLineRunner.parseJsonFilesFromInputStream()
func (r *AbstractCommandLineRunnerBase) ParseJsonFilesFromInputStream() []*AbstractCommandLineRunner_JsonFileSpec {
	var jsonFiles []*AbstractCommandLineRunner_JsonFileSpec
	reader := jutil.NewJsonReader(jutil.NewInputStreamReader(r.in, r.inputCharset))
	defer reader.Close()
	reader.BeginArray()
	for reader.HasNext() {
		jsonFiles = append(jsonFiles, parseJsonFileSpec(reader))
	}
	reader.EndArray()
	return jsonFiles
}

func parseJsonFileSpec(reader *jutil.JsonReader) (result *AbstractCommandLineRunner_JsonFileSpec) {
	isEmpty := true
	oldLenient := reader.IsLenient()
	reader.SetLenient(true)
	defer reader.SetLenient(oldLenient)
	jutil.Try(func() {
		tok := reader.Peek()
		isEmpty = false
		if tok == jutil.JsonToken_NULL {
			reader.NextNull()
			return
		}
		spec := newAbstractCommandLineRunner_JsonFileSpecDefault()
		reader.BeginObject()
		for reader.HasNext() {
			switch reader.NextName() {
			case "src":
				spec.Src = readJsonNullableString(reader)
			case "path":
				spec.Path = readJsonNullableString(reader)
			case "source_map", "sourceMap":
				spec.SourceMap = readJsonNullableString(reader)
			case "webpack_id", "webpackId":
				spec.WebpackId = readJsonNullableString(reader)
			default:
				reader.SkipValue()
			}
		}
		reader.EndObject()
		result = spec
	}, func(t jutil.Throwable) {
		switch e := t.(type) {
		case jutil.IEOFException:
			if isEmpty {
				result = nil
				return
			}
			panic(jutil.NewJsonSyntaxExceptionWithCause(e))
		case jutil.IIllegalStateException:
			panic(jutil.NewJsonSyntaxExceptionWithCause(e))
		case jutil.IIOException:
			panic(jutil.NewJsonSyntaxExceptionWithCause(e))
		default:
			panic(t)
		}
	})
	return result
}

func readJsonNullableString(reader *jutil.JsonReader) *string {
	switch reader.Peek() {
	case jutil.JsonToken_NULL:
		reader.NextNull()
		return nil
	case jutil.JsonToken_BOOLEAN:
		if reader.NextBoolean() {
			s := "true"
			return &s
		}
		s := "false"
		return &s
	default:
		s := reader.NextString()
		return &s
	}
}

// CreateInputs creates inputs from a list of files.
//
// files is a list of flag entries indicates js and zip file names;
// allowStdIn whether '-' is allowed appear as a filename to represent stdin
// (if true, '-' is only allowed to appear once); jsChunkSpecs a list chunk
// specs. It returns an array of inputs.
//
// Java: AbstractCommandLineRunner.createInputs(List<FlagEntry<JsSourceType>>, boolean, List<JsChunkSpec>)
func (r *AbstractCommandLineRunnerBase) CreateInputs(
	files []*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType],
	allowStdIn bool,
	jsChunkSpecs []*AbstractCommandLineRunner_JsChunkSpec) []*jscomp.SourceFile {
	return r.self.CreateInputsWithJsonFilesAndAllowStdIn(files /* jsonFiles= */, nil, allowStdIn, jsChunkSpecs)
}

// CreateInputsWithJsonFiles creates inputs from a list of source files and
// json files.
//
// Java: AbstractCommandLineRunner.createInputs(List<FlagEntry<JsSourceType>>, List<JsonFileSpec>, List<JsChunkSpec>)
func (r *AbstractCommandLineRunnerBase) CreateInputsWithJsonFiles(
	files []*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType],
	jsonFiles []*AbstractCommandLineRunner_JsonFileSpec,
	jsChunkSpecs []*AbstractCommandLineRunner_JsChunkSpec) []*jscomp.SourceFile {
	return r.self.CreateInputsWithJsonFilesAndAllowStdIn(files, jsonFiles /* allowStdIn= */, false, jsChunkSpecs)
}

// CreateInputsWithJsonFilesAndAllowStdIn creates inputs from a list of
// source files, zips and json files.
//
// Can be overridden by subclasses who want to pull files from different
// places.
//
// Go-only: language-adapter — a nil jsonFiles is Java's null list; an empty
// one is passed as a non-nil empty slice.
//
// Java: AbstractCommandLineRunner.createInputs(List<FlagEntry<JsSourceType>>, List<JsonFileSpec>, boolean, List<JsChunkSpec>)
func (r *AbstractCommandLineRunnerBase) CreateInputsWithJsonFilesAndAllowStdIn(
	files []*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType],
	jsonFiles []*AbstractCommandLineRunner_JsonFileSpec,
	allowStdIn bool,
	jsChunkSpecs []*AbstractCommandLineRunner_JsChunkSpec) []*jscomp.SourceFile {
	inputs := make([]*jscomp.SourceFile, 0, len(files))
	usingStdin := false
	jsChunkIndex := 0
	var jsChunkSpec *AbstractCommandLineRunner_JsChunkSpec
	if len(jsChunkSpecs) > 0 {
		jsChunkSpec = jsChunkSpecs[0]
	}
	cumulatedInputFilesExpected := math.MaxInt32
	if jsChunkSpec != nil {
		cumulatedInputFilesExpected = jsChunkSpec.GetNumInputs()
	}
	for i := 0; i < len(files); i++ {
		file := files[i]
		filename := file.value
		if file.flag == AbstractCommandLineRunner_JsSourceType_JS_ZIP {
			if r.config.typedAstListInputFilename != nil {
				panic(jscomp.NewFlagUsageException("Can't use TypedASTs with --zip."))
			}

			if "-" != filename {
				newFiles := jscomp.SourceFile_FromZipFile(filename, r.inputCharset)
				// Update the manifest maps for new zip entries.
				if rootFilename, ok := r.rootRelativePathsMap[filename]; ok {
					for _, zipEntry := range newFiles {
						zipEntryName := zipEntry.GetName()
						jutil.CheckState(strings.Contains(zipEntryName, filename))
						zipmap := strings.ReplaceAll(zipEntryName, filename, rootFilename)
						r.rootRelativePathsMap[zipEntryName] = zipmap
					}
				}

				inputs = append(inputs, newFiles...)
				if jsChunkSpec != nil {
					jsChunkSpec.numJsFiles += len(newFiles) - 1
				}
			}
		} else if "-" != filename {
			kind := rhino.StaticSourceFile_SourceKind_STRONG
			if file.flag == AbstractCommandLineRunner_JsSourceType_WEAKDEP {
				kind = rhino.StaticSourceFile_SourceKind_WEAK
			}
			inputs = append(inputs,
				jscomp.SourceFile_Builder_().
					WithKind(kind).
					WithCharset(r.inputCharset).
					WithPath(filename).
					Build())
		} else {
			if r.config.typedAstListInputFilename != nil {
				panic(jscomp.NewFlagUsageException("Can't use TypedASTs with stdin."))
			}
			if !r.config.defaultToStdin {
				panic(jscomp.NewFlagUsageException("Can't specify stdin."))
			}
			if usingStdin {
				panic(jscomp.NewFlagUsageException("Can't specify stdin twice."))
			}

			if len(r.config.outputManifests) != 0 {
				panic(jscomp.NewFlagUsageException(
					"Manifest files cannot be generated when the input is from stdin."))
			}
			if len(r.config.outputBundles) != 0 {
				panic(jscomp.NewFlagUsageException(
					"Bundle files cannot be generated when the input is from stdin."))
			}

			r.err.Println(AbstractCommandLineRunner_WAITING_FOR_INPUT_WARNING)
			inputs = append(inputs,
				jscomp.SourceFile_Builder_().
					WithPath("stdin").
					WithContentFromStream(r.in).
					WithCharset(r.inputCharset).
					Build())
			usingStdin = true
		}
		if i >= cumulatedInputFilesExpected-1 {
			jsChunkIndex++
			if jsChunkIndex < len(jsChunkSpecs) {
				jsChunkSpec = jsChunkSpecs[jsChunkIndex]
				cumulatedInputFilesExpected += jsChunkSpec.GetNumInputs()
			}
		}
	}
	if jsonFiles != nil {
		for _, jsonFile := range jsonFiles {
			inputs = append(inputs, jscomp.SourceFile_FromCode(*jsonFile.GetPath(), *jsonFile.GetSrc()))
		}
	}

	return inputs
}

// DeduplicateIjsFiles removes input --ijs files whose basename matches an
// input --js or --weakdep file. It throws FlagUsageException if there are
// both input --ijs files and chunk specs.
//
// Go-only: language-adapter — the Java method removes entries from the
// caller's list through its iterator; files is a pointer to the caller's
// slice.
//
// Java: AbstractCommandLineRunner.deduplicateIjsFiles(List<FlagEntry<JsSourceType>>, List<String>, boolean)
func (r *AbstractCommandLineRunnerBase) DeduplicateIjsFiles(
	files *[]*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType],
	moduleRoots []string,
	hasChunkSpecs bool) []*jscomp.JSError {
	var errors []*jscomp.JSError

	// First pass: collect the (module root relative) names of --js and --weakdep files.
	relativeToAbsoluteName := make(map[string]string, len(*files))
	for _, file := range *files {
		// TODO(tjgq): Handle zip files.
		if file.flag == AbstractCommandLineRunner_JsSourceType_JS || file.flag == AbstractCommandLineRunner_JsSourceType_WEAKDEP {
			absoluteName := file.value
			relativeName := r.GetModuleRootRelativeName(absoluteName, moduleRoots)
			relativeToAbsoluteName[relativeName] = absoluteName
		}
	}

	// Second pass: drop --ijs files whose (module root relative) name matches a --js or --weakdep
	// file.
	// Go-only: language-adapter — Iterator.remove() as an index loop.
	for i := 0; i < len(*files); {
		file := (*files)[i]
		if file.flag == AbstractCommandLineRunner_JsSourceType_IJS {
			if hasChunkSpecs {
				panic(jscomp.NewFlagUsageException("--ijs is incompatible with --chunk or --module."))
			}
			absoluteName := file.value
			if !strings.HasSuffix(absoluteName, ".i.js") {
				errors = append(errors, jscomp.JSError_Make(jscomp.IjsErrors_BAD_IJS_FILE_NAME, absoluteName))
				i++
				continue
			}
			relativeName := r.GetModuleRootRelativeName(absoluteName, moduleRoots)
			relativeNonIjsName := relativeName[:len(relativeName)-len(".i.js")]
			if conflictingAbsName, ok := relativeToAbsoluteName[relativeNonIjsName]; ok {
				errors = append(errors,
					jscomp.JSError_Make(
						jscomp.IjsErrors_CONFLICTING_IJS_FILE,
						conflictingAbsName,
						absoluteName))
				*files = append((*files)[:i:i], (*files)[i+1:]...)
				continue
			}
		}
		i++
	}

	return errors
}

// Java: AbstractCommandLineRunner.getModuleRootRelativeName(String, List<String>)
func (r *AbstractCommandLineRunnerBase) GetModuleRootRelativeName(filename string, moduleRoots []string) string {
	for _, moduleRoot := range moduleRoots {
		if strings.HasPrefix(filename, moduleRoot+"/") {
			return filename[len(moduleRoot)+1:]
		}
	}
	return filename
}

// CreateSourceInputs creates JS source code inputs from a list of files.
//
// Java: AbstractCommandLineRunner.createSourceInputs(List<JsChunkSpec>, List<FlagEntry<JsSourceType>>, List<JsonFileSpec>, List<String>)
func (r *AbstractCommandLineRunnerBase) CreateSourceInputs(
	jsChunkSpecs []*AbstractCommandLineRunner_JsChunkSpec,
	files []*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType],
	jsonFiles []*AbstractCommandLineRunner_JsonFileSpec,
	moduleRoots []string) []*jscomp.SourceFile {
	if r.self.IsInTestMode() {
		if r.inputsSupplierForTesting != nil {
			return r.inputsSupplierForTesting()
		}
		return nil
	}
	if len(files) == 0 && jsonFiles == nil && r.config.defaultToStdin {
		// Request to read from stdin.
		files = []*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType]{
			NewAbstractCommandLineRunner_FlagEntry(AbstractCommandLineRunner_JsSourceType_JS, "-")}
	}

	for _, error := range r.DeduplicateIjsFiles(&files, moduleRoots, len(jsChunkSpecs) != 0) {
		r.compiler.Report(error)
	}

	var result []*jscomp.SourceFile
	jutil.Try(func() {
		if jsonFiles != nil {
			result = r.CreateInputsWithJsonFiles(files, jsonFiles, jsChunkSpecs)
		} else {
			result = r.CreateInputs(files, true, jsChunkSpecs)
		}
	}, func(e jscomp.IFlagUsageException) {
		panic(jscomp.NewFlagUsageException("Bad --js flag. " + e.GetMessage()))
	})
	return result
}

// CreateExternInputs creates JS extern inputs from a list of files.
//
// Java: AbstractCommandLineRunner.createExternInputs(List<String>)
func (r *AbstractCommandLineRunnerBase) CreateExternInputs(files []string) []*jscomp.SourceFile {
	var externFiles []*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType]
	if len(files) > 0 {
		externFiles = make([]*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType], len(files))
		for i, file := range files {
			externFiles[i] = NewAbstractCommandLineRunner_FlagEntry(AbstractCommandLineRunner_JsSourceType_EXTERN, file)
		}
	}
	var result []*jscomp.SourceFile
	jutil.Try(func() {
		result = r.CreateInputs(externFiles, false, []*AbstractCommandLineRunner_JsChunkSpec{})
	}, func(e jscomp.IFlagUsageException) {
		panic(jscomp.NewFlagUsageException("Bad --externs flag. " + e.GetMessage()))
	})
	return result
}

// AbstractCommandLineRunner_CreateJsChunks creates chunk objects from a
// list of chunk specifications.
//
// specs is a list of chunk specifications, not null or empty; inputs a list
// of JS file paths, not null. It returns an array of chunk objects.
//
// Java: AbstractCommandLineRunner.createJsChunks(List<JsChunkSpec>, List<CompilerInput>)
func AbstractCommandLineRunner_CreateJsChunks(specs []*AbstractCommandLineRunner_JsChunkSpec, inputs []*jscomp.CompilerInput) []*jscomp.JSChunk {
	jutil.CheckState(specs != nil)
	jutil.CheckState(len(specs) != 0)
	jutil.CheckState(inputs != nil)

	result := make([]*jscomp.JSChunk, 0, len(specs))
	chunksByName := make(map[string]*jscomp.JSChunk, len(specs))
	numJsFilesExpected := 0
	minJsFilesRequired := 0
	for _, spec := range specs {
		if _, exists := chunksByName[spec.name]; exists {
			panic(jscomp.NewFlagUsageException("Duplicate chunk name: " + spec.name))
		}
		chunk := jscomp.NewJSChunk(spec.name)

		for _, dep := range spec.deps {
			other := chunksByName[dep]
			if other == nil {
				panic(jscomp.NewFlagUsageException(
					"Chunk '" +
						spec.name +
						"' depends on unknown chunk '" +
						dep +
						"'. Be sure to list chunks in dependency order."))
			}
			chunk.AddDependency(other)
		}

		// We will allow chunks of zero input.
		if spec.numJsFiles < 0 {
			numJsFilesExpected = -1
		} else {
			minJsFilesRequired += spec.numJsFiles
		}

		if numJsFilesExpected >= 0 {
			numJsFilesExpected += spec.numJsFiles
		}

		chunksByName[spec.name] = chunk
		result = append(result, chunk)
	}

	totalNumJsFiles := len(inputs)

	if numJsFilesExpected >= 0 || minJsFilesRequired > totalNumJsFiles {
		if minJsFilesRequired > totalNumJsFiles {
			numJsFilesExpected = minJsFilesRequired
		}

		if numJsFilesExpected > totalNumJsFiles {
			panic(jscomp.NewFlagUsageException(
				"Not enough JS files specified. Expected " +
					strconv.Itoa(numJsFilesExpected) +
					" but found " +
					strconv.Itoa(totalNumJsFiles)))
		} else if numJsFilesExpected < totalNumJsFiles {
			panic(jscomp.NewFlagUsageException(
				"Too many JS files specified. Expected " +
					strconv.Itoa(numJsFilesExpected) +
					" but found " +
					strconv.Itoa(totalNumJsFiles)))
		}
	}

	// Allocate source files to chunks in reverse order so that the first chunk
	// (presumably the base chunk) can have a size of 'auto' (-1).
	numJsFilesLeft := totalNumJsFiles
	for i := len(specs) - 1; i >= 0; i-- {
		numJsFiles := specs[i].numJsFiles
		chunk := result[i]

		// Check if the first chunk specified 'auto' for the number of files
		if i == 0 && numJsFiles == -1 {
			numJsFiles = numJsFilesLeft
		}

		chunkFiles := inputs[numJsFilesLeft-numJsFiles : numJsFilesLeft]
		for _, input := range chunkFiles {
			chunk.Add(input)
		}
		numJsFilesLeft -= numJsFiles
	}

	return result
}

// CheckChunkName validates the chunk name. Can be overridden by subclasses.
//
// Java: AbstractCommandLineRunner.checkChunkName(String)
func (r *AbstractCommandLineRunnerBase) CheckChunkName(name string) {
	if !rhino.TokenStream_IsJSIdentifier(name) {
		panic(jscomp.NewFlagUsageException("Invalid chunk name: '" + name + "'"))
	}
}

// Java: AbstractCommandLineRunner.initWithTypedAstFilesystem(List<SourceFile>, List<SourceFile>, CompilerOptions, String)
func (r *AbstractCommandLineRunnerBase) InitWithTypedAstFilesystem(
	externs []*jscomp.SourceFile,
	sources []*jscomp.SourceFile,
	options *jscomp.CompilerOptions,
	filename string) {
	jutil.Try(func() {
		typedAstListStream := jutil.NewGZIPInputStreamWithSize(jutil.NewFileInputStream(filename), abstractCommandLineRunner_GZIPPED_TYPEDAST_BUFFER_SIZE)
		defer typedAstListStream.Close()
		r.compiler.InitWithTypedAstFilesystem(externs, sources, options, typedAstListStream)
	}, func(e jutil.IIOException) {
		r.compiler.Report(jscomp.JSError_Make(AbstractCommandLineRunner_COULD_NOT_DESERIALIZE_AST, filename))
	})
}

// Java: AbstractCommandLineRunner.initChunksWithTypedAstFilesystem(List<SourceFile>, List<JSChunk>, CompilerOptions, String)
func (r *AbstractCommandLineRunnerBase) InitChunksWithTypedAstFilesystem(
	externs []*jscomp.SourceFile, chunks []*jscomp.JSChunk, options *jscomp.CompilerOptions, filename string) {
	jutil.Try(func() {
		typedAstListStream := jutil.NewGZIPInputStreamWithSize(jutil.NewFileInputStream(filename), abstractCommandLineRunner_GZIPPED_TYPEDAST_BUFFER_SIZE)
		defer typedAstListStream.Close()
		r.compiler.InitChunksWithTypedAstFilesystem(externs, chunks, options, typedAstListStream)
	}, func(e jutil.IIOException) {
		r.compiler.Report(jscomp.JSError_Make(AbstractCommandLineRunner_COULD_NOT_DESERIALIZE_AST, filename))
	})
}

// Java: AbstractCommandLineRunner.restoreState(String)
func (r *AbstractCommandLineRunnerBase) RestoreState(filename string) {
	jutil.Try(func() {
		serializedInputStream := jutil.NewFileInputStream(filename)
		defer serializedInputStream.Close()
		r.compiler.RestoreState(serializedInputStream)
	}, func(e jutil.IIOException) {
		r.compiler.Report(jscomp.JSError_Make(AbstractCommandLineRunner_COULD_NOT_DESERIALIZE_AST, filename))
	})
}

// SaveState is called at the end of compilation to save the compiler state
// if applicable.
//
// Java: AbstractCommandLineRunner.saveState()
func (r *AbstractCommandLineRunnerBase) SaveState() {
	if r.config.GetSaveCompilationStateToFilename() == nil {
		// nothing to save to.
		return
	}
	filename := *r.config.GetSaveCompilationStateToFilename()
	jutil.Try(func() {
		serializedOutputStream := jutil.NewFileOutputStream(filename)
		defer serializedOutputStream.Close()
		r.compiler.SaveState(serializedOutputStream)
	}, func(e jutil.IIOException) {
		r.compiler.Report(jscomp.JSError_Make(AbstractCommandLineRunner_COULD_NOT_SERIALIZE_AST, filename))
	})
}

// GetInputCharset queries the flag for the input charset, and returns a
// Charset object representing the selection. It throws FlagUsageException
// if flag is not a valid Charset name.
//
// Java: AbstractCommandLineRunner.getInputCharset()
func (r *AbstractCommandLineRunnerBase) GetInputCharset() *jutil.Charset {
	if r.config.charset != "" {
		if !jutil.Charset_IsSupported(r.config.charset) {
			panic(jscomp.NewFlagUsageException(r.config.charset + " is not a valid charset name."))
		}
		return jutil.Charset_ForName(r.config.charset)
	}
	return jutil.StandardCharsets_UTF_8
}

// GetLegacyOutputCharset queries the flag for the output charset.
//
// Let the outputCharset be the same as the input charset... except if we're
// reading in UTF-8 by default. By tradition, we've always output ASCII to
// avoid various hiccups with different browsers, proxies and firewalls.
//
// It returns the charset to use when writing outputs (guaranteed to be a
// supported charset) and throws FlagUsageException if flag is not a valid
// Charset name.
//
// Java: AbstractCommandLineRunner.getLegacyOutputCharset()
func (r *AbstractCommandLineRunnerBase) GetLegacyOutputCharset() *jutil.Charset {
	if r.config.charset != "" {
		if !jutil.Charset_IsSupported(r.config.charset) {
			panic(jscomp.NewFlagUsageException(r.config.charset + " is not a valid charset name."))
		}
		return jutil.Charset_ForName(r.config.charset)
	}
	return jutil.StandardCharsets_US_ASCII
}

// GetOutputCharset2 queries the flag for the output charset. Defaults to
// UTF-8. It throws FlagUsageException if flag is not a valid Charset name.
//
// Java: AbstractCommandLineRunner.getOutputCharset2()
func (r *AbstractCommandLineRunnerBase) GetOutputCharset2() *jutil.Charset {
	if r.config.charset != "" {
		if !jutil.Charset_IsSupported(r.config.charset) {
			panic(jscomp.NewFlagUsageException(r.config.charset + " is not a valid charset name."))
		}
		return jutil.Charset_ForName(r.config.charset)
	}
	return jutil.StandardCharsets_UTF_8
}

// Java: AbstractCommandLineRunner.createExterns(CompilerOptions)
func (r *AbstractCommandLineRunnerBase) CreateExterns(options *jscomp.CompilerOptions) []*jscomp.SourceFile {
	if r.self.IsInTestMode() {
		return r.externsSupplierForTesting()
	}
	return r.CreateExternInputs(r.config.externs)
}

// AbstractCommandLineRunner_CreateDefineReplacements creates a map of
// constant names to constant values from a textual description of the map.
//
// definitions is a list of overriding definitions for defines in the form
// <name>[=<val>], where <val> is a number, boolean, or single-quoted string
// without single quotes.
//
// Java: AbstractCommandLineRunner.createDefineReplacements(List<String>, CompilerOptions)
func AbstractCommandLineRunner_CreateDefineReplacements(definitions []string, options *jscomp.CompilerOptions) {
	// Parse the definitions
	for _, override := range definitions {
		assignment := jutil.StringSplitLiteralLimit(override, "=", 2)
		defName := assignment[0]

		if defName != "" {
			defValue := "true"
			if len(assignment) != 1 {
				defValue = assignment[1]
			}

			isTrue := defValue == "true"
			isFalse := defValue == "false"
			if isTrue || isFalse {
				options.SetDefineToBooleanLiteral(defName, isTrue)
				continue
			} else if len(defValue) > 1 &&
				(defValue[0] == '\'' || defValue[0] == '"') &&
				defValue[len(defValue)-1] == defValue[0] {
				// If the value starts and ends with a single quote,
				// we assume that it's a string.
				maybeStringVal := defValue[1 : len(defValue)-1]
				if strings.IndexByte(maybeStringVal, defValue[0]) == -1 {
					options.SetDefineToStringLiteral(defName, maybeStringVal)
					continue
				}
			} else {
				parsed := false
				jutil.Try(func() {
					value := jutil.DoubleParseDouble(defValue)
					options.SetDefineToDoubleLiteral(defName, value)
					parsed = true
				}, func(e jutil.INumberFormatException) {
					// do nothing, it will be caught at the end
				})
				if parsed {
					continue
				}

				if defValue != "" {
					options.SetDefineToStringLiteral(defName, defValue)
					continue
				}
			}
		}

		panic(jutil.NewRuntimeException("--define flag syntax invalid: " + override))
	}
}

// ConstructRootRelativePathsMap constructs and returns the input root path
// map. The key is the exec path of each input file, and the value is the
// corresponding root relative path.
//
// Java: AbstractCommandLineRunner.constructRootRelativePathsMap()
func (r *AbstractCommandLineRunnerBase) ConstructRootRelativePathsMap() map[string]string {
	rootRelativePathsMap := make(map[string]string, len(r.config.manifestMaps))
	for _, mapString := range r.config.manifestMaps {
		colonIndex := strings.IndexByte(mapString, ':')
		jutil.CheckState(colonIndex > 0)
		execPath := mapString[:colonIndex]
		rootRelativePath := mapString[colonIndex+1:]
		jutil.CheckState(strings.IndexByte(rootRelativePath, ':') == -1)
		rootRelativePathsMap[execPath] = rootRelativePath
	}
	return rootRelativePathsMap
}
