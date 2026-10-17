// Port of com/google/javascript/jscomp/PerformanceTrackerCodeSizeEstimator.java @ closure-compiler v20261015

/*
 * Copyright 2016 The Closure Compiler Authors.
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
	"compress/gzip"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// PerformanceTrackerCodeSizeEstimator is a class to get a raw and gzip'ed
// size estimate; it doesn't generate code.
type PerformanceTrackerCodeSizeEstimator struct {
	CodeConsumerBase

	size     int
	lastChar uint16
	// May be nil.
	output *bytes.Buffer
	// May be nil.
	//
	// Go-only: language-adapter — java.util.zip.GZIPOutputStream is zlib's
	// deflate (level 6); Go's compress/flate implements the same format with
	// a different match finder, so the compressed size (the gzSize columns
	// of the tracer report) can differ from Java's by a few bytes.
	stream      *gzip.Writer
	trackGzSize bool
}

// Java: PerformanceTrackerCodeSizeEstimator.estimate(Node, boolean)
func PerformanceTrackerCodeSizeEstimator_Estimate(jsRoot *rhino.Node, trackGzSize bool) *PerformanceTrackerCodeSizeEstimator {
	estimator := NewPerformanceTrackerCodeSizeEstimator(trackGzSize)
	CodeGenerator_ForCostEstimation(estimator).AddNode(jsRoot)
	return estimator
}

// Java: PerformanceTrackerCodeSizeEstimator(boolean)
func NewPerformanceTrackerCodeSizeEstimator(trackGzSize bool) *PerformanceTrackerCodeSizeEstimator {
	e := &PerformanceTrackerCodeSizeEstimator{}
	e.InitCodeConsumerBase(e)
	e.size = 0
	e.lastChar = 0
	e.trackGzSize = trackGzSize
	if trackGzSize {
		e.output = &bytes.Buffer{}
		e.stream = gzip.NewWriter(e.output)
	} else {
		e.output = nil
		e.stream = nil
	}
	return e
}

// Java: PerformanceTrackerCodeSizeEstimator.append(String)
func (e *PerformanceTrackerCodeSizeEstimator) Append(str string) {
	len_ := jutil.StringLength(str)
	if len_ > 0 {
		e.size += len_
		e.lastChar = jutil.StringCharAt(str, len_-1)
		if e.trackGzSize {
			if _, err := e.stream.Write(jutil.StringGetBytesUTF8(str)); err != nil {
				panic(jutil.NewRuntimeExceptionWithCause(jutil.NewIOException(err.Error())))
			}
		}
	}
}

// Java: PerformanceTrackerCodeSizeEstimator.getLastChar()
func (e *PerformanceTrackerCodeSizeEstimator) GetLastChar() uint16 {
	return e.lastChar
}

// Java: PerformanceTrackerCodeSizeEstimator.getCodeSize()
func (e *PerformanceTrackerCodeSizeEstimator) GetCodeSize() int {
	return e.size
}

// Java: PerformanceTrackerCodeSizeEstimator.getZippedCodeSize()
func (e *PerformanceTrackerCodeSizeEstimator) GetZippedCodeSize() int {
	if e.trackGzSize {
		// stream.finish(); stream.close();
		if err := e.stream.Close(); err != nil {
			panic(jutil.NewRuntimeExceptionWithCause(jutil.NewIOException(err.Error())))
		}
		return e.output.Len()
	} else {
		return 0
	}
}
