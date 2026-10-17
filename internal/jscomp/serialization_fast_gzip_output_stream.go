// Port of com/google/javascript/jscomp/serialization/FastGzipOutputStream.java @ closure-compiler v20261015

/*
 * Copyright 2026 The Closure Compiler Authors.
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
	"io"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// FastGzipOutputStream is a GZIPOutputStream that uses BEST_SPEED
// compression level and a reasonable buffer size.
//
// Go-only: language-adapter — on jutil's GZIPOutputStream (Go's
// compress/flate): the same content as Java's, possibly different
// compressed bytes (PORTING.md §7).
//
// Java: public final class FastGzipOutputStream extends GZIPOutputStream
type FastGzipOutputStream struct {
	jutil.GZIPOutputStream
}

// Java: FastGzipOutputStream.DEFAULT_BUFFER_SIZE
const fastGzipOutputStream_DEFAULT_BUFFER_SIZE = 64 * 1024

// Java: FastGzipOutputStream(OutputStream)
func NewFastGzipOutputStream(out io.Writer) *FastGzipOutputStream {
	return NewFastGzipOutputStreamWithSize(out, fastGzipOutputStream_DEFAULT_BUFFER_SIZE)
}

// Java: FastGzipOutputStream(OutputStream, int)
func NewFastGzipOutputStreamWithSize(out io.Writer, size int) *FastGzipOutputStream {
	f := &FastGzipOutputStream{}
	f.InitGZIPOutputStream(out, size, false)
	f.Def.SetLevel(jutil.Deflater_BEST_SPEED)
	return f
}
