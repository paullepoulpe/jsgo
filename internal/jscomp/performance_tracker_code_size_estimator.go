// Port of com/google/javascript/jscomp/PerformanceTrackerCodeSizeEstimator.java @ closure-compiler v20261021

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

// PerformanceTrackerCodeSizeEstimator computes a raw and optional gzip'ed
// size estimate without allocating an output code buffer.
type PerformanceTrackerCodeSizeEstimator struct {
	cg     CodeGenerator
	output *bytes.Buffer
}

// Java: PerformanceTrackerCodeSizeEstimator.estimate(Node, boolean)
func PerformanceTrackerCodeSizeEstimator_Estimate(jsRoot rhino.Node, trackGzSize bool) *PerformanceTrackerCodeSizeEstimator {
	estimator := NewPerformanceTrackerCodeSizeEstimator(trackGzSize)
	estimator.cg.AddNode(jsRoot)
	return estimator
}

// Java: PerformanceTrackerCodeSizeEstimator(boolean)
func NewPerformanceTrackerCodeSizeEstimator(trackGzSize bool) *PerformanceTrackerCodeSizeEstimator {
	e := &PerformanceTrackerCodeSizeEstimator{}
	e.cg.initDefaults()
	e.cg.cc = CodeConsumer{
		mode:               sinkSize,
		continueProcessing: true,
		trackGzSize:        trackGzSize,
	}
	if trackGzSize {
		e.output = &bytes.Buffer{}
		e.cg.cc.gzStream = gzip.NewWriter(e.output)
	}
	return e
}

// Java: PerformanceTrackerCodeSizeEstimator.getCodeSize()
func (e *PerformanceTrackerCodeSizeEstimator) GetCodeSize() int {
	return e.cg.cc.cost
}

// Java: PerformanceTrackerCodeSizeEstimator.getZippedCodeSize()
func (e *PerformanceTrackerCodeSizeEstimator) GetZippedCodeSize() int {
	if e.cg.cc.trackGzSize {
		if err := e.cg.cc.gzStream.Close(); err != nil {
			panic(jutil.NewRuntimeExceptionWithCause(jutil.NewIOException(err.Error())))
		}
		return e.output.Len()
	}
	return 0
}
