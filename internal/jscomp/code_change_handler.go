// Port of com/google/javascript/jscomp/CodeChangeHandler.java @ closure-compiler v20261021

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// CodeChangeHandler is a simple listener for code change events.
type CodeChangeHandler interface {
	// ReportChange reports a change to some unspecified node of the AST.
	//
	// Java: CodeChangeHandler.reportChange()
	ReportChange()
}

// ForbiddenChange is a change handler that throws an exception if any
// changes are made.
type ForbiddenChange struct{}

// NewForbiddenChange is new ForbiddenChange().
func NewForbiddenChange() *ForbiddenChange { return &ForbiddenChange{} }

// Java: ForbiddenChange.reportChange()
func (f *ForbiddenChange) ReportChange() {
	panic(jutil.NewIllegalStateExceptionMsg("Code changes forbidden"))
}

// RecentChange records whether the code has changed since the last reset.
type RecentChange struct {
	hasChanged bool
}

// NewRecentChange is new RecentChange().
func NewRecentChange() *RecentChange { return &RecentChange{} }

// Java: RecentChange.reportChange()
func (r *RecentChange) ReportChange() {
	r.hasChanged = true
}

// Java: RecentChange.hasCodeChanged()
func (r *RecentChange) HasCodeChanged() bool {
	return r.hasChanged
}

// Java: RecentChange.reset()
func (r *RecentChange) Reset() {
	r.hasChanged = false
}
