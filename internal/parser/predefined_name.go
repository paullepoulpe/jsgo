// Port of com/google/javascript/jscomp/parsing/parser/PredefinedName.java @ closure-compiler v20261021

/*
 * Copyright 2011 The Closure Compiler Authors.
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

package parser

// PredefinedName: the set of all non-keyword, non-reserved words used in
// javascript.
//
// Java: final class PredefinedName (static constants only; the private
// constructor has no Go counterpart).
const (
	PredefinedName_AS   = "as"
	PredefinedName_FROM = "from"
	PredefinedName_GET  = "get"
	PredefinedName_OF   = "of"
	PredefinedName_SET  = "set"
)
